package main

import (
	"bytes"
	"image"
	_ "image/png"
	"testing"
)

func TestGeneratedCreatureAtlasFindsSevenWholePoses(t *testing.T) {
	atlas, _, err := image.Decode(bytes.NewReader(creaturesPNG))
	if err != nil {
		t.Fatal(err)
	}
	rects := opaqueCells(atlas, 7)
	if len(rects) != 7 {
		t.Fatalf("opaqueCells() returned %d poses, want 7", len(rects))
	}
	for i, rect := range rects {
		if rect.Dx() < 80 || rect.Dy() < 80 {
			t.Fatalf("pose %d looks clipped: %v", i, rect)
		}
		if i > 0 && rect.Min.X+rect.Max.X <= rects[i-1].Min.X+rects[i-1].Max.X {
			t.Fatalf("pose centers are out of order at %d: %v / %v", i, rects[i-1], rect)
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
