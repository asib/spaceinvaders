package main

import "github.com/nsf/termbox-go"

const (
	fgWarn = white
	bgWarn = termbox.ColorBlack
	warn1  = "Your terminal screen is too small."
	warn2  = "Please make the screen wider/taller"
	warn3  = "and/or reduce your font size."
)

func (g *Game) DrawWarn() {
	tbprint(g.w/2-len(warn1)/2, g.h/2, fgWarn, bgWarn, warn1)
	tbprint(g.w/2-len(warn2)/2, g.h/2+1, fgWarn, bgWarn, warn2)
	tbprint(g.w/2-len(warn3)/2, g.h/2+2, fgWarn, bgWarn, warn3)
}

func (g *Game) GoWarn() {
	g.state = WarnState
	g.cfg = fgMenu
	g.cbg = bgMenu
}
