#!/usr/bin/env node
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const lessons = [
  { slug: "swap-motion", ja: { t: "交換をアニメーションにする", l: "2個のピースを動かし、揃う交換は消去へ、揃わない交換は元の場所へ戻します。", c: "交換→判定→消去→落下を分ける", d: "盤面の数字を先に変えず、移動中の割合を0から1へ進めます。接触した瞬間に揃いを調べ、結果に応じて消去または戻りを始めます。" }, en: { t: "Animate Every Swap", l: "Move two pieces, clear a matching swap, and return an invalid swap to its original cells.", c: "Split swap, check, clear, and fall", d: "Do not change the board instantly. Advance a travel ratio from zero to one, check at contact, then clear or return based on the result." } },
  { slug: "reef-director", ja: { t: "3つの海底ステージをつなぐ", l: "手数・目標点・初期盤面・背景をデータにして、難しさの違う3面とベスト記録を作ります。", c: "同じルールへ違う目標データを渡す", d: "ステージ表から名前、手数、目標点、盤面を読みます。クリア時は残り手数をボーナスにし、次の海へ進めることで、一手を大切にしてもう一度遊ぶ理由が生まれます。" }, en: { t: "Direct Three Reef Stages", l: "Store moves, score goals, opening boards, and backgrounds as data for three escalating stages and a best run.", c: "Feed different goals into the same rules", d: "Read the name, moves, target, and board from a stage table. Turn spare moves into a clear bonus and advance to the next reef, giving every move value and every run a replay target." } },
];

const esc = (s) => s.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");

function swapLab(lang) {
  const ja = lang === "ja";
  return `<section class="motion-lab match-lab" data-lab="match-swap" data-lang="${lang}" aria-labelledby="match-swap-title"><div><p class="eyebrow">SWAP LAB</p><h2 id="match-swap-title">${ja ? "盤面の数字と見た目の時間を分ける" : "Separate board rules from visual time"}</h2><p>${ja ? "有効な交換は、2個が動く→接触で3個が揃う→消える→上から落ちる。不成立の交換は同じ時間を使って元へ戻ります。" : "A valid swap travels, checks three-in-a-row on contact, clears, and lets pieces fall. An invalid swap uses the same time to return."}</p></div><div class="match-swap-board" data-match-board><div class="match-grid" data-match-grid></div><p class="lab-readout" data-match-phase aria-live="polite">${ja ? "待機" : "READY"}</p><div class="lab-actions"><button type="button" data-match-action="valid">${ja ? "揃う交換を再生" : "PLAY MATCH"}</button><button type="button" data-match-action="invalid">${ja ? "戻る交換を再生" : "PLAY RETURN"}</button><button type="button" data-match-reset>${ja ? "リセット" : "RESET"}</button></div><p class="lab-readout" data-match-readout>${ja ? "ボタンを押して、状態の順番を見よう。" : "Choose a swap and follow each phase."}</p></div></section>`;
}

function page(lesson, lang, index) {
  const q = lesson[lang];
  const other = lang === "ja" ? "en" : "ja";
  const track = lang === "ja" ? "3マッチパズル" : "Match-three Puzzle";
  const code = index ? `type stage struct { moves, goal int; board [7][6]int }\nlevel := stages[stageIndex]\nbest = max(best, score+moves*bonus)` : `swapT += 0.1\ndrawX = startX + (endX-startX)*swapT\nif swapT >= 1 { resolveOrReturn() }`;
  const lab = index === 0 ? swapLab(lang) : `<section class="motion-lab" data-lab="turn"><div><p class="eyebrow">SYSTEM LAB</p><h2>${q.c}</h2><p>${q.d}</p></div><div class="lab-stage"><div class="lab-entities" data-count="4"></div><p class="lab-readout">START</p><button type="button" class="lab-action">NEXT</button></div></section>`;
  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><title>${q.t} | Ebi Showcase</title><link rel="stylesheet" href="../../../../style.css"></head><body><header class="nav"><a class="brand" href="../../../"><span>EBI</span> SHOWCASE</a><nav><a href="../">PATH</a><a class="lang" href="../../../../${other}/tracks/match3/${lesson.slug}/">${other.toUpperCase()}</a></nav></header><main><div class="lesson-breadcrumb"><a href="../">← ${track}</a><span>POLISH ${index + 1}/2</span></div><section class="lesson-hero"><div><p class="eyebrow">★★★★☆ / PLAYABLE</p><h1>${q.t}</h1><p>${q.l}</p><div class="lesson-meta"><span>GO + EBITENGINE</span><strong>PLAYABLE</strong></div></div></section><section class="play-panel" id="play"><div class="play-copy"><p class="eyebrow">PLAY FIRST</p><h2>${q.t}</h2><p>${q.l}</p></div><div class="game-frame"><iframe class="lesson-game-frame" src="../../../../play/${lesson.slug}/" title="${q.t}" allow="autoplay; fullscreen"></iframe></div></section><section class="lesson-section"><p class="eyebrow">DEEP DIVE</p><h2>${q.c}</h2><p>${q.d}</p><div class="concept-row"><article><b>01</b><h3>RULE</h3><p>${q.l}</p></article><article><b>02</b><h3>TIME</h3><p>${q.d}</p></article><article><b>03</b><h3>REPLAY</h3><p>${lang === "ja" ? "結果を読める演出と次の目標を用意します。" : "Pair readable feedback with a next target."}</p></article></div></section>${lab}<section class="code-lesson"><div><p class="eyebrow">GO + EBITENGINE</p><h2>${q.c}</h2><p>${q.l}</p></div><pre><code>${esc(code)}</code></pre></section><section class="why-grid"><article><h3>OBSERVE</h3><p>${q.d}</p></article><article><h3>CHANGE</h3><p>${lang === "ja" ? "移動時間や目標点を変えて観察します。" : "Change travel time or score goals."}</p></article><article><h3>CHALLENGE</h3><p>${lang === "ja" ? "4つ目の海をデータだけで追加しよう。" : "Add a fourth reef using data only."}</p></article></section><nav class="lesson-pager"><a href="../">← PATH<strong>${track}</strong></a><a href="../ebi-match/">FINAL →<strong>Ebi Match</strong></a></nav></main><footer><p>© Ebi Showcase · Apache-2.0</p></footer><script src="../../../../learn.js"></script></body></html>`;
}

for (const lang of ["ja", "en"]) {
  for (let i = 0; i < lessons.length; i++) {
    const lesson = lessons[i];
    const dir = join(root, "web", lang, "tracks", "match3", lesson.slug);
    mkdirSync(dir, { recursive: true });
    writeFileSync(join(dir, "index.html"), page(lesson, lang, i));
  }
  const hubPath = join(root, "web", lang, "tracks", "match3", "index.html");
  let hub = readFileSync(hubPath, "utf8").replace(/<!-- match3-polish:start -->[\s\S]*?<!-- match3-polish:end -->\n?/g, "");
  const cards = lessons.map((lesson, i) => { const q = lesson[lang]; return `<a class="path-step" href="${lesson.slug}/"><span>0${i + 6}</span><div><h3>${q.t}</h3><p>${q.l}</p><strong>${q.c}</strong></div><b>→</b></a>`; }).join("\n");
  const final = hub.match(/<a class="path-step" href="ebi-match\/[\s\S]*?<\/a>/)?.[0];
  if (!final) throw new Error("match3 final missing");
  hub = hub.replace(final, `<!-- match3-polish:start -->\n${cards}\n<!-- match3-polish:end -->\n${final.replace("<span>06</span>", "<span>08</span>")}`).replace(/<span>6 STEPS<\/span>/, "<span>8 STEPS</span>");
  writeFileSync(hubPath, hub);
  const finalPage = join(root, "web", lang, "tracks", "match3", "ebi-match", "index.html");
  let finalHtml = readFileSync(finalPage, "utf8").replace(/STEP 06 \/ 06/g, "STEP 08 / 08").replace(/<span>STEP 06<\/span>/g, "<span>STEP 08</span>");
  if (lang === "ja") finalHtml = finalHtml.replace("12手で650点をこえよう。", "3つの海をめぐり、ベスト記録を更新しよう。").replace("同じstageで何度でも再挑戦できます。", "クリア時の残り手数はボーナスになり、3つの海を通したベスト記録を何度でも更新できます。");
  else finalHtml = finalHtml.replace(/(?:Beat|Score) 650 points in 12 moves\./, "Clear three reefs and improve your best run.").replace("Retry with the same stage data.", "Spare moves become a bonus, so replay all three reefs to improve the best run.");
  writeFileSync(finalPage, finalHtml);
}
console.log("Generated 2 match-three polish lessons in JA/EN.");
