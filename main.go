// My First Ever Go Project

package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MOVING = 1

func main() {

	tile_size := int32(16)
	columns := int32(20)
	Rows := int32(40)

	movement_speed := 5
	xOffset := int32(0)
	yOffset := int32(0)

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

	tileSet := rl.LoadTexture("assets/tiles.png")
	defer rl.UnloadTexture(tileSet)

	for !rl.WindowShouldClose() {

		leftCol, rightCol := getPieceBounds(incomingPiece)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Blue)

		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if incomingPiece[row][col] == MOVING {
					DrawTile(randomBlock, tileSet, int(xOffset)+col*16, int(yOffset)+row*16)

				}
			}
		}

		for y := int32(0); y < Rows; y++ {
			for x := int32(0); x < columns; x++ {
				rl.DrawRectangleLines(x*tile_size, y*tile_size, tile_size, tile_size, rl.Black)
			}
		}

		if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
			maxXOffset := (columns - 1 - int32(rightCol)) * tile_size
			xOffset += tile_size
			if xOffset > maxXOffset {
				xOffset = maxXOffset
			}
		}

		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
			minXOffset := int32(leftCol) * tile_size
			xOffset -= tile_size
			if xOffset < minXOffset {
				xOffset = minXOffset
			}
		}

		select {
		case <-ticker.C:
			yOffset += tile_size
			if yOffset >= 576 {
				yOffset = 576
				ticker.Stop()
			}

		default:
		}

		if rl.IsKeyDown(rl.KeyDown) || rl.IsKeyDown(rl.KeyS) {
			yOffset = yOffset + int32(movement_speed)
			if yOffset >= 576 {
				ticker.Stop()
				movement_speed = 0
				yOffset = 576
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

func DrawTile(tileid int, tileSet rl.Texture2D, x int, y int) {
	rl.DrawTextureRec(tileSet, TileRecFor(tileid), rl.NewVector2(float32(x), float32(y)), rl.White)
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
