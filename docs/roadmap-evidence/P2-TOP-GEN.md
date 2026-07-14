# P2-TOP-GEN — `gen-topdown-adventure-track.mjs` を同仕様で改修し再生成する

Status: PASS

## Required evidence

- [x] Dual panel
- [x] Unique concept-row
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/gen-topdown-adventure-track.mjs`, `scripts/authoring-lesson-helpers.mjs`, and the 16 generated JA/EN step pages.
- Behavior: every step reads its real `main.go` wrapper, labels the shared `internal/topdownadventuregame/game.go` `Update` excerpt, replaces the generic card row with DATA / UPDATE / DRAW cards, and adds a source-path rule challenge.

## Commands and results

```text
node scripts/gen-topdown-adventure-track.mjs
Generated 8-step top-down adventure track in JA/EN.
node scripts/check-lessons.mjs
Checked 474 pages (237 playable lessons). OK.
go test ./games/tracks/topdown-adventure/...
All eight packages compile successfully.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese: `web/ja/tracks/topdown-adventure/{eight-way,ebi-adventure}/index.html` includes the two source-labelled panels and localized challenge.
- English: `web/en/tracks/topdown-adventure/{eight-way,ebi-adventure}/index.html` includes the matching English panels and challenge.
- Readability / accessibility:
- Screenshots / recordings:
