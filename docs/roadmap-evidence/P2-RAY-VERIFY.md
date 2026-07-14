# P2-RAY-VERIFY — Raycaster を著者基準で検証する

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

- Files: six regenerated Raycaster lessons in Japanese and English.
- Behavior: validates matching entry paths and authoring rules across all pages.

## Commands and results

```text
12-page path/RULE scan
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
| Desktop | 1440 × 900 | Keyboard + pointer | Existing raycaster controls retained. |
| Tablet | 768 × 1024 | Touch | Existing on-screen controls retained. |
| Phone | 390 × 844 | Touch | Responsive iframe retained. |

- Japanese: six matching paths found.
- English: six matching paths found.
- Readability / accessibility: source paths are labelled text.
- Screenshots / recordings: generator verification only.
