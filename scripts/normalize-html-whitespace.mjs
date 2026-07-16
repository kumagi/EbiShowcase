#!/usr/bin/env node
/**
 * Collapse accidental runs of blank lines in authored/generated HTML.
 *
 * Several idempotent injectors remove a marked block and append its replacement
 * near </main>.  The block is replaced correctly, but a surrounding newline can
 * be left behind on every build.  After enough builds this produced hundreds of
 * empty lines in otherwise small pages.
 *
 * Whitespace is significant inside pre, textarea, script, and style, so those
 * raw-text regions are deliberately left byte-for-byte unchanged.
 */
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const checkOnly = process.argv.includes("--check");
const rawTags = new Set(["pre", "textarea", "script", "style"]);

function walk(dir, files = []) {
  for (const name of readdirSync(dir)) {
    const path = join(dir, name);
    const stat = statSync(path);
    if (stat.isDirectory()) walk(path, files);
    else if (name === "index.html" || name.endsWith(".html")) files.push(path);
  }
  return files;
}

function updateRawDepth(line, depth) {
  const tags = line.matchAll(/<\/?(pre|textarea|script|style)\b[^>]*>/gi);
  for (const match of tags) {
    const tag = match[1].toLowerCase();
    if (!rawTags.has(tag)) continue;
    if (match[0].startsWith("</")) depth = Math.max(0, depth - 1);
    else if (!match[0].endsWith("/>")) depth++;
  }
  return depth;
}

export function normalizeHTML(html) {
  const trailingNewline = html.endsWith("\n");
  const lines = html.split(/\r?\n/);
  const output = [];
  let rawDepth = 0;
  let outsideBlank = false;

  for (const line of lines) {
    if (rawDepth === 0 && /^\s*$/.test(line)) {
      if (!outsideBlank) output.push("");
      outsideBlank = true;
      continue;
    }

    output.push(line);
    outsideBlank = false;
    rawDepth = updateRawDepth(line, rawDepth);
  }

  let normalized = output.join("\n");
  if (trailingNewline && !normalized.endsWith("\n")) normalized += "\n";
  return normalized;
}

const changed = [];
for (const file of walk(join(root, "web"))) {
  const before = readFileSync(file, "utf8");
  const after = normalizeHTML(before);
  if (after === before) continue;
  changed.push(file.slice(root.length + 1));
  if (!checkOnly) writeFileSync(file, after);
}

if (checkOnly && changed.length > 0) {
  console.error(`HTML whitespace check failed: ${changed.length} file(s) contain repeated blank lines.`);
  for (const file of changed.slice(0, 20)) console.error(`- ${file}`);
  if (changed.length > 20) console.error(`- ...and ${changed.length - 20} more`);
  process.exit(1);
}

console.log(
  checkOnly
    ? "HTML whitespace check passed."
    : `Normalized repeated HTML blank lines in ${changed.length} file(s).`,
);
