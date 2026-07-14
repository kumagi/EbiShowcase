// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package raycasterui

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/raycastlogic"
	"golang.org/x/image/font/basicfont"
)

const (
	W   = 720
	H   = 480
	FOV = math.Pi / 3
)

type Variant int

const (
	FacingMove Variant = iota
	SingleRay
	DistanceStrip
	ColumnView
	TexturedView
	EbiRaycaster
)

var world = [][]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 1, 1, 0, 1, 1, 1, 0, 1, 0, 1},
	{1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1},
	{1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1},
	{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 1},
	{1, 1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1},
	{1, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

type actor struct {
	x, y  float64
	kind  int // 0 enemy, 1 key, 2 exit
	alive bool
}

type Game struct {
	variant Variant
	x, y    float64
	angle   float64
	frame   int
	shots   int
	hasKey  bool
	win     bool
	message string
	actors  []actor
	zbuf    []float64
}

func New(variant Variant) *Game {
	g := &Game{variant: variant, zbuf: make([]float64, W), message: "TURN, MOVE, AND READ THE RAYS"}
	g.reset()
	return g
}

func (g *Game) reset() {
	g.x, g.y, g.angle = 1.5, 1.5, 0
	g.frame, g.shots, g.hasKey, g.win = 0, 0, false, false
	g.actors = []actor{{5.5, 5.5, 0, true}, {8.5, 8.5, 0, true}, {8.5, 1.5, 1, true}, {10.5, 9.5, 2, true}}
}

func (g *Game) blocked(x, y float64) bool {
	ix, iy := int(x), int(y)
	return iy < 0 || iy >= len(world) || ix < 0 || ix >= len(world[iy]) || world[iy][ix] != 0
}

func (g *Game) move(amount float64) {
	nx, ny := g.x+math.Cos(g.angle)*amount, g.y+math.Sin(g.angle)*amount
	if !g.blocked(nx, g.y) {
		g.x = nx
	}
	if !g.blocked(g.x, ny) {
		g.y = ny
	}
}

func touchAt(x0, y0, x1, y1 int) bool {
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if x >= x0 && x < x1 && y >= y0 && y < y1 {
			return true
		}
	}
	return false
}

func justTouched(x0, y0, x1, y1 int) bool {
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if x >= x0 && x < x1 && y >= y0 && y < y1 {
			return true
		}
	}
	return false
}

func (g *Game) shoot() {
	g.shots, g.frame = 7, g.frame
	best, bestDepth := -1, 99.0
	for i := range g.actors {
		a := &g.actors[i]
		if a.kind != 0 || !a.alive {
			continue
		}
		p := raycastlogic.ProjectSprite(g.x, g.y, g.angle, a.x, a.y, FOV)
		if p.Depth > 0 && math.Abs(p.ScreenX) < .13 {
			wall := raycastlogic.Cast(world, g.x, g.y, math.Atan2(a.y-g.y, a.x-g.x))
			if p.Depth < wall.Distance+.2 && p.Depth < bestDepth {
				best, bestDepth = i, p.Depth
			}
		}
	}
	if best >= 0 {
		g.actors[best].alive = false
		g.message = "CLEAR SHOT!"
	} else {
		g.message = "THE RAY HIT NO ENEMY"
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || justTouched(610, 0, 720, 48) {
		g.reset()
	}
	if g.win {
		return nil
	}
	g.frame++
	turn := 0.035
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) || touchAt(0, 420, 120, H) {
		g.angle -= turn
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) || touchAt(120, 420, 240, H) {
		g.angle += turn
	}
	g.angle = raycastlogic.WrapAngle(g.angle)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) || touchAt(240, 420, 430, H) {
		g.move(.045)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.move(-.03)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || justTouched(430, 420, 720, H) {
		g.shoot()
	}
	if g.shots > 0 {
		g.shots--
	}
	if g.variant == EbiRaycaster {
		for i := range g.actors {
			a := &g.actors[i]
			if !a.alive {
				continue
			}
			d := math.Hypot(a.x-g.x, a.y-g.y)
			if a.kind == 1 && d < .55 {
				a.alive = false
				g.hasKey = true
				g.message = "GOLDEN KEY FOUND!"
			}
			if a.kind == 2 && d < .65 {
				if g.hasKey {
					g.win = true
					g.message = "MAZE ESCAPED — PRESS R"
				} else {
					g.message = "THE EXIT NEEDS THE GOLDEN KEY"
				}
			}
		}
	}
	return nil
}

func fill(screen *ebiten.Image, x, y, w, h float32, c color.Color) {
	vector.DrawFilledRect(screen, x, y, w, h, c, false)
}
func label(screen *ebiten.Image, s string, x, y int, c color.Color) {
	text.Draw(screen, s, basicfont.Face7x13, x, y, c)
}

func (g *Game) drawMap(screen *ebiten.Image, x0, y0, size float32, rays int) {
	cell := size / 12
	fill(screen, x0-4, y0-4, size+8, size+8, color.RGBA{8, 16, 38, 230})
	for y := range world {
		for x := range world[y] {
			c := color.RGBA{20, 40, 62, 255}
			if world[y][x] != 0 {
				c = color.RGBA{54, 115, 126, 255}
			}
			fill(screen, x0+float32(x)*cell, y0+float32(y)*cell, cell-1, cell-1, c)
		}
	}
	if rays > 0 {
		for i := 0; i < rays; i++ {
			t := 0.5
			if rays > 1 {
				t = float64(i) / float64(rays-1)
			}
			a := g.angle - FOV/2 + t*FOV
			hit := raycastlogic.Cast(world, g.x, g.y, a)
			x1, y1 := g.x+math.Cos(a)*hit.Distance, g.y+math.Sin(a)*hit.Distance
			vector.StrokeLine(screen, x0+float32(g.x)*cell, y0+float32(g.y)*cell, x0+float32(x1)*cell, y0+float32(y1)*cell, 1, color.RGBA{255, 205, 70, 160}, false)
		}
	}
	px, py := x0+float32(g.x)*cell, y0+float32(g.y)*cell
	vector.DrawFilledCircle(screen, px, py, 5, color.RGBA{46, 230, 200, 255}, false)
	vector.StrokeLine(screen, px, py, px+float32(math.Cos(g.angle))*18, py+float32(math.Sin(g.angle))*18, 3, color.White, false)
}

func wallColor(hit raycastlogic.Hit, textured bool, column int) color.RGBA {
	shade := uint8(205)
	if hit.Side == 1 {
		shade = 150
	}
	if !textured {
		return color.RGBA{40, shade, 205, 255}
	}
	tile := (int(hit.WallU*10) + column/9) % 2
	if tile == 0 {
		return color.RGBA{shade, 92, 72, 255}
	}
	return color.RGBA{255, shade, 74, 255}
}

func (g *Game) draw3D(screen *ebiten.Image, textured bool) {
	top, bottom := 42, 414
	fill(screen, 0, float32(top), W, float32((bottom-top)/2), color.RGBA{18, 29, 63, 255})
	fill(screen, 0, float32((top+bottom)/2), W, float32((bottom-top)/2), color.RGBA{20, 30, 34, 255})
	for x := 0; x < W; x += 2 {
		a := g.angle - FOV/2 + (float64(x)+1)/W*FOV
		hit := raycastlogic.Cast(world, g.x, g.y, a)
		d := raycastlogic.CorrectDistance(hit.Distance, a, g.angle)
		g.zbuf[x], g.zbuf[min(x+1, W-1)] = d, d
		h := raycastlogic.ProjectHeight(bottom-top, d)
		y := (top + bottom - h) / 2
		c := wallColor(hit, textured, x)
		fog := math.Min(1, .35+4/(d+3))
		c.R, c.G, c.B = uint8(float64(c.R)*fog), uint8(float64(c.G)*fog), uint8(float64(c.B)*fog)
		fill(screen, float32(x), float32(y), 2, float32(h), c)
	}
	if g.variant == EbiRaycaster {
		g.drawActors(screen, top, bottom)
	}
	if g.shots > 0 {
		centerY := float32((top + bottom) / 2)
		vector.StrokeLine(screen, W/2-12, centerY, W/2+12, centerY, 2, color.White, false)
		vector.StrokeLine(screen, W/2, centerY-12, W/2, centerY+12, 2, color.White, false)
	}
}

func (g *Game) drawActors(screen *ebiten.Image, top, bottom int) {
	type visible struct {
		i int
		p raycastlogic.Projection
	}
	var list []visible
	for i, a := range g.actors {
		if a.alive {
			p := raycastlogic.ProjectSprite(g.x, g.y, g.angle, a.x, a.y, FOV)
			if p.Depth > 0 && math.Abs(p.ScreenX) < 1.25 {
				list = append(list, visible{i, p})
			}
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].p.Depth > list[j].p.Depth })
	for _, v := range list {
		a, p := g.actors[v.i], v.p
		size := float64(bottom-top) * .7 / p.Depth
		if size > 240 {
			size = 240
		}
		cx := float64(W)/2 + p.ScreenX*float64(W)/2
		col := color.RGBA{247, 84, 104, 255}
		if a.kind == 1 {
			col = color.RGBA{255, 205, 70, 255}
		}
		if a.kind == 2 {
			col = color.RGBA{46, 230, 200, 255}
		}
		left, right := int(cx-size/2), int(cx+size/2)
		for x := left; x < right; x += 2 {
			if x < 0 || x >= W || p.Depth >= g.zbuf[x] {
				continue
			}
			half := size / 2
			fill(screen, float32(x), float32(float64(top+bottom)/2-half), 2, float32(size), col)
		}
		if cx > 20 && cx < W-20 {
			icon := "ENEMY"
			if a.kind == 1 {
				icon = "KEY"
			}
			if a.kind == 2 {
				icon = "EXIT"
			}
			label(screen, icon, int(cx)-18, (top+bottom)/2-int(size/2)-5, color.White)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 17, 40, 255})
	titles := []string{"01  DIRECTION + MOVEMENT", "02  CAST ONE RAY", "03  DISTANCE -> WALL HEIGHT", "04  ONE RAY PER COLUMN", "05  TEXTURE + FISHEYE FIX", "06  EBI RAYCASTER"}
	label(screen, titles[g.variant], 18, 27, color.RGBA{46, 230, 200, 255})
	label(screen, "ARROWS/WASD: MOVE + TURN   SPACE: RAY SHOT   R: RESET", 250, 27, color.RGBA{190, 200, 230, 255})
	switch g.variant {
	case FacingMove:
		g.drawMap(screen, 176, 62, 356, 0)
		label(screen, "THE WORLD IS STILL A 2D GRID", 242, 405, color.White)
	case SingleRay:
		g.drawMap(screen, 176, 62, 356, 1)
		hit := raycastlogic.Cast(world, g.x, g.y, g.angle)
		label(screen, fmt.Sprintf("FIRST WALL: (%d,%d)   DISTANCE: %.2f", hit.MapX, hit.MapY, hit.Distance), 210, 405, color.White)
	case DistanceStrip:
		g.drawMap(screen, 26, 84, 250, 1)
		hit := raycastlogic.Cast(world, g.x, g.y, g.angle)
		h := raycastlogic.ProjectHeight(300, hit.Distance)
		fill(screen, 420, float32(238-h/2), 90, float32(h), wallColor(hit, false, 0))
		label(screen, fmt.Sprintf("%.2f TILES AWAY", hit.Distance), 405, 405, color.White)
	case ColumnView:
		g.draw3D(screen, false)
		g.drawMap(screen, 550, 56, 145, 9)
	case TexturedView, EbiRaycaster:
		g.draw3D(screen, true)
		g.drawMap(screen, 584, 58, 108, 5)
	}
	if g.variant == EbiRaycaster {
		key := "NO KEY"
		if g.hasKey {
			key = "KEY READY"
		}
		label(screen, key+"   "+g.message, 18, 405, color.White)
		if g.win {
			fill(screen, 115, 150, 490, 120, color.RGBA{8, 18, 40, 230})
			label(screen, "MAZE ESCAPED!", 294, 205, color.RGBA{46, 230, 200, 255})
			label(screen, "PRESS R OR TAP RESET TO PLAY AGAIN", 230, 235, color.White)
		}
	}
	buttons := []struct {
		x, w int
		s    string
	}{{0, 120, "TURN LEFT"}, {120, 120, "TURN RIGHT"}, {240, 190, "MOVE"}, {430, 290, "RAY / SHOOT"}}
	for _, b := range buttons {
		fill(screen, float32(b.x), 420, float32(b.w-2), 60, color.RGBA{26, 39, 80, 255})
		label(screen, b.s, b.x+18, 455, color.White)
	}
}

func (g *Game) Layout(_, _ int) (int, int) { return W, H }

func Run(variant Variant) {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Raycaster")
	if err := ebiten.RunGame(New(variant)); err != nil {
		panic(err)
	}
}
