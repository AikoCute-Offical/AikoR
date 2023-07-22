package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var (
	version  = "1.4.0 Beta 1"
	codename = "AikoR"
	intro    = "Backend AikoR For Aiko"
)

func showVersion() {
	fmt.Printf("%s %s (%s) \n", codename, version, intro)
}

var command = &cobra.Command{
	Use: "AikoR",
	PreRun: func(_ *cobra.Command, _ []string) {
		showVersion()
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		log.Println("execute failed, error:", err)
	}
}