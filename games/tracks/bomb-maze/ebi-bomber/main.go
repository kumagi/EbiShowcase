package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"github.com/kumagi/EbiShowcase/internal/uilab"
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
type spark struct{ x, y, vx, vy, life float64 }
type game struct {
	player, enemy, exit                                  point
	soft                                                 map[point]bool
	bombs                                                []bomb
	flames                                               []flame
	items                                                []item
	power, capacity, speed, broken, frames, enemyTick    int
	enemyAlive, exitOpen, won, lost                      bool
	message                                              string
	stage, totalFrames, bestFrames, requiredBreak, shake int
	sparks                                               []spark
	audio                                                *audio.Context
	gate                                                 audiolab.Gate
	pulse                                                *shaderlab.Pulse
	cam                                                  cameralab.State
	badge                                                *ebiten.Image
}

func newGame() *game {
	g := &game{stage: 1}
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{Pos: cameralab.Vec{X: screenW / 2, Y: screenH / 2}, ViewW: screenW, ViewH: screenH}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{255, 180, 70, 255})
	g.loadStage()
	return g
}
func (g *game) loadStage() {
	players := []point{{1, 1}, {1, 7}, {7, 7}}
	enemies := []point{{7, 7}, {7, 1}, {1, 1}}
	exits := []point{{7, 1}, {1, 1}, {7, 1}}
	walls := [][]point{{{3, 1}, {5, 1}, {1, 3}, {3, 3}, {5, 3}, {7, 3}, {1, 5}, {3, 5}, {5, 5}, {7, 5}, {3, 7}, {5, 7}}, {{1, 1}, {3, 1}, {5, 1}, {3, 3}, {7, 3}, {1, 5}, {5, 5}, {7, 5}, {1, 7}, {3, 7}, {5, 7}}, {{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 3}, {5, 3}, {7, 3}, {1, 5}, {3, 5}, {7, 5}, {1, 7}, {3, 7}, {5, 7}}}
	g.player = players[g.stage-1]
	g.enemy = enemies[g.stage-1]
	g.exit = exits[g.stage-1]
	g.soft = map[point]bool{}
	for _, p := range walls[g.stage-1] {
		g.soft[p] = true
	}
	g.power = 2
	g.capacity = 1
	g.speed = 1
	g.enemyAlive = true
	g.requiredBreak = []int{6, 8, 10}[g.stage-1]
	g.bombs = nil
	g.flames = nil
	g.items = nil
	g.broken = 0
	g.frames = 0
	g.exitOpen = false
	g.won = false
	g.lost = false
	g.message = "Break walls, collect upgrades, defeat the scout, then reach EXIT."
}
func hard(p point) bool {
	return p.x < 0 || p.x >= cols || p.y < 0 || p.y >= rows || p.x == 0 || p.y == 0 || p.x == cols-1 || p.y == rows-1 || (p.x%2 == 0 && p.y%2 == 0)
}
func (g *game) blocked(p point) bool { return hard(p) || g.soft[p] || g.bombAt(p) }
func (g *game) Update() error {
	if g.won || g.lost {
		if retry() {
			if g.won && g.stage < 3 {
				g.totalFrames += g.frames
				g.stage++
				g.loadStage()
			} else {
				best := g.bestFrames
				if g.won {
					total := g.totalFrames + g.frames
					if best == 0 || total < best {
						best = total
					}
				}
				*g = *newGame()
				g.bestFrames = best
			}
		}
		return nil
	}
	g.frames++
	if g.shake > 0 {
		g.shake--
	}
	for i := len(g.sparks) - 1; i >= 0; i-- {
		p := &g.sparks[i]
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life <= 0 {
			g.sparks = append(g.sparks[:i], g.sparks[i+1:]...)
		}
	}
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
		g.play(300)
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
			g.play(110)
			g.shake = 6
			g.burst(tileCX(b.at), tileCY(b.at), 18)
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
	if !g.enemyAlive && g.broken >= g.requiredBreak {
		g.exitOpen = true
	}
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * 6.283 / float64(n)
		g.sparks = append(g.sparks, spark{x, y, math.Cos(a) * float64(1+i%3), math.Sin(a) * float64(1+i%3), 24 + float64(i%8)})
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
			g.play(720)
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
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Square, freq, .06)).Play()
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
	bgs := []color.RGBA{{7, 15, 29, 255}, {30, 20, 48, 255}, {53, 20, 29, 255}}
	s.Fill(bgs[g.stage-1])
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.frames)*2) * 5
	}
	g.drawTitle(s)
	g.drawEffectBadge(s)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("STAGE %d/3 WALLS %d/%d POWER %d BOMBS %d TIME %02d BEST %02d", g.stage, g.broken, g.requiredBreak, g.power, g.capacity, max(0, 90-g.frames/60), g.bestFrames/60), 35, 43)
	ebitenutil.DebugPrintAt(s, g.message, 28, 72)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			p := point{x, y}
			px, py := float32(boardX+x*cell)+float32(ox), float32(boardY+y*cell)
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
		trackatlas.DrawCentered(s, "flame", tileCX(f.at)+ox, tileCY(f.at), 34+math.Sin(float64(f.timer)*.7)*5)
	}
	for _, p := range g.sparks {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/14), color.RGBA{255, 180, 70, 255}, true)
	}
	trackatlas.DrawCentered(s, "hero", tileCX(g.player)+ox, tileCY(g.player), 34)
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
		msg := "STAGE CLEAR!\n\nTAP / ENTER FOR NEXT STAGE"
		if g.stage == 3 {
			msg = "ALL MAZES CLEAR!\n\nTAP / ENTER FOR A NEW RUN"
		}
		overlay(s, msg)
	}
	if g.lost {
		overlay(s, "MISSION FAILED\n\nTAP / ENTER TO RETRY")
	}
}
func (g *game) drawTitle(s *ebiten.Image) {
	if face, err := uilab.Face("en", 16); err == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(203, 5)
		text.Draw(s, "EBI BOMBER", face, op)
		return
	}
	ebitenutil.DebugPrintAt(s, "EBI BOMBER", 203, 17)
}
func (g *game) drawEffectBadge(s *ebiten.Image) {
	if g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.frames)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(screenW-34, 10)
	s.DrawImage(fx, op)
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
