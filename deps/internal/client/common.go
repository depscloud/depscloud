package client

import (
	"fmt"
	"net/url"
	"strings"
)

func getUri(format string, baseURL string, language string, organization string, module string, name string) string {
	if name != "" {
		var split []string

		if strings.EqualFold(language, "java") {
			split = strings.Split(name, ":")
		} else {
			split = strings.SplitN(name, "/", 2)
		}

		if len(split) == 1 {
			organization = "_"
			module = split[0]
		} else {
			organization = split[0]
			module = split[1]
		}
	}

	return fmt.Sprintf(format,
		baseURL,
		url.QueryEscape(language),
		url.QueryEscape(organization),
		url.QueryEscape(module))
}
