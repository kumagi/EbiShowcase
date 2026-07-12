package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	width, height = 480, 720
	worldW        = 96
	worldH        = 16
	tile          = 30
	viewCols      = width / tile
	air           = 0
	grass         = 1
	dirt          = 2
	stone         = 3
	timeLimit     = 38 * 60
)

type game struct {
	seed        int64
	tiles       [worldH][worldW]int
	surface     [worldW]int
	relics      [3]int
	found       [3]bool
	playerX     float64
	frames      int
	clear, over bool
	message     string
}

func newGame(seed int64) *game {
	g := &game{seed: seed, playerX: 2}
	g.generate()
	g.message = "Explore the generated heights and find 3 survey beacons."
	return g
}

func (g *game) generate() {
	for x := 0; x < worldW; x++ {
		surfaceY := 3 + int(valueNoise1D(g.seed, float64(x)/8)*6)
		g.surface[x] = surfaceY
		for y := 0; y < worldH; y++ {
			g.tiles[y][x] = air
			if y < surfaceY {
				continue
			}
			switch {
			case y == surfaceY:
				g.tiles[y][x] = grass
			case y <= surfaceY+3:
				g.tiles[y][x] = dirt
			default:
				g.tiles[y][x] = stone
				if valueNoise2D(g.seed+91, float64(x)/4, float64(y)/4) > 0.70 {
					g.tiles[y][x] = air
				}
			}
		}
	}
	g.relics = [3]int{
		14 + int(hash(g.seed, 11, 0)%12),
		43 + int(hash(g.seed, 22, 0)%14),
		75 + int(hash(g.seed, 33, 0)%14),
	}
}

func hash(seed int64, x, y int) uint64 {
	n := uint64(seed) ^ uint64(x)*0x632BE59BD9B4E019 ^ uint64(y)*0x9E3779B185EBCA87
	n ^= n >> 30
	n *= 0xBF58476D1CE4E5B9
	n ^= n >> 27
	n *= 0x94D049BB133111EB
	n ^= n >> 31
	return n
}

func random01(seed int64, x, y int) float64 {
	return float64(hash(seed, x, y)&0xFFFFFF) / float64(0xFFFFFF)
}

func smooth(t float64) float64 { return t * t * (3 - 2*t) }

func valueNoise1D(seed int64, x float64) float64 {
	x0 := int(math.Floor(x))
	t := smooth(x - float64(x0))
	a := random01(seed, x0, 0)
	b := random01(seed, x0+1, 0)
	return a + (b-a)*t
}

func valueNoise2D(seed int64, x, y float64) float64 {
	x0, y0 := int(math.Floor(x)), int(math.Floor(y))
	tx, ty := smooth(x-float64(x0)), smooth(y-float64(y0))
	a := random01(seed, x0, y0)
	b := random01(seed, x0+1, y0)
	c := random01(seed, x0, y0+1)
	d := random01(seed, x0+1, y0+1)
	top := a + (b-a)*tx
	bottom := c + (d-c)*tx
	return top + (bottom-top)*ty
}

func (g *game) Update() error {
	if g.clear || g.over {
		if retryPressed() {
			*g = *newGame(g.seed)
		}
		return nil
	}
	g.frames++
	if g.frames >= timeLimit {
		g.over = true
		g.message = "Survey time ended. Retry the same seed and remember the hills."
		return nil
	}

	direction := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		direction--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		direction++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*g = *newGame(g.seed)
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		*g = *newGame(g.seed + 1)
		return nil
	}
	if x, y, ok := justPressed(); ok && y >= 615 && x >= 160 && x < 320 {
		*g = *newGame(g.seed + 1)
		return nil
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 615 {
			if x < 160 {
				direction = -1
			} else if x >= 320 {
				direction = 1
			}
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 615 {
			if x < 160 {
				direction = -1
			} else if x >= 320 {
				direction = 1
			}
		}
	}
	g.playerX = math.Max(0, math.Min(worldW-1, g.playerX+direction*0.13))
	for i, x := range g.relics {
		if !g.found[i] && math.Abs(g.playerX-float64(x)) < 0.55 {
			g.found[i] = true
			g.message = fmt.Sprintf("Beacon sample %d/3 recorded at column %d.", g.foundCount(), x)
		}
	}
	if g.foundCount() == len(g.relics) {
		g.clear = true
		g.message = "All three layers surveyed in this reproducible world!"
	}
	return nil
}

func (g *game) foundCount() int {
	count := 0
	for _, found := range g.found {
		if found {
			count++
		}
	}
	return count
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{19, 38, 58, 255})
	ebitenutil.DebugPrintAt(screen, "SEEDED TERRAIN SURVEY", 165, 20)
	seconds := max(0, (timeLimit-g.frames+59)/60)
	column := int(math.Round(g.playerX))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SEED %d   TIME %02d   SAMPLES %d/3", g.seed, seconds, g.foundCount()), 128, 47)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("COLUMN %02d   HEIGHT %02d   %s", column, g.surface[column], g.message), 44, 72)

	camera := int(g.playerX) - viewCols/2
	camera = max(0, min(worldW-viewCols, camera))
	for sy := 0; sy < worldH; sy++ {
		for sx := 0; sx < viewCols; sx++ {
			wx := camera + sx
			px, py := float32(sx*tile), float32(100+sy*tile)
			switch g.tiles[sy][wx] {
			case grass:
				vector.DrawFilledRect(screen, px, py, tile, tile, color.RGBA{82, 166, 92, 255}, false)
			case dirt:
				vector.DrawFilledRect(screen, px, py, tile, tile, color.RGBA{143, 98, 62, 255}, false)
			case stone:
				vector.DrawFilledRect(screen, px, py, tile, tile, color.RGBA{92, 102, 116, 255}, false)
			}
			if g.tiles[sy][wx] != air {
				vector.StrokeRect(screen, px, py, tile, tile, 1, color.RGBA{25, 35, 45, 100}, false)
			}
		}
	}
	for i, wx := range g.relics {
		if wx < camera || wx >= camera+viewCols {
			continue
		}
		px := float32((wx-camera)*tile + tile/2)
		py := float32(100 + g.surface[wx]*tile - 12)
		c := color.RGBA{246, 187, 64, 255}
		if g.found[i] {
			c = color.RGBA{91, 201, 121, 255}
		}
		vector.StrokeCircle(screen, px, py, 10, 4, c, false)
		vector.StrokeLine(screen, px, py+10, px, py+25, 3, c, false)
	}
	playerScreenX := float32((g.playerX-float64(camera))*tile + tile/2)
	playerY := float32(100 + g.surface[column]*tile - 14)
	vector.DrawFilledCircle(screen, playerScreenX, playerY, 13, color.RGBA{238, 115, 76, 255}, false)
	vector.StrokeCircle(screen, playerScreenX, playerY, 13, 3, color.White, false)

	button(screen, 8, 620, 144, "EXPLORE LEFT", color.RGBA{53, 91, 124, 255})
	button(screen, 168, 620, 144, "NEXT SEED", color.RGBA{112, 78, 139, 255})
	button(screen, 328, 620, 144, "EXPLORE RIGHT", color.RGBA{53, 91, 124, 255})
	if g.clear || g.over {
		title := "WORLD SURVEY COMPLETE!"
		if g.over {
			title = "SURVEY TIME ENDED!"
		}
		vector.DrawFilledRect(screen, 40, 270, 400, 160, color.RGBA{5, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 150, 314)
		ebitenutil.DebugPrintAt(screen, g.message, 67, 350)
		ebitenutil.DebugPrintAt(screen, "TAP / SPACE: RETRY SAME SEED", 126, 390)
	}
}

func button(screen *ebiten.Image, x, y, w int, label string, fill color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), 66, fill, false)
	ebitenutil.DebugPrintAt(screen, label, x+w/2-len(label)*3, y+29)
}

func justPressed() (int, int, bool) {
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

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func (g *game) Layout(_, _ int) (int, int) { return width, height }

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Seeded Terrain Survey — Ebitengine")
	if err := ebiten.RunGame(newGame(7302)); err != nil {
		panic(err)
	}
}
