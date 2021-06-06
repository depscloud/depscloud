package get

import (
	"strings"

	"github.com/depscloud/api/v1beta"
	"github.com/spf13/cobra"

	"github.com/depscloud/depscloud/services/deps/internal/writer"
)

func LanguagesCommand(
	languageService v1beta.LanguageServiceClient,
	writer writer.Writer,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "languages",
		Aliases: []string{"language", "langs", "lang"},
		Short:   "Get a list of languages from the service",
		Example: strings.Join([]string{
			"deps get languages",
		}, "\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			resp, err := languageService.List(ctx, &v1beta.ListRequest{})
			if err != nil {
				return err
			}

			for _, language := range resp.Languages {
				_ = writer.Write(language)
			}

			return nil
		},
	}

	return cmd
}