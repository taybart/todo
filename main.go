package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

func main() {
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	encoding.Register()
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorDefault).
		Background(tcell.ColorDefault))
	s.EnableMouse()
	s.Clear()

	quit := make(chan struct{})

	tl := todolist{
		selected: 0,
	}
	err := tl.connectdb("todolist.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	go func() {
		adding := false
		contents := ""
		for {
			s.Clear()
			tl.show(s)
			if adding {
				_, height := s.Size()
				puts(s, tcell.StyleDefault, 1, height-1, "add:"+contents+string(tcell.RuneS9))
			}
			s.Show()
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					if adding {
						adding = false
					} else {
						close(quit)
						return
					}
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if adding {
						sz := len(contents)
						if sz > 0 {
							contents = contents[:sz-1]
						}
					}
				case tcell.KeyEnter:
					if adding {
						adding = false
						tl.push(contents)
						contents = ""
					} else {
						tl.toggleSelected()
					}
				case tcell.KeyUp:
					tl.selectNew(-1)
				case tcell.KeyDown:
					tl.selectNew(1)
				default:
					if adding {
						contents += string(ev.Rune())
					}
				}

				if !adding {
					switch ev.Rune() {
					case 'k':
						tl.selectNew(-1)
					case 'j':
						tl.selectNew(1)
					case 'a', ':':
						adding = true
					case 'd':
						tl.delselected()
					case 'q':
						close(quit)
						return
					}
				}

			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventMouse:
				switch ev.Buttons() {
				case tcell.Button1:
					for _, i := range tl.items {
						i.click(ev)
					}
				}
			}
		}
	}()

	<-quit
	tl.db.Close()
	s.Fini()
}
