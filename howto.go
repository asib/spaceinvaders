package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type AttributedText struct {
	fg termbox.Attribute
	bg termbox.Attribute
	s  string
}

const (
	fgHowto          = neonGreen
	bgHowto          = termbox.ColorBlack
	fgHowtoControl   = magenta
	instructionsWPad = 20
	instructionsHPad = 2
)

var (
	instructions = []AttributedText{
		AttributedText{fgHowto, bgHowto, `  xx
 xOOx
xxxxxx     = 10 pts
 /\/\

x    x
xxOOxx
 xxxx      = 20 pts
 /  \

  xx
xOxxOx
xxxxxx     = 30 pts
 /||\

  xxxxx
xxoxOxoxx  = ?? pts
 ##   ##


`},
		AttributedText{fgHowtoControl, bgHowto, `Left`},
		AttributedText{fgHowto, bgHowto, `/`},
		AttributedText{fgHowtoControl, bgHowto, `Right`},
		AttributedText{fgHowto, bgHowto, ` Arrow to move.

`},
		AttributedText{fgHowtoControl, bgHowto, `Space`},
		AttributedText{fgHowto, bgHowto, ` to fire.


Press ESC to close
this window.`}}
	instructionsLines  []string
	instructionsWidth  int
	instructionsHeight int
)

func init() {
	str := ""
	for _, v := range instructions {
		str += v.s
	}
	instructionsLines = strings.Split(str, "\n")
	instructionsWidth = len(instructionsLines[2])
	instructionsHeight = len(instructionsLines)
}

func (g *Game) DrawHowto() {
	g.DrawMenu()

	w, h := instructionsWidth+instructionsWPad, instructionsHeight+instructionsHPad
	x, y := g.w/2-(instructionsWidth+instructionsWPad)/2, logoY

	tbrect(x, y, w, h, fgHowto, bgHowto, true)

	x += instructionsWPad / 2
	y += instructionsHPad / 2

	leftMargin := x

	for _, l := range instructions {
		for _, c := range l.s {
			if c != '\n' {
				termbox.SetCell(x, y, c, l.fg, l.bg)
				x++
			} else {
				y++
				x = leftMargin
			}
		}
	}
}

func (g *Game) UpdateHowto() {
	g.UpdateMenu()
}

func (g *Game) HandleKeyHowto(k termbox.Key) {
	switch k {
	case termbox.KeyEsc:
		g.GoMenu()
		g.hmi = Howto
	}
}

func (g *Game) GoHowto() {
	g.state = HowtoState
	g.cfg = fgMenu
	g.cbg = bgMenu
}
