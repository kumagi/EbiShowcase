#!/usr/bin/env node
/**
 * Generate the Visual Effects Lab hub + step pages (JA and EN).
 *
 * These pages sit between the 12 core lessons and the genre tracks. They teach
 * Ebitengine's drawing pipeline (GeoM, ColorScale, Blend, SubImage, particles)
 * through hands-on toys. Re-running is safe: it overwrites the generated files.
 *
 * Feedback forms and the shared skeleton match the rest of the site; the build
 * still runs embed-lesson-sources and insert-feedback-form afterwards.
 *
 * Usage: node scripts/gen-visual-effects.mjs
 */
import { mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { advancedLessons } from "./vfx-advanced-lessons.mjs";
import { advancedLessonsPart2 } from "./vfx-advanced-lessons-part2.mjs";
import { illusionLessons } from "./vfx-illusion-lessons.mjs";

const root = new URL("..", import.meta.url).pathname;
const track = "visual-effects";

const esc = (s) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

const hub = {
  letter: "VFX",
  ja: {
    title: "ビジュアルエフェクト工房",
    eyebrow: "花火・光・歩き——見た目の道具箱",
    lead: "ここは「ゲームのルール」ではなく「見た目のかっこよさ」を練習する工房です。スタンプ・回転・色・光・歩きなど、絵を動かす道具を1つずつ触ります。後半では、すでに遊んだ LEVEL のミニゲームにキラキラを足します。全部やらなくても大丈夫——気になる見出しからでOKです。",
    pathTitle: "基本編と応用編。",
    pathLead: "基本編：画面上のスライダーで「置く・回す・染める・光らせる」を体感。応用編：ルール（得点や当たり）と演出（花火）を別の箱に入れる練習。まずはスタンプからで十分です。",
    breadcrumb: "← 全コース",
    course: "コース",
  },
  en: {
    title: "Visual Effects Lab",
    eyebrow: "Sparks, light, walk cycles—a look toolbox",
    lead: "This lab practices “how it looks,” not game rules. Touch one drawing tool at a time—stamp, rotate, tint, glow, walk. Later chapters add sparkle to the core LEVEL mini games you already played. You don’t need every step—pick a title that looks fun.",
    pathTitle: "Basics, then Advanced.",
    pathLead: "Basics: feel place/spin/tint/glow with on-page sliders. Advanced: practice keeping score/hits in one box and fireworks in another. Starting at Stamp is enough.",
    breadcrumb: "← ALL PATHS",
    course: "PATHS",
  },
};

const basicLessons = [
  {
    slug: "vfx-stamp",
    tier: "basic",
    step: "01",
    stars: "★☆☆☆☆",
    labKind: "translate",
    concept: { ja: "GeoM.Translate", en: "GeoM.Translate" },
    hubDesc: {
      ja: "1枚の絵を、好きな場所にスタンプします。描画の出発点です。",
      en: "Stamp one image anywhere. The starting point of all drawing.",
    },
    ja: {
      navConcept: "1枚の絵を置く",
      title: "スタンプと移動",
      lead: "画像を1枚、画面の好きな場所に置きます。位置を決めるのは DrawImageOptions の GeoM.Translate。エフェクトはすべて、この「置く」から始まります。",
      deepEyebrow: "DEEP DIVE / DRAW ONE IMAGE",
      deepH: "絵はどうやって<br>その場所に出る？",
      deepLead: "Ebitengine では、描く絵に「変換」を持たせてから screen.DrawImage します。いちばん基本の変換が平行移動 Translate。同じ絵を、Translate の数字を変えるだけで何個でも置けます。",
      concepts: [
        { h: "オプション", p: "描き方をまとめる箱を用意します。", code: "&op{}" },
        { h: "移動", p: "左上を (x, y) までずらします。", code: "GeoM.Translate" },
        { h: "描画", p: "その設定で画面に絵を出します。", code: "DrawImage" },
      ],
      lab: {
        eyebrow: "TRY IT / TRANSLATE",
        title: "タップした場所に置こう",
        p: "盤面をタップすると、その座標へ絵が飛びます。表示される (x, y) が、そのまま GeoM.Translate に渡す数字です。",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "同じ絵を好きな場所へ", p: "オプションを作って Translate で位置を決め、その設定で描くだけです。" },
      whys: [
        { eyebrow: "WHY OPTIONS?", h: "描き方を持ち運ぶ", p: "位置や色や合成を1つの箱にまとめて渡せます。" },
        { eyebrow: "WHY TRANSLATE?", h: "位置は変換で決まる", p: "絵そのものは動かず、置く場所だけを変えます。" },
        { eyebrow: "TRY NEXT", h: "たくさん置こう", p: "ループで Translate を変えれば、星空も並木も作れます。" },
      ],
    },
    en: {
      navConcept: "Place one image",
      title: "Stamp & Move",
      lead: "Put one image anywhere on screen. Its position comes from GeoM.Translate on DrawImageOptions. Every effect starts from this act of placing.",
      deepEyebrow: "DEEP DIVE / DRAW ONE IMAGE",
      deepH: "How does a picture<br>reach that spot?",
      deepLead: "In Ebitengine you give a drawing a transform, then call screen.DrawImage. The most basic transform is Translate. You can stamp the same image many times just by changing the Translate numbers.",
      concepts: [
        { h: "Options", p: "A small box that holds how to draw.", code: "&op{}" },
        { h: "Translate", p: "Move the top-left to (x, y).", code: "GeoM.Translate" },
        { h: "Draw", p: "Paint the image with those settings.", code: "DrawImage" },
      ],
      lab: {
        eyebrow: "TRY IT / TRANSLATE",
        title: "Place it where you tap",
        p: "Tap the board and the image jumps to that spot. The (x, y) you see is exactly what you pass to GeoM.Translate.",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "One image, any place", p: "Make an options box, set the position with Translate, then draw." },
      whys: [
        { eyebrow: "WHY OPTIONS?", h: "Carry the how", p: "Position, color, and blend all travel in one box." },
        { eyebrow: "WHY TRANSLATE?", h: "Position is a transform", p: "The image never changes—only where you place it." },
        { eyebrow: "TRY NEXT", h: "Place many", p: "Change Translate in a loop to make a starfield or a row of trees." },
      ],
    },
    code: `op := &ebiten.DrawImageOptions{}
op.GeoM.Translate(x, y) // 左上を (x, y) へ
screen.DrawImage(tenjiroh, op)`,
  },
  {
    slug: "vfx-transform",
    tier: "basic",
    step: "02",
    stars: "★★☆☆☆",
    labKind: "geom",
    concept: { ja: "GeoM.Rotate / Scale", en: "GeoM.Rotate / Scale" },
    hubDesc: {
      ja: "回転と拡大縮小。そして「どこを軸に回すか」の順番。",
      en: "Rotate and scale—and the pivot: the order of operations.",
    },
    ja: {
      navConcept: "回して伸ばす",
      title: "回転と拡大縮小、そして中心",
      lead: "同じ GeoM に Rotate と Scale を足すと、絵が回って伸び縮みします。大事なのは順番。中心を原点へ寄せてから回すと、真ん中を軸にきれいに回ります。",
      deepEyebrow: "DEEP DIVE / PIVOT",
      deepH: "なぜ回すと<br>ズレて飛んでいく？",
      deepLead: "変換は「行列のかけ算」で、書いた順に効きます。左上のまま回すと、絵は左上を軸にぐるっと大回り。先に中心を原点(-w/2,-h/2)へ動かしてから回し、最後に置きたい場所へ Translate すると、中心で回ります。Rotate の角度は度ではなくラジアンです。0.7ラジアンは 0.7×180/π ≒ 40度。Scale(2,2) は原点を中心に2倍にするので、中心回転したいときは先に中心を原点へ置きます。",
      concepts: [
        { h: "原点寄せ", p: "中心を (0,0) に持ってきます。", code: "Translate(-w/2,-h/2)" },
        { h: "回す・伸ばす", p: "その状態で回転と拡大をします。", code: "Rotate / Scale" },
        { h: "戻す", p: "置きたい場所へ運びます。", code: "Translate(cx,cy)" },
      ],
      lab: {
        eyebrow: "TRY IT / GeoM",
        title: "軸を切り替えて回そう",
        p: "回転と拡大をボタンで操作し、軸を「中心」と「角」で切り替えます。角を軸にすると、絵が大きく振り回されるのが分かります。",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "中心を軸に回す", p: "原点寄せ→回転・拡大→配置。この順番が回転の要です。" },
      whys: [
        { eyebrow: "WHY ORDER?", h: "順番で結果が変わる", p: "同じ命令でも、書く順が違えば見え方が別物になります。" },
        { eyebrow: "WHY CENTER?", h: "気持ちいい回転", p: "中心軸なら、コマや弾がその場でくるくる回ります。" },
        { eyebrow: "TRY NEXT", h: "脈打たせる", p: "Scale を sin で揺らすと、心臓のように拡大縮小します。" },
      ],
    },
    en: {
      navConcept: "Rotate & scale",
      title: "Rotate, Scale & the Pivot",
      lead: "Add Rotate and Scale to the same GeoM and the image spins and stretches. Order matters: move the center to the origin first, then rotate to spin cleanly about the middle.",
      deepEyebrow: "DEEP DIVE / PIVOT",
      deepH: "Why does spinning<br>fling it away?",
      deepLead: "Transforms are matrix multiplications applied in the order you write them. Rotate while the top-left is the origin and the image swings in a big arc. Move the center to (-w/2,-h/2) first, rotate, then Translate to its place, and it spins in the middle. Rotate uses radians, not degrees: 0.7 rad is about 40°. Scale(2,2) doubles around the current origin, so move the center first when you want a center pivot.",
      concepts: [
        { h: "To origin", p: "Bring the center to (0,0).", code: "Translate(-w/2,-h/2)" },
        { h: "Rotate/Scale", p: "Now spin and stretch.", code: "Rotate / Scale" },
        { h: "Back", p: "Carry it to where it belongs.", code: "Translate(cx,cy)" },
      ],
      lab: {
        eyebrow: "TRY IT / GeoM",
        title: "Switch the pivot and spin",
        p: "Rotate and scale with buttons, then switch the pivot between center and corner. The corner pivot flings the image around in a wide swing.",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "Spin about the center", p: "To-origin → rotate/scale → place. That order is the whole trick." },
      whys: [
        { eyebrow: "WHY ORDER?", h: "Order changes the result", p: "Same commands, different order, completely different look." },
        { eyebrow: "WHY CENTER?", h: "Satisfying spins", p: "A center pivot makes tops and bullets twirl in place." },
        { eyebrow: "TRY NEXT", h: "Make it pulse", p: "Wobble Scale with sin for a heartbeat-like breathing." },
      ],
    },
    code: `op := &ebiten.DrawImageOptions{}
op.GeoM.Translate(-w/2, -h/2) // 中心を原点へ
op.GeoM.Scale(1.7, 1.7)
op.GeoM.Rotate(0.7)           // 0.7 rad ≒ 40°（ラジアンで回す）
op.GeoM.Translate(cx, cy)     // 置きたい場所へ
screen.DrawImage(sprite, op)`,
  },
  {
    slug: "vfx-tint",
    tier: "basic",
    step: "03",
    stars: "★★☆☆☆",
    labKind: "colorscale",
    concept: { ja: "op.ColorScale", en: "op.ColorScale" },
    hubDesc: {
      ja: "同じ絵を赤く染めたり、白く光らせたり、影に落としたり。",
      en: "Tint the same image red, flash it white, or drop it to a shadow.",
    },
    ja: {
      navConcept: "色を変える・抜く",
      title: "色を変える・色を抜く",
      lead: "絵のピクセルはそのまま、ColorScale で色を「かけ算」します。赤くティント、白フラッシュ、真っ黒なシルエット。1枚の絵から3つの表現が生まれます。",
      deepEyebrow: "DEEP DIVE / COLORSCALE",
      deepH: "1枚の絵から<br>何色も作れる？",
      deepLead: "ColorScale は各色チャンネルにかける倍率です。(1,0.4,0.4,1) なら緑と青が減って赤っぽく。(6,6,6,1) なら明るさが振り切れて白フラッシュ。(0,0,0,α) なら色が抜けて影になります。4つ目のアルファは0〜1が基本で、1は元の不透明度、0.5は半分、0は透明です。",
      concepts: [
        { h: "ティント", p: "色をかけ算して染めます。", code: "Scale(1,.4,.4,1)" },
        { h: "フラッシュ", p: "明るさを上げて白に。", code: "Scale(6,6,6,1)" },
        { h: "シルエット", p: "色を0にして影に。", code: "Scale(0,0,0,.5)" },
      ],
      lab: {
        eyebrow: "TRY IT / TINT",
        title: "モードを切り替えよう",
        p: "NORMAL・TINT・FLASH・SHADOW を切り替えると、同じ四角の色だけが変わります。表示されるコードが、そのときの ColorScale です。",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "被弾の白フラッシュ", p: "敵に当たった一瞬だけ倍率を上げれば、白く光って手応えが出ます。" },
      whys: [
        { eyebrow: "WHY MULTIPLY?", h: "色は倍率で決まる", p: "各チャンネルの掛け算だから、暗くも明るくもできます。" },
        { eyebrow: "WHY DRAIN?", h: "影は色抜きコピー", p: "色を0にして半透明で下に敷けば、それだけで影です。" },
        { eyebrow: "TRY NEXT", h: "点滅させる", p: "フレームごとに倍率を上下させ、無敵時間を表現しよう。" },
      ],
    },
    en: {
      navConcept: "Recolor & drain",
      title: "Tint & Drain Color",
      lead: "Keep the pixels, multiply the color with ColorScale. Red tint, white flash, pitch-black silhouette—three looks from one image.",
      deepEyebrow: "DEEP DIVE / COLORSCALE",
      deepH: "So many colors<br>from one image?",
      deepLead: "ColorScale multiplies each color channel. (1,0.4,0.4,1) cuts green and blue for red. (6,6,6,1) blows brightness out to a white flash. (0,0,0,α) drains color into a shadow. Alpha is normally 0–1: 1 keeps opacity, 0.5 halves it, and 0 is transparent.",
      concepts: [
        { h: "Tint", p: "Multiply to dye the color.", code: "Scale(1,.4,.4,1)" },
        { h: "Flash", p: "Raise brightness to white.", code: "Scale(6,6,6,1)" },
        { h: "Silhouette", p: "Zero the color for a shadow.", code: "Scale(0,0,0,.5)" },
      ],
      lab: {
        eyebrow: "TRY IT / TINT",
        title: "Switch the modes",
        p: "Flip between NORMAL, TINT, FLASH, and SHADOW and only the color of the same square changes. The code shown is the ColorScale in use.",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "The hit flash", p: "Boost the factor for one frame on a hit and the sprite flashes white with impact." },
      whys: [
        { eyebrow: "WHY MULTIPLY?", h: "Color is a factor", p: "Per-channel multiply can darken or brighten." },
        { eyebrow: "WHY DRAIN?", h: "A shadow is a drained copy", p: "Zero the color, go translucent, lay it under: instant shadow." },
        { eyebrow: "TRY NEXT", h: "Make it blink", p: "Raise and lower the factor each frame for invincibility frames." },
      ],
    },
    code: `op := &ebiten.DrawImageOptions{}
// 赤く染める（色のかけ算）
op.ColorScale.Scale(1, 0.4, 0.4, 1)
// 白フラッシュなら明るさを振り切る
// op.ColorScale.Scale(6, 6, 6, 1)
screen.DrawImage(sprite, op)`,
  },
  {
    slug: "vfx-alpha",
    tier: "basic",
    step: "04",
    stars: "★★★☆☆",
    labKind: "opacity",
    concept: { ja: "ScaleAlpha", en: "ScaleAlpha" },
    hubDesc: {
      ja: "半透明と重ね順。薄いコピーを並べて残像を作ります。",
      en: "Transparency and draw order. Stack faint copies into an afterimage.",
    },
    ja: {
      navConcept: "透明と残像",
      title: "透明と重ね順・残像",
      lead: "アルファ(不透明度)を下げると、後ろが透けます。過去の位置に薄いコピーを重ねれば、それだけで速さを感じる残像(モーションブラー)になります。",
      deepEyebrow: "DEEP DIVE / ALPHA",
      deepH: "残像って<br>どうやって出す？",
      deepLead: "難しい計算は要りません。少し前の位置を覚えておき、古いものほど薄く（小さいアルファで）重ねて描くだけ。手前ほど濃く、後ろほど透明。重ねる順番が“流れ”を作ります。",
      concepts: [
        { h: "履歴", p: "少し前の位置を覚えます。", code: "trail []vec" },
        { h: "薄める", p: "古いほど小さいアルファに。", code: "ScaleAlpha(i/len)" },
        { h: "重ねる", p: "古い順に描いて流れを作ります。", code: "DrawImage ×N" },
      ],
      lab: {
        eyebrow: "TRY IT / ALPHA",
        title: "濃さと残像の数を変えよう",
        p: "不透明度と残像の枚数をボタンで変えます。薄いコピーが増えるほど、動きの尾が長く見えます。",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "残像トレイル", p: "位置の履歴をループで描き、古いものほどアルファを下げるだけです。" },
      whys: [
        { eyebrow: "WHY ALPHA?", h: "重なりで濃くなる", p: "半透明を重ねると、通った所ほど色が濃く残ります。" },
        { eyebrow: "WHY ORDER?", h: "手前と奥", p: "後に描いた絵が前に出ます。順番が見え方を決めます。" },
        { eyebrow: "TRY NEXT", h: "フェード演出", p: "アルファを1→0にして、やられた敵をすっと消そう。" },
      ],
    },
    en: {
      navConcept: "Alpha & afterimage",
      title: "Alpha, Order & Afterimage",
      lead: "Lower alpha (opacity) and the back shows through. Stack faint copies at past positions and you get a speed-suggesting afterimage—motion blur.",
      deepEyebrow: "DEEP DIVE / ALPHA",
      deepH: "How do you make<br>an afterimage?",
      deepLead: "No hard math. Remember recent positions and draw older ones fainter (smaller alpha). Near the head is solid, the tail is transparent. Draw order makes the flow.",
      concepts: [
        { h: "History", p: "Remember recent positions.", code: "trail []vec" },
        { h: "Fade", p: "Older = smaller alpha.", code: "ScaleAlpha(i/len)" },
        { h: "Stack", p: "Draw oldest first for the trail.", code: "DrawImage ×N" },
      ],
      lab: {
        eyebrow: "TRY IT / ALPHA",
        title: "Change opacity and trail",
        p: "Change opacity and the number of afterimages with buttons. More faint copies make a longer motion tail.",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "The trail", p: "Loop the position history and lower alpha for older copies." },
      whys: [
        { eyebrow: "WHY ALPHA?", h: "Overlap deepens", p: "Stacked translucency stays darker where the path passed." },
        { eyebrow: "WHY ORDER?", h: "Front and back", p: "The last thing drawn is on top. Order decides the look." },
        { eyebrow: "TRY NEXT", h: "Fade-outs", p: "Drive alpha 1→0 to make a defeated enemy melt away." },
      ],
    },
    code: `for i, p := range trail {
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(p.x, p.y)
    a := float32(i+1) / float32(len(trail))
    op.ColorScale.ScaleAlpha(a) // 古いほど薄い
    screen.DrawImage(sprite, op)
}`,
  },
  {
    slug: "vfx-additive",
    tier: "basic",
    step: "05",
    stars: "★★★☆☆",
    labKind: "blend",
    concept: { ja: "ebiten.BlendLighter", en: "ebiten.BlendLighter" },
    hubDesc: {
      ja: "光は足し算。加算合成で、重なるほど白く輝きます。",
      en: "Light adds up. With additive blending, overlaps glow toward white.",
    },
    ja: {
      navConcept: "加算で光らせる",
      title: "加算合成で光らせる",
      lead: "ふつうの合成は「手前が奥を隠す」。加算合成 BlendLighter は「色を足す」。だから光が重なるほど明るくなり、真ん中は白く輝きます。炎も魔法もここから。",
      deepEyebrow: "DEEP DIVE / ADDITIVE",
      deepH: "なぜ光は<br>重なると白い？",
      deepLead: "現実の光と同じで、光は足し算です。赤い光と緑の光が重なれば黄色、さらに青が乗れば白。op.Blend = ebiten.BlendLighter にすると、Ebitengine が色を上に足していきます。",
      concepts: [
        { h: "通常合成", p: "手前の絵が奥を隠します。", code: "BlendSourceOver" },
        { h: "加算合成", p: "色を足して明るくします。", code: "BlendLighter" },
        { h: "光の玉", p: "中心が濃い丸を重ねます。", code: "glow ×N" },
      ],
      lab: {
        eyebrow: "TRY IT / BLEND",
        title: "合成モードを切り替えよう",
        p: "2つの光を重ね、通常合成と加算合成を切り替えます。加算にすると、重なった所だけ明るく輝きます。",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "光の玉を重ねる", p: "色を付けた丸を BlendLighter で何枚も描くと、炎の芯のように輝きます。" },
      whys: [
        { eyebrow: "WHY ADD?", h: "明るさが積み上がる", p: "重なるほど白へ近づき、エネルギーの塊に見えます。" },
        { eyebrow: "WHY DARK BG?", h: "暗い背景で映える", p: "加算は暗い所ほど目立ちます。夜空や宇宙と相性抜群。" },
        { eyebrow: "TRY NEXT", h: "ネオンにする", p: "細い線を加算で重ね描きして、光る輪郭を作ろう。" },
      ],
    },
    en: {
      navConcept: "Additive glow",
      title: "Additive Blending Glow",
      lead: "Normal blending hides the back with the front. Additive (BlendLighter) adds color, so overlapping light gets brighter and the middle glows white. Fire and magic begin here.",
      deepEyebrow: "DEEP DIVE / ADDITIVE",
      deepH: "Why is stacked<br>light white?",
      deepLead: "Just like real light, it adds. Red plus green is yellow; add blue and it is white. Set op.Blend = ebiten.BlendLighter and Ebitengine adds colors on top of each other.",
      concepts: [
        { h: "Normal", p: "The front hides the back.", code: "BlendSourceOver" },
        { h: "Additive", p: "Add colors, get brighter.", code: "BlendLighter" },
        { h: "Light orb", p: "Stack circles bright in the center.", code: "glow ×N" },
      ],
      lab: {
        eyebrow: "TRY IT / BLEND",
        title: "Toggle the blend mode",
        p: "Overlap two lights and switch between normal and additive. With additive, only the overlap flares bright.",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "Stack light orbs", p: "Draw colored circles with BlendLighter and they glow like the core of a flame." },
      whys: [
        { eyebrow: "WHY ADD?", h: "Brightness piles up", p: "Overlaps drift toward white, reading as raw energy." },
        { eyebrow: "WHY DARK BG?", h: "Pops on dark", p: "Additive shines most on dark. Perfect for night skies and space." },
        { eyebrow: "TRY NEXT", h: "Go neon", p: "Layer thin additive lines to make a glowing outline." },
      ],
    },
    code: `op := &ebiten.DrawImageOptions{}
op.ColorScale.ScaleWithColor(lightColor)
op.Blend = ebiten.BlendLighter // 光を足し算する
screen.DrawImage(glow, op)`,
  },
  {
    slug: "vfx-walk",
    tier: "basic",
    step: "06",
    stars: "★★★☆☆",
    labKind: "sheet",
    concept: { ja: "SubImage + タイマー", en: "SubImage + timer" },
    hubDesc: {
      ja: "スプライトシートから1コマずつ切り出し、歩かせます。",
      en: "Cut one frame at a time from a sheet and make it walk.",
    },
    ja: {
      navConcept: "パラパラ歩き",
      title: "パラパラ歩き（コマ送り）",
      lead: "1枚に並んだ絵(スプライトシート)から、SubImage で1コマだけ切り出して描きます。タイマーで切り出す位置を進めれば、キャラクターが歩いて見えます。",
      deepEyebrow: "DEEP DIVE / FRAMES",
      deepH: "静止画が<br>どうして歩く？",
      deepLead: "パラパラ漫画と同じ。少しずつ違うポーズを、一定時間ごとに切り替えるだけです。SubImage はシートの中の「切り取り枠」。枠を右へずらすと、次のコマが出ます。速く歩くときは切り替えも速く。",
      concepts: [
        { h: "シート", p: "コマを横に並べた1枚の絵。", code: "sheet" },
        { h: "切り出し", p: "今のコマだけを枠で抜きます。", code: "SubImage(rect)" },
        { h: "タイマー", p: "一定間隔でコマを進めます。", code: "tick / hold" },
      ],
      lab: {
        eyebrow: "TRY IT / SUBIMAGE",
        title: "コマを送って再生しよう",
        p: "コマ送りと再生/停止を切り替えます。光っている枠が、今 SubImage で切り出しているコマです。",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "歩行アニメ", p: "タイマーでコマ番号を進め、その枠を SubImage で切り出して描きます。" },
      whys: [
        { eyebrow: "WHY A SHEET?", h: "まとめて速い", p: "1枚に並べると読み込みも切り替えも軽くなります。" },
        { eyebrow: "WHY A TIMER?", h: "速さと同期", p: "移動が速いほど切り替えも速くすると、足が滑りません。" },
        { eyebrow: "TRY NEXT", h: "向きを足す", p: "左右反転(ScaleX=-1)で、進む向きに合わせて歩かせよう。" },
      ],
    },
    en: {
      navConcept: "Frame animation",
      title: "Sprite Walk (Frame by Frame)",
      lead: "From images laid out in one sheet, cut a single frame with SubImage and draw it. Advance the cut position with a timer and the character appears to walk.",
      deepEyebrow: "DEEP DIVE / FRAMES",
      deepH: "How does a still<br>image walk?",
      deepLead: "It is a flipbook. Swap slightly different poses at a steady interval. SubImage is a crop window into the sheet; slide it right and the next frame shows. Walk faster, flip faster.",
      concepts: [
        { h: "Sheet", p: "One image with frames in a row.", code: "sheet" },
        { h: "SubImage", p: "Crop out just the current frame.", code: "SubImage(rect)" },
        { h: "Timer", p: "Advance frames at an interval.", code: "tick / hold" },
      ],
      lab: {
        eyebrow: "TRY IT / SUBIMAGE",
        title: "Step and play frames",
        p: "Step frames or play/pause. The glowing cell is the frame SubImage is cutting right now.",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "The walk cycle", p: "Advance the frame index with a timer and cut that window with SubImage." },
      whys: [
        { eyebrow: "WHY A SHEET?", h: "Batched and fast", p: "One image loads and switches more cheaply." },
        { eyebrow: "WHY A TIMER?", h: "Sync to speed", p: "Flip faster when moving faster so the feet don't slide." },
        { eyebrow: "TRY NEXT", h: "Add facing", p: "Flip horizontally (ScaleX=-1) to walk in the travel direction." },
      ],
    },
    download: {
      ja: {
        eyebrow: "無料アセット / CC0",
        h: "海老・天次郎のアトラスをダウンロード",
        p: "この工房で歩かせている主人公・海老・天次郎（Ebi Tenjiroh）の一枚絵アトラスです。歩行・走行・攻撃・やられを、下(front)・上(back)・横(side)の3方向ぶん収録。1コマ 96×96px。左向きは横向きを左右反転して使います。CC0 なので、あなたのゲームに自由に使えます。",
        note: "行=アクション×向き、列=コマ。JSONに各コマの矩形(x,y,w,h)とおすすめfpsが入っています。",
        png: "PNGを保存",
        json: "コマ表(JSON)",
        license: "ライセンス",
      },
      en: {
        eyebrow: "Free asset / CC0",
        h: "Download the Ebi Tenjiroh atlas",
        p: "The one-sheet atlas of Ebi Tenjiroh (海老・天次郎), the hero you are walking in this lab. It packs walk, run, attack, and hurt in three facings—down (front), up (back), and side—at 96×96px per frame. Left-facing = flip the side frames. It is CC0, so use it freely in your own game.",
        note: "Rows = action × facing, columns = frames. The JSON has each frame rect (x,y,w,h) and a suggested fps.",
        png: "Save PNG",
        json: "Frame map (JSON)",
        license: "License",
      },
    },
    code: `// One 96×96 cell from the Ebi Tenjiroh atlas.
const fw, fh = 96, 96
i := (tick / hold) % framesInRow      // which column (frame)
x, y := i*fw, row*fh                  // row = action + facing
rect := image.Rect(x, y, x+fw, y+fh)
frame := atlas.SubImage(rect).(*ebiten.Image)
screen.DrawImage(frame, op)`,
  },
  {
    slug: "vfx-particles",
    tier: "basic",
    step: "07",
    stars: "★★★★☆",
    labKind: "spray",
    concept: { ja: "[]particle", en: "[]particle" },
    hubDesc: {
      ja: "小さな粒をたくさん。生まれて、動いて、消えていきます。",
      en: "Many tiny dots that are born, move, and fade away.",
    },
    ja: {
      navConcept: "粒をばらまく",
      title: "粒をばらまく（パーティクル）",
      lead: "火花・煙・キラキラは、小さな粒のあつまり。粒の配列を持ち、毎フレーム「位置を進めて、寿命を減らして、薄くする」だけで、いろんな演出になります。",
      deepEyebrow: "DEEP DIVE / PARTICLES",
      deepH: "煙や火花は<br>どう作る？",
      deepLead: "1個の粒はとても単純：位置と速度と寿命を持つだけ。それを配列(スライス)にたくさん持ち、Update で動かし、寿命が0になったら消します。加算合成と拡大縮小を混ぜれば、火花にも煙にもなります。",
      concepts: [
        { h: "生成", p: "タップした場所から粒を撒きます。", code: "append(ps, p)" },
        { h: "更新", p: "位置を進め、重力を足します。", code: "x+=vx; vy+=g" },
        { h: "退場", p: "寿命が尽きた粒を消します。", code: "life--" },
      ],
      lab: {
        eyebrow: "TRY IT / BURST",
        title: "バーストして撒こう",
        p: "ボタンを押すと、粒がぱっと飛び散って消えていきます。1回で何十個もの粒が生まれては消えます。",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "粒の一生", p: "配列をまわして位置を進め、重力を足し、寿命で薄くして消します。" },
      whys: [
        { eyebrow: "WHY A SLICE?", h: "数で表現する", p: "1個は単純でも、たくさん集まると豊かに見えます。" },
        { eyebrow: "WHY LIFETIME?", h: "自動で消える", p: "寿命があるから、後始末を気にせず撒けます。" },
        { eyebrow: "TRY NEXT", h: "色を混ぜる", p: "寿命で色を赤→黄に変えると、より炎らしくなります。" },
      ],
    },
    en: {
      navConcept: "Particle system",
      title: "Scatter Particles",
      lead: "Sparks, smoke, and sparkle are crowds of tiny dots. Hold an array of particles and each frame just move them, age them, and fade them for many effects.",
      deepEyebrow: "DEEP DIVE / PARTICLES",
      deepH: "How do you build<br>smoke and sparks?",
      deepLead: "One particle is simple: position, velocity, lifetime. Keep many in a slice, move them in Update, and remove them at life zero. Mix in additive blending and scaling for sparks or smoke.",
      concepts: [
        { h: "Spawn", p: "Emit particles from the tap point.", code: "append(ps, p)" },
        { h: "Update", p: "Advance position, add gravity.", code: "x+=vx; vy+=g" },
        { h: "Retire", p: "Remove particles at end of life.", code: "life--" },
      ],
      lab: {
        eyebrow: "TRY IT / BURST",
        title: "Burst them out",
        p: "Press the button and particles fly out and fade. Dozens are born and die with one tap.",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "A particle's life", p: "Loop the slice, advance position, add gravity, fade by lifetime, remove." },
      whys: [
        { eyebrow: "WHY A SLICE?", h: "Express with numbers", p: "Each dot is simple, but the crowd looks rich." },
        { eyebrow: "WHY LIFETIME?", h: "Self-cleaning", p: "A lifetime means you can spawn freely without cleanup worries." },
        { eyebrow: "TRY NEXT", h: "Blend colors", p: "Shift color red→yellow over life for a more flame-like look." },
      ],
    },
    code: `type particle struct{ x, y, vx, vy, life float64 }

for i := range ps { // Update
    ps[i].x += ps[i].vx
    ps[i].y += ps[i].vy
    ps[i].vy += gravity
    ps[i].life--
}`,
  },
  {
    slug: "vfx-spells",
    tier: "basic",
    step: "08",
    stars: "★★★★★",
    labKind: "spellbook",
    concept: { ja: "エフェクトの合成", en: "Composed effects" },
    hubDesc: {
      ja: "総まとめ。炎・水・雷を、これまでの文法だけで作ります。",
      en: "The capstone: fire, water, and lightning from the same grammar.",
    },
    ja: {
      navConcept: "炎・水・雷",
      title: "炎・水・雷を合成する",
      lead: "最終章。粒・色・加算・半透明・線を組み合わせ、海老・天次郎が3つの魔法を放ちます。新しい命令はほとんど無し。ここまでの部品の“混ぜ方”がエフェクトの正体です。",
      deepEyebrow: "DEEP DIVE / COMPOSE",
      deepH: "魔法は<br>何でできている？",
      deepLead: "炎＝炎のスプライトを加算合成で上へ。水＝水滴スプライトを半透明＋重力で落とす。雷＝稲妻スプライトを数フレームだけ＋閃光＋粒。同じ道具でも、画像・向き・色・合成・寿命を変えるだけで、まったく違う魔法になります。（画面のキャラクターはオリジナルの海老・天次郎。当たり判定はいつも通り簡単な形です。）",
      concepts: [
        { h: "炎", p: "炎スプライトを加算で上へ。", code: "Fire PNG + BlendLighter" },
        { h: "水", p: "水滴スプライトを重力で落とす。", code: "Water PNG + gravity" },
        { h: "雷", p: "稲妻スプライトを一瞬だけ。", code: "Bolt PNG + flash" },
      ],
      lab: {
        eyebrow: "TRY IT / SPELLBOOK",
        title: "混ぜ方を名前で確認",
        p: "炎・水・雷のボタンを押すと、炎スプライト／水滴／稲妻の画像が盤面で動きます。3つ全部を見たら達成です。",
      },
      play: {
        title: "3つの魔法を唱えよう",
        p: "FIRE / WATER / THUNDER をタップして魔法を放ちます。画面に出る粒・光・線が、これまでの文法の合成です。",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "同じ文法で3つの魔法", p: "合成・色・重力を切り替えるだけで、炎・水・雷を作り分けます。" },
      whys: [
        { eyebrow: "WHY COMPOSE?", h: "部品の掛け合わせ", p: "少ない道具でも、混ぜ方の数だけ表現が増えます。" },
        { eyebrow: "WHY TENJIROH?", h: "主役が映える", p: "主人公が魔法を放つと、同じエフェクトも物語になります。" },
        { eyebrow: "TRY NEXT", h: "本格魔法へ", p: "次のショーケースで、炎・氷・雷・光・闇のかっこいい魔法を順番に試そう。" },
      ],
    },
    en: {
      navConcept: "Fire, water, lightning",
      title: "Compose: Fire, Water, Lightning",
      lead: "The grammar finale. Combine particles, color, additive, transparency, and lines as Ebi Tenjiroh casts three spells. Almost no new commands—an effect is just how you mix the parts.",
      deepEyebrow: "DEEP DIVE / COMPOSE",
      deepH: "What is a spell<br>made of?",
      deepLead: "Fire = flame sprites rising with additive blending. Water = droplet sprites falling under gravity with alpha. Lightning = bolt sprites for a few frames + a flash + sparks. Same tools; change the texture, direction, color, blend, and lifetime for wildly different spells. (The on-screen character is original Ebi Tenjiroh; hit tests stay simple shapes.)",
      concepts: [
        { h: "Fire", p: "Flame sprites rising with additive blend.", code: "Fire PNG + BlendLighter" },
        { h: "Water", p: "Droplet sprites falling with gravity.", code: "Water PNG + gravity" },
        { h: "Lightning", p: "Bolt sprites for an instant flash.", code: "Bolt PNG + flash" },
      ],
      lab: {
        eyebrow: "TRY IT / SPELLBOOK",
        title: "See each recipe by name",
        p: "Press fire, water, and lightning to animate flame, droplet, and bolt sprites on the board. See all three to complete the book.",
      },
      play: {
        title: "Cast three spells",
        p: "Tap FIRE / WATER / THUNDER to cast. The particles, glow, and lines you see are the grammar from earlier lessons, mixed together.",
      },
      codeHead: { eyebrow: "IN THE REAL GAME", h: "Three spells, one grammar", p: "Switch blend, color, and gravity to tell fire, water, and lightning apart." },
      whys: [
        { eyebrow: "WHY COMPOSE?", h: "Multiply the parts", p: "Few tools, but as many looks as ways to mix them." },
        { eyebrow: "WHY TENJIROH?", h: "The hero sells it", p: "When the hero casts, the same effect becomes a story." },
        { eyebrow: "TRY NEXT", h: "Showcase magic", p: "Next: cool fire, ice, thunder, light, and dark demos that push the same tools further." },
      ],
    },
    code: `// 炎: 加算 + 上向き + オレンジ
op.Blend = ebiten.BlendLighter
op.ColorScale.ScaleWithColor(fire)
// 水: 半透明 + 重力
op.ColorScale.ScaleAlpha(0.7)
p.vy += 0.22
// 雷: 明るい線を数フレームだけ
vector.StrokeLine(screen, ax, ay, bx, by, 5, white, false)`,
  },
  {
    slug: "vfx-magic-fire",
    tier: "basic",
    step: "09",
    stars: "★★★★★",
    labKind: "magic-fire",
    concept: { ja: "チャージ炎のレイヤー", en: "Charged fire layers" },
    hubDesc: {
      ja: "長押しで溜めて放つ。炎スプライト・火の粉・熱リング・画面フラッシュ。",
      en: "Hold to charge, release to burst: flames, embers, heat ring, screen flash.",
    },
    ja: {
      navConcept: "かっこいい炎",
      title: "かっこいい炎魔法を作ろう",
      lead: "ショーケース第1弾。ボタンを長押しして熱を溜め、離すと炎が噴き上がります。小学生向けの簡単さは捨てて、「層を重ねる」ことで迫力を出します。",
      deepEyebrow: "DEEP DIVE / FIRE",
      deepH: "炎は<br>何層ある？",
      deepLead: "①炎PNGを加算で上昇 ②火の粉（spark）を扇状に ③リングPNGで衝撃波 ④画面全体のオレンジフラッシュ。長押し中は charge が0→1に増え、離した瞬間の power で粒の数と高さが変わります。同じ BlendLighter でも、寿命・重力・色のグラデを変えるだけで「焚き火」と「爆発」が分かれます。",
      concepts: [
        { h: "チャージ", p: "押し続けている間に数値を溜める。", code: "charge = min(1, charge+0.025)" },
        { h: "炎レイヤー", p: "大きな炎＋小さな火の粉を同時に。", code: "Fire PNG + Spark PNG" },
        { h: "熱リング", p: "拡大しながら薄くなる衝撃波。", code: "Ring + ScaleMulFrom" },
      ],
      lab: {
        eyebrow: "TRY IT / FIRE LAYERS",
        title: "炎の層を見てみよう",
        p: "バーストすると、炎スプライトと火の粉が同時に上がります。ゲーム本編の「チャージ→解放」はWASMデモで体感してください。",
      },
      play: {
        title: "炎を溜めて放て",
        p: "CAST FIRE を長押しして離す（または Space）。溜めが長いほど炎が大きくなります。",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "離した瞬間にバースト", p: "charge を power にして、炎・火の粉・リングを一気に生成します。" },
      whys: [
        { eyebrow: "WHY CHARGE?", h: "操作が演出になる", p: "押し続けている時間が「溜め」に見えると、同じエフェクトでも迫力が増えます。" },
        { eyebrow: "WHY LAYERS?", h: "一枚絵では足りない", p: "大きな形・細かい粒・衝撃波を分けると、脳が「本物の炎」と感じやすくなります。" },
        { eyebrow: "TRY NEXT", h: "氷へ", p: "次は「溜めて砕く」結晶のタイミング芸です。" },
      ],
    },
    en: {
      navConcept: "Cool fire",
      title: "Make Cool Fire Magic",
      lead: "Showcase #1. Hold to heat up, release for a rising plume. We drop the “keep it elementary” rule and chase punch with stacked layers.",
      deepEyebrow: "DEEP DIVE / FIRE",
      deepH: "How many<br>layers is fire?",
      deepLead: "① Flame PNGs rising with additive blend ② Ember sparks in a fan ③ Ring PNG shockwave ④ Full-screen orange flash. While held, charge climbs 0→1; on release, power scales particle count and height. Same BlendLighter—lifetime, gravity, and color ramps decide “campfire” vs “explosion”.",
      concepts: [
        { h: "Charge", p: "Accumulate a value while held.", code: "charge = min(1, charge+0.025)" },
        { h: "Flame layers", p: "Big flames plus tiny embers together.", code: "Fire PNG + Spark PNG" },
        { h: "Heat ring", p: "A shockwave that grows and fades.", code: "Ring + ScaleMulFrom" },
      ],
      lab: {
        eyebrow: "TRY IT / FIRE LAYERS",
        title: "See the fire layers",
        p: "Burst to watch flame sprites and embers rise together. Feel charge→release in the WASM demo above.",
      },
      play: {
        title: "Charge and release fire",
        p: "Hold CAST FIRE (or Space), then release. Longer charge = bigger plume.",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "Burst on release", p: "Turn charge into power, then spawn flames, embers, and a ring at once." },
      whys: [
        { eyebrow: "WHY CHARGE?", h: "Input is the show", p: "When hold-time reads as “charging”, the same VFX feels stronger." },
        { eyebrow: "WHY LAYERS?", h: "One sprite isn’t enough", p: "Split silhouette, grit, and shockwave so the eye reads “real fire”." },
        { eyebrow: "TRY NEXT", h: "Ice", p: "Next: form a crystal, then shatter it on a beat." },
      ],
    },
    code: `// Hold → charge, release → burst
if holding { charge = min(1, charge+0.025) }
else if wasHolding {
  power := 0.35 + charge*0.65
  spawnFlames(power)  // Fire PNG + BlendLighter
  spawnEmbers(power)  // Spark PNG upward
  spawnRing(power)    // Ring grows & fades
}`,
  },
  {
    slug: "vfx-magic-ice",
    tier: "basic",
    step: "10",
    stars: "★★★★★",
    labKind: "magic-ice",
    concept: { ja: "形成→粉砕のタイミング", en: "Form → shatter timing" },
    hubDesc: {
      ja: "結晶を集めてから割る。氷スプライト・霜・粉砕・ブルーム。",
      en: "Gather a crystal, then shatter: ice shards, frost, bloom.",
    },
    ja: {
      navConcept: "かっこいい氷",
      title: "かっこいい氷魔法を作ろう",
      lead: "ショーケース第2弾。タップすると霜が集まり結晶になり、一瞬の静止のあと粉砕します。「見た目」だけでなく「タイミング」が魔法の正体です。",
      deepEyebrow: "DEEP DIVE / ICE",
      deepH: "氷は<br>二幕構成",
      deepLead: "幕1（forming）: 周囲の霜を中心へ引き寄せ、氷PNGの核を大きくする。幕2（shatter）: 結晶を消し、氷スプライトを外へ飛ばし、加算のキラキラとリングのブルームを足す。phase と timer の状態機械が、同じ粒子エンジンに「物語」を載せます。",
      concepts: [
        { h: "形成", p: "霜を中心へ lerp で引き寄せる。", code: "x += (cx-x)*0.06" },
        { h: "粉砕", p: "外向き速度＋重力で砕片を飛ばす。", code: "Ice PNG + gravity" },
        { h: "ブルーム", p: "割れた瞬間だけ光るリング。", code: "Ring + Light PNG" },
      ],
      lab: {
        eyebrow: "TRY IT / ICE SHARDS",
        title: "結晶の飛び方",
        p: "バーストで氷スプライトが外へ飛び散ります。本編の「集まる→割れる」二幕はWASMデモで。",
      },
      play: {
        title: "結晶を作って砕け",
        p: "CAST ICE（または Space）。霜が集まり、結晶が割れるまで待ってみよう。",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "phase で二幕を切り替える", p: "forming のあと shatter へ。同じ Particle 型でも役割が変わります。" },
      whys: [
        { eyebrow: "WHY TWO ACTS?", h: "溜めが印象を作る", p: "いきなり砕くより、一度見せてから壊す方が「魔法っぽい」です。" },
        { eyebrow: "WHY ICE PNG?", h: "形が属性を語る", p: "丸い火花では氷に見えません。尖ったテクスチャが属性の半分です。" },
        { eyebrow: "TRY NEXT", h: "雷へ", p: "次は枝分かれする稲妻の線芸です。" },
      ],
    },
    en: {
      navConcept: "Cool ice",
      title: "Make Cool Ice Magic",
      lead: "Showcase #2. Tap to gather frost into a crystal, pause, then shatter. Timing—not just pixels—is the spell.",
      deepEyebrow: "DEEP DIVE / ICE",
      deepH: "Ice is<br>a two-act play",
      deepLead: "Act 1 (forming): pull frost inward and grow an ice-PNG core. Act 2 (shatter): hide the crystal, fling ice sprites outward, add additive sparkles and a bloom ring. A tiny phase/timer state machine puts a story on the same particle engine.",
      concepts: [
        { h: "Form", p: "Lerp frost toward the center.", code: "x += (cx-x)*0.06" },
        { h: "Shatter", p: "Outward velocity + gravity shards.", code: "Ice PNG + gravity" },
        { h: "Bloom", p: "A ring that flashes on the break.", code: "Ring + Light PNG" },
      ],
      lab: {
        eyebrow: "TRY IT / ICE SHARDS",
        title: "Watch shards fly",
        p: "Burst to scatter ice sprites. The gather→break two-act lives in the WASM demo.",
      },
      play: {
        title: "Form and shatter",
        p: "Tap CAST ICE (or Space). Watch frost gather, then the crystal break.",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "Switch acts with phase", p: "forming → shatter. Same Particle type, different jobs." },
      whys: [
        { eyebrow: "WHY TWO ACTS?", h: "Build-up sells the break", p: "Shattering after a beat feels more magical than an instant burst." },
        { eyebrow: "WHY ICE PNG?", h: "Shape reads as element", p: "Round sparks don’t say ice—sharp textures do half the work." },
        { eyebrow: "TRY NEXT", h: "Thunder", p: "Next: branching jagged lightning lines." },
      ],
    },
    code: `// Act 1: pull frost inward
f.x += (castX - f.x) * 0.06
// Act 2: shatter shards outward
spawn(IcePNG, vx: cos(a)*sp, grav: 0.12)
spawn(Ring) // bloom on break`,
  },
  {
    slug: "vfx-magic-thunder",
    tier: "basic",
    step: "11",
    stars: "★★★★★",
    labKind: "magic-thunder",
    concept: { ja: "枝分かれ稲妻", en: "Branching lightning" },
    hubDesc: {
      ja: "ジグザグ線が枝分かれ。ボルトPNG・火花・白フラッシュ。",
      en: "Jagged branching bolts, bolt PNGs, sparks, white flash.",
    },
    ja: {
      navConcept: "かっこいい雷",
      title: "かっこいい雷魔法を作ろう",
      lead: "ショーケース第3弾。空から手もとへジグザグが走り、途中で枝が分かれます。テクスチャだけでなく vector.StrokeLine の線そのものが主役です。",
      deepEyebrow: "DEEP DIVE / THUNDER",
      deepH: "雷は<br>線の再帰",
      deepLead: "メインの太線を数セグメントに分け、各点をランダムにずらす。途中で確率的に再帰呼び出しして細い枝を生やす。その上にボルトPNGと火花を重ね、画面を白くフラッシュ。外側の青グロー→本体→白い芯の3本線で「太い稲妻」に見せます。",
      concepts: [
        { h: "ジグザグ", p: "中間点をランダムにずらす。", code: "nx += (rand-0.5)*wobble" },
        { h: "枝分かれ", p: "depth 付きで再帰生成。", code: "addBolt(..., depth+1)" },
        { h: "三重線", p: "グロー・本体・白い芯。", code: "StrokeLine × 3 widths" },
      ],
      lab: {
        eyebrow: "TRY IT / BOLT FLASH",
        title: "稲妻の一瞬",
        p: "バーストでボルトと火花が走ります。枝分かれの再帰はWASMデモで確認。",
      },
      play: {
        title: "稲妻を落とせ",
        p: "CAST THUNDER（または Space）。何度も撃って枝の形の違いを楽しもう。",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "線＋テクスチャ＋閃光", p: "StrokeLine の枝とボルトPNGを同じ瞬間に重ねます。" },
      whys: [
        { eyebrow: "WHY BRANCH?", h: "自然な複雑さ", p: "一本線より、細い枝がある方が雷らしく見えます。" },
        { eyebrow: "WHY FLASH?", h: "目が追いつく前に", p: "短い寿命＋画面フラッシュで「瞬間」を強調します。" },
        { eyebrow: "TRY NEXT", h: "光へ", p: "次は放射状のフレアです。" },
      ],
    },
    en: {
      navConcept: "Cool thunder",
      title: "Make Cool Thunder Magic",
      lead: "Showcase #3. A jagged bolt runs sky-to-hand and sprouts side branches. vector.StrokeLine—not only textures—steals the show.",
      deepEyebrow: "DEEP DIVE / THUNDER",
      deepH: "Lightning is<br>recursive lines",
      deepLead: "Split the main stroke into segments, jitter midpoints, and recursively spawn thinner branches. Layer bolt PNGs and sparks, then white-flash the screen. Three strokes (blue glow → body → white core) sell thickness.",
      concepts: [
        { h: "Jagged path", p: "Jitter midpoints randomly.", code: "nx += (rand-0.5)*wobble" },
        { h: "Branches", p: "Recurse with a depth limit.", code: "addBolt(..., depth+1)" },
        { h: "Triple stroke", p: "Glow, body, and white core.", code: "StrokeLine × 3 widths" },
      ],
      lab: {
        eyebrow: "TRY IT / BOLT FLASH",
        title: "A bolt for an instant",
        p: "Burst for bolts and sparks. Recursive branching lives in the WASM demo.",
      },
      play: {
        title: "Drop lightning",
        p: "Tap CAST THUNDER (or Space). Cast again to see new branch shapes.",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "Lines + texture + flash", p: "Stack StrokeLine branches and bolt PNGs in the same frame." },
      whys: [
        { eyebrow: "WHY BRANCH?", h: "Natural complexity", p: "Thin side forks read as lightning more than a single line." },
        { eyebrow: "WHY FLASH?", h: "Before the eye catches up", p: "Short life + screen flash sells the “instant”." },
        { eyebrow: "TRY NEXT", h: "Light", p: "Next: radial flare rays." },
      ],
    },
    code: `func addBolt(x0,y0,x1,y1, depth) {
  // jittered segments...
  if depth < 2 && rand() < 0.35 {
    addBolt(nx,ny, bx,by, depth+1) // branch
  }
}
StrokeLine(glow); StrokeLine(body); StrokeLine(core)`,
  },
  {
    slug: "vfx-magic-light",
    tier: "basic",
    step: "12",
    stars: "★★★★★",
    labKind: "magic-light",
    concept: { ja: "放射フレアとブルーム", en: "Radial flare & bloom" },
    hubDesc: {
      ja: "放射状の光線・ソフトブルーム・キラキラ粒。",
      en: "Radial rays, soft bloom, twinkling sparks.",
    },
    ja: {
      navConcept: "かっこいい光",
      title: "かっこいい光魔法を作ろう",
      lead: "ショーケース第4弾。中心から光線が走り、柔らかい光の輪が広がります。暗い背景ほど加算合成が効きます。",
      deepEyebrow: "DEEP DIVE / LIGHT",
      deepH: "光は<br>中心から外へ",
      deepLead: "①中心から放射する StrokeLine のレイ ②light.png のソフトフレアを複数枚重ねてブルーム風 ③ring.png の拡大輪 ④spark のキラキラ。暗い背景＋ BlendLighter で「光源」に見せるのがコツ。色は白〜金に寄せると神聖さが出ます。",
      concepts: [
        { h: "放射レイ", p: "角度を等分して線を伸ばす。", code: "ang = i * 2π / n" },
        { h: "ブルーム", p: "半透明のフレアを重ねる。", code: "Light PNG × layers" },
        { h: "キラキラ", p: "外向きの小さな spark。", code: "Spark + outward vel" },
      ],
      lab: {
        eyebrow: "TRY IT / FLARE",
        title: "光の広がり",
        p: "バーストでフレアと粒が広がります。放射レイの本数はWASMデモで体感。",
      },
      play: {
        title: "聖なる光を放て",
        p: "CAST LIGHT（または Space）。レイとブルームが重なる様子を見てみよう。",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "レイ＋フレア＋リング", p: "線とテクスチャを同じ原点から同時に広げます。" },
      whys: [
        { eyebrow: "WHY DARK BG?", h: "加算の味方", p: "暗い画面ほど、明るい加算が「光」に見えます。" },
        { eyebrow: "WHY MANY LAYERS?", h: "ソフトなにじみ", p: "一枚の大きな円より、薄いフレアを重ねた方が自然です。" },
        { eyebrow: "TRY NEXT", h: "闇へ", p: "最後は内側へ吸い込む渦です。" },
      ],
    },
    en: {
      navConcept: "Cool light",
      title: "Make Cool Light Magic",
      lead: "Showcase #4. Rays shoot from the center while soft bloom rings expand. Dark backgrounds make additive light sing.",
      deepEyebrow: "DEEP DIVE / LIGHT",
      deepH: "Light moves<br>out from a core",
      deepLead: "① Radial StrokeLine rays ② Stacked light.png flares for soft bloom ③ Expanding ring.png ④ Outward spark twinkles. Dark BG + BlendLighter sells a light source; gold–white tints feel sacred.",
      concepts: [
        { h: "Radial rays", p: "Even angles, growing lines.", code: "ang = i * 2π / n" },
        { h: "Bloom", p: "Stack translucent flares.", code: "Light PNG × layers" },
        { h: "Twinkles", p: "Tiny outward sparks.", code: "Spark + outward vel" },
      ],
      lab: {
        eyebrow: "TRY IT / FLARE",
        title: "Watch the flare grow",
        p: "Burst for flare and sparks. Feel ray count in the WASM demo.",
      },
      play: {
        title: "Cast holy light",
        p: "Tap CAST LIGHT (or Space). Watch rays and bloom stack.",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "Rays + flare + ring", p: "Expand lines and textures from the same origin." },
      whys: [
        { eyebrow: "WHY DARK BG?", h: "Additive’s friend", p: "On a dark screen, bright additive reads as light." },
        { eyebrow: "WHY MANY LAYERS?", h: "Soft bleed", p: "Thin stacked flares beat one hard circle." },
        { eyebrow: "TRY NEXT", h: "Dark", p: "Finale: a vortex that pulls inward." },
      ],
    },
    code: `for i := 0; i < n; i++ {
  ang := base + float64(i)*2*math.Pi/float64(n)
  StrokeLine(cx,cy, cx+cos(ang)*len, cy+sin(ang)*len, ...)
}
DrawImage(LightPNG) // soft bloom layers
DrawImage(RingPNG)  // expanding ring`,
  },
  {
    slug: "vfx-magic-dark",
    tier: "basic",
    step: "13",
    stars: "★★★★★",
    labKind: "magic-dark",
    concept: { ja: "渦と吸い込み", en: "Vortex & absorb" },
    hubDesc: {
      ja: "螺旋の闇・紫の縁・中心への吸い込み。",
      en: "Spiral void, purple fringe, pull toward the core.",
    },
    ja: {
      navConcept: "かっこいい闇",
      title: "かっこいい闇魔法を作ろう",
      lead: "ショーケース最終弾。粒が螺旋を描いて中心へ吸い込まれます。炎や光が「外へ」なのに対し、闇は「内へ」。向きを反転するだけで属性が変わります。",
      deepEyebrow: "DEEP DIVE / DARK",
      deepH: "闇は<br>内側へ向かう力",
      deepLead: "螺旋上に dark.png を配置し、毎フレーム中心方向＋接線方向の力を足す。縁には紫の spark を加算で散らし、リングで衝撃波。暗い核は通常合成、紫の縁は加算——同じテクスチャでも合成モードで「穴」と「オーラ」を分けます。",
      concepts: [
        { h: "螺旋配置", p: "極座標で初期位置を決める。", code: "SpiralOffset(t, arms, r)" },
        { h: "吸い込み", p: "中心への加速度＋接線。", code: "vx += dx*0.004" },
        { h: "縁の光", p: "紫の加算で輪郭を出す。", code: "Spark + BlendLighter" },
      ],
      lab: {
        eyebrow: "TRY IT / VOID",
        title: "闇の渦",
        p: "バーストで闇スプライトが渦を巻きます。吸い込み力はWASMデモで。",
      },
      play: {
        title: "虚無の渦を開け",
        p: "CAST DARK（または Space）。粒が中心へ巻き込まれていく様子を見てみよう。",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "螺旋＋求心力", p: "初期は螺旋、更新で中心へ引っ張ります。" },
      whys: [
        { eyebrow: "WHY INWARD?", h: "属性は向きで決まる", p: "同じ粒でも、外向きなら光・炎、内向きなら闇に見えます。" },
        { eyebrow: "WHY PURPLE?", h: "黒だけだと潰れる", p: "暗い核に彩度のある縁を足すと、形が読めます。" },
        { eyebrow: "TRY NEXT", h: "自分の属性", p: "向き・色・合成を選んで、オリジナル魔法を発明しよう。次は応用編で、コアゲームに演出を足す。" },
      ],
    },
    en: {
      navConcept: "Cool dark",
      title: "Make Cool Dark Magic",
      lead: "Showcase finale. Particles spiral into a void. Fire and light push out; dark pulls in—flip direction to flip the element.",
      deepEyebrow: "DEEP DIVE / DARK",
      deepH: "Dark is<br>a force inward",
      deepLead: "Place dark.png on a spiral, then each frame add centripetal + tangential force. Rim purple sparks with additive blend, plus a ring shockwave. Dark core uses normal blend; purple fringe uses additive—same texture, two jobs.",
      concepts: [
        { h: "Spiral seed", p: "Polar coords for spawn points.", code: "SpiralOffset(t, arms, r)" },
        { h: "Absorb", p: "Pull to center + swirl.", code: "vx += dx*0.004" },
        { h: "Fringe glow", p: "Purple additive outline.", code: "Spark + BlendLighter" },
      ],
      lab: {
        eyebrow: "TRY IT / VOID",
        title: "A dark swirl",
        p: "Burst to swirl dark sprites. Feel the absorb force in the WASM demo.",
      },
      play: {
        title: "Open a void vortex",
        p: "Tap CAST DARK (or Space). Watch particles wind into the core.",
      },
      codeHead: { eyebrow: "IN THE SHOWCASE", h: "Spiral + centripetal pull", p: "Seed on a spiral, then yank toward the center each frame." },
      whys: [
        { eyebrow: "WHY INWARD?", h: "Element is direction", p: "Same particles: outward reads as light/fire, inward as dark." },
        { eyebrow: "WHY PURPLE?", h: "Black alone collapses", p: "A saturated fringe keeps the silhouette readable." },
        { eyebrow: "TRY NEXT", h: "Advanced track", p: "Next: dress up LEVEL 01–12 games and split play from fx." },
      ],
    },
    code: `dx, dy := SpiralOffset(t, arms, 160, 4.5)
// each frame: pull + swirl
p.VX += (cx-p.X)*0.004 - (cy-p.Y)*0.0025
p.VY += (cy-p.Y)*0.004 + (cx-p.X)*0.0025
DrawImage(DarkPNG) // core normal, fringe additive`,
  },
].map((l) => ({ ...l, tier: l.tier || "basic" }));

const lessons = [...basicLessons, ...illusionLessons, ...advancedLessons, ...advancedLessonsPart2];

// --- lab markup builders ----------------------------------------------------

function btn(attr, label, variant = "") {
  const cls = "lab-button" + (variant ? " " + variant : "");
  return `<button type="button" class="${cls}" ${attr}>${label}</button>`;
}
function val(attr, label) {
  return `<div><span>${label}</span><strong ${attr}>—</strong></div>`;
}
const RESET = { ja: "リセット", en: "Reset" };

function labParts(kind, lang) {
  const R = btn("data-lab-reset", RESET[lang], "lab-button-quiet");
  switch (kind) {
    case "translate":
      return {
        controls: R,
        board: `<div class="lab-board" data-lab-board></div>`,
        values: val("data-lab-x", "X") + val("data-lab-y", "Y"),
      };
    case "geom":
      return {
        controls:
          btn("data-lab-rotl", lang === "ja" ? "回転 −" : "Rotate −") +
          btn("data-lab-rotr", lang === "ja" ? "回転 +" : "Rotate +", "lab-button-primary") +
          btn("data-lab-sdown", lang === "ja" ? "縮小" : "Scale −") +
          btn("data-lab-sup", lang === "ja" ? "拡大" : "Scale +") +
          btn("data-lab-pivot", lang === "ja" ? "軸を切替" : "Pivot") + R,
        board: `<div class="lab-board" data-lab-board></div>`,
        values: val("data-lab-angle", lang === "ja" ? "角度(ラジアン)" : "angle (rad)") + val("data-lab-scale", lang === "ja" ? "倍率" : "scale") + val("data-lab-pivot", lang === "ja" ? "軸" : "pivot"),
      };
    case "colorscale":
      return {
        controls:
          btn('data-lab-mode-set="normal"', "NORMAL", "lab-button-primary") +
          btn('data-lab-mode-set="tint"', "TINT") +
          btn('data-lab-mode-set="flash"', "FLASH") +
          btn('data-lab-mode-set="shadow"', "SHADOW") + R,
        board: `<div class="lab-board" data-lab-board></div>`,
        values: val("data-lab-mode", lang === "ja" ? "モード" : "mode") + val("data-lab-code", "ColorScale"),
      };
    case "opacity":
      return {
        controls:
          btn("data-lab-aup", lang === "ja" ? "濃く" : "More solid", "lab-button-primary") +
          btn("data-lab-adown", lang === "ja" ? "薄く" : "More clear") +
          btn("data-lab-tup", lang === "ja" ? "残像 +" : "Trail +") +
          btn("data-lab-tdown", lang === "ja" ? "残像 −" : "Trail −") + R,
        board: `<div class="lab-board" data-lab-board></div>`,
        values: val("data-lab-alpha", lang === "ja" ? "不透明度" : "alpha"),
      };
    case "blend":
      return {
        controls: btn("data-lab-toggle", lang === "ja" ? "合成を切替" : "Toggle blend", "lab-button-primary") + R,
        board: `<div class="lab-board" data-lab-board data-blend="add"></div>`,
        values: val("data-lab-mode", lang === "ja" ? "合成" : "blend"),
      };
    case "sheet":
      return {
        controls:
          btn("data-lab-step", lang === "ja" ? "コマ送り" : "Step", "lab-button-primary") +
          btn("data-lab-play", lang === "ja" ? "再生 / 停止" : "Play / Pause") + R,
        board: `<div class="lab-board" data-lab-board style="display:grid;place-items:center"><div class="lab-frames">${[0, 1, 2, 3].map((i) => `<span data-lab-cell>${i}</span>`).join("")}</div></div>`,
        values: val("data-lab-frame", lang === "ja" ? "コマ" : "frame"),
      };
    case "spray":
      return {
        controls: btn("data-lab-burst", lang === "ja" ? "バースト！" : "Burst!", "lab-button-primary") + R,
        board: `<div class="lab-board" data-lab-board></div>`,
        values: val("data-lab-count", lang === "ja" ? "総数" : "count"),
      };
    case "spellbook":
      return {
        controls:
          btn('data-lab-spell="fire"', lang === "ja" ? "炎" : "Fire", "lab-button-primary") +
          btn('data-lab-spell="water"', lang === "ja" ? "水" : "Water") +
          btn('data-lab-spell="thunder"', lang === "ja" ? "雷" : "Thunder") + R,
        board: `<div class="lab-board lab-spell-stage" data-lab-board><p class="lab-spell-hint">${lang === "ja" ? "ボタンで魔法を唱えてみて" : "Cast a spell with the buttons"}</p></div>`,
        values: val("data-lab-done", lang === "ja" ? "習得" : "learned"),
      };
    case "squash":
      return {
        controls: btn('data-lab-shape="squash"', lang === "ja" ? "潰す" : "Squash") + btn('data-lab-shape="neutral"', lang === "ja" ? "普通" : "Neutral") + btn('data-lab-shape="stretch"', lang === "ja" ? "伸ばす" : "Stretch", "lab-button-primary") + R,
        board: `<div class="lab-board lab-optical" data-lab-board><div class="lab-squash-floor"></div><div class="lab-squash-actor" data-lab-actor>EBI</div></div>`,
        values: val("data-lab-xscale", "scaleX") + val("data-lab-yscale", "scaleY"),
      };
    case "outline":
      return {
        controls: btn("data-lab-outline-down", lang === "ja" ? "細く" : "Thinner") + btn("data-lab-outline-up", lang === "ja" ? "太く" : "Thicker", "lab-button-primary") + btn("data-lab-outline-color", lang === "ja" ? "黒 / 白" : "Dark / light") + R,
        board: `<div class="lab-board lab-optical" data-lab-board><div class="lab-outline-demo" data-lab-outline>EBI</div></div>`,
        values: val("data-lab-width", lang === "ja" ? "太さ" : "width") + val("data-lab-mode", lang === "ja" ? "色" : "color"),
      };
    case "impact":
      return {
        controls: btn("data-lab-impact", lang === "ja" ? "衝撃！" : "Impact!", "lab-button-primary") + btn("data-lab-rays-down", lang === "ja" ? "線 −" : "Rays −") + btn("data-lab-rays-up", lang === "ja" ? "線 +" : "Rays +") + R,
        board: `<div class="lab-board lab-optical lab-impact-demo" data-lab-board><div class="lab-impact-core"></div><div class="lab-impact-ring"></div><div class="lab-impact-rays" data-lab-rays></div></div>`,
        values: val("data-lab-rays-value", lang === "ja" ? "集中線" : "rays"),
      };
    case "bloom":
      return {
        controls: btn("data-lab-bloom-down", lang === "ja" ? "光 −" : "Glow −") + btn("data-lab-bloom-up", lang === "ja" ? "光 +" : "Glow +", "lab-button-primary") + btn("data-lab-bloom-toggle", lang === "ja" ? "ON / OFF" : "On / off") + R,
        board: `<div class="lab-board lab-optical lab-bloom-demo" data-lab-board><div class="lab-bloom-orb" data-lab-bloom></div></div>`,
        values: val("data-lab-copies", lang === "ja" ? "コピー" : "copies") + val("data-lab-mode", "bloom"),
      };
    case "magic-fire":
    case "magic-ice":
    case "magic-thunder":
    case "magic-light":
    case "magic-dark":
      return {
        controls: btn("data-lab-burst", lang === "ja" ? "バースト！" : "Burst!", "lab-button-primary") + R,
        board: `<div class="lab-board lab-spell-stage lab-magic-stage" data-lab-board data-magic="${kind.replace("magic-", "")}"><p class="lab-spell-hint">${lang === "ja" ? "バーストで層を見てみよう" : "Burst to preview the layers"}</p></div>`,
        values: val("data-lab-count", lang === "ja" ? "粒" : "sprites"),
      };
    case "fx-split":
      return {
        controls:
          btn("data-lab-fx-ping", lang === "ja" ? "命中！（play+fx）" : "Hit! (play+fx)", "lab-button-primary") +
          btn("data-lab-fx-tick", lang === "ja" ? "fx だけ 1F進める" : "Tick FX 1 frame") +
          R,
        board: `<div class="lab-board lab-fx-split" data-lab-board>
          <div class="lab-fx-col" data-lab-play><h4>${lang === "ja" ? "play（ルール・点数）" : "play (rules / score)"}</h4><ul data-lab-play-list></ul></div>
          <div class="lab-fx-col" data-lab-fx><h4>${lang === "ja" ? "fx（見た目の粒）" : "fx (sparks)"}</h4><ul data-lab-fx-list></ul></div>
          <div class="lab-fx-stage" data-lab-fx-stage aria-hidden="true"></div>
        </div>
        <p class="lab-fx-hint" data-lab-fx-hint></p>`,
        values:
          val("data-lab-score", lang === "ja" ? "スコア(play)" : "score (play)") +
          val("data-lab-fx-n", lang === "ja" ? "粒(fx)" : "sparks (fx)"),
      };
    case "fx-meter-grade":
      return {
        controls:
          `<label class="lab-meter-slider"><span>${lang === "ja" ? "止める位置" : "stop position"}</span><input type="range" min="0" max="100" value="50" data-lab-meter-position></label>` +
          btn("data-lab-meter-judge", lang === "ja" ? "判定！" : "Judge!", "lab-button-primary") +
          btn("data-lab-meter-miss", lang === "ja" ? "Miss の位置" : "Try Miss") +
          R,
        board: `<div class="lab-board lab-grade-demo" data-lab-board>
          <div class="lab-grade-track" aria-label="${lang === "ja" ? "タイミング判定ゾーン" : "timing grade zones"}">
            <span class="lab-grade-zone lab-grade-zone-ok"></span><span class="lab-grade-zone lab-grade-zone-perfect"></span>
            <span class="lab-grade-marker" data-lab-grade-marker></span>
          </div>
          <div class="lab-grade-flow">
            <div class="lab-grade-play"><small>play</small><strong data-lab-grade-result>—</strong><span data-lab-grade-rule>${lang === "ja" ? "中心との差を計算" : "measure center distance"}</span></div>
            <b aria-hidden="true">→</b>
            <div class="lab-grade-fx"><small>fx</small><strong data-lab-grade-strength>—</strong><span>${lang === "ja" ? "判定値で強さを選ぶ" : "choose strength from grade"}</span></div>
          </div>
          <div class="lab-grade-burst" data-lab-grade-burst aria-hidden="true"></div>
          <p class="lab-grade-explain" data-lab-grade-explain aria-live="polite"></p>
        </div>`,
        values:
          val("data-lab-grade-distance", lang === "ja" ? "中心との差(play)" : "center distance (play)") +
          val("data-lab-grade-particles", lang === "ja" ? "粒の数(fx)" : "particles (fx)"),
      };
    case "fx-breakout":
      return {
        controls:
          btn("data-lab-breakout-hit", lang === "ja" ? "ボールを当てる" : "Hit a brick", "lab-button-primary") +
          btn("data-lab-breakout-tick", lang === "ja" ? "破片だけ 1F進める" : "Tick shards 1 frame") +
          R,
        board: `<div class="lab-board lab-fx-split" data-lab-board>
          <div class="lab-fx-col"><h4>${lang === "ja" ? "play（当たり判定あり）" : "play (collidable)"}</h4>
            <p>● ball</p><ul data-lab-brick-list></ul></div>
          <div class="lab-fx-col"><h4>${lang === "ja" ? "fx（当たり判定なし）" : "fx (no collision)"}</h4>
            <ul data-lab-shard-list></ul></div>
          <div class="lab-fx-stage" data-lab-shard-stage aria-hidden="true"></div>
        </div><p class="lab-fx-hint" data-lab-breakout-hint></p>`,
        values:
          val("data-lab-bricks", lang === "ja" ? "ブロック(play)" : "bricks (play)") +
          val("data-lab-shards", lang === "ja" ? "破片(fx)" : "shards (fx)"),
      };
    case "fx-snake":
      return {
        controls:
          btn("data-lab-snake-step", lang === "ja" ? "進んでエサを食べる" : "Move and eat", "lab-button-primary") +
          btn("data-lab-snake-tick", lang === "ja" ? "キラキラだけ 1F進める" : "Tick sparkles 1 frame") + R,
        board: `<div class="lab-board lab-fx-split" data-lab-board>
          <div class="lab-fx-col"><h4>${lang === "ja" ? "play（体とエサ）" : "play (body and food)"}</h4>
            <p data-lab-snake-grid></p><ul data-lab-snake-body></ul></div>
          <div class="lab-fx-col"><h4>${lang === "ja" ? "fx（捕食のキラキラ）" : "fx (eat sparkles)"}</h4><ul data-lab-snake-fx></ul></div>
          <div class="lab-fx-stage" data-lab-snake-stage aria-hidden="true"></div>
        </div><p class="lab-fx-hint" data-lab-snake-hint></p>`,
        values:
          val("data-lab-snake-length", lang === "ja" ? "体セル(play)" : "body cells (play)") +
          val("data-lab-snake-parts", lang === "ja" ? "キラキラ(fx)" : "sparkles (fx)"),
      };
    default:
      return { controls: R, board: `<div class="lab-board" data-lab-board></div>`, values: val("data-lab-x", "X") };
  }
}

function labSection(lesson, lang, idx) {
  const c = lesson[lang].lab;
  const { controls, board, values } = labParts(lesson.labKind, lang);
  const id = `lab-title-${lesson.slug}`;
  return `      <div class="motion-lab" data-lab="${lesson.labKind}" aria-labelledby="${id}">
        <div class="lab-copy">
          <p class="eyebrow">${c.eyebrow}</p>
          <h3 id="${id}">${c.title}</h3>
          <p>${c.p}</p>
          <div class="lab-controls">${controls}</div>
        </div>
        <div class="lab-visual">
          ${board}
          <div class="lab-values" aria-live="polite">${values}</div>
        </div>
      </div>`;
}

// --- page builders ----------------------------------------------------------

function downloadSection(lesson, lang) {
  if (!lesson.download) return "";
  const d = lesson.download[lang];
  const a = "../../../../assets";
  return `
      <div class="atlas-download">
        <div class="atlas-copy">
          <p class="eyebrow">${d.eyebrow}</p>
          <h3>${d.h}</h3>
          <p>${d.p}</p>
          <p class="atlas-actions">
            <a class="atlas-dl" href="${a}/ebi-boy-atlas.png" download>${d.png} ↓</a>
            <a href="${a}/ebi-boy-atlas.json" download>${d.json} ↓</a>
            <a href="${a}/ebi-boy-atlas-LICENSE.txt">${d.license}</a>
          </p>
          <p class="atlas-note">${d.note}</p>
        </div>
        <img class="atlas-preview" src="${a}/ebi-boy-atlas.png" alt="Ebi Tenjiroh sprite atlas" loading="lazy" width="192">
      </div>`;
}

function stepPage(lesson, idx, lang) {
  const t = lesson[lang];
  const other = lang === "ja" ? "en" : "ja";
  const otherLabel = lang === "ja" ? "EN" : "日本語";
  const langAttr = lang === "ja" ? 'lang="en" data-language="en"' : 'lang="ja" data-language="ja"';
  const total = lessons.length;
  const courseLabel = lang === "ja" ? "工房トップ" : "Lab home";
  const learnLabel = "HOW IT WORKS";
  const playEyebrow = lang === "ja" ? "PLAYABLE / LIVE GO + マウス" : "PLAYABLE / LIVE GO + MOUSE";
  const bridge = lang === "ja"
    ? `このコースも <a href="../../../games/tap-target/#basics">LEVEL 01</a> の <strong>Update（数字）→ Draw（絵）</strong> のくり返しの上にあります。ここでは Draw の“描き方”を一段深く扱います。`
    : `This lab also sits on the <a href="../../../games/tap-target/#basics">LEVEL 01</a> loop of <strong>Update (numbers) → Draw (pixels)</strong>. Here we dig one layer deeper into how Draw paints.`;

  const concepts = t.concepts
    .map((c, i) => `        <article>
          <span class="concept-number">${i + 1}</span>
          <h3>${c.h}</h3>
          <p>${c.p}</p>
          <code>${esc(c.code)}</code>
        </article>`)
    .join("\n");

  const whys = t.whys
    .map((w, i) => `        <article${i === 2 ? ' class="challenge"' : ""}>
          <p class="eyebrow">${w.eyebrow}</p>
          <h3>${w.h}</h3>
          <p>${w.p}</p>
        </article>`)
    .join("\n");

  // Pager.
  const prev = idx === 0
    ? (lang === "ja"
        ? `<a href="../../../games/bullet-hell/">← LEVEL 12<strong>弾幕ボス</strong></a>`
        : `<a href="../../../games/bullet-hell/">← LEVEL 12<strong>Bullet Hell Boss</strong></a>`)
    : `<a href="../${lessons[idx - 1].slug}/">← STEP ${lessons[idx - 1].step}<strong>${lessons[idx - 1][lang].title}</strong></a>`;
  const next = idx === lessons.length - 1
    ? (lang === "ja"
        ? `<a href="../../platformer/">応用トラックへ →<strong>作りたいゲームを選ぶ</strong></a>`
        : `<a href="../../platformer/">Genre tracks →<strong>Pick a game to build</strong></a>`)
    : `<a href="../${lessons[idx + 1].slug}/">STEP ${lessons[idx + 1].step} →<strong>${lessons[idx + 1][lang].title}</strong></a>`;

  return `<!doctype html>
<html lang="${lang}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
  <meta name="description" content="${esc(t.lead).slice(0, 140)}">
  <title>${t.title} | Ebi Showcase</title>
  <link rel="stylesheet" href="../../../../style.css">
</head>
<body class="overview-page">
  <header class="nav">
    <a class="brand" href="../../../"><span>EBI</span> SHOWCASE</a>
    <nav>
      <a href="../">${courseLabel}</a>
      <a href="#learn">${learnLabel}</a>
      <a class="lang" href="../../../../${other}/tracks/${track}/${lesson.slug}/" ${langAttr}>${otherLabel}</a>
    </nav>
  </header>
  <main class="overview-main">
    <div class="lesson-breadcrumb"><a href="../">← ${hub[lang].title}</a><span>${lesson.tier === "advanced" ? "ADV" : "BASIC"} ${lesson.step}</span></div>

    <section class="overview-hero">
      <p class="eyebrow">${lesson.stars} / PLAYABLE</p>
      <h1>${t.title}</h1>
      <p>${t.lead}</p>
      <div class="overview-concept"><small>${lang === "ja" ? "このステップで学ぶこと" : "What you learn here"}</small><strong>${lesson.concept[lang]}</strong></div>
    </section>

    <p class="curriculum-bridge">${bridge}</p>

    <section class="play lesson-play" id="play">
      <div class="section-head">
        <div><p class="eyebrow">${playEyebrow}</p><h2>${(t.play && t.play.title) || t.lab.title}</h2></div>
        <p>${(t.play && t.play.p) || t.lab.p}</p>
      </div>
      <div class="game-shell">
        <div class="game-top"><span>● LIVE / GO + WASM</span><span>STEP ${lesson.step}</span></div>
        <div id="game-wrap"><iframe class="lesson-game-frame" title="${t.title}" src="../../../../play/${lesson.slug}/" allow="autoplay; fullscreen"></iframe></div>
      </div>
    </section>

    <section class="physics" id="learn">
      <div class="lesson-intro">
        <p class="eyebrow">${t.deepEyebrow}</p>
        <h2>${t.deepH}</h2>
        <p class="lesson-lead">${t.deepLead}</p>
      </div>

      <div class="concept-row">
${concepts}
      </div>

${labSection(lesson, lang, idx)}

      <div class="code-lesson">
        <div>
          <p class="eyebrow">${t.codeHead.eyebrow}</p>
          <h3>${t.codeHead.h}</h3>
          <p>${t.codeHead.p}</p>
        </div>
        <pre><code>${esc(lesson.code)}</code></pre>
      </div>
${downloadSection(lesson, lang)}
      <div class="why-grid">
${whys}
      </div>
    </section>

    <nav class="lesson-pager">
      ${prev}
      ${next}
    </nav>
  </main>
  <footer><div class="brand"><span>EBI</span> SHOWCASE</div><p>Made with Go + Ebitengine.</p><a href="https://github.com/kumagi/EbiShowcase">VIEW SOURCE ↗</a></footer>
  <script src="../../../../learn.js"></script>
  <script>
    document.querySelectorAll("[data-language]").forEach(a => a.addEventListener("click", () => localStorage.setItem("ebi-language", a.dataset.language)));
  </script>
</body>
</html>
`;
}

function hubPage(lang) {
  const h = hub[lang];
  const other = lang === "ja" ? "en" : "ja";
  const otherLabel = lang === "ja" ? "EN" : "日本語";
  const basic = lessons.filter((l) => l.tier !== "advanced");
  const advanced = lessons.filter((l) => l.tier === "advanced");
  const stepLink = (l) =>
    `<a class="path-step" href="${l.slug}/"><span>${l.step}</span><div><h3>${l[lang].title}</h3><p>${l.hubDesc[lang]}</p><strong>${l.concept[lang]}</strong></div><b>→</b></a>`;
  const basicSteps = basic.map(stepLink).join("\n");
  const advSteps = advanced.map(stepLink).join("\n");
  const bridge = lang === "ja"
    ? `<p class="curriculum-bridge">共通基礎(LEVEL 01〜12)の続きです。<a href="../../games/tap-target/#basics">Update / Draw</a> のループの上に、見た目を作る道具を1つずつ足します。主人公はオリジナルの海老・天次郎（えび・てんじろう）です。</p>`
    : `<p class="curriculum-bridge">This continues the core lessons (LEVEL 01–12). On top of the <a href="../../games/tap-target/#basics">Update / Draw</a> loop, add one presentation tool at a time. The hero is original Ebi Tenjiroh (海老・天次郎 / ebi-tenjiroh).</p>`;
  const basicHead = lang === "ja"
    ? `<div class="path-intro"><p class="eyebrow">BASIC / 基本編</p><h2>描画の文法を<br>手で確かめる。</h2><p>スタンプから魔法ショーケースまで。ページ上のスライダーで「置く・回す・染める」を体感する ${basic.length} ステップ（全部やらなくてOK）。</p></div>`
    : `<div class="path-intro"><p class="eyebrow">BASIC</p><h2>Feel the drawing<br>grammar by hand.</h2><p>${basic.length} steps from Stamp to magic showcases—try the on-page sliders (you don’t need every step).</p></div>`;
  const advHead = lang === "ja"
    ? `<div class="path-intro" id="advanced"><p class="eyebrow">ADVANCED / 応用編</p><h2>12のコアゲームに<br>派手な演出を足す。</h2><p>LEVEL 01〜12 をベースに、<strong>得点や当たり（ルール）</strong>と<strong>花火や光（見た目）</strong>を別の箱に入れる練習をする ${advanced.length} 章。</p></div>`
    : `<div class="path-intro" id="advanced"><p class="eyebrow">ADVANCED</p><h2>Dress up all twelve<br>core games.</h2><p>${advanced.length} chapters on LEVEL 01–12 remakes—practice keeping score/hits in one box and fireworks in another.</p></div>`;

  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><title>${h.title} | Ebi Showcase</title><link rel="stylesheet" href="../../../style.css"></head><body><header class="nav"><a class="brand" href="../../"><span>EBI</span> SHOWCASE</a><nav><a href="../../">${lang === "ja" ? "目次" : "CURRICULUM"}</a><a class="lang" href="../../../${other}/tracks/${track}/">${otherLabel}</a></nav></header><main><div class="lesson-breadcrumb"><a href="../../">${h.breadcrumb}</a><span>${basic.length}+${advanced.length} STEPS</span></div><section class="track-hero track-visual-effects"><span class="track-letter">${hub.letter}</span><div><p class="eyebrow">${h.eyebrow}</p><h1>${h.title}</h1><p>${h.lead}</p></div>${bridge}</section><section class="path-list"><div class="path-intro"><p class="eyebrow">LEARNING PATH</p><h2>${h.pathTitle}</h2><p>${h.pathLead}</p></div>${basicHead}${basicSteps}${advHead}${advSteps}</section></main></body></html>
`;
}

function homeSection(lang) {
  const basic = lessons.filter((lesson) => lesson.tier !== "advanced");
  const advanced = lessons.filter((lesson) => lesson.tier === "advanced");
  const copy = lang === "ja"
    ? {
        eyebrow: "VISUAL EFFECTS / 表現編",
        title: "ルールを覚えたら、<br>画面に手ざわりを足そう。",
        lead: "基礎編で作った動きに、回転・色・光・アニメーション・粒を足します。前半は描画の道具を1つずつ触り、後半はLEVEL 01〜12へ演出を組み込みます。ここもすべてブラウザで操作できます。",
        basicEyebrow: "DRAWING BASICS / 描画の基本",
        basicTitle: "まずは、絵を変える13の道具。",
        basicLead: "置く、回す、染める、透かす、光らせる。気になるパネルから触って大丈夫です。",
        advancedEyebrow: "GAME FX / ゲーム演出",
        advancedTitle: "次に、12のゲームへ演出を足す。",
        advancedLead: "得点や当たり判定はそのままに、光や粒だけを別の仕組みとして重ねます。",
        concept: "今回の道具",
        open: "工房の全体を見る",
      }
    : {
        eyebrow: "VISUAL EFFECTS / PRESENTATION",
        title: "Once the rules work,<br>give the screen some feel.",
        lead: "Add rotation, color, light, animation, and particles to the movement from Basics. First touch one drawing tool at a time; then wire effects into LEVEL 01–12. Every panel runs in the browser.",
        basicEyebrow: "DRAWING BASICS",
        basicTitle: "Start with thirteen ways to change a picture.",
        basicLead: "Place, spin, tint, fade, and glow. It is fine to begin with whichever panel catches your eye.",
        advancedEyebrow: "GAME FX",
        advancedTitle: "Then dress up all twelve core games.",
        advancedLead: "Keep score and collisions intact while layering light and particles as a separate system.",
        concept: "TOOL",
        open: "View the whole lab",
      };
  const cards = (items, tier) => items.map((lesson) => `<a class="vfx-course-card vfx-course-${tier}" href="tracks/${track}/${lesson.slug}/"><div class="vfx-course-meta"><span>${tier === "basic" ? "BASIC" : "FX"} ${lesson.step}</span><span>${lesson.stars}</span></div><h3>${lesson[lang].title}</h3><p>${lesson.hubDesc[lang]}</p><div class="vfx-course-concept"><small>${copy.concept}</small><strong>${lesson.concept[lang]}</strong></div><span class="vfx-course-status">TRY IT →</span></a>`).join("\n");

  return `<!-- visual-effects-home:start -->
<section class="vfx-curriculum" id="visual-effects">
  <div class="vfx-curriculum-head"><div><p class="eyebrow">${copy.eyebrow}</p><h2>${copy.title}</h2></div><p>${copy.lead}</p></div>
  <div class="vfx-course-group"><div class="vfx-group-head"><div><p class="eyebrow">${copy.basicEyebrow}</p><h3>${copy.basicTitle}</h3></div><p>${copy.basicLead}</p></div><div class="vfx-course-grid">${cards(basic, "basic")}</div></div>
  <div class="vfx-course-group vfx-course-group-advanced"><div class="vfx-group-head"><div><p class="eyebrow">${copy.advancedEyebrow}</p><h3>${copy.advancedTitle}</h3></div><p>${copy.advancedLead}</p></div><div class="vfx-course-grid">${cards(advanced, "advanced")}</div></div>
  <a class="vfx-hub-link" href="tracks/${track}/"><span>${copy.open}</span><b>→</b></a>
</section>
<!-- visual-effects-home:end -->`;
}

function updateHome(lang) {
  const file = join(root, "web", lang, "index.html");
  let html = readFileSync(file, "utf8");
  const section = homeSection(lang);
  const generated = /<!-- visual-effects-home:start -->[\s\S]*?<!-- visual-effects-home:end -->/;
  const oldPromo = /<section class="architecture-promo vfx-promo" id="visual-effects">[\s\S]*?<\/section>/;
  if (generated.test(html)) html = html.replace(generated, section);
  else if (oldPromo.test(html)) html = html.replace(oldPromo, section);
  else throw new Error(`Visual Effects home slot not found: ${file}`);
  html = html.replace(
    lang === "ja" ? '<p class="eyebrow">YOUR LEARNING PATH</p>' : '<p class="eyebrow">YOUR LEARNING PATH</p>',
    lang === "ja" ? '<p class="eyebrow">BASIC / 基礎編</p>' : '<p class="eyebrow">BASIC / CORE LESSONS</p>',
  );
  html = html.replace(
    '<p class="eyebrow">CHOOSE YOUR NEXT ADVENTURE</p>',
    lang === "ja" ? '<p class="eyebrow">ADVANCED / 応用編</p>' : '<p class="eyebrow">ADVANCED / SPECIALIZATIONS</p>',
  );
  writeFileSync(file, html);
}

// --- write files ------------------------------------------------------------

let written = 0;
for (const lang of ["ja", "en"]) {
  const base = join(root, "web", lang, "tracks", track);
  mkdirSync(base, { recursive: true });
  writeFileSync(join(base, "index.html"), hubPage(lang));
  updateHome(lang);
  written++;
  lessons.forEach((lesson, idx) => {
    const dir = join(base, lesson.slug);
    mkdirSync(dir, { recursive: true });
    writeFileSync(join(dir, "index.html"), stepPage(lesson, idx, lang));
    written++;
  });
}
console.log(`Generated ${written} Visual Effects page(s).`);
