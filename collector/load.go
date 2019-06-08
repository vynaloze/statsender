package collector

import (
	"github.com/shirou/gopsutil/load"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
	"runtime"
)

type Load Config

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

func (l *Load) Conf() Config {
	return Config(*l)
}
