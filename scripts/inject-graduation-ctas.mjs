#!/usr/bin/env node
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const tracks = ["rhythm", "raycaster", "tower-defense", "reversi", "topdown-adventure", "platformer"];
for (const lang of ["ja", "en"]) {
  const ja = lang === "ja";
  const block = `<!-- graduation-cta:start -->\n<section class="architecture-promo setup-promo graduation-cta"><div><p class="eyebrow">MAKE / YOUR OWN GAME</p><h2>${ja ? "この型で、\n自分の1本を作ろう。" : "Use this pattern\nto make your own game."}</h2><p>${ja ? "遊んだ仕組みを一つ選び、ルールとテストを自分で足す卒業制作へ進みます。" : "Choose one played system, then add your own rule and test in a graduation project."}</p></div><a href="../../graduation/"><span>MAKE NEXT</span><strong>${ja ? "卒業制作を選ぶ" : "Choose a graduation project"}</strong><b>→</b></a></section>\n<!-- graduation-cta:end -->`;
  for (const track of tracks) {
    const path = join(root, "web", lang, "tracks", track, "index.html");
    let html = readFileSync(path, "utf8").replace(/<!-- graduation-cta:start -->[\s\S]*?<!-- graduation-cta:end -->\n?/g, "");
    html = html.replace("</main>", `${block}\n</main>`);
    writeFileSync(path, html);
  }
}
console.log("Injected graduation CTAs into 6 track hubs in JA/EN.");
