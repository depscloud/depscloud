package v1beta

import (
	"context"

	"github.com/depscloud/api"

	"github.com/jmoiron/sqlx"
)

// NewSQLDriver creates a new driver using sql semantics
func NewSQLDriver(rwdb, rodb *sqlx.DB, statements *Statements) Driver {
	return &sqlDriver{
		rwdb:       rwdb,
		rodb:       rodb,
		statements: statements,
	}
}

type sqlDriver struct {
	rwdb       *sqlx.DB
	rodb       *sqlx.DB
	statements *Statements
}

func (s *sqlDriver) Put(ctx context.Context, items []*GraphData) error {
	if s.rwdb == nil {
		return api.ErrUnsupported
	}

	tx, err := s.rwdb.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		_, err := tx.NamedExec(s.statements.InsertGraphData, item)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *sqlDriver) Delete(ctx context.Context, items []*GraphData) error {
	if s.rwdb == nil {
		return api.ErrUnsupported
	}

	tx, err := s.rwdb.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		_, err := tx.NamedExec(s.statements.DeleteGraphData, item)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *sqlDriver) List(ctx context.Context, kind string, offset, limit int) ([]*GraphData, bool, error) {
	rows, err := s.rodb.NamedQueryContext(ctx, s.statements.ListGraphData, map[string]interface{}{
		"kind":   kind,
		"offset": offset,
		"limit":  limit + 1,
	})
	if err != nil {
		return nil, false, err
	}

	results := make([]*GraphData, 0, limit)
	for i := 0; i < limit && rows.Next(); i++ {
		item := &GraphData{}
		if err := rows.StructScan(item); err != nil {
			return nil, false, err
		}

		results = append(results, item)
	}

	return results, rows.Next(), nil
}

func (s *sqlDriver) ToNeighbors(ctx context.Context, node *GraphData) ([]*GraphData, error) {
	panic("implement me")
}

func (s *sqlDriver) FromNeighbors(ctx context.Context, node *GraphData) ([]*GraphData, error) {
	panic("implement me")
}

var _ Driver = &sqlDriver{}
