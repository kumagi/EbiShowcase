package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width, height = 480, 720
	cols, rows    = 8, 8
	cell          = 50
	ox, oy        = 40, 132
	materialWood  = 0
	materialStone = 1
	toolAxe       = 0
	toolPick      = 1
	toolTorch     = 2
	goalEach      = 4
	timeLimit     = 45 * 60
)

type point struct{ x, y int }

type block struct {
	material int
	hp       int
	exists   bool
}

type tool struct {
	name      string
	preferred int
	power     int
	color     color.RGBA
}

var tools = []tool{
	{"AXE", materialWood, 3, color.RGBA{207, 126, 64, 255}},
	{"PICK", materialStone, 3, color.RGBA{105, 143, 171, 255}},
	{"TORCH", -1, 0, color.RGBA{242, 177, 62, 255}},
}

type game struct {
	blocks             [rows][cols]block
	torches            []point
	selected           int
	durability         [2]int
	wood, stone        int
	frames, worldClock int
	clear, over        bool
	message            string
}

func newGame() *game {
	g := &game{selected: toolAxe, durability: [2]int{6, 6}, message: "Match AXE to wood and PICK to stone; light the night."}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			material := (x + y) % 2
			g.blocks[y][x] = block{material: material, hp: 3, exists: true}
		}
	}
	return g
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	g.worldClock++
	if g.frames >= timeLimit {
		g.over = true
		g.message = "The expedition ended before both material goals were met."
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.selected = toolAxe
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.selected = toolPick
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) || inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.selected = toolTorch
	}
	if x, y, ok := pressPosition(); ok {
		if y >= 592 {
			g.selected = min(2, x/160)
			return nil
		}
		gx, gy := (x-ox)/cell, (y-oy)/cell
		if x >= ox && gx >= 0 && gx < cols && y >= oy && gy >= 0 && gy < rows {
			g.useAt(gx, gy)
		}
	}
	return nil
}

func (g *game) useAt(x, y int) {
	if g.selected == toolTorch {
		if len(g.torches) >= 2 {
			g.message = "Both torches are already placed."
			return
		}
		for _, p := range g.torches {
			if p.x == x && p.y == y {
				g.message = "A torch already lights that tile."
				return
			}
		}
		g.torches = append(g.torches, point{x, y})
		g.message = "Torch placed: nearby tile brightness was recalculated."
		g.checkGoal()
		return
	}
	b := &g.blocks[y][x]
	if !b.exists {
		g.message = "That tile is already mined."
		return
	}
	if g.lightAt(x, y) < 4 {
		g.message = "Too dark to mine safely. Place a torch nearby."
		return
	}
	if g.durability[g.selected] <= 0 {
		g.message = tools[g.selected].name + " is broken."
		return
	}
	data := tools[g.selected]
	power, cost := 1, 2
	if data.preferred == b.material {
		power, cost = data.power, 1
	}
	b.hp -= power
	g.durability[g.selected] = max(0, g.durability[g.selected]-cost)
	if b.hp <= 0 {
		b.exists = false
		if b.material == materialWood {
			g.wood++
		} else {
			g.stone++
		}
		g.message = fmt.Sprintf("Mined! Inventory W%d/%d S%d/%d.", g.wood, goalEach, g.stone, goalEach)
	} else {
		g.message = fmt.Sprintf("Wrong tool: only 1 damage, durability cost %d.", cost)
	}
	g.checkGoal()
	if g.clear {
		return
	}
	if (g.wood < goalEach && g.durability[toolAxe] == 0) || (g.stone < goalEach && g.durability[toolPick] == 0) {
		g.over = true
		g.message = "A needed tool broke. Match each material and retry."
	}
}

func (g *game) checkGoal() {
	if g.wood >= goalEach && g.stone >= goalEach && len(g.torches) == 2 {
		g.clear = true
		g.message = "Tool match and tile lighting mastered!"
	}
}

func (g *game) daylight() int {
	phase := float64(g.worldClock%(24*60)) / float64(24*60) * 2 * math.Pi
	return int(math.Round(5 + 4*math.Sin(phase)))
}

func (g *game) lightAt(x, y int) int {
	light := g.daylight()
	for _, torch := range g.torches {
		distance := abs(x-torch.x) + abs(y-torch.y)
		light = max(light, 9-distance*2)
	}
	return max(0, min(9, light))
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{9, 18, 32, 255})
	ebitenutil.DebugPrintAt(screen, "TOOLS & TILE LIGHT", 180, 18)
	daylight := g.daylight()
	period := "DAY"
	if daylight < 4 {
		period = "NIGHT"
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s LIGHT %d/9   TIME %02d   W %d/%d  S %d/%d", period, daylight, max(0, (timeLimit-g.frames+59)/60), g.wood, goalEach, g.stone, goalEach), 79, 45)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("AXE %d/6   PICK %d/6   TORCHES %d/2", g.durability[0], g.durability[1], len(g.torches)), 126, 70)
	ebitenutil.DebugPrintAt(screen, g.message, 55, 98)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			b := g.blocks[y][x]
			if b.exists {
				base := color.RGBA{167, 102, 57, 255}
				label := "W"
				if b.material == materialStone {
					base = color.RGBA{111, 128, 145, 255}
					label = "S"
				}
				vector.DrawFilledRect(screen, px+2, py+2, cell-4, cell-4, base, false)
				vector.StrokeRect(screen, px+4, py+4, cell-8, cell-8, float32(max(1, b.hp)), color.RGBA{235, 242, 244, 140}, false)
				ebitenutil.DebugPrintAt(screen, label+fmt.Sprint(b.hp), int(px)+18, int(py)+20)
			} else {
				vector.StrokeRect(screen, px+2, py+2, cell-4, cell-4, 1, color.RGBA{48, 67, 83, 255}, false)
			}
			light := g.lightAt(x, y)
			if light < 9 {
				alpha := uint8((9 - light) * 22)
				vector.DrawFilledRect(screen, px+1, py+1, cell-2, cell-2, color.RGBA{5, 10, 22, alpha}, false)
			}
		}
	}
	for _, p := range g.torches {
		px, py := float32(ox+p.x*cell+cell/2), float32(oy+p.y*cell+cell/2)
		vector.DrawFilledCircle(screen, px, py, 8, color.RGBA{255, 194, 63, 255}, false)
		vector.StrokeCircle(screen, px, py, 14, 3, color.RGBA{255, 221, 123, 210}, false)
	}
	for i, data := range tools {
		x := i * 160
		fill := color.RGBA{48, 72, 97, 255}
		if i == g.selected {
			fill = data.color
		}
		vector.DrawFilledRect(screen, float32(x+8), 592, 144, 86, fill, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d %s", i+1, data.name), x+54, 617)
		caption := "LIGHT SOURCE"
		if i < 2 {
			caption = fmt.Sprintf("DUR %d", g.durability[i])
		}
		ebitenutil.DebugPrintAt(screen, caption, x+45, 646)
	}
	if g.clear || g.over {
		title := "MINING CAMP COMPLETE!"
		if g.over {
			title = "EXPEDITION FAILED!"
		}
		vector.DrawFilledRect(screen, 40, 270, 400, 160, color.RGBA{5, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 151, 315)
		ebitenutil.DebugPrintAt(screen, g.message, 76, 350)
		ebitenutil.DebugPrintAt(screen, "TAP / SPACE TO RETRY", 148, 390)
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
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Tools & Tile Light — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
