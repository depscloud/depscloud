package get

import (
	"github.com/depscloud/api/v1beta"

	"github.com/spf13/cobra"
)

func addSourceFlags(cmd *cobra.Command, source *v1beta.Source) {
	flags := cmd.Flags()

	flags.StringVarP(&(source.Url), "url", "u", source.Url, "The url to the source repository")
}

func addModuleFlags(cmd *cobra.Command, module *v1beta.Module) {
	flags := cmd.Flags()

	flags.StringVarP(&(module.Language), "language", "l", module.Language, "The language of the module")
	flags.StringVarP(&(module.Name), "name", "n", module.Name, "The name of the module")
}
