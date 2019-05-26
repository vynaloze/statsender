package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/run"
)

var cmdTry = &cobra.Command{
	Use:   "try",
	Short: "Runs the application",
	Long: `Runs the application in the foreground.
Logs are printed directly to the console - useful for debugging purposes.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		run.Run()
	},
}

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Runs the application in detached mode",
	Long:  `Runs the application in the background.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("I just executed command 'run'...")
		fmt.Println("Log dir:", logDir)
	},
}

var logDir string

func addRunning() {
	rootCmd.AddCommand(cmdTry)
	cmdRun.Flags().StringVarP(&logDir, "log", "l", "/var/log/statsender", "Log directory location")
	rootCmd.AddCommand(cmdRun)
}
