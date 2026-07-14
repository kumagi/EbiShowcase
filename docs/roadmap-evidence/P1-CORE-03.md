# P1-CORE-03 — LEVEL 03 catch-stars に YOUR FIRST RULE（例: ミス連続で状態変化）を追加する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `web/{ja,en}/games/catch-stars/index.html`.
- Behavior: defines a miss-streak state transition in Update with catch reset,
  three-miss branch, named source path, and local verification.

## Commands and results

```text
go test ./games/core/catch-stars
PASS (no test files)
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing controls remain unchanged. |
| Tablet | 768 × 1024 | Touch | Existing control path remains unchanged. |
| Phone | 390 × 844 | Touch | Existing responsive shell remains unchanged. |

- Japanese: names both miss and catch insertion contexts.
- English: contains the same state/reset/verification route.
- Readability / accessibility: status is described as text and code.
- Screenshots / recordings: documentation-only change.
