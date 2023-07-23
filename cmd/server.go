package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/AikoCute-Offical/AikoR/panel"
)

var (
	cfgFile string

	serverCommand = &cobra.Command{
		Use:   "server",
		Short: "Backend AikoR For Aiko",
		Run:   serverHandle,
		Args:  cobra.NoArgs,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	serverCommand.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./aiko.yml)")
	command.AddCommand(serverCommand)
}

func serverHandle(cmd *cobra.Command, args []string) {
	showVersion()
	config := getConfig()
	panelConfig := &panel.Config{}
	if err := config.Unmarshal(panelConfig); err != nil {
		log.Panicf("Parse config file %v failed: %s \n", cfgFile, err)
	}
	p := panel.New(panelConfig)
	lastTime := time.Now()
	config.OnConfigChange(func(e fsnotify.Event) {
		if time.Now().After(lastTime.Add(3 * time.Second)) {
			fmt.Println("Config file changed:", e.Name)
			p.Close()
			runtime.GC()
			if err := config.Unmarshal(panelConfig); err != nil {
				log.Panicf("Parse config file %v failed: %s \n", cfgFile, err)
			}
			p.Start()
			lastTime = time.Now()
		}
	})
	p.Start()
	defer p.Close()
	runtime.GC()
}

func getConfig() *viper.Viper {
	config := viper.New()

	if cfgFile != "" {
		configName := path.Base(cfgFile)
		configFileExt := path.Ext(cfgFile)
		configNameOnly := strings.TrimSuffix(configName, configFileExt)
		configPath := path.Dir(cfgFile)
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

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("aiko")
		viper.SetConfigType("yml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Cannot use config file:", viper.ConfigFileUsed())
		os.Exit(1)
	}
}