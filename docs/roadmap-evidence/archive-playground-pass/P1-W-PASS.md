# P1-W-PASS — Verify all three songs/two difficulties and add W to the advanced checklist.

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

## Commands and results

```text
go test ./...
success

GOOS=js GOARCH=wasm go build ./games/tracks/rhythm/ebi-rhythm
success

bash scripts/build.sh --fast
success
OK — 208/208 gated, 29/29 VFX, and 66 home cards are linked and bilingual.
```

## Quality criteria audit

- Stages: the final config defines SUNRISE HARBOR, NEON REEF, and TEMPEST PARADE. Each has distinct BPM/tone/chart data and both EASY and HARD chart records; keyboard and direct menu cards select all six combinations.
- Animation: beat pulse, falling notes, hold tails, roll rings, hit particles, six-frame shake, and seven-frame judgement display expose intermediate motion.
- Feedback: large PERFECT/GOOD/MISS plus signed EARLY/LATE/ON TIME frames, score/combo, lane feedback, menu selection color, result rank, and sound state make every important state readable.
- Replay: finished songs return to selection through touch/Enter/Space; BEST persists separately for every song/difficulty combination.
- Desktop/tablet/phone: 1280 × 720, 768 × 1024, and 390 × 844 layouts retain all lanes and menu controls. The phone menu was visually inspected after the Noto conversion.
- Keyboard/pointer/touch: all complete actions have the three input paths; silent practice additionally guarantees a visual game loop when audio is unavailable.
- Japanese/English: generated JA/EN articles pass the matching language query to WASM. Japanese Noto Sans JP menu labels and controls were inspected at phone size; English remains the default branch.
- Tests: tap, hold, roll, grade windows, and offset compensation are deterministic `rhythmcore` tests; the production WASM final game compiles.

## Follow-up

- Regression risk: if chart timing units change, update both `NewSessionWithOffset` tests and the menu label rather than silently changing the correction meaning.
- Related task IDs: P1-W-AUDIT, P1-W-POLISH, P2-G18.
