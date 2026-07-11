package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
)

const width, height = 480, 720

type fighter struct {
	x                float64
	hp, attack, stun int
	guard            bool
}
type game struct {
	p, ai                                    fighter
	frame, round, pWins, aiWins, pInv, aiInv int
	message                                  string
	matchOver                                bool
	rng                                      *rand.Rand
}

func newGame() *game {
	g := &game{round: 1, rng: rand.New(rand.NewSource(4207))}
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
			*g = *newGame()
		}
		return nil
	}
	g.frame++
	if g.pInv > 0 {
		g.pInv--
	}
	if g.aiInv > 0 {
		g.aiInv--
	}
	left := ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	right := ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	hit := inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyX)
	g.p.guard = ebiten.IsKeyPressed(ebiten.KeyK) || ebiten.IsKeyPressed(ebiten.KeyC)
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if y > 560 {
			switch min(3, x/(width/4)) {
			case 0:
				left = true
			case 1:
				right = true
			case 2:
				hit = true
			case 3:
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
			g.p.attack = 20
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
			g.ai.attack = 20
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
			g.message = "EBI WINS THE ROUND!"
		} else {
			g.aiWins++
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
	if a.attack == 12 && *inv == 0 && math.Abs(b.x-a.x) < 105 {
		d := 14
		if b.guard {
			d = 3
			g.message = "GUARDED! CHIP DAMAGE"
		} else {
			b.stun = 18
			g.message = "CLEAN HIT!"
		}
		b.hp -= d
		*inv = 20
	}
	_ = player
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{22, 31, 48, 255})
	vector.DrawFilledRect(s, 0, 590, 480, 130, color.RGBA{59, 67, 79, 255}, false)
	draw(s, g.p, color.RGBA{45, 225, 194, 255}, true)
	draw(s, g.ai, color.RGBA{240, 75, 91, 255}, false)
	vector.DrawFilledRect(s, 30, 45, 190, 18, color.RGBA{57, 60, 76, 255}, false)
	vector.DrawFilledRect(s, 260, 45, 190, 18, color.RGBA{57, 60, 76, 255}, false)
	vector.DrawFilledRect(s, 30, 45, float32(190*max(g.p.hp, 0)/100), 18, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 450-float32(190*max(g.ai.hp, 0)/100), 45, float32(190*max(g.ai.hp, 0)/100), 18, color.RGBA{240, 75, 91, 255}, false)
	sec := max(0, 45-g.frame/60)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI %d  HP %03d      %02d      HP %03d  RIVAL %d", g.pWins, max(0, g.p.hp), sec, max(0, g.ai.hp), g.aiWins), 55, 75)
	ebitenutil.DebugPrintAt(s, g.message, 155, 115)
	labels := []string{"LEFT", "RIGHT", "ATTACK", "GUARD"}
	for i, l := range labels {
		x := float32(4 + i*119)
		vector.DrawFilledRect(s, x, 620, 113, 65, color.RGBA{70, 125 + uint8(i)*20, 200, 255}, false)
		ebitenutil.DebugPrintAt(s, l, int(x)+28, 650)
	}
	if g.matchOver {
		result := "MATCH LOST"
		if g.pWins > g.aiWins {
			result = "EBI WINS THE MATCH!"
		}
		overlay(s, result+"\n\nTAP / SPACE TO REMATCH")
	}
}
func draw(s *ebiten.Image, f fighter, c color.RGBA, right bool) {
	x := float32(f.x)
	vector.DrawFilledCircle(s, x, 500, 23, c, false)
	vector.DrawFilledRect(s, x-18, 523, 36, 80, c, false)
	if f.guard {
		vector.StrokeCircle(s, x, 540, 45, 6, color.RGBA{100, 165, 255, 255}, false)
	}
	if f.attack > 8 && f.attack < 16 {
		dx := float32(18)
		if !right {
			dx = -95
		}
		vector.DrawFilledRect(s, x+dx, 530, 77, 26, color.RGBA{255, 210, 62, 255}, false)
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
