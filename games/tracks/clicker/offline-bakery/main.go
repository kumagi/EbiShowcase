package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
)

const width, height = 480, 720
const prefix = "ebiShowcaseBakery."

type game struct {
	sweets                              float64
	total                               float64
	ovens, mixers, shops, frame         int
	cost, mixerCost, shopCost, lastSave float64
	last                                time.Time
	offline                             float64
	clear                               bool
	particles                           []crumb
	milestone                           int
}
type crumb struct{ x, y, vx, vy, life float64 }

func newGame() *game {
	g := &game{cost: 20, mixerCost: 180, shopCost: 1600, last: time.Now()}
	g.load()
	return g
}
func (g *game) rate() float64     { return float64(g.ovens)*2 + float64(g.mixers)*15 + float64(g.shops)*80 }
func (g *game) tapPower() float64 { return 1 + float64(g.ovens)*.25 + float64(g.mixers) }
func (g *game) load() {
	raw, ok := storageGet(prefix + "sweets")
	if !ok {
		return
	}
	g.sweets, _ = strconv.ParseFloat(raw, 64)
	total, _ := storageGet(prefix + "total")
	ovens, _ := storageGet(prefix + "ovens")
	mixers, _ := storageGet(prefix + "mixers")
	shops, _ := storageGet(prefix + "shops")
	g.total, _ = strconv.ParseFloat(total, 64)
	g.ovens, _ = strconv.Atoi(ovens)
	g.mixers, _ = strconv.Atoi(mixers)
	g.shops, _ = strconv.Atoi(shops)
	g.cost = 20 * math.Pow(1.25, float64(g.ovens))
	g.mixerCost = 180 * math.Pow(1.32, float64(g.mixers))
	g.shopCost = 1600 * math.Pow(1.38, float64(g.shops))
	timeValue, _ := storageGet(prefix + "time")
	stamp, _ := strconv.ParseInt(timeValue, 10, 64)
	away := math.Min(8*3600, float64(time.Now().Unix()-stamp))
	if away > 0 {
		g.offline = away * g.rate()
		g.sweets += g.offline
		g.total += g.offline
	}
}
func (g *game) save() {
	storageSet(prefix+"sweets", fmt.Sprintf("%.3f", g.sweets))
	storageSet(prefix+"total", fmt.Sprintf("%.3f", g.total))
	storageSet(prefix+"ovens", strconv.Itoa(g.ovens))
	storageSet(prefix+"mixers", strconv.Itoa(g.mixers))
	storageSet(prefix+"shops", strconv.Itoa(g.shops))
	storageSet(prefix+"time", strconv.FormatInt(time.Now().Unix(), 10))
	g.lastSave = 2
}
func (g *game) Update() error {
	now := time.Now()
	dt := math.Min(.25, now.Sub(g.last).Seconds())
	g.last = now
	g.frame++
	if !g.clear {
		made := g.rate() * dt
		g.sweets += made
		g.total += made
	}
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
	g.milestone = 0
	if g.total >= 500 {
		g.milestone = 1
	}
	if g.total >= 5000 {
		g.milestone = 2
	}
	if g.total >= 25000 {
		g.milestone = 3
	}
	if g.lastSave > 0 {
		g.lastSave -= dt
	}
	if g.frame%120 == 0 {
		g.save()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		for _, key := range []string{"sweets", "total", "ovens", "mixers", "shops", "time"} {
			storageRemove(prefix + key)
		}
		*g = *newGame()
		return nil
	}
	if g.clear {
		if any() {
			g.clear = false
			g.total = 0
			g.save()
		}
		return nil
	}
	x, y, ok := press()
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ok = true
		y = 300
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		ok = true
		x, y = 80, 560
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		ok = true
		x, y = 240, 560
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		ok = true
		x, y = 400, 560
	}
	if ok {
		g.offline = 0
		if y < 470 {
			made := g.tapPower()
			g.sweets += made
			g.total += made
			g.burst(240, 280, 12)
		} else if x < 160 && g.sweets >= g.cost {
			g.sweets -= g.cost
			g.ovens++
			g.cost = 20 * math.Pow(1.25, float64(g.ovens))
			g.burst(80, 540, 16)
		} else if x < 320 && g.sweets >= g.mixerCost {
			g.sweets -= g.mixerCost
			g.mixers++
			g.mixerCost = 180 * math.Pow(1.32, float64(g.mixers))
			g.burst(240, 540, 16)
		} else if g.sweets >= g.shopCost {
			g.sweets -= g.shopCost
			g.shops++
			g.shopCost = 1600 * math.Pow(1.38, float64(g.shops))
			g.burst(400, 540, 18)
		}
		g.save()
	}
	if g.total >= 25000 {
		g.clear = true
		g.save()
	}
	return nil
}
func (g *game) burst(x, y float64, n int) {
	for i := 0; i < n; i++ {
		a := float64(i) * math.Pi * 2 / float64(n)
		g.particles = append(g.particles, crumb{x, y, math.Cos(a) * (1 + float64(i%3)), math.Sin(a) * (1 + float64(i%3)), 28 + float64(i%8)})
	}
}
func (g *game) Draw(s *ebiten.Image) {
	backgrounds := []color.RGBA{{39, 25, 43, 255}, {39, 42, 70, 255}, {45, 65, 64, 255}, {78, 53, 35, 255}}
	s.Fill(backgrounds[g.milestone])
	ebitenutil.DebugPrintAt(s, "EBI BAKERY — SAVED IN THIS BROWSER", 110, 40)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SWEETS %s   TOTAL %s / 25.0K", short(g.sweets), short(g.total)), 120, 85)
	bob := math.Sin(float64(g.frame)*.12) * 5
	trackatlas.DrawCentered(s, "bakery", 240, 275+bob, 190+math.Sin(float64(g.frame)*.18)*5)
	for i := 0; i < min(g.ovens, 6); i++ {
		phase := math.Mod(float64(g.frame)*1.7+float64(i)*55, 390)
		trackatlas.DrawCentered(s, "coin", 45+phase, 420+math.Sin(phase*.05)*8, 16)
	}
	for _, p := range g.particles {
		vector.DrawFilledCircle(s, float32(p.x), float32(p.y), float32(2+p.life/15), color.RGBA{255, 210, 95, 255}, true)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("TAP / SPACE: BAKE +%s", short(g.tapPower())), 150, 445)
	drawShop(s, 15, "1 OVEN", g.cost, g.ovens, g.sweets >= g.cost)
	drawShop(s, 170, "2 MIXER", g.mixerCost, g.mixers, g.sweets >= g.mixerCost)
	drawShop(s, 325, "3 SHOP", g.shopCost, g.shops, g.sweets >= g.shopCost)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("PRODUCTION %s / SEC   DISTRICT %d/3", short(g.rate()), g.milestone), 105, 635)
	if g.offline > 0 {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("WELCOME BACK! OFFLINE +%s", short(g.offline)), 120, 655)
	} else if g.lastSave > 0 {
		ebitenutil.DebugPrintAt(s, "SAVED", 220, 655)
	}
	ebitenutil.DebugPrintAt(s, "R: DELETE SAVE", 185, 680)
	if g.clear {
		overlay(s, "BAKERY GOAL REACHED!\n\nTAP / SPACE TO CONTINUE")
	}
}
func drawShop(s *ebiten.Image, x float64, label string, cost float64, owned int, ready bool) {
	c := color.RGBA{55, 86, 105, 255}
	if ready {
		c = color.RGBA{45, 205, 181, 255}
	}
	vector.DrawFilledRect(s, float32(x), 485, 140, 135, c, false)
	ebitenutil.DebugPrintAt(s, label, int(x)+35, 500)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("COST %s\nOWNED %d", short(cost), owned), int(x)+25, 545)
}
func short(v float64) string {
	if v >= 1e6 {
		return fmt.Sprintf("%.2fM", v/1e6)
	}
	if v >= 1e3 {
		return fmt.Sprintf("%.2fK", v/1e3)
	}
	return fmt.Sprintf("%.1f", v)
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
func any() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 240}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Bakery — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
