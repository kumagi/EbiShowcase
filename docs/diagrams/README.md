# Diagram sources

The `*.mmd` files in this directory are the bilingual Mermaid flowchart
sources for the learning maps. Run:

```sh
npm run gen-diagrams
# (同じ処理を直接呼ぶ場合: node scripts/gen-diagrams.mjs)
```

The script writes Apache-2.0-compatible SVGs to
`web/assets/diagrams/` and inserts the corresponding figures into the static
pages. It renders the small flowchart subset used by this project without a
runtime dependency, so GitHub Actions can build the site offline. The `.mmd`
files remain valid Mermaid and can be opened with Mermaid Live or `mmdc` when
you want to experiment with a different layout.

Each generated figure has a bilingual `title`/`desc` for screen readers, while
the surrounding HTML supplies a short caption and a normal keyboard-focusable
link list where navigation is part of the diagram.

## What the new diagrams cover

The generator now covers the ideas that are easiest to lose in prose alone:

- Core games: input-to-canvas coordinates, timing windows, pipe spawning,
  reflection, bullet patterns, tile maps, enemy states, aiming vectors, snake
  slices, and Sokoban push rules.
- Visual effects: particle lifetimes, draw-layer order, alpha blending,
  lightning branches, easing, and colour-channel changes.
- Genre tracks: camera offsets, power-up states, waypoint patrols, data-driven
  stages, match-3 cascades, deck cycles, terrain/range grids, active battle
  queues, dialogue branches, merging, Reversi rays, bomb chains, lap timing,
  monster growth, survivor waves, falling blocks, RPG quests, and seeded worlds.

The route map in `scripts/gen-diagrams.mjs` is the source of truth for which
lesson receives each concept figure. A lesson may receive two or three related
figures; the device comparison remains first so the playable demo is still easy
to find.

SPDX-License-Identifier: Apache-2.0
