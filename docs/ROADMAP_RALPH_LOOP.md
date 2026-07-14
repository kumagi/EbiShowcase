# Ebi Showcase roadmap Ralph loop

This checklist is the source of truth for finishing README Phases 0–4. Phase 5
is explicitly outside the completion goal: it is reviewed only after this file
reaches 100%.

## Fixed decisions

- Go 1.25 is the supported baseline.
- Apache-2.0 remains the project license. Compatible third-party dependencies
  and assets (including OFL/BSD fonts) are allowed only with their original
  license and attribution recorded in `THIRD_PARTY_NOTICES.md`.
- Shader, audio, production text/UI, and camera presentation are integrated
  into all 25 final genre games, not only the labs.
- Graduation projects include an article, a copyable Go starter, and a complete
  reference implementation.
- Completion means Phases 0–4. New Phase 5 genres are not required.

## Loop protocol

Work on exactly one unchecked item at a time.

```sh
node scripts/roadmap-ralph-loop.mjs status
node scripts/roadmap-ralph-loop.mjs next
node scripts/roadmap-ralph-loop.mjs evidence P1-U-AUDIT
# implement and verify the task
node scripts/roadmap-ralph-loop.mjs check P1-U-AUDIT
node scripts/roadmap-ralph-loop.mjs verify
```

`check` refuses to mark an item unless its evidence file passes the task's
required checklist. Use `uncheck` when later work invalidates earlier evidence.
At phase boundaries run `verify --full`; before declaring the roadmap complete,
`complete` must pass.

## Phase 0 — Freeze the baseline

- [x] `P0-01` — Synchronize Go 1.25 across `go.mod`, README, AGENTS, and CI.
- [x] `P0-02` — Add `THIRD_PARTY_NOTICES.md` and document the compatible-license policy.
- [x] `P0-03` — Add an automated metadata check for counts, bilingual pairs, thumbnails, and final-game links.
- [x] `P0-04` — Add a fast local build mode that skips unchanged OGP generation while CI keeps the full build.
- [x] `P0-05` — Establish the roadmap evidence schema and Desktop/Tablet/Phone audit template.
- [x] `P0-06` — Link this loop from README and AGENTS and document the one-task completion protocol.
- [x] `P0-07` — Record a clean baseline: tests, 208/208, VFX 29/29, lesson structure, and GitHub Pages build.

## Phase 1 — Bring U–Y to playground quality

Each audit covers three meaningful runs/stages, animation, feedback, replay,
keyboard/pointer/touch completion, Japanese/English parity, and pure-logic tests.

- [x] `P1-U-AUDIT` — Audit Reversi and record exact gaps.
- [x] `P1-U-POLISH` — Add three CPU personalities/difficulties, move/flip animation, replay records, and complete controls.
- [x] `P1-U-PASS` — Verify Reversi on three viewports and add U to the advanced checklist.
- [x] `P1-V-AUDIT` — Audit Raycaster and record exact gaps.
- [x] `P1-V-POLISH` — Add three map missions, damage/failure, weapon feedback, grade/BEST, and complete controls.
- [x] `P1-V-PASS` — Verify Raycaster on three viewports and add V to the advanced checklist.
- [x] `P1-W-AUDIT` — Audit Rhythm, including audio unlock and timing behavior.
- [x] `P1-W-POLISH` — Add timing calibration, readable judgement feedback, robust song selection, and complete controls.
- [x] `P1-W-PASS` — Verify all three songs/two difficulties and add W to the advanced checklist.
- [x] `P1-X-AUDIT` — Audit Tower Defense and record exact gaps.
- [x] `P1-X-POLISH` — Ensure three meaningful scenarios, readable enemy intent, touch building/upgrades, and replay scoring.
- [x] `P1-X-PASS` — Verify Tower Defense on three viewports and add X to the advanced checklist.
- [x] `P1-Y-AUDIT` — Audit Top-down Adventure and record exact gaps.
- [x] `P1-Y-POLISH` — Improve movement/attack transitions, room readability, boss tells, touch tools, and replay scoring.
- [x] `P1-Y-PASS` — Verify the complete dungeon on three viewports and add Y to the advanced checklist.

## Phase 2A — Build the four technical labs

### Shader Lab

- [x] `P2-SH-01` — Create the bilingual Shader Lab hub, shared Kage helpers, and a safe fallback path.
- [x] `P2-SH-02` — Teach uniforms and time-driven palette/flash effects.
- [x] `P2-SH-03` — Teach UV distortion for water, heat haze, and impact waves.
- [x] `P2-SH-04` — Teach color separation and damage/status treatment.
- [x] `P2-SH-05` — Teach offscreen passes, downsampling, and controlled faux/real blur.
- [x] `P2-SH-06` — Add Shader Lab tests, mobile checks, diagrams, thumbnails, and bilingual completion evidence.

### Audio Lab

- [x] `P2-AU-01` — Create the bilingual Audio Lab hub and user-gesture-safe audio context.
- [x] `P2-AU-02` — Teach original pure-Go waveform and one-shot SE synthesis.
- [x] `P2-AU-03` — Teach ADSR envelopes and parameterized sound families.
- [x] `P2-AU-04` — Teach reusable SE voices/queues without per-hit allocation spikes.
- [x] `P2-AU-05` — Teach looping BGM, volume groups, ducking, pause, and resume.
- [x] `P2-AU-06` — Add Audio Lab tests, mobile checks, diagrams, thumbnails, and bilingual completion evidence.

### Text / UI Lab

- [x] `P2-UI-01` — Register selected font licenses/notices and create a `text/v2` bilingual font loader.
- [x] `P2-UI-02` — Teach measurement, alignment, wrapping, and CJK line breaking.
- [x] `P2-UI-03` — Teach reusable panels, nine-slice-style frames, gauges, and icon labels.
- [x] `P2-UI-04` — Teach keyboard/touch focus, disabled state, and action mapping.
- [x] `P2-UI-05` — Teach dialogue, menus, scroll areas, and accessible status feedback.
- [x] `P2-UI-06` — Add UI Lab tests, mobile checks, diagrams, thumbnails, and bilingual completion evidence.

### Camera Lab

- [x] `P2-CA-01` — Create a pure camera-state package and bilingual Camera Lab hub.
- [x] `P2-CA-02` — Teach follow, smoothing, clamping, and coordinate conversion.
- [x] `P2-CA-03` — Teach dead zones, look-ahead, and target switching.
- [x] `P2-CA-04` — Teach deterministic shake, hit stop coordination, and recovery.
- [x] `P2-CA-05` — Teach room transitions, framing, letterbox, and viewport-safe composition.
- [x] `P2-CA-06` — Add Camera Lab tests, mobile checks, diagrams, thumbnails, and bilingual completion evidence.

## Phase 2B — Apply all four pillars to every final genre game

Every item requires: one meaningful shader treatment, original audio feedback
and music/ambience, production `text/v2` UI instead of final-screen DebugPrint,
camera/framing feedback appropriate to the genre, three viewport audits, both
languages, and regression tests. Static-board games use framing/transition
camera work rather than artificial scrolling.

- [x] `P2-G01` — Platform action (`platformer`).
- [x] `P2-G02` — Arena survivors (`survivors`).
- [x] `P2-G03` — Idle/clicker (`clicker`).
- [x] `P2-G04` — Command RPG (`rpg`).
- [x] `P2-G05` — Platform fighter (`fighting`).
- [x] `P2-G06` — Merge physics (`merge-physics`).
- [x] `P2-G07` — Deckbuilder (`deckbuilder`).
- [x] `P2-G08` — Match-three (`match3`).
- [x] `P2-G09` — Falling blocks (`falling-blocks`).
- [x] `P2-G10` — Slingshot battle (`slingshot`).
- [x] `P2-G11` — Sandbox (`sandbox`).
- [x] `P2-G12` — Monster collection (`monster-collection`).
- [x] `P2-G13` — Falling pairs (`falling-pairs`).
- [x] `P2-G14` — Maze chase (`maze-chase`).
- [x] `P2-G15` — Bomb maze (`bomb-maze`).
- [x] `P2-G16` — Reversi (`reversi`).
- [x] `P2-G17` — Tactical RPG (`tactics`).
- [x] `P2-G18` — Active-gauge RPG (`active-rpg`).
- [x] `P2-G19` — Branching dialogue (`visual-novel`).
- [x] `P2-G20` — Top-down racing (`racing`).
- [x] `P2-G21` — Metroidvania (`metroidvania`).
- [x] `P2-G22` — Raycaster (`raycaster`).
- [x] `P2-G23` — Rhythm (`rhythm`).
- [x] `P2-G24` — Tower Defense (`tower-defense`).
- [x] `P2-G25` — Top-down Adventure (`topdown-adventure`).

## Phase 3 — Cross-cutting guides and graduation projects

- [x] `P3-SAVE-01` — Extract a versioned save model with WASM/native storage adapters and tests.
- [x] `P3-SAVE-02` — Publish a bilingual save/autosave/recovery/scene-transition guide.
- [x] `P3-ARCH-01` — Publish a bilingual project-structure guide using a real refactor into `cmd`/`internal`/`assets`.
- [x] `P3-DIST-01` — Publish a bilingual local/GitHub Pages/licensing/release checklist.
- [x] `P3-PERF-01` — Add benchmarkable allocation, pooling, culling, and spatial-query examples.
- [x] `P3-PERF-02` — Publish the bilingual performance guide with before/after measurements.
- [x] `P3-GRAD-01` — Ship the 60-second arcade brief, Go starter, tests, and reference game.
- [x] `P3-GRAD-02` — Ship the three-room exploration brief, Go starter, tests, and reference game.
- [x] `P3-GRAD-03` — Ship the three-stage puzzle brief, Go starter, tests, and reference game.
- [x] `P3-GRAD-04` — Add a graduation hub linking prerequisites, briefs, source, and release steps.

## Phase 4 — Navigation, metrics, and sustainable operations

- [x] `P4-NAV-01` — Add a first-30-minutes route from setup through the first complete game.
- [x] `P4-NAV-02` — Add interest entrances for action, puzzle, strategy, story, and presentation.
- [x] `P4-NAV-03` — Add visible progress links from lessons to labs and graduation briefs.
- [x] `P4-METRIC-01` — Report playable, advanced-quality, mobile-audit, lab, graduation, and pending-feedback metrics.
- [x] `P4-OPS-01` — Document and automate the 20-response triage and 500-response archive thresholds.
- [x] `P4-OPS-02` — Add a rotating final-game audit schedule and evidence freshness check.
- [x] `P4-RELEASE` — Run the full release audit, rebuild every WASM/OGP, and prove Phase 0–4 complete.

## Phase 5 admission review — not part of completion

After `P4-RELEASE`, record whether Stealth, Farming, Autobattler, or Twin-stick
teaches a genuinely missing calculation pattern. Do not add a genre only to
increase the count.
