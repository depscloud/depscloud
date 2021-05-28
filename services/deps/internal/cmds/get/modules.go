package get

import (
	"strings"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/services/deps/internal/writer"

	"github.com/spf13/cobra"
)

func ModulesCommand(
	sourceService v1beta.SourceServiceClient,
	moduleService v1beta.ModuleServiceClient,
	writer writer.Writer,
) *cobra.Command {
	source := &v1beta.Source{}

	cmd := &cobra.Command{
		Use:     "modules",
		Aliases: []string{"module", "mods", "mod"},
		Short:   "Get a list of modules from the service",
		Example: strings.Join([]string{
			"deps get modules",
			"deps get modules --url https://github.com/depscloud/api.git",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if source.Url != "" {
				response, err := sourceService.ListModules(ctx, &v1beta.ManagedSource{
					Source: source,
				})
				if err != nil {
					return err
				}

				for _, module := range response.Modules {
					_ = writer.Write(module)
				}

				return nil
			}

			pageToken := ""
			for {
				response, err := moduleService.List(ctx, &v1beta.ListRequest{
					PageToken: pageToken,
				})
				if err != nil {
					return err
				}

				for _, module := range response.Modules {
					_ = writer.Write(module)
				}

				if response.NextPageToken == "" {
					break
				}
				pageToken = response.NextPageToken
			}

			return nil
		},
	}

	addSourceFlags(cmd, source)

	return cmd
}
