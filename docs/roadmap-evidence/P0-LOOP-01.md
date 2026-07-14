# P0-LOOP-01 — AGENTS・AUTHORING_CHECKLIST・README 学習道筋に Update/Draw 公理の3条を正式追記し、RULE は Update 側・Draw は投影のみと明記する

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `AGENTS.md`, `docs/AUTHORING_CHECKLIST.md`, `README.md`.
- Behavior: contributor guidance, the authoring completion checklist, and the
  public roadmap introduction now agree on three axioms: Update owns input and
  state changes; Draw only projects state; identical state produces identical
  pixels. RULE challenges therefore belong in Update or its pure helpers.

## Commands and results

```text
rg -n 'Update.*(input|入力)|Draw.*(project|写)|same.*(state|game)|同じ状態' AGENTS.md docs/AUTHORING_CHECKLIST.md README.md
exit 0: each document contains the game-loop contract.

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Not applicable: contributor and README copy only. |
| Tablet | 768 × 1024 | Touch | Not applicable: contributor and README copy only. |
| Phone | 390 × 844 | Touch | Not applicable: contributor and README copy only. |

- Japanese: README states all three axioms in the learner path and ties RULE to Update.
- English: AGENTS and the checklist state the equivalent contributor contract.
- Readability / accessibility: the three rules are short numbered statements,
  followed by the concrete consequence for a learner challenge.
- Screenshots / recordings: not applicable; public lesson-page reinforcement is
  handled by the dedicated P0-LOOP-03 task.
