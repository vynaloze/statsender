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
		p.PgStatStatements,
		p.PgStatUserTables,
		p.PgStatUserIndexes,
		p.PgStatActivity,
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
	ds := &Datasource{Host: host, Port: int(port), Username: user, Password: pass,
		DbName: db}
	t, err := parseTags(tags)
	if err != nil {
		return nil, err
	}
	ds.Tags = t
	sslMode := u.Query().Get("sslmode")
	if sslMode != "" {
		ds.SslMode = &sslMode
	}
	sslKey := u.Query().Get("sslkey")
	if sslKey != "" {
		ds.SslKey = &sslKey
	}
	sslCert := u.Query().Get("sslcert")
	if sslCert != "" {
		ds.SslCert = &sslCert
	}
	sslRootCert := u.Query().Get("sslrootcert")
	if sslRootCert != "" {
		ds.SslRootCert = &sslRootCert
	}
	return ds, nil
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

func ParseSender(args []string) (*sender.Sender, error) {
	typ := args[0]
	switch typ {
	case "console":
		var s sender.Sender
		s = sender.Sout{}
		return &s, nil
	case "http":
		if len(args) < 2 {
			return nil, errors.New("sender not specified")
		}
		u, err := url.Parse(args[1])
		if err != nil {
			return nil, err
		}
		if u.Scheme == "" {
			return nil, errors.New("scheme not specified")
		}
		if u.Host == "" {
			return nil, errors.New("host not specified")
		}
		var s sender.Sender
		r, _ := strconv.ParseInt(args[2], 10, 0)
		d, _ := strconv.ParseInt(args[3], 10, 0)
		s = sender.Http{Target: args[1], MaxRetries: int(r), RetryDelay: int(d)}
		return &s, nil
	default:
		return nil, errors.New("invalid sender type - valid types: 'console', 'http'")
	}
}
