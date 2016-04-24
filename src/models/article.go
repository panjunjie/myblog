package models

import (
	"github.com/go-gorp/gorp"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

//文章
type Article struct {
	Id           int64
	Title        string
	Content      string
	CategoryId   int64
	Read         int64
	Tags         string
	Origin       string
	Status       bool
	UpCount      int64
	DownCount    int64
	CommentCount int64
	Created      time.Time
	Updated      time.Time
}

func (m *Article) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&m.Id:           "article.Id",
		&m.Title:        "article.Title",
		&m.Content:      "article.Content",
		&m.CategoryId:   "article.CategoryId",
		&m.Read:         "article.Read",
		&m.Tags:         "article.Tags",
		&m.Origin:       "article.Origin",
		&m.Status:       "article.Status",
		&m.UpCount:      "article.UpCount",
		&m.DownCount:    "article.DownCount",
		&m.CommentCount: "article.CommentCount",
		&m.Created: binding.Field{
			Form:       "article.Created",
			TimeFormat: TIMEFORMAT,
		},
		&m.Updated: binding.Field{
			Form:       "article.Updated",
			TimeFormat: TIMEFORMAT,
		},
	}
}

func (m *Article) PreInsert(s gorp.SqlExecutor) error {
	m.Created = time.Now()
	m.Updated = time.Now()
	return nil
}

func (m *Article) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = time.Now()
	return nil
}

type ArticleExt struct {
	Article
	CommentCount int
}

// 评论
type Comment struct {
	Id         int64
	ArticleId  int64
	UserId     int64
	UserName   string
	Email      string
	Replyer    string
	Content    string
	RefContent string
	UpCount    int64
	DownCount  int64
	Status     bool
	Created    time.Time
	Updated    time.Time
}

func (m *Comment) PreInsert(s gorp.SqlExecutor) error {
	m.Created = time.Now()
	m.Updated = time.Now()
	return nil
}

func (m *Comment) PreUpdate(s gorp.SqlExecutor) error {
	m.Updated = time.Now()
	return nil
}

// 评论
type CommentExt struct {
	Comment
	Title string
}
