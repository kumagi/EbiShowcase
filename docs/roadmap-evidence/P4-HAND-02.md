# P4-HAND-02 — 追加で手書きトラックを1つ（推奨: survivors または match3）同様に更新し、手順書の再利用性を確かめる

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/inject-match3-authoring.mjs`, `scripts/build.sh`, and all 16 localized Match 3 step pages.
- Behavior: a second manually authored track now follows the same durable contract as Platformer: real entry path, one observable genre rule, and a package test command after generators run.

## Commands and results

```text
node scripts/inject-match3-authoring.mjs
Injected Match 3 authoring rules into 8 steps in JA/EN.
rg 'YOUR FIRST RULE.*games/tracks/match3' web/{ja,en}/tracks/match3/*/index.html
16 localized rule/path matches.
go test ./games/tracks/match3/...
All eight packages compile; ebi-match tests pass.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese: every Match 3 step names an entry path and board-rule challenge.
- English: every corresponding step has the same path and verification flow.
- Readability / accessibility:
- Screenshots / recordings:
