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
var monsterArtFiles embed.FS

var (
	monsterArtOnce sync.Once
	monsterArt     map[string]*ebiten.Image
)

func prepareMonsterArt() {
	monsterArtOnce.Do(func() {
		monsterArt = map[string]*ebiten.Image{}
		paths := map[string]string{}
		for _, name := range []string{
			"expedition-map", "battle-tidepool", "battle-ember", "battle-kelp",
			"navigator", "reeflet", "mosshell", "cinderfin", "cloudray", "reef-lord",
			"capture-orb", "badge-tide", "badge-ember", "badge-kelp",
		} {
			paths[name] = "assets/" + name + ".png"
		}
		for name, path := range paths {
			data, err := monsterArtFiles.ReadFile(path)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			monsterArt[name] = ebiten.NewImageFromImage(decoded)
		}
	})
}

func drawMonsterBackground(dst *ebiten.Image, name string, x, y, w, h int) {
	img := monsterArt[name]
	b := img.Bounds()
	scale := max(float64(w)/float64(b.Dx()), float64(h)/float64(b.Dy()))
	target := dst.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((float64(w)-float64(b.Dx())*scale)/2, (float64(h)-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	target.DrawImage(img, op)
}

func drawMonsterSprite(dst *ebiten.Image, name string, centerX, centerY, size float64) {
	drawMonsterSpriteTone(dst, name, centerX, centerY, size, 1, 1)
}

func drawMonsterSpriteTone(dst *ebiten.Image, name string, centerX, centerY, size, brightness, alpha float64) {
	img := monsterArt[name]
	b := img.Bounds()
	scale := size / float64(max(b.Dx(), b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(centerX-float64(b.Dx())*scale/2, centerY-float64(b.Dy())*scale/2)
	op.ColorScale.Scale(float32(brightness), float32(brightness), float32(brightness), float32(alpha))
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}
