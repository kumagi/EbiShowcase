# P4-COPY-REGRESS — 「値を変えて」定型文の回帰検索を零件にし、OGP 再注入後も残っていないことを確認する

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `scripts/check-authoring-copy-regression.mjs`, `scripts/build.sh`.
- Behavior: after OGP injection, scans every HTML OGP/Twitter description for the retired Japanese promise that authoring only means “change a value.” The scan deliberately permits pedagogical body-text mentions of tuning and fails only public metadata regression.

## Commands and results

```text
node scripts/inject-ogp.mjs
OGP reinjected for 595 HTML pages.
node scripts/check-authoring-copy-regression.mjs
Authoring-copy regression check passed for 596 HTML pages.
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
