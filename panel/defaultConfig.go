package panel

import "github.com/AikoCute-Offical/AikoR/service/controller"

func getDefaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:      "none",
		AccessPath: "",
		ErrorPath:  "",
	}
}

func getDefaultConnetionConfig() *ConnetionConfig {
	return &ConnetionConfig{
		Handshake:    4,
		ConnIdle:     30,
		UplinkOnly:   2,
		DownlinkOnly: 4,
		BufferSize:   64,
	}
}

func getDefaultControllerConfig() *controller.Config {
	return &controller.Config{
		ListenIP:       "0.0.0.0",
		SendIP:         "0.0.0.0",
		UpdatePeriodic: 60,
		DNSType:        "AsIs",
	}
}

// Connecting With redis sever
func getConnectRedisConfig() *ConnectRedisConfig {
	return &ConnectRedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
}
