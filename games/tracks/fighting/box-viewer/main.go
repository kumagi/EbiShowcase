package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const width, height = 480, 720

type rect struct{ x, y, w, h float32 }
type game struct {
	p1, p2       float32
	attack, hits int
	clear        bool
}

func newGame() *game { return &game{p1: 100, p2: 340} }
func (g *game) Update() error {
	if g.clear {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	left := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	right := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	start := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > height/2 {
			if x < width/2 {
				left = true
			} else {
				right = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		_, y := ebiten.TouchPosition(id)
		if y < height/2 {
			start = true
		}
	}
	if left {
		g.p1 -= 3
	}
	if right {
		g.p1 += 3
	}
	g.p1 = max(float32(25), min(float32(430), g.p1))
	push1, push2 := rect{g.p1 - 18, 510, 36, 100}, rect{g.p2 - 18, 510, 36, 100}
	if overlap(push1, push2) {
		g.p1 = g.p2 - 36
	}
	if start && g.attack == 0 {
		g.attack = 20
	}
	if g.attack > 0 {
		g.attack--
		if g.attack == 12 {
			hit := rect{g.p1 + 15, 535, 90, 32}
			hurt := rect{g.p2 - 24, 500, 48, 105}
			if overlap(hit, hurt) {
				g.hits++
				if g.hits >= 5 {
					g.clear = true
				}
			}
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{20, 29, 45, 255})
	vector.DrawFilledRect(s, 0, 610, 480, 110, color.RGBA{54, 66, 79, 255}, false)
	drawFighter(s, g.p1, color.RGBA{45, 225, 194, 255})
	drawFighter(s, g.p2, color.RGBA{240, 75, 91, 255})
	vector.StrokeRect(s, g.p1-18, 510, 36, 100, 3, color.RGBA{90, 150, 255, 255}, false)
	vector.StrokeRect(s, g.p1-24, 500, 48, 105, 3, color.RGBA{70, 240, 110, 255}, false)
	vector.StrokeRect(s, g.p2-18, 510, 36, 100, 3, color.RGBA{90, 150, 255, 255}, false)
	vector.StrokeRect(s, g.p2-24, 500, 48, 105, 3, color.RGBA{70, 240, 110, 255}, false)
	if g.attack > 8 && g.attack < 16 {
		vector.StrokeRect(s, g.p1+15, 535, 90, 32, 4, color.RGBA{255, 75, 75, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("HITS %d/5", g.hits), 205, 35)
	ebitenutil.DebugPrintAt(s, "BLUE: PUSH   GREEN: HURT   RED: ATTACK", 100, 65)
	ebitenutil.DebugPrintAt(s, "MOVE A/D OR LOWER TOUCH   ATTACK SPACE/X OR TOP TOUCH", 55, 685)
	if g.clear {
		overlay(s, "FIVE HITS OBSERVED!\n\nTAP / SPACE TO RESET")
	}
}
func drawFighter(s *ebiten.Image, x float32, c color.RGBA) {
	vector.DrawFilledCircle(s, x, 515, 20, c, false)
	vector.DrawFilledRect(s, x-16, 535, 32, 70, c, false)
}
func overlap(a, b rect) bool { return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y }
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Hitbox Viewer — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
