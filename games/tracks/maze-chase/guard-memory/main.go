package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

const (
	width  = 480
	height = 720
)

type game struct {
	player, guard, lastSeen int
	mode, timer, score      int
	message                 string
}

func newGame() *game {
	return &game{player: 1, guard: 8, lastSeen: 8, message: "Cross the open lane, then hide behind a wall."}
}
func (g *game) move(d int) { g.player = max(1, min(9, g.player+d)) }
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.move(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.move(1)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, _ := ebiten.CursorPosition()
		if x < 240 {
			g.move(-1)
		} else {
			g.move(1)
		}
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, _ := ebiten.TouchPosition(ids[0])
		if x < 240 {
			g.move(-1)
		} else {
			g.move(1)
		}
	}
	visible := !(g.player >= 4 && g.player <= 6)
	if visible {
		g.mode = 1
		g.lastSeen = g.player
		g.timer = 90
		g.message = "CHASE: the guard stores the tile it can see."
	} else if g.mode == 1 {
		g.mode = 2
		g.message = "SEARCH: sight is lost, but lastSeen remains."
	} else if g.mode == 2 {
		g.timer--
		if g.timer <= 0 || g.guard == g.lastSeen {
			g.mode = 0
			g.message = "PATROL: memory expired, so the route resumes."
		}
	}
	if g.mode == 1 || g.mode == 2 {
		if g.guard < g.lastSeen {
			g.guard++
		} else if g.guard > g.lastSeen {
			g.guard--
		}
	} else {
		if g.guard <= 1 {
			g.score = 1
		}
		if g.guard >= 9 {
			g.score = -1
		}
		g.guard += g.score
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 18, 35, 255})
	names := []string{"PATROL", "CHASE", "SEARCH"}
	ebitenutil.DebugPrintAt(s, "GUARD MEMORY LAB", 168, 35)
	ebitenutil.DebugPrintAt(s, "SEE → REMEMBER → SEARCH → FORGET", 105, 78)
	for x := 1; x <= 9; x++ {
		px := float32(20 + x*44)
		c := color.RGBA{28, 51, 73, 255}
		if x >= 4 && x <= 6 {
			c = color.RGBA{82, 69, 94, 255}
		}
		vector.DrawFilledRect(s, px, 280, 40, 90, c, false)
	}
	vector.DrawFilledCircle(s, float32(20+g.player*44+20), 350, 16, color.RGBA{255, 194, 76, 255}, false)
	bob := float32(math.Sin(float64(g.timer)*.2) * 3)
	vector.DrawFilledCircle(s, float32(20+g.guard*44+20), 305+bob, 18, color.RGBA{231, 76, 99, 255}, false)
	vector.StrokeRect(s, float32(20+g.lastSeen*44), 275, 40, 100, 3, color.RGBA{255, 255, 255, 100}, false)
	ebitenutil.DebugPrintAt(s, "YOU", int(20+g.player*44+6), 390)
	ebitenutil.DebugPrintAt(s, "LAST SEEN", int(20+g.lastSeen*44-16), 245)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("MODE: %s   MEMORY: %02d", names[g.mode], g.timer), 145, 470)
	ebitenutil.DebugPrintAt(s, g.message, 46, 515)
	vector.DrawFilledRect(s, 10, 610, 220, 64, color.RGBA{51, 95, 139, 255}, false)
	vector.DrawFilledRect(s, 250, 610, 220, 64, color.RGBA{51, 95, 139, 255}, false)
	ebitenutil.DebugPrintAt(s, "← MOVE LEFT", 82, 638)
	ebitenutil.DebugPrintAt(s, "MOVE RIGHT →", 321, 638)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Guard Memory Lab")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
