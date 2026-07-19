package main

import (
	"bytes"
	"embed"
	"image"
	imagedraw "image/draw"
	_ "image/png"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*.png
var fightingArtFiles embed.FS

var (
	fightingArtOnce sync.Once
	fightingArt     map[string]*ebiten.Image
	fightingMotion  map[string][]*ebiten.Image
)

func prepareFightingArt() {
	fightingArtOnce.Do(func() {
		fightingArt = map[string]*ebiten.Image{}
		fightingMotion = map[string][]*ebiten.Image{}
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
		for _, side := range []string{"player", "rival"} {
			data, err := fightingArtFiles.ReadFile("assets/" + side + "-motion.png")
			if err != nil {
				panic(err)
			}
			decoded, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				panic(err)
			}
			inset := 0
			if side == "rival" {
				// The generated rival sheet contains thin white cell dividers.
				// Keep them outside every runtime frame.
				inset = 5
			}
			frames := make([]*ebiten.Image, 0, 8)
			for frame := 0; frame < 8; frame++ {
				rect := motionFrameRect(decoded.Bounds(), frame, inset)
				rect = visibleAlphaRect(decoded, rect, 8)
				cropped := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
				imagedraw.Draw(cropped, cropped.Bounds(), decoded, rect.Min, imagedraw.Src)
				frames = append(frames, ebiten.NewImageFromImage(cropped))
			}
			fightingMotion[side] = frames
		}
	})
}

func motionFrameRect(bounds image.Rectangle, frame, inset int) image.Rectangle {
	column, row := frame%4, frame/4
	minX := bounds.Min.X + bounds.Dx()*column/4 + inset
	maxX := bounds.Min.X + bounds.Dx()*(column+1)/4 - inset
	minY := bounds.Min.Y + bounds.Dy()*row/2 + inset
	maxY := bounds.Min.Y + bounds.Dy()*(row+1)/2 - inset
	return image.Rect(minX, minY, maxX, maxY)
}

func visibleAlphaRect(img image.Image, cell image.Rectangle, apron int) image.Rectangle {
	minX, minY, maxX, maxY := cell.Max.X, cell.Max.Y, cell.Min.X, cell.Min.Y
	for y := cell.Min.Y; y < cell.Max.Y; y++ {
		for x := cell.Min.X; x < cell.Max.X; x++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha == 0 {
				continue
			}
			minX = min(minX, x)
			minY = min(minY, y)
			maxX = max(maxX, x+1)
			maxY = max(maxY, y+1)
		}
	}
	if minX >= maxX || minY >= maxY {
		panic("fighting motion frame has no visible pixels")
	}
	return image.Rect(minX-apron, minY-apron, maxX+apron, maxY+apron).Intersect(cell)
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
	drawFighterImage(dst, fightingArt[name], centerX, groundY, height, 1)
}

func drawFighterMotionFrame(dst *ebiten.Image, side string, frame int, centerX, groundY, height, alpha float64) {
	frames := fightingMotion[side]
	drawFighterImage(dst, frames[max(0, min(len(frames)-1, frame))], centerX, groundY, height, alpha)
}

func drawFighterImage(dst, img *ebiten.Image, centerX, groundY, height, alpha float64) {
	b := img.Bounds()
	scale := height / float64(b.Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(centerX-float64(b.Dx())*scale/2, groundY-float64(b.Dy())*scale)
	op.ColorScale.ScaleAlpha(float32(alpha))
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(img, op)
}
