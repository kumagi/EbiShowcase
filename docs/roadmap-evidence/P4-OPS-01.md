# P4-OPS-01 — Feedback thresholds.
Status: PASS
## Implementation
- [x] Implementation
- Documented 20-response triage and 500-response archive behavior and added an executable threshold helper.
## Automated checks
- [x] Automated checks
```text
node scripts/feedback-thresholds.mjs 20
node scripts/feedback-thresholds.mjs 500
git diff --check
success
```
## Manual review
- [x] Manual review
- Archiving preserves history while restoring row 2 as the next form target.
