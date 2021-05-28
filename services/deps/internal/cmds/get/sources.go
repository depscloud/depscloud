package get

import (
	"strings"

	"github.com/depscloud/api/v1beta"
	"github.com/depscloud/depscloud/services/deps/internal/writer"

	"github.com/spf13/cobra"
)

func SourcesCommand(
	sourceService v1beta.SourceServiceClient,
	moduleService v1beta.ModuleServiceClient,
	writer writer.Writer,
) *cobra.Command {
	module := &v1beta.Module{}

	cmd := &cobra.Command{
		Use:     "sources",
		Aliases: []string{"source", "srcs", "src"},
		Short:   "Get a list of source repositories from the service",
		Example: strings.Join([]string{
			"deps get sources",
			"deps get sources -l go -o github.com -m depscloud/api",
			"deps get sources -l go -n github.com/depscloud/api",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if module.Language != "" && module.Name != "" {
				// list sources for module

				response, err := moduleService.ListSources(ctx, &v1beta.ManagedModule{
					Module: module,
				})
				if err != nil {
					return err
				}

				for _, source := range response.Sources {
					_ = writer.Write(source)
				}

				return nil
			}

			pageToken := ""
			for {
				response, err := sourceService.List(ctx, &v1beta.ListRequest{
					PageToken: pageToken,
				})
				if err != nil {
					return err
				}

				for _, source := range response.Sources {
					_ = writer.Write(source)
				}

				if response.NextPageToken == "" {
					break
				}
				pageToken = response.NextPageToken
			}

			return nil
		},
	}

	addModuleFlags(cmd, module)

	return cmd
}
