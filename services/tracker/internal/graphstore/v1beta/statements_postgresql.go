package v1beta

const postgresqlInsertGraphData = `
INSERT INTO graph_data 
(k1, k2, k3, kind, encoding, data, date_deleted, last_modified)
VALUES (:k1, :k2, :k3, :kind, :encoding, :data, NULL, :last_modified)
ON CONFLICT (k1, k2, k3) 
DO UPDATE SET
	kind = EXCLUDED.kind, 
	encoding = EXCLUDED.encoding,
	data = EXCLUDED.data,
	date_deleted = NULL,
	last_modified = EXCLUDED.last_modified
`

// PostgreSQLStatements expose statements that are specific to the PostgreSQL backend
var PostgreSQLStatements = &Statements{
	InsertGraphData: postgresqlInsertGraphData,

	// everything else is fine, no modifications required
	DeleteGraphData:    sqliteDeleteGraphData,
	ListGraphData:      sqliteListGraphData,
	SelectToNeighbor:   sqliteSelectToNeighbor,
	SelectFromNeighbor: sqliteSelectFromNeighbor,
}
