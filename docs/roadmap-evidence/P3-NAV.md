# P3-NAV — progress / choose-your-path / ホームから「MAKE（書く）」モードが PLAY と並んで見えるようにする

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

- Files: `web/{ja,en}/index.html`, `web/{ja,en}/guides/{progress,choose-your-path}/index.html`.
- Behavior: home already has the Build Track MAKE promo; progress and choose-your-path now expose a matching MAKE / WRITE NEXT panel pointing to graduation projects.

## Commands and results

```text
rg 'MAKE / (BUILD TRACK|WRITE NEXT)' web/{ja,en}/index.html web/{ja,en}/guides/{progress,choose-your-path}/index.html
PLAY-adjacent MAKE entry points appear on home, progress, and choice pages in both languages.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Promo panels use the existing responsive architecture-promo layout. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | Panel link remains a single large tap target. |

- Japanese: home Build Track and two MAKE panels link readers to graduation.
- English: matching Build Track and MAKE panels are present.
- Readability / accessibility:
- Screenshots / recordings:
