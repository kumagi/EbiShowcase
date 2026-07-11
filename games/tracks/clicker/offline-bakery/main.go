package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"strconv"
	"syscall/js"
	"time"
)

const width, height = 480, 720
const prefix = "ebiShowcaseBakery."

type game struct {
	sweets         float64
	ovens, frame   int
	cost, lastSave float64
	last           time.Time
	offline        float64
	clear          bool
}

func newGame() *game          { g := &game{cost: 20, last: time.Now()}; g.load(); return g }
func (g *game) rate() float64 { return float64(g.ovens) * 5 }
func (g *game) load() {
	store := js.Global().Get("localStorage")
	raw := store.Call("getItem", prefix+"sweets")
	if raw.Type() != js.TypeString {
		return
	}
	g.sweets, _ = strconv.ParseFloat(raw.String(), 64)
	g.ovens, _ = strconv.Atoi(store.Call("getItem", prefix+"ovens").String())
	g.cost = 20 * math.Pow(1.25, float64(g.ovens))
	stamp, _ := strconv.ParseInt(store.Call("getItem", prefix+"time").String(), 10, 64)
	away := math.Min(8*3600, float64(time.Now().Unix()-stamp))
	if away > 0 {
		g.offline = away * g.rate()
		g.sweets += g.offline
	}
}
func (g *game) save() {
	s := js.Global().Get("localStorage")
	s.Call("setItem", prefix+"sweets", fmt.Sprintf("%.3f", g.sweets))
	s.Call("setItem", prefix+"ovens", strconv.Itoa(g.ovens))
	s.Call("setItem", prefix+"time", strconv.FormatInt(time.Now().Unix(), 10))
	g.lastSave = 2
}
func (g *game) Update() error {
	now := time.Now()
	dt := math.Min(.25, now.Sub(g.last).Seconds())
	g.last = now
	g.frame++
	if !g.clear {
		g.sweets += g.rate() * dt
	}
	if g.lastSave > 0 {
		g.lastSave -= dt
	}
	if g.frame%120 == 0 {
		g.save()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		s := js.Global().Get("localStorage")
		s.Call("removeItem", prefix+"sweets")
		s.Call("removeItem", prefix+"ovens")
		s.Call("removeItem", prefix+"time")
		*g = *newGame()
		return nil
	}
	if g.clear {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	_, y, ok := press()
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		ok = true
		y = 300
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		ok = true
		y = 570
	}
	if ok {
		g.offline = 0
		if y < 470 {
			g.sweets += 1 + float64(g.ovens)
		} else if g.sweets >= g.cost {
			g.sweets -= g.cost
			g.ovens++
			g.cost = 20 * math.Pow(1.25, float64(g.ovens))
		}
		g.save()
	}
	if g.sweets >= 5000 {
		g.clear = true
		g.save()
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{39, 25, 43, 255})
	ebitenutil.DebugPrintAt(s, "EBI BAKERY — SAVED IN THIS BROWSER", 110, 40)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SWEETS %s / 5.00K", short(g.sweets)), 160, 85)
	vector.DrawFilledCircle(s, 240, 280, 100, color.RGBA{222, 143, 76, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: BAKE", 165, 405)
	vector.DrawFilledRect(s, 55, 475, 370, 145, color.RGBA{45, 205, 181, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BUY AUTO OVEN [B]\n\nCOST %s   OWNED %d\nPRODUCTION %s / SECOND", short(g.cost), g.ovens, short(g.rate())), 130, 500)
	if g.offline > 0 {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("WELCOME BACK! OFFLINE +%s", short(g.offline)), 120, 650)
	} else if g.lastSave > 0 {
		ebitenutil.DebugPrintAt(s, "SAVED", 220, 650)
	}
	ebitenutil.DebugPrintAt(s, "R: DELETE SAVE", 185, 680)
	if g.clear {
		overlay(s, "BAKERY GOAL REACHED!\n\nTAP / SPACE TO CONTINUE")
	}
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
