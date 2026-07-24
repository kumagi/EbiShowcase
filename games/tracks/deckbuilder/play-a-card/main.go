package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxfx"
	"github.com/kumagi/EbiShowcase/internal/vfxmotion"
)

const width, height = 480, 720

type card struct {
	name, kind  string
	cost, value int
	c           color.RGBA
}

var cards = []card{{"Ebi Strike", "damage", 2, 12, color.RGBA{235, 82, 78, 255}}, {"Shell Guard", "block", 1, 8, color.RGBA{76, 151, 235, 255}}, {"Herb Tea", "heal", 2, 10, color.RGBA{65, 205, 142, 255}}}

type game struct {
	hp, enemyHP, energy, block, turn int
	message                          string
	clear, over                      bool
	active                           bool
	played                           int
	flight                           vfxmotion.Proxy
	enemyReaction                    vfxmotion.Reaction
	playerReaction                   vfxmotion.Reaction
	fx                               vfxfx.System
}

func newGame() *game {
	return &game{hp: 40, enemyHP: 60, energy: 3, turn: 1, message: "Choose one card to play."}
}
func (g *game) Update() error {
	if g.active {
		g.flight.Advance()
		g.enemyReaction.Advance()
		g.playerReaction.Advance()
		g.fx.Update()
		if g.flight.Done() {
			g.finishCardMotion()
			g.active = false
		}
		return nil
	}
	if g.clear || g.over {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		choice = 2
	}
	if x, y, ok := press(); ok && y > 430 {
		choice = min(2, x/(width/3))
	}
	if choice < 0 {
		return nil
	}
	c := cards[choice]
	if g.energy < c.cost {
		g.message = "Not enough energy for " + c.name
		return nil
	}
	g.played = choice
	fromX := float64(10 + choice*157 + 73)
	fromY := 575.0
	toX, toY := 240.0, 165.0
	if c.kind != "damage" {
		toY = 330
	}
	g.flight = vfxmotion.NewProxy(choice, fromX, fromY, toX, toY, 18)
	g.active = true
	g.energy -= c.cost
	switch c.kind {
	case "damage":
		g.enemyHP -= c.value
		g.enemyReaction = vfxmotion.NewReaction(2, 5, 10)
		g.message = fmt.Sprintf("%s deals %d damage!", c.name, c.value)
	case "block":
		g.block += c.value
		g.playerReaction = vfxmotion.NewReaction(0, 4, 10)
		g.message = fmt.Sprintf("%s gives %d block!", c.name, c.value)
	case "heal":
		old := g.hp
		g.hp = min(40, g.hp+c.value)
		g.playerReaction = vfxmotion.NewReaction(0, 6, 12)
		g.message = fmt.Sprintf("%s heals %d HP!", c.name, g.hp-old)
	}
	if g.enemyHP <= 0 {
		g.clear = true
		return nil
	}
	enemyDamage := max(0, 7-g.block)
	g.hp -= enemyDamage
	if enemyDamage > 0 {
		g.playerReaction = vfxmotion.NewReaction(2, 5, 12)
	}
	g.message += fmt.Sprintf(" Enemy attacks: %d damage.", enemyDamage)
	g.block = 0
	g.energy = 3
	g.turn++
	if g.hp <= 0 {
		g.over = true
	}
	return nil
}

func (g *game) finishCardMotion() {
	c := cards[g.played]
	x, y := g.flight.Position()
	switch c.kind {
	case "damage":
		g.fx.Shockwave(x, y, 0.7, color.White, c.c)
		g.fx.Burst(x, y, 22, 3.0, c.c, true)
	case "block":
		g.fx.Ring(x, y, 0.8, c.c)
		g.fx.Burst(x, y, 14, 2.0, c.c, true)
	case "heal":
		g.fx.Shockwave(x, y, 0.55, color.White, c.c)
		g.fx.Confetti(x, y, 12)
	}
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 27, 45, 255})
	enemyX := 240.0 + g.enemyReaction.Offset(8)
	enemyColor := color.RGBA{161, 83, 197, 255}
	if g.enemyReaction.Phase() == vfxmotion.ReactionFlash {
		enemyColor = color.RGBA{250, 245, 255, 255}
	}
	vector.DrawFilledCircle(s, float32(enemyX), 165, 65, enemyColor, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("INK SLIME HP %02d/60", max(0, g.enemyHP)), 175, 75)
	playerScale := 1.0
	if g.playerReaction.Phase() != vfxmotion.ReactionDone {
		playerScale += math.Sin(float64(g.playerReaction.Frame)*0.8) * 0.08
	}
	vector.DrawFilledCircle(s, float32(240+g.playerReaction.Offset(5)), 330, float32(32*playerScale), color.RGBA{45, 225, 194, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI HP %02d/40   ENERGY %d/3   TURN %d", max(0, g.hp), g.energy, g.turn), 115, 375)
	vector.DrawFilledRect(s, 25, 400, 430, 55, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, g.message, 38, 425)
	for i, c := range cards {
		x := float32(10 + i*157)
		vector.DrawFilledRect(s, x, 480, 147, 190, c.c, false)
		vector.StrokeRect(s, x, 480, 147, 190, 3, color.White, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d  %s\n\nCOST %d\n\n%s %d", i+1, c.name, c.cost, c.kind, c.value), int(x)+18, 505)
	}
	if g.active {
		x, y := g.flight.Position()
		c := cards[g.played]
		scale := 0.75 + math.Sin(g.flight.Tween.Progress()*math.Pi)*0.2
		vector.DrawFilledRect(s, float32(x-38*scale), float32(y-50*scale), float32(76*scale), float32(100*scale), c.c, false)
		vector.StrokeRect(s, float32(x-38*scale), float32(y-50*scale), float32(76*scale), float32(100*scale), 3, color.White, false)
	}
	g.fx.Draw(s)
	if g.clear {
		overlay(s, "THE SLIME IS DEFEATED!\n\nTAP / SPACE TO RESTART")
	} else if g.over {
		overlay(s, "EBI RAN OUT OF HP!\n\nTAP / SPACE TO RETRY")
	}
}
func press() (int, int, bool) {
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
	ebiten.SetWindowTitle("Play a Card — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
