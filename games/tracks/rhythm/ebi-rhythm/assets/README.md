# Ebi Rhythm generated artwork

Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.

These original assets were generated for this repository on 2026-07-16 with
OpenAI image generation, then locally split, chroma-keyed, and alpha-trimmed.
They contain no third-party game artwork or logos.

- `stage-*.png`: three-panel prompt for original ocean-tour venues: Sunrise
  Harbor, Neon Reef, and Tempest Parade; polished contemporary Japanese mobile
  rhythm-game environment illustration with an open four-lane play area.
- `artist-*.png`: three-pose prompt for the same original adult human ocean
  artist, with aqua-violet braids, navy pearl coat, and shell-shaped keytar.
- `note-*.png`: 3x2 prompt for four colored tap gems, a pearl-comet hold note,
  and a shell-resonator roll note, designed to remain readable at small sizes.
- `difficulty-*.png`, `progress-rail.png`, `judgment-rail.png`: 2x2 prompt for
  original pearl, shell, cyan, gold, and violet rhythm-interface ornaments.

Generated source files were produced under
`~/.codex/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/` and the
working chroma/split files under `tmp/imagegen/rhythm/`. Runtime state changes
only in `Update`; `Draw` selects and composites these images without advancing
the song or changing judgment state.
