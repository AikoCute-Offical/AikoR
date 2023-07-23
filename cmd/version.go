package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	version  = "v1.4.0" //use ldflags replace
	codename = "AikoR"
	intro    = "AikoR support for Aiko"
)

var versionCommand = cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

func init() {
	command.AddCommand(&versionCommand)
}

func showVersion() {
	fmt.Printf("%s %s (%s) \n", codename, version, intro)
}