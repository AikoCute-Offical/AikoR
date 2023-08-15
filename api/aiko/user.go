package aiko

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AikoCute-Offical/AikoR/api"
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
		SetHeader("If-None-Match", c.eTags["users"]).
		ForceContentType("application/json").
		Get(path)

	// Etag identifier for a specific version of a resource. StatusCode = 304 means no changed
	if res.StatusCode() == 304 {
		return nil, errors.New(api.UserNotModified)
	}
	// update etag
	if res.Header().Get("Etag") != "" && res.Header().Get("Etag") != c.eTags["users"] {
		c.eTags["users"] = res.Header().Get("Etag")
	}

	usersResp, err := c.parseResponse(res, path, err)
	if err != nil {
		return nil, err
	}
	b, _ := usersResp.Get("users").Encode()
	json.Unmarshal(b, &users)
	if len(users) == 0 {
		return nil, errors.New("users is null")
	}

	userList := make([]api.UserInfo, len(users))
	for i := 0; i < len(users); i++ {
		u := api.UserInfo{
			UID:  users[i].Id,
			UUID: users[i].Uuid,
		}

		if c.SpeedLimit > 0 {
			u.SpeedLimit = uint64(c.SpeedLimit * 1000000 / 8)
		} else {
			u.SpeedLimit = uint64(users[i].SpeedLimit * 1000000 / 8)
		}

		if c.DeviceLimit > 0 {
			u.DeviceLimit = c.DeviceLimit
		} else {
			u.DeviceLimit = users[i].DeviceLimit
		}

		u.Email = u.UUID + "@aikopanel.user"
		if c.NodeType == "Shadowsocks" {
			u.Passwd = u.UUID
		}
		userList[i] = u
	}

	return &userList, nil
}
