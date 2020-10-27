package get

import (
	"strings"

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

func setModuleFields(module *schema.Module) *schema.Module {
	if module.Name == "" {
		return module
	}

	orgAndModule := parseName(module.Language, module.Name)

	module.Organization = orgAndModule[0]
	module.Module = orgAndModule[1]

	return module
}

func setRequestFields(req *tracker.DependencyRequest) *tracker.DependencyRequest {
	if req.Name == "" {
		return req
	}

	orgAndModule := parseName(req.Language, req.Name)

	req.Organization = orgAndModule[0]
	req.Module = orgAndModule[1]

	return req
}

func parseName(language string, name string) []string {
	var split []string

	switch language {
	case "java":
		split = strings.Split(name, ":")
		break
	case "node":
		name = strings.Replace(name, "@", "", 1)
		split = strings.SplitN(name, "/", 2)
		break
	case "js":
		name = strings.Replace(name, "@", "", 1)
		split = strings.SplitN(name, "/", 2)
		break
	default:
		split = strings.SplitN(name, "/", 2)
		break
	}

	if len(split) == 1 {
		return []string{"_", split[0]}
	} else {
		return split
	}
}
