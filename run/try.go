package run

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/robfig/cron"
	"github.com/vynaloze/statsender/config"
	"os"
	"reflect"
	"time"
)

// Try dry-runs the application, allowing to verify the provided configuration
func Try(confDir string) {
	ok := true
	fmt.Println("[INFO] reading configuration...")
	c, cErr := config.ReadConfig(confDir)
	if cErr != nil {
		err := cErr.Print()
		if err != nil {
			printError(err.Error())
		}
		os.Exit(1)
	}
	printSuccess("configuration structure is valid\n")

	fmt.Println("[INFO] testing collectors...")
	if testConnections(c) {
		printSuccess("collector structure is valid\n")
	} else {
		printError("collector structure is invalid\n")
		ok = false
	}

	fmt.Println("[INFO] testing datasources...")
	datasources, errs := connectToDatasources(c)
	for _, ds := range datasources {
		printSuccess(fmt.Sprintf("connected to %s:%d/%s", ds.DsDto.Ip, *ds.DsDto.Port, *ds.DsDto.Database))
	}
	if errs != nil {
		errs.print()
		printError("datasource structure is invalid\n")
		ok = false
	} else {
		printSuccess("datasource structure is valid\n")
	}
	fmt.Println("[INFO] testing senders...")
	var sErrs []error
	for _, s := range c.SendersToInterface() {
		err := s.Test()
		if err != nil {
			sErrs = append(sErrs, err)
		}
	}
	if sErrs != nil {
		for _, e := range sErrs {
			printError(e.Error())
		}
		printError("sender structure is invalid\n")
		ok = false
	} else {
		printSuccess("sender structure is valid\n")
	}

	if ok {
		fmt.Println("[INFO] Test complete! Looks like you are good to go!")
	} else {
		fmt.Println("[INFO] Test complete! Looks like you have some errors")
		os.Exit(1)
	}
}

func printSuccess(msg string) {
	bold := color.New(color.FgGreen).Add(color.Bold)
	reg := color.New(color.FgGreen)
	_, _ = bold.Print("[OK] ")
	_, _ = reg.Println(msg)
}

func testConnections(c *config.Config) bool {
	valid := true
	for _, s := range append(c.System.ToInterface(), c.Postgres.ToInterface()...) {
		name := reflect.TypeOf(s).String()
		switch {
		case isNilValue(s):
			printSuccess(fmt.Sprintf("%s is disabled (deleted or commented out)", name))
		case !s.Conf().Enabled:
			printSuccess(fmt.Sprintf("%s is disabled", name))
		case !config.IsCronValid(s.Conf().Cron):
			printError(fmt.Sprintf("%s has invalid cron expression", name))
			valid = false
		default:
			schedule, _ := cron.Parse(s.Conf().Cron)
			t1 := schedule.Next(time.Now())
			t2 := schedule.Next(t1)
			t3 := schedule.Next(t2)
			printSuccess(fmt.Sprintf("%s will run at:\n\t%s,\n\t%s,\n\t%s,\n\t...", name, t1.String(), t2.String(), t3.String()))
		}
	}
	return valid
}
