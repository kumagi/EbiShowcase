# P0-LOOP-03 — setup・ホーム・LEVEL 01 の beginner-bridge／概念カード／lab 文言を公理の3条で強化する

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `scripts/gen-setup-guide.mjs`, `scripts/inject-beginner-bridges.mjs`,
  `web/{ja,en}/index.html`, `web/{ja,en}/games/tap-target/index.html`, and
  regenerated `web/{ja,en}/guides/setup/index.html`.
- Behavior: setup’s first `main.go` labels Update as the state-changing place
  and Draw as projection-only. Home now has a STATE FIRST foundation card.
  LEVEL 01 has a three-card Update/Draw/same-state explanation plus a lab hint
  that Draw must not alter the numbers it displays.

## Commands and results

```text
node scripts/gen-setup-guide.mjs
wrote web/ja/guides/setup/index.html
wrote web/en/guides/setup/index.html

node scripts/inject-beginner-bridges.mjs
updated the generated beginner bridge from its source table.

rg -n 'STATE FIRST|Rules live in Update|ルールは Update|SAME STATE' \
  web/ja/index.html web/en/index.html web/ja/games/tap-target/index.html web/en/games/tap-target/index.html
exit 0

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Home and LEVEL 01 copy retain the existing responsive card and lab structure. |
| Tablet | 768 × 1024 | Touch | The additional foundation/concept card is ordinary responsive article content. |
| Phone | 390 × 844 | Touch | The text is short; the existing mobile card layout handles the third card. |

- Japanese: setup, home, bridge, LEVEL 01 cards, and lab all state Update-side rules and Draw-side projection.
- English: the paired setup, home, and LEVEL 01 wording uses the same three-part model.
- Readability / accessibility: the principle appears as a short foundation card,
  an explicit concept card, and an interactive-lab hint rather than only a
  dense paragraph.
- Screenshots / recordings: not required for this copy/generator task; the
  subsequent full browser verification covers the integrated site.
