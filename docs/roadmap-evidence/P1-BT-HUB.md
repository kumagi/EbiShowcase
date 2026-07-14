# P1-BT-HUB — Build Track ハブとホーム／progress／setup からの導線を日英で追加する（ゲート数は増やさない）。ハブ先頭に公理の3条を置く

Status: PASS

## Required evidence

- [x] Edit target
- [x] Next lines
- [x] RULE challenge
- [x] Desktop
- [x] Phone
- [x] Japanese
- [x] English
- [x] Tests

## Changes

- Files: `web/{ja,en}/build/index.html`, `web/{ja,en}/index.html`,
  `web/{ja,en}/guides/progress/index.html`, `scripts/gen-setup-guide.mjs`.
- Behavior: adds bilingual Build Track hubs with the three axioms first, all
  four STEPs, and LEVEL 01. Home, progress, and generated setup now point into
  the clone-free writing route without changing the 208 gate.

## Commands and results

```text
node scripts/gen-setup-guide.mjs
PASS

node scripts/insert-feedback-form.mjs
PASS

node scripts/inject-ogp.mjs
PASS

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Pointer | Hub exposes all four cards and LEVEL 01 link. |
| Tablet | 768 × 1024 | Touch | Shared card grid reflows to the existing responsive layout. |
| Phone | 390 × 844 | Touch | Links remain normal anchors with touch-safe shared styling. |

- Japanese: three axioms and all entry routes are present.
- English: matches the same four-step writing route and no-clone promise.
- Readability / accessibility: cards use ordinary labelled links; no new state
  is conveyed only by colour.
- Screenshots / recordings: static navigation change; no game visual changed.
