package sender

import "github.com/vynaloze/statsender/dto"

type Sender interface {
	Send(payload *dto.Stat)
}
