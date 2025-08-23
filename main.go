// My First Ever Go Project (made with lots of caffeine (about 3 monster cans) and a bit of AI)

/*  Ideas:

1 - Press Enter To Play option

*/

package main

import (
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MOVING = 1

type RowAnimation struct {
	row      int
	progress float32
}

type Particle struct {
	x, y     float32
	vx, vy   float32
	color    rl.Color
	alpha    float32
	lifetime float32
	duration float32
}

var particles []Particle

func main() {
	var currentColor rl.Color = rl.Blue
	var targetColor rl.Color = rl.Blue
	score := 0
	gameOver := false
	tile_size := int32(16)
	columns := int32(15)
	Rows := int32(30)

	startMenu := true
	var blinkVisible bool = true
	var blinkTimer float32 = 0

	animPiece := make([][]int, 4)
	for i := range animPiece {
		animPiece[i] = make([]int, 4)
	}

	animX := int32(3)
	animY := int32(0)
	animTicker := time.NewTicker(500 * time.Millisecond)
	defer animTicker.Stop()

	initTetr(int(rl.GetRandomValue(0, 6)), animPiece)

	pieceX := int32(0)
	pieceY := int32(0)
	loseBarrier := int32(2)
	level := 1 + score/500
	comboStreak := 0
	board := make([][]int, Rows)
	for i := range board {
		board[i] = make([]int, columns)
	}
	showGrid := true
	isPaused := false
	mute := false
	type floatingScore struct {
		x, y     int32
		points   int
		timer    float32
		duration float32
	}
	var floatingScores []floatingScore
	baseDrop := time.Second
	speedIncrease := 100 * time.Millisecond
	dropInterval := baseDrop - time.Duration(level-1)*speedIncrease
	if dropInterval < 100*time.Millisecond {
		dropInterval = 100 * time.Millisecond
	}
	ticker := time.NewTicker(dropInterval)
	defer ticker.Stop()

	rl.SetRandomSeed(uint32(time.Now().UnixNano()))
	random := int(rl.GetRandomValue(0, 6))
	randomBlock := random

	incomingPiece := make([][]int, 4)
	for i := range incomingPiece {
		incomingPiece[i] = make([]int, 4)
	}

	var levelColors = []rl.Color{
		rl.Blue,
		rl.Orange,
		rl.Purple,
		rl.Red,
		rl.SkyBlue,
		rl.Yellow,
	}

	var shakeDuration float32
	var shakeIntensity int32 = 5
	factor := shakeDuration / 0.3
	var shakeX = int32(float32(rl.GetRandomValue(-shakeIntensity, shakeIntensity)) * factor)
	var shakeY = int32(float32(rl.GetRandomValue(-shakeIntensity, shakeIntensity)) * factor)

	initTetr(random, incomingPiece)

	const hudHeight = 100
	rl.InitWindow(int32(columns*tile_size), int32(Rows*tile_size)+hudHeight, "GoTetris")
	defer rl.CloseWindow()
	FPS := 60

	rl.SetTargetFPS(int32(FPS))

	tileSet := rl.LoadTexture("assets/tiles.png")
	defer rl.UnloadTexture(tileSet)

	icon := rl.LoadImage("assets/logo.png")
	rl.SetWindowIcon(*icon)
	rl.UnloadImage(icon)
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	bgAudio := rl.LoadMusicStream("assets/tetris.mp3")
	rl.PlayMusicStream(bgAudio)
	rl.SetMusicVolume(bgAudio, 0.1)
	defer rl.UnloadMusicStream(bgAudio)

	scoreSound := rl.LoadSound("assets/collect-points-190037.mp3")
	rl.SetSoundVolume(scoreSound, 1)

	defer rl.UnloadSound(scoreSound)

	var animateRows []RowAnimation

	for !rl.WindowShouldClose() {
		if startMenu {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			rl.DrawText("GoTetris", 50, 100+hudHeight, 30, rl.White)
			for row := 0; row < 4; row++ {
				for col := 0; col < 4; col++ {
					if animPiece[row][col] == MOVING {
						DrawTile(int(rl.GetRandomValue(0, 6)), tileSet,
							(animX+int32(col))*tile_size, (animY+int32(row))*tile_size+hudHeight)
					}
				}
			}

			if blinkVisible {

				rl.DrawText("Press ENTER To Start", 33, 150+hudHeight, 15, rl.White)
			}

			rl.EndDrawing()

			select {
			case <-animTicker.C:
				animY++
				if animY > Rows-4 {
					animY = 0
					initTetr(int(rl.GetRandomValue(0, 6)), animPiece)
				}
			default:
			}

			blinkTimer += 1.0 / float32(FPS)
			if blinkTimer >= 0.5 {
				blinkVisible = !blinkVisible
				blinkTimer = 0
			}

			if rl.IsKeyPressed(rl.KeyEnter) {
				startMenu = false
			}
			continue
		}

		if rl.IsKeyPressed(rl.KeyM) {
			mute = !mute
			if mute {
				rl.PauseMusicStream(bgAudio)
				fmt.Println("paused")
			} else {
				rl.ResumeMusicStream(bgAudio)
				fmt.Println("unpaused")
			}
		}

		rl.UpdateMusicStream(bgAudio)

		ghostX := pieceX
		ghostY := pieceY
		for canPlace(incomingPiece, board, ghostX, ghostY+1, columns, Rows) {
			ghostY++
		}

		ghostBlock := make([][]int, 4)
		for i := range ghostBlock {
			ghostBlock[i] = make([]int, 4)
		}

		// // Lose Screen

		if loss(board, loseBarrier) {
			incomingPiece, pieceX, pieceY, randomBlock = spawnPiece()

			if !canPlace(incomingPiece, board, pieceX, pieceY, columns, Rows) {
				gameOver = true
				if gameOver {
					rl.BeginDrawing()
					rl.DrawText("Game Over!", 25, 160+hudHeight, 35, rl.White)
					rl.DrawText("Press ESC To Play Again", 17, 200+hudHeight, 17, rl.White)
					rl.DrawText(fmt.Sprintf("Score: %d", score), 65, 225+hudHeight, 25, rl.White)
					rl.EndDrawing()

					if rl.IsKeyPressed(rl.KeyEscape) {
						board, incomingPiece, pieceX, pieceY, randomBlock, score = resetGame(columns, Rows)
						gameOver = false
					}
				}
			}
		}

		if isPaused {
			rl.DrawRectangle(0, 0, columns*tile_size, Rows*tile_size, rl.NewColor(0, 0, 0, 10))
			rl.DrawText("Paused", 45, 200, 40, rl.White)
			rl.EndDrawing()
		}

		// Pause

		if !isPaused {
			if rl.IsKeyPressed(rl.KeyP) {
				isPaused = true
				fmt.Println("working2")
			}
		} else if isPaused {
			if rl.IsKeyPressed(rl.KeyP) {
				isPaused = false
				fmt.Println("working3")
			}
		}

		if !gameOver && !isPaused {
			locked := false
			leftCol, rightCol := getPieceBounds(incomingPiece)

			rl.BeginDrawing()

			// Floating Score
			for i := 0; i < len(floatingScores); i++ {
				fs := &floatingScores[i]
				alpha := uint8(255 * (1 - fs.timer/fs.duration))
				rl.DrawText(fmt.Sprintf("+%d", fs.points), fs.x, fs.y+hudHeight, 24, rl.NewColor(255, 255, 255, alpha))
				t := fs.timer / fs.duration

				fs.y -= int32(1 + 5*(1-t))
				fs.x += int32(math.Sin(float64(t*math.Pi*2)) * 2)
				fs.timer += 1.0 / float32(FPS)

				if fs.timer >= fs.duration {
					floatingScores = append(floatingScores[:i], floatingScores[i+1:]...)
					i--
				}
			}

			// Main Tetris Blocks

			for y := int32(0); y < Rows; y++ {
				for x := int32(0); x < columns; x++ {
					if board[y][x] != 0 {
						DrawTile(board[y][x]-1, tileSet, x*tile_size+shakeX, y*tile_size+hudHeight+shakeX)
					}
				}
			}

			// Shake
			if shakeDuration > 0 {
				shakeDuration -= 1.0 / float32(FPS)
				shakeX = int32(rl.GetRandomValue(-shakeIntensity, shakeIntensity))
				shakeY = int32(rl.GetRandomValue(-shakeIntensity, shakeIntensity))
			} else {
				shakeX = 0
				shakeY = 0
			}

			// Particles
			updateParticles(1.0 / float32(FPS))
			drawParticles(tile_size)

			// level / score

			newLevel := 1 + score/500
			if newLevel != level {
				level = newLevel
				dropInterval = baseDrop - time.Duration(level-1)*speedIncrease
				if dropInterval < 100*time.Millisecond {
					dropInterval = 100 * time.Millisecond
				}
				ticker.Stop()
				ticker = time.NewTicker(dropInterval)
			}

			// Animation

			for i := 0; i < len(animateRows); i++ {
				a := &animateRows[i]
				alpha := uint8(255 * (1 - a.progress))
				for x := 0; x < int(columns); x++ {
					if board[a.row][x] != 0 {
						rl.DrawRectangle(int32(x)*tile_size, int32(a.row)*tile_size, tile_size, tile_size, rl.NewColor(255, 255, 255, alpha))
					}
				}
				a.progress += 0.05
				if a.progress >= 1 {
					removeRows(board, a.row)
					animateRows = append(animateRows[:i], animateRows[i+1:]...)
					i--
				}
			}

			for row := 0; row < 4; row++ {
				for col := 0; col < 4; col++ {
					if incomingPiece[row][col] == MOVING {
						DrawTile(randomBlock, tileSet, (pieceX+int32(col))*tile_size+shakeX, (pieceY+int32(row))*tile_size+hudHeight+shakeY)

					}
				}
			}

			// BG Change

			targetColor = levelColors[(level-1)%len(levelColors)]

			t := float32(0.05)
			currentColor.R = lerp(currentColor.R, targetColor.R, t)
			currentColor.G = lerp(currentColor.G, targetColor.G, t)
			currentColor.B = lerp(currentColor.B, targetColor.B, t)
			currentColor.A = 255

			rl.ClearBackground(currentColor)

			// Grid

			if showGrid {
				for y := int32(0); y < Rows; y++ {
					for x := int32(0); x < columns; x++ {
						rl.DrawRectangleLines(x*tile_size+shakeX, y*tile_size+hudHeight+shakeY, tile_size, tile_size, rl.Black)
					}
				}
			}

			// Ghost Blocks

			for row := 0; row < 4; row++ {
				for col := 0; col < 4; col++ {
					if incomingPiece[row][col] == MOVING {
						rl.DrawRectangleLines((ghostX+int32(col))*tile_size, (ghostY+int32(row))*tile_size+hudHeight, tile_size, tile_size, rl.White)

					}
				}
			}

			// HUD
			textColor := getTextColor(currentColor)
			rl.DrawText(fmt.Sprintf("Score: %d", score), 10, 10, 20, textColor)
			rl.DrawText(fmt.Sprintf("Level: %d", level), 10, 40, 20, textColor)
			rl.DrawText(fmt.Sprintf("Combo: %d", comboStreak), 10, 70, 20, textColor)

			// KEYPRESS SECTION

			// Right Movement
			if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
				if pieceX < columns-1-int32(rightCol) {
					if canPlace(incomingPiece, board, pieceX+1, pieceY, columns, Rows) {
						pieceX++
					}
				}

			}

			// Left Movement
			if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {

				if pieceX > 0-int32(leftCol) {
					if canPlace(incomingPiece, board, pieceX-1, pieceY, columns, Rows) {
						pieceX--
					}
				}
			}

			// Soft Drop

			if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
				bottomRow := getPieceBottom(incomingPiece)
				if pieceY < Rows-0-int32(bottomRow) && canPlace(incomingPiece, board, pieceX, pieceY+1, columns, Rows) {
					pieceY++
				} else {
					lockPiece(incomingPiece, board, pieceX, pieceY, randomBlock)
					newAnimations := clearRow(board)
					cleared := len(newAnimations)
					animateRows = append(animateRows, newAnimations...)

					if cleared == 4 {
						shakeDuration = 0.3
					}

					if cleared > 0 {
						topRow := newAnimations[0].row
						pts := pointsForCleared(cleared)
						score += pointsForCleared(cleared)
						rl.PlaySound(scoreSound)
						floatingScores = append(floatingScores, floatingScore{
							x:        (columns/2 - 1) * tile_size,
							y:        int32(topRow) * tile_size,
							points:   pts,
							timer:    0,
							duration: 1.0,
						})

						comboStreak++
						if comboStreak > 1 {
							comboPoints := 50 * (comboStreak - 1)
							score += comboPoints
							floatingScores = append(floatingScores, floatingScore{
								x:        (columns/2 - 1) * tile_size,
								y:        (int32(topRow))*tile_size - 20,
								points:   comboPoints,
								timer:    0,
								duration: 1.0,
							})
						}

						if isBoardEmpty(board) {
							score += 1000
							floatingScores = append(floatingScores, floatingScore{
								x:        int32(columns/2-2) * tile_size,
								y:        int32(Rows/2) * tile_size,
								points:   1000,
								timer:    0,
								duration: 1.5,
							})
						}

					} else {
						comboStreak = 0
					}
					incomingPiece, pieceX, pieceY, randomBlock = spawnPiece()
				}

			}

			// Rotation
			if !locked {
				if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
					rotated := rotatePiece(incomingPiece)
					if canPlace(rotated, board, pieceX, pieceY, columns, Rows) {
						incomingPiece = rotated
					}
				}
			}

			// Show/Hide Grid

			if showGrid {
				if rl.IsKeyPressed(rl.KeyG) {
					showGrid = false
				}
			} else if !showGrid {
				if rl.IsKeyPressed(rl.KeyG) {
					showGrid = true
				}
			}

			// Instant Drop
			if rl.IsKeyPressed(rl.KeySpace) {
				for canPlace(incomingPiece, board, pieceX, pieceY+1, columns, Rows) {
					pieceY++
				}
			}

			// Drop 1 block per second.

			select {
			case <-ticker.C:
				bottomRow := getPieceBottom(incomingPiece)
				if pieceY < Rows-0-int32(bottomRow) && canPlace(incomingPiece, board, pieceX, pieceY+1, columns, Rows) {
					pieceY++
				} else {
					lockPiece(incomingPiece, board, pieceX, pieceY, randomBlock)
					newAnimations := clearRow(board)
					animateRows = append(animateRows, newAnimations...)
					cleared := len(newAnimations)
					if cleared == 4 {
						shakeDuration = 0.3
					}
					if cleared > 0 {
						topRow := newAnimations[0].row
						pts := pointsForCleared(cleared)
						score += pointsForCleared(cleared)
						rl.PlaySound(scoreSound)
						spawnParticles(columns, topRow, tile_size, rl.White, 20)

						floatingScores = append(floatingScores, floatingScore{
							x:        (columns/2 - 1) * tile_size,
							y:        int32(topRow) * tile_size,
							points:   pts,
							timer:    0,
							duration: 1.0,
						})
						comboStreak++
						if comboStreak > 1 {
							comboPoints := 50 * (comboStreak - 1)
							score += comboPoints
							floatingScores = append(floatingScores, floatingScore{
								x:        (columns/2 - 1) * tile_size,
								y:        int32(topRow) * tile_size,
								points:   comboPoints,
								timer:    0,
								duration: 1.0,
							})
						}
						if isBoardEmpty(board) {
							score += 1000
							floatingScores = append(floatingScores, floatingScore{
								x:        int32(columns/2-2) * tile_size,
								y:        int32(Rows/2) * tile_size,
								points:   1000,
								timer:    0,
								duration: 1.5,
							})
						}
					} else {
						comboStreak = 0
					}
					incomingPiece, pieceX, pieceY, randomBlock = spawnPiece()
				}

			default:

			}

			rl.EndDrawing()

		}

	}
}

// FUNCTIONS SECTION

func TileRecFor(tileid int) rl.Rectangle {
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
				x := xOffset + int32(col)
				y := yOffset + int32(row)

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

func lockPiece(piece [][]int, board [][]int, pieceX, pieceY int32, tileid int) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if piece[row][col] == MOVING {
				x := pieceX + int32(col)
				y := pieceY + int32(row)
				if y >= 0 && y < int32(len(board)) && x >= 0 && x < int32(len(board[0])) {
					board[y][x] = tileid + 1
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

	return newPiece, 3, 0, random
}

func loss(board [][]int, barrier int32) bool {
	for row := int32(0); row < barrier; row++ {
		for col := 0; col < len(board[0]); col++ {
			if board[row][col] != 0 {
				return true
			}
		}
	}
	return false

}

func resetGame(columns, Rows int32) ([][]int, [][]int, int32, int32, int, int) {
	score := 0

	board := make([][]int, Rows)
	for i := range board {
		board[i] = make([]int, columns)
	}

	newPiece := make([][]int, 4)
	for i := range newPiece {
		newPiece[i] = make([]int, 4)
	}
	random := int(rl.GetRandomValue(0, 6))
	initTetr(random, newPiece)

	pieceX, pieceY := int32(3), int32(0)
	return board, newPiece, pieceX, pieceY, random, score
}

func clearRow(board [][]int) []RowAnimation {

	rows := len(board)
	cols := len(board[0])
	animations := []RowAnimation{}

	for row := 0; row < rows; row++ {
		full := true
		for col := 0; col < cols; col++ {
			if board[row][col] == 0 {
				full = false
				break
			}
		}
		if full {
			animations = append(animations, RowAnimation{row: row, progress: 0})
		}
	}
	return animations
}

func pointsForCleared(rowsCleared int) int {
	switch rowsCleared {
	case 1:
		return 100
	case 2:
		return 300
	case 3:
		return 500
	case 4:
		return 800
	default:
		return 0
	}
}

func removeRows(board [][]int, rowToRemove int) {
	cols := len(board[0])
	for y := rowToRemove; y > 0; y-- {
		copy(board[y], board[y-1])
	}
	board[0] = make([]int, cols)
}

func isBoardEmpty(board [][]int) bool {
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[0]); x++ {
			if board[y][x] != 0 {
				return false
			}
		}
	}
	return true
}

func spawnParticles(columns int32, row int, tileSize int32, color rl.Color, amount int) {
	for i := 0; i < amount; i++ {
		p := Particle{
			x:        float32(rl.GetRandomValue(0, int32(columns*tileSize-tileSize))),
			y:        float32(row) * float32(tileSize),
			vx:       float32(rl.GetRandomValue(-2, 2)),
			vy:       float32(rl.GetRandomValue(-5, -1)),
			color:    color,
			alpha:    1.0,
			lifetime: 0,
			duration: 1.0,
		}
		particles = append(particles, p)
	}
}

func updateParticles(dt float32) {
	for i := 0; i < len(particles); i++ {
		p := &particles[i]
		p.x += p.vx
		p.y += p.vy
		p.lifetime += dt
		p.alpha = 1 - p.lifetime/p.duration

		if p.lifetime >= p.duration {
			particles = append(particles[:i], particles[i+1:]...)
			i--
		}
	}
}

func drawParticles(tileSize int32) {
	for _, p := range particles {
		c := p.color
		c.A = uint8(p.alpha * 255)
		rl.DrawRectangle(int32(p.x), int32(p.y), tileSize/2, tileSize/2, c)
	}
}

func getTextColor(bg rl.Color) rl.Color {
	brightness := 0.299*float32(bg.R) + 0.587*float32(bg.G) + 0.114*float32(bg.B)
	if brightness > 150 {
		return rl.Black
	}
	return rl.White
}

func lerp(a, b uint8, t float32) uint8 {
	return uint8(float32(a)*(1-t) + float32(b)*t)
}
