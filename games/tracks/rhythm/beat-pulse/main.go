package main

import "github.com/kumagi/EbiShowcase/internal/rhythmplay"

func main() {
	c := rhythmplay.MakeChart("FIRST PULSE", 90, 1, rhythmplay.Taps([]int{0, 0, 0, 0, 0, 0, 0, 0}, 1)...)
	rhythmplay.Run(rhythmplay.Config{Title: "BEAT PULSE", Subtitle: "Tap SPACE or the lane when the pulse reaches the line.", Songs: []rhythmplay.Song{{Name: "FIRST PULSE", Tone: 220, Easy: c}}})
}
