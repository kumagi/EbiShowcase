// Command gen-vfx draws soft effect sprites (fire / water / spark / bolt) with
// pure-Go software rendering and writes them for both WASM embed and the site.
//
//	internal/vfxsprites/{fire,water,spark,bolt,ice,light,dark,ring}.png
//	web/assets/vfx-*.png
//
// Run: go run ./cmd/gen-vfx
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

func main() {
	root, err := os.Getwd()
	must(err)
	fire := drawFire()
	r, g, b, a := fire.At(32, 70).RGBA()
	fmt.Printf("fire base sample #%02x%02x%02x a=%d\n", r>>8, g>>8, b>>8, a>>8)
	writeBoth(root, "fire.png", "vfx-fire.png", fire)
	writeBoth(root, "water.png", "vfx-water.png", drawWater())
	writeBoth(root, "spark.png", "vfx-spark.png", drawSpark())
	writeBoth(root, "bolt.png", "vfx-bolt.png", drawBolt())
	writeBoth(root, "ice.png", "vfx-ice.png", drawIce())
	writeBoth(root, "light.png", "vfx-light.png", drawLight())
	writeBoth(root, "dark.png", "vfx-dark.png", drawDark())
	writeBoth(root, "ring.png", "vfx-ring.png", drawRing())
	println("Generated vfx sprites: fire, water, spark, bolt, ice, light, dark, ring")
}

func writeBoth(root, embedName, webName string, img *image.RGBA) {
	writePNG(filepath.Join(root, "internal/vfxsprites", embedName), img)
	writePNG(filepath.Join(root, "web/assets", webName), img)
}

func writePNG(path string, img image.Image) {
	must(os.MkdirAll(filepath.Dir(path), 0o755))
	f, err := os.Create(path)
	must(err)
	defer f.Close()
	must(png.Encode(f, img))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func lerp(a, b, t float64) float64 { return a + (b-a)*t }

func mixRGBA(a, b color.RGBA, t float64) color.RGBA {
	t = clamp01(t)
	return color.RGBA{
		R: uint8(lerp(float64(a.R), float64(b.R), t)),
		G: uint8(lerp(float64(a.G), float64(b.G), t)),
		B: uint8(lerp(float64(a.B), float64(b.B), t)),
		A: uint8(lerp(float64(a.A), float64(b.A), t)),
	}
}

// drawFire — tip at top, yellow/white core at base → orange body → red fringe.
func drawFire() *image.RGBA {
	const W, H = 72, 112
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	core := color.RGBA{255, 245, 170, 255}
	mid := color.RGBA{255, 140, 28, 255}
	tip := color.RGBA{210, 35, 12, 255}
	outer := color.RGBA{120, 12, 4, 255}
	cx := float64(W) / 2
	for y := 0; y < H; y++ {
		ty := float64(y) / float64(H-1) // 0 = tip (top), 1 = base (bottom)
		wave := math.Sin(ty*math.Pi*3.6+0.5) * 0.08 * ty
		half := (0.10 + ty*0.38 + wave) * float64(W)
		for x := 0; x < W; x++ {
			dx := (float64(x) + 0.5 - cx) / math.Max(half, 1)
			if math.Abs(dx) > 1 {
				continue
			}
			edge := 1 - dx*dx
			heat := edge * (0.2 + 0.8*ty)
			var col color.RGBA
			switch {
			case heat > 0.78:
				col = mixRGBA(mid, core, (heat-0.78)/0.22)
			case heat > 0.42:
				col = mixRGBA(tip, mid, (heat-0.42)/0.36)
			case heat > 0.18:
				col = mixRGBA(outer, tip, (heat-0.18)/0.24)
			default:
				col = outer
			}
			a := math.Pow(edge, 1.35) * (0.2 + 0.8*ty)
			// Soft tip fade.
			a *= 0.25 + 0.75*math.Sqrt(ty)
			col.A = uint8(255 * clamp01(a))
			img.SetRGBA(x, y, col)
		}
	}
	return img
}

// drawWater — teardrop (point up) with cyan body and specular highlight.
func drawWater() *image.RGBA {
	const W, H = 56, 72
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	body := color.RGBA{90, 190, 255, 255}
	deep := color.RGBA{25, 90, 200, 255}
	hi := color.RGBA{230, 250, 255, 255}
	cx, cy := float64(W)*0.5, float64(H)*0.62
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			fx := (float64(x) + 0.5 - cx) / (float64(W) * 0.40)
			fy := (float64(y) + 0.5 - cy) / (float64(H) * 0.38)
			top := clamp01((cy - (float64(y) + 0.5)) / (cy - 2))
			fx *= 1 + top*1.6
			d := fx*fx + fy*fy
			if d > 1 {
				continue
			}
			edge := 1 - d
			col := mixRGBA(deep, body, edge)
			hx := (float64(x) + 0.5 - cx + 7) / 9
			hy := (float64(y) + 0.5 - cy + 10) / 9
			spec := math.Exp(-(hx*hx + hy*hy) * 1.5)
			col = mixRGBA(col, hi, spec*0.9)
			col.A = uint8(255 * clamp01(edge*1.2))
			img.SetRGBA(x, y, col)
		}
	}
	return img
}

// drawSpark — soft radial glow (ember / splash seed). White so ColorScale can tint.
func drawSpark() *image.RGBA {
	const R = 28
	size := R * 2
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			d := math.Hypot(float64(x-R)+0.5, float64(y-R)+0.5) / R
			if d >= 1 {
				continue
			}
			a := (1 - d) * (1 - d)
			core := math.Exp(-d * d * 6)
			r := 255
			g := 255
			b := 255
			_ = core
			img.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(255 * a)})
		}
	}
	return img
}

// drawBolt — vertical jagged lightning with bright core + cyan halo.
func drawBolt() *image.RGBA {
	const W, H = 56, 140
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	pts := make([][2]float64, 0, 14)
	x := float64(W) / 2
	for y := 6.0; y < H-6; y += 12 {
		pts = append(pts, [2]float64{x, y})
		x += (hash(int(y))*2 - 1) * 12
		if x < 12 {
			x = 12
		}
		if x > W-12 {
			x = W - 12
		}
	}
	pts = append(pts, [2]float64{float64(W) / 2, float64(H - 6)})

	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			px, py := float64(x)+0.5, float64(y)+0.5
			d := distToPoly(px, py, pts)
			if d > 12 {
				continue
			}
			core := math.Exp(-d * d * 0.7)
			halo := math.Exp(-d * d * 0.07)
			r := 200 + 55*core
			g := 220 + 35*core
			b := 255.0
			a := clamp01(halo*0.5 + core*0.95)
			img.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(255 * a)})
		}
	}
	return img
}

// drawIce — hexagonal-ish crystal shard, cyan/white.
func drawIce() *image.RGBA {
	const W, H = 64, 96
	img := image.NewRGBA(image.Rect(0, 0, W, H))
	core := color.RGBA{230, 250, 255, 255}
	mid := color.RGBA{120, 210, 255, 255}
	edgeC := color.RGBA{40, 120, 220, 255}
	cx, cy := float64(W)*0.5, float64(H)*0.55
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			fx := (float64(x) + 0.5 - cx) / (float64(W) * 0.28)
			fy := (float64(y) + 0.5 - cy) / (float64(H) * 0.42)
			// diamond / shard: |fx|+|fy|
			d := math.Abs(fx) + math.Abs(fy)*0.85
			if d > 1 {
				continue
			}
			edge := 1 - d
			col := mixRGBA(edgeC, mid, edge)
			if edge > 0.65 {
				col = mixRGBA(mid, core, (edge-0.65)/0.35)
			}
			// facet line
			if math.Abs(fx) < 0.08 {
				col = mixRGBA(col, core, 0.55)
			}
			col.A = uint8(255 * clamp01(edge*1.15))
			img.SetRGBA(x, y, col)
		}
	}
	return img
}

// drawLight — soft golden star/flare.
func drawLight() *image.RGBA {
	const R = 40
	size := R * 2
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x-R) + 0.5
			dy := float64(y-R) + 0.5
			d := math.Hypot(dx, dy) / R
			if d >= 1 {
				continue
			}
			ang := math.Atan2(dy, dx)
			ray := math.Pow(math.Abs(math.Cos(ang*4)), 8) // 8-point soft star
			a := (1 - d) * (1 - d) * (0.35 + 0.65*ray)
			img.SetRGBA(x, y, color.RGBA{255, 236, 180, uint8(255 * clamp01(a))})
		}
	}
	return img
}

// drawDark — soft purple/black wisp, denser center.
func drawDark() *image.RGBA {
	const R = 36
	size := R * 2
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			d := math.Hypot(float64(x-R)+0.5, float64(y-R)+0.5) / R
			if d >= 1 {
				continue
			}
			a := math.Pow(1-d, 1.6)
			// dark core with violet fringe
			t := d
			r := uint8(20 + 80*t)
			g := uint8(0 + 20*t)
			b := uint8(40 + 120*t)
			img.SetRGBA(x, y, color.RGBA{r, g, b, uint8(255 * a)})
		}
	}
	return img
}

// drawRing — soft annular glow for shockwaves.
func drawRing() *image.RGBA {
	const R = 64
	size := R * 2
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			d := math.Hypot(float64(x-R)+0.5, float64(y-R)+0.5) / R
			band := math.Abs(d - 0.72)
			if band > 0.22 {
				continue
			}
			a := math.Exp(-band*band*80) * (1 - d*0.3)
			img.SetRGBA(x, y, color.RGBA{255, 255, 255, uint8(255 * clamp01(a))})
		}
	}
	return img
}

func hash(n int) float64 {
	x := uint32(n)*2654435761 ^ 0x9e3779b9
	x ^= x >> 16
	return float64(x&0xffff) / 65535
}

func distToPoly(px, py float64, pts [][2]float64) float64 {
	best := 1e9
	for i := 1; i < len(pts); i++ {
		d := distSeg(px, py, pts[i-1][0], pts[i-1][1], pts[i][0], pts[i][1])
		if d < best {
			best = d
		}
	}
	return best
}

func distSeg(px, py, ax, ay, bx, by float64) float64 {
	vx, vy := bx-ax, by-ay
	wx, wy := px-ax, py-ay
	c1 := wx*vx + wy*vy
	if c1 <= 0 {
		return math.Hypot(px-ax, py-ay)
	}
	c2 := vx*vx + vy*vy
	if c2 <= c1 {
		return math.Hypot(px-bx, py-by)
	}
	t := c1 / c2
	return math.Hypot(px-(ax+t*vx), py-(ay+t*vy))
}
