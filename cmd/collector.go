package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/config"
	"github.com/vynaloze/statsender/logger"
)

var cmdCollector = &cobra.Command{
	Use:   "collector",
	Short: "Manage collectors",
	Long: `Manage collectors. 
Type 'statsender collector --help' to see more details`,
}

var cmdCollectorEnable = &cobra.Command{
	Use:   "enable <type>",
	Short: "Enables a collector",
	Long: `Enables a collector. If such collector does not exist, error is returned.
If not stated otherwise (with flag --file or --filename), sender will be saved in ${config_dir}/_collectors.hcl`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := logger.New()
		if err := config.SetCollectorEnabled(configDir, fileNameCollectors, args[0], true); err != nil {
			log.Fatal(err)
		}
	},
}

var cmdCollectorDisable = &cobra.Command{
	Use:   "disable <type>",
	Short: "Disables a collector",
	Long: `Disables a collector. If such collector does not exist, error is returned.
If not stated otherwise (with flag --file or --filename), sender will be saved in ${config_dir}/_collectors.hcl`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := logger.New()
		if err := config.SetCollectorEnabled(configDir, fileNameCollectors, args[0], false); err != nil {
			log.Fatal(err)
		}
	},
}

var cmdCollectorSetCron = &cobra.Command{
	Use:   "schedule <type> <cron_expr>",
	Short: "Changes the cron schedule of a collector",
	Long: `Changes the cron schedule of a collector. If such collector does not exist, error is returned.
If not stated otherwise (with flag --file or --filename), sender will be saved in ${config_dir}/_collectors.hcl`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid number of arguments")
		}
		if !config.IsCronValid(args[1]) {
			return errors.New("invalid cron expression")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := logger.New()
		if err := config.SetCollectorCron(configDir, fileNameCollectors, args[0], args[1]); err != nil {
			log.Fatal(err)
		}
	},
}

var fileNameCollectors string

func addCollector() {
	cmdCollectorEnable.Flags().StringVarP(&fileNameCollectors, "file", "f", "_collectors.hcl", "the name of the configuration file to update")
	cmdCollector.AddCommand(cmdCollectorEnable)
	cmdCollectorDisable.Flags().StringVarP(&fileNameCollectors, "file", "f", "_collectors.hcl", "the name of the configuration file to update")
	cmdCollector.AddCommand(cmdCollectorDisable)
	cmdCollectorSetCron.Flags().StringVarP(&fileNameCollectors, "file", "f", "_collectors.hcl", "the name of the configuration file to update")
	cmdCollector.AddCommand(cmdCollectorSetCron)
	rootCmd.AddCommand(cmdCollector)
}
