package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/bomber-coral-forge.png
var bomberBackgroundPNG []byte

//go:embed assets/bomber-characters-atlas.png
var bomberCharactersPNG []byte

//go:embed assets/bomber-bomb-flame-atlas.png
var bomberEffectsPNG []byte

//go:embed assets/bomber-walls-atlas.png
var bomberWallsPNG []byte

//go:embed assets/bomber-items-atlas.png
var bomberItemsPNG []byte

var (
	bomberArtOnce    sync.Once
	bomberBackground *ebiten.Image
	bomberCharacters [2]*ebiten.Image
	bomberEffects    [2]*ebiten.Image
	bomberWalls      [2]*ebiten.Image
	bomberItems      [3]*ebiten.Image
)

func loadBomberArt() {
	bomberArtOnce.Do(func() {
		bomberBackground = decodeBomberPNG(bomberBackgroundPNG)
		characters := decodeBomberPNG(bomberCharactersPNG)
		effects := decodeBomberPNG(bomberEffectsPNG)
		walls := decodeBomberPNG(bomberWallsPNG)
		items := decodeBomberPNG(bomberItemsPNG)
		for i := range bomberCharacters {
			bomberCharacters[i] = atlasCell(characters, i, len(bomberCharacters))
			bomberEffects[i] = atlasCell(effects, i, len(bomberEffects))
			bomberWalls[i] = atlasCell(walls, i, len(bomberWalls))
		}
		for i := range bomberItems {
			bomberItems[i] = atlasCell(items, i, len(bomberItems))
		}
	})
}

func decodeBomberPNG(data []byte) *ebiten.Image {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(src)
}

func atlasCell(atlas *ebiten.Image, index, count int) *ebiten.Image {
	w, h := atlas.Bounds().Dx()/count, atlas.Bounds().Dy()
	return atlas.SubImage(image.Rect(index*w, 0, (index+1)*w, h)).(*ebiten.Image)
}

func drawBomberSprite(dst *ebiten.Image, sprite *ebiten.Image, cx, cy, size float64) {
	if sprite == nil {
		return
	}
	w, h := sprite.Bounds().Dx(), sprite.Bounds().Dy()
	scale := size / float64(max(w, h))
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(cx, cy)
	dst.DrawImage(sprite, op)
}

func drawBomberCover(dst *ebiten.Image) {
	if bomberBackground == nil {
		return
	}
	sw, sh := bomberBackground.Bounds().Dx(), bomberBackground.Bounds().Dy()
	scale := mathMax(float64(screenW)/float64(sw), float64(screenH)/float64(sh))
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Translate(-float64(sw)/2, -float64(sh)/2)
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(screenW/2, screenH/2)
	dst.DrawImage(bomberBackground, op)
}

func mathMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
