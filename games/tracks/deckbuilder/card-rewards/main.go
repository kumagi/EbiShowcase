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

type card struct {
	name, rarity string
	damage       int
	upgraded     bool
	color        color.RGBA
}

var rewardTable = [][]card{
	{{"Needle", "COMMON", 6, false, color.RGBA{217, 91, 78, 255}}, {"Wave", "UNCOMMON", 9, false, color.RGBA{76, 148, 223, 255}}, {"Pearl", "RARE", 13, false, color.RGBA{169, 92, 212, 255}}},
	{{"Claw", "COMMON", 7, false, color.RGBA{217, 91, 78, 255}}, {"Current", "UNCOMMON", 10, false, color.RGBA{76, 148, 223, 255}}, {"Meteor", "RARE", 15, false, color.RGBA{169, 92, 212, 255}}},
}

type game struct {
	deck                 []card
	stage, score, target int
	selecting, clear     bool
	message              string
}

func newGame() *game {
	return &game{deck: []card{{"Jab", "STARTER", 5, false, color.RGBA{95, 174, 126, 255}}, {"Guard", "STARTER", 4, false, color.RGBA{95, 174, 126, 255}}}, target: 18, message: "Build power: tap TRAIN or press Space."}
}

func (g *game) Update() error {
	if g.clear {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}
	choice := -1
	for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
		if inpututil.IsKeyJustPressed(key) {
			choice = i
		}
	}
	if x, y, ok := pointerPressed(); ok {
		if g.selecting && y >= 360 && y < 570 {
			choice = min(2, x/160)
		}
		if !g.selecting && y >= 610 {
			g.train()
		}
	}
	if !g.selecting && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.train()
	}
	if g.selecting && choice >= 0 {
		rewards := rewardTable[min(g.stage, len(rewardTable)-1)]
		chosen := rewards[choice]
		g.deck = append(g.deck, chosen)
		g.stage++
		g.score = 0
		g.target += 12
		g.selecting = false
		g.message = chosen.name + " joined the deck. Train again!"
		if g.stage >= 2 {
			g.clear = true
		}
	}
	if g.selecting && inpututil.IsKeyJustPressed(ebiten.KeyU) {
		idx := len(g.deck) - 1
		g.deck[idx].damage += 3
		g.deck[idx].upgraded = true
		g.message = g.deck[idx].name + "+ upgraded by 3."
	}
	if g.selecting && inpututil.IsKeyJustPressed(ebiten.KeyR) && len(g.deck) > 1 {
		g.deck = g.deck[1:]
		g.message = "Removed the oldest card. A smaller deck is focused."
	}
	return nil
}

func (g *game) train() {
	for _, c := range g.deck {
		g.score += c.damage
	}
	if g.score >= g.target {
		g.selecting = true
		g.message = "Goal reached! Choose 1-3, U to upgrade, R to remove."
	} else {
		g.message = fmt.Sprintf("Deck made %d power. Train once more!", g.score)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 27, 45, 255})
	ebitenutil.DebugPrintAt(screen, "EBI DECK WORKSHOP", 175, 45)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("DECK %d CARDS   POWER %d/%d   REWARDS %d/2", len(g.deck), g.score, g.target, g.stage), 80, 82)
	ebitenutil.DebugPrintAt(screen, g.message, 45, 120)
	y := float32(160)
	for i, c := range g.deck {
		name := c.name
		if c.upgraded {
			name += "+"
		}
		vector.DrawFilledRect(screen, 45, y+float32(i*38), 390, 30, c.color, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%02d  %-10s  DMG %02d  %s", i+1, name, c.damage, c.rarity), 60, int(y)+8+i*38)
	}
	if g.selecting {
		rewards := rewardTable[min(g.stage, len(rewardTable)-1)]
		ebitenutil.DebugPrintAt(screen, "CHOOSE EXACTLY ONE REWARD", 145, 330)
		for i, c := range rewards {
			x := float32(i*160 + 5)
			vector.DrawFilledRect(screen, x, 360, 150, 200, c.color, false)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d  %s\n\nDMG %d\n%s", i+1, c.name, c.damage, c.rarity), int(x)+15, 390)
		}
		ebitenutil.DebugPrintAt(screen, "U: UPGRADE LAST   R: REMOVE OLDEST", 105, 585)
	} else {
		vector.DrawFilledRect(screen, 55, 610, 370, 70, color.RGBA{240, 177, 65, 255}, false)
		ebitenutil.DebugPrintAt(screen, "TRAIN DECK [SPACE / TAP]", 145, 640)
	}
	if g.clear {
		overlay(screen, "DECK EDITING COMPLETE!\n\nTAP / SPACE TO RESTART")
	}
}

func pointerPressed() (int, int, bool) {
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
func restartPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(screen *ebiten.Image, message string) {
	vector.DrawFilledRect(screen, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(screen, message, 110, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Card Rewards — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
