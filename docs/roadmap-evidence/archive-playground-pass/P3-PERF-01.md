# P3-PERF-01 — Benchmarkable performance examples.

Status: PASS

## Implementation

- [x] Implementation
- Added `internal/perflab` with a reusable output slice for culling, a simple object pool, and grid-based nearby-candidate query.

## Automated checks

- [x] Automated checks

```text
go test ./internal/perflab
go test -bench . -benchmem ./internal/perflab
git diff --check
success
```

## Manual review

- [x] Manual review
- The examples demonstrate where to measure: allocation-free frame filtering and fewer collision candidates before detailed tests.
