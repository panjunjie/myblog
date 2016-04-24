package models

import (
	"github.com/go-gorp/gorp"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

//文章
type FriendLink struct {
	Id      int64
	Name    string
	Intro   string
	Link    string
	Read    int64
	Status  bool
	Created time.Time
	Updated time.Time
}

func (m *FriendLink) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&m.Id:     "friendlink.Id",
		&m.Name:   "friendlink.Name",
		&m.Intro:  "friendlink.Intro",
		&m.Link:   "friendlink.Link",
		&m.Read:   "friendlink.Read",
		&m.Status: "friendlink.Status",
		&m.Created: binding.Field{
			Form:       "friendlink.Created",
			TimeFormat: TIMEFORMAT,
		},
		&m.Updated: binding.Field{
			Form:       "friendlink.Updated",
			TimeFormat: TIMEFORMAT,
		},
	}
}

func (m *FriendLink) PreInsert(s gorp.SqlExecutor) error {
	m.Created = time.Now()
	m.Updated = time.Now()
	return nil
}

func (m *FriendLink) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = time.Now()
	return nil
}
