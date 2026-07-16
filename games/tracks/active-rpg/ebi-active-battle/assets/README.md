# Ebi Active Battle generated artwork

These original assets are part of EbiShowcase and are licensed under
Apache-2.0. They were created with the built-in OpenAI image-generation tool
in `stylized-concept` mode on 2026-07-16. No third-party characters, logos, or
source images were used.

## Files and prompts

- `active-rpg-moonlit-arena.png` — portrait 2:3 premium painted mobile-RPG
  battle environment. Prompt: “Original moonlit coastal arena with reflective
  water, black basalt steps, ruined sea-shrine arches, luminous teal tide
  crystals, and a storm vortex; portrait battle composition with clear lower
  arena; midnight navy, turquoise and gold; no characters, text, UI, logo,
  trademark, or watermark.” Generated source:
  `$CODEX_HOME/generated_images/019f68ba-013b-7b51-8f4b-726cd03cbb36/exec-5ded3a27-8fdb-40c3-8df2-b515e49de6fb.png`.
- `active-rpg-tenjiroh.png` — Ebi Tenjiroh protagonist cutout. Prompt:
  “Original heroic anthropomorphic red shrimp adventurer with coral saber,
  pearl buckler and navy scarf; full-body three-quarter battle pose facing
  right; premium hand-painted anime mobile-game cutout; flat `#00ff00`
  chroma-key background; no text, UI, logo, trademark, watermark, cast shadow,
  or extra character.” Generated source:
  `$CODEX_HOME/generated_images/019f68ba-013b-7b51-8f4b-726cd03cbb36/exec-2546a10d-38b9-40c8-82e6-377b7b4870c6.png`.
- `active-rpg-storm-king.png` — Storm King boss cutout. Prompt: “Original
  colossal storm crab monarch with asymmetric claws, obsidian shell, turquoise
  tide core and broken gold crown; full-body three-quarter pose facing left;
  premium hand-painted anime mobile-game boss cutout; flat `#00ff00`
  chroma-key background; no text, UI, logo, trademark, watermark, cast shadow,
  or extra character.” Generated source:
  `$CODEX_HOME/generated_images/019f68ba-013b-7b51-8f4b-726cd03cbb36/exec-62e23e1c-5102-4265-8a63-ef019c562220.png`.

The two character sources were processed with the image-generation skill's
`remove_chroma_key.py` helper (`--auto-key border --soft-matte --despill`) and
then downscaled for a practical WASM payload. Their four corners are fully
transparent.

## Supporting cast sheet

`active-rpg-mage.png`, `active-rpg-shell.png`, `active-rpg-wisp.png`, and
`active-rpg-scout.png` were cropped from one coherent 2×2 generated sheet.
Prompt: “Exactly four original full-body characters in isolated equal
quadrants: human tide mage with pearl staff and white-turquoise coat; friendly
armored hermit-crab companion with pearl shell; moon wisp spirit with violet
mask and cyan flame body; human clockwork legion scout with twin short blades.
Premium hand-painted anime mobile RPG style, identical moonlit rendering,
wide green gutters, flat #00ff00 background, no shadows, text, UI, logos,
trademarks, watermarks, labels, divider lines, or overlapping quadrants.”
Generated source:
`$CODEX_HOME/generated_images/019f68ba-013b-7b51-8f4b-726cd03cbb36/exec-2e16f57f-d322-4288-ae4b-4ded7157bc0c.png`.
The sheet was chroma-keyed once, then split into exact quadrants and downscaled.
All regular encounter actors now use this same generated-art pipeline.
