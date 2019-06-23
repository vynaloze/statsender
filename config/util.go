package config

import "github.com/robfig/cron"

// IsCronValid validates the cron expression
func IsCronValid(spec string) bool {
	specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := specParser.Parse(spec)
	return err == nil
}
