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
var tacticsArtFiles embed.FS

var (
	tacticsArtOnce sync.Once
	tacticsArt     map[string]*ebiten.Image
)

func prepareTacticsArt() {
	tacticsArtOnce.Do(func() {
		tacticsArt = map[string]*ebiten.Image{}
		for name, path := range map[string]string{
			"field": "assets/tactics-coastal-highlands.png",
			"blade": "assets/tactics-blade.png",
			"bow":   "assets/tactics-bow.png",
			"enemy": "assets/tactics-enemy.png",
		} {
			data, err := tacticsArtFiles.ReadFile(path)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			tacticsArt[name] = ebiten.NewImageFromImage(decoded)
		}
	})
}

func drawTacticsCover(dst *ebiten.Image, name string, x, y, w, h float64) {
	img := tacticsArt[name]
	b := img.Bounds()
	scale := max(w/float64(b.Dx()), h/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(w-float64(b.Dx())*scale)/2, y+(h-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawTacticsUnit(dst *ebiten.Image, name string, cx, cy, w, h float64) {
	img := tacticsArt[name]
	b := img.Bounds()
	scale := min(w/float64(b.Dx()), h/float64(b.Dy()))
	dw, dh := float64(b.Dx())*scale, float64(b.Dy())*scale
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(cx-dw/2, cy-dh/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}
