package get

import (
	"context"
	"fmt"

	"github.com/deps-cloud/api/v1alpha/tracker"
	"github.com/deps-cloud/cli/internal/writer"

	"github.com/spf13/cobra"
)

func DependenciesCommand(
	dependencyClient tracker.DependencyServiceClient,
	writer writer.Writer,
) *cobra.Command {
	req := &tracker.DependencyRequest{}

	cmd := &cobra.Command{
		Use:     "dependencies",
		Short:   "Get the list of modules the given module depends on",
		Example: "deps get dependencies -l go -o github.com -m deps-cloud/api",
		RunE: func(cmd *cobra.Command, args []string) error {
			if req.Language == "" || req.Organization == "" || req.Module == "" {
				return fmt.Errorf("language, organization, and module must be provided")
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

	topologyCmd := topologyCommand(writer,
		func(req *tracker.DependencyRequest, ctx context.Context) ([]*tracker.Dependency, error) {
			resp, err := dependencyClient.ListDependencies(ctx, req)
			if err != nil {
				return nil, err
			}
			return resp.Dependencies, nil
		})

	cmd.AddCommand(topologyCmd)
	addDependencyRequestFlags(cmd, req)

	return cmd
}
