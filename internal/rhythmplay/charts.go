package rhythmplay

import (
	"math"

	"github.com/kumagi/EbiShowcase/internal/rhythmcore"
)

type Cue struct {
	Beat   float64
	Lane   int
	Kind   rhythmcore.Kind
	Length float64
	Need   int
}

func MakeChart(name string, bpm, lanes int, cues ...Cue) rhythmcore.Chart {
	framesPerBeat := 3600.0 / float64(bpm)
	notes := make([]rhythmcore.Note, 0, len(cues))
	for _, c := range cues {
		notes = append(notes, rhythmcore.Note{Lane: c.Lane, At: 90 + int(math.Round(c.Beat*framesPerBeat)), Kind: c.Kind, Duration: int(math.Round(c.Length * framesPerBeat)), Need: c.Need})
	}
	return rhythmcore.Chart{Name: name, BPM: bpm, Lanes: lanes, Notes: notes}
}

func Taps(pattern []int, spacing float64) []Cue {
	out := make([]Cue, 0, len(pattern))
	for i, lane := range pattern {
		out = append(out, Cue{Beat: float64(i) * spacing, Lane: lane})
	}
	return out
}
