// Package limiter is to control the links that go into the dispather
package limiter

import (
	"fmt"
	"github.com/AikoCute-Offical/AikoR/api"
	"github.com/juju/ratelimit"
	"sync"
	"time"
)

type UserInfo struct {
	UID         int
	SpeedLimit  uint64
	DeviceLimit int
}

type InboundInfo struct {
	Tag            string
	NodeSpeedLimit uint64
	UserInfo       *sync.Map // Key: Email value: UserInfo
	BucketHub      *sync.Map // key: Email, value: *ratelimit.Bucket
	UserOnlineIP   *sync.Map // Key: Email Value: *sync.Map: Key: IP, Value: UID or bool
	EnableReorder  bool
}

type Limiter struct {
	InboundInfo *sync.Map // Key: Tag, Value: *InboundInfo
}

func New() *Limiter {
	return &Limiter{
		InboundInfo: new(sync.Map),
	}
}

type UserIp struct {
	Uid int      `json:"Uid"`
	IPs []string `json:"Ips"`
}

func (l *Limiter) GetOnlineUserIp(tag string) ([]UserIp, error) {
	if value, ok := l.InboundInfo.Load(tag); ok {
		inboundInfo := value.(*InboundInfo)
		// Clear Speed Limiter bucket for users who are not online
		inboundInfo.BucketHub.Range(func(key, value interface{}) bool {
			if _, exists := inboundInfo.UserOnlineIP.Load(key.(string)); !exists {
				inboundInfo.BucketHub.Delete(key.(string))
			}
			return true
		})
		onlineUser := make([]UserIp, 0)
		var ipMap *sync.Map
		inboundInfo.UserOnlineIP.Range(func(key, value interface{}) bool {
			ipMap = value.(*sync.Map)
			var ip []string
			ipMap.Range(func(key, v interface{}) bool {
				if v.(bool) {
					ip = append(ip, key.(string))
				}
				return true
			})
			if len(ip) > 0 {
				if u, ok := inboundInfo.UserInfo.Load(key.(string)); ok {
					onlineUser = append(onlineUser, UserIp{
						Uid: u.(UserInfo).UID,
						IPs: ip,
					})
				}
			}
			return true
		})
		if len(onlineUser) == 0 {
			return nil, nil
		}
		return onlineUser, nil
	} else {
		return nil, fmt.Errorf("no such inbound in limiter: %s", tag)
	}
}

func (l *Limiter) UpdateOnlineUserIP(tag string, userIpList []UserIp) {
	if v, ok := l.InboundInfo.Load(tag); ok {
		inboundInfo := v.(*InboundInfo)
		//Clear old IP
		inboundInfo.UserOnlineIP.Range(func(key, value interface{}) bool {
			inboundInfo.UserOnlineIP.Delete(key)
			return true
		})
		// Update User Online IP
		for i := range userIpList {
			ipMap := new(sync.Map)
			for _, userIp := range (userIpList)[i].IPs {
				ipMap.Store(userIp, false)
			}
			inboundInfo.UserOnlineIP.Store((userIpList)[i].Uid, ipMap)
		}
	}
}

func (l *Limiter) ClearOnlineUserIP(tag string) {
	if v, ok := l.InboundInfo.Load(tag); ok {
		inboundInfo := v.(*InboundInfo)
		inboundInfo.BucketHub.Range(func(key, value interface{}) bool {
			if _, exists := inboundInfo.UserOnlineIP.Load(key.(string)); !exists {
				inboundInfo.BucketHub.Delete(key.(string))
			}
			return true
		})
		inboundInfo.UserOnlineIP.Range(func(key, value interface{}) bool {
			inboundInfo.UserOnlineIP.Delete(key)
			return true
		})
	}
}

func (l *Limiter) AddInboundLimiter(tag string, nodeSpeedLimit uint64, userList *[]api.UserInfo) error {
	inboundInfo := &InboundInfo{
		Tag:            tag,
		NodeSpeedLimit: nodeSpeedLimit,
		BucketHub:      new(sync.Map),
		UserOnlineIP:   new(sync.Map),
	}
	userMap := new(sync.Map)
	for _, u := range *userList {
		userMap.Store(fmt.Sprintf("%s|%s|%d", tag, u.Email, u.UID), UserInfo{
			UID:         u.UID,
			SpeedLimit:  u.SpeedLimit,
			DeviceLimit: u.DeviceLimit,
		})
	}
	inboundInfo.UserInfo = userMap
	l.InboundInfo.Store(tag, inboundInfo) // Replace the old inbound info
	return nil
}

func (l *Limiter) UpdateInboundLimiter(tag string, updatedUserList *[]api.UserInfo) error {

	if value, ok := l.InboundInfo.Load(tag); ok {
		inboundInfo := value.(*InboundInfo)
		// Update User info
		for _, u := range *updatedUserList {
			inboundInfo.UserInfo.Store(fmt.Sprintf("%s|%s|%d", tag, u.Email, u.UID), UserInfo{
				UID:         u.UID,
				SpeedLimit:  u.SpeedLimit,
				DeviceLimit: u.DeviceLimit,
			})
			inboundInfo.BucketHub.Delete(fmt.Sprintf("%s|%s|%d", tag, u.Email, u.UID)) // Delete old limiter bucket
		}
	} else {
		return fmt.Errorf("no such inbound in limiter: %s", tag)
	}
	return nil
}

func (l *Limiter) DeleteInboundLimiter(tag string) error {
	l.InboundInfo.Delete(tag)
	return nil
}

func (l *Limiter) GetOnlineDevice(tag string) (*[]api.OnlineUser, error) {
	onlineUser := make([]api.OnlineUser, 0)
	if value, ok := l.InboundInfo.Load(tag); ok {
		inboundInfo := value.(*InboundInfo)
		if inboundInfo.EnableReorder {
			return nil, fmt.Errorf("reorder is enabled, can not get online device")
		}
		// Clear Speed Limiter bucket for users who are not online
		inboundInfo.BucketHub.Range(func(key, value interface{}) bool {
			email := key.(string)
			if _, exists := inboundInfo.UserOnlineIP.Load(email); !exists {
				inboundInfo.BucketHub.Delete(email)
			}
			return true
		})
		inboundInfo.UserOnlineIP.Range(func(key, value interface{}) bool {
			ipMap := value.(*sync.Map)
			ipMap.Range(func(key, value interface{}) bool {
				ip := key.(string)
				uid := value.(int)
				onlineUser = append(onlineUser, api.OnlineUser{UID: uid, IP: ip})
				return true
			})
			email := key.(string)
			inboundInfo.UserOnlineIP.Delete(email) // Reset online device
			return true
		})
	} else {
		return nil, fmt.Errorf("no such inbound in limiter: %s", tag)
	}
	return &onlineUser, nil
}

func (l *Limiter) GetUserBucket(tag string, email string, ip string) (limiter *ratelimit.Bucket, SpeedLimit bool, Reject bool) {
	if value, ok := l.InboundInfo.Load(tag); ok {
		inboundInfo := value.(*InboundInfo)
		nodeLimit := inboundInfo.NodeSpeedLimit
		var userLimit uint64 = 0
		var deviceLimit = 0
		var uid = 0
		if v, ok := inboundInfo.UserInfo.Load(email); ok {
			u := v.(UserInfo)
			userLimit = u.SpeedLimit
			deviceLimit = u.DeviceLimit
			uid = u.UID
		}
		ipMap := new(sync.Map)
		if inboundInfo.EnableReorder {
			ipMap.Store(ip, true)
		} else {
			ipMap.Store(ip, uid)
		}
		// If any device is online
		if v, ok := inboundInfo.UserOnlineIP.LoadOrStore(email, ipMap); ok {
			ipMap := v.(*sync.Map)
			// If this ip is a new device
			if inboundInfo.EnableReorder {
				if _, ok := ipMap.LoadOrStore(ip, true); !ok {
					counter := 0
					ipMap.Range(func(key, value interface{}) bool {
						counter++
						return true
					})
					if counter > deviceLimit && deviceLimit > 0 {
						ipMap.Delete(ip)
						return nil, false, true
					}
				} else {
					ipMap.Store(ip, true)
				}
			} else {
				if _, ok := ipMap.LoadOrStore(ip, uid); !ok {
					counter := 0
					ipMap.Range(func(key, value interface{}) bool {
						counter++
						return true
					})
					if counter > deviceLimit && deviceLimit > 0 {
						ipMap.Delete(ip)
						return nil, false, true
					}
				} else {
					ipMap.Store(ip, uid)
				}
			}
		}
		limit := determineRate(nodeLimit, userLimit) // If need the Speed limit
		if limit > 0 {
			limiter := ratelimit.NewBucketWithQuantum(time.Second, int64(limit), int64(limit)) // Byte/s
			if v, ok := inboundInfo.BucketHub.LoadOrStore(email, limiter); ok {
				bucket := v.(*ratelimit.Bucket)
				return bucket, true, false
			} else {
				return limiter, true, false
			}
		} else {
			return nil, false, false
		}
	} else {
		newError("Get Inbound Limiter information failed").AtDebug().WriteToLog()
		return nil, false, false
	}
}

// determineRate returns the minimum non-zero rate
func determineRate(nodeLimit, userLimit uint64) (limit uint64) {
	if nodeLimit == 0 || userLimit == 0 {
		if nodeLimit > userLimit {
			return nodeLimit
		} else if nodeLimit < userLimit {
			return userLimit
		} else {
			return 0
		}
	} else {
		if nodeLimit > userLimit {
			return userLimit
		} else if nodeLimit < userLimit {
			return nodeLimit
		} else {
			return nodeLimit
		}
	}
}
