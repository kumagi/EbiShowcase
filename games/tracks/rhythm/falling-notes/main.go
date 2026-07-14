package main

import "github.com/kumagi/EbiShowcase/internal/rhythmplay"

func main() {
	c := rhythmplay.MakeChart("SHRIMP STEP", 108, 1, rhythmplay.Taps([]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, .5)...)
	rhythmplay.Run(rhythmplay.Config{Title: "FALLING NOTES", Subtitle: "Read travel time: hit when each note crosses the gold line.", Songs: []rhythmplay.Song{{Name: "SHRIMP STEP", Tone: 247, Easy: c}}})
}
