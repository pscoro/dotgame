package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type gameGraphics struct {
	name            string
	width           int
	height          int
	gridDim         int
	gameWindow      *sdl.Window
	gameRenderer    *sdl.Renderer
	backgroundColor color
	gridLineColor   color
	dotColor        color
	claimed1Color   color
	claimed2Color   color
}

type color int

func (c color) toRGB() RGB {
	rgb := RGB{
		R: (int(c) & 0xff0000) >> 16,
		G: (int(c) & 0x00ff00) >> 8,
		B: (int(c) & 0x0000ff),
	}
	return rgb
}

type RGB struct {
	R int
	G int
	B int
}

func newGameGraphics(name string, width int, height int) *gameGraphics {
	gameGraphics := gameGraphics{
		name:   name,
		width:  width,
		height: height,
	}
	gameGraphics.createSDLWindow()
	gameGraphics.createRenderer()
	gameGraphics.initColorDefaults()
	return &gameGraphics
}

func (graphics *gameGraphics) createSDLWindow() {
	window, err := sdl.CreateWindow(
		graphics.name,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(graphics.width),
		int32(graphics.height),
		sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	if window == nil {
		panic("Window is nil")
	}
	graphics.gameWindow = window
}

func (graphics *gameGraphics) createRenderer() {
	renderer, err := sdl.CreateRenderer(graphics.gameWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	if renderer == nil {
		panic("Renderer is nil")
	}
	graphics.gameRenderer = renderer
}

func (graphics *gameGraphics) initColorDefaults() {
	graphics.backgroundColor = 0xffffff
	graphics.gridLineColor = 0x000000
	graphics.dotColor = 0x000000
	graphics.claimed1Color = 0xff0000
	graphics.claimed2Color = 0x0000ff
}

func (graphics *gameGraphics) cleanup() {
	defer graphics.gameWindow.Destroy()
	defer graphics.gameRenderer.Destroy()
}

// func (graphics *gameGraphics) initEmptyBoard(gameConstants *gameConstants) {
// 	surface, err := graphics.gameWindow.GetSurface()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// background := sdl.Rect{0, 0, int32(graphics.width), int32(graphics.height)}
// 	// surface.FillRect(&background, uint32(graphics.backgroundColor))

// 	graphics.gameWindow.UpdateSurface()
// }

func (graphics *gameGraphics) createGrid(boardWidth int, boardHeight int, gridDim int) {
	gridRGB := graphics.gridLineColor.toRGB()

	graphics.gameRenderer.SetDrawColor(uint8(gridRGB.R), uint8(gridRGB.G), uint8(gridRGB.B), 1)
	for i := 0; i < boardHeight; i++ {
		for j := 0; j < boardWidth; j++ {
			gridRect := sdl.Rect{int32(j * (gridDim - 1)), int32(i * (gridDim - 1)), int32(gridDim), int32(gridDim)}
			graphics.gameRenderer.DrawRect(&gridRect)

		}
	}
}

func (graphics *gameGraphics) update(game *game) {
	gridDim := calculateGridDim(graphics.width, graphics.height, game.getGameConstants().boardWidth, game.getGameConstants().boardHeight)
	graphics.gridDim = gridDim
	backgroundRGB := graphics.backgroundColor.toRGB()
	dotRGB := graphics.dotColor.toRGB()
	claimed1RGB := graphics.claimed1Color.toRGB()
	claimed2RGB := graphics.claimed2Color.toRGB()

	backgroundRect := sdl.Rect{0, 0, int32(graphics.width), int32(graphics.width)}
	graphics.gameRenderer.Clear()
	graphics.gameRenderer.SetDrawColor(uint8(backgroundRGB.R), uint8(backgroundRGB.G), uint8(backgroundRGB.B), 1)
	graphics.gameRenderer.FillRect(&backgroundRect)
	graphics.createGrid(game.getGameConstants().boardWidth, game.getGameConstants().boardHeight, gridDim)
	for i := 0; i < game.getGameConstants().boardHeight; i++ {
		for j := 0; j < game.getGameConstants().boardWidth; j++ {
			if *game.getTileAt(j, i) == Dot {
				graphics.gameRenderer.SetDrawColor(uint8(dotRGB.R), uint8(dotRGB.G), uint8(dotRGB.B), 1)
				dotRect := sdl.Rect{int32(j*(gridDim-1) + (gridDim / 3)), int32(i*(gridDim-1) + (gridDim / 3)), int32(gridDim / 3), int32(gridDim / 3)}
				graphics.gameRenderer.FillRect(&dotRect)
			} else if *game.getTileAt(j, i) == Claimed1 {
				graphics.gameRenderer.SetDrawColor(uint8(claimed1RGB.R), uint8(claimed1RGB.G), uint8(claimed1RGB.B), 1)
				claimed1Rect := sdl.Rect{int32(j * (gridDim - 1)), int32(i * (gridDim - 1)), int32(gridDim), int32(gridDim)}
				graphics.gameRenderer.FillRect(&claimed1Rect)
			} else if *game.getTileAt(j, i) == Claimed2 {
				graphics.gameRenderer.SetDrawColor(uint8(claimed2RGB.R), uint8(claimed2RGB.G), uint8(claimed2RGB.B), 1)
				claimed2Rect := sdl.Rect{int32(j * (gridDim - 1)), int32(i * (gridDim - 1)), int32(gridDim), int32(gridDim)}
				graphics.gameRenderer.FillRect(&claimed2Rect)
			}
		}
	}
	if game.currentPlayer == Player1 {
		graphics.gameRenderer.SetDrawColor(uint8(claimed1RGB.R), uint8(claimed1RGB.G), uint8(claimed1RGB.B), 1)
		player1Rect := sdl.Rect{int32(graphics.width) - 250, 75, 100, 100}
		graphics.gameRenderer.FillRect(&player1Rect)
	} else if game.currentPlayer == Player2 {
		graphics.gameRenderer.SetDrawColor(uint8(claimed2RGB.R), uint8(claimed2RGB.G), uint8(claimed2RGB.B), 1)
		player2Rect := sdl.Rect{int32(graphics.width) - 250, 75, 100, 100}
		graphics.gameRenderer.FillRect(&player2Rect)
	}
	graphics.gameRenderer.Present()
}

func calculateGridDim(windowWidth int, windowHeight int, boardWidth int, boardHeight int) int {
	maxTilesX := windowWidth / boardWidth
	maxTilesY := windowHeight / boardHeight
	if maxTilesX < maxTilesY {
		return maxTilesX
	}
	return maxTilesY
}
