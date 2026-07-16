package main

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestReflectCircleSeparatesAndReflects(t *testing.T) {
	g := game{}
	a := ally{pos: vec{20, 0}, velocity: vec{-4, 0}}
	impact := g.reflectCircle(&a, vec{0, 0}, 24)
	if impact != 4 {
		t.Fatalf("impact = %v, want 4", impact)
	}
	wantDistance := allyRadius + 24
	if got := math.Hypot(a.pos.x, a.pos.y); math.Abs(got-wantDistance) > 1e-9 {
		t.Fatalf("separation = %v, want %v", got, wantDistance)
	}
	if a.velocity.x <= 0 {
		t.Fatalf("velocity did not reflect: %+v", a.velocity)
	}
}

func TestAliveCountsOnlyTargetsWithHP(t *testing.T) {
	g := game{enemies: []enemy{{hp: 2}, {hp: 0}, {hp: -1}, {hp: 3}}}
	if got := g.alive(); got != 2 {
		t.Fatalf("alive = %d, want 2", got)
	}
}

func TestDragPositionMovesTheHeroAndClampsPull(t *testing.T) {
	origin := vec{240, 500}
	got := dragPosition(origin, vec{40, 500})
	if math.Abs(distance(origin, got)-145) > 1e-9 {
		t.Fatalf("drag distance = %v, want 145", distance(origin, got))
	}
	if got == origin {
		t.Fatal("dragging must move the visible hero, not only its shadow")
	}
}

func TestAtlasSpritesHaveZeroOrigin(t *testing.T) {
	loadStrikeArt()
	for family, sprites := range map[string][]*ebiten.Image{
		"allies":    strikeAllies[:],
		"enemies":   strikeEnemies[:],
		"obstacles": strikeObstacles[:],
	} {
		for i, sprite := range sprites {
			if sprite.Bounds().Min.X != 0 || sprite.Bounds().Min.Y != 0 {
				t.Fatalf("%s sprite %d retained atlas origin: %v", family, i, sprite.Bounds())
			}
		}
	}
}
