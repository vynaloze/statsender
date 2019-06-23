package collector

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// Cpu collects statistics about CPU usage
type Cpu Config

type cpuPayload struct {
	UsagePercent float64 `json:"usage_percent"`
}

// Collect collects statistics from given datasource
func (c *Cpu) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	value, err := cpu.Percent(0, false)
	if err != nil {
		log.Error(err)
	}
	p := []cpuPayload{{UsagePercent: value[0]}}

	return dto.NewStat(datasource.DsDto, "cpu", p)
}

// Conf return the configuration of the collector
func (c *Cpu) Conf() Config {
	return Config(*c)
}
