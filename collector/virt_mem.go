package collector

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type VirtMem Config

func (m *VirtMem) Collect() *dto.Stat {
	log, _ := logger.New()

	stats, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
	}

	return dto.NewStat(dto.NewDatasource(), "virt_mem", stats)
}

func (m *VirtMem) Conf() Config {
	return Config(*m)
}
