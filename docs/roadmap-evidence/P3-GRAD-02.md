# P3-GRAD-02 — Three-room exploration kit.

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
- Added brief, tested room-state starter, and a reference key/three-room/exit game.

## Automated checks

- [x] Automated checks

```text
go test ./graduation/exploration-3rooms/starter
GOOS=js GOARCH=wasm go build ./graduation/exploration-3rooms/reference
git diff --check
success
```

## Manual review

- [x] Manual review
- The project gives a small, complete scene-transition exercise.
