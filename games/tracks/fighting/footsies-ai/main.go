package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
)

const width, height = 480, 720

type game struct {
	px, ax                           float64
	php, ahp, attack, aiAttack, tick int
	msg                              string
}

func newGame() *game {
	return &game{px: 100, ax: 380, php: 50, ahp: 50, msg: "Move in and out of attack range."}
}
func (g *game) Update() error {
	g.tick++
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.px -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.px += 3
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.attack = 18
	}
	d := g.ax - g.px
	if d > 125 {
		g.ax -= 1.5
	} else if d < 75 {
		g.ax += 2
	} else if g.tick%70 == 0 {
		g.aiAttack = 20
	}
	if g.attack > 0 {
		g.attack--
		if g.attack == 10 && d < 105 {
			g.ahp -= 10
			g.msg = "PLAYER HIT: range was correct."
		}
	}
	if g.aiAttack > 0 {
		g.aiAttack--
		if g.aiAttack == 10 && d < 110 {
			g.php -= 9
			g.msg = "RIVAL HIT: back away after its tell."
		}
	}
	g.px = math.Max(25, math.Min(g.ax-45, g.px))
	g.ax = math.Min(455, math.Max(g.px+45, g.ax))
	if g.php <= 0 || g.ahp <= 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*g = *newGame()
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{24, 34, 52, 255})
	vector.DrawFilledRect(s, 0, 590, 480, 130, color.RGBA{65, 70, 80, 255}, false)
	trackatlas.DrawCentered(s, "fighter-p1", g.px, 535, 125)
	trackatlas.DrawCentered(s, "fighter-p2", g.ax, 535, 125)
	vector.StrokeRect(s, float32(g.px), 500, 105, 50, 2, color.RGBA{255, 210, 62, 180}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("YOU %d   DISTANCE %.0f   RIVAL %d", g.php, g.ax-g.px, g.ahp), 115, 40)
	ebitenutil.DebugPrintAt(s, g.msg, 95, 80)
	ebitenutil.DebugPrintAt(s, "MOVE A/D   J/SPACE: POKE", 140, 670)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Footsies AI — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
