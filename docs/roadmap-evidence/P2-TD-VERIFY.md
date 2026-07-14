# P2-TD-VERIFY — Tower Defense を著者基準で検証する

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

- Files: all eight regenerated TD pages in Japanese and English.
- Behavior: validates matching entry source and RULE per page.

## Commands and results

```text
16-page entry/RULE scan
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
| Desktop | 1440 × 900 | Keyboard + pointer | Existing controls retained. |
| Tablet | 768 × 1024 | Touch | Existing placement retained. |
| Phone | 390 × 844 | Touch | Responsive iframe retained. |

- Japanese: 8/8 paths found.
- English: 8/8 paths found.
- Readability / accessibility: paths are labelled text.
- Screenshots / recordings: generation verification only.
