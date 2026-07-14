package main

import "testing"

// TODO 1 and TODO 3 in state.go: the first room contains the key.
func TestKeyCanBeCollectedOnlyInFirstRoom(t *testing.T) {
	adventure := NewAdventure()
	adventure.Step(Input{PickUpKey: true})
	if !adventure.HasKey {
		t.Fatal("picking up the room-1 key must add it to inventory")
	}
}

// TODO 4 in state.go: rooms are a bounded state machine.
func TestPathMovesThroughExactlyThreeRooms(t *testing.T) {
	adventure := State{Room: 1}
	adventure.Step(Input{NextRoom: true})
	adventure.Step(Input{NextRoom: true})
	adventure.Step(Input{NextRoom: true})
	if adventure.Room != 3 {
		t.Fatalf("room = %d, want 3", adventure.Room)
	}
}

// TODO 5 in state.go: the exit is a key-gated transition.
func TestExitNeedsKeyInThirdRoom(t *testing.T) {
	withoutKey := State{Room: 3}
	withoutKey.Step(Input{OpenExit: true})
	if withoutKey.Escaped {
		t.Fatal("the exit must stay locked without a key")
	}
	withKey := State{Room: 3, HasKey: true}
	withKey.Step(Input{OpenExit: true})
	if !withKey.Escaped {
		t.Fatal("the key must open the room-3 exit")
	}
}

// TODO 2 in state.go: a restart must erase every old state value.
func TestRestartMakesNewAdventure(t *testing.T) {
	adventure := State{Room: 3, HasKey: true, Escaped: true}
	adventure.Step(Input{Restart: true})
	if adventure != (State{Room: 1}) {
		t.Fatalf("restart = %+v, want a new room-1 adventure", adventure)
	}
}
