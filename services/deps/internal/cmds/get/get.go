package get

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
		Use:   "get <resource>",
		Short: "Retrieve information from the graph",
	}

	cmd.AddCommand(DependenciesCommand(client.Traversal(), writer))
	cmd.AddCommand(DependentsCommand(client.Traversal(), writer))
	cmd.AddCommand(LanguagesCommand(client.Languages(), writer))
	cmd.AddCommand(ModulesCommand(client.Sources(), client.Modules(), writer))
	cmd.AddCommand(SourcesCommand(client.Sources(), client.Modules(), writer))

	return cmd
}
