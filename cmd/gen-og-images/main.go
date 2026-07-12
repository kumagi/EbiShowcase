// Command gen-og-images renders per-page Open Graph images (1200×630)
// from web/assets/og/manifest.json (produced by scripts/inject-ogp.mjs).
//
//	go run ./cmd/gen-og-images
package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	ogW = 1200
	ogH = 630
)

type page struct {
	File        string `json:"file"`
	Path        string `json:"path"`
	Key         string `json:"key"`
	Lang        string `json:"lang"`
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	H1          string `json:"h1"`
	Eyebrow     string `json:"eyebrow"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type manifest struct {
	Origin string `json:"origin"`
	Pages  []page `json:"pages"`
}

func loadFaceFile(path string, size float64) (font.Face, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// TTC collections: try collection index 0 via sfnt if needed.
	f, err := opentype.Parse(raw)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
}

func trySystemFaces(size float64) font.Face {
	candidates := []string{
		"/System/Library/Fonts/Supplemental/Arial Unicode.ttf",
		"/Library/Fonts/Arial Unicode.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Bold.ttc",
		"/usr/share/fonts/opentype/noto/NotoSansCJK-Bold.ttc",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
	}
	for _, p := range candidates {
		if face, err := loadFaceFile(p, size); err == nil {
			return face
		}
	}
	return nil
}

func mustFace(ttf []byte, size float64) font.Face {
	f, err := opentype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		panic(err)
	}
	return face
}

func mix(a, b color.RGBA, t float64) color.RGBA {
	f := func(x, y uint8) uint8 { return uint8(float64(x)*(1-t) + float64(y)*t) }
	return color.RGBA{f(a.R, b.R), f(a.G, b.G), f(a.B, b.B), 255}
}

func hashColor(s string) color.RGBA {
	var h uint32
	for i := 0; i < len(s); i++ {
		h = h*33 + uint32(s[i])
	}
	// Keep hues in the site's teal / coral / indigo family.
	palette := []color.RGBA{
		{46, 201, 174, 255},
		{141, 123, 255, 255},
		{255, 138, 92, 255},
		{74, 144, 226, 255},
		{245, 199, 75, 255},
		{230, 90, 120, 255},
		{90, 190, 120, 255},
		{100, 120, 200, 255},
	}
	return palette[int(h)%len(palette)]
}

func fillGrad(img *image.RGBA, c0, c1 color.RGBA) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		t := float64(y-b.Min.Y) / float64(b.Dy())
		col := mix(c0, c1, t)
		for x := b.Min.X; x < b.Max.X; x++ {
			img.SetRGBA(x, y, col)
		}
	}
}

func fillRect(img *image.RGBA, r image.Rectangle, col color.RGBA) {
	draw.Draw(img, r, &image.Uniform{col}, image.Point{}, draw.Src)
}

func fillCircle(img *image.RGBA, cx, cy, rad int, col color.RGBA) {
	for y := cy - rad; y <= cy+rad; y++ {
		for x := cx - rad; x <= cx+rad; x++ {
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= rad*rad {
				if x >= 0 && y >= 0 && x < ogW && y < ogH {
					img.SetRGBA(x, y, col)
				}
			}
		}
	}
}

func drawString(img *image.RGBA, face font.Face, s string, x, y int, col color.RGBA) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(s)
}

func measure(face font.Face, s string) int {
	return font.MeasureString(face, s).Ceil()
}

func wrap(face font.Face, s string, maxW int) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	// Prefer word wrap when spaces exist; otherwise wrap by rune for CJK.
	if strings.Contains(s, " ") {
		words := strings.Fields(s)
		var lines []string
		cur := words[0]
		for _, w := range words[1:] {
			trial := cur + " " + w
			if measure(face, trial) <= maxW {
				cur = trial
				continue
			}
			lines = append(lines, cur)
			cur = w
		}
		lines = append(lines, cur)
		return lines
	}
	var lines []string
	var cur []rune
	for _, r := range s {
		trial := string(append(append([]rune{}, cur...), r))
		if len(cur) > 0 && measure(face, trial) > maxW {
			lines = append(lines, string(cur))
			cur = []rune{r}
			continue
		}
		cur = append(cur, r)
	}
	if len(cur) > 0 {
		lines = append(lines, string(cur))
	}
	return lines
}

func mostlyPrintable(s string) bool {
	letters := 0
	ok := 0
	for _, r := range s {
		if unicode.IsSpace(r) {
			continue
		}
		letters++
		if r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune("—–-·|/:,.'\"!?", r)) {
			ok++
		}
	}
	if letters == 0 {
		return true
	}
	return float64(ok)/float64(letters) >= 0.55
}

func displayTitle(p page, canCJK bool) string {
	t := strings.TrimSpace(p.H1)
	if t == "" {
		t = strings.Split(p.Title, "|")[0]
		t = strings.TrimSpace(t)
	}
	if canCJK || mostlyPrintable(t) {
		return t
	}
	// Fall back to a readable path label when the title is mostly CJK
	// (embedded gofont lacks glyphs; og:title meta still carries the real text).
	path := strings.Trim(p.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || path == "" {
		return "Ebi Showcase"
	}
	label := parts[len(parts)-1]
	label = strings.ReplaceAll(label, "-", " ")
	if label == "" {
		return "Ebi Showcase"
	}
	return strings.ToUpper(label[:1]) + label[1:]
}

func badge(p page) string {
	e := strings.TrimSpace(p.Eyebrow)
	if e != "" && mostlyPrintable(e) {
		if len(e) > 48 {
			return e[:48]
		}
		return e
	}
	switch p.Kind {
	case "home":
		return "EBITENGINE CURRICULUM"
	case "core":
		return "CORE LESSON"
	case "track":
		return "GENRE TRACK"
	case "vfx":
		return "VISUAL EFFECTS LAB"
	case "guide":
		return "GUIDE"
	default:
		return "EBI SHOWCASE"
	}
}

func render(p page, bold, regular, small font.Face, canCJK bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, ogW, ogH))
	accent := hashColor(p.Key + p.Kind)
	deep := color.RGBA{14, 22, 48, 255}
	mid := mix(deep, accent, 0.22)
	fillGrad(img, mid, deep)

	// Decorative orbs (海老天-ish glow)
	fillCircle(img, 1040, 120, 180, color.RGBA{accent.R, accent.G, accent.B, 40})
	fillCircle(img, 180, 520, 140, color.RGBA{46, 201, 174, 35})
	fillCircle(img, 900, 500, 90, color.RGBA{255, 138, 92, 30})

	// Left accent bar
	fillRect(img, image.Rect(0, 0, 18, ogH), accent)

	// Brand chip
	fillRect(img, image.Rect(64, 56, 280, 108), color.RGBA{20, 32, 58, 230})
	drawString(img, small, "EBI SHOWCASE", 84, 90, accent)

	// Badge / eyebrow
	b := badge(p)
	drawString(img, small, strings.ToUpper(b), 64, 170, mix(accent, color.RGBA{255, 255, 255, 255}, 0.35))

	// Title
	title := displayTitle(p, canCJK)
	titleFace := bold
	lines := wrap(titleFace, title, 1000)
	if len(lines) > 3 {
		lines = lines[:3]
		lines[2] = strings.TrimSuffix(lines[2], "…") + "…"
	}
	y := 260
	for _, line := range lines {
		drawString(img, titleFace, line, 64, y, color.RGBA{245, 250, 255, 255})
		y += 70
	}

	// Description
	desc := p.Description
	if !canCJK && !mostlyPrintable(desc) {
		desc = p.Path
		if desc == "" {
			desc = "Playable Ebitengine lessons"
		}
	}
	if len([]rune(desc)) > 110 {
		r := []rune(desc)
		desc = string(r[:107]) + "..."
	}
	for i, line := range wrap(regular, desc, 980) {
		if i >= 2 {
			break
		}
		drawString(img, regular, line, 64, 500+i*34, color.RGBA{180, 200, 220, 255})
	}

	// Lang / kind footer
	footer := fmt.Sprintf("%s  ·  %s", strings.ToUpper(p.Lang), p.Kind)
	drawString(img, small, footer, 64, 600, color.RGBA{140, 160, 190, 255})

	// Tempura-ish shrimp curve accent
	for i := 0; i < 24; i++ {
		t := float64(i) / 23
		a := -0.4 + t*2.2
		x := 1080 + int(math.Cos(a)*70)
		y := 320 + int(math.Sin(a)*42)
		fillCircle(img, x, y, 7-i/5, mix(color.RGBA{255, 214, 140, 255}, color.RGBA{255, 120, 110, 255}, t))
	}

	return img
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	manPath := filepath.Join(root, "web/assets/og/manifest.json")
	raw, err := os.ReadFile(manPath)
	if err != nil {
		panic(err)
	}
	var man manifest
	if err := json.Unmarshal(raw, &man); err != nil {
		panic(err)
	}

	bold := mustFace(gobold.TTF, 54)
	regular := mustFace(goregular.TTF, 28)
	small := mustFace(gobold.TTF, 22)
	canCJK := false
	if sys := trySystemFaces(54); sys != nil {
		bold = sys
		if sys2 := trySystemFaces(28); sys2 != nil {
			regular = sys2
		}
		if sys3 := trySystemFaces(22); sys3 != nil {
			small = sys3
		}
		canCJK = true
		fmt.Println("using system Unicode font for OG titles")
	} else {
		fmt.Println("no system CJK font; Japanese titles fall back to path labels (meta still has full text)")
	}

	outDir := filepath.Join(root, "web/assets/og")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(err)
	}

	for i, p := range man.Pages {
		img := render(p, bold, regular, small, canCJK)
		out := filepath.Join(root, "web", filepath.FromSlash(p.Image))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			panic(err)
		}
		f, err := os.Create(out)
		if err != nil {
			panic(err)
		}
		if err := png.Encode(f, img); err != nil {
			f.Close()
			panic(err)
		}
		f.Close()
		if (i+1)%50 == 0 || i+1 == len(man.Pages) {
			fmt.Printf("og images %d/%d\n", i+1, len(man.Pages))
		}
	}
}
