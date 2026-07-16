package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const width, height = 480, 720

//go:embed assets/survivors-reef-arena.png assets/survivors-ebi-spellblade.png assets/survivors-reef-hound.png assets/survivors-crown-crab.png
var survivorArtFS embed.FS

var (
	survivorArtOnce sync.Once
	survivorArt     map[string]*ebiten.Image
)

type mob struct {
	x, y, hp, r float64
	boss        bool
	kind, flash int
}

type gem struct{ x, y float64 }
type spark struct {
	x, y, vx, vy, life float64
	c                  color.RGBA
}

type game struct {
	px, py                             float64
	mobs                               []mob
	gems                               []gem
	rng                                *rand.Rand
	frame, life, inv, kills, xp, need  int
	level                              int
	speed, aura                        float64
	auraTick                           int
	bossSpawned, drafting, clear, over bool
	pickA, pickB, pickC                string
	sparks                             []spark
	shake, bestKills, area             int
	audio                              *audio.Context
	gate                               audiolab.Gate
	pulse                              *shaderlab.Pulse
	cam                                cameralab.State
	badge                              *ebiten.Image
}

func newGame() *game {
	loadSurvivorArt()
	badge := ebiten.NewImage(24, 24)
	badge.Fill(color.RGBA{255, 211, 62, 255})
	g := &game{
		px: 240, py: 360,
		rng:  rand.New(rand.NewSource(2306)),
		life: 4, speed: 3.6, aura: 46, auraTick: 18,
		level: 1, need: 4, audio: audiolab.Context(), pulse: shaderlab.NewPulse(), cam: cameralab.State{ViewW: width, ViewH: height}, badge: badge,
	}
	// Begin in the middle of a recognizable encounter instead of an empty room.
	g.seedOpeningMobs()
	return g
}

func (g *game) seedOpeningMobs() {
	for i := 0; i < 6; i++ {
		a := float64(i) * math.Pi / 3
		kind := i % 3
		hp := []float64{1, 2, 3}[kind]
		r := []float64{14, 17, 12}[kind]
		g.mobs = append(g.mobs, mob{x: 240 + math.Cos(a)*190, y: 370 + math.Sin(a)*230, hp: hp, r: r, kind: kind})
	}
}

func (g *game) resetRun() {
	best, audioContext, pulse, camera, badge := g.bestKills, g.audio, g.pulse, g.cam, g.badge
	*g = game{
		px: 240, py: 360,
		rng:  rand.New(rand.NewSource(2306)),
		life: 4, speed: 3.6, aura: 46, auraTick: 18,
		level: 1, need: 4, bestKills: best,
		audio: audioContext, pulse: pulse, cam: camera, badge: badge,
	}
	g.seedOpeningMobs()
}

func (g *game) Update() error {
	if g.clear || g.over {
		if restart() {
			g.resetRun()
		}
		return nil
	}
	if g.drafting {
		return g.updateDraft()
	}
	g.frame++
	if g.inv > 0 {
		g.inv--
	}
	g.movePlayer()
	g.cam.Pos = cameralab.Vec{X: g.px, Y: g.py}
	sec := g.frame / 60
	interval := max(10, 34-sec/2)
	if g.frame%interval == 0 && sec < 40 {
		a := g.rng.Float64() * math.Pi * 2
		dist := 280.0 + g.rng.Float64()*80
		kind := min(2, sec/14)
		hp := []float64{1, 2, 3}[kind]
		r := []float64{14, 17, 12}[kind]
		g.mobs = append(g.mobs, mob{x: g.px + math.Cos(a)*dist, y: g.py + math.Sin(a)*dist, hp: hp, r: r, kind: kind})
	}
	if sec >= 40 && !g.bossSpawned {
		g.bossSpawned = true
		g.mobs = append(g.mobs, mob{x: 240, y: 95, hp: 70, r: 40, boss: true, kind: 3})
	}
	g.area = min(2, sec/15)
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .03
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.shake > 0 {
		g.shake--
	}
	chase := 1.15 + float64(sec)*0.028
	for i := len(g.mobs) - 1; i >= 0; i-- {
		m := &g.mobs[i]
		d := math.Hypot(g.px-m.x, g.py-m.y)
		if d < 1 {
			d = 1
		}
		ms := chase
		if m.boss {
			ms = 0.85
			if g.frame%240 > 185 {
				ms = 2.65
			}
		} else if m.kind == 1 {
			ms *= .72
		} else if m.kind == 2 {
			ms *= 1.45
		}
		if m.flash > 0 {
			m.flash--
		}
		m.x += (g.px - m.x) / d * ms
		m.y += (g.py - m.y) / d * ms

		// Contact hurts first — standing still is unsafe.
		if d < m.r+14 && g.inv == 0 {
			g.life--
			g.inv = 50
			if g.life <= 0 {
				g.over = true
				return nil
			}
		}

		if d < g.aura && g.frame%g.auraTick == 0 {
			m.hp--
			m.flash = 5
			g.burst(m.x, m.y, color.RGBA{255, 211, 62, 255}, 3)
			if m.hp <= 0 {
				if m.boss {
					g.clear = true
					if g.kills > g.bestKills {
						g.bestKills = g.kills
					}
					g.burst(m.x, m.y, color.RGBA{255, 118, 92, 255}, 30)
				} else {
					g.kills++
					g.gems = append(g.gems, gem{m.x, m.y})
					g.burst(m.x, m.y, color.RGBA{91, 224, 255, 255}, 9)
				}
				g.shake = 4
				g.mobs = append(g.mobs[:i], g.mobs[i+1:]...)
			}
		}
	}
	for i := len(g.gems) - 1; i >= 0; i-- {
		gem := &g.gems[i]
		d := math.Hypot(g.px-gem.x, g.py-gem.y)
		if d < 120 {
			gem.x += (g.px - gem.x) * 0.12
			gem.y += (g.py - gem.y) * 0.12
		}
		if d < 22 {
			g.xp++
			g.play(700)
			g.burst(gem.x, gem.y, color.RGBA{128, 255, 210, 255}, 5)
			g.gems = append(g.gems[:i], g.gems[i+1:]...)
			if g.xp >= g.need {
				g.openDraft()
			}
		}
	}
	if sec >= 55 && g.bossSpawned {
		bossAlive := false
		for _, m := range g.mobs {
			if m.boss {
				bossAlive = true
			}
		}
		if !bossAlive {
			g.clear = true
		}
	}
	return nil
}

func (g *game) burst(x, y float64, c color.RGBA, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		speed := 1 + float64(i%4)*.45
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * speed, math.Sin(a) * speed, 22 + float64(i%12), c})
	}
}

func (g *game) movePlayer() {
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy++
	}
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		dx = float64(x) - g.px
		dy = float64(y) - g.py
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		dx = float64(x) - g.px
		dy = float64(y) - g.py
	}
	if dx != 0 || dy != 0 {
		l := math.Hypot(dx, dy)
		g.px += dx / l * g.speed
		g.py += dy / l * g.speed
	}
	g.px = clamp(g.px, 20, 460)
	g.py = clamp(g.py, 90, 690)
}

func (g *game) openDraft() {
	g.drafting = true
	opts := []string{"AURA+", "SPEED+", "LIFE+"}
	g.rng.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })
	g.pickA, g.pickB, g.pickC = opts[0], opts[1], opts[2]
}

func (g *game) updateDraft() error {
	choice := ""
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		choice = g.pickA
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		choice = g.pickB
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		choice = g.pickC
	}
	touch := inpututil.AppendJustPressedTouchIDs(nil)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(touch) > 0 {
		x, y := ebiten.CursorPosition()
		if len(touch) > 0 {
			x, y = ebiten.TouchPosition(touch[0])
		}
		switch {
		case y > 300 && y < 370 && x < 160:
			choice = g.pickA
		case y > 300 && y < 370 && x < 320:
			choice = g.pickB
		case y > 300 && y < 370:
			choice = g.pickC
		}
	}
	if choice == "" {
		return nil
	}
	switch choice {
	case "AURA+":
		g.aura = math.Min(90, g.aura+10)
		g.auraTick = max(10, g.auraTick-2)
	case "SPEED+":
		g.speed = math.Min(6.2, g.speed+0.45)
	case "LIFE+":
		g.life++
	}
	g.level++
	g.xp = 0
	g.need = min(12, 4+g.level)
	g.drafting = false
	g.inv = 40
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	drawSurvivorCover(s, "survivors-reef-arena", 0, 0, width, height)
	// One generated arena remains recognizable throughout the run. These very
	// light phase grades communicate escalation without replacing its detail.
	phaseWash := []color.RGBA{{15, 102, 116, 8}, {77, 34, 128, 24}, {150, 31, 67, 32}}[g.area]
	vector.DrawFilledRect(s, 0, 0, width, height, phaseWash, false)
	if g.pulse.Available() {
		fx := ebiten.NewImage(24, 24)
		if g.pulse.Draw(fx, g.badge, float32(g.frame)*.09) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(438, 12)
			s.DrawImage(fx, op)
		}
	}
	offsetX := 0.0
	if g.shake > 0 {
		offsetX = math.Sin(float64(g.frame)*2.4) * 4
	}
	for _, gem := range g.gems {
		vector.StrokeCircle(s, float32(gem.x+offsetX), float32(gem.y), 13, 2, color.RGBA{116, 255, 226, 185}, true)
		trackatlas.DrawCentered(s, "xp-gem", gem.x+offsetX, gem.y, 18+math.Sin(float64(g.frame)*.15+gem.x)*2)
	}
	for _, m := range g.mobs {
		asset := "survivors-reef-hound"
		size := m.r*3.3 + 12
		if m.boss {
			asset = "survivors-crown-crab"
			size = 142
		}
		if m.kind == 1 {
			size *= 1.12
		}
		if m.kind == 2 {
			size *= .82
		}
		bob := math.Sin(float64(g.frame)*.16+m.x) * 2
		vector.DrawFilledCircle(s, float32(m.x+offsetX), float32(m.y+m.r*.55), float32(size*.33), color.RGBA{3, 10, 25, 95}, true)
		drawSurvivorContain(s, asset, m.x+offsetX-size/2, m.y+bob-size/2, size, size, m.x < g.px, 1, m.flash > 0)
		if m.boss {
			if g.frame%240 > 150 && g.frame%240 <= 185 {
				vector.StrokeCircle(s, float32(m.x+offsetX), float32(m.y), float32(55+(g.frame%12)*2), 3, color.RGBA{255, 100, 80, 210}, true)
			}
			vector.DrawFilledRect(s, float32(m.x-45), float32(m.y-52), 90, 7, color.RGBA{50, 48, 70, 255}, false)
			vector.DrawFilledRect(s, float32(m.x-45), float32(m.y-52), float32(90*m.hp/70), 7, color.RGBA{255, 211, 62, 255}, false)
		}
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+offsetX), float32(p.y), float32(2+p.life/15), p.c, true)
	}
	pulse := g.aura + math.Sin(float64(g.frame)*.22)*4
	for i := 0; i < 10; i++ {
		a := float64(i)*math.Pi/4 + float64(g.frame)*.035
		if i >= 8 {
			a += math.Pi / 8
		}
		x := g.px + offsetX + math.Cos(a)*pulse
		y := g.py + math.Sin(a)*pulse
		vector.DrawFilledCircle(s, float32(x), float32(y), 4, color.RGBA{255, 235, 113, 210}, true)
		if g.frame%g.auraTick < 5 {
			vector.StrokeLine(s, float32(g.px+offsetX), float32(g.py), float32(x), float32(y), 3, color.RGBA{113, 246, 255, 180}, true)
		}
	}
	// The aura is an edge effect, not an opaque gameplay disc: the arena and
	// spellblade must remain readable through its center at thumbnail size.
	vector.StrokeCircle(s, float32(g.px+offsetX), float32(g.py), float32(pulse-7), 1, color.RGBA{95, 234, 255, 72}, true)
	vector.StrokeCircle(s, float32(g.px+offsetX), float32(g.py), float32(pulse-2), 3, color.RGBA{255, 221, 91, 78}, true)
	vector.StrokeCircle(s, float32(g.px+offsetX), float32(g.py), float32(pulse+3), 1, color.RGBA{255, 247, 176, 190}, true)
	heroBob := math.Sin(float64(g.frame)*.18) * 2
	vector.DrawFilledCircle(s, float32(g.px+offsetX), float32(g.py+17), 24, color.RGBA{3, 10, 25, 95}, true)
	heroAlpha := float32(1)
	if g.inv > 0 && g.inv%10 >= 5 {
		heroAlpha = .32
	}
	drawSurvivorContain(s, "survivors-ebi-spellblade", g.px+offsetX-38, g.py+heroBob-42, 76, 84, false, heroAlpha, false)
	if g.inv > 40 {
		vector.StrokeRect(s, 5, 5, 470, 710, 9, color.RGBA{255, 75, 91, 180}, false)
	}
	sec := g.frame / 60
	wave := min(3, sec/15+1)
	vector.DrawFilledRect(s, 14, 12, 452, 65, color.RGBA{5, 13, 29, 220}, false)
	vector.StrokeRect(s, 14, 12, 452, 65, 2, color.RGBA{255, 211, 62, 130}, false)
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(28, 24)
		text.Draw(s, fmt.Sprintf("TIME %02d/55  LV%d  XP %d/%d  LIFE %d  KILLS %03d", sec, g.level, g.xp, g.need, g.life, g.kills), face, op)
	} else {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("TIME %02d/55  LV%d  XP %d/%d  LIFE %d  KILLS %03d", sec, g.level, g.xp, g.need, g.life, g.kills), 28, 24)
	}
	msg := "WAVE 1 — KEEP MOVING"
	if wave == 2 {
		msg = "WAVE 2 — ENEMIES SPEED UP"
	}
	if wave == 3 {
		msg = "FINAL WAVE — DEFEAT THE BOSS"
	}
	ebitenutil.DebugPrintAt(s, msg, 120, 52)
	ebitenutil.DebugPrintAt(s, "MOVE: WASD / DRAG    GEMS → LEVEL UP PICK", 55, 685)
	if g.frame < 150 {
		alpha := uint8(232)
		if g.frame > 105 {
			alpha = uint8(max(0, 232-(g.frame-105)*5))
		}
		vector.DrawFilledRect(s, 53, 96, 374, 82, color.RGBA{4, 14, 31, alpha}, false)
		vector.StrokeRect(s, 53, 96, 374, 82, 3, color.RGBA{117, 242, 255, alpha}, false)
		ebitenutil.DebugPrintAt(s, "ABYSSAL REEF — SURVIVE 55 SECONDS", 105, 119)
		ebitenutil.DebugPrintAt(s, "KEEP MOVING • THE PEARL AURA ATTACKS", 89, 148)
	} else if !g.bossSpawned && sec >= 35 {
		pulseAlpha := uint8(150 + int(70*math.Abs(math.Sin(float64(g.frame)*.08))))
		vector.DrawFilledRect(s, 92, 96, 296, 46, color.RGBA{72, 8, 29, pulseAlpha}, false)
		vector.StrokeRect(s, 92, 96, 296, 46, 2, color.RGBA{255, 126, 108, pulseAlpha}, false)
		ebitenutil.DebugPrintAt(s, "CROWN CRAB APPROACHING", 156, 114)
	}
	if g.drafting {
		for i := 0; i < 9; i++ {
			x := float32(55 + i*46)
			vector.StrokeLine(s, 240, 360, x, 215, 4, color.RGBA{105, 240, 255, 110}, true)
		}
		vector.DrawFilledRect(s, 40, 250, 400, 180, color.RGBA{6, 18, 37, 235}, false)
		ebitenutil.DebugPrintAt(s, "LEVEL UP — PICK ONE", 150, 270)
		drawPick(s, 55, 310, g.pickA)
		drawPick(s, 175, 310, g.pickB)
		drawPick(s, 295, 310, g.pickC)
		ebitenutil.DebugPrintAt(s, "1/2/3 or TAP a card", 150, 400)
	}
	if g.clear {
		overlay(s, fmt.Sprintf("EBI SURVIVORS CLEAR!\nKILLS %d  BEST %d\n\nTAP / SPACE TO PLAY AGAIN", g.kills, g.bestKills))
	} else if g.over {
		overlay(s, "RUN ENDED!\n\nTAP / SPACE TO RETRY")
	}
}

func loadSurvivorArt() {
	survivorArtOnce.Do(func() {
		survivorArt = map[string]*ebiten.Image{}
		for _, name := range []string{"survivors-reef-arena", "survivors-ebi-spellblade", "survivors-reef-hound", "survivors-crown-crab"} {
			data, err := survivorArtFS.ReadFile("assets/" + name + ".png")
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			survivorArt[name] = ebiten.NewImageFromImage(decoded)
		}
	})
}

func survivorImage(name string) *ebiten.Image {
	return survivorArt[name]
}

func drawSurvivorCover(dst *ebiten.Image, name string, x, y, w, h float64) {
	img := survivorImage(name)
	b := img.Bounds()
	scale := math.Max(w/float64(b.Dx()), h/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(w-float64(b.Dx())*scale)/2, y+(h-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawSurvivorContain(dst *ebiten.Image, name string, x, y, w, h float64, mirror bool, alpha float32, flash bool) {
	img := survivorImage(name)
	b := img.Bounds()
	scale := math.Min(w/float64(b.Dx()), h/float64(b.Dy()))
	dw, dh := float64(b.Dx())*scale, float64(b.Dy())*scale
	op := &ebiten.DrawImageOptions{}
	if mirror {
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(x+(w+dw)/2, y+(h-dh)/2)
	} else {
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x+(w-dw)/2, y+(h-dh)/2)
	}
	op.Filter = ebiten.FilterLinear
	if flash {
		op.ColorScale.Scale(1, .28, .28, alpha)
	} else {
		op.ColorScale.ScaleAlpha(alpha)
	}
	dst.DrawImage(img, op)
}
func (g *game) play(hz float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, hz, .08)).Play()
}

func drawPick(s *ebiten.Image, x, y float64, label string) {
	vector.DrawFilledRect(s, float32(x), float32(y), 110, 55, color.RGBA{40, 70, 120, 255}, false)
	ebitenutil.DebugPrintAt(s, label, int(x)+28, int(y)+20)
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}

func clamp(v, l, h float64) float64 { return math.Max(l, math.Min(h, v)) }
func restart() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Survivors — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
