package collector

import (
	"github.com/shirou/gopsutil/net"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type Net Config

func (n *Net) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	stats, err := net.IOCounters(false)

	if err != nil {
		log.Error(err)
	}

	return dto.NewStat(datasource.DsDto, "net", stats)
}

func (n *Net) Conf() Config {
	return Config(*n)
}
