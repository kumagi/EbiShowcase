// Package rhythmcore contains deterministic, audio-independent rhythm rules.
// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
package rhythmcore

import "slices"

type Kind uint8

const (
	Tap Kind = iota
	Hold
	Roll
)

type Note struct {
	Lane     int
	At       int
	Kind     Kind
	Duration int
	Need     int
}

type Chart struct {
	Name  string
	BPM   int
	Lanes int
	Notes []Note
}

type Grade string

const (
	Perfect Grade = "PERFECT"
	Good    Grade = "GREAT"
	Miss    Grade = "MISS"
)

type Input struct {
	Lane int
	Down bool
}

type Result struct {
	Note  int
	Grade Grade
	Lane  int
	Delta int
}

type noteState struct {
	started, resolved bool
	startGrade        Grade
	delta, hits       int
}

type Session struct {
	Chart                   Chart
	Frame                   int
	Score, Combo, Best      int
	Perfects, Goods, Misses int
	// Offset shifts input timing in frames. Positive values compensate for an
	// input path that arrives late; it is deliberately part of pure state.
	Offset int
	held   []bool
	states []noteState
	events []Result
}

const PerfectWindow = 4
const GoodWindow = 9

func NewSession(chart Chart) *Session {
	return NewSessionWithOffset(chart, 0)
}

func NewSessionWithOffset(chart Chart, offset int) *Session {
	c := chart
	c.Notes = slices.Clone(chart.Notes)
	slices.SortFunc(c.Notes, func(a, b Note) int { return a.At - b.At })
	return &Session{Chart: c, Offset: offset, held: make([]bool, max(1, c.Lanes)), states: make([]noteState, len(c.Notes))}
}

func GradeDelta(delta int) Grade {
	if delta < 0 {
		delta = -delta
	}
	if delta <= PerfectWindow {
		return Perfect
	}
	if delta <= GoodWindow {
		return Good
	}
	return Miss
}

func (s *Session) timingFrame() int { return s.Frame + s.Offset }

func (s *Session) Step(inputs []Input) []Result {
	s.events = s.events[:0]
	for _, in := range inputs {
		if in.Lane < 0 || in.Lane >= len(s.held) {
			continue
		}
		s.held[in.Lane] = in.Down
		if in.Down {
			s.press(in.Lane)
		} else {
			s.release(in.Lane)
		}
	}
	for i, n := range s.Chart.Notes {
		st := &s.states[i]
		if st.resolved {
			continue
		}
		switch n.Kind {
		case Tap:
			if s.timingFrame() > n.At+GoodWindow {
				s.resolve(i, Miss, s.timingFrame()-n.At)
			}
		case Hold:
			if !st.started && s.timingFrame() > n.At+GoodWindow {
				s.resolve(i, Miss, s.timingFrame()-n.At)
			}
			if st.started && s.timingFrame() >= n.At+n.Duration {
				if s.held[n.Lane] {
					s.resolve(i, st.startGrade, st.delta)
				} else {
					s.resolve(i, Miss, st.delta)
				}
			}
		case Roll:
			if s.timingFrame() > n.At+n.Duration {
				grade := Miss
				if st.hits >= max(1, n.Need) {
					grade = Perfect
				} else if st.hits >= max(1, n.Need-2) {
					grade = Good
				}
				s.resolve(i, grade, st.hits)
			}
		}
	}
	s.Frame++
	return slices.Clone(s.events)
}

func (s *Session) press(lane int) {
	best, bestDelta := -1, GoodWindow+1
	for i, n := range s.Chart.Notes {
		st := &s.states[i]
		if st.resolved || n.Lane != lane {
			continue
		}
		if n.Kind == Roll && s.timingFrame() >= n.At-GoodWindow && s.timingFrame() <= n.At+n.Duration {
			st.hits++
			return
		}
		if st.started {
			continue
		}
		d := s.Frame + s.Offset - n.At
		ad := d
		if ad < 0 {
			ad = -ad
		}
		if ad <= GoodWindow && ad < bestDelta {
			best, bestDelta = i, ad
		}
	}
	if best < 0 {
		return
	}
	n, st := s.Chart.Notes[best], &s.states[best]
	delta := s.timingFrame() - n.At
	if n.Kind == Hold {
		st.started, st.startGrade, st.delta = true, GradeDelta(delta), delta
		return
	}
	s.resolve(best, GradeDelta(delta), delta)
}

func (s *Session) release(lane int) {
	for i, n := range s.Chart.Notes {
		st := &s.states[i]
		if n.Kind == Hold && n.Lane == lane && st.started && !st.resolved && s.Frame < n.At+n.Duration {
			s.resolve(i, Miss, st.delta)
		} else if n.Kind == Hold && n.Lane == lane && st.started && !st.resolved {
			s.resolve(i, st.startGrade, st.delta)
		}
	}
}

func (s *Session) resolve(i int, grade Grade, delta int) {
	st := &s.states[i]
	if st.resolved {
		return
	}
	st.resolved = true
	switch grade {
	case Perfect:
		s.Perfects++
		s.Combo++
		s.Score += 100 + s.Combo*2
	case Good:
		s.Goods++
		s.Combo++
		s.Score += 50 + s.Combo
	default:
		s.Misses++
		s.Combo = 0
	}
	if s.Combo > s.Best {
		s.Best = s.Combo
	}
	s.events = append(s.events, Result{Note: i, Grade: grade, Lane: s.Chart.Notes[i].Lane, Delta: delta})
}

func (s *Session) Resolved(i int) bool { return i >= 0 && i < len(s.states) && s.states[i].resolved }
func (s *Session) Started(i int) bool  { return i >= 0 && i < len(s.states) && s.states[i].started }
func (s *Session) RollHits(i int) int {
	if i < 0 || i >= len(s.states) {
		return 0
	}
	return s.states[i].hits
}
func (s *Session) Finished() bool {
	if len(s.Chart.Notes) == 0 {
		return true
	}
	for i := range s.states {
		if !s.states[i].resolved {
			return false
		}
	}
	return s.Frame > s.Chart.Notes[len(s.Chart.Notes)-1].At+2
}
