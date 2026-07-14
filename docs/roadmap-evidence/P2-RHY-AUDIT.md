# P2-RHY-AUDIT — Rhythm 全STEPの snippet / main.go / challenge 不一致を証跡表にする

Status: PASS

## Required evidence

- [x] Mismatch inventory
- [x] Edit targets
- [x] Japanese
- [x] English

## Changes

- Files: `scripts/gen-rhythm-track.mjs`, `games/tracks/rhythm/*/main.go`.
- Behavior: records the source mismatch inventory that the next generator task repairs.

## Commands and results

```text
find games/tracks/rhythm -maxdepth 2 -name main.go
PASS; seven real entry packages inspected.

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Audit only. |
| Tablet | 768 × 1024 | Touch | Audit only. |
| Phone | 390 × 844 | Touch | Audit only. |

- Japanese: the current generator emits a generic Japanese challenge.
- English: the current generator emits the same generic English challenge.
- Readability / accessibility: repair must show source paths as text.
- Screenshots / recordings: audit only.

## Mismatch inventory

| STEP | Actual editable entry | Current mismatch | Repair direction |
| --- | --- | --- | --- |
| 01 beat-pulse | `games/tracks/rhythm/beat-pulse/main.go` (`MakeChart` / `Taps`) | pseudo `pressedFrame` judge and generic challenge | entry chart → `internal/rhythmcore` judge excerpt; add four Taps values. |
| 02 falling-notes | `.../falling-notes/main.go` | position formula lacks caller/path | entry chart → `internal/rhythmplay` mapping; add chart row. |
| 03 four-lane-groove | `.../four-lane-groove/main.go` (`p`) | touch formula hides editable pattern | entry `p` → input mapper; add lane value. |
| 04 hold-notes | `.../hold-notes/main.go` (`Cue`) | pseudo hold code; no target | actual Hold cue → core held-state; add Hold cue. |
| 05 drum-roll | `.../drum-roll/main.go` (Roll `Cue`) | generic pseudocode/challenge | actual Roll cue → core event; add Roll cue. |
| 06 chart-difficulty | `.../chart-difficulty/main.go` (`easy` / `hard`) | invented song selection | actual charts → runner; add hard note. |
| 07 ebi-rhythm | `.../ebi-rhythm/main.go` (`chart` / `songs`) | no real chart builder shown | actual builder → result/session; add chart cue. |
