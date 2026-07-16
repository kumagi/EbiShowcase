#!/usr/bin/env node
/** Add a live Draw-only gallery to every genre capstone page. */
import { readFileSync, readdirSync, statSync, writeFileSync } from "node:fs";
import { dirname, join, relative, sep } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const manifest = JSON.parse(readFileSync(join(root, "web/assets/home-thumbnails/manifest.json"), "utf8"));
const marker = /<!-- capstone-renderers:start -->[\s\S]*?<!-- capstone-renderers:end -->\s*/g;

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

function escapeHTML(value) {
  return value.replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");
}

function section(lang, slug, file) {
  const ja = lang === "ja";
  const playURL = relative(dirname(file), join(root, "web", "play", `${slug}-renderer`)).split(sep).join("/") + "/";
  const code = `func (g *gallery) Update() error {
    return g.inner.Update()
}

func (g *gallery) Layout(w, h int) (int, int) {
    return g.inner.Layout(w, h)
}

func (g *gallery) Draw(screen *ebiten.Image) {
    snapshot := ebiten.NewImage(width, height)
    g.inner.Draw(snapshot)
    drawPolished(screen, snapshot)
    drawWireframe(screen, snapshot)
    drawRectangles(screen, snapshot)
}`;
  return `<!-- capstone-renderers:start -->
<section class="capstone-renderers" aria-label="${ja ? "Drawだけを替えるおまけ" : "Draw-only renderer bonus"}">
  <div class="capstone-renderers-copy">
    <p class="eyebrow">BONUS / DRAW ONLY</p>
    <h2>${ja ? "ゲームは同じ。絵だけを、別の作り方にする。" : "Same game. Different ways to make the picture."}</h2>
    <p>${ja
      ? "これは上の完成ゲームと同じパッケージから作ったWASMです。<code>Update</code>と<code>Layout</code>には手を加えず、そのまま呼びます。1回のDrawが作った同じ瞬間を、通常、輪郭線、矩形モザイクへ同時に投影します。位置、得点、敵、当たり判定、時間は1つしかありません。"
      : "This WASM is built from the same package as the finished game above. It delegates to the original <code>Update</code> and <code>Layout</code> unchanged. One Draw snapshot is projected simultaneously as polished art, edge lines, and a rectangle mosaic. Position, score, enemies, collision, and time exist only once."}</p>
  </div>
  <div class="capstone-renderers-facts">
    <article><b>01</b><h3>UPDATE</h3><p>${ja ? "入力とゲームの事実を進める。3画面で共通。" : "Advances input and game facts once for every view."}</p></article>
    <article><b>02</b><h3>LAYOUT</h3><p>${ja ? "元ゲームの大きさをそのまま返す。" : "Returns the original game's size unchanged."}</p></article>
    <article><b>03</b><h3>DRAW</h3><p>${ja ? "同じsnapshotの見せ方だけを替える。" : "Only changes how one snapshot is presented."}</p></article>
  </div>
  <div class="capstone-renderers-demo">
    <iframe data-shared-demo="capstone-renderers" src="${playURL}" loading="lazy" title="${ja ? "同じ完成ゲームを3種類のDrawで同時表示" : "The same finished game shown through three Draw styles"}" allow="autoplay; fullscreen"></iframe>
  </div>
  <div class="capstone-renderers-code">
    <div><p class="eyebrow">THE ONLY BOUNDARY</p><h3>${ja ? "UpdateとLayoutは、元のゲームへそのまま渡す。" : "Update and Layout pass straight to the original game."}</h3><p>${ja ? "矩形版はsnapshotを小さな格子へ縮め、nearestで戻します。ワイヤ版は隣の画素との差をKageシェーダーで線にします。どちらもゲーム状態を読み書きしません。" : "The rectangle view shrinks the snapshot to a small grid and restores it with nearest filtering. The wire view turns neighboring-pixel differences into lines with a Kage shader. Neither reads or writes game state."}</p></div>
    <pre><code>${escapeHTML(code)}</code></pre>
  </div>
</section>
<!-- capstone-renderers:end -->
`;
}

let count = 0;
for (const item of manifest.filter((item) => item.kind === "track")) {
  const track = item.route.split("/").at(-1);
  for (const lang of ["ja", "en"]) {
    const file = finalPage(lang, track, item.slug);
    if (!file) throw new Error(`Final page not found: ${lang} ${track}/${item.slug}`);
    let html = readFileSync(file, "utf8").replace(marker, "");
    const anchor = "</main>";
    if (!html.includes(anchor)) throw new Error(`No </main>: ${file}`);
    html = html.replace(anchor, `${section(lang, item.slug, file)}${anchor}`);
    writeFileSync(file, html);
    count++;
  }
}

console.log(`Injected Draw-only renderer galleries into ${count} capstone pages.`);
