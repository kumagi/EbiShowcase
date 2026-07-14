# Roadmap evidence (Authoring Pass)

Evidence files are created by:

```sh
node scripts/roadmap-ralph-loop.mjs evidence TASK-ID
```

The live checklist is [`docs/ROADMAP_RALPH_LOOP.md`](../ROADMAP_RALPH_LOOP.md)
(Authoring Pass, 59 tasks). Older `P*`-named evidence from the previous
playground/lab roadmap may remain under
`archive-playground-pass/` as history; it does **not**
count toward Authoring Pass completion. Do not copy those
files back to this folder with `Status: PASS` unless you are
genuinely finishing the matching Authoring task.

Keep screenshots under `docs/roadmap-evidence/assets/` and link them from the
task file. Do not check a task until every required box is checked and the
commands section records the exact commands and outcomes used for verification.

Start from [`TEMPLATE.md`](TEMPLATE.md) when writing an evidence file by hand.
The loop's `evidence TASK-ID` command creates the required-box structure for the
current task type (edit target, dual panel, holey starter, etc.).

## What “done” means here

Authoring evidence must show that a learner can:

1. find the **edit target** (file path / function),
2. perform the **RULE challenge** (not only constant tuning),
3. **verify** with `go test` or the page’s stated check.

“Build succeeded” or “iframe loads” is not enough. For playable Build Track /
verify tasks, also record Desktop and Phone (or the boxes the template requires)
plus Japanese and English.
