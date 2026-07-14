#!/usr/bin/env node
// Copyright 2026 Ebi Showcase contributors. Licensed under Apache-2.0.
// Generate the bilingual Reversi learning path and its home-page card.
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { dualLayerCodeLesson, authoringConceptRow } from "./authoring-lesson-helpers.mjs";

const root = new URL("..", import.meta.url).pathname;
const esc = (value) => String(value).replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");
function authoring(step, lang, index) {
  const ja = lang === "ja";
  const entryPath = `games/tracks/reversi/${step.slug}/main.go`;
  const entryCode = readFileSync(join(root, entryPath), "utf8").trim();
  const actions = ja
    ? ["盤面の初期配置データを1つ追加する", "合法手を検査する局面データを1つ追加する", "反転する石の局面を1つ追加する", "パスになる局面データを1つ追加する", "評価マップの1マスを変えず別の評価項を1つ足す"]
    : ["add one board setup datum", "add one position that checks a legal move", "add one flip-position case", "add one pass-position case", "add one evaluation term without changing a ScoreMap cell"];
  const cards = ja
    ? [{title:"DATA",body:"盤面、手、CPU設定が局面の真実を持つ。",code:"Board / Move / CPU"},{title:"UPDATE",body:"合法手、反転、パス、ターンを順に決める。",code:"legal → apply → turn"},{title:"DRAW",body:"盤面と評価を石、印、文字へ写す。",code:"board → pixels"}]
    : [{title:"DATA",body:"Board, move, and CPU configuration hold position truth.",code:"Board / Move / CPU"},{title:"UPDATE",body:"Legality, flips, passes, and turn are decided in order.",code:"legal → apply → turn"},{title:"DRAW",body:"Board and evaluation project into stones, marks, and text.",code:"board → pixels"}];
  return authoringConceptRow(cards) + dualLayerCodeLesson({lang,entryPath,entryCode,implementationPath:"internal/reversi (shared rules)",implementationCode:"// shared pure rules: legal move → captures → apply\ncaptures := Captures(board, player, move)",rule:{path:entryPath,location:"board / CPU configuration",action:actions[index],verify:ja?"`go test ./...` を実行し、ローカルで局面を確認する":"Run `go test ./...`, then verify the position locally"}});
}

const steps = [
  { slug: "board-grid", ja: { title: "8×8の盤面をデータで描く", lead: "64個のマスを二次元配列で持ち、クリックした場所を黒い石として表示します。", concept: "盤面と座標", deep: "画面に見えているマスを、board[y][x]というデータに置き換えます。Drawはデータを読むだけ、クリックは座標を書き込むだけです。", code: `const Size = 8\ntype Board [Size][Size]Player\nboard[y][x] = Black` }, en: { title: "Draw an 8×8 Board from Data", lead: "Store 64 cells in a two-dimensional array and show a black stone where you click.", concept: "Board data and coordinates", deep: "Turn visible cells into board[y][x] data. Draw only reads the data; a click only writes a coordinate.", code: `const Size = 8\ntype Board [Size][Size]Player\nboard[y][x] = Black` } },
  { slug: "legal-moves", ja: { title: "置ける場所を探す", lead: "8方向を調べ、相手の石を1つ以上はさんで自分の石につながる場所だけを合法手にします。", concept: "8方向の合法手判定", deep: "左・右・上下・斜めの8方向を1本ずつ進みます。相手の石の列の先に自分の石が見つかったとき、その空きマスへ置けます。", code: `for _, dir := range directions {\n    if bracketsOpponent(board, move, dir) {\n        legal = append(legal, move)\n    }\n}` }, en: { title: "Find Legal Moves", lead: "Scan eight directions and accept a square only when it brackets one or more opponent stones.", concept: "Eight-direction legality", deep: "Walk one line in each horizontal, vertical, and diagonal direction. An empty square is legal when an opponent run ends at your stone.", code: `for _, dir := range directions {\n    if bracketsOpponent(board, move, dir) {\n        legal = append(legal, move)\n    }\n}` } },
  { slug: "flip-stones", ja: { title: "はさんだ石を反転する", lead: "置いた石から8方向へ進み、はさんだ相手の列を自分の色へ変えます。", concept: "合法手と反転処理を分ける", deep: "Capturesは反転する座標の一覧を返し、Applyはその一覧を使って盤面を変更します。判定と変更を分けるとテストしやすくなります。", code: `captured := Captures(board, player, move)\nboard[move.Y][move.X] = player\nfor _, cell := range captured {\n    board[cell.Y][cell.X] = player\n}` }, en: { title: "Flip the Bracketed Stones", lead: "Walk in eight directions from a move and turn every bracketed opponent stone to your color.", concept: "Separate legality from mutation", deep: "Captures returns the coordinates to flip; Apply changes the board using that list. Separating the check from the mutation makes tests easier.", code: `captured := Captures(board, player, move)\nboard[move.Y][move.X] = player\nfor _, cell := range captured {\n    board[cell.Y][cell.X] = player\n}` } },
  { slug: "pass-and-score", ja: { title: "パスと勝敗を実装する", lead: "置ける場所がないときはパスし、両者が続けてパスするか盤面が埋まったら石の数を比べます。", concept: "ターン・パス・スコア", deep: "合法手が空ならターンだけを交代します。2回連続のパス、または満杯を終了条件にして、黒と白の石数を数えます。", code: `if len(ValidMoves(board, turn)) == 0 {\n    passes++\n    turn = Opponent(turn)\n}\nif passes == 2 || Full(board) { finish() }` }, en: { title: "Add Passes and a Score", lead: "Pass when no move exists, then compare stone counts after two consecutive passes or a full board.", concept: "Turns, passes, and score", deep: "An empty legal-move list changes only the turn. End after two consecutive passes or a full board, then count black and white stones.", code: `if len(ValidMoves(board, turn)) == 0 {\n    passes++\n    turn = Opponent(turn)\n}\nif passes == 2 || Full(board) { finish() }` } },
  { slug: "ebi-reversi", ja: { title: "3種類のCPUと対戦する", lead: "おだやか・位置重視・先読みの3種類のCPUを選び、角が高く角の隣が低い8×8スコアマップで対戦します。", concept: "位置スコアと先読みCPU", deep: "各マスに重みを置き、自分の石なら加点、相手の石なら減点します。先読みCPUは、自分の一手のあとに相手が選べる一番よい返し手まで調べます。", code: `move, ok, score := ChooseLookahead(board, White)\nif ok {\n    Apply(&board, White, move)\n}\n// score is the worst reply the opponent can force\n_ = score` }, en: { title: "Play Three CPU Personalities", lead: "Choose a friendly, positional, or look-ahead CPU and play with an 8×8 score map where corners are valuable and adjacent squares are risky.", concept: "Position scores and look-ahead CPU", deep: "Give every square a weight. Add it for your stones and subtract it for the opponent. The look-ahead CPU also checks the opponent's best reply after its move.", code: `move, ok, score := ChooseLookahead(board, White)\nif ok {\n    Apply(&board, White, move)\n}\n// score is the worst reply the opponent can force\n_ = score` } },
];

function page(step, lang, index) {
  const text = step[lang];
  const other = lang === "ja" ? "en" : "ja";
  const track = lang === "ja" ? "コンピュータと対戦できるリバーシを作ろう" : "Build Reversi with a Computer";
  const play = `../../../../play/${step.slug}/?lang=${lang}`;
  const lab = index === steps.length - 1 ? `<section class="motion-lab reversi-eval-lab" data-lab="reversi-eval" data-lang="${lang}"><div><p class="eyebrow">EVALUATION LAB</p><h2>8×8 Score Map</h2><p>${lang === "ja" ? "石を置いて、位置の重みが評価値をどう変えるか見てみましょう。角は高く、角の隣は低く設定されています。" : "Place stones and watch the position weights change the score. Corners are high; their neighbors are low."}</p></div><div class="reversi-eval-controls"><div class="reversi-eval-map" data-reversi-map aria-label="8 by 8 evaluation map"></div><p class="reversi-eval-score" data-reversi-score aria-live="polite"></p><div><button type="button" data-reversi-place="corner">CORNER</button><button type="button" data-reversi-place="risky">RISKY</button><button type="button" data-reversi-place="center">CENTER</button><button type="button" data-reversi-place="reset">RESET</button></div></div></section>` : `<section class="motion-lab" data-lab="turn" data-states="BOARD,SCAN,FLIP,SCORE"><div><p class="eyebrow">SYSTEM LAB</p><h2>${text.concept}</h2><p>${lang === "ja" ? "ボタンを押して、データが次の処理へ渡る順番を確かめます。" : "Step through the system and follow data into the next operation."}</p></div><div class="lab-stage"><div class="lab-entities" data-count="4"></div><p class="lab-readout" aria-live="polite">START</p><button type="button" class="lab-action">NEXT</button></div></section>`;
  const controls = lang === "ja" ? "タップ／クリックに加え、矢印キーまたはWASDで黄色いカーソルを動かし、EnterかSpaceで置けます。1・2・3でCPUを選び、RまたはREPLAYで新しい対局を始めます。" : "Tap or click to play. You can also move the yellow cursor with arrows or WASD, then place with Enter or Space. Choose a CPU with 1, 2, or 3; start a new match with R or REPLAY.";
  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><title>${text.title} | Ebi Showcase</title><link rel="stylesheet" href="../../../../style.css"></head><body><header class="nav"><a class="brand" href="../../../"><span>EBI</span> SHOWCASE</a><nav><a href="../">PATH</a><a class="lang" href="../../../../${other}/tracks/reversi/${step.slug}/">${other.toUpperCase()}</a></nav></header><main><div class="lesson-breadcrumb"><a href="../">← ${track}</a><span>STEP 0${index + 1} / 05</span></div><section class="lesson-hero"><div><p class="eyebrow">★★★★★ / PLAYABLE</p><h1>${text.title}</h1><p>${text.lead}</p><div class="lesson-meta"><span>GO + EBITENGINE</span><strong>PLAYABLE</strong></div></div></section><section class="play-panel" id="play"><div class="play-copy"><p class="eyebrow">PLAY FIRST</p><h2>${text.title}</h2><p>${controls}</p></div><div class="game-frame"><iframe class="lesson-game-frame" src="${play}" title="${text.title}" allow="autoplay; fullscreen"></iframe></div></section><section class="lesson-section"><p class="eyebrow">DEEP DIVE</p><h2>${text.concept}</h2><p>${text.deep}</p><div class="concept-row"><article><b>01</b><h3>SEE</h3><p>${text.lead}</p></article><article><b>02</b><h3>STATE</h3><p>${text.deep}</p></article><article><b>03</b><h3>DECIDE</h3><p>${lang === "ja" ? "見えている盤面から、次の一手の理由を言葉にします。" : "Explain the reason for the next move from the visible board."}</p></article></div></section>${lab}<section class="code-lesson"><div><p class="eyebrow">GO + EBITENGINE</p><h2>${text.concept}</h2><p>${lang === "ja" ? "実ゲームで使っている中心ルールを小さく抜き出します。" : "Here is the central rule used by the playable game."}</p></div><pre><code>${esc(text.code)}</code></pre></section><section class="why-grid"><article><h3>OBSERVE</h3><p>${text.deep}</p></article><article><h3>CHANGE</h3><p>${lang === "ja" ? "CPUを1・2・3で切り替え、同じ局面でどの手を選ぶか比べよう。" : "Switch CPUs with 1, 2, and 3, then compare their moves from the same position."}</p></article><article class="challenge"><h3>CHALLENGE</h3><p>${lang === "ja" ? "ScoreMapや先読みの深さを変え、どんなCPUが強く・楽しくなるか試そう。" : "Change ScoreMap or look-ahead depth and discover what makes a CPU strong and fun."}</p></article></section><nav class="lesson-pager"><a href="../">← PATH<strong>${track}</strong></a><a href="${index + 1 < steps.length ? `../${steps[index + 1].slug}/` : "../"}">${index + 1 < steps.length ? "NEXT →" : "COMPLETE →"}<strong>${index + 1 < steps.length ? steps[index + 1][lang].title : track}</strong></a></nav></main><footer><p>© Ebi Showcase · Apache-2.0</p></footer><script src="../../../../learn.js"></script></body></html>`;
}

function hub(lang) {
  const ja = lang === "ja";
  const other = ja ? "en" : "ja";
  const track = ja ? "コンピュータと対戦できるリバーシを作ろう" : "Build Reversi with a Computer";
  const lead = ja ? "8×8の盤面をデータで持ち、合法手・石の反転・パス・勝敗を順番に実装します。最後は、おだやか・位置重視・先読みのCPUを選んで対戦します。" : "Represent an 8×8 board, then add legal moves, flips, passes, and scores. Finish by choosing and playing friendly, positional, and look-ahead CPUs.";
  const cards = steps.map((step, i) => `<a class="path-step" href="${step.slug}/"><span>0${i + 1}</span><div><h3>${step[lang].title}</h3><p>${step[lang].lead}</p><strong>${step[lang].concept}</strong></div><b>→</b></a>`).join("");
  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><title>${track} | Ebi Showcase</title><link rel="stylesheet" href="../../../style.css"></head><body><header class="nav"><a class="brand" href="../../"><span>EBI</span> SHOWCASE</a><nav><a href="../../">${ja ? "全コース" : "ALL COURSES"}</a><a class="lang" href="../../../${other}/tracks/reversi/">${other.toUpperCase()}</a></nav></header><main><div class="lesson-breadcrumb"><a href="../../">← ${ja ? "全コース" : "ALL COURSES"}</a><span>5 PLAYABLE STEPS</span></div><section class="track-hero track-reversi"><span class="track-letter">U</span><div><p class="eyebrow">${ja ? "盤面を読む、石を返す、CPUに勝つ" : "Read the board, flip stones, beat the CPU"}</p><h1>${track}</h1><p>${lead}</p></div></section><p class="curriculum-bridge">${ja ? "戦略SLGで学んだ盤面の読み取りを、8方向の探索と位置評価へ広げます。" : "Extend board reading from the tactical RPG into eight-direction search and position evaluation."}</p><section class="path-list"><div class="path-intro"><p class="eyebrow">LEARNING PATH</p><h2>${ja ? "石を置くところから、CPUの考え方まで。" : "From placing stones to a thinking CPU."}</h2><p>${lead}</p></div>${cards}</section></main><footer><p>© Ebi Showcase · Apache-2.0</p></footer><script src="../../../learn.js"></script></body></html>`;
}

const homeCards = {
  ja: `<a class="track-card track-reversi" href="tracks/reversi/"><span class="track-letter">U</span><div><p class="eyebrow">8×8の盤面と3種類のCPU</p><h3>コンピュータと対戦できるリバーシを作ろう</h3><p>合法手、石の反転、パス、勝敗を実装し、おだやか・位置重視・先読みのCPUと対戦します。</p><strong>5 STEPS →</strong></div></a>`,
  en: `<a class="track-card track-reversi" href="tracks/reversi/"><span class="track-letter">U</span><div><p class="eyebrow">An 8×8 board and three CPUs</p><h3>Build Reversi with a Computer</h3><p>Build legal moves, flips, passes, and scores, then play friendly, positional, and look-ahead CPUs.</p><strong>5 STEPS →</strong></div></a>`,
};

for (const lang of ["ja", "en"]) {
  const base = join(root, "web", lang, "tracks", "reversi");
  mkdirSync(base, { recursive: true });
  writeFileSync(join(base, "index.html"), hub(lang));
  for (let i = 0; i < steps.length; i++) {
    const dir = join(base, steps[i].slug);
    mkdirSync(dir, { recursive: true });
    let html = page(steps[i], lang, i);
    const authored = authoring(steps[i], lang, i);
    const split = authored.indexOf('<section class="code-lesson">');
    const concept = authored.slice(0, split);
    const panels = authored.slice(split);
    html = html
      .replace(/<div class="concept-row">[\s\S]*?<\/div><\/section>/, `${concept}</section>`)
      .replace(/<section class="code-lesson">[\s\S]*?<\/section><section class="why-grid">[\s\S]*?<\/section>/, panels);
    writeFileSync(join(dir, "index.html"), html);
  }
  const file = join(root, "web", lang, "index.html");
  let html = readFileSync(file, "utf8");
  html = html.replace(/<!-- reversi-track:start -->[\s\S]*?<!-- reversi-track:end -->\n?/g, "");
  const marker = "<!-- expansion-tracks:start -->";
  html = html.replace(marker, `<!-- reversi-track:start -->\n${homeCards[lang]}\n<!-- reversi-track:end -->\n${marker}`);
  html = html.replace(lang === "ja" ? "20の専門コース" : "twenty specializations", lang === "ja" ? "21の専門コース" : "twenty-one specializations");
  writeFileSync(file, html);
}

console.log("Generated a 5-step bilingual Reversi track and home cards.");
