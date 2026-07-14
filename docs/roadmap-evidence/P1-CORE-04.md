# P1-CORE-04 — LEVEL 04 flappy に YOUR FIRST RULE（例: スコア閾で色や速さフラグ）を追加する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `web/{ja,en}/games/flappy/index.html`.
- Behavior: sets a score-threshold `fastMode` flag in Update, then uses it only
  for Update-side pipe motion while Draw only projects its visible result.

## Commands and results

```text
go test ./game
PASS (no test files)
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing flap control unchanged. |
| Tablet | 768 × 1024 | Touch | Existing touch control unchanged. |
| Phone | 390 × 844 | Touch | Existing responsive game shell unchanged. |

- Japanese: names source, score insertion, flag, movement, and verification.
- English: has the same Update-only rule.
- Readability / accessibility: Update/Draw roles are stated in text.
- Screenshots / recordings: documentation-only change.
