# P1-X-POLISH — Ensure three meaningful scenarios, readable enemy intent, touch building/upgrades, and replay scoring.

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `internal/towerdefense/logic.go`, `internal/towerdefense/logic_test.go`, `internal/towerdefenseplay/play.go`, `internal/towerdefenseplay/storage_js.go`, `internal/towerdefenseplay/storage_native.go`, and `scripts/gen-tower-defense-track.mjs`.
- Behavior: COAST WATCH, REEF CAVE, and PEARL GATE are now independent scenario records with routes, resources, lives, speed traits, wave count, goals, and optional boss. The final game selects them with 1/2/3 or top scenario buttons; intent shows the expected trait, remaining spawn count, and goal. Scenario BEST is browser-persistent and final grade derives from score/lives/coins. The language query and embedded Japanese-capable font drive a localized HUD, messages, and overlays.
- Deliberate non-goals / trade-offs: scenario names remain English proper names; gameplay intent, HUD, tower choices, messages, and results localize.

## Commands and results

```text
gofmt -w internal/towerdefense/*.go internal/towerdefenseplay/*.go
go test ./internal/towerdefense ./internal/towerdefenseplay
ok github.com/kumagi/EbiShowcase/internal/towerdefense
?  github.com/kumagi/EbiShowcase/internal/towerdefenseplay [no test files]

GOOS=js GOARCH=wasm go build ./games/tracks/tower-defense/ebi-defense
success

bash scripts/build.sh --fast
success; 208/208 playable lessons and 29/29 VFX lessons passed metadata checks.
```

## Manual review

- Desktop/tablet/phone retain Q/W/E plus visible tower cards, direct empty-ground placement, direct tower upgrade, START/retry, and scenario selection.
- Scenario buttons use both keyboard (1/2/3) and pointer/touch; selected scenario changes the route, resource/life budget, enemy speed trait, waves, and boss rule.
- Intent is shown before spawning and target rings remain in-world; completion reports a replay grade and the scenario-specific BEST.

## Follow-up

- Regression risk: scenario records must validate through `towerdefense.ValidScenario`; do not reintroduce route changes as a hard-coded wave-number special case.
- Related task IDs: P1-X-AUDIT, P1-X-PASS.
