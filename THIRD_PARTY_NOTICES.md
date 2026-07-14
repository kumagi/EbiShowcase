# Third-party notices

Ebi Showcase's own source code, lesson prose, original art, diagrams, generated
thumbnails, and generated game assets are licensed under Apache License 2.0 unless
an adjacent notice says otherwise. This file records material that keeps its
original compatible license. It is part of source and release documentation: do
not remove an entry or replace a third-party license with Apache-2.0.

When adding a dependency, font, asset, copied code, or generated output derived
from another work:

1. confirm that its license is compatible with this repository's Apache-2.0
   distribution;
2. keep the original license text whenever that license requires it;
3. add its name, version, copyright holder/source URL, license, and repository
   paths to this file; and
4. update any adjacent asset notice and this inventory in the same change.

Unknown provenance is not acceptable. A reference to a famous game is a teaching
label, not permission to copy its characters, art, sounds, levels, text, or
branding.

## Go modules compiled into the games and tools

The versions below are pinned in `go.mod`. The project's Apache-2.0 `LICENSE`
also supplies the full Apache-2.0 text for entries marked Apache-2.0. The full
BSD-3-Clause text used by the Go-derived entries appears below.

| Module | Version | License | Source |
| --- | --- | --- | --- |
| `github.com/hajimehoshi/ebiten/v2` | `v2.9.9` | Apache-2.0 | https://github.com/hajimehoshi/ebiten |
| `github.com/ebitengine/oto/v3` | `v3.4.0` | Apache-2.0 | https://github.com/ebitengine/oto |
| `github.com/ebitengine/hideconsole` | `v1.0.0` | Apache-2.0 | https://github.com/ebitengine/hideconsole |
| `github.com/ebitengine/purego` | `v0.9.0` | Apache-2.0 | https://github.com/ebitengine/purego |
| `github.com/go-text/typesetting` | `v0.3.0` | Unlicense OR BSD-3-Clause, Copyright 2021 The go-text authors | https://github.com/go-text/typesetting |
| `golang.org/x/image` | `v0.44.0` | BSD-3-Clause, Copyright 2009 The Go Authors | https://go.googlesource.com/image |
| `github.com/ebitengine/gomobile` | `v0.0.0-20250923094054-ea854a63cce1` | BSD-3-Clause, Copyright 2009 The Go Authors | https://github.com/ebitengine/gomobile |
| `github.com/jezek/xgb` | `v1.1.1` | BSD-3-Clause, Copyright 2009 The XGB Authors | https://github.com/jezek/xgb |
| `github.com/rivo/uniseg` | `v0.4.7` | MIT | https://github.com/rivo/uniseg |
| `golang.org/x/sync` | `v0.22.0` | BSD-3-Clause, Copyright 2009 The Go Authors | https://go.googlesource.com/sync |
| `golang.org/x/sys` | `v0.36.0` | BSD-3-Clause, Copyright 2009 The Go Authors | https://go.googlesource.com/sys |
| `golang.org/x/text` | `v0.40.0` | BSD-3-Clause, Copyright 2009 The Go Authors | https://go.googlesource.com/text |

### BSD-3-Clause license text

```text
Copyright 2009 The Go Authors. All rights reserved.
Copyright 2009 The XGB Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.
3. Neither the name of Google LLC, Google Inc., The Go Authors, The XGB Authors,
   nor the names of their contributors may be used to endorse or promote
   products derived from this software without specific prior written
   permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```

## Font

| Asset | License | Copyright / source | Use |
| --- | --- | --- | --- |
| `internal/ogfont/NotoSansJP.ttf` | SIL Open Font License 1.1 | The Noto Project Authors; https://github.com/google/fonts/tree/main/ofl/notosansjp | Deterministic Japanese OGP rendering by `cmd/gen-og-images` and bilingual Ebitengine `text/v2` UI through `internal/uilab` |

The full SIL Open Font License 1.1 is retained beside the font at
[`internal/ogfont/OFL.txt`](internal/ogfont/OFL.txt). See
[`internal/ogfont/README.md`](internal/ogfont/README.md) for its exact scope.

## Original and generated project assets

These are recorded to prevent their status from being confused with an external
game asset:

| Asset group | License / provenance |
| --- | --- |
| `assets/characters/`, `internal/hero/`, `internal/heroatlas/` | Original Ebi Tenjiroh character art and project-generated variants; Apache-2.0. The downloadable atlas has a separate CC0-1.0 dedication in `web/assets/ebi-boy-atlas-LICENSE.txt`. |
| `internal/vfxsprites/`, `web/assets/vfx-*.png` | Generated by this repository's `cmd/gen-vfx`; Apache-2.0. |
| `internal/trackatlas/`, `web/assets/track-atlas.*` | Generated by this repository's `cmd/gen-track-atlas`; downloadable atlas and metadata are CC0-1.0 as stated in `web/assets/track-atlas-LICENSE.txt`. |
| `web/assets/home-thumbnails/`, `web/assets/og/`, `web/assets/diagrams/` | Generated from this repository's games, lesson data, and drawing code; Apache-2.0, except OGP image exports are CC0-1.0 as stated in `web/assets/og/LICENSE.txt`. |
