package main

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/vynaloze/pgstats"
	"github.com/vynaloze/statsender/collector"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
	"github.com/vynaloze/statsender/sender"
	"os"
	"unsafe"
)

func main() {
	// Prepare logger
	var logErr error
	log, logErr := logger.New()
	if logErr != nil {
		fmt.Print(logErr)
		os.Exit(1)
	}
	// Read configuration
	c, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	// Create datasources
	var datasources []collector.Datasource
	for _, ds := range c.Datasources {
		d := dto.NewDatasource()
		s, err := pgstats.Connect(ds.DbName, ds.Username, ds.Password, pgstats.Host(ds.Host), pgstats.Port(ds.Port), pgstats.SslMode(ds.Ssl))
		if err != nil {
			log.Fatal(err)
		}
		datasources = append(datasources, collector.Datasource{DsDto: *d, Conn: *s})
	}

	// Start cron jobs and wait forever
	startCrons(datasources, c.System.toInterface(), c.Postgres.toInterface(), []sender.Sender{sender.Sout{}}) //fixme
	select {}
}

func startCrons(datasources []collector.Datasource, systemCollectors []collector.SystemCollector, postgresCollectors []collector.PostgresCollector, targets []sender.Sender) {
	log, _ := logger.New()
	crontab := cron.New()

	// System collectors
	for _, c := range systemCollectors {
		if isNilValue(c) {
			continue
		}
		err := crontab.AddFunc(c.Conf().Cron, systemJob(c, targets))
		if err != nil {
			log.Fatalf("Startup error - cron parse failed: %s", err)
		}
	}
	// Postgres collectors
	for _, ds := range datasources {
		for _, p := range postgresCollectors {
			if isNilValue(p) {
				continue
			}
			err := crontab.AddFunc(p.Conf().Cron, postgresJob(ds, p, targets))
			if err != nil {
				log.Fatalf("Startup error - cron parse failed: %s", err)
			}
		}
	}

	crontab.Start()
}

func systemJob(systemCollector collector.SystemCollector, targets []sender.Sender) func() {
	return func() {
		payload := systemCollector.Collect()

		for _, target := range targets {
			target.Send(payload)
		}
	}
}

func postgresJob(datasource collector.Datasource, postgresCollector collector.PostgresCollector, targets []sender.Sender) func() {
	return func() {
		payload := postgresCollector.Collect(datasource)

		for _, target := range targets {
			target.Send(payload)
		}
	}
}

func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}
