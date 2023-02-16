package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func handleEvents(g *game, gg *gameGraphics) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			return false
		case *sdl.MouseButtonEvent:
			if g.gameStage != Running {
				println("game not running")
				g.currentSelection = [2]int{-1, -1}
				break
			}
			if e.X > int32((gg.gridDim-1)*g.getGameConstants().boardWidth) || e.Y > int32((gg.gridDim-1)*g.getGameConstants().boardHeight) {
				println("out of bounds")
				g.currentSelection = [2]int{-1, -1}
				break
			}
			tileX, tileY := getTileMouseLocation(e.X, e.Y, int32(gg.gridDim))
			if e.Type == sdl.MOUSEBUTTONDOWN {
				g.currentSelection = [2]int{int(tileX), int(tileY)}
			}
			if g.currentSelection[0] != -1 && g.currentSelection[1] != -1 && e.Type == sdl.MOUSEBUTTONUP {
				if tileX == int32(g.currentSelection[0]-1) && tileY == int32(g.currentSelection[1]-1) ||
					tileX == int32(g.currentSelection[0]+1) && tileY == int32(g.currentSelection[1]+1) ||
					tileX == int32(g.currentSelection[0]-1) && tileY == int32(g.currentSelection[1]+1) ||
					tileX == int32(g.currentSelection[0]+1) && tileY == int32(g.currentSelection[1]-1) {
					g.makeMove(g.currentSelection[0], g.currentSelection[1], int(tileX), int(tileY))
				}
			}
			// println(e.Button)
		}
	}
	return true
}

func getTileMouseLocation(mouseX int32, mouseY int32, gridDim int32) (int32, int32) {
	tileX := mouseX / (gridDim - 1)
	tileY := mouseY / (gridDim - 1)
	return tileX, tileY
}
