package list

import (
	"github.com/jinzhu/gorm"
	// import dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"os"
	"path"
	"sort"
)

// Todo list
type Todo struct {
	selected int
	Items    []*Item
	db       *gorm.DB
}

// NewTodo create with default
func NewTodo(db string) (*Todo, error) {
	tl := &Todo{
		selected: 0,
	}

	if db != "" {
		err := tl.connectdb(path.Dir(db), path.Base(db))
		if err != nil {
			return nil, errors.Wrap(err, "Connecting to db")
		}
	} else {
		home := os.Getenv("HOME")
		p := home + "/.local/share/todo"
		err := tl.connectdb(p, "")
		if err != nil {
			return nil, errors.Wrap(err, "Connecting to db")
		}
	}
	return tl, nil
}

type byTime []*Item

func (b byTime) Len() int           { return len(b) }
func (b byTime) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byTime) Less(i, j int) bool { return b[i].CreatedAt.Before(b[j].CreatedAt) }

// Sort by time
func (tl *Todo) Sort() {
	sort.Sort(byTime(tl.Items))
}

// Sel a new item, up = -1, down = 1
func (tl *Todo) Sel(dir int) {
	if len(tl.Items) > 0 {
		for _, i := range tl.Items {
			i.IsSelected = false
		}
		tl.selected += dir
		if tl.selected > len(tl.Items)-1 {
			tl.selected = 0
		}
		if tl.selected < 0 {
			tl.selected = len(tl.Items) - 1
		}
		tl.Items[tl.selected].IsSelected = true
	}
}

// SetSel by index
func (tl *Todo) SetSel(index int) {
	if index > 0 && index < len(tl.Items) {
		tl.selected = index
		tl.Sel(0)
	}
}

// Push new item to list
func (tl *Todo) Push(contents string) {
	i := &Item{
		IsDone:     false,
		Contents:   contents,
		IsSelected: false,
	}
	if len(tl.Items) == 0 {
		i.IsSelected = true
	}
	tl.Items = append(tl.Items, i)
	tl.db.Save(i)
}

// Edit current item
func (tl *Todo) Edit(contents string) {
	tl.Items[tl.selected].Contents = contents
	tl.db.Save(tl.Items[tl.selected])
}

// Del selected item
func (tl *Todo) Del() {
	tl.db.Delete(tl.Items[tl.selected])
	sel := tl.selected
	tl.Items = append(tl.Items[:sel], tl.Items[sel+1:]...)
	tl.Sel(0)
}

// Toggle selected item
func (tl *Todo) Toggle() {
	tl.Items[tl.selected].IsDone = !tl.Items[tl.selected].IsDone
	tl.db.Save(tl.Items[tl.selected])
}

// Current item
func (tl *Todo) Current() *Item {
	return tl.Items[tl.selected]
}

// Close list
func (tl *Todo) Close() {
	tl.db.Close()
}
