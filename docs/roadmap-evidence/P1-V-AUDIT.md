# P1-V-AUDIT — Audit Raycaster and record exact gaps.

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

- Files inspected: `internal/raycasterui/raycasterui.go`, `internal/raycastlogic/raycastlogic.go`, `internal/raycastlogic/raycastlogic_test.go`, and `scripts/gen-raycaster-track.mjs`.
- Behavior: six lessons currently share a fixed 720×480 canvas. The final lesson has one 12×12 maze, two static enemies, a key and exit, center-ray shooting, a minimap, a reset button, and basic WASD/arrow plus limited touch controls. Pure DDA, fish-eye correction, and sprite projection have three unit tests.
- Deliberate non-goals / trade-offs: this is an audit only; the following polish task owns changes.

## Commands and results

```text
go test ./internal/raycastlogic
ok github.com/kumagi/EbiShowcase/internal/raycastlogic

Browser audit: `/play/ebi-raycaster/` at 1280×720 and 390×844 after a fast build.
No WASM console errors were reported.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1280 × 720 | Visual/gameplay audit | The ray-cast corridor, key, exit, minimap, shooting button and keyboard legend are visible. The fixed 3:2 game has side letterboxing but remains playable. |
| Tablet | 768 × 1024 | Code and fixed-canvas layout audit | The 720×480 logical game will scale uniformly rather than rearrange; the tiny HUD/control labels cannot become touch-friendly in portrait. |
| Phone | 390 × 844 | Visual reachability audit | The canvas is centered in a large black portrait letterbox. FPS, map labels, key/status, and the bottom buttons are much too small for reliable play. |

## Exact polish backlog

1. Replace the single fixed maze with three real mission definitions (different map, objective/enemy pattern, and completion condition), and teach mission data before final integration.
2. Add player health, enemy patrol/chase or ranged threat, damage/invulnerability/failure, and immediate red flash/shake/weapon-hit feedback.
3. Add a grade/BEST record and a complete replay loop across missions; track time, accuracy, or damage so a replay has a concrete target.
4. Redesign the final canvas around responsive logical layouts: a readable portrait HUD/minimap, large on-canvas touch controls for turn, forward, reverse/strafe, fire, and reset, with no black dead area dominating the phone screen.
5. Localize the final game HUD, instructions, messages, buttons, and article-to-WASM language query. Current play canvas strings are English even in the Japanese article.
6. Expand pure logic tests to cover mission validation, collision/damage state, deterministic enemy steps, grading, and command/raycast edge cases.

## Language and accessibility audit

- Japanese: the Japanese article exists, but its iframe does not provide a language query and the playable game is English-only.
- English: article and in-game English labels are understandable; the control explanation overstates touch completeness because reverse is keyboard-only and the 390px buttons are tiny.
- Readability: desktop contrast is good, but fixed bitmap `basicfont.Face7x13` text and a static 720×480 layout fail the phone readability target. No focus/cursor status, damage feedback, or failure-state feedback exists.

## Screenshots / recordings

- Desktop screenshot showed a readable corridor, key, minimap, and four bottom controls at 1280×720.
- Phone screenshot at 390×844 showed the whole 720×480 view reduced to a narrow center band with extensive black space above and below; retain this as the before-state for P1-V-POLISH.

## Follow-up

- Regression risk: keep DDA, distance correction, and projection calculations in `internal/raycastlogic`, independent of the responsive Ebitengine presentation layer.
- Related task IDs: P1-V-POLISH, P1-V-PASS.
