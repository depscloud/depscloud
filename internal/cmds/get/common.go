package get

import (
	"github.com/deps-cloud/api/v1alpha/schema"
	"github.com/deps-cloud/api/v1alpha/tracker"

	"github.com/spf13/cobra"
)

func addDependencyRequestFlags(cmd *cobra.Command, req *tracker.DependencyRequest) {
	flags := cmd.Flags()

	flags.StringVarP(&(req.Language), "language", "l", req.Language, "The language of the module")
	flags.StringVarP(&(req.Organization), "organization", "o", req.Organization, "The organization of the module")
	flags.StringVarP(&(req.Module), "module", "m", req.Module, "The name of the module")
}

func addSourceFlags(cmd *cobra.Command, source *schema.Source) {
	flags := cmd.Flags()

	flags.StringVarP(&(source.Url), "url", "u", source.Url, "The url to the source repository")
}

func addModuleFlags(cmd *cobra.Command, module *schema.Module) {
	flags := cmd.Flags()

	flags.StringVarP(&(module.Language), "language", "l", module.Language, "The language of the module")
	flags.StringVarP(&(module.Organization), "organization", "o", module.Organization, "The organization of the module")
	flags.StringVarP(&(module.Module), "module", "m", module.Module, "The name of the module")
}
