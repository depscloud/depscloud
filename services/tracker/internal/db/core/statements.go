package core

// Statements are used by the sqlDriver to perform operations against different backends.
type Statements struct {
	InsertGraphData string `json:"insertGraphData"`
	DeleteGraphData string `json:"deleteGraphData"`
	ListGraphData   string `json:"listGraphData"`

	// SelectOutTreeNeighbors - out-tree refers to traversing edges pointing away from the current node.
	SelectOutTreeNeighbors string `json:"selectFromNeighbors"`

	// SelectInTreeNeighbors - in-tree refers to traversing edges pointing away from the current node.
	SelectInTreeNeighbors string `json:"selectToNeighbors"`
}
