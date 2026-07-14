# P3-PUZ-ARTICLE — 日英の puzzle brief を TODO/テスト対応の手順記事にする

Status: PASS

## Required evidence

- [x] Article
- [x] TODO mapping
- [x] Japanese
- [x] English
- [x] Manual review

## Changes

- Files: `web/{ja,en}/graduation/puzzle-3stages/index.html`.
- Behavior: articles start from browser-downloaded starter files, preserve no-clone setup, and map all four tests to TODO 1–3 in the data-driven Progress code.

## Commands and results

```text
rg "Test(FreshProgressStartsAtFirstDataStage|PuzzleTargetAdvancesOneStage|LastDataStageClearsInsteadOfOverflowing|RestartMakesFreshProgress)" web/{ja,en}/graduation/puzzle-3stages/index.html
All test names are present in both localized articles.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Ordered data/progression steps are individually readable. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | Vertical sections and short commands preserve reading order. |

- Japanese: localized data-first path and exact TODO mapping are present.
- English: matching browser-download and test-first path is present.
- Readability / accessibility:
- Screenshots / recordings:
