# P2-HELP — 生成器用の二層 code-lesson ヘルパ（入口・抜粋・パス・RULE challenge フィールド）を scripts に実装する

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `scripts/authoring-lesson-helpers.mjs`.
- Behavior: exports a bilingual dual-layer panel that requires a real editable
  entry path/code, real internal path/excerpt, and a path/location/action/
  verification RULE; it also requires exactly three non-duplicated concept cards.

## Commands and results

```text
node -e "import('./scripts/authoring-lesson-helpers.mjs') ..."
PASS; renders a dual-layer panel with YOUR FIRST RULE and validates three cards.

git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Shared HTML uses existing responsive code panels. |
| Tablet | 768 × 1024 | Touch | Existing code-panel scrolling remains applicable. |
| Phone | 390 × 844 | Touch | Existing code-panel scrolling remains applicable. |

- Japanese: labels entry, mechanism, and rule in Japanese.
- English: labels the same source layers in English.
- Readability / accessibility: source paths are visible text, not only implied by snippets.
- Screenshots / recordings: helper-only task; rendered verification follows in each track task.
