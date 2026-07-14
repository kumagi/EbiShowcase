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

## Verification commands

```sh
bash scripts/build.sh
bash scripts/ralph-loop.sh verify
git diff --check
```

Browser verification must additionally exercise every final game at desktop and
mobile viewport sizes; it cannot be replaced by a canvas-exists check.
