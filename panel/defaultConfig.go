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

// getDefaultApiConfig returns the default api config.
func getDefaultApiConfig() *ApiConfig {
	return &ApiConfig{
		ApiHost:      "https://aikocute.com",
		ApiKey:       "",
		NodeID:       "",
		Timeout:      120,
		EnableVless:  false,
		EnableXTLS:   false,
		SpeedLimit:   0,
		DeviceLimit:  0,
		RuleListPath: "",
	}
}
