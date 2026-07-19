# Ebi Fighters generated artwork

These original EbiShowcase assets are Apache-2.0. They were created with the
built-in OpenAI image-generation tool in `stylized-concept` mode on 2026-07-16,
without source images, third-party characters, trademarks, or logos.

`moon-tide-arena.png` was prompted as a 2:3 portrait, side-view championship
arena above a bioluminescent sea: pearl-stone combat floor, shell arches,
jellyfish lanterns, glass wave barriers, spectators, a distant coral city,
moon and aurora, with no fighters, text, UI, logos, or watermarks. Generated
source:
`$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-0051b280-eca3-408d-ad78-8c200c9a14c7.png`.

The player source requested an exact 2x2 atlas of the same clearly human young
adult ocean duelist, always facing right: ready, coral-saber attack, hurt, and
slumped KO. His original coral, pearl, navy, and gold armor is inspired by
shrimp shells without copying a third-party character. Generated source:
`$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-2246431a-790d-4157-9816-8c0aa7edb91e.png`.

The rival source requested an exact 2x2 atlas of the same clearly human young
adult abyss-tide champion, always facing left: ready, crescent-blade attack,
hurt, and slumped KO. Her indigo, violet, silver, and cyan manta-inspired armor
is an original design. Generated source:
`$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-bc7260ea-66ab-4d8f-aaab-509894515cfa.png`.

Both atlases used a uniform `#00ff00` background and were chroma-keyed with
soft matte and despill before being split into equal quadrants. The game
preloads only these local embeds before entering the loop. `Draw` selects a
pose from fighter state and never advances combat, input, timers, or effects.

The two eight-frame motion sheets were generated with the built-in OpenAI
image-generation tool on 2026-07-20, using each fighter's ready, attack, hurt,
and KO images above as identity/style references. Both prompts requested an
exact 4x2 sequence: ready, anticipation, early step, early swing, impact,
follow-through, recovery, and return to ready. Camera, scale, baseline,
costume, facing direction, and a flat chroma-key background were held constant.
Generated sources:

- `player-motion.png`:
  `$CODEX_HOME/generated_images/019f73d0-487d-7602-a7a1-e8c8f492bc92/exec-07540675-c3f0-41b9-b916-54a9c014223c.png`
- `rival-motion.png`:
  `$CODEX_HOME/generated_images/019f73d0-487d-7602-a7a1-e8c8f492bc92/exec-c3928b21-8a79-4024-9123-6041be051881.png`

The repository copies were chroma-keyed with soft matte and despill. Runtime
code splits each sheet into eight cells and crops transparent padding without
changing the combat state. `Update` chooses an animation phase; `Draw` only
selects the corresponding frame.
