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
var depthsArtFiles embed.FS

var (
	depthsArtOnce sync.Once
	depthsArt     map[string]*ebiten.Image
)

func prepareDepthsArt() {
	depthsArtOnce.Do(func() {
		depthsArt = map[string]*ebiten.Image{}
		for name, path := range map[string]string{
			"gardens":          "assets/depths-sunken-gardens.png",
			"abyss":            "assets/depths-clockwork-abyss.png",
			"sanctum":          "assets/depths-tempest-sanctum.png",
			"platform-gardens": "assets/platform-gardens.png",
			"platform-abyss":   "assets/platform-abyss.png",
			"platform-sanctum": "assets/platform-sanctum.png",
			"tenjiroh":         "assets/depths-tenjiroh.png",
			"beetle":           "assets/depths-beetle.png",
			"guardian":         "assets/depths-guardian.png",
			"spirit":           "assets/depths-spirit.png",
		} {
			data, err := depthsArtFiles.ReadFile(path)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			depthsArt[name] = ebiten.NewImageFromImage(decoded)
		}
	})
}

func drawDepthsRegion(dst *ebiten.Image, name string, progress float64) {
	img := depthsArt[name]
	b := img.Bounds()
	// The source is a very wide establishing shot. Scaling by height and
	// sliding across it turns one painting into a parallax region panorama.
	// Fill from the HUD to the bottom edge. Platforms are transparent overlays;
	// the environment therefore remains visible below every collision ledge.
	scale := 630 / float64(b.Dy())
	dw := float64(b.Dx()) * scale
	progress = max(0, min(1, progress))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(-progress*(dw-W), 90)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawDepthsPlatform(dst *ebiten.Image, name string, cx, top, width float64) {
	img := depthsArt["platform-"+name]
	b := img.Bounds()
	scale := width / float64(b.Dx())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(cx-float64(b.Dx())*scale/2, top)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawDepthsCharacter(dst *ebiten.Image, name string, cx, cy, w, h float64, mirror bool) {
	img := depthsArt[name]
	b := img.Bounds()
	scale := min(w/float64(b.Dx()), h/float64(b.Dy()))
	dw, dh := float64(b.Dx())*scale, float64(b.Dy())*scale
	op := &ebiten.DrawImageOptions{}
	if mirror {
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(cx+dw/2, cy-dh/2)
	} else {
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(cx-dw/2, cy-dh/2)
	}
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}
