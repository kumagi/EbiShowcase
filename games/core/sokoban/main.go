package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/hero"
)

const (
	screenWidth  = 480
	screenHeight = 720
	tileSize     = 48
)

type point struct {
	x, y int
}

// 1手戻すための記録
type snapshot struct {
	player point
	boxes  []point
}

type game struct {
	player  point
	boxes   []point
	goals   []point
	walls   map[point]bool
	history []snapshot
	moves   int
	cleared bool
}

// # = 壁, @ = プレイヤー, $ = 箱, . = ゴール
var level = []string{
	"##########",
	"#        #",
	"#  . .   #",
	"#  $ $   #",
	"#   ##   #",
	"#   @    #",
	"#        #",
	"#        #",
	"##########",
}

func newGame() *game {
	g := &game{walls: map[point]bool{}}
	for y, row := range level {
		for x, c := range row {
			p := point{x, y}
			switch c {
			case '#':
				g.walls[p] = true
			case '@':
				g.player = p
			case '$':
				g.boxes = append(g.boxes, p)
			case '.':
				g.goals = append(g.goals, p)
			}
		}
	}
	return g
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func justPressedPosition() (x, y int, ok bool) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y = ebiten.CursorPosition()
		return x, y, true
	}
	ids := inpututil.AppendJustPressedTouchIDs(nil)
	if len(ids) > 0 {
		x, y = ebiten.TouchPosition(ids[0])
		return x, y, true
	}
	return 0, 0, false
}

func restartPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	if len(inpututil.AppendJustPressedTouchIDs(nil)) > 0 {
		return true
	}
	return false
}

func (g *game) boxAt(p point) int {
	for i, b := range g.boxes {
		if b == p {
			return i
		}
	}
	return -1
}

func (g *game) isGoal(p point) bool {
	for _, goal := range g.goals {
		if goal == p {
			return true
		}
	}
	return false
}

func (g *game) save() {
	boxes := append([]point(nil), g.boxes...)
	g.history = append(g.history, snapshot{player: g.player, boxes: boxes})
}

func (g *game) undo() {
	if len(g.history) == 0 || g.cleared {
		return
	}
	last := g.history[len(g.history)-1]
	g.history = g.history[:len(g.history)-1]
	g.player = last.player
	g.boxes = append([]point(nil), last.boxes...)
	if g.moves > 0 {
		g.moves--
	}
}

// 1マス動く（箱があれば押す）
func (g *game) move(d point) {
	next := point{g.player.x + d.x, g.player.y + d.y}
	if g.walls[next] {
		return
	}

	boxIndex := g.boxAt(next)
	if boxIndex >= 0 {
		beyond := point{next.x + d.x, next.y + d.y}
		if g.walls[beyond] || g.boxAt(beyond) >= 0 {
			return // 箱の先が壁か別の箱なら押せない
		}
		g.save()
		g.boxes[boxIndex] = beyond
		g.player = next
		g.moves++
	} else {
		g.save()
		g.player = next
		g.moves++
	}

	// すべてのゴールに箱があるか？
	g.cleared = true
	for _, goal := range g.goals {
		if g.boxAt(goal) < 0 {
			g.cleared = false
			break
		}
	}
}

func (g *game) readMove() point {
	d := point{}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		d.y = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		d.y = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		d.x = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		d.x = 1
	}

	// 画面中央からのタッチ方向
	if x, y, ok := justPressedPosition(); ok {
		cx, cy := screenWidth/2, screenHeight/2
		dx, dy := x-cx, y-cy
		if abs(dx) > abs(dy) {
			if dx < 0 {
				d.x = -1
			} else {
				d.x = 1
			}
		} else {
			if dy < 0 {
				d.y = -1
			} else {
				d.y = 1
			}
		}
	}
	return d
}

// --- ここから Update ---
func (g *game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *newGame()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) || inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.undo()
		return nil
	}
	if g.cleared {
		if restartPressed() {
			*g = *newGame()
		}
		return nil
	}

	d := g.readMove()
	if d.x != 0 || d.y != 0 {
		g.move(d)
	}
	return nil
}

// --- ここから Draw ---
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{7, 19, 35, 255})
	offsetY := 105

	for y := range level {
		for x := range level[y] {
			p := point{x, y}
			px := float32(x * tileSize)
			py := float32(offsetY + y*tileSize)

			floor := color.RGBA{24, 47, 65, 255}
			if (x+y)%2 == 0 {
				floor = color.RGBA{27, 53, 72, 255}
			}
			vector.DrawFilledRect(screen, px, py, tileSize, tileSize, floor, false)

			if g.walls[p] {
				vector.DrawFilledRect(screen, px+2, py+2, tileSize-4, tileSize-4, color.RGBA{74, 94, 119, 255}, false)
				vector.StrokeRect(screen, px+5, py+5, tileSize-10, tileSize-10, 2, color.RGBA{117, 146, 174, 255}, false)
			}
			if g.isGoal(p) {
				vector.DrawFilledCircle(screen, px+tileSize/2, py+tileSize/2, 10, color.RGBA{255, 211, 76, 255}, false)
				vector.StrokeCircle(screen, px+tileSize/2, py+tileSize/2, 16, 2, color.RGBA{255, 211, 76, 170}, false)
			}
		}
	}

	for _, b := range g.boxes {
		px := float32(b.x * tileSize)
		py := float32(offsetY + b.y*tileSize)
		c := color.RGBA{224, 130, 64, 255}
		if g.isGoal(b) {
			c = color.RGBA{50, 210, 151, 255}
		}
		vector.DrawFilledRect(screen, px+7, py+7, tileSize-14, tileSize-14, c, false)
		vector.StrokeRect(screen, px+12, py+12, tileSize-24, tileSize-24, 3, color.RGBA{255, 234, 190, 210}, false)
	}

	px := float32(g.player.x*tileSize + tileSize/2)
	py := float32(offsetY + g.player.y*tileSize + tileSize/2)
	hero.DrawCentered(screen, float64(px), float64(py), float64(tileSize)-4)

	ebitenutil.DebugPrintAt(screen, "SOKOBAN — PUT EVERY BOX ON A GOLD GOAL", 86, 28)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MOVES %03d", g.moves), 20, 72)
	ebitenutil.DebugPrintAt(screen, "UNDO: Z / BACKSPACE    RESET: R", 205, 72)
	if g.cleared {
		vector.DrawFilledRect(screen, 55, 286, 370, 140, color.RGBA{4, 16, 31, 240}, false)
		ebitenutil.DebugPrintAt(screen, "STAGE CLEAR!\n\nTAP / SPACE TO PLAY AGAIN", 145, 330)
	} else {
		ebitenutil.DebugPrintAt(screen, "ARROWS / WASD / TAP A DIRECTION", 126, 675)
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sokoban — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
