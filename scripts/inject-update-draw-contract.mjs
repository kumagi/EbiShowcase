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
    title: "ゲーム状態を更新する責務と、画面へ投影する責務の分離",
    update: "Update が担当",
    updateText: "キーやタッチ入力を検知し、位置・HP・ゲージ・タイマー・盤面などのゲーム状態を更新します。当たり判定やキューの処理もすべて Update が担います。",
    draw: "Draw が担当",
    drawText: "確定した game 構造体の状態を読み取り、図形・文字・エフェクトを画面ピクセルとして描画します。ゲーム状態の変更は一切行いません。",
    note: "ゲームルールは Update（または Update から呼び出す純粋関数）に実装し、Draw はその結果を画面へ投影する処理に徹します。",
    updateCode: "入力検知 → ルール計算 → 状態更新",
    drawCode: "状態読取 → 画面ピクセルへ投影",
  } : {
    eyebrow: "THE UPDATE / DRAW BOUNDARY",
    title: "Separate State Updates from Screen Projections",
    update: "Update owns state mutation",
    updateText: "Detects key and touch inputs, and updates game state—positions, HP, gauges, timers, boards, and queues. Collision resolution and queue logic live here.",
    draw: "Draw owns visual projection",
    drawText: "Reads the current game state and projects figures, text, and effects into screen pixels. It never mutates game state.",
    note: "Implement gameplay rules in Update (or pure helpers called by Update). Keep Draw strictly as a projection of that state.",
    updateCode: "detect input → evaluate rules → mutate state",
    drawCode: "read state → project pixels",
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
