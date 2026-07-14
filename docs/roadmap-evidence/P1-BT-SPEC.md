# P1-BT-SPEC — Build Track 4ステップの日英コンテンツ表（slug・次に足す行・検証・公理の焦点・レベル01との関係）を scripts または docs に固定する

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `docs/BUILD_TRACK_SPEC.md`.
- Behavior: fixes the four bilingual Build Track slugs, local browser-download
  source contract, next-line additions, verification, axiom focus, and
  LEVEL 01 relationship before page/game implementation starts.

## Commands and results

```text
git diff --check
exit 0

node scripts/roadmap-ralph-loop.mjs check P1-BT-SPEC
PASS
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Not applicable: content specification only. |
| Tablet | 768 × 1024 | Touch | Not applicable: content specification only. |
| Phone | 390 × 844 | Touch | Not applicable: content specification only. |

- Japanese: each step has a Japanese learner promise and next-line wording.
- English: each step has the matching English promise and next-line wording.
- Readability / accessibility: the source and asset download contract explicitly
  works from a fresh Go workspace without a repository clone.
- Screenshots / recordings: not applicable until the four pages are implemented.
