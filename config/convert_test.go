package config

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

var parseDSNCorrectTestTable = []struct {
	dsn      string
	tags     []string
	expected *Datasource
}{
	{
		"postgresql://user:pass@localhost:5432/db",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "localhost", Port: 5432, DbName: "db", SslMode: "",
			Tags: map[string]string{}},
	},
	{
		"postgresql://user:pass@localhost:5432/db?sslmode=disable",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "localhost", Port: 5432, DbName: "db", SslMode: "disable",
			Tags: map[string]string{}},
	},
	{
		"postgresql://user:pass@localhost:6432/other_db",
		[]string{"key=value", "foo=bar"},
		&Datasource{Username: "user", Password: "pass", Host: "localhost", Port: 6432, DbName: "other_db", SslMode: "",
			Tags: map[string]string{"key": "value", "foo": "bar"}},
	},
	{
		"user:pass@10.0.1.1:6432/db",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "10.0.1.1", Port: 6432, DbName: "db", SslMode: "",
			Tags: map[string]string{}},
	},
	{
		"postgresql://user:pass@10.0.1.1/db",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "10.0.1.1", Port: 5432, DbName: "db", SslMode: "",
			Tags: map[string]string{}},
	},
}

var parseDSNFailTestTable = []struct {
	dsn  string
	tags []string
}{
	{"opstgresql://user:pass@localhost:5432/db", []string{}},
	{"postgresql://localhost:5432/db", []string{}},
	{"postgresql://user@localhost:5432/db", []string{}},
	{"postgresql://user:pass@localhost:5432", []string{}},
	{"postgresql://user:pass@localhost:5432/db", []string{"keyvalue"}},
}

func TestParseDSN(t *testing.T) {
	for i, tt := range parseDSNCorrectTestTable {
		actual, err := ParseDSN(tt.dsn, tt.tags)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(actual, tt.expected, cmp.Options{}) {
			t.Errorf("%d: Expected '%v'; actual '%v'", i, tt.expected, actual)
		}
	}
}

func TestFailParseDSN(t *testing.T) {
	for i, tt := range parseDSNFailTestTable {
		actual, err := ParseDSN(tt.dsn, tt.tags)
		if err == nil {
			t.Errorf("%d: Expected error, but it didn't happened. Actual data: %v", i, actual)
		}
	}
}