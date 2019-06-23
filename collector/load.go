package collector

import (
	"github.com/shirou/gopsutil/load"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
	"runtime"
)

// Load collects statistics about CPU load (linux only)
type Load Config

// Collect collects statistics from given datasource
func (l *Load) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	stats, err := load.Avg()
	if err != nil {
		log.Error(err)
		if runtime.GOOS != "linux" {
			log.Error("This stat is supported only in linux")
		}
	}
	allStats := []load.AvgStat{*stats}

	return dto.NewStat(datasource.DsDto, "load", allStats)
}

// Conf return the configuration of the collector
func (l *Load) Conf() Config {
	return Config(*l)
}
