#!/usr/bin/env node
/** Mark every genre capstone as an intentional quality leap. */
import { readFileSync, readdirSync, statSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const manifest = JSON.parse(readFileSync(join(root, "web/assets/home-thumbnails/manifest.json"), "utf8"));

const ja = {
  "powerup-adventure": ["多層パララックス", "変身と4ステージ", "ゴール演出"],
  "survival-run": ["包囲から始まる戦場", "進化する武器", "ボスの予告動作"],
  "offline-bakery": ["動くベーカリー", "生産ライン", "街の成長"],
  "ebi-quest": ["探索できる世界", "3つの戦闘", "敵の次の行動"],
  "ebi-fighters": ["8コマ攻撃モーション", "後隙を狙うAI", "ガードゲージと反撃"],
  "ebi-merge": ["ガラス容器の物理", "連鎖マージ", "3つの重力課題"],
  "ebi-ascent": ["登頂ルート", "動くカード", "5階のデッキ成長"],
  "ebi-strike": ["軌道予測", "海中バンクショット", "3つの作戦盤"],
  "ebi-match": ["光る海中盤面", "入れ替えアニメ", "連鎖する3ルール"],
  "ebi-blocks": ["ネオン港の盤面", "ライン消去演出", "速度が違う3面"],
  "ebi-craft": ["昼夜の島", "採掘の手応え", "3つの開拓目標"],
  "ebi-monsters": ["地域ごとの生態", "演出つき捕獲", "進化と図鑑"],
  "ebi-chain": ["ライバル対戦", "連鎖の渦", "おじゃま予告"],
  "ebi-maze": ["海底神殿迷路", "記憶する追跡AI", "警戒リング"],
  "ebi-bomber": ["火山機関室", "連鎖する爆風", "導火線の予告"],
  "ebi-reversi": ["大会会場", "3段階のCPU", "石返しアニメ"],
  "ebi-tactics": ["地形で変わる移動", "武器ごとの射程", "3つの任務"],
  "ebi-active-battle": ["動き続ける行動ゲージ", "役割コマンド", "攻撃残像"],
  "ebi-dialogue-stage": ["時間と天候の舞台", "表情つき立ち絵", "4つの結末"],
  "ebi-racing": ["走行中からスタート", "グリップと縁石", "3コースのカップ"],
  "ebi-depths": ["巨大な探索世界", "能力で開く道", "広がる地図"],
  "ebi-raycaster": ["DDAの立体迷路", "奥行きのある敵", "3つの救出任務"],
  "ebi-rhythm": ["ライブステージ", "4レーン譜面", "3曲×2難度"],
  "ebi-defense": ["開戦済みの防衛線", "海底3戦場", "敵の意図表示"],
  "ebi-adventure": ["海底レリック神殿", "4部屋の冒険", "3段階ボス"],
};

const en = {
  "powerup-adventure": ["layered parallax", "four-stage transformation run", "goal celebration"],
  "survival-run": ["an arena already closing in", "evolving weapons", "telegraphed boss attacks"],
  "offline-bakery": ["a living bakery", "animated production lines", "a growing district"],
  "ebi-quest": ["an explorable world", "three encounters", "readable enemy intent"],
  "ebi-fighters": ["eight-frame attacks", "recovery-punishing AI", "guard gauge counterplay"],
  "ebi-merge": ["glass-container physics", "combo merges", "three gravity trials"],
  "ebi-ascent": ["a branching climb", "animated cards", "five-floor deck growth"],
  "ebi-strike": ["trajectory prediction", "undersea bank shots", "three tactical boards"],
  "ebi-match": ["a luminous reef board", "animated swaps", "three chain rules"],
  "ebi-blocks": ["a neon harbor board", "line-clear spectacle", "three speed stages"],
  "ebi-craft": ["day-and-night islands", "responsive mining", "three settlement goals"],
  "ebi-monsters": ["regional habitats", "staged captures", "evolution and research"],
  "ebi-chain": ["a visible rival", "spiraling chain FX", "incoming-garbage tells"],
  "ebi-maze": ["a sunken shrine maze", "memory-driven hunters", "alert rings"],
  "ebi-bomber": ["a volcanic engine room", "chain explosions", "fuse warnings"],
  "ebi-reversi": ["a tournament stage", "three CPU styles", "staggered disc flips"],
  "ebi-tactics": ["terrain-weighted movement", "weapon ranges", "three missions"],
  "ebi-active-battle": ["always-moving gauges", "party roles", "attack trails"],
  "ebi-dialogue-stage": ["time-and-weather stages", "expressive portraits", "four endings"],
  "ebi-racing": ["a rolling start", "grip and curbs", "a three-course cup"],
  "ebi-depths": ["a giant connected world", "ability gates", "an expanding map"],
  "ebi-raycaster": ["a DDA 3D maze", "depth-scaled enemies", "three rescue missions"],
  "ebi-rhythm": ["a live concert stage", "four-lane charts", "three songs in two modes"],
  "ebi-defense": ["an active opening defense", "three seabed battlefields", "enemy intent"],
  "ebi-adventure": ["a painted relic shrine", "a four-room quest", "a three-phase boss"],
};

function pages(dir) {
  return readdirSync(dir).flatMap((name) => {
    const path = join(dir, name);
    return statSync(path).isDirectory() ? pages(path) : name === "index.html" ? [path] : [];
  });
}

function finalPage(lang, track, slug) {
  const needle = `/play/${slug}/`;
  return pages(join(root, "web", lang, "tracks", track)).find((file) => readFileSync(file, "utf8").includes(needle));
}

function bridge(lang, slug) {
  const features = (lang === "ja" ? ja : en)[slug];
  if (!features) throw new Error(`Missing showcase copy for ${slug}`);
  const title = lang === "ja"
    ? `${features[0]}・${features[1]}・${features[2]}を、1本のゲームへ。`
    : `Unite ${features[0]}, ${features[1]}, and ${features[2]}.`;
  const body = lang === "ja"
    ? `この最終ステップでは「${features[0]}」「${features[1]}」「${features[2]}」を、同じUpdateの中で動く一つの遊びにまとめます。まず完成版を遊び、前のステップで作った小さな仕組みがどこで働いているか探してください。`
    : `This capstone combines ${features[0]}, ${features[1]}, and ${features[2]} into one playable Update loop. Play it first, then identify where each earlier lesson now works inside the finished game.`;
  return `<!-- showcase-final:start -->\n<section class="showcase-final-bridge" aria-label="${lang === "ja" ? "ショーケース最終ステップ" : "Showcase final step"}">\n  <div><p class="eyebrow">SHOWCASE FINAL / QUALITY LEAP</p><h2>${title}</h2><p>${body}</p></div>\n  <ul>${features.map((f) => `<li>${f}</li>`).join("")}</ul>\n</section>\n<!-- showcase-final:end -->\n`;
}

let count = 0;
for (const item of manifest.filter((item) => item.kind === "track")) {
  const track = item.route.split("/").at(-1);
  for (const lang of ["ja", "en"]) {
    const file = finalPage(lang, track, item.slug);
    if (!file) throw new Error(`Final page not found: ${lang} ${track} ${item.slug}`);
    let html = readFileSync(file, "utf8");
    html = html.replace(/<!-- showcase-final:start -->[\s\S]*?<!-- showcase-final:end -->\s*/g, "");
    let at = html.search(/<section\b(?=[^>]*\bid="play")[^>]*>/);
    if (at < 0) at = html.search(/<section\b[^>]*class="[^"]*\bplay\b[^"]*"[^>]*>/);
    if (at < 0) throw new Error(`Playable section not found: ${file}`);
    html = html.slice(0, at) + bridge(lang, item.slug) + html.slice(at);
    writeFileSync(file, html);
    count++;
  }
}

console.log(`Injected ${count} bilingual showcase-final bridges.`);
