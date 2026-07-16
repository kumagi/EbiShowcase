# Advanced-track quality checklist

This sheet tracks the second quality pass over every genre track. A checked row
means both languages and the playable WASM have passed the requirements below;
file count alone is not enough.

## Definition of “playground quality”

- The final game has at least three meaningfully different stages, courses,
  encounters, or runs. A palette swap does not count.
- The main character and important reactions have intermediate animation frames
  or a clearly readable procedural animation (anticipation, motion, impact, and
  recovery).
- Success, damage, collection, and failure have immediate visual feedback.
- A complete run has a score, time, grade, unlock, high score, or another reason
  to replay it.
- Keyboard, pointer, and touch controls all reach the complete game loop.
- The track contains enough playable intermediate lessons that a learner never
  has to absorb stage data, animation, feedback, and progression in one jump.
- Japanese and English lessons explain the added systems and link to the correct
  WASM.

## Definition of “showcase final quality”

The final game is also the promise shown on the home page. It must survive a
much harsher test than “the mechanics work”:

- Judge the capstone against a polished 2D commercial mobile game released in
  roughly the last three years, not against an early arcade cabinet or a
  decorated programming exercise. Procedural primitives remain useful for
  teaching hitboxes and rules, but they are not a substitute for readable
  human characters, environments, creatures, props, cards, and effects in the
  final showcase scene.
- High-resolution original raster artwork is welcome when it is embedded in the
  WASM, licensed for the repository, and driven by real gameplay state. Use it
  repeatedly: chapter/course/encounter changes, character entrances, enemies,
  cards, rewards, and action feedback—not as one detached splash illustration.
- At card size, a real WASM capture still reveals the genre, protagonist or
  main play piece, immediate danger or goal, and the most important action.
- The opening two seconds already contain a composed scene. Do not begin with
  a tiny actor in an empty field, an unpopulated board, or debug text floating
  over a flat fill.
- The scene has intentional depth: foreground framing, play field, and distant
  environment, or an equally deliberate board/table/stage composition.
- Inspect a real browser screenshot after integration. Legacy circles, clouds,
  bubbles, debug markers, and procedural overlays must not become giant opaque
  shapes that cover the new art or HUD at the actual WASM canvas size.
- Main actors are large enough to read, important states have intermediate
  poses, and impact/success/failure change silhouette as well as color.
- HUD information has a clear hierarchy and uses the bundled readable font;
  debug print is never the primary finished-game presentation.
- A human-facing genre must depict humans as humans at gameplay size: face,
  silhouette, pose, hands or held prop, and role must be readable without the
  name label. Likewise, cards and enemies must not collapse into interchangeable
  colored rectangles or circles.
- The final lesson explicitly labels the quality leap. Intermediate lessons may
  remain small and teachable; the capstone is allowed to combine art direction,
  animation, light, sound, UI, and replay systems non-linearly.
- The home thumbnail is recaptured from the improved running WASM at a meaningful
  moment. A title card, illustration, idle frame, or old capture does not pass.
- Visual polish never weakens the game-loop contract. `Draw` derives pixels
  from state and does not consume random state, update a retained render cache,
  or mutate camera/gameplay fields.
- [x] All 25 genre capstones include a synchronized renderer-freedom bonus:
  the build overlay delegates the original `Update` and `Layout` unchanged,
  calls its `Draw` once, then presents that one snapshot as the finished art,
  edge lines, and rectangle mosaic. These are examples, not a closed list of
  allowed renderers.

## Track pass

- [x] A Platform action — eight playable steps; the final game has four distinct
  stages, procedural run/idle/enemy animation, particles, scoring, power-ups,
  and a full keyboard/touch loop. Dedicated lessons now isolate animation,
  stage data, and game-feel/replay concepts before the final integration.
- [x] B Arena survivors — eight playable steps; three timed arena phases use
  distinct palettes and enemy traits, the boss telegraphs its dash, hits layer
  flash/particles/shake, upgrades interrupt the action safely, and clear runs
  preserve a best-kill target for replay.
- [x] C Idle/clicker — seven playable steps; three production lines animate
  their output, lifetime milestones change the bakery district, purchases burst
  into feedback, browser save/offline progress covers every machine, and the
  25K goal rolls into another retained-production run.
- [x] D Command RPG — nine playable steps; the short quest now contains three
  distinct encounters, readable enemy intents, attack/guard/heal roles, animated
  anticipation/contact/recovery, hit feedback, five quest states, and browser
  autosave across the world route.
- [x] E Platform fighter — nine playable steps; the rival contests spacing with
  approach/retreat/attack states, light and heavy moves have different timing
  and reach, attacks animate through lunge and recovery, hits use stop/particles/
  shake, three round stages vary presentation, and win streaks invite rematches.
- [x] F Merge physics — eight playable steps; three challenges vary gravity and
  target tier, consecutive merges build a combo bonus, merge events pulse and
  emit particles with proportional shake, danger remains readable, and best
  score gives the next physics run a target.
- [x] G Deckbuilder — nine playable steps; a five-floor run varies enemies,
  intents, backgrounds, rewards, rest/treasure choices, and deck growth; card
  plays animate through travel/contact with hit feedback, while run score and
  best score make different route and reward choices worth replaying.
- [x] H Slingshot battle — eight playable steps; dotted trajectory prediction,
  three peg/enemy layouts, direct-hit and ally-wave strategies, impact particles
  and shake, stage progression, and best total turns make bank shots replayable.
- [x] I Falling blocks — eight playable steps; three speed/goal/rule stages,
  line-clear anticipation and flash frames, particles and shake, active-piece
  breathing, combos, best score, and immediate replay polish the full loop.
- [x] J Match-three puzzle — eight playable steps; three boards add classic,
  blue-gem bonus, and chain-boost rules; swaps animate and rejected swaps return,
  clears pulse with particles/shake, and remaining-move bonus feeds best run.
- [x] K Sandbox — nine playable steps; three islands vary terrain, resources,
  enemies, and lantern goals; mining has anticipation/contact/recovery plus
  debris/flash/shake, while HP/speed scoring and session best reward mastery.
- [x] L Monster collection — nine playable steps; three-region expeditions use
  distinct encounter tables and research stamps, combat/capture sequences have
  staged animation and feedback, evolution completes the study, and best time
  rewards a faster complete expedition.
- [x] M Falling pairs — eight playable steps; three duels vary chain goals,
  miss limits, and incoming garbage, while chain pulse/particles/shake and idle
  bob clarify play; stage progression and best result support repeated runs.
- [x] N Maze chase — eight playable steps; three distinct mazes escalate from
  patrol to last-seen search to predictive ambush AI, with movement animation,
  pearl/damage/clear feedback, scoring, lives, and session best.
- [x] O Bomb maze — eight playable steps; three maze layouts vary walls, enemy,
  exit, and break goals, bombs accelerate their pulse before blast, explosions
  add debris/shake, and best total time rewards cleaner routes and chains.
- [x] P Tactical RPG — five playable steps; three missions vary terrain, enemy placement, goals, and turn limits, while weighted reach, weapon range, enemy intent, attack animation/feedback, total turns, and best record support replay.
- [x] Q Active-gauge RPG — six playable steps; three encounters preserve party HP and vary fast/guard/healer/boss roles, with READY queues, role commands, staged action animation, hit feedback, and time-plus-HP best evaluation.
- [x] R Branching dialogue — five playable steps; three chapters and nine scenes track courage/kindness/curiosity into four reachable endings, with typewriter, portrait/blink/entrance animation, choice feedback, and an ending gallery.
- [x] S Top-down racing — five playable steps; acceleration, grip, ordered gates, and racing-line AI build into a three-course cup with varied geometry, speed, laps, off-road feedback, finish particles, and best total time.
- [x] T Metroidvania — five playable steps; camera rooms, persistent exploration map, dash and high-jump gates, backtracking route design, three visual regions, hazards, unlock feedback, and best exploration time complete the world loop.
- [x] U Reversi — five playable steps build from an 8×8 board through legal moves, flips, passes, and score; the final game offers three meaningfully different CPU encounters (deterministic friendly, positional score-map, and opponent-reply scout), staggered flip/last-move animation, legal-move and result feedback, per-CPU browser BEST margins, localized keyboard/pointer/touch controls, and a board-first phone layout.
- [x] V Ray-cast maze FPS — six playable experiments build direction, DDA, wall projection, column rendering, and fisheye correction before the final integration; three mission data records vary maze routes, hazards, and target time, while weapon pulse/hit feedback, HP/damage/failure, grades, per-mission browser BEST, and localized keyboard/pointer/touch controls make the complete loop replayable on desktop and phone.
- [x] W Rhythm games — seven playable lessons progress from one timing pulse to four lanes, holds, rolls, chart difficulty, and a three-song tour; every song has EASY/HARD, a visible signed timing calibration, pulse/particles/shake/miss feedback, sound-ready and silent-practice states, per-chart browser BEST, language-aware UI, and keyboard/pointer/touch replay.
- [x] X Tower defense — eight playable lessons progress from routes and range through targeting, shots, placement, waves, and upgrades; COAST WATCH, REEF CAVE, and PEARL GATE provide distinct routes/resources/traits/boss rules, while intent text, target rings, projectiles/particles/shake, per-scenario grades and BEST, and language-aware keyboard/pointer/touch controls complete the replay loop.
- [x] Y Top-down adventure — eight playable lessons culminate in a four-room data route (key gate, crawler seal, tool seals, guardian); readable sword timing, hit recovery, DASH/STORM tells, persistent mobile HUD/controls, S/A/B/C results, locally saved BEST, and localized Japanese/English UI make the dungeon replayable.

## Contemporary mobile-art pass

The checked Track pass above records mechanics, teaching steps, and replay
quality. This second list is deliberately separate and stricter: it records the
2026-07 pass against the contemporary mobile-art definition above. Never infer
these checks from compilation or from an asset existing on disk; inspect a real
WASM screenshot and reject mixed-quality legacy primitives.

- [x] A Platform action — generated floating-reef world, runner, slug, and
  collision-matched terrain atlas, with objective HUD and a real opening capture.
- [x] B Arena survivors — generated arena, hero, mixed opening swarm and boss,
  transparent aura rings, survival objective, and a real opening capture.
- [x] C Idle/clicker — generated palace patisserie, chef, pastry target and
  three production-line illustrations, state-driven purchase/target feedback,
  and a real opening capture.
- [x] D Command RPG — the playable WASM opens on a legible pearl-kingdom map
  with distinct village, bridge, tower, party art, and generated battle cast.
- [x] E Platform fighter — two full human fighters, the moon-tide arena, and
  transparent hurt/attack outlines make both the fantasy and the lesson readable.
- [x] F Merge physics — the pearl nursery and seven cropped creature tiers are
  used by the live physics board, NEXT preview, target HUD, and opening layout.
- [x] G Deckbuilder — generated environment, five generated enemies, three
  distinct card illustrations, integrated combat UI, and a real opening capture.
- [x] H Slingshot battle — the reef coliseum, launch heroes, three target
  families, and stage obstacles are all present in the live trajectory game.
- [x] I Falling blocks — generated cargo-tower stage, seven material-distinct
  block faces, seeded challenge board, HOLD/NEXT previews, and a real opening
  capture that shows an immediate placement problem.
- [x] J Match-three puzzle — generated royal vault, five silhouette-distinct
  relic pieces, state-driven swap/clear/goal effects, and a real opening capture
  with no legacy overlay obscuring the HUD.
- [x] K Sandbox — three generated islands and the hero, resources, tools,
  workshop, crawler, and tide beacon all communicate the craft loop in WASM.
- [x] L Monster collection — the live expedition opens with a human navigator,
  party roster, three illustrated habitats, encounter species, orb, and badges.
- [x] M Falling pairs — the live duel uses a polished reef arena, readable
  creature pairs, a visible human rival, seeded chains, and stage-specific goals.
- [x] N Maze chase — a generated pearl labyrinth, readable large cast portraits,
  tile-scale pursuers, shell pearls, and decorated walls all ship in the WASM.
- [x] O Bomb maze — the coral forge, hero, scout, bombs, flames, walls, items,
  and an opening chain-reaction problem are all visible and state-driven.
- [x] P Tactical RPG — generated coastal battlefield, readable human blade,
  bow, and legion units, terrain-cost overlays, move/attack ranges, and a real
  opening capture.
- [x] Q Active-gauge RPG — generated moonlit arena and a unified generated cast
  for every ally, normal enemy, and boss, with gauge/impact feedback and a real
  battle capture.
- [x] R Branching dialogue — three generated chapter backgrounds, three readable
  human cast members, MIO expression poses, integrated dialogue UI, and a real
  opening/choice capture.
- [x] S Top-down racing — three generated coastal courses, generated player and
  rival vehicles, unobtrusive route/gate guidance, touch controls, and a real
  opening race capture.
- [x] T Metroidvania — three generated regions, generated hero/enemies/guardians,
  region-matched collision platforms, combat gates, and a real opening capture.
- [x] U Reversi — the CPU match now ships as a readable portrait board inside
  a generated pearl-observatory tournament scene, with live score-map teaching UI.
- [x] V Ray-cast maze FPS — three generated wall materials are sampled by the
  live ray columns, while generated guardian, key, and exit sprites retain depth.
- [x] W Rhythm games — three illustrated venues, a clearly human performer,
  distinct tap/hold/roll art, difficulty crests, rail, and results ship in WASM.
- [x] X Tower defense — the pearl-gate battlefield, three tower families, three
  enemy classes, base, and projectile art preserve path/range teaching overlays.
- [x] Y Top-down adventure — the shipped WASM uses the relic-temple scene,
  cropped hero, key, treasure chest, and guardian artwork at readable play sizes.

## Verification commands

```sh
bash scripts/build.sh
bash scripts/ralph-loop.sh verify
git diff --check
```

Browser verification must additionally exercise every final game at desktop and
mobile viewport sizes; it cannot be replaced by a canvas-exists check.
