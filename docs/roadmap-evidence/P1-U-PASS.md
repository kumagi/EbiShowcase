# P1-U-PASS — Verify Reversi on three viewports and add U to the advanced checklist.

Status: PASS

## Required evidence

- [x] Stages
- [x] Animation
- [x] Feedback
- [x] Replay
- [x] Keyboard
- [x] Pointer
- [x] Touch
- [x] Japanese
- [x] English
- [x] Tests

## Changes

- Files: `docs/ADVANCED_QUALITY_CHECKLIST.md`, `docs/roadmap-evidence/P1-U-PASS.md`.
- Behavior: U is now recorded as an advanced-quality track. Its three complete, replayable CPU encounters are FRIENDLY (first legal move), POSITION (one-ply score map), and SCOUT (bounded reply-aware search), not palette-only changes. The five learning steps lead from board data to the full encounter.

## Commands and results

```text
go test ./...
success

bash scripts/build.sh --fast
success
OK — 208/208 gated, 29/29 VFX, and 66 home cards are linked and bilingual.

node scripts/roadmap-ralph-loop.mjs verify-task P1-U-POLISH
P1-U-POLISH: evidence OK
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1280 × 720 | Pointer selected a legal dot; CPU replied | Landscape board, CPU controls, score map, replay control, and localized status are legible. |
| Tablet | 768 × 1024 | Layout / reachability review | Responsive board-first composition remains readable without horizontal overflow. |
| Phone | 390 × 844 | Touch legal move, then CPU response; keyboard `R`, Arrow Right, Space | Touch move changed BLACK 02/WHITE 02 to BLACK 03/WHITE 03 after a WHITE reply; the cursor moves by keyboard and reset is available. |

## Quality criteria audit

- Stages: five playable lessons (`board-grid`, `legal-moves`, `flip-stones`, `pass-and-score`, `ebi-reversi`) isolate the systems. The final encounter has three deliberately different CPU behaviors, selectable by button or 1/2/3.
- Animation: captured stones use staggered scale/recolor flip frames; the latest move pulses with a gold outline.
- Feedback: blue legal dots, yellow keyboard cursor, turn/count HUD, CPU hint, score map, selected difficulty, final result, and REPLAY feedback make state changes visible.
- Replay: R/Enter/REPLAY restarts a game; each difficulty saves its own best winning margin in browser local storage.
- Japanese: `/play/ebi-reversi/?lang=ja` was inspected; title, CPU rule, status, map note, and replay hint were Japanese.
- English: generated English article points its iframe to `?lang=en`; labels use the English branch of the same WASM source.
- Tests: pure rules, capture, evaluation, and deterministic look-ahead selection pass through `go test ./...`; the final WASM package compiles under `GOOS=js GOARCH=wasm` during the site build.

## Follow-up

- Regression risk: changes to the final Reversi UI must retain the small `reversi` rules package as the authoritative testable logic layer.
- Related task IDs: P1-U-AUDIT, P1-U-POLISH, P2-G16.
