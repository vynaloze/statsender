package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/config"
	"github.com/vynaloze/statsender/logger"
)

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initializes the application config",
	Long:  `Generates default configuration files and examples of usage.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := logger.New()
		err := config.InitConfig(configDir)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var cmdTry = &cobra.Command{
	Use:   "try",
	Short: "Tests the application",
	Long: `Checks if the configuration can be parsed 
and if it is possible to establish connections to databases and http endpoints.
Displays stats that will be collected and their intervals`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		//run.Run(configDir)
	},
}

func addGeneral() {
	rootCmd.AddCommand(cmdInit)
	rootCmd.AddCommand(cmdTry)
}
