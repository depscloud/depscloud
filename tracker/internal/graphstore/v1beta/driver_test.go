package v1beta_test

import (
	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1beta"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Resolve_dne(t *testing.T) {
	_, err := v1beta.Resolve("dne", "", "")
	require.Error(t, err)
	require.Equal(t, "failed to resolve driver: dne", err.Error())
}

func Test_Resolve_sqlite_noaddress(t *testing.T) {
	_, err := v1beta.Resolve("sqlite", "", "")
	require.Error(t, err)
	require.Equal(t, "must provide one storage address", err.Error())
}

func Test_Resolve_sqlite_readwrite(t *testing.T) {
	_, err := v1beta.Resolve("sqlite", "file::memory:?cache=shared", "")
	require.Nil(t, err)
}

func Test_Resolve_sqlite_readonly(t *testing.T) {
	_, err := v1beta.Resolve("sqlite", "", "file::memory:?cache=shared&mode=ro")
	require.Nil(t, err)
}
