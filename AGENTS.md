# Ebi Showcase

A static web showcase of mini-games built with Go + Ebitengine, compiled to WebAssembly and played in the browser. Deployed to GitHub Pages as a purely static site. There is no backend, database, or API.

See `README.md` (Japanese) for the canonical local-run steps.

## Cursor Cloud specific instructions

### Services & how to run them
This product is a single static site. End-to-end testing means: build the WASM bundle, serve `dist/` over HTTP, and open it in a browser (WASM cannot load from `file://`).

- Build: `bash scripts/build.sh` — compiles every game (`game/`, `games/core/*`, `games/tracks/*`) to `dist/play/<id>/game.wasm` and assembles `dist/`. `dist/` is git-ignored and regenerated; rebuild after editing any Go game or `web/` asset.
- Serve: `python3 -m http.server 8080 --directory dist`, then open `http://localhost:8080` (or a specific game at `http://localhost:8080/play/<id>/`, e.g. `play/flappy/`).

### Non-obvious caveats
- Native `go build` / `go test ./...` FAILS for the game packages because Ebitengine's desktop backend needs X11/OpenGL/ALSA C headers that are not installed. This is expected — games are only ever compiled for `GOOS=js GOARCH=wasm` (no C deps). Do not try to run the games natively.
  - For Go unit tests use `go test ./examples/...` (the only native-testable package).
  - For linting game packages with `go vet`, prefix with the WASM target: `GOOS=js GOARCH=wasm go vet ./game/... ./games/...`.
- `scripts/build.sh` copies `wasm_exec.js` from `$(go env GOROOT)/lib/wasm/wasm_exec.js` (Go 1.24 path; older Go versions used `misc/wasm/`).
- Curriculum tooling is optional and needs only Node built-ins (no `package.json`/`npm install`): `node scripts/curriculum.mjs summary` and `bash scripts/ralph-loop.sh {next|status|verify}`.
