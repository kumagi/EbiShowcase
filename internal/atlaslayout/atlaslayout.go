// Package atlaslayout is the single source of truth for the 海老・天次郎
// (Ebi Tenjiroh) texture atlas: frame size and the ordered list of animations
// (action + facing).
//
// It has no dependencies so both the offline generator (cmd/gen-atlas) and the
// runtime loader (internal/heroatlas) can share the exact same layout.
package atlaslayout

// Frame size and grid width, in pixels.
const (
	FrameW = 96
	FrameH = 96
	Cols   = 4
)

// Anim is one animation strip: Frames cells starting at column 0 of Row.
type Anim struct {
	Name   string // e.g. "walk-side"
	Action string // idle, walk, run, attack, hurt
	Dir    string // down, up, side  (left = draw side flipped horizontally)
	Row    int
	Frames int
	FPS    int
}

// Anims lists every strip in row order. Row equals the slice index.
var Anims = []Anim{
	{"idle-down", "idle", "down", 0, 2, 4},
	{"idle-up", "idle", "up", 1, 2, 4},
	{"idle-side", "idle", "side", 2, 2, 4},
	{"walk-down", "walk", "down", 3, 4, 10},
	{"walk-up", "walk", "up", 4, 4, 10},
	{"walk-side", "walk", "side", 5, 4, 10},
	{"run-down", "run", "down", 6, 4, 14},
	{"run-up", "run", "up", 7, 4, 14},
	{"run-side", "run", "side", 8, 4, 14},
	{"attack-down", "attack", "down", 9, 3, 14},
	{"attack-up", "attack", "up", 10, 3, 14},
	{"attack-side", "attack", "side", 11, 3, 14},
	{"hurt-down", "hurt", "down", 12, 2, 8},
	{"hurt-up", "hurt", "up", 13, 2, 8},
	{"hurt-side", "hurt", "side", 14, 2, 8},
}

// SheetW and SheetH are the full atlas pixel dimensions.
func SheetW() int { return Cols * FrameW }
func SheetH() int { return len(Anims) * FrameH }

// Rect returns the pixel rectangle (x, y, w, h) of one frame.
func Rect(row, col int) (x, y, w, h int) {
	return col * FrameW, row * FrameH, FrameW, FrameH
}

// Find returns the animation with the given name, and whether it was found.
func Find(name string) (Anim, bool) {
	for _, a := range Anims {
		if a.Name == name {
			return a, true
		}
	}
	return Anim{}, false
}
