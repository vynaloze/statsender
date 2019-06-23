package dto

import (
	"github.com/vynaloze/statsender/logger"
	"net"
	"os"
)

// Datasource represents a single data source from which a stat was gathered
type Datasource struct {
	Ip       string            `json:"ip"`
	Hostname *string           `json:"hostname"`
	Port     *int              `json:"port"`
	Database *string           `json:"database"`
	Tags     map[string]string `json:"tags"`
}

// NewSystemDsDto creates a new dto.Datasource from a System datasource
func NewSystemDsDto() *Datasource {
	ds := Datasource{}
	ds.Ip = getOutboundIP()
	ds.Hostname = getHostname()
	return &ds
}

// NewPostgresDsDto creates a new dto.Datasource from a Postgres datasource
func NewPostgresDsDto(host string, port int, dbname string, tags map[string]string) *Datasource {
	ds := Datasource{}
	if host == "localhost" || host == "127.0.0.1" {
		ds.Ip = getOutboundIP()
		ds.Hostname = getHostname()
	} else {
		ds.Ip = host
		ds.Hostname = nil
	}
	ds.Port = &port
	ds.Database = &dbname
	ds.Tags = tags
	return &ds
}

// Get preferred outbound ip of this machine
func getOutboundIP() string {
	log, _ := logger.New()
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func getHostname() *string {
	log, _ := logger.New()
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &name
}
