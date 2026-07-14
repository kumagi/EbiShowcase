# P3-PUZ-VERIFY — puzzle-3stages の著者フローを検証する

Status: PASS

## Required evidence

- [x] Article
- [x] Starter
- [x] Tests
- [x] Reference game
- [x] Japanese
- [x] English
- [x] Mobile

## Changes

- Files: `graduation/puzzle-3stages/{starter,reference}/`, `web/{ja,en}/graduation/puzzle-3stages/index.html`.
- Behavior: starter intentionally leaves progression TODOs red; reference applies the same StageData / Progress API and is green. Articles cover browser download, test mapping, and reference-last use.

## Commands and results

```text
go test ./graduation/puzzle-3stages/starter
Expected red: three unfinished progression tests fail.
go test ./graduation/puzzle-3stages/reference
Passes.
go build ./graduation/puzzle-3stages/starter ./graduation/puzzle-3stages/reference
Both programs compile.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Data, advance, and reset are separate readable sections. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | The articles preserve the sequential test-first reading order. |

- Japanese: localized path and exact test/TODO mapping are present.
- English: same no-clone authoring flow is present.
- Readability / accessibility:
- Screenshots / recordings:
