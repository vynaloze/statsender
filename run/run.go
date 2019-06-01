package run

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/vynaloze/pgstats"
	"github.com/vynaloze/statsender/collector"
	"github.com/vynaloze/statsender/config"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
	"github.com/vynaloze/statsender/sender"
	"os"
	"unsafe"
)

func Run(confDir string) {
	// Prepare logger
	var logErr error
	log, logErr := logger.New()
	if logErr != nil {
		fmt.Print(logErr)
		os.Exit(1)
	}
	log.Debug("reading configuration")
	c, cErr := config.ReadConfig(confDir)
	if cErr != nil {
		log.Fatal(cErr)
		os.Exit(1)
	}
	log.Debug("connecting to datasources")
	var datasources []collector.Datasource
	for _, ds := range c.Datasources {
		d := dto.NewPostgresDsDto(ds.Host, ds.Port, ds.DbName, ds.Tags)
		s, err := pgstats.Connect(ds.DbName, ds.Username, ds.Password, optionalParams(ds)...)
		if err != nil {
			log.Fatal(err)
		}
		datasources = append(datasources, collector.Datasource{DsDto: d, Conn: s})
	}

	log.Debug("starting collector jobs in the background")
	startCrons(datasources, c.System.ToInterface(), c.Postgres.ToInterface(), c.SendersToInterface())
	// wait forever
	select {}
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

func startCrons(datasources []collector.Datasource, systemCollectors []collector.Collector, postgresCollectors []collector.Collector, targets []sender.Sender) {
	log, _ := logger.New()
	crontab := cron.New()

	// System collectors
	for _, c := range systemCollectors {
		if isNilValue(c) || !c.Conf().Enabled {
			log.Debugf("skipping disabled collector %+v", c)
			continue
		}
		log.Debugf("scheduling collector %+v", c)
		err := crontab.AddFunc(c.Conf().Cron, newJob(collector.Datasource{DsDto: dto.NewSystemDsDto(), Conn: nil}, c, targets))
		if err != nil {
			log.Fatalf("Startup error - cron parse failed: %s", err)
		}
	}
	// Postgres collectors
	for _, ds := range datasources {
		log.Debugf("scheduling collectors for datasource %+v", ds)
		for _, p := range postgresCollectors {
			if isNilValue(p) || !p.Conf().Enabled {
				log.Debugf("skipping disabled collector %+v", p)
				continue
			}
			log.Debugf("scheduling collector %+v", p)
			err := crontab.AddFunc(p.Conf().Cron, newJob(ds, p, targets))
			if err != nil {
				log.Fatalf("Startup error - cron parse failed: %s", err)
			}
		}
	}

	crontab.Start()
}

func newJob(datasource collector.Datasource, collector collector.Collector, targets []sender.Sender) func() {
	return func() {
		payload := collector.Collect(&datasource)

		for _, target := range targets {
			go target.Send(payload)
		}
	}
}

func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}
