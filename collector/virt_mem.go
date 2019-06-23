package collector

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// VirtMem collects statistics about virtual memory usage
type VirtMem Config

type memPayload struct {
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
}

// Collect collects statistics from given datasource
func (m *VirtMem) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	stats, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
	}
	payload := []memPayload{{stats.Total, stats.Available}}

	return dto.NewStat(datasource.DsDto, "virt_mem", payload)
}

// Conf return the configuration of the collector
func (m *VirtMem) Conf() Config {
	return Config(*m)
}
