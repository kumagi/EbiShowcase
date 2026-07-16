#!/usr/bin/env node
/**
 * Build-only overlay for capstone renderer galleries.
 *
 * The replacement changes only the call that starts Ebitengine:
 *
 *   ebiten.RunGame(game) -> renderfreedom.Run(game)
 *
 * Update, Draw, Layout, game state, and every rule remain the original source.
 */
import { existsSync, mkdirSync, readFileSync, readdirSync, rmSync, statSync, writeFileSync } from "node:fs";
import { dirname, join, relative, resolve, sep } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(join(dirname(fileURLToPath(import.meta.url)), ".."));
const cache = join(root, ".cache", "ebi-showcase", "render-freedom-overlay");
const overlayFile = join(cache, "overlay.json");
const rendererImport = '"github.com/kumagi/EbiShowcase/internal/renderfreedom"';

function walk(dir) {
  if (!existsSync(dir)) return [];
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const file = join(dir, entry.name);
    if (entry.isDirectory() && [".git", ".cache", "dist", "node_modules"].includes(entry.name)) return [];
    return entry.isDirectory() ? walk(file) : entry.name.endsWith(".go") ? [file] : [];
  });
}

function addImport(source, file) {
  if (source.includes(rendererImport)) return source;
  if (/import\s*\(\s*\n/.test(source)) {
    return source.replace(/import\s*\(\s*\n/, (match) => `${match}\t${rendererImport}\n`);
  }
  const single = source.match(/import\s+("[^"]+")/);
  if (single) {
    return source.replace(single[0], `import (\n\t${single[1]}\n\t${rendererImport}\n)`);
  }
  throw new Error(`Cannot inject renderer import: ${relative(root, file)}`);
}

function prepare() {
  rmSync(cache, { recursive: true, force: true });
  mkdirSync(cache, { recursive: true });
  const replacements = {};
  for (const file of walk(root)) {
    const rel = relative(root, file).split(sep).join("/");
    if (rel.startsWith(".cache/") || rel.startsWith("dist/") || rel.startsWith("node_modules/") || rel.startsWith("internal/renderfreedom/")) continue;
    const source = readFileSync(file, "utf8");
    if (!source.includes("ebiten.RunGame(")) continue;
    let patched = source.replaceAll("ebiten.RunGame(", "renderfreedom.Run(");
    patched = addImport(patched, file);
    const out = join(cache, "src", rel);
    mkdirSync(dirname(out), { recursive: true });
    writeFileSync(out, patched);
    replacements[resolve(file)] = resolve(out);
  }
  if (Object.keys(replacements).length === 0) throw new Error("No ebiten.RunGame call sites found");
  writeFileSync(overlayFile, JSON.stringify({ Replace: replacements }, null, 2) + "\n");
  process.stdout.write(overlayFile);
}

function capstones() {
  const manifest = JSON.parse(readFileSync(join(root, "web/assets/home-thumbnails/manifest.json"), "utf8"));
  const tracks = manifest.filter((item) => item.kind === "track");
  for (const item of tracks) {
    const track = item.route.split("/").at(-1);
    const dir = join(root, "games", "tracks", track, item.slug);
    if (!existsSync(dir) || !statSync(dir).isDirectory()) {
      throw new Error(`Capstone package missing: ${track}/${item.slug}`);
    }
    process.stdout.write(`${item.slug}\t${dir}\n`);
  }
}

const command = process.argv[2] || "prepare";
if (command === "prepare") prepare();
else if (command === "capstones") capstones();
else throw new Error("usage: render-freedom-overlays.mjs [prepare|capstones]");
