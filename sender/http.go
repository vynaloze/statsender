package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
	"net/http"
	"time"
)

type Http struct {
	URL        string `hcl:"url"`
	Endpoint   string `hcl:"endpoint"`
	RetryDelay int    `hcl:"retryDelay"`
	MaxRetries int    `hcl:"maxRetries"`
	retries    int
}

func (h Http) Send(payload *dto.Stat) {
	log, _ := logger.New()

	p, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
	}
	data := string(p)
	log.Debugf("Target:%s; Data: %s", h.URL+h.Endpoint, data)

	response, err := http.Post(h.URL+h.Endpoint, "application/json", bytes.NewBuffer(p))
	log.Debug("Response: %s", response)

	var failure string
	if err != nil {
		failure = err.Error()
	} else if response.StatusCode != 200 {
		failure = fmt.Sprint(response)
	}

	if failure != "" {
		if h.retries >= h.MaxRetries {
			log.Errorf("Failed to forward payload %d times. Abort mission, I repeat ABORT MISSION!", h.retries)
			log.Error(err)
		} else {
			log.Warnf("Failed to forward payload (%d try). Will try again in %d seconds.\nError: %s", h.retries+1, h.RetryDelay, err)
			h.forwardWithDelay(payload)
		}
	}

}

func (h Http) forwardWithDelay(payload *dto.Stat) {
	h.retries++
	for range time.Tick(time.Duration(h.RetryDelay) * time.Second) {
		h.Send(payload)
	}
}
