package v1beta

const mysqlInsertGraphData = `
INSERT INTO graph_data 
(k1, k2, k3, kind, encoding, data, date_deleted, last_modified)
VALUES (:k1, :k2, :k3, :kind, :encoding, :data, NULL, :last_modified)
ON DUPLICATE KEY UPDATE
	kind = :kind,
	encoding = :encoding,
	graph_item_data = :graph_item_data, 
	date_deleted = NULL,
	last_modified = :last_modified;
`

// MySQLStatements expose statements that are specific to the mysql backend
var MySQLStatements = &Statements{
	InsertGraphData: mysqlInsertGraphData,

	// everything else is fine, no modifications required
	DeleteGraphData:    sqliteDeleteGraphData,
	ListGraphData:      sqliteListGraphData,
	SelectToNeighbor:   sqliteSelectToNeighbor,
	SelectFromNeighbor: sqliteSelectFromNeighbor,
}
