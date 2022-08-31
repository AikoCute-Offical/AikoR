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
		ListenIP:        "0.0.0.0",
		SendIP:          "0.0.0.0",
		UpdatePeriodic:  60,
		DNSType:         "AsIs",
		DisableSniffing: true,
	}
}

func getDefaultNodesConfig() *NodesConfig {
	return &NodesConfig{
		PanelType:        "",
		ApiConfig:        nil,
		ControllerConfig: nil,
	}
}

func getDefaultConfig() *Config {
	return &Config{
		LogConfig:          getDefaultLogConfig(),
		DnsConfigPath:      "",
		InboundConfigPath:  "",
		OutboundConfigPath: "",
		RouteConfigPath:    "",
		ConnetionConfig:    getDefaultConnetionConfig(),
		NodesConfig:        []*NodesConfig{},
	}
}

func getDefaultCertConfig() *CertConfig {
	return &CertConfig{
		CertMode:         "none",
		RejectUnknownSni: false,
		CertDomain:       "",
		CertFile:         "",
		KeyFile:          "",
		Provider:         "",
		Email:            "",
		DNSEnv:           map[string]string{},
	}
}

func getDefaultFallBackConfig() *FallBackConfig {
	return &FallBackConfig{
		SNI:              "",
		Alpn:             "",
		Path:             "",
		Dest:             "",
		ProxyProtocolVer: 0,
	}
}

func getDefaultDNSEnv() *DNSEnv {
	return &DNSEnv{
		CLOUDFLARE_EMAIL:   "",
		CLOUDFLARE_API_KEY: "",
	}
}
