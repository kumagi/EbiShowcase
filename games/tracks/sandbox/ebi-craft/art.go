package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*.png
var craftArtFiles embed.FS

var (
	craftArtOnce sync.Once
	craftArt     map[string]*ebiten.Image
)

func prepareCraftArt() {
	craftArtOnce.Do(func() {
		craftArt = map[string]*ebiten.Image{}
		paths := map[string]string{
			"island-moss":    "assets/island-moss.png",
			"island-crystal": "assets/island-crystal.png",
			"island-ember":   "assets/island-ember.png",
		}
		for _, name := range []string{"hero-idle", "hero-mine", "pickaxe", "wood", "stone", "crystal", "beacon", "crawler", "workshop"} {
			paths[name] = "assets/" + name + ".png"
		}
		for name, path := range paths {
			data, err := craftArtFiles.ReadFile(path)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			craftArt[name] = ebiten.NewImageFromImage(decoded)
		}
	})
}

func drawIslandArt(dst *ebiten.Image, stage int) {
	name := []string{"island-moss", "island-crystal", "island-ember"}[stage]
	img := craftArt[name]
	b := img.Bounds()
	scale := max(float64(width)/float64(b.Dx()), float64(height)/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((width-float64(b.Dx())*scale)/2, (height-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawCraftSprite(dst *ebiten.Image, name string, centerX, centerY, size float64) {
	img := craftArt[name]
	b := img.Bounds()
	scale := size / float64(max(b.Dx(), b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(centerX-float64(b.Dx())*scale/2, centerY-float64(b.Dy())*scale/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}
