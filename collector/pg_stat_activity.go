package collector

import (
	"github.com/vynaloze/pgstats/nullable"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// PgStatActivity collects statistics about current server process activity
type PgStatActivity Config

type pgStatActivityPayload struct {
	Pid             int64           `json:"pid"`
	Usename         nullable.String `json:"usename"`
	ApplicationName nullable.String `json:"application_name"`
	ClientAddr      nullable.String `json:"client_addr"`
	Datname         nullable.String `json:"datname"`
	State           nullable.String `json:"state"`
	WaitEventType   nullable.String `json:"wait_event_type"`
	WaitEvent       nullable.String `json:"wait_event"`
	QueryStart      nullable.Time   `json:"query_start"`
	Query           nullable.String `json:"query"`
}

// Collect collects statistics from given datasource
func (p *PgStatActivity) Collect(datasource *Datasource) *dto.Stat {
	log, _ := logger.New()

	s, err := datasource.PgStats.PgStatActivity()
	if err != nil {
		log.Error(err)
	}
	var payload []pgStatActivityPayload
	for _, r := range s {
		payload = append(payload, pgStatActivityPayload{
			r.Pid,
			r.Usename,
			r.ApplicationName,
			r.ClientAddr,
			r.Datname,
			r.State,
			r.WaitEventType,
			r.WaitEvent,
			r.QueryStart,
			r.Query,
		})
	}

	return dto.NewStat(datasource.DsDto, "pg_stat_activity", payload)
}

// Conf return the configuration of the collector
func (p *PgStatActivity) Conf() Config {
	return Config(*p)
}
