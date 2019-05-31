package config

import "github.com/robfig/cron"

func IsCronValid(spec string) bool {
	specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := specParser.Parse(spec)
	return err == nil
}
