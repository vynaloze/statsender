package config

import (
	"github.com/vynaloze/statsender/collector"
	"github.com/vynaloze/statsender/sender"
)

type Config struct {
	Debug       *bool         `hcl:"debug"`
	Datasources []Datasource  `hcl:"datasource,block"`
	System      *System       `hcl:"system,block"`
	Postgres    *Postgres     `hcl:"postgres,block"`
	Sout        *sender.Sout  `hcl:"console,block"`
	Http        []sender.Http `hcl:"http,block"`
}

type Datasource struct {
	Host     string            `hcl:"host"`
	Port     int               `hcl:"port"`
	Username string            `hcl:"username"`
	Password string            `hcl:"password"`
	DbName   string            `hcl:"dbname"`
	SslMode  string            `hcl:"sslmode"`
	Tags     map[string]string `hcl:"tags"`
}

type System struct {
	Cpu       *collector.Cpu       `hcl:"cpu,block"`
	VirtMem   *collector.VirtMem   `hcl:"virtual_memory,block"`
	SwapMem   *collector.SwapMem   `hcl:"swap_memory,block"`
	DiskIo    *collector.DiskIo    `hcl:"disk_io,block"`
	DiskUsage *collector.DiskUsage `hcl:"disk_usage,block"`
	Net       *collector.Net       `hcl:"network,block"`
	Load      *collector.Load      `hcl:"load,block"`
}

type Postgres struct {
	PgStatUserIndexes *collector.PgStatUserIndexes `hcl:"pg_stat_user_indexes,block"`
}
