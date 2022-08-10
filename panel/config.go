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

type ApiConfig struct {
	ApiHost      string `mapstructure:"ApiHost"`
	ApiKey       string `mapstructure:"ApiKey"`
	NodeID       string `mapstructure:"NodeID"`
	Timeout      uint32 `mapstructure:"Timeout"`
	EnableVless  bool   `mapstructure:"EnableVless"`
	EnableXTLS   bool   `mapstructure:"EnableXTLS"`
	SpeedLimit   uint32 `mapstructure:"SpeedLimit"`
	DeviceLimit  uint32 `mapstructure:"DeviceLimit"`
	RuleListPath string `mapstructure:"RuleListPath"`
}
