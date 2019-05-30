package main

import (
	"github.com/asdine/storm"
	"github.com/gdamore/tcell"
	"github.com/google/uuid"
	"os"
)

type item struct {
	ID         uuid.UUID
	IsDone     bool
	Contents   string
	Tags       []string
	x, y       int
	isSelected bool
}

func (i item) show(s tcell.Screen) {
	status := "✗" // notdone
	style := tcell.StyleDefault.Foreground(tcell.ColorRed)
	if i.IsDone {
		style = style.Foreground(tcell.ColorGreen)
		status = "✓" // done
	}
	puts(s, style, i.x, i.y, status)

	style = tcell.StyleDefault
	if i.isSelected {
		style = style.Underline(true).Bold(true)
	}
	puts(s, style, i.x+2, i.y, i.Contents)
}

func (i *item) click(ev *tcell.EventMouse) {
	_, y := ev.Position()
	if i.y == y {
		i.IsDone = !i.IsDone
	}
}

type todolist struct {
	selected int
	items    []*item
	db       *storm.DB
}

func (tl *todolist) connectdb(name string) error {
	home := os.Getenv("HOME")
	if _, err := os.Stat(home + "/.local/share/todo"); os.IsNotExist(err) {
		err := os.MkdirAll(home+"/.local/share/todo", os.ModePerm)
		if err != nil {
			return err
		}
	}
	db, err := storm.Open(home + "/.local/share/todo/" + name)
	if err != nil {
		return err
	}
	tl.db = db
	err = tl.db.Init(&item{})
	if err != nil {
		return err
	}

	tl.restore()
	return nil
}

func (tl todolist) show(s tcell.Screen) {
	for y, i := range tl.items {
		i.y = y
		i.show(s)
	}
}

func (tl *todolist) selectNew(dir int) {
	for _, i := range tl.items {
		i.isSelected = false
	}
	tl.selected += dir
	if tl.selected > len(tl.items)-1 {
		tl.selected = 0
	}
	if tl.selected < 0 {
		tl.selected = len(tl.items) - 1
	}
	tl.items[tl.selected].isSelected = true
}

func (tl *todolist) push(contents string) {
	i := &item{
		ID:       uuid.New(),
		IsDone:   false,
		Contents: contents,
		x:        1, y: len(tl.items),
		isSelected: false,
	}
	if len(tl.items) == 0 {
		i.isSelected = true
	}
	tl.items = append(tl.items, i)
	tl.db.Save(i)
}
func (tl *todolist) delselected() {
	tl.db.DeleteStruct(tl.items[tl.selected])
	sel := tl.selected
	tl.selectNew(-1)
	tl.items = append(tl.items[:sel], tl.items[sel+1:]...)
}

func (tl *todolist) toggleSelected() {
	tl.items[tl.selected].IsDone = !tl.items[tl.selected].IsDone
	tl.db.Save(tl.items[tl.selected])
}

func (tl *todolist) restore() {
	tl.db.All(&tl.items)
	for i := range tl.items {
		tl.items[i].x = 1
		tl.items[i].y = i
		tl.items[i].isSelected = false
	}
	if len(tl.items) > 0 {
		tl.items[0].isSelected = true
	}
}

func (tl *todolist) save() {
	for _, i := range tl.items {
		tl.db.Save(i)
	}
}
