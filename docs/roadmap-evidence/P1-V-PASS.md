# P1-V-PASS — Verify Raycaster on three viewports and add V to the advanced checklist.

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

- Files: `docs/ADVANCED_QUALITY_CHECKLIST.md`, `docs/roadmap-evidence/P1-V-PASS.md`.
- Behavior: V is recorded as an advanced-quality track. The final game has three selectable, data-driven escape missions and supports a full replay loop with HP/failure, grades, and per-mission browser BEST time.

## Commands and results

```text
go test ./...
success

GOOS=js GOARCH=wasm go build ./games/tracks/raycaster/ebi-raycaster
success

bash scripts/build.sh --fast
success
OK — 208/208 gated, 29/29 VFX, and 66 home cards are linked and bilingual.
```

## Quality criteria audit

- Stages: six playable lessons lead from movement to DDA, projected strips, columns, and texture/fisheye correction. Final mission buttons 1/2/3 select SUNSET HALL, TEAL VAULT, and NIGHT SPIRAL, whose maps, routes, enemies, and target times differ.
- Animation: a shot runs for seven frames with an expanding gold muzzle/reticle pulse; damage has a 14-frame red flash.
- Feedback: clear/miss messages, minimap, crosshair, HP/key/time HUD, damage flash, system-down state, mission grade, and selected-mission color all expose the game state.
- Replay: R, desktop reset, and portrait REPLAY/RESET begin a new mission; BEST time is preserved per mission in local storage.
- Desktop: 1280 × 720 audit found readable FPS, minimap, HUD, keyboard controls, and visible on-canvas controls.
- Tablet: 768 × 1024 uses the portrait HUD rather than scaling the prior 720×480 canvas into unreadable letterboxing.
- Phone: 390 × 844 screenshot verified a game panel plus readable MISSION/HP/KEY/TIME/message, three mission buttons, replay control, and five large touch targets. Pressing FIRE produced `THE RAY HIT NO ENEMY`, verifying a touch action reaches the game loop.
- Japanese: `/play/ebi-raycaster/index.html?lang=ja` was visually checked after switching to embedded Noto Sans JP; Japanese labels render instead of glyph boxes. English uses the same WASM with `?lang=en` and was checked in the previous phone audit.
- Tests: pure DDA/projection, Mission validation, and Grade tests pass; the WebAssembly final package builds.

## Follow-up

- Regression risk: preserve the language query in generated iframes and the embedded Japanese-capable font whenever the game HUD changes.
- Related task IDs: P1-V-AUDIT, P1-V-POLISH, P2-G17.
