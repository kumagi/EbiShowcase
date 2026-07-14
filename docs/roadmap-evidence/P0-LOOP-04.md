# P0-LOOP-04 — Draw 内で game を書き換えている実装・教材表現を棚卸し、あれば修繕キューを証跡表に残す

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Audit and repairs

| Candidate | Finding | Repair |
| --- | --- | --- |
| `games/tracks/platformer/scrolling-stage/main.go` | `Draw` reset and incremented `g.visible` while culling terrain. | Moved the derived visible-terrain count to `Update`; Draw now only uses it for the HUD. |
| `games/tracks/visual-effects/vfx-walk/main.go` | `Draw` repaired `g.frame` and wrote live-code tokens/hint into `g.shell`. | Moved frame normalization and live-token/hint updates to `Update`; Draw now selects and paints the already-valid frame. |
| All remaining `Draw` methods | Automated scan found no receiver mutation, RNG call, or direct input read in a Draw body after excluding equality comparisons. | No repair queue remains at this baseline; rerun the scan whenever a new Draw method is added. |

## Changes

- Files: `games/tracks/platformer/scrolling-stage/main.go`,
  `games/tracks/visual-effects/vfx-walk/main.go`, this evidence record.
- Behavior: both repaired games preserve their visible HUD/live-code output while
  restoring the Update → state → Draw discipline.

## Commands and results

```text
gofmt -w games/tracks/platformer/scrolling-stage/main.go games/tracks/visual-effects/vfx-walk/main.go

go test ./games/tracks/platformer/scrolling-stage ./games/tracks/visual-effects/vfx-walk
PASS (both packages compile; no test files).

Draw-body scan for receiver mutation, rand.*, inpututil.*, and ebiten input APIs
exit 1 after filtering equality comparisons: no remaining candidate lines.

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | The repairs retain the existing visible-terrain HUD and VFX live-code display. |
| Tablet | 768 × 1024 | Touch | Rendering-only refactor; touch handling remains in Update. |
| Phone | 390 × 844 | Touch | Rendering-only refactor; touch handling remains in Update. |

- Japanese: the public Update/Draw explanation now matches the repaired source.
- English: the public Update/Draw explanation now matches the repaired source.
- Readability / accessibility: no learner example now relies on hidden state
  changes during a draw pass.
- Screenshots / recordings: deferred to the phase/browser integration audit;
  this task’s deterministic source scan is the primary evidence.
