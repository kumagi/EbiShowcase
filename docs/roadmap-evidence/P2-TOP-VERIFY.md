# P2-TOP-VERIFY — Top-down Adventure を著者基準で検証する

Status: PASS

## Required evidence

- [x] Edit target
- [x] RULE challenge
- [x] Desktop
- [x] Phone
- [x] Japanese
- [x] English
- [x] Tests

## Changes

- Files: generated Top-down Adventure pages and `scripts/gen-topdown-adventure-track.mjs`.
- Behavior: the visible edit target is a real path under `games/tracks/topdown-adventure`; all paths are paired with a real shared engine excerpt and a one-rule verification instruction.

## Commands and results

```text
node scripts/check-lessons.mjs
474 pages / 237 playable lessons passed the required article skeleton.
go test ./games/tracks/topdown-adventure/...
All 8 step packages compile successfully.
bash scripts/build.sh --fast
All generators, OGP assets, and the 208/208 gated lesson check passed.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | The standard responsive lesson frame and code panels preserve their source labels. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | The shared single-column breakpoint retains the game iframe and sequential code panels. |

- Japanese: checked `eight-way`: localized entry label, shared source panel, and rule challenge are present.
- English: checked `ebi-adventure`: matching English labels, source path, and rule challenge are present.
- Readability / accessibility:
- Screenshots / recordings:
