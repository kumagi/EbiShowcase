# Power-up Adventure generated art

The following original raster assets were generated for this repository with
the built-in OpenAI image generation tool and are distributed under the
repository's Apache License 2.0. They intentionally avoid existing game
franchises and third-party characters.

## `assets/platformer-sunrise-archipelago.png`

Prompt:

> Create a premium modern commercial mobile-game environment background for a
> whimsical ocean-cliff platforming adventure starring a brave shrimp hero.
> Background scenery only. Portrait 2:3 composition with a luminous coastal
> archipelago at sunrise, layered turquoise cliffs, misty waterfalls, floating
> islands, jungle canopy and a glowing sky portal in the upper-right. Preserve
> an open playable middle area. Highly polished hand-painted 2D game art with
> strong far/middle/near depth. No characters, collidable-looking platforms,
> UI, text, logos, watermark or imitation of an existing franchise.

Generated source: 1024x1536 RGB PNG. Final optimized asset: 720x1080 RGB PNG.

## `assets/platformer-ebi-runner.png`

Prompt:

> Create exactly one full-body brave anthropomorphic coral-orange shrimp child
> hero with large teal eyes, curved antennae, cream explorer scarf, navy shell
> backpack and gold cuffs, running toward the right. Premium hand-painted 2D
> mobile game character with a crisp silhouette readable at 48 pixels. Place it
> on a perfectly uniform `#ff00ff` chroma-key background with generous padding.
> No shadow, floor, text, logo, watermark, extra character or franchise copy.

Generated source: 1024x1536 RGB PNG. Chroma-key removal used the imagegen skill's
`remove_chroma_key.py` helper with a soft matte and despill. Final optimized
asset: 512x768 RGBA PNG.

## `assets/platformer-reef-slug.png`

Prompt:

> Create exactly one full-body mischievous tropical rock-slug enemy: a squat
> teal-and-purple armored sea slug with a craggy coral shell, luminous amber
> eyes, tiny fangs and four stubby feet, charging toward the left. Premium
> hand-painted 2D mobile game character with a crisp, playful silhouette. Place
> it on a perfectly uniform `#00ff00` chroma-key background with generous
> padding. No shadow, floor, text, logo, watermark, extra character or
> franchise copy.

Generated source: 1536x1024 RGB PNG. Chroma-key removal used the imagegen skill's
`remove_chroma_key.py` helper with a soft matte and despill. Final optimized
asset: 768x512 RGBA PNG.

## `assets/platformer-reef-platform.png`

Prompt:

> Create one isolated wide floating tropical reef-island platform tile with a
> straight readable walkable top, thick emerald grass, tiny white flowers,
> warm sandstone and turquoise coral rock underneath, hanging roots, coral
> tips and naturally tapered ends. Use a wide 3:1 silhouette and premium
> hand-painted mobile platformer rendering matched to a luminous tropical
> archipelago. Place it on a perfectly uniform `#ff00ff` chroma-key background.
> No character, enemy, UI, text, logo, watermark, cast shadow or franchise
> imitation.

Generated at 2048x768 RGB, chroma-keyed with a soft matte and despill, then
optimized to 1024x341 RGBA. The game-local terrain atlas cuts this original
plate into a flowered left cap, repeatable reef center, planted right cap and
complete tapered small-island silhouette. Long ground is composed from the
three cap/center pieces while short ledges use the compact island treatment;
their collision rectangles are unchanged.
