package v1beta_test

import (
	"os"
	"testing"

	"github.com/depscloud/depscloud/services/tracker/internal/graphstore/v1beta"

	"github.com/stretchr/testify/require"
)

func TestSQLDriver(t *testing.T) {
	storageAddress := "sqldriver_test.db?cache=shared"
	storageReadOnlyAddress := "sqldriver_test.db?cache=shared&mode=ro"

	defer os.Remove("sqldriver_test.db")

	driver, err := v1beta.Resolve("sqlite", storageAddress, storageReadOnlyAddress)
	require.Nil(t, err)

	testServer(t, driver)
}
