package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const width, height = 480, 720

type item struct {
	name, slot             string
	price, attack, defense int
}
type game struct {
	gold, potions int
	sword, shield bool
	message       string
	clear         bool
}

var shop = []item{{"Coral Sword", "weapon", 60, 8, 0}, {"Shell Shield", "armor", 50, 0, 6}, {"Herb Potion", "bag", 20, 0, 0}}

func newGame() *game { return &game{gold: 150, message: "Buy equipment and two potions."} }
func (g *game) Update() error {
	if g.clear {
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
	if x, y, ok := press(); ok && y > 270 && y < 590 {
		choice = min(2, (y-285)/95)
		_ = x
	}
	if choice >= 0 {
		it := shop[choice]
		owned := (choice == 0 && g.sword) || (choice == 1 && g.shield)
		if owned {
			g.message = it.name + " is already equipped."
		} else if g.gold < it.price {
			g.message = "Not enough gold."
		} else {
			g.gold -= it.price
			switch choice {
			case 0:
				g.sword = true
			case 1:
				g.shield = true
			case 2:
				g.potions++
			}
			g.message = "Bought " + it.name + "!"
		}
	}
	if g.sword && g.shield && g.potions >= 2 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{24, 29, 48, 255})
	ebitenutil.DebugPrintAt(s, "EBI EQUIPMENT SHOP", 160, 35)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("GOLD %03d", g.gold), 200, 75)
	atk, def := 5, 3
	if g.sword {
		atk += 8
	}
	if g.shield {
		def += 6
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EQUIPMENT  WEAPON: %-12s  ARMOR: %-12s", yes(g.sword, "Coral Sword"), yes(g.shield, "Shell Shield")), 55, 115)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STATS  ATTACK %d   DEFENSE %d   BAG: POTION x%d/2", atk, def, g.potions), 70, 145)
	for i, it := range shop {
		y := float32(285 + i*95)
		c := color.RGBA{45, 205, 181, 255}
		if g.gold < it.price {
			c = color.RGBA{76, 82, 100, 255}
		}
		vector.DrawFilledRect(s, 45, y, 390, 75, c, false)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%d  %-14s  %3d GOLD  [%s]", i+1, it.name, it.price, it.slot), 75, int(y)+30)
	}
	vector.DrawFilledRect(s, 35, 590, 410, 75, color.RGBA{6, 18, 37, 235}, false)
	ebitenutil.DebugPrintAt(s, g.message, 65, 620)
	if g.clear {
		overlay(s, "SHOPPING COMPLETE!\n\nTAP / SPACE TO SHOP AGAIN")
	}
}
func yes(v bool, n string) string {
	if v {
		return n
	}
	return "None"
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
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Inventory Shop — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
