package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/depscloud/depscloud/internal/mux"
)

type httpTroubleshootClient struct {
	client  *http.Client
	baseURL string
}

func (s *httpTroubleshootClient) GetServerVersion() (*mux.Version, error) {
	uri := fmt.Sprintf("%s/version", s.baseURL)

	r, err := s.client.Get(uri)
	if err != nil {
		return nil, err
	}
	version := &mux.Version{}
	if err := json.NewDecoder(r.Body).Decode(version); err != nil {
		return nil, err
	}
	return version, nil
}

func (s *httpTroubleshootClient) GetHealth() (string, error) {
	uri := fmt.Sprintf("%s/health", s.baseURL)

	r, err := s.client.Get(uri)
	if err != nil {
		return "", err
	}
	healthString, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(healthString), nil
}
