#!/usr/bin/env node
import { readdirSync, readFileSync, statSync } from "node:fs";
import { join } from "node:path";
const root = new URL("..", import.meta.url).pathname;
const files = [];
function walk(dir) { for (const name of readdirSync(dir)) { const path = join(dir, name); statSync(path).isDirectory() ? walk(path) : path.endsWith(".html") && files.push(path); } }
walk(join(root, "web"));
// Old generic OGP copy promised only “change a value”. Authoring OGP must name
// a deliberate Go rule instead. Body text may still discuss tuning pedagogically.
const legacy = /<meta (?:property="og:description"|name="twitter:description")[^>]*値を変えて[^>]*>/;
const offenders = files.filter((path) => legacy.test(readFileSync(path, "utf8")));
if (offenders.length) { console.error(`Legacy authoring OGP copy in:\n${offenders.join("\n")}`); process.exit(1); }
console.log(`Authoring-copy regression check passed for ${files.length} HTML pages.`);
