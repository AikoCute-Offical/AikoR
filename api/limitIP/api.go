package limitip

type IpRecorder interface {
	SyncOnlineIp(Ips []UserIpList) ([]UserIpList, error)
}
