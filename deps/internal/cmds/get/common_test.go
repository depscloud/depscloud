package get

import (
	"testing"

	"github.com/depscloud/api/v1alpha/schema"
	"github.com/depscloud/api/v1alpha/tracker"

	"github.com/stretchr/testify/require"
)

func Test_validateDependencyRequest(t *testing.T) {
	err := validateDependencyRequest(&tracker.DependencyRequest{})
	require.Error(t, err)

	err = validateDependencyRequest(&tracker.DependencyRequest{Language: "go"})
	require.Error(t, err)

	err = validateDependencyRequest(&tracker.DependencyRequest{Language: "go", Name: "testing"})
	require.NoError(t, err)

	err = validateDependencyRequest(&tracker.DependencyRequest{Language: "go", Module: "testing"})
	require.Error(t, err)

	err = validateDependencyRequest(&tracker.DependencyRequest{Language: "go", Organization: "_", Module: "testing"})
	require.NoError(t, err)
}

func Test_isEmpty(t *testing.T) {
	require.True(t, isEmpty(&schema.Module{}))

	require.False(t, isEmpty(&schema.Module{
		Language: "test",
	}))

	require.False(t, isEmpty(&schema.Module{
		Name: "test",
	}))

	require.False(t, isEmpty(&schema.Module{
		Organization: "test",
	}))

	require.False(t, isEmpty(&schema.Module{
		Module: "test",
	}))
}

func Test_validateModule(t *testing.T) {
	err := validateModule(&schema.Module{})
	require.Error(t, err)

	err = validateModule(&schema.Module{Language: "go"})
	require.Error(t, err)

	err = validateModule(&schema.Module{Language: "go", Name: "testing"})
	require.NoError(t, err)

	err = validateModule(&schema.Module{Language: "go", Module: "testing"})
	require.Error(t, err)

	err = validateModule(&schema.Module{Language: "go", Organization: "_", Module: "testing"})
	require.NoError(t, err)
}
