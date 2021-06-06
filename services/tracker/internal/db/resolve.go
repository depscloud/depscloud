package db

import (
	"fmt"

	"github.com/depscloud/depscloud/services/tracker/internal/db/core"
	dbmysql "github.com/depscloud/depscloud/services/tracker/internal/db/mysql"
	dbpostgresql "github.com/depscloud/depscloud/services/tracker/internal/db/postgresql"
	dbsqlite "github.com/depscloud/depscloud/services/tracker/internal/db/sqlite"

	"github.com/jmoiron/sqlx"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

const (
	MySQLDriverName      = "mysql"
	SQLiteDriverName     = "sqlite3"
	PostgreSQLDriverName = "postgres"
)

// Resolve takes in connection criteria and returns the appropriate driver.
func Resolve(driver, storageAddress, storageReadOnlyAddress string) (name string, rw *gorm.DB, ro *gorm.DB, err error) {
	var rwDialector gorm.Dialector
	var roDialector gorm.Dialector

	switch driver {
	case "mysql":
		name = MySQLDriverName
		rwDialector = mysql.Open(storageAddress)
		roDialector = mysql.Open(storageReadOnlyAddress)
		break
	case "sqlite", "sqlite3":
		name = SQLiteDriverName
		rwDialector = sqlite.Open(storageAddress)
		roDialector = sqlite.Open(storageReadOnlyAddress)
		break
	case "postgres", "postgresql", "pgx":
		name = PostgreSQLDriverName
		rwDialector = postgres.Open(storageAddress)
		roDialector = postgres.Open(storageReadOnlyAddress)
		break
	default:
		return "", nil, nil, fmt.Errorf("failed to resolve driver: %s", driver)
	}

	if len(storageAddress) > 0 {
		rw, err = gorm.Open(rwDialector, &gorm.Config{})
		if err != nil {
			return "", nil, nil, err
		}
	}

	ro = rw
	if len(storageReadOnlyAddress) > 0 {
		ro, err = gorm.Open(roDialector, &gorm.Config{})
		if err != nil {
			return "", nil, nil, err
		}
	}

	if rw == nil && ro == nil {
		return "", nil, nil, fmt.Errorf("must provide one storage address")
	}

	return name, rw, ro, nil
}

// ToSQLX converts the provided gorm db to a SQLX one.
func ToSQLX(name string, gormDB *gorm.DB) (*sqlx.DB, error) {
	if gormDB == nil {
		return nil, nil
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, name), nil
}

// StatementsFor uses the provided information to lookup statements for the driver.
func StatementsFor(driver, version string) *core.Statements {
	switch driver {
	case MySQLDriverName:
		if version == "v1alpha" {
			return dbmysql.V1Alpha
		} else {
			return dbmysql.V1Beta
		}
	case SQLiteDriverName:
		if version == "v1alpha" {
			return dbsqlite.V1Alpha
		} else {
			return dbsqlite.V1Beta
		}
	case PostgreSQLDriverName:
		if version == "v1alpha" {
			return dbpostgresql.V1Alpha
		} else {
			return dbpostgresql.V1Beta
		}
	}
	return nil
}
