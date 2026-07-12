package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW    = 480
	screenH    = 720
	tileSize   = 32
	chunkTiles = 8
	chunkSize  = tileSize * chunkTiles
	timeLimit  = 75 * 60
)

type chunkPos struct{ x, y int }

type chunk struct {
	tiles [chunkTiles][chunkTiles]int
}

type game struct {
	playerX, playerY float64
	chunks           map[chunkPos]*chunk
	visited          map[chunkPos]bool
	beacons, frames  int
	message          string
	won, lost        bool
}

func newGame() *game {
	g := &game{
		playerX: 16,
		playerY: 16,
		chunks:  map[chunkPos]*chunk{},
		visited: map[chunkPos]bool{},
		message: "Explore new chunks and collect gold beacons!",
	}
	g.ensureAround(g.playerChunk())
	g.visitCurrentChunk()
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
	if g.frames >= timeLimit {
		g.lost = true
		g.message = "Expedition time ended. Try a shorter route!"
		return nil
	}

	dx, dy := movementInput()
	if dx != 0 || dy != 0 {
		length := math.Hypot(dx, dy)
		g.playerX += dx / length * 3.5
		g.playerY += dy / length * 3.5
	}
	center := g.playerChunk()
	g.ensureAround(center)
	g.visitCurrentChunk()
	g.collectBeacon()
	if len(g.visited) >= 12 && g.beacons >= 5 {
		g.won = true
		g.message = "Map complete: explored 12 chunks and found 5 beacons!"
	}
	return nil
}

func (g *game) playerChunk() chunkPos {
	tileX := int(math.Floor(g.playerX / tileSize))
	tileY := int(math.Floor(g.playerY / tileSize))
	return chunkPos{floorDiv(tileX, chunkTiles), floorDiv(tileY, chunkTiles)}
}

func (g *game) ensureAround(center chunkPos) {
	// Radius two covers the 480x505 viewport with one safety chunk around it.
	for cy := center.y - 2; cy <= center.y+2; cy++ {
		for cx := center.x - 2; cx <= center.x+2; cx++ {
			key := chunkPos{cx, cy}
			if _, exists := g.chunks[key]; !exists {
				g.chunks[key] = generateChunk(key)
			}
		}
	}
}

func generateChunk(key chunkPos) *chunk {
	c := &chunk{}
	for y := 0; y < chunkTiles; y++ {
		for x := 0; x < chunkTiles; x++ {
			worldX := key.x*chunkTiles + x
			worldY := key.y*chunkTiles + y
			n := hash(worldX, worldY) % 10
			switch {
			case n < 2:
				c.tiles[y][x] = 0 // shallow water
			case n < 8:
				c.tiles[y][x] = 1 // grass
			default:
				c.tiles[y][x] = 2 // stone
			}
		}
	}
	// Every chunk has one deterministic beacon. Its tile is modified in the
	// cached chunk when collected, proving that cached chunk state persists.
	bx := int(hash(key.x, key.y*17) % chunkTiles)
	by := int(hash(key.x*31, key.y) % chunkTiles)
	c.tiles[by][bx] = 3
	return c
}

func hash(x, y int) uint64 {
	n := uint64(int64(x))*0x9e3779b185ebca87 ^ uint64(int64(y))*0xc2b2ae3d27d4eb4f ^ 0x72e4a91d
	n ^= n >> 30
	n *= 0xbf58476d1ce4e5b9
	n ^= n >> 27
	n *= 0x94d049bb133111eb
	return n ^ (n >> 31)
}

func (g *game) visitCurrentChunk() {
	key := g.playerChunk()
	if !g.visited[key] {
		g.visited[key] = true
		g.message = fmt.Sprintf("Entered new chunk (%d,%d). Cached forever!", key.x, key.y)
	}
}

func (g *game) collectBeacon() {
	tileX := int(math.Floor(g.playerX / tileSize))
	tileY := int(math.Floor(g.playerY / tileSize))
	key := chunkPos{floorDiv(tileX, chunkTiles), floorDiv(tileY, chunkTiles)}
	localX := positiveMod(tileX, chunkTiles)
	localY := positiveMod(tileY, chunkTiles)
	c := g.chunks[key]
	if c != nil && c.tiles[localY][localX] == 3 {
		c.tiles[localY][localX] = 1
		g.beacons++
		g.message = fmt.Sprintf("Beacon collected! %d/5 (chunk stayed cached)", g.beacons)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{9, 20, 34, 255})
	center := g.playerChunk()
	seconds := max(0, (timeLimit-g.frames+59)/60)
	ebitenutil.DebugPrintAt(screen, "CHUNK EXPEDITION", 181, 18)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("VISITED %02d/12   BEACONS %d/5   TIME %02d", len(g.visited), g.beacons, seconds), 104, 44)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CAMERA CHUNK (%d,%d)   CACHE %d", center.x, center.y, len(g.chunks)), 105, 69)

	vector.DrawFilledRect(screen, 0, 88, screenW, 505, color.RGBA{20, 39, 54, 255}, false)
	minTileX := int(math.Floor((g.playerX-screenW/2)/tileSize)) - 1
	maxTileX := int(math.Floor((g.playerX+screenW/2)/tileSize)) + 1
	minTileY := int(math.Floor((g.playerY-252)/tileSize)) - 1
	maxTileY := int(math.Floor((g.playerY+252)/tileSize)) + 1
	minChunkX, maxChunkX := floorDiv(minTileX, chunkTiles), floorDiv(maxTileX, chunkTiles)
	minChunkY, maxChunkY := floorDiv(minTileY, chunkTiles), floorDiv(maxTileY, chunkTiles)
	drawnChunks := 0
	for cy := minChunkY; cy <= maxChunkY; cy++ {
		for cx := minChunkX; cx <= maxChunkX; cx++ {
			key := chunkPos{cx, cy}
			c := g.chunks[key]
			if c == nil {
				continue
			}
			drawnChunks++
			g.drawChunk(screen, key, c)
		}
	}
	vector.DrawFilledCircle(screen, screenW/2, 340, 12, color.RGBA{238, 95, 79, 255}, false)
	vector.StrokeCircle(screen, screenW/2, 340, 12, 3, color.White, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("DRAWN NOW %d / CACHED %d", drawnChunks, len(g.chunks)), 151, 570)

	drawButton(screen, 8, "LEFT")
	drawButton(screen, 126, "UP")
	drawButton(screen, 244, "DOWN")
	drawButton(screen, 362, "RIGHT")
	ebitenutil.DebugPrintAt(screen, g.message, 63, 687)
	if g.won {
		overlay(screen, "EXPEDITION COMPLETE!\n\nTAP / ENTER TO RETRY")
	}
	if g.lost {
		overlay(screen, "TIME UP\n\nTAP / ENTER TO RETRY")
	}
}

func (g *game) drawChunk(screen *ebiten.Image, key chunkPos, c *chunk) {
	for y := 0; y < chunkTiles; y++ {
		for x := 0; x < chunkTiles; x++ {
			worldX := float64((key.x*chunkTiles + x) * tileSize)
			worldY := float64((key.y*chunkTiles + y) * tileSize)
			px := float32(worldX - g.playerX + screenW/2)
			py := float32(worldY - g.playerY + 340)
			colors := [...]color.RGBA{{42, 111, 153, 255}, {67, 145, 83, 255}, {99, 104, 116, 255}, {241, 190, 55, 255}}
			vector.DrawFilledRect(screen, px, py, tileSize, tileSize, colors[c.tiles[y][x]], false)
			vector.StrokeRect(screen, px, py, tileSize, tileSize, 1, color.RGBA{10, 25, 35, 70}, false)
			if c.tiles[y][x] == 3 {
				vector.DrawFilledCircle(screen, px+tileSize/2, py+tileSize/2, 6, color.White, false)
			}
		}
	}
	left := float32(float64(key.x*chunkSize) - g.playerX + screenW/2)
	top := float32(float64(key.y*chunkSize) - g.playerY + 340)
	vector.StrokeRect(screen, left, top, chunkSize, chunkSize, 3, color.RGBA{245, 206, 82, 190}, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d,%d", key.x, key.y), int(left)+5, int(top)+5)
}

func movementInput() (float64, float64) {
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		dy++
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if y >= 610 {
			dx, dy = directionForButton(x)
		}
	}
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		x, y := ebiten.TouchPosition(ids[0])
		if y >= 610 {
			dx, dy = directionForButton(x)
		}
	}
	return dx, dy
}

func directionForButton(x int) (float64, float64) {
	switch {
	case x < 120:
		return -1, 0
	case x < 240:
		return 0, -1
	case x < 360:
		return 0, 1
	default:
		return 1, 0
	}
}

func drawButton(screen *ebiten.Image, x int, label string) {
	vector.DrawFilledRect(screen, float32(x), 610, 110, 58, color.RGBA{52, 83, 119, 255}, false)
	ebitenutil.DebugPrintAt(screen, label, x+36, 635)
}

func retryPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyR) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || len(ebiten.AppendTouchIDs(nil)) > 0
}

func floorDiv(n, d int) int {
	q, r := n/d, n%d
	if r < 0 {
		q--
	}
	return q
}

func positiveMod(n, d int) int {
	r := n % d
	if r < 0 {
		r += d
	}
	return r
}

func overlay(screen *ebiten.Image, text string) {
	vector.DrawFilledRect(screen, 45, 270, 390, 150, color.RGBA{4, 14, 31, 247}, false)
	vector.StrokeRect(screen, 45, 270, 390, 150, 4, color.RGBA{243, 189, 70, 255}, false)
	ebitenutil.DebugPrintAt(screen, text, 112, 327)
}

func (g *game) Layout(_, _ int) (int, int) { return screenW, screenH }

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Chunk Expedition — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
