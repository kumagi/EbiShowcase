# P3-PERF-02 — Bilingual performance guide.

Status: PASS

## Implementation

- [x] Implementation
- Published matching Japanese and English guides for reusable slices, pooling, spatial grids, and benchmark interpretation.

## Automated checks

- [x] Automated checks

```text
go test -bench . -benchmem ./internal/perflab
rg -n "allocs/op|Cull|Grid" web/ja/guides/performance/index.html web/en/guides/performance/index.html
git diff --check
success
```

## Manual review

- [x] Manual review
- Pages identify the measured commands and separate allocation, culling, and broad-phase search concerns.
