# Three-room exploration — graduation starter

Make a three-room adventure: collect a key in room 1, walk through room 2,
then open the room-3 exit. This is a red-test starter. Do not open `reference/`
until your tests pass.

Start in an empty Go workspace; do not clone this repository. Create a folder,
run `go mod init example.com/three-rooms` and `go get github.com/hajimehoshi/ebiten/v2`,
then save the three `starter/` files from their GitHub **Raw** browser view.

Run `go test`. Each failure names its TODO: TODO 3 collects the key, TODO 4
moves among bounded rooms, TODO 5 opens the gated exit, and TODO 2 restarts.
Keep those rules in `state.go`; `main.go` only maps input and draws state.
