# P1-CORE-05 — LEVEL 05–06（pong / breakout）に各1つの YOUR FIRST RULE を追加する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `web/{ja,en}/games/{pong,breakout}/index.html`.
- Behavior: Pong gets a named first-to-three Update state; Breakout gets a
  named last-life Update rescue state. Both relegate tuning to secondary work.

## Commands and results

```text
go test ./games/core/pong ./games/core/breakout
PASS (no test files)
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing games and controls unchanged. |
| Tablet | 768 × 1024 | Touch | Existing pointer/touch controls unchanged. |
| Phone | 390 × 844 | Touch | Existing responsive game shell unchanged. |

- Japanese: each game names its source, Update condition, state, and check.
- English: carries the equivalent two Update-side authoring routes.
- Readability / accessibility: win/rescue states are described with labels and code.
- Screenshots / recordings: documentation-only change.
