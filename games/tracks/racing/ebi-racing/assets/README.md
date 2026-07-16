# Ebi Circuit generated artwork

These original EbiShowcase assets are Apache-2.0. They were created with the
built-in OpenAI image-generation tool in `stylized-concept` mode on 2026-07-16,
without source images, third-party characters, trademarks, or logos.

The three environment prompts requested exact top-down, 2:3 portrait racing
venues with a calm central oval reserved for the collision-course overlay:

- `coral-grand-prix-v2.png`: an evening pearl-city championship venue with one
  exact continuous oval, cyan water, gold lighting, grandstands, waterfalls,
  and a clean road for the existing physics and checkpoints. Generated on
  2026-07-16 and used as the first course after the finished-game quality pass.

- `coral-coast.png`: a sunrise turquoise coral coast with white spectator
  terraces, festival tents, sea plants, lantern buoys, and distant yachts.
  Generated source:
  `$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-d3f735e8-cb03-4878-a010-4027460800a6.png`.
- `reef-temple.png`: a bioluminescent pearl-stone reef temple with cyan
  waterfalls, shell grandstands, anemone gardens, and glowing fish. Generated
  source:
  `$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-deb25a2c-b306-4ca9-aef3-019c944a7c72.png`.
- `storm-citadel.png`: an obsidian storm citadel above cloud vortices, with
  cyan lightning rods, gold machinery, and crimson banners. Generated source:
  `$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-ad1476c2-0e18-4c99-a27f-443920c3941d.png`.

The vehicle prompt requested an exact orthographic top-down two-vehicle atlas:
Ebi Tenjiroh's coral-red, pearl-white, and gold shrimp-inspired hydro racer on
the left; an indigo, cyan, and violet manta-inspired rival on the right; both
pointing upward on a uniform `#00ff00` background with no shadows, text, UI,
logos, trademarks, or watermarks. Generated source:
`$CODEX_HOME/generated_images/019f6944-c568-7c92-94ec-77d35eaf1291/exec-624cd191-67e4-4579-998b-172004dbc06c.png`.
It was chroma-keyed with soft matte and despill, split, alpha-trimmed, and saved
as `player-car.png` and `rival-car.png`.

The game preloads only these local embeds before entering the loop. `Draw`
projects course, car, gate, and race state without changing the simulation.
