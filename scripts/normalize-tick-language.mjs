#!/usr/bin/env node
/**
 * Keep simulation cadence (Update ticks) distinct from rendered or sprite frames.
 * SPDX-License-Identifier: Apache-2.0
 */
import { readdirSync, readFileSync, writeFileSync } from "node:fs";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(join(fileURLToPath(import.meta.url), "..", ".."));

function walk(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const path = join(dir, entry.name);
    return entry.isDirectory() ? walk(path) : entry.name.endsWith(".html") ? [path] : [];
  });
}

const literal = [
  ["1フレーム進める", "1 tick進める"],
  ["1F進める", "1 tick進める"],
  ["1フレームずつ", "1 tickずつ"],
  ["1フレーム後", "1 tick後"],
  ["1フレームだけ", "1 tickだけ"],
  ["1フレームで", "1 tickで"],
  ["1フレームに", "1 tickに"],
  ["毎フレームの", "tickごとの"],
  ["毎フレームではなく", "tickごとではなく"],
  ["毎フレーム", "tickごとに"],
  ["フレームごとの", "tickごとの"],
  ["フレームごとに", "tickごとに"],
  ["フレームごと", "tickごと"],
  ["次のフレーム", "次のtick"],
  ["次フレーム", "次のtick"],
  ["前のフレーム", "前のtick"],
  ["前フレーム", "前のtick"],
  ["今フレーム", "今のtick"],
  ["同じフレーム", "同じtick"],
  ["接触フレーム", "接触tick"],
  ["命中フレーム", "命中tick"],
  ["リリースされたフレーム", "リリースを検出したtick"],
  ["フレーム更新", "tick更新"],
  ["フレーム差", "tick差"],
  ["予定フレーム", "予定tick"],
  ["押したフレーム", "押したtick"],
  ["待ちフレーム", "待ちtick"],
  ["数フレーム", "数tick"],
  ["このフレーム", "このtick"],
  ["フレーム係数", "tick係数"],
  ["個/フレーム", "個/tick"],
  ["Advance 1 frame", "Advance 1 tick"],
  ["Advance one frame", "Advance one tick"],
  ["Next frame", "Next tick"],
  ["next frame", "next tick"],
  ["previous frame", "previous tick"],
  ["Every frame", "Every tick"],
  ["every frame", "every tick"],
  ["Each frame", "Each tick"],
  ["each frame", "each tick"],
  ["once per frame", "once per tick"],
  ["per frame", "per tick"],
  ["one frame", "one tick"],
  ["across frames", "across ticks"],
  ["frame updates", "tick updates"],
  ["frame update", "tick update"],
  ["frame loop", "game loop"],
  ["frame difference", "tick difference"],
  ["chart frame", "chart tick"],
  ["pressed frame", "pressed tick"],
  ["wait frames", "wait ticks"],
  ["a few frames", "a few ticks"],
  ["few frames", "few ticks"],
  ["ten frames", "ten ticks"],
  ["first frame", "first tick"],
  ["single-frame", "single-tick"],
  ["contact frame", "contact tick"],
  ["contact frames", "contact ticks"],
  ["PLAY EACH FRAME", "RUN THE GAME"],
];

let changed = 0;
for (const file of walk(join(root, "web"))) {
  const route = file.slice(join(root, "web").length).replaceAll("\\", "/");
  // These pages talk about pictures inside a sprite sheet, where “frame” is
  // intentionally the right word rather than an Update tick.
  if (route.includes("/visual-effects/vfx-walk/") || route === "/en/index.html" || route === "/ja/index.html") continue;

  const before = readFileSync(file, "utf8");
  let after = before
    .replaceAll("Cut one frame at a time from a sheet", "Cut __SPRITE_FRAME__ at a time from a sheet")
    .replaceAll("同じ数字を毎フレーム描く", "Drawが呼ばれるたび同じ数字を描く")
    .replaceAll("Draw は毎フレームその答え", "Draw は呼ばれるたびその答え")
    .replaceAll("the same numbers are drawn every frame", "Draw presents the same numbers whenever it is called")
    .replaceAll("Draw projects that answer into pixels every frame", "Draw projects that answer into pixels whenever it is called");
  for (const [from, to] of literal) after = after.replaceAll(from, to);
  after = after
    .replaceAll("tickごとにの", "tickごとの")
    .replaceAll("tickごとにに", "tickごとに")
    .replaceAll("Cut __SPRITE_FRAME__ at a time from a sheet", "Cut one frame at a time from a sheet")
    .replace(/(\d+)フレームごと/g, "$1 tickごと")
    .replace(/(\d+)フレームに1回/g, "$1 tickに1回")
    .replace(/(\d+)フレームで1マス/g, "$1 tickで1マス")
    .replace(/(\d+)-frame(?=\s+(?:cooldown|window|delay|timer))/gi, "$1-tick")
    .replace(/(\d+) frames(?=\s+(?:so|until|before|after|for|—|-))/gi, "$1 ticks");
  if (after !== before) {
    writeFileSync(file, after);
    changed++;
  }
}

console.log(`Normalized Update cadence to tick terminology in ${changed} page(s).`);
