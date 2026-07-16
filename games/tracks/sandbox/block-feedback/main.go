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

const (
	width  = 480
	height = 720
)

type chip struct {
	x, y, vx, vy float64
	life         int
}
type game struct {
	hp, phase, combo, score, shake int
	chips                          []chip
	message                        string
}

func newGame() *game {
	return &game{hp: 4, message: "Tap the block: anticipation comes before contact."}
}
func (g *game) hit() {
	if g.phase == 0 && g.hp > 0 {
		g.phase = 24
	}
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		if g.hp <= 0 {
			*g = *newGame()
		} else {
			g.hit()
		}
	}
	if g.phase > 0 {
		g.phase--
		if g.phase == 10 {
			g.hp--
			g.combo++
			g.score += 100 * g.combo
			g.shake = 8
			for i := 0; i < 14; i++ {
				a := float64(i) * .7
				g.chips = append(g.chips, chip{240, 350, math.Cos(a) * 3, math.Sin(a)*3 - 1, 32})
			}
			g.message = fmt.Sprintf("CONTACT! combo x%d — impact happens on one tick.", g.combo)
			if g.hp == 0 {
				g.message = "Block cleared! Tap to replay and beat your rhythm."
			}
		}
	}
	for i := len(g.chips) - 1; i >= 0; i-- {
		p := &g.chips[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .1
		p.life--
		if p.life <= 0 {
			g.chips = append(g.chips[:i], g.chips[i+1:]...)
		}
	}
	if g.shake > 0 {
		g.shake--
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{16, 34, 50, 255})
	ebitenutil.DebugPrintAt(s, "BLOCK REACTION LAB", 165, 35)
	ebitenutil.DebugPrintAt(s, "WIND-UP  →  CONTACT  →  RECOVERY", 107, 78)
	x, y := float32(240), float32(350)
	if g.phase > 10 {
		x -= float32(math.Sin(float64(g.phase-10)/14*math.Pi) * 38)
	}
	if g.shake > 0 {
		x += float32(g.shake%3-1) * 3
	}
	squash := float32(0)
	if g.phase <= 10 && g.phase > 5 {
		squash = 8
	}
	vector.DrawFilledRect(s, x-65, y-65+squash/2, 130, 130-squash, color.RGBA{107, 165, 183, 255}, false)
	vector.StrokeRect(s, x-65, y-65+squash/2, 130, 130-squash, 5, color.RGBA{225, 240, 240, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("HP %d/4", g.hp), 218, 345)
	for _, p := range g.chips {
		vector.DrawFilledCircle(s, float32(p.x), float32(p.y), 4, color.RGBA{255, 200, 83, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, g.message, 43, 515)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SCORE %04d", g.score), 196, 560)
	vector.DrawFilledRect(s, 70, 610, 340, 64, color.RGBA{216, 132, 50, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE TO SWING", 154, 637)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Block Reaction Lab")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
