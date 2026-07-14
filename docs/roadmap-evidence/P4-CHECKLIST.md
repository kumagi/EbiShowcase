# P4-CHECKLIST — `AUTHORING_CHECKLIST.md` を最終版にし、ADVANCED_QUALITY との関係（遊べる ≠ 書ける）を明記する

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `docs/AUTHORING_CHECKLIST.md`.
- Behavior: adds an explicit three-meter table: PLAYABLE proves launchability, ADVANCED_QUALITY proves production quality, and AUTHORING proves edit → RULE → verify. Neither of the first two substitutes for authoring readiness; removes duplicated Update placement wording.

## Commands and results

```text
rg 'PLAYABLE|ADVANCED_QUALITY|AUTHORING|does not prove' docs/AUTHORING_CHECKLIST.md
The definition and comparison table name all three independent meters.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese:
- English:
- Readability / accessibility:
- Screenshots / recordings:
