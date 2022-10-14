package xflash

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/AikoCute-Offical/AikoR/api"
	"github.com/bitly/go-simplejson"
	"github.com/go-resty/resty/v2"
)

// Client APIClient create a api client to the panel.
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
	ConfigResp    *simplejson.Json
	access        sync.Mutex
}

// New creat a api instance
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
			log.Printf("request %s failed: %s", req.URL, v.Err)
		}
	})
	client.SetBaseURL(apiConfig.APIHost)
	// Create Key for each requests
	client.SetQueryParam("key", apiConfig.Key)
	client.SetQueryParams(map[string]string{
		"node_id": strconv.Itoa(apiConfig.NodeID),
		"token":   apiConfig.Key,
	})

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

func readLocalRuleList(path string) (LocalRuleList []api.DetectRule) {

	LocalRuleList = make([]api.DetectRule, 0)
	if path != "" {
		// open the file
		file, err := os.Open(path)

		//handle errors while opening
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
			return make([]api.DetectRule, 0)
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

// GetNodeInfo will pull NodeInfo Config from sspanel
func (c *APIClient) GetNodeInfo() (nodeInfo *NodeInfo, err error) {
	var path = "/api/v1/server/trojan/config"
	res, err := c.client.R().
		ForceContentType("application/json").
		SetQueryParam("node_id", strconv.Itoa(c.NodeID)).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return nil, fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}

	var repNodeInfo RepNodeInfo
	if err := json.Unmarshal(res.Body(), &repNodeInfo); err != nil {
		return nil, fmt.Errorf("parse node info failed: %s", err)
	}

	if len(repNodeInfo.Message) > 0 {
		return nil, fmt.Errorf("api error, message: %s", repNodeInfo.Message)
	}
	return repNodeInfo.Data, nil
}

// GetUserList will pull user form sspanel
func (c *APIClient) GetUserList() (UserList *[]UserInfo, err error) {
	var path = "/api/v1/server/trojan/users"
	res, err := c.client.R().
		SetQueryParam("node_id", strconv.Itoa(c.NodeID)).
		ForceContentType("application/json").
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return nil, fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}

	var repUserList RepUserList
	if err := json.Unmarshal(res.Body(), &repUserList); err != nil {
		return nil, fmt.Errorf("parse node info failed: %s", err)
	}

	if len(repUserList.Message) > 0 {
		return nil, fmt.Errorf("api error, message: %s", repUserList.Message)
	}

	return repUserList.Data, nil
}

// ReportUserTraffic reports the user traffic
func (c *APIClient) ReportUserTraffic(userTraffic []*UserTraffic) error {
	var path = "/api/v1/server/trojan/submit"

	res, err := c.client.R().
		SetQueryParam("node_id", strconv.Itoa(c.NodeID)).
		SetBody(userTraffic).
		ForceContentType("application/json").
		Post(path)
	if err != nil {
		return fmt.Errorf("request %s failed: %s", c.assembleURL(path), err)
	}

	if res.StatusCode() > 400 {
		body := res.Body()
		return fmt.Errorf("request %s failed: %s, %s", c.assembleURL(path), string(body), err)
	}

	var repUserTraffic RepUserTraffic
	if err := json.Unmarshal(res.Body(), &repUserTraffic); err != nil {
		return fmt.Errorf("parse node info failed: %s", err)
	}
	if len(repUserTraffic.Message) > 0 {
		return fmt.Errorf("api error, message: %s", repUserTraffic.Message)
	}
	return nil
}
