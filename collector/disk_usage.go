package collector

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type DiskUsage Config

func (d *DiskUsage) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Error(err)
	}

	var allStats []disk.UsageStat
	for _, p := range partitions {
		stats, err := disk.Usage(p.Mountpoint)
		if err != nil {
			log.Error(err)
		}
		allStats = append(allStats, *stats)
	}

	return dto.NewStat(datasource.DsDto, "disk_usage", allStats)
}

func (d *DiskUsage) Conf() Config {
	return Config(*d)
}
