package run

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/vynaloze/statsender/collector"
	"github.com/vynaloze/statsender/config"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
	"github.com/vynaloze/statsender/sender"
	"os"
)

// Run starts the main flow of the application and waits infinitely
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
		err := cErr.Print()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
	log.Debug("connecting to datasources")
	datasources, errs := connectToDatasources(c)
	if errs != nil {
		errs.print()
		os.Exit(1)
	}

	log.Debug("starting collector jobs in the background")
	startCrons(datasources, c.System.ToInterface(), c.Postgres.ToInterface(), c.SendersToInterface())
	// wait forever
	select {}
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
		err := crontab.AddFunc(c.Conf().Cron, newJob(collector.Datasource{DsDto: dto.NewSystemDsDto(), PgStats: nil, ConnectionString: nil}, c, targets))
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
