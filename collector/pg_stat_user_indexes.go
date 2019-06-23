package collector

import (
	"github.com/vynaloze/pgstats/nullable"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type PgStatUserIndexes Config

type pgStatUserIndexesPayload struct {
	Table       string         `json:"table"`
	Index       string         `json:"index"`
	IdxScan     nullable.Int64 `json:"idx_scan"`
	IdxTupRead  nullable.Int64 `json:"idx_tup_read"`
	IdxTupFetch nullable.Int64 `json:"idx_tup_fetch"`
}

func (p *PgStatUserIndexes) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	s, err := datasource.PgStats.PgStatUserIndexes()
	if err != nil {
		log.Error(err)
	}
	var payload []pgStatUserIndexesPayload
	for _, r := range s {
		payload = append(payload, pgStatUserIndexesPayload{
			r.Schemaname + "." + r.Relname,
			r.Indexrelname,
			r.IdxScan,
			r.IdxTupRead,
			r.IdxTupFetch,
		})
	}

	return dto.NewStat(datasource.DsDto, "pg_stat_user_indexes", payload)
}

func (p *PgStatUserIndexes) Conf() Config {
	return Config(*p)
}
