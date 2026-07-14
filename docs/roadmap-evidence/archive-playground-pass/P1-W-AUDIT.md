# P1-W-AUDIT — Audit Rhythm, including audio unlock and timing behavior.

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

- Files inspected: `internal/rhythmcore/core.go`, `internal/rhythmcore/core_test.go`, `internal/rhythmplay/game.go`, `games/tracks/rhythm/ebi-rhythm/main.go`, and `scripts/gen-rhythm-track.mjs`.
- Behavior: Ebi Rhythm Tour has three original generated-tone charts and EASY/HARD selection. START creates an `audio.Player` and calls Play after keyboard/pointer input; the pure timing core handles taps, holds, rolls, grades, combo, and score. The player has particles and shake for hits, session-only score BEST, and a fixed 480×720 canvas.
- Deliberate non-goals / trade-offs: this audit does not modify the rhythm system.

## Commands and results

```text
go test ./internal/rhythmcore
ok github.com/kumagi/EbiShowcase/internal/rhythmcore

Browser audit: `/play/ebi-rhythm/index.html` at 390×844.
The three-song menu, difficulty buttons, and START button rendered without console errors.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1280 × 720 | Source/UI audit | Four keyboard lanes and mouse lanes are implemented; the play scene has timing line, notes, particles, and shake. |
| Tablet | 768 × 1024 | Fixed-canvas layout audit | Portrait aspect suits the game, but UI text remains bitmap English and no calibration/offset control is available. |
| Phone | 390 × 844 | Menu rendering inspected | Three songs, EASY/HARD, and START are reachable, though the 480-wide canvas is scaled down and a top letterbox remains. |

## Exact polish backlog

1. Add a first-gesture audio state: explicit `SOUND READY` / `TAP TO ENABLE SOUND`, a retry when browser audio resumes late, and a visible silent-mode fallback rather than assuming `Play` succeeded.
2. Add timing calibration (early/late offset, reset, visible milliseconds/frames) and show each hit's signed timing result rather than only the grade.
3. Improve final-game readability: large centered judgment, lane press/hold state, a miss flash, song/difficulty/START selection feedback, and a phone-safe play layout without a dead top band.
4. Make all three songs and both difficulties reliably selectable by pointer/touch and keyboard, persist per-song/difficulty BEST locally, and give a finished run a clear replay/next-song path.
5. Localize game menu/HUD/results and pass article language into the WASM URL; current Japanese article still embeds English-only playable UI.
6. Extend core tests with calibration/offset behavior and add testable score/BEST key rules; retain tap/hold/roll coverage.

## Language and accessibility audit

- Japanese: articles are Japanese but the game has no query language support and every menu/HUD/result string uses English bitmap glyphs.
- English: playable text is legible at desktop scale. Audio intent is described in the article, but the player cannot see whether browser playback was actually unlocked.
- Readability: high color contrast and large notes are good; timing feedback lasts briefly and has no signed early/late indication. Touch lane controls exist but lack press-state emphasis and calibration guidance.

## Screenshots / recordings

- Phone menu screenshot at 390×844 showed all three song cards, EASY/HARD, and START; it also showed the top letterbox and small fixed-font labels that P1-W-POLISH must address.

## Follow-up

- Regression risk: preserve `rhythmcore` as audio-independent deterministic logic; audio/HUD must remain an adapter around it.
- Related task IDs: P1-W-POLISH, P1-W-PASS.
