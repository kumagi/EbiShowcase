package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
)

const W, H = 480, 720

type stage struct {
	name                 string
	enemies, speed, goal int
	bg                   color.RGBA
}

var stages = []stage{{"MOONLIT SHORE", 2, 7, 3, color.RGBA{12, 30, 52, 255}}, {"CLOCKWORK CAVE", 3, 10, 4, color.RGBA{34, 30, 48, 255}}, {"TEMPEST THRONE", 2, 15, 5, color.RGBA{44, 18, 44, 255}}}

type game struct{ stage, wins, best, burst int }

func (g *game) Update() error {
	if g.burst > 0 {
		g.burst--
	}
	pressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
	if pressed {
		g.wins++
		g.burst = 15
		if g.wins >= stages[g.stage].goal {
			grade := 1000 + g.stage*500
			if grade > g.best {
				g.best = grade
			}
			g.stage = (g.stage + 1) % 3
			g.wins = 0
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	st := stages[g.stage]
	s.Fill(st.bg)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENCOUNTER %d/3 / %s", g.stage+1, st.name), 135, 65)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ENEMIES %d  SPEED %d  GOAL %d", st.enemies, st.speed, st.goal), 125, 100)
	for i := 0; i < st.enemies; i++ {
		trackatlas.DrawCentered(s, []string{"ghost-patrol", "leaf-guard", "boss-crab"}[g.stage], 130+float64(i)*110, 300, 95)
	}
	vector.DrawFilledRect(s, 90, 465, 300, 20, color.RGBA{35, 48, 68, 255}, false)
	vector.DrawFilledRect(s, 90, 465, float32(300*g.wins/st.goal), 20, color.RGBA{245, 185, 65, 255}, false)
	if g.burst > 0 {
		vector.StrokeCircle(s, 240, 350, float32(40+(15-g.burst)*6), 5, color.RGBA{255, 205, 70, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("VICTORIES %d/%d   BEST %d", g.wins, st.goal, g.best), 130, 520)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: WIN TEST BATTLE", 120, 650)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
