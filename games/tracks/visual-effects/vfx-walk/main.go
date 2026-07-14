// vfx-walk — STEP 06: SubImage frame animation via live Go + mouse.
package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/internal/heroatlas"
	"github.com/kumagi/EbiShowcase/internal/vfxlive"
)

var actions = []string{"idle", "walk", "run", "attack", "hurt"}
var facings = []string{"down", "up", "side"}

type game struct {
	shell *vfxlive.Shell
	frame int
	hold  float64
	x     float64
}

func newGame() *game {
	return &game{
		x: 240,
		shell: vfxlive.New(
			"SubImage frames",
			[]string{
				"name := \"{anim}\"",
				"frames := atlas.Anim(name) // SubImage strips",
				"hold := 60 / {fps}",
				"screen.DrawImage(frames[{frame}], op)",
				"op.GeoM.Scale({flip}, 1) // side facing flip",
			},
			&vfxlive.Param{Key: "action", Label: "action", Value: 1, Min: 0, Max: 4, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "facing", Label: "facing", Value: 2, Min: 0, Max: 2, Step: 1, Format: "%.0f"},
			&vfxlive.Param{Key: "play", Label: "playing", Value: 1, Bool: true},
			&vfxlive.Param{Key: "speed", Label: "speed", Value: 1, Min: 0.25, Max: 2.5, Format: "%.2f"},
		),
	}
}

func (g *game) animName() string {
	ai := int(g.shell.Get("action") + 0.5)
	fi := int(g.shell.Get("facing") + 0.5)
	if ai < 0 {
		ai = 0
	}
	if ai >= len(actions) {
		ai = len(actions) - 1
	}
	if fi < 0 {
		fi = 0
	}
	if fi >= len(facings) {
		fi = len(facings) - 1
	}
	return actions[ai] + "-" + facings[fi]
}

func (g *game) Update() error {
	g.shell.Update()
	name := g.animName()
	frames := heroatlas.Anim(name)
	fps := float64(heroatlas.FPS(name)) * g.shell.Get("speed")
	if fps < 1 {
		fps = 1
	}
	if g.shell.Bool("play") && len(frames) > 0 {
		g.hold++
		if g.hold >= 60/fps {
			g.hold = 0
			g.frame = (g.frame + 1) % len(frames)
		}
	}
	if len(frames) == 0 || g.frame >= len(frames) {
		g.frame = 0
	}
	if actions[int(g.shell.Get("action")+0.5)] == "walk" || actions[int(g.shell.Get("action")+0.5)] == "run" {
		sp := 1.2 * g.shell.Get("speed")
		if actions[int(g.shell.Get("action")+0.5)] == "run" {
			sp *= 1.8
		}
		g.x += sp
		if g.x > 420 {
			g.x = 60
		}
	}
	g.shell.SetToken("anim", name)
	g.shell.SetToken("fps", fmt.Sprintf("%d", int(float64(heroatlas.FPS(name))*g.shell.Get("speed"))))
	g.shell.SetToken("frame", fmt.Sprintf("%d", g.frame))
	g.shell.SetToken("flip", "1")
	g.shell.Hint = "action 0-4 · facing 0-2 · toggle playing"
	return nil
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{10, 14, 28, 255})
	g.shell.FillStage(s, color.RGBA{14, 18, 36, 255})
	_, sy, _, sh := g.shell.Stage()
	cy := sy + sh/2 + 30

	name := g.animName()
	frames := heroatlas.Anim(name)
	flip := 1.0
	if int(g.shell.Get("facing")+0.5) == 2 {
		// demonstrate flip token; learner can imagine left vs right
		flip = 1
	}
	if len(frames) > 0 {
		fr := frames[g.frame]
		b := fr.Bounds()
		sc := 2.4
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
		op.GeoM.Scale(sc*flip, sc)
		op.GeoM.Translate(g.x, cy)
		op.Filter = ebiten.FilterNearest
		s.DrawImage(fr, op)
	}

	g.shell.Draw(s)
}

func (g *game) Layout(_, _ int) (int, int) { return vfxlive.Width, vfxlive.Height }

func main() {
	ebiten.SetWindowSize(vfxlive.Width, vfxlive.Height)
	ebiten.SetWindowTitle("Live Go: Walk frames — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
