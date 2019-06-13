package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/taybart/todo/list"
)

const (
	up   = -1
	down = 1
)

const (
	_ = iota
	neutral
	additem
	edititem
	addtag
	search
)

var s tcell.Screen

func main() {

	quickone := flag.Bool("q", false, "quick note")
	msg := flag.String("m", "", "message to add")
	db := flag.String("db", "", "db location")
	flag.Parse()

	tl, err := list.NewTodo(*db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tl.Close()
	if *quickone {
		if *msg == "" {
			fmt.Fprintf(os.Stderr, "You must specify a message with -m\n")
			os.Exit(1)
		}
		tl.Push(*msg)
		fmt.Println("Added:", *msg)
		os.Exit(0)
	}

	s, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer s.Fini()

	encoding.Register()
	if err = s.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorDefault).
		Background(tcell.ColorDefault))
	s.EnableMouse()
	s.Clear()

	quit := make(chan struct{})

	opts := options{
		showTime: false,
		padding:  1,
	}

	go func() {
		state := neutral
		contents := ""
		for {
			s.Clear()
			showlist(opts, tl)
			cmdstr := contents + string(tcell.RuneS9)
			switch state {
			case additem:
				cmdstr = "a: " + cmdstr
			case edititem:
				cmdstr = "e: " + cmdstr
			case addtag:
				cmdstr = "tag: " + cmdstr
			case search:
				cmdstr = "/ " + cmdstr
			default:
				cmdstr = ""
			}
			_, height := s.Size()
			puts(tcell.StyleDefault, 1, height-1, cmdstr)
			s.Show()
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					if state > 1 {
						state = neutral
					} else {
						close(quit)
						return
					}
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					sz := len(contents)
					if sz > 0 {
						contents = contents[:sz-1]
					}
				case tcell.KeyEnter:
					switch state {
					case additem:
						tl.Push(contents)
					case edititem:
						tl.Edit(contents)
					case addtag:
						tl.AddTag(contents)
					case neutral:
						tl.Toggle()
					}
					state = neutral
					contents = ""
				case tcell.KeyUp:
					tl.Sel(up)
				case tcell.KeyDown:
					tl.Sel(down)
				case tcell.KeyCtrlT:
					opts.showTime = !opts.showTime
				default:
					if state != neutral {
						contents += string(ev.Rune())
					}
				}

				if state == neutral {
					switch ev.Rune() {
					case 'k':
						tl.Sel(up)
					case 'j':
						tl.Sel(down)
					case 'a', ':':
						state = additem
					case 'e':
						state = edititem
						contents = tl.Current().Contents
					case 't':
						state = addtag
					case 'd':
						tl.Del()
					case '/':
						state = search
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
					for row := range tl.Items {
						clickitem(ev, opts, row, tl)
					}
				}
			}
		}
	}()
	<-quit
}
