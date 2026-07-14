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
const OGP_IMAGE_VERSION = "20260713-font1";

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

function enrichDescription(description, lang) {
  const source = String(description || "").replace(/\s+/g, " ").trim();
  const tail = lang === "ja"
    ? "遊べるデモを動かし、短いGoコードを読み、値を変えて結果を確かめながら、自分のEbitengineゲームへ応用します。"
    : "Play the demo, read a short Go example, change a value, and apply the idea to your own Ebitengine game step by step.";
  const extra = lang === "ja"
    ? "キーボードとタッチの両方で試せます。"
    : "The demo works with keyboard and touch.";
  let result = source || (lang === "ja" ? "遊べるミニゲームで学ぶEbitengineのレッスン。" : "Learn Ebitengine through a playable mini-game lesson.");
  if (result.length < 120) result = `${result} ${tail}`;
  if (result.length < 120) result = `${result} ${extra}`;
  return result.slice(0, 160);
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
  if (pagePath.includes("/guides/")) return "guide";
  if (pagePath.includes("/tracks/visual-effects")) return "vfx";
  if (pagePath.includes("/tracks/")) return "track";
  if (pagePath.includes("/games/")) return "core";
  return "site";
}

function extract(html, pagePath) {
  const lang = pick(html, /<html[^>]*\blang="([^"]+)"/i) || (pagePath.startsWith("ja") ? "ja" : "en");
  const title = normalizeTitle(pick(html, /<title>([^<]*)<\/title>/i) || "Ebi Showcase");
  let description = pick(html, /<meta\s+name="description"\s+content="([^"]*)"/i);
  if (!description) {
    description = pick(
      html,
      /<p class="lead"[^>]*>([\s\S]*?)<\/p>/i,
      /<section class="overview-hero"[\s\S]*?<h1[^>]*>[\s\S]*?<\/h1>\s*<p>([\s\S]*?)<\/p>/i,
      /class="[^"]*track-hero[^"]*"[\s\S]*?<h1[^>]*>[\s\S]*?<\/h1>\s*<p>([\s\S]*?)<\/p>/i,
      /<p class="lesson-lead"[^>]*>([\s\S]*?)<\/p>/i,
      /<meta\s+name="description"\s+content='([^']*)'/i,
    );
  }
  if (!description) {
    description =
      lang === "ja"
        ? "遊べるミニゲームで学ぶ Ebitengine ショーケース。"
        : "Learn Ebitengine through playable mini games.";
  }
  description = enrichDescription(description, lang);
  const h1 = pick(html, /<h1[^>]*>([\s\S]*?)<\/h1>/i) || title.split("|")[0].trim();
  const eyebrow = pick(
    html,
    /<p class="eyebrow"[^>]*>([\s\S]*?)<\/p>/i,
    /<div class="lesson-breadcrumb"[\s\S]*?<span>([\s\S]*?)<\/span>/i,
  );
  return { lang, title, description, h1, eyebrow, kind: classify(pagePath) };
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
    imageAlt: info.h1 || info.title,
  });

  let next = ensureTitle(html, info.title);
  next = ensureDescription(next, info.description);
  next = inject(next, block);
  next = ensureCanonical(next, pageURL);
  next = ensurePagerRelations(next);
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
