package types

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/deps-cloud/dts/pkg/store"
)

func key(params... string) []byte {
	hash := sha256.New()
	for _, param := range params {
		hash.Write([]byte(param))
	}
	return hash.Sum(nil)
}

// Source defines the data on the `source` node
type Source struct {
	URL string `json:"url"`
}

// SourceKey encodes the Source into it's corresponding key
func SourceKey(source *Source) []byte {
	return key(source.URL)
}

// Manages defines the edge between a `source` and it's `modules`
type Manages struct {
	Language string `json:"language"`
	System   string `json:"system"`
	Version  string `json:"version"`
}

// Module defines a library within the graph
type Module struct {
	Language     string `json:"language"`
	Organization string `json:"organization"`
	Module       string `json:"module"`
}

// ModuleKey encodes the Module into it's corresponding key
func ModuleKey(module *Module) []byte {
	return key(module.Language, module.Organization, module.Module)
}

// Depends defines the edge between two modules
type Depends struct {
	VersionConstraint string   `json:"version_constraint"`
	Scopes            []string `json:"scopes"`
}

// DataType defines the data type of the GraphItem
type DataType = string

const (
	// SourceType represents a Source
	SourceType DataType = "source"
	// ManagesType represents a Manage
	ManagesType DataType = "manages"
	// ModuleType represents a Module
	ModuleType DataType = "module"
	// DependsType represents a Depends
	DependsType DataType = "depends"
)

// Encode will serialize the provided item into a []byte for storage
func Encode(item interface{}) ([]byte, store.GraphItemEncoding, error) {
	body, err := json.Marshal(item)
	return body, store.EncodingJSON, err
}

// Decode handles decoding the provided graph item into the corresponding type
func Decode(graphItem *store.GraphItem) (interface{}, error) {
	itemType := graphItem.GraphItemType
	enc := graphItem.Encoding

	var item interface{}

	if itemType == SourceType {
		item = &Source{}
	} else if itemType == ManagesType {
		item = &Manages{}
	} else if itemType == ModuleType {
		item = &Module{}
	} else if itemType == DependsType {
		item = &Depends{}
	} else {
		return nil, fmt.Errorf("unrecognized node type")
	}

	var err error
	if enc == store.EncodingJSON {
		err = json.Unmarshal(graphItem.GraphItemData, item)
	} else {
		return nil, fmt.Errorf("unrecognized encoding")
	}

	return item, err
}
