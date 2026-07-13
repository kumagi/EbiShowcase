#!/usr/bin/env node
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const lessons = [
  { slug: "camera-rooms", ja: { t: "世界座標と追従カメラ", l: "画面4枚分の世界を歩き、世界座標からカメラ座標を引いて一部分だけ描きます。", c: "世界と画面を別の座標にする", d: "主人公は世界のどこにいるかを持ち、DrawだけがcameraXを引きます。差の10%ずつ追えば滑らかなカメラになります。" }, en: { t: "World Coordinates and a Follow Camera", l: "Walk through a four-screen world and draw only one slice by subtracting camera position.", c: "Separate world and screen coordinates", d: "The player stores a world position; only Draw subtracts cameraX. Following ten percent of the gap creates a smooth camera." } },
  { slug: "exploration-map", ja: { t: "歩いた部屋だけ広がる地図", l: "部屋番号を訪問済み集合へ記録し、未探索と探索済みを色分けします。", c: "探索状態を部屋IDで保存する", d: "現在位置を部屋幅で割ってroom IDを求めます。一度訪れたIDをboolへ残せば、戻っても地図は消えません。" }, en: { t: "A Map that Grows as You Explore", l: "Store visited room IDs and color unexplored and discovered rooms differently.", c: "Save exploration by room ID", d: "Divide world position by room width for the room ID. Keeping visited IDs as booleans preserves the map when the player returns." } },
  { slug: "dash-gate", ja: { t: "能力で開くダッシュゲート", l: "最初は通れない封印を置き、ダッシュ取得後に同じ場所を突破できるようにします。", c: "能力フラグが地形の意味を変える", d: "dashがfalseなら壁として止め、trueならX入力で大きな速度を与えます。同じ地形が取得前は障害、取得後は通路になります。" }, en: { t: "A Dash-gated Passage", l: "Block a seal at first, then let the player burst through after finding dash.", c: "Ability flags change what terrain means", d: "Stop at the wall while dash is false; with dash true, X grants a large velocity. The same terrain changes from obstacle to route." } },
  { slug: "ability-route", ja: { t: "能力で昔の道を読み替える", l: "地図、ダッシュ、高跳びを順に取得し、以前見た障害へ戻って奥へ進みます。", c: "取得順と戻り道を設計する", d: "能力と対応するゲートを交互に配置します。新能力を得た直後に古い場所を思い出させると、世界が広がったと感じられます。" }, en: { t: "Reinterpret Old Routes with Abilities", l: "Collect the double jump, then return to an earlier high ledge to reach the relic.", c: "Design acquisition order and backtracking", d: "A low route gives the ability first. Returning to the old ledge makes the same space newly playable." } },
];

const esc = (s) => s.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");

function labMarkup(slug, lang, q) {
  const ja = lang === "ja";
  if (slug === "ability-route") {
    return `<section class="motion-lab metroid-lab" data-lab="ability-gate" data-lang="${lang}" aria-labelledby="ability-lab-title"><div><p class="eyebrow">ABILITY LAB</p><h2 id="ability-lab-title">${ja ? "二段ジャンプで高台へ" : "Reach the high ledge with a double jump"}</h2><p>${ja ? "左の高台は最初のジャンプだけでは届きません。まず右の羽を拾い、同じジャンプをもう一度押して高台へ戻ります。" : "The left ledge is too high for one jump. Collect the wings on the right, then press jump again in mid-air to return to the ledge."}</p></div><div class="ability-lab-board" data-lab-board><div class="ability-lab-world"><span class="ability-lab-player" data-lab-player></span><span class="ability-lab-ledge" data-lab-ledge>${ja ? "高台" : "LEDGE"}</span><span class="ability-lab-wings" data-lab-wings>✦</span><span class="ability-lab-relic">◆</span></div><p class="lab-readout" data-lab-status aria-live="polite">${ja ? "右の羽を拾おう" : "Find the wings on the right"}</p><div class="lab-actions"><button type="button" data-ability-move="left">← ${ja ? "戻る" : "BACK"}</button><button type="button" data-ability-jump>${ja ? "ジャンプ" : "JUMP"}</button><button type="button" data-ability-move="right">${ja ? "進む" : "GO"} →</button><button type="button" data-ability-reset>${ja ? "最初から" : "RESET"}</button></div><p class="lab-readout" data-lab-ability>${ja ? "能力: なし" : "ABILITY: none"}</p></div></section>`;
  }
  if (slug === "exploration-map") {
    return `<section class="motion-lab metroid-lab" data-lab="metroid-map" data-lang="${lang}" aria-labelledby="map-lab-title"><div><p class="eyebrow">MAP LAB</p><h2 id="map-lab-title">${ja ? "進んだ部屋だけ地図に残す" : "Keep only visited rooms on the map"}</h2><p>${ja ? "次の部屋・前の部屋を押すと、現在位置と訪問済みの部屋が同時に変わります。未探索の部屋はまだ見えません。" : "Move to the next or previous room. The current room and the visited map change together; unexplored rooms stay hidden."}</p></div><div class="map-lab-board" data-lab-board><div class="map-lab-world"><div class="map-lab-path" data-map-path></div></div><div class="map-lab-controls"><button type="button" data-map-move="back">← ${ja ? "前の部屋" : "BACK"}</button><button type="button" data-map-move="next">${ja ? "次の部屋" : "NEXT"} →</button><button type="button" data-map-reset>${ja ? "地図をリセット" : "RESET MAP"}</button></div><p class="lab-readout" data-map-readout aria-live="polite">${ja ? "部屋ID 0 / 訪問 1" : "room ID 0 / visited 1"}</p></div></section>`;
  }
  return `<section class="motion-lab" data-lab="turn"><div><p class="eyebrow">SYSTEM LAB</p><h2>${q.c}</h2><p>${q.d}</p></div><div class="lab-stage"><div class="lab-entities" data-count="4"></div><p class="lab-readout">START</p><button type="button" class="lab-action">NEXT</button></div></section>`;
}

function page(x, lang, index) {
  const q = x[lang];
  const other = lang === "ja" ? "en" : "ja";
  const track = lang === "ja" ? "巨大マップ探索アクション" : "Large-map Exploration";
  const codes = [
    "cameraX += (playerX-screenCenter-cameraX)*0.1\nscreenX := worldX-cameraX",
    "room := int(playerX/roomWidth)\nvisited[room] = true",
    "if !hasDash && playerX > sealX { playerX = sealX }\nif hasDash && pressedX { velocityX = dashSpeed }",
    "if hasWings && pressedJump && !onGround { jumpsLeft-- }\nif playerX < oldLedge { relic = true }",
  ];
  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><title>${q.t} | Ebi Showcase</title><link rel="stylesheet" href="../../../../style.css"></head><body><header class="nav"><a class="brand" href="../../../"><span>EBI</span> SHOWCASE</a><nav><a href="../">PATH</a><a class="lang" href="../../../../${other}/tracks/metroidvania/${x.slug}/">${other.toUpperCase()}</a></nav></header><main><div class="lesson-breadcrumb"><a href="../">← ${track}</a><span>STEP ${index + 1}/5</span></div><section class="lesson-hero"><div><p class="eyebrow">★★★★☆ / PLAYABLE</p><h1>${q.t}</h1><p>${q.l}</p><div class="lesson-meta"><span>GO + EBITENGINE</span><strong>PLAYABLE</strong></div></div></section><section class="play-panel"><div class="play-copy"><p class="eyebrow">PLAY FIRST</p><h2>${q.t}</h2><p>${q.l}</p></div><div class="game-frame"><iframe class="lesson-game-frame" src="../../../../play/${x.slug}/" title="${q.t}" allow="autoplay; fullscreen"></iframe></div></section><section class="lesson-section"><p class="eyebrow">DEEP DIVE</p><h2>${q.c}</h2><p>${q.d}</p><div class="concept-row"><article><b>01</b><h3>WORLD</h3><p>${q.l}</p></article><article><b>02</b><h3>STATE</h3><p>${q.d}</p></article><article><b>03</b><h3>RETURN</h3><p>${lang === "ja" ? "前に見た場所へ意味を足します。" : "Add new meaning to an old place."}</p></article></div></section>${labMarkup(x.slug, lang, q)}<section class="code-lesson"><div><p class="eyebrow">GO + EBITENGINE</p><h2>${q.c}</h2><p>${q.l}</p></div><pre><code>${esc(codes[index])}</code></pre></section><section class="why-grid"><article><h3>OBSERVE</h3><p>${q.d}</p></article><article><h3>CHANGE</h3><p>${lang === "ja" ? "部屋幅や能力値を変えます。" : "Change room width or ability strength."}</p></article><article><h3>CHALLENGE</h3><p>${lang === "ja" ? "新能力と対応ゲートを追加しよう。" : "Add another ability and matching gate."}</p></article></section><nav class="lesson-pager"><a href="../">← PATH<strong>${track}</strong></a><a href="../ebi-depths/">FINAL →<strong>Ebi Depths</strong></a></nav></main><footer><p>© Ebi Showcase · Apache-2.0</p></footer><script src="../../../../learn.js"></script></body></html>`;
}

for (const lang of ["ja", "en"]) {
  for (let i = 0; i < lessons.length; i++) {
    const lesson = lessons[i];
    const dir = join(root, "web", lang, "tracks", "metroidvania", lesson.slug);
    mkdirSync(dir, { recursive: true });
    writeFileSync(join(dir, "index.html"), page(lesson, lang, i));
  }
  const hubPath = join(root, "web", lang, "tracks", "metroidvania", "index.html");
  let hub = readFileSync(hubPath, "utf8");
  const final = hub.match(/<a class="path-step" href="ebi-depths\/">[\s\S]*?<\/a>/)?.[0];
  if (!final) throw new Error("metroidvania final missing");
  const cards = lessons.map((lesson, i) => {
    const q = lesson[lang];
    return `<a class="path-step" href="${lesson.slug}/"><span>0${i + 1}</span><div><h3>${q.t}</h3><p>${q.l}</p><strong>${q.c}</strong></div><b>→</b></a>`;
  }).join("\n");
  hub = hub.replace(final, `${cards}\n${final.replace("<span>01</span>", "<span>05</span>")}`).replace(/<span>1 PLAYABLE GAME<\/span>/, "<span>5 PLAYABLE GAMES</span>");
  writeFileSync(hubPath, hub);
}
console.log("Generated metroidvania track expansion in JA/EN.");
