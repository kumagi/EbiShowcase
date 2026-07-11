package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type game struct {
	playerX, aiX float64
	ballX, ballY float64
	vx, vy       float64
	player, ai   int
}

func newGame() *game {
	g := &game{playerX: width / 2, aiX: width / 2}
	g.serve(1)
	return g
}

func (g *game) serve(direction float64) {
	g.ballX, g.ballY = width/2, height/2
	g.vx, g.vy = 3.2, 4.2*direction
}

func (g *game) Update() error {
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
	g.playerX = math.Max(55, math.Min(width-55, g.playerX))
	if g.aiX < g.ballX-8 {
		g.aiX += 3.5
	} else if g.aiX > g.ballX+8 {
		g.aiX -= 3.5
	}

	g.ballX += g.vx
	g.ballY += g.vy
	if g.ballX < 14 || g.ballX > width-14 {
		g.vx = -g.vx
		g.ballX = math.Max(14, math.Min(width-14, g.ballX))
	}
	if g.vy > 0 && g.ballY >= 646 && g.ballY <= 662 && math.Abs(g.ballX-g.playerX) < 66 {
		g.ballY, g.vy = 646, -math.Abs(g.vy)*1.025
		g.vx += (g.ballX - g.playerX) * .055
	}
	if g.vy < 0 && g.ballY <= 74 && g.ballY >= 58 && math.Abs(g.ballX-g.aiX) < 66 {
		g.ballY, g.vy = 74, math.Abs(g.vy)*1.015
		g.vx += (g.ballX - g.aiX) * .04
	}
	g.vx = math.Max(-8, math.Min(8, g.vx))
	if g.ballY < -20 {
		g.player++
		g.serve(1)
	}
	if g.ballY > height+20 {
		g.ai++
		g.serve(-1)
	}
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{7, 20, 38, 255})
	for y := 10; y < height; y += 30 {
		vector.DrawFilledRect(s, width/2-2, float32(y), 4, 16, color.RGBA{45, 226, 194, 80}, false)
	}
	vector.DrawFilledRect(s, float32(g.aiX-60), 50, 120, 14, color.RGBA{255, 105, 79, 255}, false)
	vector.DrawFilledRect(s, float32(g.playerX-60), 656, 120, 14, color.RGBA{45, 226, 194, 255}, false)
	vector.DrawFilledCircle(s, float32(g.ballX), float32(g.ballY), 12, color.White, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("CPU %02d", g.ai), 24, 22)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("YOU %02d", g.player), 400, 688)
	ebitenutil.DebugPrintAt(s, "ARROWS / DRAG", 184, 350)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Pong — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
