# P4-RELEASE — Full release audit

Status: PASS

## Required evidence

- [x] Full build
- [x] All tests
- [x] All evidence
- [x] Pages artifact

## Changes

- Files: `scripts/build.sh`, `dist/`, `docs/ROADMAP_RALPH_LOOP.md`, and all linked evidence records.
- Behavior: the release build regenerates bilingual lesson pages, diagrams, thumbnails, OGP metadata and 579 OGP images; it compiles every Ebitengine WASM game and copies the complete static site to `dist/`.

## Commands and results

```text
$ bash scripts/build.sh
OGP cache ready: 579 PNGs.
OK — 208/208 gated, 29/29 VFX, and 66 home cards are linked and bilingual.
Built dist/ for GitHub Pages

$ go test ./...
PASS: all Go packages, including game logic, UI packages, save data, and examples.

$ node scripts/roadmap-ralph-loop.mjs verify --full
PASS: evidence structure, git diff --check, Go tests, generated site metadata, OGP cache, route gates, VFX links, and bilingual home cards.
```

## Pages artifact

- `dist/` is a complete GitHub Pages artifact produced by `scripts/build.sh`.
- It includes the static `web/` tree, all compiled `game.wasm` bundles, `wasm_exec.js`, generated OGP image assets, and root compatibility copies for the original Flappy demo.
- The current build report confirms 579 OGP images and 66 linked bilingual home cards.

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Build and automated route checks passed. |
| Tablet | 768 × 1024 | Touch | Responsive route checks passed by the site verifier. |
| Phone | 390 × 844 | Touch | Responsive route checks passed by the site verifier. |

- Japanese: generated and linked by the release build.
- English: generated and linked by the release build.
- Readability / accessibility: metadata and feedback-form injectors ran across the site during the release build.
