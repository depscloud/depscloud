package graphstore

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// Statements defines the SQL statements that are used by the GraphStore. Each
// statement should use named parameters.
type Statements struct {
	CreateGraphDataTable                  string `json:"createGraphDataTable"`
	InsertGraphData                       string `json:"insertGraphData"`
	DeleteGraphData                       string `json:"deleteGraphData"`
	ListGraphData                         string `json:"listGraphData"`
	SelectGraphDataUpstreamDependencies   string `json:"selectGraphDataUpstreamDependencies"`
	SelectGraphDataDownstreamDependencies string `json:"selectGraphDataDownstreamDependencies"`
}

// sqlStatements compatible to Sqlite and MySQL
const sqlStatements = `
createGraphDataTable: |
  CREATE TABLE IF NOT EXISTS dts_graphdata(
      graph_item_type VARCHAR(55),
      k1 CHAR(64),
      k2 CHAR(64),
      k3 CHAR(64),
      encoding TINYINT,
      graph_item_data TEXT,
      last_modified DATETIME,
      date_deleted DATETIME DEFAULT NULL,
      PRIMARY KEY (graph_item_type, k1, k2, k3)
  );

insertGraphData: |
  REPLACE INTO dts_graphdata 
  (graph_item_type, k1, k2, k3, encoding, graph_item_data, last_modified, date_deleted)
  VALUES (:graph_item_type, :k1, :k2, :k3, :encoding, :graph_item_data, :last_modified, NULL);

deleteGraphData: |
  UPDATE dts_graphdata
  SET date_deleted = :date_deleted
  WHERE (graph_item_type = :graph_item_type and k1 = :k1 and k2 = :k2);

listGraphData: |
  SELECT graph_item_type, k1, k2, encoding, graph_item_data
  FROM dts_graphdata
  WHERE graph_item_type = :graph_item_type 
  LIMIT :limit OFFSET :offset;

selectGraphDataUpstreamDependencies: |
  SELECT g1.graph_item_type, g1.k1, g1.k2, g1.encoding, g1.graph_item_data,
          g2.graph_item_type, g2.k1, g2.k2, g2.k3, g2.encoding, g2.graph_item_data
  FROM dts_graphdata AS g1
  INNER JOIN dts_graphdata AS g2 ON g1.k1 = g2.k2
  WHERE g2.k1 = :key 
  AND g2.graph_item_type IN (:edge_types) 
  AND g2.k1 != g2.k2 
  AND g2.date_deleted IS NULL
  AND g1.k1 = g1.k2 
  AND g1.date_deleted IS NULL;

selectGraphDataDownstreamDependencies: |
  SELECT g1.graph_item_type, g1.k1, g1.k2, g1.encoding, g1.graph_item_data,
          g2.graph_item_type, g2.k1, g2.k2, g2.k3, g2.encoding, g2.graph_item_data
  FROM dts_graphdata AS g1
  INNER JOIN dts_graphdata AS g2 ON g1.k2 = g2.k1
  WHERE g2.k2 = :key 
  AND g2.graph_item_type IN (:edge_types) 
  AND g2.k1 != g2.k2 
  AND g2.date_deleted IS NULL
  AND g1.k1 = g1.k2 
  AND g1.date_deleted IS NULL;
`

// sqlStatements for PostgreSQL
const postgresStatements = `
createGraphDataTable: |
  CREATE TABLE IF NOT EXISTS dts_graphdata(
      graph_item_type varchar(55),
      k1 char(64),
      k2 char(64),
      k3 char(64),
      encoding smallint,
      graph_item_data text,
      last_modified timestamp,
      date_deleted timestamp DEFAULT NULL,
      PRIMARY KEY (graph_item_type, k1, k2, k3)
  ); 

insertGraphData: |
  INSERT INTO dts_graphdata 
  (graph_item_type, k1, k2, k3, encoding, graph_item_data, last_modified)
  VALUES (:graph_item_type, :k1, :k2, :k3, :encoding, :graph_item_data, :last_modified)
  ON CONFLICT (graph_item_type, k1, k2, k3) 
  DO UPDATE SET graph_item_data = EXCLUDED.graph_item_data, 
                encoding = EXCLUDED.encoding, 
                last_modified = EXCLUDED.last_modified

deleteGraphData: |
  UPDATE dts_graphdata
  SET date_deleted = :date_deleted
  WHERE (graph_item_type = :graph_item_type and k1 = :k1 and k2 = :k2);

listGraphData: |
  SELECT graph_item_type, k1, k2, encoding, graph_item_data
  FROM dts_graphdata
  WHERE graph_item_type = :graph_item_type 
  LIMIT :limit OFFSET :offset;

selectGraphDataUpstreamDependencies: |
  SELECT g1.graph_item_type, g1.k1, g1.k2, g1.encoding, g1.graph_item_data,
          g2.graph_item_type, g2.k1, g2.k2, g2.k3, g2.encoding, g2.graph_item_data
  FROM dts_graphdata AS g1
  INNER JOIN dts_graphdata AS g2 ON g1.k1 = g2.k2
  WHERE g2.k1 = :key 
  AND g2.graph_item_type IN (:edge_types) 
  AND g2.k1 != g2.k2 
  AND g2.date_deleted IS NULL
  AND g1.k1 = g1.k2 
  AND g1.date_deleted IS NULL;

selectGraphDataDownstreamDependencies: |
  SELECT g1.graph_item_type, g1.k1, g1.k2, g1.encoding, g1.graph_item_data,
          g2.graph_item_type, g2.k1, g2.k2, g2.k3, g2.encoding, g2.graph_item_data
  FROM dts_graphdata AS g1
  INNER JOIN dts_graphdata AS g2 ON g1.k2 = g2.k1
  WHERE g2.k2 = :key 
  AND g2.graph_item_type IN (:edge_types) 
  AND g2.k1 != g2.k2 
  AND g2.date_deleted IS NULL
  AND g1.k1 = g1.k2 
  AND g1.date_deleted IS NULL;
`

// LoadStatementsFile loads an external yaml file containing SQL statements
func LoadStatementsFile(yamlFile string) (*Statements, error) {
	contents, err := ioutil.ReadFile(yamlFile)

	if err != nil {
		return nil, err
	}

	return LoadStatements(contents)
}

// LoadStatements parses contents into their corresponding statements
func LoadStatements(contents []byte) (*Statements, error) {
	statements := &Statements{}

	if err := yaml.Unmarshal(contents, statements); err != nil {
		return nil, err
	}

	return statements, nil
}

// DefaultStatementsFor the given database driver
func DefaultStatementsFor(driver string) (*Statements, error) {
	var rawStatements string
	switch driver {
	case postgres:
		rawStatements = postgresStatements

	case mysql, sqlite:
		rawStatements = sqlStatements
	}

	return LoadStatements([]byte(rawStatements))
}
