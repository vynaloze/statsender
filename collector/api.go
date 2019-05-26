package collector

import (
	"github.com/vynaloze/pgstats"
	"github.com/vynaloze/statsender/dto"
)

type Collector interface {
	Collect(datasource *Datasource) *dto.Stat
	Conf() Config
}

type Config struct {
	Cron string `hcl:"cron"`
}

type Datasource struct {
	DsDto *dto.Datasource
	Conn  *pgstats.PgStats
}
