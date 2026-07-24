# Genre presentation curriculum

The Visual Effects Advanced course (A01–A12) is the shared presentation
foundation for all twenty-five genre specializations. A genre course does not
start over with “add a particle” or treat animation as optional final polish.
It starts with the following contract already in force:

1. gameplay commits one unambiguous result;
2. a value, event, plan, or snapshot crosses into presentation once;
3. presentation advances its own finite clocks;
4. `Draw` reads state and never commits rules;
5. unit tests cover the rule/result seam without a GPU.

[`scripts/genre-presentation-map.mjs`](../scripts/genre-presentation-map.mjs)
maps each specialization to two or three distinct A01–A12 patterns. It also
defines, in Japanese and English, the genre rule that is genuinely new, the
visible motion baseline, and the test seam. The build injects those contracts
into every track hub and lesson page.

## Course order

A specialization should now be read in four layers:

1. **Foundation reference** — the mapped A01–A12 patterns are assumed, not
   re-taught.
2. **Genre rule** — each step isolates one new rule such as wall kicks,
   target priority, or legal-move search.
3. **Readable action** — the playable demo presents that rule with
   anticipation, action, impact, and recovery appropriate to its genre.
4. **Test seam** — rule results remain deterministic when poses, particles,
   screen shake, and renderer quality are removed.

Small intermediate games may intentionally use simpler art than a capstone.
They may not use instantaneous teleportation, deletion, or number changes as
an excuse to omit the presentation model being referenced. The final game must
combine the mapped patterns under a bounded effect budget and preserve a
complete replay loop.

## Regression gate

`node scripts/check-genre-presentation.mjs` verifies:

- all curriculum tracks have a presentation map and no unknown map entries
  exist;
- every one of the 184 playable specialization lessons exposes an explicit
  motion transform or feedback system, so a rule-only demo cannot silently
  replace the assumed presentation foundation;
- every capstone and its local shared implementation contains explicit
  presentation clocks, motion transforms, and feedback systems;
- `Draw` does not mutate retained game state;
- every Japanese and English genre page contains exactly one generated
  foundation bridge.

The gate is structural, not a replacement for browser play. Real desktop and
mobile runs must still confirm that motion is readable, effects do not cover
the play field, touch can finish the loop, and reduced presentation quality
does not alter rules.
