# P2-TD-AUDIT — Tower Defense 全STEPの不一致表を作る

Status: PASS

## Required evidence

- [x] Mismatch inventory
- [x] Edit targets
- [x] Japanese
- [x] English

## Changes

- Files: `scripts/gen-tower-defense-track.mjs`, eight TD entry packages.
- Behavior: records generic formula/challenge pages versus real Config/data entries.

## Commands and results

```text
find games/tracks/tower-defense -maxdepth 2 -name main.go
PASS — eight entries.
git diff --check
exit 0
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | Audit only. |
| Tablet | 768 × 1024 | Touch | Audit only. |
| Phone | 390 × 844 | Touch | Audit only. |

- Japanese: generic tuning challenges have no file/function target.
- English: same missing authoring path.
- Readability / accessibility: repair must expose literal paths.
- Screenshots / recordings: audit only.

## Mismatch inventory

All eight pages (`path-patrol`, `range-circle`, `front-target`, `projectile-hit`,
`tower-economy`, `wave-data`, `tower-upgrades`, `ebi-defense`) use a generic
formula and generic challenge. Their actual editable entry is respectively
`games/tracks/tower-defense/<slug>/main.go`, mostly a shared-engine Config/data
wrapper. The generator must show that whole entry first, then the shared TD
mechanism; the RULE targets one path/wave/tower data row rather than tuning a
number. This applies identically to Japanese and English.
