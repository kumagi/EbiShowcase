package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const (
	width  = 480
	height = 720
)

type stage struct {
	name                  string
	guards, sight, pearls int
	rule                  string
	c                     color.RGBA
}

var stages = []stage{{"PEARL GARDEN", 1, 5, 16, "One patrol guard teaches safe routes.", color.RGBA{50, 92, 142, 255}}, {"CROSSROAD VAULT", 2, 7, 20, "Two guards remember the last seen tile.", color.RGBA{89, 69, 151, 255}}, {"MOON LABYRINTH", 3, 9, 24, "An ambusher aims three tiles ahead.", color.RGBA{138, 61, 101, 255}}}

type game struct {
	stage, pearls, score, best int
	done                       bool
}

var best int

func (g *game) step() {
	if g.done {
		*g = game{best: best}
		return
	}
	g.pearls++
	g.score += 50 + g.stage*10
	if g.pearls >= stages[g.stage].pearls {
		g.score += 500
		g.stage++
		g.pearls = 0
		if g.stage >= 3 {
			g.stage = 2
			g.done = true
			if g.score > best {
				best = g.score
			}
			g.best = best
		}
	}
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		g.step()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	d := stages[g.stage]
	s.Fill(color.RGBA{8, 18, 35, 255})
	ebitenutil.DebugPrintAt(s, "MAZE MARATHON DIRECTOR", 142, 34)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAGE %d/3  %s", g.stage+1, d.name), 135, 75)
	ebitenutil.DebugPrintAt(s, d.rule, 80, 108)
	for y := 0; y < 7; y++ {
		for x := 0; x < 9; x++ {
			wall := (x+y+g.stage)%4 == 0 || x == 0 || y == 0 || x == 8 || y == 6
			c := color.RGBA{20, 38, 59, 255}
			if wall {
				c = d.c
			}
			vector.DrawFilledRect(s, float32(42+x*44), float32(165+y*44), 40, 40, c, false)
		}
	}
	for i := 0; i < d.guards; i++ {
		vector.DrawFilledCircle(s, float32(100+i*130), 460, 17, color.RGBA{238, 78, 104, 255}, false)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("DATA guards:%d sight:%d pearls:%d", d.guards, d.sight, d.pearls), 106, 515)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("COLLECTED %02d/%02d  SCORE %05d", g.pearls, d.pearls, g.score), 110, 554)
	vector.DrawFilledRect(s, 70, 610, 340, 64, color.RGBA{230, 164, 58, 255}, false)
	label := "TAP / SPACE: COLLECT"
	if g.done {
		label = "3 MAZES CLEAR — REPLAY"
	}
	ebitenutil.DebugPrintAt(s, label, 135, 637)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Maze Marathon Director")
	if err := ebiten.RunGame(&game{best: best}); err != nil {
		panic(err)
	}
}
