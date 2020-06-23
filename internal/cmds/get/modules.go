package get

import (
	"strings"

	"github.com/deps-cloud/api/v1alpha/schema"
	"github.com/deps-cloud/api/v1alpha/tracker"
	"github.com/deps-cloud/cli/internal/writer"

	"github.com/spf13/cobra"
)

func ModulesCommand(
	modulesClient tracker.ModuleServiceClient,
	writer writer.Writer,
) *cobra.Command {
	source := &schema.Source{}

	cmd := &cobra.Command{
		Use:   "modules",
		Short: "Get a list of modules from the service",
		Example: strings.Join([]string{
			"deps get modules",
			"deps get modules --url https://github.com/deps-cloud/api.git",
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
