# P1-X-PASS — Verify Tower Defense on three viewports and add X to the advanced checklist.

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
go test ./internal/towerdefense ./internal/towerdefenseplay
success
GOOS=js GOARCH=wasm go build ./games/tracks/tower-defense/ebi-defense
success
bash scripts/build.sh --fast
success; 208/208 playable lessons and 29/29 VFX lessons pass metadata checks.
```

## Quality criteria audit

- Stages: COAST WATCH, REEF CAVE, and PEARL GATE have different routes, lives/coins, speed traits, wave count, and boss condition; 1/2/3 and visible scenario cards select them.
- Animation/feedback: moving enemies/projectiles, tower shots, target ring, particles, camera shake, HP bars, intent line, selected tower/scenario color, score, grade, and result overlay expose state.
- Replay: scenario-specific BEST is saved locally; R/tap/Enter replays a completed or failed defense.
- Desktop/tablet/phone: the portrait game fits 768×1024 and 390×844 with full-width tower cards and START; desktop keeps pointer range and route readability.
- Keyboard/pointer/touch: Q/W/E + 1/2/3, direct tower/scenario cards, empty-ground placement, tower upgrade, START, and retry cover the complete loop.
- Japanese/English: generator appends the language query to each WASM iframe; Japanese HUD/messages/results render with embedded Noto Sans JP and English uses the default branch.
- Tests: path interpolation, front-target selection, scenario validation, result grades, and WASM build pass.

## Follow-up

- Regression risk: scenario switching must reset route, resources, lives, and BEST key together.
- Related task IDs: P1-X-AUDIT, P1-X-POLISH, P2-G19.
