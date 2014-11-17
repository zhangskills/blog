package models

import (
	"time"
)

type Blog struct {
	Id      int64     `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	ViewNum int       `json:"viewNum"`
	Created time.Time `json:"created" xorm:"created"`
	Updated time.Time `json:"updated" xorm:"updated"`
	Deleted time.Time `json:"-" xorm:"deleted"`

	Tags []*Tag `json:"tags" xorm:"-"`
}

type Tag struct {
	Id      int64     `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created" xorm:"created"`
}

type BlogTag struct {
	Id      int64
	BlogId  int64
	TagId   int64
	Created time.Time `xorm:"created"`
}

type KeyCount struct {
	Key   string
	Count int64
}
