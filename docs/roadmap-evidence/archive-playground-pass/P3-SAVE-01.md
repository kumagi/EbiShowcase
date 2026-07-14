# P3-SAVE-01 — Versioned save model and adapters.

Status: PASS

## Implementation

- [x] Implementation
- Added `internal/savegame`: a versioned JSON envelope with timestamp, typed payload decoding, autosave helper, and explicit unsupported-version rejection.
- Added native `MemoryStore` plus a WASM `localStorage` adapter behind the same `Store` interface. The adapter copies byte slices so callers cannot mutate stored values accidentally.

## Automated checks

- [x] Automated checks
```text
go test ./internal/savegame
GOOS=js GOARCH=wasm go test -c -o /tmp/ebishowcase-savegame.wasm ./internal/savegame
git diff --check
success
```

## Manual review

- [x] Manual review
- The format contains `version`, `updated_at`, and isolated `data`; changing a game payload therefore requires an intentional migration instead of silently decoding an unknown schema.
