// Package limiter is to control the links that go into the dispather
package limiter

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"

	"github.com/AikoCute-Offical/AikoR/api"
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
	BucketHub      *sync.Map // key: Email, value: *rate.Limiter
	UserOnlineIP   *sync.Map // Key: Email Value: *sync.Map: Key: IP, Value: UID
}

type Limiter struct {
	InboundInfo *sync.Map // Key: Tag, Value: *InboundInfo
	r           *redis.Client
	g           struct {
		redislimit int
		timeout    int
		expiry     int
	}
}

func New() *Limiter {
	return &Limiter{
		InboundInfo: new(sync.Map),
	}
}

func (l *Limiter) AddInboundLimiter(tag string, nodeSpeedLimit uint64, userList *[]api.UserInfo, Redis *RedisConfig) error {
	// global limit
	if Redis.RedisLimit > 0 {
		l.r = redis.NewClient(&redis.Options{
			Addr:     Redis.RedisAddr,
			Password: Redis.RedisPassword,
			DB:       Redis.RedisDB,
		})
		l.g.redislimit = Redis.RedisLimit
		l.g.timeout = Redis.Timeout
		l.g.expiry = Redis.Expiry
	}

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
			// Update old limiter bucket
			limit := determineRate(inboundInfo.NodeSpeedLimit, u.SpeedLimit)
			if limit > 0 {
				if bucket, ok := inboundInfo.BucketHub.Load(fmt.Sprintf("%s|%s|%d", tag, u.Email, u.UID)); ok {
					limiter := bucket.(*rate.Limiter)
					limiter.SetLimit(rate.Limit(limit))
					limiter.SetBurst(int(limit))
				}
			} else {
				inboundInfo.BucketHub.Delete(fmt.Sprintf("%s|%s|%d", tag, u.Email, u.UID))
			}
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

func (l *Limiter) GetUserBucket(tag string, email string, ip string) (limiter *rate.Limiter, SpeedLimit bool, Reject bool) {
	if value, ok := l.InboundInfo.Load(tag); ok {
		inboundInfo := value.(*InboundInfo)
		nodeLimit := inboundInfo.NodeSpeedLimit
		var (
			userLimit        uint64 = 0
			deviceLimit, uid int
		)

		if v, ok := inboundInfo.UserInfo.Load(email); ok {
			u := v.(UserInfo)
			uid = u.UID
			userLimit = u.SpeedLimit
			deviceLimit = u.DeviceLimit
		}

		// Global device limit
		if l.g.redislimit > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(l.g.timeout))
			defer cancel()
			uidString := strconv.Itoa(uid)
			// If any device is online
			if exists, err := l.r.Exists(ctx, uidString).Result(); err != nil {
				newError(fmt.Sprintf("Redis: %v", err)).AtError().WriteToLog()
			} else if exists == 0 { // No user is online
				l.r.SAdd(ctx, uidString, ip)
				l.r.Expire(ctx, uidString, time.Second*time.Duration(l.g.expiry))
			} else {
				// If this ip is a new device
				if online, err := l.r.SIsMember(ctx, uidString, ip).Result(); err != nil {
					newError(fmt.Sprintf("Redis: %v", err)).AtError().WriteToLog()
				} else if !online {
					l.r.SAdd(ctx, uidString, ip)
					if l.r.SCard(ctx, uidString).Val() > int64(l.g.redislimit) {
						l.r.SRem(ctx, uidString, ip)
						return nil, false, true
					}
				}
			}
		}

		// Local device limit
		ipMap := new(sync.Map)
		ipMap.Store(ip, uid)
		// If any device is online
		if v, ok := inboundInfo.UserOnlineIP.LoadOrStore(email, ipMap); ok {
			ipMap := v.(*sync.Map)
			// If this ip is a new device
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
			}
		}
		limit := determineRate(nodeLimit, userLimit) // If need the Speed limit
		if limit > 0 {
			limiter := rate.NewLimiter(rate.Limit(limit), int(limit)) // Byte/s
			if v, ok := inboundInfo.BucketHub.LoadOrStore(email, limiter); ok {
				bucket := v.(*rate.Limiter)
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
