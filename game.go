package main

import (
	"math/rand"
)

type game struct {
	gameConstants    gameConstants
	gameStage        gameStage
	currentPlayer    player
	evaluate         bool
	currentSelection [2]int
	movedDot         [2]int
	board            [][]tile
}

type player int

const (
	Player1 player = iota
	Player2
)

type gameConstants struct {
	boardWidth  int // number of grid tiles wide
	boardHeight int // number of grid tiles tall
	dotCoverage float32
}

type tile int

const (
	Empty tile = iota
	Dot
	Claimed1
	Claimed2
)

type gameStage int

const (
	Init gameStage = iota
	Running
	Finished
)

func newGame(boardWidth int, boardHeight int, dotCoverage float32) *game {
	gc := gameConstants{
		boardWidth:  boardWidth,
		boardHeight: boardHeight,
		dotCoverage: dotCoverage,
	}
	g := game{
		gameConstants:    gc,
		gameStage:        Init,
		evaluate:         false,
		currentSelection: [2]int{-1, -1},
	}
	g.board = make([][]tile, boardHeight)
	for i := 0; i < boardHeight; i++ {
		g.board[i] = make([]tile, boardWidth)
	}
	return &g
}

func (g *game) getGameConstants() *gameConstants {
	return &g.gameConstants
}

func (g *game) cleanup() {
	// defers here
}

func (g *game) init() {
	g.genDots()
}

func (g *game) startGame() {
	g.gameStage = Running
}

func (g *game) genDots() {
	for i := 0; i < g.gameConstants.boardHeight; i++ {
		for j := 0; j < g.gameConstants.boardWidth; j++ {
			if rand.Float32() < g.gameConstants.dotCoverage {
				g.board[i][j] = Dot
			}
		}
	}
	g.removeSquares()
}

func (g *game) removeSquares() {
	madeChanges := false
	for i := 0; i < g.gameConstants.boardHeight-1; i++ {
		for j := 0; j < g.gameConstants.boardWidth-1; j++ {
			if g.board[i][j] == Dot && g.board[i+1][j] == Dot && g.board[i][j+1] == Dot && g.board[i+1][j+1] == Dot {
				madeChanges = true
				switch rand.Intn(4) {
				case 0:
					g.board[i][j] = Empty
					break
				case 1:
					g.board[i+1][j] = Empty
					break
				case 2:
					g.board[i][j+1] = Empty
					break
				case 3:
					g.board[i+1][j+1] = Empty
					break
				}
			}
		}
	}
	if madeChanges {
		g.removeSquares()
	}
}

func (g *game) getTileAt(x int, y int) *tile {
	return &g.board[x][y]
}

// func (g *game) update() {
// 	madeChanges := false
// 	for i := 0; i < g.gameConstants.boardHeight-1; i++ {
// 		for j := 0; j < g.gameConstants.boardWidth-1; j++ {
// 			if g.board[i][j] == Dot && g.board[i+1][j] == Dot && g.board[i][j+1] == Dot && g.board[i+1][j+1] == Dot {
// 				madeChanges = true
// 				// temporary fix, player/claim is inverted because current player was already updated
// 				if g.currentPlayer == Player1 {
// 					g.board[i][j] = Claimed2
// 					g.board[i+1][j] = Claimed2
// 					g.board[i][j+1] = Claimed2
// 					g.board[i+1][j+1] = Claimed2
// 				} else if g.currentPlayer == Player2 {
// 					g.board[i][j] = Claimed1
// 					g.board[i+1][j] = Claimed1
// 					g.board[i][j+1] = Claimed1
// 					g.board[i+1][j+1] = Claimed1
// 				} else {
// 					panic("lol")
// 				}
// 			}
// 		}
// 	}
// 	if madeChanges {
// 		g.update()
// 	}
// }

func (g *game) update(i int, j int) {
	toClaim := [][2]int{}
	newDots := [][2]int{}

	if g.board[i][j] != Dot {
		println(g.board[i][j])
		panic("movedDot != dot")
	}
	// check right-down box
	// x o
	// o o
	if len(g.board) > i+1 && len(g.board[0]) > j+1 && g.board[i+1][j] == Dot && g.board[i][j+1] == Dot && g.board[i+1][j+1] == Dot {
		toClaim = append(toClaim, [2]int{i, j})
		toClaim = append(toClaim, [2]int{i + 1, j})
		toClaim = append(toClaim, [2]int{i, j + 1})
		toClaim = append(toClaim, [2]int{i + 1, j + 1})
		if g.currentPlayer == Player1 {
			newDots = append(newDots, [2]int{i - 1, j + 1})
		} else if g.currentPlayer == Player2 {
			newDots = append(newDots, [2]int{i + 1, j - 1})
		}
	}
	// check left-down box
	// o x
	// o o
	if len(g.board) > i+1 && j-1 >= 0 && g.board[i+1][j] == Dot && g.board[i][j-1] == Dot && g.board[i+1][j-1] == Dot {
		toClaim = append(toClaim, [2]int{i, j})
		toClaim = append(toClaim, [2]int{i + 1, j})
		toClaim = append(toClaim, [2]int{i, j - 1})
		toClaim = append(toClaim, [2]int{i + 1, j - 1})
		if g.currentPlayer == Player1 {
			newDots = append(newDots, [2]int{i + 1, j + 1})
		} else if g.currentPlayer == Player2 {
			newDots = append(newDots, [2]int{i - 1, j - 1})
		}
	}
	// check right-up box
	// o o
	// x o
	if i-1 >= 0 && len(g.board[0]) > j+1 && g.board[i-1][j] == Dot && g.board[i][j+1] == Dot && g.board[i-1][j+1] == Dot {
		toClaim = append(toClaim, [2]int{i, j})
		toClaim = append(toClaim, [2]int{i - 1, j})
		toClaim = append(toClaim, [2]int{i, j + 1})
		toClaim = append(toClaim, [2]int{i - 1, j + 1})
		if g.currentPlayer == Player1 {
			newDots = append(newDots, [2]int{i - 1, j - 1})
		} else if g.currentPlayer == Player2 {
			newDots = append(newDots, [2]int{i + 1, j + 1})
		}
	}
	// check left-up box
	// o o
	// o x
	if i-1 >= 0 && j-1 >= 0 && g.board[i-1][j] == Dot && g.board[i][j-1] == Dot && g.board[i-1][j-1] == Dot {
		toClaim = append(toClaim, [2]int{i, j})
		toClaim = append(toClaim, [2]int{i - 1, j})
		toClaim = append(toClaim, [2]int{i, j - 1})
		toClaim = append(toClaim, [2]int{i - 1, j - 1})
		if g.currentPlayer == Player1 {
			newDots = append(newDots, [2]int{i + 1, j - 1})
		} else if g.currentPlayer == Player2 {
			newDots = append(newDots, [2]int{i - 1, j + 1})
		}
	}

	for i := 0; i < len(newDots); i++ {
		if g.board[newDots[i][0]][newDots[i][1]] == Empty {
			println("updating from placing new dot on empty")
			g.board[newDots[i][0]][newDots[i][1]] = Dot
			g.update(newDots[i][0], newDots[i][1])
		}
		if g.currentPlayer == Player1 {
			if g.board[newDots[i][0]][newDots[i][1]] == Claimed2 {
				g.board[newDots[i][0]][newDots[i][1]] = Empty
			} else if g.board[newDots[i][0]][newDots[i][1]] == Claimed1 {
				println("updating from placing new dot on claimed")
				g.board[newDots[i][0]][newDots[i][1]] = Dot
				g.update(newDots[i][0], newDots[i][1])
			}
		} else if g.currentPlayer == Player2 {
			if g.board[newDots[i][0]][newDots[i][1]] == Claimed2 {
				println("updating from placing new dot on claimed")
				g.board[newDots[i][0]][newDots[i][1]] = Dot
				g.update(newDots[i][0], newDots[i][1])
			} else if g.board[newDots[i][0]][newDots[i][1]] == Claimed1 {
				g.board[newDots[i][0]][newDots[i][1]] = Empty
			}
		}
	}

	for i := 0; i < len(toClaim); i++ {
		if g.currentPlayer == Player1 {
			g.board[toClaim[i][0]][toClaim[i][1]] = Claimed2
		} else if g.currentPlayer == Player2 {
			g.board[toClaim[i][0]][toClaim[i][1]] = Claimed1
		}
	}
}

func (g *game) makeMove(prevTileX int, prevTileY int, newTileX int, newTileY int) {
	if g.board[prevTileX][prevTileY] != Dot || g.board[newTileX][newTileY] != Empty {
		return
	}
	g.board[prevTileX][prevTileY] = Empty
	g.board[newTileX][newTileY] = Dot
	g.movedDot = [2]int{newTileX, newTileY}
	if g.currentPlayer == Player1 {
		g.currentPlayer = Player2
	} else {
		g.currentPlayer = Player1
	}
	g.evaluate = true
}
