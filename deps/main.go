package main

import (
	"fmt"
	"log"

	"github.com/depscloud/depscloud/deps/internal/client"
	"github.com/depscloud/depscloud/deps/internal/cmds/completion"
	"github.com/depscloud/depscloud/deps/internal/cmds/debug"
	"github.com/depscloud/depscloud/deps/internal/cmds/get"
	"github.com/depscloud/depscloud/deps/internal/writer"
	"github.com/depscloud/depscloud/internal/v"

	"github.com/spf13/cobra"
)

const long = `
deps provides command line access to information stored in a deps.cloud API.

  # configure for private deployments
  export DEPSCLOUD_BASE_URL="https://api.deps.cloud"

  # list available sources
  deps get sources
  deps get sources -l go -o github.com -m depscloud/api
  deps get sources -l go -n github.com/depscloud/api

  # list modules for a source
  deps get modules --url https://github.com/depscloud/api.git

  # list dependents a module
  deps get dependents -l go -o github.com -m depscloud/api
  deps get dependents -l go -n github.com/depscloud/api

  # list dependencies of a module
  deps get dependencies -l go -o github.com -m depscloud/api
  deps get dependencies -l go -n github.com/depscloud/api
`

// variables set by build using -X ldflag
var version string
var commit string
var date string

func main() {
	version := v.Info{Version: version, Commit: commit, Date: date}

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
