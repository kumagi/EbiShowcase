// Package heroatlas loads the shared 海老・天次郎 (Ebi Tenjiroh) texture atlas
// and hands out per-frame sub-images for each animation (walk, run, attack,
// hurt, idle) and facing (down, up, side). Left-facing is the side frames
// drawn flipped.
package heroatlas

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/internal/atlaslayout"
)

//go:embed ebi-boy-atlas.png
var atlasPNG []byte

var sheet *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(atlasPNG))
	if err != nil {
		log.Fatal(err)
	}
	sheet = ebiten.NewImageFromImage(img)
}

// FrameW and FrameH are the size of one frame.
const (
	FrameW = atlaslayout.FrameW
	FrameH = atlaslayout.FrameH
)

// Sheet returns the whole atlas image (for showing the sheet itself).
func Sheet() *ebiten.Image { return sheet }

// frame returns the sub-image at (row, col).
func frame(row, col int) *ebiten.Image {
	x, y, w, h := atlaslayout.Rect(row, col)
	return sheet.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
}

// Anim returns the ordered frames for a named animation such as "walk-side".
// It returns nil if the name is unknown.
func Anim(name string) []*ebiten.Image {
	a, ok := atlaslayout.Find(name)
	if !ok {
		return nil
	}
	out := make([]*ebiten.Image, a.Frames)
	for i := 0; i < a.Frames; i++ {
		out[i] = frame(a.Row, i)
	}
	return out
}

// FPS returns the suggested frames-per-second for a named animation.
func FPS(name string) int {
	if a, ok := atlaslayout.Find(name); ok {
		return a.FPS
	}
	return 8
}
