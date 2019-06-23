// Package collector contains all statistic collectors
package collector

import (
	"github.com/vynaloze/pgstats"
	"github.com/vynaloze/statsender/dto"
)

// All collectors must satisfy collector.Collector interface
type Collector interface {
	Collect(datasource *Datasource) *dto.Stat
	Conf() Config
}

// Config represents base configuration of each collector
type Config struct {
	Cron    string `hcl:"cron"`
	Enabled bool   `hcl:"enabled"`
}

// Datasource contains information from where the collector should gather statistics
type Datasource struct {
	DsDto            *dto.Datasource
	PgStats          *pgstats.PgStats
	ConnectionString *string
}
