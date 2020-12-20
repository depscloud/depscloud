package v1beta

import (
	"context"

	"github.com/depscloud/api"

	"github.com/jmoiron/sqlx"
)

type sqlDriver struct {
	rwdb       *sqlx.DB
	rodb       *sqlx.DB
	statements *Statements
}

func (s *sqlDriver) Put(ctx context.Context, items []*GraphData) error {
	if s.rwdb == nil {
		return api.ErrUnsupported
	}

	tx, err := s.rwdb.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		_, err := tx.NamedExecContext(ctx, s.statements.InsertGraphData, item)
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

	tx, err := s.rwdb.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		_, err := tx.NamedExecContext(ctx, s.statements.DeleteGraphData, item)
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
	defer rows.Close()

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

func (s *sqlDriver) neighbors(ctx context.Context, statement string, keys []string) ([]*GraphData, error) {
	query, args, err := sqlx.Named(statement, map[string]interface{}{
		"keys": keys,
	})
	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	// use DB specific bindvar type
	query = s.rodb.Rebind(query)

	rows, err := s.rodb.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*GraphData, 0)
	for rows.Next() {
		item := &GraphData{}
		if err := rows.StructScan(item); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, nil
}

func (s *sqlDriver) NeighborsTo(ctx context.Context, toKeys []string) ([]*GraphData, error) {
	return s.neighbors(ctx, s.statements.SelectToNeighbor, toKeys)
}

func (s *sqlDriver) NeighborsFrom(ctx context.Context, fromKeys []string) ([]*GraphData, error) {
	return s.neighbors(ctx, s.statements.SelectFromNeighbor, fromKeys)
}

var _ Driver = &sqlDriver{}
