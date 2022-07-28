package panel

import (
	"log"
	"sync"

	"github.com/AikoCute-Offical/AikoR/api"
	"github.com/AikoCute-Offical/AikoR/api/sspanel"
	_ "github.com/AikoCute-Offical/AikoR/main/distro/all"
	"github.com/AikoCute-Offical/AikoR/service"
	"github.com/AikoCute-Offical/AikoR/service/controller"
	"github.com/xtls/xray-core/app/dispatcher"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/stats"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
)

// Panel Structure
type Panel struct {
	access      sync.Mutex
	panelConfig *Config
	Server      *core.Instance
	Service     []service.Service
	Running     bool
}

func New(panelConfig *Config) *Panel {
	p := &Panel{panelConfig: panelConfig}
	return p
}

func (p *Panel) loadCore(c *LogConfig) *core.Instance {
	// Log Config
	logConfig := &conf.LogConfig{
		LogLevel:  c.Level,
		AccessLog: c.AccessPath,
		ErrorLog:  c.ErrorPath,
	}
	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&stats.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
			serial.ToTypedMessage(logConfig.Build()),
		},
	}
	server, err := core.New(config)
	if err != nil {
		log.Panicf("failed to create instance: %s", err)
	}
	log.Printf("Xray Core Version: %s", core.Version())

	return server
}

// Start Start the panel
func (p *Panel) Start() {
	p.access.Lock()
	defer p.access.Unlock()
	log.Print("Start the panel..")
	// Load Core
	server := p.loadCore(p.panelConfig.LogConfig)
	if err := server.Start(); err != nil {
		log.Panicf("Failed to start instance: %s", err)
	}
	p.Server = server
	// Load Nodes config
	for _, nodeConfig := range p.panelConfig.NodesConfig {
		var apiClient api.API
		if nodeConfig.PanelType == "SSpanel" {
			apiClient = sspanel.New(nodeConfig.ApiConfig)
		} else {
			log.Panicf("Unsupport panel type: %s", nodeConfig.PanelType)
		}
		var controllerService service.Service
		// Regist controller service
		controllerService = controller.New(server, apiClient, nodeConfig.ControllerConfig)
		p.Service = append(p.Service, controllerService)

	}

	// Start all the service
	for _, s := range p.Service {
		err := s.Start()
		if err != nil {
			log.Panicf("Panel Start fialed: %s", err)
		}
	}
	p.Running = true
	return
}

// Close Close the panel
func (p *Panel) Close() {
	p.access.Lock()
	defer p.access.Unlock()
	for _, s := range p.Service {
		err := s.Close()
		if err != nil {
			log.Panicf("Panel Close fialed: %s", err)
		}
	}
	p.Server.Close()
	p.Running = false
	return
}
