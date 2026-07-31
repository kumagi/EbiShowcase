#!/usr/bin/env node
/**
 * Keep heavy WASM demos out of the initial rendering window and make the
 * shared lesson script non-parser-blocking.
 */
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { dirname, join, relative } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const webRoot = join(root, "web");
const start = "<!-- page-delivery:start -->";
const end = "<!-- page-delivery:end -->";
const markedBlock = new RegExp(`^[ \\t]*${start}[\\s\\S]*?^[ \\t]*${end}[ \\t]*\\r?\\n?`, "gm");

function walk(dir, files = []) {
  for (const name of readdirSync(dir)) {
    const file = join(dir, name);
    const stat = statSync(file);
    if (stat.isDirectory()) walk(file, files);
    else if (name.endsWith(".html")) files.push(file);
  }
  return files;
}

function href(from, to) {
  return relative(dirname(from), to).replaceAll("\\", "/");
}

let updated = 0;
let deferredFrames = 0;
for (const file of walk(webRoot)) {
  if (file.endsWith("/game.html")) continue;
  const before = readFileSync(file, "utf8");
  let next = before;

  next = next.replace(/<iframe\b[^>]*>/gi, (tag) => {
    if (/\bdata-game-src=/i.test(tag)) return tag;
    const src = tag.match(/\bsrc="([^"]*\/play\/[^"]*)"/i)?.[1];
    if (!src) return tag;
    deferredFrames++;
    let result = tag.replace(/\s+loading="[^"]*"/i, "");
    result = result.replace(/\s+src="[^"]*"/i, ` loading="lazy" data-game-src="${src}"`);
    return result;
  });

  next = next.replace(/<script\b([^>]*\bsrc="[^"]*learn\.js"[^>]*)>\s*<\/script>/gi, (tag, attrs) => {
    if (/\bdefer\b/i.test(attrs)) return tag;
    return `<script${attrs} defer></script>`;
  });

  if (next.includes("data-game-src=") && !next.includes(start)) {
    const bootHref = href(file, join(webRoot, "page-boot.js"));
    const block = `${start}\n  <script src="${bootHref}" defer></script>\n  ${end}`;
    const learnTag = next.match(/<script\b[^>]*\bsrc="[^"]*learn\.js"[^>]*>\s*<\/script>/i)?.[0];
    next = learnTag
      ? next.replace(learnTag, `${block}\n  ${learnTag}`)
      : next.replace(/<\/body>/i, `  ${block}\n</body>`);
  } else if (!next.includes("data-game-src=") && next.includes(start)) {
    next = next.replace(markedBlock, "");
  }

  if (next !== before) {
    writeFileSync(file, next);
    updated++;
  }
}

console.log(`Optimized initial delivery in ${updated} pages; deferred ${deferredFrames} WASM frames.`);
