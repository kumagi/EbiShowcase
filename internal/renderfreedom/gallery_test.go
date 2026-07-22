package renderfreedom

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type fakeGame struct {
	updates, draws     int
	w, h, drawW, drawH int
}

func (f *fakeGame) Update() error { f.updates++; return nil }
func (f *fakeGame) Draw(screen *ebiten.Image) {
	f.draws++
	f.drawW, f.drawH = screen.Bounds().Dx(), screen.Bounds().Dy()
}
func (f *fakeGame) Layout(w, h int) (int, int) {
	f.w, f.h = w+7, h+9
	return f.w, f.h
}

func TestGalleryDrawsOriginalOnceAtNativeSize(t *testing.T) {
	screen := ebiten.NewImage(480, 720)
	inner := &fakeGame{}
	g := &gallery{inner: inner}
	g.Draw(screen)
	if inner.draws != 1 || inner.drawW != 480 || inner.drawH != 720 {
		t.Fatalf("Draw calls/native size = (%d, %d, %d), want (1, 480, 720)", inner.draws, inner.drawW, inner.drawH)
	}
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
