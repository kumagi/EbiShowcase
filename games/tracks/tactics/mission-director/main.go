// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
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
	W  = 480
	H  = 720
	N  = 7
	T  = 58
	OX = 37
	OY = 100
)

type pt struct{ x, y int }
type mission struct {
	name        string
	start, goal pt
	limit       int
	walls       []pt
	c           color.RGBA
}

var ms = []mission{{"REACH THE FLAG", pt{0, 6}, pt{6, 0}, 14, []pt{{2, 4}, {3, 4}, {4, 4}}, color.RGBA{27, 75, 91, 255}}, {"RESCUE SCOUT", pt{6, 6}, pt{0, 1}, 12, []pt{{4, 5}, {4, 4}, {4, 3}, {2, 2}}, color.RGBA{72, 51, 91, 255}}, {"ESCAPE FORT", pt{3, 3}, pt{6, 6}, 9, []pt{{2, 3}, {4, 3}, {3, 2}, {3, 4}}, color.RGBA{38, 82, 60, 255}}}

type game struct {
	stage             int
	p                 pt
	turn, total, best int
	won, lost         bool
	message           string
}

func newGame() *game { g := &game{}; g.load(0); return g }
func (g *game) load(i int) {
	g.stage = i
	g.p = ms[i].start
	g.turn = 0
	g.won = false
	g.lost = false
	g.message = ms[i].name + ": reach the gold objective."
}
func (g *game) blocked(p pt) bool {
	for _, w := range ms[g.stage].walls {
		if p == w {
			return true
		}
	}
	return false
}
func (g *game) Update() error {
	if g.won || g.lost {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || pressed() {
			b := g.best
			*g = *newGame()
			g.best = b
		}
		return nil
	}
	d := pt{}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		d.x = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		d.x = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		d.y = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		d.y = 1
	}
	if x, y, ok := press(); ok && y > 570 {
		switch x / 120 {
		case 0:
			d.x = -1
		case 1:
			d.y = -1
		case 2:
			d.y = 1
		case 3:
			d.x = 1
		}
	}
	if d != (pt{}) {
		n := pt{g.p.x + d.x, g.p.y + d.y}
		if n.x >= 0 && n.x < N && n.y >= 0 && n.y < N && !g.blocked(n) {
			g.p = n
			g.turn++
			g.total++
			if g.p == ms[g.stage].goal {
				if g.stage == len(ms)-1 {
					g.won = true
					if g.best == 0 || g.total < g.best {
						g.best = g.total
					}
				} else {
					g.load(g.stage + 1)
				}
			} else if g.turn >= ms[g.stage].limit {
				g.lost = true
			}
		}
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(ms[g.stage].c)
	ebitenutil.DebugPrintAt(s, "MISSION DATA DIRECTOR", 154, 24)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MISSION %d/3 %s  TURN %d/%d TOTAL %d BEST %d", g.stage+1, ms[g.stage].name, g.turn, ms[g.stage].limit, g.total, g.best), 35, 54)
	ebitenutil.DebugPrintAt(s, g.message, 82, 80)
	for y := 0; y < N; y++ {
		for x := 0; x < N; x++ {
			p := pt{x, y}
			c := color.RGBA{52, 78, 91, 255}
			if g.blocked(p) {
				c = color.RGBA{91, 88, 104, 255}
			}
			if p == ms[g.stage].goal {
				c = color.RGBA{220, 166, 60, 255}
			}
			vector.DrawFilledRect(s, float32(OX+x*T+1), float32(OY+y*T+1), T-2, T-2, c, false)
		}
	}
	vector.DrawFilledCircle(s, float32(OX+g.p.x*T+T/2), float32(OY+g.p.y*T+T/2), 18, color.RGBA{72, 184, 198, 255}, false)
	for i, n := range []string{"LEFT", "UP", "DOWN", "RIGHT"} {
		vector.DrawFilledRect(s, float32(i*120+4), 580, 112, 76, color.RGBA{45, 79, 117, 255}, false)
		ebitenutil.DebugPrintAt(s, n, i*120+40, 613)
	}
	if g.won {
		over(s, "3 MISSIONS CLEAR!\nTAP / ENTER TO REPLAY")
	} else if g.lost {
		over(s, "TURN LIMIT!\nTAP / ENTER TO RETRY")
	}
}
func over(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 50, 270, 380, 130, color.RGBA{5, 13, 28, 244}, false)
	ebitenutil.DebugPrintAt(s, t, 135, 320)
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
func pressed() bool                        { _, _, ok := press(); return ok }
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Mission Data Director")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
