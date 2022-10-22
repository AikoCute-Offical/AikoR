package limitip

import "github.com/AikoCute-Offical/AikoR/app/mydispatcher"

type IpRecorder interface {
	SyncOnlineIp(Ips []mydispatcher.UserIpList) ([]mydispatcher.UserIpList, error)
}
