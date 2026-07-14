# P1-W-POLISH — Add timing calibration, readable judgement feedback, robust song selection, and complete controls.

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `internal/rhythmcore/core.go`, `internal/rhythmcore/core_test.go`, `internal/rhythmplay/game.go`, `internal/rhythmplay/storage_js.go`, `internal/rhythmplay/storage_native.go`, and `scripts/gen-rhythm-track.mjs`.
- Behavior: timing correction is now a pure, signed frame offset passed to `NewSessionWithOffset`; the menu offers persistent `OFFSET − / +` controls. A grade now gives a high-contrast signed EARLY/LATE/ON TIME reading. Menu and result text use embedded Noto Sans JP and the generated iframe supplies a language query. The game stores BEST per song/difficulty, shows SOUND READY / SOUND ON state, and provides a visible keyboard/touch-accessible silent-practice fallback that retains all visual judgement and gameplay.
- Deliberate non-goals / trade-offs: calibration is in frames rather than milliseconds because the chart engine is frame-based; the menu intentionally presents the frame unit so learners can connect calibration directly to the tested rule.

## Commands and results

```text
gofmt -w internal/rhythmcore/*.go internal/rhythmplay/*.go
go test ./internal/rhythmcore ./internal/rhythmplay
ok github.com/kumagi/EbiShowcase/internal/rhythmcore
?  github.com/kumagi/EbiShowcase/internal/rhythmplay [no test files]

GOOS=js GOARCH=wasm go build ./games/tracks/rhythm/ebi-rhythm
success

node scripts/gen-rhythm-track.mjs
success; JA final article embeds /play/ebi-rhythm/?lang=ja.
```

## Manual review

- Phone 390 × 844: Japanese menu was visually inspected. Song cards, timing-adjustment controls, EASY/HARD, START, and the sound-state cue were readable and Japanese glyphs rendered correctly.
- Pointer/touch: song card, difficulty, timing adjustment, silent-practice strip, START, and each play lane have direct pointer/touch hit regions. Keyboard equivalents are Left/Right, Up/Down/X, [ ], M, Enter/Space, and D/F/J/K.
- Audio: the initial state is `SOUND READY`; START runs from a user gesture and starts the original synthesized beat. `M`/the visible silent-practice strip provides a guaranteed visual-only loop if sound cannot be used.
- Replay: the completion overlay returns to song selection by tap, Enter, or Space. BEST score is stored per `song/difficulty` key in browser local storage.

## Follow-up

- Regression risk: preserve `NewSessionWithOffset` and its test; UI calibration must only set the core offset rather than changing judgement thresholds ad hoc.
- Related task IDs: P1-W-AUDIT, P1-W-PASS.
