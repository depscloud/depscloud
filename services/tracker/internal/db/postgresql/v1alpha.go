package dbpostgresql

import (
	"github.com/depscloud/depscloud/services/tracker/internal/db/core"
	dbsqlite "github.com/depscloud/depscloud/services/tracker/internal/db/sqlite"
)

const v1alphaInsertGraphData = `
INSERT INTO dts_graphdata 
	(graph_item_type, k1, k2, k3, encoding, graph_item_data, last_modified)
VALUES
	(:graph_item_type, :k1, :k2, :k3, :encoding, :graph_item_data, :last_modified)
ON CONFLICT (graph_item_type, k1, k2, k3) 
DO UPDATE SET
	graph_item_data = EXCLUDED.graph_item_data, 
	encoding = EXCLUDED.encoding, 
	date_deleted = NULL,
	last_modified = EXCLUDED.last_modified;
`

// V1Alpha expose statements that are specific to the V1Alpha PostgreSQL backend.
var V1Alpha = &core.Statements{
	InsertGraphData:        v1alphaInsertGraphData,
	DeleteGraphData:        dbsqlite.V1Alpha.DeleteGraphData,
	ListGraphData:          dbsqlite.V1Alpha.ListGraphData,
	SelectInTreeNeighbors:  dbsqlite.V1Alpha.SelectInTreeNeighbors,
	SelectOutTreeNeighbors: dbsqlite.V1Alpha.SelectOutTreeNeighbors,
}
