package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/logger"
	"github.com/vynaloze/statsender/run"
	"os"
	"os/exec"
	"path/filepath"
)

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Runs the application in detached mode",
	Long:  `Runs the application in the background.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := logger.New()
		if detached {
			err := os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			logger.OutputToFile(filepath.Join(logDir, "statsender.log"))
			run.Run(configDir)
		} else {
			cmd := exec.Command(os.Args[0], append(os.Args[1:], "--detached")...)
			err := cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

var logDir string
var detached bool

func addRunning() {
	cmdRun.Flags().StringVarP(&logDir, "log", "l", "logs", "log directory location")
	cmdRun.Flags().BoolVar(&detached, "detached", false, "starts the application in detached mode")
	_ = cmdRun.Flags().MarkHidden("detached")
	rootCmd.AddCommand(cmdRun)
}
