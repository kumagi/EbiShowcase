#!/usr/bin/env node
/**
 * Put the non-negotiable Update/Draw boundary beside every genre lesson's
 * representative code. This is deliberately a structural guide rather than
 * invented per-game source: actual helper names remain on each code panel.
 * SPDX-License-Identifier: Apache-2.0
 */
import { readdirSync, readFileSync, writeFileSync } from "node:fs";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(join(fileURLToPath(import.meta.url), "..", ".."));

function pages(dir) {
  const result = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const file = join(dir, entry.name);
    if (entry.isDirectory()) result.push(...pages(file));
    else if (entry.name === "index.html") result.push(file);
  }
  return result;
}

function contract(ja) {
  const copy = ja ? {
    eyebrow: "UPDATE と DRAW の境界",
    title: "数字を決める場所と、絵にする場所を分けよう",
    update: "Update が担当",
    updateText: "キー・タッチを読み、位置・HP・ゲージ・タイマー・盤面を変えます。判定やキュー操作もこちらです。",
    draw: "Draw が担当",
    drawText: "できあがった game の値を読んで、絵・文字・エフェクトを置きます。game の値は変えません。",
    note: "このページのルールは Update 側（または Update が呼ぶ純粋な関数）へ置き、Draw はその結果を表示するだけにします。",
    updateCode: "入力 → ルール → game を更新",
    drawCode: "game を読む → ピクセルへ投影",
  } : {
    eyebrow: "THE UPDATE / DRAW BOUNDARY",
    title: "Decide the numbers first, then paint them",
    update: "Update owns",
    updateText: "Read keys and touch, then change positions, HP, gauges, timers, boards, and queues. Rules and decisions live here too.",
    draw: "Draw owns",
    drawText: "Read the finished game state and place pictures, text, and effects. It never changes game state.",
    note: "For this lesson, put a new rule in Update (or a pure helper called by Update). Draw only shows the result of that rule.",
    updateCode: "input → rule → change game",
    drawCode: "read game → project pixels",
  };
  return `<!-- update-draw-contract:start -->\n<section class="update-draw-contract" aria-label="${copy.eyebrow}">\n  <p class="eyebrow">${copy.eyebrow}</p>\n  <h2>${copy.title}</h2>\n  <div class="concept-row">\n    <article><span class="concept-number">1</span><h3>${copy.update}</h3><p>${copy.updateText}</p><code>${copy.updateCode}</code></article>\n    <article><span class="concept-number">2</span><h3>${copy.draw}</h3><p>${copy.drawText}</p><code>${copy.drawCode}</code></article>\n  </div>\n  <p class="update-draw-note">${copy.note}</p>\n</section>\n<!-- update-draw-contract:end -->`;
}

let changed = 0;
for (const language of ["ja", "en"]) {
  for (const file of pages(join(root, "web", language, "tracks"))) {
    const html = readFileSync(file, "utf8");
    const clean = html.replace(/<!-- update-draw-contract:start -->[\s\S]*?<!-- update-draw-contract:end -->\n?/g, "");
    if (!clean.includes("<section class=\"code-lesson\"")) continue;
    const next = clean.replace("<section class=\"code-lesson\"", `${contract(language === "ja")}\n<section class="code-lesson"`);
    if (next !== html) {
      writeFileSync(file, next);
      changed++;
    }
  }
}
console.log(`Injected Update/Draw boundary cards into ${changed} genre lesson pages.`);
