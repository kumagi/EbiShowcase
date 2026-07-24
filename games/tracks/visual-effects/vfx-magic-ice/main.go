// vfx-magic-ice — STEP 10: form→shatter with live Go knobs.
package main

import (
	"fmt"
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

type game struct {
	shell   *vfxlive.Shell
	parts   []vfxmagic.Particle
	rng     *rand.Rand
	phase   int
	timer   float64
	crystal float64
	flash   float64
	bloom   float64
}

func newGame() *game {
	return &game{
		rng: rand.New(rand.NewSource(10)),
		shell: vfxlive.New(
			"Aqua freeze→crack→shatter",
			[]string{
				"phase := {phase} // 0 idle 1 freeze 2 shatter",
				"frost += (target-frost) * {pull}",
				"drawIceShell(aqua); drawCracks()",
				"shatter(n={count}, grav={grav}, bloom={bloom})",
			},
			&vfxlive.Param{Key: "pull", Label: "pull", Value: 0.06, Min: 0.02, Max: 0.15, Format: "%.2f"},
			&vfxlive.Param{Key: "count", Label: "shards", Value: 48, Min: 12, Max: 90, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "grav", Label: "grav", Value: 0.12, Min: 0.02, Max: 0.3, Format: "%.2f"},
			&vfxlive.Param{Key: "bloom", Label: "bloom", Value: 1.2, Min: 0.4, Max: 2.5, Format: "%.1f"},
			&vfxlive.Param{Key: "cast", Label: "CAST", Value: 0, Bool: true},
		),
	}
}

func (g *game) start() {
	if g.phase != 0 {
		return
	}
	g.phase = 1
	g.timer = 0
	g.crystal = 0
	g.flash = 0.22
	g.bloom = 0
	g.parts = nil
	cx, cy := g.targetPosition()
	for i := 0; i < 34; i++ {
		a := g.rng.Float64() * 2 * math.Pi
		r := 60 + g.rng.Float64()*105
		life := 80.0
		src := vfxsprites.Spark
		scale := 0.14 + g.rng.Float64()*0.18
		if i%5 == 0 {
			src = vfxsprites.Ice
			scale = 0.22 + g.rng.Float64()*0.2
		}
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx + math.Cos(a)*r, Y: cy + math.Sin(a)*r*0.55,
			Rot: a, Life: life, Max: life, Scale: scale, Add: true,
			Tint: color.RGBA{105, 238, 255, 255}, Src: src,
		})
	}
}

func (g *game) shatter() {
	g.phase = 2
	g.timer = 0
	g.flash = 0.85
	g.bloom = g.shell.Get("bloom")
	n := int(g.shell.Get("count") + 0.5)
	grav := g.shell.Get("grav")
	cx, cy := g.targetPosition()
	g.parts = nil
	for i := 0; i < n; i++ {
		life := 50 + g.rng.Float64()*50
		a := g.rng.Float64() * 2 * math.Pi
		sp := 2.5 + g.rng.Float64()*6
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy, VX: math.Cos(a) * sp, VY: math.Sin(a)*sp*0.85 - 1,
			Rot: a, Spin: (g.rng.Float64() - 0.5) * 0.15,
			Life: life, Max: life, Scale: 0.4 + g.rng.Float64()*0.7,
			Add: false, Grav: grav, Tint: color.RGBA{92, 230, 250, 255}, Src: vfxsprites.Ice,
			FadeFrom: 0.2, FadeTo: 1, ScaleMulFrom: 1.1, ScaleMulTo: 0.4,
		})
	}
	for i := 0; i < 54; i++ {
		life := 25 + g.rng.Float64()*35
		a := g.rng.Float64() * 2 * math.Pi
		sp := 1 + g.rng.Float64()*4
		g.parts = append(g.parts, vfxmagic.Particle{
			X: cx, Y: cy, VX: math.Cos(a) * sp, VY: math.Sin(a)*sp - 1,
			Life: life, Max: life, Scale: 0.18 + g.rng.Float64()*0.16, Add: true,
			Tint: color.RGBA{130, 245, 255, 255}, Src: vfxsprites.Spark,
		})
	}
}

func (g *game) Update() error {
	g.shell.Update()
	if g.shell.Bool("cast") {
		g.start()
		g.shell.Param("cast").Value = 0
	}
	g.shell.SetToken("phase", fmt.Sprintf("%d", g.phase))
	cx, cy := g.targetPosition()
	pull := g.shell.Get("pull")

	switch g.phase {
	case 1:
		g.timer++
		g.crystal = math.Min(1, g.timer/46)
		for i := range g.parts {
			g.parts[i].X += (cx - g.parts[i].X) * pull
			g.parts[i].Y += (cy - g.parts[i].Y) * pull
			g.parts[i].Rot += 0.08
		}
		// Hold the complete shell long enough for the learner to read "frozen",
		// then show cracks before the shards fly.
		if g.timer >= 78 {
			g.shatter()
		}
	case 2:
		g.timer++
		if g.bloom > 0 {
			g.bloom -= 0.015
		}
		if g.timer > 100 && len(g.parts) < 5 {
			g.phase = 0
		}
	}
	if g.flash > 0 {
		g.flash -= 0.05
	}
	if g.phase == 2 {
		alive := g.parts[:0]
		for i := range g.parts {
			p := g.parts[i]
			if p.Update() {
				alive = append(alive, p)
			}
		}
		g.parts = alive
	}
	return nil
}

func (g *game) targetPosition() (float64, float64) {
	_, sy, _, sh := g.shell.Stage()
	return 240, sy + sh*0.65
}

func drawIceShard(dst *ebiten.Image, x, base, height, width, growth float64, fill, edge color.RGBA) {
	growth = math.Max(0, math.Min(1, growth))
	height *= growth
	if height < 2 {
		return
	}
	top := base - height
	for y := top; y < base; y += 2 {
		t := (y - top) / height
		half := width * 0.5
		if t < 0.2 {
			half *= t / 0.2
		}
		if t > 0.82 {
			half *= 1 - (t-0.82)*0.35
		}
		vector.DrawFilledRect(dst, float32(x-half), float32(y), float32(half*2), 2.2, fill, false)
	}
	vector.StrokeLine(dst, float32(x), float32(top), float32(x-width*0.48), float32(base), 2, edge, false)
	vector.StrokeLine(dst, float32(x), float32(top), float32(x+width*0.48), float32(base), 2, edge, false)
	vector.StrokeLine(dst, float32(x), float32(top), float32(x-width*0.1), float32(base), 1, color.RGBA{225, 255, 255, edge.A}, false)
}

func drawSnowflake(dst *ebiten.Image, x, y, radius float64, alpha uint8) {
	for i := 0; i < 3; i++ {
		a := float64(i) * math.Pi / 3
		dx, dy := math.Cos(a)*radius, math.Sin(a)*radius
		vector.StrokeLine(dst, float32(x-dx), float32(y-dy), float32(x+dx), float32(y+dy),
			1.4, color.RGBA{184, 248, 255, alpha}, false)
	}
}

func (g *game) drawFrozenShell(dst *ebiten.Image, strength float64, cracked bool) {
	if strength <= 0 {
		return
	}
	strength = math.Min(1, strength)
	cx, cy := g.targetPosition()
	fill := color.RGBA{52, 205, 232, uint8(82 + 58*strength)}
	edge := color.RGBA{135, 242, 255, uint8(155 + 90*strength)}
	base := cy + 77
	shards := []struct {
		dx, height, width float64
	}{
		{-68, 74, 34},
		{-43, 112, 47},
		{-15, 132, 54},
		{17, 124, 52},
		{46, 108, 44},
		{70, 70, 32},
	}
	for i, shard := range shards {
		delay := float64(i%3) * 0.08
		growth := math.Max(0, math.Min(1, (strength-delay)/(1-delay)))
		drawIceShard(dst, cx+shard.dx, base, shard.height, shard.width, growth, fill, edge)
	}
	drawSnowflake(dst, cx-88, cy-48, 15*strength, uint8(210*strength))
	drawSnowflake(dst, cx+88, cy-20, 12*strength, uint8(190*strength))
	drawSnowflake(dst, cx+4, cy-83, 10*strength, uint8(180*strength))
	if cracked {
		crack := color.RGBA{235, 255, 255, uint8(235 * strength)}
		vector.StrokeLine(dst, float32(cx-8), float32(cy-78), float32(cx+5), float32(cy-28), 2, crack, false)
		vector.StrokeLine(dst, float32(cx+5), float32(cy-28), float32(cx-23), float32(cy+8), 2, crack, false)
		vector.StrokeLine(dst, float32(cx+5), float32(cy-28), float32(cx+31), float32(cy-4), 2, crack, false)
		vector.StrokeLine(dst, float32(cx-23), float32(cy+8), float32(cx-11), float32(cy+49), 1.5, crack, false)
		vector.StrokeLine(dst, float32(cx+31), float32(cy-4), float32(cx+17), float32(cy+46), 1.5, crack, false)
	}
}

func (g *game) drawFrostVignette(dst *ebiten.Image, strength float64) {
	if strength <= 0 {
		return
	}
	_, sy, _, sh := g.shell.Stage()
	alpha := uint8(105 * math.Min(1, strength))
	col := color.RGBA{105, 231, 250, alpha}
	for i := 0; i < 7; i++ {
		y := sy + 24 + float64(i)*sh/7
		reach := 18 + float64((i*17)%31)
		vector.StrokeLine(dst, 0, float32(y), float32(reach), float32(y-13), 2, col, false)
		vector.StrokeLine(dst, 480, float32(y), float32(480-reach), float32(y+11), 2, col, false)
	}
	vector.StrokeLine(dst, 0, float32(sy+2), 480, float32(sy+2), 4, color.RGBA{155, 246, 255, alpha}, false)
	vector.StrokeLine(dst, 0, float32(sy+sh-2), 480, float32(sy+sh-2), 4, color.RGBA{155, 246, 255, alpha}, false)
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{2, 15, 28, 255})
	g.shell.FillStage(s, color.RGBA{3, 28, 46, 255})
	vfxmagic.SoftFlash(s, g.flash, 82, 224, 255)
	_, sy, _, sh := g.shell.Stage()
	cx, cy := g.targetPosition()
	hero.DrawBottomCentered(s, cx, sy+sh*0.75+40, 120)
	freezeStrength := 0.0
	if g.phase == 1 {
		freezeStrength = g.crystal
	} else if g.phase == 2 && g.timer < 16 {
		freezeStrength = 1 - g.timer/16
	}
	g.drawFrostVignette(s, freezeStrength)
	g.drawFrozenShell(s, freezeStrength, g.phase == 1 && g.timer >= 55)
	if g.bloom > 0 {
		vfxmagic.DrawSprite(s, vfxsprites.Ring, cx, cy, 0, g.bloom*2, g.bloom*0.8,
			color.RGBA{72, 225, 250, 255}, float32(0.58*g.bloom), true)
	}
	for i := range g.parts {
		g.parts[i].Draw(s)
	}
	g.shell.Hint = "freeze target → cracks → aqua shatter  ·  tap CAST"
	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Ice magic — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
