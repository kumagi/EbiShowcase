// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
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
	W  = 480
	H  = 720
	N  = 7
	T  = 58
	OX = 37
	OY = 110
)

type pt struct{ x, y int }
type node struct {
	p pt
	c int
}
type pq []node

func (q pq) Len() int           { return len(q) }
func (q pq) Less(i, j int) bool { return q[i].c < q[j].c }
func (q pq) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *pq) Push(v any)        { *q = append(*q, v.(node)) }
func (q *pq) Pop() any          { o := *q; v := o[len(o)-1]; *q = o[:len(o)-1]; return v }

var land = [N][N]int{{0, 0, 1, 1, 0, 2, 0}, {0, 1, 1, 0, 0, 2, 0}, {0, 0, 0, 0, 1, 2, 0}, {2, 2, 0, 1, 1, 0, 0}, {0, 0, 0, 0, 0, 1, 0}, {0, 1, 2, 2, 0, 1, 0}, {0, 0, 0, 1, 0, 0, 0}}

type game struct {
	hero                       pt
	reach                      map[pt]int
	budget, round, score, best int
	message                    string
}

func newGame() *game {
	g := &game{hero: pt{0, 6}, budget: 6, message: "Blue cells fit inside 6 movement points."}
	g.calc()
	return g
}
func (g *game) calc() {
	g.reach = map[pt]int{}
	q := &pq{{g.hero, 0}}
	heap.Init(q)
	for q.Len() > 0 {
		n := heap.Pop(q).(node)
		if old, ok := g.reach[n.p]; ok && old <= n.c {
			continue
		}
		g.reach[n.p] = n.c
		for _, d := range []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			p := pt{n.p.x + d.x, n.p.y + d.y}
			if p.x >= 0 && p.x < N && p.y >= 0 && p.y < N {
				c := n.c + 1 + land[p.y][p.x]
				if c <= g.budget {
					heap.Push(q, node{p, c})
				}
			}
		}
	}
}
func (g *game) Update() error {
	if x, y, ok := press(); ok && x >= OX && x < OX+N*T && y >= OY && y < OY+N*T {
		p := pt{(x - OX) / T, (y - OY) / T}
		if c, yes := g.reach[p]; yes && p != g.hero {
			g.score += g.budget - c + 1
			g.hero = p
			g.round++
			g.message = fmt.Sprintf("Cost %d. Saved %d point(s)!", c, g.budget-c)
			if g.score > g.best {
				g.best = g.score
			}
			if g.round%3 == 0 {
				g.budget = max(3, g.budget-1)
			}
			g.calc()
		} else {
			g.message = "That tile costs too much from here."
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{12, 25, 37, 255})
	ebitenutil.DebugPrintAt(s, "WEIGHTED TERRAIN TREK", 158, 25)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MOVE %d  TRIPS %d  SCORE %d  BEST %d", g.budget, g.round, g.score, g.best), 102, 55)
	ebitenutil.DebugPrintAt(s, g.message, 78, 80)
	cs := []color.RGBA{{103, 151, 77, 255}, {49, 103, 66, 255}, {89, 88, 101, 255}}
	for y := 0; y < N; y++ {
		for x := 0; x < N; x++ {
			p := pt{x, y}
			c := cs[land[y][x]]
			if _, ok := g.reach[p]; ok {
				c = color.RGBA{64, 137, 158, 255}
			}
			vector.DrawFilledRect(s, float32(OX+x*T+1), float32(OY+y*T+1), T-2, T-2, c, false)
			ebitenutil.DebugPrintAt(s, fmt.Sprint(1+land[y][x]), OX+x*T+5, OY+y*T+5)
		}
	}
	x := float32(OX + g.hero.x*T + T/2)
	y := float32(OY + g.hero.y*T + T/2)
	vector.DrawFilledCircle(s, x, y, 19, color.RGBA{244, 174, 64, 255}, false)
	vector.StrokeCircle(s, x, y, 23, 3, color.White, false)
	ebitenutil.DebugPrintAt(s, "Tap a blue destination. Forest=2, mountain=3.", 78, 560)
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
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Weighted Terrain Trek")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
