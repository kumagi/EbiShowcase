# P4-CORE-LATE — LEVEL 07–12 すべてに YOUR FIRST RULE（編集先パス付き）を追加する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Japanese
- [x] English
- [x] Automated checks

## Changes

- Files: `scripts/inject-core-authoring-links.mjs`, and generated `web/{ja,en}/games/{snake,space-shooter,sokoban,platformer,dungeon,bullet-hell}/index.html`.
- Behavior: LEVEL 07–12 now identify their real `games/core/<slug>/main.go` entry point, name one genre-specific first rule, and give a package-specific `go test` verification command.

## Commands and results

```text
node scripts/inject-core-authoring-links.mjs
Regenerated authoring panels for core levels 01–12.
rg 'YOUR FIRST RULE.*games/core/(snake|space-shooter|sokoban|platformer|dungeon|bullet-hell)/main.go' web/{ja,en}/games/{snake,space-shooter,sokoban,platformer,dungeon,bullet-hell}/index.html
12 localized rule/path matches.
go test ./games/core/{snake,space-shooter,sokoban,platformer,dungeon,bullet-hell}
All six packages compile; space-shooter tests pass.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese: six pages have source paths and concrete Japanese rule prompts.
- English: matching six pages have English prompts and the same source paths.
- Readability / accessibility:
- Screenshots / recordings:
