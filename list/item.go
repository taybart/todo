package list

import (
	"github.com/jinzhu/gorm"
	// import dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Item in list
type Item struct {
	gorm.Model
	IsDone     bool
	Contents   string
	Tags       []*Tag `gorm:"many2many:item_tag;"`
	IsSelected bool   `gorm:"-"`
}

// Tag for projects
type Tag struct {
	gorm.Model
	Tag   string
	Items []*Item `gorm:"many2many:item_tag;"`
}

func (t *Tag) String() string {
	return t.Tag
}
