# P1-U-POLISH — Add three CPU personalities/difficulties, move/flip animation, replay records, and complete controls.

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review
- [x] Three CPU personalities/difficulties
- [x] Move / flip animation
- [x] Replay record
- [x] Keyboard, pointer, and touch controls

## Changes

- Files: `internal/reversi/reversi.go`, `internal/reversi/reversi_test.go`, `internal/reversiui/reversiui.go`, `internal/reversiui/storage_js.go`, `internal/reversiui/storage_native.go`, `internal/ogfont/ogfont.go`, and `scripts/gen-reversi-track.mjs`.
- Behavior: the final Reversi lesson has three deterministic opponents: FRIENDLY takes the first legal move, POSITION greedily evaluates the score map, and SCOUT chooses a move after checking the opponent's strongest immediate reply. Captured stones shrink and recolor in staggered flip animations; the last move pulses. The per-opponent best winning margin is stored in browser local storage. The board now has a landscape and a board-first portrait layout, localized in-game text, on-canvas difficulty/replay buttons, and keyboard cursor placement as well as mouse/touch input.
- Deliberate non-goals / trade-offs: SCOUT is deliberately only two plies (CPU move plus the opponent reply). It is understandable, deterministic, and responsive in WebAssembly; deeper search belongs to a later AI lesson.

## Commands and results

```text
gofmt -w internal/reversi/reversi.go internal/reversi/reversi_test.go internal/ogfont/ogfont.go internal/reversiui/*.go
go test ./internal/reversi ./internal/reversiui
ok github.com/kumagi/EbiShowcase/internal/reversi
?  github.com/kumagi/EbiShowcase/internal/reversiui [no test files]

GOOS=js GOARCH=wasm go build ./games/tracks/reversi/ebi-reversi
success

node scripts/gen-reversi-track.mjs && go test ./... && bash scripts/build.sh --fast
success; all Go packages pass, 208/208 playable lessons and 29/29 VFX lessons pass metadata checks.
```

## Viewport and input audit

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1280 × 720 | Pointer, difficulty buttons | The board is large and readable beside the CPU panel; Japanese labels and all three CPU buttons render without overlap. |
| Tablet | 768 × 1024 | Responsive layout inspected | The same logical game changes to a vertically stacked, board-first layout with no horizontal overflow. |
| Phone | 390 × 844 | Touch: placed a legal C3 move; CPU replied | Board occupies almost the full width, legal dots and cursor are touchable, score map moves below the board, and text remains readable. |

Keyboard was also checked at 390 × 844: `R` reset the match, Arrow Right moved the yellow cursor, and Space attempted placement only when the selected cell was legal.

## Language and accessibility audit

- Japanese: `/play/ebi-reversi/?lang=ja` rendered Japanese title, status, CPU hint, score-map explanation, and replay instruction.
- English: generated English final article embeds `/play/ebi-reversi/?lang=en`; the same complete labels are selected by the WASM query parameter.
- Readability: high-contrast legal dots, cursor, selected difficulty fill, last-move pulse, visible REPLAY button, and no unintended phone horizontal scroll were checked. The score map uses a separate portrait row so it cannot collide with board rows.

## Screenshots / recordings

- Desktop game canvas was inspected at 1280 × 720 after WASM startup.
- Phone game canvas was inspected at 390 × 844 before and after a touch move and CPU response; this captured the board-first responsive state rather than a title screen.

## Follow-up

- Regression risk: `syscall/js` local-storage access is behind a WASM build tag; the native test build intentionally uses a no-op implementation.
- Related task IDs: P1-U-AUDIT, P1-U-PASS.
