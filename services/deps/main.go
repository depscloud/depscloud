package main

import (
	"fmt"
	"log"

	"github.com/depscloud/depscloud/internal/appconf"
	"github.com/depscloud/depscloud/services/deps/internal/client"
	"github.com/depscloud/depscloud/services/deps/internal/cmds/completion"
	"github.com/depscloud/depscloud/services/deps/internal/cmds/debug"
	"github.com/depscloud/depscloud/services/deps/internal/cmds/get"
	"github.com/depscloud/depscloud/services/deps/internal/writer"

	"github.com/spf13/cobra"

	_ "google.golang.org/grpc/health"
)

const long = `
deps provides command line access to information stored in a deps.cloud API.

  # configure for private deployments
  export DEPSCLOUD_BASE_URL="https://api.deps.cloud"

  # list available sources
  deps get sources
  deps get sources -l go -n github.com/depscloud/api

  # list modules for a source
  deps get modules --url https://github.com/depscloud/api.git

  # list dependents of a module
  deps get dependents -l go -n github.com/depscloud/api

  # list dependencies for a module
  deps get dependencies -l go -n github.com/depscloud/api
`

func main() {
	version := appconf.Current()

	c := client.DefaultClient()
	w := writer.Default

	cmd := &cobra.Command{
		Use:  "deps",
		Long: long,
	}

	cmd.AddCommand(completion.Command())
	cmd.AddCommand(get.Command(c, w))

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Output version information",
		RunE: func(_ *cobra.Command, args []string) error {
			versionString := fmt.Sprintf("%s %s", cmd.Use, version)
			fmt.Println(versionString)
			return nil
		},
	})

	cmd.AddCommand(debug.Command(version))

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
