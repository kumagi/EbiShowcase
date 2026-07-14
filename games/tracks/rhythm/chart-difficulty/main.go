package main

import "github.com/kumagi/EbiShowcase/internal/rhythmplay"

func main() {
	easy := rhythmplay.MakeChart("CHART LAB EASY", 110, 4, rhythmplay.Taps([]int{0, 1, 2, 3, 0, 1, 2, 3, 0, 2, 1, 3}, 1)...)
	hard := rhythmplay.MakeChart("CHART LAB HARD", 142, 4, rhythmplay.Taps([]int{0, 1, 2, 3, 0, 2, 1, 3, 0, 1, 3, 2, 0, 3, 1, 2, 0, 2, 3, 1, 0, 1, 2, 3}, .5)...)
	rhythmplay.Run(rhythmplay.Config{Title: "CHART DIFFICULTY", Subtitle: "The same runner reads Easy and Hard note data.", Difficulty: true, Songs: []rhythmplay.Song{{Name: "DATA DANCE", Tone: 330, Easy: easy, Hard: hard}}})
}
