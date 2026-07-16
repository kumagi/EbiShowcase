package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
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
	roundOver, roundWinner                   int
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
	prepareFightingArt()
	b := ebiten.NewImage(20, 20)
	b.Fill(color.RGBA{255, 100, 80, 255})
	g := &game{round: 1, rng: rand.New(rand.NewSource(4207)), audio: audiolab.Context(), pulse: shaderlab.NewPulse(), cam: cameralab.State{ViewW: width, ViewH: height}, badge: b}
	g.resetRound()
	return g
}
func (g *game) resetRound() {
	g.p = fighter{x: 120, hp: 100}
	g.ai = fighter{x: 360, hp: 100}
	g.frame = 0
	g.pInv, g.aiInv = 0, 0
	g.hitstop, g.shake = 0, 0
	g.roundOver, g.roundWinner = 0, 0
	g.sparks = nil
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
	if g.roundOver > 0 {
		g.roundOver--
		if g.roundOver == 0 {
			if g.pWins >= 2 || g.aiWins >= 2 {
				g.matchOver = true
			} else {
				g.round++
				g.resetRound()
			}
		}
		return nil
	}
	left := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	right := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	hit := inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	heavy := inpututil.IsKeyJustPressed(ebiten.KeyU) || inpututil.IsKeyJustPressed(ebiten.KeyZ)
	g.p.guard = ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyC)
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 610 {
			switch min(4, x/(width/5)) {
			case 0:
				left = true
			case 1:
				right = true
			case 4:
				g.p.guard = true
			}
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 610 {
			switch min(4, x/(width/5)) {
			case 2:
				hit = true
			case 3:
				heavy = true
			}
		}
	}
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 610 {
			switch min(4, x/(width/5)) {
			case 0:
				left = true
			case 1:
				right = true
			case 4:
				g.p.guard = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y >= 610 {
			switch min(4, x/(width/5)) {
			case 2:
				hit = true
			case 3:
				heavy = true
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
			g.roundWinner = 1
			g.pWins++
			g.streak++
			if g.streak > g.bestStreak {
				g.bestStreak = g.streak
			}
			g.message = "EBI WINS THE ROUND!"
		} else {
			g.roundWinner = -1
			g.aiWins++
			g.streak = 0
			g.message = "RIVAL WINS THE ROUND!"
		}
		// Keep the finished state on screen so the generated KO pose and the
		// hit/hurt boxes can be inspected before the next round begins.
		g.roundOver = 90
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
	s.Fill(color.RGBA{8, 17, 34, 255})
	drawFightingArena(s)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.frame)*2) * 5
	}
	vector.DrawFilledRect(s, 0, 0, 480, 132, color.RGBA{5, 11, 26, 205}, false)
	vector.StrokeLine(s, 0, 130, 480, 130, 3, color.RGBA{244, 185, 74, 145}, false)
	draw(s, g.p, "player", true, ox, g.roundWinner < 0 && (g.roundOver > 0 || g.matchOver))
	draw(s, g.ai, "rival", false, ox, g.roundWinner > 0 && (g.roundOver > 0 || g.matchOver))
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 211, 62, 255}, true)
	}
	vector.DrawFilledRect(s, 30, 45, 190, 18, color.RGBA{57, 60, 76, 255}, false)
	vector.DrawFilledRect(s, 260, 45, 190, 18, color.RGBA{57, 60, 76, 255}, false)
	vector.DrawFilledRect(s, 30, 45, float32(190*max(g.p.hp, 0)/100), 18, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 450-float32(190*max(g.ai.hp, 0)/100), 45, float32(190*max(g.ai.hp, 0)/100), 18, color.RGBA{240, 75, 91, 255}, false)
	vector.StrokeRect(s, 28, 43, 194, 22, 2, color.RGBA{255, 255, 255, 120}, false)
	vector.StrokeRect(s, 258, 43, 194, 22, 2, color.RGBA{255, 255, 255, 120}, false)
	sec := max(0, 45-g.frame/60)
	if f, e := uilab.Face("en", 16); e == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(55, 75)
		text.Draw(s, fmt.Sprintf("EBI %d  HP %03d      %02d      HP %03d  RIVAL %d", g.pWins, max(0, g.p.hp), sec, max(0, g.ai.hp), g.aiWins), f, op)
	} else {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI %d  HP %03d      %02d      HP %03d  RIVAL %d", g.pWins, max(0, g.p.hp), sec, max(0, g.ai.hp), g.aiWins), 55, 75)
	}
	ebitenutil.DebugPrintAt(s, g.message, 155, 115)
	if g.pulse.Available() {
		fx := ebiten.NewImage(20, 20)
		if g.pulse.Draw(fx, g.badge, float32(g.frame)*.1) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(440, 12)
			s.DrawImage(fx, op)
		}
	}
	labels := []string{"LEFT", "RIGHT", "JAB", "HEAVY", "GUARD"}
	vector.DrawFilledRect(s, 0, 610, width, 110, color.RGBA{4, 11, 25, 178}, false)
	for i, l := range labels {
		x := float32(3 + i*95)
		vector.DrawFilledRect(s, x, 620, 90, 65, color.RGBA{70, 120 + uint8(i)*18, 200, 255}, false)
		vector.StrokeRect(s, x, 620, 90, 65, 2, color.RGBA{255, 255, 255, 90}, false)
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
func draw(s *ebiten.Image, f fighter, fighterName string, right bool, offset float64, ko bool) {
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
	pose := "ready"
	if ko {
		pose = "ko"
	} else if f.stun > 0 {
		pose = "hurt"
	} else if f.attack > 0 {
		pose = "attack"
	}
	vector.DrawFilledCircle(s, float32(x), 588, 42, color.RGBA{3, 8, 18, 100}, true)
	drawFighterPose(s, fighterName+"-"+pose, x, 605, 270)
	if f.guard {
		vector.StrokeCircle(s, float32(x), 500, 63, 6, color.RGBA{100, 205, 255, 235}, true)
	}
	// Hurt and attack rectangles remain a separate teaching overlay. The
	// generated pose can change without changing these collision rules.
	hurtColor := color.NRGBA{74, 224, 255, 170}
	if !right {
		hurtColor = color.NRGBA{236, 94, 187, 170}
	}
	vector.DrawFilledRect(s, float32(x-27), 432, 54, 158, color.NRGBA{hurtColor.R, hurtColor.G, hurtColor.B, 22}, false)
	vector.StrokeRect(s, float32(x-27), 432, 54, 158, 2, hurtColor, false)
	if f.attack > 8 && f.attack < 16 {
		dx := 18.0
		if !right {
			dx = -95
		}
		vector.DrawFilledRect(s, float32(x+dx), 500, 77, 42, color.NRGBA{255, 210, 62, 28}, false)
		vector.StrokeRect(s, float32(x+dx), 500, 77, 42, 3, color.NRGBA{255, 210, 62, 235}, false)
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
