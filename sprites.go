package main

import "github.com/nsf/termbox-go"

const (
	bulletSprite = "."
	fgBullet     = white
	bgBullet     = termbox.ColorBlack

	playerSprite = `  /\
OOxxOO
OXOOXO`
	playerSpriteWidth  = 6
	playerSpriteHeight = 3
	fgPlayer           = neonGreen
	bgPlayer           = termbox.ColorBlack
	livesSprite        = `⏣ `

	alienSpriteWidth  = 8
	alienSpriteHeight = 4

	ufoSpriteWidth  = 9
	ufoSpriteHeight = 3
)

var (
	smAlienSprite = [2]string{`   xx
  xOOx
 xxxxxx
  /\/\`, `   xx
  xOOx
 xxxxxx
  \  /`}

	mdAlienSprite = [2]string{` x    x
 xxOOxx
  xxxx
  /  \`, `
 xxOOxx
x xxxx x
  /  \`}

	lgAlienSprite = [2]string{`   xx
 xOxxOx
 xxxxxx
  /||\`, `   xx
 xOxxOx
 xxxxxx
  \||/`}

	ufoSprite = `  xxxxx
xxoxOxoxx
 ##   ##`
)
