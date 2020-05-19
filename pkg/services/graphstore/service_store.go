package graphstore

import (
	"context"
	"time"

	"github.com/deps-cloud/api"
	"github.com/deps-cloud/api/v1alpha/store"

	"github.com/jmoiron/sqlx"

	"github.com/sirupsen/logrus"
)

// NewSQLGraphStore constructs a new GraphStore with a sql driven backend. Current
// queries support sqlite3 but should be able to work on mysql as well.
func NewSQLGraphStore(rwdb, rodb *sqlx.DB, statements *Statements) (store.GraphStoreServer, error) {
	if rwdb != nil {
		if _, err := rwdb.Exec(statements.CreateGraphDataTable); err != nil {
			return nil, err
		}
	}

	return &graphStore{
		rwdb:       rwdb,
		rodb:       rodb,
		statements: statements,
	}, nil
}

type graphStore struct {
	rwdb       *sqlx.DB
	rodb       *sqlx.DB
	statements *Statements
}

var _ store.GraphStoreServer = &graphStore{}

func (gs *graphStore) Put(ctx context.Context, req *store.PutRequest) (*store.PutResponse, error) {
	if gs.rwdb == nil {
		return nil, api.ErrUnsupported
	}

	if len(req.GetItems()) == 0 {
		return &store.PutResponse{}, nil
	}

	timestamp := time.Now()
	errors := make([]error, 0)

	tx, err := gs.rwdb.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, item := range req.GetItems() {
		_, err := tx.NamedExec(gs.statements.InsertGraphData, map[string]interface{}{
			"graph_item_type": item.GetGraphItemType(),
			"k1":              Base64encode(item.GetK1()),
			"k2":              Base64encode(item.GetK2()),
			"k3":              Base64encode(item.GetK3()),
			"encoding":        item.GetEncoding(),
			"graph_item_data": string(item.GetGraphItemData()),
			"last_modified":   timestamp,
		})

		if err != nil {
			errors = append(errors, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if len(errors) > 0 {
		for _, err := range errors {
			logrus.Errorf(err.Error())
		}
		return nil, api.ErrPartialInsertion
	}

	return &store.PutResponse{}, nil
}

func (gs *graphStore) Delete(ctx context.Context, req *store.DeleteRequest) (*store.DeleteResponse, error) {
	if gs.rwdb == nil {
		return nil, api.ErrUnsupported
	}

	if len(req.GetItems()) == 0 {
		return &store.DeleteResponse{}, nil
	}

	timestamp := time.Now()
	errors := make([]error, 0)

	tx, err := gs.rwdb.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, key := range req.GetItems() {
		_, err := tx.NamedExec(gs.statements.DeleteGraphData, map[string]interface{}{
			"date_deleted":    timestamp,
			"graph_item_type": key.GetGraphItemType(),
			"k1":              Base64encode(key.GetK1()),
			"k2":              Base64encode(key.GetK2()),
		})
		if err != nil {
			errors = append(errors, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if len(errors) > 0 {
		for _, err := range errors {
			logrus.Errorf(err.Error())
		}
		return nil, api.ErrPartialDeletion
	}

	return &store.DeleteResponse{}, nil
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func (gs *graphStore) List(ctx context.Context, req *store.ListRequest) (*store.ListResponse, error) {
	graphItemType := req.GetType()
	page := max(req.GetPage(), 1)

	limit := max(min(req.GetCount(), 100), 10)
	offset := (page - 1) * limit

	rows, err := gs.rodb.NamedQuery(gs.statements.ListGraphData, map[string]interface{}{
		"graph_item_type": graphItemType,
		"limit":           limit,
		"offset":          offset,
	})
	if err != nil {
		return nil, err
	}

	items, err := readGraphItems(rows)
	if err != nil {
		return nil, err
	}

	return &store.ListResponse{
		Items: items,
	}, nil
}

func (gs *graphStore) FindUpstream(ctx context.Context, req *store.FindRequest) (*store.FindResponse, error) {
	query, args, err := sqlx.Named(gs.statements.SelectGraphDataUpstreamDependencies, map[string]interface{}{
		"key":        Base64encode(req.GetKey()),
		"edge_types": req.GetEdgeTypes(),
	})
	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	rows, err := gs.rodb.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	pairs, err := readGraphItemPairs(rows)
	if err != nil {
		return nil, err
	}

	return &store.FindResponse{
		Pairs: pairs,
	}, nil
}

func (gs *graphStore) FindDownstream(ctx context.Context, req *store.FindRequest) (*store.FindResponse, error) {
	query, args, err := sqlx.Named(gs.statements.SelectGraphDataDownstreamDependencies, map[string]interface{}{
		"key":        Base64encode(req.GetKey()),
		"edge_types": req.GetEdgeTypes(),
	})
	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	rows, err := gs.rodb.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	pairs, err := readGraphItemPairs(rows)
	if err != nil {
		return nil, err
	}

	return &store.FindResponse{
		Pairs: pairs,
	}, nil
}

func readGraphItems(rows *sqlx.Rows) ([]*store.GraphItem, error) {
	defer rows.Close()

	results := make([]*store.GraphItem, 0)

	for rows.Next() {
		var (
			t    string
			k1   string
			k2   string
			enc  store.GraphItemEncoding
			data string
		)

		if err := rows.Scan(&t, &k1, &k2, &enc, &data); err != nil {
			return nil, err
		}

		k1Bytes, _ := Base64decode(k1)
		k2Bytes, _ := Base64decode(k2)

		item := &store.GraphItem{
			GraphItemType: t,
			K1:            k1Bytes,
			K2:            k2Bytes,
			Encoding:      enc,
			GraphItemData: []byte(data),
		}

		results = append(results, item)
	}

	return results, nil
}

func readGraphItemPairs(rows *sqlx.Rows) ([]*store.GraphItemPair, error) {
	defer rows.Close()

	results := make([]*store.GraphItemPair, 0)

	for rows.Next() {
		var (
			nodeType string
			nodeK1   string
			nodeK2   string
			nodeEnc  store.GraphItemEncoding
			nodeData string
			edgeType string
			edgeK1   string
			edgeK2   string
			edgeK3   string
			edgeEnc  store.GraphItemEncoding
			edgeData string
		)

		if err := rows.Scan(&nodeType, &nodeK1, &nodeK2, &nodeEnc, &nodeData, &edgeType, &edgeK1, &edgeK2, &edgeK3, &edgeEnc, &edgeData); err != nil {
			return nil, err
		}

		nodeK1Bytes, _ := Base64decode(nodeK1)
		nodeK2Bytes, _ := Base64decode(nodeK2)
		edgeK1Bytes, _ := Base64decode(edgeK1)
		edgeK2Bytes, _ := Base64decode(edgeK2)
		edgeK3Bytes, _ := Base64decode(edgeK3)

		pair := &store.GraphItemPair{
			Edge: &store.GraphItem{
				GraphItemType: edgeType,
				K1:            edgeK1Bytes,
				K2:            edgeK2Bytes,
				K3:            edgeK3Bytes,
				Encoding:      edgeEnc,
				GraphItemData: []byte(edgeData),
			},
			Node: &store.GraphItem{
				GraphItemType: nodeType,
				K1:            nodeK1Bytes,
				K2:            nodeK2Bytes,
				Encoding:      nodeEnc,
				GraphItemData: []byte(nodeData),
			},
		}

		results = append(results, pair)
	}

	return results, nil
}
