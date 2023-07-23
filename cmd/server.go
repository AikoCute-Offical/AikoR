package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/AikoCute-Offical/AikoR/panel"
)

var (
	configFile string
	version    = "1.3.4 Beta 2"
	codename   = "AikoR"
	intro      = "Backend AikoR For Aiko"
)

var serverCommand = cobra.Command{
	Use:     "server",
	Short:   "Backend AikoR For Aiko",
	Long:    `AikoR is a backend service for Aiko.`,
	Version: fmt.Sprintf("%s %s (%s)", codename, version, intro),
	Run: func(cmd *cobra.Command, args []string) {
		config := getConfig()
		panelConfig := &panel.Config{}
		if err := config.Unmarshal(panelConfig); err != nil {
			log.Panicf("Parse config file %v failed: %s \n", configFile, err)
		}
		p := panel.New(panelConfig)
		lastTime := time.Now()
		config.OnConfigChange(func(e fsnotify.Event) {
			if time.Now().After(lastTime.Add(3 * time.Second)) {
				fmt.Println("Config file changed:", e.Name)
				p.Close()
				runtime.GC()
				if err := config.Unmarshal(panelConfig); err != nil {
					log.Panicf("Parse config file %v failed: %s \n", configFile, err)
				}
				p.Start()
				lastTime = time.Now()
			}
		})
		p.Start()
		defer p.Close()

		runtime.GC()

		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	serverCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file for AikoR")
	serverCommand.PersistentFlags().Bool("watch", true, "watch file path change")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("aiko")
		viper.SetConfigType("yml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Config file error: %s \n", err)
	}

	viper.WatchConfig()
}

func getConfig() *viper.Viper {
	config := viper.New()

	if configFile != "" {
		configName := path.Base(configFile)
		configFileExt := path.Ext(configFile)
		configNameOnly := strings.TrimSuffix(configName, configFileExt)
		configPath := path.Dir(configFile)
		config.SetConfigName(configNameOnly)
		config.SetConfigType(strings.TrimPrefix(configFileExt, "."))
		config.AddConfigPath(configPath)
		os.Setenv("XRAY_LOCATION_ASSET", configPath)
		os.Setenv("XRAY_LOCATION_CONFIG", configPath)
	} else {
		config.SetConfigName("aiko")
		config.SetConfigType("yml")
		config.AddConfigPath(".")
	}

	if err := config.ReadInConfig(); err != nil {
		log.Panicf("Config file error: %s \n", err)
	}

	config.WatchConfig()

	return config
}