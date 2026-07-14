# P1-Y-AUDIT — Audit Top-down Adventure and record exact gaps.

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

- Files inspected: `internal/topdownadventurelogic/logic.go`, `internal/topdownadventurelogic/logic_test.go`, `internal/topdownadventuregame/game.go`, and `scripts/gen-topdown-adventure-track.mjs`.
- Behavior: final dungeon uses four stages (key, sealed fights, tools, boss), 18-frame sword attacks, damage flash/shake, three tools, boss phases, touch movement/attack/tool controls, session BEST, and transition frames. Pure geometry/room/boss logic has tests.
- Deliberate non-goals / trade-offs: audit only.

## Commands and results

```text
go test ./internal/topdownadventurelogic
ok github.com/kumagi/EbiShowcase/internal/topdownadventurelogic
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1280 × 720 | Code/UI audit | Keyboard and on-canvas movement, attack, and tool controls reach complete dungeon logic. |
| Tablet | 768 × 1024 | Layout audit | Portrait canvas fits, but fixed bitmap UI remains small in article iframes. |
| Phone | 390 × 844 | Source/input audit | Touch zones exist, but the final UI lacks a clear mobile-first status/control hierarchy and all strings remain English. |

## Exact polish backlog

1. Make four room data records explicitly communicate each goal, transition, and clear condition; make room transitions readable rather than only a stage integer.
2. Improve movement/attack readability with idle/walk/attack/recovery state feedback and boss tell/telegraph before dash/storm damage.
3. Give phone controls large labeled regions and a persistent HP/score/room/tool HUD; ensure attack/tool/retry are obvious.
4. Persist final BEST locally and grade complete dungeon runs by HP/time/score.
5. Localize HUD, messages, tool names, boss phases, controls, and results with language query and Japanese-capable font.
6. Extend pure tests to room data validation and run grading while keeping attack/room/boss rules isolated.

## Language and accessibility audit

- Japanese: lesson exists but iframe has no language query and game labels/messages are English-only.
- English: systems are clear in code; phone labels need more contrast/size.
- Readability: attack/flash/shake are present, but boss phase text lacks a strong pre-attack tell and mobile controls do not communicate pressed state strongly.

## Follow-up

- Regression risk: retain pure topdownadventurelogic as the authority for attack/room/boss rules.
- Related task IDs: P1-Y-POLISH, P1-Y-PASS.
