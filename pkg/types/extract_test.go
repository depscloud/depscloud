package types_test

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	desapi "github.com/deps-cloud/des/api"
	dtsapi "github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/pkg/store"
	"github.com/deps-cloud/dts/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestExtractSource(t *testing.T) {
	source := types.ExtractSource(&dtsapi.SourceInformation{
		Url: "git@github.com:deps-cloud/dts.git",
	})

	require.Equal(t, types.SourceType, source.GraphItemType)
	require.Equal(t, "e71RQQ0WHgOJmesXrO09WYmvqSwO4nlJDkcJ10adwcM", store.Base64encode(source.K1))
	require.Equal(t, source.K1, source.K2)
	require.Equal(t, store.EncodingJSON, source.Encoding)
	require.Equal(t, `{"url":"git@github.com:deps-cloud/dts.git"}`, string(source.GraphItemData))
}

func TestExtractManagesModule(t *testing.T) {
	key, err := store.Base64decode("e71RQQ0WHgOJmesXrO09WYmvqSwO4nlJDkcJ10adwcM")
	require.Nil(t, err)

	manages, module := types.ExtractManagesModule(key, &desapi.DependencyManagementFile{
		Language:     proto.String("golang"),
		System:       proto.String("vgo"),
		Organization: proto.String("github.com"),
		Module:       proto.String("deps-cloud/dts"),
		Version:      proto.String("1.0.0"),
	})

	require.Equal(t, types.ManagesType, manages.GraphItemType)
	require.Equal(t, "e71RQQ0WHgOJmesXrO09WYmvqSwO4nlJDkcJ10adwcM", store.Base64encode(manages.K1))
	require.Equal(t, "o+7tCib28rEm+GSEWbPbZG0Tx/hEhvR6jIwC0BcUQSc", store.Base64encode(manages.K2))
	require.Equal(t, store.EncodingJSON, manages.Encoding)
	require.Equal(t, `{"language":"golang","system":"vgo","version":"1.0.0"}`, string(manages.GraphItemData))

	require.Equal(t, types.ModuleType, module.GraphItemType)
	require.Equal(t, "o+7tCib28rEm+GSEWbPbZG0Tx/hEhvR6jIwC0BcUQSc", store.Base64encode(module.K1))
	require.Equal(t, module.K1, module.K2)
	require.Equal(t, store.EncodingJSON, module.Encoding)
	require.Equal(t, `{"language":"golang","organization":"github.com","module":"deps-cloud/dts"}`, string(module.GraphItemData))
}

func TestExtractDependsModule(t *testing.T) {
	key, err := store.Base64decode("o+7tCib28rEm+GSEWbPbZG0Tx/hEhvR6jIwC0BcUQSc")
	require.Nil(t, err)

	depends, module := types.ExtractDependsModule("golang", key, &desapi.Dependency{
		Organization:      proto.String("github.com"),
		Module:            proto.String("deps-cloud/des"),
		VersionConstraint: proto.String("1.0"),
		Scopes:            make([]string, 0),
	})

	require.Equal(t, types.DependsType, depends.GraphItemType)
	require.Equal(t, "o+7tCib28rEm+GSEWbPbZG0Tx/hEhvR6jIwC0BcUQSc", store.Base64encode(depends.K1))
	require.Equal(t, "S62SeqTwUoB7IovlnGLgNsLzPlcVXkawPhZoppgaKnw", store.Base64encode(depends.K2))
	require.Equal(t, store.EncodingJSON, depends.Encoding)
	require.Equal(t, `{"version_constraint":"1.0","scopes":[]}`, string(depends.GraphItemData))

	require.Equal(t, types.ModuleType, module.GraphItemType)
	require.Equal(t, "S62SeqTwUoB7IovlnGLgNsLzPlcVXkawPhZoppgaKnw", store.Base64encode(module.K1))
	require.Equal(t, module.K1, module.K2)
	require.Equal(t, store.EncodingJSON, module.Encoding)
	require.Equal(t, `{"language":"golang","organization":"github.com","module":"deps-cloud/des"}`, string(module.GraphItemData))
}
