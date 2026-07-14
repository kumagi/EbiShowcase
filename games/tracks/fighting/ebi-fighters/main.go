package main

import (
	"fmt"
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
	"image/color"
	"math"
	"math/rand"
)

const width, height = 480, 720

type fighter struct {
	x                            float64
	hp, attack, attackKind, stun int
	guard                        bool
}
type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	p, ai                                    fighter
	frame, round, pWins, aiWins, pInv, aiInv int
	hitstop, shake, streak, bestStreak       int
	sparks                                   []spark
	message                                  string
	matchOver                                bool
	rng                                      *rand.Rand
	audio                                    *audio.Context
	gate                                     audiolab.Gate
	pulse                                    *shaderlab.Pulse
	cam                                      cameralab.State
	badge                                    *ebiten.Image
}

func newGame() *game {
	b := ebiten.NewImage(20, 20)
	b.Fill(color.RGBA{255, 100, 80, 255})
	g := &game{round: 1, rng: rand.New(rand.NewSource(4207)), audio: audio.NewContext(audiolab.SampleRate), pulse: shaderlab.NewPulse(), cam: cameralab.State{ViewW: width, ViewH: height}, badge: b}
	g.resetRound()
	return g
}
func (g *game) resetRound() {
	g.p = fighter{x: 120, hp: 100}
	g.ai = fighter{x: 360, hp: 100}
	g.frame = 0
	g.message = fmt.Sprintf("ROUND %d — FIGHT!", g.round)
}
func (g *game) Update() error {
	if g.matchOver {
		if any() {
			best, streak := g.bestStreak, g.streak
			*g = *newGame()
			g.bestStreak, g.streak = best, streak
		}
		return nil
	}
	if g.hitstop > 0 {
		g.hitstop--
		return nil
	}
	g.frame++
	g.cam.Pos = cameralab.Vec{X: (g.p.x + g.ai.x) / 2, Y: 360}
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
	if g.pInv > 0 {
		g.pInv--
	}
	if g.aiInv > 0 {
		g.aiInv--
	}
	left := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	right := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	hit := inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	heavy := inpututil.IsKeyJustPressed(ebiten.KeyU) || inpututil.IsKeyJustPressed(ebiten.KeyZ)
	g.p.guard = ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyC)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > 560 {
			switch min(4, x/(width/5)) {
			case 0:
				left = true
			case 1:
				right = true
			case 2:
				hit = true
			case 3:
				heavy = true
			case 4:
				g.p.guard = true
			}
		}
	}
	if g.p.stun > 0 {
		g.p.stun--
	} else {
		if left {
			g.p.x -= 3
		}
		if right {
			g.p.x += 3
		}
		if hit && g.p.attack == 0 {
			g.p.attack, g.p.attackKind = 20, 1
		}
		if heavy && g.p.attack == 0 {
			g.p.attack, g.p.attackKind = 32, 2
		}
	}
	if g.ai.stun > 0 {
		g.ai.stun--
	} else {
		dist := g.ai.x - g.p.x
		g.ai.guard = false
		if dist > 100 {
			g.ai.x -= 1.8
		} else if g.frame%55 == 0 {
			if g.rng.Intn(3) == 0 {
				g.ai.attack, g.ai.attackKind = 32, 2
			} else {
				g.ai.attack, g.ai.attackKind = 20, 1
			}
		} else if g.rng.Intn(100) < 3 {
			g.ai.guard = true
		}
	}
	g.p.x = max(25, min(g.ai.x-42, g.p.x))
	g.ai.x = min(455, max(g.p.x+42, g.ai.x))
	g.updateAttack(&g.p, &g.ai, &g.aiInv, true)
	g.updateAttack(&g.ai, &g.p, &g.pInv, false)
	seconds := 45 - g.frame/60
	if g.p.hp <= 0 || g.ai.hp <= 0 || seconds <= 0 {
		if g.p.hp > g.ai.hp {
			g.pWins++
			g.streak++
			if g.streak > g.bestStreak {
				g.bestStreak = g.streak
			}
			g.message = "EBI WINS THE ROUND!"
		} else {
			g.aiWins++
			g.streak = 0
			g.message = "RIVAL WINS THE ROUND!"
		}
		if g.pWins >= 2 || g.aiWins >= 2 {
			g.matchOver = true
		} else {
			g.round++
			g.resetRound()
		}
	}
	return nil
}
func (g *game) updateAttack(a, b *fighter, inv *int, player bool) {
	if a.attack <= 0 {
		return
	}
	a.attack--
	active := 12
	reach := 105.0
	d := 14
	if a.attackKind == 2 {
		active = 16
		reach = 125
		d = 24
	}
	if a.attack == active && *inv == 0 && math.Abs(b.x-a.x) < reach {
		if b.guard {
			d = max(2, d/5)
			g.message = "GUARDED! CHIP DAMAGE"
		} else {
			b.stun = 18
			g.message = "CLEAN HIT!"
		}
		b.hp -= d
		g.gate.Arm(true)
		g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Noise, 190, .08)).Play()
		g.hitstop = 5
		if a.attackKind == 2 {
			g.hitstop = 9
			g.shake = 7
		}
		g.burst((a.x+b.x)/2, 520, 8+a.attackKind*5)
		*inv = 20
	}
	_ = player
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 24 + float64(i%8)})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	if g.pulse.Available() {
		fx := ebiten.NewImage(20, 20)
		if g.pulse.Draw(fx, g.badge, float32(g.frame)*.1) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(440, 12)
			s.DrawImage(fx, op)
		}
	}
	bg := []color.RGBA{{22, 31, 48, 255}, {45, 27, 54, 255}, {25, 52, 57, 255}}
	s.Fill(bg[(g.round-1)%3])
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.frame)*2) * 5
	}
	for i := 0; i < 7; i++ {
		vector.DrawFilledCircle(s, float32((i*83+g.round*17)%480)+float32(ox), float32(180+(i%3)*90), float32(8+g.round*2), color.RGBA{255, 220, 120, 24}, true)
	}
	vector.DrawFilledRect(s, float32(ox), 590, 480, 130, color.RGBA{59, 67, 79, 255}, false)
	draw(s, g.p, "fighter-p1", true, ox)
	draw(s, g.ai, "fighter-p2", false, ox)
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	vector.DrawFilledRect(s, 30, 45, 190, 18, color.RGBA{57, 60, 76, 255}, false)
	vector.DrawFilledRect(s, 260, 45, 190, 18, color.RGBA{57, 60, 76, 255}, false)
	vector.DrawFilledRect(s, 30, 45, float32(190*max(g.p.hp, 0)/100), 18, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 450-float32(190*max(g.ai.hp, 0)/100), 45, float32(190*max(g.ai.hp, 0)/100), 18, color.RGBA{240, 75, 91, 255}, false)
	sec := max(0, 45-g.frame/60)
	if f, e := uilab.Face("en", 16); e == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(55, 75)
		text.Draw(s, fmt.Sprintf("EBI %d  HP %03d      %02d      HP %03d  RIVAL %d", g.pWins, max(0, g.p.hp), sec, max(0, g.ai.hp), g.aiWins), f, op)
	} else {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI %d  HP %03d      %02d      HP %03d  RIVAL %d", g.pWins, max(0, g.p.hp), sec, max(0, g.ai.hp), g.aiWins), 55, 75)
	}
	ebitenutil.DebugPrintAt(s, g.message, 155, 115)
	labels := []string{"LEFT", "RIGHT", "JAB", "HEAVY", "GUARD"}
	for i, l := range labels {
		x := float32(3 + i*95)
		vector.DrawFilledRect(s, x, 620, 90, 65, color.RGBA{70, 120 + uint8(i)*18, 200, 255}, false)
		ebitenutil.DebugPrintAt(s, l, int(x)+15, 650)
	}
	if g.matchOver {
		result := "MATCH LOST"
		if g.pWins > g.aiWins {
			result = "EBI WINS THE MATCH!"
		}
		overlay(s, fmt.Sprintf("%s\nWIN STREAK %d  BEST %d\n\nTAP / SPACE TO REMATCH", result, g.streak, g.bestStreak))
	}
}
func draw(s *ebiten.Image, f fighter, sprite string, right bool, offset float64) {
	x := f.x + offset
	if f.attack > 0 {
		progress := 0.0
		if f.attackKind == 1 {
			progress = math.Sin(float64(20-f.attack)/20*math.Pi) * 18
		} else {
			progress = math.Sin(float64(32-f.attack)/32*math.Pi) * 30
		}
		if !right {
			progress = -progress
		}
		x += progress
	}
	trackatlas.DrawCentered(s, sprite, x, 540, 130)
	if f.guard {
		vector.StrokeCircle(s, float32(x), 540, 45, 6, color.RGBA{100, 165, 255, 255}, false)
	}
	if f.attack > 8 && f.attack < 16 {
		dx := 18.0
		if !right {
			dx = -95
		}
		vector.DrawFilledRect(s, float32(x+dx), 530, 77, 26, color.RGBA{255, 210, 62, 255}, false)
	}
}
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 115, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Fighters — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
