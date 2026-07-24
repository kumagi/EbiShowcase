package main

const RoundTicks = 60 * 60

type Input struct {
	Action  bool
	Restart bool
}

type Round struct {
	Score int
	Ticks int
	Over  bool
}

func NewRound() Round { return Round{} }

func (r *Round) Step(in Input) {
	if in.Restart {
		*r = NewRound()
		return
	}
	if r.Over {
		return
	}
	if in.Action {
		r.Score += 10
	}
	r.Ticks++
	if r.Ticks >= RoundTicks {
		r.Over = true
	}
}

func (r Round) SecondsLeft() int {
	left := 60 - r.Ticks/60
	if left < 0 {
		return 0
	}
	return left
}
