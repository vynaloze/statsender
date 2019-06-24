package config

import (
	"github.com/google/go-cmp/cmp"
	"github.com/vynaloze/statsender/sender"
	"testing"
)

var sslDisable = "disable"
var sslCert = "/foo/bar.crt"

var parseDSNCorrectTestTable = []struct {
	dsn      string
	tags     []string
	expected *Datasource
}{
	{
		"postgresql://user:pass@localhost:5432/db",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "localhost", Port: 5432, DbName: "db",
			Tags: map[string]string{}},
	},
	{
		"postgresql://user:pass@localhost:5432/db?sslmode=disable",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "localhost", Port: 5432, DbName: "db", SslMode: &sslDisable,
			Tags: map[string]string{}},
	},
	{
		"postgresql://user:pass@localhost:5432/db?sslmode=disable&sslcert=/foo/bar.crt",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "localhost", Port: 5432, DbName: "db", SslMode: &sslDisable,
			SslCert: &sslCert, Tags: map[string]string{}},
	},
	{
		"postgresql://user:pass@localhost:6432/other_db",
		[]string{"key=value", "foo=bar"},
		&Datasource{Username: "user", Password: "pass", Host: "localhost", Port: 6432, DbName: "other_db",
			Tags: map[string]string{"key": "value", "foo": "bar"}},
	},
	{
		"user:pass@10.0.1.1:6432/db",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "10.0.1.1", Port: 6432, DbName: "db",
			Tags: map[string]string{}},
	},
	{
		"postgresql://user:pass@10.0.1.1/db",
		[]string{},
		&Datasource{Username: "user", Password: "pass", Host: "10.0.1.1", Port: 5432, DbName: "db",
			Tags: map[string]string{}},
	},
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

func TestFailParseDSN(t *testing.T) {
	for i, tt := range parseDSNFailTestTable {
		actual, err := ParseDSN(tt.dsn, tt.tags)
		if err == nil {
			t.Errorf("%d: Expected error, but it didn't happened. Actual data: %v", i, actual)
		}
	}
}

var parseSenderCorrectTestTable = []struct {
	typ      string
	spec     string
	flags    []string
	expected sender.Sender
}{
	{
		"console",
		"",
		[]string{},
		sender.Sout{},
	},
	{
		"console",
		"should be ignored",
		[]string{},
		sender.Sout{},
	},
	{
		"http",
		"http://10.0.1.2:8080/stats",
		[]string{"3", "7"},
		sender.Http{Target: "http://10.0.1.2:8080/stats", MaxRetries: 3, RetryDelay: 7},
	},
	{
		"http",
		"http://localhost:8080/stats",
		[]string{"3", "7"},
		sender.Http{Target: "http://localhost:8080/stats", MaxRetries: 3, RetryDelay: 7},
	},
	{
		"http",
		"http://10.0.1.2:8080/stats/v2",
		[]string{"3", "7"},
		sender.Http{Target: "http://10.0.1.2:8080/stats/v2", MaxRetries: 3, RetryDelay: 7},
	},
	{
		"http",
		"http://10.0.1.2/stats",
		[]string{"3", "7"},
		sender.Http{Target: "http://10.0.1.2/stats", MaxRetries: 3, RetryDelay: 7},
	},
	{
		"http",
		"http://10.0.1.2:8080",
		[]string{"3", "7"},
		sender.Http{Target: "http://10.0.1.2:8080", MaxRetries: 3, RetryDelay: 7},
	},
	{
		"http",
		"http://10.0.1.2",
		[]string{"3", "7"},
		sender.Http{Target: "http://10.0.1.2", MaxRetries: 3, RetryDelay: 7},
	},
}

func TestParseSender(t *testing.T) {
	for i, tt := range parseSenderCorrectTestTable {
		actual, err := ParseSender(append([]string{tt.typ, tt.spec}, tt.flags...))
		if err != nil {
			t.Error(err)
		}
		if *actual != tt.expected {
			t.Errorf("%d: Expected '%v'; actual '%v'", i, tt.expected, *actual)
		}
	}
}

var parseSenderFailTestTable = []struct {
	typ  string
	spec string
}{
	{"unknown", ""},
	{"http", "/stats"},
	{"http", "10.0.1.2:8080/stats"},
}

func TestFailParseSender(t *testing.T) {
	for i, tt := range parseSenderFailTestTable {
		actual, err := ParseSender([]string{tt.typ, tt.spec})
		if err == nil {
			t.Errorf("%d: Expected error, but it didn't happened. Actual data: %v", i, *actual)
		}
	}
}
