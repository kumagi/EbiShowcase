#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
rm -rf "$ROOT/dist"
mkdir -p "$ROOT/dist"
cp -R "$ROOT/web/." "$ROOT/dist/"
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$ROOT/dist/game.wasm" "$ROOT/game"
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$ROOT/dist/wasm_exec.js"
echo "Built dist/ for GitHub Pages"
