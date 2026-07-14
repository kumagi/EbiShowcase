package towerdefenseplay

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/towerdefense"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const W, H = 480, 720

type Config struct {
	Step  int
	Title string
}
type enemy struct {
	id                               int
	progress, hp, maxHP, speed, slow float64
	boss                             bool
}
type tower struct {
	pos                   towerdefense.Vec
	kind, level, cooldown int
}
type projectile struct {
	pos                  towerdefense.Vec
	target, damage, kind int
	life                 int
}
type particle struct {
	x, y, vx, vy, life float64
	c                  color.RGBA
}
type game struct {
	cfg                                                                                   Config
	path                                                                                  towerdefense.Path
	enemies                                                                               []enemy
	towers                                                                                []tower
	shots                                                                                 []projectile
	particles                                                                             []particle
	nextID, wave, spawnLeft, spawnTimer, lives, coins, score, best, selected, tick, shake int
	running, won, lost                                                                    bool
	message                                                                               string
}

var maps = [][]towerdefense.Vec{
	{{-20, 150}, {120, 150}, {120, 300}, {350, 300}, {350, 470}, {500, 470}},
	{{-20, 450}, {105, 450}, {105, 220}, {275, 220}, {275, 390}, {500, 390}},
}

func newGame(cfg Config) *game {
	g := &game{cfg: cfg, lives: 10, coins: 160, selected: 0, message: "Press SPACE or START to launch the wave."}
	g.path = towerdefense.NewPath(maps[0])
	if cfg.Step >= 2 && cfg.Step <= 4 {
		g.towers = []tower{{towerdefense.Vec{240, 375}, 0, 1, 0}}
	}
	if cfg.Step <= 4 {
		g.running = true
		g.beginWave()
	}
	return g
}
func (g *game) beginWave() {
	if g.running && g.spawnLeft > 0 {
		return
	}
	g.wave++
	g.running = true
	g.spawnLeft = 3 + g.wave
	if g.cfg.Step <= 3 {
		g.spawnLeft = 3
	}
	if g.cfg.Step == 4 {
		g.spawnLeft = 5
	}
	if g.cfg.Step == 8 && g.wave == 6 {
		g.spawnLeft = 1
	}
	g.spawnTimer = 1
	g.message = fmt.Sprintf("WAVE %d: protect the pearl!", g.wave)
}
func (g *game) maxWaves() int {
	switch {
	case g.cfg.Step <= 4:
		return 1
	case g.cfg.Step == 5:
		return 1
	case g.cfg.Step <= 7:
		return 3
	default:
		return 6
	}
}
func (g *game) Update() error {
	g.tick++
	if g.shake > 0 {
		g.shake--
	}
	g.updateParticles()
	if g.won || g.lost {
		if retryPressed() {
			best := g.best
			*g = *newGame(g.cfg)
			g.best = best
		}
		return nil
	}
	if x, y, ok := press(); ok {
		if y >= 640 {
			if !g.running {
				g.beginWave()
			}
		} else if y >= 570 && g.cfg.Step >= 7 {
			g.selected = min(2, x/160)
		} else if y < 555 && g.cfg.Step >= 5 {
			g.placeOrUpgrade(towerdefense.Vec{float64(x), float64(y)})
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && !g.running {
		g.beginWave()
	}
	if g.cfg.Step >= 7 {
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.selected = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyW) {
			g.selected = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			g.selected = 2
		}
	}
	if g.running {
		g.updateSpawn()
		g.updateEnemies()
		if g.cfg.Step >= 3 {
			g.updateTowers()
		}
		if g.cfg.Step >= 4 {
			g.updateShots()
		}
		g.checkWave()
	}
	return nil
}
func (g *game) updateSpawn() {
	if g.spawnLeft <= 0 {
		return
	}
	g.spawnTimer--
	if g.spawnTimer > 0 {
		return
	}
	boss := g.cfg.Step == 8 && g.wave == 6
	hp := 28 + float64(g.wave*8)
	speed := .75 + float64(g.wave)*.08
	if boss {
		hp = 320
		speed = .55
	}
	g.enemies = append(g.enemies, enemy{g.nextID, 0, hp, hp, speed, 0, boss})
	g.nextID++
	g.spawnLeft--
	g.spawnTimer = max(35, 80-g.wave*5)
}
func (g *game) updateEnemies() {
	for i := len(g.enemies) - 1; i >= 0; i-- {
		e := &g.enemies[i]
		factor := 1.0
		if e.slow > 0 {
			factor = .55
			e.slow--
		}
		e.progress += e.speed * factor
		if e.progress >= g.path.Total {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			if g.cfg.Step <= 3 {
				g.score++
				continue
			}
			g.lives--
			g.shake = 8
			if g.lives <= 0 {
				g.lost = true
				g.message = "The pearl gate fell."
			}
		}
	}
}
func (g *game) updateTowers() {
	for i := range g.towers {
		t := &g.towers[i]
		if t.cooldown > 0 {
			t.cooldown--
			continue
		}
		radius := 105.0 + float64(t.level*8)
		targets := make([]towerdefense.Target, len(g.enemies))
		for j, e := range g.enemies {
			targets[j] = towerdefense.Target{g.path.Position(e.progress), e.progress, e.hp > 0}
		}
		idx := towerdefense.SelectFront(t.pos, radius, targets)
		if idx < 0 {
			continue
		}
		if g.cfg.Step == 3 {
			g.message = "The outlined runner is furthest along inside range."
			continue
		}
		kind := t.kind
		damage := 10 + t.level*4
		if kind == 2 {
			damage = 8 + t.level*3
		}
		g.shots = append(g.shots, projectile{t.pos, g.enemies[idx].id, damage, kind, 90})
		t.cooldown = 45 - t.level*4
		if kind == 1 {
			t.cooldown += 8
		}
	}
}
func (g *game) updateShots() {
	for i := len(g.shots) - 1; i >= 0; i-- {
		p := &g.shots[i]
		idx := g.enemyIndex(p.target)
		if idx < 0 {
			g.shots = append(g.shots[:i], g.shots[i+1:]...)
			continue
		}
		target := g.path.Position(g.enemies[idx].progress)
		dx, dy := target.X-p.pos.X, target.Y-p.pos.Y
		d := math.Hypot(dx, dy)
		if d < 13 {
			g.hit(idx, *p, target)
			g.shots = append(g.shots[:i], g.shots[i+1:]...)
			continue
		}
		p.pos.X += dx / d * 7
		p.pos.Y += dy / d * 7
		p.life--
		if p.life <= 0 {
			g.shots = append(g.shots[:i], g.shots[i+1:]...)
		}
	}
}
func (g *game) hit(idx int, p projectile, pos towerdefense.Vec) {
	if p.kind == 2 {
		for i := range g.enemies {
			q := g.path.Position(g.enemies[i].progress)
			if towerdefense.Distance(pos, q) < 55 {
				g.enemies[i].hp -= float64(p.damage)
			}
		}
	} else {
		g.enemies[idx].hp -= float64(p.damage)
	}
	if p.kind == 1 {
		g.enemies[idx].slow = 90
	}
	g.shake = 3
	for n := 0; n < 10; n++ {
		a := float64(n) * math.Pi / 5
		g.particles = append(g.particles, particle{pos.X, pos.Y, math.Cos(a) * 2.5, math.Sin(a) * 2.5, 22, color.RGBA{255, 195, 65, 255}})
	}
	for i := len(g.enemies) - 1; i >= 0; i-- {
		if g.enemies[i].hp <= 0 {
			reward := 18
			if g.enemies[i].boss {
				reward = 150
			}
			g.coins += reward
			g.score += reward * 10
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
		}
	}
}
func (g *game) enemyIndex(id int) int {
	for i, e := range g.enemies {
		if e.id == id {
			return i
		}
	}
	return -1
}
func (g *game) placeOrUpgrade(pos towerdefense.Vec) {
	for i := range g.towers {
		if towerdefense.Distance(pos, g.towers[i].pos) < 30 {
			cost := 45 + g.towers[i].level*25
			if g.cfg.Step >= 7 && g.coins >= cost && g.towers[i].level < 3 {
				g.coins -= cost
				g.towers[i].level++
				g.burst(pos, color.RGBA{100, 230, 180, 255})
				g.message = "Tower upgraded!"
			}
			return
		}
	}
	cost := 60 + g.selected*20
	if g.coins < cost || g.nearPath(pos) {
		g.message = "Need coins and open ground."
		return
	}
	g.coins -= cost
	g.towers = append(g.towers, tower{pos, g.selected, 1, 0})
	g.burst(pos, color.RGBA{95, 205, 255, 255})
	g.message = "Tower placed. Tap it later to upgrade."
}
func (g *game) nearPath(pos towerdefense.Vec) bool {
	for d := 0.0; d <= g.path.Total; d += 12 {
		if towerdefense.Distance(pos, g.path.Position(d)) < 42 {
			return true
		}
	}
	return pos.Y < 100 || pos.Y > 550
}
func (g *game) checkWave() {
	if g.spawnLeft > 0 || len(g.enemies) > 0 || len(g.shots) > 0 {
		return
	}
	g.running = false
	if g.cfg.Step <= 3 {
		g.won = true
		g.message = "Path and targeting observation complete!"
		return
	}
	if g.wave >= g.maxWaves() {
		g.won = true
		grade := g.score + g.lives*100 + g.coins
		if grade > g.best {
			g.best = grade
		}
		g.message = "Every wave cleared!"
		return
	}
	if g.cfg.Step == 8 && g.wave == 3 {
		g.path = towerdefense.NewPath(maps[1])
		g.towers = nil
		g.coins += 140
		g.message = "MAP 2 OPEN! Build for the new route."
	} else {
		g.coins += 40
		g.message = "Wave clear! Build, upgrade, then START."
	}
}
func (g *game) updateParticles() {
	for i := len(g.particles) - 1; i >= 0; i-- {
		p := &g.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += .04
		p.life--
		if p.life <= 0 {
			g.particles = append(g.particles[:i], g.particles[i+1:]...)
		}
	}
}
func (g *game) burst(pos towerdefense.Vec, c color.RGBA) {
	for i := 0; i < 14; i++ {
		a := float64(i) * math.Pi / 7
		g.particles = append(g.particles, particle{pos.X, pos.Y, math.Cos(a) * 2, math.Sin(a) * 2, 24, c})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	bg := color.RGBA{12, 29, 44, 255}
	if g.cfg.Step == 8 && g.wave >= 3 {
		bg = color.RGBA{31, 22, 49, 255}
	}
	s.Fill(bg)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2.2) * float64(g.shake)
	}
	g.drawPath(s, ox)
	ebitenutil.DebugPrintAt(s, g.cfg.Title, 150, 18)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("WAVE %d/%d  LIVES %d  COINS %d  SCORE %d", g.wave, g.maxWaves(), g.lives, g.coins, g.score), 70, 46)
	ebitenutil.DebugPrintAt(s, g.message, 55, 72)
	for i, t := range g.towers {
		r := 105 + float64(t.level*8)
		if g.cfg.Step == 2 || g.cfg.Step == 3 || i == len(g.towers)-1 {
			vector.StrokeCircle(s, float32(t.pos.X+ox), float32(t.pos.Y), float32(r), 2, color.RGBA{95, 185, 225, 100}, true)
		}
		sprite := []string{"scout", "species-2", "king-crab"}[t.kind]
		trackatlas.DrawCentered(s, sprite, t.pos.X+ox, t.pos.Y, 50+float64(t.level*4))
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("L%d", t.level), int(t.pos.X+ox)-7, int(t.pos.Y)+24)
	}
	targetIdx := -1
	if g.cfg.Step >= 3 && len(g.towers) > 0 {
		targets := make([]towerdefense.Target, len(g.enemies))
		for i, e := range g.enemies {
			targets[i] = towerdefense.Target{g.path.Position(e.progress), e.progress, true}
		}
		targetIdx = towerdefense.SelectFront(g.towers[0].pos, 115, targets)
	}
	for i, e := range g.enemies {
		p := g.path.Position(e.progress)
		size := 42.0
		if e.boss {
			size = 78
		}
		trackatlas.DrawCentered(s, map[bool]string{true: "boss-crab", false: "slug"}[e.boss], p.X+ox, p.Y, size)
		vector.DrawFilledRect(s, float32(p.X-22+ox), float32(p.Y-30), 44, 5, color.RGBA{40, 45, 55, 255}, false)
		vector.DrawFilledRect(s, float32(p.X-22+ox), float32(p.Y-30), float32(44*e.hp/e.maxHP), 5, color.RGBA{235, 85, 80, 255}, false)
		if i == targetIdx {
			vector.StrokeCircle(s, float32(p.X+ox), float32(p.Y), 30, 3, color.RGBA{255, 220, 90, 255}, true)
		}
	}
	for _, p := range g.shots {
		vector.DrawFilledCircle(s, float32(p.pos.X+ox), float32(p.pos.Y), 6, color.RGBA{255, 205, 65, 255}, true)
	}
	for _, p := range g.particles {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/10), p.c, true)
	}
	if g.cfg.Step >= 7 {
		labels := []string{"[Q] BOLT $60", "[W] SLOW $80", "[E] SPLASH $100"}
		for i, l := range labels {
			c := color.RGBA{44, 68, 92, 255}
			if i == g.selected {
				c = color.RGBA{176, 98, 50, 255}
			}
			vector.DrawFilledRect(s, float32(i*160+4), 570, 152, 58, c, false)
			ebitenutil.DebugPrintAt(s, l, i*160+18, 595)
		}
	}
	vector.DrawFilledRect(s, 40, 640, 400, 55, color.RGBA{180, 98, 48, 255}, false)
	label := "START NEXT WAVE / SPACE"
	if g.running {
		label = "WAVE IN PROGRESS"
	}
	ebitenutil.DebugPrintAt(s, label, 145, 663)
	if g.won {
		overlay(s, fmt.Sprintf("DEFENSE COMPLETE! BEST %d\n\nTAP / R: PLAY AGAIN", g.best))
	}
	if g.lost {
		overlay(s, "THE PEARL WAS TAKEN\n\nTAP / R: RETRY")
	}
}
func (g *game) drawPath(s *ebiten.Image, ox float64) {
	for i := 1; i < len(g.path.Points); i++ {
		a, b := g.path.Points[i-1], g.path.Points[i]
		vector.StrokeLine(s, float32(a.X+ox), float32(a.Y), float32(b.X+ox), float32(b.Y), 46, color.RGBA{70, 76, 82, 255}, true)
		vector.StrokeLine(s, float32(a.X+ox), float32(a.Y), float32(b.X+ox), float32(b.Y), 3, color.RGBA{215, 183, 105, 255}, true)
	}
	trackatlas.DrawCentered(s, "pearl", 455+ox, g.path.Points[len(g.path.Points)-1].Y, 45)
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
func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, text string) {
	vector.DrawFilledRect(s, 35, 260, 410, 165, color.RGBA{4, 12, 24, 244}, false)
	vector.StrokeRect(s, 35, 260, 410, 165, 4, color.RGBA{245, 190, 65, 255}, false)
	ebitenutil.DebugPrintAt(s, text, 90, 320)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func Run(cfg Config) {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle(cfg.Title + " — Ebitengine")
	if err := ebiten.RunGame(newGame(cfg)); err != nil {
		panic(err)
	}
}
