# P3-SAVE-02 — Bilingual save guide.

Status: PASS

## Implementation

- [x] Implementation
- Published `/ja/guides/save/` and `/en/guides/save/` with the same four promises: versioned data, meaningful autosave points, safe recovery, and scene/model separation.
- Each page gives a concrete `savegame.Autosave` call and recovery order linked to the new shared implementation.

## Automated checks

- [x] Automated checks

```text
go test ./internal/savegame
rg -n "version|Autosave|scene" web/ja/guides/save/index.html web/en/guides/save/index.html
git diff --check
success
```

## Manual review

- [x] Manual review
- The Japanese and English pages cover equivalent concepts and are usable as standalone system-design lessons.
