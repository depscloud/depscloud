package get

import (
	"fmt"
	"strings"

	"github.com/deps-cloud/api/v1alpha/schema"
	"github.com/deps-cloud/api/v1alpha/tracker"
	"github.com/deps-cloud/cli/internal/writer"

	"github.com/spf13/cobra"
)

func SourcesCommand(
	sourcesClient tracker.SourceServiceClient,
	modulesClient tracker.ModuleServiceClient,
	writer writer.Writer,
) *cobra.Command {
	module := &schema.Module{}

	cmd := &cobra.Command{
		Use:   "sources",
		Short: "Get a list of source repositories from the service",
		Example: strings.Join([]string{
			"depscloud-cli get sources",
			"depscloud-cli get sources -l go -o github.com -m deps-cloud/api",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if module.Language != "" && module.Organization != "" && module.Module != "" {
				response, err := modulesClient.ListSources(ctx, module)
				if err != nil {
					return err
				}

				for _, source := range response.Sources {
					_ = writer.Write(source)
				}

				return nil
			} else if module.Language != "" || module.Organization != "" || module.Module != "" {
				return fmt.Errorf("language, organization, and module must be provided")
			}

			pageSize := 100

			for i := 1; true; i++ {
				response, err := sourcesClient.List(ctx, &tracker.ListRequest{
					Page:  int32(i),
					Count: int32(pageSize),
				})
				if err != nil {
					return err
				}

				for _, source := range response.Sources {
					_ = writer.Write(source)
				}

				if len(response.Sources) < pageSize {
					break
				}
			}

			return nil
		},
	}

	addModuleFlags(cmd, module)

	return cmd
}
