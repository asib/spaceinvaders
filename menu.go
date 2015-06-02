package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

type Point struct {
	x int
	y int
}

type Star struct {
	Point
	vx int
	vy int
}

const (
	menuPad         = 10
	fgMenu          = red
	bgMenu          = termbox.ColorBlack
	fgMenuHighlight = termbox.ColorBlack
	bgMenuHighlight = neonGreen
	fgStar          = termbox.ColorWhite
	bgStar          = termbox.ColorBlack
	numStars        = 50
	starSymbol      = "."
	logoY           = 3
	logo            = `_______________  ________________        ____              _______________  _______________
I             I  I              I       /    \             I             I  I             I
I             I  I    -------   I      /  /\  \            I             I  I             I
I     ---------  I    I     I   I     /  /  \  \           I             I  I             I
I     I          I    -------   I    /  /    \  \          I       -------  I       -------
I     I          I              I   /  /______\  \         I       I        I       I      
I     I          I      ---------  /              \        I       I        I       -------
I     ---------  I      I         /       __       \       I       I        I             I
I             I  I      I        /       /  \       \      I       I        I             I
I             I  I      I       /       /    \       \     I       I        I       -------
I             I  I      I      /       /      \       \    I       I        I       I      
------        I  I      I     /       /        \       \   I       -------  I       -------
I             I  I      I    /       /          \       \  I             I  I             I
I             I  I      I   I        I          I       I  I             I  I             I
I             I  I      I   I        I          I       I  I             I  I             I
---------------  --------   ----------          ---------  ---------------  ---------------
________  _____    ____ __          __          ____   --------| ---------      ___________
I      I  I    \   I  I \ \        / / /\       I   \  I   ----| I ----- I      I         I
--    --  I     \  I  I  \ \      / / /  \      I I\ \ I   I     I I   I I      I     -----
  I  I    I   I\ \ I  I   \ \    / / / oo \     I I/ / I   ----| I -----  \     I         I
  I  I    I   I \ \I  I    \ \  / / /  __  \    I   /  I   ----| I   ____  \    I         I
--    --  I   I  \ \  I     \ \/ / /  /  \  \   I  /   I   I     I   I   \  \   ------    I
I      I  I   I   \   I      \  / /  /    \  \  I /    I   ----| I   I    \  \  I         I
--------  -----    ----       \/ /__/      \__\ I/     --------| -----     ---  -----------`
)

const (
	FirstMenuItem     = 0
	Play          int = iota - 1
	Highscores
	Howto
	NumMenuItems
)

var (
	menuItems      = map[int]string{Play: "PLAY", Highscores: "HIGHSCORES", Howto: "HOWTO"}
	logoLines      = strings.Split(logo, "\n")
	logoLineLength = len(logoLines[0])
	logoHeight     = len(logoLines)
	stars          = make([]*Star, 0, numStars)
)

func PrintLogo(x, y int, fg, bg termbox.Attribute, lines []string) {
	for _, line := range lines {
		tbprint(x, y, fg, bg, line)
		y++
	}
}

func (g *Game) DrawMenu() {
	x := g.w/2 - logoLineLength/2
	y := logoY
	PrintLogo(x, y, fgMenu, bgMenu, logoLines)

	length := 0
	i := 0
	for _, v := range menuItems {
		length += len(v)
		if i+1 != len(menuItems) {
			length += menuPad
		}
		i++
	}

	x = g.w/2 - length/2
	y += logoHeight + 5
	for i := FirstMenuItem; i < NumMenuItems; i++ {
		v := menuItems[i]
		if i == g.hmi {
			tbprint(x, y, fgMenuHighlight, bgMenuHighlight, v)
		} else {
			tbprint(x, y, fgMenu, bgMenu, v)
		}
		x += len(v) + menuPad
	}

	for _, s := range stars {
		tbprint(s.x, s.y, fgStar, bgStar, starSymbol)
	}
}

func (g *Game) UpdateMenu() {
	if len(stars) != cap(stars) && g.fc%3 == 0 {
		n := len(stars)
		stars = stars[0 : n+1]
		stars[n] = NewStar(g.w, rand.Intn(g.h))
	}

	for i, s := range stars {
		s.x += s.vx
		s.y += s.vy

		if s.x < 0 || s.x > g.w || s.y < 0 || s.y > g.h {
			stars[i] = NewStar(g.w, rand.Intn(g.h))
		}
	}
}

func NewStar(x, y int) *Star {
	rand.Seed(time.Now().UTC().UnixNano())

	vx, vy := -1*(1+rand.Intn(3)), 0
	/*
	 *for vx == 0 || vy == 0 {
	 *  vx = rand.Intn(3) - 1
	 *  vy = rand.Intn(3) - 1
	 *}
	 */

	return &Star{Point{x, y}, vx, vy}
}

func (g *Game) HandleKeyMenu(k termbox.Key) {
	switch k {
	case termbox.KeyArrowLeft:
		// because of Go's bad mod operator, have to add the length here
		g.hmi = (g.hmi - 1 + NumMenuItems) % NumMenuItems
	case termbox.KeyArrowRight:
		g.hmi = (g.hmi + 1) % NumMenuItems
	case termbox.KeyEnter:
		switch g.hmi {
		case Howto:
			g.GoHowto()
		case Play:
			g.GoPlay()
		}
	}
}

func (g *Game) GoMenu() {
	g.state = MenuState
	g.cfg = fgMenu
	g.cbg = bgMenu
	g.hmi = FirstMenuItem
}
