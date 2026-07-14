# P1-X-AUDIT — Audit Tower Defense and record exact gaps.

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

- Files inspected: `internal/towerdefense/logic.go`, `internal/towerdefense/logic_test.go`, `internal/towerdefenseplay/play.go`, `games/tracks/tower-defense/ebi-defense/main.go`, and `scripts/gen-tower-defense-track.mjs`.
- Behavior: final Ebi Pearl Defense has a first route for waves 1–3, a second route after wave 3, six waves including a boss, three tower kinds, upgrades, score/lives/coins, particles/shake, and a session-only BEST. Keyboard Q/W/E and bottom touch buttons select tower type; taps on open ground place/upgrade, and START launches waves.
- Deliberate non-goals / trade-offs: audit only; later polish owns behavior changes.

## Commands and results

```text
go test ./internal/towerdefense
ok github.com/kumagi/EbiShowcase/internal/towerdefense

Browser audit: `/play/ebi-defense/index.html` at 390×844.
The final game canvas rendered with no WASM console errors.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1280 × 720 | Code/UI audit | Two routes, wave state, tower kinds, and boss loop exist; buttons and map are usable. |
| Tablet | 768 × 1024 | Fixed-canvas layout audit | Portrait canvas fits, but bitmap status/button labels remain too small in a page iframe. |
| Phone | 390 × 844 | Final game menu visual audit | Path, tower-type buttons, and START are reachable, but the whole 480×720 canvas is reduced into a center band with unused black space; text is tiny. |

## Exact polish backlog

1. Turn the existing route transition into three explicit scenario records (coast, cave, boss gate) with distinct goals/starting resources/enemy traits, not merely a map swap halfway through one run.
2. Show enemy intent before it reaches the gate: upcoming type/speed/armor or trait, remaining spawns, boss warning, and a visual urgency cue for the front runner.
3. Replace fixed bitmap labels with responsive desktop/portrait HUD and localized Noto Sans JP text. Add language query to generated iframe URLs.
4. Make every mobile action obvious and large: tower type, placement, upgrade, START, retry, and scenario selection. Make selected range/cost/upgrade result visible without relying on tiny debug text.
5. Persist scenario BEST scores locally and display grade/BEST at completion; retain coins/lives/score as an explainable replay result.
6. Add pure tests for scenario validation and result grade/key logic; retain path and target tests.

## Language and accessibility audit

- Japanese: article text is Japanese but game text (`WAVE`, `LIVES`, `START`, tower buttons, messages, overlays) is English-only; iframe does not provide a language query.
- English: controls are understandable at desktop scale but button labels are bitmap-small on phone.
- Readability: enemy HP bars and selected range help, but no clear intent/urgency state tells learners why a particular enemy should be prioritized. Color alone distinguishes tower choice.

## Screenshots / recordings

- Phone screenshot at 390×844 captured the final game in its pre-wave state: path and bottom tower/START controls work but text and content are reduced by the fixed canvas; use as before-state for P1-X-POLISH.

## Follow-up

- Regression risk: preserve pure `Path.Position` and `SelectFront`; scenario/HUD changes must remain adapters around those rules.
- Related task IDs: P1-X-POLISH, P1-X-PASS.
