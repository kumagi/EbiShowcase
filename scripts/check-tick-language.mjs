#!/usr/bin/env node
/** Fail when learner-facing copy calls an Update tick a frame. */
// SPDX-License-Identifier: Apache-2.0
import { readdirSync, readFileSync } from "node:fs";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const root = resolve(join(fileURLToPath(import.meta.url), "..", ".."));
function walk(dir) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    const path = join(dir, entry.name);
    return entry.isDirectory() ? walk(path) : entry.name.endsWith(".html") ? [path] : [];
  });
}

const forbidden = [
  /1フレーム進|1F進|1フレーム後|毎フレーム|\d+フレーム(?:ごと|に1回|で1マス)|次のフレーム|前フレーム|接触フレーム|命中フレーム|フレーム更新|フレーム差|予定フレーム|押したフレーム|待ちフレーム|数フレーム/,
  /Advance (?:1|one) frame|Next frame|every frame|each frame|once per frame|next frame|previous frame|one frame|per frame|frame updates?|frame loop|frame difference|chart frame|pressed frame|wait frames|(?:a )?few frames|ten frames|first frame|single-frame|contact frames?/i,
];
const failures = [];
for (const file of walk(join(root, "web"))) {
  const route = file.slice(join(root, "web").length).replaceAll("\\", "/");
  if (route.includes("/visual-effects/vfx-walk/") || route === "/en/index.html" || route === "/ja/index.html") continue;
  const lines = readFileSync(file, "utf8").split("\n");
  lines.forEach((line, index) => {
    const cadenceCopy = line.replaceAll("Cut one frame at a time from a sheet", "Cut one sprite picture at a time from a sheet");
    if (forbidden.some((pattern) => pattern.test(cadenceCopy))) failures.push(`${route}:${index + 1}`);
  });
}
if (failures.length) {
  console.error(`Update cadence is still described as a frame in ${failures.length} location(s):`);
  failures.slice(0, 60).forEach((failure) => console.error(`  ${failure}`));
  process.exit(1);
}
console.log("Tick terminology check passed: Update cadence is not described as a frame.");
