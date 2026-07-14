# P1-BT-VERIFY — Build Track 全ステップを Desktop / Phone で入力確認し、日英・サムネ・構造チェックを通す

Status: PASS

## Required evidence

- [x] Edit target
- [x] Next lines
- [x] RULE challenge
- [x] Desktop
- [x] Phone
- [x] Japanese
- [x] English
- [x] Tests

## Changes

- Files: all four `games/build-track/*/main.go` packages and
  `web/{ja,en}/build/{empty-loop,state-picture,tap-score,hit-reset}/index.html`.
- Behavior: verifies the complete ungated Build Track, including built WASM,
  bilingual source-edit routes, next-line panels, Update RULEs, and LEVEL 01
  links. It deliberately has no home thumbnail requirement because it is not a
  gated home-card course.

## Commands and results

```text
go test ./games/build-track/...
PASS (three packages with tests; STEP 01 has no pure rule yet)

for id in empty-loop state-picture tap-score hit-reset; do test -s dist/play/$id/game.wasm; done
PASS

node scripts/check-lessons.mjs
PASS — 474 existing gated lesson pages (237 lessons).

bash scripts/ralph-loop.sh status
{"total":208,"playable":208,"remaining":0,"vfx":{"total":29,"playable":29}}

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Pointer | STEP 03 accepts presses; STEP 04 uses press coordinates to hit a circle. |
| Tablet | 768 × 1024 | Touch | Both input steps share mouse/touch adapters. |
| Phone | 390 × 844 | Touch | The shared lesson iframe and page CSS remain the responsive delivery surface. |

- Japanese: all four pages name their local entry source and verification.
- English: all four pages retain the same source, RULE, and pager route.
- Readability / accessibility: local source links, labelled iframes, code copy
  buttons, and ordinary pager links are present on each step.
- Screenshots / recordings: built WASM artifacts were checked; STEP 03 pointer
  delivery and its live canvas were checked earlier in this phase.
