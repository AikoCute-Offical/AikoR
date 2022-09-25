package panel

import (
	"github.com/AikoCute-Offical/AikoR/api"
	"github.com/AikoCute-Offical/AikoR/service/controller"
)

type Config struct {
	LogConfig          *LogConfig       `mapstructure:"Log"`
	DnsConfigPath      string           `mapstructure:"DnsConfigPath"`
	InboundConfigPath  string           `mapstructure:"InboundConfigPath"`
	OutboundConfigPath string           `mapstructure:"OutboundConfigPath"`
	RouteConfigPath    string           `mapstructure:"RouteConfigPath"`
	ConnetionConfig    *ConnetionConfig `mapstructure:"ConnetionConfig"`
	NodesConfig        []*NodesConfig   `mapstructure:"Nodes"`
}

type NodesConfig struct {
	PanelType        string             `mapstructure:"PanelType"`
	ApiConfig        *api.Config        `mapstructure:"ApiConfig"`
	ControllerConfig *controller.Config `mapstructure:"ControllerConfig"`
}

type LogConfig struct {
	Level      string `mapstructure:"Level"`
	AccessPath string `mapstructure:"AccessPath"`
	ErrorPath  string `mapstructure:"ErrorPath"`
}

type ConnetionConfig struct {
	Handshake    uint32 `mapstructure:"handshake"`
	ConnIdle     uint32 `mapstructure:"connIdle"`
	UplinkOnly   uint32 `mapstructure:"uplinkOnly"`
	DownlinkOnly uint32 `mapstructure:"downlinkOnly"`
	BufferSize   int32  `mapstructure:"bufferSize"`
}

type CertConfig struct {
	CertMode         string            `mapstructure:"CertMode"` // none, file, http, dns
	RejectUnknownSni bool              `mapstructure:"RejectUnknownSni"`
	CertDomain       string            `mapstructure:"CertDomain"`
	CertFile         string            `mapstructure:"CertFile"`
	KeyFile          string            `mapstructure:"KeyFile"`
	Provider         string            `mapstructure:"Provider"` // alidns, cloudflare, gandi, godaddy....
	Email            string            `mapstructure:"Email"`
	DNSEnv           map[string]string `mapstructure:"DNSEnv"`
}

type FallBackConfig struct {
	SNI              string `mapstructure:"SNI"`
	Alpn             string `mapstructure:"Alpn"`
	Path             string `mapstructure:"Path"`
	Dest             string `mapstructure:"Dest"`
	ProxyProtocolVer uint64 `mapstructure:"ProxyProtocolVer"`
}

type IpReportConfig struct {
	Url          string `mapstructure:"Url"`
	Token        string `mapstructure:"Token"`
	Periodic     int    `mapstructure:"Periodic"`
	Timeout      int    `mapstructure:"Timeout"`
	EnableIpSync bool   `mapstructure:"EnableIpSync"`
}

type ControllerConfig struct {
	ListenIP             string            `mapstructure:"ListenIP"`
	SendIP               string            `mapstructure:"SendIP"`
	UpdatePeriodic       int               `mapstructure:"UpdatePeriodic"`
	EnableDNS            bool              `mapstructure:"EnableDNS"`
	DNSType              string            `mapstructure:"DNSType"`
	DisableUploadTraffic bool              `mapstructure:"DisableUploadTraffic"`
	DisableGetRule       bool              `mapstructure:"DisableGetRule"`
	EnableProxyProtocol  bool              `mapstructure:"EnableProxyProtocol"`
	EnableFallback       bool              `mapstructure:"EnableFallback"`
	DisableIVCheck       bool              `mapstructure:"DisableIVCheck"`
	DisableSniffing      bool              `mapstructure:"DisableSniffing"`
	AikoRNginx           bool              `mapstructure:"AikoRNginx"`
	FallBackConfigs      []*FallBackConfig `mapstructure:"FallBackConfigs"`
	EnableIpRecorder     bool              `mapstructure:"EnableIpRecorder"`
	IpRecorderConfig     *IpReportConfig   `mapstructure:"IpRecorderConfig"`
	CertConfig           *CertConfig       `mapstructure:"CertConfig"`
}

type ApiConfig struct {
	APIHost             string  `mapstructure:"ApiHost"`
	NodeID              int     `mapstructure:"NodeID"`
	Key                 string  `mapstructure:"ApiKey"`
	NodeType            string  `mapstructure:"NodeType"`
	EnableVless         bool    `mapstructure:"EnableVless"`
	EnableXTLS          bool    `mapstructure:"EnableXTLS"`
	Timeout             int     `mapstructure:"Timeout"`
	SpeedLimit          float64 `mapstructure:"SpeedLimit"`
	DeviceLimit         int     `mapstructure:"DeviceLimit"`
	RuleListPath        string  `mapstructure:"RuleListPath"`
	DisableCustomConfig bool    `mapstructure:"DisableCustomConfig"`
}

type NodeConfig struct {
	ApiConfig        *ApiConfig        `mapstructure:"ApiConfig"`
	ControllerConfig *ControllerConfig `mapstructure:"ControllerConfig"`
}

type DNSEnv struct {
	CLOUDFLARE_EMAIL    string `mapstructure:"CLOUDFLARE_EMAIL"`
	CLOUDFLARE_API_KEY  string `mapstructure:"CLOUDFLARE_API_KEY"`
	ALICLOUD_ACCESS_KEY string `mapstructure:"ALICLOUD_ACCESS_KEY"`
	ALICLOUD_SECRET_KEY string `mapstructure:"ALICLOUD_SECRET_KEY"`
}
