// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"math"
)

const (
	w = 480
	h = 720
)

type region struct {
	name    string
	species [2]int
	goal    int
	c       color.RGBA
}

var regions = []region{{"TIDEPOOL", [2]int{1, 0}, 1, color.RGBA{32, 92, 116, 255}}, {"EMBER COVE", [2]int{2, 1}, 2, color.RGBA{119, 59, 48, 255}}, {"KELP FOREST", [2]int{3, 2}, 3, color.RGBA{38, 100, 66, 255}}}

type game struct {
	visits                             [3]int
	stamps                             [3]bool
	region, wild, reveal, frames, best int
	done                               bool
}

func (g *game) Update() error {
	g.frames++
	if g.done {
		if pressed() || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			b := g.best
			*g = game{best: b}
		}
		return nil
	}
	if g.reveal > 0 {
		g.reveal--
		if g.reveal == 0 && g.wild == regions[g.region].species[0] {
			g.stamps[g.region] = true
			if g.count() == 3 {
				g.done = true
				if g.best == 0 || g.frames < g.best {
					g.best = g.frames
				}
			}
		}
		return nil
	}
	choice := -1
	for i, k := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if inpututil.IsKeyJustPressed(k) {
			choice = i
		}
	}
	if x, y, ok := press(); ok && y > 560 {
		choice = min(2, x/160)
	}
	if choice >= 0 {
		g.region = choice
		r := regions[choice]
		g.wild = r.species[g.visits[choice]%2]
		g.visits[choice]++
		g.reveal = 60
	}
	return nil
}
func (g *game) count() int {
	n := 0
	for _, v := range g.stamps {
		if v {
			n++
		}
	}
	return n
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{8, 21, 38, 255})
	ebitenutil.DebugPrintAt(s, "REGION EXPEDITION", 172, 24)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("RESEARCH STAMPS %d/3  TIME %.1fs  BEST %.1fs", g.count(), float64(g.frames)/60, float64(g.best)/60), 80, 52)
	for i, r := range regions {
		x := 12 + i*156
		vector.DrawFilledRect(s, float32(x), 100, 144, 370, r.c, false)
		ebitenutil.DebugPrintAt(s, r.name, x+30, 120)
		trackatlas.DrawCentered(s, trackatlas.Species(r.species[0]), float64(x+72), 210, 90)
		ebitenutil.DebugPrintAt(s, "RARE TARGET", x+34, 266)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("VISITS %d", g.visits[i]), x+40, 298)
		label := "STAMP --"
		if g.stamps[i] {
			label = "STAMP STAR"
		}
		ebitenutil.DebugPrintAt(s, label, x+35, 332)
		ebitenutil.DebugPrintAt(s, "TABLE alternates", x+21, 375)
		ebitenutil.DebugPrintAt(s, "common / target", x+22, 399)
	}
	if g.reveal > 0 {
		p := 1 - float64(g.reveal)/60
		size := 72 + math.Sin(p*math.Pi)*70
		vector.DrawFilledRect(s, 50, 180, 380, 220, color.RGBA{5, 13, 28, 235}, false)
		trackatlas.DrawCentered(s, trackatlas.Species(g.wild), 240, 270, size)
		ebitenutil.DebugPrintAt(s, "WHO APPEARED?", 183, 350)
	}
	for i, r := range regions {
		vector.DrawFilledRect(s, float32(i*160+5), 570, 150, 80, r.c, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d %s", i+1, r.name), i*160+31, 606)
	}
	if g.done {
		vector.DrawFilledRect(s, 45, 250, 390, 150, color.RGBA{6, 15, 31, 245}, false)
		ebitenutil.DebugPrintAt(s, "ALL REGION STAMPS!", 159, 295)
		ebitenutil.DebugPrintAt(s, "TAP / ENTER TO RACE AGAIN", 130, 340)
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
func pressed() bool                        { _, _, ok := press(); return ok }
func (g *game) Layout(int, int) (int, int) { return w, h }
func main() {
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Region Expedition — Ebitengine")
	if err := ebiten.RunGame(&game{}); err != nil {
		panic(err)
	}
}
