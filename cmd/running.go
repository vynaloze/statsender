package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/run"
)

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Runs the application in detached mode",
	Long:  `Runs the application in the background.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		run.Run(configDir)
		// todo detached mode
	},
}

var logDir string

func addRunning() {
	cmdRun.Flags().StringVarP(&logDir, "log", "l", "logs", "log directory location")
	rootCmd.AddCommand(cmdRun)
}
