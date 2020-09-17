package v1beta

const sqliteInsertGraphData = `
REPLACE INTO graph_data
(k1, k2, k3, kind, encoding, data, date_deleted, last_modified)
VALUES (:k1, :k2, :k3, :kind, :encoding, :data, NULL, :last_modified);
`

const sqliteDeleteGraphData = `
UPDATE graph_data
SET date_deleted = :date_deleted
WHERE k1 = :k1
AND k2 = :k2
AND k3 = :k3;
`

const sqliteListGraphData = `
SELECT k1, k2, k3, kind, encoding, data
FROM graph_data
WHERE kind = :kind
AND date_deleted IS NULL
LIMIT :limit OFFSET :offset;
`

const sqliteSelectFromNeighbor = `
SELECT g1.k1, g1.k2, g1.k3, g1.kind, g1.encoding, g1.data,
	   g2.k1, g2.k2, g2.k3, g2.kind, g2.encoding, g2.data
FROM graph_data AS g1
INNER JOIN graph_data AS g2 ON g1.k1 = g2.k2
WHERE g2.k1 IN (:keys)
AND g2.k1 != g2.k2
AND g2.date_deleted IS NULL
AND g1.k1 = g1.k2
AND g1.date_deleted IS NULL;
`

const sqliteSelectToNeighbor = `
SELECT g1.k1, g1.k2, g1.k3, g1.kind, g1.encoding, g1.data,
	   g2.k1, g2.k2, g2.k3, g2.kind, g2.encoding, g2.data
FROM graph_data AS g1
INNER JOIN graph_data AS g2 ON g1.k2 = g2.k1
WHERE g2.k2 IN (:keys)
AND g2.k1 != g2.k2
AND g2.date_deleted IS NULL
AND g1.k1 = g1.k2
AND g1.date_deleted IS NULL;
`

// SQLiteStatements expose statements that are specific to the SQLite backend
var SQLiteStatements = &Statements{
	InsertGraphData:    sqliteInsertGraphData,
	DeleteGraphData:    sqliteDeleteGraphData,
	ListGraphData:      sqliteListGraphData,
	SelectToNeighbor:   sqliteSelectToNeighbor,
	SelectFromNeighbor: sqliteSelectFromNeighbor,
}
