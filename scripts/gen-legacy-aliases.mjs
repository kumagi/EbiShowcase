#!/usr/bin/env node
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const aliases = [
  ["games/dungeon", "games/dungeon-crawler"],
  ...["tap", "meter", "catch", "flight", "pong", "breakout", "snake", "shooter", "sokoban", "platform", "dungeon", "bullethell"]
    .map((slug) => [`tracks/visual-effects/vfx-fx-${slug}`, `tracks/visual-effects/fx-${slug}`]),
  ["tracks/fighting/frame-attack", "tracks/fighting/three-punch"],
  ["tracks/fighting/hit-reaction", "tracks/fighting/hit-reaction-dojyo"],
  ["tracks/merge-physics/merge-rule", "tracks/merge-physics/same-merge"],
  ["tracks/bomb-maze/escape-ai", "tracks/maze-chase/escape-ai"],
  ["tracks/racing/ebi-racing", "tracks/racing/ebi-circuit"],
];
for (const lang of ["ja", "en"]) {
  for (const [sourceRoute, targetRoute] of aliases) {
    const source = join(root, "web", lang, sourceRoute, "index.html");
    const target = join(root, "web", lang, targetRoute, "index.html");
    let html = readFileSync(source, "utf8")
      .replaceAll(`/${lang}/${sourceRoute}/`, `/${lang}/${targetRoute}/`)
      .replaceAll(`${lang}/${sourceRoute}/`, `${lang}/${targetRoute}/`)
      .replace(/[ \t]+$/gm, "");
    mkdirSync(dirname(target), { recursive: true });
    writeFileSync(target, html);
  }
}
console.log(`Generated ${aliases.length} full legacy aliases in JA/EN.`);
