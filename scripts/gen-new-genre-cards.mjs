#!/usr/bin/env node
// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
// Keep the three post-raycaster genre cards and public track count bilingual.
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const cards = {
  ja: `<a class="track-card track-rhythm" href="tracks/rhythm/"><span class="track-letter">W</span><div><p class="eyebrow">音の時間をゲームルールにする</p><h3>リズムゲーム</h3><p>拍を光らせるところから、判定窓、譜面、長押し、難度選択のあるライブへ進みます。</p><strong>PLAYABLE PATH →</strong></div></a>
<a class="track-card track-tower-defense" href="tracks/tower-defense/"><span class="track-letter">X</span><div><p class="eyebrow">道・射程・標的を考える</p><h3>タワーディフェンス</h3><p>敵の経路、射程、先頭標的、配置、強化、ボスウェーブを一つずつ組み立てます。</p><strong>PLAYABLE PATH →</strong></div></a>
<a class="track-card track-topdown-adventure" href="tracks/topdown-adventure/"><span class="track-letter">Y</span><div><p class="eyebrow">剣と鍵で部屋を切り開く</p><h3>見下ろし型アクションアドベンチャー</h3><p>8方向移動から始め、剣、無敵時間、鍵と仕掛け、多段階ボスの迷宮を完成させます。</p><strong>PLAYABLE PATH →</strong></div></a>`,
  en: `<a class="track-card track-rhythm" href="tracks/rhythm/"><span class="track-letter">W</span><div><p class="eyebrow">Turn musical time into game rules</p><h3>Rhythm Game</h3><p>Grow a flashing beat into judgement windows, charts, hold notes, difficulty, and a complete live set.</p><strong>PLAYABLE PATH →</strong></div></a>
<a class="track-card track-tower-defense" href="tracks/tower-defense/"><span class="track-letter">X</span><div><p class="eyebrow">Reason about paths, range, and targets</p><h3>Tower Defense</h3><p>Build enemy paths, range, lead targeting, placement, upgrades, and a boss wave one system at a time.</p><strong>PLAYABLE PATH →</strong></div></a>
<a class="track-card track-topdown-adventure" href="tracks/topdown-adventure/"><span class="track-letter">Y</span><div><p class="eyebrow">Open rooms with a blade and keys</p><h3>Top-down Action Adventure</h3><p>Start with eight-way movement, then add a blade, invulnerability, keys, tools, and a multi-phase dungeon boss.</p><strong>PLAYABLE PATH →</strong></div></a>`,
};

for (const lang of ["ja", "en"]) {
  const file = join(root, "web", lang, "index.html");
  let html = readFileSync(file, "utf8");
  html = html.replace(/<!-- new-genre-tracks:start -->[\s\S]*?<!-- new-genre-tracks:end -->\n?/g, "");
  const marker = "<!-- raycaster-track:end -->";
  if (!html.includes(marker)) throw new Error(`raycaster marker missing in ${lang} home`);
  html = html.replace(marker, `${marker}\n<!-- new-genre-tracks:start -->\n${cards[lang]}\n<!-- new-genre-tracks:end -->`);
  if (lang === "ja") {
    html = html.replaceAll("12の基礎ゲーム＋22の専門トラック", "12の基礎ゲーム＋25の専門トラック")
      .replaceAll("22の専門コース", "25の専門コース")
      .replaceAll("22専門トラック", "25専門トラック")
      .replace("12本の共通基礎（ループから弾幕まで）を終えたら、25の専門コースへ。どれも LEVEL 01 の Update / Draw の続きです。完成形をいきなり作らず、仕組みをひとつずつ足します。", "基礎で使った移動、当たり判定、時間、状態管理を、作りたいジャンルに合わせて組み合わせます。完成形をいきなり作らず、小さく遊べる仕組みをひとつずつ足します。")
      .replace("Update / Draw のゲームループと、タッチ座標の当たり判定を学びます。", "丸を押した場所を調べ、当たったら得点する最初のルールを学びます。")
      .replace("動くバーを狙った場所で止め、フレームと状態を理解します。", "動くバーを狙った場所で止め、時間の進み方とゲームの状態を学びます。");
  } else {
    html = html.replaceAll("12 core games + 22 specializations", "12 core games + 25 specializations")
      .replaceAll("twenty-two specializations", "twenty-five specializations")
      .replaceAll("22 SPECIALIST TRACKS", "25 SPECIALIST TRACKS")
      .replace("Learn the Update / Draw game loop, then tap targets with coordinate hit tests.", "Check where the player pressed, then learn the first rule: a hit adds a point.")
      .replace("Stop a moving marker in the target and understand frames and state.", "Stop a moving marker in the target and learn how time and game state advance.");
  }
  writeFileSync(file, html);
}

console.log("Generated three bilingual genre cards and updated the track count to 25.");
