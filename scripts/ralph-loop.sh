#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

case "${1:-next}" in
  next)
    node scripts/curriculum.mjs next
    ;;
  status)
    node scripts/curriculum.mjs summary
    ;;
  verify)
    bash scripts/build.sh
    node scripts/curriculum.mjs summary
    test -f dist/play/tap-target/game.wasm
    test -f dist/play/timing-meter/game.wasm
    test -f dist/play/catch-stars/game.wasm
    test -f dist/play/flappy/game.wasm
    echo "Ralph verification passed"
    ;;
  *)
    echo "usage: $0 {next|status|verify}" >&2
    exit 2
    ;;
esac
