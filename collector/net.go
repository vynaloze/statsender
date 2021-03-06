package collector

import (
	"github.com/shirou/gopsutil/net"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// Net collects statistics about network I/O
type Net Config

type netPayload struct {
	BytesIn  uint64 `json:"bytes_in"`
	BytesOut uint64 `json:"bytes_out"`
}

// Collect collects statistics from given datasource
func (n *Net) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	stats, err := net.IOCounters(false)
	if err != nil {
		log.Error(err)
	}
	payload := []netPayload{{stats[0].BytesRecv, stats[0].BytesSent}}

	return dto.NewStat(datasource.DsDto, "net", payload)
}

// Conf return the configuration of the collector
func (n *Net) Conf() Config {
	return Config(*n)
}
