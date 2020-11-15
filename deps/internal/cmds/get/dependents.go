package get

import (
	"strings"

	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/deps/internal/writer"

	"github.com/spf13/cobra"
)

func DependentsCommand(
	dependencyClient tracker.DependencyServiceClient,
	searchClient tracker.SearchServiceClient,
	writer writer.Writer,
) *cobra.Command {
	req := &tracker.DependencyRequest{}

	cmd := &cobra.Command{
		Use:     "dependents",
		Aliases: []string{"dependent"},
		Short:   "Get the list of modules that depend on the given module",
		Example: strings.Join([]string{
			"deps get dependents -l go -o github.com -m depscloud/api",
			"deps get dependents -l go -n github.com/depscloud/api",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateDependencyRequest(req); err != nil {
				return err
			}

			ctx := cmd.Context()

			response, err := dependencyClient.ListDependents(ctx, req)
			if err != nil {
				return err
			}

			for _, dependent := range response.Dependents {
				_ = writer.Write(dependent)
			}

			return nil
		},
	}

	topologyCmd := topologyCommand(writer, searchClient, func(depRequest *tracker.DependencyRequest) *tracker.SearchRequest {
		return &tracker.SearchRequest{
			DependentsOf: depRequest,
		}
	})

	cmd.AddCommand(topologyCmd)
	addDependencyRequestFlags(cmd, req)

	return cmd
}
