// Package perflab contains small, benchmarkable examples used by the
// performance guide. They are deliberately independent of Ebitengine.
package perflab

type Point struct{ X, Y float64 }
type Rect struct{ X, Y, W, H float64 }

func (r Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.X+r.W && p.Y >= r.Y && p.Y < r.Y+r.H
}

// Cull appends only visible points into dst so a caller can reuse its backing
// array every tick.
func Cull(dst []Point, all []Point, view Rect) []Point {
	dst = dst[:0]
	for _, p := range all {
		if view.Contains(p) {
			dst = append(dst, p)
		}
	}
	return dst
}

type Pool struct{ free []Point }

func (p *Pool) Acquire() Point {
	if len(p.free) == 0 {
		return Point{}
	}
	v := p.free[len(p.free)-1]
	p.free = p.free[:len(p.free)-1]
	return v
}
func (p *Pool) Release(v Point) { p.free = append(p.free, v) }

// Grid narrows collision candidates before precise shape tests.
type Grid struct {
	Cell  float64
	cells map[[2]int][]Point
}

func NewGrid(cell float64) *Grid { return &Grid{Cell: cell, cells: map[[2]int][]Point{}} }
func (g *Grid) Add(p Point) {
	k := [2]int{int(p.X / g.Cell), int(p.Y / g.Cell)}
	g.cells[k] = append(g.cells[k], p)
}
func (g *Grid) Nearby(p Point) []Point {
	result := []Point{}
	cx, cy := int(p.X/g.Cell), int(p.Y/g.Cell)
	for y := cy - 1; y <= cy+1; y++ {
		for x := cx - 1; x <= cx+1; x++ {
			result = append(result, g.cells[[2]int{x, y}]...)
		}
	}
	return result
}
