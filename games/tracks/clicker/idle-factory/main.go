package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"github.com/kumagi/EbiShowcase/internal/uilab"
	"image/color"
	"time"
)

const width, height = 480, 720

type game struct {
	sweets   float64
	machines int
	cost     float64
	last     time.Time
	clear    bool
	audio    *audio.Context
	gate     audiolab.Gate
	pulse    *shaderlab.Pulse
	cam      cameralab.State
	badge    *ebiten.Image
}

func newGame() *game {
	b := ebiten.NewImage(24, 24)
	b.Fill(color.RGBA{255, 210, 90, 255})
	return &game{cost: 15, last: time.Now(), audio: audiolab.Context(), pulse: shaderlab.NewPulse(), cam: cameralab.State{Pos: cameralab.Vec{240, 360}, ViewW: width, ViewH: height}, badge: b}
}
func (g *game) Update() error {
	now := time.Now()
	dt := now.Sub(g.last).Seconds()
	g.last = now
	if dt > .25 {
		dt = .25
	}
	if !g.clear {
		g.sweets += float64(g.machines) * dt
	}
	if g.clear {
		if anyPress() {
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
		if y < 470 {
			g.sweets++
			g.play(600)
		} else if g.sweets >= g.cost {
			g.sweets -= g.cost
			g.machines++
			g.cost += 10
			g.play(360)
		}
	}
	if g.sweets >= 200 {
		g.clear = true
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{26, 30, 48, 255})
	if g.pulse.Available() {
		fx := ebiten.NewImage(24, 24)
		if g.pulse.Draw(fx, g.badge, float32(g.sweets)*.06) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(430, 16)
			s.DrawImage(fx, op)
		}
	}
	if f, e := uilab.Face("en", 18); e == nil {
		op := &text.DrawOptions{}
		op.GeoM.Translate(150, 55)
		text.Draw(s, fmt.Sprintf("SWEETS %06.1f / 200", g.sweets), f, op)
	} else {
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("SWEETS %06.1f / 200", g.sweets), 150, 55)
	}
	vector.DrawFilledCircle(s, 240, 275, 100, color.RGBA{222, 143, 76, 255}, false)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: BAKE +1", 155, 405)
	vector.DrawFilledRect(s, 60, 480, 360, 145, color.RGBA{45, 196, 174, 255}, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("BUY AUTO OVEN [B]\n\nCOST %.0f     OWNED %d\nPRODUCTION %.1f SWEETS / SECOND", g.cost, g.machines, float64(g.machines)), 135, 505)
	if g.clear {
		overlay(s, "200 SWEETS PRODUCED!\n\nTAP / SPACE TO RESTART")
	}
}
func (g *game) play(hz float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Sine, hz, .08)).Play()
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
func anyPress() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 240}, false)
	ebitenutil.DebugPrintAt(s, msg, 130, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Idle Factory — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
