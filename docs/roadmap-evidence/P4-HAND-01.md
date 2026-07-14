# P4-HAND-01 — platformer トラック全STEPを二層コードまたは同等の編集先明示＋RULE challenge に更新する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/inject-platformer-authoring.mjs`, `scripts/build.sh`, and all 16 localized Platformer step pages.
- Behavior: every step names its real runnable `games/tracks/platformer/<slug>/main.go` entry, asks for one observable rule change, and gives its package test command. Injection runs after generators during the site build.

## Commands and results

```text
node scripts/inject-platformer-authoring.mjs
Injected platformer authoring rules into 8 steps in JA/EN.
rg 'YOUR FIRST RULE.*games/tracks/platformer' web/{ja,en}/tracks/platformer/*/index.html
16 localized rule/path matches.
go test ./games/tracks/platformer/...
All eight packages compile.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese: every platformer step has a Japanese edit path and verification action.
- English: every matching step has the same real path and English action.
- Readability / accessibility:
- Screenshots / recordings:
