#!/usr/bin/env node
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const lessons = [
  {
    slug: "block-feedback",
    ja: {
      title: "ブロックが壊れる一瞬を作る",
      lead: "道具を振りかぶり、当て、破片が飛び、元の姿勢へ戻るまでを時間に分けます。",
      deep: "数字をすぐ減らすだけでは、何が起きたか目で追えません。攻撃を「予備動作・接触・戻り」に分け、接触した1 tickだけでHP、破片、揺れ、得点を同時に変えます。",
      concept: ["予備動作で次の出来事を知らせる", "接触フレームで結果を一度だけ決める", "破片とつぶれで素材の重さを伝える"],
      challenge: "石と木で破片の飛び方を変えてみよう。",
    },
    en: {
      title: "Build the Moment a Block Breaks",
      lead: "Split a tool swing into anticipation, contact, flying chips, and recovery.",
      deep: "If a number changes instantly, the player cannot follow what happened. Divide the action into wind-up, contact, and recovery, then change HP, chips, shake, and score together on exactly one contact frame.",
      concept: ["Use anticipation to announce what comes next", "Resolve the result once on the contact frame", "Use chips and squash to communicate material weight"],
      challenge: "Give stone and wood different chip motion.",
    },
  },
  {
    slug: "island-director",
    ja: {
      title: "島データから3つの冒険を作る",
      lead: "地形、資源、敵、必要な明かりをデータにまとめ、同じルールで違う島を順番に遊びます。",
      deep: "ステージごとにプログラムを丸ごとコピーせず、島の名前、色、資源数、敵数、目標をひとまとまりのデータにします。クリアしたら番号を1つ進めて次のデータを読みます。",
      concept: ["島ごとの差をstage構造体へ入れる", "クリア条件をデータから組み立てる", "3島の合計得点を次の挑戦の目標にする"],
      challenge: "雨で歩きにくい第4の島を考えてみよう。",
    },
    en: {
      title: "Build Three Adventures from Island Data",
      lead: "Store terrain, resources, enemies, and required lights as data, then play different islands with the same rules.",
      deep: "Do not copy the whole program for every stage. Group an island's name, colors, resources, enemies, and goal into one data record. After a clear, advance one index and load the next record.",
      concept: ["Put island differences in a stage struct", "Build clear conditions from data", "Use the three-island total as the next replay target"],
      challenge: "Design a fourth rainy island that slows movement.",
    },
  },
];
const finalLesson = {
  slug: "ebi-craft",
  ja: {
    title: "Ebi Craft Expedition",
    lead: "苔の野営地、結晶洞窟、火の島を巡り、採取・道具・建築・敵を一つの冒険へつなぎます。",
    deep: "完成版は3つの島を同じgame状態で順番に読み込みます。島ごとに資源の位置、必要なランタン、夜の敵が変わります。採掘には予備動作と接触があり、破片・点滅・揺れが結果を伝えます。残りHPと速さが得点になるため、クリア後にも近道を考えて再挑戦できます。",
    concept: ["3島の地形・資源・敵をデータから読み込む", "採掘と建築をアニメーションと粒子でつなぐ", "合計得点とベスト記録で別の採取順を試す"],
    challenge: "一度通ったマスだけ見える探索地図を追加しよう。",
  },
  en: {
    title: "Ebi Craft Expedition",
    lead: "Cross Moss Camp, Crystal Cave, and Ember Isle, joining gathering, tools, building, and enemies into one adventure.",
    deep: "The final game loads three islands in sequence through the same game state. Resource positions, required lanterns, and night enemies change on each island. Mining has anticipation and contact, while chips, flash, and shake communicate the result. Remaining HP and speed become score, so a clear invites a faster route on the next expedition.",
    concept: ["Load three islands' terrain, resources, and enemies from data", "Connect mining and building with animation and particles", "Try a different gathering route against the best total"],
    challenge: "Add an exploration map that reveals only visited tiles.",
  },
};

const esc = (s) => s.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");

function page(item, lang, index) {
  const q = item[lang], other = lang === "ja" ? "en" : "ja";
  const track = lang === "ja" ? "2Dブロックサンドボックス" : "2D Block Sandbox";
  const labels = lang === "ja" ? { back: "学習コース", play: "まず遊ぶ", observe: "観察", change: "変える", next: "次へ" } : { back: "Learning path", play: "Play first", observe: "Observe", change: "Change", next: "Next" };
  const code = index === 0 ? `phase--\nif phase == contactFrame {\n    block.HP--\n    spawnChips()\n    shake = 8\n}` : `stage := islands[stageIndex]\nload(stage.resources)\nif cleared(stage.goal) {\n    stageIndex++\n}`;
  const next = index === 0 ? `<a href="../island-director/">${labels.next} →<strong>${lessons[1][lang].title}</strong></a>` : index === 1 ? `<a href="../ebi-craft/">FINAL →<strong>Ebi Craft Expedition</strong></a>` : `<a href="../">PATH →<strong>${track}</strong></a>`;
  const stepLabel = index < 2 ? `POLISH ${index + 1}/2` : "STEP 09 / FINAL";
  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><title>${q.title} | Ebi Showcase</title><link rel="stylesheet" href="../../../../style.css"></head><body><header class="nav"><a class="brand" href="../../../"><span>EBI</span> SHOWCASE</a><nav><a href="../">PATH</a><a class="lang" href="../../../../${other}/tracks/sandbox/${item.slug}/">${other.toUpperCase()}</a></nav></header><main><div class="lesson-breadcrumb"><a href="../">← ${labels.back}</a><span>${stepLabel}</span></div><section class="lesson-hero"><div><p class="eyebrow">★★★★☆ / PLAYABLE</p><h1>${q.title}</h1><p>${q.lead}</p><div class="lesson-meta"><span>GO + EBITENGINE</span><strong>PLAYABLE</strong></div></div></section><section class="play-panel" id="play"><div class="play-copy"><p class="eyebrow">${labels.play}</p><h2>${q.title}</h2><p>${q.lead}</p></div><div class="game-frame"><iframe class="lesson-game-frame" src="../../../../play/${item.slug}/" title="${q.title}" allow="autoplay; fullscreen"></iframe></div></section><section class="lesson-section"><p class="eyebrow">DEEP DIVE</p><h2>${q.title}</h2><p>${q.deep}</p><div class="concept-row">${q.concept.map((c, i) => `<article><b>0${i + 1}</b><h3>${["STATE", "UPDATE", "FEEDBACK"][i]}</h3><p>${c}</p></article>`).join("")}</div></section><section class="motion-lab" data-lab="turn"><div><p class="eyebrow">SYSTEM LAB</p><h2>${q.title}</h2><p>${q.deep}</p></div><div class="lab-stage"><div class="lab-entities" data-count="4"></div><p class="lab-readout">START</p><button type="button" class="lab-action">NEXT</button></div></section><section class="code-lesson"><div><p class="eyebrow">GO + EBITENGINE</p><h2>${q.concept[1]}</h2><p>${q.lead}</p></div><pre><code>${esc(code)}</code></pre></section><section class="why-grid"><article><h3>${labels.observe}</h3><p>${q.deep}</p></article><article><h3>${labels.change}</h3><p>${q.concept[2]}</p></article><article><h3>CHALLENGE</h3><p>${q.challenge}</p></article></section><nav class="lesson-pager"><a href="../">← PATH<strong>${track}</strong></a>${next}</nav></main><footer><p>© Ebi Showcase · Apache-2.0</p></footer><script src="../../../../learn.js"></script></body></html>`;
}

for (const lang of ["ja", "en"]) {
  for (let i = 0; i < lessons.length; i++) {
    const item = lessons[i], dir = join(root, "web", lang, "tracks", "sandbox", item.slug);
    mkdirSync(dir, { recursive: true });
    writeFileSync(join(dir, "index.html"), page(item, lang, i));
  }
  writeFileSync(join(root, "web", lang, "tracks", "sandbox", "ebi-craft", "index.html"), page(finalLesson, lang, 2));
  const hub = join(root, "web", lang, "tracks", "sandbox", "index.html");
  let html = readFileSync(hub, "utf8").replace(/<!-- sandbox-polish:start -->[\s\S]*?<!-- sandbox-polish:end -->\n?/g, "");
  const cards = lessons.map((item, i) => { const q = item[lang]; return `<a class="path-step" href="${item.slug}/"><span>0${i + 7}</span><div><h3>${q.title}</h3><p>${q.lead}</p><strong>${q.concept[1]}</strong></div><b>→</b></a>`; }).join("\n");
  const final = html.match(/<a class="path-step" href="ebi-craft\/">[\s\S]*?<\/a>/)?.[0];
  if (!final) throw new Error("sandbox final card missing");
  const fq = finalLesson[lang];
  const finalCard = `<a class="path-step" href="ebi-craft/"><span>09</span><div><h3>${fq.title}</h3><p>${fq.lead}</p><strong>${fq.concept[2]}</strong></div><b>→</b></a>`;
  html = html.replace(final, `<!-- sandbox-polish:start -->\n${cards}\n<!-- sandbox-polish:end -->\n${finalCard}`).replace(/<span>\d+ STEPS<\/span>/, "<span>9 STEPS</span>");
  writeFileSync(hub, html);
}
console.log("Generated 2 sandbox polish lessons in JA/EN.");
