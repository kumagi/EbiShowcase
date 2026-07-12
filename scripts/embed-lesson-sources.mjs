#!/usr/bin/env node
/**
 * Embed Go lesson sources into HTML placeholders.
 *
 * Fills <code data-embed-slot> inside .full-code[data-embed-source="..."].
 *
 * Usage:
 *   node scripts/embed-lesson-sources.mjs
 */
import { readFileSync, writeFileSync, existsSync, readdirSync, statSync } from "node:fs";
import { join, relative } from "node:path";

const root = new URL("..", import.meta.url).pathname;

function walkHtml(dir, out = []) {
  if (!existsSync(dir)) return out;
  for (const name of readdirSync(dir)) {
    const full = join(dir, name);
    const st = statSync(full);
    if (st.isDirectory()) walkHtml(full, out);
    else if (name.endsWith(".html")) out.push(full);
  }
  return out;
}

function escapeHtml(text) {
  return text
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

function embedInHtml(html, fileLabel) {
  const errors = [];
  let changed = false;
  let slots = 0;

  const marker = 'data-embed-source="';
  let searchFrom = 0;
  let result = "";
  let cursor = 0;

  while (true) {
    const attrAt = html.indexOf(marker, searchFrom);
    if (attrAt < 0) break;

    const openDiv = html.lastIndexOf("<div", attrAt);
    if (openDiv < 0 || openDiv < cursor) {
      errors.push(`${fileLabel}: malformed full-code around offset ${attrAt}`);
      break;
    }

    const pathStart = attrAt + marker.length;
    const pathEnd = html.indexOf('"', pathStart);
    if (pathEnd < 0) {
      errors.push(`${fileLabel}: unclosed data-embed-source`);
      break;
    }
    const sourcePath = html.slice(pathStart, pathEnd);
    slots++;

    const abs = join(root, sourcePath);
    if (!existsSync(abs)) {
      errors.push(`${fileLabel}: missing ${sourcePath}`);
      searchFrom = pathEnd + 1;
      continue;
    }
    const source = readFileSync(abs, "utf8");
    if (!source.trim()) {
      errors.push(`${fileLabel}: empty ${sourcePath}`);
      searchFrom = pathEnd + 1;
      continue;
    }

    const slotOpen = html.indexOf("data-embed-slot", pathEnd);
    if (slotOpen < 0) {
      errors.push(`${fileLabel}: no data-embed-slot for ${sourcePath}`);
      break;
    }
    const codeOpenEnd = html.indexOf(">", slotOpen);
    if (codeOpenEnd < 0) {
      errors.push(`${fileLabel}: broken code tag for ${sourcePath}`);
      break;
    }
    const codeClose = html.indexOf("</code>", codeOpenEnd);
    if (codeClose < 0) {
      errors.push(`${fileLabel}: missing </code> for ${sourcePath}`);
      break;
    }

    const escaped = escapeHtml(source.replace(/\n$/, "") + "\n");
    const old = html.slice(codeOpenEnd + 1, codeClose);
    if (old !== escaped) changed = true;
    const before = html.slice(cursor, codeOpenEnd + 1);
    result += before + escaped;
    cursor = codeClose;
    searchFrom = codeClose + 7;
  }

  if (changed) result += html.slice(cursor);
  return { html: changed ? result : html, changed, slots, errors };
}

const files = [
  ...walkHtml(join(root, "web", "ja")),
  ...walkHtml(join(root, "web", "en")),
];

let updated = 0;
let slots = 0;
const errors = [];

for (const file of files) {
  const raw = readFileSync(file, "utf8");
  if (!raw.includes("data-embed-source=")) continue;
  const label = relative(root, file);
  const { html, changed, slots: n, errors: errs } = embedInHtml(raw, label);
  slots += n;
  errors.push(...errs);
  if (changed && errs.length === 0) {
    writeFileSync(file, html);
    updated++;
  }
}

if (errors.length) {
  console.error("embed-lesson-sources failed:\n" + errors.map((e) => `  ${e}`).join("\n"));
  process.exit(1);
}

console.log(`Embedded sources into ${updated} page(s) (${slots} slot(s) scanned).`);
