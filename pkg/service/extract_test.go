package service_test

import (
	"testing"

	desapi "github.com/deps-cloud/des/api"
	dtsapi "github.com/deps-cloud/dts/api"
	"github.com/deps-cloud/dts/api/v1alpha/store"
	"github.com/deps-cloud/dts/pkg/service"
	"github.com/deps-cloud/dts/pkg/services/graphstore"
	"github.com/deps-cloud/dts/pkg/types"

	"github.com/gogo/protobuf/proto"

	"github.com/stretchr/testify/require"
)

func TestExtractSource(t *testing.T) {
	source := service.ExtractSource(&dtsapi.SourceInformation{
		Url: "git@github.com:deps-cloud/dts.git",
	})

	require.Equal(t, types.SourceType, source.GraphItemType)
	require.Equal(t, "lxgT4MU082iVaZ+2ZV1/QpHpBOtTCF+lSjUuXaxAMjE", graphstore.Base64encode(source.K1))
	require.Equal(t, source.K1, source.K2)
	require.Equal(t, store.GraphItemEncoding_JSON, source.Encoding)
	require.Equal(t, `{"url":"git@github.com:deps-cloud/dts.git"}`, string(source.GraphItemData))
}

func TestExtractManagesModule(t *testing.T) {
	key, err := graphstore.Base64decode("lxgT4MU082iVaZ+2ZV1/QpHpBOtTCF+lSjUuXaxAMjE")
	require.Nil(t, err)

	manages, module := service.ExtractManagesModule(key, &desapi.DependencyManagementFile{
		Language:     proto.String("golang"),
		System:       proto.String("vgo"),
		Organization: proto.String("github.com"),
		Module:       proto.String("deps-cloud/dts"),
		Version:      proto.String("1.0.0"),
	})

	require.Equal(t, types.ManagesType, manages.GraphItemType)
	require.Equal(t, "lxgT4MU082iVaZ+2ZV1/QpHpBOtTCF+lSjUuXaxAMjE", graphstore.Base64encode(manages.K1))
	require.Equal(t, "/nUJf/8RK3/nO0spQcKBYkcKqTAnt29L0op3kvGILTE", graphstore.Base64encode(manages.K2))
	require.Equal(t, store.GraphItemEncoding_JSON, manages.Encoding)
	require.Equal(t, `{"language":"golang","system":"vgo","version":"1.0.0"}`, string(manages.GraphItemData))

	require.Equal(t, types.ModuleType, module.GraphItemType)
	require.Equal(t, "/nUJf/8RK3/nO0spQcKBYkcKqTAnt29L0op3kvGILTE", graphstore.Base64encode(module.K1))
	require.Equal(t, module.K1, module.K2)
	require.Equal(t, store.GraphItemEncoding_JSON, module.Encoding)
	require.Equal(t, `{"language":"golang","organization":"github.com","module":"deps-cloud/dts"}`, string(module.GraphItemData))
}

func TestExtractDependsModule(t *testing.T) {
	key, err := graphstore.Base64decode("/nUJf/8RK3/nO0spQcKBYkcKqTAnt29L0op3kvGILTE")
	require.Nil(t, err)

	depends, module := service.ExtractDependsModule("golang", key, &desapi.Dependency{
		Organization:      proto.String("github.com"),
		Module:            proto.String("deps-cloud/des"),
		VersionConstraint: proto.String("1.0"),
		Scopes:            make([]string, 0),
	})

	require.Equal(t, types.DependsType, depends.GraphItemType)
	require.Equal(t, "/nUJf/8RK3/nO0spQcKBYkcKqTAnt29L0op3kvGILTE", graphstore.Base64encode(depends.K1))
	require.Equal(t, "+Cc+G+AqhS2O82R1scJvAntPbKI+sfg9DIR6oqiaqho", graphstore.Base64encode(depends.K2))
	require.Equal(t, store.GraphItemEncoding_JSON, depends.Encoding)
	require.Equal(t, `{"version_constraint":"1.0"}`, string(depends.GraphItemData))

	require.Equal(t, types.ModuleType, module.GraphItemType)
	require.Equal(t, "+Cc+G+AqhS2O82R1scJvAntPbKI+sfg9DIR6oqiaqho", graphstore.Base64encode(module.K1))
	require.Equal(t, module.K1, module.K2)
	require.Equal(t, store.GraphItemEncoding_JSON, module.Encoding)
	require.Equal(t, `{"language":"golang","organization":"github.com","module":"deps-cloud/des"}`, string(module.GraphItemData))
}

func TestExtractModuleKey(t *testing.T) {
	key := service.ExtractModuleKeyFromRequest(&dtsapi.Request{
		Language:     "go",
		Organization: "github.com",
		Module:       "deps-cloud/dts",
	})

	require.Equal(t, "uSxVZFE+IIxni/bvq8F62lHsD24CcMOK/u3ivEJdess", graphstore.Base64encode(key))
}
