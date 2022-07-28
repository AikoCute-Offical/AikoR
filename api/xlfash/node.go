package xflash

import (
	"bufio"
	"bytes"
	md52 "crypto/md5"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/xtls/xray-core/infra/conf"
)

type DetectRule struct {
	ID      int
	Pattern *regexp.Regexp
}
type DetectResult struct {
	UID    int
	RuleID int
}

// readLocalRuleList reads the local rule list file
func readLocalRuleList(path string) (LocalRuleList []DetectRule) {
	LocalRuleList = make([]DetectRule, 0)
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
			LocalRuleList = append(LocalRuleList, DetectRule{
				ID:      -1,
				Pattern: regexp.MustCompile(fileScanner.Text()),
			})
		}
		// handle first encountered error while reading
		if err := fileScanner.Err(); err != nil {
			log.Fatalf("Error while reading file: %s", err)
			return []DetectRule{}
		}
		file.Close()
	}

	return LocalRuleList
}

type NodeInfo struct {
	DeviceLimit int
	SpeedLimit  uint64
	NodeType    string
	NodeId      int
	TLSType     string
	EnableVless bool
	EnableTls   bool
	//EnableSS2022 bool
	V2ray  *V2rayConfig
	Trojan *TrojanConfig
	SS     *SSConfig
}

type SSConfig struct {
	Port              int    `json:"port"`
	TransportProtocol string `json:"transportProtocol"`
	CypherMethod      string `json:"cypher"`
}
type V2rayConfig struct {
	Inbounds []conf.InboundDetourConfig `json:"inbounds"`
	Routing  *struct {
		Rules json.RawMessage `json:"rules"`
	} `json:"routing"`
}

type Rule struct {
	Type        string   `json:"type"`
	InboundTag  string   `json:"inboundTag,omitempty"`
	OutboundTag string   `json:"outboundTag"`
	Domain      []string `json:"domain,omitempty"`
	Protocol    []string `json:"protocol,omitempty"`
}

type TrojanConfig struct {
	LocalPort         int           `json:"local_port"`
	Password          []interface{} `json:"password"`
	TransportProtocol string
	Ssl               struct {
		Sni string `json:"sni"`
	} `json:"ssl"`
}

// GetNodeInfo will pull NodeInfo Config from sspanel
func (c *Client) GetNodeInfo() (nodeInfo *NodeInfo, err error) {
	var path string
	var res *resty.Response
	switch c.NodeType {
	case "V2ray":
		path = "/api/v1/server/Deepbwork/config"
	case "Trojan":
		path = "/api/v1/server/TrojanTidalab/config"
	case "Shadowsocks":
		if nodeInfo, err = c.ParseSSNodeResponse(); err == nil {
			return nodeInfo, nil
		} else {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported Node type: %s", c.NodeType)
	}
	res, err = c.client.R().
		SetQueryParam("local_port", "1").
		ForceContentType("application/json").
		Get(path)
	err = c.checkResponse(res, path, err)
	if err != nil {
		return nil, err
	}
	c.access.Lock()
	defer c.access.Unlock()
	switch c.NodeType {
	case "V2ray":
		i := bytes.Index(res.Body(), []byte("outbo"))
		md := md52.Sum(res.Body()[:i])
		nodeNotIsChange := true
		if c.NodeInfoRspMd5 == [16]byte{} {
			nodeNotIsChange = false
			c.NodeInfoRspMd5 = md
		} else {
			if c.NodeInfoRspMd5 != md {
				nodeNotIsChange = false
				c.NodeInfoRspMd5 = md
			}
		}
		md2 := md52.Sum(res.Body()[i:])
		ruleIsChange := false
		if c.NodeRuleRspMd5 != md2 {
			ruleIsChange = true
			c.NodeRuleRspMd5 = md2
		}
		nodeInfo, err = c.ParseV2rayNodeResponse(res.Body(), nodeNotIsChange, ruleIsChange)
	case "Trojan":
		md := md52.Sum(res.Body())
		if c.NodeInfoRspMd5 != [16]byte{} {
			if c.NodeInfoRspMd5 == md {
				return nil, nil
			}
		}
		c.NodeInfoRspMd5 = md
		nodeInfo, err = c.ParseTrojanNodeResponse(res.Body())
	}
	return nodeInfo, nil
}

func (c *Client) GetNodeRule() (*[]DetectRule, *[]string, error) {
	ruleList := c.LocalRuleList
	if c.NodeType != "V2ray" || c.RemoteRuleCache == nil {
		return &ruleList, nil, nil
	}
	// V2board only support the rule for v2ray
	// fix: reuse config response
	c.access.Lock()
	defer c.access.Unlock()
	if len(*c.RemoteRuleCache) >= 2 {
		for i, rule := range (*c.RemoteRuleCache)[1].Domain {
			ruleListItem := DetectRule{
				ID:      i,
				Pattern: regexp.MustCompile(rule),
			}
			ruleList = append(ruleList, ruleListItem)
		}
	}
	var protocolList []string
	if len(*c.RemoteRuleCache) >= 3 {
		for _, str := range (*c.RemoteRuleCache)[2].Protocol {
			protocolList = append(protocolList, str)
		}
	}
	c.RemoteRuleCache = nil
	return &ruleList, &protocolList, nil
}

// ParseTrojanNodeResponse parse the response for the given nodeinfor format
func (c *Client) ParseTrojanNodeResponse(body []byte) (*NodeInfo, error) {
	node := &NodeInfo{Trojan: &TrojanConfig{}}
	var err = json.Unmarshal(body, node.Trojan)
	if err != nil {
		return nil, fmt.Errorf("unmarshal nodeinfo error: %s", err)
	}
	node.SpeedLimit = uint64(c.SpeedLimit * 1000000 / 8)
	node.DeviceLimit = c.DeviceLimit
	node.NodeId = c.NodeID
	node.NodeType = c.NodeType
	node.Trojan.TransportProtocol = "tcp"
	return node, nil
}

// ParseSSNodeResponse parse the response for the given nodeinfor format
func (c *Client) ParseSSNodeResponse() (*NodeInfo, error) {
	var port int
	var method string
	userInfo, err := c.GetUserList()
	if err != nil {
		return nil, err
	}
	if len(*userInfo) > 0 {
		port = (*userInfo)[0].Port
		method = (*userInfo)[0].Cipher
	} else {
		return nil, fmt.Errorf("shadowsocks node need a active user")
	}

	if err != nil {
		return nil, err
	}
	node := &NodeInfo{
		SpeedLimit:  uint64(c.SpeedLimit * 1000000 / 8),
		DeviceLimit: c.DeviceLimit,
		//EnableSS2022: c.EnableSS2022,
		NodeType: c.NodeType,
		NodeId:   c.NodeID,
		SS: &SSConfig{
			Port:              port,
			TransportProtocol: "tcp",
			CypherMethod:      method,
		},
	}
	return node, nil
}

// ParseV2rayNodeResponse parse the response for the given nodeinfor format
func (c *Client) ParseV2rayNodeResponse(body []byte, notParseNode, parseRule bool) (*NodeInfo, error) {
	if notParseNode && !parseRule {
		return nil, nil
	}
	node := &NodeInfo{V2ray: &V2rayConfig{}}
	err := json.Unmarshal(body, node.V2ray)
	if err != nil {
		return nil, fmt.Errorf("unmarshal nodeinfo error: %s", err)
	}
	if parseRule {
		c.RemoteRuleCache = &[]Rule{}
		err := json.Unmarshal(node.V2ray.Routing.Rules, c.RemoteRuleCache)
		if err != nil {
			log.Println(err)
		}
		if notParseNode {
			return nil, nil
		}
	}
	node.V2ray.Routing = nil
	node.SpeedLimit = uint64(c.SpeedLimit * 1000000 / 8)
	node.DeviceLimit = c.DeviceLimit
	node.NodeType = c.NodeType
	node.NodeId = c.NodeID
	if c.EnableXTLS {
		node.TLSType = "xtls"
	} else {
		node.TLSType = "tls"
	}
	node.EnableVless = c.EnableVless
	node.EnableTls = node.V2ray.Inbounds[0].StreamSetting.Security == "tls"
	return node, nil
}
