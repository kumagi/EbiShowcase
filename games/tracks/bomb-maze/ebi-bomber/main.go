package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
)

const (
	screenW   = 480
	screenH   = 720
	cols      = 9
	rows      = 9
	cell      = 48
	boardX    = 24
	boardY    = 104
	fuse      = 150
	blastLife = 28
)

type point struct{ x, y int }
type bomb struct {
	at    point
	timer int
}
type flame struct {
	at    point
	timer int
}
type item struct {
	at   point
	kind int
}
type game struct {
	player, enemy, exit                               point
	soft                                              map[point]bool
	bombs                                             []bomb
	flames                                            []flame
	items                                             []item
	power, capacity, speed, broken, frames, enemyTick int
	enemyAlive, exitOpen, won, lost                   bool
	message                                           string
}

func newGame() *game {
	g := &game{player: point{1, 1}, enemy: point{7, 7}, exit: point{7, 1}, soft: map[point]bool{}, power: 2, capacity: 1, speed: 1, enemyAlive: true, message: "Break walls, collect upgrades, defeat the scout, then reach EXIT."}
	for _, p := range []point{{3, 1}, {5, 1}, {1, 3}, {3, 3}, {5, 3}, {7, 3}, {1, 5}, {3, 5}, {5, 5}, {7, 5}, {3, 7}, {5, 7}} {
		g.soft[p] = true
	}
	return g
}
func hard(p point) bool {
	return p.x < 0 || p.x >= cols || p.y < 0 || p.y >= rows || p.x == 0 || p.y == 0 || p.x == cols-1 || p.y == rows-1 || (p.x%2 == 0 && p.y%2 == 0)
}
func (g *game) blocked(p point) bool { return hard(p) || g.soft[p] || g.bombAt(p) }
func (g *game) Update() error {
	if g.won || g.lost {
		if retry() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if d, ok := inputDir(); ok {
		for n := 0; n < g.speed; n++ {
			p := point{g.player.x + d.x, g.player.y + d.y}
			if g.blocked(p) {
				break
			}
			g.player = p
			g.collect()
		}
	}
	if placePressed() && len(g.bombs) < g.capacity && !g.bombAt(g.player) {
		g.bombs = append(g.bombs, bomb{g.player, fuse})
		g.message = "Bomb armed. Move to a safe corridor."
	}
	g.updateBombs()
	g.updateFlames()
	g.moveEnemy()
	if g.enemyAlive && g.player == g.enemy {
		g.lost = true
		g.message = "The scout caught Ebi Tenjiroh."
	}
	if g.exitOpen && g.player == g.exit {
		g.won = true
		g.message = "EXIT reached with every system working together!"
	}
	if g.frames >= 90*60 {
		g.lost = true
		g.message = "Time up."
	}
	return nil
}
func (g *game) updateBombs() {
	next := g.bombs[:0]
	for _, b := range g.bombs {
		b.timer--
		if b.timer <= 0 {
			for _, p := range g.blast(b.at) {
				g.flames = append(g.flames, flame{p, blastLife})
				if g.soft[p] {
					delete(g.soft, p)
					g.broken++
					if (p.x*7+p.y*11)%3 == 0 {
						g.items = append(g.items, item{p, (p.x + p.y) % 3})
					}
				}
			}
		} else {
			next = append(next, b)
		}
	}
	g.bombs = next
}
func (g *game) updateFlames() {
	next := g.flames[:0]
	for _, f := range g.flames {
		f.timer--
		if f.at == g.player {
			g.lost = true
			g.message = "Ebi Tenjiroh was caught in the blast."
		}
		if g.enemyAlive && f.at == g.enemy {
			g.enemyAlive = false
			g.message = "Scout defeated. Find upgrades and open EXIT."
		}
		if f.timer > 0 {
			next = append(next, f)
		}
	}
	g.flames = next
	if !g.enemyAlive && g.broken >= 8 {
		g.exitOpen = true
	}
}
func (g *game) blast(o point) []point {
	r := []point{o}
	for _, d := range []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		for n := 1; n <= g.power; n++ {
			p := point{o.x + d.x*n, o.y + d.y*n}
			if hard(p) {
				break
			}
			r = append(r, p)
			if g.soft[p] {
				break
			}
		}
	}
	return r
}
func (g *game) moveEnemy() {
	if !g.enemyAlive {
		return
	}
	g.enemyTick++
	if g.enemyTick < 18 {
		return
	}
	g.enemyTick = 0
	danger := map[point]bool{}
	for _, b := range g.bombs {
		if b.timer < 80 {
			for _, p := range g.blast(b.at) {
				danger[p] = true
			}
		}
	}
	best := g.enemy
	bestScore := 999
	for _, d := range []point{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		p := point{g.enemy.x + d.x, g.enemy.y + d.y}
		if g.blocked(p) || danger[p] {
			continue
		}
		score := abs(p.x-g.player.x) + abs(p.y-g.player.y)
		if score < bestScore {
			bestScore = score
			best = p
		}
	}
	g.enemy = best
}
func (g *game) collect() {
	out := g.items[:0]
	for _, it := range g.items {
		if it.at == g.player {
			switch it.kind {
			case 0:
				g.power = min(4, g.power+1)
			case 1:
				g.capacity = min(3, g.capacity+1)
			case 2:
				g.speed = 2
			}
			g.message = "Upgrade collected: blast / capacity / movement improved."
		} else {
			out = append(out, it)
		}
	}
	g.items = out
}
func (g *game) bombAt(p point) bool {
	for _, b := range g.bombs {
		if b.at == p {
			return true
		}
	}
	return false
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{7, 15, 29, 255})
	ebitenutil.DebugPrintAt(s, "EBI BOMBER", 203, 17)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("WALLS %d/8  POWER %d  BOMBS %d  SPEED %d  TIME %02d", g.broken, g.power, g.capacity, g.speed, max(0, 90-g.frames/60)), 60, 43)
	ebitenutil.DebugPrintAt(s, g.message, 28, 72)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			p := point{x, y}
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			c := color.RGBA{25, 47, 51, 255}
			vector.DrawFilledRect(s, px+1, py+1, cell-2, cell-2, c, false)
			if hard(p) {
				trackatlas.Draw(s, "tile-wall", float64(px+1), float64(py+1), float64(cell-2))
			} else if g.soft[p] {
				trackatlas.Draw(s, "tile-crate", float64(px+1), float64(py+1), float64(cell-2))
			}
		}
	}
	if g.exitOpen {
		trackatlas.DrawTinted(s, "tile-exit", tileCX(g.exit), tileCY(g.exit), cell-2, 1, 1, 1, 1)
		ebitenutil.DebugPrintAt(s, "EXIT", int(tileCX(g.exit))-12, int(tileCY(g.exit))-5)
	} else {
		trackatlas.DrawTinted(s, "tile-exit", tileCX(g.exit), tileCY(g.exit), cell-2, 0.5, 0.5, 0.5, 1)
		ebitenutil.DebugPrintAt(s, "LOCK", int(tileCX(g.exit))-12, int(tileCY(g.exit))-5)
	}
	upgradeSprites := []string{"upgrade-blast", "upgrade-cap", "upgrade-spd"}
	for _, it := range g.items {
		trackatlas.DrawCentered(s, upgradeSprites[it.kind], tileCX(it.at), tileCY(it.at), 30)
	}
	for _, b := range g.bombs {
		trackatlas.DrawCentered(s, "bomb", tileCX(b.at), tileCY(b.at), 30)
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%.1f", float64(b.timer)/60), int(tileCX(b.at))-12, int(tileCY(b.at))-5)
	}
	for _, f := range g.flames {
		trackatlas.DrawCentered(s, "flame", tileCX(f.at), tileCY(f.at), 34)
	}
	trackatlas.DrawCentered(s, "hero", tileCX(g.player), tileCY(g.player), 34)
	if g.enemyAlive {
		trackatlas.DrawCentered(s, "scout", tileCX(g.enemy), tileCY(g.enemy), 34)
	}
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT", "BOMB"}
	for i, l := range labels {
		vector.DrawFilledRect(s, float32(i*96+3), 600, 90, 62, color.RGBA{46, 78, 114, 255}, false)
		ebitenutil.DebugPrintAt(s, l, i*96+27, 627)
	}
	ebitenutil.DebugPrintAt(s, "Arrows/WASD + Space | tap controls", 107, 688)
	if g.won {
		overlay(s, "STAGE CLEAR!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(s, "MISSION FAILED\n\nTAP / ENTER TO RETRY")
	}
}
func tileCX(p point) float64 { return float64(boardX + p.x*cell + cell/2) }
func tileCY(p point) float64 { return float64(boardY + p.y*cell + cell/2) }
func inputDir() (point, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		return point{-1, 0}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return point{0, -1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return point{0, 1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		return point{1, 0}, true
	}
	if x, y, ok := press(); ok && y >= 600 && x < 384 {
		return [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}[min(3, x/96)], true
	}
	return point{}, false
}
func placePressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		return true
	}
	x, y, ok := press()
	return ok && y >= 600 && x >= 384
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
func retry() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func overlay(s *ebiten.Image, t string) {
	vector.DrawFilledRect(s, 40, 270, 400, 160, color.RGBA{4, 12, 27, 245}, false)
	ebitenutil.DebugPrintAt(s, t, 125, 330)
}
func (g *game) Layout(int, int) (int, int) { return screenW, screenH }
func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Ebi Bomber")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
