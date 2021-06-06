package search

import (
	"fmt"
	"strings"

	"github.com/depscloud/api/v1beta"
	"github.com/spf13/cobra"

	"github.com/depscloud/depscloud/services/deps/internal/writer"
)

func SourcesCommand(
	sourceService v1beta.SourceServiceClient,
	writer writer.Writer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sources",
		Aliases: []string{"source", "srcs", "src"},
		Short:   "Search sources the system knows about",
		Example: strings.Join([]string{
			"  deps search sources depscloud/api",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing search string")
			}
			ctx := cmd.Context()

			resp, err := sourceService.Search(ctx, &v1beta.SourcesSearchRequest{
				Like: &v1beta.Source{
					Url: args[0],
				},
			})
			if err != nil {
				return err
			}

			for _, source := range resp.Sources {
				_ = writer.Write(source)
			}

			return nil
		},
	}

	return cmd
}
