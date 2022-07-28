package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"

	"github.com/AikoCute-Offical/AikoR/panel"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	configFile   = flag.String("config", "", "Config file for AikoR.")
	printVersion = flag.Bool("version", false, "show version")
)

var (
	version  = "0.0.1"
	codename = "AikoR"
	intro    = "A Xray backend that supports many panels"
)

func showVersion() {
	fmt.Printf("%s %s (%s) \n", codename, version, intro)
}

func getConfig() *viper.Viper {
	config := viper.New()

	// Set custom path and name
	if *configFile != "" {
		configName := path.Base(*configFile)
		configFileExt := path.Ext(*configFile)
		configNameOnly := strings.TrimSuffix(configName, configFileExt)
		configPath := path.Dir(*configFile)
		config.SetConfigName(configNameOnly)
		config.SetConfigType(strings.TrimPrefix(configFileExt, "."))
		config.AddConfigPath(configPath)
	} else {
		// Set default config path
		config.SetConfigName("aiko")
		config.SetConfigType("yml")
		config.AddConfigPath(".")

	}

	if err := config.ReadInConfig(); err != nil {
		log.Panicf("Fatal error config file: %s \n", err)
	}

	config.WatchConfig() // Watch the config

	return config
}

func main() {
	flag.Parse()
	showVersion()
	if *printVersion {
		return
	}

	config := getConfig()
	panelConfig := &panel.Config{}
	config.Unmarshal(panelConfig)
	p := panel.New(panelConfig)
	config.OnConfigChange(func(e fsnotify.Event) {
		// Hot reload function
		fmt.Println("Config file changed:", e.Name)
		p.Close()
		config.Unmarshal(panelConfig)
		p = panel.New(panelConfig)
		p.Start()
	})
	p.Start()
	defer p.Close()

	//Explicitly triggering GC to remove garbage from config loading.
	runtime.GC()
	// Running backend
	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	}
}
