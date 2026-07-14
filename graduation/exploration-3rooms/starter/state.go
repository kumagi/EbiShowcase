package main

// Input is the small, testable meaning of a button press. main.go translates
// keyboard, pointer, or touch input to this value.
type Input struct {
	PickUpKey bool
	NextRoom  bool
	OpenExit  bool
	Restart   bool
}

// State is the whole adventure. Keep the room, inventory, and ending state
// together so every transition can be tested without a window.
type State struct {
	Room    int
	HasKey  bool
	Escaped bool
}

func NewAdventure() State {
	// TODO 1: Start in room 1 with no key and no ending.
	return State{}
}

func (s *State) Step(in Input) {
	if in.Restart {
		// TODO 2: Replace the state with NewAdventure().
		return
	}
	if s.Escaped {
		return
	}

	// TODO 3: Pick up the key only in room 1.
	// TODO 4: Move from room 1 → 2 → 3, never beyond room 3.
	// TODO 5: Open the exit only in room 3 and only with the key.
}
