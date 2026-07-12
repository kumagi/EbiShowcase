package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
	"github.com/kumagi/EbiShowcase/internal/trackatlas"
	"image/color"
	"strconv"
	"syscall/js"
)

const width, height, tile = 480, 720, 48
const saveKey = "ebiShowcaseQuest."

type game struct {
	x, y, quest, hp, bossHP, scene int
	companion, clear               bool
	message                        string
}

func newGame() *game {
	g := &game{x: 1, y: 10, hp: 60, bossHP: 90, message: "Meet Momo in the southwest village."}
	g.load()
	return g
}
func (g *game) load() {
	s := js.Global().Get("localStorage")
	q := s.Call("getItem", saveKey+"quest")
	if q.Type() != js.TypeString {
		return
	}
	g.quest, _ = strconv.Atoi(q.String())
	g.x, _ = strconv.Atoi(s.Call("getItem", saveKey+"x").String())
	g.y, _ = strconv.Atoi(s.Call("getItem", saveKey+"y").String())
	g.hp, _ = strconv.Atoi(s.Call("getItem", saveKey+"hp").String())
	g.companion = g.quest > 0
	g.setMessage()
}
func (g *game) save() {
	s := js.Global().Get("localStorage")
	s.Call("setItem", saveKey+"quest", strconv.Itoa(g.quest))
	s.Call("setItem", saveKey+"x", strconv.Itoa(g.x))
	s.Call("setItem", saveKey+"y", strconv.Itoa(g.y))
	s.Call("setItem", saveKey+"hp", strconv.Itoa(g.hp))
}
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		s := js.Global().Get("localStorage")
		for _, k := range []string{"quest", "x", "y", "hp"} {
			s.Call("removeItem", saveKey+k)
		}
		*g = *newGame()
		return nil
	}
	if g.clear {
		if any() {
			g.resetSave()
			*g = *newGame()
		}
		return nil
	}
	if g.scene == 1 {
		return g.battle()
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
		g.x = max(0, min(9, g.x+dx))
		g.y = max(0, min(11, g.y+dy))
		switch {
		case g.quest == 0 && g.x == 2 && g.y == 9:
			g.quest = 1
			g.companion = true
			g.message = "Momo joined! Find the crystal in the northeast."
		case g.quest == 1 && g.x == 8 && g.y == 2:
			g.quest = 2
			g.message = "Crystal found! Enter the dark tower."
		case g.quest == 2 && g.x == 8 && g.y == 1:
			g.scene = 1
			g.bossHP = 90
			g.message = "Shadow Crab blocks the way!"
		case g.quest == 3 && g.x == 1 && g.y == 10:
			g.quest = 4
			g.clear = true
			g.message = "The village is safe! Quest complete."
		}
		g.save()
	}
	return nil
}
func (g *game) resetSave() {
	s := js.Global().Get("localStorage")
	for _, k := range []string{"quest", "x", "y", "hp"} {
		s.Call("removeItem", saveKey+k)
	}
}
func (g *game) battle() error {
	choice := -1
	if inpututil.IsKeyJustPressed(ebiten.Key1) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		choice = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		choice = 1
	}
	if x, y, ok := press(); ok && y > 530 {
		choice = min(1, x/(width/2))
	}
	if choice < 0 {
		return nil
	}
	if choice == 0 {
		d := 12
		if g.companion {
			d += 7
		}
		g.bossHP -= d
		g.message = fmt.Sprintf("Ebi and Momo deal %d damage!", d)
	} else {
		g.hp = min(60, g.hp+15)
		g.message = "Momo's song restores 15 HP!"
	}
	if g.bossHP <= 0 {
		g.quest = 3
		g.scene = 0
		g.x, g.y = 8, 2
		g.message = "Boss defeated! Return to the southwest village."
		g.save()
		return nil
	}
	g.hp -= 11
	if g.hp <= 0 {
		g.hp = 60
		g.scene = 0
		g.x, g.y = 2, 9
		g.message = "The party escaped and recovered in the village."
	}
	g.save()
	return nil
}
func (g *game) setMessage() {
	switch g.quest {
	case 0:
		g.message = "Meet Momo in the southwest village."
	case 1:
		g.message = "Find the crystal in the northeast."
	case 2:
		g.message = "Enter the dark tower beside the crystal."
	case 3:
		g.message = "Return to the southwest village."
	case 4:
		g.message = "Quest complete."
	}
}
func (g *game) Draw(s *ebiten.Image) {
	if g.scene == 1 {
		g.drawBattle(s)
		return
	}
	s.Fill(color.RGBA{12, 28, 40, 255})
	oy := 74
	for y := 0; y < 12; y++ {
		for x := 0; x < 10; x++ {
			t := "tile-grass"
			if x > 5 {
				t = "tile-cobble"
			}
			trackatlas.Draw(s, t, float64(x*tile), float64(oy+y*tile), tile)
		}
	}
	trackatlas.DrawCentered(s, "npc", 2*tile+24, float64(oy+9*tile+24), 32)
	if g.quest == 1 {
		vector.DrawFilledCircle(s, 8*tile+24, float32(oy+2*tile+24), 11, color.RGBA{72, 205, 255, 255}, false)
	}
	vector.DrawFilledRect(s, 8*tile+7, float32(oy+tile+7), 34, 34, color.RGBA{55, 45, 75, 255}, false)
	hero.DrawCentered(s, float64(g.x*tile+24), float64(oy+g.y*tile+24), 34)
	if g.companion {
		trackatlas.DrawCentered(s, "ally", float64(g.x*tile+10), float64(oy+g.y*tile+38), 22)
	}
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("EBI QUEST  HP %02d/60  QUEST %d/4", g.hp, g.quest), 120, 25)
	ebitenutil.DebugPrintAt(s, g.message, 60, 50)
	ebitenutil.DebugPrintAt(s, "AUTOSAVED IN THIS BROWSER   R: DELETE SAVE", 85, 690)
	if g.clear {
		overlay(s, "EBI QUEST COMPLETE!\n\nTAP / SPACE TO PLAY AGAIN")
	}
}
func (g *game) drawBattle(s *ebiten.Image) {
	s.Fill(color.RGBA{18, 24, 43, 255})
	trackatlas.DrawCentered(s, "king-crab", 240, 210, 156)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("SHADOW CRAB HP %02d/90", max(0, g.bossHP)), 160, 105)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("PARTY HP %02d/60", g.hp), 185, 415)
	ebitenutil.DebugPrintAt(s, g.message, 60, 480)
	vector.DrawFilledRect(s, 30, 550, 200, 80, color.RGBA{45, 225, 194, 255}, false)
	vector.DrawFilledRect(s, 250, 550, 200, 80, color.RGBA{255, 190, 75, 255}, false)
	ebitenutil.DebugPrintAt(s, "1 TEAM ATTACK", 80, 585)
	ebitenutil.DebugPrintAt(s, "2 MOMO'S SONG", 295, 585)
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
func overlay(s *ebiten.Image, msg string) {
	vector.DrawFilledRect(s, 55, 280, 370, 150, color.RGBA{6, 18, 37, 245}, false)
	ebitenutil.DebugPrintAt(s, msg, 125, 330)
}
func (g *game) Layout(_, _ int) (int, int) { return width, height }
func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Ebi Quest — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
