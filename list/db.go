package list

import (
	"github.com/jinzhu/gorm"
	// import dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

func (tl *Todo) connectdb(path, name string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	n := "td.db"
	if name != "" {
		n = name
	}
	db, err := gorm.Open("sqlite3", path+"/"+n)
	if err != nil {
		return err
	}
	tl.db = db
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Item{})
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
	if tag.ID == 0 {
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
