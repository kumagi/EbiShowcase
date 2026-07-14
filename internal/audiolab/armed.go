// Package audiolab separates browser-safe user consent from sound generation.
package audiolab

// Gate is a tiny state machine: no audio work begins until an explicit game
// input arms it. It is pure Go so autoplay behaviour can be tested.
type Gate struct{ armed bool }

func (g *Gate) Arm(gesture bool) bool {
	if gesture {
		g.armed = true
	}
	return g.armed
}
func (g *Gate) Armed() bool { return g.armed }
