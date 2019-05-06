package store

// GraphItemEncoding identifies the encoding that was used for
// serializing the GraphItemData
type GraphItemEncoding = uint8

const (
	// EncodingRaw means that no encoding was used. Mainly for testing.
	EncodingRaw GraphItemEncoding = 0
	// EncodingJSON means that the data was json encoded.
	EncodingJSON GraphItemEncoding = 1
	// EncodingBase64 means that the data was base64 encoded.
	EncodingBase64 GraphItemEncoding = 2
)

// GraphItemEncodings returns all possible encodings that can be used.
func GraphItemEncodings() []GraphItemEncoding {
	return []GraphItemEncoding{
		EncodingRaw,
		EncodingJSON,
		EncodingBase64,
	}
}

// GraphItem defines an item that is stored withing the graph store.
// This can represent either a node or an edge. The format is flexible
// enough to support either structures with minimal overhead. Implementation
// is loosely based off of Dropbox's edgestore system.
//
// K1 == K2 implies a node
// k1 != K2 implies an edge
//
// Primary key: { GraphItemType, K1, K2 }
type GraphItem struct {
	GraphItemType string `json:"graph_type"`
	K1 []byte `json:"k1"`
	K2 []byte `json:"k2"`
	Encoding GraphItemEncoding `json:"encoding"`
	GraphItemData []byte `json:"graph_item_data"`
}

// PrimaryKey defines the primary key to a GraphItem
type PrimaryKey struct {
	GraphItemType string `json:"graph_type"`
	K1 []byte `json:"k1"`
	K2 []byte `json:"k2"`
}

// GraphStore is an interface that allows the backing data storage to be replaced.
// This helps for testing and provides flexibility for alternative implementations.
type GraphStore interface {
	Put(item *GraphItem) error
	FindByPrimary(key *PrimaryKey) (*GraphItem, error)

	FindUpstream(key []byte) ([]*GraphItem, error)
	FindDownstream(key []byte) ([]*GraphItem, error)
}
