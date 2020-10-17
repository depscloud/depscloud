package get

import (
	"fmt"
	"strings"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/depscloud/deps/internal/writer"

	"github.com/spf13/cobra"
)

func SourcesCommand(
	sourcesClient tracker.SourceServiceClient,
	modulesClient tracker.ModuleServiceClient,
	writer writer.Writer,
) *cobra.Command {
	module := &schema.Module{}

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

			if module.Language != "" && ((module.Organization != "" && module.Module != "") || module.Name != "") {
				response, err := modulesClient.ListSources(ctx, setModuleFields(module))
				if err != nil {
					return err
				}

				for _, source := range response.Sources {
					_ = writer.Write(source)
				}

				return nil
			} else if module.Language != "" || module.Organization != "" || module.Module != "" && module.Name != "" {
				return fmt.Errorf("language + name or language + organization + module must be provided")
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
