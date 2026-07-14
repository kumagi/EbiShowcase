# P3-DIST-01 — Bilingual distribution guide.

Status: PASS

## Implementation

- [x] Implementation
- Published Japanese and English release checklists covering local tests, WASM compilation, GitHub Pages deployment, Apache-2.0/third-party notices, and release metadata.

## Automated checks

- [x] Automated checks

```text
rg -n "WASM|Pages|Apache|THIRD_PARTY" web/ja/guides/release/index.html web/en/guides/release/index.html
git diff --check
success
```

## Manual review

- [x] Manual review
- The two pages give equivalent, actionable publishing order.
