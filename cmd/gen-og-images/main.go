// Command gen-og-images renders per-page Open Graph images (1200×630)
// from web/assets/og/manifest.json (produced by scripts/inject-ogp.mjs).
//
// Each card uses a real WASM gameplay capture from home-thumbnails. Track
// lessons without a dedicated capture show the finished game they build into.
//
//	go run ./cmd/gen-og-images
package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/webp"
)

const (
	ogW = 1200
	ogH = 630
)

type page struct {
	File         string `json:"file"`
	Path         string `json:"path"`
	Key          string `json:"key"`
	Lang         string `json:"lang"`
	Kind         string `json:"kind"`
	Title        string `json:"title"`
	H1           string `json:"h1"`
	Eyebrow      string `json:"eyebrow"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Hook         string `json:"hook"`
	Action       string `json:"action"`
	Preview      string `json:"preview"`
	PreviewMode  string `json:"previewMode"`
	PreviewLabel string `json:"previewLabel"`
}

type manifest struct {
	Origin string `json:"origin"`
	Pages  []page `json:"pages"`
}

type faces struct {
	title   font.Face
	body    font.Face
	small   font.Face
	micro   font.Face
	canCJK  bool
	closeFn func()
}

func weightScore(name, wanted string) int {
	name = strings.ToLower(name)
	if wanted == "bold" {
		for i, label := range []string{"black", "heavy", "extra bold", "extrabold", "bold", "semi bold", "semibold", "demibold", "medium", "regular", "light", "thin"} {
			if strings.Contains(name, label) {
				return 20 - i
			}
		}
		return 0
	}
	for i, label := range []string{"regular", "normal", "book", "medium", "light", "semi bold", "semibold", "bold", "thin"} {
		if strings.Contains(name, label) {
			return 20 - i
		}
	}
	return 0
}

func loadFaceFile(path string, size float64, wanted string) (font.Face, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var selected *sfnt.Font
	if collection, collectionErr := opentype.ParseCollection(raw); collectionErr == nil {
		bestScore := -1
		for i := 0; i < collection.NumFonts(); i++ {
			candidate, fontErr := collection.Font(i)
			if fontErr != nil {
				continue
			}
			name, _ := candidate.Name(nil, sfnt.NameIDTypographicSubfamily)
			if name == "" {
				name, _ = candidate.Name(nil, sfnt.NameIDSubfamily)
			}
			score := weightScore(name, wanted)
			if selected == nil || score > bestScore {
				selected = candidate
				bestScore = score
			}
		}
	}
	if selected == nil {
		selected, err = opentype.Parse(raw)
		if err != nil {
			return nil, err
		}
	}
	return opentype.NewFace(selected, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
}

func supportsJapanese(face font.Face) bool {
	for _, r := range "日本語海老天次郎" {
		if _, ok := face.GlyphAdvance(r); !ok {
			return false
		}
	}
	return true
}

func trySystemFace(size float64) font.Face {
	candidates := []string{
		"/System/Library/Fonts/ヒラギノ角ゴシック W6.ttc",
		"/System/Library/Fonts/ヒラギノ角ゴシック W3.ttc",
		"/System/Library/Fonts/Supplemental/Arial Unicode.ttf",
		"/Library/Fonts/Arial Unicode.ttf",
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Bold.ttc",
		"/usr/share/opentype/noto/NotoSansCJK-Bold.ttc",
	}
	for _, path := range candidates {
		if face, err := loadFaceFile(path, size, "bold"); err == nil && supportsJapanese(face) {
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

func loadFaces(root string) faces {
	result := faces{
		title:  mustFace(gobold.TTF, 52),
		body:   mustFace(goregular.TTF, 27),
		small:  mustFace(gobold.TTF, 20),
		micro:  mustFace(gobold.TTF, 16),
		canCJK: false,
	}
	bundled := filepath.Join(root, "internal/ogfont/NotoSansJP.ttf")
	if face, err := loadFaceFile(bundled, 52, "bold"); err == nil && supportsJapanese(face) {
		result.title = face
		result.body, _ = loadFaceFile(bundled, 27, "regular")
		result.small, _ = loadFaceFile(bundled, 20, "bold")
		result.micro, _ = loadFaceFile(bundled, 16, "bold")
		result.canCJK = true
		fmt.Println("using bundled Noto Sans JP for OG text")
		return result
	}
	if face := trySystemFace(52); face != nil {
		result.title = face
		result.body = trySystemFace(27)
		result.small = trySystemFace(20)
		result.micro = trySystemFace(16)
		result.canCJK = true
		fmt.Println("using system Unicode font for OG text")
		return result
	}
	fmt.Println("no CJK font; Japanese card text falls back to route labels")
	return result
}

func mix(a, b color.RGBA, t float64) color.RGBA {
	f := func(x, y uint8) uint8 { return uint8(float64(x)*(1-t) + float64(y)*t) }
	return color.RGBA{f(a.R, b.R), f(a.G, b.G), f(a.B, b.B), 255}
}

func hashColor(s string) color.RGBA {
	var hash uint32
	for i := 0; i < len(s); i++ {
		hash = hash*33 + uint32(s[i])
	}
	palette := []color.RGBA{
		{46, 230, 200, 255},
		{141, 123, 255, 255},
		{255, 138, 92, 255},
		{74, 144, 226, 255},
		{245, 199, 75, 255},
		{230, 90, 120, 255},
	}
	return palette[int(hash)%len(palette)]
}

func cardPalette(accent color.RGBA) color.Palette {
	deep := color.RGBA{9, 16, 39, 255}
	top := mix(deep, accent, .26)
	result := append(color.Palette{}, palette.WebSafe...)
	for i := 0; i < 24; i++ {
		result = append(result, mix(top, deep, float64(i)/23))
	}
	for i := 0; i < 8; i++ {
		result = append(result, mix(accent, color.RGBA{255, 255, 255, 255}, float64(i)/10))
	}
	result = append(result,
		color.RGBA{247, 250, 255, 255},
		color.RGBA{192, 207, 234, 255},
		color.RGBA{126, 148, 185, 255},
		color.RGBA{243, 247, 255, 255},
		color.RGBA{30, 42, 75, 255},
		color.RGBA{10, 20, 47, 255},
		color.RGBA{8, 18, 39, 255},
		color.RGBA{15, 24, 52, 255},
	)
	return result
}

func fillGradient(img *image.RGBA, top, bottom color.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		t := float64(y-bounds.Min.Y) / float64(bounds.Dy())
		col := mix(top, bottom, t)
		draw.Draw(img, image.Rect(0, y, ogW, y+1), image.NewUniform(col), image.Point{}, draw.Src)
	}
}

func roundedMask(width, height, radius int) *image.Alpha {
	mask := image.NewAlpha(image.Rect(0, 0, width, height))
	r2 := radius * radius
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cx, cy := x, y
			if x >= radius && x < width-radius || y >= radius && y < height-radius {
				mask.SetAlpha(x, y, color.Alpha{A: 255})
				continue
			}
			if x >= width-radius {
				cx = width - radius - 1
			} else if x >= radius {
				cx = x
			} else {
				cx = radius
			}
			if y >= height-radius {
				cy = height - radius - 1
			} else if y >= radius {
				cy = y
			} else {
				cy = radius
			}
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= r2 {
				mask.SetAlpha(x, y, color.Alpha{A: 255})
			}
		}
	}
	return mask
}

func fillRounded(img *image.RGBA, rect image.Rectangle, radius int, col color.RGBA) {
	mask := roundedMask(rect.Dx(), rect.Dy(), radius)
	draw.DrawMask(img, rect, image.NewUniform(col), image.Point{}, mask, image.Point{}, draw.Over)
}

func drawString(img *image.RGBA, face font.Face, value string, x, y int, col color.RGBA) {
	d := &font.Drawer{Dst: img, Src: image.NewUniform(col), Face: face, Dot: fixed.P(x, y)}
	d.DrawString(value)
}

func drawStrongString(img *image.RGBA, face font.Face, value string, x, y int, col color.RGBA) {
	for oy := -1; oy <= 1; oy++ {
		for ox := -1; ox <= 1; ox++ {
			drawString(img, face, value, x+ox, y+oy, col)
		}
	}
}

func measure(face font.Face, value string) int {
	return font.MeasureString(face, value).Ceil()
}

func wrap(face font.Face, value string, maxWidth int) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	if strings.Contains(value, " ") {
		words := strings.Fields(value)
		lines := make([]string, 0, 3)
		current := words[0]
		for _, word := range words[1:] {
			trial := current + " " + word
			if measure(face, trial) <= maxWidth {
				current = trial
				continue
			}
			lines = append(lines, current)
			current = word
		}
		return append(lines, current)
	}
	lines := make([]string, 0, 3)
	current := make([]rune, 0, len([]rune(value)))
	for _, r := range value {
		trial := string(append(append([]rune{}, current...), r))
		if len(current) > 0 && measure(face, trial) > maxWidth {
			lines = append(lines, string(current))
			current = []rune{r}
			continue
		}
		current = append(current, r)
	}
	if len(current) > 0 {
		lines = append(lines, string(current))
	}
	return lines
}

func truncateRunes(value string, max int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= max {
		return string(runes)
	}
	return strings.TrimRight(string(runes[:max-1]), " 、,;:") + "…"
}

func mostlyPrintable(value string) bool {
	letters, printable := 0, 0
	for _, r := range value {
		if unicode.IsSpace(r) {
			continue
		}
		letters++
		if r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r) || strings.ContainsRune("—–-·|/:,.'\"!?", r)) {
			printable++
		}
	}
	return letters == 0 || float64(printable)/float64(letters) >= .55
}

func fallbackLabel(p page) string {
	path := strings.Trim(p.Path, "/")
	if path == "" {
		return "Ebi Showcase"
	}
	parts := strings.Split(path, "/")
	label := strings.ReplaceAll(parts[len(parts)-1], "-", " ")
	if label == "" {
		return "Ebi Showcase"
	}
	return strings.ToUpper(label[:1]) + label[1:]
}

func displayText(value string, p page, canCJK bool) string {
	if canCJK || mostlyPrintable(value) {
		return strings.TrimSpace(value)
	}
	return fallbackLabel(p)
}

func badge(p page) string {
	eyebrow := strings.TrimSpace(p.Eyebrow)
	if eyebrow != "" && mostlyPrintable(eyebrow) {
		return truncateRunes(strings.ToUpper(eyebrow), 42)
	}
	switch p.Kind {
	case "home":
		return "PLAYABLE EBITENGINE CURRICULUM"
	case "core":
		return "CORE LESSON"
	case "track":
		return "GENRE TRACK"
	case "vfx":
		return "VISUAL EFFECTS LAB"
	case "guide":
		return "PRACTICAL GUIDE"
	case "build":
		return "BUILD TRACK"
	case "graduation":
		return "GRADUATION PROJECT"
	default:
		return "EBI SHOWCASE"
	}
}

func loadPreview(root string, p page) (image.Image, error) {
	path := filepath.Join(root, "web", filepath.FromSlash(p.Preview))
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s preview: %w", p.Path, err)
	}
	defer file.Close()
	preview, err := webp.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("%s preview decode: %w", p.Path, err)
	}
	return preview, nil
}

func drawPreview(img *image.RGBA, source image.Image, p page, face font.Face, accent color.RGBA) {
	panel := image.Rect(802, 24, 1174, 596)
	fillRounded(img, panel.Add(image.Pt(0, 10)), 30, color.RGBA{0, 0, 0, 80})
	fillRounded(img, panel, 30, color.RGBA{243, 247, 255, 255})
	fillRounded(img, image.Rect(813, 35, 1163, 545), 23, mix(accent, color.RGBA{15, 24, 52, 255}, .84))

	shotRect := image.Rect(828, 48, 1148, 528)
	scaled := image.NewRGBA(image.Rect(0, 0, shotRect.Dx(), shotRect.Dy()))
	xdraw.CatmullRom.Scale(scaled, scaled.Bounds(), source, source.Bounds(), draw.Src, nil)
	mask := roundedMask(shotRect.Dx(), shotRect.Dy(), 18)
	draw.DrawMask(img, shotRect, scaled, image.Point{}, mask, image.Point{}, draw.Over)

	label := truncateRunes(strings.ToUpper(p.PreviewLabel), 33)
	labelWidth := measure(face, label)
	drawStrongString(img, face, label, panel.Min.X+(panel.Dx()-labelWidth)/2, 570, color.RGBA{30, 42, 75, 255})
}

func render(root string, p page, f faces) (*image.RGBA, error) {
	preview, err := loadPreview(root, p)
	if err != nil {
		return nil, err
	}

	img := image.NewRGBA(image.Rect(0, 0, ogW, ogH))
	accent := hashColor(p.Key + p.Kind)
	deep := color.RGBA{9, 16, 39, 255}
	fillGradient(img, mix(deep, accent, .26), deep)

	for y := 18; y < ogH; y += 34 {
		for x := 18; x < 790; x += 34 {
			fillRounded(img, image.Rect(x, y, x+3, y+3), 1, color.RGBA{accent.R, accent.G, accent.B, 34})
		}
	}
	fillRounded(img, image.Rect(0, 0, 16, ogH), 0, accent)
	fillRounded(img, image.Rect(54, 42, 276, 92), 15, color.RGBA{10, 20, 47, 218})
	drawStrongString(img, f.small, "EBI SHOWCASE", 76, 75, accent)

	drawString(img, f.micro, badge(p), 58, 142, mix(accent, color.RGBA{255, 255, 255, 255}, .42))

	title := displayText(p.H1, p, f.canCJK)
	titleLines := wrap(f.title, title, 690)
	if len(titleLines) > 3 {
		titleLines = titleLines[:3]
		titleLines[2] = truncateRunes(titleLines[2], 22)
	}
	y := 223
	for _, line := range titleLines {
		drawStrongString(img, f.title, line, 58, y, color.RGBA{247, 250, 255, 255})
		y += 61
	}

	hook := displayText(p.Hook, p, f.canCJK)
	hookLines := wrap(f.body, hook, 680)
	if len(hookLines) > 2 {
		hookLines = hookLines[:2]
		hookLines[1] = truncateRunes(hookLines[1], 34)
	}
	hookY := y + 28
	if hookY < 410 {
		hookY = 410
	}
	for _, line := range hookLines {
		drawString(img, f.body, line, 58, hookY, color.RGBA{192, 207, 234, 255})
		hookY += 36
	}

	action := displayText(p.Action, p, f.canCJK)
	action = truncateRunes(action, 45)
	pillWidth := measure(f.small, action) + 40
	if pillWidth > 690 {
		pillWidth = 690
	}
	fillRounded(img, image.Rect(56, 520, 56+pillWidth, 568), 24, accent)
	drawStrongString(img, f.small, action, 76, 552, color.RGBA{8, 18, 39, 255})

	route := "/" + strings.Trim(p.Path, "/") + "/"
	if p.Path == "" {
		route = "/"
	}
	drawString(img, f.micro, strings.ToUpper(p.Lang)+"  ·  "+route, 58, 605, color.RGBA{126, 148, 185, 255})
	drawPreview(img, preview, p, f.micro, accent)
	return img, nil
}

func selected(key string) bool {
	filter := strings.TrimSpace(os.Getenv("OGP_ONLY"))
	if filter == "" {
		return true
	}
	for _, candidate := range strings.Split(filter, ",") {
		if strings.TrimSpace(candidate) == key {
			return true
		}
	}
	return false
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	raw, err := os.ReadFile(filepath.Join(root, "web/assets/og/manifest.json"))
	if err != nil {
		panic(err)
	}
	var man manifest
	if err := json.Unmarshal(raw, &man); err != nil {
		panic(err)
	}

	f := loadFaces(root)
	outDir := filepath.Join(root, "web/assets/og")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		panic(err)
	}

	total := 0
	for _, p := range man.Pages {
		if selected(p.Key) {
			total++
		}
	}
	rendered := 0
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	for _, p := range man.Pages {
		if !selected(p.Key) {
			continue
		}
		img, err := render(root, p, f)
		if err != nil {
			panic(err)
		}
		out := filepath.Join(root, "web", filepath.FromSlash(p.Image))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			panic(err)
		}
		file, err := os.Create(out)
		if err != nil {
			panic(err)
		}
		paletted := image.NewPaletted(img.Bounds(), cardPalette(hashColor(p.Key+p.Kind)))
		draw.Draw(paletted, paletted.Bounds(), img, image.Point{}, draw.Src)
		if err := encoder.Encode(file, paletted); err != nil {
			file.Close()
			panic(err)
		}
		if err := file.Close(); err != nil {
			panic(err)
		}
		rendered++
		if rendered%50 == 0 || rendered == total {
			fmt.Printf("og images %d/%d\n", rendered, total)
		}
	}
}
