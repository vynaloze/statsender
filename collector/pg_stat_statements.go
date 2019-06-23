package collector

import (
	"database/sql"
	// import to register driver
	_ "github.com/lib/pq"
	"github.com/vynaloze/pgstats/nullable"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// PgStatStatements collects execution statistics of all SQL statements executed by a server
type PgStatStatements Config

type pgStatStatementsPayload struct {
	Query          string          `json:"query"`
	User           nullable.String `json:"user"`
	Calls          int64           `json:"calls"`
	Rows           int64           `json:"rows"`
	AvgTime        float64         `json:"avg_time"`
	MinTime        float64         `json:"min_time"`
	MaxTime        float64         `json:"max_time"`
	SharedBlksHit  int64           `json:"shared_blks_hit"`
	SharedBlksRead int64           `json:"shared_blks_read"`
	LocalBlksHit   int64           `json:"local_blks_hit"`
	LocalBlksRead  int64           `json:"local_blks_read"`
}

// Collect collects statistics from given datasource
func (p *PgStatStatements) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	db, err := sql.Open("postgres", *datasource.ConnectionString)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Error(err)
	}

	query := "select s.query,u.rolname,s.calls,s.rows,s.mean_time,s.min_time,s.max_time," +
		" s.shared_blks_hit,s.shared_blks_read,s.local_blks_hit,s.local_blks_read" +
		" from pg_stat_statements s left join pg_authid u on s.userid = u.oid"
	rows, err := db.Query(query)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()

	var payload []pgStatStatementsPayload
	for rows.Next() {
		p := new(pgStatStatementsPayload)
		err := rows.Scan(&p.Query, &p.User, &p.Calls, &p.Rows, &p.AvgTime, &p.MinTime, &p.MaxTime,
			&p.SharedBlksHit, &p.SharedBlksRead, &p.LocalBlksHit, &p.LocalBlksRead)
		if err != nil {
			log.Error(err)
		}
		payload = append(payload, *p)
	}

	return dto.NewStat(datasource.DsDto, "pg_stat_statements", payload)
}

// Conf return the configuration of the collector
func (p *PgStatStatements) Conf() Config {
	return Config(*p)
}
