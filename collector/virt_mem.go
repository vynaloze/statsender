package collector

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type VirtMem Config

func (m *VirtMem) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	stats, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
	}
	allStats := []mem.VirtualMemoryStat{*stats}

	return dto.NewStat(datasource.DsDto, "virt_mem", allStats)
}

func (m *VirtMem) Conf() Config {
	return Config(*m)
}
