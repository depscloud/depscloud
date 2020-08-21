package get

import (
	"strings"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/cli/internal/writer"

	"github.com/spf13/cobra"
)

func ModulesCommand(
	modulesClient tracker.ModuleServiceClient,
	writer writer.Writer,
) *cobra.Command {
	source := &schema.Source{}

	cmd := &cobra.Command{
		Use:   "modules",
		Aliases: []string{"module", "mods", "mod"},
		Short: "Get a list of modules from the service",
		Example: strings.Join([]string{
			"deps get modules",
			"deps get modules --url https://github.com/depscloud/api.git",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if source.Url != "" {
				response, err := modulesClient.ListManaged(ctx, source)
				if err != nil {
					return err
				}

				for _, module := range response.Modules {
					_ = writer.Write(module)
				}

				return nil
			}

			pageSize := 100

			for i := 1; true; i++ {
				response, err := modulesClient.List(ctx, &tracker.ListRequest{
					Page:  int32(i),
					Count: int32(pageSize),
				})
				if err != nil {
					return err
				}

				for _, module := range response.Modules {
					_ = writer.Write(module)
				}

				if len(response.Modules) < pageSize {
					break
				}
			}

			return nil
		},
	}

	addSourceFlags(cmd, source)

	return cmd
}
