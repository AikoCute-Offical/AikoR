package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"

	"github.com/AikoCute-Offical/AikoR/common/limiter"
	"github.com/go-resty/resty/v2"

	"github.com/AikoCute-Offical/AikoR/api"
	"github.com/AikoCute-Offical/AikoR/app/mydispatcher"
	"github.com/AikoCute-Offical/AikoR/common/legocmd"
	"github.com/AikoCute-Offical/AikoR/common/serverstatus"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/task"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/features/inbound"
	"github.com/xtls/xray-core/features/outbound"
	"github.com/xtls/xray-core/features/routing"
	"github.com/xtls/xray-core/features/stats"
)

type LimitInfo struct {
	end              int64
	originSpeedLimit uint64
}

type Controller struct {
	server                  *core.Instance
	config                  *Config
	clientInfo              api.ClientInfo
	apiClient               api.API
	nodeInfo                *api.NodeInfo
	Tag                     string
	userList                *[]api.UserInfo
	nodeInfoMonitorPeriodic *task.Periodic
	userReportPeriodic      *task.Periodic
	onlineIpReportPeriodic  *task.Periodic
	limitedUsers            map[api.UserInfo]LimitInfo
	warnedUsers             map[api.UserInfo]int
	panelType               string
	ihm                     inbound.Manager
	ohm                     outbound.Manager
	stm                     stats.Manager
	dispatcher              *mydispatcher.DefaultDispatcher
}

// New return a Controller service with default parameters.
func New(server *core.Instance, api api.API, config *Config, panelType string) *Controller {
	controller := &Controller{
		server:     server,
		config:     config,
		apiClient:  api,
		panelType:  panelType,
		ihm:        server.GetFeature(inbound.ManagerType()).(inbound.Manager),
		ohm:        server.GetFeature(outbound.ManagerType()).(outbound.Manager),
		stm:        server.GetFeature(stats.ManagerType()).(stats.Manager),
		dispatcher: server.GetFeature(routing.DispatcherType()).(*mydispatcher.DefaultDispatcher),
	}
	return controller
}

// Start implement the Start() function of the service interface
func (c *Controller) Start() error {
	c.clientInfo = c.apiClient.Describe()
	// First fetch Node Info
	newNodeInfo, err := c.apiClient.GetNodeInfo()
	if err != nil {
		return err
	}
	c.nodeInfo = newNodeInfo
	c.Tag = c.buildNodeTag()
	// Add new tag
	err = c.addNewTag(newNodeInfo)
	if err != nil {
		log.Panic(err)
		return err
	}
	// Update user
	userInfo, err := c.apiClient.GetUserList()
	if err != nil {
		return err
	}

	err = c.addNewUser(userInfo, newNodeInfo)
	if err != nil {
		return err
	}
	//sync controller userList
	c.userList = userInfo

	// Add Limiter
	if err := c.AddInboundLimiter(c.Tag, newNodeInfo.SpeedLimit, userInfo); err != nil {
		log.Print(err)
	}
	// Add Rule Manager
	if !c.config.DisableGetRule {
		if ruleList, err := c.apiClient.GetNodeRule(); err != nil {
			log.Printf("Get rule list filed: %s", err)
		} else if len(*ruleList) > 0 {
			if err := c.UpdateRule(c.Tag, *ruleList); err != nil {
				log.Print(err)
			}
		}
	}
	c.nodeInfoMonitorPeriodic = &task.Periodic{
		Interval: time.Duration(c.config.UpdatePeriodic) * time.Second,
		Execute:  c.nodeInfoMonitor,
	}
	c.userReportPeriodic = &task.Periodic{
		Interval: time.Duration(c.config.UpdatePeriodic) * time.Second,
		Execute:  c.userInfoMonitor,
	}
	if c.config.DynamicSpeedLimitConfig == nil {
		c.config.DynamicSpeedLimitConfig = &DynamicSpeedLimitConfig{0, 0, 0, 0}
	}
	if c.config.DynamicSpeedLimitConfig.Limit > 0 {
		c.limitedUsers = make(map[api.UserInfo]LimitInfo)
		c.warnedUsers = make(map[api.UserInfo]int)
	}
	log.Printf("[%s: %d] Start monitor node status", c.nodeInfo.NodeType, c.nodeInfo.NodeID)
	// delay to start nodeInfoMonitor
	go func() {
		time.Sleep(time.Duration(c.config.UpdatePeriodic) * time.Second)
		_ = c.nodeInfoMonitorPeriodic.Start()
	}()

	log.Printf("[%s: %d] Start report node status", c.nodeInfo.NodeType, c.nodeInfo.NodeID)
	// delay to start userReport
	go func() {
		time.Sleep(time.Duration(c.config.UpdatePeriodic) * time.Second)
		_ = c.userReportPeriodic.Start()
	}()

	if c.config.EnableIpRecorder {
		c.onlineIpReportPeriodic = &task.Periodic{
			Interval: time.Duration(c.config.IpRecorderConfig.Periodic) * time.Second,
			Execute:  c.onlineIpReport,
		}
		go func() {
			time.Sleep(time.Duration(c.config.IpRecorderConfig.Periodic) * time.Second)
			_ = c.onlineIpReportPeriodic.Start()
		}()
		log.Printf("[%s: %d] Start report online ip", c.nodeInfo.NodeType, c.nodeInfo.NodeID)
	}
	return nil
}

// Close implement the Close() function of the service interface
func (c *Controller) Close() error {
	if c.nodeInfoMonitorPeriodic != nil {
		err := c.nodeInfoMonitorPeriodic.Close()
		if err != nil {
			log.Panicf("node info periodic close failed: %s", err)
		}
	}

	if c.nodeInfoMonitorPeriodic != nil {
		err := c.userReportPeriodic.Close()
		if err != nil {
			log.Panicf("user report periodic close failed: %s", err)
		}
	}
	if c.onlineIpReportPeriodic != nil {
		err := c.onlineIpReportPeriodic.Close()
		if err != nil {
			log.Panicf("online ip report periodic close failed: %s", err)
		}
	}
	return nil
}

func (c *Controller) nodeInfoMonitor() (err error) {
	// First fetch Node Info
	newNodeInfo, err := c.apiClient.GetNodeInfo()
	if err != nil {
		log.Print(err)
		return nil
	}

	// Update User
	newUserInfo, err := c.apiClient.GetUserList()
	if err != nil {
		log.Print(err)
		return nil
	}

	var nodeInfoChanged = false
	// If nodeInfo changed
	if !reflect.DeepEqual(c.nodeInfo, newNodeInfo) {
		// Remove old tag
		oldtag := c.Tag
		err := c.removeOldTag(oldtag)
		if err != nil {
			log.Print(err)
			return nil
		}
		if c.nodeInfo.NodeType == "Shadowsocks-Plugin" {
			err = c.removeOldTag(fmt.Sprintf("dokodemo-door_%s+1", c.Tag))
		}
		if err != nil {
			log.Print(err)
			return nil
		}
		// Add new tag
		c.nodeInfo = newNodeInfo
		c.Tag = c.buildNodeTag()
		err = c.addNewTag(newNodeInfo)
		if err != nil {
			log.Print(err)
			return nil
		}
		nodeInfoChanged = true
		// Remove Old limiter
		if err = c.DeleteInboundLimiter(oldtag); err != nil {
			log.Print(err)
			return nil
		}
	}

	// Check Rule
	if !c.config.DisableGetRule {
		if ruleList, err := c.apiClient.GetNodeRule(); err != nil {
			log.Printf("Get rule list filed: %s", err)
		} else if len(*ruleList) > 0 {
			if err := c.UpdateRule(c.Tag, *ruleList); err != nil {
				log.Print(err)
			}
		}
	}

	// Check Cert
	if c.nodeInfo.EnableTLS && (c.config.CertConfig.CertMode == "dns" || c.config.CertConfig.CertMode == "http") {
		lego, err := legocmd.New()
		if err != nil {
			log.Print(err)
		}
		// Xray-core supports the OcspStapling certification hot renew
		_, _, err = lego.RenewCert(c.config.CertConfig.CertDomain, c.config.CertConfig.Email, c.config.CertConfig.CertMode, c.config.CertConfig.Provider, c.config.CertConfig.DNSEnv)
		if err != nil {
			log.Print(err)
		}
	}

	if nodeInfoChanged {
		err = c.addNewUser(newUserInfo, newNodeInfo)
		if err != nil {
			log.Print(err)
			return nil
		}
		// Add Limiter
		if err := c.AddInboundLimiter(c.Tag, newNodeInfo.SpeedLimit, newUserInfo); err != nil {
			log.Print(err)
			return nil
		}
	} else {
		deleted, added := compareUserList(c.userList, newUserInfo)
		if len(deleted) > 0 {
			deletedEmail := make([]string, len(deleted))
			for i, u := range deleted {
				deletedEmail[i] = fmt.Sprintf("%s|%s|%d", c.Tag, u.Email, u.UID)
			}
			err := c.removeUsers(deletedEmail, c.Tag)
			if err != nil {
				log.Print(err)
			}
		}
		if len(added) > 0 {
			err = c.addNewUser(&added, c.nodeInfo)
			if err != nil {
				log.Print(err)
			}
			// Update Limiter
			if err := c.UpdateInboundLimiter(c.Tag, &added); err != nil {
				log.Print(err)
			}
		}
		log.Printf("[%s: %d] %d user deleted, %d user added", c.nodeInfo.NodeType, c.nodeInfo.NodeID, len(deleted), len(added))
	}
	c.userList = newUserInfo
	return nil
}

func (c *Controller) removeOldTag(oldtag string) (err error) {
	err = c.removeInbound(oldtag)
	if err != nil {
		return err
	}
	err = c.removeOutbound(oldtag)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) addNewTag(newNodeInfo *api.NodeInfo) (err error) {
	if newNodeInfo.NodeType != "Shadowsocks-Plugin" {
		inboundConfig, err := InboundBuilder(c.config, newNodeInfo, c.Tag)
		if err != nil {
			return err
		}
		err = c.addInbound(inboundConfig)
		if err != nil {

			return err
		}
		outBoundConfig, err := OutboundBuilder(c.config, newNodeInfo, c.Tag)
		if err != nil {

			return err
		}
		err = c.addOutbound(outBoundConfig)
		if err != nil {

			return err
		}

	} else {
		return c.addInboundForSSPlugin(*newNodeInfo)
	}
	return nil
}

func (c *Controller) addInboundForSSPlugin(newNodeInfo api.NodeInfo) (err error) {
	// Shadowsocks-Plugin require a seaperate inbound for other TransportProtocol likes: ws, grpc
	fakeNodeInfo := newNodeInfo
	fakeNodeInfo.TransportProtocol = "tcp"
	fakeNodeInfo.EnableTLS = false
	// Add a regular Shadowsocks inbound and outbound
	inboundConfig, err := InboundBuilder(c.config, &fakeNodeInfo, c.Tag)
	if err != nil {
		return err
	}
	err = c.addInbound(inboundConfig)
	if err != nil {

		return err
	}
	outBoundConfig, err := OutboundBuilder(c.config, &fakeNodeInfo, c.Tag)
	if err != nil {

		return err
	}
	err = c.addOutbound(outBoundConfig)
	if err != nil {

		return err
	}
	// Add an inbound for upper streaming protocol
	fakeNodeInfo = newNodeInfo
	fakeNodeInfo.Port++
	fakeNodeInfo.NodeType = "dokodemo-door"
	dokodemoTag := fmt.Sprintf("dokodemo-door_%s+1", c.Tag)
	inboundConfig, err = InboundBuilder(c.config, &fakeNodeInfo, dokodemoTag)
	if err != nil {
		return err
	}
	err = c.addInbound(inboundConfig)
	if err != nil {

		return err
	}
	outBoundConfig, err = OutboundBuilder(c.config, &fakeNodeInfo, dokodemoTag)
	if err != nil {

		return err
	}
	err = c.addOutbound(outBoundConfig)
	if err != nil {

		return err
	}
	return nil
}

func (c *Controller) addNewUser(userInfo *[]api.UserInfo, nodeInfo *api.NodeInfo) (err error) {
	users := make([]*protocol.User, 0)
	if nodeInfo.NodeType == "V2ray" {
		if nodeInfo.EnableVless {
			users = c.buildVlessUser(userInfo)
		} else {
			var alterID uint16 = 0
			if (c.panelType == "V2board" || c.panelType == "Xflash" || c.panelType == "V2Raysocks" || c.panelType == "AikoVPN") && len(*userInfo) > 0 {
				// use latest userInfo
				alterID = (*userInfo)[0].AlterID
			} else {
				alterID = nodeInfo.AlterID
			}
			users = c.buildVmessUser(userInfo, alterID)
		}
	} else if nodeInfo.NodeType == "Trojan" {
		users = c.buildTrojanUser(userInfo)
	} else if nodeInfo.NodeType == "Shadowsocks" {
		users = c.buildSSUser(userInfo, nodeInfo.CypherMethod)
	} else if nodeInfo.NodeType == "Shadowsocks-Plugin" {
		users = c.buildSSPluginUser(userInfo)
	} else {
		return fmt.Errorf("unsupported node type: %s", nodeInfo.NodeType)
	}
	err = c.addUsers(users, c.Tag)
	if err != nil {
		return err
	}
	log.Printf("[%s: %d] Added %d new users", c.nodeInfo.NodeType, c.nodeInfo.NodeID, len(*userInfo))
	return nil
}

func compareUserList(old, new *[]api.UserInfo) (deleted, added []api.UserInfo) {
	msrc := make(map[api.UserInfo]byte) //Index by source array
	mall := make(map[api.UserInfo]byte) //Indexing all elements of source + destination

	var set []api.UserInfo //intersection

	//1.source array to build map
	for _, v := range *old {
		msrc[v] = 0
		mall[v] = 0
	}
	//2. In the target array, if it cannot be stored, that is, repeated elements, all sets that cannot be stored are unions
	for _, v := range *new {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) { // The length changes, that is, it can be stored
			l = len(mall)
		} else { // Can't save, enter union
			set = append(set, v)
		}
	}
	//3. Traverse the intersection, look for it in the union, delete it from the union if found, and delete it, it is the complement (ie union-intersection = all changed elements)
	for _, v := range set {
		delete(mall, v)
	}
	//4. At this point, mall is a complement, all elements are found in the source, and if found, they are deleted. If they are not found, they must be found in the destination array, that is, newly added
	for v := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}

	return deleted, added
}

func limitUser(c *Controller, user api.UserInfo, silentUsers *[]api.UserInfo) {
	c.limitedUsers[user] = LimitInfo{
		end:              time.Now().Unix() + int64(c.config.DynamicSpeedLimitConfig.LimitDuration*60),
		originSpeedLimit: user.SpeedLimit,
	}
	log.Printf("    User: %s Speed: %d End: %s", user.Email, user.SpeedLimit, time.Unix(c.limitedUsers[user].end, 0).Format("01-02 15:04:05"))
	user.SpeedLimit = uint64(c.config.DynamicSpeedLimitConfig.LimitSpeed) * 1024 * 1024 / 8
	*silentUsers = append(*silentUsers, user)
}

func (c *Controller) userInfoMonitor() (err error) {
	// Get server status
	CPU, Mem, Disk, Uptime, err := serverstatus.GetSystemInfo()
	if err != nil {
		log.Print(err)
	}
	// Unlock users
	if c.config.DynamicSpeedLimitConfig.Limit > 0 && len(c.limitedUsers) > 0 {
		log.Printf("Limited users:")
		toReleaseUsers := make([]api.UserInfo, 0)
		for user, limitInfo := range c.limitedUsers {
			if time.Now().Unix() > limitInfo.end {
				user.SpeedLimit = limitInfo.originSpeedLimit
				toReleaseUsers = append(toReleaseUsers, user)
				log.Printf("    User: %s Speed: %d End: nil (Unlimit)", user.Email, user.SpeedLimit)
				delete(c.limitedUsers, user)
			} else {
				log.Printf("    User: %s Speed: %d End: %s", user.Email, user.SpeedLimit, time.Unix(c.limitedUsers[user].end, 0).Format("01-02 15:04:05"))
			}
		}
		if len(toReleaseUsers) > 0 {
			if err := c.UpdateInboundLimiter(c.Tag, &toReleaseUsers); err != nil {
				log.Print(err)
			}
		}
	}
	err = c.apiClient.ReportNodeStatus(
		&api.NodeStatus{
			CPU:    CPU,
			Mem:    Mem,
			Disk:   Disk,
			Uptime: Uptime,
		})
	if err != nil {
		log.Print(err)
	}

	// Get User traffic
	var userTraffic []api.UserTraffic
	var upCounterList []stats.Counter
	var downCounterList []stats.Counter
	AutoSpeedLimit := int64(c.config.DynamicSpeedLimitConfig.Limit)
	UpdatePeriodic := int64(c.config.UpdatePeriodic)
	limitedUsers := make([]api.UserInfo, 0)
	for _, user := range *c.userList {
		up, down, upCounter, downCounter := c.getTraffic(c.buildUserTag(&user))
		if up > 0 || down > 0 {
			// Over speed users
			if AutoSpeedLimit > 0 {
				if down > AutoSpeedLimit*1024*1024*UpdatePeriodic/8 {
					if _, ok := c.limitedUsers[user]; !ok {
						if c.config.DynamicSpeedLimitConfig.WarnTimes == 0 {
							limitUser(c, user, &limitedUsers)
						} else {
							c.warnedUsers[user] += 1
							if c.warnedUsers[user] > c.config.DynamicSpeedLimitConfig.WarnTimes {
								limitUser(c, user, &limitedUsers)
								delete(c.warnedUsers, user)
							}
						}
					}
				} else {
					delete(c.warnedUsers, user)
				}
			}
			userTraffic = append(userTraffic, api.UserTraffic{
				UID:      user.UID,
				Email:    user.Email,
				Upload:   up,
				Download: down})

			if upCounter != nil {
				upCounterList = append(upCounterList, upCounter)
			}
			if downCounter != nil {
				downCounterList = append(downCounterList, downCounter)
			}
		} else {
			delete(c.warnedUsers, user)
		}
	}

	if len(limitedUsers) > 0 {
		if err := c.UpdateInboundLimiter(c.Tag, &limitedUsers); err != nil {
			log.Print(err)
		}
	}

	if len(userTraffic) > 0 {
		var err error // Define an empty error
		if !c.config.DisableUploadTraffic {
			err = c.apiClient.ReportUserTraffic(&userTraffic)
		}
		// If report traffic error, not clear the traffic
		if err != nil {
			log.Print(err)
		} else {
			c.resetTraffic(&upCounterList, &downCounterList)
		}
	}

	if !c.config.EnableIpRecorder {
		//ClearOnlineIp
		c.ClearOnlineIp(c.Tag)
	}
	// Report Illegal user
	if detectResult, err := c.GetDetectResult(c.Tag); err != nil {
		log.Print(err)
	} else if len(*detectResult) > 0 {
		if err = c.apiClient.ReportIllegal(detectResult); err != nil {
			log.Print(err)
		} else {
			log.Printf("[%s: %d] Report %d illegal behaviors", c.nodeInfo.NodeType, c.nodeInfo.NodeID, len(*detectResult))
		}

	}
	runtime.GC()
	return nil
}
func (c *Controller) onlineIpReport() (err error) {
	onlineIp, err := c.dispatcher.Limiter.GetOnlineUserIp(c.Tag)
	if err != nil {
		log.Print(err)
		return nil
	}
	rsp, err := resty.New().SetTimeout(time.Duration(c.config.IpRecorderConfig.Timeout) * time.Second).
		R().
		SetBody(onlineIp).
		Post(c.config.IpRecorderConfig.Url +
			"/api/v1/SyncOnlineIp?token=" +
			c.config.IpRecorderConfig.Token)
	if err != nil {
		log.Print(err)
		c.dispatcher.Limiter.ClearOnlineUserIP(c.Tag)
		return nil
	}
	log.Printf("[%s: %d] report %d online Ip", c.nodeInfo.NodeType, c.nodeInfo.NodeID, len(onlineIp))
	if rsp.StatusCode() == 200 {
		onlineIp = []limiter.UserIp{}
		err := json.Unmarshal(rsp.Body(), &onlineIp)
		if err != nil {
			log.Print(err)
			c.dispatcher.Limiter.ClearOnlineUserIP(c.Tag)
			return nil
		}
		if c.config.IpRecorderConfig.EnableIpSync {
			c.dispatcher.Limiter.UpdateOnlineUserIP(c.Tag, onlineIp)
			log.Printf("[%s: %d] update %d online Ip", c.nodeInfo.NodeType, c.nodeInfo.NodeID, len(onlineIp))
		}
	} else {
		c.dispatcher.Limiter.ClearOnlineUserIP(c.Tag)
	}
	runtime.GC()
	return nil
}

func (c *Controller) buildNodeTag() string {
	return fmt.Sprintf("%s_%s_%d", c.nodeInfo.NodeType, c.config.ListenIP, c.nodeInfo.Port)
}
