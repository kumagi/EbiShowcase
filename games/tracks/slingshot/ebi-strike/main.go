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

const (
	screenW     = 480
	screenH     = 720
	allyRadius  = 19.0
	enemyRadius = 27.0
	stopSpeed   = 0.11
	friction    = 0.982
	maxTurns    = 8
)

type vec struct{ x, y float64 }

type ally struct {
	pos      vec
	velocity vec
}

type enemy struct {
	pos      vec
	hp       int
	cooldown int
}
type spark struct{ x, y, vx, vy, life float64 }

type game struct {
	allies                              [2]ally
	enemies                             []enemy
	active, turns                       int
	dragging, moving                    bool
	dragNow                             vec
	allyEffectUsed                      bool
	pulseAt                             vec
	pulseFrames                         int
	message                             string
	won, lost                           bool
	stage, totalTurns, bestTurns, shake int
	pillars                             []vec
	sparks                              []spark
}

func newGame() *game {
	g := &game{stage: 1}
	g.loadStage()
	return g
}
func (g *game) loadStage() {
	g.allies = [2]ally{{pos: vec{125, 525}}, {pos: vec{355, 525}}}
	enemySets := [][]enemy{{{pos: vec{105, 185}, hp: 2}, {pos: vec{370, 205}, hp: 2}, {pos: vec{240, 345}, hp: 3}}, {{pos: vec{80, 160}, hp: 2}, {pos: vec{240, 220}, hp: 4}, {pos: vec{400, 160}, hp: 2}, {pos: vec{240, 430}, hp: 3}}, {{pos: vec{75, 165}, hp: 3}, {pos: vec{405, 165}, hp: 3}, {pos: vec{120, 355}, hp: 3}, {pos: vec{360, 355}, hp: 3}, {pos: vec{240, 260}, hp: 5}}}
	pillarSets := [][]vec{{{240, 205}, {145, 375}, {350, 390}}, {{160, 255}, {320, 255}, {160, 420}, {320, 420}}, {{240, 170}, {90, 270}, {390, 270}, {240, 390}}}
	g.enemies = enemySets[g.stage-1]
	g.pillars = pillarSets[g.stage-1]
	g.active = 0
	g.turns = 0
	g.dragging = false
	g.moving = false
	g.won = false
	g.lost = false
	g.message = "Drag the glowing ally backward and release."
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			if g.won && g.stage < 3 {
				g.totalTurns += g.turns
				g.stage++
				g.loadStage()
			} else {
				best := g.bestTurns
				if g.won {
					total := g.totalTurns + g.turns
					if best == 0 || total < best {
						best = total
					}
				}
				*g = *newGame()
				g.bestTurns = best
			}
		}
		return nil
	}
	if g.pulseFrames > 0 {
		g.pulseFrames--
	}
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	for i := range g.enemies {
		if g.enemies[i].cooldown > 0 {
			g.enemies[i].cooldown--
		}
	}

	if !g.moving {
		g.updateAim()
		return nil
	}

	a := &g.allies[g.active]
	a.pos.x += a.velocity.x
	a.pos.y += a.velocity.y
	g.bounceWalls(a)
	g.bouncePillars(a)
	g.hitEnemies(a)
	g.hitAlly(a)
	a.velocity.x *= friction
	a.velocity.y *= friction
	if math.Hypot(a.velocity.x, a.velocity.y) < stopSpeed {
		a.velocity = vec{}
		g.moving = false
		g.endTurn()
	}
	return nil
}

func (g *game) updateAim() {
	a := &g.allies[g.active]
	if !g.dragging {
		x, y, ok := pressPosition()
		if ok && distance(vec{float64(x), float64(y)}, a.pos) <= allyRadius+18 {
			g.dragging = true
			g.dragNow = vec{float64(x), float64(y)}
		}
		return
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.dragNow = vec{float64(x), float64(y)}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		g.dragNow = vec{float64(x), float64(y)}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) || len(inpututil.AppendJustReleasedTouchIDs(nil)) > 0 {
		pull := vec{a.pos.x - g.dragNow.x, a.pos.y - g.dragNow.y}
		length := math.Hypot(pull.x, pull.y)
		g.dragging = false
		if length < 14 {
			g.message = "Pull farther before release."
			return
		}
		if length > 145 {
			pull.x *= 145 / length
			pull.y *= 145 / length
		}
		a.velocity = vec{pull.x * 0.105, pull.y * 0.105}
		g.turns++
		g.moving = true
		g.allyEffectUsed = false
		g.message = "Moving: bounce, attack, then slow to rest."
	}
}

func (g *game) bounceWalls(a *ally) {
	if a.pos.x < 28 {
		a.pos.x = 28
		a.velocity.x = math.Abs(a.velocity.x) * 0.88
	}
	if a.pos.x > 452 {
		a.pos.x = 452
		a.velocity.x = -math.Abs(a.velocity.x) * 0.88
	}
	if a.pos.y < 105 {
		a.pos.y = 105
		a.velocity.y = math.Abs(a.velocity.y) * 0.88
	}
	if a.pos.y > 575 {
		a.pos.y = 575
		a.velocity.y = -math.Abs(a.velocity.y) * 0.88
	}
}

func (g *game) bouncePillars(a *ally) {
	for _, pillar := range g.pillars {
		g.reflectCircle(a, pillar, 24)
	}
}

func (g *game) reflectCircle(a *ally, center vec, otherRadius float64) float64 {
	dx, dy := a.pos.x-center.x, a.pos.y-center.y
	d := math.Hypot(dx, dy)
	minimum := allyRadius + otherRadius
	if d >= minimum {
		return 0
	}
	if d == 0 {
		dx, d = 1, 1
	}
	nx, ny := dx/d, dy/d
	a.pos = vec{center.x + nx*minimum, center.y + ny*minimum}
	dot := a.velocity.x*nx + a.velocity.y*ny
	impact := math.Abs(dot)
	if dot < 0 {
		a.velocity.x -= 1.85 * dot * nx
		a.velocity.y -= 1.85 * dot * ny
	}
	return impact
}

func (g *game) hitEnemies(a *ally) {
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.hp <= 0 {
			continue
		}
		impact := g.reflectCircle(a, e.pos, enemyRadius)
		if impact >= 1.2 && e.cooldown == 0 {
			e.hp--
			e.cooldown = 22
			g.shake = 4
			g.burst(e.pos.x, e.pos.y, 10)
			g.message = fmt.Sprintf("Direct contact! Enemy HP %d.", e.hp)
		}
	}
	g.checkWin()
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 26 + float64(i%8)})
	}
}

func (g *game) hitAlly(a *ally) {
	if g.allyEffectUsed {
		return
	}
	other := &g.allies[1-g.active]
	if g.reflectCircle(a, other.pos, allyRadius) == 0 {
		return
	}
	g.allyEffectUsed = true
	g.pulseAt = other.pos
	g.pulseFrames = 30
	hits := 0
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.hp > 0 && distance(e.pos, other.pos) <= 165 {
			e.hp--
			hits++
		}
	}
	g.message = fmt.Sprintf("ALLY WAVE! %d enemy hit(s).", hits)
	g.checkWin()
}

func (g *game) checkWin() {
	for _, e := range g.enemies {
		if e.hp > 0 {
			return
		}
	}
	g.won = true
	g.moving = false
	g.message = "Every reef guardian is defeated!"
}

func (g *game) endTurn() {
	if g.won {
		return
	}
	if g.turns >= maxTurns {
		g.lost = true
		g.message = "No turns left. Plan more ally waves!"
		return
	}
	g.active = 1 - g.active
	g.message = fmt.Sprintf("Turn ended at rest. Ally %d is ready.", g.active+1)
}

func (g *game) Draw(screen *ebiten.Image) {
	bgs := []color.RGBA{{10, 22, 41, 255}, {28, 25, 55, 255}, {54, 25, 35, 255}}
	screen.Fill(bgs[g.stage-1])
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.pulseFrames+g.turns)*2) * 5
	}
	ebitenutil.DebugPrintAt(screen, "EBI STRIKE / REEF RESCUE", 149, 21)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("STAGE %d/3 TURN %d/%d ALLY %d ENEMIES %d BEST %d", g.stage, g.turns, maxTurns, g.active+1, g.alive(), g.bestTurns), 55, 49)
	ebitenutil.DebugPrintAt(screen, g.message, 78, 74)
	vector.DrawFilledRect(screen, 18, 88, 444, 505, color.RGBA{25, 55, 72, 255}, false)
	vector.StrokeRect(screen, 18, 88, 444, 505, 4, color.RGBA{83, 137, 158, 255}, false)

	for _, p := range g.pillars {
		trackatlas.DrawCentered(screen, "peg", p.x+ox, p.y, 48)
	}
	for _, e := range g.enemies {
		if e.hp <= 0 {
			continue
		}
		if e.cooldown > 0 {
			trackatlas.DrawTinted(screen, "leaf-guard", e.pos.x+ox, e.pos.y, enemyRadius*2, 1.3, 1.1, 0.7, 1)
		} else {
			trackatlas.DrawCentered(screen, "leaf-guard", e.pos.x+ox, e.pos.y, enemyRadius*2)
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP%d", e.hp), int(e.pos.x)-11, int(e.pos.y)-5)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(screen, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	if g.pulseFrames > 0 {
		r := float32(45 + (30-g.pulseFrames)*4)
		vector.StrokeCircle(screen, float32(g.pulseAt.x), float32(g.pulseAt.y), r, 5, color.RGBA{250, 210, 72, 210}, false)
	}
	for i, a := range g.allies {
		trackatlas.DrawCentered(screen, "ally", a.pos.x, a.pos.y, allyRadius*2)
		if i == g.active && !g.moving {
			vector.StrokeCircle(screen, float32(a.pos.x), float32(a.pos.y), allyRadius+2, 4, color.RGBA{252, 205, 68, 255}, false)
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("A%d", i+1), int(a.pos.x)-7, int(a.pos.y)-5)
	}
	if g.dragging {
		a := g.allies[g.active]
		vector.StrokeLine(screen, float32(a.pos.x), float32(a.pos.y), float32(g.dragNow.x), float32(g.dragNow.y), 5, color.RGBA{246, 184, 64, 255}, false)
		vector.StrokeLine(screen, float32(a.pos.x), float32(a.pos.y), float32(a.pos.x+a.pos.x-g.dragNow.x), float32(a.pos.y+a.pos.y-g.dragNow.y), 2, color.White, false)
		pull := vec{a.pos.x - g.dragNow.x, a.pos.y - g.dragNow.y}
		for i := 1; i <= 6; i++ {
			t := float64(i) * .55
			vector.DrawFilledCircle(screen, float32(a.pos.x+pull.x*.105*t), float32(a.pos.y+pull.y*.105*t), 3, color.RGBA{255, 255, 255, 150}, true)
		}
	}
	ebitenutil.DebugPrintAt(screen, "DRAG ACTIVE ALLY BACKWARD, THEN RELEASE", 89, 620)
	ebitenutil.DebugPrintAt(screen, "Touch ally = wave  |  Touch enemy = direct hit", 74, 650)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FRICTION %.3f  |  STOP %.2f", friction, stopSpeed), 151, 680)
	if g.won {
		msg := "STAGE CLEAR!\n\nTAP / ENTER FOR NEXT STAGE"
		if g.stage == 3 {
			msg = "REEF RESCUED!\n\nTAP / ENTER FOR A NEW RUN"
		}
		overlay(screen, msg)
	}
	if g.lost {
		overlay(screen, "OUT OF TURNS\n\nTAP / ENTER TO RETRY")
	}
}

func (g *game) alive() int {
	n := 0
	for _, e := range g.enemies {
		if e.hp > 0 {
			n++
		}
	}
	return n
}

func distance(a, b vec) float64 { return math.Hypot(a.x-b.x, a.y-b.y) }

func pressPosition() (int, int, bool) {
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

func retryPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	_, _, ok := pressPosition()
	return ok
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 45, 270, 390, 150, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 45, 270, 390, 150, 4, color.RGBA{244, 189, 68, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 120, 326)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ebi Strike — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
