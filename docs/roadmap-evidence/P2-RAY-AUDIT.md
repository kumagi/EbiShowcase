# P2-RAY-AUDIT — Raycaster 全STEPの不一致表を作る

Status: PASS

## Required evidence

- [x] Mismatch inventory
- [x] Edit targets
- [x] Japanese
- [x] English

## Changes

- Files: `scripts/gen-raycaster-track.mjs`, `games/tracks/raycaster/*/main.go`.
- Behavior: records six thin entries and generic/illustrative article snippets.

## Commands and results

```text
find games/tracks/raycaster -maxdepth 2 -name main.go
PASS — six entries inspected.
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Audit only. |
| Tablet | 768 × 1024 | Touch | Audit only. |
| Phone | 390 × 844 | Touch | Audit only. |

- Japanese: generic CHANGE/CHALLENGE has no entry location.
- English: same missing authoring route.
- Readability / accessibility: the repair will expose literal paths.
- Screenshots / recordings: audit only.

## Mismatch inventory

| STEP | Entry source | Existing mismatch | Repair |
| --- | --- | --- | --- |
| facing-move | `games/tracks/raycaster/facing-move/main.go` | angle formula is not the entry chart/config | show entry plus `internal/raycaster*` movement excerpt; add a map/config rule. |
| single-ray | `.../single-ray/main.go` | DDA pseudocode lacks editable target | show entry plus cast excerpt; change one wall cell. |
| distance-strip | `.../distance-strip/main.go` | projection formula lacks source path | show entry plus strip mapper; add one wall/material data rule. |
| column-view | `.../column-view/main.go` | column loop is illustrative only | show entry plus column renderer; add one map opening. |
| textured-view | `.../textured-view/main.go` | correction formula lacks caller | show entry plus texture mapper; add material data. |
| ebi-raycaster | `.../ebi-raycaster/main.go` | Mission sample does not identify configuration entry | show entry plus mission runner; add enemy/key row. |
