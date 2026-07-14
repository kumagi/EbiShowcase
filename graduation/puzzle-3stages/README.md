# Three-stage puzzle — graduation starter

Three puzzle records share one completion rule. Start in a fresh Go workspace:
create a folder, run `go mod init example.com/three-puzzles` and `go get github.com/hajimehoshi/ebiten/v2`, then save the `starter/` files from their GitHub Raw browser view. No clone is required.

`go test` begins red. TODO 1 creates progress, TODO 3 uses each `StageData.Target` to advance, and TODO 2 resets. Keep data and progression in `progress.go`; reserve `main.go` for input and drawing. Compare `reference/` only after green tests.
