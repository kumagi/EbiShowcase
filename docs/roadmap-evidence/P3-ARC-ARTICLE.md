# P3-ARC-ARTICLE — 日英の arcade-60 解説を、テスト名と TODO 位置が対応する手順記事にする

Status: PASS

## Required evidence

- [x] Article
- [x] TODO mapping
- [x] Japanese
- [x] English
- [x] Manual review

## Changes

- Files: `web/ja/graduation/arcade-60/index.html`, `web/en/graduation/arcade-60/index.html`.
- Behavior: both articles begin in a fresh Go workspace (no clone), link to the browser-download starter, and pair the three red test names with TODO 4, TODO 5, and TODO 3 in the exact repair order.

## Commands and results

```text
rg "TestActionAddsOneStar|TestTimeLimitEndsRound|TestRestartMakesFreshRound" web/{ja,en}/graduation/arcade-60/index.html
All six article/test references are present.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Each test/TODO step is an independent readable section. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | Sequential sections and short command lines avoid a required side-by-side layout. |

- Japanese: localized instructions use the same exact test names and source locations.
- English: matching article explains red → green workflow and reference-last policy.
- Readability / accessibility:
- Screenshots / recordings:
