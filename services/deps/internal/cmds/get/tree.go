package get

import (
	"fmt"
	"io"

	"github.com/depscloud/api/v1beta"
	"github.com/spf13/cobra"

	"github.com/depscloud/depscloud/services/deps/internal/writer"
)

type Dependents struct {
	DependentsOf *v1beta.Dependency   `json:"dependents_of"`
	Dependents   []*v1beta.Dependency `json:"dependents"`
}

type Dependencies struct {
	DependenciesFor *v1beta.Dependency   `json:"dependencies_for"`
	Dependencies    []*v1beta.Dependency `json:"dependencies"`
}

func treeCommand(
	writer writer.Writer,
	searchService v1beta.TraversalServiceClient,
	requestConverter convertRequest,
) *cobra.Command {
	req := &v1beta.Module{}

	cmd := &cobra.Command{
		Use:     "tree",
		Aliases: []string{"subtree"},
		Short:   "Get the associated tree using the provided module as the root",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			call, err := searchService.BreadthFirstSearch(ctx)
			if err != nil {
				return err
			}

			if err := call.Send(requestConverter(req)); err != nil {
				return err
			}

			for resp, err := call.Recv(); true; resp, err = call.Recv() {
				if err == io.EOF {
					break
				} else if err != nil {
					return err
				}

				var item interface{}

				if dependentsOf := resp.GetRequest().GetDependentsOf(); dependentsOf != nil {
					item = &Dependents{
						DependentsOf: dependentsOf,
						Dependents:   resp.GetDependents(),
					}
				} else if dependenciesFor := resp.GetRequest().GetDependenciesFor(); dependenciesFor != nil {
					item = &Dependencies{
						DependenciesFor: dependenciesFor,
						Dependencies:    resp.GetDependencies(),
					}
				} else {
					return fmt.Errorf("failed to determine query")
				}

				_ = writer.Write(item)
			}

			return nil
		},
	}

	addModuleFlags(cmd, req)

	return cmd
}
