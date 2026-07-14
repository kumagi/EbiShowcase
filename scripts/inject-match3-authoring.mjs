#!/usr/bin/env node
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
const root = new URL("..", import.meta.url).pathname;
const steps = ["grid-swap", "swap-motion", "find-matches", "clear-and-fall", "cascade", "special-pieces", "reef-director", "ebi-match"];
for (const lang of ["ja", "en"]) for (const slug of steps) {
  const ja = lang === "ja", path = join(root, "web", lang, "tracks", "match3", slug, "index.html"), entry = `games/tracks/match3/${slug}/main.go`;
  const block = `<!-- match3-authoring:start -->\n<section class="next-lines core-authoring-links"><p class="eyebrow">YOUR FIRST RULE</p><h3>${ja ? "本物の盤面ルールを1つ足す" : "Add one real board rule"}</h3><p>${ja ? `<code>${entry}</code> を開き、色・消去数・連鎖条件のどれかを1つ変える。<code>go test ./games/tracks/match3/${slug}</code> を実行し、盤面で結果を確かめよう。` : `Open <code>${entry}</code> and change one color, clear count, or cascade condition. Run <code>go test ./games/tracks/match3/${slug}</code>, then check the board result.`}</p></section>\n<!-- match3-authoring:end -->`;
  let html = readFileSync(path, "utf8").replace(/<!-- match3-authoring:start -->[\s\S]*?<!-- match3-authoring:end -->\n?/g, "");
  html = html.replace("</main>", `${block}\n</main>`); writeFileSync(path, html);
}
console.log("Injected Match 3 authoring rules into 8 steps in JA/EN.");
