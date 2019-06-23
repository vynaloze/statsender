// Package sender contains all statistic senders
package sender

import "github.com/vynaloze/statsender/dto"

// All senders must satisfy sender.Sender interface
type Sender interface {
	Send(payload *dto.Stat)
	Test() error
}
