package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const (
	width, height = 480, 720
	cols, rows    = 10, 8
	cell          = 42
	ox, oy        = 30, 105
	ground        = 0
	wood          = 1
	stone         = 2
	crystal       = 3
	lantern       = 4
)

type island struct {
	name, hint string
	sky, soil  color.RGBA
	resources  [][3]int
	goal       [3]int
	enemies    int
}

var islands = []island{
	{"MOSS CAMP", "Gather wood + stone, craft a pickaxe, mine crystal,\nthen build the tide beacon.", color.RGBA{31, 76, 91, 255}, color.RGBA{55, 102, 76, 255}, [][3]int{{1, 1, wood}, {2, 3, wood}, {4, 5, stone}, {7, 2, stone}, {8, 6, crystal}}, [3]int{1, 1, 1}, 1},
	{"CRYSTAL CAVE", "Build the tide beacon while two armored crawlers\npatrol the crystal clearing.", color.RGBA{27, 35, 76, 255}, color.RGBA{55, 61, 95, 255}, [][3]int{{1, 6, wood}, {3, 2, wood}, {5, 5, stone}, {8, 1, stone}, {6, 3, crystal}, {9, 7, crystal}}, [3]int{1, 1, 2}, 2},
	{"EMBER ISLE", "Gather enough material and raise two tide beacons\nbefore the long night.", color.RGBA{92, 43, 46, 255}, color.RGBA{104, 68, 55, 255}, [][3]int{{0, 1, wood}, {2, 6, wood}, {3, 3, wood}, {5, 1, stone}, {7, 6, stone}, {9, 2, stone}, {4, 7, crystal}, {8, 4, crystal}}, [3]int{2, 2, 2}, 3},
}

type particle struct {
	x, y, vx, vy float64
	life         int
	c            color.RGBA
}
type crawler struct {
	x, y  int
	phase float64
}

type game struct {
	stage, px, py, frames, stageFrames           int
	tiles                                        [rows][cols]int
	bag                                          [3]int
	placed, goalLanterns                         int
	pickaxe                                      bool
	hp, score, best, combo                       int
	harvestPhase, targetX, targetY, shake, flash int
	crawlers                                     []crawler
	particles                                    []particle
	clear, over                                  bool
	message                                      string
	rng                                          *rand.Rand
	audio                                        *audio.Context
	gate                                         audiolab.Gate
	pulse                                        *shaderlab.Pulse
	cam                                          cameralab.State
	badge                                        *ebiten.Image
}

func newGame() *game {
	prepareCraftArt()
	g := &game{hp: 5, rng: rand.New(rand.NewSource(44)), best: sessionBest}
	g.audio = audiolab.Context()
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: width / 2, Y: height / 2}, ViewW: width, ViewH: height}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{255, 210, 74, 255})
	g.loadIsland(0)
	return g
}

var sessionBest int

func (g *game) loadIsland(n int) {
	g.stage, g.px, g.py, g.stageFrames = n, 0, rows/2, 0
	g.tiles, g.bag, g.placed, g.pickaxe = [rows][cols]int{}, [3]int{}, 0, false
	g.crawlers, g.particles = nil, nil
	data := islands[n]
	for _, r := range data.resources {
		g.tiles[r[1]][r[0]] = r[2]
	}
	for i := 0; i < data.enemies; i++ {
		g.crawlers = append(g.crawlers, crawler{cols - 1 - i, rows - 1 - i%2, float64(i)})
	}
	g.goalLanterns = 1
	if n == 2 {
		g.goalLanterns = 2
	}
	g.message = data.hint
}

func (g *game) move(dx, dy int) {
	if g.harvestPhase > 0 {
		return
	}
	nx, ny := g.px+dx, g.py+dy
	if nx >= 0 && nx < cols && ny >= 0 && ny < rows {
		g.px, g.py = nx, ny
	}
}

func (g *game) beginHarvest() {
	k := g.tiles[g.py][g.px]
	if k < wood || k > crystal {
		g.message = "Stand on wood, stone, or crystal first."
		return
	}
	if k == crystal && !g.pickaxe {
		g.message = "Crystal needs a pickaxe: craft with 1 wood + 1 stone."
		return
	}
	g.harvestPhase, g.targetX, g.targetY = 18, g.px, g.py
}

func (g *game) resolveHarvest() {
	k := g.tiles[g.targetY][g.targetX]
	if k < wood || k > crystal {
		return
	}
	g.tiles[g.targetY][g.targetX] = ground
	g.play(580)
	g.bag[k-1]++
	g.combo++
	g.score += 60 + g.combo*10
	g.shake, g.flash = 8, 5
	colors := []color.RGBA{{232, 154, 72, 255}, {173, 191, 205, 255}, {100, 229, 237, 255}}
	for i := 0; i < 14; i++ {
		a := float64(i) * 0.75
		g.particles = append(g.particles, particle{float64(ox + g.px*cell + cell/2), float64(oy + g.py*cell + cell/2), math.Cos(a) * 2.6, math.Sin(a)*2.6 - 1, 30, colors[k-1]})
	}
	g.message = fmt.Sprintf("Collected! Combo x%d — W%d S%d C%d", g.combo, g.bag[0], g.bag[1], g.bag[2])
}

func (g *game) craft() {
	if !g.pickaxe {
		if g.bag[0] < 1 || g.bag[1] < 1 {
			g.message = "Pickaxe needs 1 wood + 1 stone."
			return
		}
		g.bag[0]--
		g.bag[1]--
		g.pickaxe = true
		g.play(430)
		g.score += 100
		g.message = "Pickaxe crafted! Now crystal can break."
		return
	}
	if g.tiles[g.py][g.px] != ground {
		g.message = "Build on an empty tile."
		return
	}
	if g.bag[0] < 1 || g.bag[1] < 1 || g.bag[2] < 1 {
		g.message = "Lantern needs 1 wood + 1 stone + 1 crystal."
		return
	}
	g.bag[0]--
	g.bag[1]--
	g.bag[2]--
	g.tiles[g.py][g.px] = lantern
	g.play(750)
	g.placed++
	g.score += 300
	g.shake = 5
	for i := 0; i < 20; i++ {
		a := float64(i) * .6
		g.particles = append(g.particles, particle{float64(ox + g.px*cell + 21), float64(oy + g.py*cell + 21), math.Cos(a) * 2, math.Sin(a) * 2, 36, color.RGBA{255, 210, 74, 255}})
	}
	if g.placed >= g.goalLanterns {
		g.finishIsland()
	} else {
		g.message = "One light burns. Build the final lantern!"
	}
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .06)).Play()
}

func (g *game) finishIsland() {
	g.score += max(0, 1800-g.stageFrames/2) + g.hp*100
	if g.stage == len(islands)-1 {
		g.clear = true
		if g.score > sessionBest {
			sessionBest = g.score
		}
		g.best = sessionBest
		g.message = "Three islands shine! Try again for a higher expedition score."
		return
	}
	g.loadIsland(g.stage + 1)
	g.message = "New island unlocked — its recipe and danger are tougher."
}

func (g *game) updateCrawlers() {
	if g.stageFrames < 8*60 {
		return
	}
	for i := range g.crawlers {
		c := &g.crawlers[i]
		c.phase += .12
		lit := false
		for y := max(0, c.y-2); y <= min(rows-1, c.y+2); y++ {
			for x := max(0, c.x-2); x <= min(cols-1, c.x+2); x++ {
				if g.tiles[y][x] == lantern {
					lit = true
				}
			}
		}
		if lit {
			continue
		}
		if g.frames%(34-i*3) == 0 {
			if abs(c.x-g.px) > abs(c.y-g.py) {
				if c.x < g.px {
					c.x++
				} else if c.x > g.px {
					c.x--
				}
			} else {
				if c.y < g.py {
					c.y++
				} else if c.y > g.py {
					c.y--
				}
			}
		}
		if c.x == g.px && c.y == g.py && g.frames%55 == 0 {
			g.hp--
			g.combo = 0
			g.flash = 10
			g.shake = 10
			g.message = "Crawler hit! A lantern makes a safe zone."
		}
	}
	if g.hp <= 0 {
		g.over = true
		g.message = "The expedition ended. Retry and light safe zones sooner."
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
	g.stageFrames++
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
		g.beginHarvest()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.craft()
	}
	if x, y, ok := pressPosition(); ok {
		if x >= ox && x < ox+cols*cell && y >= oy && y < oy+rows*cell {
			tx, ty := (x-ox)/cell, (y-oy)/cell
			if abs(tx-g.px)+abs(ty-g.py) == 1 {
				g.move(tx-g.px, ty-g.py)
			}
		} else if y >= 575 {
			if x < 240 {
				g.beginHarvest()
			} else {
				g.craft()
			}
		}
	}
	if g.harvestPhase > 0 {
		g.harvestPhase--
		if g.harvestPhase == 8 {
			g.resolveHarvest()
		}
	}
	g.updateCrawlers()
	for i := len(g.particles) - 1; i >= 0; i-- {
		p := &g.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .08
		p.life--
		if p.life <= 0 {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
	}
	if g.shake > 0 {
		g.shake--
	}
	if g.flash > 0 {
		g.flash--
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	d := islands[g.stage]
	screen.Fill(d.sky)
	dx, dy := 0, 0
	if g.shake > 0 {
		// Deterministic screen shake: Draw never consumes gameplay randomness.
		dx = (g.frames%5 - 2)
		dy = ((g.frames/2)%5 - 2)
	}
	world := ebiten.NewImage(width, height)
	world.Fill(color.RGBA{0, 0, 0, 0})
	drawIslandArt(world, g.stage)
	vector.DrawFilledRect(world, 0, 0, width, 96, color.RGBA{3, 10, 23, 190}, false)
	g.drawTitle(world, d.name)
	g.drawEffectBadge(world)
	toolSprite := "workshop"
	if g.pickaxe {
		toolSprite = "pickaxe"
	}
	drawCraftSprite(world, toolSprite, 72, 62, 46)
	drawCraftSprite(world, "beacon", 432, 61, 52)
	ebitenutil.DebugPrintAt(world, fmt.Sprintf("HP %d  SCORE %05d  TIDE BEACON %d/%d", g.hp, g.score, g.placed, g.goalLanterns), 98, 43)
	ebitenutil.DebugPrintAt(world, fmt.Sprintf("BAG W%d S%d C%d  TOOL:%v", g.bag[0], g.bag[1], g.bag[2], map[bool]string{false: "HAND", true: "PICK"}[g.pickaxe]), 108, 68)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			px, py := float32(ox+x*cell), float32(oy+y*cell)
			vector.StrokeRect(world, px+1, py+1, cell-2, cell-2, 1, color.RGBA{220, 244, 238, 42}, false)
			if k := g.tiles[y][x]; k != ground {
				if k == lantern {
					vector.DrawFilledCircle(world, px+cell/2, py+cell/2, 32+float32(math.Sin(float64(g.frames)*.09))*3, color.RGBA{255, 211, 91, 34}, true)
				}
				drawCraftSprite(world, craftSprite(k), float64(px+cell/2), float64(py+cell/2), cell+6)
			}
		}
	}
	vector.StrokeRect(world, ox-5, oy-5, cols*cell+10, rows*cell+10, 3, color.RGBA{207, 238, 199, 115}, false)
	for _, c := range g.crawlers {
		bounce := math.Sin(c.phase) * 3
		drawCraftSprite(world, "crawler", float64(ox+c.x*cell+21), float64(oy+c.y*cell+21)+bounce, 48)
	}
	py := float64(oy + g.py*cell + 21)
	if g.harvestPhase > 0 {
		py -= math.Sin(float64(g.harvestPhase)/18*math.Pi) * 8
	}
	heroSprite := "hero-idle"
	if g.harvestPhase > 0 {
		heroSprite = "hero-mine"
	}
	drawCraftSprite(world, heroSprite, float64(ox+g.px*cell+21), py, 58)
	for _, p := range g.particles {
		vector.DrawFilledCircle(world, float32(p.x), float32(p.y), float32(2+p.life%3), p.c, false)
	}
	ebitenutil.DebugPrintAt(world, g.message, 38, 470)
	button(world, 10, "HARVEST [X]", color.RGBA{170, 101, 54, 255})
	button(world, 250, "CRAFT / BUILD [C]", color.RGBA{209, 148, 49, 255})
	ebitenutil.DebugPrintAt(world, "Tap a neighboring tile to move / arrows or WASD", 65, 655)
	ebitenutil.DebugPrintAt(world, "Finish all 3 islands. Faster + more HP = more score.", 55, 680)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dx), float64(dy))
	screen.DrawImage(world, op)
	if g.flash > 0 {
		vector.DrawFilledRect(screen, 0, 0, width, height, color.RGBA{255, 255, 255, 55}, false)
	}
	if g.clear || g.over {
		vector.DrawFilledRect(screen, 28, 235, 424, 220, color.RGBA{5, 13, 28, 245}, false)
		title := "EXPEDITION COMPLETE!"
		if g.over {
			title = "EXPEDITION LOST"
		}
		ebitenutil.DebugPrintAt(screen, title, 151, 275)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE %05d   BEST %05d", g.score, max(g.best, sessionBest)), 145, 316)
		ebitenutil.DebugPrintAt(screen, g.message, 47, 353)
		ebitenutil.DebugPrintAt(screen, "TAP / ENTER TO EXPLORE AGAIN", 125, 409)
	}
}

func (g *game) drawTitle(screen *ebiten.Image, island string) {
	label := fmt.Sprintf("EBI CRAFT  ISLAND %d/3  %s", g.stage+1, island)
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(105, 4)
		text.Draw(screen, label, face, op)
		return
	}
	ebitenutil.DebugPrintAt(screen, label, 105, 16)
}
func (g *game) drawEffectBadge(screen *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.frames)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(width-32, 10)
	screen.DrawImage(fx, op)
}

func craftSprite(k int) string {
	switch k {
	case wood:
		return "wood"
	case stone:
		return "stone"
	case crystal:
		return "crystal"
	default:
		return "beacon"
	}
}
func button(s *ebiten.Image, x int, label string, c color.RGBA) {
	vector.DrawFilledRect(s, float32(x), 575, 220, 58, c, false)
	ebitenutil.DebugPrintAt(s, label, x+47, 598)
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
	ebiten.SetWindowTitle("Ebi Craft Expedition")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
