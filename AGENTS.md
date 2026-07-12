# AGENTS.md — Ebi Showcase contributor guide

## Purpose

Ebi Showcase is a static, bilingual learning site that teaches 2D game development through small browser experiences. Every core or genre lesson is a pair:

1. a genuinely playable mini game written in Go with Ebitengine and compiled to WebAssembly; and
2. a Japanese and English article that explains the game's main technique in language a primary-school student can follow.

The long-term completion gate is `105/105` playable core + genre entries. A core/genre lesson page or design mockup without a working Ebitengine game is not complete. The separately counted Visual Effects Lab is the intentional exception: its Ebitengine programs are focused interactive drawing toys rather than full games.

## Start here: repository mental model

Before editing, identify which of these four surfaces the request belongs to. They have different completion rules and different sources of truth.

1. **Core curriculum (LEVEL 01–12)** — twelve small, complete games that establish the shared vocabulary. LEVEL 01 is the one canonical introduction to `Update`, `Draw`, and optional `Layout`; later pages build on it instead of teaching the engine from zero again.
2. **Visual Effects Lab** — eight focused drawing toys between core and genre tracks. These deliberately do not need game goals because they isolate the rendering pipeline. Their pages are generated from `scripts/gen-visual-effects.mjs`.
3. **Genre tracks** — together with the 12 core lessons, these make up the 105-gate curriculum. Each step must be a genuine game or game-system exercise, and each track ends in an integrated genre game. Track pages are normally hand-authored.
4. **Cross-cutting guides and site infrastructure** — for example the game-data guide, bilingual entry pages, shared WASM loader, feedback form, stylesheet, and learning labs. These explain or support the curriculum but are not Playable-count entries.

The public learning path is therefore:

`LEVEL 01 game loop → core mechanics → optional Visual Effects Lab → genre specialization → data/architecture thinking`

Do not flatten these into one undifferentiated list. A rendering toy is not a missing game, a guide is not PLAYABLE, and a genre step must not reintroduce `Update`/`Draw` as if LEVEL 01 did not exist.

## New-agent orientation checklist

For a context-free task, do these in order before changing files:

1. Read this file completely.
2. Run `git status --short`; assume every existing change belongs to the user or another agent.
3. Run `bash scripts/ralph-loop.sh status` and `bash scripts/ralph-loop.sh next` for curriculum work.
4. Identify whether the target is core, VFX-generated, genre-track, or infrastructure using the mental model above.
5. Read the target JA page, EN page, preceding lesson, following lesson, and track hub before writing. Preserve the teaching progression and pager links.
6. Read the corresponding Go game and one recently completed neighboring game for local conventions. Do not copy a generic game and merely recolor it.
7. Check whether the file is generated. Search its header and `scripts/` before hand-editing.
8. Make the smallest coherent scoped change, then use the verification ladder in this guide.

Use stable curriculum IDs such as `tracks/falling-blocks/falling-cell`, not remembered numeric order. Inserting a separately counted group such as VFX changes displayed `order` values even though the 105-gate sequence is unchanged.

## Non-negotiable principles

- Game logic, input, simulation, collision, and rendering belong in Go + Ebitengine. HTML, CSS, and JavaScript are only for the static showcase and shared WASM loader.
- Everything in the repository is Apache License 2.0 compatible. Do not add assets, fonts, code, or copied game names/artwork with incompatible or unclear licensing.
- Prefer original simple shapes; the shared protagonist 海老・天次郎 (Ebi Tenjiroh) sprite is allowed as the curriculum hero. Ebi-themed presentation otherwise. Famous games are references for mechanics, not permission to copy their art, text, characters, music, maps, or branding.
- Every game must work with keyboard/mouse on PC and touch input on phones and tablets.
- GitHub Pages project URLs must work. Use relative links; do not assume the site is hosted at `/`.
- Japanese and English are independent, shareable URLs. Keep both versions in sync.
- `dist/` is generated and ignored. Never hand-edit or commit it.
- Preserve unrelated working-tree changes. Several agents may continue the curriculum over time.

## Curriculum and Ralph Loop

The curriculum order is derived from the Japanese home page and track hubs by `scripts/curriculum.mjs`.

```sh
bash scripts/ralph-loop.sh status  # playable/total count
bash scripts/ralph-loop.sh next    # first missing game in curriculum order
bash scripts/ralph-loop.sh verify  # rebuild all WASM and run structural checks
```

Continue from the entry returned by `next`; do not skip ahead merely because a later game is easier. An entry becomes structurally PLAYABLE when its expected `main.go` exists:

- core Flappy: `game/main.go`
- other core lessons: `games/core/<slug>/main.go`
- track lessons: `games/tracks/<track>/<slug>/main.go`

The build creates a separate game at `dist/play/<slug>/` for every implementation. Slugs must therefore remain unique across tracks.

## Visual Effects Lab (between core and tracks)

- The Visual Effects Lab is a distinct group that sits between the 12 core lessons and the genre tracks. It teaches Ebitengine's drawing pipeline (`GeoM` translate/rotate/scale, `ColorScale`, `Blend`/`BlendLighter`, `SubImage` frame animation, particle slices) through hands-on toys that are deliberately "less than games".
- Its games live at `games/tracks/visual-effects/<slug>/main.go` with `vfx-` slugs, and its pages at `web/{ja,en}/tracks/visual-effects/`. It appears in `curriculum.mjs` ordered right after LEVEL 12.
- It is counted **separately** from the `105/105` gate. `curriculum.mjs` exports `gated` (core + genre tracks) which drives `next`/`status`; the summary reports `vfx: { total, playable }` on the side. Do not fold the Visual Effects Lab into the 105 total.
- The Lab's bilingual hub and step pages are generated by `scripts/gen-visual-effects.mjs` from a content table (keeps JA/EN in sync and skeleton-compliant). Edit the table and re-run it rather than hand-editing the generated pages; the build still runs `embed-lesson-sources` and `insert-feedback-form` afterward. Interactive lab demos use shared handlers in `web/learn.js` (`data-lab` kinds: translate, geom, colorscale, opacity, blend, sheet, spray, spellbook).
- Keep its pedagogical promise visible: these are **hands-on drawing toys, deliberately less than games**. A VFX step should isolate one visual operation, provide immediate manipulation, and connect it to an in-game use. Do not add artificial win/lose rules merely to resemble genre lessons.
- The progression is intentional: place → transform/pivot → tint → alpha/draw order → additive blend → sprite-sheet animation → particles → composed spells. New material must preserve prerequisite order or explicitly justify a new insertion.

## 海老・天次郎 (Ebi Tenjiroh) texture atlas (downloadable asset)

- The curriculum protagonist is **海老・天次郎 (Ebi Tenjiroh)** — Japanese articles use 海老・天次郎; English articles use Ebi Tenjiroh. Do not write “Ebi Boy”.
- The character sprite atlas is generated by pure-Go software rendering (no GPU) so frames stay pixel-aligned, consistent, and transparent. Layout is one shared source of truth in `internal/atlaslayout` (frame size 96×96, 15 strips: idle/walk/run/attack/hurt × down/up/side). Left-facing = draw the `side` frames flipped horizontally.
- Regenerate with `go run ./cmd/gen-atlas`. It writes `internal/heroatlas/ebi-boy-atlas.png` (embedded into WASM games; filename kept for asset stability), plus the downloadable `web/assets/ebi-boy-atlas.{png,json}` and `web/assets/ebi-boy-atlas-LICENSE.txt`. Re-run whenever the character art or layout changes, then rebuild.
- Runtime access is `internal/heroatlas`: `Anim("walk-side")` returns the ordered `*ebiten.Image` frames (via `SubImage`), `FPS(name)` the suggested cadence, `Sheet()` the whole image. `vfx-walk` consumes it as the STEP 06 demo.
- Soft effect sprites (fire / water / spark / bolt) are generated by `go run ./cmd/gen-vfx` into `internal/vfxsprites/` (WASM embed) and `web/assets/vfx-*.png` (HTML labs). `vfx-spells` and the spellbook lab use these textures with additive blending and motion—not flat colored dots.
- The generated atlas + metadata are dedicated to the public domain (CC0-1.0); the surrounding source stays Apache-2.0. The download link and layout explanation live on the `vfx-walk` lesson (`download` field in the generator).

## Definition of done for one lesson

- A distinct, playable Ebitengine game exists at the expected path. It should demonstrate the lesson's named mechanic, not merely reskin a generic template.
- The game has a goal, failure or retry behavior where appropriate, clear feedback, and understandable controls.
- Keyboard and touch controls are both implemented. Pointer controls should use Ebitengine input APIs, not page-side JavaScript.
- `web/ja/<route>/index.html` and `web/en/<route>/index.html` embed `../../../../play/<slug>/` (or the correct relative depth).
- Both articles explain the mechanic step by step, including a small representative Go snippet. Japanese should avoid unexplained jargon and use concrete examples a child can picture.
- Core lessons (curriculum order 1–12) also embed the full `main.go` with a short “look here first” map and copy button (`data-embed-source`). Track lessons use the code-lesson snippet only; point learners to the repo for the full file.
- Shared hero art may live in `internal/hero`. Lesson copy must name the on-screen character 海老・天次郎 / Ebi Tenjiroh when that sprite is drawn, and must not claim the HTML listing is a fully self-contained single file if shared packages are imported.
- The relevant hub/card says PLAYABLE or otherwise clearly identifies the finished lesson when the surrounding design supports statuses.
- Go formatting, JS/WASM compilation, and the Ralph verification pass.
- A real browser loads one Canvas without console errors. Check a desktop viewport and a phone-size viewport; the lesson page must not overflow horizontally.

## Teaching and writing contract

The first pages are the style reference, especially `web/{ja,en}/games/tap-target/index.html` and the generated VFX STEP 01 page.

- Lead with a concrete thing the learner can see or do, then name the abstraction. Japanese copy should be understandable to an upper-primary-school learner without talking down to them.
- Teach one main idea and at most one closely related supporting idea per step. State what came from the previous lesson and what the new piece adds.
- Use the pattern **play → observe/predict → manipulate a small lab → inspect representative Go → explain why → offer one challenge**.
- Use concrete metaphors only when they map accurately to the code. LEVEL 01's flipbook explanation is the canonical model: `Update` changes numbers/state; `Draw` renders the current state. Do not later claim Draw runs rules or Update paints the screen.
- Progressive disclosure is intentional. LEVEL 01 foregrounds Update and Draw while putting Layout and build details in an optional section. Keep advanced implementation details out of the learner's first conceptual step unless needed to understand the mechanic.
- A `DEEP DIVE` must explain the named mechanic, not merely describe controls. The three `concept-row` cards should form an ordered explanation. The motion lab must let the learner change or step through the same concept. The code snippet must correspond to the real Go implementation.
- Explain data and state transitions, not just visible results: where values live, when they change, and why update order matters.
- Famous titles may appear only as genre signposts. Page titles, characters, art, story, terminology, and assets must remain original and Ebi-themed.
- Keep JA and EN semantically aligned, but write natural prose in each language rather than word-for-word translation.
- Accessibility is part of the lesson: meaningful iframe titles, button `type`, labels/ARIA where the lab needs them, readable status feedback, and no information conveyed only by color.

## Article consistency (current curriculum policy)

- Teach Ebitengine from the game loop first: LEVEL 01 owns Update / Draw / (optional) Layout. Later lessons review or extend that frame; they do not reintroduce the engine as if from zero.
- When a lesson draws the shared hero sprite (`internal/hero`), Japanese articles must call out **海老・天次郎** and English articles **Ebi Tenjiroh**, and both must say that hit tests remain simple shapes.
- Full-source blurbs must not claim `internal/hero` if that package is unused. Shape-only games say the listing is the entire game; hero games say main logic plus shared art.
- Track hubs and STEP 01 pages should link back to LEVEL 01 so genre paths feel like continuations, not a second beginner track.

## HTML page anatomy

Playable lesson pages should retain this order unless a lesson has a clear reason to differ:

1. shared navigation, bilingual link, breadcrumb, and step count;
2. hero with difficulty and `PLAYABLE` status;
3. playable WASM iframe with concise keyboard/touch instructions;
4. `DEEP DIVE` introduction;
5. ordered `concept-row`;
6. concept-matched `motion-lab` using `web/learn.js`;
7. representative `code-lesson` taken from the real implementation;
8. `why-grid`, including a learner challenge;
9. previous/next pager;
10. generated feedback section near the end of `<main>`;
11. shared footer and `learn.js`.

Core LEVEL 01–12 additionally include `code-focus`, full embedded `main.go`, and a copy button. VFX pages follow the same explanatory skeleton but are generated and may label themselves `INTERACTIVE` instead of `PLAYABLE`.

All URLs must remain relative to the current page depth. Test project-site hosting mentally (`https://USER.github.io/REPO/`), not only domain-root hosting.

## Implementation conventions

- Keep each mini game self-contained in one package unless shared code clearly improves several games.
- Use a logical canvas of `480 × 720` unless the mechanic strongly benefits from another size. The shared loader scales it responsively.
- Prefer deterministic `rand.New(rand.NewSource(...))` seeds for lessons so behavior is reproducible.
- Support retry without reloading the page.
- Use simple Ebitengine/vector shapes when art is unnecessary. If assets are introduced, document their source and license.
- Keep update rules explicit and teachable. Small structs, slices, state enums, and data tables are preferred over clever abstractions.
- Run `gofmt` on every changed Go file.

## Useful commands

The project currently requires Go 1.24 or later.

```sh
go mod download
gofmt -w games/tracks/<track>/<slug>/main.go
GOOS=js GOARCH=wasm go test ./games/tracks/<track>/<slug>
bash scripts/build.sh
bash scripts/ralph-loop.sh verify
python3 -m http.server 8080 --directory dist
```

Open lesson pages through HTTP, not `file://`, because browsers will not load WASM correctly from local files.

## Generated files and build side effects

Know the owner before editing:

- `dist/` is disposable build output and is ignored. Never patch it.
- VFX hub/lesson HTML is owned by the content table in `scripts/gen-visual-effects.mjs`. Regenerate instead of hand-editing those pages.
- Core full-source slots are owned by `data-embed-source` plus `scripts/embed-lesson-sources.mjs`. Edit the Go source and surrounding explanation, not the generated code inside `data-embed-slot`.
- Feedback markup is owned by `scripts/insert-feedback-form.mjs`. Do not hand-edit one page's copied form; change the injector when site-wide behavior or wording changes.
- Character and effect PNG/JSON assets are owned by `cmd/gen-atlas`, `cmd/gen-vfx`, and their layout packages. Regenerate them; do not paint over generated binaries.

`scripts/build.sh` is not read-only with respect to source: before creating `dist/`, it refreshes embedded source blocks and feedback sections inside `web/`. Expect a build to modify tracked HTML when generated content was stale. Inspect those diffs; do not blindly discard or stage them.

When adding a new hand-authored page, it is acceptable to omit copied feedback markup initially, but run the injector or full build before declaring it done. Keep a valid `</main>` so injection has an unambiguous location.

## Feedback workflow

- Every Japanese and English content page has a shared, site-styled feedback form at the end of the page. It posts to the Google Form `formResponse` endpoint; do not replace it with a visible iframe when adding or restructuring pages.
- The feedback form is progressive enhancement: native form fields submit without JavaScript, while pages that load `web/learn.js` get localized sending/sent/failed status without navigation. Playable lesson pages already require `learn.js`; do not remove it.
- `scripts/insert-feedback-form.mjs` is the markup and endpoint source of truth. It is intentionally idempotent and derives a hidden page path so feedback can be traced to the lesson. Change the generator once rather than editing hundreds of copies.
- Feedback responses are stored in the linked Google Sheet `1r6jYssPE7AdluEqJ1nqyzzRWrOH8Ncp-zUyK4xnn0lw`.
- The local triage tool uses OAuth. The OAuth client JSON belongs under `.secrets/` and must never be committed or pasted into chat. The first run opens a browser for Google consent and stores a refresh token in `.secrets/feedback-token.json`.
- Use these commands from the repository root:

```sh
node scripts/feedback-sheet.mjs list      # show responses and row numbers
node scripts/feedback-sheet.mjs check 12  # mark row 12 with ✅
node scripts/feedback-sheet.mjs delete 12 # delete row 12 permanently
```

- Run `check` before `delete` when the response should remain as a record. Row 1 is the header; use the row number printed by `list`.

Do not expose, commit, print, or paste OAuth client files, refresh tokens, or sheet contents. Listing feedback is a read of external user-provided content; do it only when the task asks for feedback triage.

## Verification ladder

Use proportionate checks, in this order, and do not treat an earlier layer as proof of a later one:

1. `gofmt` the changed Go file and run its targeted `GOOS=js GOARCH=wasm go test`.
2. Run `git diff --check` on the assigned files.
3. Run `bash scripts/ralph-loop.sh verify` once the working batch is integrated. This rebuilds every WASM target, injects generated content, and checks the bilingual article skeleton.
4. Inspect `bash scripts/ralph-loop.sh status` and `next`; remember VFX has its own count.
5. Serve `dist/` over HTTP and use a real browser. For each new lesson, prove one Canvas loads, exercise representative controls, inspect console errors created during that run, and verify a 390×844 viewport has no horizontal overflow and the iframe remains usable.

The structural checker proves required markup exists; it does not prove the game is distinct, the article is correct, controls work, or the page looks good. Review those directly.

## Parallel-agent workflow

Parallel work is encouraged for independent lesson directories, but all agents share one working tree.

- The coordinating agent assigns exact game and JA/EN page paths. A worker edits only those paths and runs targeted tests; it does not run whole-repository generators, builds, commits, or browser sessions unless explicitly assigned.
- Never assign work by numeric `order` alone. Assign the stable curriculum ID and slug because generated groups can shift order numbers.
- Workers must report changed files, implemented mechanic, control coverage, and exact test results.
- The coordinator reviews the actual diff, then runs shared generators/full verification once no worker is writing affected files.
- Do not run `scripts/build.sh`, `gen-visual-effects`, feedback injection, global formatting, or broad mechanical rewrites while other agents are editing overlapping HTML.
- Preserve changes outside the assignment even if they look unrelated or incomplete. Escalate conflicts instead of resetting them.

## Git hygiene

- `git add .` is expected to be safe because generated `dist/` and local/editor files are ignored. Check `git status` before committing anyway.
- Commit source games, bilingual pages, scripts, and licensed assets together when they form one coherent curriculum milestone.
- Do not push unless the user explicitly asks. The normal handoff is a verified commit that the repository owner can push.
- Do not rewrite, reset, or discard another contributor's changes.
- Before `git add .`, inspect the complete status. It is intended to be safe because generated `dist/`, credentials, and local tooling outputs are ignored, but build-time updates to tracked `web/` files are real source changes and must be understood.

## Completion standard

Do not claim the showcase is complete from file counts alone. `105/105` must also mean every listed page embeds the correct built WASM, each game is actually playable, the named technical topic is explained in both languages, and desktop/mobile browser checks have passed.
