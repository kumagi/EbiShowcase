package main

import (
	"bytes"
	"image"
	_ "image/png"
	"testing"
)

func TestGeneratedCreatureAtlasHasSevenTransparentGutterCells(t *testing.T) {
	atlas, _, err := image.Decode(bytes.NewReader(creaturesPNG))
	if err != nil {
		t.Fatal(err)
	}
	rects := alphaColumnCells(atlas, 7)
	if len(rects) != 7 {
		t.Fatalf("alphaColumnCells() returned %d poses, want 7", len(rects))
	}
	for i, rect := range rects {
		if rect.Dx() < 80 || rect.Dy() < 80 {
			t.Fatalf("pose %d looks clipped: %v", i, rect)
		}
		if i > 0 && rect.Min.X <= rects[i-1].Max.X {
			t.Fatalf("pose %d overlaps its neighbor: %v / %v", i, rects[i-1], rect)
		}
	}
}

func TestCreatureAtlasGuttersAreActuallyTransparent(t *testing.T) {
	atlas, _, err := image.Decode(bytes.NewReader(creaturesPNG))
	if err != nil {
		t.Fatal(err)
	}
	rects := alphaColumnCells(atlas, 7)
	for i := 1; i < len(rects); i++ {
		for x := rects[i-1].Max.X; x < rects[i].Min.X; x++ {
			for y := atlas.Bounds().Min.Y; y < atlas.Bounds().Max.Y; y++ {
				_, _, _, a := atlas.At(x, y).RGBA()
				if a >= 0x1000 {
					t.Fatalf("opaque neighbor fragment in gutter %d at (%d,%d)", i, x, y)
				}
			}
		}
	}
}

func TestCreatureSpritesHaveZeroOrigin(t *testing.T) {
	loadGeneratedArt()
	for i, sprite := range creatureSprites {
		if sprite == nil || sprite.Bounds().Min.X != 0 || sprite.Bounds().Min.Y != 0 {
			t.Fatalf("creature %d retained atlas coordinates: %v", i, sprite)
		}
	}
}
