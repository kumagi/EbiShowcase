# P1-CORE-01 — LEVEL 01 tap-target の主 challenge を YOUR FIRST RULE（Update 側）にし、編集ファイルと挿入位置を本文に書き、公理の3条を本文で再掲する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `web/ja/games/tap-target/index.html`, `web/en/games/tap-target/index.html`.
- Behavior: replaces tuning-first copy with a named Update-side combo rule,
  exact insertion point, reset condition, local verification, and the three
  Update/Draw axioms in the article body.

## Commands and results

```text
go test ./games/core/tap-target
PASS (no test files)

node scripts/inject-ogp.mjs
PASS

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Pointer | Existing target game remains unchanged and playable. |
| Tablet | 768 × 1024 | Touch | Existing touch path remains unchanged. |
| Phone | 390 × 844 | Touch | Existing responsive lesson shell remains unchanged. |

- Japanese: exact path, fields, Update hit branch, reset, and verification are named.
- English: has the equivalent rule and three promises.
- Readability / accessibility: the three promises are a text block, not an
  image or colour-only instruction.
- Screenshots / recordings: no runtime behaviour changed in this documentation task.
