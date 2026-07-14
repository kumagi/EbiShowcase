# P2-HAND-PATTERN — 手書き良例（platformer 系）から二層＋RULE の適用パッチ手順を docs に短い手順書として残す

Status: PASS

## Required evidence

- [x] Implementation
- [x] Automated checks
- [x] Manual review

## Changes

- Files: `docs/HAND_AUTHORED_AUTHORING_PATTERN.md`.
- Behavior: documents the repeatable patch sequence for manually maintained lessons: real entry path, labelled internal mechanism, DATA/UPDATE/DRAW cards, one concrete rule, bilingual parity, and verification.

## Commands and results

```text
go test ./games/tracks/platformer/...
Platformer step packages compile.
node scripts/check-lessons.mjs
Checked 474 pages (237 playable lessons). OK.
git diff --check
No whitespace errors.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | The documented order matches the existing lesson reading flow: play, concepts, code, challenge. |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | The guide requires sequential panels, so no side-by-side code comparison is needed on a narrow screen. |

- Japanese: the guide gives Japanese authoring instructions and preserves literal UI labels and commands.
- English: the prescribed page labels include their English counterparts, so translators retain the same source contract.
- Readability / accessibility:
- Screenshots / recordings:
