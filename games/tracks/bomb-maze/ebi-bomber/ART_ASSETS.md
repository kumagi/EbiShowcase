# Ebi Bomber art assets

All raster assets in `assets/` are original project artwork generated for this
capstone and distributed under the repository's Apache License 2.0. They do not
copy characters, logos, layouts, or locations from an existing game.

## Production notes

- `bomber-coral-forge.png` — portrait premium mobile-game environment: an
  original submerged coral-forge labyrinth chamber, pearlstone central arena,
  cyan water channels, brass machinery, warm furnaces and an open touch-control
  strip; no characters, UI, text, logos or copied locations.
- `bomber-characters-atlas.png` — two isolated original characters: Ebi
  Tenjiroh as a coral-orange shrimp demolition explorer, and a midnight-blue
  bioluminescent reef-lizard scout; polished three-quarter top-down mobile-game
  rendering with silhouettes readable at gameplay size.
- `bomber-bomb-flame-atlas.png` — an ornate pearl-and-obsidian clockwork bomb
  and a four-direction coral fire burst with a pearl core, orange flame and cyan
  edge light.
- `bomber-walls-atlas.png` — an indestructible pearlstone/brass forge pillar and
  a visibly breakable ivory-shell barricade with coral braces and orange cracks.
- `bomber-items-atlas.png` — three isolated upgrades: a blast-range flame pearl,
  a two-bomb capacity satchel, and aqua current fins for movement speed.

The four atlases were generated on uniform chroma backgrounds, converted to
transparent PNGs with the repository image-generation helper, then resized for
WASM delivery. Prompts explicitly required contemporary Japanese mobile-game
quality, original designs, crisp small-scale readability, no text, no logos, no
watermarks, no retro pixel treatment and no copyrighted resemblance.

Generated source images are retained in the local Codex generation store; the
optimized PNGs embedded by `art.go` are the canonical repository assets.

## Prompt record

The production prompts requested these exact subjects and constraints:

- Background: “an original submerged coral-forge labyrinth chamber built
  around a square pearlstone arena”; dark ocean-temple machinery, cyan water
  channels, brass pipes, furnaces, coral vents and shell arches; polished
  contemporary Japanese mobile-game environment illustration; high
  three-quarter top-down camera; a centered unobstructed board and open lower
  touch-control strip; no characters, bombs, walls, items, UI, text, logos,
  watermark, retro pixels or copyrighted resemblance.
- Characters: “exactly two original compact full-body characters in equal
  cells”: Ebi Tenjiroh with teal eyes, pearl armor, cream scarf, brass utility
  belt and fuse lighter; a sleek abyssal reef lizard with cyan bioluminescent
  stripes, coral-red fins and gold eyes; uniform `#ff00ff` chroma background;
  no extra figures, casts shadows, text, logos, watermark or copied designs.
- Bomb/effect: “exactly two isolated objects in equal cells”: a spherical
  pearl-and-obsidian clockwork bomb with brass bands, orange fuse and cyan runes;
  a cross-shaped coral fire burst with four orange arms and a restrained cyan
  edge; uniform `#00ff00` chroma background; no generic cartoon bomb, flat fire,
  text, logo, watermark or copyrighted resemblance.
- Walls: “exactly two isolated square-footprint maze obstacles”: an
  indestructible pearlstone forge pillar with brass armor and cyan core; a
  breakable shell-crate barricade with coral braces and glowing fracture seams;
  uniform `#00ff00` chroma background; no plain rectangles, generic wooden
  boxes, text, logo, watermark or copied designs.
- Items: “exactly three isolated collectible upgrades”: a coral-orange
  blast-range pearl in a four-point flame medallion; a navy-and-brass satchel
  holding two pearl bombs; aqua crystal current fins with gold trim; uniform
  `#00ff00` chroma background; no generic flat icons, retro pixels, extra
  objects, text, logo, watermark or copyrighted resemblance.
