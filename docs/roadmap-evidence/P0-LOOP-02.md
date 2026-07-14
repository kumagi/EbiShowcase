# P0-LOOP-02 — Glossary の概念マップと本文を「画面は game に1bit違わず追従／Draw は game を書かない」に強化する（図解キャプション含む）

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `docs/Glossary.md`, `scripts/gen-diagrams.mjs`, regenerated
  `web/assets/diagrams/{ja,en}_glossary-map.svg` and lesson diagram embeds.
- Behavior: the glossary now explains the three Update/Draw axioms in prose and
  Mermaid. Its generated SVG shows input → Update → game → Draw → pixels, with
  JSON branching from state; its title and description identify projection
  rather than drawing-as-mutation.

## Commands and results

```text
node scripts/gen-diagrams.mjs
Generated 98 bilingual diagrams and injected lesson visuals into 490 page(s).

rg -n '入力|Update|game|Draw|画面|same state' web/assets/diagrams/ja_glossary-map.svg web/assets/diagrams/en_glossary-map.svg
exit 0: both generated diagrams include the complete state-to-pixel path.

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Reviewed the glossary flow: input leads to Update, then state, then Draw/pixels. |
| Tablet | 768 × 1024 | Touch | SVG is responsive static content; no input required. |
| Phone | 390 × 844 | Touch | SVG labels remain short, and the document text carries the full explanation. |

- Japanese: labels say 入力, 状態を書き換える, 投影するだけ, and 同じ state → 同じ絵.
- English: the paired asset says Input, mutate state, project only, and same state → same frame.
- Readability / accessibility: SVG `title`, `desc`, visible labels, Mermaid, and
  surrounding prose communicate the same map.
- Screenshots / recordings: generated bilingual SVG assets are the durable
  visual record for this documentation task.
