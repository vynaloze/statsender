package main

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/pkg/errors"
	"github.com/vynaloze/statsender/collector"
	"github.com/vynaloze/statsender/logger"
	"os"
	"path/filepath"
)

type config struct {
	Debug       *bool        `hcl:"debug"`
	Datasources []datasource `hcl:"datasource,block"`
	System      system       `hcl:"system,block"`
	Postgres    postgres     `hcl:"postgres,block"`
	// todo targets
}

type datasource struct {
	Host     string `hcl:"host"`
	Port     int    `hcl:"port"`
	Username string `hcl:"username"`
	Password string `hcl:"password"`
	DbName   string `hcl:"dbname"`
	Ssl      string `hcl:"ssl"`
	// todo tags
}

type system struct {
	Cpu       *collector.Cpu       `hcl:"cpu,block"`
	VirtMem   *collector.VirtMem   `hcl:"virtual_memory,block"`
	SwapMem   *collector.SwapMem   `hcl:"swap_memory,block"`
	DiskIo    *collector.DiskIo    `hcl:"disk_io,block"`
	DiskUsage *collector.DiskUsage `hcl:"disk_usage,block"`
	Net       *collector.Net       `hcl:"network_io,block"`
	Load      *collector.Load      `hcl:"load,block"`
}

type postgres struct {
	PgStatUserIndexes *collector.PgStatUserIndexes `hcl:"pg_stat_user_indexes,block"`
}

func readConfig() (*config, error) {
	log, _ := logger.New()

	// Find, parse, merge and decode all files
	parser := hclparse.NewParser()
	var files []*hcl.File
	var diags hcl.Diagnostics

	for _, fn := range allHclFiles() {
		f, moreDiags := parser.ParseHCLFile(fn)
		files = append(files, f)
		diags = append(diags, moreDiags...)
	}

	body := hcl.MergeFiles(files)
	var c config
	moreDiags := gohcl.DecodeBody(body, nil, &c)
	diags = append(diags, moreDiags...)

	// Check for errors
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,      // writer to send messages to
			parser.Files(), // the parser's file cache, for source snippets
			78,             // wrapping width
			true,           // generate colored/highlighted output
		)
		err := wr.WriteDiagnostics(diags)
		if err != nil {
			log.Error(err)
		}
		return nil, errors.New("invalid configuration - see diagnostics above")
	}

	// Set log level
	if c.Debug == nil {
		logger.SetDebug(false)
	} else {
		logger.SetDebug(*c.Debug)
	}

	return &c, nil
}

func allHclFiles() []string {
	log, _ := logger.New()
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var files []string
	err = filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ".hcl" {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func (s *system) toInterface() []collector.SystemCollector {
	return []collector.SystemCollector{
		s.Cpu,
		s.VirtMem,
		s.SwapMem,
		s.DiskIo,
		s.DiskUsage,
		s.Net,
		s.Load,
	}
}

func (p *postgres) toInterface() []collector.PostgresCollector {
	return []collector.PostgresCollector{
		p.PgStatUserIndexes,
	}
}
