package main

type Input struct{ PickUpKey, NextRoom, OpenExit, Restart bool }
type State struct {
	Room            int
	HasKey, Escaped bool
}

func NewAdventure() State { return State{Room: 1} }
func (s *State) Step(in Input) {
	if in.Restart {
		*s = NewAdventure()
		return
	}
	if s.Escaped {
		return
	}
	if in.PickUpKey && s.Room == 1 {
		s.HasKey = true
	}
	if in.NextRoom && s.Room < 3 {
		s.Room++
	}
	if in.OpenExit && s.Room == 3 && s.HasKey {
		s.Escaped = true
	}
}
