package v1beta

import (
	"database/sql"
	"time"

	"gorm.io/gorm/schema"
)

// Encoding describes how the data is stored
type Encoding = uint8

const (
	// EncodingUnspecified signals the data format was not specified.
	EncodingUnspecified Encoding = iota
	// EncodingJSON stores the data as JSON.
	EncodingJSON
	// EncodingProtocolBuffers stores the data as encoded protocol buffers.
	EncodingProtocolBuffers
)

const defaultCollectionName = "graph_data"

// GraphData describes how the data is stored
type GraphData struct {
	CollectionName string		 `json:"-"`
	K1             string        `json:"k1"           db:"k1"            gorm:"column:k1;type:varchar(64);primaryKey;index:secondary,priority:2"`
	K2             string        `json:"k2"           db:"k2"            gorm:"column:k2;type:varchar(64);primaryKey;index:secondary,priority:1"`
	K3             string        `json:"k3"           db:"k3"            gorm:"column:k3;type:varchar(64);primaryKey;index:secondary,priority:3"`
	Kind           string        `json:"kind"         db:"kind"          gorm:"column:kind;type:varchar(55);index:kind"`
	Encoding       Encoding      `json:"encoding"     db:"encoding"      gorm:"column:encoding"`
	Data           string        `json:"data"         db:"data"          gorm:"column:data;type:text"`
	DateDeleted    *sql.NullTime `json:"dateDeleted"  db:"date_deleted"  gorm:"column:date_deleted;index:date_deleted"`
	LastModified   time.Time     `json:"lastModified" db:"last_modified" gorm:"column:last_modified;index:last_modified"`
}

// TableName returns the name of the collection or the default if not set.
func (d *GraphData) TableName() string {
	if d.CollectionName != "" {
		return d.CollectionName
	}
	return defaultCollectionName
}

var _ schema.Tabler = &GraphData{}
