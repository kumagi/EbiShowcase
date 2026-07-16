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
var fightingArtFiles embed.FS

var (
	fightingArtOnce sync.Once
	fightingArt     map[string]*ebiten.Image
)

func prepareFightingArt() {
	fightingArtOnce.Do(func() {
		fightingArt = map[string]*ebiten.Image{}
		paths := map[string]string{"arena": "assets/moon-tide-arena.png"}
		for _, side := range []string{"player", "rival"} {
			for _, pose := range []string{"ready", "attack", "hurt", "ko"} {
				paths[side+"-"+pose] = "assets/" + side + "-" + pose + ".png"
			}
		}
		for name, path := range paths {
			data, err := fightingArtFiles.ReadFile(path)
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			fightingArt[name] = ebiten.NewImageFromImage(decoded)
		}
	})
}

func drawFightingArena(dst *ebiten.Image) {
	img := fightingArt["arena"]
	b := img.Bounds()
	scale := max(float64(width)/float64(b.Dx()), float64(height)/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((width-float64(b.Dx())*scale)/2, (height-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

func drawFighterPose(dst *ebiten.Image, name string, centerX, groundY, height float64) {
	img := fightingArt[name]
	b := img.Bounds()
	scale := height / float64(b.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(centerX-float64(b.Dx())*scale/2, groundY-float64(b.Dy())*scale)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}
