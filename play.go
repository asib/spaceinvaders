package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
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
	playerSpriteHere         = 11

	scoreText      = "Score: "
	scorex, scorey = 10, 1

	livesText        = "Lives: "
	livesRightOffset = 0

	ufoMoveEvery = 3
	ufoIndex     = 666
	ufoReward    = 100

	alienBulletSpeed   = 1
	alienShootValMax   = 100
	alienShootVal      = 1
	alienPadVertical   = 1
	alienPadHorizontal = 3

	// rwd is reward
	rwdSm  = 10
	rwdMd  = 20
	rwdLg  = 30
	rwdUfo = 100

	alienStartx, alienStarty = 10, 7

	numBarricades = 4

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

	ufo      *RegEntity
	ufoTimer <-chan time.Time

	aliens           []*Alien
	alienBullets     []*Bullet
	alienSpriteIndex int
	alienMoveEvery   uint8
	alienv           [2]int
	rowsSm           = 2
	rowsMd           = 2
	rowsLg           = 1
	aliensHorizontal = 1

	barricadePositions [][]int

	lvl int
)

func (g *Game) DrawPlay() {
	tbprintsprite(player.x, player.y, player.fg, player.bg, player.sprite)

	for i := range barricadePositions {
		for j := range barricadePositions[i] {
			if barricadePositions[i][j] != nonIndex {
				tbprint(i, j, fgBarricade, bgBarricade, "x")
			}
		}
	}

	for _, b := range alienBullets {
		if b != nil {
			tbprintsprite(b.x, b.y, b.fg, b.bg, b.sprite)
		}
	}

	if ufo != nil {
		tbprintsprite(ufo.x, ufo.y, ufo.fg, ufo.bg, ufo.sprite)
	}

	for _, a := range aliens {
		if a != nil {
			tbprintsprite(a.x, a.y, a.fg, a.bg, a.sprite[alienSpriteIndex])
		}
	}

	if player.bullet != nil {
		tbprintsprite(player.bullet.x, player.bullet.y,
			player.bullet.fg, player.bullet.bg, player.bullet.sprite)
	}

	tbprint(scorex, scorey, fgPlayText, bgPlayText, scoreText+fmt.Sprintf("%d", player.score))

	livesStr := livesText + strings.Replace(strings.Repeat(livesSprite, player.lives), "⏣ ", "⏣  ", -1)
	livesx = g.w - livesRightOffset - len(livesStr)
	livesy = scorey
	tbprint(livesx, livesy, fgPlayText, bgPlayText, livesStr)
}

func (g *Game) PlayerPositions() [][]int {
	screen := make([][]int, g.w)
	for i := range screen {
		screen[i] = make([]int, g.h)
		for j := range screen[i] {
			screen[i][j] = nonIndex
		}
	}

	x, y := player.x, player.y
	initx := x
	lines := strings.Split(player.sprite, "\n")
	for _, l := range lines {
		for _, c := range l {
			if c != ' ' {
				screen[x][y] = playerSpriteHere
			}

			x++
		}

		y++
		x = initx
	}

	return screen
}

func (g *Game) AlienPositions() [][]int {
	screen := make([][]int, g.w)
	for i := range screen {
		screen[i] = make([]int, g.h)
		for j := range screen[i] {
			screen[i][j] = nonIndex
		}
	}

	if ufo != nil && ufo.x > 0 && ufo.x < g.w-ufoSpriteWidth {
		x, y := ufo.x, ufo.y
		initx := x
		lines := strings.Split(ufo.sprite, "\n")
		for _, l := range lines {
			for _, c := range l {
				if c != ' ' {
					screen[x][y] = ufoIndex
				}
				x++
			}
			y++
			x = initx
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

func (g *Game) WipeBullets() {
	player.bullet = nil
	alienBullets = make([]*Bullet, int(aliensHorizontal*numRows/10))
}

func (g *Game) barricadeYPos() int {
	return (g.h - playerSpriteBottomOffset - playerSpriteHeight) - barricadeSpriteHeight - 2
}

func (g *Game) genBarricades() [][]int {
	screen := make([][]int, g.w)
	for i := range screen {
		screen[i] = make([]int, g.h)
		for j := range screen[i] {
			screen[i][j] = nonIndex
		}
	}

	gap := (g.w - 4*barricadeSpriteWidth) / 5
	x, y := gap, g.barricadeYPos()
	inity := y
	lines := strings.Split(barricadeSprite, "\n")
	for i := 0; i < numBarricades; i++ {
		initx := gap*(i+1) + barricadeSpriteWidth*i
		x = initx
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
		y = inity
	}

	return screen
}

func (g *Game) wipePlay() {
	// calculate rows/cols
	i := 0
	for x := 0; x < (g.w / 2); x, i = x+(alienSpriteWidth+alienPadHorizontal), i+1 {
		aliensHorizontal = i
	}
	i = 0
	for y := alienStarty; y < (g.barricadeYPos() - alienSpriteHeight - alienPadVertical); y, i = y+(alienSpriteHeight+alienPadVertical), i+1 {
		if i == 5 {
			break
		}
	}
	switch {
	case i <= 3:
		rowsSm = 1
		rowsMd = 1
	case i == 4:
		rowsSm = 2
		rowsMd = 1
	default:
		rowsSm = 2
		rowsMd = 2
	}

	player = nil
	aliens = make([]*Alien, aliensHorizontal*numRows)
	alienBullets = make([]*Bullet, int(aliensHorizontal*numRows/10))
	alienSpriteIndex = 0
	alienMoveEvery = uint8(15)
	alienv = rightMove

	barricadePositions = g.genBarricades()

	lvl = 1
}

func (g *Game) drawGetName(name string, showLenWarn bool) {
	const (
		msg                  = "You set a new highscore!"
		prompt               = "Please enter a name 3-10 characters long:"
		lenWarn              = "Name must be 3-10 characters long!"
		getNameHeight        = 8
		getNameHeightLenWarn = 11
		getNameWidthPad      = 4
		fgGetName            = white
		bgGetName            = termbox.ColorBlack
		fgGetNameName        = neonGreen
	)
	g.DrawPlay()
	if len(name) < 10 {
		name += "_"
	}
	w, h := len(prompt)+getNameWidthPad, getNameHeight
	if showLenWarn {
		h = getNameHeightLenWarn
	}
	x, y := g.w/2-w/2, g.h/2-h/2
	tbrect(x, y, w, h, fgGetName, bgGetName, true)

	// prompt
	x += getNameWidthPad/2 + 1
	y += 2
	tbprint(x, y, fgGetName, bgGetName, msg)
	y += 2
	tbprint(x, y, fgGetName, bgGetName, prompt)

	// name
	x = g.w/2 - len(name)/2
	y += 2
	tbprint(x, y, fgGetNameName, bgGetName, name)

	// showLenWarn ?
	if showLenWarn {
		x = g.w/2 - len(lenWarn)/2
		y += 3
		tbprint(x, y, fgGetName, bgGetName, lenWarn)
	}

	termbox.Flush()
}

func (g *Game) getName() string {
	name := ""
	showLenWarn := false
	for {
		g.drawGetName(name, showLenWarn)
		if len(name) >= 10 {
			showLenWarn = true
		}

		select {
		case ev := <-g.evq:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEnter:
					if len(name) < 3 {
						showLenWarn = true
					} else {
						return name
					}
				case 0:
					if len(name) < 10 {
						name += string(ev.Ch)
					}
				case termbox.KeyBackspace:
					fallthrough
				case termbox.KeyBackspace2:
					if len(name) > 0 {
						name = name[:len(name)-1]
					}
				}
			default:
			}
		default:
		}
	}
}

func (g *Game) checkHighscores() {
	if len(g.highscores) < maxHighscores || player.score > g.highscores[len(g.highscores)-1].score {
		name := g.getName()
		g.highscores = append(g.highscores, &Highscore{player.score, name})
		sort.Sort(sort.Reverse(ByScore(g.highscores)))
		if len(g.highscores) > maxHighscores {
			g.highscores = append([]*Highscore(nil), g.highscores[:maxHighscores]...)
		}

		// write highscores
		data := ""
		for i, score := range g.highscores {
			data += fmt.Sprintf("%s%s%d", score.name, highscoreSeparator, score.score)
			if i != len(g.highscores)-1 {
				data += "\n"
			}
		}
		ioutil.WriteFile(highscoreFilename, []byte(data), 0666)
	}
}

func (g *Game) gameOver() {
	g.FreezeFlash("GAME OVER")
	g.checkHighscores()
	g.wipePlay()
	g.GoMenu()
	// TODO: finish this function, need to add highscore stuff
}

func newUfoTimer() <-chan time.Time {
	return time.After(time.Duration(rand.Intn(20)+15) * time.Second)
}

func newUfo() *RegEntity {
	return &RegEntity{Entity{0 - ufoSpriteWidth, 4, fgUfo, bgUfo}, ufoSprite}
}

func (g *Game) UpdatePlay() {
	playerPos := g.PlayerPositions()
	for b := range alienBullets {
		if alienBullets[b] != nil {
			alienBullets[b].y += alienBulletSpeed

			if alienBullets[b].y >= g.h {
				alienBullets[b] = nil
			} else {
				x, y := alienBullets[b].x, alienBullets[b].y
				if playerPos[x][y] != nonIndex {
					player.lives -= 1
					if player.lives == 0 {
						g.gameOver()
						return
					}
					g.WipeBullets()
					// TODO: clear the UFO if it's there
					// g.ClearUFO()
					g.FreezeFlash(lvlFlash())
					return
				} else if barricadePositions[x][y] != nonIndex {
					alienBullets[b] = nil
					barricadePositions[x][y] = nonIndex
				}
			}
		}
	}

	if player.bullet != nil {
		player.bullet.y += player.bullet.vy
		if player.bullet.y < 0 {
			player.bullet = nil
		} else {
			screen := g.AlienPositions()
			x, y := player.bullet.x, player.bullet.y
			if screen[x][y] != nonIndex {
				player.bullet = nil
				if screen[x][y] == ufoIndex {
					player.score += ufoReward
					ufo = nil
				} else {
					player.score += aliens[screen[x][y]].reward
					aliens[screen[x][y]] = nil
				}
			} else if barricadePositions[x][y] != nonIndex {
				player.bullet = nil
			}
		}
	}

	rand.Seed(time.Now().UnixNano())

	if ufo != nil && g.fc%ufoMoveEvery == 0 {
		ufo.x += 1
		if ufo.x > g.w {
			ufo = nil
			ufoTimer = newUfoTimer()
		}
	}

	if g.fc%alienMoveEvery == 0 {
		alienSpriteIndex = (alienSpriteIndex + 1) % 2

		downFlag := false
		xval := 999999999 // some meaningless number
		levelComplete := true
		for i := 0; i < len(aliens); i++ {
			if aliens[i] != nil {
				levelComplete = false
				aliens[i].x += alienv[0]
				aliens[i].y += alienv[1]

				if aliens[i].x <= 0 || aliens[i].x+alienSpriteWidth >= g.w {
					downFlag = true
					xval = aliens[i].x
				}

				if aliens[i].y >= player.y-playerSpriteHeight {
					g.gameOver()
					return
				}

				// try firing
				if rand.Intn(alienShootValMax) == 6 {
					for j := range alienBullets {
						if alienBullets[j] == nil {
							alienBullets[j] = NewBullet(aliens[i].x+alienSpriteWidth/2,
								aliens[i].y, alienBulletSpeed)
							break
						}
					}
				}
			}
		}

		if levelComplete && ufo == nil {
			lvl += 1
			if alienMoveEvery != 2 {
				alienMoveEvery--
			}
			g.BeginNextLevel()
			g.WipeBullets()
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

	select {
	case <-ufoTimer:
		ufo = newUfo()
		ufoTimer = nil
	default:
		if ufoTimer == nil && ufo == nil {
			ufoTimer = newUfoTimer()
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
			aliens[arrayOffset+i*aliensHorizontal+j] = NewAlien(x, y, fg, bg, sprite, reward)
			x += spriteW + alienPadHorizontal
		}
		x = startx
		y += spriteH + alienPadVertical
	}

	return y, arrayOffset + rows*cols
}

func lvlFlash() string {
	return fmt.Sprintf("Level %d", lvl)
}

func (g *Game) FreezeFlash(m string) {
	g.Draw()

	tbprint(g.w/2-len(m)/2, g.h/2, fgPlayText, bgPlayText, m)
	termbox.Flush()

	time.Sleep(flashDuration)
}

func (g *Game) BeginNextLevel() {
	g.FreezeFlash(lvlFlash())

	x, y, offset := alienStartx, alienStarty, 0
	y, offset = makeAliens(x, y, rowsLg, aliensHorizontal,
		alienSpriteWidth, alienSpriteHeight, rwdLg,
		offset, fgAlien, bgAlien, lgAlienSprite)
	y, offset = makeAliens(x, y, rowsMd, aliensHorizontal,
		alienSpriteWidth, alienSpriteHeight, rwdMd,
		offset, fgAlien, bgAlien, mdAlienSprite)
	makeAliens(x, y, rowsSm, aliensHorizontal,
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
		g.wipePlay()
		player = &Player{
			RegEntity{Entity{startx, starty, fgPlayer, bgPlayer}, playerSprite},
			0,
			initLives,
			nil,
		}

		g.BeginNextLevel()
	}
}
