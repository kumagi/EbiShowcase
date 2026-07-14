# P2-RAY-GEN — `gen-raycaster-track.mjs` を同仕様で改修し再生成する

Status: PASS

## Required evidence

- [x] Dual panel
- [x] Unique concept-row
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/gen-raycaster-track.mjs`, `scripts/authoring-lesson-helpers.mjs`,
  and regenerated Raycaster pages.
- Behavior: each page shows its real entry, a labelled shared-engine excerpt,
  data/update/draw concepts, and a concrete map/mission RULE.

## Commands and results

```text
node scripts/gen-raycaster-track.mjs
PASS
go test ./games/tracks/raycaster/...
PASS
node scripts/check-lessons.mjs
PASS
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing maze controls retained. |
| Tablet | 768 × 1024 | Touch | Existing on-screen controls retained. |
| Phone | 390 × 844 | Touch | Existing responsive iframe retained. |

- Japanese: labels paths and rules in Japanese.
- English: labels the matching path and rule in English.
- Readability / accessibility: the source layers are literal labelled text.
- Screenshots / recordings: article-generation change; gameplay unchanged.
