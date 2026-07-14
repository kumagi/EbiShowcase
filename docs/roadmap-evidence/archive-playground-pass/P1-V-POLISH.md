# P1-V-POLISH — Add three map missions, damage/failure, weapon feedback, grade/BEST, and complete controls.

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `internal/raycastlogic/raycastlogic.go`, `internal/raycastlogic/raycastlogic_test.go`, `internal/raycasterui/raycasterui.go`, `internal/raycasterui/storage_js.go`, `internal/raycasterui/storage_native.go`, and `scripts/gen-raycaster-track.mjs`.
- Behavior: the final lesson reads three Mission data records (SUNSET HALL, TEAL VAULT, and NIGHT SPIRAL), each with its own map, route, enemies, key, exit, and target time. The player has three HP; close enemies cause timed damage and a red impact flash, while shooting gets a center reticle and hit/miss feedback. A successful clear is graded from time, damage, and shots; each mission stores a browser-local BEST time. The portrait game uses a game-view panel plus a large native HUD, mission buttons, replay, and five touch controls.
- Deliberate non-goals / trade-offs: enemies are static guard hazards rather than pathfinding agents; that keeps the lesson centered on DDA/raycasting and makes mission data understandable before adding AI.

## Commands and results

```text
gofmt -w internal/raycastlogic/*.go internal/raycasterui/*.go
go test ./internal/raycastlogic ./internal/raycasterui
ok github.com/kumagi/EbiShowcase/internal/raycastlogic
?  github.com/kumagi/EbiShowcase/internal/raycasterui [no test files]

GOOS=js GOARCH=wasm go build ./games/tracks/raycaster/ebi-raycaster
success

bash scripts/build.sh --fast
success; 208/208 playable lessons and 29/29 VFX lessons passed metadata checks.
```

## Manual review

- Desktop: final game has readable ray-cast view, minimap, HP/key/time, three keyboard-selected missions, reset, shooting feedback, and mouse/touch-compatible controls.
- Phone (390 × 844): inspected after a production build. The game view occupies the upper panel; MISSION, HP/key/time, state message, three mission buttons, REPLAY/RESET, and five full-height touch controls are legible below it. There is no longer a large empty lower black letterbox.
- Tablet (768 × 1024): the portrait layout uses the same board-first HUD and controls without horizontal overflow.
- Japanese/English: generator now supplies `?lang=ja` and `?lang=en` to the WASM iframe. Key, hit, damage, exit, and restart state use the language branch; the final article explains three mission records, grading, and saved BEST.

## Follow-up

- Regression risk: keep Mission validation and Grade in `raycastlogic`; UI changes must not move that logic into Ebitengine callbacks.
- Related task IDs: P1-V-AUDIT, P1-V-PASS.
