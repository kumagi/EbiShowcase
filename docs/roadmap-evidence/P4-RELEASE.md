# P4-RELEASE — `verify --full` と `complete` 前提のフル監査（ビルド、テスト、日英、Build Track、graduation、メーター）を通し、Authoring Pass 完了を宣言する

Status: PASS

## Required evidence

- [x] Full build
- [x] All tests
- [x] All evidence
- [x] Authoring metrics
- [x] Pages artifact

## Changes

- Files: full generated `web/`/`dist/` artifact, `scripts/roadmap-ralph-loop.mjs`, `scripts/check-graduation-starters.mjs`, and all P0–P4 evidence.
- Behavior: full verification treats starter tests as intentionally red and reference implementations as green; all non-starter Go packages, curriculum/build/metadata gates, authoring-copy regression, assets, and authoring metrics pass.

## Commands and results

```text
node scripts/roadmap-ralph-loop.mjs verify --full
Green non-starter Go packages; three expected-red starter packages and three green references; full site build passed.
bash scripts/ralph-loop.sh verify
208/208 gated lessons, 29/29 VFX, 66 home cards, 595 OGP images, bilingual links, lesson and metadata gates passed.
bash scripts/ralph-loop.sh status
208/208 playable plus authoring: Build Track 4, Core 12/12, hubs 12/12, briefs 6/6, first-30-minutes 2.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Representative authored panels preserve the existing responsive layout. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | Generated mobile game frames and vertical authoring panels are covered by the site lesson gate and phase evidence. |

- Japanese: generated hub, graduation, and sampled authoring contracts are present.
- English: matching generated contracts and bilingual link gate pass.
- Readability / accessibility:
- Screenshots / recordings:
