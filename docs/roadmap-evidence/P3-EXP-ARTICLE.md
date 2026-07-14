# P3-EXP-ARTICLE — 日英の exploration brief を TODO/テスト対応の手順記事にする

Status: PASS

## Required evidence

- [x] Article
- [x] TODO mapping
- [x] Japanese
- [x] English
- [x] Manual review

## Changes

- Files: `web/{ja,en}/graduation/exploration-3rooms/index.html`.
- Behavior: both articles use fresh-workspace browser download instructions and map all four red test names to TODO 2–5 in `state.go`.

## Commands and results

```text
rg "Test(KeyCanBeCollectedOnlyInFirstRoom|PathMovesThroughExactlyThreeRooms|ExitNeedsKeyInThirdRoom|RestartMakesNewAdventure)" web/{ja,en}/graduation/exploration-3rooms/index.html
All eight localized test references are present.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Four ordered state-machine sections are independently readable. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | Each step is vertical and keeps commands short. |

- Japanese: all test identifiers and TODO mappings appear in Japanese instructions.
- English: matching download and implementation sequence appears in English.
- Readability / accessibility:
- Screenshots / recordings:
