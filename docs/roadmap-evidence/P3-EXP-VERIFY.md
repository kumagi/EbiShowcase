# P3-EXP-VERIFY — exploration-3rooms の著者フローを検証する

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

- Files: `graduation/exploration-3rooms/{starter,reference}/`, `web/{ja,en}/graduation/exploration-3rooms/index.html`.
- Behavior: starter starts with four purposeful red tests; reference implements the identical State API and is green. Both articles give a no-clone, test-first authoring route.

## Commands and results

```text
go test ./graduation/exploration-3rooms/starter
Expected red: all four TODO-mapped state tests fail before implementation.
go test ./graduation/exploration-3rooms/reference
Passes.
go build ./graduation/exploration-3rooms/starter ./graduation/exploration-3rooms/reference
Both Ebitengine programs compile.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Ordered source/article sections keep each state rule independent. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | The article's vertical steps avoid a horizontal comparison requirement. |

- Japanese: localized article includes all source paths, test identifiers, and TODO locations.
- English: matching authoring route is present in the English article.
- Readability / accessibility:
- Screenshots / recordings:
