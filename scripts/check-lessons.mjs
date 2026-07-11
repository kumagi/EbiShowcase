#!/usr/bin/env node
/**
 * Structural quality gate for playable lesson articles.
 * Checks JA+EN pages for the Flappy-class skeleton:
 * DEEP DIVE, concept-row, motion-lab, code-lesson, why-grid, learn.js, play iframe.
 *
 * Usage:
 *   node scripts/check-lessons.mjs
 *   node scripts/check-lessons.mjs --max-order 40
 *   node scripts/check-lessons.mjs --json
 */
import { readFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import { curriculum } from "./curriculum.mjs";

const root = new URL("..", import.meta.url).pathname;
const args = process.argv.slice(2);
const json = args.includes("--json");
const maxIdx = args.indexOf("--max-order");
const maxOrder = maxIdx >= 0 ? Number(args[maxIdx + 1]) : Infinity;

const required = [
  { name: "DEEP DIVE", re: /DEEP DIVE/ },
  { name: "concept-row", re: /class="concept-row"/ },
  { name: "motion-lab", re: /class="motion-lab"|data-lab=/ },
  { name: "code-lesson", re: /class="code-lesson"/ },
  { name: "why-grid", re: /class="why-grid"/ },
  { name: "learn.js", re: /learn\.js/ },
];

function checkPage(html, slug) {
  const missing = [];
  for (const rule of required) {
    if (!rule.re.test(html)) missing.push(rule.name);
  }
  const playOk =
    html.includes(`play/${slug}/`) ||
    html.includes("game.html") || // core flappy
    html.includes(`play/${slug}"`);
  if (!playOk) missing.push(`iframe play/${slug}`);
  return missing;
}

const failures = [];
let checked = 0;

for (const entry of curriculum) {
  if (!entry.playable) continue;
  if (entry.order > maxOrder) continue;
  checked++;
  for (const lang of ["ja", "en"]) {
    const file = join(root, "web", lang, entry.route, "index.html");
    if (!existsSync(file)) {
      failures.push({ order: entry.order, id: entry.id, lang, missing: ["FILE"] });
      continue;
    }
    const html = readFileSync(file, "utf8");
    const missing = checkPage(html, entry.slug);
    if (missing.length) {
      failures.push({ order: entry.order, id: entry.id, lang, missing });
    }
  }
}

const summary = {
  checked,
  pages: checked * 2,
  failures: failures.length,
  ok: failures.length === 0,
};

if (json) {
  console.log(JSON.stringify({ summary, failures }, null, 2));
} else {
  console.log(`Checked ${summary.pages} pages (${checked} playable lessons).`);
  if (failures.length === 0) {
    console.log("OK — all playable lessons meet the article skeleton.");
  } else {
    console.log(`FAIL — ${failures.length} page(s) missing required pieces:\n`);
    for (const f of failures) {
      console.log(`  #${f.order} ${f.lang} ${f.id}: ${f.missing.join(", ")}`);
    }
  }
}

process.exit(failures.length === 0 ? 0 : 1);
