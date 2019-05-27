package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{Use: "statsender"}

var verbose bool
var configDir string

func Execute() {
	rootCmd.Version = "0.0.1"
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&configDir, "config", "c", "/etc/statsender/conf.d", "sets configuration directory location")

	addRunning()
	addDs()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
