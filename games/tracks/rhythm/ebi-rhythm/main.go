package main

import (
	"github.com/kumagi/EbiShowcase/internal/rhythmcore"
	"github.com/kumagi/EbiShowcase/internal/rhythmplay"
)

func chart(name string, bpm int, pattern []int, hard bool) rhythmcore.Chart {
	cues := rhythmplay.Taps(pattern, map[bool]float64{false: 1, true: .5}[hard])
	if hard {
		cues = append(cues, rhythmplay.Cue{Beat: 4, Lane: 0, Kind: rhythmcore.Hold, Length: 2}, rhythmplay.Cue{Beat: 8, Lane: 3, Kind: rhythmcore.Roll, Length: 2, Need: 7})
	}
	return rhythmplay.MakeChart(name, bpm, 4, cues...)
}
func main() {
	songs := []rhythmplay.Song{{Name: "SUNRISE HARBOR", Tone: 220, Easy: chart("SUNRISE EASY", 100, []int{0, 1, 2, 3, 0, 2, 1, 3, 0, 1, 2, 3}, false), Hard: chart("SUNRISE HARD", 132, []int{0, 1, 2, 3, 0, 2, 1, 3, 0, 3, 2, 1, 0, 2, 3, 1, 0, 1, 3, 2}, true)}, {Name: "NEON REEF", Tone: 294, Easy: chart("NEON EASY", 116, []int{0, 2, 1, 3, 0, 1, 2, 3, 1, 2, 0, 3}, false), Hard: chart("NEON HARD", 148, []int{0, 2, 1, 3, 0, 1, 3, 2, 0, 3, 1, 2, 0, 2, 3, 1, 0, 1, 2, 3}, true)}, {Name: "TEMPEST PARADE", Tone: 370, Easy: chart("TEMPEST EASY", 126, []int{0, 1, 3, 2, 0, 2, 1, 3, 0, 3, 2, 1}, false), Hard: chart("TEMPEST HARD", 160, []int{0, 1, 3, 2, 0, 2, 1, 3, 0, 3, 1, 2, 0, 1, 2, 3, 0, 2, 3, 1, 0, 3, 2, 1}, true)}}
	rhythmplay.Run(rhythmplay.Config{Title: "EBI RHYTHM TOUR", Subtitle: "Three visual songs, two charts, taps, holds, and rolls.", Difficulty: true, Songs: songs})
}
