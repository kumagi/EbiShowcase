#!/usr/bin/env node
/**
 * Resolve every clickable learning card on the bilingual home pages to the
 * WASM game it represents, then inject a shared WebP preview into both homes.
 *
 * Core and VFX cards resolve to their own lesson. Genre cards resolve to the
 * final playable lesson in that track, so their image shows the integrated
 * game learners are working toward.
 *
 * Usage:
 *   node scripts/home-thumbnails.mjs list
 *   node scripts/home-thumbnails.mjs inject
 */
import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const assetsDir = join(root, "web/assets/home-thumbnails");
const cardPattern = /<a class="([^"]+)" href="([^"]+)">/g;

function cardKind(classes) {
  const names = classes.split(/\s+/);
  if (names.includes("vfx-course-card")) return "vfx";
  if (names.includes("track-card")) return "track";
  if (names.includes("course-card")) return "core";
  return "";
}

function pageFile(lang, route) {
  return join(root, "web", lang, route, "index.html");
}

function playSlug(html) {
  return html.match(/<iframe[^>]+src="[^"]*\/play\/([^/]+)\//)?.[1] || "";
}

function resolvePlayable(route) {
  if (route === "games/flappy") return "flappy";
  const file = pageFile("ja", route);
  if (!existsSync(file)) throw new Error(`Home card target is missing: ${route}`);
  const html = readFileSync(file, "utf8");
  const direct = playSlug(html);
  if (direct) return direct;

  const steps = [...html.matchAll(/class="path-step" href="([^/]+)\//g)].map((m) => m[1]);
  for (const step of steps.reverse()) {
    const child = join(route, step);
    const childFile = pageFile("ja", child);
    if (!existsSync(childFile)) continue;
    const slug = playSlug(readFileSync(childFile, "utf8"));
    if (slug) return slug;
  }
  throw new Error(`No playable completion target found for home card: ${route}`);
}

export function collectHomeThumbnails() {
  const home = readFileSync(join(root, "web/ja/index.html"), "utf8");
  const seen = new Set();
  const items = [];
  for (const match of home.matchAll(cardPattern)) {
    const [, classes, rawHref] = match;
    const kind = cardKind(classes);
    if (!kind || seen.has(rawHref)) continue;
    seen.add(rawHref);
    const route = rawHref.replace(/^\.\//, "").replace(/\/$/, "");
    const file = `${route.replaceAll("/", "--")}.webp`;
    items.push({ kind, href: rawHref, route, slug: resolvePlayable(route), file });
  }
  return items;
}

function inject(lang, items) {
  const file = join(root, "web", lang, "index.html");
  let html = readFileSync(file, "utf8");
  const titles = new Map();
  for (const item of items) {
    const hrefAt = html.indexOf(`href="${item.href}"`);
    const cardEnd = html.indexOf("</a>", hrefAt);
    const titleStart = html.indexOf("<h3>", hrefAt);
    const titleEnd = html.indexOf("</h3>", titleStart);
    if (hrefAt < 0 || cardEnd < 0 || titleStart < 0 || titleStart > cardEnd || titleEnd > cardEnd) {
      throw new Error(`No thumbnail title for ${lang} home card: ${item.href}`);
    }
    const rawTitle = html.slice(titleStart + 4, titleEnd);
    const title = rawTitle
      .replace(/<br\s*\/?>/gi, " ")
      .replace(/<[^>]+>/g, "")
      .replace(/&amp;/g, "&")
      .replace(/\s+/g, " ")
      .trim();
    titles.set(item.href, title);
  }
  html = html.replace(/<span class="home-card-shot"><img[^>]*>(?:<span class="home-card-shot-label">[\s\S]*?<\/span>)?<\/span>/g, "");
  const byHref = new Map(items.map((item) => [item.href, item]));
  html = html.replace(cardPattern, (tag, classes, href) => {
    if (!cardKind(classes)) return tag;
    const item = byHref.get(href);
    if (!item) throw new Error(`No thumbnail mapping for ${lang} home card: ${href}`);
    const title = titles.get(href);
    if (!title) throw new Error(`No thumbnail title for ${lang} home card: ${href}`);
    const prefix = lang === "ja"
      ? { core: "このゲームを作る", vfx: "この絵を作る", track: "完成ゲーム" }[item.kind]
      : { core: "BUILD THIS GAME", vfx: "CREATE THIS EFFECT", track: "FINAL GAME" }[item.kind];
    return `${tag}<span class="home-card-shot"><img src="../assets/home-thumbnails/${item.file}" width="480" height="720" loading="lazy" decoding="async" alt=""><span class="home-card-shot-label"><small>${prefix}</small><strong>${title}</strong></span></span>`;
  });
  writeFileSync(file, html);
}

const command = process.argv[2] || "list";
const items = collectHomeThumbnails();
mkdirSync(assetsDir, { recursive: true });
writeFileSync(join(assetsDir, "manifest.json"), `${JSON.stringify(items, null, 2)}\n`);

if (command === "list") {
  process.stdout.write(`${JSON.stringify(items, null, 2)}\n`);
} else if (command === "inject") {
  const missing = items.filter((item) => !existsSync(join(assetsDir, item.file)));
  if (missing.length) throw new Error(`Missing ${missing.length} home thumbnail(s). Capture them before injecting.\n${missing.map((item) => `${item.slug}\t${item.file}`).join("\n")}`);
  inject("ja", items);
  inject("en", items);
  console.log(`Injected ${items.length} WebP previews into each home page.`);
} else {
  throw new Error(`Unknown command: ${command}`);
}
