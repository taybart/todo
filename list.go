package main

import (
	"github.com/gdamore/tcell"
	"github.com/taybart/todo/list"
)

type options struct {
	showTime bool
	padding  int
}

func showitem(opts options, x, y int, i *list.Item) {
	status := "✗" // notdone
	style := tcell.StyleDefault.Foreground(tcell.ColorRed)
	if i.IsDone {
		style = style.Foreground(tcell.ColorGreen)
		status = "✓" // done
	}
	puts(style, opts.padding, y+opts.padding, status)

	str := i.Contents
	if len(i.Tags) > 0 {
		str += " => "
		for i, t := range i.Tags {
			if i != 0 {
				str += ", "
			}
			str += "#" + t.String()
		}
	}
	style = tcell.StyleDefault
	if i.IsSelected {
		style = style.Underline(true).Bold(true)
	}
	puts(style, x+2+opts.padding, y+opts.padding, str)
	if opts.showTime {
		width, _ := s.Size()
		created := i.CreatedAt.Format("2006-01-02 15:04:05")
		puts(style, width-(len(created)+opts.padding), y+opts.padding, created)
	}
}

func showlist(opts options, tl *list.Todo) {
	for y, i := range tl.Items {
		showitem(opts, 1, y, i)
	}
}

func clickitem(ev *tcell.EventMouse, opts options, i int, tl *list.Todo) {
	_, y := ev.Position()
	if i+opts.padding == y {
		tl.Items[i].IsDone = !tl.Items[i].IsDone
		tl.SetSel(i)
		tl.Save()
	}
}
