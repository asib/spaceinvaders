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

	smAlienSprite = `  xx
 xOOx
xxxxxx
 /\/\`
	medAlienSprite = `x    x
xxOOxx
 xxxx
 /  \`
	lgAlienSprite = `  xx
xOxxOx
xxxxxx
 /||\`
	ufoSprite = `  xxxxx
xxoxOxoxx
 ##   ##`
)
