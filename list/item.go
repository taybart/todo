package list

// Item in list
type Item struct {
	Base
	IsDone     bool
	Contents   string
	Tags       []*Tag `gorm:"many2many:item_tag;"`
	IsSelected bool   `gorm:"-"`
}

// Tag for projects
type Tag struct {
	Base
	Tag   string
	Items []*Item `gorm:"many2many:item_tag;"`
}

func (t *Tag) String() string {
	return t.Tag
}
