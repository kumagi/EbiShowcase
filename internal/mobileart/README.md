# Mobile showcase artwork

This package contains original raster artwork used by the Ebitengine capstone
games. All files in this directory are licensed under Apache-2.0 with the rest
of EbiShowcase. The images are embedded in the WASM binaries; they are not
remote dependencies.

`mobileart.go` deliberately exposes drawing helpers only. Game rules, choices,
timers, combat, and movement remain normal Go state advanced by `Update`.
Changing one of these images changes presentation without changing the game.
The game calls `Preload` in its constructor, so the first `Draw` does not lazily
decode images or mutate a retained rendering cache.

## Visual-novel set

The initial Moon Lantern Stories set was created with the built-in OpenAI image
generation tool in `stylized-concept` mode on 2026-07-16. The generated source
PNGs were copied into `assets/`; flat green backgrounds on the character sheets
were removed locally with the image-generation skill's chroma-key helper.

- `visual-novel-harbor.png`: premium mobile-game background; moon-lantern
  festival harbor, turquoise water, lighthouse, silver ribbon clue; portrait
  staging composition; no characters, UI, text, logo, or watermark.
- `visual-novel-clock.png`: premium mobile-game background; hidden astronomical
  archive with immense brass moon clock, books, stairs, and aqua doorway;
  portrait staging composition; no characters, UI, text, logo, or watermark.
- `visual-novel-observatory.png`: premium mobile-game background; moonlit glass
  ocean observatory, brass telescope and coral planters; portrait staging
  composition; no characters, UI, text, logo, or watermark.
- `navigator.png`: original young adult ocean-observatory navigator, navy and
  pearl uniform, translucent star chart, moonlit rim light.
- `navigator-worried.png`, `navigator-surprised.png`, `navigator-joy.png`: a
  three-panel expression pass made with `navigator.png` as the design reference,
  then split after chroma removal. Scene data selects these real facial and hand
  poses; the labels do not merely claim that an expression changed.
- `researcher.png`: original young adult clockwork researcher, teal waistcoat,
  charcoal coat, moon-lantern mechanism.
- `keeper.png`: original older moon-clock keeper, violet astronomical robes,
  clock key and star ledger.

The character prompts required one centered human, readable face and hands,
complete silhouette, original design, and a uniform `#00ff00` background. They
explicitly excluded chibi, pixel-art, flat-vector, logos, watermarks, and
copyrighted-character resemblance.
