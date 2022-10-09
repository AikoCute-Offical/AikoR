package aiko

type UserTraffic struct {
	UID      int   `json:"user_id"`
	Upload   int64 `json:"u"`
	Download int64 `json:"d"`
	Count    int64 `json:"n"`
}

type OnlineUser struct {
	UID int    `json:"user_id"`
	IP  string `json:"ip"`
}
