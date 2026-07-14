# P3-EXP-STARTER — `graduation/exploration-3rooms/starter` を状態・鍵・遷移の穴あき＋テストへ拡充する

Status: PASS

## Required evidence

- [x] Holey starter
- [x] Failing tests
- [x] TODO mapping
- [x] Automated checks

## Changes

- Files: `graduation/exploration-3rooms/{starter,reference}/{main.go,state.go,state_test.go}`, `graduation/exploration-3rooms/README.md`.
- Behavior: the starter separates testable room/inventory/ending state from Ebitengine input and drawing. Four red tests map 1:1 to restart, key, transition, and exit TODOs; reference uses the same API with green tests.

## Commands and results

```text
go test ./graduation/exploration-3rooms/starter
Expected red: key, bounded path, gated exit, and restart tests fail.
go test ./graduation/exploration-3rooms/reference
Passes.
go build ./graduation/exploration-3rooms/starter ./graduation/exploration-3rooms/reference
Both Ebitengine programs compile.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese: README directs a fresh Go workspace and lists the Japanese concept mapping.
- English: source identifiers and commands remain readable and downloadable without a clone.
- Readability / accessibility:
- Screenshots / recordings:
