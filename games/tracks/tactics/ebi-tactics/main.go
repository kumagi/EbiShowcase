package main

import (
	"container/heap"
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const (
	W    = 480
	H    = 720
	cols = 8
	rows = 8
	tile = 52
	ox   = 32
	oy   = 112
)

type pt struct{ x, y int }
type unit struct {
	p                   pt
	hp, move, reach     int
	name                string
	enemy, moved, acted bool
}
type node struct {
	p    pt
	cost int
}
type pq []node

func (q pq) Len() int           { return len(q) }
func (q pq) Less(i, j int) bool { return q[i].cost < q[j].cost }
func (q pq) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *pq) Push(x any)        { *q = append(*q, x.(node)) }
func (q *pq) Pop() any          { o := *q; n := o[len(o)-1]; *q = o[:len(o)-1]; return n }

var terrain = [rows][cols]int{{0, 0, 1, 1, 0, 0, 2, 0}, {0, 1, 1, 0, 0, 2, 2, 0}, {0, 0, 0, 0, 1, 0, 0, 0}, {2, 2, 0, 1, 1, 0, 2, 0}, {0, 0, 0, 0, 0, 0, 1, 0}, {0, 1, 2, 2, 0, 1, 1, 0}, {0, 0, 0, 1, 0, 0, 0, 0}, {1, 1, 0, 0, 0, 2, 0, 0}}

type mission struct {
	name      string
	turnLimit int
	terrain   [rows][cols]int
	enemies   []pt
}

var missions = []mission{
	{"FOREST GATE", 7, terrain, []pt{{7, 0}, {6, 3}, {7, 7}}},
	{"MOUNTAIN PASS", 8, [rows][cols]int{{2, 2, 1, 0, 0, 1, 2, 2}, {2, 1, 1, 0, 2, 1, 1, 2}, {1, 1, 0, 0, 2, 2, 0, 1}, {0, 0, 0, 1, 1, 0, 0, 0}, {0, 2, 2, 1, 0, 0, 1, 0}, {0, 1, 2, 2, 0, 1, 1, 0}, {0, 0, 0, 1, 0, 0, 2, 0}, {1, 1, 0, 0, 0, 2, 2, 0}}, []pt{{7, 0}, {5, 2}, {7, 4}, {6, 7}}},
	{"RIVER FORT", 9, [rows][cols]int{{0, 0, 0, 1, 1, 0, 0, 0}, {0, 2, 0, 1, 1, 0, 2, 0}, {0, 2, 0, 0, 0, 0, 2, 0}, {0, 2, 2, 1, 1, 2, 2, 0}, {0, 0, 0, 1, 1, 0, 0, 0}, {1, 1, 0, 0, 0, 0, 1, 1}, {0, 0, 0, 2, 2, 0, 0, 0}, {0, 1, 0, 0, 0, 0, 1, 0}}, []pt{{7, 0}, {4, 1}, {7, 5}, {4, 6}}},
}

type spark struct {
	x, y, vx, vy float64
	life         int
}

type game struct {
	units                    []unit
	selected, turn           int
	cursor                   pt
	reach                    map[pt]int
	message                  string
	won, lost                bool
	mission                  int
	totalTurns               int
	bestTurns                int
	frame, attackAnim, shake int
	attackFrom, attackTo     pt
	sparks                   []spark
	rng                      *rand.Rand
	scene                    *ebiten.Image
	audio                    *audio.Context
	gate                     audiolab.Gate
	pulse                    *shaderlab.Pulse
	cam                      cameralab.State
	badge                    *ebiten.Image
}

func newGame() *game {
	g := &game{selected: 0, cursor: pt{0, 7}, rng: rand.New(rand.NewSource(1601)), scene: ebiten.NewImage(W, H)}
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: W / 2, Y: H / 2}, ViewW: W, ViewH: H}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{255, 211, 83, 255})
	g.loadMission(0)
	g.recalc()
	return g
}

func (g *game) loadMission(index int) {
	g.mission = index
	g.turn = 0
	terrain = missions[index].terrain
	g.units = []unit{{pt{0, 7}, 10, 5, 1, "BLADE", false, false, false}, {pt{1, 7}, 8, 5, 2, "BOW", false, false, false}}
	for i, p := range missions[index].enemies {
		g.units = append(g.units, unit{p, 4 + i%2, 4, 1, fmt.Sprintf("E%d", i+1), true, false, false})
	}
	g.selected = 0
	g.cursor = g.units[0].p
	g.message = fmt.Sprintf("MISSION %d: %s. Read terrain and weapon range.", index+1, missions[index].name)
}
func cost(p pt) int    { return 1 + terrain[p.y][p.x] }
func inside(p pt) bool { return p.x >= 0 && p.x < cols && p.y >= 0 && p.y < rows }
func (g *game) occupied(p pt, ignore int) bool {
	for i, u := range g.units {
		if i != ignore && u.hp > 0 && u.p == p {
			return true
		}
	}
	return false
}
func (g *game) recalc() {
	g.reach = map[pt]int{}
	if g.selected < 0 || g.selected >= len(g.units) {
		return
	}
	u := g.units[g.selected]
	q := &pq{{u.p, 0}}
	heap.Init(q)
	for q.Len() > 0 {
		n := heap.Pop(q).(node)
		if old, ok := g.reach[n.p]; ok && old <= n.cost {
			continue
		}
		g.reach[n.p] = n.cost
		for _, d := range []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			p := pt{n.p.x + d.x, n.p.y + d.y}
			nc := n.cost
			if inside(p) {
				nc += cost(p)
			}
			if inside(p) && nc <= u.move && !g.occupied(p, g.selected) {
				heap.Push(q, node{p, nc})
			}
		}
	}
}
func (g *game) Update() error {
	g.frame++
	if g.attackAnim > 0 {
		g.attackAnim--
	}
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .06
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.won || g.lost {
		if retry() {
			best := g.bestTurns
			*g = *newGame()
			g.bestTurns = best
		}
		return nil
	}
	cx, cy := g.cursor.x, g.cursor.y
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		cx--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		cx++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		cy--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		cy++
	}
	if inside(pt{cx, cy}) {
		g.cursor = pt{cx, cy}
	}
	if x, y, ok := press(); ok {
		if x >= ox && x < ox+cols*tile && y >= oy && y < oy+rows*tile {
			p := pt{(x - ox) / tile, (y - oy) / tile}
			g.choose(p)
		} else if y >= 600 {
			g.waitUnit()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		for step := 1; step <= 2; step++ {
			next := (g.selected + step) % 2
			if g.units[next].hp > 0 && !g.units[next].acted {
				g.selected = next
				break
			}
		}
		g.cursor = g.units[g.selected].p
		g.recalc()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.choose(g.cursor)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.waitUnit()
	}
	return nil
}
func (g *game) choose(p pt) {
	for i := range g.units {
		if g.units[i].hp > 0 && g.units[i].p == p {
			if g.units[i].enemy {
				g.attack(i)
			} else if !g.units[i].moved {
				if i == g.selected {
					g.units[i].moved = true
					g.message = "Stayed on this tile. Attack or WAIT."
					return
				}
				g.selected = i
				g.cursor = p
				g.recalc()
			}
			return
		}
	}
	if c, ok := g.reach[p]; ok && !g.units[g.selected].moved && !g.units[g.selected].acted {
		g.units[g.selected].p = p
		g.units[g.selected].moved = true
		g.message = fmt.Sprintf("Moved through terrain cost %d. Now choose an enemy in range.", c)
		g.recalc()
	}
}
func (g *game) attack(target int) {
	u := &g.units[g.selected]
	e := &g.units[target]
	d := abs(u.p.x-e.p.x) + abs(u.p.y-e.p.y)
	if u.moved && !u.acted && d <= u.reach {
		e.hp -= 3
		g.play(680)
		g.attackFrom, g.attackTo = u.p, e.p
		g.attackAnim = 22
		g.shake = 8
		g.burst(e.p)
		g.message = fmt.Sprintf("%s attacks %d tiles away!", u.name, d)
		if e.hp <= 0 {
			g.message += " Enemy defeated."
		}
		u.acted = true
		g.enemyTurnIfDone()
	}
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Square, freq, .06)).Play()
}

func (g *game) waitUnit() {
	if g.selected < 0 || g.selected >= 2 || g.units[g.selected].acted {
		return
	}
	g.units[g.selected].moved = true
	g.units[g.selected].acted = true
	g.message = g.units[g.selected].name + " waits and ends its action."
	g.enemyTurnIfDone()
}
func (g *game) burst(p pt) {
	for i := 0; i < 24; i++ {
		a := g.rng.Float64() * math.Pi * 2
		s := 1 + g.rng.Float64()*3
		g.sparks = append(g.sparks, spark{float64(ox + p.x*tile + tile/2), float64(oy + p.y*tile + tile/2), math.Cos(a) * s, math.Sin(a)*s - 1, 20 + g.rng.Intn(18)})
	}
}
func (g *game) enemyTurnIfDone() {
	if (g.units[0].hp > 0 && !g.units[0].acted) || (g.units[1].hp > 0 && !g.units[1].acted) {
		return
	}
	g.turn++
	g.totalTurns++
	for i := 2; i < len(g.units); i++ {
		e := &g.units[i]
		if e.hp <= 0 {
			continue
		}
		best := 0
		if dist(e.p, g.units[1].p) < dist(e.p, g.units[0].p) {
			best = 1
		}
		d := stepToward(e.p, g.units[best].p)
		n := pt{e.p.x + d.x, e.p.y + d.y}
		if inside(n) && !g.occupied(n, i) {
			e.p = n
		}
		if dist(e.p, g.units[best].p) <= 1 {
			g.units[best].hp -= 2
		}
	}
	g.units[0].moved = false
	g.units[1].moved = false
	g.units[0].acted = false
	g.units[1].acted = false
	g.selected = 0
	if g.units[0].hp <= 0 {
		g.selected = 1
	}
	g.cursor = g.units[0].p
	g.recalc()
	alive := 0
	for i := 2; i < len(g.units); i++ {
		if g.units[i].hp > 0 {
			alive++
		}
	}
	if alive == 0 {
		if g.mission < len(missions)-1 {
			g.loadMission(g.mission + 1)
			g.recalc()
			return
		}
		g.won = true
		if g.bestTurns == 0 || g.totalTurns < g.bestTurns {
			g.bestTurns = g.totalTurns
		}
		return
	}
	if g.units[0].hp <= 0 && g.units[1].hp <= 0 || g.turn >= missions[g.mission].turnLimit {
		g.lost = true
	}
}
func dist(a, b pt) int { return abs(a.x-b.x) + abs(a.y-b.y) }
func stepToward(a, b pt) pt {
	if abs(b.x-a.x) > abs(b.y-a.y) {
		if b.x > a.x {
			return pt{1, 0}
		}
		return pt{-1, 0}
	}
	if b.y > a.y {
		return pt{0, 1}
	}
	return pt{0, -1}
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func (g *game) Draw(s *ebiten.Image) {
	g.scene.Clear()
	g.drawScene(g.scene)
	op := &ebiten.DrawImageOptions{}
	if g.shake > 0 {
		op.GeoM.Translate(float64((g.frame%3-1)*3), float64(((g.frame/2)%3-1)*2))
	}
	s.DrawImage(g.scene, op)
}
func (g *game) drawScene(s *ebiten.Image) {
	s.Fill([]color.RGBA{{14, 24, 35, 255}, {28, 20, 39, 255}, {10, 37, 43, 255}}[g.mission])
	g.drawTitle(s)
	g.drawEffectBadge(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MISSION %d/3 %-13s TURN %d/%d TOTAL %d BEST %d", g.mission+1, missions[g.mission].name, g.turn+1, missions[g.mission].turnLimit, g.totalTurns, g.bestTurns), 38, 42)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SELECT %s  MOVE %d  RANGE %d", g.units[g.selected].name, g.units[g.selected].move, g.units[g.selected].reach), 135, 61)
	ebitenutil.DebugPrintAt(s, g.message, 32, 82)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			p := pt{x, y}
			c := []color.RGBA{{102, 151, 78, 255}, {54, 112, 70, 255}, {94, 92, 104, 255}}[terrain[y][x]]
			if _, ok := g.reach[p]; ok {
				c = color.RGBA{76, 139, 158, 255}
			}
			vector.DrawFilledRect(s, float32(ox+x*tile+1), float32(oy+y*tile+1), tile-2, tile-2, c, false)
			ebitenutil.DebugPrintAt(s, fmt.Sprint(cost(p)), ox+x*tile+4, oy+y*tile+4)
		}
	}
	vector.StrokeRect(s, float32(ox+g.cursor.x*tile+2), float32(oy+g.cursor.y*tile+2), tile-4, tile-4, 3, color.RGBA{255, 255, 255, 255}, false)
	for i, u := range g.units {
		if u.hp <= 0 {
			continue
		}
		c := color.RGBA{231, 93, 75, 255}
		if u.enemy {
			c = color.RGBA{126, 76, 154, 255}
		}
		x, y := float32(ox+u.p.x*tile+tile/2), float32(oy+u.p.y*tile+tile/2)
		if g.attackAnim > 10 && u.p == g.attackFrom {
			progress := float32(22-g.attackAnim) / 12
			x += float32(g.attackTo.x-g.attackFrom.x) * tile * progress * .55
			y += float32(g.attackTo.y-g.attackFrom.y) * tile * progress * .55
		}
		y += float32(math.Sin(float64(g.frame+i*9)*.12)) * 2
		vector.DrawFilledCircle(s, x, y, 18, c, false)
		if i == g.selected {
			vector.StrokeCircle(s, x, y, 22, 3, color.White, false)
		}
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s%d", u.name[:1], u.hp), int(x)-12, int(y)-5)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x), float32(p.y), 3, color.RGBA{255, 211, 83, uint8(min(255, p.life*10))}, false)
	}
	vector.DrawFilledRect(s, 150, 600, 180, 38, color.RGBA{52, 82, 118, 255}, false)
	ebitenutil.DebugPrintAt(s, "WAIT [W]", 210, 615)
	ebitenutil.DebugPrintAt(s, "Tap unit → tile → enemy | TAB switches ally", 74, 662)
	if g.won {
		overlay(s, fmt.Sprintf("3 MISSIONS CLEARED!\nTOTAL TURNS %d  BEST %d\nTAP / ENTER TO REPLAY", g.totalTurns, g.bestTurns))
	}
	if g.lost {
		overlay(s, "MISSION FAILED\n\nTAP / ENTER TO RETRY")
	}
}
func (g *game) drawTitle(s *ebiten.Image) {
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(194, 6)
		text.Draw(s, "EBI TACTICS", face, op)
		return
	}
	ebitenutil.DebugPrintAt(s, "EBI TACTICS", 194, 18)
}
func (g *game) drawEffectBadge(s *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.frame)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(W-34, 10)
	s.DrawImage(fx, op)
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
	ebitenutil.DebugPrintAt(s, t, 125, 328)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Tactics")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
