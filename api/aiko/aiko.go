package aiko

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/go-resty/resty/v2"

	"github.com/AikoCute-Offical/AikoR/api"
)

// APIClient create an api client to the panel.
type APIClient struct {
	client        *resty.Client
	APIHost       string
	NodeID        int
	Key           string
	NodeType      string
	EnableVless   bool
	EnableXTLS    bool
	SpeedLimit    float64
	DeviceLimit   int
	LocalRuleList []api.DetectRule
	resp          atomic.Value
	eTag          string
}

// New create an api instance
func New(apiConfig *api.Config) *APIClient {
	client := resty.New()
	client.SetRetryCount(3)
	if apiConfig.Timeout > 0 {
		client.SetTimeout(time.Duration(apiConfig.Timeout) * time.Second)
	} else {
		client.SetTimeout(5 * time.Second)
	}
	client.OnError(func(req *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			// v.Response contains the last response from the server
			// v.Err contains the original error
			log.Print(v.Err)
		}
	})
	client.SetBaseURL(apiConfig.APIHost)
	// Create Key for each requests
	client.SetQueryParams(map[string]string{
		"node_id":   strconv.Itoa(apiConfig.NodeID),
		"node_type": strings.ToLower(apiConfig.NodeType),
		"token":     apiConfig.Key,
	})
	// Read local rule list
	localRuleList := readLocalRuleList(apiConfig.RuleListPath)
	apiClient := &APIClient{
		client:        client,
		NodeID:        apiConfig.NodeID,
		Key:           apiConfig.Key,
		APIHost:       apiConfig.APIHost,
		NodeType:      apiConfig.NodeType,
		EnableVless:   apiConfig.EnableVless,
		EnableXTLS:    apiConfig.EnableXTLS,
		SpeedLimit:    apiConfig.SpeedLimit,
		DeviceLimit:   apiConfig.DeviceLimit,
		LocalRuleList: localRuleList,
	}
	return apiClient
}

// readLocalRuleList reads the local rule list file
func readLocalRuleList(path string) (LocalRuleList []api.DetectRule) {
	LocalRuleList = make([]api.DetectRule, 0)

	if path != "" {
		// open the file
		file, err := os.Open(path)

		// handle errors while opening
		if err != nil {
			log.Printf("Error when opening file: %s", err)
			return LocalRuleList
		}

		fileScanner := bufio.NewScanner(file)

		// read line by line
		for fileScanner.Scan() {
			LocalRuleList = append(LocalRuleList, api.DetectRule{
				ID:      -1,
				Pattern: regexp.MustCompile(fileScanner.Text()),
			})
		}
		// handle first encountered error while reading
		if err := fileScanner.Err(); err != nil {
			log.Fatalf("Error while reading file: %s", err)
			return
		}

		file.Close()
	}

	return LocalRuleList
}

// Describe return a description of the client
func (c *APIClient) Describe() api.ClientInfo {
	return api.ClientInfo{APIHost: c.APIHost, NodeID: c.NodeID, Key: c.Key, NodeType: c.NodeType}
}

// Debug set the client debug for client
func (c *APIClient) Debug() {
	c.client.SetDebug(true)
}

func (c *APIClient) assembleURL(path string) string {
	return c.APIHost + path
}

func (c *APIClient) parseResponse(res *resty.Response, path string, err error) (*simplejson.Json, error) {
	if err != nil {
		return nil, fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 399 {
		body := res.Body()
		return nil, fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}
	rtn, err := simplejson.NewJson(res.Body())
	if err != nil {
		return nil, fmt.Errorf("ret %s invalid", res.String())
	}
	return rtn, nil
}

// GetNodeInfo will pull NodeInfo Config from panel
func (c *APIClient) GetNodeInfo() (nodeInfo *api.NodeInfo, err error) {
	path := "/api/v1/server/UniProxy/config"

	res, err := c.client.R().
		ForceContentType("application/json").
		Get(path)

	response, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}

	c.resp.Store(response)

	switch c.NodeType {
	case "V2ray":
		nodeInfo, err = c.parseV2rayNodeResponse(response)
	case "Trojan":
		nodeInfo, err = c.parseTrojanNodeResponse(response)
	case "Shadowsocks":
		nodeInfo, err = c.parseSSNodeResponse(response)
	default:
		return nil, fmt.Errorf("unsupported Node type: %s", c.NodeType)
	}

	if err != nil {
		res, _ := response.MarshalJSON()
		return nil, fmt.Errorf("Parse node info failed: %s, \nError: %s", string(res), err)
	}

	return nodeInfo, nil
}

// GetUserList will pull user form panel
func (c *APIClient) GetUserList() (UserList *[]api.UserInfo, err error) {
	path := "/api/v1/server/UniProxy/user"

	switch c.NodeType {
	case "V2ray", "Trojan", "Shadowsocks":
		break
	default:
		return nil, fmt.Errorf("unsupported Node type: %s", c.NodeType)
	}

	res, err := c.client.R().
		SetHeader("If-None-Match", c.eTag).
		ForceContentType("application/json").
		Get(path)

	// Etag identifier for a specific version of a resource. StatusCode = 304 means no changed
	if res.StatusCode() == 304 {
		return nil, errors.New("users no change")
	}
	// update etag
	if res.Header().Get("Etag") != "" && res.Header().Get("Etag") != c.eTag {
		c.eTag = res.Header().Get("Etag")
	}

	response, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}

	numOfUsers := len(response.Get("users").MustArray())
	userList := make([]api.UserInfo, numOfUsers)
	for i := 0; i < numOfUsers; i++ {
		user := response.Get("users").GetIndex(i)
		u := api.UserInfo{
			UID:  user.Get("id").MustInt(),
			UUID: user.Get("uuid").MustString(),
		}

		// Support 1.7.1 speed limit
		if c.SpeedLimit > 0 {
			u.SpeedLimit = uint64(c.SpeedLimit * 1000000 / 8)
		} else {
			u.SpeedLimit = uint64(user.Get("speed_limit").MustUint64() * 1000000 / 8)
		}

		u.DeviceLimit = c.DeviceLimit // todo waiting v2board send configuration
		u.Email = u.UUID + "@v2board.user"
		if c.NodeType == "Shadowsocks" {
			u.Passwd = u.UUID
		}
		userList[i] = u
	}

	return &userList, nil
}

// ReportUserTraffic reports the user traffic
func (c *APIClient) ReportUserTraffic(userTraffic *[]api.UserTraffic) error {
	path := "/api/v1/server/UniProxy/push"

	// json structure: {uid1: [u, d], uid2: [u, d], uid1: [u, d], uid3: [u, d]}
	data := make(map[int][]int64, len(*userTraffic))
	for _, traffic := range *userTraffic {
		data[traffic.UID] = []int64{traffic.Upload, traffic.Download}
	}

	res, err := c.client.R().
		SetBody(data).
		ForceContentType("application/json").
		Post(path)
	_, err = c.parseResponse(res, path, err)
	if err != nil {
		return err
	}

	return nil
}

// GetNodeRule implements the API interface
func (c *APIClient) GetNodeRule() (*[]api.DetectRule, error) {
	ruleList := c.LocalRuleList

	nodeInfoResponse := c.resp.Load().(*simplejson.Json)
	for i, rule := range nodeInfoResponse.Get("routes").MustArray() {
		r := rule.(map[string]any)
		if r["action"] == "block" {
			ruleListItem := api.DetectRule{
				ID:      i,
				Pattern: regexp.MustCompile(strings.TrimPrefix(r["match"].(string), "regexp:")),
			}
			ruleList = append(ruleList, ruleListItem)
		}
	}

	return &ruleList, nil
}

// ReportNodeStatus implements the API interface
func (c *APIClient) ReportNodeStatus(nodeStatus *api.NodeStatus) (err error) {
	return nil
}

// ReportNodeOnlineUsers implements the API interface
func (c *APIClient) ReportNodeOnlineUsers(onlineUserList *[]api.OnlineUser) error {
	return nil
}

// ReportIllegal implements the API interface
func (c *APIClient) ReportIllegal(detectResultList *[]api.DetectResult) error {
	return nil
}

// parseTrojanNodeResponse parse the response for the given nodeInfo format
func (c *APIClient) parseTrojanNodeResponse(nodeInfoResponse *simplejson.Json) (*api.NodeInfo, error) {
	var TLSType = "tls"
	if c.EnableXTLS {
		TLSType = "xtls"
	}

	// Create GeneralNodeInfo
	nodeInfo := &api.NodeInfo{
		NodeType:          c.NodeType,
		NodeID:            c.NodeID,
		Port:              uint32(nodeInfoResponse.Get("server_port").MustUint64()),
		TransportProtocol: "tcp",
		EnableTLS:         true,
		TLSType:           TLSType,
		Host:              nodeInfoResponse.Get("host").MustString(),
		ServiceName:       nodeInfoResponse.Get("server_name").MustString(),
	}
	return nodeInfo, nil
}

// parseSSNodeResponse parse the response for the given nodeInfo format
func (c *APIClient) parseSSNodeResponse(nodeInfoResponse *simplejson.Json) (*api.NodeInfo, error) {
	// Create GeneralNodeInfo
	return &api.NodeInfo{
		NodeType:          c.NodeType,
		NodeID:            c.NodeID,
		Port:              uint32(nodeInfoResponse.Get("server_port").MustUint64()),
		TransportProtocol: "tcp",
		CypherMethod:      nodeInfoResponse.Get("cipher").MustString(),
		ServerKey:         nodeInfoResponse.Get("server_key").MustString(), // shadowsocks2022 share key
	}, nil
}

// parseV2rayNodeResponse parse the response for the given nodeInfo format
func (c *APIClient) parseV2rayNodeResponse(nodeInfoResponse *simplejson.Json) (*api.NodeInfo, error) {
	var (
		TLSType                 = "tls"
		path, host, serviceName string
		header                  json.RawMessage
		enableTLS               bool
		alterID                 uint16 = 0
	)

	if c.EnableXTLS {
		TLSType = "xtls"
	}

	transportProtocol := nodeInfoResponse.Get("network").MustString()

	switch transportProtocol {
	case "ws":
		path = nodeInfoResponse.Get("networkSettings").Get("path").MustString()
		host = nodeInfoResponse.Get("networkSettings").Get("headers").Get("Host").MustString()
	case "grpc":
		if data, ok := nodeInfoResponse.Get("networkSettings").CheckGet("serviceName"); ok {
			serviceName = data.MustString()
		}
	case "tcp":
		if data, ok := nodeInfoResponse.Get("networkSettings").CheckGet("headers"); ok {
			if httpHeader, err := data.MarshalJSON(); err != nil {
				return nil, err
			} else {
				header = httpHeader
			}
		}
	}

	if nodeInfoResponse.Get("tls").MustInt() == 1 {
		enableTLS = true
	}

	// Create GeneralNodeInfo
	return &api.NodeInfo{
		NodeType:          c.NodeType,
		NodeID:            c.NodeID,
		Port:              uint32(nodeInfoResponse.Get("server_port").MustUint64()),
		AlterID:           alterID,
		TransportProtocol: transportProtocol,
		EnableTLS:         enableTLS,
		TLSType:           TLSType,
		Path:              path,
		Host:              host,
		EnableVless:       c.EnableVless,
		ServiceName:       serviceName,
		Header:            header,
	}, nil
}
