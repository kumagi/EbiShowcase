package main

import (
	"github.com/kumagi/EbiShowcase/internal/rhythmcore"
	"github.com/kumagi/EbiShowcase/internal/rhythmplay"
)

func main() {
	c := rhythmplay.MakeChart("LONG TIDE", 100, 2, rhythmplay.Cue{Beat: 0, Lane: 0}, rhythmplay.Cue{Beat: 1, Lane: 1}, rhythmplay.Cue{Beat: 2, Lane: 0, Kind: rhythmcore.Hold, Length: 2}, rhythmplay.Cue{Beat: 5, Lane: 1}, rhythmplay.Cue{Beat: 6, Lane: 1, Kind: rhythmcore.Hold, Length: 2}, rhythmplay.Cue{Beat: 9, Lane: 0}, rhythmplay.Cue{Beat: 10, Lane: 1})
	rhythmplay.Run(rhythmplay.Config{Title: "HOLD THE TIDE", Subtitle: "Press at the head and keep D or K held through the tail.", Songs: []rhythmplay.Song{{Name: "LONG TIDE", Tone: 196, Easy: c}}})
}
