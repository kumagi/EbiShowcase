#!/usr/bin/env node
/**
 * Inject lab-go snippets into every motion-lab (JA+EN) from try-it-go-catalog.mjs.
 * Uses a nesting-aware scanner so nested </div> do not truncate labs.
 */
import { readFileSync, writeFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import { curriculum } from "./curriculum.mjs";
import { catalog, lineCaptions } from "./try-it-go-catalog.mjs";

const root = new URL("..", import.meta.url).pathname;

function esc(s) {
  return String(s)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

function goBlock(kind) {
  const entry = catalog[kind];
  if (!entry) return "";
  const lines = entry.lines
    .map((l) => `<span data-go-line="${esc(l.id)}">${esc(l.code)}</span>`)
    .join("\n");
  return `<div class="lab-go-block">
<p class="eyebrow">TRY IT · GO</p>
<pre class="lab-go" tabindex="0" aria-label="Go snippet">${lines}</pre>
<p class="lab-go-caption" data-lab-go-caption>—</p>
</div>`;
}

function annotateButtons(body, kind, lang) {
  const entry = catalog[kind];
  if (!entry?.buttons) return body;
  const caps = lineCaptions[lang] || lineCaptions.ja;
  let html = body;
  html = html.replace(/\sdata-go-focus="[^"]*"/g, "");
  html = html.replace(/\sdata-go-caption="[^"]*"/g, "");
  html = html.replace(/\sdata-go-clear(?=[\s>])/g, "");
  // Longer attribute names first so data-lab-fx does not eat data-lab-fx-ping.
  // Consume optional ="…" so data-lab-dir="N" is not split into caption="…"="N".
  const attrs = Object.entries(entry.buttons).sort((a, b) => b[0].length - a[0].length);
  for (const [attr, spec] of attrs) {
    const re = new RegExp(
      `(<button\\b[^>]*\\b${attr}(?:="[^"]*")?)(?=[\\s>/])([^>]*>)`,
      "gi",
    );
    html = html.replace(re, (full, start, end) => {
      let attrs = start;
      if (spec.clear) {
        attrs += " data-go-clear";
        return `${attrs}${end}`;
      }
      const ids = (spec.ids || []).join(" ");
      attrs += ` data-go-focus="${ids}"`;
      if (spec.ids?.length) {
        const bits = spec.ids.map((id) => caps[id]).filter(Boolean);
        if (bits.length) attrs += ` data-go-caption="${esc(bits.join(" → "))}"`;
      }
      return `${attrs}${end}`;
    });
  }
  return html;
}

/** Find motion-lab elements with correct nesting. */
function findMotionLabs(html) {
  const labs = [];
  const startRe = /<(div|section)(\s[^>]*\bclass="[^"]*\bmotion-lab\b[^"]*"[^>]*)>/gi;
  let m;
  while ((m = startRe.exec(html))) {
    const el = m[1];
    const attrs = m[2];
    const openEnd = m.index + m[0].length;
    let i = openEnd;
    let depth = 1;
    const openTag = new RegExp(`<${el}\\b`, "gi");
    const closeTag = new RegExp(`</${el}>`, "gi");
    while (i < html.length && depth > 0) {
      openTag.lastIndex = i;
      closeTag.lastIndex = i;
      const nextOpen = openTag.exec(html);
      const nextClose = closeTag.exec(html);
      if (!nextClose) break;
      if (nextOpen && nextOpen.index < nextClose.index) {
        depth++;
        i = nextOpen.index + nextOpen[0].length;
      } else {
        depth--;
        i = nextClose.index + nextClose[0].length;
        if (depth === 0) {
          labs.push({
            start: m.index,
            openEnd,
            closeStart: nextClose.index,
            end: i,
            el,
            attrs,
            body: html.slice(openEnd, nextClose.index),
          });
        }
      }
    }
  }
  return labs;
}

function processPage(html, lang) {
  const labs = findMotionLabs(html);
  if (!labs.length) return { html, updated: 0 };
  // rewrite from end so indices stay valid
  let out = html;
  let updated = 0;
  for (let n = labs.length - 1; n >= 0; n--) {
    const lab = labs[n];
    let kind = (/data-lab="([^"]+)"/.exec(lab.attrs) || [])[1] || null;
    if (!kind && /data-lab-bird/.test(lab.body)) kind = "flappy";
    if (!kind || !catalog[kind]) continue;

    let attrs = lab.attrs;
    if (kind === "flappy" && !/data-lab=/.test(attrs)) {
      attrs = attrs.replace(
        /\bclass="([^"]*\bmotion-lab\b[^"]*)"/,
        'class="$1" data-lab="flappy"',
      );
    }

    let body = lab.body.replace(/<div class="lab-go-block">[\s\S]*?<\/div>\s*/g, "");
    body = annotateButtons(body, kind, lang);
    body = body.replace(/\s+$/, "") + "\n" + goBlock(kind) + "\n";
    const replacement = `<${lab.el}${attrs}>${body}</${lab.el}>`;
    out = out.slice(0, lab.start) + replacement + out.slice(lab.end);
    updated++;
  }
  return { html: out, updated };
}

let pages = 0;
let labs = 0;
const missing = new Set();

for (const item of curriculum) {
  for (const lang of ["ja", "en"]) {
    const path = join(root, "web", lang, item.route, "index.html");
    if (!existsSync(path)) continue;
    const raw = readFileSync(path, "utf8");
    for (const m of raw.matchAll(/data-lab="([^"]+)"/g)) {
      if (!catalog[m[1]]) missing.add(m[1]);
    }
    const { html, updated } = processPage(raw, lang);
    if (updated > 0) {
      writeFileSync(path, html);
      pages++;
      labs += updated;
    }
  }
}

console.log(`Wired ${labs} lab-go blocks across ${pages} pages.`);
if (missing.size) console.log("Missing catalog kinds:", [...missing].sort().join(", "));
else console.log("Catalog covers all data-lab kinds found.");
