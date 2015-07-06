package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	fgPlay     = termbox.ColorBlack
	bgPlay     = termbox.ColorBlack
	fgPlayText = neonGreen
	bgPlayText = termbox.ColorBlack
	fgAlien    = white
	bgAlien    = termbox.ColorBlack

	playerSpriteBottomOffset = 2
	initLives                = 5
	playerMoveSpeed          = 2
	playerBulletSpeed        = -1

	scoreText      = "Score: "
	scorex, scorey = 10, 1

	livesText        = "Lives: "
	livesRightOffset = 0

	alienPadVertical    = 1
	alienPadHorizontal  = 3
	maxAliensHorizontal = 11
	rowsSm              = 2
	rowsMd              = 2
	rowsLg              = 1

	// rwd is reward
	rwdSm  = 10
	rwdMd  = 20
	rwdLg  = 30
	rwdUfo = 100

	alienStartx, alienStarty = 10, 3

	flashDuration = 1 * time.Second

	nonIndex = -50
)

type Entity struct {
	x, y   int
	fg, bg termbox.Attribute
}

type AnimatedEntity struct {
	Entity
	sprite [2]string
}

type RegEntity struct {
	Entity
	sprite string
}

type Bullet struct {
	RegEntity
	vy int
}

type Player struct {
	RegEntity
	score, lives int
	bullet       *Bullet
}

type Alien struct {
	AnimatedEntity
	reward int
}

func NewBullet(x, y, vy int) *Bullet {
	return &Bullet{RegEntity{Entity{x, y, fgBullet, bgBullet}, bulletSprite}, vy}
}

func NewAlien(x, y int, fg, bg termbox.Attribute, sprite [2]string, reward int) *Alien {
	return &Alien{AnimatedEntity{Entity{x, y, fg, bg}, sprite}, reward}
}

var (
	livesx, livesy int

	startx, starty int
	player         *Player

	numRows = rowsSm + rowsMd + rowsLg

	rightMove = [2]int{1, 0}
	leftMove  = [2]int{-1, 0}
	downMove  = [2]int{0, 1}

	// +1 on the end because we can have 1 active UFO at any given time
	aliens           = make([]*Alien, maxAliensHorizontal*numRows+1)
	alienSpriteIndex = 0
	alienMoveEvery   = uint8(15)
	alienv           = rightMove

	lvl = 1
)

func (g *Game) DrawPlay() {
	tbprint(scorex, scorey, fgPlayText, bgPlayText, scoreText+fmt.Sprintf("%d", player.score))

	livesStr := livesText + strings.Replace(strings.Repeat(livesSprite, player.lives), "⏣ ", "⏣  ", -1)
	livesx = g.w - livesRightOffset - len(livesStr)
	livesy = scorey
	tbprint(livesx, livesy, fgPlayText, bgPlayText, livesStr)

	tbprintsprite(player.x, player.y, player.fg, player.bg, player.sprite)

	for _, a := range aliens {
		if a != nil {
			tbprintsprite(a.x, a.y, a.fg, a.bg, a.sprite[alienSpriteIndex])
		}
	}

	if player.bullet != nil {
		tbprintsprite(player.bullet.x, player.bullet.y,
			player.bullet.fg, player.bullet.bg, player.bullet.sprite)
	}
}

func (g *Game) AlienPositions() [][]int {
	screen := make([][]int, g.w)
	for i := range screen {
		screen[i] = make([]int, g.h)
		for j := range screen[i] {
			screen[i][j] = nonIndex
		}
	}

	for i, a := range aliens {
		if a != nil {
			x, y := a.x, a.y
			initx := x
			lines := strings.Split(a.sprite[alienSpriteIndex], "\n")
			for _, l := range lines {
				for _, c := range l {
					if c != ' ' {
						screen[x][y] = i
					}

					x++
				}
				y++
				x = initx
			}
		}
	}

	return screen
}

func (g *Game) UpdatePlay() {
	if player.bullet != nil {
		player.bullet.y += player.bullet.vy
		if player.bullet.y < 0 {
			player.bullet = nil
		} else {
			screen := g.AlienPositions()
			x, y := player.bullet.x, player.bullet.y
			if screen[x][y] != nonIndex {
				player.bullet = nil
				player.score += aliens[screen[x][y]].reward
				aliens[screen[x][y]] = nil
			}
		}
	}

	if g.fc%alienMoveEvery == 0 {
		alienSpriteIndex = (alienSpriteIndex + 1) % 2

		downFlag := false
		xval := 999999999 // some meaningless number
		for i := 0; i < len(aliens); i++ {
			if aliens[i] != nil {
				aliens[i].x += alienv[0]
				aliens[i].y += alienv[1]

				if aliens[i].x <= 0 || aliens[i].x+alienSpriteWidth >= g.w {
					downFlag = true
					xval = aliens[i].x
				}
			}
		}

		switch {
		case alienv == downMove:
			if xval == 0 {
				alienv = rightMove
			} else {
				alienv = leftMove
			}
		case downFlag:
			alienv = downMove
		}
	}
}

func (g *Game) HandleKeyPlay(k termbox.Key) {
	switch k {
	case termbox.KeyArrowRight:
		player.x += playerMoveSpeed
	case termbox.KeyArrowLeft:
		player.x -= playerMoveSpeed
	case termbox.KeySpace:
		if player.bullet == nil {
			player.bullet = NewBullet(player.x+playerSpriteWidth/2, player.y,
				playerBulletSpeed)
		}
	}

	switch {
	case player.x+playerSpriteWidth > g.w:
		player.x = g.w - playerSpriteWidth
	case player.x < 0:
		player.x = 0
	}
}

func makeAliens(x, y, rows, cols, spriteW, spriteH, reward, arrayOffset int,
	fg, bg termbox.Attribute, sprite [2]string) (int, int) {
	startx := x

	// create aliens
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			aliens[arrayOffset+i*maxAliensHorizontal+j] = NewAlien(x, y, fg, bg, sprite, reward)
			x += spriteW + alienPadHorizontal
		}
		x = startx
		y += spriteH + alienPadVertical
	}

	return y, arrayOffset + rows*cols
}

func (g *Game) BeginNextLevel() {
	g.Draw()

	flash := fmt.Sprintf("Level %d", lvl)
	tbprint(g.w/2-len(flash)/2, g.h/2, fgPlayText, bgPlayText, flash)
	termbox.Flush()

	time.Sleep(flashDuration)

	x, y, offset := alienStartx, alienStarty, 0
	y, offset = makeAliens(x, y, rowsLg, maxAliensHorizontal,
		alienSpriteWidth, alienSpriteHeight, rwdLg,
		offset, fgAlien, bgAlien, lgAlienSprite)
	y, offset = makeAliens(x, y, rowsMd, maxAliensHorizontal,
		alienSpriteWidth, alienSpriteHeight, rwdMd,
		offset, fgAlien, bgAlien, mdAlienSprite)
	makeAliens(x, y, rowsSm, maxAliensHorizontal,
		alienSpriteWidth, alienSpriteHeight, rwdSm,
		offset, fgAlien, bgAlien, smAlienSprite)
}

func (g *Game) GoPlay() {
	g.state = PlayState
	g.cfg = fgPlay
	g.cbg = bgPlay

	startx = g.w/2 - playerSpriteWidth/2
	starty = g.h - playerSpriteBottomOffset - playerSpriteHeight

	if player == nil {
		player = &Player{
			RegEntity{Entity{startx, starty, fgPlayer, bgPlayer}, playerSprite},
			0,
			initLives,
			nil,
		}

		g.BeginNextLevel()
	}
}
