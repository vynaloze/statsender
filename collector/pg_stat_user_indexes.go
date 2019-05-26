package collector

import (
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type PgStatUserIndexes Config

func (p *PgStatUserIndexes) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	s, err := datasource.Conn.PgStatUserIndexes()
	if err != nil {
		log.Error(err)
	}

	return dto.NewStat(datasource.DsDto, "pg_stat_user_indexes", s)
}

func (p *PgStatUserIndexes) Conf() Config {
	return Config(*p)
}
