package collector

import (
	"github.com/vynaloze/pgstats/nullable"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type PgStatUserTables Config

type pgStatUserTablesPayload struct {
	Table            string         `json:"table"`
	SeqScan          nullable.Int64 `json:"seq_scan"`
	SeqTupFetch      nullable.Int64 `json:"seq_tup_fetch"`
	IdxScan          nullable.Int64 `json:"idx_scan"`
	IdxTupFetch      nullable.Int64 `json:"idx_tup_fetch"`
	LiveTup          nullable.Int64 `json:"live_tup"`
	DeadTup          nullable.Int64 `json:"dead_tup"`
	InsTup           nullable.Int64 `json:"ins_tup"`
	UpdTup           nullable.Int64 `json:"upd_tup"`
	DelTup           nullable.Int64 `json:"del_tup"`
	LastVacuum       nullable.Time  `json:"last_vacuum"`
	LastAutovacuum   nullable.Time  `json:"last_autovacuum"`
	LastAnalyze      nullable.Time  `json:"last_analyze"`
	LastAutoanalyze  nullable.Time  `json:"last_autoanalyze"`
	VacuumCount      nullable.Int64 `json:"vacuum_count"`
	AutovacuumCount  nullable.Int64 `json:"autovacuum_count"`
	AnalyzeCount     nullable.Int64 `json:"analyze_count"`
	AutoanalyzeCount nullable.Int64 `json:"autoanalyze_count"`
}

func (p *PgStatUserTables) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	s, err := datasource.PgStats.PgStatUserTables()
	if err != nil {
		log.Error(err)
	}
	var payload []pgStatUserTablesPayload
	for _, r := range s {
		payload = append(payload, pgStatUserTablesPayload{
			r.Schemaname + "." + r.Relname,
			r.SeqScan,
			r.SeqTupRead,
			r.IdxScan,
			r.IdxTupFetch,
			r.NLiveTup,
			r.NDeadTup,
			r.NTupIns,
			r.NTupUpd,
			r.NTupDel,
			r.LastVacuum,
			r.LastAutovacuum,
			r.LastAnalyze,
			r.LastAutoanalyze,
			r.VacuumCount,
			r.AutovacuumCount,
			r.AnalyzeCount,
			r.AutoanalyzeCount,
		})
	}

	return dto.NewStat(datasource.DsDto, "pg_stat_user_tables", payload)
}

func (p *PgStatUserTables) Conf() Config {
	return Config(*p)
}
