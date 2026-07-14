#!/usr/bin/env node
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const steps = ["tiny-platformer", "moving-platforms", "patrol-enemies", "scrolling-stage", "stage-data", "run-animation", "game-feel", "powerup-adventure"];
for (const lang of ["ja", "en"]) {
  const ja = lang === "ja";
  for (const slug of steps) {
    const path = join(root, "web", lang, "tracks", "platformer", slug, "index.html");
    const entry = `games/tracks/platformer/${slug}/main.go`;
    const action = ja ? "ルールを1つ追加し、動きが変わることを確かめる" : "add one rule and verify that play changes";
    const block = `<!-- platformer-authoring:start -->\n<section class="next-lines core-authoring-links"><p class="eyebrow">YOUR FIRST RULE</p><h3>${ja ? "実在する入口から、1ルールを書く" : "Write one rule from the real entry point"}</h3><p>${ja ? `<code>${entry}</code> を開き、${action}。<code>go test ./games/tracks/platformer/${slug}</code> を実行してから遊び直そう。` : `Open <code>${entry}</code>, ${action}. Run <code>go test ./games/tracks/platformer/${slug}</code>, then play again.`}</p></section>\n<!-- platformer-authoring:end -->`;
    let html = readFileSync(path, "utf8").replace(/<!-- platformer-authoring:start -->[\s\S]*?<!-- platformer-authoring:end -->\n?/g, "");
    html = html.replace("</main>", `${block}\n</main>`);
    writeFileSync(path, html);
  }
}
console.log("Injected platformer authoring rules into 8 steps in JA/EN.");
