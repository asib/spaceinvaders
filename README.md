# Space Invaders
Terminal Space Invaders game written in Go.

#### Installation
```sh
go get -u github.com/asib/spaceinvaders
cd $GOPATH/src/github.com/asib/spaceinvaders
go build
```

You can also simply download a binary for your OS/Arch from [here](https://github.com/asib/spaceinvaders/releases/tag/v1.0).

#### Controls

* Use the arrow keys to move left/right, spacebar to fire.
* Press `q` at any time to quit.

The game will adjust the number of "invaders" to (roughly) fit your terminal's screen size.
This means you can make the game more/less difficult by making your screen bigger/smaller.
__Just make sure you don't resize the screen once you've started playing__, else the game will crash.

__If you're having trouble fitting all the graphics onto your terminal screen, even when it's maximised, lower your font size__.
