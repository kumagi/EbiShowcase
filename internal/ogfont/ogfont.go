// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
// Package ogfont exposes the bundled OFL-1.1 Noto Sans JP file to the small
// number of Ebitengine surfaces that must render Japanese text in WASM.
package ogfont

import _ "embed"

// NotoSansJP is redistributed under SIL Open Font License 1.1. See OFL.txt
// and THIRD_PARTY_NOTICES.md; it is not relicensed as Apache-2.0.
//
//go:embed NotoSansJP.ttf
var NotoSansJP []byte
