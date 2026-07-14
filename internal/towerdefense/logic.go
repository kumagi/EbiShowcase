// Package towerdefense contains deterministic geometry and targeting rules.
// It deliberately has no Ebitengine dependency, so game rules can be tested.
package towerdefense

import "math"

type Vec struct{ X, Y float64 }

func Distance(a, b Vec) float64 { return math.Hypot(a.X-b.X, a.Y-b.Y) }

type Path struct {
	Points  []Vec
	lengths []float64
	Total   float64
}

func NewPath(points []Vec) Path {
	p := Path{Points: append([]Vec(nil), points...)}
	for i := 1; i < len(points); i++ {
		l := Distance(points[i-1], points[i])
		p.lengths = append(p.lengths, l)
		p.Total += l
	}
	return p
}

// Position converts distance travelled into a point on the waypoint path.
func (p Path) Position(progress float64) Vec {
	if len(p.Points) == 0 {
		return Vec{}
	}
	if progress <= 0 {
		return p.Points[0]
	}
	for i, length := range p.lengths {
		if progress <= length {
			t := progress / length
			a, b := p.Points[i], p.Points[i+1]
			return Vec{a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t}
		}
		progress -= length
	}
	return p.Points[len(p.Points)-1]
}

type Target struct {
	Pos      Vec
	Progress float64
	Alive    bool
}

// SelectFront returns the in-range target furthest along the route.
func SelectFront(origin Vec, radius float64, targets []Target) int {
	best, bestProgress := -1, -1.0
	for i, target := range targets {
		if target.Alive && Distance(origin, target.Pos) <= radius && target.Progress > bestProgress {
			best, bestProgress = i, target.Progress
		}
	}
	return best
}
