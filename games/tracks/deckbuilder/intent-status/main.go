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

type intent struct {
	name   string
	damage int
	block  int
	color  color.RGBA
}

var intents = []intent{
	{"ATTACK 12", 12, 0, color.RGBA{232, 88, 76, 255}},
	{"DEFEND + ATTACK 6", 6, 8, color.RGBA{74, 151, 228, 255}},
	{"HEAVY ATTACK 18", 18, 0, color.RGBA{181, 82, 205, 255}},
}

type game struct {
	hp, enemyHP, energy, block int
	poison, weak, enemyBlock   int
	turn, intentIndex          int
	message                    string
	clear, over                bool
}

func newGame() *game {
	return &game{hp: 48, enemyHP: 90, energy: 3, turn: 1, message: "Read the intent, then choose your cards."}
}

func (g *game) Update() error {
	if g.clear || g.over {
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
	end := inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if x, y, ok := pointerPressed(); ok {
		if y >= 640 {
			end = true
		} else if y >= 470 {
			choice = min(2, x/160)
		}
	}
	if choice >= 0 {
		cost := []int{1, 1, 2}[choice]
		if g.energy < cost {
			g.message = "Not enough energy."
		} else {
			g.energy -= cost
			switch choice {
			case 0:
				damage := max(0, 7-g.enemyBlock)
				g.enemyBlock = max(0, g.enemyBlock-7)
				g.enemyHP -= damage
				g.message = fmt.Sprintf("Jab dealt %d after enemy block.", damage)
			case 1:
				g.block += 9
				g.message = "Shell added 9 block."
			case 2:
				g.poison += 4
				g.weak = 2
				g.message = "Toxin added 4 poison and 2 Weak."
			}
			if g.enemyHP <= 0 {
				g.clear = true
			}
		}
	}
	if end && !g.clear {
		g.resolveEnemyTurn()
	}
	return nil
}

func (g *game) resolveEnemyTurn() {
	in := intents[g.intentIndex]
	g.enemyHP -= g.poison
	if g.poison > 0 {
		g.poison--
	}
	damage := in.damage
	if g.weak > 0 {
		damage = damage * 3 / 4
		g.weak--
	}
	taken := max(0, damage-g.block)
	g.hp -= taken
	g.enemyBlock = in.block
	g.block = 0
	g.energy = 3
	g.turn++
	g.intentIndex = (g.intentIndex + 1) % len(intents)
	g.message = fmt.Sprintf("Poison ticked; enemy dealt %d. New intent shown!", taken)
	if g.enemyHP <= 0 {
		g.clear = true
	}
	if g.hp <= 0 {
		g.over = true
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 27, 45, 255})
	in := intents[g.intentIndex]
	vector.DrawFilledCircle(screen, 240, 150, 62, color.RGBA{216, 103, 75, 255}, false)
	vector.DrawFilledRect(screen, 112, 235, 256, 55, in.color, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("OMEN CRAB  HP %02d/90  BLOCK %d", max(0, g.enemyHP), g.enemyBlock), 135, 70)
	ebitenutil.DebugPrintAt(screen, "NEXT: "+in.name, 165, 257)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TURN %d  HP %02d/48  ENERGY %d/3  BLOCK %d", g.turn, max(0, g.hp), g.energy, g.block), 92, 340)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ENEMY STATUS: POISON %d  WEAK %d", g.poison, g.weak), 125, 375)
	ebitenutil.DebugPrintAt(screen, g.message, 55, 415)
	cards := []struct {
		name, body string
		cost       int
		c          color.RGBA
	}{
		{"JAB", "DMG 7", 1, color.RGBA{226, 91, 76, 255}}, {"SHELL", "BLOCK 9", 1, color.RGBA{69, 147, 228, 255}}, {"TOXIN", "POISON 4\nWEAK 2", 2, color.RGBA{91, 183, 108, 255}},
	}
	for i, c := range cards {
		x := float32(i*160 + 5)
		shade := c.c
		if g.energy < c.cost {
			shade = color.RGBA{72, 78, 92, 255}
		}
		vector.DrawFilledRect(screen, x, 470, 150, 155, shade, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d  %s\nCOST %d\n%s", i+1, c.name, c.cost, c.body), int(x)+16, 495)
	}
	vector.DrawFilledRect(screen, 55, 642, 370, 55, color.RGBA{240, 177, 65, 255}, false)
	ebitenutil.DebugPrintAt(screen, "END TURN [E]", 185, 665)
	if g.clear {
		overlay(screen, "STATUS MASTERED!\n\nTAP / SPACE TO RESTART")
	}
	if g.over {
		overlay(screen, "THE OMEN WON!\n\nTAP / SPACE TO RETRY")
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
	ebitenutil.DebugPrintAt(screen, message, 120, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Intent and Status — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
