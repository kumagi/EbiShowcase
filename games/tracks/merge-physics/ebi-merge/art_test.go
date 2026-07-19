package main

import (
	"bytes"
	"image"
	_ "image/png"
	"math"
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

func TestCreatureCropsKeepTransparentApron(t *testing.T) {
	atlas, _, err := image.Decode(bytes.NewReader(creaturesPNG))
	if err != nil {
		t.Fatal(err)
	}
	for i, rect := range alphaColumnCells(atlas, 7) {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			for _, y := range []int{rect.Min.Y, rect.Max.Y - 1} {
				_, _, _, a := atlas.At(x, y).RGBA()
				if a != 0 {
					t.Fatalf("creature %d alpha touches horizontal crop edge at (%d,%d)", i, x, y)
				}
			}
		}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for _, x := range []int{rect.Min.X, rect.Max.X - 1} {
				_, _, _, a := atlas.At(x, y).RGBA()
				if a != 0 {
					t.Fatalf("creature %d alpha touches vertical crop edge at (%d,%d)", i, x, y)
				}
			}
		}
	}
}

func TestMergeRadiiRoughlyPreserveArea(t *testing.T) {
	for tier := 1; tier < len(radii); tier++ {
		ratio := radii[tier] * radii[tier] / (radii[tier-1] * radii[tier-1])
		if math.Abs(ratio-2) > .12 {
			t.Fatalf("tier %d area ratio = %.3f, want approximately 2", tier, ratio)
		}
	}
}

func TestEveryMergeScoresAndHigherTiersScoreMore(t *testing.T) {
	previous := 0
	for tier := 1; tier <= maxTier; tier++ {
		points := mergePoints(tier, 1)
		if points <= previous {
			t.Fatalf("tier %d points = %d, previous = %d", tier, points, previous)
		}
		previous = points
	}
}
