package v2board

import "encoding/json"

type UserTraffic struct {
	UID      int   `json:"user_id"`
	Upload   int64 `json:"u"`
	Download int64 `json:"d"`
}

type OnlineUser struct {
	UID int
	IP  string
}

type V2RayUserInfo struct {
	Uuid    string `json:"uuid"`
	Email   string `json:"email"`
	AlterId int    `json:"alter_id"`
}
type TrojanUserInfo struct {
	Password string `json:"password"`
}

type PostData struct {
	Data interface{} `json:"data"`
}

type ReportIllegal struct {
	ID  int `json:"list_id"`
	UID int `json:"user_id"`
}

type IllegalItem struct {
	ID  int `json:"list_id"`
	UID int `json:"user_id"`
}

type Response struct {
	Ret  uint            `json:"ret"`
	Data json.RawMessage `json:"data"`
}

type SystemLoad struct {
	Uptime string `json:"uptime"`
	Load   string `json:"load"`
}
