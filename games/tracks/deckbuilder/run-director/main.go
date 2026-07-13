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

type node struct {
	name           string
	danger, reward int
}
type game struct {
	floor, score int
	route        [5]int
}

var floors = [][]node{{{"SCOUT", 1, 1}, {"REST", 0, 0}}, {{"MAGE", 2, 2}, {"TREASURE", 1, 3}}, {{"ELITE", 3, 4}, {"REST", 0, 0}}, {{"KNIGHT", 4, 5}, {"SHOP", 1, 2}}, {{"BOSS", 5, 8}, {"BOSS", 5, 8}}}

func (g *game) Update() error {
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y > 420 {
			choice = x / (width / 2)
		}
	}
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y > 420 {
			choice = x / (width / 2)
		}
	}
	if choice >= 0 {
		g.route[g.floor] = choice
		n := floors[g.floor][choice]
		g.score += n.reward*10 - n.danger*2
		g.floor++
		if g.floor >= 5 {
			g.floor = 0
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{21, 31, 49, 255})
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("FLOOR %d/5  RUN SCORE %d", g.floor+1, g.score), 150, 40)
	for i := 0; i < 5; i++ {
		y := 110 + float32(i)*60
		c := color.RGBA{55, 65, 88, 255}
		if i < g.floor {
			c = color.RGBA{45, 205, 181, 255}
		}
		vector.DrawFilledCircle(s, 240, y, 22, c, true)
		if i < g.floor {
			ebitenutil.DebugPrintAt(s, floors[i][g.route[i]].name, 280, int(y)-5)
		}
	}
	for i, n := range floors[g.floor] {
		x := 20 + float32(i)*230
		vector.DrawFilledRect(s, x, 430, 210, 160, color.RGBA{70 + uint8(i)*60, 100, 145, 255}, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d %s\n\nDANGER %d\nREWARD %d", i+1, n.name, n.danger, n.reward), int(x)+45, 465)
	}
	ebitenutil.DebugPrintAt(s, "CHOOSE A ROUTE: 1 / 2 / TAP", 125, 670)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Run Director — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
