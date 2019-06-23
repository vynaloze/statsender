package collector

import (
	"github.com/shirou/gopsutil/mem"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// SwapMem collects statistics about swap memory usage
type SwapMem Config

// Collect collects statistics from given datasource
func (m *SwapMem) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	stats, err := mem.SwapMemory()
	if err != nil {
		log.Error(err)
	}
	payload := []memPayload{{stats.Total, stats.Free}}

	return dto.NewStat(datasource.DsDto, "swap_mem", payload)
}

// Conf return the configuration of the collector
func (m *SwapMem) Conf() Config {
	return Config(*m)
}
