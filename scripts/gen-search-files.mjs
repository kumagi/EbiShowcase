#!/usr/bin/env node
/** Generate the public sitemap and project-level robots file from web pages. */
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join, relative } from "node:path";
import { SITE_ORIGIN, absoluteURL } from "./site-origin.mjs";

const root = new URL("..", import.meta.url).pathname;
const webRoot = join(root, "web");

function walk(dir, files = []) {
  for (const name of readdirSync(dir).sort()) {
    if (name === "node_modules" || name === "play") continue;
    const full = join(dir, name);
    const stat = statSync(full);
    if (stat.isDirectory()) walk(full, files);
    else if (name === "index.html") files.push(full);
  }
  return files;
}

function isIndexable(html) {
  return !/<meta\s+[^>]*name=["']robots["'][^>]*content=["'][^"']*noindex/i.test(html) &&
    !/<meta\s+[^>]*content=["'][^"']*noindex[^"']*["'][^>]*name=["']robots["']/i.test(html);
}

function pageURL(file) {
  const rel = relative(webRoot, file).replace(/\\/g, "/");
  const directory = rel === "index.html" ? "" : rel.slice(0, -"index.html".length);
  return absoluteURL(directory);
}

function escapeXML(value) {
  return String(value)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&apos;");
}

const urls = [...new Set(
  walk(webRoot)
    .filter((file) => isIndexable(readFileSync(file, "utf8")))
    .map(pageURL),
)].sort();

const sitemap = [
  '<?xml version="1.0" encoding="UTF-8"?>',
  '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">',
  ...urls.map((url) => `  <url><loc>${escapeXML(url)}</loc></url>`),
  "</urlset>",
  "",
].join("\n");

const sitePath = new URL(`${SITE_ORIGIN}/`).pathname;
const robots = [
  "User-agent: *",
  `Allow: ${sitePath}`,
  "",
  `Sitemap: ${absoluteURL("sitemap.xml")}`,
  "",
].join("\n");

writeFileSync(join(webRoot, "sitemap.xml"), sitemap);
writeFileSync(join(webRoot, "robots.txt"), robots);
console.log(`Generated sitemap.xml with ${urls.length} indexable pages and robots.txt.`);
