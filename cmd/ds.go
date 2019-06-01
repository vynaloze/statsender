package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/config"
	"github.com/vynaloze/statsender/logger"
)

var cmdDs = &cobra.Command{
	Use:   "datasource",
	Short: "Manage datasources",
	Long: `Manage datasources. 
Type 'statsender datasource --help' to see more details`,
}

var cmdDsAdd = &cobra.Command{
	Use:   "add <DSN>",
	Short: "Adds a new datasource",
	Long: `Adds a new datasource. 
Valid <DSN> format: '[postgresql://]login:password@host[:port]/dbname[?param1=value1&...]'
Optional tags are provided as flags: --tag key1=value1 --tag key2=value2 ...
If not stated otherwise (with flag --file or --filename), datasource will be saved in ${config_dir}/_datasources.hcl`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid number of arguments")
		}
		if _, err := config.ParseDSN(args[0], tags); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := logger.New()
		ds, _ := config.ParseDSN(args[0], tags) // unhandled error because we verify it in Args
		log.Debugf("Parsed datasource: %v", *ds)
		if err := config.AddDatasource(configDir, fileNameDs, *ds); err != nil {
			log.Fatal(err)
		}
	},
}

var tags []string
var fileNameDs string

func addDs() {
	cmdDsAdd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "optional tag of a datasource, in format key=value")
	cmdDsAdd.Flags().StringVarP(&fileNameDs, "file", "f", "_datasources.hcl", "the name of the configuration file to update")
	cmdDs.AddCommand(cmdDsAdd)
	rootCmd.AddCommand(cmdDs)
}
