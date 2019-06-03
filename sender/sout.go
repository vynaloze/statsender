package sender

import (
	"encoding/json"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

type Sout struct {
}

func (_ Sout) Send(payload *dto.Stat) {
	log, _ := logger.New()
	p, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
	}
	data := string(p)
	log.Info(data)
}

func (_ Sout) Test() error {
	return nil
}
