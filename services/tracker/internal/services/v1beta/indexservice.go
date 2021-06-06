package v1beta

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Index defines a generic structure that allows us to index some data. Right now, the only use case is
// to index some Node data to help improve some of the usability of the tooling. There's a discoverability
// element that's currently missing... The general idea here is:
//
//   Module => Name => "github.com/depscloud/depscloud" => K1 (node)
//   Depends => Version => "v0.1.0" => K3 (edge)
//
type Index struct {
	Kind  string `gorm:"column:kind;varchar(255);primaryKey;"`
	Field string `gorm:"column:field;varchar(255);primaryKey;"`
	Value string `gorm:"column:value;varchar(4096);primaryKey;index;"`
	Key   string `gorm:"column:key;varchar(64);primaryKey;"`
}

func (i *Index) TableName() string {
	return "graph_data_index"
}

func NewSQLIndexService(rw, ro *gorm.DB) IndexService {
	if rw != nil {
		_ = rw.AutoMigrate(&Index{})
	}

	return &indexService{
		rw: rw,
		ro: ro,
	}
}

type IndexService interface {
	Index(ctx context.Context, fields []*Index) error
	Distinct(ctx context.Context, filter *Index) ([]string, error)
	Query(ctx context.Context, filter *Index) ([]*Index, error)
}

type indexService struct {
	rw *gorm.DB
	ro *gorm.DB
}

func (i *indexService) Index(ctx context.Context, fields []*Index) error {
	return i.rw.Transaction(func(tx *gorm.DB) error {
		return tx.
			WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "kind"},
					{Name: "field"},
					{Name: "value"},
					{Name: "key"},
				},
				DoNothing: true,
			}).
			Create(&fields).
			Error
	}, &sql.TxOptions{
		ReadOnly: false,
	})
}

func (i *indexService) Distinct(ctx context.Context, filter *Index) ([]string, error) {
	results := make([]string, 0)

	err := i.ro.Transaction(func(tx *gorm.DB) error {
		return tx.
			WithContext(ctx).
			Model(&Index{}).
			Where(filter).
			Distinct("value").
			Find(&results).
			Error
	}, &sql.TxOptions{
		ReadOnly: true,
	})

	return results, err
}

func (i *indexService) Query(ctx context.Context, filter *Index) ([]*Index, error) {
	results := make([]*Index, 0)

	// convert Value to valueLike for wildcard search
	valueLike := "%" + filter.Value + "%"
	filter.Value = ""

	err := i.ro.Transaction(func(tx *gorm.DB) error {
		return tx.
			WithContext(ctx).
			Where(filter).
			Where("value LIKE ?", valueLike).
			Find(&results).
			Error
	}, &sql.TxOptions{
		ReadOnly: true,
	})

	return results, err
}
