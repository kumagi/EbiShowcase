# P3-PUZ-STARTER — `graduation/puzzle-3stages/starter` をデータ駆動進行の穴あき＋テストへ拡充する

Status: PASS

## Required evidence

- [x] Holey starter
- [x] Failing tests
- [x] TODO mapping
- [x] Automated checks

## Changes

- Files: `graduation/puzzle-3stages/{starter,reference}/{main.go,progress.go,progress_test.go}`, `graduation/puzzle-3stages/README.md`.
- Behavior: stage names and targets are data; the starter leaves reset and data-driven advancement as TODOs. Its red tests prove target counting, final-stage clearing, and reset, while reference exposes the same Progress API with green tests.

## Commands and results

```text
go test ./graduation/puzzle-3stages/starter
Expected red: target advance, final-stage clear, and restart fail before TODO completion.
go test ./graduation/puzzle-3stages/reference
Passes.
go build ./graduation/puzzle-3stages/starter ./graduation/puzzle-3stages/reference
Both programs compile.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese: README describes data-driven stages and the TODO route.
- English: Go symbols and commands make the same fresh-workspace flow usable.
- Readability / accessibility:
- Screenshots / recordings:
