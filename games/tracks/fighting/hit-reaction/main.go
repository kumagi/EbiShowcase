package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const width, height = 480, 720

type game struct {
	p1, p2, v2                  float64
	attack, hitstop, stun, hits int
	clear                       bool
}

func newGame() *game { return &game{p1: 135, p2: 290} }
func (g *game) Update() error {
	if g.clear {
		if press() {
			*g = *newGame()
		}
		return nil
	}
	if g.hitstop > 0 {
		g.hitstop--
		return nil
	}
	if g.stun > 0 {
		g.stun--
		g.p2 += g.v2
		g.v2 *= .86
	} else {
		g.p2 += (290 - g.p2) * .04
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
	g.p1 = math.Max(30, math.Min(g.p2-40, g.p1))
	if start && g.attack == 0 {
		g.attack = 20
	}
	if g.attack > 0 {
		g.attack--
		if g.attack == 12 && g.p2-g.p1 < 115 {
			g.hitstop = 8
			g.stun = 25
			g.v2 = 7
			g.hits++
			if g.hits >= 5 {
				g.clear = true
			}
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	bg := color.RGBA{19, 28, 44, 255}
	if g.hitstop > 0 {
		bg = color.RGBA{80, 78, 95, 255}
	}
	s.Fill(bg)
	vector.DrawFilledRect(s, 0, 610, 480, 110, color.RGBA{55, 66, 78, 255}, false)
	draw(s, float32(g.p1), color.RGBA{45, 225, 194, 255})
	draw(s, float32(g.p2), color.RGBA{240, 75, 91, 255})
	if g.attack > 8 && g.attack < 16 {
		vector.DrawFilledRect(s, float32(g.p1+18), 530, 82, 25, color.RGBA{255, 210, 62, 255}, false)
	}
	state := "READY"
	if g.hitstop > 0 {
		state = fmt.Sprintf("HIT STOP %d", g.hitstop)
	} else if g.stun > 0 {
		state = fmt.Sprintf("HIT STUN %d  KNOCKBACK %.1f", g.stun, g.v2)
	}
	ebitenutil.DebugPrintAt(s, state, 170, 70)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("HITS %d/5", g.hits), 205, 105)
	ebitenutil.DebugPrintAt(s, "MOVE: A/D OR LOWER TOUCH   ATTACK: SPACE/X OR TOP", 65, 685)
	if g.clear {
		overlay(s, "REACTIONS OBSERVED!\n\nTAP / SPACE TO RESET")
	}
}
func draw(s *ebiten.Image, x float32, c color.RGBA) {
	vector.DrawFilledCircle(s, x, 505, 22, c, false)
	vector.DrawFilledRect(s, x-17, 527, 34, 78, c, false)
}
func press() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Hit Reaction Dojo — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
