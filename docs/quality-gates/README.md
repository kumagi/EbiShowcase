# Quality gates

This directory is the **single catalog** of quality expectations shared by:

- humans and coding agents (`AGENTS.md` points here)
- deterministic scripts (`scripts/check-quality-gates.mjs`)
- LLM review lenses (`ai_feedback_crawler.py --lens …`)

## Files

| File | Role |
| --- | --- |
| `catalog.json` | Machine-readable gates (`id`, family, severity, check type, applies_to, prompt_hint) |
| `README.md` | This overview |

Detailed checklists such as `docs/AUTHORING_CHECKLIST.md` and
`docs/ADVANCED_QUALITY_CHECKLIST.md` remain the human-readable expansions.
When you invent a new “documents should …” rule, **add a gate id first**, then
teach it in prose.

## Three meters

1. **PLAYABLE** — demo launches (208 gated + VFX side count)
2. **ADVANCED_QUALITY** — finished genre polish / replay / mobile
3. **AUTHORING** — learner can edit a RULE and verify it

A page can pass playable while still failing authoring gates.

## Check types

| `check` | Who runs it |
| --- | --- |
| `deterministic` | `node scripts/check-quality-gates.mjs` |
| `llm` | Local/remote review model with a small lens from the catalog |
| `human` | Browser / judgment (Desktop · Tablet · Phone) |

## Commands

```sh
# List every gate
node scripts/check-quality-gates.mjs --list

# Run deterministic checks (fail severity exits 1; warn prints only)
node scripts/check-quality-gates.mjs

# Structure + site only (safe for ralph verify)
node scripts/check-quality-gates.mjs --family structure,site,brand

# Authoring progress (mostly warn until the Authoring Pass finishes)
node scripts/check-quality-gates.mjs --family authoring,loop,pedagogy --sample 24

# Emit LLM lenses for the feedback crawler
node scripts/check-quality-gates.mjs --lenses loop,authoring --json
```

`ai_feedback_crawler.py` chooses two random `llm` gates when `--lens` is
omitted. Supplying `--lens loop,authoring` fixes the review to those families.

Wire-up:

- `bash scripts/ralph-loop.sh verify` runs the structure/site/brand subset.
- Full curriculum skeleton remains in `scripts/check-lessons.mjs` and
  `scripts/check-site-metadata.mjs` (deeper checks this catalog summarises).

## Adding a gate

1. Append an object to `catalog.json` with a stable `id` (`family.topic`).
2. If it is scriptable, implement it in `check-quality-gates.mjs`.
3. If it needs judgment, set `"check": "llm"` and write a concrete `prompt_hint`.
4. Link `sources` to the prose doc that explains *why*.
