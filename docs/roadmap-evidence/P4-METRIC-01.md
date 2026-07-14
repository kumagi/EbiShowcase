# P4-METRIC-01 — Metrics report.
Status: PASS
## Implementation
- [x] Implementation
- Added `scripts/report-metrics.mjs` for playable tracks, advanced/mobile counts, labs, graduation projects, and unchecked feedback.
## Automated checks
- [x] Automated checks
```text
node scripts/report-metrics.mjs
git diff --check
success
```
## Manual review
- [x] Manual review
- Output is machine-readable JSON for README/dashboard reuse.
