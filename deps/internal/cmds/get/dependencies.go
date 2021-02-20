package get

import (
	"strings"

	"github.com/depscloud/api/v1beta"

	"github.com/depscloud/depscloud/deps/internal/writer"

	"github.com/spf13/cobra"
)

func DependenciesCommand(
	traversalService v1beta.TraversalServiceClient,
	writer writer.Writer,
) *cobra.Command {
	req := &v1beta.Module{}

	cmd := &cobra.Command{
		Use:     "dependencies",
		Aliases: []string{"dependency"},
		Short:   "Get the list of modules the given module depends on",
		Example: strings.Join([]string{
			"deps get dependencies -l go -o github.com -m depscloud/api",
			"deps get dependencies -l go -n github.com/depscloud/api",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			response, err := traversalService.GetDependencies(ctx, &v1beta.Dependency{
				Module: req,
			})
			if err != nil {
				return err
			}

			for _, dependency := range response.Dependencies {
				_ = writer.Write(dependency)
			}

			return nil
		},
	}

	topologyCmd := topologyCommand(writer, traversalService, func(module *v1beta.Module) *v1beta.SearchRequest {
		return &v1beta.SearchRequest{
			DependenciesFor: &v1beta.Dependency{
				Module: module,
			},
		}
	})

	cmd.AddCommand(topologyCmd)
	addModuleFlags(cmd, req)

	return cmd
}
