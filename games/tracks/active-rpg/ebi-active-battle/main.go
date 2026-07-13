package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const W, H, maxGauge = 480, 720, 1000

type actor struct {
	name, sprite                        string
	hp, maxHP, speed, gauge, side, role int
	ready                               bool
}
type encounter struct {
	name    string
	bg      color.RGBA
	enemies []actor
}

var encounters = []encounter{
	{"MOONLIT SHORE", color.RGBA{12, 28, 50, 255}, []actor{{"WISP", "ghost-patrol", 30, 30, 12, 0, 1, 0, false}, {"CRAB", "king-crab", 48, 48, 6, 0, 1, 1, false}}},
	{"CLOCKWORK CAVE", color.RGBA{32, 30, 48, 255}, []actor{{"SCOUT", "scout", 38, 38, 15, 0, 1, 0, false}, {"GUARD", "leaf-guard", 62, 62, 7, 0, 1, 1, false}, {"MENDER", "species-2", 34, 34, 10, 0, 1, 2, false}}},
	{"TEMPEST THRONE", color.RGBA{42, 18, 43, 255}, []actor{{"STORM KING", "boss-crab", 130, 130, 9, 0, 1, 3, false}, {"SPARK", "ghost-chase", 42, 42, 16, 0, 1, 0, false}}},
}

type spark struct {
	x, y, vx, vy, life float64
	c                  color.RGBA
}
type motion struct {
	active                              bool
	source, target, kind, timer, damage int
}
type game struct {
	actors                                                                  []actor
	ready, stage, frames, totalFrames, best, shake, flashTarget, flashTimer int
	message                                                                 string
	won, lost, stageClear                                                   bool
	motion                                                                  motion
	sparks                                                                  []spark
}

func party() []actor {
	return []actor{{"TENJIROH", "hero", 54, 54, 9, 0, 0, 0, false}, {"MAGE", "ally", 38, 38, 13, 0, 0, 1, false}, {"SHELL", "pet", 72, 72, 6, 0, 0, 2, false}}
}
func newGame() *game {
	g := &game{ready: -1, best: 0, flashTarget: -1}
	g.loadStage(0, party())
	return g
}
func (g *game) loadStage(n int, allies []actor) {
	g.stage = n
	g.actors = append([]actor{}, allies...)
	for _, e := range encounters[n].enemies {
		g.actors = append(g.actors, e)
	}
	g.ready = -1
	g.frames = 0
	g.stageClear = false
	g.motion = motion{}
	g.flashTarget, g.flashTimer = -1, 0
	g.message = "Speed fills every gauge. Watch who becomes READY."
}
func (g *game) Update() error {
	g.totalFrames++
	g.frames++
	if g.shake > 0 {
		g.shake--
	}
	if g.flashTimer > 0 {
		g.flashTimer--
		if g.flashTimer == 0 {
			g.flashTarget = -1
		}
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .04
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.won || g.lost || g.stageClear {
		if retry() {
			if g.stageClear {
				allies := []actor{}
				for _, a := range g.actors {
					if a.side == 0 && a.hp > 0 {
						a.hp = min(a.maxHP, a.hp+14)
						a.gauge = 0
						a.ready = false
						allies = append(allies, a)
					}
				}
				g.loadStage(g.stage+1, allies)
			} else {
				best := g.best
				*g = *newGame()
				g.best = best
			}
		}
		return nil
	}
	if g.motion.active {
		g.updateMotion()
		return nil
	}
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
				g.startEnemy(i)
				return nil
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
		if x, y, ok := press(); ok && y >= 590 {
			action = min(2, x/160)
		}
		if action >= 0 {
			g.startPlayer(action)
		}
	}
	g.check()
	return nil
}
func (g *game) startPlayer(kind int) {
	source := g.ready
	target := g.firstAlive(1)
	damage := 9 + g.actors[source].speed/2
	if kind == 1 {
		damage = 18
		g.actors[source].hp = max(1, g.actors[source].hp-3)
	}
	if kind == 2 {
		target = g.lowestAlly()
		damage = -15
	}
	g.motion = motion{true, source, target, kind, 34, damage}
	g.message = g.actors[source].name + " prepares an action..."
	g.actors[source].gauge = 0
	g.actors[source].ready = false
	g.ready = -1
}
func (g *game) startEnemy(i int) {
	a := &g.actors[i]
	target := g.lowestAlly()
	damage := 7 + a.speed/3
	if a.role == 1 {
		damage = 13
	}
	if a.role == 2 {
		target = g.lowestEnemy()
		damage = -11
	}
	if a.role == 3 && a.hp < a.maxHP/2 {
		damage = 19
	}
	g.motion = motion{true, i, target, 3, 36, damage}
	g.message = a.name + " telegraphs its move!"
	a.gauge = 0
	a.ready = false
}
func (g *game) updateMotion() {
	m := &g.motion
	m.timer--
	if m.timer == 16 {
		t := &g.actors[m.target]
		if m.damage < 0 {
			t.hp = min(t.maxHP, t.hp-m.damage)
			g.message = g.actors[m.source].name + " restores " + t.name
		} else {
			t.hp -= m.damage
			g.message = fmt.Sprintf("%s hits %s for %d!", g.actors[m.source].name, t.name, m.damage)
		}
		g.flashTarget = m.target
		g.flashTimer = 8
		g.shake = 6
		if m.kind == 1 || m.kind == 3 {
			g.shake = 10
		}
		x, y := g.actorPos(m.target)
		for i := 0; i < 18; i++ {
			a := float64(i) * math.Pi / 9
			g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * 3, math.Sin(a) * 3, 28, color.RGBA{255, 198, 75, 255}})
		}
	}
	if m.timer <= 0 {
		m.active = false
		g.pickReady()
		g.check()
	}
}
func (g *game) pickReady() {
	g.ready = -1
	for i, a := range g.actors {
		if a.side == 0 && a.hp > 0 && a.ready {
			g.ready = i
			return
		}
	}
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
func (g *game) lowestEnemy() int {
	best := -1
	for i, a := range g.actors {
		if a.side == 1 && a.hp > 0 && (best < 0 || a.hp < g.actors[best].hp) {
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
		score := 3000 - g.totalFrames/3
		for _, a := range g.actors {
			if a.side == 0 {
				score += max(0, a.hp) * 10
			}
		}
		if score > g.best {
			g.best = score
		}
		if g.stage == len(encounters)-1 {
			g.won = true
		} else {
			g.stageClear = true
		}
	}
	if ally == 0 || g.frames > 100*60 {
		g.lost = true
	}
}
func (g *game) actorPos(i int) (float64, float64) {
	a := g.actors[i]
	slot := 0
	for j := 0; j < i; j++ {
		if g.actors[j].side == a.side {
			slot++
		}
	}
	x := 115.0
	if a.side == 1 {
		x = 365
	}
	return x, 155 + float64(slot)*135
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(encounters[g.stage].bg)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.totalFrames)*2.3) * float64(g.shake)
	}
	vector.DrawFilledCircle(s, float32(80+g.stage*150), 90, 90, color.RGBA{50, 120, 140, 35}, true)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ACTIVE BATTLE %d/3 / %s", g.stage+1, encounters[g.stage].name), 125, 18)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TIME %02d  BEST %d", g.totalFrames/60, g.best), 170, 43)
	ebitenutil.DebugPrintAt(s, g.message, 45, 68)
	for i, a := range g.actors {
		x, y := g.actorPos(i)
		if a.hp <= 0 {
			continue
		}
		lunge := 0.0
		if g.motion.active && g.motion.source == i {
			p := 1 - math.Abs(float64(g.motion.timer-17))/17
			lunge = p * 28
			if a.side == 1 {
				lunge = -lunge
			}
		}
		bob := math.Sin(float64(g.totalFrames)*.08+float64(i)) * 2
		scale := 78.0
		if i == g.flashTarget {
			scale = 88
		}
		trackatlas.DrawCentered(s, a.sprite, x+lunge+ox, y+bob, scale)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s HP %d/%d", a.name, max(0, a.hp), a.maxHP), int(x)-58, int(y)+45)
		vector.DrawFilledRect(s, float32(x-58), float32(y+62), 116, 11, color.RGBA{20, 28, 42, 255}, false)
		vector.DrawFilledRect(s, float32(x-58), float32(y+62), float32(116*a.gauge/maxGauge), 11, color.RGBA{245, 187, 65, 255}, false)
		if a.ready {
			ebitenutil.DebugPrintAt(s, "READY!", int(x)-20, int(y)+79)
		}
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/9), p.c, true)
	}
	labels := [3]string{"[1] ATTACK", "[2] BURST", "[3] HEAL"}
	for i, l := range labels {
		c := color.RGBA{45, 65, 90, 255}
		if g.ready >= 0 && !g.motion.active {
			c = color.RGBA{185, 102, 48, 255}
		}
		vector.DrawFilledRect(s, float32(i*160+5), 590, 150, 70, c, false)
		ebitenutil.DebugPrintAt(s, l, i*160+28, 620)
	}
	ebitenutil.DebugPrintAt(s, "Keys 1-3 or tap when an ally is READY", 105, 685)
	if g.stageClear {
		overlay(s, fmt.Sprintf("ENCOUNTER CLEAR! BEST %d\n\nTAP / ENTER: NEXT BATTLE", g.best))
	}
	if g.won {
		overlay(s, fmt.Sprintf("TEMPEST CONQUERED! BEST %d\n\nTAP / ENTER: NEW RUN", g.best))
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
	vector.DrawFilledRect(s, 35, 265, 410, 165, color.RGBA{4, 12, 25, 242}, false)
	vector.StrokeRect(s, 35, 265, 410, 165, 4, color.RGBA{245, 188, 65, 255}, false)
	ebitenutil.DebugPrintAt(s, t, 90, 325)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Active Battle")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
