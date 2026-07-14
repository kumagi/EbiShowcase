# P1-CORE-LINK — Core 前半ページと Build Track・testing ガイドを相互リンクし、first-30 草案の仮リンクを置く

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/inject-core-authoring-links.mjs`, `scripts/build.sh`, and
  Core LEVEL 01–06 JA/EN pages.
- Behavior: a build-owned, idempotent block links every early Core lesson to
  Build Track, Testing, and the first-30-minutes writing route.

## Commands and results

```text
node scripts/inject-core-authoring-links.mjs
PASS
node scripts/insert-feedback-form.mjs
PASS
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Pointer | Three writing links appear after lesson content. |
| Tablet | 768 × 1024 | Touch | Shared responsive link panel applies. |
| Phone | 390 × 844 | Touch | Standard anchors preserve touch navigation. |

- Japanese: points to Build Track, unit testing, and first-30-minutes.
- English: points to the same three routes.
- Readability / accessibility: labelled text links are used.
- Screenshots / recordings: navigation-only change.
