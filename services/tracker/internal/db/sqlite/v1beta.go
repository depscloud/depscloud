package dbsqlite

import (
	"github.com/depscloud/depscloud/services/tracker/internal/db/core"
)

const v1betaInsertGraphData = `
REPLACE INTO graph_data
(k1, k2, k3, kind, encoding, data, date_deleted, last_modified)
VALUES (:k1, :k2, :k3, :kind, :encoding, :data, NULL, :last_modified);
`

const v1betaDeleteGraphData = `
UPDATE graph_data
SET date_deleted = :date_deleted
WHERE k1 = :k1
AND k2 = :k2
AND k3 = :k3;
`

const v1betaListGraphData = `
SELECT k1, k2, k3, kind, encoding, data
FROM graph_data
WHERE kind = :kind
AND date_deleted IS NULL
LIMIT :limit OFFSET :offset;
`

const v1betaSelectFromNeighbor = `
SELECT g1.k1, g1.k2, g1.k3, g1.kind, g1.encoding, g1.data
FROM graph_data AS g1
INNER JOIN graph_data AS g2 ON g1.k1 = g2.k2
WHERE g2.k1 IN (:keys)
  AND g2.k1 != g2.k2
  AND g2.date_deleted IS NULL
  AND g1.k1 = g1.k2
  AND g1.date_deleted IS NULL

UNION ALL

SELECT k1, k2, k3, kind, encoding, data
FROM graph_data
WHERE k1 in (:keys)
  AND k1 != k2
  AND date_deleted is NULL;
`

const v1betaSelectToNeighbor = `
SELECT g1.k1, g1.k2, g1.k3, g1.kind, g1.encoding, g1.data
FROM graph_data AS g1
INNER JOIN graph_data AS g2 ON g1.k2 = g2.k1
WHERE g2.k2 IN (:keys)
  AND g2.k1 != g2.k2
  AND g2.date_deleted IS NULL
  AND g1.k1 = g1.k2
  AND g1.date_deleted IS NULL

UNION ALL

SELECT k1, k2, k3, kind, encoding, data
FROM graph_data
WHERE k2 in (:keys)
  AND k1 != k2
  AND date_deleted is NULL;
`

// V1Beta expose statements that are specific to the V1Beta SQLite backend.
var V1Beta = &core.Statements{
	InsertGraphData:        v1betaInsertGraphData,
	DeleteGraphData:        v1betaDeleteGraphData,
	ListGraphData:          v1betaListGraphData,
	SelectInTreeNeighbors:  v1betaSelectToNeighbor,
	SelectOutTreeNeighbors: v1betaSelectFromNeighbor,
}
