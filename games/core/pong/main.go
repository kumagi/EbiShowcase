package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 480
	screenHeight = 720
	paddleHalfW  = 60
	ballRadius   = 12
)

type game struct {
	playerX     float64
	cpuX        float64
	ballX       float64
	ballY       float64
	ballVX      float64 // ボールの横の速さ
	ballVY      float64 // ボールの縦の速さ
	playerScore int
	cpuScore    int
}

func newGame() *game {
	g := &game{
		playerX: screenWidth / 2,
		cpuX:    screenWidth / 2,
	}
	g.serve(1) // 下向きにサーブ
	return g
}

// direction: +1 は下へ、-1 は上へ
func (g *game) serve(direction float64) {
	g.ballX = screenWidth / 2
	g.ballY = screenHeight / 2
	g.ballVX = 3.2
	g.ballVY = 4.2 * direction
}

func (g *game) movePlayer() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.playerX -= 6
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerX += 6
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		g.playerX = float64(x)
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		g.playerX = float64(x)
	}
	g.playerX = math.Max(55, math.Min(screenWidth-55, g.playerX))
}

// CPU はボールの X にゆっくりついていく
func (g *game) moveCPU() {
	if g.cpuX < g.ballX-8 {
		g.cpuX += 3.5
	} else if g.cpuX > g.ballX+8 {
		g.cpuX -= 3.5
	}
}

func (g *game) bounceBall() {
	g.ballX += g.ballVX
	g.ballY += g.ballVY

	// 左右の壁
	if g.ballX < 14 || g.ballX > screenWidth-14 {
		g.ballVX = -g.ballVX
		g.ballX = math.Max(14, math.Min(screenWidth-14, g.ballX))
	}

	// 自分のパドル（下）
	hitPlayer := g.ballVY > 0 &&
		g.ballY >= 646 && g.ballY <= 662 &&
		math.Abs(g.ballX-g.playerX) < 66
	if hitPlayer {
		g.ballY = 646
		g.ballVY = -math.Abs(g.ballVY) * 1.025
		g.ballVX += (g.ballX - g.playerX) * 0.055
	}

	// CPU のパドル（上）
	hitCPU := g.ballVY < 0 &&
		g.ballY <= 74 && g.ballY >= 58 &&
		math.Abs(g.ballX-g.cpuX) < 66
	if hitCPU {
		g.ballY = 74
		g.ballVY = math.Abs(g.ballVY) * 1.015
		g.ballVX += (g.ballX - g.cpuX) * 0.04
	}

	g.ballVX = math.Max(-8, math.Min(8, g.ballVX))
}

// --- ここから Update ---
func (g *game) Update() error {
	g.movePlayer()
	g.moveCPU()
	g.bounceBall()

	// 上に抜けたらプレイヤーの得点
	if g.ballY < -20 {
		g.playerScore++
		g.serve(1)
	}
	// 下に抜けたら CPU の得点
	if g.ballY > screenHeight+20 {
		g.cpuScore++
		g.serve(-1)
	}
	return nil
}

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 20, 38, 255})

	// 中央の点線
	for y := 10; y < screenHeight; y += 30 {
		vector.DrawFilledRect(screen, screenWidth/2-2, float32(y), 4, 16, color.RGBA{45, 226, 194, 80}, false)
	}

	vector.DrawFilledRect(screen, float32(g.cpuX-paddleHalfW), 50, 120, 14, color.RGBA{255, 105, 79, 255}, false)
	vector.DrawFilledRect(screen, float32(g.playerX-paddleHalfW), 656, 120, 14, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledCircle(screen, float32(g.ballX), float32(g.ballY), ballRadius, color.White, false)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CPU %02d", g.cpuScore), 24, 22)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("YOU %02d", g.playerScore), 400, 688)
	ebitenutil.DebugPrintAt(screen, "ARROWS / DRAG", 184, 350)
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
