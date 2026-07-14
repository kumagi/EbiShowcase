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
    node scripts/check-lessons.mjs
    node scripts/check-site-metadata.mjs
    echo "Ralph verification passed"
    ;;
  lessons)
    node scripts/check-lessons.mjs "${@:2}"
    ;;
  roadmap)
    node scripts/roadmap-ralph-loop.mjs "${@:2}"
    ;;
  *)
    echo "usage: $0 {next|status|verify|lessons|roadmap ...}" >&2
    exit 2
    ;;
esac
