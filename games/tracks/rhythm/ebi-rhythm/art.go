// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
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
var artFS embed.FS

var artCache sync.Map

func art(name string) *ebiten.Image {
	if value, ok := artCache.Load(name); ok {
		return value.(*ebiten.Image)
	}
	b, err := artFS.ReadFile("assets/" + name + ".png")
	if err != nil {
		panic(err)
	}
	src, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	img := ebiten.NewImageFromImage(src)
	artCache.Store(name, img)
	return img
}

func drawCover(dst *ebiten.Image, img *ebiten.Image) {
	dw, dh := dst.Bounds().Dx(), dst.Bounds().Dy()
	iw, ih := img.Bounds().Dx(), img.Bounds().Dy()
	s := max(float64(dw)/float64(iw), float64(dh)/float64(ih))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s, s)
	op.GeoM.Translate((float64(dw)-float64(iw)*s)/2, (float64(dh)-float64(ih)*s)/2)
	dst.DrawImage(img, op)
}

func drawContain(dst *ebiten.Image, img *ebiten.Image, x, y, w, h float64, alpha float32) {
	iw, ih := float64(img.Bounds().Dx()), float64(img.Bounds().Dy())
	s := min(w/iw, h/ih)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s, s)
	op.GeoM.Translate(x+(w-iw*s)/2, y+(h-ih*s)/2)
	op.ColorScale.ScaleAlpha(alpha)
	dst.DrawImage(img, op)
}
