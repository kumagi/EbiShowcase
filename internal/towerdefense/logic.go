// Package towerdefense contains deterministic geometry and targeting rules.
// It deliberately has no Ebitengine dependency, so game rules can be tested.
package towerdefense

import "math"

type Vec struct{ X, Y float64 }

// Scenario is the data boundary between a reusable defense rule set and a
// meaningful playable mission. UI code reads it; it does not copy wave rules.
type Scenario struct {
	Name, Goal string
	Route      []Vec
	Waves      int
	Coins      int
	Lives      int
	SpeedScale float64
	Boss       bool
}

func ValidScenario(s Scenario) bool {
	return s.Name != "" && s.Goal != "" && len(s.Route) >= 2 && NewPath(s.Route).Total > 0 && s.Waves > 0 && s.Coins >= 0 && s.Lives > 0 && s.SpeedScale > 0
}

// ResultGrade turns a completed defense into a small replay target.
func ResultGrade(score, lives, coins int) string {
	total := score + lives*100 + coins
	switch {
	case lives >= 8 && total >= 1500:
		return "S"
	case lives >= 5 && total >= 900:
		return "A"
	case lives >= 2:
		return "B"
	default:
		return "C"
	}
}

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
