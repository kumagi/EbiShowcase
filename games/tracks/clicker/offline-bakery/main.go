package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/uilab"
)

const width, height = 480, 720
const prefix = "ebiShowcaseBakery."

// Only this capstone's local generated artwork enters its WASM binary.
//
//go:embed assets/bakery-pearl-palace.png assets/bakery-chef.png assets/bakery-equipment-atlas.png
var bakeryArtFS embed.FS

var (
	bakeryArtOnce sync.Once
	bakeryArt     map[string]*ebiten.Image
	bakeryGear    [3]*ebiten.Image
	bakeryFace14  *text.GoTextFace
	bakeryFace16  *text.GoTextFace
	bakeryFace20  *text.GoTextFace
)

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
	previousMilestone, milestoneFlash   int
	tapFlash, buyFlash, buyLine         int
}
type crumb struct{ x, y, vx, vy, life float64 }

func newGame() *game {
	loadBakeryArt()
	g := &game{cost: 20, mixerCost: 180, shopCost: 1600, last: time.Now()}
	g.load()
	g.milestone = milestoneFor(g.total)
	g.previousMilestone = g.milestone
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
	if g.tapFlash > 0 {
		g.tapFlash--
	}
	if g.buyFlash > 0 {
		g.buyFlash--
	}
	if g.milestoneFlash > 0 {
		g.milestoneFlash--
	}
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
	g.milestone = milestoneFor(g.total)
	if g.milestone > g.previousMilestone {
		g.previousMilestone = g.milestone
		g.milestoneFlash = 150
		g.burst(240, 315, 30)
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
			g.tapFlash = 12
		} else if x < 160 && g.sweets >= g.cost {
			g.sweets -= g.cost
			g.ovens++
			g.cost = 20 * math.Pow(1.25, float64(g.ovens))
			g.burst(80, 540, 16)
			g.buyLine, g.buyFlash = 0, 24
		} else if x < 320 && g.sweets >= g.mixerCost {
			g.sweets -= g.mixerCost
			g.mixers++
			g.mixerCost = 180 * math.Pow(1.32, float64(g.mixers))
			g.burst(240, 540, 16)
			g.buyLine, g.buyFlash = 1, 24
		} else if g.sweets >= g.shopCost {
			g.sweets -= g.shopCost
			g.shops++
			g.shopCost = 1600 * math.Pow(1.38, float64(g.shops))
			g.burst(400, 540, 18)
			g.buyLine, g.buyFlash = 2, 24
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
	drawBakeryCover(s, bakeryArt["palace"], 0, 0, width, height)
	grades := []color.RGBA{{7, 28, 48, 20}, {17, 64, 76, 22}, {85, 38, 92, 28}, {135, 77, 22, 34}}
	vector.DrawFilledRect(s, 0, 0, width, height, grades[g.milestone], false)

	// The HUD turns the artwork into a readable game objective immediately.
	vector.DrawFilledRect(s, 10, 10, 460, 92, color.RGBA{5, 17, 32, 224}, true)
	vector.StrokeRect(s, 10, 10, 460, 92, 2, color.RGBA{255, 220, 145, 188}, true)
	drawCenteredLabel(s, "PEARL PALACE PATISSERIE", 240, 22, bakeryFace16, color.RGBA{255, 235, 188, 255})
	drawCenteredLabel(s, fmt.Sprintf("SWEETS %s   •   +%s / SEC", short(g.sweets), short(g.rate())), 240, 48, bakeryFace16, color.White)
	goalName, goalCost := g.nextTarget()
	drawCenteredLabel(s, fmt.Sprintf("NEXT TARGET  %s  %s", goalName, short(goalCost)), 240, 74, bakeryFace14, color.RGBA{126, 244, 255, 255})

	// The chef and the pastry on the tray are the manual-click product. Tap
	// feedback is a presentation-only transform driven by tapFlash from Update.
	bob := math.Sin(float64(g.frame)*.08) * 3
	chefScale := 1.0
	if g.tapFlash > 0 {
		chefScale = 1.035
	}
	chefW, chefH := 285*chefScale, 342*chefScale
	drawBakeryContain(s, bakeryArt["chef"], 98-(chefW-285)/2, 95+bob-(chefH-342)/2, chefW, chefH, false, 1)
	trayX, trayY := float32(328), float32(246+bob)
	for i := 0; i < 3; i++ {
		r := float32(34+i*8) + float32(math.Sin(float64(g.frame)*.12+float64(i))*2)
		vector.StrokeCircle(s, trayX, trayY, r, float32(4-i), color.RGBA{255, 224, 116, uint8(150 - i*40)}, true)
	}
	if g.tapFlash > 0 {
		vector.StrokeCircle(s, trayX, trayY, float32(52+(12-g.tapFlash)*3), 4, color.RGBA{125, 247, 255, 210}, true)
	}
	drawCenteredLabel(s, fmt.Sprintf("TAP THE PEARL TART  +%s", short(g.tapPower())), 240, 422, bakeryFace16, color.RGBA{255, 245, 213, 255})

	for _, p := range g.particles {
		vector.DrawFilledCircle(s, float32(p.x), float32(p.y), float32(2+p.life/15), color.RGBA{255, 218, 116, 230}, true)
	}

	// The three generated machines are both visible production lines and the
	// exact items purchased by the lower cards. Locked lines remain recognizable.
	drawLine(s, 0, 14, 448, g.ovens, g.frame, g.buyLine == 0 && g.buyFlash > 0)
	drawLine(s, 1, 170, 448, g.mixers, g.frame, g.buyLine == 1 && g.buyFlash > 0)
	drawLine(s, 2, 326, 448, g.shops, g.frame, g.buyLine == 2 && g.buyFlash > 0)
	drawShop(s, 0, 15, "1  OVEN", g.cost, g.ovens, g.sweets >= g.cost, g.buyLine == 0 && g.buyFlash > 0)
	drawShop(s, 1, 170, "2  MIXER", g.mixerCost, g.mixers, g.sweets >= g.mixerCost, g.buyLine == 1 && g.buyFlash > 0)
	drawShop(s, 2, 325, "3  BOUTIQUE", g.shopCost, g.shops, g.sweets >= g.shopCost, g.buyLine == 2 && g.buyFlash > 0)

	progress := math.Min(1, g.total/25000)
	vector.DrawFilledRect(s, 15, 633, 450, 14, color.RGBA{4, 16, 30, 220}, true)
	vector.DrawFilledRect(s, 18, 636, float32(444*progress), 8, color.RGBA{99, 231, 242, 235}, true)
	drawCenteredLabel(s, fmt.Sprintf("DISTRICT %d/3   •   TOTAL %s / 25.0K", g.milestone, short(g.total)), 240, 651, bakeryFace14, color.White)
	if g.offline > 0 {
		drawCenteredLabel(s, fmt.Sprintf("WELCOME BACK  •  OFFLINE +%s", short(g.offline)), 240, 674, bakeryFace14, color.RGBA{255, 230, 145, 255})
	} else if g.lastSave > 0 {
		drawCenteredLabel(s, "SAVED IN THIS BROWSER", 240, 674, bakeryFace14, color.RGBA{181, 247, 222, 255})
	} else {
		drawCenteredLabel(s, "R: DELETE SAVE", 240, 674, bakeryFace14, color.RGBA{205, 214, 225, 255})
	}
	if g.frame < 120 {
		alpha := uint8(235)
		if g.frame > 82 {
			alpha = uint8(max(0, 235-(g.frame-82)*6))
		}
		vector.DrawFilledRect(s, 54, 350, 372, 58, color.RGBA{4, 17, 34, alpha}, true)
		vector.StrokeRect(s, 54, 350, 372, 58, 3, color.RGBA{255, 221, 128, alpha}, true)
		drawCenteredLabel(s, "TAP THE TART  •  BUY AN OVEN AT 20", 240, 368, bakeryFace16, color.RGBA{255, 244, 215, alpha})
	}
	if g.milestoneFlash > 0 {
		alpha := uint8(min(240, 90+g.milestoneFlash))
		vector.DrawFilledRect(s, 76, 300, 328, 78, color.RGBA{21, 52, 72, alpha}, true)
		vector.StrokeRect(s, 76, 300, 328, 78, 3, color.RGBA{255, 224, 133, alpha}, true)
		drawCenteredLabel(s, fmt.Sprintf("DISTRICT %d UNLOCKED", g.milestone), 240, 323, bakeryFace20, color.RGBA{255, 240, 193, alpha})
		drawCenteredLabel(s, "THE PEARL PALACE SHINES BRIGHTER", 240, 351, bakeryFace14, color.RGBA{133, 242, 255, alpha})
	}
	if g.clear {
		overlay(s, "BAKERY GOAL REACHED!\n\nTAP / SPACE TO CONTINUE")
	}
}

func drawLine(s *ebiten.Image, line int, x, y float64, owned, frame int, flash bool) {
	alpha := float32(.32)
	if owned > 0 {
		alpha = .95
	}
	bob := 0.0
	if owned > 0 {
		bob = math.Sin(float64(frame)*.11+float64(line)) * 2.5
		vector.DrawFilledCircle(s, float32(x+70), float32(y+42), 42, color.RGBA{80, 225, 244, 25}, true)
	}
	if flash {
		vector.StrokeCircle(s, float32(x+70), float32(y+42), 47, 5, color.RGBA{255, 235, 133, 220}, true)
	}
	drawBakeryContain(s, bakeryGear[line], x+9, y-6+bob, 122, 88, false, alpha)
	vector.DrawFilledRect(s, float32(x+91), float32(y+59), 40, 22, color.RGBA{4, 17, 32, 220}, true)
	drawCenteredLabel(s, fmt.Sprintf("×%d", owned), x+111, y+61, bakeryFace14, color.White)
}

func drawShop(s *ebiten.Image, line int, x float64, label string, cost float64, owned int, ready, flash bool) {
	fill := color.RGBA{8, 27, 45, 236}
	border := color.RGBA{132, 159, 174, 170}
	if ready {
		fill = color.RGBA{11, 68, 78, 242}
		border = color.RGBA{114, 244, 224, 240}
	}
	if flash {
		border = color.RGBA{255, 226, 124, 255}
	}
	vector.DrawFilledRect(s, float32(x), 526, 140, 99, fill, true)
	vector.StrokeRect(s, float32(x), 526, 140, 99, 3, border, true)
	drawBakeryContain(s, bakeryGear[line], x+5, 530, 55, 54, false, 1)
	drawLabel(s, label, x+58, 539, bakeryFace14, color.RGBA{255, 239, 205, 255})
	drawLabel(s, fmt.Sprintf("COST %s", short(cost)), x+58, 563, bakeryFace14, color.White)
	drawLabel(s, fmt.Sprintf("OWNED %d", owned), x+58, 588, bakeryFace14, color.RGBA{154, 236, 244, 255})
}

func loadBakeryArt() {
	bakeryArtOnce.Do(func() {
		bakeryArt = make(map[string]*ebiten.Image, 3)
		for key, filename := range map[string]string{
			"palace":    "bakery-pearl-palace.png",
			"chef":      "bakery-chef.png",
			"equipment": "bakery-equipment-atlas.png",
		} {
			data, err := bakeryArtFS.ReadFile("assets/" + filename)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			bakeryArt[key] = ebiten.NewImageFromImage(decoded)
		}
		atlas := bakeryArt["equipment"]
		for i := range bakeryGear {
			bakeryGear[i] = atlas.SubImage(image.Rect(i*512, 0, (i+1)*512, 512)).(*ebiten.Image)
		}
		bakeryFace14, _ = uilab.Face("en", 14)
		bakeryFace16, _ = uilab.Face("en", 16)
		bakeryFace20, _ = uilab.Face("en", 20)
	})
}

func drawBakeryCover(dst, img *ebiten.Image, x, y, w, h float64) {
	b := img.Bounds()
	scale := math.Max(w/float64(b.Dx()), h/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Min.X), -float64(b.Min.Y))
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(w-float64(b.Dx())*scale)/2, y+(h-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawBakeryContain(dst, img *ebiten.Image, x, y, w, h float64, mirror bool, alpha float32) {
	b := img.Bounds()
	scale := math.Min(w/float64(b.Dx()), h/float64(b.Dy()))
	dw, dh := float64(b.Dx())*scale, float64(b.Dy())*scale
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Min.X), -float64(b.Min.Y))
	if mirror {
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(x+(w+dw)/2, y+(h-dh)/2)
	} else {
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x+(w-dw)/2, y+(h-dh)/2)
	}
	op.ColorScale.ScaleAlpha(alpha)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawLabel(dst *ebiten.Image, label string, x, y float64, face *text.GoTextFace, c color.Color) {
	if face == nil {
		ebitenutil.DebugPrintAt(dst, label, int(x), int(y))
		return
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, label, face, op)
}

func drawCenteredLabel(dst *ebiten.Image, label string, centerX, y float64, face *text.GoTextFace, c color.Color) {
	if face == nil {
		ebitenutil.DebugPrintAt(dst, label, int(centerX)-len(label)*3, int(y))
		return
	}
	w, _ := text.Measure(label, face, 0)
	drawLabel(dst, label, centerX-w/2, y, face, c)
}

func milestoneFor(total float64) int {
	switch {
	case total >= 25000:
		return 3
	case total >= 5000:
		return 2
	case total >= 500:
		return 1
	default:
		return 0
	}
}

func (g *game) nextTarget() (string, float64) {
	name, cost := "OVEN", g.cost
	if g.mixerCost < cost {
		name, cost = "MIXER", g.mixerCost
	}
	if g.shopCost < cost {
		name, cost = "BOUTIQUE", g.shopCost
	}
	return name, cost
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
