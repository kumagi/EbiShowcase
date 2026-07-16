# Ebi Depths generated artwork

These original EbiShowcase assets are Apache-2.0. They were created with the
built-in OpenAI image-generation tool in `stylized-concept` mode on 2026-07-16
without source images, third-party characters, trademarks, or logos.

The environment source was prompted as an exact three-strip atlas: “SUNKEN
GARDENS with bioluminescent plants and sea ruins; CLOCKWORK ABYSS with brass
gears, violet crystal machinery and bridges; TEMPEST SANCTUM with black ocean
cliffs, storm citadel and cyan lightning. Three equal wide side-view strips,
coherent premium hand-painted anime mobile-metroidvania style, black separator
gutters, no characters, text, UI, logo, trademark, watermark, grid or labels.”
Generated source:
`$CODEX_HOME/generated_images/019f68ba-013b-7b51-8f4b-726cd03cbb36/exec-9704e02c-bbf9-4241-8bfd-66414a811d58.png`.
It was split into `depths-sunken-gardens.png`,
`depths-clockwork-abyss.png`, and `depths-tempest-sanctum.png`.

The cast source was prompted as an exact 2×2 atlas: “Ebi Tenjiroh, athletic
red shrimp explorer with navy scarf, coral saber and pearl lantern; hostile
obsidian abyss beetle with cyan eyes; ancient clockwork guardian knight with
violet core, halberd and shield; luminous winged tide relic spirit. Full-body
side-view characters, coherent premium hand-painted anime mobile-metroidvania
style, flat #00ff00 background, wide gutters, no crossing, shadows, labels,
text, UI, logos, trademarks or watermarks.” Generated source:
`$CODEX_HOME/generated_images/019f68ba-013b-7b51-8f4b-726cd03cbb36/exec-627ff0b5-3a91-4e2a-aa30-72d782ba66e8.png`.
It was chroma-keyed with soft matte and despill, split into exact quadrants,
and saved as `depths-tenjiroh.png`, `depths-beetle.png`,
`depths-guardian.png`, and `depths-spirit.png`.

The platform source requested exactly three equal side-view ledges on a uniform
`#00ff00` background: sunken limestone ruins with roots and flowers; clockwork
abyss rock with gears, pipes, and amethyst; tempest-sanctum obsidian with gold
citadel fragments and lightning cracks. Each has a level walkable top, a deep
foreground rock silhouette, no characters, text, UI, logos, trademarks, or
watermarks. Generated source:
`$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-29cab039-93e5-451e-ac9c-228cc1e8a0f6.png`.
It was chroma-keyed with soft matte and despill, split into thirds,
alpha-trimmed, and saved as `platform-gardens.png`, `platform-abyss.png`, and
`platform-sanctum.png`.

The game preloads only these local embeds before entering the loop. `Draw`
reads them without changing exploration, combat, map, or ability state.
