// Package dto defines DataTransferObjects, used to communicate between collectors and senders
package dto

import (
	"time"
)

// Stat represents a single batch of statistics,
// collected by single collector and forwarded to one or many senders
type Stat struct {
	Timestamp  int64       `json:"timestamp"`
	Datasource Datasource  `json:"datasource"`
	Id         string      `json:"id"`
	Payload    interface{} `json:"payload"`
}

// NewStat creates a new dto.Stat
func NewStat(datasource *Datasource, id string, payload interface{}) *Stat {
	return &Stat{
		Timestamp:  time.Now().Unix(),
		Datasource: *datasource,
		Id:         id,
		Payload:    payload,
	}
}
