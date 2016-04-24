package models

import (
	"encoding/gob"
)

const (
	TIMEFORMAT = "2006-01-02 15:04:05"
)

func init() {
	//注册 以备session使用复杂的数据结构
	gob.Register(&Account{})
	gob.Register(&Category{})
	gob.Register(&Article{})
	gob.Register(&Comment{})
	gob.Register(&FriendLink{})
	gob.Register(&Advert{})
	gob.Register(&Statistics{})
	gob.Register(&Settings{})
}
