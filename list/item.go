package list

import (
	"time"
)

// Item in list
type Item struct {
	ID         uint64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	IsDone     bool
	Contents   string
	Tags       []*Tag `gorm:"many2many:item_tag;"`
	IsSelected bool   `gorm:"-"`
}

// Tag for projects
type Tag struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Tag       string
	Items     []*Item `gorm:"many2many:item_tag;"`
}

func (t *Tag) String() string {
	return t.Tag
}
