package types

// Source defines the data on the `source` node
type Source struct {
	URL string `json:"url"`
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
