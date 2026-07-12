package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const width, height = 480, 720

type species struct {
	name       string
	baseHP     int
	baseAttack int
	color      color.RGBA
}

var speciesBook = []species{
	{"REEF SPROUT", 45, 8, color.RGBA{92, 190, 113, 255}},
	{"REEF BLOOM", 68, 20, color.RGBA{231, 129, 191, 255}},
}

type monster struct {
	speciesID int
	exp, hp   int
}

type game struct {
	partner     monster
	rivalHP     int
	turns       int
	clear, over bool
	message     string
}

func newGame() *game {
	return &game{partner: monster{hp: speciesBook[0].baseHP}, rivalHP: 60, message: "Train to level 4, evolve, then challenge the reef rival."}
}

func levelFromExp(exp int) int {
	level := 1
	for level < 9 && exp >= (level+1)*(level+1)*10 {
		level++
	}
	return level
}

func (g *game) maxHP() int {
	s := speciesBook[g.partner.speciesID]
	return s.baseHP + (levelFromExp(g.partner.exp)-1)*3
}

func (g *game) attack() int {
	s := speciesBook[g.partner.speciesID]
	return s.baseAttack + levelFromExp(g.partner.exp)*2
}

func (g *game) evolveIfReady() {
	if g.partner.speciesID == 0 && levelFromExp(g.partner.exp) >= 4 {
		// Only the species definition changes. Exp and current HP remain individual state.
		g.partner.speciesID = 1
		g.message = "Evolution! Species definition changed; EXP and current HP stayed."
	}
}

func (g *game) act(kind int) {
	if g.clear || g.over {
		return
	}
	g.turns++
	switch kind {
	case 0:
		g.partner.exp += 30
		g.partner.hp -= 6
		g.message = "+30 EXP, but training cost 6 HP."
		g.evolveIfReady()
	case 1:
		g.partner.hp = min(g.maxHP(), g.partner.hp+12)
		g.message = "Rest restored 12 HP; EXP did not change."
	case 2:
		damage := g.attack()
		g.rivalHP = max(0, g.rivalHP-damage)
		g.message = fmt.Sprintf("Challenge dealt %d using current species stats.", damage)
		if g.rivalHP == 0 {
			g.clear = true
			g.message = "Reef rival defeated — growth plan complete!"
			return
		}
		g.partner.hp -= 14
	}
	if g.partner.hp <= 0 || g.turns >= 15 {
		g.partner.hp = max(0, g.partner.hp)
		g.over = true
		g.message = "Training plan failed. Balance EXP, rest, and challenges."
	}
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.act(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.act(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.act(2)
	}
	if x, y, ok := pressPosition(); ok && y >= 565 {
		g.act(min(2, x/160))
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{13, 24, 43, 255})
	s := speciesBook[g.partner.speciesID]
	level := levelFromExp(g.partner.exp)
	ebitenutil.DebugPrintAt(screen, "REEF GROWTH TRIAL", 180, 22)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TURN %02d/15", g.turns), 202, 50)
	vector.DrawFilledCircle(screen, 145, 190, 67, s.color, true)
	vector.StrokeCircle(screen, 145, 190, 67, 4, color.White, true)
	ebitenutil.DebugPrintAt(screen, s.name, 96, 278)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("LV %d  EXP %d", level, g.partner.exp), 108, 306)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %d/%d  ATK %d", g.partner.hp, g.maxHP(), g.attack()), 90, 334)
	vector.DrawFilledCircle(screen, 355, 190, 60, color.RGBA{109, 122, 172, 255}, true)
	ebitenutil.DebugPrintAt(screen, "REEF RIVAL", 316, 278)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %d/60", g.rivalHP), 327, 306)
	next := (level + 1) * (level + 1) * 10
	if level >= 9 {
		next = g.partner.exp
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CURVE: next level needs %d EXP", next), 126, 390)
	ebitenutil.DebugPrintAt(screen, "EVOLUTION: level >= 4", 153, 423)
	ebitenutil.DebugPrintAt(screen, g.message, 50, 505)
	drawButton(screen, 10, "TRAIN [1/T]", color.RGBA{67, 151, 111, 255})
	drawButton(screen, 170, "REST [2/R]", color.RGBA{72, 151, 190, 255})
	drawButton(screen, 330, "CHALLENGE [3/C]", color.RGBA{205, 103, 79, 255})
	ebitenutil.DebugPrintAt(screen, "EXP is kept when speciesID changes.", 112, 665)
	if g.clear || g.over {
		title := "GROWTH TRIAL CLEAR!"
		if g.over {
			title = "PARTNER EXHAUSTED!"
		}
		vector.DrawFilledRect(screen, 38, 260, 404, 170, color.RGBA{5, 13, 28, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 151, 308)
		ebitenutil.DebugPrintAt(screen, g.message, 65, 347)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 394)
	}
}

func drawButton(screen *ebiten.Image, x int, label string, c color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), 565, 140, 66, c, false)
	ebitenutil.DebugPrintAt(screen, label, x+20, 592)
}

func pressPosition() (int, int, bool) {
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
func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Reef Growth Trial — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
