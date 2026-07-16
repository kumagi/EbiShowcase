// Copyright 2026 Ebi Showcase contributors
// SPDX-License-Identifier: Apache-2.0

// Renderer Freedom is one autoplaying game state projected through several
// independent Draw functions. Update never knows which view is on screen.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 960
	screenHeight = 540
	boardSize    = 7
	moveTicks    = 18
)

type cell struct{ x, y int }

type move struct{ dx, dy int }

var route = []move{{1, 0}, {1, 0}, {0, 1}, {1, 0}, {0, -1}}

type game struct {
	player, playerFrom, playerTo cell
	box, boxFrom, boxTo          cell
	goal                         cell
	step, moveFrame, pause       int
	tick, clears                 int
	movingBox                    bool
}

func newGame() *game {
	g := &game{goal: cell{4, 2}}
	g.resetBoard()
	return g
}

func (g *game) resetBoard() {
	g.player = cell{1, 3}
	g.box = cell{3, 3}
	g.playerFrom, g.playerTo = g.player, g.player
	g.boxFrom, g.boxTo = g.box, g.box
	g.step, g.moveFrame, g.pause = 0, 0, 25
	g.movingBox = false
}

func blocked(c cell) bool {
	if c.x <= 0 || c.y <= 0 || c.x >= boardSize-1 || c.y >= boardSize-1 {
		return true
	}
	return (c == cell{2, 2}) || (c == cell{5, 4})
}

func add(c cell, m move) cell { return cell{c.x + m.dx, c.y + m.dy} }

func (g *game) beginMove(m move) {
	next := add(g.player, m)
	if blocked(next) {
		return
	}
	g.playerFrom, g.playerTo = g.player, next
	g.boxFrom, g.boxTo = g.box, g.box
	g.movingBox = false
	if next == g.box {
		boxNext := add(g.box, m)
		if blocked(boxNext) {
			return
		}
		g.boxTo = boxNext
		g.movingBox = true
	}
	g.moveFrame = 1
}

func (g *game) Update() error {
	g.tick++
	if g.pause > 0 {
		g.pause--
		return nil
	}
	if g.moveFrame > 0 {
		g.moveFrame++
		if g.moveFrame > moveTicks {
			g.player = g.playerTo
			if g.movingBox {
				g.box = g.boxTo
			}
			// Commit both interpolation endpoints. Leaving playerFrom at the
			// previous cell would make the following pause draw one step back,
			// then appear to teleport when the next move begins.
			g.playerFrom, g.playerTo = g.player, g.player
			g.boxFrom, g.boxTo = g.box, g.box
			g.moveFrame = 0
			g.step++
			g.pause = 8
		}
		return nil
	}
	if g.box == g.goal {
		g.clears++
		g.pause = 90
		g.resetBoard()
		g.pause = 90
		return nil
	}
	if g.step >= len(route) {
		g.resetBoard()
		return nil
	}
	g.beginMove(route[g.step])
	return nil
}

func (g *game) positions() (px, py, bx, by float64) {
	t := 0.0
	if g.moveFrame > 0 {
		t = math.Min(1, float64(g.moveFrame)/moveTicks)
		t = t * t * (3 - 2*t)
	}
	px = float64(g.playerFrom.x) + float64(g.playerTo.x-g.playerFrom.x)*t
	py = float64(g.playerFrom.y) + float64(g.playerTo.y-g.playerFrom.y)*t
	bx = float64(g.boxFrom.x) + float64(g.boxTo.x-g.boxFrom.x)*t
	by = float64(g.boxFrom.y) + float64(g.boxTo.y-g.boxFrom.y)*t
	return
}

type view struct{ x, y, w, h float32 }

func drawPanel(screen *ebiten.Image, v view, bg color.Color) {
	vector.DrawFilledRect(screen, v.x, v.y, v.w, v.h, bg, false)
	vector.StrokeRect(screen, v.x+.5, v.y+.5, v.w-1, v.h-1, 2, color.RGBA{80, 108, 160, 150}, false)
}

func boardGeometry(v view) (ox, oy, tile float32) {
	tile = float32(math.Min(float64((v.w-24)/boardSize), float64((v.h-44)/boardSize)))
	ox = v.x + (v.w-tile*boardSize)/2
	oy = v.y + 31
	return
}

func (g *game) drawTerminal(screen *ebiten.Image, v view) {
	drawPanel(screen, v, color.RGBA{4, 20, 14, 255})
	ebitenutil.DebugPrintAt(screen, "DRAW A", int(v.x+12), int(v.y+10))
	ox, oy, tile := boardGeometry(v)
	px, py, bx, by := g.positions()
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			glyph := "."
			if blocked(cell{x, y}) {
				glyph = "#"
			}
			if (cell{x, y}) == g.goal {
				glyph = "+"
			}
			ebitenutil.DebugPrintAt(screen, glyph, int(ox+float32(x)*tile+tile*.38), int(oy+float32(y)*tile+tile*.3))
		}
	}
	ebitenutil.DebugPrintAt(screen, "@", int(ox+float32(px)*tile+tile*.38), int(oy+float32(py)*tile+tile*.3))
	ebitenutil.DebugPrintAt(screen, "[]", int(ox+float32(bx)*tile+tile*.25), int(oy+float32(by)*tile+tile*.3))
}

func (g *game) drawBlueprint(screen *ebiten.Image, v view) {
	drawPanel(screen, v, color.RGBA{8, 35, 60, 255})
	ebitenutil.DebugPrintAt(screen, "DRAW B", int(v.x+12), int(v.y+10))
	ox, oy, tile := boardGeometry(v)
	line := color.RGBA{63, 199, 235, 190}
	for i := 0; i <= boardSize; i++ {
		vector.StrokeLine(screen, ox+float32(i)*tile, oy, ox+float32(i)*tile, oy+tile*boardSize, 1, line, false)
		vector.StrokeLine(screen, ox, oy+float32(i)*tile, ox+tile*boardSize, oy+float32(i)*tile, 1, line, false)
	}
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			if blocked(cell{x, y}) {
				vector.StrokeRect(screen, ox+float32(x)*tile+3, oy+float32(y)*tile+3, tile-6, tile-6, 2, color.RGBA{102, 226, 255, 255}, false)
			}
		}
	}
	px, py, bx, by := g.positions()
	vector.StrokeCircle(screen, ox+(float32(px)+.5)*tile, oy+(float32(py)+.5)*tile, tile*.28, 3, color.RGBA{255, 239, 125, 255}, false)
	vector.StrokeRect(screen, ox+float32(bx)*tile+tile*.2, oy+float32(by)*tile+tile*.2, tile*.6, tile*.6, 3, color.RGBA{255, 128, 173, 255}, false)
	vector.StrokeCircle(screen, ox+(float32(g.goal.x)+.5)*tile, oy+(float32(g.goal.y)+.5)*tile, tile*.38, 2, color.RGBA{94, 255, 194, 255}, false)
}

func (g *game) drawIllustrated(screen *ebiten.Image, v view) {
	drawPanel(screen, v, color.RGBA{18, 19, 48, 255})
	ebitenutil.DebugPrintAt(screen, "DRAW C", int(v.x+12), int(v.y+10))
	ox, oy, tile := boardGeometry(v)
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			c := color.RGBA{40, 48, 82, 255}
			if (x+y)%2 == 0 {
				c = color.RGBA{47, 57, 96, 255}
			}
			vector.DrawFilledRect(screen, ox+float32(x)*tile+1, oy+float32(y)*tile+1, tile-2, tile-2, c, false)
			if blocked(cell{x, y}) {
				vector.DrawFilledRect(screen, ox+float32(x)*tile+4, oy+float32(y)*tile+4, tile-8, tile-8, color.RGBA{82, 96, 137, 255}, false)
			}
		}
	}
	gx := ox + (float32(g.goal.x)+.5)*tile
	gy := oy + (float32(g.goal.y)+.5)*tile
	for i := 0; i < 6; i++ {
		a := float64(g.tick)*.035 + float64(i)*math.Pi/3
		r := tile * (.28 + .07*float32(math.Sin(float64(g.tick)*.05+float64(i))))
		vector.DrawFilledCircle(screen, gx+float32(math.Cos(a))*r, gy+float32(math.Sin(a))*r, 2.5, color.RGBA{92, 255, 193, 220}, false)
	}
	px, py, bx, by := g.positions()
	pcx, pcy := ox+(float32(px)+.5)*tile, oy+(float32(py)+.5)*tile
	vector.DrawFilledCircle(screen, pcx+3, pcy+5, tile*.3, color.RGBA{0, 0, 0, 80}, false)
	vector.DrawFilledCircle(screen, pcx, pcy, tile*.3, color.RGBA{255, 188, 84, 255}, false)
	vector.DrawFilledCircle(screen, pcx-tile*.1, pcy-tile*.04, 2.5, color.RGBA{30, 31, 55, 255}, false)
	vector.DrawFilledCircle(screen, pcx+tile*.1, pcy-tile*.04, 2.5, color.RGBA{30, 31, 55, 255}, false)
	vector.DrawFilledRect(screen, ox+float32(bx)*tile+tile*.18+3, oy+float32(by)*tile+tile*.18+5, tile*.64, tile*.64, color.RGBA{0, 0, 0, 80}, false)
	vector.DrawFilledRect(screen, ox+float32(bx)*tile+tile*.18, oy+float32(by)*tile+tile*.18, tile*.64, tile*.64, color.RGBA{239, 104, 121, 255}, false)
	vector.StrokeRect(screen, ox+float32(bx)*tile+tile*.26, oy+float32(by)*tile+tile*.26, tile*.48, tile*.48, 2, color.RGBA{255, 218, 155, 255}, false)
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{8, 12, 29, 255})
	if screen.Bounds().Dy() > screen.Bounds().Dx() {
		ebitenutil.DebugPrintAt(screen, "ONE GAME STATE", 18, 15)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("UPDATE TICK %05d", g.tick), 18, 35)
		g.drawTerminal(screen, view{18, 62, 504, 270})
		g.drawBlueprint(screen, view{18, 345, 504, 270})
		g.drawIllustrated(screen, view{18, 628, 504, 270})
		ebitenutil.DebugPrintAt(screen, "RULES LIVE IN UPDATE; DRAW CHANGES THE VIEW", 92, 925)
		return
	}
	ebitenutil.DebugPrintAt(screen, "ONE GAME STATE", 22, 17)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("UPDATE TICK %05d   SAME PLAYER / BOX / GOAL IN EVERY VIEW", g.tick), 22, 38)
	g.drawTerminal(screen, view{18, 72, 296, 420})
	g.drawBlueprint(screen, view{332, 72, 296, 420})
	g.drawIllustrated(screen, view{646, 72, 296, 420})
	ebitenutil.DebugPrintAt(screen, "RULES LIVE HERE; DRAW ONLY CHANGES THE VIEW", 288, 510)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideHeight > outsideWidth {
		return 540, 960
	}
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("One game state, many views")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
