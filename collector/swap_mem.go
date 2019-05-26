package collector

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type SwapMem Config

func (m *SwapMem) Collect() *dto.Stat {
	log, _ := logger.New()

	stats, err := mem.SwapMemory()
	if err != nil {
		log.Error(err)
	}

	return dto.NewStat(dto.NewDatasource(), "swap_mem", stats)
}

func (m *SwapMem) Conf() Config {
	return Config(*m)
}