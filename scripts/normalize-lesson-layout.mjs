#!/usr/bin/env node
/**
 * Convert alternate lesson skeletons (lesson-hero / play-panel / game-stage / …)
 * into the standard overview-page layout used by the rest of the curriculum.
 *
 * Usage:
 *   node scripts/normalize-lesson-layout.mjs
 *   node scripts/normalize-lesson-layout.mjs web/ja/tracks/bomb-maze/timed-bomb/index.html
 */
import { readFileSync, writeFileSync, readdirSync, statSync } from "node:fs";
import { dirname, join, relative } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");

const DEFAULT_TARGETS = [
  "web/ja/tracks/bomb-maze/timed-bomb/index.html",
  "web/ja/tracks/bomb-maze/cross-blast/index.html",
  "web/ja/tracks/bomb-maze/breakable-walls/index.html",
  "web/ja/tracks/bomb-maze/chain-explosion/index.html",
  "web/ja/tracks/bomb-maze/escape-ai/index.html",
  "web/ja/tracks/bomb-maze/ebi-bomber/index.html",
  "web/en/tracks/bomb-maze/timed-bomb/index.html",
  "web/en/tracks/bomb-maze/cross-blast/index.html",
  "web/en/tracks/bomb-maze/breakable-walls/index.html",
  "web/en/tracks/bomb-maze/chain-explosion/index.html",
  "web/en/tracks/bomb-maze/escape-ai/index.html",
  "web/en/tracks/bomb-maze/ebi-bomber/index.html",
  "web/ja/tracks/maze-chase/buffered-turn/index.html",
  "web/ja/tracks/maze-chase/junction-ai/index.html",
  "web/en/tracks/maze-chase/buffered-turn/index.html",
  "web/en/tracks/maze-chase/junction-ai/index.html",
  "web/ja/tracks/falling-pairs/chain-score/index.html",
  "web/en/tracks/falling-pairs/chain-score/index.html",
];

function grab(html, re) {
  const m = html.match(re);
  return m ? m[1].trim() : "";
}

function grabAll(html, re) {
  return [...html.matchAll(re)].map((m) => m[1].trim());
}

function first(html, ...res) {
  for (const re of res) {
    const v = grab(html, re);
    if (v) return v;
  }
  return "";
}

function unwrapArticles(block) {
  if (!block) return "";
  // Convert <article><b>01</b>… to concept-number style when needed.
  return block.replace(/<article>\s*<b>(\d+)<\/b>/g, (_, n) => {
    const num = String(Number(n));
    return `<article><span class="concept-number">${num}</span>`;
  });
}

function normalizeConcepts(html) {
  const block = first(
    html,
    /<div class="concept-row">([\s\S]*?)<\/div>\s*(?:<\/section>|<section|<div class="motion-lab"|<section class="motion-lab")/,
    /<div class="concept-row">([\s\S]*?)<\/div>/,
  );
  if (!block) return "";
  return `<div class="concept-row">\n${unwrapArticles(block)}\n</div>`;
}

function normalizeMotionLab(html) {
  // Prefer section.motion-lab, else div.motion-lab
  const section = html.match(/<section class="motion-lab"([\s\S]*?)<\/section>/);
  if (section) {
    let inner = section[0]
      .replace(/^<section /, "<div ")
      .replace(/<\/section>$/, "</div>");
    // Drop duplicate trailing formula paragraph that polish left behind
    return cleanLab(inner);
  }
  const div = html.match(/<div class="motion-lab"([\s\S]*?)<\/div>\s*(?:<div class="formula"|<\/section>|<section class="code-lesson"|<div class="formula")/);
  if (div) {
    // match may be too greedy; find by stack
    const start = html.indexOf('<div class="motion-lab"');
    if (start < 0) return "";
    const end = findMatchingClose(html, start);
    if (end < 0) return "";
    return cleanLab(html.slice(start, end));
  }
  // Fallback: open tag to next formula/code-lesson
  const start = html.search(/<(?:div|section) class="motion-lab"/);
  if (start < 0) return "";
  const end = findMatchingClose(html, start);
  if (end < 0) return "";
  return cleanLab(html.slice(start, end).replace(/^<section /, "<div ").replace(/<\/section>$/, "</div>"));
}

function findMatchingClose(html, start) {
  const openMatch = html.slice(start).match(/^<(div|section)\b[^>]*>/);
  if (!openMatch) return -1;
  const tag = openMatch[1];
  let i = start + openMatch[0].length;
  let depth = 1;
  const openRe = new RegExp(`<${tag}\\b[^>]*>`, "g");
  const closeRe = new RegExp(`</${tag}>`, "g");
  while (i < html.length && depth > 0) {
    openRe.lastIndex = i;
    closeRe.lastIndex = i;
    const o = openRe.exec(html);
    const c = closeRe.exec(html);
    if (!c) return -1;
    if (o && o.index < c.index) {
      depth++;
      i = o.index + o[0].length;
    } else {
      depth--;
      i = c.index + c[0].length;
      if (depth === 0) return i;
    }
  }
  return -1;
}

function cleanLab(lab) {
  // Ensure it's a div
  lab = lab.replace(/^<section\b/, "<div").replace(/<\/section>\s*$/, "</div>");
  return lab.trim();
}

function normalizeFormula(html) {
  const start = html.indexOf('<div class="formula">');
  if (start < 0) return "";
  const end = findMatchingClose(html, start);
  if (end < 0) return "";
  let block = html.slice(start, end);
  // Drop an immediate duplicate <p>…</p></div> that polish sometimes left after the real formula
  const after = html.slice(end, end + 400);
  const dup = after.match(/^(\s*<p>([\s\S]*?)<\/p>)\s*<\/div>/);
  if (dup && block.includes(`<p>${dup[2]}</p>`)) {
    // the extra closing </div> was orphaned; ignore it by only emitting our matched block
  }
  // Deduplicate last two identical paragraphs inside the formula
  const paras = [...block.matchAll(/<p>([\s\S]*?)<\/p>/g)].map((x) => x[1]);
  if (paras.length >= 2 && paras[paras.length - 1] === paras[paras.length - 2]) {
    block = block.replace(/<p>([\s\S]*?)<\/p>(\s*)<\/div>\s*$/, "$2</div>");
  }
  return block.trim();
}

function normalizeCodeLesson(html) {
  // May be <section class="code-lesson"> or <div class="code-lesson">
  const start = html.search(/<(?:section|div) class="code-lesson"/);
  if (start < 0) return "";
  const end = findMatchingClose(html, start);
  if (end < 0) return "";
  let block = html.slice(start, end);
  block = block.replace(/^<section\b/, "<div").replace(/<\/section>$/, "</div>");
  // Downgrade inner h2 to h3 when present (overview style uses h3)
  block = block.replace(/<h2>/g, "<h3>").replace(/<\/h2>/g, "</h3>");
  return block;
}

function normalizeWhyGrid(html) {
  const start = html.search(/<(?:section|div) class="why-grid"/);
  if (start < 0) return "";
  const end = findMatchingClose(html, start);
  if (end < 0) return "";
  let block = html.slice(start, end);
  block = block.replace(/^<section\b/, "<div").replace(/<\/section>$/, "</div>");
  return block;
}

function stepLabelFromBreadcrumb(crumb) {
  const m = crumb.match(/STEP\s+(\d+)\s*\/\s*(\d+)/i);
  if (m) return `STEP ${m[1]}`;
  return "LESSON";
}

function normalizePage(relPath) {
  const abs = join(root, relPath);
  const html = readFileSync(abs, "utf8");

  if (html.includes('body class="overview-page"') && html.includes("overview-hero") && html.includes("lesson-play")) {
    // Already standard — still scrub duplicate formula tails if any
    if (!html.includes('class="lesson-hero"') && !html.includes('class="play-panel"') && !html.includes('class="game-stage"')) {
      return { relPath, skipped: true, reason: "already overview" };
    }
  }

  const lang = grab(html, /<html[^>]*\blang="([^"]+)"/) || (relPath.includes("/en/") ? "en" : "ja");
  const title = grab(html, /<title>([\s\S]*?)<\/title>/);
  const description = grab(html, /<meta name="description" content="([^"]*)"/);
  const langHref = first(
    html,
    /<a class="lang"[^>]*href="([^"]+)"/,
    /href="(\.\.\/\.\.\/\.\.\/\.\.\/(?:ja|en)\/[^"]+)"/,
  );
  const breadcrumb = grab(html, /<div class="lesson-breadcrumb">([\s\S]*?)<\/div>/);
  const h1 = grab(html, /<h1>([\s\S]*?)<\/h1>/);

  // Hero lead: first <p> after h1 inside hero / deep-dive start
  let heroLead = first(
    html,
    /<h1>[^<]*<\/h1>\s*<p>([\s\S]*?)<\/p>/,
    /<section class="lesson-hero">[\s\S]*?<h1>[^<]*<\/h1>\s*<p>([\s\S]*?)<\/p>/,
  );

  const stars = first(
    html,
    /<div class="lesson-meta">[\s\S]*?<span>([^<]*★[^<]*)<\/span>/,
    /DIFFICULTY\s*(★[★☆]*)/,
    /(★{1,5}☆{0,5})/,
  ) || "★★★☆☆";
  const starClean = stars.replace(/DIFFICULTY\s*/i, "").replace(/\s*PLAYABLE.*/i, "").trim();
  const eyebrow = `${starClean.includes("★") ? starClean : "★★★☆☆"} / PLAYABLE`;

  // Concept: prefer overview-concept, else the topic eyebrow in lesson-hero (not PLAYABLE status).
  let concept = first(
    html,
    /<div class="overview-concept">[\s\S]*?<strong>([\s\S]*?)<\/strong>/,
  );
  if (!concept) {
    const heroEye = first(
      html,
      /<section class="lesson-hero">[\s\S]*?<p class="eyebrow">([\s\S]*?)<\/p>/,
    );
    if (heroEye && !/^playable$/i.test(heroEye) && !heroEye.includes("★")) {
      concept = heroEye;
    }
  }

  // Play section copy
  const playH2 = first(
    html,
    /<section class="play-panel"[\s\S]*?<h2>([\s\S]*?)<\/h2>/,
    /<section class="game-stage"[\s\S]*?<h2>([\s\S]*?)<\/h2>/,
    /<div class="play-copy">[\s\S]*?<h2>([\s\S]*?)<\/h2>/,
    /<div class="game-copy">[\s\S]*?<h2>([\s\S]*?)<\/h2>/,
  );
  const playP = first(
    html,
    /<section class="play-panel"[\s\S]*?<h2>[^<]*<\/h2>\s*<p>([\s\S]*?)<\/p>/,
    /<section class="game-stage"[\s\S]*?<h2>[^<]*<\/h2>\s*<p>([\s\S]*?)<\/p>/,
    /<div class="play-copy">[\s\S]*?<h2>[^<]*<\/h2>\s*<p>([\s\S]*?)<\/p>/,
    /<div class="game-copy">[\s\S]*?<h2>[^<]*<\/h2>\s*<p>([\s\S]*?)<\/p>/,
  );
  const iframeSrc = first(
    html,
    /<iframe[^>]*\bsrc="([^"]+)"/,
  );
  const iframeTitle = first(
    html,
    /<iframe[^>]*\btitle="([^"]*)"/,
  ) || h1;

  // Deep dive
  const diveEyebrow = first(
    html,
    /<(?:section class="lesson-section"|section class="deep-dive"|div class="lesson-intro")[^>]*>[\s\S]*?<p class="eyebrow">([\s\S]*?)<\/p>/,
    /<p class="eyebrow">(DEEP DIVE[\s\S]*?)<\/p>/,
  ) || "DEEP DIVE";
  const diveH2 = first(
    html,
    /<(?:section class="lesson-section"|section class="deep-dive")[\s\S]*?<h2>([\s\S]*?)<\/h2>/,
    /<div class="lesson-intro">[\s\S]*?<h2>([\s\S]*?)<\/h2>/,
  );
  const diveLead = first(
    html,
    /<(?:section class="lesson-section"|section class="deep-dive")[\s\S]*?<h2>[^<]*<\/h2>\s*<p>([\s\S]*?)<\/p>/,
    /<p class="lesson-lead">([\s\S]*?)<\/p>/,
  );

  const concepts = normalizeConcepts(html);
  const motionLab = normalizeMotionLab(html);
  let formula = normalizeFormula(html);
  // Strip duplicate paragraph that appears immediately after formula close in source
  const dupAfter = html.match(/<\/div>\s*<p>([\s\S]*?)<\/p>\s*<\/div>\s*<(?:section|div) class="code-lesson"/);
  if (dupAfter && formula.includes(dupAfter[1])) {
    // already handled inside formula usually
  }

  const codeLesson = normalizeCodeLesson(html);
  const whyGrid = normalizeWhyGrid(html);
  const pager = grab(html, /(<nav class="lesson-pager">[\s\S]*?<\/nav>)/);
  const feedback = grab(html, /(<section class="feedback-section"[\s\S]*?<\/section>)/);

  const trackHref = lang === "ja" ? "../" : "../";
  const trackLabel = lang === "ja" ? "コース" : "PATH";
  const howLabel = "HOW IT WORKS";
  const langLabel = lang === "ja" ? "EN" : "日本語";
  const langAttr = lang === "ja" ? ' lang="en" data-language="en"' : "";
  const playEyebrow = "PLAYABLE / WEBASSEMBLY";
  const stepBadge = stepLabelFromBreadcrumb(breadcrumb);
  const liveBadge = "● LIVE / GO + WASM";

  const conceptBlock = concept
    ? `<div class="overview-concept"><small>${lang === "ja" ? "このステップで学ぶこと" : "In this step"}</small><strong>${concept}</strong></div>`
    : "";

  const castNote = lang === "ja"
    ? `<p class="cast-note">動きの計算は <a href="../../../games/tap-target/#basics">LEVEL 01</a> の状態境界に沿ってUpdateへ書き、Drawはその結果を自由に投影します。</p>`
    : `<p class="cast-note">Frame math belongs in Update's state boundary from <a href="../../../games/tap-target/#basics">LEVEL 01</a>; Draw may project it in any style.</p>`;

  const bridge = lang === "ja"
    ? `<p class="curriculum-bridge">共通基礎の続きです。ゲームの1コマが初めてなら、先に<a href="../../../games/tap-target/#basics">LEVEL 01 の Update / Draw</a>を見てください。</p>`
    : `<p class="curriculum-bridge">This continues the shared basics. If the frame loop is new, start with <a href="../../../games/tap-target/#basics">LEVEL 01 Update / Draw</a>.</p>`;

  const descMeta = description
    ? `  <meta name="description" content="${description}">\n`
    : "";

  const out = `<!doctype html>
<html lang="${lang}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
${descMeta}  <title>${title}</title>
  <link rel="stylesheet" href="../../../../style.css">
</head>
<body class="overview-page">
  <header class="nav"><a class="brand" href="../../../"><span>EBI</span> SHOWCASE</a><nav><a href="${trackHref}">${trackLabel}</a><a href="#learn">${howLabel}</a><a class="lang" href="${langHref}"${langAttr}>${langLabel}</a></nav></header>
  <main class="overview-main">
    <div class="lesson-breadcrumb">${breadcrumb}</div>

    <section class="overview-hero">
      <p class="eyebrow">${eyebrow}</p>
      <h1>${h1}</h1>
      <p>${heroLead}</p>
      ${conceptBlock}
      ${castNote}
    </section>
    ${bridge}

    <section class="play lesson-play" id="play">
      <div class="section-head"><div><p class="eyebrow">${playEyebrow}</p><h2>${playH2}</h2></div><p>${playP}</p></div>
      <div class="game-shell"><div class="game-top"><span>${liveBadge}</span><span>${stepBadge}</span></div><div id="game-wrap"><iframe class="lesson-game-frame" title="${iframeTitle}" src="${iframeSrc}" allow="autoplay; fullscreen"></iframe></div></div>
    </section>

    <section class="physics" id="learn">
      <div class="lesson-intro">
        <p class="eyebrow">${diveEyebrow}</p>
        <h2>${diveH2}</h2>
        <p class="lesson-lead">${diveLead}</p>
      </div>

      ${concepts}

      ${motionLab}

      ${formula}

      ${codeLesson}

      ${whyGrid}
    </section>

    ${pager}

    ${feedback}
  </main>
  <footer><div class="brand"><span>EBI</span> SHOWCASE</div><p>Made with Go + Ebitengine.</p><a href="https://github.com/kumagi/EbiShowcase">VIEW SOURCE ↗</a></footer>
  <script src="../../../../learn.js"></script>
</body>
</html>
`;

  // Sanity: required pieces
  const missing = [];
  if (!h1) missing.push("h1");
  if (!iframeSrc) missing.push("iframe");
  if (!playH2) missing.push("playH2");
  if (!diveH2) missing.push("diveH2");
  if (!pager) missing.push("pager");
  if (!feedback) missing.push("feedback");
  if (missing.length) {
    return { relPath, error: `missing ${missing.join(",")}` };
  }

  writeFileSync(abs, out);
  return { relPath, ok: true, bytes: out.length };
}

const args = process.argv.slice(2);
const targets = args.length ? args.map((a) => relative(root, a.startsWith("/") ? a : join(process.cwd(), a))) : DEFAULT_TARGETS;

let ok = 0;
let fail = 0;
let skip = 0;
for (const t of targets) {
  const r = normalizePage(t);
  if (r.error) {
    console.error("FAIL", r.relPath, r.error);
    fail++;
  } else if (r.skipped) {
    console.log("skip", r.relPath, r.reason);
    skip++;
  } else {
    console.log("ok  ", r.relPath, r.bytes);
    ok++;
  }
}
console.log(`done ok=${ok} skip=${skip} fail=${fail}`);
process.exit(fail ? 1 : 0);
