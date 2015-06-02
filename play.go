package main

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	fgPlay     = termbox.ColorBlack
	bgPlay     = termbox.ColorBlack
	fgPlayText = neonGreen
	bgPlayText = termbox.ColorBlack

	playerSpriteBottomOffset = 2
	initLives                = 5
	playerMoveSpeed          = 2
	playerBulletSpeed        = -1

	scoreText      = "Score: "
	scorex, scorey = 10, 1

	livesText        = "Lives: "
	livesRightOffset = 0
)

type Entity struct {
	x, y   int
	fg, bg termbox.Attribute
	sprite string
}

type Bullet struct {
	Entity
	vy int
}

type Player struct {
	Entity
	score, lives int
	bullet       *Bullet
}

type Alien struct {
	Entity
	reward int
}

func NewBullet(x, y, vy int) *Bullet {
	return &Bullet{Entity{x, y, fgBullet, bgBullet, bulletSprite}, vy}
}

func NewAlien(x, y int, fg, bg termbox.Attribute, sprite string, reward int) *Alien {
	return &Alien{Entity{x, y, fg, bg, sprite}, reward}
}

var (
	livesx, livesy int

	startx, starty int
	player         *Player
)

func (g *Game) DrawPlay() {
	tbprint(scorex, scorey, fgPlayText, bgPlayText, scoreText+fmt.Sprintf("%d", player.score))

	livesStr := livesText + strings.Replace(strings.Repeat(livesSprite, player.lives), "⏣ ", "⏣  ", -1)
	livesx = g.w - livesRightOffset - len(livesStr)
	livesy = scorey
	tbprint(livesx, livesy, fgPlayText, bgPlayText, livesStr)

	tbprintsprite(player.x, player.y, player.fg, player.bg, player.sprite)

	if player.bullet != nil {
		tbprintsprite(player.bullet.x, player.bullet.y,
			player.bullet.fg, player.bullet.bg, player.bullet.sprite)
	}
}

func (g *Game) UpdatePlay() {
	if player.bullet != nil {
		player.bullet.y += player.bullet.vy
		if player.bullet.y < 0 {
			player.bullet = nil
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

func (g *Game) GoPlay() {
	g.state = PlayState
	g.cfg = fgPlay
	g.cbg = bgPlay

	startx = g.w/2 - playerSpriteWidth/2
	starty = g.h - playerSpriteBottomOffset - playerSpriteHeight

	if player == nil {
		player = &Player{
			Entity{startx, starty, fgPlayer, bgPlayer, playerSprite},
			0,
			initLives,
			nil,
		}
	}
}
