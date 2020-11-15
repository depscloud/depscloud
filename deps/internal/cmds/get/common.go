package get

import (
	"fmt"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/spf13/cobra"
)

func addDependencyRequestFlags(cmd *cobra.Command, req *tracker.DependencyRequest) {
	flags := cmd.Flags()

	flags.StringVarP(&(req.Language), "language", "l", req.Language, "The language of the module")
	flags.StringVarP(&(req.Organization), "organization", "o", req.Organization, "The organization of the module")
	flags.StringVarP(&(req.Module), "module", "m", req.Module, "The name of the module")
	flags.StringVarP(&(req.Name), "name", "n", req.Name, "The name of the module")
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
	flags.StringVarP(&(module.Name), "name", "n", module.Name, "The name of the module")
}

func validateDependencyRequest(req *tracker.DependencyRequest) error {
	if req.GetLanguage() == "" {
		return fmt.Errorf("a language must be provided")
	} else if req.GetOrganization() == "" && req.GetModule() == "" && req.GetName() == "" {
		return fmt.Errorf("a name must be provided")
	} else if req.GetName() == "" && (req.GetOrganization() == "" || req.GetModule() == "") {
		return fmt.Errorf("both an organization and module must be provided. [deprecated, please use name]")
	}
	return nil
}

func isEmpty(req *schema.Module) bool {
	return req.GetLanguage() == "" &&
		req.GetOrganization() == "" &&
		req.GetModule() == "" &&
		req.GetName() == ""
}

func validateModule(req *schema.Module) error {
	if req.GetLanguage() == "" {
		return fmt.Errorf("a language must be provided")
	} else if req.GetOrganization() == "" && req.GetModule() == "" && req.GetName() == "" {
		return fmt.Errorf("a name must be provided")
	} else if req.GetName() == "" && (req.GetOrganization() == "" || req.GetModule() == "") {
		return fmt.Errorf("both an organization and module must be provided. [deprecated, please use name]")
	}
	return nil
}
