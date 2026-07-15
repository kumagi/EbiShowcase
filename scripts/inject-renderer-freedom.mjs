#!/usr/bin/env node
/** Add the renderer-independence lesson beside the loop contract. */
import { readdirSync, readFileSync, writeFileSync } from "node:fs";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(join(fileURLToPath(import.meta.url), "..", ".."));
const marker = /<!-- renderer-freedom:start -->[\s\S]*?<!-- renderer-freedom:end -->\n?/g;

function walk(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const file = join(dir, entry.name);
    return entry.isDirectory() ? walk(file) : entry.name === "index.html" ? [file] : [];
  });
}

function panel(ja) {
  return `<!-- renderer-freedom:start -->
<section class="renderer-freedom" aria-label="${ja ? "Drawの独立性" : "Renderer independence"}">
  <p class="eyebrow">${ja ? "同じ STATE、4つの DRAW" : "ONE STATE, FOUR DRAWS"}</p>
  <h2>${ja ? "描き方を替えても、ゲームのルールは替わらない。" : "Change the picture, not the game."}</h2>
  <p>${ja ? "Updateが作った同じ <code>game</code> を、ASCII、ワイヤーフレーム、ドット絵、美麗なスプライトの4通りで描けます。4分割画面へ同時に描いても、Update・入力・物理は一つのままです。DrawはUpdateの直後に必ず呼ばれる相棒ではなく、状態を読む差し替え可能な投影です。" : "The same <code>game</code> made by Update can be rendered as ASCII, wireframe, pixel art, or polished sprites. Draw it four ways in four panels and Update, input, and physics remain one unchanged system. Draw is a replaceable projection, not a partner that must run immediately after Update."}</p>
  <div class="renderer-freedom-grid" aria-hidden="true"><span><b>ASCII</b><code>+---+\n| @ |\n+---+</code></span><span><b>WIREFRAME</b><i>◇───◇<br>╲ ○ ╱</i></span><span><b>PIXEL ART</b><i>▪▪●▪▪<br>▪▪▪▪▪</i></span><span><b>SPRITES</b><i>✦  EBI  ✦</i></span></div>
</section>
<!-- renderer-freedom:end -->`;
}

let changed = 0;
for (const lang of ["ja", "en"]) {
  for (const file of walk(join(root, "web", lang))) {
    const html = readFileSync(file, "utf8");
    const clean = html.replace(marker, "");
    if (!clean.includes("<!-- update-draw-contract:end -->") && !clean.includes("<!-- END BEGINNER BRIDGE -->")) continue;
    const anchor = clean.includes("<!-- update-draw-contract:end -->")
      ? "<!-- update-draw-contract:end -->"
      : "<!-- END BEGINNER BRIDGE -->";
    const next = clean.replace(anchor, `${anchor}\n${panel(lang === "ja")}`);
    if (next !== html) { writeFileSync(file, next); changed++; }
  }
}
console.log(`Injected renderer-independence panels into ${changed} pages.`);
