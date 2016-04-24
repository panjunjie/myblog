package models

import (
	"github.com/go-gorp/gorp"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

//广告实体
type Advert struct {
	Id       int64
	MetaLink string
	Ad1      string
	Ad2      string
	Ad3      string
	Ad4      string
	Ad5      string
	Created  time.Time
	Updated  time.Time
}

func (m *Advert) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&m.Id:       "advert.Id",
		&m.MetaLink: "advert.MetaLink",
		&m.Ad1:      "advert.Ad1",
		&m.Ad2:      "advert.Ad2",
		&m.Ad3:      "advert.Ad3",
		&m.Ad4:      "advert.Ad4",
		&m.Ad5:      "advert.Ad5",
		&m.Created: binding.Field{
			Form:       "advert.Created",
			TimeFormat: TIMEFORMAT,
		},
		&m.Updated: binding.Field{
			Form:       "advert.Updated",
			TimeFormat: TIMEFORMAT,
		},
	}
}

func (m *Advert) PreInsert(s gorp.SqlExecutor) error {
	m.Created = time.Now()
	m.Updated = time.Now()
	return nil
}

func (m *Advert) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = time.Now()
	return nil
}
