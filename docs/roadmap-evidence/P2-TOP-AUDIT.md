# P2-TOP-AUDIT — Top-down Adventure 全STEPの不一致表を作る

Status: PASS

## Required evidence

- [x] Mismatch inventory
- [x] Edit targets
- [x] Japanese
- [x] English

## Changes

- Files: `scripts/gen-topdown-adventure-track.mjs`; generated `web/{ja,en}/tracks/topdown-adventure/*/index.html`.
- Behavior: all eight step pages used a display-only excerpt plus generic cards. The runnable entry is a thin `Run(lesson)` wrapper, but the lesson did not label it or show the shared `Update` implementation.

| Step | Previous mismatch | Edit target |
| --- | --- | --- |
| 01 eight-way | Normalization excerpt was not tied to its actual entry/shared loop. | Entry + movement/update layer |
| 02 sword-reach | Attack-box rule lacked a traceable Update location. | Entry + shared Update layer |
| 03 hurt-recovery | Recovery timer appeared as isolated pseudo-code. | Entry + shared Update layer |
| 04 room-clear | Room phase explanation did not identify data versus loop. | Entry + room/update layer |
| 05 key-treasure | Key rule lacked a concrete source route. | Entry + shared Update layer |
| 06 tool-puzzles | Tool rule was shown without the data/update/draw split. | Entry + shared Update layer |
| 07 guardian-phases | Phase code was not presented as part of the runnable game. | Entry + shared Update layer |
| 08 ebi-adventure | Stage switch implied a separate implementation for each room. | Entry + shared Update layer |

## Commands and results

```text
find games/tracks/topdown-adventure -maxdepth 2 -name main.go
All 8 entry points call internal/topdownadventuregame.Run with their lesson number.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese: eight pages have a localized title, explanation, and editing challenge.
- English: the matching eight pages have independent English copy.
- Readability / accessibility:
- Screenshots / recordings:
