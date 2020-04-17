package get

import (
	"github.com/deps-cloud/cli/internal/http"
	"github.com/deps-cloud/cli/internal/writer"

	"github.com/spf13/cobra"
)

func Command(
	client http.Client,
	writer writer.Writer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <resource>",
		Short: "Retrieve information from the graph",
	}

	cmd.AddCommand(DependenciesCommand(client.Dependencies(), writer))
	cmd.AddCommand(DependentsCommand(client.Dependencies(), writer))
	cmd.AddCommand(ModulesCommand(client.Modules(), writer))
	cmd.AddCommand(SourcesCommand(client.Sources(), client.Modules(), writer))
	cmd.AddCommand(TopologyCommand(client.Dependencies(), writer))

	return cmd
}
