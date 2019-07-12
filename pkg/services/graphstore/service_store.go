package graphstore

import (
	"context"
	"time"

	"github.com/deps-cloud/tracker/api"
	"github.com/deps-cloud/tracker/api/v1alpha/store"

	"github.com/jmoiron/sqlx"

	"github.com/sirupsen/logrus"
)

const createGraphDataTable = `CREATE TABLE IF NOT EXISTS dts_graphdata(
 	graph_item_type VARCHAR(55),
	k1 CHAR(64),
	k2 CHAR(64),
	encoding TINYINT,
	graph_item_data TEXT,
	last_modified DATETIME,
	date_deleted DATETIME DEFAULT NULL,
	PRIMARY KEY (graph_item_type, k1, k2)
);`

// TODO: move away from this replace operation since it does a delete and insert
const insertGraphData = `REPLACE INTO dts_graphdata 
(graph_item_type, k1, k2, encoding, graph_item_data, last_modified, date_deleted)
VALUES (:graph_item_type, :k1, :k2, :encoding, :graph_item_data, :last_modified, NULL);`

const deleteGraphData = `UPDATE dts_graphdata
SET date_deleted = :date_deleted
WHERE (graph_item_type = :graph_item_type and k1 = :k1 and k2 = :k2);`

const listGraphData = `SELECT
graph_item_type, k1, k2, encoding, graph_item_data
FROM dts_graphdata
WHERE graph_item_type = :graph_item_type 
LIMIT :limit OFFSET :offset;
`

const selectGraphDataUpstreamDependencies = `SELECT
g1.graph_item_type, g1.k1, g1.k2, g1.encoding, g1.graph_item_data,
g2.graph_item_type, g2.k1, g2.k2, g2.encoding, g2.graph_item_data
FROM dts_graphdata AS g1
INNER JOIN dts_graphdata AS g2 ON g1.k1 = g2.k2
WHERE g2.k1 = :key 
AND g2.graph_item_type IN (:edge_types) 
AND g2.k1 != g2.k2 
AND g2.date_deleted IS NULL
AND g1.k1 = g1.k2 
AND g1.date_deleted IS NULL;`

const selectGraphDataDownstreamDependencies = `SELECT
g1.graph_item_type, g1.k1, g1.k2, g1.encoding, g1.graph_item_data,
g2.graph_item_type, g2.k1, g2.k2, g2.encoding, g2.graph_item_data
FROM dts_graphdata AS g1
INNER JOIN dts_graphdata AS g2 ON g1.k2 = g2.k1
WHERE g2.k2 = :key 
AND g2.graph_item_type IN (:edge_types) 
AND g2.k1 != g2.k2 
AND g2.date_deleted IS NULL
AND g1.k1 = g1.k2 
AND g1.date_deleted IS NULL;`

// NewSQLGraphStore constructs a new GraphStore with a sql driven backend. Current
// queries support sqlite3 but should be able to work on mysql as well.
func NewSQLGraphStore(rwdb, rodb *sqlx.DB) (store.GraphStoreServer, error) {
	if rwdb != nil {
		if _, err := rwdb.Exec(createGraphDataTable); err != nil {
			return nil, err
		}
	}

	return &graphStore{
		rwdb: rwdb,
		rodb: rodb,
	}, nil
}

type graphStore struct {
	rwdb *sqlx.DB
	rodb *sqlx.DB
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

	for _, item := range req.GetItems() {
		_, err := tx.NamedExec(insertGraphData, map[string]interface{}{
			"graph_item_type": item.GetGraphItemType(),
			"k1":              Base64encode(item.GetK1()),
			"k2":              Base64encode(item.GetK2()),
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

	for _, key := range req.GetItems() {
		_, err := tx.NamedExec(deleteGraphData, map[string]interface{}{
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

	rows, err := gs.rodb.NamedQuery(listGraphData, map[string]interface{}{
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
	query, args, err := sqlx.Named(selectGraphDataUpstreamDependencies, map[string]interface{}{
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
	query, args, err := sqlx.Named(selectGraphDataDownstreamDependencies, map[string]interface{}{
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
			edgeEnc  store.GraphItemEncoding
			edgeData string
		)

		if err := rows.Scan(&nodeType, &nodeK1, &nodeK2, &nodeEnc, &nodeData, &edgeType, &edgeK1, &edgeK2, &edgeEnc, &edgeData); err != nil {
			return nil, err
		}

		nodeK1Bytes, _ := Base64decode(nodeK1)
		nodeK2Bytes, _ := Base64decode(nodeK2)
		edgeK1Bytes, _ := Base64decode(edgeK1)
		edgeK2Bytes, _ := Base64decode(edgeK2)

		pair := &store.GraphItemPair{
			Edge: &store.GraphItem{
				GraphItemType: edgeType,
				K1:            edgeK1Bytes,
				K2:            edgeK2Bytes,
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
