#!/usr/bin/env node
// Keep every Core LEVEL connected to the clone-free authoring path.
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const slugs = ["tap-target", "timing-meter", "catch-stars", "flappy", "pong", "breakout", "snake", "space-shooter", "sokoban", "platformer", "dungeon", "bullet-hell"];
const lateRules = {
  snake: ["food score rule", "add one food-score rule"],
  "space-shooter": ["enemy wave rule", "add one enemy-wave rule"],
  sokoban: ["box-goal rule", "add one box-goal rule"],
  platformer: ["jump or platform rule", "add one jump or platform rule"],
  dungeon: ["room or key rule", "add one room or key rule"],
  "bullet-hell": ["bullet pattern rule", "add one bullet-pattern rule"],
};

for (const lang of ["ja", "en"]) {
  const ja = lang === "ja";
  for (const slug of slugs) {
    const rule = lateRules[slug];
    const ruleLine = rule ? (ja ? `YOUR FIRST RULE: <code>games/core/${slug}/main.go</code> を開き、${rule[0]}を1つ足す。<code>go test ./games/core/${slug}</code> を実行してから遊び直そう。` : `YOUR FIRST RULE: open <code>games/core/${slug}/main.go</code>, ${rule[1]}, then run <code>go test ./games/core/${slug}</code> before playing again.`) : "";
    const block = `<!-- core-authoring-links:start -->\n<section class="next-lines core-authoring-links"><p class="eyebrow">WRITE NEXT</p><h3>${ja ? "遊んだら、1ルールを書いて確かめよう" : "After playing, write and verify one rule"}</h3><p>${ruleLine || (ja ? "空の Go フォルダから書くなら Build Track、ルールをテストするならユニットテスト入門、30分で自分の1本を始めるなら first-30-minutes へ進みます。" : "Use Build Track to write from an empty Go folder, Testing to prove a rule, or first-30-minutes to begin your own game in half an hour.")}</p><p><a href="../../build/">Build Track →</a> · <a href="../../guides/testing/">${ja ? "ユニットテスト入門" : "Testing guide"} →</a> · <a href="../../guides/first-30-minutes/">first-30-minutes →</a></p></section>\n<!-- core-authoring-links:end -->`;
    const file = join(root, "web", lang, "games", slug, "index.html");
    let html = readFileSync(file, "utf8");
    html = html.replace(/<!-- core-authoring-links:start -->[\s\S]*?<!-- core-authoring-links:end -->\n?/, "");
    html = html.replace("</main>", `${block}\n</main>`);
    writeFileSync(file, html);
  }
}
