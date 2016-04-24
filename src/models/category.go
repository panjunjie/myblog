package models

import (
	"github.com/go-gorp/gorp"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

type Category struct {
	Id        int64
	Name      string
	Seq       int
	Intro     string
	Status    bool
	ItemCount int
	Created   time.Time
	Updated   time.Time
}

func (m *Category) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&m.Id:        "category.Id",
		&m.Name:      "category.Name",
		&m.Seq:       "category.Seq",
		&m.Intro:     "category.Intro",
		&m.Status:    "category.Status",
		&m.ItemCount: "category.ItemCount",
		&m.Created: binding.Field{
			Form:       "category.Created",
			TimeFormat: TIMEFORMAT,
		},
		&m.Updated: binding.Field{
			Form:       "category.Updated",
			TimeFormat: TIMEFORMAT,
		},
	}
}

func (m *Category) PreInsert(s gorp.SqlExecutor) error {
	m.Created = time.Now()
	m.Updated = time.Now()
	return nil
}

func (m *Category) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = time.Now()
	return nil
}
