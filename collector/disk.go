package collector

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type Disk Config

type diskPayload struct {
	BytesAvailable uint64 `json:"bytes_available"`
	BytesTotal     uint64 `json:"bytes_total"`
	Reads          uint64 `json:"reads"`
	Writes         uint64 `json:"writes"`
	BytesRead      uint64 `json:"bytes_read"`
	BytesWrite     uint64 `json:"bytes_write"`
}

func (d *Disk) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	payload := diskPayload{}

	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Error(err)
	}
	for _, p := range partitions {
		usageStats, err := disk.Usage(p.Mountpoint)
		if err != nil {
			log.Error(err)
		}
		payload.BytesAvailable += usageStats.Free
		payload.BytesTotal += usageStats.Total
	}

	ioStats, err := disk.IOCounters()
	if err != nil {
		log.Error(err)
	}
	for _, s := range ioStats {
		payload.Reads += s.ReadCount
		payload.Writes += s.WriteCount
		payload.BytesRead += s.ReadBytes
		payload.BytesWrite += s.WriteBytes
	}

	return dto.NewStat(datasource.DsDto, "disk", []diskPayload{payload})
}

func (d *Disk) Conf() Config {
	return Config(*d)
}
