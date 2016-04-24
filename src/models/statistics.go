package models

import (
	"github.com/go-gorp/gorp"
	"time"
)

//统计实体
type Statistics struct {
	Id       int64
	Category int //Category 1:统计博客数评论数访问数点赞数；2：归档博客；3：热门博客；4：热评博客；5：点赞博客排名；
	Title1   string
	Title2   string
	Title3   string
	Title4   string
	Param1   string
	Param2   string
	Param3   string
	Param4   string
	Created  time.Time
	Updated  time.Time
}

func (m *Statistics) PreInsert(s gorp.SqlExecutor) error {
	m.Created = time.Now()
	m.Updated = time.Now()
	return nil
}

func (m *Statistics) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = time.Now()
	return nil
}
