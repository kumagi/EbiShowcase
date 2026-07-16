#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
go run ./cmd/gen-favicon "$ROOT"
FAST_BUILD=0
if [[ "${1:-}" == "--fast" ]]; then
  FAST_BUILD=1
  shift
fi
if [[ "$#" -ne 0 ]]; then
  echo "usage: $0 [--fast]" >&2
  exit 2
fi
# A deployment must always refresh OGP metadata and images. --fast is only a
# local convenience; CI deliberately ignores it.
if [[ -n "${CI:-}" ]]; then
  FAST_BUILD=0
fi
node "$ROOT/scripts/gen-setup-guide.mjs"
node "$ROOT/scripts/gen-testing-guide.mjs"
node "$ROOT/scripts/gen-visual-effects.mjs"
node "$ROOT/scripts/gen-expansion-tracks.mjs"
node "$ROOT/scripts/gen-tactics-polish.mjs"
node "$ROOT/scripts/gen-reversi-track.mjs"
node "$ROOT/scripts/gen-raycaster-track.mjs"
node "$ROOT/scripts/gen-rhythm-track.mjs"
node "$ROOT/scripts/gen-tower-defense-track.mjs"
node "$ROOT/scripts/gen-topdown-adventure-track.mjs"
node "$ROOT/scripts/gen-shader-lab.mjs"
node "$ROOT/scripts/gen-audio-lab.mjs"
node "$ROOT/scripts/gen-camera-lab.mjs"
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
node "$ROOT/scripts/gen-new-genre-cards.mjs"
node "$ROOT/scripts/gen-legacy-aliases.mjs"
node "$ROOT/scripts/embed-lesson-sources.mjs"
node "$ROOT/scripts/insert-feedback-form.mjs"
node "$ROOT/scripts/inject-beginner-bridges.mjs"
node "$ROOT/scripts/inject-core-authoring-links.mjs"
node "$ROOT/scripts/inject-graduation-ctas.mjs"
node "$ROOT/scripts/inject-platformer-authoring.mjs"
node "$ROOT/scripts/inject-match3-authoring.mjs"
node "$ROOT/scripts/inject-update-draw-contract.mjs"
node "$ROOT/scripts/inject-renderer-freedom.mjs"
node "$ROOT/scripts/inject-showcase-finals.mjs"
node "$ROOT/scripts/inject-capstone-renderer-extras.mjs"
node "$ROOT/scripts/normalize-loop-language.mjs"
node "$ROOT/scripts/normalize-tick-language.mjs"
node "$ROOT/scripts/inject-feedback-code-examples.mjs"
node "$ROOT/scripts/gen-feedback-teaching-notes.mjs"
node "$ROOT/scripts/gen-diagrams.mjs"
node "$ROOT/scripts/home-thumbnails.mjs" inject
node "$ROOT/scripts/inject-ogp.mjs"
node "$ROOT/scripts/normalize-html-whitespace.mjs"
node "$ROOT/scripts/normalize-html-whitespace.mjs" --check
node "$ROOT/scripts/check-authoring-copy-regression.mjs"
node "$ROOT/scripts/check-tick-language.mjs"
OGP_STATE="$ROOT/.cache/ebi-showcase/ogp-inputs.sha256"
OGP_FINGERPRINT="$(node "$ROOT/scripts/ogp-cache.mjs" fingerprint)"
if [[ "$FAST_BUILD" == "1" ]] && [[ -f "$OGP_STATE" ]] && [[ "$(<"$OGP_STATE")" == "$OGP_FINGERPRINT" ]] && node "$ROOT/scripts/ogp-cache.mjs" verify; then
  echo "OGP images unchanged — skipped PNG regeneration (--fast)."
else
  go run ./cmd/gen-og-images "$ROOT"
  mkdir -p "$(dirname "$OGP_STATE")"
  node "$ROOT/scripts/ogp-cache.mjs" fingerprint > "$OGP_STATE"
  node "$ROOT/scripts/ogp-cache.mjs" verify
fi
node "$ROOT/scripts/check-site-metadata.mjs"
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

# Build Track lessons are intentionally outside the 208 playable curriculum
# gate, but their browser demos are built alongside the published games.
while IFS= read -r main; do
  package="$(dirname "$main")"
  id="$(basename "$package")"
  build_game "$id" "$package"
done < <(find "$ROOT/games/build-track" -mindepth 2 -maxdepth 2 -name main.go -print 2>/dev/null | sort)

while IFS= read -r main; do
  package="$(dirname "$main")"
  id="$(basename "$package")"
  build_game "$id" "$package"
done < <(find "$ROOT/games/tracks" -mindepth 3 -maxdepth 3 -name main.go -print 2>/dev/null | sort)

# Build a synchronized renderer gallery for every genre capstone. The overlay
# replaces only ebiten.RunGame(...); original Update and Layout source stays
# byte-for-byte unchanged.
RENDER_OVERLAY="$(node "$ROOT/scripts/render-freedom-overlays.mjs" prepare)"
while IFS=$'\t' read -r id package; do
  out="$ROOT/dist/play/${id}-renderer"
  mkdir -p "$out"
  GOOS=js GOARCH=wasm go build -overlay="$RENDER_OVERLAY" -trimpath -ldflags="-s -w" -o "$out/game.wasm" "$package"
  cp "$WASM_EXEC" "$out/wasm_exec.js"
  cp "$ROOT/web/game.html" "$out/index.html"
done < <(node "$ROOT/scripts/render-freedom-overlays.mjs" capstones)

# Keep legacy URLs working for the first published game.
cp "$ROOT/dist/play/flappy/game.wasm" "$ROOT/dist/game.wasm"
cp "$WASM_EXEC" "$ROOT/dist/wasm_exec.js"
echo "Built dist/ for GitHub Pages"
