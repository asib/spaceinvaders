package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	instructions = `  xx
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


Left/Right Arrow
to move.


Press ESC to close
this window.`
	fgHowto          = 0x53
	bgHowto          = termbox.ColorBlack
	instructionsWPad = 20
	instructionsHPad = 2
)

var (
	instructionsLines  = strings.Split(instructions, "\n")
	instructionsWidth  = len(instructionsLines[2])
	instructionsHeight = len(instructionsLines)
)

func (g *Game) DrawHowto() {
	g.DrawMenu()

	w, h := instructionsWidth+instructionsWPad, instructionsHeight+instructionsHPad
	x, y := g.w/2-(instructionsWidth+instructionsWPad)/2, logoY

	tbrect(x, y, w, h, fgHowto, bgHowto, true)

	x += instructionsWPad / 2
	y += instructionsHPad / 2

	for _, l := range instructionsLines {
		tbprint(x, y, fgHowto, bgHowto, l)
		y++
	}
}

func (g *Game) UpdateHowto() {
	g.UpdateMenu()
}

func (g *Game) GoHowto() {
	g.state = HowtoState
}
