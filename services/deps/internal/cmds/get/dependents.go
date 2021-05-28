package get

import (
	"strings"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/services/deps/internal/writer"

	"github.com/spf13/cobra"
)

func DependentsCommand(
	traversalService v1beta.TraversalServiceClient,
	writer writer.Writer,
) *cobra.Command {
	req := &v1beta.Module{}

	cmd := &cobra.Command{
		Use:     "dependents",
		Aliases: []string{"dependent"},
		Short:   "Get the list of modules that depend on the given module",
		Example: strings.Join([]string{
			"deps get dependents -l go -o github.com -m depscloud/api",
			"deps get dependents -l go -n github.com/depscloud/api",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			response, err := traversalService.GetDependents(ctx, &v1beta.Dependency{
				Module: req,
			})
			if err != nil {
				return err
			}

			for _, dependent := range response.Dependents {
				_ = writer.Write(dependent)
			}

			return nil
		},
	}

	topologyCmd := topologyCommand(writer, traversalService, func(module *v1beta.Module) *v1beta.SearchRequest {
		return &v1beta.SearchRequest{
			DependentsOf: &v1beta.Dependency{
				Module: module,
			},
		}
	})

	cmd.AddCommand(topologyCmd)
	addModuleFlags(cmd, req)

	return cmd
}
