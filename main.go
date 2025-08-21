// My First Ever Go Project

package main

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MOVING = 1

func main() {

	rl.SetRandomSeed(uint32(time.Now().UnixNano()))
	random := int(rl.GetRandomValue(0, 6))
	randomBlock := random

	incomingPiece := make([][]int, 4)
	for i := range incomingPiece {
		incomingPiece[i] = make([]int, 4)
	}

	rl.InitWindow(450, 800, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	initTetr(random, incomingPiece)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		tileSet := rl.LoadTexture("assets/tiles.png")
		rl.DrawTexture(tileSet, 100, 200, rl.White)

		for i := 0; i < 7; i++ {
			DrawTile(i, tileSet, i*16, 0)
		}

		for row := 0; row < 4; row++ {
			for col := 0; col < 4; col++ {
				if incomingPiece[row][col] == MOVING {

					DrawTile(randomBlock, tileSet, 200+col*16, 100+row*16)
				}
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
	switch random {
	case 0: // O
		incomingPiece[1][1] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[1][2] = MOVING
		incomingPiece[2][2] = MOVING

	case 1: // L
		incomingPiece[1][0] = MOVING
		incomingPiece[1][1] = MOVING
		incomingPiece[1][2] = MOVING
		incomingPiece[2][2] = MOVING

	case 2: // J
		incomingPiece[1][2] = MOVING
		incomingPiece[2][0] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[2][2] = MOVING

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
		incomingPiece[1][1] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[2][2] = MOVING
		incomingPiece[3][2] = MOVING

	case 6: // Z
		incomingPiece[1][2] = MOVING
		incomingPiece[2][2] = MOVING
		incomingPiece[2][1] = MOVING
		incomingPiece[3][1] = MOVING
	}

}
