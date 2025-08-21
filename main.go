// My First Ever Go Project

package main

import rl "github.com/gen2brain/raylib-go/raylib"

// a tetromino is a 4x4 grid represented by a 2d array
type Tetromino [2][6]uint8

// 7 Tetrominoes
var Tetrominoes [7][4]Tetromino

func main() {
	rl.InitWindow(450, 800, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		tileSet := rl.LoadTexture("assets/tiles.png")
		rl.DrawTexture(tileSet, 100, 200, rl.White)
		for i := range 7 {
			DrawTile(i, tileSet, i*16, 0)
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

func initTetr() {
	Tetrominoes = Tetromino{
		{}, //I
		{}, //Z
		{}, //S
		{}, //T
		{}, //L
		{}, //J
		{}  //O

	
	}
}
