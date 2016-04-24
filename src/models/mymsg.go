package models

import (
	//"github.com/go-gorp/gorp"
	//"github.com/mholt/binding"
	//"net/http"
	"time"
)

// 给我的留言
type MyMsg struct {
	Id         int64
	UserName   string
	Email      string
	Content    string
	RefContent string
	Status     bool
	Created    time.Time
	Updated    time.Time
}
