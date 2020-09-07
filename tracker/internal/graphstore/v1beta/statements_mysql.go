package v1beta

const mysqlInsertGraphData = `
REPLACE INTO graph_data
(k1, k2, k3, kind, encoding, data, date_deleted, last_modified)
VALUES (:k1, :k2, :k3, :kind, :encoding, :data, NULL, :last_modified);
INSERT INTO dts_graphdata 
(graph_item_type, k1, k2, k3, encoding, graph_item_data, last_modified, date_deleted)
VALUES (:graph_item_type, :k1, :k2, :k3, :encoding, :graph_item_data, :last_modified, NULL)
ON DUPLICATE KEY UPDATE
encoding = :encoding,
graph_item_data = :graph_item_data, 
last_modified = :last_modified,
date_deleted = NULL;
`

// MySQLStatements expose statements that are specific to the mysql backend
var MySQLStatements = &Statements{
	InsertGraphData: mysqlInsertGraphData,

	// everything else is fine, no modifications required
	DeleteGraphData:     sqliteDeleteGraphData,
	ListGraphData:       sqliteListGraphData,
	SelectToNeighbors:   sqliteSelectToNeighbors,
	SelectFromNeighbors: sqliteSelectFromNeighbors,
}
