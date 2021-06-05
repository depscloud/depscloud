package dbmysql

import (
	"github.com/depscloud/depscloud/services/tracker/internal/db/core"
	dbsqlite "github.com/depscloud/depscloud/services/tracker/internal/db/sqlite"
)

const v1betaInsertGraphData = `
INSERT INTO graph_data 
	(k1, k2, k3, kind, encoding, data, date_deleted, last_modified)
VALUES
	(:k1, :k2, :k3, :kind, :encoding, :data, NULL, :last_modified)
ON DUPLICATE KEY UPDATE
	kind = :kind,
	encoding = :encoding,
	data = :data, 
	date_deleted = NULL,
	last_modified = :last_modified;
`

// V1Beta expose statements that are specific to the V1Beta MySQL backend.
var V1Beta = &core.Statements{
	InsertGraphData:        v1betaInsertGraphData,
	DeleteGraphData:        dbsqlite.V1Beta.DeleteGraphData,
	ListGraphData:          dbsqlite.V1Beta.ListGraphData,
	SelectInTreeNeighbors:  dbsqlite.V1Beta.SelectInTreeNeighbors,
	SelectOutTreeNeighbors: dbsqlite.V1Beta.SelectOutTreeNeighbors,
}
