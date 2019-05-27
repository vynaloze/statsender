package config

import (
	"github.com/vynaloze/statsender/collector"
	"github.com/vynaloze/statsender/sender"
	"net/url"
	"strconv"
	"strings"
)

func (s *System) ToInterface() []collector.Collector {
	return []collector.Collector{
		s.Cpu,
		s.VirtMem,
		s.SwapMem,
		s.DiskIo,
		s.DiskUsage,
		s.Net,
		s.Load,
	}
}

func (p *Postgres) ToInterface() []collector.Collector {
	return []collector.Collector{
		p.PgStatUserIndexes,
	}
}

func (c *Config) SendersToInterface() []sender.Sender {
	var s []sender.Sender
	if c.Sout != nil {
		s = append(s, c.Sout)
	}
	for _, h := range c.Http {
		s = append(s, h)
	}
	return s
}

func ParseDSN(dsn string) (*Datasource, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	host := u.Hostname()
	var port int64
	if u.Port() == "" {
		port = 5432
	} else {
		port, err = strconv.ParseInt(u.Port(), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	user := u.User.Username()
	pass, _ := u.User.Password()
	db := strings.TrimPrefix(u.Path, "/")
	// todo support not only sslmode
	sslMode := u.Query().Get("sslmode")

	return &Datasource{Host: host, Port: int(port), Username: user, Password: pass,
		DbName: db, Ssl: sslMode}, nil
}
