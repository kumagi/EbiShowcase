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
    screen.DrawImage(snapshot, nil) // native-size playable input surface
    drawEdgeInset(screen, snapshot)
    drawASCIIGameMock(screen)
}`;
  return `<!-- capstone-renderers:start -->
<section class="capstone-renderers" aria-label="${ja ? "Drawだけを替えるおまけ" : "Draw-only renderer bonus"}">
  <div class="capstone-renderers-copy">
    <p class="eyebrow">BONUS / DRAW ONLY</p>
    <h2>${ja ? "ゲームは同じ。絵だけを、別の作り方にする。" : "Same game. Different ways to make the picture."}</h2>
    <p>${ja
      ? "これは上の完成ゲームと同じパッケージから作ったWASMです。<code>Update</code>と<code>Layout</code>には手を加えず、そのまま呼びます。元ゲームを等倍のまま操作でき、右には同じ瞬間の輪郭線と、文字だけでゲームを描くASCII画面のモックを並べます。ゲームのルールと見た目は別々に設計できます。"
      : "This WASM is built from the same package as the finished game above. It delegates to the original <code>Update</code> and <code>Layout</code> unchanged. The native-size game stays playable beside an edge view of the same instant and a mock game drawn entirely with ASCII characters. Rules and presentation can be designed separately."}</p>
  </div>
  <div class="capstone-renderers-facts">
    <article><b>01</b><h3>UPDATE</h3><p>${ja ? "入力とゲームの事実を進める。3画面で共通。" : "Advances input and game facts once for every view."}</p></article>
    <article><b>02</b><h3>LAYOUT</h3><p>${ja ? "元ゲームの大きさをそのまま返す。" : "Returns the original game's size unchanged."}</p></article>
    <article><b>03</b><h3>DRAW</h3><p>${ja ? "画像・輪郭線・ASCIIなど、描き方を選ぶ。" : "Chooses pixels, edge lines, ASCII, or another presentation."}</p></article>
  </div>
  <div class="capstone-renderers-demo">
    <iframe data-shared-demo="capstone-renderers" src="${playURL}" loading="lazy" title="${ja ? "同じ完成ゲームを3種類のDrawで同時表示" : "The same finished game shown through three Draw styles"}" allow="autoplay; fullscreen"></iframe>
  </div>
  <div class="capstone-renderers-code">
    <div><p class="eyebrow">THE ONLY BOUNDARY</p><h3>${ja ? "UpdateとLayoutは、元のゲームへそのまま渡す。" : "Update and Layout pass straight to the original game."}</h3><p>${ja ? "等倍snapshotが入力と一致するゲーム画面です。輪郭版はKageで隣の画素との差だけを線にします。ASCIIモックは、同じようなゲーム場面を文字だけでも伝えられることを示します。描画はルールを変更しません。" : "The native snapshot is the input-aligned playable surface. The edge inset uses Kage to keep only neighboring-pixel differences. The ASCII mock shows how a comparable game scene can communicate with characters alone. Rendering does not change the rules."}</p></div>
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
