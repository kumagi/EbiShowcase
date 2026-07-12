package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW = 480
	screenH = 720
)

var encounterNames = [...]string{"MOSSFIN", "CORALOO", "BUBBLOOM"}

type game struct {
	playerHP            int
	wildHP              int
	sleeping            bool
	basicOrbs, reefOrbs int
	sleepPowder         int
	turns               int
	encounter           int
	storage             []string
	rng                 *rand.Rand
	lastRoll            int
	message             string
	won, lost           bool
}

func newGame() *game {
	return &game{
		playerHP:    100,
		wildHP:      100,
		basicOrbs:   5,
		reefOrbs:    3,
		sleepPowder: 3,
		rng:         rand.New(rand.NewSource(8505)),
		lastRoll:    -1,
		message:     "Lower HP, add SLEEP, then choose an orb.",
	}
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	choice := -1
	for i, key := range [...]ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4} {
		if inpututil.IsKeyJustPressed(key) {
			choice = i
		}
	}
	if x, y, ok := pressPosition(); ok && y >= 480 {
		row := (y - 480) / 92
		col := x / 240
		choice = row*2 + col
	}
	if choice >= 0 && choice < 4 {
		g.takeAction(choice)
	}
	return nil
}

func (g *game) takeAction(choice int) {
	if g.turns >= 7 {
		g.lost = true
		g.message = "The wild creature fled after too many turns."
		return
	}
	g.turns++
	switch choice {
	case 0:
		before := g.wildHP
		g.wildHP = max(1, g.wildHP-22)
		g.message = fmt.Sprintf("Gentle strike: HP %d -> %d.", before, g.wildHP)
		g.wildResponse(11)
	case 1:
		if g.sleepPowder == 0 {
			g.turns--
			g.message = "No sleep powder remains."
			return
		}
		g.sleepPowder--
		roll := g.rng.Intn(100)
		if roll < 75 {
			g.sleeping = true
			g.message = fmt.Sprintf("Sleep powder succeeded (%d < 75).", roll)
		} else {
			g.message = fmt.Sprintf("Sleep powder missed (%d >= 75).", roll)
			g.wildResponse(13)
		}
	case 2:
		if g.basicOrbs == 0 {
			g.turns--
			g.message = "No basic orbs remain."
			return
		}
		g.basicOrbs--
		g.tryCapture(0)
	case 3:
		if g.reefOrbs == 0 {
			g.turns--
			g.message = "No reef orbs remain."
			return
		}
		g.reefOrbs--
		g.tryCapture(15)
	}
	if !g.won && !g.lost && g.playerHP <= 0 {
		g.lost = true
		g.message = "Your field team ran out of HP."
	}
	if !g.won && !g.lost && g.basicOrbs+g.reefOrbs == 0 {
		g.lost = true
		g.message = "All capture tools were used before storage filled."
	}
}

func (g *game) captureChance(itemBonus int) int {
	hpBonus := (100 - g.wildHP) * 55 / 100
	statusBonus := 0
	if g.sleeping {
		statusBonus = 20
	}
	return min(95, 15+hpBonus+statusBonus+itemBonus)
}

func (g *game) tryCapture(itemBonus int) {
	chance := g.captureChance(itemBonus)
	g.lastRoll = g.rng.Intn(100)
	if g.lastRoll < chance {
		name := encounterNames[g.encounter]
		g.storage = append(g.storage, name)
		g.message = fmt.Sprintf("CAPTURED %s! Roll %d < %d%%. Saved in slot %d.", name, g.lastRoll, chance, len(g.storage))
		if len(g.storage) == 3 {
			g.won = true
			return
		}
		g.encounter++
		g.wildHP = 100
		g.sleeping = false
		g.turns = 0
		g.lastRoll = -1
		return
	}
	g.message = fmt.Sprintf("Capture failed: roll %d >= %d%%. It counterattacks!", g.lastRoll, chance)
	g.sleeping = false
	g.wildResponse(18)
}

func (g *game) wildResponse(damage int) {
	if g.sleeping {
		g.message += " Sleeping creature cannot counterattack."
		return
	}
	g.playerHP = max(0, g.playerHP-damage)
	g.message += fmt.Sprintf(" Counterattack %d damage.", damage)
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{9, 20, 39, 255})
	ebitenutil.DebugPrintAt(screen, "CAPTURE PROBABILITY CAMP", 158, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ENCOUNTER %d/3   TURN %d/7   TEAM HP %d", g.encounter+1, g.turns, g.playerHP), 111, 45)
	vector.DrawFilledRect(screen, 20, 78, 440, 380, color.RGBA{26, 57, 69, 255}, false)

	name := encounterNames[min(g.encounter, len(encounterNames)-1)]
	ebitenutil.DebugPrintAt(screen, name, 211, 99)
	vector.DrawFilledCircle(screen, 240, 202, 72, color.RGBA{91, 180, 143, 255}, false)
	vector.StrokeCircle(screen, 240, 202, 72, 5, color.White, false)
	state := "AWAKE"
	if g.sleeping {
		state = "SLEEP (+20%)"
	}
	ebitenutil.DebugPrintAt(screen, state, 202, 196)
	drawBar(screen, 75, 286, 330, g.wildHP, color.RGBA{224, 88, 84, 255})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("WILD HP %d/100", g.wildHP), 188, 312)

	hpBonus := (100 - g.wildHP) * 55 / 100
	statusBonus := 0
	if g.sleeping {
		statusBonus = 20
	}
	basicChance := g.captureChance(0)
	reefChance := g.captureChance(15)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FORMULA: BASE 15 + HP %d + STATUS %d + ITEM", hpBonus, statusBonus), 84, 344)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BASIC %02d%%                         REEF %02d%%", basicChance, reefChance), 88, 373)
	drawBar(screen, 60, 397, 150, basicChance, color.RGBA{77, 157, 207, 255})
	drawBar(screen, 270, 397, 150, reefChance, color.RGBA{237, 177, 61, 255})
	rollText := "LAST ROLL: —"
	if g.lastRoll >= 0 {
		rollText = fmt.Sprintf("LAST ROLL: %d", g.lastRoll)
	}
	ebitenutil.DebugPrintAt(screen, rollText, 188, 425)

	labels := [...]string{"[1] GENTLE HIT", "[2] SLEEP", "[3] BASIC ORB", "[4] REEF ORB"}
	details := [...]string{"-22 HP / risk 11", fmt.Sprintf("75%% / left %d", g.sleepPowder), fmt.Sprintf("+0%% / left %d", g.basicOrbs), fmt.Sprintf("+15%% / left %d", g.reefOrbs)}
	for i, label := range labels {
		row, col := i/2, i%2
		x, y := col*240+6, row*92+480
		vector.DrawFilledRect(screen, float32(x), float32(y), 228, 82, color.RGBA{52, 84, 122, 255}, false)
		ebitenutil.DebugPrintAt(screen, label, x+14, y+18)
		ebitenutil.DebugPrintAt(screen, details[i], x+14, y+47)
	}
	ebitenutil.DebugPrintAt(screen, g.message, 35, 670)
	storage := "STORAGE ["
	for i := 0; i < 3; i++ {
		if i < len(g.storage) {
			storage += g.storage[i]
		} else {
			storage += "EMPTY"
		}
		if i < 2 {
			storage += " | "
		}
	}
	storage += "]"
	ebitenutil.DebugPrintAt(screen, storage, 69, 699)
	if g.won {
		overlay(screen, "STORAGE FILLED!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "CAPTURE EXPEDITION FAILED\n\nTAP / ENTER TO RETRY")
	}
}

func drawBar(screen *ebiten.Image, x, y, width, value int, c color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(width), 14, color.RGBA{43, 55, 71, 255}, false)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(width*value/100), 14, c, false)
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	_, _, ok := pressPosition()
	return ok
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 42, 270, 396, 155, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 42, 270, 396, 155, 4, color.RGBA{243, 188, 69, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 93, 327)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Capture Probability Camp — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
