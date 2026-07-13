#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
node "$ROOT/scripts/gen-expansion-tracks.mjs"
node "$ROOT/scripts/gen-tactics-polish.mjs"
node "$ROOT/scripts/gen-active-rpg-polish.mjs"
node "$ROOT/scripts/gen-visual-novel-polish.mjs"
node "$ROOT/scripts/gen-racing-track.mjs"
node "$ROOT/scripts/gen-metroidvania-track.mjs"
node "$ROOT/scripts/gen-platformer-polish.mjs"
node "$ROOT/scripts/gen-survivors-polish.mjs"
node "$ROOT/scripts/gen-clicker-polish.mjs"
node "$ROOT/scripts/gen-rpg-polish.mjs"
node "$ROOT/scripts/gen-fighting-polish.mjs"
node "$ROOT/scripts/gen-merge-polish.mjs"
node "$ROOT/scripts/gen-deckbuilder-polish.mjs"
node "$ROOT/scripts/gen-slingshot-polish.mjs"
node "$ROOT/scripts/gen-falling-blocks-polish.mjs"
node "$ROOT/scripts/gen-match3-polish.mjs"
node "$ROOT/scripts/gen-sandbox-polish.mjs"
node "$ROOT/scripts/gen-bomb-maze-polish.mjs"
node "$ROOT/scripts/gen-monster-polish.mjs"
node "$ROOT/scripts/gen-falling-pairs-polish.mjs"
node "$ROOT/scripts/gen-maze-chase-polish.mjs"
node "$ROOT/scripts/embed-lesson-sources.mjs"
node "$ROOT/scripts/insert-feedback-form.mjs"
node "$ROOT/scripts/inject-beginner-bridges.mjs"
node "$ROOT/scripts/inject-ogp.mjs"
go run ./cmd/gen-og-images "$ROOT"
rm -rf "$ROOT/dist"
mkdir -p "$ROOT/dist"
cp -R "$ROOT/web/." "$ROOT/dist/"
WASM_EXEC="$(go env GOROOT)/lib/wasm/wasm_exec.js"

build_game() {
  local id="$1"
  local package="$2"
  local out="$ROOT/dist/play/$id"
  mkdir -p "$out"
  GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$out/game.wasm" "$package"
  cp "$WASM_EXEC" "$out/wasm_exec.js"
  cp "$ROOT/web/game.html" "$out/index.html"
}

build_game "flappy" "$ROOT/game"
while IFS= read -r main; do
  package="$(dirname "$main")"
  id="${package#"$ROOT/games/core/"}"
  build_game "$id" "$package"
done < <(find "$ROOT/games/core" -mindepth 2 -maxdepth 2 -name main.go -print | sort)

while IFS= read -r main; do
  package="$(dirname "$main")"
  id="$(basename "$package")"
  build_game "$id" "$package"
done < <(find "$ROOT/games/tracks" -mindepth 3 -maxdepth 3 -name main.go -print 2>/dev/null | sort)

# Keep legacy URLs working for the first published game.
cp "$ROOT/dist/play/flappy/game.wasm" "$ROOT/dist/game.wasm"
cp "$WASM_EXEC" "$ROOT/dist/wasm_exec.js"
echo "Built dist/ for GitHub Pages"
