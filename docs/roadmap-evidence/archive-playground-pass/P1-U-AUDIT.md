# P1-U-AUDIT — Audit Reversi and record exact gaps.

Status: PASS

## Required evidence

- [x] Current behavior
- [x] Gap list
- [x] Desktop
- [x] Tablet
- [x] Phone
- [x] Japanese
- [x] English

## Changes

- Files: `internal/reversi/reversi.go`, `internal/reversi/reversi_test.go`, `internal/reversiui/reversiui.go`, `games/tracks/reversi/ebi-reversi/main.go`, and the JA/EN final articles.
- Behavior: the player controls BLACK by clicking/tapping blue legal-move dots; WHITE waits 18 frames then uses a one-ply positional score-map evaluation. Rules, captures, corner weighting, and corner choice have four unit tests.
- Deliberate non-goals / trade-offs: this audit does not yet modify behavior; the following polish task owns the fixes.

## Commands and results

```text
go test ./internal/reversi
ok github.com/kumagi/EbiShowcase/internal/reversi

Browser audit: `/play/ebi-reversi/` at 1280×720, 768×1024, and 390×844; one legal board move was made at desktop/tablet scale and WHITE responded after its short wait.
Canvas had no horizontal DOM overflow at the audited sizes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1280 × 720 (game canvas) | Pointer move on a blue legal-move dot | Rules and CPU response work; game is letterboxed with large black side regions; all in-game text is English debug text; score-map labels overlap its bottom row. |
| Tablet | 768 × 1024 | Pointer-equivalent move at scaled board coordinate | Board accepts the move and CPU responds. Vertical letterboxing is large, the right-side score map is small, and its explanatory labels overlap. |
| Phone | 390 × 844 | Visual reachability inspected; touch handler exists in code | No horizontal overflow, but the 900×720 composition is heavily scaled down. Status/map text is too small to read comfortably and no on-canvas buttons offer a recovery/setting path. |

- Japanese: article is Japanese and explains the score map, but the playable game uses only English strings (`REVERSI 05 / CPU EVALUATION`, `Your turn`, `MAP EVAL`).
- English: article and game strings are understandable, but the article/game language parity is accidental rather than selected by the shared loader.
- Readability / accessibility: `lastFlips` is recorded but not drawn, so moves have no visible flip animation; no regular-play keyboard action exists; Space restarts only after game over; there are no CPU personalities/difficulties, run history, score/BEST, or stage/mode selection.
- Screenshots / recordings: browser captures inspected at 1280×720, 768×1024, and 390×844 during the audit. The phone capture visibly shows the unreadably small right-side information panel; retain these observations as the before-state for P1-U-POLISH.

## Exact polish backlog

1. Split rule state from presentation state so each captured stone has a timed flip/scale animation and the CPU move receives a clear highlight.
2. Offer three named CPU personalities/difficulties that teach a progression (deterministic first/legal, positional greedy, and bounded look-ahead or mobility-aware), while keeping deterministic tests.
3. Add three replayable match goals/modes and a local BEST/result record so finishing a board has a reason to replay.
4. Replace the fixed 900×720 debug layout with a responsive board-first presentation: readable phone HUD, no overlapping score-map labels, and visible restart/difficulty controls.
5. Support keyboard selection/confirm during play as well as mouse and touch; localize all in-game labels from the page/game language.
