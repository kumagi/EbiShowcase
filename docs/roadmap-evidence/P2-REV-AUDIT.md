# P2-REV-AUDIT — Reversi 全STEPの不一致表を作る

Status: PASS

## Required evidence

- [x] Mismatch inventory
- [x] Edit targets
- [x] Japanese
- [x] English

## Changes

- Files: `scripts/gen-reversi-track.mjs`, five Reversi entry packages.
- Behavior: records generic board/CPU snippets versus actual entry wrappers.

## Commands and results

```text
find games/tracks/reversi -maxdepth 2 -name main.go
PASS — five entries.
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Audit only. |
| Tablet | 768 × 1024 | Touch | Audit only. |
| Phone | 390 × 844 | Touch | Audit only. |

- Japanese: generic ScoreMap/change challenge has no named entry location.
- English: same missing authoring route.
- Readability / accessibility: repair exposes literal source paths.
- Screenshots / recordings: audit only.

## Mismatch inventory

All five entries (`board-grid`, `legal-moves`, `flip-stones`, `pass-and-score`,
`ebi-reversi`) are actual `games/tracks/reversi/<slug>/main.go` wrappers around
shared Reversi logic. Current articles show illustrative rules and generic
challenges. Repair: entry source first, labelled shared Reversi mechanism
second, then a RULE that adds a board/move/evaluation data case and verification.
