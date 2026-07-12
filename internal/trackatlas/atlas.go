// Package trackatlas loads the shared 応用編 (15 genre tracks) texture atlas
// and hands out SubImage sprites. Theme: 海老天 (ebi tempura).
//
// Typical use:
//
//	trackatlas.DrawCentered(screen, "pearl", x, y, 16)
//	trackatlas.Draw(screen, "tile-wall", px, py, tileSize)
package trackatlas

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kumagi/EbiShowcase/internal/tracklayout"
)

//go:embed track-atlas.png
var atlasPNG []byte

var (
	sheet   *ebiten.Image
	sprites map[string]*ebiten.Image
)

func init() {
	img, _, err := image.Decode(bytes.NewReader(atlasPNG))
	if err != nil {
		log.Fatal(err)
	}
	sheet = ebiten.NewImageFromImage(img)
	sprites = make(map[string]*ebiten.Image, len(tracklayout.Sprites))
	for _, sp := range tracklayout.Sprites {
		x, y, w, h := tracklayout.Rect(sp.Row, sp.Col)
		sprites[sp.Name] = sheet.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
	}
}

// FrameW and FrameH are the size of one atlas cell.
const (
	FrameW = tracklayout.FrameW
	FrameH = tracklayout.FrameH
)

// Sheet returns the whole atlas image.
func Sheet() *ebiten.Image { return sheet }

// Get returns a SubImage for the named sprite, or nil if unknown.
func Get(name string) *ebiten.Image {
	return sprites[name]
}

// Draw draws the sprite with its top-left at (x, y), scaled to size×size pixels.
func Draw(dst *ebiten.Image, name string, x, y, size float64) {
	img := sprites[name]
	if img == nil || size <= 0 {
		return
	}
	b := img.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(size/float64(b.Dx()), size/float64(b.Dy()))
	op.GeoM.Translate(x, y)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

// DrawCentered draws the sprite centered on (cx, cy), scaled to size×size.
func DrawCentered(dst *ebiten.Image, name string, cx, cy, size float64) {
	Draw(dst, name, cx-size/2, cy-size/2, size)
}

// DrawTinted draws a centered sprite with a ColorScale tint (r,g,b,a in 0..1).
func DrawTinted(dst *ebiten.Image, name string, cx, cy, size float64, r, g, b, a float32) {
	img := sprites[name]
	if img == nil || size <= 0 {
		return
	}
	bb := img.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(size/float64(bb.Dx()), size/float64(bb.Dy()))
	op.GeoM.Translate(cx-size/2, cy-size/2)
	op.ColorScale.Scale(r, g, b, a)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

// Gem returns the gem sprite name for color index 0..4 (match3 / pairs).
func Gem(kind int) string {
	names := []string{"gem-red", "gem-blue", "gem-yellow", "gem-green", "gem-purple"}
	if kind < 0 || kind >= len(names) {
		return "gem-trash"
	}
	return names[kind]
}

// Merge returns the merge-tier sprite name for tiers 1..7.
func Merge(tier int) string {
	if tier < 1 {
		tier = 1
	}
	if tier > 7 {
		tier = 7
	}
	return []string{"merge-1", "merge-2", "merge-3", "merge-4", "merge-5", "merge-6", "merge-7"}[tier-1]
}

// Species returns species-0..3 or species-evo.
func Species(i int) string {
	if i < 0 {
		i = 0
	}
	if i > 4 {
		i = 4
	}
	return []string{"species-0", "species-1", "species-2", "species-3", "species-evo"}[i]
}
