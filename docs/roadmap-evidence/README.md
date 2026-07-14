# Roadmap evidence

Evidence files are created by:

```sh
node scripts/roadmap-ralph-loop.mjs evidence TASK-ID
```

Keep screenshots under `docs/roadmap-evidence/assets/` and link them from the
task file. Do not check a task until every required box is checked and the
commands section records the exact commands and outcomes used for verification.

Start from [`TEMPLATE.md`](TEMPLATE.md) when writing an evidence file by hand.
The loop's `evidence TASK-ID` command creates the same required-box structure;
then add the viewport table, language/accessibility findings, and durable
screenshot or recording links from the template. A documentation-only task must
explicitly say why a gameplay viewport/input check is not applicable.

For playable work, record all three representative surfaces: Desktop
(1440×900, keyboard + pointer), Tablet (768×1024, touch), and Phone (390×844,
touch). Record Japanese and English separately. "Build succeeded" is not a
substitute for a completed control or replay loop.

Integration tasks (`P2-G01` through `P2-G25`) must prove Shader, Audio, Text/UI,
Camera, Desktop, Tablet, Phone, Japanese, English, and Tests. Phase 1 pass tasks
must prove Stages, Animation, Feedback, Replay, Keyboard, Pointer, Touch,
Japanese, English, and Tests.
