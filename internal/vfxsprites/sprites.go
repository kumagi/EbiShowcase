// Package vfxsprites holds generated soft-particle textures for Visual Effects Lab.
// Regenerate with: go run ./cmd/gen-vfx
package vfxsprites

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed fire.png
var firePNG []byte

//go:embed water.png
var waterPNG []byte

//go:embed spark.png
var sparkPNG []byte

//go:embed bolt.png
var boltPNG []byte

//go:embed ice.png
var icePNG []byte

//go:embed light.png
var lightPNG []byte

//go:embed dark.png
var darkPNG []byte

//go:embed ring.png
var ringPNG []byte

var (
	Fire  *ebiten.Image
	Water *ebiten.Image
	Spark *ebiten.Image
	Bolt  *ebiten.Image
	Ice   *ebiten.Image
	Light *ebiten.Image
	Dark  *ebiten.Image
	Ring  *ebiten.Image
)

func init() {
	Fire = mustDecode(firePNG)
	Water = mustDecode(waterPNG)
	Spark = mustDecode(sparkPNG)
	Bolt = mustDecode(boltPNG)
	Ice = mustDecode(icePNG)
	Light = mustDecode(lightPNG)
	Dark = mustDecode(darkPNG)
	Ring = mustDecode(ringPNG)
}

func mustDecode(png []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(png))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}
