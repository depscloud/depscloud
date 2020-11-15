package v1alpha

import (
	"net/url"
	"strings"

	"github.com/depscloud/api/v1alpha/schema"
)

func setModuleFields(module *schema.Module) *schema.Module {
	if module.Name == "" {
		return module
	}

	orgAndModule := parseName(module.Language, module.Name)
	module.Organization = orgAndModule[0]
	module.Module = orgAndModule[1]
	return module
}

func parseName(language string, name string) []string {
	var split []string

	switch language {
	case "jsonnet":
		u, err := url.Parse(name)
		if err != nil {
			return []string{"", ""}
		}

		module := strings.TrimSuffix(u.Path, ".git")
		module = strings.TrimPrefix(module, "/")

		return []string{u.Host, module}
	case "java":
		split = strings.Split(name, ":")
		break
	case "node", "js":
		name = strings.TrimPrefix(name, "@")
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
