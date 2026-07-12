package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width, height = 480, 720
	goalXP        = 9
	maxActions    = 8
)

type species struct {
	name   string
	maxHP  int
	attack int
	color  color.RGBA
}

// speciesTable is definition data. Training never changes these values.
var speciesTable = []species{
	{"REEFLET", 10, 3, color.RGBA{75, 187, 207, 255}},
	{"MOSSHELL", 14, 2, color.RGBA{89, 181, 112, 255}},
}

type creature struct {
	nickname  string
	speciesID int
	hp        int
	xp        int
}

type game struct {
	creatures   [2]creature
	selected    int
	actions     int
	frames      int
	clear, over bool
	message     string
}

func newGame() *game {
	return &game{
		creatures: [2]creature{
			{nickname: "AQUA", speciesID: 0, hp: 10},
			{nickname: "BUBBLE", speciesID: 0, hp: 6},
		},
		message: "Same REEFLET species, different HP. Train both to 9 XP.",
	}
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.frames >= 45*60 {
		g.over = true
		g.message = "Training time ended. Compare both instance states."
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.selected = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.selected = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.train()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.rest()
	}
	if x, y, ok := pressPosition(); ok {
		if y >= 180 && y < 430 {
			if x < width/2 {
				g.selected = 0
			} else {
				g.selected = 1
			}
		} else if y >= 500 && y < 610 {
			if x < width/2 {
				g.train()
			} else {
				g.rest()
			}
		}
	}
	return nil
}

func (g *game) train() {
	c := &g.creatures[g.selected]
	if c.xp >= goalXP {
		g.message = c.nickname + " already reached the XP goal. Select the other twin."
		return
	}
	c.hp -= 2
	c.xp += 3
	g.actions++
	g.message = fmt.Sprintf("%s trained: HP -2, XP +3. Species data stayed fixed.", c.nickname)
	if c.hp <= 0 {
		c.hp = 0
		g.over = true
		g.message = c.nickname + " fainted. Rest an injured instance before training."
		return
	}
	g.checkGoal()
}

func (g *game) rest() {
	c := &g.creatures[g.selected]
	definition := speciesTable[c.speciesID]
	before := c.hp
	c.hp = min(definition.maxHP, c.hp+3)
	g.actions++
	g.message = fmt.Sprintf("%s rested: HP %d -> %d (capped by species MaxHP %d).", c.nickname, before, c.hp, definition.maxHP)
	g.checkGoal()
}

func (g *game) checkGoal() {
	if g.creatures[0].xp >= goalXP && g.creatures[1].xp >= goalXP {
		g.clear = true
		g.message = "Both individuals grew while sharing one unchanged species definition!"
		return
	}
	if g.actions >= maxActions {
		g.over = true
		g.message = "Eight actions used. Read each individual's HP before choosing."
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{11, 23, 41, 255})
	ebitenutil.DebugPrintAt(screen, "SPECIES & INSTANCE TRAINING", 148, 20)
	definition := speciesTable[0]
	vector.DrawFilledRect(screen, 40, 52, 400, 100, color.RGBA{30, 55, 75, 255}, false)
	ebitenutil.DebugPrintAt(screen, "SHARED / IMMUTABLE SPECIES DEFINITION", 119, 65)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ID 0  %s   MAX HP %d   ATTACK %d", definition.name, definition.maxHP, definition.attack), 114, 98)
	ebitenutil.DebugPrintAt(screen, "Both cards point to SpeciesID 0", 138, 126)

	for i, c := range g.creatures {
		x := 25 + i*230
		fill := color.RGBA{32, 55, 76, 255}
		if i == g.selected {
			fill = color.RGBA{46, 91, 116, 255}
		}
		vector.DrawFilledRect(screen, float32(x), 180, 205, 250, fill, false)
		vector.StrokeRect(screen, float32(x), 180, 205, 250, 4, definition.color, false)
		vector.DrawFilledCircle(screen, float32(x+102), 248, 45, definition.color, false)
		vector.StrokeCircle(screen, float32(x+102), 248, 45, 4, color.White, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d: %s", i+1, c.nickname), x+69, 310)
		ebitenutil.DebugPrintAt(screen, "SpeciesID 0 -> REEFLET", x+34, 340)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP %02d/%02d", c.hp, definition.maxHP), x+70, 370)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("XP %02d/%02d", c.xp, goalXP), x+74, 397)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ACTIONS %d/%d   TIME %02d", g.actions, maxActions, max(0, 45-g.frames/60)), 153, 458)
	ebitenutil.DebugPrintAt(screen, g.message, 51, 479)
	button(screen, 20, 510, 210, "TRAIN [T / SPACE]", color.RGBA{220, 102, 75, 255})
	button(screen, 250, 510, 210, "REST [R / H]", color.RGBA{82, 171, 112, 255})
	ebitenutil.DebugPrintAt(screen, "Select card with tap or 1 / 2", 143, 635)
	if g.clear || g.over {
		title := "TWIN TRAINING COMPLETE!"
		if g.over {
			title = "TRAINING PLAN FAILED!"
		}
		vector.DrawFilledRect(screen, 40, 270, 400, 160, color.RGBA{5, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 145, 315)
		ebitenutil.DebugPrintAt(screen, g.message, 67, 350)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 390)
	}
}

func button(screen *ebiten.Image, x, y, w int, label string, fill color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), 100, fill, false)
	ebitenutil.DebugPrintAt(screen, label, x+w/2-len(label)*3, y+46)
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
	ebiten.SetWindowTitle("Species & Instance Training — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
