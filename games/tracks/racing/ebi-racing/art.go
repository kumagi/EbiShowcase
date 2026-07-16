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
var racingArtFiles embed.FS

var (
	racingArtOnce sync.Once
	racingArt     map[string]*ebiten.Image
)

func prepareRacingArt() {
	racingArtOnce.Do(func() {
		racingArt = map[string]*ebiten.Image{}
		for name, path := range map[string]string{
			"coast":  "assets/coral-grand-prix-v2.png",
			"reef":   "assets/reef-temple.png",
			"storm":  "assets/storm-citadel.png",
			"player": "assets/player-car.png",
			"rival":  "assets/rival-car.png",
		} {
			data, err := racingArtFiles.ReadFile(path)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			racingArt[name] = ebiten.NewImageFromImage(decoded)
		}
	})
}

func drawCourseArt(dst *ebiten.Image, stage int) {
	name := []string{"coast", "reef", "storm"}[stage-1]
	img := racingArt[name]
	b := img.Bounds()
	scale := max(float64(W)/float64(b.Dx()), float64(H)/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((W-float64(b.Dx())*scale)/2, (H-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawRaceCar(dst *ebiten.Image, name string, x, y, angle, height float64) {
	img := racingArt[name]
	b := img.Bounds()
	scale := height / float64(b.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(x, y)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}
