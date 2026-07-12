// Package vfxsprites holds soft effect textures for fire, water, and lightning.
// Images are generated offline by cmd/gen-vfx and embedded here for WASM games.
package vfxsprites

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

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

var (
	Fire  *ebiten.Image
	Water *ebiten.Image
	Spark *ebiten.Image
	Bolt  *ebiten.Image
)

func init() {
	Fire = mustDecode(firePNG)
	Water = mustDecode(waterPNG)
	Spark = mustDecode(sparkPNG)
	Bolt = mustDecode(boltPNG)
}

func mustDecode(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
