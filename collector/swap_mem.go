package collector

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type SwapMem Config

func (m *SwapMem) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	stats, err := mem.SwapMemory()
	if err != nil {
		log.Error(err)
	}
	allStats := []mem.SwapMemoryStat{*stats}

	return dto.NewStat(datasource.DsDto, "swap_mem", allStats)
}

func (m *SwapMem) Conf() Config {
	return Config(*m)
}
