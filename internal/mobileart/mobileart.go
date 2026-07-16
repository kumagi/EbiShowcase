// Package mobileart provides high-resolution original artwork for the genre
// capstones. Gameplay remains normal Go state; these images are presentation
// layers selected by Draw without mutating the game.
package mobileart

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*.png
var files embed.FS

var (
	once    sync.Once
	images  map[string]*ebiten.Image
	loadErr error
)

func load() {
	images = map[string]*ebiten.Image{}
	entries, err := fs.ReadDir(files, "assets")
	if err != nil {
		loadErr = err
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := files.ReadFile("assets/" + entry.Name())
		if err != nil {
			loadErr = err
			return
		}
		decoded, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			loadErr = fmt.Errorf("mobileart %s: %w", entry.Name(), err)
			return
		}
		name := entry.Name()[:len(entry.Name())-len(".png")]
		images[name] = ebiten.NewImageFromImage(decoded)
	}
}

// Preload decodes the embedded artwork before the game loop starts. Capstones
// call this from their constructor so Draw never performs lazy cache mutation.
func Preload() {
	once.Do(load)
	if loadErr != nil {
		panic(loadErr)
	}
}

// Get returns one immutable embedded image. Unknown names return nil.
func Get(name string) *ebiten.Image {
	Preload()
	return images[name]
}

// DrawCover fills a destination rectangle while preserving aspect ratio. Any
// overflow is clipped by the destination image or its parent viewport.
func DrawCover(dst *ebiten.Image, name string, x, y, w, h float64) {
	img := Get(name)
	if img == nil || w <= 0 || h <= 0 {
		return
	}
	b := img.Bounds()
	scale := max(w/float64(b.Dx()), h/float64(b.Dy()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x+(w-float64(b.Dx())*scale)/2, y+(h-float64(b.Dy())*scale)/2)
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}

// DrawContain fits a complete image inside a rectangle and optionally mirrors
// it. It is suitable for transparent character portraits and large props.
func DrawContain(dst *ebiten.Image, name string, x, y, w, h float64, mirror bool) {
	DrawContainAlpha(dst, name, x, y, w, h, mirror, 1)
}

// DrawContainAlpha is DrawContain with an opacity multiplier. Keeping this
// transform in Draw makes fades presentation-only; the story state is still
// advanced exclusively by Update.
func DrawContainAlpha(dst *ebiten.Image, name string, x, y, w, h float64, mirror bool, alpha float32) {
	img := Get(name)
	if img == nil || w <= 0 || h <= 0 {
		return
	}
	b := img.Bounds()
	scale := min(w/float64(b.Dx()), h/float64(b.Dy()))
	dw, dh := float64(b.Dx())*scale, float64(b.Dy())*scale
	op := &ebiten.DrawImageOptions{}
	if mirror {
		op.GeoM.Scale(-scale, scale)
		op.GeoM.Translate(x+(w+dw)/2, y+(h-dh)/2)
	} else {
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x+(w-dw)/2, y+(h-dh)/2)
	}
	op.Filter = ebiten.FilterLinear
	op.ColorScale.ScaleAlpha(alpha)
	dst.DrawImage(img, op)
}
