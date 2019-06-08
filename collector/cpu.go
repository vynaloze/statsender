package collector

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type Cpu Config

type cpuPayload struct {
	UsagePercent float64 `json:"usage_percent"`
}

func (c *Cpu) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	value, err := cpu.Percent(0, false)
	if err != nil {
		log.Error(err)
	}
	p := []cpuPayload{{UsagePercent: value[0]}}

	return dto.NewStat(datasource.DsDto, "cpu", p)
}

func (c *Cpu) Conf() Config {
	return Config(*c)
}
