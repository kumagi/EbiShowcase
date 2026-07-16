# Ebi Craft generated artwork

These original EbiShowcase assets are Apache-2.0. They were created with the
built-in OpenAI image-generation tool in `stylized-concept` mode on 2026-07-16,
without source images, third-party characters, trademarks, or logos.

The world source requested an exact three-panel top-down atlas with calm central
crafting clearings and detailed scenery around the edges: a turquoise moss camp
island, a moonlit indigo crystal geode, and an ember-lit volcanic island. It
contains no characters, resources, text, UI, logos, or watermarks. Generated
source:
`$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-db0e65a2-9b75-4acd-982d-022393c8d229.png`.
It was split into `island-moss.png`, `island-crystal.png`, and
`island-ember.png`.

The sprite source requested an exact 3x3 top-down atlas on a uniform `#00ff00`
background, in this order: Ebi Tenjiroh idle, Ebi Tenjiroh mining, coral
pickaxe; driftwood, pearl stone, tide crystal; completed tide lighthouse,
armored crawler crab, and workshop scaffold. Every design is original, with no
text, UI, logos, trademarks, watermarks, extra sprites, or human characters.
Generated source:
`$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-cbc7f267-809b-41d4-8290-30c1b0c208ee.png`.
It was chroma-keyed with soft matte and despill, then split into equal cells.

The game preloads only these local embeds before entering the loop. `Draw`
projects world, tile, inventory, enemy, and animation state without changing
movement, crafting, harvesting, score, timers, or randomness.
