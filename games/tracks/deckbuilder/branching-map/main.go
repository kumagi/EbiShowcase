package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type node struct {
	x, y           float32
	kind           string
	danger, reward int
	next           []int
}

var nodes = []node{{240, 620, "START", 0, 0, []int{1, 2}}, {135, 500, "BATTLE", 8, 8, []int{3, 4}}, {345, 500, "REST", 0, 0, []int{4, 5}}, {90, 365, "SHOP", 0, 0, []int{6}}, {240, 365, "ELITE", 15, 18, []int{6, 7}}, {390, 365, "BATTLE", 10, 10, []int{7}}, {155, 225, "REST", 0, 0, []int{8}}, {325, 225, "SHOP", 0, 0, []int{8}}, {240, 90, "BOSS", 22, 30, nil}}

type game struct {
	current, hp, gold, power, steps int
	visited                         []bool
	message                         string
	clear, over                     bool
}

func newGame() *game {
	v := make([]bool, len(nodes))
	v[0] = true
	return &game{hp: 35, gold: 5, power: 12, visited: v, message: "Choose a connected node. Keys 1-2 or tap."}
}
func (g *game) Update() error {
	if g.clear || g.over {
		if restart() {
			*g = *newGame()
		}
		return nil
	}
	options := nodes[g.current].next
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if x, y, ok := press(); ok {
		for i, id := range options {
			n := nodes[id]
			dx, dy := float32(x)-n.x, float32(y)-n.y
			if dx*dx+dy*dy < 42*42 {
				choice = i
				break
			}
		}
	}
	if choice >= 0 && choice < len(options) {
		g.visit(options[choice])
	}
	return nil
}
func (g *game) visit(id int) {
	g.current = id
	g.visited[id] = true
	g.steps++
	n := nodes[id]
	switch n.kind {
	case "BATTLE", "ELITE", "BOSS":
		taken := max(0, n.danger-g.power/2)
		g.hp -= taken
		g.gold += n.reward
		g.power += 2
		g.message = fmt.Sprintf("%s: lost %d HP, gained %d gold.", n.kind, taken, n.reward)
	case "REST":
		g.hp = min(35, g.hp+10)
		g.message = "REST: recovered 10 HP."
	case "SHOP":
		if g.gold >= 8 {
			g.gold -= 8
			g.power += 5
			g.message = "SHOP: bought +5 power for 8 gold."
		} else {
			g.message = "SHOP: need 8 gold, saved it."
		}
	}
	if n.kind == "BOSS" {
		if g.hp > 0 {
			g.clear = true
		} else {
			g.over = true
		}
	} else if g.hp <= 0 {
		g.over = true
	}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	ebitenutil.DebugPrintAt(s, "BRANCHING REEF", 185, 25)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("HP %02d/35   GOLD %02d   POWER %02d   STEPS %d", max(0, g.hp), g.gold, g.power, g.steps), 100, 50)
	for i, n := range nodes {
		for _, to := range n.next {
			t := nodes[to]
			vector.StrokeLine(s, n.x, n.y, t.x, t.y, 4, color.RGBA{65, 83, 112, 255}, false)
		}
		c := color.RGBA{65, 83, 112, 255}
		if g.visited[i] {
			c = color.RGBA{52, 208, 174, 255}
		}
		for _, id := range nodes[g.current].next {
			if id == i {
				c = color.RGBA{240, 177, 65, 255}
			}
		}
		vector.DrawFilledCircle(s, n.x, n.y, 35, c, false)
		ebitenutil.DebugPrintAt(s, n.kind, int(n.x)-24, int(n.y)-5)
	}
	ebitenutil.DebugPrintAt(s, g.message, 50, 670)
	if g.clear {
		overlay(s, "ROUTE COMPLETE!\n\nTAP / SPACE TO RESTART")
	} else if g.over {
		overlay(s, "THE ROUTE ENDED!\n\nTAP / SPACE TO RETRY")
	}
}
func press() (int, int, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x, y, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, m string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, m, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Branching Map — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
