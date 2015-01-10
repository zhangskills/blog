package models

import (
	"time"
)

type Blog struct {
	Id      int64      `json:"id,omitempty"`
	Title   string     `json:"title,omitempty"`
	Content string     `json:"content,omitempty"`
	ViewNum int        `json:"viewNum,omitempty"`
	Created *time.Time `json:"created,omitempty" xorm:"created"`
	Updated *time.Time `json:"updated,omitempty" xorm:"updated"`
	Deleted *time.Time `json:"-" xorm:"deleted"`
	Tags    []*Tag     `json:"tags,omitempty" xorm:"-"`
}

type Tag struct {
	Id      int64      `json:"id,omitempty"`
	Name    string     `json:"name,omitempty"`
	Created *time.Time `json:"created,omitempty" xorm:"created"`
}

type BlogTag struct {
	Id      int64
	BlogId  int64
	TagId   int64
	Created *time.Time `xorm:"created"`
}

type KeyCount struct {
	Key   string
	Count int64
}
