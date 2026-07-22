# Quality gates

This directory is the **single catalog** of quality expectations shared by:

- humans and coding agents (`AGENTS.md` points here)
- deterministic scripts (`scripts/check-quality-gates.mjs`)
- LLM review lenses (`ai_feedback_crawler.py --lens …`)

## Files

| File | Role |
| --- | --- |
| `catalog.json` | Machine-readable gates and shared LLM policy (`fail_when`, `do_not_flag`, `evidence_required`) |
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

`ai_feedback_crawler.py` always audits one gate per run and, by default, asks
for three independent votes per page. It accepts a finding only when a majority
both reports an actionable verdict and quotes the same page evidence. With no `--lens`, it
chooses from all LLM gates. A family such as `--lens pedagogy` narrows that
random choice; an exact id such as `--lens pedagogy.code-matches-impl` fixes it.

Cheap local models receive the catalog's shared `llm_review_policy` plus the
structured conditions for each selected gate. A fail/warn requires a literal
quote from the supplied page; missing or ambiguous evidence is a pass. Keep
`do_not_flag` concrete because it is the main defense against repetitive false
positives from generic Update/Draw text and inferred repository structure.

Wire-up:

- `bash scripts/ralph-loop.sh verify` runs the structure/site/brand subset.
- Full curriculum skeleton remains in `scripts/check-lessons.mjs` and
  `scripts/check-site-metadata.mjs` (deeper checks this catalog summarises).

## Adding a gate

1. Append an object to `catalog.json` with a stable `id` (`family.topic`).
2. If it is scriptable, implement it in `check-quality-gates.mjs`.
3. If it needs judgment, set `"check": "llm"` and write concrete `fail_when`,
   `do_not_flag`, and `evidence_required` fields. Keep `prompt_hint` as a short overview.
4. Link `sources` to the prose doc that explains *why*.
