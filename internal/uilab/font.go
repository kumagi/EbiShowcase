// Package uilab contains production-oriented UI helpers shared by labs.
package uilab

import (
	"bytes"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kumagi/EbiShowcase/internal/ogfont"
	"golang.org/x/text/language"
	"sync"
)

var once sync.Once
var source *text.GoTextFaceSource
var sourceErr error

// Face uses the bundled Noto Sans JP for both Japanese and Latin glyphs, so a
// UI never switches to a missing-glyph fallback midway through a sentence.
func Face(lang string, size float64) (*text.GoTextFace, error) {
	once.Do(func() { source, sourceErr = text.NewGoTextFaceSource(bytes.NewReader(ogfont.NotoSansJP)) })
	if sourceErr != nil {
		return nil, sourceErr
	}
	tag := language.English
	if lang == "ja" {
		tag = language.Japanese
	}
	return &text.GoTextFace{Source: source, Size: size, Language: tag}, nil
}
