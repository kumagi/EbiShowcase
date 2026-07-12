// Command gen-track-atlas draws the shared 応用編 texture atlas (海老天 theme)
// with pure-Go software rendering. Output:
//
//	internal/trackatlas/track-atlas.png
//	web/assets/track-atlas.png
//	web/assets/track-atlas.json
//	web/assets/track-atlas-LICENSE.txt
//
// Run: go run ./cmd/gen-track-atlas
package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/kumagi/EbiShowcase/internal/tracklayout"
)

const supersample = 3

var (
	batter     = color.RGBA{255, 214, 140, 255}
	batterLite = color.RGBA{255, 236, 190, 255}
	batterDark = color.RGBA{220, 150, 70, 255}
	shrimpPink = color.RGBA{255, 150, 130, 255}
	shrimpDeep = color.RGBA{230, 90, 90, 255}
	teal       = color.RGBA{46, 201, 174, 255}
	tealDark   = color.RGBA{28, 140, 130, 255}
	coral      = color.RGBA{255, 120, 110, 255}
	gold       = color.RGBA{245, 199, 75, 255}
	ink        = color.RGBA{28, 40, 62, 255}
	foam       = color.RGBA{240, 248, 255, 255}
	sea        = color.RGBA{40, 90, 140, 255}
	seaDark    = color.RGBA{22, 50, 85, 255}
	kelp       = color.RGBA{50, 150, 90, 255}
	wood       = color.RGBA{160, 110, 70, 255}
	stone      = color.RGBA{120, 130, 150, 255}
	purple     = color.RGBA{170, 100, 210, 255}
	redGem     = color.RGBA{230, 70, 90, 255}
	blueGem    = color.RGBA{70, 140, 230, 255}
	yellowGem  = color.RGBA{240, 200, 60, 255}
	greenGem   = color.RGBA{70, 190, 110, 255}
	trashGray  = color.RGBA{140, 145, 155, 255}
)

func mix(a, b color.RGBA, t float64) color.RGBA {
	f := func(x, y uint8) uint8 { return uint8(float64(x)*(1-t) + float64(y)*t) }
	return color.RGBA{f(a.R, b.R), f(a.G, b.G), f(a.B, b.B), 255}
}

type canvas struct {
	img *image.RGBA
	s   float64
	ox  float64
	oy  float64
}

func (c *canvas) set(xp, yp int, col color.RGBA) {
	if xp < 0 || yp < 0 || xp >= c.img.Bounds().Dx() || yp >= c.img.Bounds().Dy() {
		return
	}
	c.img.SetRGBA(xp, yp, col)
}

func (c *canvas) ellipse(cx, cy, rx, ry float64, col color.RGBA) {
	cx += c.ox
	cy += c.oy
	x0 := int((cx - rx) * c.s)
	x1 := int((cx + rx) * c.s)
	y0 := int((cy - ry) * c.s)
	y1 := int((cy + ry) * c.s)
	for yp := y0; yp <= y1; yp++ {
		for xp := x0; xp <= x1; xp++ {
			fx := (float64(xp) + 0.5) / c.s
			fy := (float64(yp) + 0.5) / c.s
			nx := (fx - cx) / rx
			ny := (fy - cy) / ry
			if nx*nx+ny*ny <= 1 {
				c.set(xp, yp, col)
			}
		}
	}
}

func (c *canvas) circle(cx, cy, r float64, col color.RGBA) { c.ellipse(cx, cy, r, r, col) }

func (c *canvas) roundRect(x0, y0, x1, y1, rad float64, col color.RGBA) {
	for yp := int((y0 + c.oy) * c.s); yp <= int((y1+c.oy)*c.s); yp++ {
		for xp := int((x0 + c.ox) * c.s); xp <= int((x1+c.ox)*c.s); xp++ {
			fx := (float64(xp)+0.5)/c.s - c.ox
			fy := (float64(yp)+0.5)/c.s - c.oy
			ix := math.Max(x0+rad, math.Min(x1-rad, fx))
			iy := math.Max(y0+rad, math.Min(y1-rad, fy))
			dx, dy := fx-ix, fy-iy
			if dx*dx+dy*dy <= rad*rad {
				c.set(xp, yp, col)
			}
		}
	}
}

func (c *canvas) limb(x0, y0, x1, y1, rad float64, col color.RGBA) {
	steps := int(math.Hypot(x1-x0, y1-y0)*c.s) + 1
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		c.circle(x0+(x1-x0)*t, y0+(y1-y0)*t, rad, col)
	}
}

func (c *canvas) crispDots(cx, cy, r float64, n int, col color.RGBA) {
	for i := 0; i < n; i++ {
		a := float64(i) * 2 * math.Pi / float64(n)
		c.circle(cx+math.Cos(a)*r*0.55, cy+math.Sin(a)*r*0.4, 1.2, col)
	}
}

func drawShrimp(c *canvas, cx, cy, scale float64, body, tip color.RGBA) {
	// Curved tempura shrimp body
	for i := 0; i < 8; i++ {
		t := float64(i) / 7
		a := -0.9 + t*2.4
		x := cx + math.Cos(a)*14*scale
		y := cy + math.Sin(a)*10*scale - 2*scale
		r := (5.5 - t*2.2) * scale
		col := mix(body, tip, t*0.6)
		c.circle(x, y, r, col)
		c.circle(x-r*0.3, y-r*0.3, r*0.35, batterLite)
	}
	// Tail fan
	c.ellipse(cx-12*scale, cy+6*scale, 5*scale, 3*scale, tip)
	c.crispDots(cx, cy-2*scale, 10*scale, 5, batterDark)
}

func drawHero(c *canvas) {
	// Hooded Tenjiroh-ish silhouette in cell
	c.roundRect(16, 22, 32, 40, 4, teal)
	c.circle(24, 16, 9, color.RGBA{243, 200, 162, 255})
	c.roundRect(14, 8, 34, 18, 6, teal)
	c.circle(21, 15, 1.6, ink)
	c.circle(27, 15, 1.6, ink)
	c.roundRect(18, 20, 30, 24, 2, color.RGBA{255, 138, 92, 255})
	c.roundRect(17, 38, 23, 44, 2, color.RGBA{52, 80, 122, 255})
	c.roundRect(25, 38, 31, 44, 2, color.RGBA{52, 80, 122, 255})
}

func drawAlly(c *canvas) {
	drawShrimp(c, 24, 26, 0.85, shrimpPink, shrimpDeep)
	c.circle(28, 18, 2, ink)
}

func drawNPC(c *canvas) {
	c.ellipse(24, 28, 12, 14, kelp)
	c.circle(24, 14, 8, color.RGBA{230, 210, 170, 255})
	c.roundRect(16, 6, 32, 12, 3, wood)
	c.circle(21, 14, 1.5, ink)
	c.circle(27, 14, 1.5, ink)
}

func drawPet(c *canvas) {
	c.ellipse(24, 28, 14, 10, color.RGBA{255, 170, 190, 255})
	c.circle(24, 18, 8, color.RGBA{255, 190, 200, 255})
	c.circle(21, 17, 1.5, ink)
	c.circle(27, 17, 1.5, ink)
	c.ellipse(14, 22, 4, 6, color.RGBA{255, 150, 170, 255})
	c.ellipse(34, 22, 4, 6, color.RGBA{255, 150, 170, 255})
}

func drawFighter(c *canvas, body color.RGBA) {
	c.circle(24, 14, 7, color.RGBA{243, 200, 162, 255})
	c.roundRect(17, 20, 31, 36, 3, body)
	c.limb(18, 24, 10, 30, 2.2, body)
	c.limb(30, 24, 38, 30, 2.2, body)
	c.limb(20, 36, 18, 44, 2.2, ink)
	c.limb(28, 36, 30, 44, 2.2, ink)
	c.circle(22, 13, 1.3, ink)
	c.circle(26, 13, 1.3, ink)
}

func drawSlime(c *canvas) {
	c.ellipse(24, 28, 16, 12, color.RGBA{120, 220, 150, 255})
	c.ellipse(24, 22, 14, 10, color.RGBA{160, 240, 180, 255})
	c.circle(19, 22, 2.5, ink)
	c.circle(29, 22, 2.5, ink)
	c.circle(18, 21, 1, foam)
	c.circle(28, 21, 1, foam)
}

func drawCrab(c *canvas, col color.RGBA) {
	c.ellipse(24, 26, 14, 10, col)
	c.circle(18, 18, 5, col)
	c.circle(30, 18, 5, col)
	c.circle(17, 17, 1.5, ink)
	c.circle(31, 17, 1.5, ink)
	c.limb(10, 24, 4, 18, 2, col)
	c.limb(38, 24, 44, 18, 2, col)
	c.limb(12, 32, 8, 40, 1.8, col)
	c.limb(36, 32, 40, 40, 1.8, col)
}

func drawGhost(c *canvas, col color.RGBA) {
	c.ellipse(24, 22, 13, 14, col)
	c.roundRect(11, 22, 37, 38, 2, col)
	for i := 0; i < 4; i++ {
		c.circle(14+float64(i)*6.5, 38, 3.5, col)
	}
	c.circle(19, 20, 2.5, foam)
	c.circle(29, 20, 2.5, foam)
	c.circle(19, 20, 1.2, ink)
	c.circle(29, 20, 1.2, ink)
}

func drawScout(c *canvas) {
	c.roundRect(14, 18, 34, 40, 4, color.RGBA{80, 120, 180, 255})
	c.circle(24, 14, 8, color.RGBA{243, 200, 162, 255})
	c.roundRect(16, 8, 32, 14, 3, color.RGBA{50, 70, 110, 255})
	c.circle(21, 14, 1.4, ink)
	c.circle(27, 14, 1.4, ink)
}

func drawLeaf(c *canvas) {
	c.ellipse(24, 24, 14, 16, kelp)
	c.ellipse(24, 24, 8, 10, color.RGBA{90, 190, 120, 255})
	c.limb(24, 10, 24, 38, 1.5, wood)
	c.circle(20, 22, 1.5, ink)
	c.circle(28, 22, 1.5, ink)
}

func drawSlug(c *canvas) {
	c.ellipse(24, 30, 16, 10, purple)
	c.circle(16, 22, 5, purple)
	c.circle(14, 20, 1.5, foam)
	c.circle(18, 20, 1.5, foam)
	c.circle(14, 20, 0.8, ink)
	c.circle(18, 20, 0.8, ink)
}

func drawSwarm(c *canvas) {
	c.circle(24, 24, 14, coral)
	c.circle(18, 20, 2, ink)
	c.circle(30, 20, 2, ink)
	c.circle(24, 28, 4, shrimpDeep)
	c.crispDots(24, 24, 10, 6, batter)
}

func drawSpecies(c *canvas, idx int) {
	cols := []color.RGBA{
		color.RGBA{100, 200, 220, 255},
		color.RGBA{180, 210, 100, 255},
		color.RGBA{240, 140, 90, 255},
		color.RGBA{160, 140, 240, 255},
		gold,
	}
	col := cols[idx%len(cols)]
	c.ellipse(24, 26, 15, 12, col)
	c.circle(24, 16, 8, mix(col, foam, 0.3))
	c.circle(20, 15, 1.8, ink)
	c.circle(28, 15, 1.8, ink)
	if idx == 4 {
		c.ellipse(24, 8, 10, 4, gold)
	}
}

func drawBomb(c *canvas) {
	c.circle(24, 28, 12, ink)
	c.circle(20, 24, 3, color.RGBA{80, 90, 110, 255})
	c.limb(24, 16, 28, 8, 1.5, wood)
	c.circle(28, 7, 3, coral)
}

func drawFlame(c *canvas) {
	c.ellipse(24, 28, 10, 14, color.RGBA{255, 160, 50, 255})
	c.ellipse(24, 24, 6, 10, color.RGBA{255, 220, 80, 255})
	c.ellipse(24, 20, 3, 5, foam)
}

func drawOrb(c *canvas, col color.RGBA) {
	c.circle(24, 24, 14, col)
	c.circle(18, 18, 4, mix(col, foam, 0.55))
	c.circle(28, 28, 3, mix(col, ink, 0.2))
}

func drawStar(c *canvas, col color.RGBA) {
	for i := 0; i < 5; i++ {
		a := -math.Pi/2 + float64(i)*2*math.Pi/5
		c.limb(24, 24, 24+math.Cos(a)*16, 24+math.Sin(a)*16, 3, col)
	}
	c.circle(24, 24, 6, mix(col, foam, 0.4))
}

func drawFlag(c *canvas) {
	c.roundRect(14, 8, 18, 42, 1, wood)
	c.roundRect(18, 10, 38, 24, 2, coral)
	c.roundRect(18, 16, 34, 20, 1, foam)
}

func drawGem(c *canvas, col color.RGBA) {
	c.roundRect(10, 10, 38, 38, 6, col)
	c.roundRect(14, 14, 28, 24, 3, mix(col, foam, 0.45))
	c.circle(30, 30, 3, mix(col, ink, 0.25))
	// batter crunch rim
	c.crispDots(24, 24, 16, 8, batter)
}

func drawPeg(c *canvas) {
	c.circle(24, 24, 12, stone)
	c.circle(24, 24, 7, mix(stone, foam, 0.3))
}

func drawAura(c *canvas) {
	c.circle(24, 24, 18, color.RGBA{255, 220, 100, 90})
	c.circle(24, 24, 14, color.RGBA{255, 230, 140, 120})
	c.circle(24, 24, 6, gold)
}

func drawMerge(c *canvas, tier int) {
	r := 6.0 + float64(tier)*1.8
	cols := []color.RGBA{shrimpPink, batter, coral, gold, teal, purple, color.RGBA{255, 100, 160, 255}}
	col := cols[(tier-1)%len(cols)]
	c.circle(24, 24, r, col)
	c.circle(24-r*0.3, 24-r*0.3, r*0.35, batterLite)
	if tier >= 3 {
		drawShrimp(c, 24, 26, 0.25+float64(tier)*0.05, shrimpPink, shrimpDeep)
	}
}

func drawPulse(c *canvas) {
	c.circle(24, 24, 18, color.RGBA{100, 220, 255, 80})
	c.circle(24, 24, 12, color.RGBA{100, 220, 255, 120})
	c.circle(24, 24, 4, foam)
}

func drawTile(c *canvas, kind string) {
	switch kind {
	case "grass":
		c.roundRect(2, 2, 46, 46, 2, kelp)
		c.roundRect(2, 2, 46, 14, 2, color.RGBA{90, 190, 110, 255})
	case "grass-dark":
		c.roundRect(2, 2, 46, 46, 2, color.RGBA{30, 100, 60, 255})
		c.roundRect(2, 2, 46, 14, 2, kelp)
	case "cobble":
		c.roundRect(2, 2, 46, 46, 2, stone)
		c.roundRect(6, 6, 22, 22, 2, mix(stone, foam, 0.15))
		c.roundRect(26, 26, 42, 42, 2, mix(stone, ink, 0.1))
	case "water":
		c.roundRect(2, 2, 46, 46, 2, sea)
		c.ellipse(16, 18, 10, 4, mix(sea, foam, 0.35))
		c.ellipse(32, 30, 8, 3, mix(sea, foam, 0.25))
	case "wall":
		c.roundRect(2, 2, 46, 46, 3, seaDark)
		c.roundRect(8, 8, 40, 40, 2, color.RGBA{60, 90, 130, 255})
		c.crispDots(24, 24, 14, 5, coral)
	case "crate":
		c.roundRect(6, 6, 42, 42, 3, wood)
		c.roundRect(10, 10, 38, 38, 2, mix(wood, batter, 0.2))
		c.limb(10, 10, 38, 38, 1.5, batterDark)
		c.limb(38, 10, 10, 38, 1.5, batterDark)
	case "wood":
		c.roundRect(8, 10, 40, 38, 3, wood)
		c.limb(12, 18, 36, 18, 1, batterDark)
		c.limb(12, 28, 36, 28, 1, batterDark)
	case "stone":
		c.roundRect(8, 12, 40, 40, 4, stone)
		c.circle(20, 22, 3, mix(stone, foam, 0.2))
	case "glass":
		c.roundRect(10, 8, 38, 40, 4, color.RGBA{140, 220, 230, 220})
		c.roundRect(14, 12, 26, 22, 2, foam)
	case "lantern":
		c.roundRect(18, 8, 30, 14, 2, wood)
		c.roundRect(14, 14, 34, 36, 4, gold)
		c.roundRect(18, 18, 30, 30, 2, color.RGBA{255, 240, 180, 255})
		c.roundRect(20, 36, 28, 42, 1, wood)
	case "exit":
		c.roundRect(6, 6, 42, 42, 3, gold)
		c.roundRect(12, 12, 36, 36, 2, seaDark)
		c.roundRect(18, 18, 30, 42, 2, color.RGBA{40, 60, 90, 255})
	case "platform":
		c.roundRect(2, 18, 46, 36, 3, wood)
		c.roundRect(2, 14, 46, 22, 2, kelp)
	case "cell":
		c.roundRect(4, 4, 44, 44, 4, color.RGBA{30, 45, 70, 255})
		c.roundRect(8, 8, 40, 40, 3, color.RGBA{40, 58, 88, 255})
	}
}

func drawCard(c *canvas, kind string) {
	base := color.RGBA{45, 60, 90, 255}
	accent := coral
	switch kind {
	case "block":
		accent = blueGem
	case "skill":
		accent = purple
	}
	c.roundRect(8, 4, 40, 44, 4, base)
	c.roundRect(12, 8, 36, 20, 2, accent)
	c.roundRect(14, 24, 34, 38, 2, mix(base, foam, 0.15))
}

func drawUI(c *canvas, kind string) {
	switch kind {
	case "btn":
		c.roundRect(6, 10, 42, 38, 8, sea)
		c.roundRect(10, 14, 38, 34, 6, mix(sea, foam, 0.2))
	case "btn-accent":
		c.roundRect(6, 10, 42, 38, 8, color.RGBA{230, 130, 60, 255})
		c.roundRect(10, 14, 38, 34, 6, gold)
	case "panel":
		c.roundRect(4, 4, 44, 44, 4, color.RGBA{20, 32, 52, 240})
		c.roundRect(8, 8, 40, 40, 3, color.RGBA{35, 50, 75, 255})
	case "modal":
		c.roundRect(2, 8, 46, 40, 4, color.RGBA{8, 18, 36, 245})
		c.roundRect(6, 12, 42, 36, 3, color.RGBA{25, 40, 65, 255})
	}
}

func drawRoute(c *canvas, treasure bool) {
	c.circle(24, 24, 16, seaDark)
	if treasure {
		c.circle(24, 24, 10, gold)
		c.roundRect(18, 18, 30, 30, 2, batter)
	} else {
		c.ellipse(24, 28, 10, 6, kelp)
		c.circle(24, 18, 6, color.RGBA{230, 210, 170, 255})
	}
}

func drawBlock(c *canvas) {
	c.roundRect(6, 6, 42, 42, 4, batter)
	c.roundRect(10, 10, 38, 38, 3, batterLite)
	c.crispDots(24, 24, 12, 6, batterDark)
}

func drawBakery(c *canvas) {
	c.roundRect(8, 20, 40, 42, 3, wood)
	c.roundRect(14, 8, 34, 22, 2, coral)
	c.circle(24, 28, 10, batter)
	drawShrimp(c, 24, 30, 0.45, shrimpPink, shrimpDeep)
}

func paint(name string, c *canvas) {
	switch name {
	case "hero":
		drawHero(c)
	case "ally":
		drawAlly(c)
	case "npc":
		drawNPC(c)
	case "pet":
		drawPet(c)
	case "fighter-p1":
		drawFighter(c, teal)
	case "fighter-p2":
		drawFighter(c, coral)
	case "slime":
		drawSlime(c)
	case "king-crab", "boss-crab":
		drawCrab(c, coral)
	case "ghost-patrol":
		drawGhost(c, gold)
	case "ghost-chase":
		drawGhost(c, coral)
	case "ghost-search":
		drawGhost(c, purple)
	case "scout":
		drawScout(c)
	case "leaf-guard":
		drawLeaf(c)
	case "slug":
		drawSlug(c)
	case "swarm":
		drawSwarm(c)
	case "species-0":
		drawSpecies(c, 0)
	case "species-1":
		drawSpecies(c, 1)
	case "species-2":
		drawSpecies(c, 2)
	case "species-3":
		drawSpecies(c, 3)
	case "species-evo":
		drawSpecies(c, 4)
	case "bomb":
		drawBomb(c)
	case "flame":
		drawFlame(c)
	case "capture-orb":
		drawOrb(c, color.RGBA{100, 200, 255, 255})
	case "pearl":
		drawOrb(c, gold)
	case "coin":
		c.circle(24, 24, 12, gold)
		c.circle(24, 24, 8, batter)
		c.circle(24, 24, 4, gold)
	case "xp-gem":
		drawOrb(c, teal)
	case "power-star":
		drawStar(c, teal)
	case "upgrade-blast":
		drawOrb(c, gold)
	case "upgrade-cap":
		drawOrb(c, blueGem)
	case "upgrade-spd":
		drawOrb(c, greenGem)
	case "flag":
		drawFlag(c)
	case "gem-red":
		drawGem(c, redGem)
	case "gem-blue":
		drawGem(c, blueGem)
	case "gem-yellow":
		drawGem(c, yellowGem)
	case "gem-green":
		drawGem(c, greenGem)
	case "gem-purple":
		drawGem(c, purple)
	case "gem-trash":
		drawGem(c, trashGray)
	case "peg":
		drawPeg(c)
	case "aura":
		drawAura(c)
	case "merge-1":
		drawMerge(c, 1)
	case "merge-2":
		drawMerge(c, 2)
	case "merge-3":
		drawMerge(c, 3)
	case "merge-4":
		drawMerge(c, 4)
	case "merge-5":
		drawMerge(c, 5)
	case "merge-6":
		drawMerge(c, 6)
	case "merge-7":
		drawMerge(c, 7)
	case "pulse":
		drawPulse(c)
	case "tile-grass":
		drawTile(c, "grass")
	case "tile-grass-dark":
		drawTile(c, "grass-dark")
	case "tile-cobble":
		drawTile(c, "cobble")
	case "tile-water":
		drawTile(c, "water")
	case "tile-wall":
		drawTile(c, "wall")
	case "tile-crate":
		drawTile(c, "crate")
	case "tile-wood":
		drawTile(c, "wood")
	case "tile-stone":
		drawTile(c, "stone")
	case "tile-glass":
		drawTile(c, "glass")
	case "tile-lantern":
		drawTile(c, "lantern")
	case "tile-exit":
		drawTile(c, "exit")
	case "tile-platform":
		drawTile(c, "platform")
	case "tile-cell":
		drawTile(c, "cell")
	case "card-attack":
		drawCard(c, "attack")
	case "card-block":
		drawCard(c, "block")
	case "card-skill":
		drawCard(c, "skill")
	case "ui-btn":
		drawUI(c, "btn")
	case "ui-btn-accent":
		drawUI(c, "btn-accent")
	case "ui-panel":
		drawUI(c, "panel")
	case "ui-modal":
		drawUI(c, "modal")
	case "route-rest":
		drawRoute(c, false)
	case "route-treasure":
		drawRoute(c, true)
	case "block-cell":
		drawBlock(c)
	case "bakery":
		drawBakery(c)
	default:
		c.circle(24, 24, 10, foam)
	}
}

func downsample(src *image.RGBA, w, h, s int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var r, g, b, a, n int
			for dy := 0; dy < s; dy++ {
				for dx := 0; dx < s; dx++ {
					c := src.RGBAAt(x*s+dx, y*s+dy)
					if c.A == 0 {
						continue
					}
					r += int(c.R)
					g += int(c.G)
					b += int(c.B)
					a += int(c.A)
					n++
				}
			}
			if n == 0 {
				continue
			}
			dst.SetRGBA(x, y, color.RGBA{uint8(r / n), uint8(g / n), uint8(b / n), uint8(a / n)})
		}
	}
	return dst
}

func main() {
	fw, fh := tracklayout.FrameW, tracklayout.FrameH
	sw, sh := tracklayout.SheetW(), tracklayout.SheetH()
	hi := image.NewRGBA(image.Rect(0, 0, sw*supersample, sh*supersample))

	for _, sp := range tracklayout.Sprites {
		c := &canvas{
			img: hi,
			s:   float64(supersample),
			ox:  float64(sp.Col * fw),
			oy:  float64(sp.Row * fh),
		}
		paint(sp.Name, c)
	}

	sheet := downsample(hi, sw, sh, supersample)

	writePNG := func(path string) {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			panic(err)
		}
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := png.Encode(f, sheet); err != nil {
			panic(err)
		}
	}
	writePNG("internal/trackatlas/track-atlas.png")
	writePNG("web/assets/track-atlas.png")

	type frameMeta struct {
		Name string `json:"name"`
		X    int    `json:"x"`
		Y    int    `json:"y"`
		W    int    `json:"w"`
		H    int    `json:"h"`
	}
	meta := struct {
		FrameW  int         `json:"frameW"`
		FrameH  int         `json:"frameH"`
		Cols    int         `json:"cols"`
		SheetW  int         `json:"sheetW"`
		SheetH  int         `json:"sheetH"`
		Sprites []frameMeta `json:"sprites"`
		Theme   string      `json:"theme"`
	}{
		FrameW: fw, FrameH: fh, Cols: tracklayout.Cols,
		SheetW: sw, SheetH: sh,
		Theme: "ebi-tempura / 海老天 — shared atlas for all 15 genre tracks",
	}
	for _, sp := range tracklayout.Sprites {
		x, y, w, h := tracklayout.Rect(sp.Row, sp.Col)
		meta.Sprites = append(meta.Sprites, frameMeta{sp.Name, x, y, w, h})
	}
	jb, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile("web/assets/track-atlas.json", jb, 0o644); err != nil {
		panic(err)
	}
	license := `Ebi Showcase Track Atlas (track-atlas.png + track-atlas.json)
Dedicated to the public domain under CC0 1.0 Universal.

Generated by pure-Go software rendering (cmd/gen-track-atlas).
Theme: 海老天 (ebi tempura) props for the 15 genre tracks.
`
	if err := os.WriteFile("web/assets/track-atlas-LICENSE.txt", []byte(license), 0o644); err != nil {
		panic(err)
	}
	println("wrote track atlas", sw, "x", sh, "sprites", len(tracklayout.Sprites))
}
