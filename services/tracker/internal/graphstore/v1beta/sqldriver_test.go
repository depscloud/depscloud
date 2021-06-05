package v1beta_test

import (
	"os"
	"testing"

	"github.com/depscloud/depscloud/services/tracker/internal/db"
	"github.com/depscloud/depscloud/services/tracker/internal/graphstore/v1beta"

	"github.com/stretchr/testify/require"
)

func TestSQLDriver(t *testing.T) {
	storageAddress := "sqldriver_test.db?cache=shared"
	storageReadOnlyAddress := "sqldriver_test.db?cache=shared&mode=ro"

	defer os.Remove("sqldriver_test.db")

	name, rw, ro, err := db.Resolve("sqlite", storageAddress, storageReadOnlyAddress)
	require.Nil(t, err)

	err = rw.AutoMigrate(&v1beta.GraphData{})
	require.Nil(t, err)

	sqlxRW, err := db.ToSQLX(name, rw)
	require.Nil(t, err)

	sqlxRO, err := db.ToSQLX(name, ro)
	require.Nil(t, err)

	statements := db.StatementsFor(name, "v1beta")

	driver := v1beta.NewSQLDriver(sqlxRW, sqlxRO, statements)
	testServer(t, driver)
}
