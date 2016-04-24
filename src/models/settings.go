package models

import (
	"github.com/go-gorp/gorp"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

//博客设置实体
type Settings struct {
	Id            int64
	Favicon       string
	BlogName      string
	BlogIntro     string
	EmailSmtp     string
	EmailSmtpPort int
	EmailSender   string
	EmailPwd      string
	Org           string
	Donate        string
	ICPNO         string
	Created       time.Time
	Updated       time.Time
}

func (m *Settings) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&m.Id:            "settings.Id",
		&m.Favicon:       "settings.Favicon",
		&m.BlogName:      "settings.BlogName",
		&m.BlogIntro:     "settings.BlogIntro",
		&m.EmailSmtp:     "settings.EmailSmtp",
		&m.EmailSmtpPort: "settings.EmailSmtpPort",
		&m.EmailSender:   "settings.EmailSender",
		&m.EmailPwd:      "settings.EmailPwd",
		&m.Org:           "settings.Org",
		&m.Donate:        "settings.Donate",
		&m.ICPNO:         "settings.ICPNO",
		&m.Created: binding.Field{
			Form:       "settings.Created",
			TimeFormat: TIMEFORMAT,
		},
		&m.Updated: binding.Field{
			Form:       "settings.Updated",
			TimeFormat: TIMEFORMAT,
		},
	}
}

func (m *Settings) PreInsert(s gorp.SqlExecutor) error {
	m.Created = time.Now()
	m.Updated = time.Now()
	return nil
}

func (m *Settings) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = time.Now()
	return nil
}
