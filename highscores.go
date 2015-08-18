package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

const (
	fgHighscores       = neonGreen
	bgHighscores       = termbox.ColorBlack
	highscoresWidthPad = 7
	scorePad           = 10
	namePad            = 10
	highscoresHeight   = maxHighscores + 6
	title              = "HIGHSCORES"
	prompt             = "Press ESC to exit"
)

func (g *Game) DrawHighscores() {
	g.DrawMenu()

	w, h := scorePad+1+namePad+2*highscoresWidthPad, highscoresHeight
	x, y := g.w/2-w/2, logoY
	tbrect(x, y, w, h, fgHighscores, bgHighscores, true)

	y += 2
	tbprint(g.w/2-len(title)/2, y, fgHighscores, bgHighscores, title)

	y += 2
	x += highscoresWidthPad
	for _, hs := range g.highscores {
		tbprint(x, y, fgHighscores, bgHighscores, fmt.Sprintf("%-10s %010d", hs.name, hs.score))
		y++
	}
	for i := 0; i < maxHighscores-len(g.highscores); i++ {
		tbprint(x, y, fgHighscores, bgHighscores, fmt.Sprintf("%-10s %010d", "?????", 0))
		y++
	}

	y++
	x = g.w/2 - len(prompt)/2
	p1 := "Press "
	p2 := "ESC "
	p3 := "to exit"
	tbprint(x, y, fgHighscores, bgHighscores, p1)
	x += len(p1)
	tbprint(x, y, magenta, bgHighscores, p2)
	x += len(p2)
	tbprint(x, y, fgHighscores, bgHighscores, p3)
}

func (g *Game) UpdateHighscores() {
	g.UpdateMenu()
}

func (g *Game) HandleKeyHighscores(k termbox.Key) {
	switch k {
	case termbox.KeyEsc:
		g.GoMenu()
		g.hmi = Highscores
	}
}

func (g *Game) GoHighscores() {
	g.state = HighscoresState
	g.cfg = fgMenu
	g.cbg = bgMenu
}
