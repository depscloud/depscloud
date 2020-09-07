package v1beta

// Statements are used by the sqlDriver to perform operations against different backends.
type Statements struct {
	InsertGraphData    string `json:"insertGraphData"`
	DeleteGraphData    string `json:"deleteGraphData"`
	ListGraphData      string `json:"listGraphData"`
	SelectFromNeighbor string `json:"selectFromNeighbors"`
	SelectToNeighbor   string `json:"selectToNeighbors"`
}
