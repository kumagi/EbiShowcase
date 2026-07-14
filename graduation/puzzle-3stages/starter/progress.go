package main

// StageData changes puzzle targets without copying the progression rule.
type StageData struct {
	Name   string
	Target int
}

var Stages = []StageData{{"Spark", 2}, {"Wave", 3}, {"Crown", 4}}

type Input struct {
	Solve   bool
	Restart bool
}
type Progress struct {
	Stage, Presses int
	Clear          bool
}

func NewProgress() Progress { // TODO 1: Start at stage 0, not clear.
	return Progress{}
}
func (p *Progress) Step(in Input) {
	if in.Restart { // TODO 2: reset with NewProgress.
		return
	}
	if p.Clear || !in.Solve {
		return
	}
	// TODO 3: count a solve press, compare it to Stages[p.Stage].Target,
	// then advance exactly one data-driven stage or set Clear after the last one.
}
func (p Progress) Current() StageData {
	if p.Stage >= len(Stages) {
		return Stages[len(Stages)-1]
	}
	return Stages[p.Stage]
}
