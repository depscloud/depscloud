package main

import (
	"github.com/deps-cloud/cli/internal/cmds/completion"
	"github.com/deps-cloud/cli/internal/cmds/get"
	"github.com/deps-cloud/cli/internal/http"
	"github.com/deps-cloud/cli/internal/writer"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const long = `
depscloud-cli provides command line access to information stored in a deps.cloud API.

  # configure for private deployments
  export DEPSCLOUD_BASE_URL="https://api.deps.cloud"

  # list available sources
  depscloud-cli get sources

  # list modules for a source
  depscloud-cli get modules --url https://github.com/deps-cloud/api.git

  # list dependents a module
  depscloud-cli get dependents -l go -o github.com -m deps-cloud/api

  # list dependencies of a module
  depscloud-cli get dependencies -l go -o github.com -m deps-cloud/api
`

func main() {
	client := http.DefaultClient()
	writer := writer.Default

	cmd := &cobra.Command{
		Use:  "depscloud-cli",
		Long: long,
	}

	cmd.AddCommand(completion.Command())
	cmd.AddCommand(get.Command(client, writer))

	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
