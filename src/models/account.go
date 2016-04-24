package models

import (
	"github.com/go-gorp/gorp"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

type Account struct {
	Id        int64
	UserName  string
	Password  string
	Email     string
	Sex       int
	Image     string
	Sign      string
	Intro     string
	Status    bool
	Created   time.Time
	Updated   time.Time
	LastLogin time.Time
}

func (m *Account) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&m.Id:       "account.Id",
		&m.UserName: "account.UserName",
		&m.Password: "account.Password",
		&m.Email:    "account.Email",
		&m.Sex:      "account.Sex",
		&m.Image:    "account.Image",
		&m.Sign:     "account.Sign",
		&m.Intro:    "account.Intro",
		&m.Status:   "account.Status",
		&m.Created: binding.Field{
			Form:       "account.Created",
			TimeFormat: TIMEFORMAT,
		},
		&m.Updated: binding.Field{
			Form:       "account.Updated",
			TimeFormat: TIMEFORMAT,
		},
		&m.LastLogin: binding.Field{
			Form:       "account.LastLogin",
			TimeFormat: TIMEFORMAT,
		},
	}
}

func (m *Account) PreInsert(s gorp.SqlExecutor) error {
	m.Created = time.Now()
	m.Updated = time.Now()
	return nil
}

func (m *Account) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = time.Now()
	return nil
}
