# P2-REV-VERIFY — Reversi を著者基準で検証する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Desktop
- [x] Phone
- [x] Japanese
- [x] English
- [x] Tests

## Changes

- Files: five generated Reversi lessons in both languages.
- Behavior: validates matching entry path and RULE per lesson.

## Commands and results

```text
10-page entry/RULE scan
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
| Desktop | 1440 × 900 | Keyboard + pointer | Existing controls retained. |
| Tablet | 768 × 1024 | Touch | Existing board tap retained. |
| Phone | 390 × 844 | Touch | Responsive iframe retained. |

- Japanese: 5/5 matching paths found.
- English: 5/5 matching paths found.
- Readability / accessibility: paths are labelled text.
- Screenshots / recordings: generation verification only.
