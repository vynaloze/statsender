package sender

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/vynaloze/statsender/dto"
	"github.com/vynaloze/statsender/logger"
	"io/ioutil"
	"net/http"
	"time"
)

// Http sender sends the gathered statistics to the defined remote endpoint
type Http struct {
	Target     string    `hcl:"target"`
	RetryDelay int       `hcl:"retryDelay"`
	MaxRetries int       `hcl:"maxRetries"`
	RootCAs    *[]string `hcl:"rootCAs"`
	retries    int
}

// Send sends the dto.Stat in a JSON format to the defined remote endpoint.
// In case of failure, it will be retired with a delay
func (h Http) Send(payload *dto.Stat) {
	log, _ := logger.New()

	p, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
	}
	data := string(p)
	log.Debugf("Target:%s; Data: %s", h.Target, data)

	response, err := h.client().Post(h.Target, "application/json", bytes.NewBuffer(p))
	log.Debug("Response: %s", response)

	var failure string
	if err != nil {
		failure = err.Error()
	} else if response.StatusCode < 200 || response.StatusCode > 299 {
		failure = fmt.Sprint(response)
	}

	if failure != "" {
		if h.retries >= h.MaxRetries {
			log.Errorf("Failed to forward payload %d times. Abort mission, I repeat ABORT MISSION!", h.retries)
			log.Error(err)
		} else {
			log.Warnf("Failed to forward payload (%d try). Will try again in %d seconds.\nError: %v", h.retries+1, h.RetryDelay, err)
			h.forwardWithDelay(payload)
		}
	}

}

// Test tests the if the connection to the defined remote endpoint can be established
func (h Http) Test() error {
	_, err := h.client().Post(h.Target, "application/json", nil)
	return err
}

func (h Http) client() *http.Client {
	log, _ := logger.New()
	caCertPool := x509.NewCertPool()

	if h.RootCAs != nil {
		for _, certFile := range *h.RootCAs {
			log.Debugf("Reading certificate from %s", certFile)
			cert, err := ioutil.ReadFile(certFile)
			if err != nil {
				log.Fatal("Error reading certificate: ", err)
			}
			caCertPool.AppendCertsFromPEM(cert)
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
}

func (h Http) forwardWithDelay(payload *dto.Stat) {
	h.retries++
	for range time.Tick(time.Duration(h.RetryDelay) * time.Second) {
		h.Send(payload)
	}
}
