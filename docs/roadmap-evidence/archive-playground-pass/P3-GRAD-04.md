# P3-GRAD-04 — Graduation hub.

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
- Added bilingual graduation hub with prerequisites, each brief/source location, and project/release guide links.

## Automated checks

- [x] Automated checks

```text
rg -n "arcade|exploration|puzzle|release" web/ja/graduation/index.html web/en/graduation/index.html
git diff --check
success
```

## Manual review

- [x] Manual review
- A learner can choose a bounded project and find its source plus release steps from one page.
