#!/usr/bin/env node
import {
  existsSync,
  readFileSync,
  readdirSync,
  statSync,
} from "node:fs";
import { join, relative } from "node:path";
import { curriculum } from "./curriculum.mjs";
import { genrePresentationMap, validateGenrePresentationMap } from "./genre-presentation-map.mjs";

const root = new URL("..", import.meta.url).pathname;
const mapped = validateGenrePresentationMap();
const trackEntries = curriculum.filter((item) => item.group === "track");
const trackIDs = new Set(trackEntries.map((item) => item.route.split("/")[1]));

for (const id of trackIDs) {
  if (!mapped.has(id)) throw new Error(`genre presentation map misses curriculum track: ${id}`);
}
for (const id of mapped) {
  if (!trackIDs.has(id)) throw new Error(`genre presentation map has unknown curriculum track: ${id}`);
}

function collectLocalSource(file, seen = new Set()) {
  if (seen.has(file) || !existsSync(file)) return "";
  seen.add(file);
  const source = readFileSync(file, "utf8");
  let combined = `\n// source: ${relative(root, file)}\n${source}`;
  for (const match of source.matchAll(/"github\.com\/kumagi\/EbiShowcase\/(internal\/[^"]+)"/g)) {
    const imported = join(root, match[1]);
    if (!existsSync(imported) || !statSync(imported).isDirectory()) continue;
    for (const name of readdirSync(imported)) {
      if (!name.endsWith(".go") || name.endsWith("_test.go")) continue;
      combined += collectLocalSource(join(imported, name), seen);
    }
  }
  return combined;
}

const capstones = new Map();
for (const entry of trackEntries) capstones.set(entry.route.split("/")[1], entry);

const signals = {
  clock: /\b(frame|frames|timer|tick|phase|anim\w*|tween\w*|reaction|transition|pulse|flash|shake|life)\b/gi,
  motion: /Lerp|Ease|Sin\(|Cos\(|GeoM\.(Rotate|Scale)|Offset\(|bob|squash|interpol/gi,
  feedback: /particle|spark|burst|flash|shake|ring|trail|debris|confetti|shockwave/gi,
  replay: /\b(best|score|grade|replay|restart|retry)\b/gi,
};

function count(source, pattern) {
  return (source.match(pattern) || []).length;
}

function drawBodies(source) {
  return [...source.matchAll(/func \([^)]*\) Draw\([^]*?\n\}/g)].map((match) => match[0]);
}

const rows = [];
for (const lesson of trackEntries) {
  const source = collectLocalSource(lesson.source);
  const motion = count(source, signals.motion);
  const feedback = count(source, signals.feedback);
  if (motion+feedback === 0) {
    throw new Error(`${lesson.route}: playable lesson has no explicit motion or feedback architecture`);
  }
}

for (const definition of genrePresentationMap) {
  const capstone = capstones.get(definition.id);
  const source = collectLocalSource(capstone.source);
  const row = {
    id: definition.id,
    clock: count(source, signals.clock),
    motion: count(source, signals.motion),
    feedback: count(source, signals.feedback),
    replay: count(source, signals.replay),
  };
  if (row.clock < 8 || row.motion < 4 || row.feedback < 8) {
    throw new Error(`${definition.id} capstone lacks presentation evidence: ${JSON.stringify(row)}`);
  }
  const drawMutations = drawBodies(source).flatMap((body) =>
    [...body.matchAll(/\bg\.[A-Za-z_]\w*\s*(?:\+\+|--|[+\-*/]?=(?!=))/g)].map((match) => match[0]));
  if (drawMutations.length) {
    throw new Error(`${definition.id}: Draw mutates retained game state: ${drawMutations.slice(0, 4).join(", ")}`);
  }
  rows.push(row);
}

let pages = 0;
for (const lang of ["ja", "en"]) {
  for (const definition of genrePresentationMap) {
    const dir = join(root, "web", lang, "tracks", definition.id);
    const stack = [dir];
    while (stack.length) {
      const current = stack.pop();
      for (const name of readdirSync(current)) {
        const path = join(current, name);
        if (statSync(path).isDirectory()) stack.push(path);
        else if (name === "index.html") {
          const html = readFileSync(path, "utf8");
          const markers = (html.match(/<!-- genre-presentation-bridge:start -->/g) || []).length;
          if (markers !== 1) throw new Error(`${relative(root, path)}: expected one genre presentation bridge, got ${markers}`);
          pages++;
        }
      }
    }
  }
}

console.log(`Genre presentation gate: ${trackEntries.length} playable lessons + ${rows.length} capstones + ${pages} bilingual pages.`);
