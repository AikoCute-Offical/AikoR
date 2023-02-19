package dev

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AikoCute-Offical/AikoR/api"
	"github.com/bitly/go-simplejson"
	"github.com/go-resty/resty/v2"
)

// GetUserList will pull user form panel
func (c *APIClient) GetUserList() (UserList *[]api.UserInfo, err error) {
	var users []*user
	path := "/api/v1/server/Aiko/user"

	switch c.NodeType {
	case "V2ray", "Trojan", "Shadowsocks":
		break
	default:
		return nil, fmt.Errorf("unsupported node type: %s", c.NodeType)
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

	usersResp, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}
	b, _ := usersResp.Get("users").Encode()
	json.Unmarshal(b, &users)

	userList := make([]api.UserInfo, len(users))
	for i := 0; i < len(users); i++ {
		u := api.UserInfo{
			UID:  users[i].Id,
			UUID: users[i].Uuid,
		}

		// Support 1.7.1 speed limit
		if c.SpeedLimit > 0 {
			u.SpeedLimit = uint64(c.SpeedLimit * 1000000 / 8)
		} else {
			u.SpeedLimit = uint64(users[i].SpeedLimit * 1000000 / 8)
		}

		// Support 1.7.1 device limit
		if c.DeviceLimit > 0 {
			u.DeviceLimit = c.DeviceLimit
		} else {
			u.DeviceLimit = users[i].DeviceLimit
		}

		u.Email = u.UUID + "@v2board.user"
		if c.NodeType == "Shadowsocks" {
			u.Passwd = u.UUID
		}
		userList[i] = u
	}

	return &userList, nil
}

func (c *APIClient) parseResponse(res *resty.Response, path string, err error) (*simplejson.Json, error) {
	if err != nil {
		return nil, fmt.Errorf("request %s failed: %v", c.assembleURL(path), err)
	}

	if res.StatusCode() > 399 {
		return nil, fmt.Errorf("request %s failed: %s, %v", c.assembleURL(path), res.String(), err)
	}

	rtn, err := simplejson.NewJson(res.Body())
	if err != nil {
		return nil, fmt.Errorf("ret %s invalid", res.String())
	}

	return rtn, nil
}
