// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package raycasterui

import (
	"fmt"
	"image/color"
	"math"
	"sort"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/audiolab"
	"github.com/kumagi/EbiShowcase/internal/cameralab"
	"github.com/kumagi/EbiShowcase/internal/ogfont"
	"github.com/kumagi/EbiShowcase/internal/raycastlogic"
	"github.com/kumagi/EbiShowcase/internal/shaderlab"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

// missions deliberately differ in routes and enemy pressure. The final game
// selects these data rather than copying rules for each maze.
var missions = []raycastlogic.Mission{
	{Name: "SUNSET HALL", Grid: world, StartX: 1.5, StartY: 1.5, StartAngle: 0, KeyX: 8.5, KeyY: 1.5, ExitX: 10.5, ExitY: 9.5, Enemies: []raycastlogic.Point{{5.5, 5.5}, {8.5, 8.5}}, GoalTime: 60},
	{Name: "TEAL VAULT", Grid: [][]int{{1, 1, 1, 1, 1, 1, 1, 1}, {1, 0, 0, 0, 1, 0, 0, 1}, {1, 0, 1, 0, 1, 0, 1, 1}, {1, 0, 1, 0, 0, 0, 0, 1}, {1, 0, 1, 1, 1, 1, 0, 1}, {1, 0, 0, 0, 0, 1, 0, 1}, {1, 1, 1, 1, 0, 0, 0, 1}, {1, 1, 1, 1, 1, 1, 1, 1}}, StartX: 1.5, StartY: 1.5, StartAngle: 0, KeyX: 6.5, KeyY: 3.5, ExitX: 5.5, ExitY: 6.5, Enemies: []raycastlogic.Point{{3.5, 3.5}, {2.5, 5.5}, {6.5, 5.5}}, GoalTime: 75},
	{Name: "NIGHT SPIRAL", Grid: [][]int{{1, 1, 1, 1, 1, 1, 1, 1, 1}, {1, 0, 0, 0, 0, 0, 0, 0, 1}, {1, 0, 1, 1, 1, 1, 1, 0, 1}, {1, 0, 1, 0, 0, 0, 1, 0, 1}, {1, 0, 1, 0, 1, 0, 1, 0, 1}, {1, 0, 1, 0, 1, 0, 0, 0, 1}, {1, 0, 1, 0, 1, 1, 1, 1, 1}, {1, 0, 0, 0, 0, 0, 0, 0, 1}, {1, 1, 1, 1, 1, 1, 1, 1, 1}}, StartX: 1.5, StartY: 7.5, StartAngle: -math.Pi / 2, KeyX: 3.5, KeyY: 3.5, ExitX: 7.5, ExitY: 1.5, Enemies: []raycastlogic.Point{{2.5, 1.5}, {5.5, 3.5}, {7.5, 5.5}}, GoalTime: 90},
}

type actor struct {
	x, y  float64
	kind  int // 0 enemy, 1 key, 2 exit
	alive bool
}

var (
	rayFontOnce sync.Once
	rayFontBase *opentype.Font
	rayFontErr  error
	rayFaces    = map[float64]font.Face{}
)

func rayFace(size float64) font.Face {
	rayFontOnce.Do(func() { rayFontBase, rayFontErr = opentype.Parse(ogfont.NotoSansJP) })
	if rayFontErr != nil {
		panic(rayFontErr)
	}
	if face := rayFaces[size]; face != nil {
		return face
	}
	face, err := opentype.NewFace(rayFontBase, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		panic(err)
	}
	rayFaces[size] = face
	return face
}

type Game struct {
	variant    Variant
	x, y       float64
	angle      float64
	frame      int
	shots      int
	totalShots int
	hasKey     bool
	win        bool
	message    string
	actors     []actor
	zbuf       []float64
	mission    int
	grid       [][]int
	hp         int
	damage     int
	flash      int
	ticks      int
	best       int
	canvas     *ebiten.Image
	viewW      int
	viewH      int
	lang       string
	audio      *audio.Context
	gate       audiolab.Gate
	pulse      *shaderlab.Pulse
	cam        cameralab.State
	badge      *ebiten.Image
}

func New(variant Variant) *Game {
	g := &Game{variant: variant, zbuf: make([]float64, W), canvas: ebiten.NewImage(W, H), viewW: W, viewH: H, lang: browserLanguage(), message: "TURN, MOVE, AND READ THE RAYS"}
	g.audio = audio.NewContext(audiolab.SampleRate)
	g.pulse = shaderlab.NewPulse()
	g.cam = cameralab.State{ViewW: W, ViewH: H}
	g.badge = ebiten.NewImage(20, 20)
	g.badge.Fill(color.RGBA{46, 230, 200, 255})
	g.reset()
	return g
}

func (g *Game) reset() {
	m := missions[g.mission]
	g.best = storedBest(fmt.Sprintf("ebiShowcase.raycaster.best.%d", g.mission))
	g.grid = m.Grid
	g.x, g.y, g.angle = m.StartX, m.StartY, m.StartAngle
	g.frame, g.shots, g.totalShots, g.hasKey, g.win, g.ticks, g.damage, g.flash = 0, 0, 0, false, false, 0, 0, 0
	g.hp = 3
	g.message = g.tr("TURN, MOVE, AND READ THE RAYS", "旋回・移動・光線を読み取ろう")
	g.actors = make([]actor, 0, len(m.Enemies)+2)
	for _, p := range m.Enemies {
		g.actors = append(g.actors, actor{p.X, p.Y, 0, true})
	}
	g.actors = append(g.actors, actor{m.KeyX, m.KeyY, 1, true}, actor{m.ExitX, m.ExitY, 2, true})
}

func (g *Game) tr(en, ja string) string {
	if g.lang == "ja" {
		return ja
	}
	return en
}

func (g *Game) blocked(x, y float64) bool {
	ix, iy := int(x), int(y)
	return iy < 0 || iy >= len(g.grid) || ix < 0 || ix >= len(g.grid[iy]) || g.grid[iy][ix] != 0
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

func (g *Game) toVirtual(x, y int) (int, int) {
	if g.viewW <= 0 || g.viewH <= 0 {
		return x, y
	}
	if g.portrait() && y >= 560 {
		return x * W / g.viewW, 420 + (y-560)*60/(g.viewH-560)
	}
	return x * W / g.viewW, y * H / g.viewH
}

func (g *Game) portrait() bool { return g.viewH > g.viewW }

func (g *Game) rawJustPressed(x0, y0, x1, y1 int) bool {
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if x >= x0 && x < x1 && y >= y0 && y < y1 {
			return true
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return x >= x0 && x < x1 && y >= y0 && y < y1
	}
	return false
}

func (g *Game) touchAt(x0, y0, x1, y1 int) bool {
	for _, id := range ebiten.AppendTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		x, y = g.toVirtual(x, y)
		if x >= x0 && x < x1 && y >= y0 && y < y1 {
			return true
		}
	}
	return false
}

func (g *Game) justTouched(x0, y0, x1, y1 int) bool {
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		x, y = g.toVirtual(x, y)
		if x >= x0 && x < x1 && y >= y0 && y < y1 {
			return true
		}
	}
	return false
}

func (g *Game) pressedAt(x0, y0, x1, y1 int) bool {
	if g.justTouched(x0, y0, x1, y1) {
		return true
	}
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}
	x, y := ebiten.CursorPosition()
	x, y = g.toVirtual(x, y)
	return x >= x0 && x < x1 && y >= y0 && y < y1
}

func (g *Game) shoot() {
	g.playSE(680)
	g.shots, g.frame = 7, g.frame
	g.totalShots++
	best, bestDepth := -1, 99.0
	for i := range g.actors {
		a := &g.actors[i]
		if a.kind != 0 || !a.alive {
			continue
		}
		p := raycastlogic.ProjectSprite(g.x, g.y, g.angle, a.x, a.y, FOV)
		if p.Depth > 0 && math.Abs(p.ScreenX) < .13 {
			wall := raycastlogic.Cast(g.grid, g.x, g.y, math.Atan2(a.y-g.y, a.x-g.x))
			if p.Depth < wall.Distance+.2 && p.Depth < bestDepth {
				best, bestDepth = i, p.Depth
			}
		}
	}
	if best >= 0 {
		g.actors[best].alive = false
		g.message = g.tr("CLEAR SHOT!", "命中！")
	} else {
		g.message = g.tr("THE RAY HIT NO ENEMY", "敵には当たりませんでした")
	}
}

func (g *Game) playSE(freq float64) {
	g.gate.Arm(true)
	g.audio.NewPlayerF32FromBytes(audiolab.OneShot(audiolab.Square, freq, .05)).Play()
}

func (g *Game) Update() error {
	if g.variant == EbiRaycaster {
		for i, key := range []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3} {
			if inpututil.IsKeyJustPressed(key) {
				g.mission = i
				g.reset()
				return nil
			}
		}
		if g.portrait() {
			for i := 0; i < 3; i++ {
				if g.rawJustPressed(18+i*118, 390, 18+(i+1)*118, 434) {
					g.mission = i
					g.reset()
					return nil
				}
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || g.pressedAt(610, 0, 720, 48) || (g.portrait() && g.rawJustPressed(18, 462, 372, 510)) {
		g.reset()
	}
	if g.win {
		return nil
	}
	g.frame++
	g.ticks++
	if g.flash > 0 {
		g.flash--
	}
	turn := 0.035
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) || g.touchAt(0, 420, 120, H) {
		g.angle -= turn
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) || g.touchAt(120, 420, 240, H) {
		g.angle += turn
	}
	g.angle = raycastlogic.WrapAngle(g.angle)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) || g.touchAt(240, 420, 360, H) {
		g.move(.045)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) || g.touchAt(360, 420, 430, H) {
		g.move(-.03)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || g.pressedAt(430, 420, 720, H) {
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
				g.message = g.tr("GOLDEN KEY FOUND!", "金の鍵を手に入れた！")
			}
			if a.kind == 2 && d < .65 {
				if g.hasKey {
					g.win = true
					grade := raycastlogic.Grade(g.ticks/60, g.damage, g.totalShots, len(missions[g.mission].Enemies))
					if g.best == 0 || g.ticks < g.best {
						g.best = g.ticks
					}
					storeBest(fmt.Sprintf("ebiShowcase.raycaster.best.%d", g.mission), g.best)
					g.message = g.tr("MAZE ESCAPED! GRADE "+grade+" — PRESS R", "脱出成功！ 評価 "+grade+" — Rで再挑戦")
				} else {
					g.message = g.tr("THE EXIT NEEDS THE GOLDEN KEY", "出口には金の鍵が必要です")
				}
			}
			if a.kind == 0 && d < .72 && g.frame%45 == 0 && g.hp > 0 {
				g.hp--
				g.damage++
				g.flash = 14
				g.message = g.tr("HIT! FIND COVER OR FIRE", "ダメージ！ 隠れるか撃とう")
				if g.hp == 0 {
					g.message = g.tr("SYSTEM DOWN — PRESS R", "システム停止 — Rで再挑戦")
					g.win = true
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
	text.Draw(screen, s, rayFace(14), x, y, c)
}

func (g *Game) drawMap(screen *ebiten.Image, x0, y0, size float32, rays int) {
	cell := size / float32(max(len(g.grid), len(g.grid[0])))
	fill(screen, x0-4, y0-4, size+8, size+8, color.RGBA{8, 16, 38, 230})
	for y := range g.grid {
		for x := range g.grid[y] {
			c := color.RGBA{20, 40, 62, 255}
			if g.grid[y][x] != 0 {
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
			hit := raycastlogic.Cast(g.grid, g.x, g.y, a)
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
		hit := raycastlogic.Cast(g.grid, g.x, g.y, a)
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
		pulse := float32(8 + (7-g.shots)*5)
		vector.StrokeCircle(screen, W/2, centerY, pulse, 3, color.RGBA{255, 211, 112, 220}, false)
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

func (g *Game) drawVirtual(screen *ebiten.Image) {
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
		hit := raycastlogic.Cast(g.grid, g.x, g.y, g.angle)
		label(screen, fmt.Sprintf("FIRST WALL: (%d,%d)   DISTANCE: %.2f", hit.MapX, hit.MapY, hit.Distance), 210, 405, color.White)
	case DistanceStrip:
		g.drawMap(screen, 26, 84, 250, 1)
		hit := raycastlogic.Cast(g.grid, g.x, g.y, g.angle)
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
		key := g.tr("NO KEY", "鍵なし")
		if g.hasKey {
			key = g.tr("KEY READY", "鍵あり")
		}
		label(screen, fmt.Sprintf("M%d %s  HP:%d  %s", g.mission+1, missions[g.mission].Name, g.hp, key), 18, 389, color.RGBA{255, 211, 112, 255})
		label(screen, key+"   "+g.message, 18, 405, color.White)
		best := "--:--"
		if g.best > 0 {
			best = fmt.Sprintf("%02d:%02d", g.best/3600, g.best/60%60)
		}
		label(screen, "1/2/3: MISSION   "+fmt.Sprintf("TIME %02d:%02d  BEST %s", g.ticks/3600, g.ticks/60%60, best), 365, 405, color.RGBA{184, 211, 233, 255})
		if g.win {
			fill(screen, 115, 150, 490, 120, color.RGBA{8, 18, 40, 230})
			label(screen, "MAZE ESCAPED!", 294, 205, color.RGBA{46, 230, 200, 255})
			label(screen, "PRESS R OR TAP RESET TO PLAY AGAIN", 230, 235, color.White)
		}
	}
	if g.flash > 0 {
		fill(screen, 0, 42, W, 372, color.RGBA{220, 45, 70, uint8(g.flash * 9)})
	}
	buttons := []struct {
		x, w int
		s    string
	}{{0, 120, "TURN LEFT"}, {120, 120, "TURN RIGHT"}, {240, 120, "FORWARD"}, {360, 70, "BACK"}, {430, 290, "RAY / SHOOT"}}
	for _, b := range buttons {
		fill(screen, float32(b.x), 420, float32(b.w-2), 60, color.RGBA{26, 39, 80, 255})
		label(screen, b.s, b.x+18, 455, color.White)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.cam.ViewW, g.cam.ViewH = float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	g.drawEffectBadge(screen)
	if g.portrait() {
		g.drawPortrait(screen)
		return
	}
	g.drawVirtual(g.canvas)
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(sw)/W, float64(sh)/H)
	screen.DrawImage(g.canvas, op)
}

func (g *Game) drawEffectBadge(screen *ebiten.Image) {
	if g.variant != EbiRaycaster || g.pulse == nil || !g.pulse.Available() {
		return
	}
	fx := ebiten.NewImage(20, 20)
	if !g.pulse.Draw(fx, g.badge, float32(g.ticks)*.08) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screen.Bounds().Dx()-34), 10)
	screen.DrawImage(fx, op)
}

func (g *Game) drawPortrait(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 17, 40, 255})
	g.drawVirtual(g.canvas)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screen.Bounds().Dx())/W, 280.0/H)
	screen.DrawImage(g.canvas, op)
	fill(screen, 0, 280, float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy()-280), color.RGBA{10, 17, 40, 255})
	mission := missions[g.mission]
	label(screen, g.tr("MISSION ", "ミッション ")+fmt.Sprint(g.mission+1)+": "+mission.Name, 18, 310, color.RGBA{46, 230, 200, 255})
	key := g.tr("NO KEY", "鍵なし")
	if g.hasKey {
		key = g.tr("KEY READY", "鍵あり")
	}
	label(screen, fmt.Sprintf("HP %d/3   %s   TIME %02d:%02d", g.hp, key, g.ticks/3600, g.ticks/60%60), 18, 334, color.RGBA{255, 211, 112, 255})
	label(screen, g.message, 18, 360, color.White)
	for i := 0; i < 3; i++ {
		x := float32(18 + i*118)
		c := color.RGBA{26, 39, 80, 255}
		if i == g.mission {
			c = color.RGBA{46, 130, 150, 255}
		}
		fill(screen, x, 390, 110, 44, c)
		label(screen, fmt.Sprintf("%d %s", i+1, []string{"HALL", "VAULT", "SPIRAL"}[i]), int(x+9), 416, color.White)
	}
	fill(screen, 18, 462, 354, 48, color.RGBA{39, 82, 112, 255})
	label(screen, g.tr("REPLAY / RESET", "もう一度 / リセット"), 132, 491, color.White)
	label(screen, g.tr("TOUCH CONTROLS", "タッチ操作"), 18, 540, color.RGBA{184, 211, 233, 255})
	buttons := []struct {
		x, w  int
		label string
	}{
		{0, 65, g.tr("LEFT", "左")}, {65, 65, g.tr("RIGHT", "右")}, {130, 65, g.tr("GO", "進む")}, {195, 38, g.tr("BACK", "戻る")}, {233, 157, g.tr("FIRE", "撃つ")},
	}
	for _, b := range buttons {
		fill(screen, float32(b.x), 560, float32(b.w-2), 160, color.RGBA{26, 39, 80, 255})
		label(screen, b.label, b.x+max(4, (b.w-len(b.label)*7)/2), 646, color.White)
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	if outsideH > outsideW {
		// A portrait canvas uses the available height rather than preserving the
		// old desktop aspect ratio; Draw scales the virtual game into it.
		g.viewW, g.viewH = 390, 720
		return g.viewW, g.viewH
	}
	g.viewW, g.viewH = W, H
	return W, H
}

func Run(variant Variant) {
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Ebi Raycaster")
	if err := ebiten.RunGame(New(variant)); err != nil {
		panic(err)
	}
}
