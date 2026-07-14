package main

import "testing"

func TestKeyCanBeCollectedOnlyInFirstRoom(t *testing.T) {
	s := NewAdventure()
	s.Step(Input{PickUpKey: true})
	if !s.HasKey {
		t.Fatal("key")
	}
}
func TestPathMovesThroughExactlyThreeRooms(t *testing.T) {
	s := State{Room: 1}
	for range 3 {
		s.Step(Input{NextRoom: true})
	}
	if s.Room != 3 {
		t.Fatal(s.Room)
	}
}
func TestExitNeedsKeyInThirdRoom(t *testing.T) {
	s := State{Room: 3, HasKey: true}
	s.Step(Input{OpenExit: true})
	if !s.Escaped {
		t.Fatal("exit")
	}
}
func TestRestartMakesNewAdventure(t *testing.T) {
	s := State{Room: 3, HasKey: true, Escaped: true}
	s.Step(Input{Restart: true})
	if s != (State{Room: 1}) {
		t.Fatal(s)
	}
}
