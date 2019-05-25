package dto

import (
	"time"
)

type Stat struct {
	Timestamp  int64       `json:"timestamp"`
	Datasource Datasource  `json:"datasource"`
	Id         string      `json:"id"`
	Payload    interface{} `json:"payload"`
}

func NewStat(datasource *Datasource, id string, payload interface{}) *Stat {
	return &Stat{
		Timestamp:  time.Now().Unix(),
		Datasource: *datasource,
		Id:         id,
		Payload:    payload,
	}
}
