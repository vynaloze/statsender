package config

import (
	"github.com/pkg/errors"
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

func ParseDSN(dsn string, tags []string) (*Datasource, error) {
	if !strings.HasPrefix(dsn, "postgresql://") {
		dsn = strings.Join([]string{"postgresql://", dsn}, "")
	}
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	host := u.Hostname()
	if host == "" {
		return nil, errors.New("hostname not set")
	}
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
	if user == "" {
		return nil, errors.New("username not set")
	}
	pass, set := u.User.Password()
	if !set {
		return nil, errors.New("password not set")
	}
	db := strings.TrimPrefix(u.Path, "/")
	if db == "" {
		return nil, errors.New("dbname not set")
	}
	// todo support not only sslmode
	sslMode := u.Query().Get("sslmode")

	t, err := parseTags(tags)
	return &Datasource{Host: host, Port: int(port), Username: user, Password: pass,
		DbName: db, SslMode: sslMode, Tags: t}, err
}

func parseTags(tags []string) (map[string]string, error) {
	m := make(map[string]string)
	for _, t := range tags {
		s := strings.Split(t, "=")
		if len(s) != 2 {
			return nil, errors.New("invalid format of a tag")
		}
		m[s[0]] = s[1]
	}
	return m, nil
}
