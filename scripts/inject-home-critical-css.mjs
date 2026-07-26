#!/usr/bin/env node
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const criticalCSS = readFileSync(join(root, "scripts/home-critical.css"), "utf8").trim();
const inlineCriticalCSS = criticalCSS
  .replace(/\/\*[\s\S]*?\*\//g, "")
  .replace(/\s+/g, " ")
  .replace(/\s*([{}:;,>])\s*/g, "$1")
  .replace(/;}/g, "}")
  .trim();
const start = "<!-- home-critical-css:start -->";
const end = "<!-- home-critical-css:end -->";
const markedBlock = new RegExp(`${start}[\\s\\S]*?${end}`);
const blockingStylesheet = '  <link rel="stylesheet" href="../style.css">';

if (!criticalCSS.includes(".catalog-hero") || !criticalCSS.includes(".catalog-character")) {
  throw new Error("home critical CSS is missing the catalog hero");
}

const block = `${start}
  <style>${inlineCriticalCSS}</style>
  <link rel="stylesheet" href="../style.css" media="print" onload="this.media='all';this.onload=null">
  <noscript><link rel="stylesheet" href="../style.css"></noscript>
  ${end}`;

for (const lang of ["ja", "en"]) {
  const path = join(root, `web/${lang}/index.html`);
  const before = readFileSync(path, "utf8");
  const after = markedBlock.test(before)
    ? before.replace(markedBlock, block)
    : before.replace(blockingStylesheet, `  ${block}`);

  if (after === before && !before.includes(start)) {
    throw new Error(`stylesheet insertion point not found: web/${lang}/index.html`);
  }
  if ((after.match(/home-critical-css:start/g) || []).length !== 1) {
    throw new Error(`critical CSS marker count is invalid: web/${lang}/index.html`);
  }
  writeFileSync(path, after);
}

console.log("Inlined non-blocking home styles in JA/EN.");
