# P4-NAV-02 — Interest entrances.

Status: PASS

## Implementation

- [x] Implementation
- Added matching Japanese/English action, puzzle, strategy, story, and presentation entrances with a first game and learning chain.

## Automated checks

- [x] Automated checks

```text
rg -n "ACTION|PUZZLE|STRATEGY|STORY|PRESENTATION" web/ja/guides/choose-your-path/index.html web/en/guides/choose-your-path/index.html
git diff --check
success
```

## Manual review

- [x] Manual review
- Each route points to a playable first lesson instead of a generic category page.
