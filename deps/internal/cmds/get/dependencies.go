package get

import (
	"strings"

	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/deps/internal/writer"

	"github.com/spf13/cobra"
)

func DependenciesCommand(
	dependencyClient tracker.DependencyServiceClient,
	searchClient tracker.SearchServiceClient,
	writer writer.Writer,
) *cobra.Command {
	req := &tracker.DependencyRequest{}

	cmd := &cobra.Command{
		Use:     "dependencies",
		Aliases: []string{"dependency"},
		Short:   "Get the list of modules the given module depends on",
		Example: strings.Join([]string{
			"deps get dependencies -l go -o github.com -m depscloud/api",
			"deps get dependencies -l go -n github.com/depscloud/api",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDependencyRequest(req); err != nil {
				return err
			}

			ctx := cmd.Context()
			response, err := dependencyClient.ListDependencies(ctx, req)
			if err != nil {
				return err
			}

			for _, dependency := range response.Dependencies {
				_ = writer.Write(dependency)
			}

			return nil
		},
	}

	topologyCmd := topologyCommand(writer, searchClient, func(depRequest *tracker.DependencyRequest) *tracker.SearchRequest {
		return &tracker.SearchRequest{
			DependenciesOf: depRequest,
		}
	})

	cmd.AddCommand(topologyCmd)
	addDependencyRequestFlags(cmd, req)

	return cmd
}
