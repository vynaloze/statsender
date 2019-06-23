package collector

import (
	"github.com/vynaloze/pgstats/nullable"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// PgStatArchiver collects statistics about WAL archiver process's activity
type PgStatArchiver Config

type pgStatArchiverPayload struct {
	ArchivedCount    nullable.Int64 `json:"archived_count"`
	LastArchivedTime nullable.Time  `json:"last_archived_time"`
	FailedCount      nullable.Int64 `json:"failed_count"`
	LastFailedTime   nullable.Time  `json:"last_failed_time"`
}

// Collect collects statistics from given datasource
func (p *PgStatArchiver) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	s, err := datasource.PgStats.PgStatArchiver()
	if err != nil {
		log.Error(err)
	}
	payload := pgStatArchiverPayload{
		s.ArchivedCount,
		s.LastArchivedTime,
		s.FailedCount,
		s.LastFailedTime,
	}

	return dto.NewStat(datasource.DsDto, "pg_stat_archiver", []pgStatArchiverPayload{payload})
}

// Conf return the configuration of the collector
func (p *PgStatArchiver) Conf() Config {
	return Config(*p)
}
