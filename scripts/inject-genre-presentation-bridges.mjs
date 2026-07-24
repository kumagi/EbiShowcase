#!/usr/bin/env node
import { existsSync, readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join, relative } from "node:path";
import {
  genrePresentationMap,
  presentationPatterns,
  validateGenrePresentationMap,
} from "./genre-presentation-map.mjs";
import { curriculum } from "./curriculum.mjs";

const root = new URL("..", import.meta.url).pathname;
const checkOnly = process.argv.includes("--check");
const expected = validateGenrePresentationMap();
const start = "<!-- genre-presentation-bridge:start -->";
const end = "<!-- genre-presentation-bridge:end -->";
const oldBlock = new RegExp(`${start.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}[\\s\\S]*?${end.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}`, "g");

function indexFiles(dir) {
  const files = [];
  for (const name of readdirSync(dir)) {
    const path = join(dir, name);
    if (statSync(path).isDirectory()) files.push(...indexFiles(path));
    else if (name === "index.html") files.push(path);
  }
  return files;
}

function patternLinks(entry, lang, prefix) {
  return entry.patterns.map((id) => {
    const pattern = presentationPatterns[id];
    return `<a href="${prefix}visual-effects/${pattern.slug}/"><b>${id}</b><span>${pattern[lang]}</span></a>`;
  }).join("");
}

function hubBridge(entry, lang) {
  const c = entry[lang];
  const prefix = "../";
  if (lang === "ja") {
    return `${start}<section class="genre-presentation-bridge" aria-labelledby="genre-motion-${entry.id}">
      <div class="genre-presentation-head"><div><p class="eyebrow">SPECIALIZATION / 好きなジャンルから開始</p><h2 id="genre-motion-${entry.id}">VFXは必要な時に戻る。<br>まず${c.name}のルールを作る。</h2></div><p>${c.focus}<br><strong>下の3つは演出で迷った時の参考ページです。未受講でも、このコースから始められます。</strong></p></div>
      <nav class="genre-pattern-links" aria-label="必要な時に参照するVisual Effects設計">${patternLinks(entry, lang, prefix)}</nav>
      <div class="genre-presentation-contract">
        <article><p class="eyebrow">MOTION BASELINE</p><h3>動いて見えることも完成条件</h3><p>${c.motion}</p></article>
        <article><p class="eyebrow">TEST SEAM</p><h3>派手さを外してもルールは同じ</h3><p>${c.test}</p></article>
      </div>
    </section>${end}`;
  }
  return `${start}<section class="genre-presentation-bridge" aria-labelledby="genre-motion-${entry.id}">
    <div class="genre-presentation-head"><div><p class="eyebrow">SPECIALIZATION / START WITH A GENRE</p><h2 id="genre-motion-${entry.id}">Return to VFX when needed.<br>Start with the rules of ${c.name}.</h2></div><p>${c.focus}<br><strong>The three links below are references for when presentation gets tricky. You can start this path without completing them first.</strong></p></div>
    <nav class="genre-pattern-links" aria-label="Visual Effects references to use when needed">${patternLinks(entry, lang, prefix)}</nav>
    <div class="genre-presentation-contract">
      <article><p class="eyebrow">MOTION BASELINE</p><h3>Readable motion is part of done</h3><p>${c.motion}</p></article>
      <article><p class="eyebrow">TEST SEAM</p><h3>The same rules survive without spectacle</h3><p>${c.test}</p></article>
    </div>
  </section>${end}`;
}

function lessonBridge(entry, lang) {
  const c = entry[lang];
  const prefix = "../../";
  const intro = lang === "ja"
    ? `<strong>演出の言葉が初めてでも、この章から進めます。</strong><span>${c.focus} 下のAリンクは、必要になった時だけ戻る参考ページです。</span>`
    : `<strong>You can start this chapter even if the presentation terms are new.</strong><span>${c.focus} The A-links below are references to revisit only when you need them.</span>`;
  const contract = lang === "ja"
    ? `<div class="genre-presentation-inline-contract"><p><b>MOTION BASELINE</b>${c.motion}</p><p><b>TEST SEAM</b>${c.test}</p></div>`
    : `<div class="genre-presentation-inline-contract"><p><b>MOTION BASELINE</b>${c.motion}</p><p><b>TEST SEAM</b>${c.test}</p></div>`;
  return `${start}<aside class="genre-presentation-inline">${intro}<nav aria-label="${lang === "ja" ? "参照するVFX設計" : "VFX design references"}">${patternLinks(entry, lang, prefix)}</nav>${contract}</aside>${end}`;
}

function insertAfterOpeningSection(html, block) {
  const mainAt = html.search(/<main(?:\s|>)/);
  if (mainAt < 0) throw new Error("missing <main>");
  const sectionEnd = html.indexOf("</section>", mainAt);
  if (sectionEnd >= 0) {
    const at = sectionEnd + "</section>".length;
    return `${html.slice(0, at)}${block}${html.slice(at)}`;
  }
  const mainOpenEnd = html.indexOf(">", mainAt);
  return `${html.slice(0, mainOpenEnd + 1)}${block}${html.slice(mainOpenEnd + 1)}`;
}

function refreshHome(lang) {
  const file = join(root, "web", lang, "index.html");
  let html = readFileSync(file, "utf8");
  const header = lang === "ja"
    ? `<div class="specializations-head"><p class="eyebrow">SPECIALIZATIONS / ゲーム制作編</p><h2>全部終える必要はない。<br>作りたいゲームから選ぼう。</h2><p>LEVEL 01でUpdateとDrawの役割が分かったら、気になるジャンルへ進めます。難しい演出が出てきた時だけ、各ページから対応するVFX教材へ戻れます。25個から選べない時は、<a href="guides/choose-your-path/">5つの遊び方から選ぶガイド</a>を使ってください。</p></div>`
    : `<div class="specializations-head"><p class="eyebrow">SPECIALIZATIONS / BUILD GAMES</p><h2>You do not have to finish everything.<br>Choose the game you want.</h2><p>Once LEVEL 01 makes the Update/Draw roles clear, you can enter any genre. When presentation gets tricky, each page links back to the relevant VFX lesson. If twenty-five choices feel like too many, use the <a href="guides/choose-your-path/">five-style path chooser</a>.</p></div>`;
  html = html.replace(/<div class="specializations-head">[\s\S]*?<\/div><div class="track-grid">/, `${header}<div class="track-grid">`);

  const counts = new Map();
  for (const item of curriculum.filter((item) => item.group === "track")) {
    const id = item.route.split("/")[1];
    counts.set(id, (counts.get(id) || 0) + 1);
  }
  for (const entry of genrePresentationMap) {
    const patternNames = entry.patterns.map((id) => `${id} ${presentationPatterns[id][lang]}`).join(lang === "ja" ? "・" : ", ");
    const cardCopy = lang === "ja"
      ? `${entry[lang].name}固有のルールを一歩ずつ作ります。演出で迷った時は ${patternNames} を参考にできます。`
      : `Build the rules unique to ${entry[lang].name} one step at a time. Use ${patternNames} only when you need presentation help.`;
    const cardPattern = new RegExp(`<a class="track-card track-${entry.id}"[\\s\\S]*?<\\/a>`);
    const match = html.match(cardPattern);
    if (!match) throw new Error(`missing home card for ${lang}/${entry.id}`);
    const updated = match[0]
      .replace(/<div><p class="eyebrow">[\s\S]*?<\/p>/, `<div><p class="eyebrow">${lang === "ja" ? "好きなジャンルから始める" : "START WITH A GENRE"}</p>`)
      .replace(/<h3>[\s\S]*?<\/h3><p>[\s\S]*?<\/p><strong>[\s\S]*?<\/strong>/,
        `<h3>${entry[lang].name}</h3><p>${cardCopy}</p><strong>${counts.get(entry.id)} STEPS →</strong>`);
    html = html.replace(match[0], updated);
  }

  const current = readFileSync(file, "utf8");
  if (html !== current) {
    if (checkOnly) throw new Error(`genre specialization home copy is stale: ${relative(root, file)}`);
    writeFileSync(file, html);
    changed++;
  }
}

let changed = 0;
let checked = 0;
for (const lang of ["ja", "en"]) {
  const tracksRoot = join(root, "web", lang, "tracks");
  const diskTracks = readdirSync(tracksRoot)
    .filter((name) => name !== "visual-effects" && statSync(join(tracksRoot, name)).isDirectory());
  const missing = diskTracks.filter((id) => !expected.has(id));
  if (missing.length) throw new Error(`missing presentation map entries: ${missing.join(", ")}`);

  for (const entry of genrePresentationMap) {
    const dir = join(tracksRoot, entry.id);
    if (!existsSync(dir)) throw new Error(`missing track directory: ${dir}`);
    for (const file of indexFiles(dir)) {
      checked++;
      const clean = readFileSync(file, "utf8").replace(oldBlock, "");
      const isHub = relative(dir, file) === "index.html";
      const block = isHub ? hubBridge(entry, lang) : lessonBridge(entry, lang);
      const next = insertAfterOpeningSection(clean, block);
      if (next !== readFileSync(file, "utf8")) {
        if (checkOnly) throw new Error(`genre presentation bridge is stale: ${relative(root, file)}`);
        writeFileSync(file, next);
        changed++;
      }
    }
  }
  refreshHome(lang);
}

console.log(checkOnly
  ? `Verified ${checked} genre page presentation bridge(s).`
  : `Injected genre presentation bridges into ${changed}/${checked} page(s).`);
