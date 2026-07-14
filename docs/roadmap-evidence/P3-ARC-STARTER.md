# P3-ARC-STARTER — `graduation/arcade-60/starter` を穴あきの60秒ゲーム骨格＋失敗するテスト群へ拡充する

Status: PASS

## Required evidence

- [x] Holey starter
- [x] Failing tests
- [x] TODO mapping
- [x] Automated checks

## Changes

- Files: `graduation/arcade-60/starter/{main.go,rules.go,rules_test.go}`, `graduation/arcade-60/README.md`.
- Behavior: the starter runs as an Ebitengine window while `Round.Step` remains pure and unfinished. Keyboard/pointer input enters at `main.go`; scoring, timer, finish, and restart are TODOs in `rules.go`; `Draw` has a deliberate result-UI TODO.

## Commands and results

```text
go test ./graduation/arcade-60/starter
Expected red: TestActionAddsOneStar, TestTimeLimitEndsRound, TestRestartMakesFreshRound fail.
go build ./graduation/arcade-60/starter
Passes: the starter is runnable while its rule tests are deliberately red.
git diff --check
Passes.
```

## Manual checks

| Surface | Representative viewport | Input completed | Result / issue |
| --- | --- | --- | --- |
| Desktop | 1440 × 900 | Keyboard + pointer | |
| Tablet | 768 × 1024 | Touch | |
| Phone | 390 × 844 | Touch | |

- Japanese: README is intentionally concise English source material; the bilingual step-by-step article is the next task.
- English: README starts from a fresh Go workspace, explicitly avoiding a repository clone, and maps every red test to one TODO.
- Readability / accessibility:
- Screenshots / recordings:
