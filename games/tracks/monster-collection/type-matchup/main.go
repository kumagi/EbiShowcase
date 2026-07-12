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

type affinity int

const (
	tide affinity = iota
	flame
	leaf
)

var names = []string{"TIDE", "FLAME", "LEAF"}
var colors = []color.RGBA{{72, 166, 220, 255}, {231, 102, 70, 255}, {92, 185, 103, 255}}

// Rows are attack types and columns are defender types.
var matchup = [3][3]float64{
	{1, 2, 0.5}, // Tide beats Flame.
	{0.5, 1, 2}, // Flame beats Leaf.
	{2, 0.5, 1}, // Leaf beats Tide.
}

var defensePlan = []affinity{leaf, tide, flame, leaf, tide, flame, leaf, tide}
var counterPlan = []affinity{flame, tide, leaf, tide, flame, leaf, flame, tide}

type game struct {
	playerHP, enemyHP int
	turn              int
	clear, over       bool
	message           string
	lastMultiplier    float64
}

func newGame() *game {
	return &game{playerHP: 65, enemyHP: 110, message: "Read defense and the predicted counter, then choose a type."}
}

func (g *game) attack(chosen affinity) {
	if g.clear || g.over {
		return
	}
	defense := defensePlan[g.turn%len(defensePlan)]
	counter := counterPlan[g.turn%len(counterPlan)]
	mult := matchup[chosen][defense]
	damage := int(14 * mult)
	g.enemyHP = max(0, g.enemyHP-damage)
	g.lastMultiplier = mult
	if g.enemyHP == 0 {
		g.clear = true
		g.message = fmt.Sprintf("%s dealt %d (x%.1f). Victory!", names[chosen], damage, mult)
		return
	}
	// The chosen attack type also becomes this turn's defensive stance.
	counterMult := matchup[counter][chosen]
	counterDamage := int(10 * counterMult)
	g.playerHP = max(0, g.playerHP-counterDamage)
	g.message = fmt.Sprintf("%s x%.1f = %d. Counter %s x%.1f = %d.", names[chosen], mult, damage, names[counter], counterMult, counterDamage)
	g.turn++
	if g.playerHP == 0 || g.turn >= len(defensePlan) {
		g.over = true
		if g.playerHP == 0 {
			g.message = "Your partner fainted. Predict both sides of the table!"
		} else {
			g.message = "Eight turns passed. Use strong matchups more often!"
		}
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
		g.attack(tide)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) || inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.attack(flame)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.attack(leaf)
	}
	if x, y, ok := pressPosition(); ok && y >= 565 {
		g.attack(affinity(min(2, x/160)))
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{13, 22, 41, 255})
	ebitenutil.DebugPrintAt(screen, "AFFINITY FORECAST DUEL", 158, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TURN %d/8", g.turn+1), 208, 48)
	defense := defensePlan[min(g.turn, len(defensePlan)-1)]
	counter := counterPlan[min(g.turn, len(counterPlan)-1)]
	vector.DrawFilledCircle(screen, 120, 165, 55, color.RGBA{91, 205, 220, 255}, true)
	vector.StrokeCircle(screen, 120, 165, 55, 4, color.White, true)
	ebitenutil.DebugPrintAt(screen, "YOUR PARTNER", 75, 230)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %02d/65", g.playerHP), 91, 254)
	vector.DrawFilledCircle(screen, 360, 165, 55, colors[defense], true)
	vector.StrokeCircle(screen, 360, 165, 55, 4, color.White, true)
	ebitenutil.DebugPrintAt(screen, "MIMIC GUARD", 318, 230)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %03d/110", g.enemyHP), 326, 254)

	vector.DrawFilledRect(screen, 40, 290, 400, 92, color.RGBA{31, 49, 74, 255}, false)
	ebitenutil.DebugPrintAt(screen, "ENEMY DEFENSE", 68, 310)
	ebitenutil.DebugPrintAt(screen, names[defense], 75, 342)
	ebitenutil.DebugPrintAt(screen, "PREDICTED COUNTER", 260, 310)
	ebitenutil.DebugPrintAt(screen, names[counter], 307, 342)

	ebitenutil.DebugPrintAt(screen, "TABLE: ATTACK ROW -> DEFENSE COLUMN", 105, 408)
	ebitenutil.DebugPrintAt(screen, "          TIDE   FLAME  LEAF", 115, 436)
	for a := 0; a < 3; a++ {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%-6s    x%.1f    x%.1f    x%.1f", names[a], matchup[a][0], matchup[a][1], matchup[a][2]), 90, 462+a*24)
	}
	ebitenutil.DebugPrintAt(screen, g.message, 48, 540)
	drawButton(screen, 10, "TIDE [1/T]", tide)
	drawButton(screen, 170, "FLAME [2/F]", flame)
	drawButton(screen, 330, "LEAF [3/L]", leaf)
	ebitenutil.DebugPrintAt(screen, "Your chosen attack also becomes your counter defense.", 58, 660)
	if g.clear || g.over {
		title := "MATCHUP MASTERED!"
		if g.over {
			title = "PARTNER FAINTED!"
		}
		vector.DrawFilledRect(screen, 38, 270, 404, 170, color.RGBA{5, 13, 28, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 157, 316)
		ebitenutil.DebugPrintAt(screen, g.message, 63, 352)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 399)
	}
}

func drawButton(screen *ebiten.Image, x int, label string, kind affinity) {
	vector.DrawFilledRect(screen, float32(x), 565, 140, 66, colors[kind], false)
	ebitenutil.DebugPrintAt(screen, label, x+24, 592)
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
	ebiten.SetWindowTitle("Affinity Forecast Duel — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
