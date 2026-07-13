package main

import (
	"container/heap"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
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
	p               pt
	hp, move, reach int
	name            string
	enemy, moved    bool
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

type game struct {
	units          []unit
	selected, turn int
	cursor         pt
	reach          map[pt]int
	message        string
	won, lost      bool
}

func newGame() *game {
	g := &game{units: []unit{{pt{0, 7}, 10, 5, 1, "BLADE", false, false}, {pt{1, 7}, 8, 5, 2, "BOW", false, false}, {pt{7, 0}, 4, 4, 1, "RED", true, false}, {pt{6, 3}, 5, 4, 1, "BLUE", true, false}, {pt{7, 7}, 5, 4, 1, "GOLD", true, false}}, selected: 0, cursor: pt{0, 7}, message: "Select a unit. Forest costs 2; mountain costs 3."}
	g.recalc()
	return g
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
	if g.won || g.lost {
		if retry() {
			*g = *newGame()
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
	if x, y, ok := press(); ok && x >= ox && x < ox+cols*tile && y >= oy && y < oy+rows*tile {
		p := pt{(x - ox) / tile, (y - oy) / tile}
		g.choose(p)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.selected = (g.selected + 1) % 2
		g.cursor = g.units[g.selected].p
		g.recalc()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.choose(g.cursor)
	}
	return nil
}
func (g *game) choose(p pt) {
	for i := range g.units {
		if g.units[i].hp > 0 && g.units[i].p == p {
			if g.units[i].enemy {
				g.attack(i)
			} else if !g.units[i].moved {
				g.selected = i
				g.cursor = p
				g.recalc()
			}
			return
		}
	}
	if c, ok := g.reach[p]; ok && !g.units[g.selected].moved {
		g.units[g.selected].p = p
		g.units[g.selected].moved = true
		g.message = fmt.Sprintf("Moved through terrain cost %d. Now choose an enemy in range.", c)
		g.enemyTurnIfDone()
		g.recalc()
	}
}
func (g *game) attack(target int) {
	u := &g.units[g.selected]
	e := &g.units[target]
	d := abs(u.p.x-e.p.x) + abs(u.p.y-e.p.y)
	if u.moved && d <= u.reach {
		e.hp -= 3
		g.message = fmt.Sprintf("%s attacks %d tiles away!", u.name, d)
		if e.hp <= 0 {
			g.message += " Enemy defeated."
		}
		g.enemyTurnIfDone()
	}
}
func (g *game) enemyTurnIfDone() {
	if !g.units[0].moved || !g.units[1].moved {
		return
	}
	g.turn++
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
	g.selected = 0
	g.cursor = g.units[0].p
	g.recalc()
	alive := 0
	for i := 2; i < len(g.units); i++ {
		if g.units[i].hp > 0 {
			alive++
		}
	}
	if alive == 0 {
		g.won = true
	}
	if g.units[0].hp <= 0 && g.units[1].hp <= 0 || g.turn >= 8 {
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
	s.Fill(color.RGBA{14, 24, 35, 255})
	ebitenutil.DebugPrintAt(s, "EBI TACTICS", 194, 18)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TURN %d/8  SELECT %s  MOVE %d  RANGE %d", g.turn+1, g.units[g.selected].name, g.units[g.selected].move, g.units[g.selected].reach), 80, 45)
	ebitenutil.DebugPrintAt(s, g.message, 32, 72)
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
		vector.DrawFilledCircle(s, x, y, 18, c, false)
		if i == g.selected {
			vector.StrokeCircle(s, x, y, 22, 3, color.White, false)
		}
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s%d", u.name[:1], u.hp), int(x)-12, int(y)-5)
	}
	ebitenutil.DebugPrintAt(s, "Tap a unit/tile/enemy | TAB switches ally | Enter confirms", 55, 650)
	if g.won {
		overlay(s, "TACTICAL VICTORY!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(s, "MISSION FAILED\n\nTAP / ENTER TO RETRY")
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
