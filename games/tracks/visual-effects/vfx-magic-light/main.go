// vfx-magic-light — STEP 12: radial flare with live Go knobs.
package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
	"github.com/kumagi/EbiShowcase/internal/vfxmagic"
	"github.com/kumagi/EbiShowcase/internal/vfxsprites"
)

type ray struct {
	ang, len, life, max float64
	width               float32
}

type game struct {
	shell  *vfxlive.Shell
	rays   []ray
	parts  []vfxmagic.Particle
	rng    *rand.Rand
	flash  float64
	t      float64
	ritual float64
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(12)),
		shell: vfxlive.New(
			"Sacred pillar→halo→benediction",
			[]string{
				"DrawLightPillar(whiteGold)",
				"DrawHalo(rings=3, rays={rays})",
				"mote.y += fallSlowly() // sparks={sparks}",
				"Benediction(length={len}, bloom={bloom})",
			},
			&vfxlive.Param{Key: "rays", Label: "rays", Value: 12, Min: 4, Max: 24, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "len", Label: "length", Value: 160, Min: 60, Max: 280, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "bloom", Label: "bloom", Value: 1.4, Min: 0.5, Max: 2.8, Format: "%.1f"},
			&vfxlive.Param{Key: "sparks", Label: "sparks", Value: 50, Min: 10, Max: 100, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "cast", Label: "CAST", Value: 0, Bool: true},
		),
	}
}

func (g *game) cast() {
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.49
	n := int(g.shell.Get("rays") + 0.5)
	length := g.shell.Get("len")
	base := -math.Pi / 2
	g.rays = nil
	for i := 0; i < n; i++ {
		life := 72.0
		rayLength := length
		if i%2 == 1 {
			rayLength *= 0.68
		}
		g.rays = append(g.rays, ray{
			ang: base + float64(i)*math.Pi*2/float64(n),
			len: rayLength, life: life, max: life, width: 2.4,
		})
	}
	g.parts = nil
	bloom := g.shell.Get("bloom")
	for i := 0; i < 5; i++ {
		life := 54 + float64(i)*7
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy, Rot: float64(i) * math.Pi / 8,
			Life: life, Max: life, Scale: bloom * (0.52 + float64(i)*0.2),
			Add: true, Tint: color.RGBA{255, 242, 194, 255}, Src: vfxsprites.Light,
			FadeFrom: 0.03, FadeTo: 0.82, ScaleMulFrom: 1.45, ScaleMulTo: 0.65,
		})
	}
	ns := int(g.shell.Get("sparks") + 0.5)
	for i := 0; i < ns; i++ {
		life := 55 + g.rng.Float64()*55
		g.parts = append(g.parts, vfxmagic.Particle{
			X:  cx + (g.rng.Float64()-0.5)*190,
			Y:  sy + 8 + g.rng.Float64()*sh*0.72,
			VX: (g.rng.Float64() - 0.5) * 0.16, VY: 0.12 + g.rng.Float64()*0.42,
			Spin: (g.rng.Float64() - 0.5) * 0.025,
			Life: life, Max: life, Scale: 0.11 + g.rng.Float64()*0.2, Add: true,
			Tint: color.RGBA{255, 228, 142, 255}, Src: vfxsprites.Light,
			FadeFrom: 0.02, FadeTo: 0.74, ScaleMulFrom: 0.7, ScaleMulTo: 1.2,
		})
	}
	g.flash = 0.52
	g.ritual = 96
}

func (g *game) Update() error {
	g.t += 0.05
	g.shell.Update()
	if g.shell.Bool("cast") {
		g.cast()
		g.shell.Param("cast").Value = 0
	}
	if g.flash > 0 {
		g.flash -= 0.035
	}
	if g.ritual > 0 {
		g.ritual--
	}
	rays := g.rays[:0]
	for _, r := range g.rays {
		r.life--
		r.len *= 1.006
		if r.life > 0 {
			rays = append(rays, r)
		}
	}
	g.rays = rays
	alive := g.parts[:0]
	for i := range g.parts {
		p := g.parts[i]
		if p.Update() {
			alive = append(alive, p)
		}
	}
	g.parts = alive
	return nil
}

func drawSacredPillar(dst *ebiten.Image, cx, top, floor, strength float64) {
	if strength <= 0 {
		return
	}
	layers := []struct {
		width float64
		alpha uint8
	}{
		{220, 14},
		{154, 24},
		{96, 38},
		{42, 62},
	}
	for _, layer := range layers {
		a := uint8(float64(layer.alpha) * strength)
		vector.DrawFilledRect(dst, float32(cx-layer.width/2), float32(top), float32(layer.width), float32(floor-top),
			color.RGBA{255, 235, 165, a}, false)
	}
	edge := color.RGBA{255, 238, 176, uint8(95 * strength)}
	vector.StrokeLine(dst, float32(cx-110), float32(top), float32(cx-62), float32(floor), 2, edge, false)
	vector.StrokeLine(dst, float32(cx+110), float32(top), float32(cx+62), float32(floor), 2, edge, false)
}

func drawSacredHalo(dst *ebiten.Image, cx, cy, pulse, strength float64) {
	if strength <= 0 {
		return
	}
	gold := color.RGBA{255, 224, 126, uint8(220 * strength)}
	white := color.RGBA{255, 255, 238, uint8(235 * strength)}
	for i, radius := range []float64{38, 58, 78} {
		r := radius + pulse*float64(i+1)*1.5
		width := float32(3 - i/2)
		vector.StrokeCircle(dst, float32(cx), float32(cy), float32(r), width, gold, false)
	}
	// A fixed cross and paired side rays make the silhouette ceremonial rather
	// than another spinning explosion.
	vector.StrokeLine(dst, float32(cx), float32(cy-94), float32(cx), float32(cy+94), 4, gold, false)
	vector.StrokeLine(dst, float32(cx-94), float32(cy), float32(cx+94), float32(cy), 4, gold, false)
	vector.StrokeLine(dst, float32(cx), float32(cy-82), float32(cx), float32(cy+82), 1.5, white, false)
	vector.StrokeLine(dst, float32(cx-82), float32(cy), float32(cx+82), float32(cy), 1.5, white, false)
	for i := 0; i < 5; i++ {
		y := cy - 50 + float64(i)*24
		reach := 82 + float64(4-i)*14
		vector.StrokeLine(dst, float32(cx-22), float32(cy+25), float32(cx-reach), float32(y), 2, gold, false)
		vector.StrokeLine(dst, float32(cx+22), float32(cy+25), float32(cx+reach), float32(y), 2, gold, false)
	}
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{9, 8, 25, 255})
	g.shell.FillStage(s, color.RGBA{14, 12, 38, 255})
	vfxmagic.SoftFlash(s, g.flash, 255, 242, 195)
	_, sy, _, sh := g.shell.Stage()
	cx, cy := 240.0, sy+sh*0.49
	floor := sy + sh*0.84
	strength := math.Min(1, g.ritual/26)
	drawSacredPillar(s, cx, sy, floor, strength)
	for _, r := range g.rays {
		f := r.life / r.max
		x1 := cx + math.Cos(r.ang)*r.len
		y1 := cy + math.Sin(r.ang)*r.len
		vector.StrokeLine(s, float32(cx), float32(cy), float32(x1), float32(y1), r.width*2.4, color.RGBA{255, 193, 66, uint8(55 * f)}, false)
		vector.StrokeLine(s, float32(cx), float32(cy), float32(x1), float32(y1), r.width, color.RGBA{255, 238, 174, uint8(185 * f)}, false)
	}
	drawSacredHalo(s, cx, cy, (math.Sin(g.t*2)+1)*0.5, strength)
	if strength > 0 {
		vfxmagic.DrawSprite(s, vfxsprites.Ring, cx, floor-4, 0, 2.2, 0.38,
			color.RGBA{255, 224, 130, 255}, float32(0.38*strength), true)
	}
	hero.DrawBottomCentered(s, cx, sy+sh*0.78+30, 120)
	idleGlow := 0.08
	if strength > 0 {
		idleGlow = 0.58 * strength
	}
	vfxmagic.DrawSprite(s, vfxsprites.Light, cx, cy, 0, 1.25, 1.25,
		color.RGBA{255, 246, 216, 255}, float32(idleGlow), true)
	for i := range g.parts {
		g.parts[i].Draw(s)
	}
	g.shell.Hint = "pillar → halo → falling gold motes  ·  tap CAST"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Light magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
