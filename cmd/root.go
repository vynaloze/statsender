package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/logger"
	"os"
)

var rootCmd = &cobra.Command{Use: "statsender"}

var verbose bool
var configDir string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Version = "0.0.1"
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&configDir, "config", "c", "conf", "sets configuration directory location")
	addRunning()
	addDs()
	addSender()
	addCollector()
	addGeneral()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	logger.SetDebug(verbose)
}
