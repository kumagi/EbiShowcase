// Package hero loads the shared Ebi Showcase protagonist sprite.
package hero

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed ebi-boy.png
var tenjirohPNG []byte

var sprite *ebiten.Image

func init() {
	img, _, err := image.Decode(bytes.NewReader(tenjirohPNG))
	if err != nil {
		log.Fatal(err)
	}
	sprite = ebiten.NewImageFromImage(img)
}

// Image returns the shared protagonist sprite (full body, transparent PNG).
func Image() *ebiten.Image {
	return sprite
}

// DrawCentered draws the hero scaled so its height is heightPx,
// centered on (cx, cy) in screen coordinates.
func DrawCentered(dst *ebiten.Image, cx, cy, heightPx float64) {
	if sprite == nil || heightPx <= 0 {
		return
	}
	b := sprite.Bounds()
	sw := float64(b.Dx())
	sh := float64(b.Dy())
	scale := heightPx / sh
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(cx-sw*scale/2, cy-sh*scale/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(sprite, op)
}

// DrawBottomCentered draws the hero standing on (footX, footY),
// scaled to heightPx.
func DrawBottomCentered(dst *ebiten.Image, footX, footY, heightPx float64) {
	if sprite == nil || heightPx <= 0 {
		return
	}
	b := sprite.Bounds()
	sw := float64(b.Dx())
	sh := float64(b.Dy())
	scale := heightPx / sh
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(footX-sw*scale/2, footY-sh*scale)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(sprite, op)
}
