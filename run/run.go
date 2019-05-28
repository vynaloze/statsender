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
	// Read configuration
	c, cErr := config.ReadConfig(confDir)
	if cErr != nil {
		fmt.Print(cErr)
		os.Exit(1)
	}
	// Create datasources
	var datasources []collector.Datasource
	for _, ds := range c.Datasources {
		d := dto.NewPostgresDsDto(ds.Host, ds.Port, ds.DbName, ds.Tags)
		s, err := pgstats.Connect(ds.DbName, ds.Username, ds.Password, pgstats.Host(ds.Host), pgstats.Port(ds.Port), pgstats.SslMode(ds.SslMode))
		if err != nil {
			log.Fatal(err)
		}
		datasources = append(datasources, collector.Datasource{DsDto: d, Conn: s})
	}

	// Start cron jobs and wait forever
	startCrons(datasources, c.System.ToInterface(), c.Postgres.ToInterface(), c.SendersToInterface())
	select {}
}

func startCrons(datasources []collector.Datasource, systemCollectors []collector.Collector, postgresCollectors []collector.Collector, targets []sender.Sender) {
	log, _ := logger.New()
	crontab := cron.New()

	// System collectors
	for _, c := range systemCollectors {
		if isNilValue(c) {
			continue
		}
		err := crontab.AddFunc(c.Conf().Cron, newJob(collector.Datasource{DsDto: dto.NewSystemDsDto(), Conn: nil}, c, targets))
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
