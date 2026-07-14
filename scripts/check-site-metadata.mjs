#!/usr/bin/env node
/**
 * Metadata gate for the published curriculum: counts, bilingual pages, home
 * thumbnails, and final-game targets for every genre card.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Usage: node scripts/check-site-metadata.mjs [--json]
 */
import { existsSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { curriculum, gated } from "./curriculum.mjs";
import { collectHomeThumbnails } from "./home-thumbnails.mjs";

const root = new URL("..", import.meta.url).pathname;
const json = process.argv.includes("--json");
const expected = { gated: 208, vfx: 29, core: 12, track: 25, cards: 66 };
const failures = [];

function fail(message) {
  failures.push(message);
}

function pagePath(lang, route) {
  return join(root, "web", lang, route, "index.html");
}

function readPage(lang, route) {
  const file = pagePath(lang, route);
  if (!existsSync(file)) {
    fail(`${lang}/${route}: missing page`);
    return "";
  }
  return readFileSync(file, "utf8");
}

function isWebP(file) {
  if (!existsSync(file)) return false;
  const bytes = readFileSync(file);
  return bytes.length >= 12 && bytes.subarray(0, 4).toString("ascii") === "RIFF" && bytes.subarray(8, 12).toString("ascii") === "WEBP";
}

const vfx = curriculum.filter((entry) => entry.group === "vfx");
const tracks = curriculum.filter((entry) => entry.group === "track");
if (gated.length !== expected.gated) fail(`gated count: expected ${expected.gated}, got ${gated.length}`);
if (vfx.length !== expected.vfx) fail(`VFX count: expected ${expected.vfx}, got ${vfx.length}`);
if (new Set(tracks.map((entry) => entry.id.split("/")[1])).size !== expected.track) fail(`track count: expected ${expected.track}, got ${new Set(tracks.map((entry) => entry.id.split("/")[1])).size}`);

for (const entry of curriculum) {
  if (!entry.playable) fail(`${entry.id}: missing Go implementation at ${entry.source}`);
  for (const lang of ["ja", "en"]) readPage(lang, entry.route);
}

const cards = collectHomeThumbnails();
const byKind = Object.groupBy(cards, ({ kind }) => kind);
for (const [kind, count] of Object.entries({ core: expected.core, vfx: expected.vfx, track: expected.track })) {
  if ((byKind[kind] || []).length !== count) fail(`${kind} home-card count: expected ${count}, got ${(byKind[kind] || []).length}`);
}
if (cards.length !== expected.cards) fail(`home-card total: expected ${expected.cards}, got ${cards.length}`);

const manifestPath = join(root, "web/assets/home-thumbnails/manifest.json");
if (!existsSync(manifestPath)) {
  fail("home thumbnail manifest is missing");
} else {
  const manifest = JSON.parse(readFileSync(manifestPath, "utf8"));
  const actual = JSON.stringify(manifest.map(({ kind, href, route, slug, file }) => ({ kind, href, route, slug, file })));
  const wanted = JSON.stringify(cards.map(({ kind, href, route, slug, file }) => ({ kind, href, route, slug, file })));
  if (actual !== wanted) fail("home thumbnail manifest is stale; run node scripts/home-thumbnails.mjs inject");
}

for (const card of cards) {
  const thumbnail = join(root, "web/assets/home-thumbnails", card.file);
  if (!isWebP(thumbnail)) fail(`${card.route}: missing or invalid WebP thumbnail ${card.file}`);
  for (const lang of ["ja", "en"]) {
    const home = readPage(lang, "");
    if (!home.includes(`src="../assets/home-thumbnails/${card.file}"`)) fail(`${lang}: card ${card.route} does not embed ${card.file}`);
  }
  if (card.kind !== "track") continue;
  const [, track] = card.route.split("/");
  const source = join(root, "games", "tracks", track, card.slug, "main.go");
  if (!existsSync(source)) fail(`${card.route}: final-game source missing at games/tracks/${track}/${card.slug}/main.go`);
  const finalRoute = `tracks/${track}/${card.slug}`;
  for (const lang of ["ja", "en"]) {
    const page = readPage(lang, finalRoute);
    if (page && !page.includes(`/play/${card.slug}/`)) fail(`${lang}/${finalRoute}: final-game iframe does not target /play/${card.slug}/`);
  }
}

const summary = {
  expected,
  curriculum: { total: curriculum.length, gated: gated.length, vfx: vfx.length, playable: curriculum.filter((entry) => entry.playable).length },
  cards: cards.length,
  failures: failures.length,
  ok: failures.length === 0,
};

if (json) {
  console.log(JSON.stringify({ summary, failures }, null, 2));
} else if (failures.length) {
  console.log(`FAIL — ${failures.length} metadata problem(s):`);
  for (const message of failures) console.log(`  - ${message}`);
} else {
  console.log(`OK — ${summary.curriculum.gated}/${expected.gated} gated, ${summary.curriculum.vfx}/${expected.vfx} VFX, and ${cards.length} home cards are linked and bilingual.`);
}

process.exit(failures.length ? 1 : 0);
