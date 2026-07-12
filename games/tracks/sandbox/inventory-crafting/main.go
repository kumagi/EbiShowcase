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

type recipe struct {
	name               string
	wood, stone, fiber int
	rope               int
}

var recipes = []recipe{
	{name: "ROPE", fiber: 3},
	{name: "PICKAXE", wood: 2, stone: 2, rope: 1},
	{name: "BEACON", wood: 2, stone: 3, rope: 1},
}

type game struct {
	wood, stone, fiber, rope int
	woodNodes, stoneNodes    int
	fiberNodes               int
	pickaxe, beacon          bool
	actions                  int
	frames                   int
	clear, over              bool
	message                  string
}

func newGame() *game {
	return &game{woodNodes: 2, stoneNodes: 3, fiberNodes: 2, message: "Gather bundles, then craft two ropes and two tools."}
}

func (g *game) gather(kind int) {
	switch kind {
	case 0:
		if g.woodNodes == 0 {
			g.message = "The driftwood patch is empty."
			return
		}
		g.woodNodes--
		g.wood += 2
		g.message = "Gathered a stack of 2 wood."
	case 1:
		if g.stoneNodes == 0 {
			g.message = "The stone patch is empty."
			return
		}
		g.stoneNodes--
		g.stone += 2
		g.message = "Gathered a stack of 2 stone."
	case 2:
		if g.fiberNodes == 0 {
			g.message = "The sea-grass patch is empty."
			return
		}
		g.fiberNodes--
		g.fiber += 3
		g.message = "Gathered a stack of 3 fiber."
	}
	g.actions++
	g.checkEnd()
}

func (g *game) craft(index int) {
	r := recipes[index]
	if g.wood < r.wood || g.stone < r.stone || g.fiber < r.fiber || g.rope < r.rope {
		g.message = "Recipe does not match the current inventory."
		return
	}
	g.wood -= r.wood
	g.stone -= r.stone
	g.fiber -= r.fiber
	g.rope -= r.rope
	switch index {
	case 0:
		g.rope++
	case 1:
		if g.pickaxe {
			g.message = "You already made the pickaxe."
			// Restore ingredients because duplicate outputs are not useful here.
			g.wood += r.wood
			g.stone += r.stone
			g.rope += r.rope
			return
		}
		g.pickaxe = true
	case 2:
		if g.beacon {
			g.message = "You already made the beacon."
			g.wood += r.wood
			g.stone += r.stone
			g.rope += r.rope
			return
		}
		g.beacon = true
	}
	g.actions++
	g.message = r.name + " crafted: ingredients consumed, output added."
	g.checkEnd()
}

func (g *game) checkEnd() {
	if g.pickaxe && g.beacon {
		g.clear = true
		g.message = "Workshop complete: pickaxe and beacon crafted!"
	} else if g.actions >= 13 {
		g.over = true
		g.message = "Action limit reached. Follow the recipe counts exactly!"
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
	if g.frames >= 50*60 {
		g.over = true
		g.message = "Workshop closed. Read the recipes and try again!"
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.gather(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.gather(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.gather(2)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.craft(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.craft(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.craft(2)
	}
	if x, y, ok := pressPosition(); ok {
		if y >= 110 && y < 220 {
			g.gather(min(2, x/160))
		} else if y >= 420 && y < 545 {
			g.craft(min(2, x/160))
		}
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 24, 38, 255})
	ebitenutil.DebugPrintAt(screen, "TIDEPOOL CRAFT TABLE", 166, 22)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ACTIONS %02d/13   TIME %02d", g.actions, max(0, 50-g.frames/60)), 142, 50)
	drawGather(screen, 10, "WOOD [1]", "+2", g.woodNodes, color.RGBA{194, 126, 66, 255})
	drawGather(screen, 170, "STONE [2]", "+2", g.stoneNodes, color.RGBA{127, 149, 166, 255})
	drawGather(screen, 330, "FIBER [3]", "+3", g.fiberNodes, color.RGBA{79, 177, 121, 255})
	ebitenutil.DebugPrintAt(screen, "INVENTORY STACKS", 180, 250)
	items := []struct {
		name string
		qty  int
	}{{"WOOD", g.wood}, {"STONE", g.stone}, {"FIBER", g.fiber}, {"ROPE", g.rope}}
	for i, item := range items {
		x := 30 + i*110
		vector.DrawFilledRect(screen, float32(x), 280, 90, 78, color.RGBA{35, 58, 78, 255}, false)
		ebitenutil.DebugPrintAt(screen, item.name, x+20, 300)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("x %d", item.qty), x+31, 329)
	}
	ebitenutil.DebugPrintAt(screen, "RECIPES — TAP TO CRAFT", 154, 390)
	drawRecipe(screen, 10, "ROPE [R]", "3 fiber", false)
	drawRecipe(screen, 170, "PICK [P]", "2W 2S 1R", g.pickaxe)
	drawRecipe(screen, 330, "BEACON [B]", "2W 3S 1R", g.beacon)
	ebitenutil.DebugPrintAt(screen, g.message, 54, 580)
	ebitenutil.DebugPrintAt(screen, "Gather -> inventory -> match -> consume -> output", 70, 620)
	if g.clear || g.over {
		title := "WORKSHOP COMPLETE!"
		if g.over {
			title = "RECIPE PLAN FAILED!"
		}
		vector.DrawFilledRect(screen, 38, 252, 404, 166, color.RGBA{5, 14, 27, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 150, 298)
		ebitenutil.DebugPrintAt(screen, g.message, 70, 334)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 381)
	}
}

func drawGather(screen *ebiten.Image, x int, name, gain string, left int, c color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), 110, 140, 110, c, false)
	ebitenutil.DebugPrintAt(screen, name, x+30, 137)
	ebitenutil.DebugPrintAt(screen, gain+"  PATCHES "+fmt.Sprint(left), x+24, 180)
}

func drawRecipe(screen *ebiten.Image, x int, name, cost string, done bool) {
	c := color.RGBA{69, 91, 116, 255}
	if done {
		c = color.RGBA{70, 158, 113, 255}
	}
	vector.DrawFilledRect(screen, float32(x), 420, 140, 125, c, false)
	ebitenutil.DebugPrintAt(screen, name, x+32, 447)
	ebitenutil.DebugPrintAt(screen, cost, x+30, 485)
	if done {
		ebitenutil.DebugPrintAt(screen, "DONE", x+50, 518)
	}
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
	ebiten.SetWindowTitle("Tidepool Craft Table — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
