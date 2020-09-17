package v1beta

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"

	"github.com/jmoiron/sqlx"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

const (
	mysqlDriverName = "mysql"
	sqliteDriverName = "sqlite3"
	postgresqlDriverName = "postgres"
)

// Driver represents a generic interface for storing a graph.
type Driver interface {
	Put(ctx context.Context, data []*GraphData) error
	Delete(ctx context.Context, data []*GraphData) error
	List(ctx context.Context, kind string, offset, limit int) ([]*GraphData, bool, error)
	NeighborsTo(ctx context.Context, toKeys []string) ([]*GraphData, error)
	NeighborsFrom(ctx context.Context, fromKeys []string) ([]*GraphData, error)
}

// Resolve takes in connection criteria and returns the appropriate driver
func Resolve(driver, storageAddress, storageReadOnlyAddress string) (Driver, error) {
	var statements *Statements
	var dialectorRW gorm.Dialector
	var dialectorRO gorm.Dialector

	switch driver {
	case "mysql":
		driver = mysqlDriverName
		statements = MySQLStatements
		dialectorRW = mysql.Open(storageAddress)
		dialectorRO = mysql.Open(storageReadOnlyAddress)
		break
	case "sqlite", "sqlite3":
		driver = sqliteDriverName
		statements = SQLiteStatements
		dialectorRW = sqlite.Open(storageAddress)
		dialectorRO = sqlite.Open(storageReadOnlyAddress)
		break
	case "postgres", "postgresql", "pgx":
		driver = postgresqlDriverName
		statements = PostgreSQLStatements
		dialectorRW = postgres.Open(storageAddress)
		dialectorRO = postgres.Open(storageReadOnlyAddress)
		break
	default:
		return nil, fmt.Errorf("failed to resolve driver: %s", driver)
	}

	var dbrw *sql.DB
	if len(storageAddress) > 0 {
		gormRW, err := gorm.Open(dialectorRW, &gorm.Config{})
		if err != nil {
			return nil, err
		}

		// only migrate on RW connections
		err = gormRW.AutoMigrate(&GraphData{})
		if err != nil {
			return nil, err
		}

		dbrw, err = gormRW.DB()
		if err != nil {
			return nil, err
		}
	}

	dbro := dbrw
	if len(storageReadOnlyAddress) > 0 {
		gormRO, err := gorm.Open(dialectorRO, &gorm.Config{})
		if err != nil {
			return nil, err
		}

		dbro, err = gormRO.DB()
		if err != nil {
			return nil, err
		}
	}

	if dbrw == nil && dbro == nil {
		return nil, fmt.Errorf("must provide one storage address")
	}

	return &sqlDriver{
		rwdb:       sqlx.NewDb(dbrw, driver),
		rodb:       sqlx.NewDb(dbro, driver),
		statements: statements,
	}, nil
}
