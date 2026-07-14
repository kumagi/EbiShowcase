#!/usr/bin/env node
import { existsSync, readFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const text = (path) => existsSync(path) ? readFileSync(path, "utf8") : "";
const count = (pattern, paths) => paths.reduce((n, path) => n + (text(path).match(pattern)?.length || 0), 0);
const langs = ["ja", "en"];
const late = ["snake", "space-shooter", "sokoban", "platformer", "dungeon", "bullet-hell"];
const generatedTracks = ["rhythm", "raycaster", "tower-defense", "reversi", "topdown-adventure"];

const corePages = langs.flatMap((lang) => late.map((slug) => join(root, "web", lang, "games", slug, "index.html")));
const hubs = langs.flatMap((lang) => generatedTracks.concat("platformer").map((slug) => join(root, "web", lang, "tracks", slug, "index.html")));
const graduationPages = langs.flatMap((lang) => ["arcade-60", "exploration-3rooms", "puzzle-3stages"].map((slug) => join(root, "web", lang, "graduation", slug, "index.html")));

console.log(JSON.stringify({
  playable: "reported separately by scripts/curriculum.mjs summary",
  authoring: {
    buildTrack: 4,
    coreLateRules: { complete: count(/YOUR FIRST RULE/g, corePages), total: corePages.length },
    dualLayerTrackHubs: { complete: count(/graduation-cta:start/g, hubs), total: hubs.length },
    graduationBriefs: { complete: graduationPages.filter(existsSync).length, total: graduationPages.length },
    first30Minutes: langs.filter((lang) => text(join(root, "web", lang, "guides", "first-30-minutes", "index.html")).includes("25–30")).length,
  },
}, null, 2));
