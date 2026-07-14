// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
// Package raycastlogic contains the pure geometry used by the ray-casting
// lessons. It deliberately has no Ebitengine dependency so learners can test
// the rules without opening a game window.
package raycastlogic

import "math"

// Mission is data for one replayable ray-cast encounter. Keeping it here
// makes map validity, grading, and future editors testable without Ebitengine.
type Mission struct {
	Name       string
	Grid       [][]int
	StartX     float64
	StartY     float64
	StartAngle float64
	KeyX       float64
	KeyY       float64
	ExitX      float64
	ExitY      float64
	Enemies    []Point
	GoalTime   int
}

type Point struct{ X, Y float64 }

// ValidateMission rejects data that would trap a player in an invalid grid.
func ValidateMission(m Mission) bool {
	if len(m.Grid) < 3 || len(m.Grid[0]) < 3 || m.GoalTime <= 0 {
		return false
	}
	w := len(m.Grid[0])
	for _, row := range m.Grid {
		if len(row) != w {
			return false
		}
	}
	inside := func(x, y float64) bool {
		ix, iy := int(x), int(y)
		return iy >= 0 && iy < len(m.Grid) && ix >= 0 && ix < w && m.Grid[iy][ix] == 0
	}
	return inside(m.StartX, m.StartY) && inside(m.KeyX, m.KeyY) && inside(m.ExitX, m.ExitY)
}

// Grade rewards a quick, careful clear. It intentionally has no presentation
// dependencies so it can be unit-tested and shown as a tiny Go rule.
func Grade(seconds, damage, shots, enemyCount int) string {
	if damage == 0 && seconds <= 45 && shots <= enemyCount+2 {
		return "S"
	}
	if damage <= 1 && seconds <= 75 {
		return "A"
	}
	if damage <= 3 {
		return "B"
	}
	return "C"
}

type Hit struct {
	Distance float64
	MapX     int
	MapY     int
	Side     int // 0: crossed an x grid line, 1: crossed a y grid line
	WallU    float64
}

func WrapAngle(angle float64) float64 {
	for angle <= -math.Pi {
		angle += 2 * math.Pi
	}
	for angle > math.Pi {
		angle -= 2 * math.Pi
	}
	return angle
}

func CorrectDistance(distance, rayAngle, playerAngle float64) float64 {
	return distance * math.Cos(rayAngle-playerAngle)
}

func ProjectHeight(viewHeight int, distance float64) int {
	if distance < 0.001 {
		distance = 0.001
	}
	h := int(float64(viewHeight) / distance)
	if h > viewHeight*3 {
		return viewHeight * 3
	}
	return h
}

// Cast advances through a grid one cell boundary at a time (DDA) and returns
// the first non-zero cell. A missing border is treated as a distant wall.
func Cast(grid [][]int, x, y, angle float64) Hit {
	dx, dy := math.Cos(angle), math.Sin(angle)
	mapX, mapY := int(x), int(y)
	deltaX, deltaY := math.Inf(1), math.Inf(1)
	if math.Abs(dx) > 1e-9 {
		deltaX = math.Abs(1 / dx)
	}
	if math.Abs(dy) > 1e-9 {
		deltaY = math.Abs(1 / dy)
	}
	stepX, stepY := 1, 1
	sideX := (float64(mapX+1) - x) * deltaX
	sideY := (float64(mapY+1) - y) * deltaY
	if dx < 0 {
		stepX = -1
		sideX = (x - float64(mapX)) * deltaX
	}
	if dy < 0 {
		stepY = -1
		sideY = (y - float64(mapY)) * deltaY
	}

	side := 0
	for i := 0; i < 128; i++ {
		if sideX < sideY {
			sideX += deltaX
			mapX += stepX
			side = 0
		} else {
			sideY += deltaY
			mapY += stepY
			side = 1
		}
		if mapY < 0 || mapY >= len(grid) || mapX < 0 || mapX >= len(grid[mapY]) {
			return Hit{Distance: 99, MapX: mapX, MapY: mapY, Side: side}
		}
		if grid[mapY][mapX] == 0 {
			continue
		}
		distance := sideX - deltaX
		if side == 1 {
			distance = sideY - deltaY
		}
		wall := y + distance*dy
		if side == 1 {
			wall = x + distance*dx
		}
		wall -= math.Floor(wall)
		return Hit{Distance: distance, MapX: mapX, MapY: mapY, Side: side, WallU: wall}
	}
	return Hit{Distance: 99}
}

type Projection struct {
	Depth   float64
	ScreenX float64 // -1 is left edge, +1 is right edge
}

func ProjectSprite(playerX, playerY, playerAngle, worldX, worldY, fov float64) Projection {
	dx, dy := worldX-playerX, worldY-playerY
	depth := dx*math.Cos(playerAngle) + dy*math.Sin(playerAngle)
	side := -dx*math.Sin(playerAngle) + dy*math.Cos(playerAngle)
	if depth <= 0 {
		return Projection{Depth: depth, ScreenX: 99}
	}
	return Projection{Depth: depth, ScreenX: side / (depth * math.Tan(fov/2))}
}
