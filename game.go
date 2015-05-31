package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func tbrect(x, y, w, h int, fg, bg termbox.Attribute, border bool) {
	end := " " + strings.Repeat("_", w)
	if border {
		tbprint(x, y-1, fg, bg, end)
	}

	s := strings.Repeat(" ", w)
	if border {
		s = fmt.Sprintf("%c%s%c", '|', s, '|')
	}

	for i := 0; i < h; i++ {
		tbprint(x, y, fg, bg, s)
		y++
	}

	if border {
		tbprint(x, y, fg, bg, end)
	}
}

const (
	fgDefault = termbox.ColorRed
	bgDefault = termbox.ColorYellow
)

// GameState is used as an enum
type GameState uint8

const (
	MenuState GameState = iota
	HowtoState
)

type Game struct {
	state GameState
	evq   chan termbox.Event
	timer <-chan time.Time
	// frame counter
	fc uint8
	// highlighted menu item
	hmi int
	w   int
	h   int
	fg  termbox.Attribute
	bg  termbox.Attribute
}

func NewGame() *Game {
	return &Game{
		evq:   make(chan termbox.Event),
		timer: time.Tick(33 * time.Millisecond),
		fc:    1,
	}
}

// Tick allows us to rate limit the FPS
func (g *Game) Tick() {
	<-g.timer
	g.fc++
	if g.fc > 30 {
		g.fc = 1
	}
}

func (g *Game) Listen() {
	go func() {
		for {
			g.evq <- termbox.PollEvent()
		}
	}()
}

func (g *Game) HandleKey(k termbox.Key) {
	switch g.state {
	case MenuState:
		switch k {
		case termbox.KeyArrowLeft:
			// because of Go's bad mod operator, have to add the length here
			g.hmi = (g.hmi - 1 + NumMenuItems) % NumMenuItems
		case termbox.KeyArrowRight:
			g.hmi = (g.hmi + 1) % NumMenuItems
		case termbox.KeyEnter:
			if g.hmi == Howto {
				g.GoHowto()
			}
		}
	case HowtoState:
		switch k {
		case termbox.KeyEsc:
			g.GoMenu()
			g.hmi = Howto
		}
	}
}

func (g *Game) FitScreen() {
	termbox.Clear(g.fg, g.bg)
	g.w, g.h = termbox.Size()
	g.Draw()
}

func (g *Game) Draw() {
	termbox.Clear(g.fg, g.bg)

	switch g.state {
	case MenuState:
		g.DrawMenu()
	case HowtoState:
		g.DrawHowto()
	}

	termbox.Flush()
}

func (g *Game) Update() {
	g.Tick()

	switch g.state {
	case MenuState:
		g.UpdateMenu()
	case HowtoState:
		g.UpdateHowto()
	}

	return
}

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatalln(err)
	}
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	g := NewGame()
	g.Listen()
	g.GoMenu()
	g.FitScreen()

main:
	for {
		select {
		case ev := <-g.evq:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case 0:
					if ev.Ch == 'q' {
						break main
					}
				default:
					g.HandleKey(ev.Key)
				}
			case termbox.EventResize:
				g.FitScreen()
			}
		default:
		}

		g.Update()
		g.Draw()
	}
}
