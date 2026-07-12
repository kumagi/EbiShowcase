package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"image/color"
	"math/rand"
)

const width, height, tile = 480, 720, 48
const (
	mapScene = iota
	battleScene
)

type enemy struct {
	name       string
	hp, attack int
	c          color.RGBA
}
type game struct {
	x, y, scene, hp, wins, steps int
	foe                          enemy
	message                      string
	clear, over                  bool
	rng                          *rand.Rand
}

var tables = [][]enemy{{{"Leaf Slime", 18, 4, color.RGBA{80, 205, 105, 255}}, {"Mush Bug", 22, 5, color.RGBA{187, 90, 185, 255}}}, {{"Sand Crab", 25, 6, color.RGBA{225, 154, 75, 255}}, {"Dust Bat", 20, 7, color.RGBA{133, 91, 184, 255}}}, {{"Snow Eel", 28, 7, color.RGBA{96, 190, 229, 255}}, {"Ice Slime", 32, 6, color.RGBA{145, 215, 245, 255}}}}

func newGame() *game {
	return &game{x: 1, y: 10, hp: 50, rng: rand.New(rand.NewSource(3406)), message: "Reach the castle after three victories."}
}
func (g *game) Update() error {
	if g.clear || g.over {
		if any() {
			*g = *newGame()
		}
		return nil
	}
	if g.scene == battleScene {
		return g.updateBattle()
	}
	dx, dy := 0, 0
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		dx = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		dx = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		dy = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		dy = 1
	}
	if x, y, ok := press(); ok {
		rx, ry := x-(g.x*tile+24), y-(74+g.y*tile+24)
		if abs(rx) > abs(ry) {
			if rx < 0 {
				dx = -1
			} else {
				dx = 1
			}
		} else {
			if ry < 0 {
				dy = -1
			} else {
				dy = 1
			}
		}
	}
	if dx != 0 || dy != 0 {
		g.x = clampInt(g.x+dx, 0, 9)
		g.y = clampInt(g.y+dy, 0, 11)
		g.steps++
		if g.x == 8 && g.y == 1 && g.wins >= 3 {
			g.clear = true
			return nil
		}
		if g.steps%5 == 0 && g.rng.Intn(100) < 65 {
			biome := g.biome()
			list := tables[biome]
			g.foe = list[g.rng.Intn(len(list))]
			g.scene = battleScene
			g.message = "A wild " + g.foe.name + " appeared!"
		}
	}
	return nil
}
func (g *game) updateBattle() error {
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if x, y, ok := press(); ok && y > 540 {
		choice = min(1, x/(width/2))
	}
	if choice < 0 {
		return nil
	}
	if choice == 1 {
		if g.rng.Intn(100) < 60 {
			g.scene = mapScene
			g.message = "Escaped back to the map."
			return nil
		}
		g.message = "Could not escape!"
	} else {
		d := 8 + g.rng.Intn(7)
		g.foe.hp -= d
		g.message = fmt.Sprintf("Ebi attacks for %d!", d)
		if g.foe.hp <= 0 {
			g.wins++
			g.scene = mapScene
			g.message = "Victory! Continue toward the castle."
			return nil
		}
	}
	g.hp -= g.foe.attack
	if g.hp <= 0 {
		g.over = true
	}
	return nil
}
func (g *game) biome() int {
	if g.x < 4 {
		return 0
	}
	if g.x < 7 {
		return 1
	}
	return 2
}
func (g *game) Draw(s *ebiten.Image) {
	if g.scene == battleScene {
		g.drawBattle(s)
		return
	}
	s.Fill(color.RGBA{12, 28, 40, 255})
	oy := 74
	for y := 0; y < 12; y++ {
		for x := 0; x < 10; x++ {
			c := color.RGBA{87, 169, 91, 255}
			if x >= 4 {
				c = color.RGBA{211, 170, 83, 255}
			}
			if x >= 7 {
				c = color.RGBA{177, 220, 235, 255}
			}
			vector.DrawFilledRect(s, float32(x*tile), float32(oy+y*tile), tile, tile, c, false)
		}
	}
	gate := color.RGBA{130, 73, 82, 255}
	if g.wins >= 3 {
		gate = color.RGBA{45, 225, 194, 255}
	}
	vector.DrawFilledRect(s, 8*tile+6, float32(oy+tile+6), 36, 36, gate, false)
	hero.DrawCentered(s, float64(g.x*tile+24), float64(oy+g.y*tile+24), 34)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("WORLD MAP  HP %02d/50  VICTORIES %d/3", g.hp, g.wins), 95, 25)
	ebitenutil.DebugPrintAt(s, g.message, 70, 50)
	ebitenutil.DebugPrintAt(s, "GREEN / DESERT / SNOW HAVE DIFFERENT ENEMIES", 80, 690)
	if g.clear {
		overlay(s, "CASTLE REACHED!\n\nTAP / SPACE TO TRAVEL AGAIN")
	}
}
func (g *game) drawBattle(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 28, 48, 255})
	vector.DrawFilledCircle(s, 240, 230, 70, g.foe.c, false)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("%s  HP %02d", g.foe.name, max(0, g.foe.hp)), 175, 120)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI HP %02d/50", max(0, g.hp)), 190, 400)
	ebitenutil.DebugPrintAt(s, g.message, 60, 480)
	vector.DrawFilledRect(s, 35, 550, 195, 80, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 250, 550, 195, 80, color.RGBA{224, 159, 72, 255}, false)
	ebitenutil.DebugPrintAt(s, "1 ATTACK", 95, 585)
	ebitenutil.DebugPrintAt(s, "2 RUN", 320, 585)
	if g.over {
		overlay(s, "PARTY DEFEATED!\n\nTAP / SPACE TO RESTART")
	}
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
func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func clampInt(v, l, h int) int { return max(l, min(h, v)) }
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("World Encounters — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
