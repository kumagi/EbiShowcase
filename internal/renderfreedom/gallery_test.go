package renderfreedom

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type fakeGame struct {
	updates int
	w, h    int
}

func (f *fakeGame) Update() error      { f.updates++; return nil }
func (f *fakeGame) Draw(*ebiten.Image) {}
func (f *fakeGame) Layout(w, h int) (int, int) {
	f.w, f.h = w+7, h+9
	return f.w, f.h
}

func TestGalleryDelegatesUpdateAndLayout(t *testing.T) {
	inner := &fakeGame{}
	g := &gallery{inner: inner}
	if err := g.Update(); err != nil || inner.updates != 1 {
		t.Fatalf("Update delegation = (%v, %d), want (nil, 1)", err, inner.updates)
	}
	w, h := g.Layout(100, 200)
	if w != 107 || h != 209 {
		t.Fatalf("Layout delegation = (%d, %d), want (107, 209)", w, h)
	}
}

func TestWireframeShaderCompiles(t *testing.T) {
	if _, err := ebiten.NewShader(wireframeSource); err != nil {
		t.Fatalf("wireframe shader: %v", err)
	}
}
