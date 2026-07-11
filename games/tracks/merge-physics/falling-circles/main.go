package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

const width, height = 480, 720

type ball struct {
	x, y, vx, vy, r float64
	c               color.RGBA
	settled         bool
}
type game struct {
	balls []ball
	rng   *rand.Rand
	drops int
	clear bool
}

func newGame() *game { return &game{rng: rand.New(rand.NewSource(4301))} }
func (g *game) Update() error {
	if g.clear {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	if x, ok := drop(); ok && g.drops < 12 {
		r := float64(18 + g.rng.Intn(24))
		g.balls = append(g.balls, ball{float64(x), 90, 0, 0, r, color.RGBA{uint8(90 + g.rng.Intn(150)), uint8(90 + g.rng.Intn(150)), uint8(90 + g.rng.Intn(150)), 255}, false})
		g.drops++
	}
	settled := 0
	for i := range g.balls {
		b := &g.balls[i]
		if b.settled {
			settled++
			continue
		}
		b.vy += .55
		b.x += b.vx
		b.y += b.vy
		if b.x-b.r < 25 {
			b.x = 25 + b.r
			b.vx = math.Abs(b.vx) * .7
		}
		if b.x+b.r > 455 {
			b.x = 455 - b.r
			b.vx = -math.Abs(b.vx) * .7
		}
		if b.y+b.r > 660 {
			b.y = 660 - b.r
			b.vy = -b.vy * .45
			b.vx *= .96
			if math.Abs(b.vy) < .8 {
				b.vy = 0
				b.settled = true
				settled++
			}
		}
	}
	if g.drops == 12 && settled == 12 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 28, 44, 255})
	vector.DrawFilledRect(s, 20, 70, 440, 610, color.RGBA{35, 44, 61, 255}, false)
	vector.StrokeRect(s, 20, 70, 440, 610, 5, color.RGBA{115, 130, 150, 255}, false)
	for _, b := range g.balls {
		vector.DrawFilledCircle(s, float32(b.x), float32(b.y), float32(b.r), b.c, false)
		vector.StrokeCircle(s, float32(b.x), float32(b.y), float32(b.r), 2, color.White, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DROPS %02d/12   FIXED STEP 1/60 SECOND", g.drops), 105, 25)
	ebitenutil.DebugPrintAt(s, "CLICK / TAP TO DROP A CIRCLE", 135, 690)
	if g.clear {
		overlay(s, "ALL CIRCLES SETTLED!\n\nTAP / SPACE TO RESET")
	}
}
func drop() (int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		return max(50, min(430, x)), true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		return max(50, min(430, x)), true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return 240, true
	}
	return 0, false
}
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Falling Circles — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
