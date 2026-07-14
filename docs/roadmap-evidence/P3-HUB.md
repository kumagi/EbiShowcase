# P3-HUB — `web/{ja,en}/graduation/` を前提・クローン手順・4コマンド・3 brief の本文化し、カード置き場で終わらせない

Status: PASS

## Required evidence

- [x] Article
- [x] Starter
- [x] Tests
- [x] Reference game
- [x] Japanese
- [x] English
- [x] Mobile

## Changes

- Files: `web/{ja,en}/graduation/index.html`.
- Behavior: replaces the former card shelf with a no-clone fresh-workspace article, four setup commands, red→green→play workflow, and links to all three detailed briefs. The roadmap's old clone wording is superseded by the explicit project decision to avoid clone assumptions.

## Commands and results

```text
rg 'arcade-60|exploration-3rooms|puzzle-3stages|go mod init|go get|go test|go run' web/{ja,en}/graduation/index.html
Both hubs contain all three briefs and the four-command authoring route.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Setup article, three linked briefs, and workflow cards maintain reading order. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | Brief links use the existing vertically responsive path list. |

- Japanese: all instructions and three briefs are localized.
- English: matching no-clone setup and brief links are localized.
- Readability / accessibility:
- Screenshots / recordings:
