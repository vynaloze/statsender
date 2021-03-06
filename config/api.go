// Package config handles reading and generating configuration files
package config

import (
	"github.com/vynaloze/statsender/collector"
	"github.com/vynaloze/statsender/sender"
)

// Config represents the top-level configuration
type Config struct {
	Datasources []Datasource  `hcl:"datasource,block"`
	System      *System       `hcl:"system,block"`
	Postgres    *Postgres     `hcl:"postgres,block"`
	Sout        *sender.Sout  `hcl:"console,block"`
	Http        []sender.Http `hcl:"http,block"`
}

// Datasource represents a single PostgreSQL database connection info
type Datasource struct {
	Host        string            `hcl:"host"`
	Port        int               `hcl:"port"`
	Username    string            `hcl:"username"`
	Password    string            `hcl:"password"`
	DbName      string            `hcl:"dbname"`
	SslMode     *string           `hcl:"sslmode"`
	SslKey      *string           `hcl:"sslkey"`
	SslCert     *string           `hcl:"sslcert"`
	SslRootCert *string           `hcl:"sslrootcert"`
	Tags        map[string]string `hcl:"tags"`
}

// System represents configuration of system collectors
type System struct {
	Cpu     *collector.Cpu     `hcl:"cpu,block"`
	VirtMem *collector.VirtMem `hcl:"virtual_memory,block"`
	SwapMem *collector.SwapMem `hcl:"swap_memory,block"`
	Disk    *collector.Disk    `hcl:"disk,block"`
	Net     *collector.Net     `hcl:"network,block"`
	Load    *collector.Load    `hcl:"load,block"`
}

// Postgres represents configuration of PostgreSQL collectors
type Postgres struct {
	PgStatStatements  *collector.PgStatStatements  `hcl:"pg_stat_statements,block"`
	PgStatUserTables  *collector.PgStatUserTables  `hcl:"pg_stat_user_tables,block"`
	PgStatUserIndexes *collector.PgStatUserIndexes `hcl:"pg_stat_user_indexes,block"`
	PgStatActivity    *collector.PgStatActivity    `hcl:"pg_stat_activity,block"`
	PgLocks           *collector.PgLocks           `hcl:"pg_locks,block"`
	PgStatArchiver    *collector.PgStatArchiver    `hcl:"pg_stat_archiver,block"`
}
