package main

import "github.com/nsf/termbox-go"

const (
	fgWarn = white
	bgWarn = termbox.ColorBlack
	warn1  = "Your terminal screen is too small."
	warn2  = "Please make the screen wider/taller"
	warn3  = "and/or reduce your font size."
	warn4  = "Press "
	warn5  = "Space/Enter "
	warn6  = "to retry loading."
)

func (g *Game) HandleKeyWarn(k termbox.Key) {
	switch k {
	case termbox.KeyEnter:
		fallthrough
	case termbox.KeySpace:
		if g.checkSize() {
			g.GoMenu()
		}
	}
}

func (g *Game) DrawWarn() {
	y := g.h / 2
	tbprint(g.w/2-len(warn1)/2, y, fgWarn, bgWarn, warn1)
	y++
	tbprint(g.w/2-len(warn2)/2, y, fgWarn, bgWarn, warn2)
	y++
	tbprint(g.w/2-len(warn3)/2, y, fgWarn, bgWarn, warn3)
	y++
	x := g.w/2 - (len(warn4)+len(warn5)+len(warn6))/2
	tbprint(x, y, fgWarn, bgWarn, warn4)
	x += len(warn4)
	tbprint(x, y, magenta, bgWarn, warn5)
	x += len(warn5)
	tbprint(x, y, fgWarn, bgWarn, warn6)
}

func (g *Game) GoWarn() {
	g.state = WarnState
	g.cfg = fgMenu
	g.cbg = bgMenu
}
