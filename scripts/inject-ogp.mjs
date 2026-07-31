#!/usr/bin/env node
/**
 * Inject per-page Open Graph / Twitter Card meta into every content HTML under web/.
 * Also writes web/assets/og/manifest.json for cmd/gen-og-images.
 *
 * Idempotent: replaces a marked OGP block between <!-- ogp:start --> and <!-- ogp:end -->.
 *
 * SITE_ORIGIN (default https://kumagi.github.io/EbiShowcase) must be absolute for crawlers.
 */
import { readdirSync, readFileSync, writeFileSync, mkdirSync, statSync } from "node:fs";
import { join, relative, dirname } from "node:path";
import { SITE_ORIGIN, absoluteURL } from "./site-origin.mjs";

const root = new URL("..", import.meta.url).pathname;
const webRoot = join(root, "web");
// Bump when the rendered card design/font changes so social crawlers do not
// keep showing a cached image at the otherwise stable per-page URL.
const OGP_IMAGE_VERSION = "20260731-gameplay1";
const thumbnailRoot = join(webRoot, "assets", "home-thumbnails");
const thumbnailItems = JSON.parse(readFileSync(join(thumbnailRoot, "manifest.json"), "utf8"));
const thumbnailsBySlug = new Map(thumbnailItems.map((item) => [item.slug, item]));
const thumbnailsByRoute = new Map(thumbnailItems.map((item) => [item.route, item]));

function walk(dir, files = []) {
  for (const name of readdirSync(dir)) {
    if (name === "play" || name === "node_modules") continue;
    const full = join(dir, name);
    const st = statSync(full);
    if (st.isDirectory()) walk(full, files);
    else if (name.endsWith(".html")) files.push(full);
  }
  return files;
}

function decodeEntities(s) {
  return String(s || "")
    .replace(/&amp;/g, "&")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/&nbsp;/g, " ")
    .replace(/<br\s*\/?>/gi, " ")
    .replace(/<[^>]+>/g, "")
    .replace(/\s+/g, " ")
    .trim();
}

function escAttr(s) {
  return String(s)
    .replace(/&/g, "&amp;")
    .replace(/"/g, "&quot;")
    .replace(/</g, "&lt;");
}

function escText(s) {
  return String(s)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}

function normalizeTitle(raw) {
  const source = String(raw || "Ebi Showcase").trim();
  if (!source || /^Ebi Showcase$/i.test(source)) return "Ebi Showcase";
  const name = source
    .replace(/^Ebi Showcase\s*[—–-]\s*/i, "")
    .replace(/\s*(?:\||[—–-])\s*Ebi Showcase\s*$/i, "")
    .trim();
  return name ? `Ebi Showcase – ${name}` : "Ebi Showcase";
}

function cleanDescription(description, lang) {
  // Never let a previous generated meta tail become the next build's source.
  let source = String(description || "");
  const stale = lang === "ja"
    ? [
      /遊べるデモを動かし、短いGoコードを読み、値を変えて[\s\S]*$/u,
      /遊べるデモを動かし、Goで1つルールを足して[\s\S]*$/u,
      /キーボードとタッチの両方で試せます。?[\s\S]*$/u,
    ]
    : [
      /Play the demo, read a short Go example, change a value, and apply[\s\S]*$/iu,
      /Play the demo, add one Go rule, verify it, and apply[\s\S]*$/iu,
      /The demo works with keyboard and touch\.?[\s\S]*$/iu,
    ];
  for (const pattern of stale) source = source.replace(pattern, "");
  source = source.replace(/\s+/g, " ").trim();
  return source;
}

function truncate(value, max) {
  const chars = Array.from(String(value || ""));
  if (chars.length <= max) return chars.join("");
  return chars.slice(0, max - 1).join("").replace(/[、,;:\s]+$/u, "") + "…";
}

function firstSentence(value) {
  const source = String(value || "").trim();
  const sentence = source.match(/^.*?[。.!?](?:\s|$)/u)?.[0] || source;
  return sentence.trim();
}

function pick(html, ...res) {
  for (const re of res) {
    const m = html.match(re);
    if (m?.[1]) return decodeEntities(m[1]);
  }
  return "";
}

function pagePathFromFile(file) {
  let rel = relative(webRoot, file).replace(/\\/g, "/");
  if (rel.endsWith("/index.html")) rel = rel.slice(0, -"/index.html".length);
  else if (rel === "index.html") rel = "";
  else if (rel.endsWith(".html")) rel = rel.slice(0, -".html".length);
  return rel;
}

function ogImageKey(pagePath) {
  if (!pagePath || pagePath === "") return "root";
  return pagePath.replace(/\//g, "--");
}

function counterpartPath(pagePath, lang) {
  if (!pagePath) return lang;
  if (pagePath === "ja" || pagePath === "en") return lang;
  if (pagePath.startsWith("ja/")) return lang + pagePath.slice(2);
  if (pagePath.startsWith("en/")) return lang + pagePath.slice(2);
  return pagePath;
}

function classify(pagePath) {
  if (!pagePath || pagePath === "root" || pagePath === "ja" || pagePath === "en") return "home";
  if (pagePath.includes("/graduation")) return "graduation";
  if (pagePath.includes("/build")) return "build";
  if (pagePath.includes("/labs/")) return "guide";
  if (pagePath.includes("/guides/")) return "guide";
  if (pagePath.includes("/tracks/visual-effects")) return "vfx";
  if (pagePath.includes("/tracks/")) return "track";
  if (pagePath.includes("/games/")) return "core";
  return "site";
}

function extract(html, pagePath) {
  const lang = pick(html, /<html[^>]*\blang="([^"]+)"/i) || (pagePath.startsWith("ja") ? "ja" : "en");
  const title = normalizeTitle(pick(html, /<title>([^<]*)<\/title>/i) || "Ebi Showcase");
  const h1 = pick(html, /<h1[^>]*>([\s\S]*?)<\/h1>/i) || title.split("|")[0].trim();
  const eyebrow = pick(
    html,
    /<p class="eyebrow"[^>]*>([\s\S]*?)<\/p>/i,
    /<div class="lesson-breadcrumb"[\s\S]*?<span>([\s\S]*?)<\/span>/i,
  );
  const concept = pick(
    html,
    /<div class="overview-concept"[\s\S]*?<strong>([\s\S]*?)<\/strong>/i,
    /<div class="lesson-meta"[\s\S]*?<strong>([\s\S]*?)<\/strong>/i,
    /<section class="play[^"]*"[\s\S]*?<h2[^>]*>([\s\S]*?)<\/h2>/i,
    /<section class="play-panel"[\s\S]*?<h2[^>]*>([\s\S]*?)<\/h2>/i,
    /<section class="test-rule-strip"[\s\S]*?<strong>([\s\S]*?)<\/strong>/i,
  );
  let description = pick(
    html,
    /<(?:section|div)\s+class="[^"]*(?:overview-hero|lesson-hero|track-hero|data-hero|test-hero|test-step-hero)[^"]*"[\s\S]*?<h1[^>]*>[\s\S]*?<\/h1>[\s\S]*?<p(?![^>]*\beyebrow\b)[^>]*>([\s\S]*?)<\/p>/i,
    /<p class="lead"[^>]*>([\s\S]*?)<\/p>/i,
    /<p class="lesson-lead"[^>]*>([\s\S]*?)<\/p>/i,
    /<meta\s+name="description"\s+content="([^"]*)"/i,
    /<meta\s+name="description"\s+content='([^']*)'/i,
  );
  description = cleanDescription(description, lang);
  if (!description) {
    description = lang === "ja"
      ? `${h1}を、実際に動くEbitengineの画面とGoコードで学びます。`
      : `Explore ${h1} through a working Ebitengine experience and the Go code behind it.`;
  } else if (Array.from(description).length < 35) {
    const subject = (concept || h1).replace(/[。.!?]+$/u, "").trim();
    const addition = lang === "ja"
      ? `「${subject}」を、実際に動く画面とGoコードで確かめます。`
      : `Try ${subject} in a working game, then see how the Go rule is built.`;
    if (!description.includes(subject)) description = `${description} ${addition}`;
  }
  const pageLabel = title.replace(/^Ebi Showcase\s*[–—-]\s*/i, "").trim();
  if (pageLabel && !description.includes(pageLabel)) description = `${pageLabel} — ${description}`;
  description = truncate(description, 155);
  const hook = concept
    ? (lang === "ja" ? `${concept}を、遊んで解き明かす。` : `Play with ${concept}—then build it.`)
    : firstSentence(description);
  return {
    lang,
    title,
    description,
    h1,
    eyebrow,
    concept,
    hook: truncate(hook, 64),
    kind: classify(pagePath),
  };
}

function unlocalizedPath(pagePath) {
  return pagePath.replace(/^(?:ja|en)(?:\/|$)/, "");
}

function directPlayableSlug(html) {
  for (const match of html.matchAll(/<iframe\b[^>]*>/gi)) {
    const tag = match[0];
    if (!/\blesson-game-frame\b/i.test(tag)) continue;
    const source = tag.match(/\b(?:src|data-game-src)="[^"]*\/play\/([^/"]+)\/?"/i)?.[1];
    if (source) return source;
  }
  return "";
}

const guidePreviewSlugs = new Map([
  ["guides/performance", "space-shooter"],
  ["guides/game-data", "dungeon"],
  ["guides/save", "tap-target"],
  ["guides/setup", "tap-target"],
  ["guides/first-30-minutes", "tap-target"],
  ["guides/choose-your-path", "rpg"],
  ["labs/shader", "vfx-faux-bloom"],
  ["labs/audio", "vfx-spells"],
  ["labs/camera", "platformer"],
  ["graduation/arcade-60", "bullet-hell"],
  ["graduation/exploration-3rooms", "metroidvania"],
  ["graduation/puzzle-3stages", "match3"],
  ["build", "tap-target"],
]);

function resolvePreview(html, pagePath, kind, lang) {
  const route = unlocalizedPath(pagePath);
  const directSlug = directPlayableSlug(html);
  let item = directSlug ? thumbnailsBySlug.get(directSlug) : null;
  let mode = item ? "exact" : "";

  if (!item && route.startsWith("guides/testing/")) {
    const lessonSlug = route.split("/").at(-1);
    item = thumbnailsBySlug.get(lessonSlug);
    if (item) mode = "related";
  }
  if (!item) {
    const track = route.match(/^tracks\/([^/]+)/)?.[1];
    if (track) {
      item = thumbnailsByRoute.get(`tracks/${track}`);
      if (item) mode = route === `tracks/${track}` ? "exact" : "capstone";
    }
  }
  if (!item) {
    const mappedSlug =
      guidePreviewSlugs.get(route) ||
      [...guidePreviewSlugs].find(([prefix]) => route.startsWith(`${prefix}/`))?.[1];
    item = thumbnailsBySlug.get(mappedSlug);
    if (item) mode = "related";
  }
  if (!item) {
    item = thumbnailsByRoute.get(route) || thumbnailsBySlug.get("tap-target");
    mode = thumbnailsByRoute.has(route) ? "exact" : "related";
  }

  const labels = lang === "ja"
    ? { exact: "このページの実ゲーム", capstone: "作っていく完成ゲーム", related: "関連する実ゲーム" }
    : { exact: "REAL GAME ON THIS PAGE", capstone: "THE FINAL GAME YOU'LL BUILD", related: "RELATED REAL GAME" };
  const actions = lang === "ja"
    ? {
      exact: "今すぐ遊べる · Goで仕組みを作る",
      capstone: "このSTEPを遊ぶ · 完成ゲームへつなぐ",
      related: kind === "guide" ? "読んで試す · 自分のゲームへ持ち帰る" : "仕組みを学ぶ · 実ゲームで確かめる",
    }
    : {
      exact: "PLAY NOW · BUILD THE RULE IN GO",
      capstone: "PLAY THIS STEP · BUILD TOWARD THE FINAL",
      related: kind === "guide" ? "READ · TRY · USE IT IN YOUR GAME" : "LEARN THE IDEA · SEE IT IN A REAL GAME",
    };
  if (kind === "home") {
    labels.related = lang === "ja" ? "ブラウザで動くゲーム教材" : "PLAYABLE EBITENGINE LESSONS";
    actions.related = lang === "ja"
      ? "208レッスンから選ぶ · ブラウザですぐ遊ぶ"
      : "CHOOSE A LESSON · PLAY IT IN YOUR BROWSER";
  }
  return {
    preview: `assets/home-thumbnails/${item.file}`,
    previewMode: mode,
    previewLabel: labels[mode],
    action: actions[mode],
  };
}

function buildOgBlock(meta) {
  const locale = meta.lang === "ja" ? "ja_JP" : "en_US";
  const localeAlt = meta.lang === "ja" ? "en_US" : "ja_JP";
  const rows = [
    `  <!-- ogp:start -->`,
    `  <meta property="og:site_name" content="Ebi Showcase">`,
    `  <meta property="og:type" content="website">`,
    `  <meta property="og:title" content="${escAttr(meta.title)}">`,
    `  <meta property="og:description" content="${escAttr(meta.description)}">`,
    `  <meta property="og:url" content="${escAttr(meta.pageURL)}">`,
    `  <meta property="og:image" content="${escAttr(meta.imageURL)}">`,
    `  <meta property="og:image:width" content="1200">`,
    `  <meta property="og:image:height" content="630">`,
    `  <meta property="og:image:alt" content="${escAttr(meta.imageAlt)}">`,
    `  <meta property="og:locale" content="${locale}">`,
    `  <meta property="og:locale:alternate" content="${localeAlt}">`,
    `  <meta name="twitter:card" content="summary_large_image">`,
    `  <meta name="twitter:title" content="${escAttr(meta.title)}">`,
    `  <meta name="twitter:description" content="${escAttr(meta.description)}">`,
    `  <meta name="twitter:image" content="${escAttr(meta.imageURL)}">`,
    `  <meta name="twitter:image:alt" content="${escAttr(meta.imageAlt)}">`,
    `  <!-- ogp:end -->`,
  ];
  return rows.join("\n");
}

function stripOldOgp(html) {
  return html
    .replace(/\n?[ \t]*<!-- ogp:start -->[\s\S]*?<!-- ogp:end -->\n?/g, "\n")
    .replace(/\n?[ \t]*<meta\s+property="og:[^"]+"\s+content="[^"]*"\s*\/?>\n?/gi, "")
    .replace(/\n?[ \t]*<meta\s+name="twitter:[^"]+"\s+content="[^"]*"\s*\/?>\n?/gi, "");
}

function ensureCanonical(html, pageURL) {
  if (/rel="canonical"/i.test(html)) {
    return html.replace(
      /<link\s+rel="canonical"\s+href="[^"]*"\s*\/?>/i,
      `<link rel="canonical" href="${escAttr(pageURL)}">`,
    );
  }
  return html.replace(/<\/title>/i, `</title>\n  <link rel="canonical" href="${escAttr(pageURL)}">`);
}

function ensureTitle(html, title) {
  return html.replace(/<title>[^<]*<\/title>/i, `<title>${escText(title)}</title>`);
}

function ensureDescription(html, description) {
  const tag = `<meta name="description" content="${escAttr(description)}">`;
  if (/<meta\s+name="description"\s+content="[^"]*"\s*\/?\s*>/i.test(html)) {
    return html.replace(/<meta\s+name="description"\s+content="[^"]*"\s*\/?\s*>/i, tag);
  }
  return html.replace(/<meta\s+charset="[^"]*"\s*\/?\s*>/i, (match) => `${match}\n  ${tag}`);
}

function faviconHref(file) {
  return relative(dirname(file), join(webRoot, "assets", "favicon.png")).replace(/\\/g, "/");
}

function ensureFavicon(html, href) {
  const block = [
    "  <!-- favicon:start -->",
    `  <link rel="icon" type="image/png" href="${escAttr(href)}">`,
    `  <link rel="apple-touch-icon" href="${escAttr(href)}">`,
    "  <!-- favicon:end -->",
  ].join("\n");
  const clean = html.replace(/\n?[ \t]*<!-- favicon:start -->[\s\S]*?<!-- favicon:end -->\n?/g, "\n");
  return clean.replace(/<\/head>/i, `${block}\n</head>`);
}

function ensurePagerRelations(html) {
  return html.replace(/<nav\s+class="lesson-pager"[^>]*>([\s\S]*?)<\/nav>/gi, (whole, body) => {
    const nextBody = body.replace(/<a\b([^>]*)>([\s\S]*?)<\/a>/gi, (anchor, attrs, content) => {
      const label = decodeEntities(content).replace(/\s+/g, " ").trim();
      const relation = /(?:→|NEXT|FINAL|COMPLETE|次|完了)/i.test(label) ? "next" : "prev";
      const cleanAttrs = attrs.replace(/\s+rel="[^"]*"/i, "");
      return `<a${cleanAttrs} rel="${relation}">${content}</a>`;
    });
    return whole.replace(body, nextBody);
  });
}

function inject(html, block) {
  let next = stripOldOgp(html);
  if (/<\/title>/i.test(next)) {
    return next.replace(/<\/title>/i, `</title>\n  ${block.replace(/\n/g, "\n  ")}`);
  }
  if (/<head[^>]*>/i.test(next)) {
    return next.replace(/<head[^>]*>/i, (m) => `${m}\n${block}`);
  }
  return block + "\n" + next;
}

const files = walk(webRoot).filter((f) => !f.endsWith("/game.html"));
const manifest = {
  origin: SITE_ORIGIN,
  generatedAt: new Date().toISOString(),
  pages: [],
};

let updated = 0;
for (const file of files) {
  const pagePath = pagePathFromFile(file);
  const html = readFileSync(file, "utf8");
  const info = extract(html, pagePath);
  const preview = resolvePreview(html, pagePath, info.kind, info.lang);
  const key = ogImageKey(pagePath || "root");
  const pageURL = absoluteURL(pagePath ? `${pagePath}/` : "");
  // Root language gate uses trailing path without forcing index
  const imagePath = `assets/og/${key}.png`;
  const imageURL = `${absoluteURL(imagePath)}?v=${OGP_IMAGE_VERSION}`;

  const block = buildOgBlock({
    title: info.title,
    description: info.description,
    pageURL,
    imageURL,
    lang: info.lang,
    imageAlt: `${info.h1 || info.title} — ${preview.previewLabel}`,
  });

  let next = ensureTitle(html, info.title);
  next = ensureDescription(next, info.description);
  next = inject(next, block);
  next = ensureCanonical(next, pageURL);
  next = ensurePagerRelations(next);
  next = ensureFavicon(next, faviconHref(file));
  if (next !== html) {
    writeFileSync(file, next);
    updated++;
  }

  manifest.pages.push({
    file: relative(root, file).replace(/\\/g, "/"),
    path: pagePath || "",
    key,
    lang: info.lang,
    kind: info.kind,
    title: info.title,
    h1: info.h1,
    eyebrow: info.eyebrow,
    description: info.description,
    image: imagePath,
    hook: info.hook,
    action: preview.action,
    preview: preview.preview,
    previewMode: preview.previewMode,
    previewLabel: preview.previewLabel,
  });
}

const ogDir = join(webRoot, "assets", "og");
mkdirSync(ogDir, { recursive: true });
writeFileSync(join(ogDir, "manifest.json"), JSON.stringify(manifest, null, 2) + "\n");
writeFileSync(
  join(ogDir, "README.md"),
  `# OGP images

Generated by \`go run ./cmd/gen-og-images\` from \`manifest.json\`.
Do not hand-edit PNGs. Re-run after \`node scripts/inject-ogp.mjs\`.

Site origin: \`${SITE_ORIGIN}\`
`,
);

console.log(`OGP injected into ${updated}/${files.length} HTML files (origin ${SITE_ORIGIN})`);
console.log(`Manifest: web/assets/og/manifest.json (${manifest.pages.length} pages)`);
