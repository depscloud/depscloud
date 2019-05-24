package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
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

const insertGraphData = `INSERT OR REPLACE INTO dts_graphdata 
(graph_item_type, k1, k2, encoding, graph_item_data, last_modified, date_deleted)
VALUES %s;`

const deleteGraphData = `UPDATE dts_graphdata
SET date_deleted = ?
WHERE %s;`

const selectGraphDataPrimaryKey = `SELECT
graph_item_type, k1, k2, encoding, graph_item_data
FROM dts_graphdata 
WHERE graph_item_type = ? and k1 = ? and k2 = ?;`

const selectGraphDataUpstreamDependencies = `SELECT
graph_item_type, k1, k2, encoding, graph_item_data
FROM dts_graphdata 
WHERE k1 IN (SELECT k2 FROM dts_graphdata WHERE k1 = ? and k1 != k2 and date_deleted is NULL)
AND k1 == k2 and date_deleted is NULL;`

const selectGraphDataDownstreamDependencies = `SELECT
graph_item_type, k1, k2, encoding, graph_item_data
FROM dts_graphdata
WHERE k2 IN (SELECT k1 FROM dts_graphdata WHERE k2 = ? and k1 != k2 and date_deleted is NULL)
AND k1 == k2 and date_deleted is NULL;`

// NewSQLGraphStore constructs a new GraphStore with a sql driven backend. Current
// queries support sqlite3 but should be able to work on mysql as well.
func NewSQLGraphStore(db *sql.DB) (GraphStore, error) {
	if _, err := db.Exec(createGraphDataTable); err != nil {
		return nil, err
	}

	return &sqlGraphStore{
		db: db,
	}, nil
}

var _ GraphStore = &sqlGraphStore{}

type sqlGraphStore struct {
	db *sql.DB
}

func (gs *sqlGraphStore) Put(items []*GraphItem) error {
	timestamp := time.Now()

	blocks := make([]string, 0, len(items))
	args := make([]interface{}, 0, len(items)*6)

	for _, item := range items {
		blocks = append(blocks, "(?, ?, ?, ?, ?, ?, NULL)")
		args = append(args, item.GraphItemType, Base64encode(item.K1), Base64encode(item.K2),
			item.Encoding, string(item.GraphItemData), timestamp)
	}

	statement := fmt.Sprintf(insertGraphData, strings.Join(blocks, ", "))
	_, err := gs.db.Exec(statement, args...)

	return err
}

func (gs *sqlGraphStore) Delete(keys []*PrimaryKey) error {
	timestamp := time.Now()

	blocks := make([]string, 0, len(keys))
	args := make([]interface{}, 0, 1+(len(keys)*3))
	args = append(args, timestamp)

	for _, key := range keys {
		blocks = append(blocks, "(graph_item_type = ? and k1 = ? and k2 = ?)")
		args = append(args, key.GraphItemType, Base64encode(key.K1), Base64encode(key.K2))
	}

	statement := fmt.Sprintf(deleteGraphData, strings.Join(blocks, " OR "))
	_, err := gs.db.Exec(statement, args...)

	return err
}

func (gs *sqlGraphStore) FindByPrimary(key *PrimaryKey) (*GraphItem, error) {
	rows, err := gs.db.Query(selectGraphDataPrimaryKey, key.GraphItemType, Base64encode(key.K1), Base64encode(key.K2))
	if err != nil {
		return nil, err
	}

	res, err := readFullyAndClose(rows)
	if err != nil {
		return nil, err
	}

	if len(res) != 1 {
		return nil, fmt.Errorf("failed to read record using primary key: %v", key)
	}

	return res[0], nil
}

func (gs *sqlGraphStore) FindUpstream(key []byte) ([]*GraphItem, error) {
	rows, err := gs.db.Query(selectGraphDataUpstreamDependencies, Base64encode(key))
	if err != nil {
		return nil, err
	}

	res, err := readFullyAndClose(rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (gs *sqlGraphStore) FindDownstream(key []byte) ([]*GraphItem, error) {
	rows, err := gs.db.Query(selectGraphDataDownstreamDependencies, Base64encode(key))
	if err != nil {
		return nil, err
	}

	res, err := readFullyAndClose(rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func readFullyAndClose(rows *sql.Rows) ([]*GraphItem, error) {
	defer rows.Close()

	results := make([]*GraphItem, 0)

	for rows.Next() {
		var (
			graphItemType string
			k1            string
			k2            string
			encoding      uint8
			graphItemData string
		)

		if err := rows.Scan(&graphItemType, &k1, &k2, &encoding, &graphItemData); err != nil {
			return nil, err
		}

		k1bytes, err := Base64decode(k1)
		if err != nil {
			return nil, err
		}

		k2bytes, err := Base64decode(k2)
		if err != nil {
			return nil, err
		}

		results = append(results, &GraphItem{
			GraphItemType: graphItemType,
			K1:            k1bytes,
			K2:            k2bytes,
			Encoding:      encoding,
			GraphItemData: []byte(graphItemData),
		})
	}

	return results, nil
}
