package v1beta_test

import (
	"testing"

	"github.com/depscloud/depscloud/tracker/internal/graphstore/v1beta"

	"github.com/stretchr/testify/require"
)

func TestSQLDriver(t *testing.T) {
	driver, err := v1beta.Resolve("sqlite", "file::memory:?cache=shared", "")
	require.Nil(t, err)
	testServer(t, driver)
}
