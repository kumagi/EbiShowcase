package main

import (
	"github.com/kumagi/EbiShowcase/internal/rhythmcore"
	"github.com/kumagi/EbiShowcase/internal/rhythmplay"
)

func main() {
	c := rhythmplay.MakeChart("SHELL ROLL", 124, 2, rhythmplay.Cue{Beat: 0, Lane: 0}, rhythmplay.Cue{Beat: 1, Lane: 1}, rhythmplay.Cue{Beat: 2, Lane: 0, Kind: rhythmcore.Roll, Length: 2, Need: 6}, rhythmplay.Cue{Beat: 5, Lane: 1, Kind: rhythmcore.Roll, Length: 2, Need: 7}, rhythmplay.Cue{Beat: 8, Lane: 0}, rhythmplay.Cue{Beat: 8.5, Lane: 1})
	rhythmplay.Run(rhythmplay.Config{Title: "DRUM ROLL REEF", Subtitle: "Repeatedly tap inside each pink roll window.", Songs: []rhythmplay.Song{{Name: "SHELL ROLL", Tone: 294, Easy: c}}})
}
