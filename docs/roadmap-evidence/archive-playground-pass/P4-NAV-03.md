# P4-NAV-03 — Visible progress links.
Status: PASS
## Implementation
- [x] Implementation
- Added bilingual Play → Deepen → Make progression maps with links to genre tracks, technical labs, and graduation projects.
## Automated checks
- [x] Automated checks
```text
rg -n "PLAY|DEEPEN|MAKE|Graduation" web/ja/guides/progress/index.html web/en/guides/progress/index.html
git diff --check
success
```
## Manual review
- [x] Manual review
- The map makes the next destination visible after a first playable lesson.
