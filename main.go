package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	initSDL()
	defer cleanupSDL()

	graphics := newGameGraphics("Dots", 1000, 750)
	defer graphics.cleanup()

	game := newGame(50, 50, 0.4)
	defer game.cleanup()
	// gameConstants := game.getGameConstants()
	// graphics.initEmptyBoard(gameConstants)
	game.init()
	game.startGame()
	running := true
	for running {
		running = handleEvents(game, graphics)
		if game.evaluate {
			i := game.movedDot[0]
			j := game.movedDot[1]
			game.update(i, j)
			game.evaluate = false
		}
		graphics.update(game)
		sdl.Delay(64)
	}
}

func initSDL() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
}

func cleanupSDL() {
	defer sdl.Quit()
}
