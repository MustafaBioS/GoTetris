// My First Ever Go Project

package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MOVING = 1

func main() {
	score := 0
	tile_size := int32(16)
	columns := int32(20)
	Rows := int32(40)

	pieceX := int32(0)
	pieceY := int32(0)

	board := make([][]int, Rows)
	for i := range board {
		board[i] = make([]int, columns)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	rl.SetRandomSeed(uint32(time.Now().UnixNano()))
	random := int(rl.GetRandomValue(0, 6))
	randomBlock := random

	incomingPiece := make([][]int, 4)
	for i := range incomingPiece {
		incomingPiece[i] = make([]int, 4)
	}

	initTetr(random, incomingPiece)

	rl.InitWindow(320, 640, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	rl.DrawText(fmt.Sprintf("Score: %d", score), 180, 500, 60, rl.White)

	tileSet := rl.LoadTexture("assets/tiles.png")
	defer rl.UnloadTexture(tileSet)

	for !rl.WindowShouldClose() {

		leftCol, rightCol := getPieceBounds(incomingPiece)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)

		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if incomingPiece[row][col] == MOVING {

					DrawTile(randomBlock, tileSet, (pieceX+int32(col))*tile_size, (pieceY+int32(row))*tile_size)

				}
			}
		}

		for y := int32(0); y < Rows; y++ {
			for x := int32(0); x < columns; x++ {
				rl.DrawRectangleLines(x*tile_size, y*tile_size, tile_size, tile_size, rl.Black)
			}
		}

		// Right Movement
		if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
			if pieceX < columns-1-int32(rightCol) {
				pieceX++
			}

		}

		// Left Movement
		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
			if pieceX > 0-int32(leftCol) {
				pieceX--
			}
		}
		select {
		case <-ticker.C:
			bottomRow := getPieceBottom(incomingPiece)
			if pieceY < Rows-1-int32(bottomRow) && canPlace(incomingPiece, board, pieceX, pieceY+1, columns, Rows) {
				pieceY++
			} else {
				lockPiece(incomingPiece, board, pieceX, pieceY)
				score++
				fmt.Print(score)

				incomingPiece, pieceX, pieceY, randomBlock = spawnPiece()

				ticker.Reset(time.Second)
			}

		default:

		}

		// Soft Drop

		if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
			bottomRow := getPieceBottom(incomingPiece)
			if pieceY < Rows-1-int32(bottomRow) {
				pieceY++
			} else {
				pieceY = Rows - 1 - int32(bottomRow)
				ticker.Stop()
			}
		}

		// Rotation

		if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
			rotated := rotatePiece(incomingPiece)
			incomingPiece = rotated
			if canPlace(rotated, board, pieceX, pieceY, columns, Rows) {
				incomingPiece = rotated
			}
		}
		rl.EndDrawing()
	}

}

func TileRecFor(tileid int) rl.Rectangle {
	//TODO
	if tileid > 6 {
		panic("tileid must be between 0 and 6")
	}
	return rl.NewRectangle(float32(tileid*16), 0, 16, 16)
}

func DrawTile(tileid int, tileSet rl.Texture2D, px, py int32) {
	rl.DrawTextureRec(tileSet, TileRecFor(tileid), rl.NewVector2(float32(px), float32(py)), rl.White)
}

func initTetr(random int, incomingPiece [][]int) {
	for i := range incomingPiece {
		for j := range incomingPiece[i] {
			incomingPiece[i][j] = 0
		}
	}

	switch random {
	case 0: // O
		incomingPiece[1][1] = MOVING
		incomingPiece[1][2] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[2][2] = MOVING

	case 1: // L
		incomingPiece[1][1] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[3][1] = MOVING
		incomingPiece[3][2] = MOVING

	case 2: // J
		incomingPiece[1][2] = MOVING
		incomingPiece[2][2] = MOVING
		incomingPiece[3][2] = MOVING
		incomingPiece[3][1] = MOVING

	case 3: // I
		incomingPiece[0][1] = MOVING
		incomingPiece[1][1] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[3][1] = MOVING

	case 4: // T
		incomingPiece[1][0] = MOVING
		incomingPiece[1][1] = MOVING
		incomingPiece[1][2] = MOVING
		incomingPiece[2][1] = MOVING

	case 5: // S
		incomingPiece[1][2] = MOVING
		incomingPiece[1][3] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[2][2] = MOVING

	case 6: // Z
		incomingPiece[1][1] = MOVING
		incomingPiece[1][2] = MOVING
		incomingPiece[2][2] = MOVING
		incomingPiece[2][3] = MOVING
	}

}

func getPieceBounds(piece [][]int) (left, right int) {
	left = 4
	right = -1

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if piece[row][col] == MOVING {
				if col < left {
					left = col
				}
				if col > right {
					right = col
				}
			}
		}

	}
	return
}

func rotatePiece(piece [][]int) [][]int {
	size := len(piece)
	rotated := make([][]int, size)
	for i := range rotated {
		rotated[i] = make([]int, size)
	}

	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			rotated[col][size-1-row] = piece[row][col]
		}
	}
	return rotated
}

func canPlace(piece [][]int, board [][]int, xOffset, yOffset int32, columns, rows int32) bool {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if piece[row][col] == MOVING {
				x := xOffset/int32(16) + int32(col)
				y := yOffset/int32(16) + int32(row)

				if x < 0 || x >= columns || y < 0 || y >= rows {
					return false
				}

				if board[y][x] != 0 {
					return false
				}
			}
		}
	}
	return true
}

func getPieceBottom(piece [][]int) int {
	bottom := -1
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if piece[row][col] == MOVING && row > bottom {
				bottom = row
			}
		}
	}
	return bottom
}

func lockPiece(piece [][]int, board [][]int, pieceX, pieceY int32) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if piece[row][col] == MOVING {
				x := pieceX + int32(col)
				y := pieceY + int32(row)
				if y >= 0 && y < int32(len(board)) && x >= 0 && x < int32(len(board[0])) {
					board[y][x] = 1
				}
			}
		}
	}
}

func spawnPiece() ([][]int, int32, int32, int) {
	newPiece := make([][]int, 4)
	for i := range newPiece {
		newPiece[i] = make([]int, 4)
	}
	random := int(rl.GetRandomValue(0, 6))
	initTetr(random, newPiece)
	return newPiece, 0, 0, random
}
