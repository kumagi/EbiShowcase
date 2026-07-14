# P1-Y-POLISH — Make the top-down adventure readable, replayable, and localized.

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `internal/topdownadventurelogic/logic.go`, `internal/topdownadventurelogic/logic_test.go`, `internal/topdownadventuregame/game.go`, `internal/topdownadventuregame/storage_js.go`, `internal/topdownadventuregame/storage_native.go`, and `scripts/gen-topdown-adventure-track.mjs`.
- Behavior: the final dungeon now has four explicit room records (key gate, crawler seal, tool seals, guardian) and presents the current room goal on every transition. Sword attacks retain anticipation and contact frames; the DASH phase adds a red expanding warning ring before its speed burst, while STORM retains its expanding danger ring. The persistent HUD and six 74×70 touch controls expose HP, score, room, selected tool, attack, and tool selection at all times.
- Replay: completed runs receive S/A/B/C from pure score/HP/frame rules and save BEST by lesson in browser local storage.
- Language: every embedded adventure lesson passes `lang=ja` or `lang=en`; Noto Sans JP renders Japanese HUD, controls, room goals, boss tells, overlays, and tool names.

## Commands and results

```text
gofmt -w internal/topdownadventurelogic/*.go internal/topdownadventuregame/*.go
go test ./internal/topdownadventurelogic ./internal/topdownadventuregame
ok github.com/kumagi/EbiShowcase/internal/topdownadventurelogic
?  github.com/kumagi/EbiShowcase/internal/topdownadventuregame [no test files]

GOOS=js GOARCH=wasm go build ./games/tracks/topdown-adventure/ebi-adventure
success

bash scripts/build.sh --fast && go test ./...
success
```

## Manual review

- Desktop keeps WASD/arrows, X/Space attack, Q/C tool selection, readable action rectangle, and the DASH boss tell.
- Tablet and phone keep the full-width portrait HUD plus six direct touch regions for movement, attack, and tool selection; tap/Enter/R retries a result.
- Japanese and English final lesson iframes pass their language query, and the Japanese branch uses the embedded Noto face rather than missing-glyph boxes.

## Follow-up

- Regression risk: keep `RoomSpec`, `ValidDungeonRoute`, and `RunGrade` free of Ebitengine so route/grade changes remain unit-testable.
- Related task IDs: P1-Y-AUDIT, P1-Y-PASS.
