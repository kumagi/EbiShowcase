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
	W        = 480
	H        = 720
	maxGauge = 1000
)

type actor struct {
	name                          string
	hp, maxHP, speed, gauge, side int
	ready                         bool
}
type game struct {
	actors    []actor
	ready     int
	frames    int
	message   string
	won, lost bool
}

func newGame() *game {
	return &game{actors: []actor{{"EBI", 42, 42, 8, 0, 0, false}, {"MAGE", 30, 30, 11, 0, 0, false}, {"SHELL", 55, 55, 5, 0, 0, false}, {"WISP", 34, 34, 9, 0, 1, false}, {"GOLEM", 65, 65, 4, 0, 1, false}}, ready: -1, message: "Speed is added to every action gauge each frame."}
}
func (g *game) Update() error {
	if g.won || g.lost {
		if retry() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	for i := range g.actors {
		a := &g.actors[i]
		if a.hp <= 0 || a.ready {
			continue
		}
		a.gauge += a.speed
		if a.gauge >= maxGauge {
			a.gauge = maxGauge
			a.ready = true
			if a.side == 0 && g.ready < 0 {
				g.ready = i
			} else if a.side == 1 {
				g.enemyAct(i)
			}
		}
	}
	if g.ready >= 0 {
		action := -1
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			action = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			action = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			action = 2
		}
		if x, y, ok := press(); ok && y >= 580 {
			action = min(2, x/160)
		}
		if action >= 0 {
			g.playerAct(action)
		}
	}
	g.check()
	return nil
}
func (g *game) playerAct(action int) {
	a := &g.actors[g.ready]
	switch action {
	case 0:
		t := g.firstAlive(1)
		g.actors[t].hp -= 8 + a.speed/2
		g.message = a.name + " attacks."
	case 1:
		t := g.firstAlive(1)
		g.actors[t].hp -= 14
		a.hp -= 3
		g.message = a.name + " uses a strong skill."
	case 2:
		t := g.lowestAlly()
		g.actors[t].hp = min(g.actors[t].maxHP, g.actors[t].hp+12)
		g.message = a.name + " heals " + g.actors[t].name
	}
	a.gauge = 0
	a.ready = false
	g.ready = -1
	for i := range g.actors {
		if g.actors[i].side == 0 && g.actors[i].ready {
			g.ready = i
			break
		}
	}
}
func (g *game) enemyAct(i int) {
	a := &g.actors[i]
	t := g.lowestAlly()
	g.actors[t].hp -= 6 + a.speed/2
	a.gauge = 0
	a.ready = false
	g.message = a.name + " acts as soon as its gauge fills."
}
func (g *game) firstAlive(side int) int {
	for i, a := range g.actors {
		if a.side == side && a.hp > 0 {
			return i
		}
	}
	return 0
}
func (g *game) lowestAlly() int {
	best := -1
	for i, a := range g.actors {
		if a.side == 0 && a.hp > 0 && (best < 0 || a.hp < g.actors[best].hp) {
			best = i
		}
	}
	return max(0, best)
}
func (g *game) check() {
	ally, enemy := 0, 0
	for _, a := range g.actors {
		if a.hp > 0 {
			if a.side == 0 {
				ally++
			} else {
				enemy++
			}
		}
	}
	if enemy == 0 {
		g.won = true
	}
	if ally == 0 || g.frames > 90*60 {
		g.lost = true
	}
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 18, 35, 255})
	ebitenutil.DebugPrintAt(s, "ACTIVE GAUGE BATTLE", 168, 18)
	ebitenutil.DebugPrintAt(s, "Gauge += speed every Update. Full gauge = READY.", 68, 48)
	ebitenutil.DebugPrintAt(s, g.message, 42, 75)
	for i, a := range g.actors {
		x := 45
		if a.side == 1 {
			x = 270
		}
		y := 120 + (i%3)*130
		c := color.RGBA{54, 126, 168, 255}
		if a.side == 1 {
			c = color.RGBA{153, 70, 91, 255}
		}
		vector.DrawFilledRect(s, float32(x), float32(y), 165, 100, c, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s  HP %d/%d", a.name, max(0, a.hp), a.maxHP), x+12, y+15)
		vector.DrawFilledRect(s, float32(x+12), float32(y+48), 140, 14, color.RGBA{25, 35, 50, 255}, false)
		vector.DrawFilledRect(s, float32(x+12), float32(y+48), float32(140*a.gauge/maxGauge), 14, color.RGBA{244, 190, 65, 255}, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("SPD %d  %s", a.speed, map[bool]string{true: "READY", false: "WAIT"}[a.ready]), x+12, y+72)
	}
	labels := [...]string{"[1] ATTACK", "[2] SKILL", "[3] HEAL"}
	for i, l := range labels {
		c := color.RGBA{48, 78, 115, 255}
		if g.ready >= 0 {
			c = color.RGBA{176, 104, 55, 255}
		}
		vector.DrawFilledRect(s, float32(i*160+5), 580, 150, 70, c, false)
		ebitenutil.DebugPrintAt(s, l, i*160+28, 611)
	}
	ebitenutil.DebugPrintAt(s, "Keys 1-3 or tap when an ally says READY", 105, 682)
	if g.won {
		overlay(s, "BATTLE CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(s, "PARTY DEFEATED\n\nTAP / ENTER TO RETRY")
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
func retry() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 45, 270, 390, 150, color.RGBA{4, 12, 24, 240}, false)
	ebitenutil.DebugPrintAt(s, t, 135, 328)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Active Battle")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
