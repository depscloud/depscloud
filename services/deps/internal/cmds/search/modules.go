package search

import (
	"fmt"
	"strings"

	"github.com/depscloud/api/v1beta"
	"github.com/spf13/cobra"

	"github.com/depscloud/depscloud/services/deps/internal/writer"
)

func ModulesCommand(
	moduleService v1beta.ModuleServiceClient,
	writer writer.Writer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "modules",
		Aliases: []string{"module", "mods", "mod"},
		Short:   "Search modules the system knows about",
		Example: strings.Join([]string{
			"  deps search modules depscloud/api",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing search string")
			}
			ctx := cmd.Context()

			resp, err := moduleService.Search(ctx, &v1beta.ModulesSearchRequest{
				Like: &v1beta.Module{
					Name: args[0],
				},
			})
			if err != nil {
				return err
			}

			for _, module := range resp.Modules {
				_ = writer.Write(module)
			}

			return nil
		},
	}

	return cmd
}