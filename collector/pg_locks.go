package collector

import (
	"database/sql"
	// import to register driver
	_ "github.com/lib/pq"
	"github.com/vynaloze/pgstats/nullable"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// PgLocks collects statistics about locks
type PgLocks Config

type pgLocksPayload struct {
	GrantedLocks    int64       `json:"granted_locks"`
	NotGrantedLocks int64       `json:"not_granted_locks"`
	LocksInfo       []locksInfo `json:"locks_info"`
}

type locksInfo struct {
	BlockedPid          int64           `json:"blocked_pid"`
	BlockedUser         nullable.String `json:"blocked_user"`
	BlockedApplication  nullable.String `json:"blocked_application"`
	BlockedClient       nullable.String `json:"blocked_client"`
	BlockedDatabase     nullable.String `json:"blocked_database"`
	BlockedWaitingSince nullable.Time   `json:"blocked_waiting_since"`
	BlockedQuery        nullable.String `json:"blocked_query"`

	BlockingPid         int64           `json:"blocking_pid"`
	BlockingUser        nullable.String `json:"blocking_user"`
	BlockingApplication nullable.String `json:"blocking_application"`
	BlockingClient      nullable.String `json:"blocking_client"`
	BlockingDatabase    nullable.String `json:"blocking_database"`
	BlockingLockMode    nullable.String `json:"blocking_lock_mode"`
	BlockingQuery       nullable.String `json:"blocking_query"`
}

// Collect collects statistics from given datasource
func (p *PgLocks) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	db, err := sql.Open("postgres", *datasource.ConnectionString)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Error(err)
	}

	payload := pgLocksPayload{}
	if err = payload.getLocksCount(db); err != nil {
		log.Error(err)
	}
	if err = payload.getLocksInfo(db); err != nil {
		log.Error(err)
	}

	return dto.NewStat(datasource.DsDto, "pg_locks", []pgLocksPayload{payload})
}

// Conf return the configuration of the collector
func (p *PgLocks) Conf() Config {
	return Config(*p)
}

func (p *pgLocksPayload) getLocksCount(db *sql.DB) error {
	query := "select granted, count(granted) from pg_locks group by granted"
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		granted := new(bool)
		count := new(int64)
		err := rows.Scan(granted, count)
		if err != nil {
			return err
		}
		if *granted {
			p.GrantedLocks = *count
		} else {
			p.NotGrantedLocks = *count
		}
	}
	return nil
}

func (p *pgLocksPayload) getLocksInfo(db *sql.DB) error {
	query := `SELECT
    blocked_locks.pid     				AS blocked_pid,
    blocked_activity.usename  			AS blocked_user,
    blocked_activity.application_name 	AS blocked_application,
    blocked_activity.client_addr		AS blocked_client,
    blocked_activity.datname			AS blocked_database,
    blocked_activity.query_start		AS blocked_waiting_since,
    blocked_activity.query    			AS blocked_statement,

    blocking_locks.pid    				AS blocking_pid,
    blocking_activity.usename 			AS blocking_user,
    blocking_activity.application_name 	AS blocking_application,
    blocking_activity.client_addr		AS blocking_client,
    blocking_activity.datname			AS blocking_database,
    blocking_locks.mode					AS blocking_lock_type,
    blocking_activity.query    			AS blocking_statement

FROM  pg_catalog.pg_locks         blocked_locks
          JOIN pg_catalog.pg_stat_activity blocked_activity  ON blocked_activity.pid = blocked_locks.pid
          JOIN pg_catalog.pg_locks         blocking_locks
               ON blocking_locks.locktype = blocked_locks.locktype
                   AND blocking_locks.DATABASE IS NOT DISTINCT FROM blocked_locks.DATABASE
                   AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
                   AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
                   AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
                   AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
                   AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
                   AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
                   AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
                   AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
                   AND blocking_locks.pid != blocked_locks.pid

          JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.GRANTED;`

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var infos []locksInfo
	for rows.Next() {
		i := new(locksInfo)
		err := rows.Scan(&i.BlockedPid, &i.BlockedUser, &i.BlockedApplication, &i.BlockedClient, &i.BlockedDatabase, &i.BlockedWaitingSince, &i.BlockedQuery,
			&i.BlockingPid, &i.BlockingUser, &i.BlockingApplication, &i.BlockingClient, &i.BlockingDatabase, &i.BlockingLockMode, &i.BlockingQuery)
		if err != nil {
			return err
		}
		infos = append(infos, *i)
	}

	p.LocksInfo = infos
	return nil
}
