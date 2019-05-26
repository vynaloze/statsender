package collector

import (
	"github.com/shirou/gopsutil/disk"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type DiskIo Config

func (d *DiskIo) Collect() *dto.Stat {
	log, _ := logger.New()

	stats, err := disk.IOCounters()

	if err != nil {
		log.Error(err)
	}

	return dto.NewStat(dto.NewDatasource(), "disk_io", stats)
}

func (d *DiskIo) Conf() Config {
	return Config(*d)
}
