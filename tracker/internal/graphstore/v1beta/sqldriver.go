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

func (s *sqlDriver) neighbors(ctx context.Context, key, statement string) ([]*GraphData, error) {
	rows, err := s.rodb.NamedQueryContext(ctx, statement, map[string]interface{}{
		"keys": key,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*GraphData, 0)
	for rows.Next() {
		var (
			g1K1       string
			g1K2       string
			g1K3       string
			g1Kind     string
			g1Encoding Encoding
			g1Data     string

			g2K1       string
			g2K2       string
			g2K3       string
			g2Kind     string
			g2Encoding Encoding
			g2Data     string
		)

		err := rows.Scan(&g1K1, &g1K2, &g1K3, &g1Kind, &g1Encoding, &g1Data,
			&g2K1, &g2K2, &g2K3, &g2Kind, &g2Encoding, &g2Data)

		if err != nil {
			return nil, err
		}

		g1 := &GraphData{
			K1:       g1K1,
			K2:       g1K2,
			K3:       g1K3,
			Kind:     g1Kind,
			Encoding: g1Encoding,
			Data:     g1Data,
		}

		g2 := &GraphData{
			K1:       g2K1,
			K2:       g2K2,
			K3:       g2K3,
			Kind:     g2Kind,
			Encoding: g2Encoding,
			Data:     g2Data,
		}

		results = append(results, g1, g2)
	}

	return results, nil
}

func (s *sqlDriver) NeighborsTo(ctx context.Context, to *GraphData) ([]*GraphData, error) {
	return s.neighbors(ctx, to.K1, s.statements.SelectToNeighbor)
}

func (s *sqlDriver) NeighborsFrom(ctx context.Context, from *GraphData) ([]*GraphData, error) {
	return s.neighbors(ctx, from.K1, s.statements.SelectFromNeighbor)
}

var _ Driver = &sqlDriver{}
