# P4-NAV-01 — First 30-minute route.

Status: PASS

## Implementation

- [x] Implementation
- Added matching Japanese/English routes from setup through a testable, buildable first graduation game.

## Automated checks

- [x] Automated checks

```text
rg -n "30|setup|graduation" web/ja/guides/first-30-minutes/index.html web/en/guides/first-30-minutes/index.html
git diff --check
success
```

## Manual review

- [x] Manual review
- Four timed steps prevent a newcomer from choosing among every lesson at once.
