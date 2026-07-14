# P2-RHY-VERIFY — Rhythm 全STEPで編集先が本文から辿れ、Desktop/Phone と日英を確認する

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

- Files: all regenerated `web/{ja,en}/tracks/rhythm/*/index.html` pages.
- Behavior: verifies each of seven steps exposes the matching source entry and a RULE.

## Commands and results

```text
for each JA/EN rhythm page: require matching source path and `YOUR FIRST RULE`
PASS — 14 pages
go test ./games/tracks/rhythm/...
PASS
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing rhythm controls retained. |
| Tablet | 768 × 1024 | Touch | Existing touch lanes retained. |
| Phone | 390 × 844 | Touch | Existing responsive iframe retained. |

- Japanese: 7/7 matching entry paths found.
- English: 7/7 matching entry paths found.
- Readability / accessibility: each code source path is searchable text.
- Screenshots / recordings: generator verification; no gameplay regression.
