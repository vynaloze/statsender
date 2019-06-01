package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/config"
	"github.com/vynaloze/statsender/logger"
)

var cmdSender = &cobra.Command{
	Use:   "sender",
	Short: "Manage senders",
	Long: `Manage senders. 
Type 'statsender sender --help' to see more details`,
}

var cmdSenderAdd = &cobra.Command{
	Use:   "add <type> [<spec>]",
	Short: "Adds a new sender",
	Long: `Adds a new sender.
Valid <type>s: 'console', 'http'
In case of 'http', <spec> looks like: 'http[s]://host[:port][/endpoint]'
If not stated otherwise (with flag --file or --filename), sender will be saved in ${config_dir}/_senders.hcl`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return errors.New("invalid number of arguments")
		}
		if _, err := config.ParseSender(args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := logger.New()
		log.Debug("parsing argument")
		s, _ := config.ParseSender(args) // unhandled error because we verify it in Args
		log.Debugf("parsed sender: %+v", *s)
		if err := config.AddSender(configDir, fileNameSender, *s); err != nil {
			log.Fatal(err)
		}
	},
}

var fileNameSender string

func addSender() {
	cmdSenderAdd.Flags().StringVarP(&fileNameSender, "file", "f", "_senders.hcl", "the name of the configuration file to update")
	cmdSender.AddCommand(cmdSenderAdd)
	rootCmd.AddCommand(cmdSender)
}
