package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Only this capstone's three images enter its WASM. Keeping the embed local
// avoids pulling every other genre's high-resolution artwork into the binary.
//
//go:embed assets/*.png
var activeRPGArtFiles embed.FS

var (
	activeRPGArtOnce sync.Once
	activeRPGArt     map[string]*ebiten.Image
)

func loadActiveRPGArt() {
	activeRPGArt = map[string]*ebiten.Image{}
	for name, path := range map[string]string{
		"arena":    "assets/active-rpg-moonlit-arena.png",
		"tenjiroh": "assets/active-rpg-tenjiroh.png",
		"storm":    "assets/active-rpg-storm-king.png",
		"mage":     "assets/active-rpg-mage.png",
		"shell":    "assets/active-rpg-shell.png",
		"wisp":     "assets/active-rpg-wisp.png",
		"scout":    "assets/active-rpg-scout.png",
	} {
		data, err := activeRPGArtFiles.ReadFile(path)
		if err != nil {
			panic(err)
		}
		decoded, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		activeRPGArt[name] = ebiten.NewImageFromImage(decoded)
	}
}

func activeRPGImage(name string) *ebiten.Image {
	activeRPGArtOnce.Do(loadActiveRPGArt)
	return activeRPGArt[name]
}

// prepareActiveRPGArt completes the one-time decoding before the game loop.
// Draw then only reads immutable images and remains a pure projection of state.
func prepareActiveRPGArt() { activeRPGArtOnce.Do(loadActiveRPGArt) }

func drawActiveRPGCover(dst *ebiten.Image, name string, x, y, w, h float64) {
	img := activeRPGImage(name)
	b := img.Bounds()
	scale := max(w/float64(b.Dx()), h/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(w-float64(b.Dx())*scale)/2, y+(h-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawActiveRPGContain(dst *ebiten.Image, name string, x, y, w, h, scalePulse float64) {
	img := activeRPGImage(name)
	b := img.Bounds()
	scale := min(w/float64(b.Dx()), h/float64(b.Dy())) * scalePulse
	dw, dh := float64(b.Dx())*scale, float64(b.Dy())*scale
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(w-dw)/2, y+(h-dh)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}
