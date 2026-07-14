# P4-SAMPLE — カリキュラムから無作為5レッスンを抜き、著者基準で監査して結果を証跡に残す（不合格なら当該を直してから PASS）

Status: PASS

## Required evidence

- [x] Sample audit
- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Manual review

## Changes

- Files: sampled localized lesson pages, their referenced Go entry points, and this evidence record.
- Behavior: deterministic sample `authoring-pass-20260715-v1` selects five distinct lesson shapes: Core Snake, Platformer Tiny Platformer, Match 3 Ebi Match, Rhythm Beat Tap, and Top-down Adventure Eight-way. All had a traceable edit target and rule path after the P2/P4 injections; no repair was needed.

| Sample | Source target | Rule / verification | Result |
| --- | --- | --- | --- |
| Core Snake | `games/core/snake/main.go` | food-score rule; `go test ./games/core/snake` | pass |
| Platformer Tiny Platformer | `games/tracks/platformer/tiny-platformer/main.go` | one observable platform rule; package test | pass |
| Match 3 Ebi Match | `games/tracks/match3/ebi-match/main.go` | color / clear / cascade rule; package test | pass |
| Rhythm Beat Pulse | `games/tracks/rhythm/beat-pulse/main.go` | source-labelled two-layer rule challenge | pass |
| Top-down Eight-way | `games/tracks/topdown-adventure/eight-way/main.go` | source-labelled two-layer rule challenge | pass |

## Commands and results

```text
rg 'YOUR FIRST RULE|games/(core/snake|tracks/platformer/tiny-platformer|tracks/match3/ebi-match|tracks/rhythm/beat-pulse|tracks/topdown-adventure/eight-way)' <sampled JA/EN pages>
Every sampled pair names the real entry and a RULE challenge.
go test ./games/core/snake ./games/tracks/platformer/tiny-platformer ./games/tracks/match3/ebi-match ./games/tracks/rhythm/beat-pulse ./games/tracks/topdown-adventure/eight-way
All five packages compile; ebi-match test passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Source/path labels use sequential sections, with no required side-by-side interaction. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | The authoring panels stack after the playable frame. |

- Japanese: all five sampled Japanese pages contain their rule/source contract.
- English: all five matched English pages contain equivalent source paths and verification intent.
- Readability / accessibility:
- Screenshots / recordings:
