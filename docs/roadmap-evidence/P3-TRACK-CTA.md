# P3-TRACK-CTA — 主要トラックハブ（少なくとも U–Y と platformer）に「この型で自分の1本へ」→ graduation への CTA を付ける

Status: PASS

## Required evidence

- [x] Article
- [x] Starter
- [x] Tests
- [x] Reference game
- [x] Japanese
- [x] English
- [x] Mobile

## Changes

- Files: `scripts/inject-graduation-ctas.mjs`, `scripts/build.sh`, and twelve generated track hubs.
- Behavior: idempotently inserts a localized MAKE CTA before the main closing tag of Rhythm, Raycaster, Tower Defense, Reversi, Top-down Adventure (U–Y), and Platformer hubs. Each CTA resolves to the language-matching graduation hub.

## Commands and results

```text
node scripts/inject-graduation-ctas.mjs
Injected graduation CTAs into 6 track hubs in JA/EN.
rg 'graduation-cta|MAKE / YOUR OWN GAME' web/{ja,en}/tracks/{rhythm,raycaster,tower-defense,reversi,topdown-adventure,platformer}/index.html
36 marker matches: start/end marker plus label for all 12 hubs.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing architecture-promo layout preserves a readable full-width CTA. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | Existing promo link is one large touch target. |

- Japanese: six hubs link to `../../graduation/` with localized CTA copy.
- English: the matching six hubs link to the English graduation route.
- Readability / accessibility:
- Screenshots / recordings:
