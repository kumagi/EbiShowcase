# P4-METRIC — playable と別に authoring 進捗を表示するスクリプトまたは `ralph-loop.sh` 拡張を入れる

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `scripts/report-authoring-metrics.mjs`, `scripts/ralph-loop.sh`.
- Behavior: `ralph-loop.sh status` keeps the frozen playable curriculum count and additionally prints independent authoring counts for Build Track, late Core rules, authored hubs, graduation briefs, and first-30-minutes.

## Commands and results

```text
bash scripts/ralph-loop.sh status
playable: 208/208; authoring: Build Track 4, Core 12/12, hubs 12/12, briefs 6/6, first-30-minutes 2.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese:
- English:
- Readability / accessibility:
- Screenshots / recordings:
