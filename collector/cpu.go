package collector

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type Cpu Config

func (c *Cpu) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	value, err := cpu.Percent(0, false)
	if err != nil {
		log.Error(err)
	}

	return dto.NewStat(datasource.DsDto, "cpu", value[0])
}

func (c *Cpu) Conf() Config {
	return Config(*c)
}
