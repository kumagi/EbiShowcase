package main

import "github.com/kumagi/EbiShowcase/internal/rhythmplay"

func main() {
	p := []int{0, 1, 2, 3, 0, 2, 1, 3, 0, 1, 3, 2, 0, 2, 1, 3}
	c := rhythmplay.MakeChart("FOUR CURRENT", 116, 4, rhythmplay.Taps(p, .5)...)
	rhythmplay.Run(rhythmplay.Config{Title: "FOUR-LANE GROOVE", Subtitle: "D F J K and four touch lanes turn one beat into a pattern.", Songs: []rhythmplay.Song{{Name: "FOUR CURRENT", Tone: 262, Easy: c}}})
}
