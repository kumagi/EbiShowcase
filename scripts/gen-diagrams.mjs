/**
 * Generate the bilingual curriculum diagrams and place them into lesson pages.
 *
 * The source files are Mermaid flowcharts under docs/diagrams/.  The renderer
 * intentionally supports the small flowchart subset used here so the GitHub
 * Pages build stays dependency-free; the generated SVGs remain readable and
 * can also be regenerated with mermaid-cli when a richer layout is desired.
 * SPDX-License-Identifier: Apache-2.0
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const webRoot = path.join(root, "web");
const docsDiagramRoot = path.join(root, "docs", "diagrams");
const assetRoot = path.join(webRoot, "assets", "diagrams");
const esc = (value) => String(value).replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;").replaceAll('"', "&quot;");
const fileSafe = (value) => value.replaceAll("/", "--").replaceAll(/[^a-zA-Z0-9_-]/g, "-");

const coreLinks = [
  ["LEVEL 01 · タッチ", "LEVEL 01 · Tap", "games/tap-target/"],
  ["LEVEL 02 · 時間", "LEVEL 02 · Timing", "games/timing-meter/"],
  ["LEVEL 03 · 衝突", "LEVEL 03 · Collision", "games/catch-stars/"],
  ["LEVEL 04 · 加速度", "LEVEL 04 · Acceleration", "games/flappy/"],
  ["LEVEL 05 · 反射", "LEVEL 05 · Reflection", "games/pong/"],
  ["LEVEL 06 · 配列", "LEVEL 06 · Arrays", "games/breakout/"],
  ["LEVEL 07 · 履歴", "LEVEL 07 · History", "games/snake/"],
  ["LEVEL 08 · エンティティ", "LEVEL 08 · Entities", "games/space-shooter/"],
  ["LEVEL 09 · マップ", "LEVEL 09 · Map", "games/sokoban/"],
  ["LEVEL 10 · カメラ", "LEVEL 10 · Camera", "games/platformer/"],
  ["LEVEL 11 · AI", "LEVEL 11 · AI", "games/dungeon/"],
  ["LEVEL 12 · 弾幕", "LEVEL 12 · Bullets", "games/bullet-hell/"],
];

const diagrams = [
  {
    id: "curriculum-flow",
    title: { ja: "学びの階段", en: "Learning staircase" },
    description: { ja: "LEVEL 01から専門トラックへ進む学習の流れ", en: "The path from LEVEL 01 to specialist tracks" },
    nodes: [
      { id: "start", ja: "LEVEL 01\n数字 → 絵", en: "LEVEL 01\nNumbers → picture" },
      { id: "core", ja: "LEVEL 02–12\nゲームの部品", en: "LEVEL 02–12\nGame building blocks" },
      { id: "vfx", ja: "Visual Effects Lab\n見た目を磨く", en: "Visual Effects Lab\nPolish the picture" },
      { id: "tracks", ja: "25専門トラック\n大きなゲーム", en: "25 specialist tracks\nBigger games" },
    ],
    edges: [["start", "core"], ["core", "vfx"], ["vfx", "tracks"]],
  },
  {
    id: "glossary-map",
    title: { ja: "用語のつながり", en: "How the terms connect" },
    description: { ja: "入力をUpdateでgameへ入れ、Drawが画面へ投影する地図", en: "A map where Update turns input into game state and Draw projects it" },
    nodes: [
      { id: "input", ja: "入力\nキー・タッチ", en: "Input\nkeys / touch" },
      { id: "update", ja: "Update\n状態を書き換える", en: "Update\nmutate state" },
      { id: "state", ja: "game\nposition / score", en: "game\nposition / score" },
      { id: "draw", ja: "Draw\n投影するだけ", en: "Draw\nproject only" },
      { id: "pixels", ja: "画面\n同じ state → 同じ絵", en: "Pixels\nsame state → same frame" },
      { id: "data", ja: "JSON\n面を増やす", en: "JSON\nadd stages" },
    ],
    edges: [["input", "update"], ["update", "state"], ["state", "draw"], ["draw", "pixels"], ["state", "data"]],
  },
  {
    id: "core-loop",
    title: { ja: "LEVEL 01 ゲームループ", en: "LEVEL 01 game loop" },
    description: { ja: "入力から状態が進み、Drawがその状態を表示する役割分担", en: "Update advances state from input; Draw independently presents that state" },
    nodes: [
      { id: "input", ja: "入力\nタッチ", en: "Input\ntap" },
      { id: "update", ja: "Update\n位置を計算", en: "Update\ncalculate position" },
      { id: "draw", ja: "Draw\n現在の絵", en: "Draw\ncurrent picture" },
      { id: "repeat", ja: "設定された刻みで\nUpdateを続ける", en: "keep Update ticking\nat its configured rate" },
    ],
    edges: [["input", "update"], ["update", "repeat"], ["update", "draw"]],
  },
  {
    id: "velocity-vector",
    title: { ja: "速度ベクトル", en: "Velocity vector" },
    description: { ja: "位置に速度を足し、重力で速度を変える", en: "Add velocity to position and change it with gravity" },
    nodes: [
      { id: "gravity", ja: "重力\nvy += 0.42", en: "Gravity\nvy += 0.42" },
      { id: "velocity", ja: "速度\nvx / vy", en: "Velocity\nvx / vy" },
      { id: "position", ja: "位置\nx += vx", en: "Position\nx += vx" },
      { id: "draw", ja: "Draw\nなめらかな軌道", en: "Draw\nsmooth arc" },
    ],
    edges: [["gravity", "velocity"], ["velocity", "position"], ["position", "draw"]],
  },
  {
    id: "collision-rect",
    kind: "collision",
    title: { ja: "矩形の重なり", en: "Rectangle overlap" },
    description: { ja: "カゴと星の四角形が重なったかを調べる", en: "Ask whether the basket and star rectangles overlap" },
    nodes: [{ id: "basket", ja: "カゴ", en: "basket" }, { id: "star", ja: "星", en: "star" }, { id: "overlap", ja: "重なった？", en: "overlap?" }], edges: [["basket", "overlap"], ["star", "overlap"]],
  },
  {
    id: "spritesheet-frames",
    kind: "spritesheet",
    title: { ja: "スプライトシートのフレーム", en: "Sprite-sheet frames" },
    description: { ja: "一枚の横長画像からフレーム番号を選ぶ", en: "Choose a frame number from one strip" },
    nodes: [{ id: "sheet", ja: "一枚のシート", en: "one sheet" }, { id: "frame", ja: "frame番号", en: "frame number" }, { id: "image", ja: "表示する絵", en: "picture to draw" }], edges: [["sheet", "frame"], ["frame", "image"]],
  },
  {
    id: "input-map",
    title: { ja: "入力マッピング", en: "Input mapping" },
    description: { ja: "キーボードとタッチを同じゲーム操作へつなぐ", en: "Connect keyboard and touch to the same action" },
    nodes: [
      { id: "keyboard", ja: "← ↑ ↓ →\nキーボード", en: "← ↑ ↓ →\nkeyboard" },
      { id: "touch", ja: "画面をタップ\nタッチ", en: "tap the screen\ntouch" },
      { id: "action", ja: "同じ入力\nmove / fire", en: "same action\nmove / fire" },
      { id: "state", ja: "Update\n状態を変える", en: "Update\nchange state" },
    ],
    edges: [["keyboard", "action"], ["touch", "action"], ["action", "state"]],
  },
  {
    id: "score-flow",
    title: { ja: "得点の流れ", en: "Score flow" },
    description: { ja: "判定結果を点数と表示へ渡す", en: "Pass a result into score and feedback" },
    nodes: [
      { id: "event", ja: "イベント\nhit / miss", en: "Event\nhit / miss" },
      { id: "rule", ja: "ルール\npointsFor…", en: "Rule\npointsFor…" },
      { id: "score", ja: "score += points", en: "score += points" },
      { id: "feedback", ja: "Draw\n数字・色・音", en: "Draw\nnumber / color / sound" },
    ],
    edges: [["event", "rule"], ["rule", "score"], ["score", "feedback"]],
  },
  {
    id: "geom-chain",
    title: { ja: "GeoMの変換チェーン", en: "GeoM transform chain" },
    description: { ja: "中心を原点へ移し、回し、画面へ戻す", en: "Move to the pivot, rotate, then move back" },
    nodes: [
      { id: "origin", ja: "中心を\n原点へ", en: "move pivot\nto origin" },
      { id: "rotate", ja: "Rotate\n回転", en: "Rotate\nturn" },
      { id: "scale", ja: "Scale\n大きさ", en: "Scale\nsize" },
      { id: "screen", ja: "画面の\n場所へ", en: "move to\nscreen" },
    ],
    edges: [["origin", "rotate"], ["rotate", "scale"], ["scale", "screen"]],
  },
  {
    id: "jump-arc",
    kind: "jump",
    title: { ja: "ジャンプの軌道", en: "Jump arc" },
    description: { ja: "初速度と重力で放物線を作る", en: "Make an arc from initial speed and gravity" },
    nodes: [{ id: "speed", ja: "初速度", en: "initial speed" }, { id: "gravity", ja: "重力", en: "gravity" }, { id: "arc", ja: "軌道", en: "arc" }], edges: [["speed", "arc"], ["gravity", "arc"]],
  },
  {
    id: "puzzle-rotation",
    kind: "rotation",
    title: { ja: "パズルの回転中心", en: "Puzzle rotation pivot" },
    description: { ja: "タイルの中心を軸に90度回す", en: "Turn tiles 90 degrees around their center" },
    nodes: [{ id: "tile", ja: "タイル", en: "tile" }, { id: "pivot", ja: "中心", en: "pivot" }, { id: "turn", ja: "90°回転", en: "90° turn" }], edges: [["tile", "pivot"], ["pivot", "turn"]],
  },
  {
    id: "bullet-pool",
    title: { ja: "弾の再利用", en: "Bullet reuse" },
    description: { ja: "画面外の弾を消さずにプールへ戻す", en: "Return off-screen bullets to a reusable pool" },
    nodes: [
      { id: "pool", ja: "空きスロット\n[]Bullet", en: "free slot\n[]Bullet" },
      { id: "fire", ja: "発射\nactive=true", en: "fire\nactive=true" },
      { id: "move", ja: "Update\n位置を進める", en: "Update\nmove" },
      { id: "reuse", ja: "画面外\n空きへ戻す", en: "off-screen\nreturn" },
    ],
    edges: [["pool", "fire"], ["fire", "move"], ["move", "reuse"], ["reuse", "pool"]],
  },
  {
    id: "bpm-timing",
    title: { ja: "BPMとフレーム", en: "BPM and frames" },
    description: { ja: "1分の拍数を秒とフレームへ変換する", en: "Turn beats per minute into seconds and frames" },
    nodes: [
      { id: "bpm", ja: "BPM\n120拍/分", en: "BPM\n120 beats/min" },
      { id: "seconds", ja: "60 ÷ BPM\n0.5秒/拍", en: "60 ÷ BPM\n0.5 sec/beat" },
      { id: "frames", ja: "0.5 × 60\n30 tick", en: "0.5 × 60\n30 ticks" },
      { id: "meter", ja: "メーター\n同期して動く", en: "meter\nmove in sync" },
    ],
    edges: [["bpm", "seconds"], ["seconds", "frames"], ["frames", "meter"]],
  },
  {
    id: "input-coordinate",
    title: { ja: "入力と座標の変換", en: "Input to canvas coordinates" },
    description: { ja: "画面上のタップをゲーム内の座標へ変換する", en: "Turn a screen tap into game-space coordinates" },
    nodes: [
      { id: "pointer", ja: "指・マウス\nclientX / clientY", en: "finger / mouse\nclientX / clientY" },
      { id: "scale", ja: "表示倍率\n幅・高さ", en: "display scale\nwidth / height" },
      { id: "canvas", ja: "ゲーム座標\nx / y", en: "game coordinates\nx / y" },
      { id: "target", ja: "的を判定\nUpdateへ", en: "check the target\ninto Update" },
    ],
    edges: [["pointer", "scale"], ["scale", "canvas"], ["canvas", "target"]],
  },
  {
    id: "timing-window",
    title: { ja: "タイミング判定の窓", en: "Timing judgement window" },
    description: { ja: "針と中心の差をPerfect・OK・Missへ分ける", en: "Classify the distance from the marker as Perfect, OK, or Miss" },
    nodes: [
      { id: "marker", ja: "針の位置\nmarker", en: "marker position\nmarker" },
      { id: "center", ja: "中心との差\nabs(marker-center)", en: "distance to center\nabs(marker-center)" },
      { id: "window", ja: "判定の窓\n±4 / ±12", en: "judgement windows\n±4 / ±12" },
      { id: "feedback", ja: "演出\n色・音・粒", en: "feedback\ncolor / sound / sparks" },
    ],
    edges: [["marker", "center"], ["center", "window"], ["window", "feedback"]],
  },
  {
    id: "flappy-pipes",
    title: { ja: "パイプを一定間隔で出す", en: "Spawn pipes on a rhythm" },
    description: { ja: "タイマー、すき間、スクロール、得点を一つの流れで見る", en: "Connect the timer, gap, scrolling, and score" },
    nodes: [
      { id: "timer", ja: "spawnTimer\n1 tickごとに減る", en: "spawnTimer\ndecrease each tick" },
      { id: "gap", ja: "すき間を決める\n乱数 + clamp", en: "choose the gap\nrandom + clamp" },
      { id: "spawn", ja: "上下のパイプ\n同時に追加", en: "top + bottom pipes\nadd together" },
      { id: "scroll", ja: "左へ流す\nx -= speed", en: "scroll left\nx -= speed" },
      { id: "score", ja: "通過したら\nscore++", en: "when passed\nscore++" },
    ],
    edges: [["timer", "gap"], ["gap", "spawn"], ["spawn", "scroll"], ["scroll", "score"], ["score", "timer"]],
  },
  {
    id: "bounce-angle",
    title: { ja: "反射角を決める", en: "Choose a bounce angle" },
    description: { ja: "衝突した場所と法線から次の速度を作る", en: "Build the next velocity from the hit point and normal" },
    nodes: [
      { id: "ball", ja: "ボール\nvx / vy", en: "ball\nvx / vy" },
      { id: "hit", ja: "衝突位置\n上・左・右", en: "hit position\ntop / left / right" },
      { id: "normal", ja: "法線\n反対向き", en: "normal\npoint outward" },
      { id: "reflect", ja: "速度を反射\nvy = -vy", en: "reflect velocity\nvy = -vy" },
    ],
    edges: [["ball", "hit"], ["hit", "normal"], ["normal", "reflect"], ["reflect", "ball"]],
  },
  {
    id: "bullet-spawn",
    title: { ja: "弾幕を時間で組み立てる", en: "Build a bullet pattern over time" },
    description: { ja: "時計とパターンから弾を生み、命中後に爆発へつなぐ", en: "Use a clock and pattern to create bullets, then trigger a burst" },
    nodes: [
      { id: "clock", ja: "frame\n時間", en: "frame\ntime" },
      { id: "pattern", ja: "円形・螺旋\n角度を計算", en: "circle / spiral\ncalculate angle" },
      { id: "bullets", ja: "弾リスト\n[]Bullet", en: "bullet list\n[]Bullet" },
      { id: "hit", ja: "当たり判定\n敵・自機", en: "collision\nenemy / player" },
      { id: "burst", ja: "爆発FX\n粒を生む", en: "burst FX\nemit particles" },
    ],
    edges: [["clock", "pattern"], ["pattern", "bullets"], ["bullets", "hit"], ["hit", "burst"]],
  },
  {
    id: "tile-map",
    title: { ja: "タイルマップを再生する", en: "Play a tile map" },
    description: { ja: "文字やJSONのマップを、衝突判定と画面の絵へ分けて使う", en: "Split a text or JSON map into collision rules and pictures" },
    nodes: [
      { id: "data", ja: "map.json\n0=床 1=壁", en: "map.json\n0=floor 1=wall" },
      { id: "grid", ja: "2次元slice\nmap[y][x]", en: "2D slice\nmap[y][x]" },
      { id: "rule", ja: "移動できる？\n壁なら止める", en: "can move?\nstop at a wall" },
      { id: "draw", ja: "タイルを描く\n同じマスへ", en: "draw tiles\nat the same cells" },
    ],
    edges: [["data", "grid"], ["grid", "rule"], ["grid", "draw"]],
  },
  {
    id: "enemy-state",
    title: { ja: "敵AIの状態遷移", en: "Enemy AI state machine" },
    description: { ja: "距離とHPを見て、巡回・追跡・攻撃・退場を切り替える", en: "Switch between patrol, chase, attack, and defeat using distance and HP" },
    nodes: [
      { id: "patrol", ja: "PATROL\n見回る", en: "PATROL\nwalk around" },
      { id: "chase", ja: "CHASE\n近づく", en: "CHASE\nmove closer" },
      { id: "attack", ja: "ATTACK\n攻撃する", en: "ATTACK\nstrike" },
      { id: "hurt", ja: "HURT\nのけぞる", en: "HURT\nreact" },
      { id: "defeat", ja: "DEFEAT\n消える", en: "DEFEAT\nleave" },
    ],
    edges: [["patrol", "chase"], ["chase", "attack"], ["attack", "hurt"], ["hurt", "chase"], ["hurt", "defeat"]],
  },
  {
    id: "aim-vector",
    title: { ja: "ねらいをベクトルへ", en: "Turn aiming into a vector" },
    description: { ja: "目標との差分を正規化し、弾の速度と回転へ使う", en: "Normalize the target delta and use it for bullet speed and rotation" },
    nodes: [
      { id: "pointer", ja: "マウス位置\n目標", en: "pointer\ntarget" },
      { id: "delta", ja: "dx / dy\n差分", en: "dx / dy\ndelta" },
      { id: "normalize", ja: "長さ1へ\n正規化", en: "length 1\nnormalize" },
      { id: "shot", ja: "弾の速度\n角度", en: "bullet velocity\nangle" },
    ],
    edges: [["pointer", "delta"], ["delta", "normalize"], ["normalize", "shot"]],
  },
  {
    id: "snake-body",
    title: { ja: "スライスで体を動かす", en: "Move a body slice" },
    description: { ja: "頭を先頭へ足し、尾を外して、食べた時だけ長くする", en: "Prepend a head, remove the tail, and grow only after eating" },
    nodes: [
      { id: "head", ja: "次の頭\n方向を足す", en: "next head\nadd direction" },
      { id: "append", ja: "append\n先頭へ", en: "append\nat the front" },
      { id: "tail", ja: "尾を削る\n通常移動", en: "trim the tail\nnormal move" },
      { id: "grow", ja: "食べた？\n尾を残す", en: "ate food?\nkeep the tail" },
    ],
    edges: [["head", "append"], ["append", "tail"], ["tail", "grow"], ["grow", "head"]],
  },
  {
    id: "push-rule",
    title: { ja: "箱を押せるか決める", en: "Decide whether a box can be pushed" },
    description: { ja: "次のマスとその先を調べてから、プレイヤーと箱を動かす", en: "Check the next cell and the cell beyond before moving player and box" },
    nodes: [
      { id: "input", ja: "方向入力\n← ↑ → ↓", en: "direction input\n← ↑ → ↓" },
      { id: "next", ja: "次のマス\n床？箱？", en: "next cell\nfloor? box?" },
      { id: "beyond", ja: "箱の先\n空いている？", en: "beyond the box\nempty?" },
      { id: "move", ja: "同時に移動\nmapを更新", en: "move together\nupdate map" },
    ],
    edges: [["input", "next"], ["next", "beyond"], ["beyond", "move"], ["move", "input"]],
  },
  {
    id: "particle-lifecycle",
    title: { ja: "パーティクルの一生", en: "A particle's life" },
    description: { ja: "粒を生み、速度で動かし、透明にして消す", en: "Emit a particle, move it with velocity, fade it, and remove it" },
    nodes: [
      { id: "emit", ja: "生まれる\n位置・速度", en: "emit\nposition / velocity" },
      { id: "move", ja: "動く\nx += vx", en: "move\nx += vx" },
      { id: "fade", ja: "薄くなる\nalpha--", en: "fade\nalpha--" },
      { id: "remove", ja: "life=0\nリストから除く", en: "life=0\nremove from list" },
    ],
    edges: [["emit", "move"], ["move", "fade"], ["fade", "remove"]],
  },
  {
    id: "layer-stack",
    title: { ja: "描画レイヤーの順番", en: "Draw layers in order" },
    description: { ja: "背景、世界、エフェクト、UIを奥から手前へ重ねる", en: "Stack background, world, effects, and UI from back to front" },
    nodes: [
      { id: "back", ja: "背景\n空・床", en: "background\nsky / floor" },
      { id: "world", ja: "世界\n敵・主人公", en: "world\nenemies / hero" },
      { id: "fx", ja: "FX\n光・火花", en: "FX\nglow / sparks" },
      { id: "ui", ja: "UI\n点数・説明", en: "UI\nscore / guide" },
    ],
    edges: [["back", "world"], ["world", "fx"], ["fx", "ui"]],
  },
  {
    id: "alpha-blend",
    title: { ja: "透明度を重ねる", en: "Layer alpha" },
    description: { ja: "元の絵に透明度とブレンド方法を掛けて、残像や光を作る", en: "Combine source pixels, alpha, and a blend mode to make trails and light" },
    nodes: [
      { id: "source", ja: "元の絵\nRGBA", en: "source picture\nRGBA" },
      { id: "alpha", ja: "透明度\n0〜1", en: "alpha\n0 to 1" },
      { id: "blend", ja: "Blend\n通常・加算", en: "Blend\nnormal / lighter" },
      { id: "result", ja: "残像・光\n見た目だけ", en: "trail / glow\nvisual only" },
    ],
    edges: [["source", "alpha"], ["alpha", "blend"], ["blend", "result"]],
  },
  {
    id: "lightning-chain",
    title: { ja: "雷を枝分かれさせる", en: "Branch a lightning bolt" },
    description: { ja: "線を伸ばし、途中から枝を作り、最後に粒を散らす", en: "Extend a line, branch it midway, then scatter sparks at the end" },
    nodes: [
      { id: "seed", ja: "始点\n終点", en: "start\nend" },
      { id: "line", ja: "折れ線\n少し揺らす", en: "polyline\njitter points" },
      { id: "branch", ja: "枝分かれ\n再帰", en: "branch\nrecurse" },
      { id: "spark", ja: "着弾火花\nparticle", en: "impact spark\nparticle" },
    ],
    edges: [["seed", "line"], ["line", "branch"], ["branch", "spark"]],
  },
  {
    id: "easing-curve",
    title: { ja: "イージングで気持ちよく", en: "Make motion feel good with easing" },
    description: { ja: "急に動かさず、始まりと終わりの速さを変えて見た目を整える", en: "Vary the speed at the beginning and end to make motion readable" },
    nodes: [
      { id: "start", ja: "開始\nt=0", en: "start\nt=0" },
      { id: "ease", ja: "easeInOut\nゆっくり", en: "easeInOut\nslow start" },
      { id: "peak", ja: "中央\n速く", en: "middle\nfast" },
      { id: "settle", ja: "停止\nt=1", en: "settle\nt=1" },
    ],
    edges: [["start", "ease"], ["ease", "peak"], ["peak", "settle"]],
  },
  {
    id: "color-channel",
    title: { ja: "色を数値で変える", en: "Change color channels" },
    description: { ja: "ColorScaleのRGBA倍率で、点滅・毒・回復の色を作る", en: "Use ColorScale's RGBA multipliers for flashes, poison, and healing" },
    nodes: [
      { id: "base", ja: "元の色\nR G B A", en: "base color\nR G B A" },
      { id: "scale", ja: "倍率\n0〜2", en: "multiplier\n0 to 2" },
      { id: "mode", ja: "状態\n通常・毒・回復", en: "state\nnormal / poison / heal" },
      { id: "draw", ja: "Draw\n色が変わる", en: "Draw\ncolor changes" },
    ],
    edges: [["base", "scale"], ["scale", "mode"], ["mode", "draw"]],
  },
  {
    id: "camera-offset",
    title: { ja: "カメラのオフセット", en: "Camera offset" },
    description: { ja: "主人公の世界座標からカメラを引いて画面座標を作る", en: "Subtract the camera from world coordinates to get screen coordinates" },
    nodes: [
      { id: "world", ja: "世界座標\nplayerX", en: "world position\nplayerX" },
      { id: "focus", ja: "追従範囲\n画面中央", en: "follow zone\nscreen center" },
      { id: "camera", ja: "cameraX\n境界でclamp", en: "cameraX\nclamp at edges" },
      { id: "screen", ja: "画面座標\nworld - camera", en: "screen position\nworld - camera" },
    ],
    edges: [["world", "focus"], ["focus", "camera"], ["camera", "screen"]],
  },
  {
    id: "powerup-state",
    title: { ja: "パワーアップの状態", en: "Power-up states" },
    description: { ja: "取得、強化中、時間切れを状態として管理する", en: "Manage pickup, powered-up, and timeout as explicit states" },
    nodes: [
      { id: "normal", ja: "NORMAL\n通常", en: "NORMAL\nordinary" },
      { id: "pickup", ja: "PICKUP\nアイテム取得", en: "PICKUP\ncollect item" },
      { id: "powered", ja: "POWERED\n効果タイマー", en: "POWERED\neffect timer" },
      { id: "timeout", ja: "TIMEOUT\n元へ戻る", en: "TIMEOUT\nreturn" },
    ],
    edges: [["normal", "pickup"], ["pickup", "powered"], ["powered", "timeout"], ["timeout", "normal"]],
  },
  {
    id: "patrol-route",
    title: { ja: "Waypointを巡回する", en: "Patrol between waypoints" },
    description: { ja: "目的地へ向かい、近づいたら次のWaypointへ切り替える", en: "Move toward a waypoint and switch when close enough" },
    nodes: [
      { id: "waypoint", ja: "Waypoint\nA → B → C", en: "waypoints\nA → B → C" },
      { id: "delta", ja: "目標との差\ndx / dy", en: "target delta\ndx / dy" },
      { id: "move", ja: "正規化して移動\n速度を掛ける", en: "normalize and move\napply speed" },
      { id: "turn", ja: "近づいたら\n次へ", en: "when close\nchoose next" },
    ],
    edges: [["waypoint", "delta"], ["delta", "move"], ["move", "turn"], ["turn", "waypoint"]],
  },
  {
    id: "stage-data",
    title: { ja: "ステージをデータから増やす", en: "Add stages through data" },
    description: { ja: "同じゲームルールに、別のJSONを読み込んで新しい面を作る", en: "Reuse one game rule with another JSON file to make a new stage" },
    nodes: [
      { id: "json", ja: "stage-01.json\nstage-02.json", en: "stage-01.json\nstage-02.json" },
      { id: "decode", ja: "Decode\n構造体へ", en: "Decode\ninto structs" },
      { id: "director", ja: "StageDirector\n配置する", en: "StageDirector\nplace entities" },
      { id: "play", ja: "同じルール\n別の面", en: "same rules\nnew stage" },
    ],
    edges: [["json", "decode"], ["decode", "director"], ["director", "play"]],
  },
  {
    id: "match3-cascade",
    title: { ja: "マッチ3の連鎖", en: "Match-3 cascade" },
    description: { ja: "入れ替え、判定、消去、落下を繰り返して連鎖を作る", en: "Repeat swap, find, clear, and fall to create a chain" },
    nodes: [
      { id: "swap", ja: "入れ替え\n2マス", en: "swap\ntwo cells" },
      { id: "find", ja: "3つ以上？\nmatch", en: "three or more?\nmatch" },
      { id: "clear", ja: "消す\nscore++", en: "clear\nscore++" },
      { id: "fall", ja: "落とす\n空きを埋める", en: "fall\nfill gaps" },
    ],
    edges: [["swap", "find"], ["find", "clear"], ["clear", "fall"], ["fall", "find"]],
  },
  {
    id: "deck-cycle",
    title: { ja: "デッキを一周させる", en: "Cycle a deck" },
    description: { ja: "山札から引き、手札から使い、捨て札を混ぜて再利用する", en: "Draw, play from hand, then shuffle the discard pile back into the deck" },
    nodes: [
      { id: "draw", ja: "山札\nDraw", en: "draw pile\nDraw" },
      { id: "hand", ja: "手札\n選ぶ", en: "hand\nchoose" },
      { id: "play", ja: "Play\n効果を解決", en: "Play\nresolve effect" },
      { id: "discard", ja: "捨て札\nDiscard", en: "discard pile\nDiscard" },
      { id: "shuffle", ja: "空なら混ぜる\nShuffle", en: "when empty\nShuffle" },
    ],
    edges: [["draw", "hand"], ["hand", "play"], ["play", "discard"], ["discard", "shuffle"], ["shuffle", "draw"]],
  },
  {
    id: "range-grid",
    title: { ja: "地形と射程を重ねる", en: "Overlay terrain and range" },
    description: { ja: "移動コスト、攻撃射程、見通しを同じマス目で調べる", en: "Check movement cost, weapon range, and line of sight on one grid" },
    nodes: [
      { id: "terrain", ja: "地形\n平地・森・山", en: "terrain\nplain / forest / mountain" },
      { id: "cost", ja: "移動コスト\n1・2・3", en: "move cost\n1 / 2 / 3" },
      { id: "range", ja: "射程\n弓なら2マス", en: "range\nbow reaches 2" },
      { id: "los", ja: "見通し\n壁で止まる", en: "line of sight\nstops at wall" },
    ],
    edges: [["terrain", "cost"], ["cost", "range"], ["range", "los"]],
  },
  {
    id: "active-queue",
    title: { ja: "行動ゲージのキュー", en: "Active battle queue" },
    description: { ja: "素早さを足し、満タンのキャラを行動キューへ送る", en: "Add speed each tick and send full gauges to the action queue" },
    nodes: [
      { id: "speed", ja: "素早さ\n毎tick加算", en: "speed\nadd each tick" },
      { id: "gauge", ja: "行動ゲージ\n0 → 100", en: "action gauge\n0 → 100" },
      { id: "ready", ja: "READY\nキューへ", en: "READY\nqueue" },
      { id: "resolve", ja: "行動解決\nゲージを戻す", en: "resolve action\nreset gauge" },
    ],
    edges: [["speed", "gauge"], ["gauge", "ready"], ["ready", "resolve"], ["resolve", "speed"]],
  },
  {
    id: "dialogue-branch",
    title: { ja: "会話の選択肢とフラグ", en: "Dialogue choices and flags" },
    description: { ja: "台詞を読み、選択肢でフラグを変え、次の行を選ぶ", en: "Read a line, set a flag through a choice, and select the next line" },
    nodes: [
      { id: "line", ja: "台詞\nchapter.json", en: "dialogue\nchapter.json" },
      { id: "choice", ja: "選択肢\nA / B", en: "choices\nA / B" },
      { id: "flag", ja: "フラグ\ntrust += 1", en: "flag\ntrust += 1" },
      { id: "next", ja: "次の台詞\n分岐する", en: "next line\nbranch" },
    ],
    edges: [["line", "choice"], ["choice", "flag"], ["flag", "next"], ["next", "line"]],
  },
  {
    id: "merge-step",
    title: { ja: "同じ物を合体させる", en: "Merge matching objects" },
    description: { ja: "落下、接触、同じ種類の判定、合体、スコアの順に進める", en: "Advance through falling, contact, matching, merging, and scoring" },
    nodes: [
      { id: "drop", ja: "落とす\n重力", en: "drop\ngravity" },
      { id: "touch", ja: "接触\n円の距離", en: "contact\ncircle distance" },
      { id: "same", ja: "同じ種類？\nid一致", en: "same type?\nid matches" },
      { id: "merge", ja: "合体\n半径UP", en: "merge\nradius up" },
      { id: "score", ja: "演出と得点\n次の物を出す", en: "feedback and score\nspawn next" },
    ],
    edges: [["drop", "touch"], ["touch", "same"], ["same", "merge"], ["merge", "score"]],
  },
  {
    id: "reversi-legal",
    title: { ja: "リバーシの合法手", en: "Reversi legal move" },
    description: { ja: "8方向へ進み、相手の列の先に自分の石があるか調べる", en: "Scan eight directions and look for your stone beyond an opponent run" },
    nodes: [
      { id: "candidate", ja: "空きマス\n置いてみる", en: "empty cell\ntry a stone" },
      { id: "ray", ja: "8方向へ\n隣を調べる", en: "eight directions\ncheck neighbor" },
      { id: "bracket", ja: "相手の列\n自分で閉じる", en: "opponent run\nclosed by yours" },
      { id: "flip", ja: "合法手\n石を反転", en: "legal move\nflip stones" },
    ],
    edges: [["candidate", "ray"], ["ray", "bracket"], ["bracket", "flip"]],
  },
  {
    id: "bomb-chain",
    title: { ja: "爆弾の連鎖", en: "Bomb chain reaction" },
    description: { ja: "導火線、十字の爆風、壊れる壁、別の爆弾を順に処理する", en: "Process the fuse, cross blast, breakable walls, and other bombs in order" },
    nodes: [
      { id: "fuse", ja: "タイマー\n0まで減る", en: "fuse timer\ncount down" },
      { id: "blast", ja: "十字爆風\n上下左右", en: "cross blast\nfour directions" },
      { id: "wall", ja: "壊れる壁\n通路を開く", en: "breakable wall\nopen a path" },
      { id: "chain", ja: "別の爆弾\n即起爆", en: "another bomb\ndetonate now" },
      { id: "danger", ja: "危険マス\n逃げる", en: "danger cells\nescape" },
    ],
    edges: [["fuse", "blast"], ["blast", "wall"], ["wall", "chain"], ["chain", "danger"]],
  },
  {
    id: "race-lap",
    title: { ja: "ラップを記録する", en: "Record a racing lap" },
    description: { ja: "入力、コースのゲート、ラップタイム、自己ベストをつなぐ", en: "Connect steering, course gates, lap time, and the personal best" },
    nodes: [
      { id: "steer", ja: "入力\n加速・旋回", en: "input\naccelerate / steer" },
      { id: "gate", ja: "ゲート\n順番に通る", en: "gates\npass in order" },
      { id: "lap", ja: "ラップ\n秒を計る", en: "lap\nmeasure seconds" },
      { id: "best", ja: "BEST更新\n保存する", en: "update BEST\nsave it" },
    ],
    edges: [["steer", "gate"], ["gate", "lap"], ["lap", "best"], ["best", "steer"]],
  },
  {
    id: "monster-growth",
    title: { ja: "モンスターを育てる", en: "Grow a monster" },
    description: { ja: "経験値、レベル、技、進化をデータでつなぐ", en: "Connect experience, level, moves, and evolution through data" },
    nodes: [
      { id: "battle", ja: "戦闘\n経験値", en: "battle\nexperience" },
      { id: "level", ja: "レベルUP\n能力値", en: "level up\nstats" },
      { id: "move", ja: "技を覚える\n技データ", en: "learn a move\nmove data" },
      { id: "evolve", ja: "条件達成\n進化", en: "condition met\nevolve" },
    ],
    edges: [["battle", "level"], ["level", "move"], ["move", "evolve"]],
  },
  {
    id: "survivor-wave",
    title: { ja: "ウェーブを時間で管理する", en: "Manage survivor waves by time" },
    description: { ja: "経過時間から敵数、出現間隔、報酬選択を切り替える", en: "Use elapsed time to change enemy count, spawn interval, and reward choices" },
    nodes: [
      { id: "time", ja: "経過時間\nseconds", en: "elapsed time\nseconds" },
      { id: "wave", ja: "ウェーブ\n難度UP", en: "wave\ndifficulty up" },
      { id: "spawn", ja: "敵を出す\n間隔を短く", en: "spawn enemies\nshorter interval" },
      { id: "draft", ja: "報酬選択\n武器を強化", en: "reward draft\nupgrade weapon" },
    ],
    edges: [["time", "wave"], ["wave", "spawn"], ["spawn", "draft"], ["draft", "time"]],
  },
  {
    id: "falling-blocks",
    title: { ja: "落ちものの一手", en: "One move of falling blocks" },
    description: { ja: "回転、壁キック、落下、固定、ライン消去を順に処理する", en: "Process rotation, wall kicks, falling, locking, and line clears" },
    nodes: [
      { id: "piece", ja: "落下中\nmino", en: "falling\nmino" },
      { id: "rotate", ja: "回転\n90°", en: "rotate\n90°" },
      { id: "kick", ja: "壁キック\n候補を試す", en: "wall kick\ntry offsets" },
      { id: "lock", ja: "固定\n盤面へ", en: "lock\ninto board" },
      { id: "clear", ja: "列を消す\nscore", en: "clear lines\nscore" },
    ],
    edges: [["piece", "rotate"], ["rotate", "kick"], ["kick", "lock"], ["lock", "clear"]],
  },
  {
    id: "rpg-quest",
    title: { ja: "クエストをイベントで進める", en: "Advance a quest with events" },
    description: { ja: "イベント、条件、フラグ、報酬をデータとして再生する", en: "Replay events, conditions, flags, and rewards from data" },
    nodes: [
      { id: "event", ja: "イベント\n話す・倒す", en: "event\ntalk / defeat" },
      { id: "condition", ja: "条件\nflag / item", en: "condition\nflag / item" },
      { id: "flag", ja: "進行フラグ\nstate", en: "progress flag\nstate" },
      { id: "reward", ja: "報酬\n次の目的", en: "reward\nnext goal" },
    ],
    edges: [["event", "condition"], ["condition", "flag"], ["flag", "reward"], ["reward", "event"]],
  },
  {
    id: "sandbox-seed",
    title: { ja: "シードから世界を作る", en: "Build a world from a seed" },
    description: { ja: "同じシードなら同じ地形を再生できる", en: "The same seed can replay the same terrain" },
    nodes: [
      { id: "seed", ja: "seed\n乱数の出発点", en: "seed\nrandom start" },
      { id: "chunk", ja: "chunk\n必要な範囲だけ", en: "chunk\nload nearby" },
      { id: "noise", ja: "地形ルール\n高さ・資源", en: "terrain rule\nheight / resource" },
      { id: "world", ja: "配置\n遊べる世界", en: "placement\nplayable world" },
    ],
    edges: [["seed", "chunk"], ["chunk", "noise"], ["noise", "world"]],
  },
  {
    id: "capture-flow",
    title: { ja: "捕獲シーケンス", en: "Capture sequence" },
    description: { ja: "弱らせる、投げる、揺れる、成功判定を順番に演出する", en: "Stage weaken, throw, shake, and success checks in order" },
    nodes: [
      { id: "weaken", ja: "弱らせる\nHPを減らす", en: "weaken\nlower HP" },
      { id: "throw", ja: "投げる\n軌道FX", en: "throw\ntrajectory FX" },
      { id: "shake", ja: "揺れる\nタイマー", en: "shake\ntimer" },
      { id: "result", ja: "捕獲成功？\n再挑戦", en: "captured?\ntry again" },
    ],
    edges: [["weaken", "throw"], ["throw", "shake"], ["shake", "result"], ["result", "weaken"]],
  },
];

const routeDiagrams = new Map([
  ["games/tap-target", ["core-loop", "input-coordinate"]],
  ["games/flappy", ["velocity-vector", "flappy-pipes"]],
  ["games/catch-stars", ["collision-rect", "score-flow"]],
  ["games/breakout", ["score-flow", "bounce-angle"]],
  ["games/space-shooter", ["input-map", "aim-vector", "bullet-pool"]],
  ["games/timing-meter", ["bpm-timing", "timing-window"]],
  ["games/pong", ["bounce-angle"]],
  ["games/platformer", ["jump-arc", "tile-map"]],
  ["games/dungeon-crawler", ["tile-map", "enemy-state"]],
  ["games/dungeon", ["tile-map", "enemy-state"]],
  ["games/bullet-hell", ["bullet-spawn", "layer-stack"]],
  ["games/snake", ["snake-body"]],
  ["games/sokoban", ["push-rule", "tile-map"]],
  ["tracks/visual-effects/vfx-walk", ["spritesheet-frames"]],
  ["tracks/visual-effects/vfx-transform", ["geom-chain"]],
  ["tracks/visual-effects/vfx-particles", ["particle-lifecycle"]],
  ["tracks/visual-effects/vfx-magic-fire", ["particle-lifecycle", "layer-stack"]],
  ["tracks/visual-effects/fx-dungeon", ["layer-stack"]],
  ["tracks/visual-effects/vfx-spells", ["layer-stack", "alpha-blend"]],
  ["tracks/visual-effects/vfx-alpha", ["alpha-blend"]],
  ["tracks/visual-effects/vfx-additive", ["alpha-blend"]],
  ["tracks/visual-effects/vfx-magic-thunder", ["lightning-chain"]],
  ["tracks/visual-effects/vfx-squash", ["easing-curve"]],
  ["tracks/visual-effects/vfx-impact-lines", ["easing-curve"]],
  ["tracks/visual-effects/vfx-tint", ["color-channel"]],
  ["tracks/visual-effects/vfx-fx-shooter", ["bullet-pool", "layer-stack"]],
  ["tracks/visual-effects/vfx-fx-bullethell", ["bullet-spawn"]],
  ["tracks/visual-effects/vfx-fx-snake", ["snake-body"]],
  ["tracks/visual-effects/vfx-fx-sokoban", ["push-rule"]],
  ["tracks/visual-effects/vfx-fx-breakout", ["bounce-angle"]],
  ["tracks/visual-effects/fx-pong", ["bounce-angle"]],
  ["tracks/visual-effects/fx-meter", ["timing-window"]],
  ["tracks/visual-effects/fx-tap", ["particle-lifecycle"]],
  ["tracks/visual-effects/vfx-fx-catch", ["particle-lifecycle"]],
  ["tracks/platformer/scrolling-stage", ["camera-offset"]],
  ["tracks/platformer/powerup-adventure", ["powerup-state"]],
  ["tracks/platformer/patrol-enemies", ["patrol-route"]],
  ["tracks/platformer/moving-platforms", ["jump-arc", "camera-offset"]],
  ["tracks/platformer/stage-data", ["stage-data"]],
  ["tracks/platformer/run-animation", ["spritesheet-frames"]],
  ["tracks/match3/grid-swap", ["puzzle-rotation", "match3-cascade"]],
  ["tracks/match3/find-matches", ["match3-cascade"]],
  ["tracks/match3/cascade", ["match3-cascade"]],
  ["tracks/deckbuilder/deck-cycle", ["deck-cycle"]],
  ["tracks/deckbuilder/card-motion", ["easing-curve"]],
  ["tracks/tactics/terrain-reach", ["range-grid"]],
  ["tracks/tactics/weapon-range", ["range-grid"]],
  ["tracks/active-rpg/gauge-race", ["bpm-timing", "active-queue"]],
  ["tracks/active-rpg/ready-queue", ["active-queue"]],
  ["tracks/visual-novel/dialogue-window", ["dialogue-branch"]],
  ["tracks/visual-novel/choice-flags", ["dialogue-branch"]],
  ["tracks/merge-physics/merge-rule", ["merge-step"]],
  ["tracks/merge-physics/circle-collision", ["merge-step"]],
  ["tracks/reversi/legal-moves", ["reversi-legal"]],
  ["tracks/reversi/ebi-reversi", ["reversi-legal"]],
  ["tracks/bomb-maze/timed-bomb", ["bomb-chain"]],
  ["tracks/bomb-maze/chain-explosion", ["bomb-chain"]],
  ["tracks/racing/ebi-racing", ["race-lap"]],
  ["tracks/racing/lap-gates", ["race-lap"]],
  ["tracks/monster-collection/growth-evolution", ["monster-growth"]],
  ["tracks/monster-collection/capture-sequence", ["capture-flow"]],
  ["tracks/survivors/wave-director", ["survivor-wave"]],
  ["tracks/survivors/auto-turret", ["aim-vector", "survivor-wave"]],
  ["tracks/falling-blocks/rotation-kicks", ["falling-blocks"]],
  ["tracks/falling-blocks/lock-lines", ["falling-blocks"]],
  ["tracks/rpg/ebi-quest", ["rpg-quest"]],
  ["tracks/rpg/dialogue-flags", ["rpg-quest"]],
  ["tracks/sandbox/chunk-world", ["sandbox-seed"]],
  ["tracks/sandbox/terrain-generation", ["sandbox-seed"]],
]);

function nodeLabel(node, lang) {
  return (node[lang] || node.en).split("\n");
}

function mermaidSource(def, lang) {
  const lines = ["flowchart LR"];
  for (const node of def.nodes) {
    const label = (node[lang] || node.en).replaceAll('"', "'").replaceAll("\n", "<br/>");
    lines.push(`  ${node.id}["${label}"]`);
  }
  for (const [from, to] of def.edges) lines.push(`  ${from} --> ${to}`);
  return `${lines.join("\n")}\n`;
}

function svgShell(width, height, title, description, content) {
  return `<!-- SPDX-License-Identifier: Apache-2.0 -->\n<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${width} ${height}" role="img" aria-labelledby="title desc" focusable="false"><title id="title">${esc(title)}</title><desc id="desc">${esc(description)}</desc><defs><marker id="arrow" markerWidth="9" markerHeight="9" refX="7" refY="4.5" orient="auto"><path d="M0,0 L9,4.5 L0,9 z" fill="#6a58e6"/></marker><linearGradient id="card" x1="0" x2="1" y1="0" y2="1"><stop stop-color="#ffffff"/><stop offset="1" stop-color="#eef1ff"/></linearGradient></defs>${content}</svg>\n`;
}

function multilineText(lines, x, y, options = {}) {
  const size = options.size || 18;
  const color = options.color || "#2b3158";
  const weight = options.weight || 800;
  const lineHeight = options.lineHeight || 25;
  const anchor = options.anchor || "middle";
  return `<text x="${x}" y="${y}" text-anchor="${anchor}" font-family="system-ui,-apple-system,'Hiragino Sans',sans-serif" font-size="${size}" font-weight="${weight}" fill="${color}">${lines.map((line, i) => `<tspan x="${x}" dy="${i ? lineHeight : 0}">${esc(line)}</tspan>`).join("")}</text>`;
}

function graphSvg(def, lang) {
  const width = Math.max(900, def.nodes.length * 215 + 80);
  const cols = Math.min(4, Math.max(2, def.nodes.length));
  const cardW = 178;
  const cardH = 90;
  const gapX = (width - 80 - cols * cardW) / Math.max(1, cols - 1);
  const positions = new Map();
  def.nodes.forEach((node, index) => {
    const row = Math.floor(index / cols);
    const col = index % cols;
    positions.set(node.id, { x: 40 + col * (cardW + gapX), y: 82 + row * 150 });
  });
  const rows = Math.ceil(def.nodes.length / cols);
  const height = 115 + rows * 150;
  let content = `<rect width="${width}" height="${height}" rx="28" fill="#f4f6ff"/>`;
  for (const [from, to] of def.edges) {
    const a = positions.get(from); const b = positions.get(to);
    if (!a || !b) continue;
    const x1 = a.x + cardW; const y1 = a.y + cardH / 2;
    const x2 = b.x; const y2 = b.y + cardH / 2;
    const bend = Math.max(30, Math.abs(x2 - x1) * .3);
    content += `<path d="M${x1} ${y1} C${x1 + bend} ${y1},${x2 - bend} ${y2},${x2} ${y2}" fill="none" stroke="#6a58e6" stroke-width="4" marker-end="url(#arrow)"/>`;
  }
  def.nodes.forEach((node, index) => {
    const p = positions.get(node.id);
    const fill = index % 3 === 0 ? "#dffaf4" : index % 3 === 1 ? "#e8e2ff" : "#ffe4ef";
    content += `<g><rect x="${p.x}" y="${p.y}" width="${cardW}" height="${cardH}" rx="18" fill="url(#card)" stroke="${fill}" stroke-width="8"/><circle cx="${p.x + 23}" cy="${p.y + 23}" r="8" fill="#2ee6c8"/>${multilineText(nodeLabel(node, lang), p.x + cardW / 2, p.y + 42, {size: 17, lineHeight: 23})}</g>`;
  });
  return svgShell(width, height, def.title[lang], def.description[lang], content);
}

function collisionSvg(def, lang) {
  const title = def.title[lang]; const description = def.description[lang];
  const labels = lang === "ja" ? ["カゴ", "星", "重なった？"] : ["basket", "star", "overlap?"];
  const rules = lang === "ja"
    ? ["星の左 < カゴの右", "星の右 > カゴの左", "星の上 < カゴの下", "星の下 > カゴの上"]
    : ["star left < basket right", "star right > basket left", "star top < basket bottom", "star bottom > basket top"];
  const content = `<rect width="920" height="390" rx="28" fill="#f4f6ff"/><rect x="80" y="78" width="250" height="130" rx="18" fill="#dffaf4" stroke="#0a8f7f" stroke-width="6"/><rect x="255" y="128" width="250" height="130" rx="18" fill="#ffe4ef" stroke="#e56891" stroke-width="6"/><path d="M510 166 C565 166 585 166 625 166" fill="none" stroke="#6a58e6" stroke-width="5" marker-end="url(#arrow)"/>${multilineText([labels[0]], 205, 145, {size: 23})}${multilineText([labels[1]], 380, 199, {size: 23})}${multilineText([labels[2]], 745, 176, {size: 22, color: "#6a58e6"})}<g font-family="monospace" font-size="15" fill="#2b3158"><text x="80" y="305">✓ ${esc(rules[0])}</text><text x="485" y="305">✓ ${esc(rules[1])}</text><text x="80" y="345">✓ ${esc(rules[2])}</text><text x="485" y="345">✓ ${esc(rules[3])}</text></g><text x="80" y="375" font-family="monospace" font-size="14" fill="#6a58e6">AABB: all four conditions must be true</text>`;
  return svgShell(920, 390, title, description, content);
}

function spritesheetSvg(def, lang) {
  const title = def.title[lang]; const description = def.description[lang];
  const frameLabel = lang === "ja" ? "frame++ → 表示する画像を選ぶ" : "frame++ → choose the picture";
  let content = `<rect width="920" height="300" rx="28" fill="#f4f6ff"/><text x="55" y="42" font-family="monospace" font-size="15" fill="#6a58e6">SubImage(sheet, frame * frameWidth, 0, frameWidth, frameHeight)</text>`;
  for (let i = 0; i < 6; i++) {
    const x = 55 + i * 125;
    content += `<rect x="${x}" y="85" width="110" height="110" rx="12" fill="${i === 2 ? "#ffe4ef" : "#e8e2ff"}" stroke="${i === 2 ? "#e56891" : "#8d7bff"}" stroke-width="${i === 2 ? 6 : 3}"/><circle cx="${x + 55}" cy="140" r="${22 + i * 2}" fill="#2ee6c8"/><text x="${x + 55}" y="228" text-anchor="middle" font-family="monospace" font-size="15" fill="#2b3158">frame ${i}</text>`;
  }
  content += `<path d="M${55 + 2 * 125 + 110} 140 C610 140 635 140 700 140" fill="none" stroke="#6a58e6" stroke-width="5" marker-end="url(#arrow)"/>${multilineText([frameLabel], 805, 132, {size: 18, color: "#6a58e6"})}`;
  return svgShell(920, 300, title, description, content);
}

function jumpSvg(def, lang) {
  const title = def.title[lang]; const description = def.description[lang];
  const labels = lang === "ja" ? ["初速度", "重力", "なめらかな放物線"] : ["initial speed", "gravity", "smooth arc"];
  const content = `<rect width="920" height="300" rx="28" fill="#f4f6ff"/><path d="M90 235 C240 235 220 55 460 55 S690 235 830 235" fill="none" stroke="#6a58e6" stroke-width="8"/><path d="M90 235 L830 235" stroke="#9aa6ce" stroke-width="3" stroke-dasharray="8 8"/>${multilineText([labels[0]], 115, 210, {size: 17, anchor: "start"})}${multilineText([labels[1]], 440, 78, {size: 17})}${multilineText([labels[2]], 700, 210, {size: 19, color: "#6a58e6"})}<text x="55" y="28" font-family="monospace" font-size="15" fill="#6a58e6"><tspan x="55" dy="0">y += vy</tspan><tspan x="55" dy="19">vy += gravity</tspan></text>`;
  return svgShell(920, 300, title, description, content);
}

function rotationSvg(def, lang) {
  const title = def.title[lang]; const description = def.description[lang];
  const labels = lang === "ja" ? ["回す前", "中心", "90°回転"] : ["before", "pivot", "90° turn"];
  const content = `<rect width="920" height="300" rx="28" fill="#f4f6ff"/><rect x="125" y="95" width="120" height="120" rx="12" fill="#dffaf4" stroke="#0a8f7f" stroke-width="5"/><circle cx="185" cy="155" r="7" fill="#e56891"/><path d="M300 155 H575" stroke="#6a58e6" stroke-width="5" marker-end="url(#arrow)"/><rect x="650" y="95" width="120" height="120" rx="12" fill="#ffe4ef" stroke="#e56891" stroke-width="5" transform="rotate(90 710 155)"/><circle cx="710" cy="155" r="7" fill="#e56891"/>${multilineText([labels[0]], 185, 250, {size: 18})}${multilineText([labels[1]], 185, 155, {size: 14, color: "#e56891"})}${multilineText([labels[2]], 710, 250, {size: 18})}<text x="55" y="40" font-family="monospace" font-size="15" fill="#6a58e6">Translate(-pivot) → Rotate → Translate(pivot)</text>`;
  return svgShell(920, 300, title, description, content);
}

function deviceSvg(lang) {
  const ja = lang === "ja";
  const title = ja ? "PC・タブレット・スマホの表示" : "PC, tablet, and phone layout";
  const description = ja ? "同じゲーム画面が幅に合わせて縮み、操作しやすい場所へ並ぶ比較図" : "The same game canvas scales down and keeps controls reachable";
  const labels = ja ? [["PC", "横に説明 + ゲーム"], ["タブレット", "2列をゆったり"], ["スマホ", "1列 + 大きな操作"]] : [["PC", "copy + game side by side"], ["Tablet", "two roomy columns"], ["Phone", "one column + large controls"]];
  let content = `<rect width="1100" height="360" rx="30" fill="#f4f6ff"/>`;
  const widths = [300, 250, 180]; const x = [45, 425, 775];
  widths.forEach((w, i) => {
    const h = i === 0 ? 190 : i === 1 ? 225 : 270;
    const y = 75 + (270 - h);
    content += `<g><rect x="${x[i]}" y="${y}" width="${w}" height="${h}" rx="${i === 2 ? 24 : 16}" fill="#131a3f" stroke="#8d7bff" stroke-width="5"/><rect x="${x[i] + 18}" y="${y + 22}" width="${w - 36}" height="${h - 80}" rx="10" fill="#2e3a72"/><rect x="${x[i] + 28}" y="${y + h - 45}" width="${w - 56}" height="13" rx="7" fill="#2ee6c8"/>${multilineText(labels[i], x[i] + w / 2, y + h + 32, {size: 17, color: "#2b3158", lineHeight: 22})}</g>`;
  });
  return svgShell(1100, 360, title, description, content);
}

function writeDiagram(def, lang) {
  const svg = def.kind === "collision" ? collisionSvg(def, lang) : def.kind === "spritesheet" ? spritesheetSvg(def, lang) : def.kind === "jump" ? jumpSvg(def, lang) : def.kind === "rotation" ? rotationSvg(def, lang) : graphSvg(def, lang);
  fs.writeFileSync(path.join(assetRoot, `${lang}_${def.id}.svg`), svg);
  if (def.nodes.length) fs.writeFileSync(path.join(docsDiagramRoot, `${def.id}.${lang}.mmd`), mermaidSource(def, lang));
}

function figure(relativeAsset, alt, caption, extraClass = "") {
  return `<figure class="diagram-figure ${extraClass}"><img src="${relativeAsset}" alt="${esc(alt)}" loading="lazy" decoding="async"><figcaption>${esc(caption)}</figcaption></figure>`;
}

function curriculumBlock(lang) {
  const ja = lang === "ja";
  const links = coreLinks.map(([jaLabel, enLabel, href]) => `<a href="${href}"><span>${ja ? jaLabel : enLabel}</span><b>→</b></a>`).join("");
  return `<!-- curriculum-map:start -->\n<section class="curriculum-map" aria-labelledby="curriculum-map-title"><div><p class="eyebrow">CURRICULUM MAP / 学びの地図</p><h2 id="curriculum-map-title">${ja ? "小さな遊びから、大きなゲームへ。" : "From one tiny game to a bigger one."}</h2><p>${ja ? "最初は丸を押すだけ。次のレベルでは時間、当たり判定、動きへ進み、最後に作りたいジャンルを選びます。" : "Start by pressing one target. Later levels add timing, collisions, and movement before you choose a genre to build."}</p></div><figure class="diagram-figure curriculum-map-figure"><img src="../assets/diagrams/${lang}_curriculum-flow.svg" alt="${esc(ja ? "LEVEL 01から専門トラックまでの学習の流れ" : "Learning flow from LEVEL 01 to specialist tracks")}" loading="lazy" decoding="async"><figcaption>${ja ? "矢印はおすすめの順番です。迷ったらLEVEL 01から、気になるゲームがあればそこから始めても大丈夫です。" : "Arrows show a suggested order. Start at LEVEL 01 if unsure, or jump to a game that interests you."}</figcaption></figure><nav class="curriculum-map-links" aria-label="${ja ? "カリキュラムへのリンク" : "Curriculum links"}">${links}<a href="tracks/visual-effects/"><span>${ja ? "VISUAL EFFECTS LAB · 見た目" : "VISUAL EFFECTS LAB · polish"}</span><b>→</b></a><a href="#specializations"><span>${ja ? "25専門トラック · 大きなゲーム" : "25 SPECIALIST TRACKS · bigger games"}</span><b>→</b></a></nav></section>\n<!-- curriculum-map:end -->`;
}

function injectHome(file, lang) {
  let html = fs.readFileSync(file, "utf8");
  const block = curriculumBlock(lang);
  const marker = /<!-- curriculum-map:start -->[\s\S]*?<!-- curriculum-map:end -->/;
  html = html.replace(marker, "");
  const anchor = "<!-- home-next-steps:start -->";
  html = html.replace(anchor, `${block}\n${anchor}`);
  fs.writeFileSync(file, html);
}

function pageRoute(file, lang) {
  const langRoot = path.join(webRoot, lang);
  return path.relative(langRoot, path.dirname(file)).replaceAll(path.sep, "/");
}

function injectLesson(file, lang) {
  let html = fs.readFileSync(file, "utf8");
  const route = pageRoute(file, lang);
  const isLesson = /<iframe[^>]+(?:play\/|lesson-game-frame)/i.test(html) || /<section[^>]+class="[^"]*\bplay(?:-panel|\s|\b)/i.test(html);
  if (!isLesson) return false;
  // Device scaling is introduced once at LEVEL 01. Repeating this generic
  // diagram hides the lesson-specific visuals that follow.
  const blocks = [];
  if (route === "games/tap-target") {
    const deviceAsset = path.relative(path.dirname(file), path.join(assetRoot, `${lang}_layout-comparison.svg`)).replaceAll(path.sep, "/");
    blocks.push(figure(deviceAsset, lang === "ja" ? "PC・タブレット・スマホでのゲーム表示比較" : "Game layout comparison on PC, tablet, and phone", lang === "ja" ? "同じキャンバスが画面幅に合わせて縮み、操作ボタンは押しやすい位置に並びます。" : "The same canvas scales to the available width while controls stay easy to reach.", "lesson-device-diagram"));
  }
  const diagramIds = routeDiagrams.get(route) || [];
  for (const diagramId of diagramIds) {
    const def = diagrams.find((item) => item.id === diagramId);
    if (!def) continue;
    const asset = path.relative(path.dirname(file), path.join(assetRoot, `${lang}_${diagramId}.svg`)).replaceAll(path.sep, "/");
    blocks.push(figure(asset, def.title[lang], lang === "ja" ? `${def.title[lang]}の見取り図です。コードを読む前に、数字と絵の関係を目で追ってみましょう。` : `A visual map of ${def.title[lang]}. Trace the numbers and pictures before reading the code.`, "lesson-concept-diagram"));
  }
  const block = `<!-- diagram-lesson:start -->\n${blocks.join("\n")}\n<!-- diagram-lesson:end -->`;
  const marker = /<!-- diagram-lesson:start -->[\s\S]*?<!-- diagram-lesson:end -->/;
  if (route === "games/tap-target") {
    html = html.replace(marker, "");
    const anchor = "<!-- feedback-note:games/tap-target:start -->";
    if (!html.includes(anchor)) return false;
    html = html.replace(anchor, `${block}\n${anchor}`);
  }
  else {
    html = html.replace(marker, "");
    const feedbackEnd = [...html.matchAll(/<!-- feedback-note:[^>]+:end -->/g)].at(0);
    const beginnerEnd = "<!-- END BEGINNER BRIDGE -->";
    if (feedbackEnd) {
      const at = feedbackEnd.index + feedbackEnd[0].length;
      html = html.slice(0, at) + `\n${block}` + html.slice(at);
    } else if (html.includes(beginnerEnd)) {
      html = html.replace(beginnerEnd, `${beginnerEnd}\n${block}`);
    } else {
      const play = html.search(/<section[^>]*class=["'][^"']*play[^"']*["'][^>]*>/i);
      if (play < 0) return false;
      const close = html.indexOf("</section>", play);
      if (close < 0) return false;
      const at = close + "</section>".length;
      html = html.slice(0, at) + `\n${block}` + html.slice(at);
    }
  }
  fs.writeFileSync(file, html);
  return true;
}

fs.mkdirSync(assetRoot, { recursive: true });
fs.mkdirSync(docsDiagramRoot, { recursive: true });
for (const def of diagrams) for (const lang of ["ja", "en"]) writeDiagram(def, lang);
for (const lang of ["ja", "en"]) fs.writeFileSync(path.join(assetRoot, `${lang}_layout-comparison.svg`), deviceSvg(lang));
injectHome(path.join(webRoot, "ja", "index.html"), "ja");
injectHome(path.join(webRoot, "en", "index.html"), "en");

let injected = 0;
for (const lang of ["ja", "en"]) {
  const langRoot = path.join(webRoot, lang);
  const walk = (dir) => {
    for (const name of fs.readdirSync(dir)) {
      const file = path.join(dir, name);
      if (fs.statSync(file).isDirectory()) walk(file);
      else if (name === "index.html" && injectLesson(file, lang)) injected++;
    }
  };
  walk(langRoot);
}
console.log(`Generated ${diagrams.length * 2 + 2} bilingual diagrams and injected lesson visuals into ${injected} page(s).`);
