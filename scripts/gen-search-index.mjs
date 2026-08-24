#!/usr/bin/env node
// SPDX-License-Identifier: Apache-2.0
/** Build the bilingual client-side search index from generated lesson pages. */
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { dirname, join, relative } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const webRoot = join(root, "web");

const GROUPS = {
  ja: {
    games: { label: "基礎編 LEVEL 01-12", order: 1 },
    build: { label: "Build Track", order: 2 },
    "tracks/visual-effects": { label: "Visual Effects Lab", order: 3 },
    tracks: { label: "ゲーム制作編", order: 4 },
    labs: { label: "ラボ", order: 5 },
    guides: { label: "ガイド", order: 6 },
    graduation: { label: "卒業制作", order: 7 },
  },
  en: {
    games: { label: "Basics LEVEL 01-12", order: 1 },
    build: { label: "Build Track", order: 2 },
    "tracks/visual-effects": { label: "Visual Effects Lab", order: 3 },
    tracks: { label: "Genre Tracks", order: 4 },
    labs: { label: "Labs", order: 5 },
    guides: { label: "Guides", order: 6 },
    graduation: { label: "Graduation", order: 7 },
  },
};

function walk(dir, files = []) {
  for (const name of readdirSync(dir).sort()) {
    if (name === "node_modules" || name === "play" || name === "assets") continue;
    const full = join(dir, name);
    const stat = statSync(full);
    if (stat.isDirectory()) walk(full, files);
    else if (name === "index.html") files.push(full);
  }
  return files;
}

const stripTags = (html) =>
  html
    .replace(/<style[\s\S]*?<\/style>/gi, "")
    .replace(/<script[\s\S]*?<\/script>/gi, "")
    .replace(/<[^>]+>/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/\s+/g, " ")
    .trim();

function classify(route, lang) {
  if (route.startsWith("games/")) return GROUPS[lang].games.label;
  if (route.startsWith("build/")) return GROUPS[lang].build.label;
  if (route.startsWith("tracks/visual-effects/")) return GROUPS[lang]["tracks/visual-effects"].label;
  if (route.startsWith("tracks/")) return GROUPS[lang].tracks.label;
  if (route.startsWith("labs/")) return GROUPS[lang].labs.label;
  if (route.startsWith("guides/")) return GROUPS[lang].guides.label;
  if (route.startsWith("graduation/")) return GROUPS[lang].graduation.label;
  return null;
}

function isHub(route) {
  // a hub is the second-level directory itself (…/<group>/<name>/) with no deeper step
  const parts = route.split("/").filter(Boolean);
  if (parts[0] === "tracks") return parts.length <= 2;
  return parts.length <= 1;
}

function extract(html, route, lang) {
  const title =
    html.match(/<meta\s+property="og:title"\s+content="([^"]+)"/i)?.[1]?.replace(/^Ebi Showcase – /, "") ||
    stripTags(html.match(/<h1[^>]*>([\s\S]*?)<\/h1>/i)?.[1] || "") ||
    route;
  const desc =
    html.match(/<meta\s+name="description"\s+content="([^"]*)"/i)?.[1] ||
    stripTags(html.match(/<meta\s+property="og:description"\s+content="([^"]*)"/i)?.[1] || "");
  const stars = (html.match(/<p class="eyebrow">([^<]*★[^<]*)<\/p>/i)?.[1] || "").trim();
  let concept = "";
  for (const pattern of [
    /class="overview-concept"[^>]*>\s*<small>[^<]*<\/small>\s*<strong>([\s\S]*?)<\/strong>/i,
    /DEEP DIVE\s*<\/p>\s*<h2[^>]*>([\s\S]*?)<\/h2>/i,
    /class="course-concept">\s*<small>[^<]*<\/small>\s*<strong>([\s\S]*?)<\/strong>/i,
  ]) {
    const hit = stripTags((html.match(pattern) || [])[1] || "");
    if (hit) {
      concept = hit;
      break;
    }
  }
  const playable = /iframe[^>]*(play\/|lesson-game-frame|data-game-src)/i.test(html);
  return { title, desc, stars, concept, playable };
}

const index = [];
for (const lang of ["ja", "en"]) {
  const langRoot = join(webRoot, lang);
  for (const file of walk(langRoot)) {
    const route = relative(langRoot, dirname(file)).replaceAll("\\", "/");
    if (route === "" || route === "search") continue; // home and the search page itself
    const group = classify(route, lang);
    if (!group) continue;
    const html = readFileSync(file, "utf8");
    if (/content=["'][^"']*noindex/i.test(html)) continue;
    const { title, desc, stars, concept, playable } = extract(html, route, lang);
    const entry = {
      path: `/${lang}/${route}/`,
      t: title,
      g: group,
      hub: isHub(route),
      p: playable,
    };
    if (desc) entry.d = desc.slice(0, 160);
    if (stars) entry.s = stars.split("/")[0].trim();
    if (concept && concept !== title) entry.c = concept.slice(0, 80);
    index.push(entry);
  }
}

index.sort((a, b) => a.path.localeCompare(b.path));
writeFileSync(join(webRoot, "assets", "search-index.json"), JSON.stringify(index));
console.log(`Search index written: ${index.length} entries (${new Set(index.map((e) => e.path.split("/")[1])).size} languages).`);
