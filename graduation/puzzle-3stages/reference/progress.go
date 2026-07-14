package main

type StageData struct {
	Name   string
	Target int
}

var Stages = []StageData{{"Spark", 2}, {"Wave", 3}, {"Crown", 4}}

type Input struct{ Solve, Restart bool }
type Progress struct {
	Stage, Presses int
	Clear          bool
}

func NewProgress() Progress { return Progress{} }
func (p *Progress) Step(in Input) {
	if in.Restart {
		*p = NewProgress()
		return
	}
	if p.Clear || !in.Solve {
		return
	}
	p.Presses++
	if p.Presses < Stages[p.Stage].Target {
		return
	}
	p.Presses = 0
	if p.Stage == len(Stages)-1 {
		p.Clear = true
		return
	}
	p.Stage++
}
func (p Progress) Current() StageData {
	if p.Stage >= len(Stages) {
		return Stages[len(Stages)-1]
	}
	return Stages[p.Stage]
}
