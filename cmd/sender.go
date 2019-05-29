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
Type 'statsender sender' to see more details`,
}

var cmdSenderAdd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new senders",
	Long: `Adds a new senders. 
If not stated otherwise (with flag --file or --filename), it will be saved in <config_dir>/_senders.hcl`,
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
		s, _ := config.ParseSender(args) // unhandled error because we verify it in Args
		log.Debugf("Parsed sender: %v", *s)
		if err := config.AddSender(configDir, fileNameSender, *s); err != nil {
			log.Fatal(err)
		}
	},
}

var fileNameSender string

func addSender() {
	cmdSenderAdd.Flags().StringVarP(&fileNameSender, "file", "f", "_sender.hcl", "the name of the configuration file to update")
	cmdSender.AddCommand(cmdSenderAdd)
	rootCmd.AddCommand(cmdSender)
}
