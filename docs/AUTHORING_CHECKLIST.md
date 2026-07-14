# Authoring checklist

`PLAYABLE` answers “can this be played now?”  `AUTHORING` answers “can a
learner make one deliberate change and prove what it did?”  They are separate
meters. A page can remain playable while it is not yet authoring-ready; do not
reduce the 208/208 playable gate to express authoring backlog.

`ADVANCED_QUALITY` is a third, independent production-quality pass. It asks
whether a finished genre game has satisfying replay, controls, feedback,
mobile behavior, and bilingual presentation. It does **not** prove that a
learner can find a source file, add a rule, or verify it. Conversely, a clear
authoring lesson may still need advanced-quality polish.

| Meter | It proves | It does not prove |
| --- | --- | --- |
| PLAYABLE | A current demo launches and can be played. | The game is polished or editable by a learner. |
| ADVANCED_QUALITY | A finished genre game meets replay/control/feedback quality. | Its learning path has a real edit → RULE → verify route. |
| AUTHORING | A learner can make and verify one deliberate rule change. | The game has every production-quality feature. |

## Authoring Pass freeze

Until the current Authoring Pass is released, do not add a genre, increase the
208/208 playable gate, or treat a new demo as a substitute for completing an
existing authoring task. Record a promising genre as a future candidate, then
finish the learner's edit → RULE → verify path in the current curriculum first.
New-genre approval happens only after the P4 release review.

## Update / Draw axiom (required vocabulary)

Every authoring-ready page must be consistent with:

1. `Update` owns all state mutation and all input.
2. `Draw` only projects the current `game` onto the screen.
3. The screen follows `game` bit-for-bit; `Draw` never writes `game`.

RULE challenges add logic on the Update side (or a pure function Update calls).
Do not set a primary challenge that mutates state inside `Draw`.

## Every authoring-ready lesson

- [ ] Restates or finger-points the Update/Draw axiom when the lesson touches the
  loop (LEVEL 01 states it fully; later pages may be brief).
- [ ] Names the exact file to open as a repository-relative path, or names the
  file the learner creates in a fresh local Go workspace.
- [ ] Names one primary **RULE** challenge: a branch, counter, state change, or
  data row. Changing a speed, color, or other constant is optional TUNING, not
  the primary task.
- [ ] Places that RULE in `Update` or a pure function `Update` calls. `Draw`
  projects the resulting `game` state only: it never reads input, changes state,
  consumes randomness, or decides a rule result.
- [ ] Names the function, method, or data table to edit and says what behavior
  should change.
- [ ] Gives one verification method: a focused `go test`, a local `go run`, or
  a visible before/after browser result.
- [ ] Shows REAL GO from the actual implementation. A thin wrapper shows the
  editable entry file first and a short, path-labelled `internal/` excerpt
  second; it never presents an invented loop as project source.
- [ ] Preserves the three game-loop axioms: Update owns input/state changes,
  Draw maps state to pixels, and the same `game` state produces the same frame.
- [ ] Keeps Japanese and English semantically aligned and uses relative links.

## A complete track

- [ ] Every step has a traceable edit entry and a RULE challenge appropriate to
  that step’s data, update order, or draw mapping.
- [ ] The three concept cards divide the explanation into data shape, update
  order, and draw mapping rather than repeating the lead paragraph.
- [ ] The hub says what one small authored rule unlocks next and links onward to
  graduation or a related authoring route.
- [ ] The final game remains keyboard, pointer, and touch playable; authoring
  work must not regress the existing playable contract.

## A graduation project

- [ ] The brief names the local files to create or download; cloning this
  repository is not a prerequisite.
- [ ] The starter deliberately contains named TODOs for input, a rule/state,
  result feedback, and a testable completion condition.
- [ ] A test initially fails for each assigned TODO and becomes green after the
  learner’s implementation.
- [ ] The article maps every TODO to its test name and says how to run it.
- [ ] A separate reference implementation is clearly labelled as a final
  comparison, not a copy-first answer.

## Evidence and release checks

- Record the edited paths, named RULE, verification command, and Japanese/
  English review in the corresponding roadmap evidence file.
- Run targeted tests first, then `git diff --check`; phase boundaries also run
  `node scripts/roadmap-ralph-loop.mjs verify --full`.
- Use `docs/ADVANCED_QUALITY_CHECKLIST.md` for replay, controls, feedback, and
  production quality. Passing it does not automatically pass this checklist.
