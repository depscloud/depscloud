package debug

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/depscloud/depscloud/internal/appconf"
	"github.com/depscloud/depscloud/services/deps/internal/client"

	"github.com/mjpitz/go-gracefully/check"

	"github.com/spf13/cobra"
)

func Command(version appconf.V) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Output information helpful for debugging",
		RunE: func(_ *cobra.Command, args []string) error {

			systemInfo := client.GetSystemInfo()
			debugClient := httpDebugClient{client: http.DefaultClient, baseURL: systemInfo.BaseURL}
			serverVersion, versionErr := debugClient.GetServerVersion()
			healthString, healthErr := debugClient.GetHealth()

			// Printing Client environment variables
			fmt.Println(fmt.Sprintf("System V: %s", systemInfo))
			// Printing Client version information
			fmt.Println(fmt.Sprintf("Client Version: %s", version))

			// Printing Server version information
			if versionErr != nil {
				fmt.Println(fmt.Sprintf("Error While retrieving server version"), versionErr)
			} else {
				fmt.Println(fmt.Sprintf("Server Version: %s", serverVersion))
			}

			// Printing Server Health information
			if healthErr != nil {
				fmt.Println(fmt.Sprintf("Error While retrieving server health"), healthErr)
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

func (s *httpDebugClient) GetServerVersion() (version appconf.V, err error) {
	uri := fmt.Sprintf("%s/version", s.baseURL)

	r, err := s.client.Get(uri)
	if err != nil {
		return version, err
	}

	if err := json.NewDecoder(r.Body).Decode(&version); err != nil {
		return version, err
	}

	return version, err
}

func (s *httpDebugClient) GetHealth() (string, error) {
	uri := fmt.Sprintf("%s/health", s.baseURL)

	resp, err := s.client.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r := check.Result{}
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	healthString, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	return string(healthString), nil
}
