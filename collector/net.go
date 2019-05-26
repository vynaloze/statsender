package collector

import (
	"github.com/shirou/gopsutil/net"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type Net Config

func (n *Net) Collect() *dto.Stat {
	log, _ := logger.New()

	statsSlice, err := net.IOCounters(false)

	if err != nil {
		log.Error(err)
	}

	stats := statsSlice[0]

	if stats.Name != "all" {
		log.Error("Not aggregated network statistics. This should never happened")
	}

	return dto.NewStat(dto.NewDatasource(), "net", stats)
}

func (n *Net) Conf() Config {
	return Config(*n)
}
