package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/lessonlogic"
)

const (
	screenWidth  = 480
	screenHeight = 720
	brickW       = 50.0
	brickH       = 24.0
	ballR        = 10.0
)

type brick struct {
	x, y  float64
	alive bool
	row   int
}

type game struct {
	paddleX float64
	ballX   float64
	ballY   float64
	ballVX  float64
	ballVY  float64
	bricks  []brick
	score   int
	lives   int
}

func newGame() *game {
	g := &game{
		paddleX: screenWidth / 2,
		lives:   3,
	}
	for row := 0; row < 6; row++ {
		for col := 0; col < 8; col++ {
			g.bricks = append(g.bricks, brick{
				x:     22 + float64(col)*56,
				y:     90 + float64(row)*32,
				alive: true,
				row:   row,
			})
		}
	}
	g.serve()
	return g
}

func (g *game) serve() {
	g.ballX = screenWidth / 2
	g.ballY = 590
	g.ballVX = 3.4
	g.ballVY = -4.5
}

func (g *game) movePaddle() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.paddleX -= 6
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.paddleX += 6
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		g.paddleX = float64(x)
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		g.paddleX = float64(x)
	}
	g.paddleX = math.Max(55, math.Min(screenWidth-55, g.paddleX))
}

func (g *game) bounceWorld() {
	g.ballX += g.ballVX
	g.ballY += g.ballVY

	// 左右の壁
	if g.ballX < 12 || g.ballX > screenWidth-12 {
		g.ballVX = -g.ballVX
		g.ballX = math.Max(12, math.Min(screenWidth-12, g.ballX))
	}
	// 天井
	if g.ballY < 52 {
		g.ballVY = math.Abs(g.ballVY)
	}
	// パドル
	hitPaddle := g.ballVY > 0 &&
		g.ballY > 630 && g.ballY < 650 &&
		math.Abs(g.ballX-g.paddleX) < 64
	if hitPaddle {
		g.ballY = 630
		g.ballVY = -math.Abs(g.ballVY)
		g.ballVX += (g.ballX - g.paddleX) * 0.06
	}
}

func (g *game) hitBricks() {
	for i := range g.bricks {
		b := &g.bricks[i]
		if !b.alive {
			continue
		}
		hit := g.ballX+ballR > b.x &&
			g.ballX-ballR < b.x+brickW &&
			g.ballY+ballR > b.y &&
			g.ballY-ballR < b.y+brickH
		if hit {
			b.alive = false
			g.score += 10
			g.ballVY = -g.ballVY
			break
		}
	}
}

func (g *game) aliveBrickCount() int {
	n := 0
	for _, b := range g.bricks {
		if b.alive {
			n++
		}
	}
	return n
}

// --- ここから Update ---
func (g *game) Update() error {
	g.movePaddle()
	g.bounceWorld()
	g.hitBricks()

	// ボールが下に落ちた
	if g.ballY > screenHeight+15 {
		var gameOver bool
		g.lives, gameOver = lessonlogic.SpendLife(g.lives)
		if gameOver {
			*g = *newGame()
		} else {
			g.serve()
		}
	}

	// ブロックを全部壊したらクリア → 最初から
	if g.aliveBrickCount() == 0 {
		*g = *newGame()
	}
	return nil
}

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})

	colors := []color.RGBA{
		{255, 105, 79, 255},
		{255, 145, 72, 255},
		{255, 210, 72, 255},
		{45, 226, 194, 255},
		{71, 161, 220, 255},
		{151, 113, 230, 255},
	}
	for _, b := range g.bricks {
		if !b.alive {
			continue
		}
		vector.DrawFilledRect(screen, float32(b.x), float32(b.y), brickW, brickH, colors[b.row], false)
	}

	vector.DrawFilledRect(screen, float32(g.paddleX-58), 640, 116, 14, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledCircle(screen, float32(g.ballX), float32(g.ballY), ballR, color.White, false)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %04d   LIVES %d", g.score, g.lives), 160, 25)
	ebitenutil.DebugPrintAt(screen, "ARROWS / DRAG", 184, 690)
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Breakout — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
