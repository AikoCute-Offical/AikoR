package dev

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AikoCute-Offical/AikoR/api"
)

// Describe return a description of the client
func (c *APIClient) Describe() api.ClientInfo {
	return api.ClientInfo{APIHost: c.APIHost, NodeID: c.NodeID, Key: c.Key, NodeType: c.NodeType}
}

// Debug set the client debug for client
func (c *APIClient) Debug() {
	c.client.SetDebug(true)
}

// GetNodeInfo will pull NodeInfo Config from panel
func (c *APIClient) GetNodeInfo() (nodeInfo *api.NodeInfo, err error) {
	server := new(serverConfig)
	path := "/api/v1/server/Aiko/config"

	res, err := c.client.R().
		ForceContentType("application/json").
		Get(path)

	nodeInfoResp, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}
	b, _ := nodeInfoResp.Encode()
	json.Unmarshal(b, server)

	if server.ServerPort == 0 {
		return nil, errors.New("server port must > 0")
	}

	c.resp.Store(server)

	switch c.NodeType {
	case "V2ray":
		nodeInfo, err = c.parseV2rayNodeResponse(server)
	case "Trojan":
		nodeInfo, err = c.parseTrojanNodeResponse(server)
	case "Shadowsocks":
		nodeInfo, err = c.parseSSNodeResponse(server)
	default:
		return nil, fmt.Errorf("unsupported node type: %s", c.NodeType)
	}

	if err != nil {
		return nil, fmt.Errorf("parse node info failed: %s, \nError: %v", res.String(), err)
	}

	return nodeInfo, nil
}
