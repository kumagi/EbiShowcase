# 60-second arcade — graduation starter

Make a tiny game with one action, a score, a 60-second limit, and a result
screen. This is a **red-test starter**, not a finished game. Do not read
`reference/` until your tests are green.

## Start without cloning this repository

Create an empty Go workspace, initialize it, then download the three files in
`starter/` from this repository's GitHub file view (use **Raw** → save as).

```sh
mkdir arcade-60 && cd arcade-60
go mod init example.com/arcade-60
go get github.com/hajimehoshi/ebiten/v2
# save starter/main.go, starter/rules.go, and starter/rules_test.go here
go test
go run .
```

The first `go test` is expected to fail. Open `rules.go` and complete the TODO
named in each failing test:

| Test | TODO | Rule to implement |
| --- | --- | --- |
| `TestActionAddsOneStar` | TODO 4 | one action gives exactly 10 points |
| `TestTimeLimitEndsRound` | TODO 5 | frame 3600 ends the round |
| `TestRestartMakesFreshRound` | TODO 3 | restart creates a clean round |

Then add your own rule, such as a combo multiplier or a missed-action penalty,
and write its test before changing `Round.Step`. `main.go` is only the input
and display edge; the rules should stay testable in `rules.go`.

When all tests pass, open `reference/main.go` only to compare structure—not to
copy it. Apache-2.0.
