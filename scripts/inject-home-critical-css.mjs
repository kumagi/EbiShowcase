#!/usr/bin/env node
import { readFileSync, readdirSync, statSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const criticalCSS = readFileSync(join(root, "scripts/home-critical.css"), "utf8").trim();
const contentCriticalCSS = readFileSync(join(root, "scripts/content-critical.css"), "utf8").trim();
function minify(css) {
  return css
  .replace(/\/\*[\s\S]*?\*\//g, "")
  .replace(/\s+/g, " ")
  .replace(/\s*([{}:;,>])\s*/g, "$1")
  .replace(/;}/g, "}")
  .trim();
}
const inlineCriticalCSS = minify(criticalCSS);
const inlineContentCriticalCSS = minify(contentCriticalCSS);
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

if (!contentCriticalCSS.includes(".overview-hero") || !contentCriticalCSS.includes(".lesson-hero") || !contentCriticalCSS.includes(".test-hero")) {
  throw new Error("content critical CSS is missing a supported hero");
}

function walk(dir, files = []) {
  for (const name of readdirSync(dir)) {
    const file = join(dir, name);
    const stat = statSync(file);
    if (stat.isDirectory()) walk(file, files);
    else if (name.endsWith(".html")) files.push(file);
  }
  return files;
}

const contentStart = "<!-- page-critical-css:start -->";
const contentEnd = "<!-- page-critical-css:end -->";
const contentMarkedBlock = new RegExp(`${contentStart}[\\s\\S]*?${contentEnd}`);
let contentPages = 0;

for (const file of walk(join(root, "web"))) {
  if (file.endsWith("/game.html") || file.endsWith("/ja/index.html") || file.endsWith("/en/index.html") || file.endsWith("/web/index.html")) continue;
  const before = readFileSync(file, "utf8");
  const existing = before.match(contentMarkedBlock)?.[0] || "";
  const href =
    existing.match(/href="((?:\.\.\/)*style\.css)"/)?.[1] ||
    before.match(/<link\s+rel="stylesheet"\s+href="((?:\.\.\/)*style\.css)">/)?.[1];
  if (!href) continue;
  const block = `${contentStart}
  <style>${inlineContentCriticalCSS}</style>
  <link rel="stylesheet" href="${href}" media="print" onload="this.media='all';this.onload=null">
  <noscript><link rel="stylesheet" href="${href}"></noscript>
  ${contentEnd}`;
  const after = existing
    ? before.replace(contentMarkedBlock, block)
    : before.replace(
      new RegExp(`\\s*<link\\s+rel="stylesheet"\\s+href="${href.replaceAll(".", "\\.")}">`),
      `\n  ${block}`,
    );
  if (after === before && !existing) throw new Error(`stylesheet insertion point not found: ${file}`);
  writeFileSync(file, after);
  contentPages++;
}

console.log(`Inlined non-blocking critical styles in JA/EN home and ${contentPages} content pages.`);
