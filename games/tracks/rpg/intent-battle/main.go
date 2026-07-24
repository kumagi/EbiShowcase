package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
)

const width, height = 480, 720

type game struct {
	hp, enemy, turn int
	message         string
	over            bool
	active          bool
	guard           bool
	actionTween     vfxmotion.Tween
	enemyReaction   vfxmotion.Reaction
	playerReaction  vfxmotion.Reaction
	fx              vfxfx.System
}

func newGame() *game { return &game{hp: 50, enemy: 70, message: "Read NEXT, then choose."} }
func (g *game) Update() error {
	if g.active {
		g.actionTween.Advance()
		g.enemyReaction.Advance()
		g.playerReaction.Advance()
		g.fx.Update()
		if g.actionTween.Done() {
			g.active = false
		}
		return nil
	}
	if g.over {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
			*g = *newGame()
		}
		return nil
	}
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y > 520 {
			choice = x / (width / 2)
		}
	}
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y > 520 {
			choice = x / (width / 2)
		}
	}
	if choice < 0 {
		return nil
	}
	guard := choice == 1
	g.guard = guard
	g.actionTween = vfxmotion.NewTween(24)
	g.active = true
	if !guard {
		g.enemy -= 12
		g.enemyReaction = vfxmotion.NewReaction(2, 5, 12)
		g.fx.Shockwave(240, 210, 0.65, color.White, color.RGBA{255, 120, 80, 255})
		g.fx.Burst(240, 210, 20, 2.8, color.RGBA{255, 130, 80, 255}, true)
	} else {
		g.fx.Ring(120, 410, 0.75, color.RGBA{100, 190, 255, 255})
	}
	intent := g.turn % 3
	damage := []int{8, 20, 5}[intent]
	if guard {
		damage = (damage + 1) / 2
	}
	g.hp -= damage
	if damage > 0 {
		g.playerReaction = vfxmotion.NewReaction(2, 5, 12)
		g.fx.Burst(120, 410, 12, 2.1, color.RGBA{255, 90, 90, 255}, true)
		g.fx.FlashScreen(0.3, 255, 70, 70)
	}
	g.message = fmt.Sprintf("Enemy dealt %d. Guard=%v", damage, guard)
	g.turn++
	if g.enemy <= 0 || g.hp <= 0 {
		g.over = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{28, 31, 58, 255})
	bossX := 300.0 + g.enemyReaction.Offset(9)
	if g.enemyReaction.Phase() == vfxmotion.ReactionFlash {
		trackatlas.DrawTinted(s, "boss-crab", bossX, 210, 145, 1, 1, .35, 1)
	} else {
		trackatlas.DrawCentered(s, "boss-crab", bossX, 210, 145)
	}
	heroX := 120.0 + g.playerReaction.Offset(6)
	pose := hero.Pose{}
	if g.active {
		t := vfxmotion.EaseInOutCubic(g.actionTween.Progress())
		if g.guard {
			squash := math.Sin(t*math.Pi) * 0.15
			pose.ScaleX, pose.ScaleY = 1+squash, 1-squash
		} else {
			heroX += math.Sin(t*math.Pi) * 90
			pose.Rotation = -math.Sin(t*math.Pi) * 0.14
		}
	}
	hero.DrawCenteredPose(s, heroX, 410, 82, pose)
	intent := []string{"NORMAL 8", "HEAVY 20 — GUARD!", "QUICK 5"}[g.turn%3]
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("PARTY %d/50   ENEMY %d/70", max(0, g.hp), max(0, g.enemy)), 135, 40)
	ebitenutil.DebugPrintAt(s, "NEXT: "+intent, 160, 370)
	ebitenutil.DebugPrintAt(s, g.message, 130, 440)
	vector.DrawFilledRect(s, 30, 530, 200, 90, color.RGBA{45, 205, 181, 255}, false)
	vector.DrawFilledRect(s, 250, 530, 200, 90, color.RGBA{80, 145, 225, 255}, false)
	g.fx.Draw(s)
	ebitenutil.DebugPrintAt(s, "1 ATTACK", 95, 570)
	ebitenutil.DebugPrintAt(s, "2 GUARD", 315, 570)
	if g.over {
		ebitenutil.DebugPrintAt(s, "BATTLE END — SPACE / TAP", 135, 670)
	}
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Enemy Intent Battle — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
