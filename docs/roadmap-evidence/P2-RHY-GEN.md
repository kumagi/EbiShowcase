# P2-RHY-GEN — `gen-rhythm-track.mjs` を二層パネル・一意 concept-row・RULE challenge 対応に改修し再生成する

Status: PASS

## Required evidence

- [x] Dual panel
- [x] Unique concept-row
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/gen-rhythm-track.mjs`, `scripts/authoring-lesson-helpers.mjs`,
  and regenerated `web/{ja,en}/tracks/rhythm/` pages.
- Behavior: each rhythm page now shows its actual entry `main.go`, a labelled
  internal mechanism excerpt, a data/update/draw concept row, and a STEP-specific RULE.

## Commands and results

```text
node scripts/gen-rhythm-track.mjs
PASS
go test ./games/tracks/rhythm/...
PASS
rg "EDIT THIS ENTRY|編集する入口|YOUR FIRST RULE" generated pages
PASS
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Existing rhythm demos unchanged. |
| Tablet | 768 × 1024 | Touch | Existing touch lanes unchanged. |
| Phone | 390 × 844 | Touch | Existing iframe shell unchanged. |

- Japanese: labels editable source and rule paths in Japanese.
- English: labels the same entry/mechanism/rule sequence in English.
- Readability / accessibility: source locations are explicit text and code panels scroll.
- Screenshots / recordings: generated article change; gameplay unchanged.
