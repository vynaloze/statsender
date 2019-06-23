package sender

import (
	"encoding/json"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
)

// Sout sender prints the gathered statistics to the standard output.
// It is useful for testing/debugging purposes
type Sout struct {
}

// Send prints the dto.Stat in a JSON format to the standard output
func (Sout) Send(payload *dto.Stat) {
	log, _ := logger.New()
	p, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
	}
	data := string(p)
	log.Info(data)
}

// Test tests the connection - in case of Sout sender, it is always without errors
func (Sout) Test() error {
	return nil
}
