package search

import (
	"github.com/depscloud/depscloud/services/deps/internal/client"
	"github.com/depscloud/depscloud/services/deps/internal/writer"

	"github.com/spf13/cobra"
)

func Command(
	client client.Client,
	writer writer.Writer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <resource>",
		Short: "Search information in the graph",
	}

	cmd.AddCommand(ModulesCommand(client.Modules(), writer))
	cmd.AddCommand(SourcesCommand(client.Sources(), writer))

	return cmd
}
