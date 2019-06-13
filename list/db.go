package list

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	// import dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"path"
	"time"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New()
	return scope.SetColumn("ID", uuid)
}

func (tl *Todo) connectdb(loc string) error {
	p := path.Dir(loc)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return err
		}
	}
	db, err := gorm.Open("sqlite3", loc)
	if err != nil {
		return err
	}
	tl.db = db
	db.AutoMigrate(&Tag{}, &Item{})
	tl.Restore()
	return nil
}

// Restore list from db
func (tl *Todo) Restore() {
	tl.db.Preload("Tags").Find(&tl.Items)
	for i := range tl.Items {
		tl.Items[i].IsSelected = false
	}
	if len(tl.Items) > 0 {
		tl.Items[0].IsSelected = true
	}
}

// Save list to db
func (tl *Todo) Save() {
	for _, i := range tl.Items {
		tl.db.Save(i)
	}
}

// AddTag what it says on the tin
func (tl *Todo) AddTag(t string) {
	tag := Tag{}
	tl.db.First(&tag, "tag = ?", t)
	if tag.ID == uuid.Nil {
		tag = Tag{Tag: t}
	}
	tag.Items = append(tag.Items, tl.Items[tl.selected])
	tl.db.Save(&tag)
	tl.Items[tl.selected].Tags = append(tl.Items[tl.selected].Tags, &tag)
	// tl.db.Save(tl.Items[tl.selected])
}

// GetItemsByTag what it says on the tin
func (tl *Todo) GetItemsByTag(t string) []Item {
	var items []Item
	tag := Tag{}

	tl.db.First(&tag, "tag = ?", t)
	tl.db.Model(&tag).Related(&items, "Items")

	return items
}
