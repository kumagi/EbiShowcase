package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

const W, H = 480, 720

type unit struct {
	name         string
	gauge, speed int
	ready        bool
}
type game struct {
	units []unit
	queue []int
	acted int
}

func (g *game) Update() error {
	for i := range g.units {
		u := &g.units[i]
		if u.ready {
			continue
		}
		u.gauge += u.speed
		if u.gauge >= 1000 {
			u.gauge = 1000
			u.ready = true
			g.queue = append(g.queue, i)
		}
	}
	pressed := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
	if pressed && len(g.queue) > 0 {
		i := g.queue[0]
		g.queue = g.queue[1:]
		g.units[i].gauge = 0
		g.units[i].ready = false
		g.acted++
	}
	return nil
}
func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 25, 45, 255})
	ebitenutil.DebugPrintAt(s, "READY QUEUE", 190, 55)
	for i, u := range g.units {
		y := 150 + i*105
		vector.DrawFilledRect(s, 65, float32(y), 350, 70, color.RGBA{35, 55, 78, 255}, false)
		vector.DrawFilledRect(s, 75, float32(y+42), float32(330*u.gauge/1000), 14, color.RGBA{245, 185, 65, 255}, false)
		state := "WAIT"
		if u.ready {
			state = "READY"
		}
		ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s  SPD %d  %s", u.name, u.speed, state), 85, y+15)
	}
	q := "EMPTY"
	if len(g.queue) > 0 {
		q = ""
		for _, i := range g.queue {
			q += g.units[i].name + " > "
		}
	}
	ebitenutil.DebugPrintAt(s, "QUEUE: "+q, 80, 520)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("ACTIONS %d", g.acted), 190, 565)
	ebitenutil.DebugPrintAt(s, "TAP / SPACE: resolve first READY", 120, 655)
}
func (g *game) Layout(int, int) (int, int) { return W, H }
func main() {
	ebiten.SetWindowSize(W, H)
	if err := ebiten.RunGame(&game{units: []unit{{"TENJIROH", 850, 8, false}, {"MAGE", 780, 13, false}, {"SHELL", 920, 5, false}}}); err != nil {
		panic(err)
	}
}
