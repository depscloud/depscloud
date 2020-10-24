package debug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/depscloud/depscloud/deps/internal/client"
	"github.com/depscloud/depscloud/internal/mux"

	"github.com/spf13/cobra"
)

func Command(version mux.Version) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Output information helpful for debugging",
		RunE: func(_ *cobra.Command, args []string) error {
			systemInfo := client.GetSystemInfo()
			debugClient := httpDebugClient{client: http.DefaultClient, baseURL: systemInfo.BaseURL}

			// Printing Client environment variables
			fmt.Println(systemInfo)
			// Printing Client version information
			fmt.Println(fmt.Sprintf("Client Version: %s", version))

			// Printing Server version information
			serverVersion, err := debugClient.GetServerVersion()
			if err != nil {
				fmt.Println(fmt.Sprintf("Error While retrieving server version"))
			} else {
				fmt.Println(fmt.Sprintf("Server Version: %s", serverVersion))
			}
			// Printing Server Health information
			healthString, err := debugClient.GetHealth()
			if err != nil {
				fmt.Println(fmt.Sprintf("Error While retrieving server health"))
			} else {
				fmt.Println(fmt.Sprintf("Server Health: %s", healthString))
			}
			return nil
		},
	}
	return cmd
}

type httpDebugClient struct {
	client  *http.Client
	baseURL string
}

func (s *httpDebugClient) GetServerVersion() (*mux.Version, error) {
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

func (s *httpDebugClient) GetHealth() (string, error) {
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
