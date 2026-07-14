#!/usr/bin/env node
/**
 * Deterministic OGP input fingerprint and generated-image verifier.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Usage:
 *   node scripts/ogp-cache.mjs fingerprint
 *   node scripts/ogp-cache.mjs verify
 */
import { createHash } from "node:crypto";
import { existsSync, readdirSync, readFileSync, statSync } from "node:fs";
import { join, relative } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const command = process.argv[2] || "fingerprint";

function walk(dir, predicate, files = []) {
  for (const name of readdirSync(dir).sort()) {
    const file = join(dir, name);
    const stat = statSync(file);
    if (stat.isDirectory()) walk(file, predicate, files);
    else if (predicate(file)) files.push(file);
  }
  return files;
}

function rendererInputs() {
  const generator = walk(join(root, "cmd", "gen-og-images"), (file) => file.endsWith(".go"));
  const font = walk(join(root, "internal", "ogfont"), () => true);
  return [
    ...generator,
    ...font,
    join(root, "go.mod"),
    join(root, "go.sum"),
    join(root, "scripts", "inject-ogp.mjs"),
    join(root, "scripts", "site-origin.mjs"),
  ].sort();
}

function fingerprint() {
  const hash = createHash("sha256");
  hash.update(`SITE_ORIGIN=${process.env.SITE_ORIGIN || "https://kumagi.github.io/EbiShowcase"}\n`);
  // inject-ogp writes a timestamp on every run, while the image renderer only
  // consumes the rest of this manifest. Hash that semantic payload rather than
  // the whole HTML tree: generator formatting changes such as extra blank
  // lines must not force hundreds of identical PNG encodes.
  const manifestPath = join(root, "web", "assets", "og", "manifest.json");
  if (!existsSync(manifestPath)) throw new Error("OGP manifest is missing; run inject-ogp first");
  const manifest = JSON.parse(readFileSync(manifestPath, "utf8"));
  delete manifest.generatedAt;
  hash.update("web/assets/og/manifest.json\0");
  hash.update(JSON.stringify(manifest));
  hash.update("\0");
  for (const file of rendererInputs()) {
    hash.update(`${relative(root, file).replace(/\\/g, "/")}\0`);
    hash.update(readFileSync(file));
    hash.update("\0");
  }
  return hash.digest("hex");
}

function isPNG(file) {
  if (!existsSync(file)) return false;
  const bytes = readFileSync(file);
  return bytes.length >= 8 && bytes.subarray(0, 8).equals(Buffer.from([0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a]));
}

function verify() {
  const manifestPath = join(root, "web", "assets", "og", "manifest.json");
  if (!existsSync(manifestPath)) throw new Error("OGP manifest is missing");
  const manifest = JSON.parse(readFileSync(manifestPath, "utf8"));
  if (!Array.isArray(manifest.pages) || manifest.pages.length === 0) throw new Error("OGP manifest has no pages");
  const missing = manifest.pages.filter((page) => !page.image || !isPNG(join(root, "web", page.image)));
  if (missing.length) throw new Error(`Missing or invalid OGP PNGs (${missing.length}): ${missing.slice(0, 8).map((page) => page.image || page.path).join(", ")}`);
  console.log(`OGP cache ready: ${manifest.pages.length} PNGs.`);
}

if (command === "fingerprint") console.log(fingerprint());
else if (command === "verify") verify();
else throw new Error("usage: ogp-cache.mjs {fingerprint|verify}");
