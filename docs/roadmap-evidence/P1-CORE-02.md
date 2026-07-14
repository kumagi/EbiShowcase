# P1-CORE-02 — LEVEL 02 timing-meter に YOUR FIRST RULE（例: Perfect 連続で bonus）を追加する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `web/{ja,en}/games/timing-meter/index.html`.
- Behavior: makes a named Perfect-streak bonus the primary Update-side rule;
  speed and band width no longer lead the learner task.

## Commands and results

```text
go test ./games/core/timing-meter
PASS (no test files)

node scripts/inject-ogp.mjs
PASS

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing timing controls remain unchanged. |
| Tablet | 768 × 1024 | Touch | Existing tap control remains unchanged. |
| Phone | 390 × 844 | Touch | Existing responsive shell remains unchanged. |

- Japanese: exact field and insertion context are named.
- English: equivalent field, condition, reset, reward, and test/run route.
- Readability / accessibility: behaviour is described in prose and code, not colour.
- Screenshots / recordings: documentation-only change.
