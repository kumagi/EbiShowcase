package towerdefenseplay

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/ogfont"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/towerdefense"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const W, H = 480, 720

type Config struct {
	Step                                int
	Title                               string
	BackgroundPNG, TowersPNG, BattlePNG []byte
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
	scenario                                                                              int
	running, won, lost                                                                    bool
	message                                                                               string
	lang                                                                                  string
	audio                                                                                 *audio.Context
	gate                                                                                  audiolab.Gate
	pulse                                                                                 *shaderlab.Pulse
	cam                                                                                   cameralab.State
	badge                                                                                 *ebiten.Image
	background, towerArt, battleArt                                                       *ebiten.Image
}

var (
	tdFontOnce sync.Once
	tdFontBase *opentype.Font
	tdFontErr  error
	tdFaces    = map[float64]font.Face{}
)

func tdFace(size float64) font.Face {
	tdFontOnce.Do(func() { tdFontBase, tdFontErr = opentype.Parse(ogfont.NotoSansJP) })
	if tdFontErr != nil {
		panic(tdFontErr)
	}
	if f := tdFaces[size]; f != nil {
		return f
	}
	f, err := opentype.NewFace(tdFontBase, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		panic(err)
	}
	tdFaces[size] = f
	return f
}
func tdLabel(s *ebiten.Image, v string, x, y int, c color.Color, size float64) {
	text.Draw(s, v, tdFace(size), x, y, c)
}

var maps = [][]towerdefense.Vec{
	{{-20, 150}, {120, 150}, {120, 300}, {350, 300}, {350, 470}, {500, 470}},
	{{-20, 450}, {105, 450}, {105, 220}, {275, 220}, {275, 390}, {500, 390}},
}

var scenarios = []towerdefense.Scenario{
	{Name: "COAST WATCH", Goal: "Three quick tide waves", Route: maps[0], Waves: 3, Coins: 180, Lives: 10, SpeedScale: .9},
	{Name: "REEF CAVE", Goal: "Slow the armored cave swarm", Route: maps[1], Waves: 4, Coins: 220, Lives: 8, SpeedScale: 1.15},
	{Name: "PEARL GATE", Goal: "Survive the final boss gate", Route: maps[0], Waves: 6, Coins: 260, Lives: 10, SpeedScale: 1, Boss: true},
}

func newGame(cfg Config) *game {
	s := scenarios[0]
	g := &game{cfg: cfg, lives: s.Lives, coins: s.Coins, selected: 0, lang: browserLanguage()}
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{ViewW: W, ViewH: H}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{46, 230, 200, 255})
	if len(cfg.BackgroundPNG) > 0 {
		g.background, g.towerArt, g.battleArt = decodeArt(cfg.BackgroundPNG), decodeArt(cfg.TowersPNG), decodeArt(cfg.BattlePNG)
	}
	g.message = g.tr("Press SPACE or START to launch the wave.", "SPACEかSTARTでウェーブを始めよう。")
	g.path = towerdefense.NewPath(s.Route)
	g.best = storedBest(g.bestKey())
	if cfg.Step >= 2 && cfg.Step <= 4 {
		g.towers = []tower{{towerdefense.Vec{240, 375}, 0, 1, 0}}
	}
	if cfg.Step == 8 {
		// The capstone opens on an active, readable defense line. Learners can
		// still sell nothing and freely add/upgrade towers with the remaining
		// budget, but the first thumbnail is already a battle.
		g.towers = []tower{{towerdefense.Vec{82, 238}, 0, 1, 12}, {towerdefense.Vec{245, 385}, 1, 1, 26}}
		g.coins -= 100
		g.beginWave()
	}
	if cfg.Step <= 4 {
		g.running = true
		g.beginWave()
	}
	return g
}

func decodeArt(data []byte) *ebiten.Image {
	im, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(im)
}
func (g *game) tr(en, ja string) string {
	if g.lang == "ja" {
		return ja
	}
	return en
}
func (g *game) currentScenario() towerdefense.Scenario { return scenarios[g.scenario] }
func (g *game) bestKey() string                        { return fmt.Sprintf("ebiShowcase.defense.best.%d", g.scenario) }
func (g *game) resetScenario(index int) {
	index = min(len(scenarios)-1, max(0, index))
	*g = *newGame(g.cfg)
	g.scenario = index
	s := g.currentScenario()
	g.path = towerdefense.NewPath(s.Route)
	g.lives, g.coins, g.best = s.Lives, s.Coins, storedBest(g.bestKey())
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
	if g.cfg.Step == 8 && g.currentScenario().Boss && g.wave == g.maxWaves() {
		g.spawnLeft = 1
	}
	g.spawnTimer = 1
	g.message = fmt.Sprintf(g.tr("WAVE %d: protect the pearl!", "第%d波：真珠を守ろう！"), g.wave)
}
func (g *game) maxWaves() int {
	if g.cfg.Step == 8 {
		return g.currentScenario().Waves
	}
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
	if g.cfg.Step == 8 {
		for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
			if inpututil.IsKeyJustPressed(key) {
				g.resetScenario(i)
				return nil
			}
		}
	}
	if x, y, ok := press(); ok {
		if g.cfg.Step == 8 && y >= 86 && y < 128 {
			g.resetScenario(min(2, x/160))
			return nil
		} else if y >= 640 {
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
	boss := g.cfg.Step == 8 && g.currentScenario().Boss && g.wave == g.maxWaves()
	hp := 28 + float64(g.wave*8)
	speed := (.75 + float64(g.wave)*.08) * g.currentScenario().SpeedScale
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
				g.message = g.tr("The pearl gate fell.", "真珠の門が破られた。")
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
			g.message = g.tr("The outlined runner is furthest along inside range.", "黄色い輪の敵が射程内で最も先へ進んでいます。")
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
				g.play(680)
				g.towers[i].level++
				g.burst(pos, color.RGBA{100, 230, 180, 255})
				g.message = g.tr("Tower upgraded!", "塔を強化した！")
			}
			return
		}
	}
	cost := 60 + g.selected*20
	if g.coins < cost || g.nearPath(pos) {
		g.message = g.tr("Need coins and open ground.", "コインと空いた地面が必要です。")
		return
	}
	g.coins -= cost
	g.play(520)
	g.towers = append(g.towers, tower{pos, g.selected, 1, 0})
	g.burst(pos, color.RGBA{95, 205, 255, 255})
	g.message = g.tr("Tower placed. Tap it later to upgrade.", "塔を置いた。あとでタップすると強化できます。")
}
func (g *game) play(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, freq, .05)).Play()
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
		g.message = g.tr("Path and targeting observation complete!", "経路と標的の観察完了！")
		return
	}
	if g.wave >= g.maxWaves() {
		g.won = true
		grade := g.score + g.lives*100 + g.coins
		if grade > g.best {
			g.best = grade
			storeBest(g.bestKey(), g.best)
		}
		g.message = g.tr("Every wave cleared!", "すべてのウェーブを守り切った！")
		return
	}
	g.coins += 40
	g.message = g.tr("Wave clear! Build, upgrade, then START.", "ウェーブ完了！ 建設・強化してSTART。")
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
	if g.background != nil {
		drawCover(s, g.background)
	} else {
		s.Fill(bg)
	}
	vector.DrawFilledRect(s, 0, 0, W, H, color.RGBA{2, 10, 26, 42}, false)
	// Underwater depth layers: bubbles, reef silhouettes and a soft vignette.
	for i := 0; i < 24; i++ {
		x := float32((i*83 + g.tick/4) % W)
		y := float32(130 + (i*57-g.tick/3+720)%430)
		vector.StrokeCircle(s, x, y, float32(2+i%5), 1, color.RGBA{116, 218, 226, 65}, true)
	}
	if g.background == nil {
		for i := 0; i < 9; i++ {
			x := float32(i*65 - 18)
			h := float32(32 + (i*37)%75)
			vector.DrawFilledRect(s, x, 558-h, 32, h, color.RGBA{18, 62, 66, 180}, false)
			vector.DrawFilledCircle(s, x+16, 558-h, 19, color.RGBA{26, 86, 79, 190}, false)
		}
	}
	g.drawEffectBadge(s)
	ox := 0.0
	if g.shake > 0 {
		ox = math.Sin(float64(g.tick)*2.2) * float64(g.shake)
	}
	g.drawPath(s, ox)
	tdLabel(s, g.cfg.Title, 138, 26, color.RGBA{46, 230, 200, 255}, 16)
	tdLabel(s, fmt.Sprintf("%s %d/%d  %s %d  %s %d  %s %d", g.tr("WAVE", "波"), g.wave, g.maxWaves(), g.tr("LIVES", "ライフ"), g.lives, g.tr("COINS", "コイン"), g.coins, g.tr("SCORE", "得点"), g.score), 40, 53, color.White, 13)
	tdLabel(s, g.message, 38, 79, color.RGBA{255, 211, 112, 255}, 13)
	if g.cfg.Step == 8 {
		for i := range scenarios {
			c := color.RGBA{38, 62, 92, 255}
			if i == g.scenario {
				c = color.RGBA{45, 144, 151, 255}
			}
			vector.DrawFilledRect(s, float32(i*160+3), 88, 154, 36, c, false)
			ebitenutil.DebugPrintAt(s, fmt.Sprintf("[%d] %s", i+1, []string{"COAST", "CAVE", "GATE"}[i]), i*160+14, 101)
		}
		trait := "RUSH"
		if g.currentScenario().SpeedScale > 1 {
			trait = "FAST"
		}
		if g.currentScenario().Boss {
			trait = "BOSS AT FINAL WAVE"
		}
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("INTENT: %s · NEXT %d · %s", trait, g.spawnLeft, g.currentScenario().Goal), 48, 140)
	}
	for i, t := range g.towers {
		r := 105 + float64(t.level*8)
		if g.cfg.Step == 2 || g.cfg.Step == 3 || i == len(g.towers)-1 {
			vector.StrokeCircle(s, float32(t.pos.X+ox), float32(t.pos.Y), float32(r), 2, color.RGBA{95, 185, 225, 100}, true)
		}
		if g.towerArt != nil {
			drawAtlas(s, g.towerArt, 3, t.kind, t.pos.X+ox, t.pos.Y, 58+float64(t.level*5))
		} else {
			trackatlas.DrawCentered(s, []string{"scout", "species-2", "king-crab"}[t.kind], t.pos.X+ox, t.pos.Y, 50+float64(t.level*4))
		}
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
		unit := 0
		if g.wave >= 3 {
			unit = 1
		}
		if e.boss {
			unit = 2
		}
		if g.battleArt != nil {
			drawAtlas(s, g.battleArt, 7, unit, p.X+ox, p.Y, size)
		} else {
			trackatlas.DrawCentered(s, map[bool]string{true: "boss-crab", false: "slug"}[e.boss], p.X+ox, p.Y, size)
		}
		vector.DrawFilledRect(s, float32(p.X-22+ox), float32(p.Y-30), 44, 5, color.RGBA{40, 45, 55, 255}, false)
		vector.DrawFilledRect(s, float32(p.X-22+ox), float32(p.Y-30), float32(44*e.hp/e.maxHP), 5, color.RGBA{235, 85, 80, 255}, false)
		if i == targetIdx {
			vector.StrokeCircle(s, float32(p.X+ox), float32(p.Y), 30, 3, color.RGBA{255, 220, 90, 255}, true)
		}
	}
	for _, p := range g.shots {
		if g.battleArt != nil {
			drawAtlas(s, g.battleArt, 7, 4+p.kind, p.pos.X+ox, p.pos.Y, 22)
		} else {
			vector.DrawFilledCircle(s, float32(p.pos.X+ox), float32(p.pos.Y), 6, color.RGBA{255, 205, 65, 255}, true)
		}
	}
	for _, p := range g.particles {
		vector.DrawFilledCircle(s, float32(p.x+ox), float32(p.y), float32(2+p.life/10), p.c, true)
	}
	if g.cfg.Step >= 7 {
		labels := []string{g.tr("[Q] BOLT $60", "[Q] 弾 $60"), g.tr("[W] SLOW $80", "[W] 鈍足 $80"), g.tr("[E] SPLASH $100", "[E] 範囲 $100")}
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
	label := g.tr("START NEXT WAVE / SPACE", "次の波をSTART / SPACE")
	if g.running {
		label = g.tr("WAVE IN PROGRESS", "ウェーブ進行中")
	}
	ebitenutil.DebugPrintAt(s, label, 145, 663)
	if g.won {
		g.overlay(s, fmt.Sprintf(g.tr("DEFENSE COMPLETE! GRADE %s\nBEST %d\n\nTAP / R: PLAY AGAIN", "防衛成功！ 評価 %s\nBEST %d\n\nタップ / R：もう一度"), towerdefense.ResultGrade(g.score, g.lives, g.coins), g.best))
	}
	if g.lost {
		g.overlay(s, g.tr("THE PEARL WAS TAKEN\n\nTAP / R: RETRY", "真珠が奪われた\n\nタップ / R：再挑戦"))
	}
}
func (g *game) drawEffectBadge(s *ebiten.Image) {
	if g.cfg.Step != 8 || g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.tick)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(W-34, 10)
	s.DrawImage(fx, op)
}
func (g *game) drawPath(s *ebiten.Image, ox float64) {
	for i := 1; i < len(g.path.Points); i++ {
		a, b := g.path.Points[i-1], g.path.Points[i]
		vector.StrokeLine(s, float32(a.X+ox), float32(a.Y), float32(b.X+ox), float32(b.Y), 46, color.RGBA{70, 76, 82, 255}, true)
		vector.StrokeLine(s, float32(a.X+ox), float32(a.Y), float32(b.X+ox), float32(b.Y), 3, color.RGBA{215, 183, 105, 255}, true)
		vector.StrokeLine(s, float32(a.X+ox), float32(a.Y-16), float32(b.X+ox), float32(b.Y-16), 2, color.RGBA{116, 218, 226, 80}, true)
	}
	if g.battleArt != nil {
		drawAtlas(s, g.battleArt, 7, 3, 455+ox, g.path.Points[len(g.path.Points)-1].Y, 58)
	} else {
		trackatlas.DrawCentered(s, "pearl", 455+ox, g.path.Points[len(g.path.Points)-1].Y, 45)
	}
}

func drawAtlas(dst, atlas *ebiten.Image, columns, index int, cx, cy, size float64) {
	w, h := atlas.Bounds().Dx()/columns, atlas.Bounds().Dy()
	src := atlas.SubImage(image.Rect(index*w, 0, (index+1)*w, h)).(*ebiten.Image)
	scale := size / float64(max(w, h))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(cx-float64(w)*scale/2, cy-float64(h)*scale/2)
	dst.DrawImage(src, op)
}
func drawCover(dst, src *ebiten.Image) {
	w, h := float64(src.Bounds().Dx()), float64(src.Bounds().Dy())
	scale := math.Max(W/w, H/h)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((W-w*scale)/2, (H-h*scale)/2)
	dst.DrawImage(src, op)
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
func (g *game) overlay(s *ebiten.Image, text string) {
	vector.DrawFilledRect(s, 35, 260, 410, 165, color.RGBA{4, 12, 24, 244}, false)
	vector.StrokeRect(s, 35, 260, 410, 165, 4, color.RGBA{245, 190, 65, 255}, false)
	for i, line := range strings.Split(text, "\n") {
		tdLabel(s, line, 90, 320+i*30, color.White, 18)
	}
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func Run(cfg Config) {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle(cfg.Title + " — Ebitengine")
	if err := ebiten.RunGame(newGame(cfg)); err != nil {
		panic(err)
	}
}
