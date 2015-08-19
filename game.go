package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/simulatedsimian/joystick"
)

type Highscore struct {
	score int
	name  string
}

type ByScore []*Highscore

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].score < a[j].score }

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

// print a multi-line sprite
func tbprintsprite(x, y int, fg, bg termbox.Attribute, sprite string) {
	lines := strings.Split(sprite, "\n")
	for _, l := range lines {
		tbprint(x, y, fg, bg, l)
		y++
	}
}

const (
	highscoreFilename  = "hs"
	highscoreSeparator = ":"
	maxHighscores      = 5
	fgDefault          = termbox.ColorRed
	bgDefault          = termbox.ColorYellow
	fps                = 30
)

// GameState is used as an enum
type GameState uint8

const (
	MenuState GameState = iota
	HowtoState
	PlayState
	HighscoresState
	WarnState
)

type Game struct {
	highscores []*Highscore

	state GameState
	evq   chan termbox.Event
	timer <-chan time.Time

	js joystick.Joystick

	// frame counter
	fc uint8

	// highlighted menu item
	hmi int
	w   int
	h   int

	// fg and bg colors used when termbox.Clear() is called
	cfg termbox.Attribute
	cbg termbox.Attribute
}

func NewGame() *Game {
	return &Game{
		highscores: make([]*Highscore, 0),
		evq:        make(chan termbox.Event),
		timer:      time.Tick(time.Duration(1000/fps) * time.Millisecond),
		fc:         1,
	}
}

// Tick allows us to rate limit the FPS
func (g *Game) Tick() {
	<-g.timer
	g.fc++
	if g.fc > fps {
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
		g.HandleKeyMenu(k)
	case HowtoState:
		g.HandleKeyHowto(k)
	case PlayState:
		g.HandleKeyPlay(k)
	case HighscoresState:
		g.HandleKeyHighscores(k)
	case WarnState:
		g.HandleKeyWarn(k)
	}
}

func (g *Game) FitScreen() {
	termbox.Clear(g.cfg, g.cbg)
	g.w, g.h = termbox.Size()
	g.Draw()
}

func (g *Game) Draw() {
	termbox.Clear(g.cfg, g.cbg)

	switch g.state {
	case MenuState:
		g.DrawMenu()
	case HowtoState:
		g.DrawHowto()
	case PlayState:
		g.DrawPlay()
	case HighscoresState:
		g.DrawHighscores()
	case WarnState:
		g.DrawWarn()
	}

	termbox.Flush()
}

func (g *Game) ReadJoystick() {
	if g.js != nil {
		jstate, err := g.js.Read()
		if err == nil {
			if jstate.Buttons&1 != 0 {
				g.HandleKey(termbox.KeySpace)
			}
			if jstate.AxisData[0] < -10000 {
				g.HandleKey(termbox.KeyArrowLeft)
			}
			if jstate.AxisData[0] > 10000 {
				g.HandleKey(termbox.KeyArrowRight)
			}
		}
	}
}

func (g *Game) Update() {
	g.Tick()

	switch g.state {
	case MenuState:
		g.UpdateMenu()
	case HowtoState:
		g.UpdateHowto()
	case PlayState:
		g.ReadJoystick()
		g.UpdatePlay()
	case HighscoresState:
		g.UpdateHighscores()
	}

	return
}

func (g *Game) loadHighscores() {
	data, err := ioutil.ReadFile(highscoreFilename)
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(data), "\n")
	for _, l := range lines {
		parts := strings.Split(l, highscoreSeparator)
		if i, err := strconv.Atoi(parts[1]); err == nil {
			g.highscores = append(g.highscores, &Highscore{i, parts[0]})
		} else {
			log.Fatalln(err)
		}
	}

	sort.Sort(sort.Reverse(ByScore(g.highscores)))
}

func (g *Game) checkSize() bool {
	if g.w < logoLineLength+8 || g.h < (logoY+logoHeight+5+2) {
		return false
	}
	return true
}

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatalln(err)
	}
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	f, err := os.Create("diwe.log")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(f)

	g := NewGame()

	js, _ := joystick.Open(0)
	g.js = js

	if _, err := os.Stat(highscoreFilename); err == nil {
		g.loadHighscores()
	}

	g.Listen()
	g.FitScreen()
	if g.checkSize() {
		g.GoMenu()
	} else {
		g.GoWarn()
	}
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
