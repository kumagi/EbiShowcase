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

// Pose changes only how the hero is drawn. Collision boxes and gameplay
// positions must continue to use their own state.
type Pose struct {
	ScaleX   float64
	ScaleY   float64
	Rotation float64
	Alpha    float32
}

func normalizedPose(p Pose) Pose {
	if p.ScaleX == 0 {
		p.ScaleX = 1
	}
	if p.ScaleY == 0 {
		p.ScaleY = 1
	}
	if p.Alpha == 0 {
		p.Alpha = 1
	}
	return p
}

// DrawCentered draws the hero scaled so its height is heightPx,
// centered on (cx, cy) in screen coordinates.
func DrawCentered(dst *ebiten.Image, cx, cy, heightPx float64) {
	DrawCenteredPose(dst, cx, cy, heightPx, Pose{})
}

// DrawCenteredPose draws around a visual center with squash, stretch, and
// rotation. The supplied center remains stable while the pose changes.
func DrawCenteredPose(dst *ebiten.Image, cx, cy, heightPx float64, pose Pose) {
	if sprite == nil || heightPx <= 0 {
		return
	}
	pose = normalizedPose(pose)
	b := sprite.Bounds()
	sw := float64(b.Dx())
	sh := float64(b.Dy())
	scale := heightPx / sh
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-sw/2, -sh/2)
	op.GeoM.Scale(scale*pose.ScaleX, scale*pose.ScaleY)
	op.GeoM.Rotate(pose.Rotation)
	op.GeoM.Translate(cx, cy)
	op.ColorScale.ScaleAlpha(pose.Alpha)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(sprite, op)
}

// DrawBottomCentered draws the hero standing on (footX, footY),
// scaled to heightPx.
func DrawBottomCentered(dst *ebiten.Image, footX, footY, heightPx float64) {
	DrawBottomCenteredPose(dst, footX, footY, heightPx, Pose{})
}

// DrawBottomCenteredPose keeps the feet anchored while applying a visual pose.
func DrawBottomCenteredPose(dst *ebiten.Image, footX, footY, heightPx float64, pose Pose) {
	if sprite == nil || heightPx <= 0 {
		return
	}
	pose = normalizedPose(pose)
	b := sprite.Bounds()
	sw := float64(b.Dx())
	sh := float64(b.Dy())
	scale := heightPx / sh
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-sw/2, -sh)
	op.GeoM.Scale(scale*pose.ScaleX, scale*pose.ScaleY)
	op.GeoM.Rotate(pose.Rotation)
	op.GeoM.Translate(footX, footY)
	op.ColorScale.ScaleAlpha(pose.Alpha)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(sprite, op)
}
