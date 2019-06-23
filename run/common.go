package run

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/vynaloze/pgstats"
	"github.com/vynaloze/statsender/collector"
	"github.com/vynaloze/statsender/config"
	"github.com/vynaloze/statsender/dto"
	"unsafe"
)

type datasourceError struct {
	err error
	ds  config.Datasource
}

type datasourceErrors []datasourceError

func (d datasourceErrors) print() {
	for _, e := range d {
		printError(fmt.Sprintf("could not connect to %s:%d/%s (%s)", e.ds.Host, e.ds.Port, e.ds.DbName, e.err.Error()))
	}
}

func connectToDatasources(c *config.Config) ([]collector.Datasource, datasourceErrors) {
	var datasources []collector.Datasource
	var errs datasourceErrors
	for _, ds := range c.Datasources {
		d := dto.NewPostgresDsDto(ds.Host, ds.Port, ds.DbName, ds.Tags)
		s, err := pgstats.Connect(ds.DbName, ds.Username, ds.Password, optionalParams(ds)...)
		if err != nil {
			errs = append(errs, datasourceError{err: err, ds: ds})
		} else {
			datasources = append(datasources, collector.Datasource{DsDto: d, PgStats: s, ConnectionString: prepareConnectionString(ds)})
		}
	}
	return datasources, errs
}

func optionalParams(ds config.Datasource) []pgstats.Option {
	o := []pgstats.Option{pgstats.Host(ds.Host), pgstats.Port(ds.Port)}
	if ds.SslMode != nil {
		o = append(o, pgstats.SslMode(*ds.SslMode))
	}
	if ds.SslCert != nil {
		o = append(o, pgstats.SslCert(*ds.SslCert))
	}
	if ds.SslRootCert != nil {
		o = append(o, pgstats.SslRootCert(*ds.SslRootCert))
	}
	if ds.SslKey != nil {
		o = append(o, pgstats.SslKey(*ds.SslKey))
	}
	return o
}

func prepareConnectionString(ds config.Datasource) *string {
	baseString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d",
		ds.DbName, ds.Username, ds.Password, ds.Host, ds.Port)
	if ds.SslMode != nil {
		baseString += " sslmode=" + *ds.SslMode
	}
	if ds.SslCert != nil {
		baseString += " sslcert=" + *ds.SslCert
	}
	if ds.SslRootCert != nil {
		baseString += " sslrootcert=" + *ds.SslRootCert
	}
	if ds.SslKey != nil {
		baseString += " sslkey=" + *ds.SslKey
	}
	return &baseString
}

func printError(msg string) {
	bold := color.New(color.FgRed).Add(color.Bold)
	reg := color.New(color.FgRed)
	_, _ = bold.Print("[FAIL] ")
	_, _ = reg.Println(msg)
}

func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}
