# P3-ARCH-01 — Bilingual project-structure guide.

Status: PASS

## Implementation

- [x] Implementation
- Added a compiling `examples/project-structure` refactor with a thin `cmd/reef-run`, Ebitengine presentation under `internal/game`, tested pure logic under `internal/rules`, and an `assets` ownership/readme boundary.
- Published equivalent Japanese and English structure guides.

## Automated checks

- [x] Automated checks

```text
go test ./examples/project-structure/internal/rules
GOOS=js GOARCH=wasm go build ./examples/project-structure/cmd/reef-run
git diff --check
success
```

## Manual review

- [x] Manual review
- The guide points directly to a real refactor rather than a diagram-only recommendation.
