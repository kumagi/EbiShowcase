// Package vfxlive is the shared "live Go + mouse knobs" shell for Visual Effects Lab.
// Each demo shows a Go snippet whose numeric tokens update as the learner drags
// on-screen sliders or taps toggles; the same values drive the drawing stage.
package vfxlive

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kumagi/EbiShowcase/internal/vfxui"
)

const (
	Width  = 480
	Height = 720
)

// Param is one mouse-tunable value shown both in the Go panel and as a control.
type Param struct {
	Key    string  // placeholder token in Lines, e.g. "angle"
	Label  string  // slider label
	Value  float64 // current
	Min    float64
	Max    float64
	Step   float64 // 0 = continuous; >0 snaps
	Format string  // printf format; default "%.2f"
	Bool   bool    // render as ON/OFF toggle instead of slider
}

func (p *Param) Text() string {
	if p.Bool {
		if p.Value >= 0.5 {
			return "true"
		}
		return "false"
	}
	f := p.Format
	if f == "" {
		if p.Step >= 1 && p.Step == math.Trunc(p.Step) {
			f = "%.0f"
		} else {
			f = "%.2f"
		}
	}
	return fmt.Sprintf(f, p.Value)
}

func (p *Param) BoolOn() bool { return p.Bool && p.Value >= 0.5 }

func (p *Param) Set(v float64) {
	if p.Bool {
		if v >= 0.5 {
			p.Value = 1
		} else {
			p.Value = 0
		}
		return
	}
	if v < p.Min {
		v = p.Min
	}
	if v > p.Max {
		v = p.Max
	}
	if p.Step > 0 {
		v = math.Round(v/p.Step) * p.Step
	}
	p.Value = v
}

func (p *Param) Toggle() {
	if p.BoolOn() {
		p.Value = 0
	} else {
		p.Value = 1
	}
}

// Shell owns the code panel + bottom controls and reserves a middle stage rect.
type Shell struct {
	Title    string
	Hint     string
	Lines    []string // Go-ish lines; use {key} for live tokens
	Params   []*Param
	Tokens   map[string]string // display-only {key} values (not sliders)
	dragging int               // param index, or -1
	codeH    int
	ctrlH    int
}

// New builds a shell. Call Layout() after changing Params count if needed.
func New(title string, lines []string, params ...*Param) *Shell {
	s := &Shell{Title: title, Lines: lines, Params: params, Tokens: map[string]string{}, dragging: -1}
	s.recompute()
	return s
}

// SetToken sets a display-only code token (not a mouse control).
func (s *Shell) SetToken(key, text string) {
	if s.Tokens == nil {
		s.Tokens = map[string]string{}
	}
	s.Tokens[key] = text
}

func (s *Shell) recompute() {
	// ~16px per code line + title + padding.
	s.codeH = 28 + len(s.Lines)*16 + 12
	if s.codeH < 100 {
		s.codeH = 100
	}
	if s.codeH > 220 {
		s.codeH = 220
	}
	n := len(s.Params)
	s.ctrlH = 36 + n*36 + 28
	if s.Hint != "" {
		s.ctrlH += 18
	}
	if s.ctrlH < 100 {
		s.ctrlH = 100
	}
}

// Stage returns the middle rectangle reserved for the visual demo.
func (s *Shell) Stage() (x, y, w, h float64) {
	s.recompute()
	y0 := float64(s.codeH)
	h0 := float64(Height - s.codeH - s.ctrlH)
	return 0, y0, Width, h0
}

// Get returns a param by key (panics if missing — demos are static).
func (s *Shell) Get(key string) float64 {
	for _, p := range s.Params {
		if p.Key == key {
			return p.Value
		}
	}
	panic("vfxlive: missing param " + key)
}

func (s *Shell) Param(key string) *Param {
	for _, p := range s.Params {
		if p.Key == key {
			return p
		}
	}
	panic("vfxlive: missing param " + key)
}

func (s *Shell) Bool(key string) bool { return s.Param(key).BoolOn() }

// Update handles slider drag / toggle taps. Returns true if a control ate the press.
func (s *Shell) Update() (ate bool) {
	s.recompute()
	ctrlTop := float64(Height - s.ctrlH)

	if x, y, ok := vfxui.JustPressed(); ok {
		for i, p := range s.Params {
			rx, ry, rw, rh := s.controlRect(i)
			if x >= rx && x <= rx+rw && y >= ry && y <= ry+rh {
				ate = true
				if p.Bool {
					p.Toggle()
				} else {
					s.dragging = i
					s.setFromX(i, x)
				}
				return
			}
		}
		// Ignore presses in chrome so stage taps stay clean.
		if y < float64(s.codeH) || y >= ctrlTop {
			ate = true
		}
	}

	if s.dragging >= 0 {
		if x, _, ok := vfxui.Held(); ok {
			s.setFromX(s.dragging, x)
			ate = true
		} else {
			s.dragging = -1
		}
	}
	return ate
}

func (s *Shell) setFromX(i int, x float64) {
	p := s.Params[i]
	rx, _, rw, _ := s.controlRect(i)
	trackX := rx + 110
	trackW := rw - 130
	t := (x - trackX) / trackW
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	p.Set(p.Min + t*(p.Max-p.Min))
}

func (s *Shell) controlRect(i int) (x, y, w, h float64) {
	ctrlTop := float64(Height - s.ctrlH)
	return 12, ctrlTop + 28 + float64(i)*36, Width - 24, 32
}

// Draw paints code panel + controls. Call after filling the stage area yourself,
// or call DrawChrome before/after stage as you prefer.
func (s *Shell) Draw(dst *ebiten.Image) {
	s.recompute()
	s.drawCode(dst)
	s.drawControls(dst)
}

func (s *Shell) drawCode(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 0, 0, Width, float32(s.codeH), color.RGBA{8, 12, 24, 255}, false)
	vector.StrokeRect(dst, 0, 0, Width, float32(s.codeH), 2, color.RGBA{40, 60, 100, 255}, false)
	ebitenutil.DebugPrintAt(dst, "LIVE GO  /  "+s.Title, 12, 8)
	y := 30
	for _, line := range s.Lines {
		rendered := s.expand(line)
		// Dim comment-like lines.
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			ebitenutil.DebugPrintAt(dst, rendered, 12, y)
		} else {
			ebitenutil.DebugPrintAt(dst, rendered, 12, y)
		}
		y += 16
		if y > s.codeH-8 {
			break
		}
	}
}

func (s *Shell) expand(line string) string {
	out := line
	for k, v := range s.Tokens {
		out = strings.ReplaceAll(out, "{"+k+"}", v)
	}
	for _, p := range s.Params {
		out = strings.ReplaceAll(out, "{"+p.Key+"}", p.Text())
	}
	return out
}

func (s *Shell) drawControls(dst *ebiten.Image) {
	ctrlTop := float64(Height - s.ctrlH)
	vector.DrawFilledRect(dst, 0, float32(ctrlTop), Width, float32(s.ctrlH), color.RGBA{12, 16, 32, 255}, false)
	vector.StrokeRect(dst, 0, float32(ctrlTop), Width, float32(s.ctrlH), 2, color.RGBA{40, 60, 100, 255}, false)
	ebitenutil.DebugPrintAt(dst, "DRAG SLIDERS / TAP TOGGLES  —  code updates live", 12, int(ctrlTop)+8)
	for i, p := range s.Params {
		rx, ry, rw, rh := s.controlRect(i)
		active := s.dragging == i || (p.Bool && p.BoolOn())
		fill := color.RGBA{28, 38, 64, 255}
		if active {
			fill = color.RGBA{36, 56, 96, 255}
		}
		vector.DrawFilledRect(dst, float32(rx), float32(ry), float32(rw), float32(rh), fill, false)
		edge := color.RGBA{70, 100, 150, 255}
		if active {
			edge = color.RGBA{120, 240, 220, 255}
		}
		vector.StrokeRect(dst, float32(rx), float32(ry), float32(rw), float32(rh), 2, edge, false)

		if p.Bool {
			state := "OFF"
			if p.BoolOn() {
				state = "ON"
			}
			ebitenutil.DebugPrintAt(dst, fmt.Sprintf("%-10s  [%s]  tap to toggle", p.Label, state), int(rx)+10, int(ry)+8)
			continue
		}

		ebitenutil.DebugPrintAt(dst, p.Label, int(rx)+8, int(ry)+8)
		trackX := rx + 110
		trackW := rw - 130
		trackY := ry + rh/2
		vector.StrokeLine(dst, float32(trackX), float32(trackY), float32(trackX+trackW), float32(trackY), 3, color.RGBA{60, 80, 120, 255}, false)
		t := 0.0
		if p.Max > p.Min {
			t = (p.Value - p.Min) / (p.Max - p.Min)
		}
		kx := trackX + t*trackW
		vector.DrawFilledRect(dst, float32(kx-6), float32(trackY-8), 12, 16, color.RGBA{120, 240, 220, 255}, false)
		ebitenutil.DebugPrintAt(dst, p.Text(), int(rx+rw)-48, int(ry)+8)
	}
	if s.Hint != "" {
		ebitenutil.DebugPrintAt(dst, s.Hint, 12, Height-22)
	}
}

// FillStage paints the stage background.
func (s *Shell) FillStage(dst *ebiten.Image, c color.Color) {
	_, y, w, h := s.Stage()
	vector.DrawFilledRect(dst, 0, float32(y), float32(w), float32(h), c, false)
}

// Parse is a tiny helper for demos that want int from a float param.
func ParseInt(v float64) int {
	i, _ := strconv.Atoi(fmt.Sprintf("%.0f", v))
	return i
}
