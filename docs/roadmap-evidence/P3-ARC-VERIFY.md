# P3-ARC-VERIFY — arcade-60 を starter→テスト緑→reference 照合まで通し、Mobile 幅の記事可読性を確認する

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

- Files: `graduation/arcade-60/{starter,reference}/`, `web/{ja,en}/graduation/arcade-60/index.html`.
- Behavior: the starter intentionally presents three red tests; reference implements the same `Round` API and makes the equivalent three tests green. Articles tell readers to compare only after their own tests pass.

## Commands and results

```text
go test ./graduation/arcade-60/starter
Expected red: all three TODO-mapped tests fail before author work.
go test ./graduation/arcade-60/reference
Passes: the matching score, timer, and restart contract is green.
go build ./graduation/arcade-60/starter ./graduation/arcade-60/reference
Both Ebitengine programs compile.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Article has separate download, score, timer, and restart sections. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | Short commands and vertical sections remain readable without a side-by-side comparison. |

- Japanese: article preserves exact test identifiers and TODO locations.
- English: matching fresh-workspace and reference-last flow is present.
- Readability / accessibility:
- Screenshots / recordings:
