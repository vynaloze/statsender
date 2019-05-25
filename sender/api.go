package sender

import "github.com/vynaloze/statsender/dto"

type Payload string

type Sender interface {
	Send(payload *dto.Stat)
}
