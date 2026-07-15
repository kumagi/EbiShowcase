#!/usr/bin/env node
/** Remove wording that presents Update and Draw as an obligatory handshake. */
import { readdirSync, readFileSync, writeFileSync } from "node:fs";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(join(fileURLToPath(import.meta.url), "..", ".."));
function walk(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((e) => {
    const p = join(dir, e.name);
    return e.isDirectory() ? walk(p) : e.name === "index.html" ? [p] : [];
  });
}
const replacements = [
  ["This is the shared Update → Draw skeleton for the whole path; the numbered steps below add one system at a time.", "This is the shared game-state boundary: rules live in Update and the renderer is freely replaceable. The numbered steps below add one system at a time."],
  ["The skeleton is still the <a href=\"../tap-target/#basics\">LEVEL 01</a> <strong>Update → Draw</strong> loop.", "The state boundary still follows <a href=\"../tap-target/#basics\">LEVEL 01</a>: Update owns rules and Draw is a replaceable projection."],
  ["ゲームの骨格は <a href=\"../tap-target/#basics\">LEVEL 01</a> の <strong>Update（数字）→ Draw（絵）</strong> です。", "ゲームの状態境界は <a href=\"../tap-target/#basics\">LEVEL 01</a> と同じです。ルールはUpdate、Drawの描き方は自由です。"],
  ["ebiten.RunGame(g) // Update → Draw → repeat", "ebiten.RunGame(g) // Update ticks; Draw projects state"],
  ["Step through Update → Draw", "Step through the state boundary"],
  ["Update → Draw を1コマずつ", "UpdateとDrawの境界を確かめる"],
  ["One frame = one Update + one Draw", "The game state advances in Update; Draw is an independent projection"],
  ["This track still runs on the <a href=\"../../../games/tap-target/#basics\">LEVEL 01</a> <strong>Update → Draw</strong> loop—each step adds one new piece.", "This track shares the <a href=\"../../../games/tap-target/#basics\">LEVEL 01</a> state boundary: rules advance in Update, while Draw is a replaceable projection—each step adds one new piece."],
  ["This track still runs on the <a href=\"../../../games/tap-target/#basics\">LEVEL 01</a> <strong>Update → Draw</strong> loop—each step adds one new piece. Here the drawing and the hit test are both circles.", "This track shares the <a href=\"../../../games/tap-target/#basics\">LEVEL 01</a> state boundary: Update owns the rules and Draw is replaceable—each step adds one new piece. Here the drawing and the hit test are both circles."],
  ["Update / Draw のループの上に", "Updateが進める状態の上に、Drawは自由な投影として"],
  ["Update / Draw の上に", "Updateが進める状態の上に、Drawは自由な投影として"],
  ["LEVEL 01 の Update / Draw に、このジャンル固有の仕組みを統合します。", "LEVEL 01の状態境界に、このジャンル固有の仕組みを統合します。ルールはUpdate、見た目は自由なDrawです。"],
  ["LEVEL 01 Update / Draw</a> loop", "LEVEL 01 state boundary</a>"],
  ["LEVEL 01 の Update / Draw</a> と同じ枠の中に書きます。", "LEVEL 01の状態境界に沿って書きます。ルールはUpdate、描画は自由です。"],
  ["LEVEL 01 の Update / Draw</a> と同じ枠です。", "LEVEL 01の状態境界です。"],
  ["ここまでの「入力 → Update の RULE → Draw」の型", "ここまでの「入力をUpdateでルールへ変換し、Drawは状態を投影する」型"],
  ["the \"input → Update → Draw\" pattern", "the pattern where Update decides state and Draw projects it"],
  ["Update / Draw</a> loop", "Update state boundary</a>"],
  ["Update / Draw</a> のループ", "Updateが進める状態</a>"],
];
let changed = 0;
for (const file of walk(join(root, "web"))) {
  const before = readFileSync(file, "utf8");
  let after = before;
  for (const [from, to] of replacements) after = after.replaceAll(from, to);
  if (after !== before) { writeFileSync(file, after); changed++; }
}
console.log(`Normalized Update/Draw independence wording in ${changed} pages.`);
