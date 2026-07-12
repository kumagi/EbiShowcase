package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW = 480
	screenH = 720
	cols    = 9
	rows    = 9
	cell    = 48
	boardX  = 24
	boardY  = 106
	fuse    = 105
)

type point struct{ x, y int }
type tile int
type item int

const (
	floor tile = iota
	hardWall
	softWall
)

const (
	noItem item = iota
	fireItem
	bombItem
	speedItem
)

type bomb struct {
	at    point
	timer int
}

type flame struct {
	at    point
	timer int
}

type game struct {
	board             [rows][cols]tile
	items             map[point]item
	bombs             []bomb
	flames            []flame
	player            point
	power, capacity   int
	speed             int
	moveCooldown      int
	destroyed, target int
	frames            int
	message           string
	won, lost         bool
}

func newGame() *game {
	g := &game{
		items:    map[point]item{},
		player:   point{1, 1},
		power:    1,
		capacity: 1,
		target:   10,
		message:  "Break 10 wooden walls. Some hide useful items!",
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			switch {
			case x == 0 || y == 0 || x == cols-1 || y == rows-1 || (x%2 == 0 && y%2 == 0):
				g.board[y][x] = hardWall
			case !(x <= 2 && y <= 2) && (x*5+y*7)%4 != 0:
				g.board[y][x] = softWall
			}
		}
	}
	return g
}

func (g *game) Update() error {
	if g.won || g.lost {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	g.frames++
	if g.moveCooldown > 0 {
		g.moveCooldown--
	}
	if d, ok := inputDir(); ok && g.moveCooldown == 0 {
		n := point{g.player.x + d.x, g.player.y + d.y}
		if g.board[n.y][n.x] == floor && !g.bombAt(n) {
			g.player = n
			g.collect(n)
			g.moveCooldown = max(3, 9-g.speedLevel()*2)
		}
	}
	if placePressed() && len(g.bombs) < g.capacity && !g.bombAt(g.player) {
		g.bombs = append(g.bombs, bomb{at: g.player, timer: fuse})
		g.message = "Bomb set. Move away before the fuse reaches zero!"
	}
	g.updateBombs()
	g.updateFlames()
	if g.frames >= 75*60 {
		g.lost = true
		g.message = "Time up. Retry and choose a shorter route."
	}
	return nil
}

func (g *game) updateBombs() {
	next := g.bombs[:0]
	for _, b := range g.bombs {
		b.timer--
		if b.timer > 0 {
			next = append(next, b)
			continue
		}
		g.explode(b.at)
	}
	g.bombs = next
}

func (g *game) explode(origin point) {
	g.addFlame(origin)
	for _, d := range []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		for distance := 1; distance <= g.power; distance++ {
			p := point{origin.x + d.x*distance, origin.y + d.y*distance}
			if g.board[p.y][p.x] == hardWall {
				break
			}
			g.addFlame(p)
			if g.board[p.y][p.x] == softWall {
				g.board[p.y][p.x] = floor
				g.destroyed++
				if drop := deterministicDrop(p); drop != noItem {
					g.items[p] = drop
					g.message = "A wall became floor and revealed an item!"
				} else {
					g.message = "A wall became floor. This one was empty."
				}
				if g.destroyed >= g.target {
					g.won = true
					g.message = "Ten walls changed into floor. Goal complete!"
				}
				break
			}
		}
	}
}

func deterministicDrop(p point) item {
	// A fixed coordinate hash makes the lesson repeatable without global randomness.
	switch (p.x*17 + p.y*31 + 7) % 5 {
	case 0:
		return fireItem
	case 1:
		return bombItem
	case 2:
		return speedItem
	default:
		return noItem
	}
}

func (g *game) addFlame(p point) { g.flames = append(g.flames, flame{at: p, timer: 25}) }

func (g *game) updateFlames() {
	next := g.flames[:0]
	for _, f := range g.flames {
		f.timer--
		if f.at == g.player {
			g.lost = true
			g.message = "The blast caught you. Retry and leave an escape path."
		}
		if f.timer > 0 {
			next = append(next, f)
		}
	}
	g.flames = next
}

func (g *game) collect(p point) {
	switch g.items[p] {
	case fireItem:
		g.power = min(3, g.power+1)
		g.message = "FIRE: future blasts reach one tile farther."
	case bombItem:
		g.capacity = min(3, g.capacity+1)
		g.message = "BOMB: you can keep one more bomb active."
	case speedItem:
		g.speed = min(2, g.speed+1)
		g.message = "SPEED: movement cooldown became shorter."
	}
	if g.items[p] != noItem {
		delete(g.items, p)
	}
}

func (g *game) speedLevel() int {
	return g.speed
}

func (g *game) bombAt(p point) bool {
	for _, b := range g.bombs {
		if b.at == p {
			return true
		}
	}
	return false
}

func (g *game) Draw(s *ebiten.Image) {
	s.Fill(color.RGBA{7, 15, 29, 255})
	ebitenutil.DebugPrintAt(s, "BREAKABLE WALL WORKSHOP", 145, 18)
	ebitenutil.DebugPrintAt(s, fmt.Sprintf("WALLS %02d/%02d  FIRE %d  BOMBS %d  TIME %02d", g.destroyed, g.target, g.power, g.capacity, max(0, 75-g.frames/60)), 78, 45)
	ebitenutil.DebugPrintAt(s, g.message, 30, 73)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			p := point{x, y}
			px, py := float32(boardX+x*cell), float32(boardY+y*cell)
			c := color.RGBA{24, 48, 48, 255}
			if g.board[y][x] == hardWall {
				c = color.RGBA{72, 93, 125, 255}
			} else if g.board[y][x] == softWall {
				c = color.RGBA{151, 91, 54, 255}
			}
			vector.DrawFilledRect(s, px+1, py+1, cell-2, cell-2, c, false)
			if g.board[y][x] == softWall {
				vector.StrokeRect(s, px+7, py+7, cell-14, cell-14, 2, color.RGBA{226, 155, 82, 255}, false)
			}
			if it := g.items[p]; it != noItem {
				drawItem(s, px+cell/2, py+cell/2, it)
			}
		}
	}
	for _, b := range g.bombs {
		x := float32(boardX + b.at.x*cell + cell/2)
		y := float32(boardY + b.at.y*cell + cell/2)
		vector.DrawFilledCircle(s, x, y, 14, color.RGBA{20, 23, 34, 255}, false)
		vector.DrawFilledCircle(s, x+8, y-10, 3, color.RGBA{255, 205, 81, 255}, false)
	}
	for _, f := range g.flames {
		x := float32(boardX + f.at.x*cell + cell/2)
		y := float32(boardY + f.at.y*cell + cell/2)
		vector.DrawFilledCircle(s, x, y, 20, color.RGBA{255, 158, 52, 225}, false)
		vector.DrawFilledCircle(s, x, y, 9, color.RGBA{255, 244, 135, 255}, false)
	}
	px := float32(boardX + g.player.x*cell + cell/2)
	py := float32(boardY + g.player.y*cell + cell/2)
	vector.DrawFilledCircle(s, px, py, 14, color.RGBA{235, 91, 76, 255}, false)
	vector.DrawFilledCircle(s, px-5, py-3, 2, color.White, false)
	vector.DrawFilledCircle(s, px+5, py-3, 2, color.White, false)
	drawControls(s)
	if g.won {
		drawOverlay(s, "WALL WORK COMPLETE!\n\nTAP / ENTER TO RETRY")
	} else if g.lost {
		drawOverlay(s, "BLASTED OR TIME UP\n\nTAP / ENTER TO RETRY")
	}
}

func drawItem(s *ebiten.Image, x, y float32, it item) {
	colors := map[item]color.RGBA{
		fireItem:  {255, 93, 61, 255},
		bombItem:  {126, 185, 255, 255},
		speedItem: {119, 235, 155, 255},
	}
	vector.DrawFilledCircle(s, x, y, 12, colors[it], false)
	label := map[item]string{fireItem: "F", bombItem: "B", speedItem: "S"}[it]
	ebitenutil.DebugPrintAt(s, label, int(x)-3, int(y)-5)
}

func drawControls(s *ebiten.Image) {
	labels := [...]string{"LEFT", "UP", "DOWN", "RIGHT", "BOMB"}
	for i, label := range labels {
		vector.DrawFilledRect(s, float32(i*96+3), 600, 90, 62, color.RGBA{45, 78, 113, 255}, false)
		ebitenutil.DebugPrintAt(s, label, i*96+27, 627)
	}
	ebitenutil.DebugPrintAt(s, "Arrows/WASD + Space, or tap buttons", 100, 688)
}

func inputDir() (point, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		return point{-1, 0}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		return point{0, -1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return point{0, 1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		return point{1, 0}, true
	}
	if x, y, ok := press(); ok && y >= 600 && x < 384 {
		return [...]point{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}[min(3, x/96)], true
	}
	return point{}, false
}

func placePressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyX) {
		return true
	}
	x, y, ok := press()
	return ok && y >= 600 && x >= 384
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

func retryPressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(inpututil.AppendJustPressedTouchIDs(nil)) > 0
}

func drawOverlay(s *ebiten.Image, text string) {
	vector.DrawFilledRect(s, 40, 270, 400, 160, color.RGBA{4, 12, 27, 245}, false)
	ebitenutil.DebugPrintAt(s, text, 115, 330)
}

func (g *game) Layout(int, int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Breakable Wall Workshop")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
