package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const (
	width, height = 480, 720
	worldW        = 12
	worldH        = 8
	chunkW        = 4
	cell          = 38
	ox, oy        = 12, 92
	ground        = 0
	woodTile      = 1
	stoneTile     = 2
	crystalTile   = 3
	lanternTile   = 4
)

type snapshot struct {
	tiles                [worldH][worldW]int
	px, py               int
	wood, stone, crystal int
	lanterns             [3]bool
	pickaxe              bool
}

type game struct {
	tiles                   [worldH][worldW]int
	px, py                  int
	wood, stone, crystal    int
	lanterns                [3]bool
	pickaxe                 bool
	creatureX, creatureY    int
	creatureAlive           bool
	hp, frames, attackTimer int
	saved                   snapshot
	hasSave                 bool
	clear, over             bool
	message                 string
}

func newGame() *game {
	g := &game{px: 0, py: 4, hp: 5, creatureX: 11, creatureY: 7, message: "Gather one set in each chunk; craft and place 3 lanterns."}
	for chunk := 0; chunk < 3; chunk++ {
		base := chunk * chunkW
		g.tiles[1+chunk%2][base+1] = woodTile
		g.tiles[4][base+2] = stoneTile
		g.tiles[6-chunk%2][base+3] = crystalTile
	}
	// Chunk 1 also contains the extra wood and stone used to craft the pickaxe.
	g.tiles[2][0], g.tiles[3][0] = woodTile, stoneTile
	return g
}

func (g *game) night() bool { return (g.frames/900)%2 == 1 }

func (g *game) move(dx, dy int) {
	nx, ny := g.px+dx, g.py+dy
	if nx < 0 || nx >= worldW || ny < 0 || ny >= worldH {
		return
	}
	g.px, g.py = nx, ny
}

func (g *game) harvest() {
	kind := g.tiles[g.py][g.px]
	if kind == crystalTile && !g.pickaxe {
		g.message = "Crystal needs a pickaxe. Gather wood and stone, then craft."
		return
	}
	switch kind {
	case woodTile:
		g.wood++
	case stoneTile:
		g.stone++
	case crystalTile:
		g.crystal++
	default:
		g.message = "No resource under Ebi Tenjiroh. Move onto a resource tile."
		return
	}
	g.tiles[g.py][g.px] = ground
	g.message = "Harvested into inventory; this tile is now ground."
}

func (g *game) craftOrPlace() {
	if !g.pickaxe {
		if g.wood < 1 || g.stone < 1 {
			g.message = "Pickaxe recipe needs 1 wood + 1 stone."
			return
		}
		g.wood--
		g.stone--
		g.pickaxe = true
		g.message = "Pickaxe crafted! Crystal tiles can now be harvested."
		return
	}
	g.placeLantern()
}

func (g *game) placeLantern() {
	chunk := g.px / chunkW
	if g.lanterns[chunk] {
		g.message = "This chunk already has a lantern."
		return
	}
	if g.tiles[g.py][g.px] != ground {
		g.message = "Harvest this tile before building here."
		return
	}
	if g.wood < 1 || g.stone < 1 || g.crystal < 1 {
		g.message = "Recipe needs 1 wood + 1 stone + 1 crystal."
		return
	}
	g.wood--
	g.stone--
	g.crystal--
	g.tiles[g.py][g.px] = lanternTile
	g.lanterns[chunk] = true
	g.message = fmt.Sprintf("Lantern crafted and placed in chunk %d!", chunk+1)
	if g.lanterns[0] && g.lanterns[1] && g.lanterns[2] {
		g.clear = true
		g.message = "All chunks are linked by light — Ebi Craft complete!"
	}
}

func (g *game) saveOrLoad() {
	if !g.hasSave {
		g.saved = snapshot{g.tiles, g.px, g.py, g.wood, g.stone, g.crystal, g.lanterns, g.pickaxe}
		g.hasSave = true
		g.message = "Checkpoint saved: tile diffs, position, and inventory copied."
		return
	}
	g.tiles, g.px, g.py = g.saved.tiles, g.saved.px, g.saved.py
	g.wood, g.stone, g.crystal = g.saved.wood, g.saved.stone, g.saved.crystal
	g.lanterns = g.saved.lanterns
	g.pickaxe = g.saved.pickaxe
	g.message = "Checkpoint restored from the saved world snapshot."
}

func (g *game) updateCreature() {
	if !g.night() {
		g.creatureAlive = false
		return
	}
	if !g.creatureAlive {
		g.creatureX, g.creatureY, g.creatureAlive = 11, 7, true
	}
	for y := max(0, g.creatureY-2); y <= min(worldH-1, g.creatureY+2); y++ {
		for x := max(0, g.creatureX-2); x <= min(worldW-1, g.creatureX+2); x++ {
			if g.tiles[y][x] == lanternTile {
				g.creatureAlive = false
				g.message = "Lantern light sent the night crawler away."
				return
			}
		}
	}
	if g.frames%36 == 0 {
		if g.creatureX < g.px {
			g.creatureX++
		} else if g.creatureX > g.px {
			g.creatureX--
		} else if g.creatureY < g.py {
			g.creatureY++
		} else if g.creatureY > g.py {
			g.creatureY--
		}
	}
	distance := abs(g.creatureX-g.px) + abs(g.creatureY-g.py)
	if distance <= 1 {
		g.attackTimer++
		if g.attackTimer >= 50 {
			g.hp--
			g.attackTimer = 0
			g.message = "Night crawler hit! Reach lantern light."
		}
	} else {
		g.attackTimer = 0
	}
	if g.hp <= 0 {
		g.over = true
		g.message = "The night crawler won. Build light sooner!"
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
	if g.frames >= 90*60 {
		g.over = true
		g.message = "The long night arrived before the light network was ready."
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.move(-1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.move(1, 0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.move(0, -1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.move(0, 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.harvest()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.craftOrPlace()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		g.saveOrLoad()
	}
	if x, y, ok := pressPosition(); ok {
		if y >= oy && y < oy+worldH*cell && x >= ox && x < ox+worldW*cell {
			tx, ty := (x-ox)/cell, (y-oy)/cell
			dx, dy := tx-g.px, ty-g.py
			if abs(dx)+abs(dy) == 1 {
				g.move(dx, dy)
			}
		} else if y >= 545 {
			switch {
			case x < 160:
				g.harvest()
			case x < 320:
				g.craftOrPlace()
			default:
				g.saveOrLoad()
			}
		}
	}
	g.updateCreature()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	bg := color.RGBA{21, 52, 68, 255}
	if g.night() {
		bg = color.RGBA{8, 18, 39, 255}
	}
	screen.Fill(bg)
	phase := "DAY"
	if g.night() {
		phase = "NIGHT"
	}
	ebitenutil.DebugPrintAt(screen, "EBI CRAFT — THREE CHUNKS", 142, 18)
	tool := "HAND"
	if g.pickaxe {
		tool = "PICK"
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s HP%d W%d S%d C%d TOOL:%s", phase, g.hp, g.wood, g.stone, g.crystal, tool), 105, 47)
	for y := 0; y < worldH; y++ {
		for x := 0; x < worldW; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			base := color.RGBA{45 + uint8((x/chunkW)*8), 86, 75, 255}
			vector.DrawFilledRect(screen, px+1, py+1, cell-2, cell-2, base, false)
			kind := g.tiles[y][x]
			if kind != ground {
				trackatlas.Draw(screen, tileSprite(kind), float64(px+3), float64(py+3), float64(cell-6))
			}
		}
	}
	for c := 1; c < 3; c++ {
		x := float32(ox + c*chunkW*cell)
		vector.StrokeLine(screen, x, oy, x, oy+worldH*cell, 4, color.RGBA{245, 195, 75, 180}, false)
	}
	if g.creatureAlive {
		trackatlas.DrawCentered(screen, "slug", float64(ox+g.creatureX*cell+cell/2), float64(oy+g.creatureY*cell+cell/2), 28)
	}
	trackatlas.DrawCentered(screen, "hero", float64(ox+g.px*cell+cell/2), float64(oy+g.py*cell+cell/2), 30)
	for c := 0; c < 3; c++ {
		status := "DARK"
		if g.lanterns[c] {
			status = "LIT"
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CHUNK %d %s", c+1, status), 30+c*155, 420)
	}
	ebitenutil.DebugPrintAt(screen, g.message, 50, 482)
	button(screen, 10, "HARVEST [X]", color.RGBA{177, 114, 59, 255})
	button(screen, 170, "CRAFT+PLACE [C]", color.RGBA{211, 151, 59, 255})
	label := "SAVE [K]"
	if g.hasSave {
		label = "LOAD [K]"
	}
	button(screen, 330, label, color.RGBA{62, 125, 151, 255})
	ebitenutil.DebugPrintAt(screen, "Tap a neighboring cell to move / arrows or WASD", 68, 632)
	ebitenutil.DebugPrintAt(screen, "Pick: 1W+1S / Lantern: 1W+1S+1C", 92, 665)
	if g.clear || g.over {
		title := "EBI CRAFT COMPLETE!"
		if g.over {
			title = "WORLD LOST TO NIGHT!"
		}
		vector.DrawFilledRect(screen, 36, 250, 408, 168, color.RGBA{5, 13, 28, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 146, 297)
		ebitenutil.DebugPrintAt(screen, g.message, 68, 334)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO RETRY", 146, 382)
	}
}

func tileSprite(kind int) string {
	switch kind {
	case woodTile:
		return "tile-wood"
	case stoneTile:
		return "tile-stone"
	case crystalTile:
		return "tile-glass"
	default:
		return "tile-lantern"
	}
}

func button(screen *ebiten.Image, x int, label string, c color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), 545, 140, 58, c, false)
	ebitenutil.DebugPrintAt(screen, label, x+15, 568)
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

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Craft — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
