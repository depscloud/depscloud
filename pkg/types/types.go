package types

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
