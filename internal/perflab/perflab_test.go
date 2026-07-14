package perflab

import "testing"

func sample(n int) []Point {
	r := make([]Point, n)
	for i := range r {
		r[i] = Point{float64(i % 1000), float64(i / 1000)}
	}
	return r
}
func TestCullAndGrid(t *testing.T) {
	all := sample(100)
	got := Cull(make([]Point, 0, 20), all, Rect{10, 0, 10, 1})
	if len(got) != 10 {
		t.Fatalf("got %d", len(got))
	}
	g := NewGrid(10)
	for _, p := range all {
		g.Add(p)
	}
	if len(g.Nearby(Point{15, 0})) == 0 {
		t.Fatal("no candidates")
	}
}
func BenchmarkCullReuse(b *testing.B) {
	all := sample(10000)
	dst := make([]Point, 0, 1000)
	view := Rect{300, 0, 200, 10}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst = Cull(dst, all, view)
	}
	_ = dst
}
func BenchmarkGridNearby(b *testing.B) {
	all := sample(10000)
	g := NewGrid(32)
	for _, p := range all {
		g.Add(p)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.Nearby(Point{500, 5})
	}
}
