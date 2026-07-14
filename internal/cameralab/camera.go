// Package cameralab contains deterministic camera policy without Ebitengine.
package cameralab

type Vec struct{ X, Y float64 }
type State struct {
	Pos          Vec
	ViewW, ViewH float64
}

func (s *State) Follow(target Vec, smooth, worldW, worldH float64) {
	s.Pos.X += (target.X - s.Pos.X) * smooth
	s.Pos.Y += (target.Y - s.Pos.Y) * smooth
	minX, maxX := s.ViewW/2, worldW-s.ViewW/2
	minY, maxY := s.ViewH/2, worldH-s.ViewH/2
	if s.Pos.X < minX {
		s.Pos.X = minX
	}
	if s.Pos.X > maxX {
		s.Pos.X = maxX
	}
	if s.Pos.Y < minY {
		s.Pos.Y = minY
	}
	if s.Pos.Y > maxY {
		s.Pos.Y = maxY
	}
}
func (s *State) FollowDeadZone(target, velocity Vec, halfW, halfH, look float64) {
	want := Vec{target.X + velocity.X*look, target.Y + velocity.Y*look}
	dx, dy := want.X-s.Pos.X, want.Y-s.Pos.Y
	if dx > halfW {
		s.Pos.X += dx - halfW
	}
	if dx < -halfW {
		s.Pos.X += dx + halfW
	}
	if dy > halfH {
		s.Pos.Y += dy - halfH
	}
	if dy < -halfH {
		s.Pos.Y += dy + halfH
	}
}

func (s State) WorldToScreen(p Vec) Vec {
	return Vec{p.X - s.Pos.X + s.ViewW/2, p.Y - s.Pos.Y + s.ViewH/2}
}
func (s State) ScreenToWorld(p Vec) Vec {
	return Vec{p.X + s.Pos.X - s.ViewW/2, p.Y + s.Pos.Y - s.ViewH/2}
}
