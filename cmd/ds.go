package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/vynaloze/statsender/config"
	"github.com/vynaloze/statsender/logger"
	"strings"
)

var cmdDs = &cobra.Command{
	Use:   "datasource",
	Short: "Manage datasources",
	Long: `Manage datasources. 
Type 'statsender datasource' to see more details`,
}

var cmdDsAdd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new datasource",
	Long: `Adds a new datasource. 
If not stated otherwise (with flag --file or --filename), it will be saved in <config_dir>/_ds.hcl`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid number of arguments")
		}
		if _, err := config.ParseDSN(args[0]); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := logger.New()
		ds, _ := config.ParseDSN(args[0])
		ds.Tags = parseTags(tags)
		log.Debugf("Parsed datasource: %v", *ds)
		if err := config.AddDatasource(configDir, fileName, *ds); err != nil {
			log.Fatal(err)
		}
	},
}

var tags []string
var fileName string

func addDs() {
	cmdDsAdd.Flags().StringSliceVar(&tags, "tag", []string{}, "optional tag of a datasource, in format key=value")
	cmdDsAdd.Flags().StringVarP(&fileName, "file", "f", "_ds.hcl", "the name of the configuration file to update")
	cmdDs.AddCommand(cmdDsAdd)
	rootCmd.AddCommand(cmdDs)
}

func parseTags(tags []string) map[string]string {
	log, _ := logger.New()
	m := make(map[string]string)
	for _, t := range tags {
		s := strings.Split(t, "=")
		if len(s) != 2 {
			log.Fatal("Invalid format of a tag") //todo nicer message, without scary stacktrace
		}
		m[s[0]] = s[1]
	}
	return m
}
