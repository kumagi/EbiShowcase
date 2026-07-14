# P2-TD-GEN — `gen-tower-defense-track.mjs` を同仕様で改修し再生成する

Status: PASS

## Required evidence

- [x] Dual panel
- [x] Unique concept-row
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/gen-tower-defense-track.mjs`, shared authoring helper, and
  regenerated TD pages.
- Behavior: entries, shared-engine excerpts, authoring concept rows, and
  map/path/wave RULEs are generated per step.

## Commands and results

```text
node scripts/gen-tower-defense-track.mjs
PASS
go test ./games/tracks/tower-defense/...
PASS
node scripts/check-lessons.mjs
PASS
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing TD controls retained. |
| Tablet | 768 × 1024 | Touch | Existing tower placement retained. |
| Phone | 390 × 844 | Touch | Existing responsive iframe retained. |

- Japanese: literal entry paths and rules generated.
- English: matching entry paths and rules generated.
- Readability / accessibility: source layers are labelled text.
- Screenshots / recordings: generation-only change.
