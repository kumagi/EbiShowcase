// Command gen-atlas draws the 海老・天次郎 (Ebi Tenjiroh) texture atlas with pure-Go software
// rendering (no GPU), so the frames are pixel-aligned, consistent, and
// transparent. It writes:
//
//	internal/heroatlas/ebi-boy-atlas.png   (embedded into the WASM games)
//	web/assets/ebi-boy-atlas.png           (downloadable)
//	web/assets/ebi-boy-atlas.json          (frame + animation metadata)
//	web/assets/ebi-boy-atlas-LICENSE.txt   (CC0 dedication)
//
// Run: go run ./cmd/gen-atlas
package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/kumagi/EbiShowcase/internal/atlaslayout"
)

const supersample = 4

// --- palette ----------------------------------------------------------------

var (
	skin       = color.RGBA{243, 200, 162, 255}
	hair       = color.RGBA{35, 46, 74, 255}
	hood       = color.RGBA{46, 201, 174, 255}
	hoodShade  = color.RGBA{32, 150, 130, 255}
	scarf      = color.RGBA{255, 138, 92, 255}
	pants      = color.RGBA{52, 80, 122, 255}
	pantsShade = color.RGBA{38, 60, 96, 255}
	shoe       = color.RGBA{232, 238, 250, 255}
	eyeCol     = color.RGBA{22, 30, 52, 255}
	slashCol   = color.RGBA{190, 245, 255, 255}
)

func mix(a, b color.RGBA, t float64) color.RGBA {
	f := func(x, y uint8) uint8 { return uint8(float64(x)*(1-t) + float64(y)*t) }
	return color.RGBA{f(a.R, b.R), f(a.G, b.G), f(a.B, b.B), 255}
}

// --- canvas (works in 96-space, scaled onto the supersampled buffer) --------

type canvas struct {
	img *image.RGBA
	s   float64
	ox  float64 // cell origin in 96-space
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
	px0 := int((x0 + c.ox) * c.s)
	px1 := int((x1 + c.ox) * c.s)
	py0 := int((y0 + c.oy) * c.s)
	py1 := int((y1 + c.oy) * c.s)
	for yp := py0; yp <= py1; yp++ {
		for xp := px0; xp <= px1; xp++ {
			fx := (float64(xp)+0.5)/c.s - c.ox
			fy := (float64(yp)+0.5)/c.s - c.oy
			ix := math.Max(x0+rad, math.Min(x1-rad, fx))
			iy := math.Max(y0+rad, math.Min(y1-rad, fy))
			dx := fx - ix
			dy := fy - iy
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

// --- pose -------------------------------------------------------------------

type pose struct {
	dx, bob, headDX      float64
	legDX, legLift       [2]float64
	armDX, armDY         [2]float64
	blink                bool
	slashX, slashY, slsh float64
	hurt                 bool
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func computePose(action, dir string, f, n int) pose {
	var p pose
	ph := 2 * math.Pi * float64(f) / float64(n)
	sw := math.Sin(ph)
	side := dir == "side"
	switch action {
	case "idle":
		p.bob = (math.Sin(ph) + 1) / 2 * 2
	case "walk", "run":
		amp, lift := 8.0, 4.0
		if action == "run" {
			amp, lift = 13, 7
			p.headDX = boolf(side, 4)
			p.armDY[0], p.armDY[1] = -4, -4
		}
		if side {
			p.legDX[1] = sw * amp
			p.legDX[0] = -sw * amp
			p.legLift[1] = maxf(0, sw) * lift
			p.legLift[0] = maxf(0, -sw) * lift
			p.armDX[1] = -sw * amp
			p.armDX[0] = sw * amp
		} else {
			p.legLift[0] = maxf(0, sw) * lift
			p.legLift[1] = maxf(0, -sw) * lift
			p.legDX[0] = -2
			p.legDX[1] = 2
			p.armDX[0] = sw * 3
			p.armDX[1] = -sw * 3
		}
		p.bob = math.Abs(math.Sin(ph)) * boolf(action == "run", 2, 1.5)
	case "attack":
		switch f {
		case 0: // windup
			p.armDX[1] = boolf(side, -10, 0)
			p.armDY[1] = boolf(!side, -8, 0)
			p.dx = -2
		case 1: // strike
			p.dx = boolf(side, 5, 0)
			if side {
				p.armDX[1] = 16
				p.armDY[1] = -6
				p.slashX, p.slashY = 70, 42
			} else if dir == "up" {
				p.armDY[1] = -14
				p.slashX, p.slashY = 52, 14
			} else {
				p.armDY[1] = 14
				p.slashX, p.slashY = 52, 74
			}
			p.slsh = 1
		case 2: // recover
			p.armDX[1] = boolf(side, 6, 0)
			p.dx = 1
		}
	case "hurt":
		p.hurt = true
		p.blink = true
		if f == 0 {
			p.dx = -6
			p.headDX = -3
		} else {
			p.dx = -2
			p.headDX = -1
		}
	}
	return p
}

func boolf(cond bool, yes float64, no ...float64) float64 {
	if cond {
		return yes
	}
	if len(no) > 0 {
		return no[0]
	}
	return 0
}

// --- character --------------------------------------------------------------

func drawChar(c *canvas, action, dir string, f, n int) {
	p := computePose(action, dir, f, n)
	tint := func(col color.RGBA) color.RGBA {
		if p.hurt {
			return mix(col, color.RGBA{255, 90, 80, 255}, 0.5)
		}
		return col
	}
	dx := p.dx
	by := -p.bob
	side := dir == "side"

	feetY, hipY := 90.0, 63.0
	hipL, hipR := 44.0, 52.0
	footL, footR := 42.0, 54.0
	shL, shR := 37.0, 59.0
	handLx, handRx, handY := 35.0, 61.0, 63.0
	headX, headY := 48.0, 30.0

	pcol := tint(pants)
	pcolFar := tint(pantsShade)
	hcol := tint(hood)
	hcolFar := tint(hoodShade)

	// Legs (index 0 far in side view).
	legCol0, legCol1 := pcol, pcol
	if side {
		legCol0 = pcolFar
	}
	c.limb(dx+hipL, by+hipY, dx+footL+p.legDX[0], by+feetY-p.legLift[0], 5, legCol0)
	c.ellipse(dx+footL+p.legDX[0], by+feetY-p.legLift[0]+1, 6, 3.5, tint(shoe))
	c.limb(dx+hipR, by+hipY, dx+footR+p.legDX[1], by+feetY-p.legLift[1], 5, legCol1)
	c.ellipse(dx+footR+p.legDX[1], by+feetY-p.legLift[1]+1, 6, 3.5, tint(shoe))

	// Far arm (behind torso) in side view.
	if side {
		c.limb(dx+shL, by+46, dx+handLx+p.armDX[0], by+handY+p.armDY[0], 4, hcolFar)
		c.circle(dx+handLx+p.armDX[0], by+handY+p.armDY[0], 3.5, tint(skin))
	}

	// Torso (hoodie) + scarf.
	c.roundRect(dx+35, by+42, dx+61, by+68, 9, hcol)
	c.roundRect(dx+37, by+58, dx+59, by+68, 6, hcolFar) // lower shade
	c.roundRect(dx+40, by+41, dx+56, by+47, 3, tint(scarf))

	// Front arm(s).
	if side {
		c.limb(dx+shR, by+46, dx+handRx+p.armDX[1], by+handY+p.armDY[1], 4.2, hcol)
		c.circle(dx+handRx+p.armDX[1], by+handY+p.armDY[1], 3.6, tint(skin))
	} else {
		c.limb(dx+shL, by+46, dx+handLx+p.armDX[0], by+handY+p.armDY[0], 4, hcol)
		c.circle(dx+handLx+p.armDX[0], by+handY+p.armDY[0], 3.5, tint(skin))
		c.limb(dx+shR, by+46, dx+handRx+p.armDX[1], by+handY+p.armDY[1], 4, hcol)
		c.circle(dx+handRx+p.armDX[1], by+handY+p.armDY[1], 3.5, tint(skin))
	}

	// Head.
	hx := dx + headX + p.headDX
	hy := by + headY
	switch dir {
	case "up":
		// Back of the head: hair only, small neck.
		c.circle(hx, hy, 15, tint(hair))
		c.circle(hx, hy-13, 5, tint(hair)) // tuft
	case "side":
		c.circle(hx-2, hy-1, 15, tint(hair))       // hair behind
		c.circle(hx+3, hy+2, 12, tint(skin))       // face
		c.ellipse(hx+13, hy+3, 2.6, 3, tint(skin)) // nose bump
		c.circle(hx-4, hy-11, 5, tint(hair))       // tuft
		if p.blink {
			c.roundRect(hx+4, hy, hx+9, hy+1.6, 0.8, eyeCol)
		} else {
			c.circle(hx+6, hy+0.5, 2, eyeCol)
		}
	default: // down
		c.circle(hx, hy-1, 15, tint(hair)) // hair
		c.circle(hx, hy+3, 13, tint(skin)) // face
		c.circle(hx, hy-12, 5, tint(hair)) // tuft
		if p.blink {
			c.roundRect(hx-8, hy+3, hx-3, hy+4.6, 0.8, eyeCol)
			c.roundRect(hx+3, hy+3, hx+8, hy+4.6, 0.8, eyeCol)
		} else {
			c.circle(hx-5, hy+3, 2, eyeCol)
			c.circle(hx+5, hy+3, 2, eyeCol)
		}
	}

	// Slash / impact for the attack strike frame.
	if p.slsh > 0 {
		sx, sy := dx+p.slashX, by+p.slashY
		c.circle(sx, sy, 8, slashCol)
		c.circle(sx+6, sy-6, 5, slashCol)
		c.circle(sx+9, sy-11, 3, slashCol)
		c.circle(sx-5, sy+6, 4, slashCol)
	}
}

// --- assembly ---------------------------------------------------------------

func main() {
	root := repoRoot()
	W := atlaslayout.SheetW() * supersample
	H := atlaslayout.SheetH() * supersample
	super := &canvas{img: image.NewRGBA(image.Rect(0, 0, W, H)), s: supersample}

	for _, a := range atlaslayout.Anims {
		for f := 0; f < a.Frames; f++ {
			x, y, _, _ := atlaslayout.Rect(a.Row, f)
			super.ox = float64(x)
			super.oy = float64(y)
			drawChar(super, a.Action, a.Dir, f, a.Frames)
		}
	}

	final := downsample(super.img, supersample)

	// PNG outputs.
	writePNG(filepath.Join(root, "internal/heroatlas/ebi-boy-atlas.png"), final)
	writePNG(filepath.Join(root, "web/assets/ebi-boy-atlas.png"), final)

	// JSON metadata.
	writeJSON(filepath.Join(root, "web/assets/ebi-boy-atlas.json"), final.Bounds())

	// License.
	must(os.WriteFile(filepath.Join(root, "web/assets/ebi-boy-atlas-LICENSE.txt"),
		[]byte(licenseText), 0o644))

	println("Generated Ebi Tenjiroh atlas:", atlaslayout.SheetW(), "x", atlaslayout.SheetH(),
		"(", len(atlaslayout.Anims), "strips )")
}

func downsample(src *image.RGBA, s int) *image.RGBA {
	W := src.Bounds().Dx() / s
	H := src.Bounds().Dy() / s
	dst := image.NewRGBA(image.Rect(0, 0, W, H))
	n := float64(s * s)
	for j := 0; j < H; j++ {
		for i := 0; i < W; i++ {
			var sr, sg, sb, covered float64
			for oy := 0; oy < s; oy++ {
				for ox := 0; ox < s; ox++ {
					c := src.RGBAAt(i*s+ox, j*s+oy)
					if c.A > 0 {
						sr += float64(c.R)
						sg += float64(c.G)
						sb += float64(c.B)
						covered++
					}
				}
			}
			if covered == 0 {
				continue
			}
			dst.SetRGBA(i, j, color.RGBA{
				R: uint8(sr / covered),
				G: uint8(sg / covered),
				B: uint8(sb / covered),
				A: uint8(255 * covered / n),
			})
		}
	}
	return dst
}

type jsonAnim struct {
	Name   string   `json:"name"`
	Action string   `json:"action"`
	Dir    string   `json:"dir"`
	FPS    int      `json:"fps"`
	Frames [][4]int `json:"frames"` // [x,y,w,h]
}

func writeJSON(path string, b image.Rectangle) {
	type doc struct {
		Image       string     `json:"image"`
		FrameWidth  int        `json:"frameWidth"`
		FrameHeight int        `json:"frameHeight"`
		SheetWidth  int        `json:"sheetWidth"`
		SheetHeight int        `json:"sheetHeight"`
		License     string     `json:"license"`
		Note        string     `json:"note"`
		Animations  []jsonAnim `json:"animations"`
	}
	d := doc{
		Image:       "ebi-boy-atlas.png",
		FrameWidth:  atlaslayout.FrameW,
		FrameHeight: atlaslayout.FrameH,
		SheetWidth:  b.Dx(),
		SheetHeight: b.Dy(),
		License:     "CC0-1.0",
		Note:        "Ebi Tenjiroh (海老・天次郎) sprite atlas. Left-facing = draw the 'side' frames flipped horizontally.",
	}
	for _, a := range atlaslayout.Anims {
		ja := jsonAnim{Name: a.Name, Action: a.Action, Dir: a.Dir, FPS: a.FPS}
		for f := 0; f < a.Frames; f++ {
			x, y, w, h := atlaslayout.Rect(a.Row, f)
			ja.Frames = append(ja.Frames, [4]int{x, y, w, h})
		}
		d.Animations = append(d.Animations, ja)
	}
	out, err := json.MarshalIndent(d, "", "  ")
	must(err)
	must(os.WriteFile(path, append(out, '\n'), 0o644))
}

func writePNG(path string, img image.Image) {
	must(os.MkdirAll(filepath.Dir(path), 0o755))
	f, err := os.Create(path)
	must(err)
	defer f.Close()
	must(png.Encode(f, img))
}

// repoRoot returns the current working directory, which must be the repo root
// (run as `go run ./cmd/gen-atlas`).
func repoRoot() string {
	dir, err := os.Getwd()
	must(err)
	return dir
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

const licenseText = `Ebi Tenjiroh (海老・天次郎) Texture Atlas (ebi-boy-atlas.png / .json)

Copyright Ebi Showcase contributors.

To the extent possible under law, the authors have dedicated this sprite atlas
to the public domain under the Creative Commons CC0 1.0 Universal dedication.
You may copy, modify, distribute, and use it, even commercially, without asking
permission. https://creativecommons.org/publicdomain/zero/1.0/

This dedication covers the generated image and metadata only. The Ebi Showcase
source code remains under the Apache License 2.0.
`
