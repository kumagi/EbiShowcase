#!/usr/bin/env node
/**
 * Insert full-code placeholders into intro lesson articles (order 1–12).
 * Idempotent: skips pages that already have data-embed-source.
 */
import { readFileSync, writeFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import { curriculum } from "./curriculum.mjs";

const root = new URL("..", import.meta.url).pathname;

function sourcePath(entry) {
  if (entry.slug === "flappy") return "game/main.go";
  return `games/core/${entry.slug}/main.go`;
}

function block(lang, src) {
  if (lang === "ja") {
    return `
      <p class="full-code-intro">下のコードが、このゲームの<strong>全部</strong>です。1つのファイルで動いています。コピーして、手元のエディタで開けます。</p>
      <div class="full-code" data-embed-source="${src}">
        <div class="full-code-head">
          <span>main.go · ぜんぶのコード</span>
          <button type="button" class="full-code-copy" data-copy data-copied-label="コピーした！">コピー</button>
        </div>
        <pre><code data-embed-slot></code></pre>
      </div>
`;
  }
  return `
      <p class="full-code-intro">The code below is the <strong>entire</strong> game in one file. Copy it and open it in your editor.</p>
      <div class="full-code" data-embed-source="${src}">
        <div class="full-code-head">
          <span>main.go · full source</span>
          <button type="button" class="full-code-copy" data-copy data-copied-label="Copied!">Copy</button>
        </div>
        <pre><code data-embed-slot></code></pre>
      </div>
`;
}

let updated = 0;
for (const entry of curriculum) {
  if (!entry.playable || entry.order > 12) continue;
  const src = sourcePath(entry);
  for (const lang of ["ja", "en"]) {
    const file = join(root, "web", lang, entry.route, "index.html");
    if (!existsSync(file)) {
      console.error("missing", file);
      process.exit(1);
    }
    let html = readFileSync(file, "utf8");
    if (html.includes("data-embed-source=")) continue;

    const marker = '<div class="why-grid">';
    const at = html.indexOf(marker);
    if (at < 0) {
      console.error("no why-grid in", file);
      process.exit(1);
    }
    html = html.slice(0, at) + block(lang, src) + "\n    " + html.slice(at);
    writeFileSync(file, html);
    updated++;
  }
}
console.log(`Inserted full-code placeholders into ${updated} page(s).`);
