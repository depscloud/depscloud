package v1alpha

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseName(t *testing.T) {
	tests := []struct {
		language             string
		name                 string
		expectedOrganization string
		expectedModule       string
	}{
		{"jsonnet", "https://github.com/depscloud/deploy.git", "github.com", "depscloud/deploy"},
		{"go", "github.com/depscloud/api", "github.com", "depscloud/api"},
		{"java", "com.google.guava:guava", "com.google.guava", "guava"},
		{"node", "@depscloud/api", "depscloud", "api"},
		{"php", "symfony/console", "symfony", "console"},
		{"rust", "bytes", "_", "bytes"},
	}

	for _, tc := range tests {
		result := parseName(tc.language, tc.name)

		require.Equal(t, []string{tc.expectedOrganization, tc.expectedModule}, result)
	}
}
