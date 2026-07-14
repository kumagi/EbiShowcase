# Build Track content specification

This is the source-of-truth content table for the four ungated Build Track
steps. It sits between the setup guide's empty local Go window and the finished
LEVEL 01 game. A reader starts with a fresh local Go workspace; cloning this
repository is never a prerequisite.

## Reader setup contract

Each step page gives a browser-download link for its exact `main.go` source and
any needed asset. The reader creates their own folder, saves the download as
`main.go`, and runs the commands shown on the page:

```sh
mkdir ebi-build-step
cd ebi-build-step
go mod init example.com/ebi-build-step
go get github.com/hajimehoshi/ebiten/v2
go run .
```

The files inside this repository remain named in every page as a public-source
map, but the lesson must work after downloading just the listed file(s). All
four initial steps use Ebitengine's shape drawing and need no binary asset.

## URL and source contract

| Step | Slug / bilingual URL | Local source supplied by browser download | Next pager | LEVEL 01 relationship |
| --- | --- | --- | --- | --- |
| 01 | `build-track/empty-loop` · `/ja/build/empty-loop/` / `/en/build/empty-loop/` | `games/build-track/empty-loop/main.go` | STEP 02 `state-picture` | Shows the empty loop that LEVEL 01 fills with a target and score. |
| 02 | `build-track/state-picture` · `/ja/build/state-picture/` / `/en/build/state-picture/` | `games/build-track/state-picture/main.go` | STEP 03 `tap-score` | Makes LEVEL 01's target position a visible `game` field before it moves. |
| 03 | `build-track/tap-score` · `/ja/build/tap-score/` / `/en/build/tap-score/` | `games/build-track/tap-score/main.go` | STEP 04 `hit-reset` | Isolates LEVEL 01's input and score rule before collision. |
| 04 | `build-track/hit-reset` · `/ja/build/hit-reset/` / `/en/build/hit-reset/` | `games/build-track/hit-reset/main.go` | LEVEL 01 `tap-target` | Is the small, hand-written predecessor of the finished LEVEL 01 target game. |

## Bilingual step content table

| Step | Japanese learner promise / next lines | English learner promise / next lines | Verification | Axiom focus |
| --- | --- | --- | --- | --- |
| 01 | **空の窓を動かす。** `type game struct{}`、`func (g *game) Update() error { return nil }`、`func (g *game) Draw(s *ebiten.Image) { s.Fill(...) }`、`ebiten.RunGame(&game{})` を順に足す。 | **Run an empty window.** Add `type game struct{}`, `Update`, `Draw` with `Fill`, then `ebiten.RunGame(&game{})`. | `go run .` opens one solid-colour window; close it normally. | `Update` changes state; `Draw` paints. Neither has a rule yet. |
| 02 | **数字を絵に写す。** `game` に `x, y int` を足し、`Draw` の `DrawRect` が `g.x, g.y` を読む行を足す。反例として `Draw` の `g.x++` は禁止枠に示す。 | **Project numbers into a picture.** Add `x, y int` to `game`, then make `DrawRect` read `g.x, g.y`. Show `g.x++` in `Draw` as a prohibited counterexample. | Change initial `x` once, run again, and see only the square's starting position change. | The position's truth lives in `game`; `Draw` only projects it. |
| 03 | **押したら得点が増える。** `Update` のタップ判定の中へ `g.score++` を足す。Draw は `g.score` を文字にするだけ。 | **Tap to add a score.** Add `g.score++` inside the tap branch in `Update`. Draw only turns `g.score` into text. | Tap/click the square three times and observe score 3; a focused `go test` checks the pure increment helper. | Input and score mutation belong in `Update`. |
| 04 | **当たったら標的を移す。** `Update` で `math.Hypot` を使う当たり判定の行と、命中時に `g.targetX, g.targetY = ...` を足す。 | **Move the target after a hit.** In `Update`, add a `math.Hypot` hit-test line and, on success, `g.targetX, g.targetY = ...`. | Click/tap the circle: score rises and it appears at its next deterministic position; `go test ./...` proves hit/miss behavior. | A RULE belongs in `Update`; equal `game` state produces an equal frame. |

## Required page anatomy

Every bilingual page must include, in this order:

1. a playable embedded WASM demo (ungated; it does not change the 208 count);
2. the complete currently-working `main.go` for that step, with a direct
   browser download link;
3. a visually distinct **NEXT LINES / 次に足す行** panel containing only the
   small addition for the next action, not the full later answer;
4. one short **RULE / ルール** challenge that says file, function, behavior,
   and verification;
5. one short reminder of the relevant Update/Draw axiom;
6. previous/next pager and a link to LEVEL 01, with LEVEL 01 linking back to
   this route.

The fourth step's `math.Hypot` line is the learner-authored collision decision.
It is in `Update` (or a pure helper Update calls), never in `Draw`.
