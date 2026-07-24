package main

// RoundTicks is 60 seconds at Ebitengine's default 60 updates per second.
const RoundTicks = 60 * 60

// Input is deliberately small so the game rules can be tested without
// starting a window. main.go converts keyboard, pointer, and touch input to it.
type Input struct {
	Action  bool
	Restart bool
}

// Round is all state that changes while the 60-second game is running.
// TODO 1: Keep the score, elapsed ticks, and finished state here.
type Round struct {
	Score int
	Ticks int
	Over  bool
}

func NewRound() Round {
	// TODO 2: Return a fresh playable round.
	return Round{}
}

// Step owns the rules. Draw must only show the Round after this function has
// decided what happened.
func (r *Round) Step(in Input) {
	if in.Restart {
		// TODO 3: Replace the current state with NewRound().
		return
	}
	if r.Over {
		return
	}

	// TODO 4: When Action is true, add 10 score points exactly once.
	// TODO 5: Advance Ticks once per Step and set Over at RoundTicks.
}

func (r Round) SecondsLeft() int {
	left := 60 - r.Ticks/60
	if left < 0 {
		return 0
	}
	return left
}
