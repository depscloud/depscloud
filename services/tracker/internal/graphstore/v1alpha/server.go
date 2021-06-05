package v1alpha

import (
	"context"
	"time"

	"github.com/depscloud/api"
	"github.com/depscloud/api/v1alpha/store"
	"github.com/depscloud/depscloud/services/tracker/internal/db"
	"github.com/depscloud/depscloud/services/tracker/internal/db/core"

	"github.com/jmoiron/sqlx"
)

const sqliteCreateTable = `
CREATE TABLE IF NOT EXISTS dts_graphdata(
  graph_item_type VARCHAR(55),
  k1 CHAR(64),
  k2 CHAR(64),
  k3 VARCHAR(64),
  encoding TINYINT,
  graph_item_data TEXT,
  last_modified DATETIME,
  date_deleted DATETIME DEFAULT NULL,
  PRIMARY KEY (graph_item_type, k1, k2, k3)
);
CREATE INDEX IF NOT EXISTS secondary ON dts_graphdata(graph_item_type, k2, k1, k3);
CREATE INDEX IF NOT EXISTS date_deleted ON dts_graphdata(date_deleted);
`

const mysqlCreateTable = `
CREATE TABLE IF NOT EXISTS dts_graphdata(
  graph_item_type VARCHAR(55),
  k1 CHAR(64),
  k2 CHAR(64),
  k3 VARCHAR(64),
  encoding TINYINT,
  graph_item_data TEXT,
  last_modified DATETIME,
  date_deleted DATETIME DEFAULT NULL,
  PRIMARY KEY (graph_item_type, k1, k2, k3),
  KEY secondary (graph_item_type, k2, k1, k3),
  KEY (date_deleted)
);
`

const postgresqlCreateTable = `
CREATE TABLE IF NOT EXISTS dts_graphdata(
  graph_item_type VARCHAR(55),
  k1 CHAR(64),
  k2 CHAR(64),
  k3 VARCHAR(64),
  encoding SMALLINT,
  graph_item_data TEXT,
  last_modified TIMESTAMP,
  date_deleted TIMESTAMP DEFAULT NULL,
  PRIMARY KEY (graph_item_type, k1, k2, k3)
);
CREATE INDEX IF NOT EXISTS secondary ON dts_graphdata(graph_item_type, k2, k1, k3);
CREATE INDEX IF NOT EXISTS date_deleted ON dts_graphdata(date_deleted);
`

// NewSQLGraphStore constructs a new GraphStore with a sql driven backend. Current
// queries support sqlite3 but should be able to work on mysql as well.
func NewSQLGraphStore(rwdb, rodb *sqlx.DB, statements *core.Statements) (store.GraphStoreServer, error) {
	if rwdb != nil {
		createTable := sqliteCreateTable

		switch rwdb.DriverName() {
		case db.SQLiteDriverName:
			createTable = sqliteCreateTable
		case db.MySQLDriverName:
			createTable = mysqlCreateTable
		case db.PostgreSQLDriverName:
			createTable = postgresqlCreateTable
		}

		if _, err := rwdb.Exec(createTable); err != nil {
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
	store.UnsafeGraphStoreServer

	rwdb       *sqlx.DB
	rodb       *sqlx.DB
	statements *core.Statements
}

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
			"k3":              Base64encode(key.GetK3()),
		})
		if err != nil {
			errors = append(errors, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if len(errors) > 0 {
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
	keys := make([]string, len(req.GetKeys()))
	for i, key := range req.GetKeys() {
		keys[i] = Base64encode(key)
	}

	query, args, err := sqlx.Named(gs.statements.SelectOutTreeNeighbors, map[string]interface{}{
		"keys":       keys,
		"edge_types": req.GetEdgeTypes(),
		"node_types": req.GetNodeTypes(),
	})
	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	// transform the query to the DB specific bindvar type
	query = gs.rodb.Rebind(query)

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
	keys := make([]string, len(req.GetKeys()))
	for i, key := range req.GetKeys() {
		keys[i] = Base64encode(key)
	}

	query, args, err := sqlx.Named(gs.statements.SelectInTreeNeighbors, map[string]interface{}{
		"keys":       keys,
		"edge_types": req.GetEdgeTypes(),
		"node_types": req.GetNodeTypes(),
	})
	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	// transform the query to the DB specific bindvar type
	query = gs.rodb.Rebind(query)

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

var _ store.GraphStoreServer = &graphStore{}
