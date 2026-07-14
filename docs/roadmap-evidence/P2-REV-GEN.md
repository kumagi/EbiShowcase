# P2-REV-GEN — `gen-reversi-track.mjs` を同仕様で改修し再生成する

Status: PASS

## Required evidence

- [x] Dual panel
- [x] Unique concept-row
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/gen-reversi-track.mjs`, shared helper, regenerated Reversi pages.
- Behavior: every step has real entry source, shared-rules excerpt, distinct
  data/update/draw concepts, and a board/CPU RULE challenge.

## Commands and results

```text
node scripts/gen-reversi-track.mjs
PASS
go test ./games/tracks/reversi/...
PASS
node scripts/check-lessons.mjs
PASS
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing board controls retained. |
| Tablet | 768 × 1024 | Touch | Existing board tap retained. |
| Phone | 390 × 844 | Touch | Responsive iframe retained. |

- Japanese: source paths and rules generated.
- English: matching source paths and rules generated.
- Readability / accessibility: source layers are labelled text.
- Screenshots / recordings: generation-only change.
