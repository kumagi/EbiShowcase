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
	leftWall      = 35.0
	rightWall     = 445.0
	topWall       = 115.0
	bottomWall    = 555.0
	radius        = 15.0
	launchSpeed   = 8.0
	restitution   = 0.92
	shotLimit     = 420
	maxMisses     = 6
)

type mission struct {
	angle   int
	bounces int
	targetX float64
	targetY float64
}

type game struct {
	missions              []mission
	missionIndex          int
	angle                 int
	x, y, vx, vy          float64
	bounces, shotFrames   int
	misses                int
	flying, cleared, over bool
	message               string
}

func newGame() *game {
	g := &game{}
	for i, angle := range []int{35, 65, 120} {
		required := i + 1
		tx, ty := targetFor(angle, required)
		g.missions = append(g.missions, mission{angle: angle, bounces: required, targetX: tx, targetY: ty})
	}
	g.resetAim()
	g.message = "Match the angle and hit the pearl after the exact bounce count."
	return g
}

func targetFor(angle, required int) (float64, float64) {
	x, y := 240.0, 520.0
	radians := float64(angle) * math.Pi / 180
	vx, vy := math.Cos(radians)*launchSpeed, -math.Sin(radians)*launchSpeed
	bounces, after := 0, 0
	for frame := 0; frame < 3000; frame++ {
		x, y = x+vx, y+vy
		vx, vy, x, y, bounces = reflectCircle(vx, vy, x, y, bounces)
		if bounces == required {
			after++
			if after == 18 {
				return x, y
			}
		} else if bounces > required {
			break
		}
	}
	return 240, 300
}

func (g *game) resetAim() {
	g.x, g.y = 240, 520
	g.vx, g.vy = 0, 0
	g.angle = 45
	g.bounces, g.shotFrames = 0, 0
	g.flying = false
}

func (g *game) Update() error {
	if g.cleared || g.over {
		if retryPressed() {
			*g = *newGame()
		}
		return nil
	}
	if !g.flying {
		minus := inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA)
		plus := inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD)
		launch := inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
		if x, y, ok := pressPosition(); ok && y >= 615 {
			switch {
			case x < 155:
				minus = true
			case x < 325:
				launch = true
			default:
				plus = true
			}
		}
		if minus {
			g.angle -= 5
			if g.angle < 20 {
				g.angle = 160
			}
		}
		if plus {
			g.angle += 5
			if g.angle > 160 {
				g.angle = 20
			}
		}
		if launch {
			radians := float64(g.angle) * math.Pi / 180
			g.vx = math.Cos(radians) * launchSpeed
			g.vy = -math.Sin(radians) * launchSpeed
			g.flying = true
			g.message = "Watch which velocity component flips at each wall."
		}
		return nil
	}

	g.x, g.y = g.x+g.vx, g.y+g.vy
	g.vx, g.vy, g.x, g.y, g.bounces = reflectCircle(g.vx, g.vy, g.x, g.y, g.bounces)
	g.shotFrames++
	m := g.missions[g.missionIndex]
	dx, dy := g.x-m.targetX, g.y-m.targetY
	if g.bounces == m.bounces && dx*dx+dy*dy <= 24*24 {
		if g.missionIndex == len(g.missions)-1 {
			g.cleared = true
			g.message = "Three exact reflection paths completed!"
			return nil
		}
		g.missionIndex++
		g.resetAim()
		g.message = "Direct hit! The next pearl needs one more bounce."
		return nil
	}
	if g.bounces > m.bounces || g.shotFrames >= shotLimit {
		g.misses++
		if g.misses >= maxMisses {
			g.over = true
			g.message = "Six paths missed. Reset and compare the angles again."
			return nil
		}
		g.resetAim()
		g.message = fmt.Sprintf("Wrong path. Misses %d/%d — adjust by 5 degrees.", g.misses, maxMisses)
	}
	return nil
}

func reflectCircle(vx, vy, x, y float64, bounces int) (float64, float64, float64, float64, int) {
	if x-radius < leftWall {
		x = leftWall + radius
		vx = math.Abs(vx) * restitution
		bounces++
	} else if x+radius > rightWall {
		x = rightWall - radius
		vx = -math.Abs(vx) * restitution
		bounces++
	}
	if y-radius < topWall {
		y = topWall + radius
		vy = math.Abs(vy) * restitution
		bounces++
	} else if y+radius > bottomWall {
		y = bottomWall - radius
		vy = -math.Abs(vy) * restitution
		bounces++
	}
	return vx, vy, x, y, bounces
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{11, 22, 40, 255})
	ebitenutil.DebugPrintAt(screen, "PEARL BOUNCE LAB", 184, 25)
	m := g.missions[g.missionIndex]
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("MISSION %d/3   HIT AFTER EXACTLY %d BOUNCE(S)", g.missionIndex+1, m.bounces), 93, 53)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ANGLE %03d DEG   BOUNCES %d   MISSES %d/%d", g.angle, g.bounces, g.misses, maxMisses), 107, 77)
	vector.StrokeRect(screen, leftWall, topWall, rightWall-leftWall, bottomWall-topWall, 5, color.RGBA{76, 151, 180, 255}, false)
	vector.DrawFilledCircle(screen, float32(m.targetX), float32(m.targetY), 24, color.RGBA{239, 178, 67, 80}, false)
	vector.StrokeCircle(screen, float32(m.targetX), float32(m.targetY), 24, 4, color.RGBA{255, 209, 113, 255}, false)
	if !g.flying {
		g.drawPreview(screen)
	}
	vector.DrawFilledCircle(screen, float32(g.x), float32(g.y), radius, color.RGBA{238, 245, 248, 255}, false)
	vector.StrokeCircle(screen, float32(g.x), float32(g.y), radius, 3, color.RGBA{89, 192, 220, 255}, false)
	if g.flying {
		vector.StrokeLine(screen, float32(g.x), float32(g.y), float32(g.x+g.vx*5), float32(g.y+g.vy*5), 3, color.RGBA{255, 221, 125, 255}, false)
	}
	ebitenutil.DebugPrintAt(screen, g.message, 57, 583)
	button(screen, 20, 620, 130, "-5 DEG", color.RGBA{49, 91, 126, 255})
	button(screen, 165, 620, 150, "LAUNCH", color.RGBA{230, 168, 60, 255})
	button(screen, 330, 620, 130, "+5 DEG", color.RGBA{49, 91, 126, 255})
	if g.cleared || g.over {
		title := "REFLECTION MASTERED!"
		if g.over {
			title = "BOUNCE LAB CLOSED!"
		}
		vector.DrawFilledRect(screen, 42, 270, 396, 160, color.RGBA{4, 14, 29, 247}, false)
		ebitenutil.DebugPrintAt(screen, title, 158, 315)
		ebitenutil.DebugPrintAt(screen, g.message, 79, 350)
		ebitenutil.DebugPrintAt(screen, "TAP / SPACE TO RETRY", 148, 390)
	}
}

func (g *game) drawPreview(screen *ebiten.Image) {
	radians := float64(g.angle) * math.Pi / 180
	x, y := g.x, g.y
	vx, vy := math.Cos(radians)*launchSpeed, -math.Sin(radians)*launchSpeed
	bounces := 0
	for i := 0; i < 90; i++ {
		nx, ny := x+vx, y+vy
		vx, vy, nx, ny, bounces = reflectCircle(vx, vy, nx, ny, bounces)
		if i%5 < 3 {
			vector.StrokeLine(screen, float32(x), float32(y), float32(nx), float32(ny), 2, color.RGBA{146, 189, 205, 125}, false)
		}
		x, y = nx, ny
	}
}

func button(screen *ebiten.Image, x, y, w int, label string, fill color.RGBA) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), 66, fill, false)
	ebitenutil.DebugPrintAt(screen, label, x+w/2-len(label)*3, y+29)
}

func pressPosition() (int, int, bool) {
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
	ebiten.SetWindowTitle("Pearl Bounce Lab — Ebitengine")
	if err := ebiten.RunGame(newGame()); err != nil {
		panic(err)
	}
}
