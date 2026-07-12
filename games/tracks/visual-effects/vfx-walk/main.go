// vfx-walk — Visual Effects Lab STEP 06.
// Frame animation from a real texture atlas: pick an action and a facing, then
// flip frames with SubImage on a timer. Uses the downloadable Ebi Tenjiroh atlas.
package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/heroatlas"
)

const width, height = 480, 720

var actions = []string{"idle", "walk", "run", "attack", "hurt"}
var facings = []string{"down", "up", "side", "side"} // last = side flipped (left)

type game struct {
	actionIdx int
	faceIdx   int
	frame     int
	tick      int
	playing   bool
	x         float64
	seen      map[string]bool
	buttons   []button
	clear     bool
}

type button struct {
	x, y, w, h float64
	label      string
}

func (b button) hit(px, py float64) bool {
	return px >= b.x && px <= b.x+b.w && py >= b.y && py <= b.y+b.h
}

func newGame() *game {
	g := &game{actionIdx: 1, playing: true, x: width / 2, seen: map[string]bool{"walk": true}}
	w, gap := 138.0, 12.0
	x := (width - (w*3 + gap*2)) / 2
	for _, l := range []string{"ACTION", "FACING", "PLAY/STOP"} {
		g.buttons = append(g.buttons, button{x, 636, w, 54, l})
		x += w + gap
	}
	return g
}

func (g *game) animName() string {
	return actions[g.actionIdx] + "-" + facings[g.faceIdx]
}

func (g *game) flipped() bool { return g.faceIdx == 3 }

func pressXY() (float64, float64, bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return float64(x), float64(y), true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		return float64(x), float64(y), true
	}
	return 0, 0, false
}

func anyStart() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Update() error {
	if g.clear {
		if anyStart() {
			*g = *newGame()
		}
		return nil
	}
	px, py, tapped := pressXY()
	if tapped {
		switch {
		case g.buttons[0].hit(px, py):
			g.actionIdx = (g.actionIdx + 1) % len(actions)
			g.frame, g.tick = 0, 0
		case g.buttons[1].hit(px, py):
			g.faceIdx = (g.faceIdx + 1) % len(facings)
			g.frame, g.tick = 0, 0
		case g.buttons[2].hit(px, py):
			g.playing = !g.playing
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.actionIdx = (g.actionIdx + 1) % len(actions)
		g.frame, g.tick = 0, 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.faceIdx = (g.faceIdx + 1) % len(facings)
		g.frame, g.tick = 0, 0
	}
	g.seen[actions[g.actionIdx]] = true
	if len(g.seen) >= len(actions) {
		g.clear = true
	}

	frames := heroatlas.Anim(g.animName())
	if len(frames) == 0 {
		return nil
	}
	hold := 60 / heroatlas.FPS(g.animName())
	if hold < 1 {
		hold = 1
	}
	stepNow := false
	if g.playing {
		g.tick++
		if g.tick >= hold {
			g.tick = 0
			g.frame = (g.frame + 1) % len(frames)
			stepNow = true
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.frame = (g.frame + 1) % len(frames)
		stepNow = true
	}
	// Only the "side" facings travel across the ground.
	moving := actions[g.actionIdx] == "walk" || actions[g.actionIdx] == "run"
	if g.playing && moving && (g.faceIdx == 2 || g.faceIdx == 3) {
		sp := 1.6
		if actions[g.actionIdx] == "run" {
			sp = 3.4
		}
		if g.flipped() {
			g.x -= sp
		} else {
			g.x += sp
		}
		if g.x < 60 {
			g.x = width - 60
		}
		if g.x > width-60 {
			g.x = 60
		}
	} else if !moving {
		g.x = width / 2
	}
	_ = stepNow
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 26, 44, 255})
	vector.DrawFilledRect(s, 0, 470, width, 130, color.RGBA{26, 38, 60, 255}, false)

	frames := heroatlas.Anim(g.animName())
	if len(frames) == 0 {
		return
	}

	// Current strip: the frames of this animation, active one highlighted.
	cell := 64.0
	startX := (width - cell*float64(len(frames))) / 2
	stripY := 92.0
	for i, fr := range frames {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(cell/heroatlas.FrameW, cell/heroatlas.FrameH)
		op.GeoM.Translate(startX+float64(i)*cell, stripY)
		s.DrawImage(fr, op)
		c := color.RGBA{80, 100, 140, 255}
		if i == g.frame {
			c = color.RGBA{120, 240, 220, 255}
		}
		vector.StrokeRect(s, float32(startX+float64(i)*cell), float32(stripY), float32(cell), float32(cell), 3, c, false)
	}

	// Big animated character.
	scale := 2.6
	fr := frames[g.frame]
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	if g.flipped() {
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(g.x+heroatlas.FrameW*scale/2, 300-heroatlas.FrameH*scale/2)
	} else {
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(g.x-heroatlas.FrameW*scale/2, 300-heroatlas.FrameH*scale/2)
	}
	s.DrawImage(fr, op)

	face := facings[g.faceIdx]
	if g.flipped() {
		face = "side (flipped = left)"
	}
	ebitenutil.DebugPrintAt(s, "ATLAS FRAMES (SubImage) - ACTIVE ONE GLOWS", 92, 24)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("anim %s   frame %d/%d   fps %d", g.animName(), g.frame+1, len(frames), heroatlas.FPS(g.animName())), 16, 52)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("facing: %s   (tried %d/5 actions)", face, len(g.seen)), 16, 616)
	for _, b := range g.buttons {
		vector.DrawFilledRect(s, float32(b.x), float32(b.y), float32(b.w), float32(b.h), color.RGBA{40, 54, 82, 235}, false)
		vector.StrokeRect(s, float32(b.x), float32(b.y), float32(b.w), float32(b.h), 3, color.RGBA{96, 116, 158, 255}, false)
		ebitenutil.DebugPrintAt(s, b.label, int(b.x+b.w/2-float64(len(b.label))*3), int(b.y+b.h/2-8))
	}
	if g.clear {
		overlay(s, "YOU TRIED EVERY ACTION!\n\nONE ATLAS + SubImage + TIMER\n= WALK, RUN, ATTACK, HURT.\nTAP / SPACE TO RESET")
	}
}

func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 45, 250, 390, 170, color.RGBA{8, 16, 32, 245}, false)
	vector.StrokeRect(s, 45, 250, 390, 170, 3, color.RGBA{120, 240, 220, 255}, false)
	ebitenutil.DebugPrintAt(s, msg, 85, 285)
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Sprite Atlas — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
