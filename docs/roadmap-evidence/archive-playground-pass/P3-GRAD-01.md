# P3-GRAD-01 — 60-second arcade kit.

Status: PASS

## Implementation

- [x] Implementation
- [x] Article
- [x] Starter
- [x] Tests
- [x] Reference game
- [x] Japanese
- [x] English
- [x] Mobile
- Added an Apache-2.0 brief, copyable pure-rule starter/test, and a complete 60-second reference game with score, timer, movement, and restart.

## Automated checks

- [x] Automated checks

```text
go test ./graduation/arcade-60/starter
GOOS=js GOARCH=wasm go build ./graduation/arcade-60/reference
git diff --check
success
```

## Manual review

- [x] Manual review
- The brief gives a bounded first graduation project and makes reference code optional rather than the first step.
