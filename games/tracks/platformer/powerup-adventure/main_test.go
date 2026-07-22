package main

import (
	"math"
	"testing"
)

func TestTerrainSurfaceMapsToCollisionTop(t *testing.T) {
	for _, renderedHeight := range []float64{38, 92, 132} {
		walkY := 640.0
		drawY := terrainDrawY(walkY, renderedHeight)
		mappedSurface := drawY + terrainSurfaceRow*renderedHeight/terrainSourceHeight
		if math.Abs(mappedSurface-walkY) > 0.001 {
			t.Fatalf("height %.0f maps surface to %.3f, want %.3f", renderedHeight, mappedSurface, walkY)
		}
	}
}

func TestSolidTerrainStopsHorizontalEntry(t *testing.T) {
	grounds := []rect{{x: 100, y: 90, w: 80, h: 100}}
	player := rect{x: 94, y: 110, w: 10, h: 20}
	got, vx := resolveSolidTerrainSides(player, 88, 6, grounds)
	if got.x != 90 || vx != 0 {
		t.Fatalf("solid side resolution = (x %.1f, vx %.1f), want (90, 0)", got.x, vx)
	}
}

func TestThinLedgeDoesNotBlockHorizontalEntry(t *testing.T) {
	grounds := []rect{{x: 100, y: 90, w: 80, h: 18}}
	player := rect{x: 94, y: 92, w: 10, h: 12}
	got, vx := resolveSolidTerrainSides(player, 88, 6, grounds)
	if got.x != player.x || vx != 6 {
		t.Fatalf("one-way ledge resolution = (x %.1f, vx %.1f), want unchanged", got.x, vx)
	}
}
