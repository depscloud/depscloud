package dbsqlite

import (
	"github.com/depscloud/depscloud/services/tracker/internal/db/core"
)

const v1alphaInsertGraphData = `
REPLACE INTO dts_graphdata 
(graph_item_type, k1, k2, k3, encoding, graph_item_data, last_modified, date_deleted)
VALUES (:graph_item_type, :k1, :k2, :k3, :encoding, :graph_item_data, :last_modified, NULL);
`

const v1alphaDeleteGraphData = `
UPDATE dts_graphdata
SET date_deleted = :date_deleted
WHERE (graph_item_type = :graph_item_type and k1 = :k1 and k2 = :k2 and k3 = :k3);
`

const v1alphaListGraphData = `
SELECT graph_item_type, k1, k2, encoding, graph_item_data
FROM dts_graphdata
WHERE graph_item_type = :graph_item_type 
LIMIT :limit OFFSET :offset;
`

const v1alphaSelectFromNeighbor = `
SELECT g1.graph_item_type, g1.k1, g1.k2, g1.encoding, g1.graph_item_data,
       g2.graph_item_type, g2.k1, g2.k2, g2.k3, g2.encoding, g2.graph_item_data
FROM dts_graphdata AS g1
INNER JOIN dts_graphdata AS g2 ON g1.k1 = g2.k2
WHERE g2.k1 IN (:keys) 
AND g2.graph_item_type IN (:edge_types) 
AND g2.k1 != g2.k2 
AND g2.date_deleted IS NULL
AND g1.graph_item_type IN (:node_types)
AND g1.k1 = g1.k2 
AND g1.date_deleted IS NULL;
`

const v1alphaSelectToNeighbor = `
SELECT g1.graph_item_type, g1.k1, g1.k2, g1.encoding, g1.graph_item_data,
	   g2.graph_item_type, g2.k1, g2.k2, g2.k3, g2.encoding, g2.graph_item_data
FROM dts_graphdata AS g1
INNER JOIN dts_graphdata AS g2 ON g1.k2 = g2.k1
WHERE g2.k2 IN (:keys) 
AND g2.graph_item_type IN (:edge_types) 
AND g2.k1 != g2.k2 
AND g2.date_deleted IS NULL
AND g1.graph_item_type IN (:node_types)
AND g1.k1 = g1.k2 
AND g1.date_deleted IS NULL;
`

// V1Alpha expose statements that are specific to the V1Alpha SQLite backend.
var V1Alpha = &core.Statements{
	InsertGraphData:        v1alphaInsertGraphData,
	DeleteGraphData:        v1alphaDeleteGraphData,
	ListGraphData:          v1alphaListGraphData,
	SelectInTreeNeighbors:  v1alphaSelectToNeighbor,
	SelectOutTreeNeighbors: v1alphaSelectFromNeighbor,
}
