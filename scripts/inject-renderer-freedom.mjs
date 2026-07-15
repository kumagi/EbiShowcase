#!/usr/bin/env node
/** Add the renderer-independence lesson beside the loop contract. */
import { readdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join, relative, resolve, sep } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(join(fileURLToPath(import.meta.url), "..", ".."));
const marker = /<!-- renderer-freedom:start -->[\s\S]*?<!-- renderer-freedom:end -->\n?/g;

function walk(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const file = join(dir, entry.name);
    return entry.isDirectory() ? walk(file) : entry.name === "index.html" ? [file] : [];
  });
}

function panel(ja, file) {
  const playURL = relative(dirname(file), join(root, "web", "play", "renderer-freedom")).split(sep).join("/") + "/";
  return `<!-- renderer-freedom:start -->
<section class="renderer-freedom" aria-label="${ja ? "Drawの独立性" : "Renderer independence"}">
  <p class="eyebrow">${ja ? "ひとつのゲーム、自由な見せ方" : "ONE GAME, OPEN-ENDED PRESENTATION"}</p>
  <h2>${ja ? "描き方を替えても、ゲームのルールは替わらない。" : "Change the picture, not the game."}</h2>
  <p>${ja ? "Updateと<code>game</code>構造体がゲームの事実を司ります。下の並んだ画面は、同じ自動プレイの位置・箱・ゴールを同時に見ています。違うのはDrawだけで、このデモ以外の見せ方にも自由に差し替えられます。" : "Update and the <code>game</code> struct own the facts of the game. The views below observe the same autoplaying position, box, and goal at the same time. Only Draw differs, and these views can be replaced by any other presentation."}</p>
  <div class="renderer-freedom-demo"><iframe data-shared-demo="renderer-freedom" title="${ja ? "同じゲーム状態を異なるDrawで表示する自動プレイデモ" : "Autoplay demo showing one game state through different Draw functions"}" src="${playURL}" loading="lazy" allow="autoplay; fullscreen"></iframe></div>
</section>
<!-- renderer-freedom:end -->`;
}

let changed = 0;
for (const lang of ["ja", "en"]) {
  for (const file of walk(join(root, "web", lang))) {
    const html = readFileSync(file, "utf8");
    // Removing a generated panel must not leave indentation-only lines that
    // fail git diff checks after a full rebuild.
    const clean = html.replace(marker, "").replace(/^[ \t]+$/gm, "");
    const route = relative(join(root, "web", lang), file).split(sep).join("/");
    if (route !== "games/tap-target/index.html") {
      if (clean !== html) { writeFileSync(file, clean); changed++; }
      continue;
    }
    const anchor = '<section class="physics" id="learn">';
    const next = clean.replace(anchor, `${panel(lang === "ja", file)}\n${anchor}`);
    if (next !== html) { writeFileSync(file, next); changed++; }
  }
}
console.log(`Injected renderer-independence panels into ${changed} pages.`);
