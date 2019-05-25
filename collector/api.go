package collector

import (
	"github.com/vynaloze/pgstats"
	"github.com/vynaloze/statsender/dto"
)

type SystemCollector interface {
	Collect() *dto.Stat
	Conf() Config
}

type PostgresCollector interface {
	Collect(datasource Datasource) *dto.Stat
	Conf() Config
}

type Config struct {
	Cron string `hcl:"cron"`
}

type Datasource struct {
	// all conn params here?
	DsDto dto.Datasource
	Conn  pgstats.PgStats
}
