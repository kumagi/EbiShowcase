# AGENTS.md — Ebi Showcase contributor guide

## Purpose

Ebi Showcase is a static, bilingual learning site that teaches 2D game development through small games that run directly in the browser. Every lesson is a pair:

1. a genuinely playable mini game written in Go with Ebitengine and compiled to WebAssembly; and
2. a Japanese and English article that explains the game's main technique in language a primary-school student can follow.

The long-term completion gate is `105/105` playable curriculum entries. A lesson page or design mockup without a working Ebitengine game is not complete.

## Non-negotiable principles

- Game logic, input, simulation, collision, and rendering belong in Go + Ebitengine. HTML, CSS, and JavaScript are only for the static showcase and shared WASM loader.
- Everything in the repository is Apache License 2.0 compatible. Do not add assets, fonts, code, or copied game names/artwork with incompatible or unclear licensing.
- Prefer original simple shapes; the shared Ebi Boy sprite is allowed as the curriculum protagonist. Ebi-themed presentation otherwise. Famous games are references for mechanics, not permission to copy their art, text, characters, music, maps, or branding.
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

## Definition of done for one lesson

- A distinct, playable Ebitengine game exists at the expected path. It should demonstrate the lesson's named mechanic, not merely reskin a generic template.
- The game has a goal, failure or retry behavior where appropriate, clear feedback, and understandable controls.
- Keyboard and touch controls are both implemented. Pointer controls should use Ebitengine input APIs, not page-side JavaScript.
- `web/ja/<route>/index.html` and `web/en/<route>/index.html` embed `../../../../play/<slug>/` (or the correct relative depth).
- Both articles explain the mechanic step by step, including a small representative Go snippet. Japanese should avoid unexplained jargon and use concrete examples a child can picture.
- Core lessons (curriculum order 1–12) also embed the full `main.go` with a short “look here first” map and copy button (`data-embed-source`). Track lessons use the code-lesson snippet only; point learners to the repo for the full file.
- Shared hero art may live in `internal/hero`. Lesson copy must say the on-screen character is Ebi Boy when that sprite is drawn, and must not claim the HTML listing is a fully self-contained single file if shared packages are imported.
- The relevant hub/card says PLAYABLE or otherwise clearly identifies the finished lesson when the surrounding design supports statuses.
- Go formatting, JS/WASM compilation, and the Ralph verification pass.
- A real browser loads one Canvas without console errors. Check a desktop viewport and a phone-size viewport; the lesson page must not overflow horizontally.

## Article consistency (current curriculum policy)

- Teach Ebitengine from the game loop first: LEVEL 01 owns Update / Draw / (optional) Layout. Later lessons review or extend that frame; they do not reintroduce the engine as if from zero.
- When a lesson draws the shared hero sprite (`internal/hero`), articles in both languages must call out **Ebi Boy** and say that hit tests remain simple shapes.
- Full-source blurbs must not claim `internal/hero` if that package is unused. Shape-only games say the listing is the entire game; hero games say main logic plus shared art.
- Track hubs and STEP 01 pages should link back to LEVEL 01 so genre paths feel like continuations, not a second beginner track.

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

## Git hygiene

- `git add .` is expected to be safe because generated `dist/` and local/editor files are ignored. Check `git status` before committing anyway.
- Commit source games, bilingual pages, scripts, and licensed assets together when they form one coherent curriculum milestone.
- Do not push unless the user explicitly asks. The normal handoff is a verified commit that the repository owner can push.
- Do not rewrite, reset, or discard another contributor's changes.

## Completion standard

Do not claim the showcase is complete from file counts alone. `105/105` must also mean every listed page embeds the correct built WASM, each game is actually playable, the named technical topic is explained in both languages, and desktop/mobile browser checks have passed.
