#!/usr/bin/env node
/**
 * Generates the bilingual, non-playable unit-testing course.
 * Every lesson reads the real LEVEL 01–12 Update method and the real tested
 * pure functions, so the article cannot silently drift from the repository.
 * SPDX-License-Identifier: Apache-2.0
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const logicFiles = [
  fs.readFileSync(path.join(root, "internal/lessonlogic/rules.go"), "utf8"),
  fs.readFileSync(path.join(root, "internal/lessonlogic/core_updates.go"), "utf8"),
].join("\n");
const testFiles = [
  fs.readFileSync(path.join(root, "internal/lessonlogic/rules_test.go"), "utf8"),
  fs.readFileSync(path.join(root, "internal/lessonlogic/core_updates_test.go"), "utf8"),
].join("\n");

const lessons = [
  {
    slug: "tap-target", level: "01", source: "games/core/tap-target/main.go",
    pure: ["PointInCircle"], tests: ["TestPointInCircle"],
    phases: { ja: ["開始前と時間切れを先に処理", "残り時間を1 tick減らす", "押した瞬間だけ円との距離を判定", "命中なら得点・半径・的の位置を更新"], en: ["Handle not-started and time-up states first", "Subtract one tick from the timer", "Test the circle only on a fresh press", "On hit, update score, radius, and target position"] },
    ja: { title: "光る丸：入力と当たり判定を切り離す", lead: "Ebitengineが読むポインター入力と、数字だけで答えが出る円判定を分けます。", idea: "CursorPositionやJustPressedを偽物にする必要はありません。Updateが入力を読み、PointInCircleへ普通の座標を渡せば、境界だけを高速にテストできます。", manual: "押し心地、丸の見やすさ、制限時間の体感は実際に遊んで確認します。Drawの丸そのものはユニットテストしません。", cases: [["中心", "100,80", "true"], ["円周上", "130,80", "true"], ["1px外", "131,80", "false"]] },
    en: { title: "Tap Target: separate input from hit testing", lead: "Split pointer input owned by Ebitengine from a circle rule that only needs numbers.", idea: "Do not fake CursorPosition or JustPressed. Let Update read them and pass plain coordinates to PointInCircle, then test the boundary quickly.", manual: "Play-test tap feel, target readability, and timer feel. Do not unit-test the circle drawn by Draw.", cases: [["center", "100,80", "true"], ["on edge", "130,80", "true"], ["1 px outside", "131,80", "false"]] },
  },
  {
    slug: "timing-meter", level: "02", source: "games/core/timing-meter/main.go",
    pure: ["BouncedPosition", "TimingScore"], tests: ["TestBouncedPosition", "TestTimingScore"],
    phases: { ja: ["新しい押下で停止／次ラウンドを切り替える", "停止時に中心からの距離を採点", "停止中なら状態を進めず終了", "移動中だけ位置と速度を反射計算"], en: ["Use a fresh press to stop or start the next round", "Score distance from center when stopping", "Return without advancing while stopped", "Only while moving, advance and reflect position"] },
    ja: { title: "タイミングメーター：境界値を表にする", lead: "反射とPERFECT/GREAT/GOOD/MISSの境目を、二つの純粋関数として守ります。", idea: "見た目の線を止めるテストではなく、8と8.1、28と28.1のように境目の両側を直接渡します。", manual: "線の速さが楽しいか、中央が読みやすいかは人が確認します。Drawのピクセル位置は採点ルールではありません。", cases: [["PERFECT最後", "8", "100"], ["GREAT最初", "8.1", "50"], ["右端で反射", "99+3", "100,-3"]] },
    en: { title: "Timing Meter: make a table of boundaries", lead: "Protect reflection and PERFECT/GREAT/GOOD/MISS edges with two pure functions.", idea: "Do not test by watching a line stop. Pass both sides of each edge directly: 8 and 8.1, 28 and 28.1.", manual: "A person still judges whether speed feels fun and the center is readable. Draw pixels do not own scoring.", cases: [["last PERFECT", "8", "100"], ["first GREAT", "8.1", "50"], ["reflect at right", "99+3", "100,-3"]] },
  },
  {
    slug: "catch-stars", level: "03", source: "games/core/catch-stars/main.go",
    pure: ["ClassifyFallingObject"], tests: ["TestClassifyFallingObject"], support: "falling",
    phases: { ja: ["ゲームオーバー時は再開だけ受ける", "カゴを動かし出現タイマーを進める", "各星を落下させる", "取得／落下／継続を分類", "残す星だけでスライスを作り直す"], en: ["When over, accept restart only", "Move the basket and advance the spawn timer", "Advance each star downward", "Classify caught, missed, or kept", "Rebuild the slice with kept stars only"] },
    ja: { title: "星キャッチ：ループ中の三つの結果を型にする", lead: "星1個の結果を、残す・取った・落としたの三択へ分離します。", idea: "スライスの削除、得点、残機を一度にテストせず、まず1個の星がどの結果になるかを固定します。", manual: "星の量、落下の気持ちよさ、カゴの見た目はプレイ確認です。", cases: [["画面内", "y=719", "Keep"], ["画面下", "y=721", "Missed"], ["取得", "caught=true", "Caught"]] },
    en: { title: "Catch Stars: give the loop three typed outcomes", lead: "Separate one star's result into keep, caught, or missed.", idea: "Do not test slice deletion, score, and lives all at once. First freeze the answer for one star.", manual: "Play-test star density, falling feel, and basket presentation.", cases: [["visible", "y=719", "Keep"], ["below", "y=721", "Missed"], ["caught", "caught=true", "Caught"]] },
  },
  {
    slug: "flappy", level: "04", source: "game/main.go",
    pure: ["IntegrateGravity"], tests: ["TestIntegrateGravity"],
    phases: { ja: ["フレーム番号と押下を読む", "開始前／ゲームオーバーを処理", "押下で上向き速度を設定", "重力→速度→位置の順に積分", "パイプ移動・得点・再利用", "衝突とBEST更新"], en: ["Read tick count and fresh press", "Handle pre-start and game-over", "Apply flap velocity on press", "Integrate gravity, velocity, then position", "Move, score, and recycle pipes", "Resolve collision and best score"] },
    ja: { title: "Flappy：加速度の1 tickを固定する", lead: "重力を足してから位置を動かす順序を、数値の入出力としてテストします。", idea: "60回待つテストではなく、位置100・速度-7.4へ重力0.42を1回適用し、次の位置と速度を確認します。", manual: "羽ばたきの気持ちよさ、パイプの見やすさ、アニメーションはプレイ確認です。", cases: [["羽ばたき直後", "100,-7.4,+0.42", "93.02,-6.98"], ["落下中", "100,2,+0.42", "102.42,2.42"]] },
    en: { title: "Flappy: freeze one tick of acceleration", lead: "Test the order gravity→velocity→position as numeric input and output.", idea: "Do not wait 60 ticks. Apply 0.42 once to position 100 and velocity -7.4, then compare both next values.", manual: "Play-test flap feel, pipe readability, and animation.", cases: [["after flap", "100,-7.4,+0.42", "93.02,-6.98"], ["falling", "100,2,+0.42", "102.42,2.42"]] },
  },
  {
    slug: "pong", level: "05", source: "games/core/pong/main.go",
    pure: ["ExitScore"], tests: ["TestExitScore"],
    phases: { ja: ["プレイヤーのパドルを動かす", "CPUをボールへ追従させる", "ボール移動と壁／パドル反射", "上端／下端の通過を得点へ変換", "得点後に新しいサーブを作る"], en: ["Move the player paddle", "Make the CPU follow the ball", "Move and bounce off walls and paddles", "Convert top/bottom exit into a score", "Create a new serve after scoring"] },
    ja: { title: "Pong：画面外と得点を分ける", lead: "ボールのY座標から、得点なし・上側・下側を返す規則を取り出します。", idea: "serveには乱数があるため、境界判定と同じテストへ混ぜません。まず誰の得点かを決定的に確認します。", manual: "CPUの強さ、反射音、パドルへの当たり方はプレイ確認です。", cases: [["場内", "360", "得点なし"], ["上へ通過", "-20.1", "player"], ["下へ通過", "740.1", "CPU"]] },
    en: { title: "Pong: separate leaving the field from serving", lead: "Extract a rule that maps ball Y to no score, upper exit, or lower exit.", idea: "Serve uses randomness, so do not mix it into the boundary test. Decide who scored deterministically first.", manual: "Play-test CPU strength, bounce audio, and paddle contact feel.", cases: [["inside", "360", "no score"], ["exits top", "-20.1", "player"], ["exits bottom", "740.1", "CPU"]] },
  },
  {
    slug: "space-shooter", level: "06", source: "games/core/space-shooter/main.go",
    pure: ["AimVelocity"], tests: ["TestAimVelocity"],
    phases: { ja: ["終了状態・無敵時間・背景を更新", "プレイヤー入力と射撃", "味方弾／敵弾を移動", "敵AIと照準射撃", "弾×敵を逆順で解決", "敵弾×自機と画面外敵を解決", "全滅時に次ウェーブ生成"], en: ["Update end state, invincibility, and background", "Read player movement and shooting", "Move friendly and enemy bullets", "Run enemy AI and aimed shots", "Resolve shots vs enemies backwards", "Resolve enemy shots and escaped enemies", "Spawn a wave when none remain"] },
    ja: { title: "シューティング：照準ベクトルを純粋にする", lead: "敵から自機への差分を、指定した速さの弾ベクトルへ変換します。", idea: "敵スライスや弾生成を渡さず、dx・dy・speedだけを渡します。3-4-5の三角形なら期待値を手計算できます。", manual: "敵の出現量、弾の読みやすさ、爆発の手応えはプレイ確認です。", cases: [["3-4-5", "3,4,speed10", "6,8"], ["同じ位置", "0,0,speed10", "0,0"]] },
    en: { title: "Shooter: make aiming a pure vector rule", lead: "Turn enemy-to-player displacement into a bullet vector with a requested speed.", idea: "Pass only dx, dy, and speed—not enemy slices or spawning. A 3-4-5 triangle has an answer you can calculate by hand.", manual: "Play-test enemy density, bullet readability, and explosion feedback.", cases: [["3-4-5", "3,4,speed10", "6,8"], ["same point", "0,0,speed10", "0,0"]] },
  },
  {
    slug: "breakout", level: "07", source: "games/core/breakout/main.go",
    pure: ["SpendLife"], tests: ["TestSpendLife"],
    phases: { ja: ["パドル移動", "壁とパドルの反射", "ブロック衝突と削除", "落球なら残機を減らす", "残機ありはserve、0なら全体reset", "全ブロック破壊でもreset"], en: ["Move paddle", "Bounce off world and paddle", "Resolve and remove brick hits", "Spend a life when the ball falls", "Serve if lives remain; reset at zero", "Reset after clearing every brick"] },
    ja: { title: "ブロック崩し：残機の境目を小さくする", lead: "落球後の残機とゲームオーバー判定を、描画もボールも知らない関数へ移します。", idea: "newGameやserveまでユニットテストに含めず、1→0だけがゲームオーバーになる約束を守ります。", manual: "反射角、ブロック配置、クリア演出はプレイ確認です。", cases: [["残機3", "3", "2,false"], ["最後の1機", "1", "0,true"]] },
    en: { title: "Breakout: shrink the life boundary", lead: "Move post-miss lives and game-over decision into a function that knows no drawing or ball.", idea: "Do not include newGame and serve in this unit. Protect the promise that 1→0 ends the game.", manual: "Play-test bounce angles, brick layout, and clear presentation.", cases: [["three lives", "3", "2,false"], ["last life", "1", "0,true"]] },
  },
  {
    slug: "snake", level: "08", source: "games/core/snake/main.go",
    pure: ["SnakeStepInterval"], tests: ["TestSnakeStepInterval"],
    phases: { ja: ["終了時は再開だけ受ける", "次方向を入力バッファへ保存", "tickカウンターを進める", "得点から移動間隔を計算", "移動tickでなければ終了", "頭・餌・尾・自己衝突を1マス進める"], en: ["When over, accept restart only", "Store next direction in the input buffer", "Advance tick counter", "Convert score into move interval", "Return on a non-move tick", "Advance head, food, tail, and self collision one cell"] },
    ja: { title: "Snake：速さの式を時間制御から外す", lead: "得点から『何tickごとに1マス進むか』だけを計算する関数にします。", idea: "frame%waitのタイミングはUpdateが担当し、waitの下限4と難易度曲線は表テストで確認します。", manual: "曲がる操作感、盤面サイズ、餌の見つけやすさはプレイ確認です。", cases: [["開始", "score0", "10"], ["得点3", "score3", "9"], ["高速上限", "score18/99", "4"]] },
    en: { title: "Snake: remove the speed formula from time control", lead: "Make a function that only converts score into ticks per grid move.", idea: "Update owns frame%wait timing; a table test protects the floor of 4 and the difficulty curve.", manual: "Play-test turning feel, board size, and food visibility.", cases: [["start", "score0", "10"], ["score three", "score3", "9"], ["speed floor", "score18/99", "4"]] },
  },
  {
    slug: "sokoban", level: "09", source: "games/core/sokoban/main.go",
    pure: ["AdvanceTween"], tests: ["TestAdvanceTween"],
    phases: { ja: ["Rで全体リセット", "停止中だけUndo", "クリア後は再開だけ", "移動中は補間だけ進め入力を無視", "完了時にplayer／boxを確定しクリア判定", "停止中だけ新しい移動を計画"], en: ["Reset everything with R", "Allow undo only while idle", "After clear, accept restart only", "While moving, advance tween and ignore input", "On completion, commit player/box and check clear", "Only while idle, plan a new move"] },
    ja: { title: "倉庫番：補間の完了を状態遷移にする", lead: "0〜1のprogressを進め、完了したかを返す関数へ分離します。", idea: "Drawはprogressを読むだけです。playerやboxの確定はUpdateに残し、0.98+0.14が必ず1で完了することをテストします。", manual: "滑らかさ、箱の重さ、ゴールの見やすさはプレイ確認です。", cases: [["途中", "0.28+0.14", "0.42,false"], ["完了", "0.98+0.14", "1,true"]] },
    en: { title: "Sokoban: make tween completion a transition", lead: "Extract a function that advances 0..1 progress and reports completion.", idea: "Draw only reads progress. Update still commits player and box; the test proves 0.98+0.14 clamps to 1 and completes.", manual: "Play-test smoothness, box weight, and goal readability.", cases: [["moving", "0.28+0.14", "0.42,false"], ["complete", "0.98+0.14", "1,true"]] },
  },
  {
    slug: "platformer", level: "10", source: "games/core/platformer/main.go",
    pure: ["HorizontalVelocity", "VerticalVelocity"], tests: ["TestHorizontalVelocity", "TestVerticalVelocity"],
    phases: { ja: ["終了時は再開だけ", "左右・ジャンプ入力を読む", "加速／摩擦／上限とジャンプ／重力を計算", "X移動後にX衝突を解消", "Y移動後にY衝突と接地を解消", "コイン・落下・ゴールを判定", "カメラ状態を追従"], en: ["When finished, accept restart only", "Read left, right, and jump", "Calculate acceleration/friction/cap and jump/gravity", "Move X then resolve X collision", "Move Y then resolve Y collision and grounding", "Resolve coins, falling, and goal", "Advance camera follow state"] },
    ja: { title: "Platformer：入力→速度と衝突を段階化する", lead: "横速度と縦速度を先に純粋計算し、そのあとUpdateがX衝突・Y衝突の順に状態へ適用します。", idea: "巨大なUpdateを丸ごと呼ばず、左右同時入力、空中ジャンプ、落下上限という境界を個別に固定します。", manual: "ジャンプの高さ、床の角での感触、カメラ酔いは必ずプレイ確認します。", cases: [["右加速", "vx0,right", "0.65"], ["摩擦", "vx5,no input", "3.9"], ["接地ジャンプ", "jump+ground", "jumpSpeed+gravity"], ["空中ジャンプ", "jump+air", "無視"]] },
    en: { title: "Platformer: stage input→velocity before collision", lead: "Compute horizontal and vertical velocity purely, then let Update apply X collision followed by Y collision.", idea: "Do not call the giant Update. Freeze opposite inputs, air jump, and fall cap independently.", manual: "Always play-test jump height, corner feel, and camera comfort.", cases: [["right accel", "vx0,right", "0.65"], ["friction", "vx5,no input", "3.9"], ["ground jump", "jump+ground", "jumpSpeed+gravity"], ["air jump", "jump+air", "ignored"]] },
  },
  {
    slug: "dungeon", level: "11", source: "games/core/dungeon/main.go",
    pure: ["EnemyMode", "AimVelocity"], tests: ["TestEnemyMode", "TestAimVelocity"],
    phases: { ja: ["終了時は再開だけ", "プレイヤーをX→Yに移動", "攻撃タイマーと新規攻撃入力", "敵ごとに距離から状態を遷移", "徘徊または追跡速度を決定", "敵をX→Yに移動", "剣命中・接触ダメージ", "全滅＋出口でクリア"], en: ["When finished, accept restart only", "Move player on X then Y", "Advance attack timer and fresh attack input", "For each enemy, transition mode from distance", "Choose wander or chase velocity", "Move enemy on X then Y", "Resolve sword hit and contact damage", "Clear after all enemies and exit"] },
    ja: { title: "ダンジョン：AIのヒステリシスをテストする", lead: "165未満で追跡、230超で徘徊へ戻り、その間は今の状態を保つ規則を分離します。", idea: "敵を実際に歩かせず、距離200で現在状態が保持されることをテストすると、境界でのガタガタ切替を防げます。", manual: "敵の怖さ、剣の間合い、迷路の読みやすさはプレイ確認です。", cases: [["接近", "164.9", "chase"], ["中間・徘徊中", "200,current0", "wander"], ["中間・追跡中", "200,current1", "chase"], ["遠い", "230.1", "wander"]] },
    en: { title: "Dungeon: test AI hysteresis", lead: "Extract: chase below 165, return to wander above 230, keep the current mode between them.", idea: "Without walking an enemy, test that distance 200 preserves its current mode and prevents boundary flicker.", manual: "Play-test enemy threat, sword reach, and maze readability.", cases: [["near", "164.9", "chase"], ["band while wandering", "200,current0", "wander"], ["band while chasing", "200,current1", "chase"], ["far", "230.1", "wander"]] },
  },
  {
    slug: "bullet-hell", level: "12", source: "games/core/bullet-hell/main.go",
    pure: ["CircleHit", "Outside"], tests: ["TestCircleHit", "TestOutside"],
    phases: { ja: ["終了時は再開だけ", "tickと無敵時間を更新", "プレイヤー移動", "自機弾・リング・扇弾を周期生成", "全弾を逆順に移動", "画面外判定", "自機弾×ボス／敵弾×自機", "弾削除と勝敗更新"], en: ["When finished, accept restart only", "Advance tick and invincibility", "Move player", "Periodically spawn player, ring, and fan shots", "Move every bullet backwards", "Test outside bounds", "Resolve player shot vs boss and enemy shot vs player", "Remove bullets and update win/lose"] },
    ja: { title: "弾幕：大量ループの判定を1発へ縮める", lead: "何百発のスライスではなく、円2個の重なりと1点の画面外判定をテストします。", idea: "Updateは逆順削除とHP更新を担当します。純粋関数には1発分の数字だけを渡し、接触境界と余白30を固定します。", manual: "弾幕の美しさ、避けられる隙間、被弾の見やすさはプレイ確認です。", cases: [["円が重なる", "distance9,sum10", "true"], ["円が接するだけ", "distance10,sum10", "false"], ["余白上", "x=-30", "inside"], ["余白外", "x=-30.1", "outside"]] },
    en: { title: "Bullet Hell: shrink a huge loop to one bullet", lead: "Test two-circle overlap and one point outside bounds—not a slice of hundreds of bullets.", idea: "Update owns reverse deletion and HP. Give pure functions one bullet's numbers and freeze contact and the 30 px margin.", manual: "Play-test pattern beauty, dodge gaps, and hit readability.", cases: [["circles overlap", "distance9,sum10", "true"], ["circles only touch", "distance10,sum10", "false"], ["on margin", "x=-30", "inside"], ["past margin", "x=-30.1", "outside"]] },
  },
];

if (lessons.length !== 12 || lessons.some((lesson, index) => lesson.level !== String(index + 1).padStart(2, "0"))) {
  throw new Error("testing course must cover LEVEL 01–12 exactly once and in order");
}

function extractFunction(source, signatureStart) {
  const start = source.indexOf(signatureStart);
  if (start < 0) throw new Error(`missing function: ${signatureStart}`);
  const brace = source.indexOf("{", start);
  let depth = 0;
  for (let i = brace; i < source.length; i++) {
    if (source[i] === "{") depth++;
    if (source[i] === "}" && --depth === 0) return source.slice(start, i + 1);
  }
  throw new Error(`unterminated function: ${signatureStart}`);
}

function namedFunction(source, name) {
  return extractFunction(source, `func ${name}(`);
}

function updateFunction(file) {
  return extractFunction(fs.readFileSync(path.join(root, file), "utf8"), "func (g *game) Update() error");
}

function supportCode(kind) {
  if (kind !== "falling") return "";
  const start = logicFiles.indexOf("type FallingOutcome int");
  const end = logicFiles.indexOf("// ClassifyFallingObject", start);
  return logicFiles.slice(start, end).trim();
}

function goFile(functions, support = "") {
  const body = [support, ...functions].filter(Boolean).join("\n\n");
  const imports = body.includes("math.") ? '\nimport "math"\n' : "";
  return `package lessonlogic${imports}\n${body}\n`;
}

function testFile(functions) {
  const body = functions.join("\n\n");
  const imports = body.includes("math.") ? 'import (\n\t"math"\n\t"testing"\n)' : 'import "testing"';
  return `package lessonlogic\n\n${imports}\n\n${body}\n`;
}

const esc = (s) => String(s).replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");

function shell({ lang, depth, title, desc, body }) {
  const ja = lang === "ja";
  const prefix = "../".repeat(depth);
  const other = ja ? "en" : "ja";
  const route = depth === 3 ? "guides/testing/" : `guides/testing/${title.slug}/`;
  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><meta name="description" content="${esc(desc)}"><title>${esc(title.text)} | Ebi Showcase</title><link rel="stylesheet" href="${prefix}style.css"></head><body class="testing-guide"><header class="nav"><a class="brand" href="${prefix}${lang}/"><span>EBI</span> SHOWCASE</a><nav><a href="${depth === 3 ? "./" : "../"}">TESTING</a><a class="lang" href="${prefix}${other}/${route}" lang="${other}">${ja ? "EN" : "日本語"}</a></nav></header><main>${body}</main><footer><span>EBI SHOWCASE</span><span>GO + EBITENGINE</span><span>APACHE-2.0</span></footer><script src="${prefix}learn.js"></script></body></html>`;
}

function hub(lang) {
  const ja = lang === "ja";
  const c = ja ? {
    title: "LEVEL 01〜12のUpdateをユニットテストする", desc: "Drawを原則テストせず、EbitengineのUpdateから純粋なゲームルールを分離してLEVEL 01〜12をGoで網羅的にテストする教材。",
    h1: "Drawを試す前に、<br><em>Updateのルールを守る。</em>",
    lead: "Drawは現在のgameを画面へ投影するだけなので、原則としてユニットテストの対象にしません。入力・時間・乱数・衝突・得点を持つUpdateを、薄い接続役と純粋なルールへ分けます。下の12章では、実際のUpdateを一行も省略せず読み、そのまま実行できるテストを書きます。",
    input: "Updateが入力を読む", pure: "純粋関数へ普通の値を渡す", state: "Updateがgameへ結果を書く", noDraw: "Drawはgameを読むだけ",
    why: "なぜDraw()は原則テストしない？", whyBody: "同じgame・素材・画面サイズなら同じ絵を出す、という参照透過なDrawを守れば、ルールの正しさはgameの状態で確認できます。ピクセル比較はフォント・GPU・OS差で壊れやすく、当たり判定や得点の答えにもなりません。独自シェーダーや画像生成器を守りたい場合は別枠のゴールデン画像テストを使いますが、それはUpdateのユニットテストとは分けます。",
    update: "Update()をどう構造化する？", updateBody: "①Ebitengineから入力を読む、②普通の数値・小さな構造体へ変換、③純粋関数で次の答えを計算、④gameへ結果を保存、の順にします。純粋関数はEbitengineをimportせず、時刻や乱数も引数で受け取ります。自分のGoプロジェクトにinternal/lessonlogicを作ればよく、このサイトのcloneは不要です。",
    cards: "12本すべてを、実コードで分解する", cardsLead: "各ページに、現行のUpdate全文、責務の全手順、抽出した関数の完全なGoファイル、完全な_test.go、境界ケース、残すべき手動確認を掲載します。", read: "全文を読む →",
  } : {
    title: "Unit-test Update across LEVEL 01–12", desc: "A complete guide to leaving Draw untested by default and extracting pure rules from every LEVEL 01–12 Update for Go unit tests.",
    h1: "Protect Update rules<br><em>before testing pictures.</em>",
    lead: "Draw projects the current game and is normally not a unit-test target. Split input, time, randomness, collision, and scoring in Update into a thin adapter plus pure rules. Each of the 12 chapters prints the real Update without omissions and a runnable complete test.",
    input: "Update reads input", pure: "pass plain values to pure rules", state: "Update writes results to game", noDraw: "Draw only reads game",
    why: "Why not unit-test Draw() by default?", whyBody: "When Draw is a referentially transparent projection, the same game, assets, and viewport produce the same picture; rule correctness is visible in game state. Pixel comparisons are brittle across fonts, GPUs, and operating systems and do not prove collision or scoring. Protect a custom shader or image generator with a separate golden-image suite when needed—not with Update unit tests.",
    update: "How should Update() be structured?", updateBody: "① Read Ebitengine input. ② Convert it to plain values or small structs. ③ Calculate the next answer in pure functions. ④ Store results in game. Pure functions do not import Ebitengine and receive time/random values as arguments. Create internal/lessonlogic in your own Go project; cloning this site is not required.",
    cards: "Decompose all twelve real games", cardsLead: "Every page includes the full current Update, every responsibility in order, a complete pure Go file, a complete _test.go, boundary cases, and the manual checks that remain.", read: "READ ALL →",
  };
  const cards = lessons.map((lesson) => { const t = lesson[lang]; return `<a class="test-course-card" href="${lesson.slug}/"><span>${lesson.level}</span><h3>${t.title}</h3><p>${t.lead}</p><strong>${c.read}</strong></a>`; }).join("");
  const body = `<section class="test-hero"><p class="eyebrow">SPECIAL GUIDE / UNIT TESTING</p><h1>${c.h1}</h1><p>${c.lead}</p></section><section class="test-boundary" aria-label="Update and Draw testing boundary"><div><small>01 / INPUT ADAPTER</small><b>${c.input}</b></div><i>→</i><div class="is-pure"><small>02 / TEST HERE</small><b>${c.pure}</b><em>go test / no window</em></div><i>→</i><div><small>03 / STATE</small><b>${c.state}</b><em>${c.noDraw}</em></div></section><section class="test-explain test-principles"><div><p class="eyebrow">THE DEFAULT</p><h2>${c.why}</h2><p>${c.whyBody}</p></div><div><p class="eyebrow">THE STRUCTURE</p><h2>${c.update}</h2><p>${c.updateBody}</p></div></section><section class="test-terminal"><code>$ go test ./internal/lessonlogic</code><strong>✓ rules, not pixels</strong></section><section class="test-course"><p class="eyebrow">LEVEL 01–12 / COMPLETE MAP</p><h2>${c.cards}</h2><p>${c.cardsLead}</p><div class="test-course-grid">${cards}</div></section><nav class="test-guide-links"><a href="tap-target/">${ja ? "LEVEL 01から始める →" : "START LEVEL 01 →"}</a><a href="../../">${ja ? "← ホームへ" : "← HOME"}</a></nav>`;
  return shell({ lang, depth: 3, title: { text: c.title }, desc: c.desc, body });
}

function lessonPage(lang, index) {
  const lesson = lessons[index];
  const t = lesson[lang];
  const ja = lang === "ja";
  const update = updateFunction(lesson.source);
  const pureFunctions = lesson.pure.map((name) => namedFunction(logicFiles, name));
  const tests = lesson.tests.map((name) => namedFunction(testFiles, name));
  const pure = goFile(pureFunctions, supportCode(lesson.support));
  const test = testFile(tests);
  const rows = t.cases.map((r) => `<tr><th>${esc(r[0])}</th><td><code>${esc(r[1])}</code></td><td><code>${esc(r[2])}</code></td></tr>`).join("");
  const phases = lesson.phases[lang].map((phase, i) => `<li><span>${String(i + 1).padStart(2, "0")}</span><p>${phase}</p></li>`).join("");
  const prev = index === 0 ? "../" : `../${lessons[index - 1].slug}/`;
  const next = index === lessons.length - 1 ? "../" : `../${lessons[index + 1].slug}/`;
  const c = ja ? {
    complete: "現行Updateの全手順", completeBody: "省略記号はありません。入力から早期return、ループ、勝敗まで、現在ゲームで動く順番です。この中から抽出した純粋関数を直接テストします。",
    update: "REAL GO / Update()全文", pure: "EXTRACTED / 完全な純粋ロジック", test: "TEST / 完全な_test.go", table: "境界を先に決める", manual: "テストのあとも人が確認すること", run: "自分のGoプロジェクトに上のファイルを作って実行", back: "← 前へ", forward: "次へ →",
  } : {
    complete: "Every step in the current Update", completeBody: "There are no ellipses. This is the running order from input and early returns through loops and win/lose. Unit-test the extracted pure rules directly.",
    update: "REAL GO / complete Update()", pure: "EXTRACTED / complete pure logic", test: "TEST / complete _test.go", table: "Choose boundaries first", manual: "What a person still checks", run: "Create these files in your own Go project, then run", back: "← PREVIOUS", forward: "NEXT →",
  };
  const logicPath = index < 2 ? "internal/lessonlogic/rules.go" : "internal/lessonlogic/core_updates.go";
  const testPath = index < 2 ? "internal/lessonlogic/rules_test.go" : "internal/lessonlogic/core_updates_test.go";
  const body = `<section class="test-step-hero"><a href="../">TESTING GUIDE</a><p class="eyebrow">LEVEL ${lesson.level} / 12</p><h1>${t.title}</h1><p>${t.lead}</p></section><section class="test-rule-strip"><span>EBITENGINE INPUT</span><i>→</i><strong>PURE RULE + go test</strong><i>→</i><span>GAME STATE</span></section><section class="test-explain"><div><p class="eyebrow">WHY THIS SEAM?</p><h2>${t.idea}</h2><p><strong>Draw:</strong> ${t.manual}</p></div><div class="test-case-table"><p class="eyebrow">${c.table}</p><table><tbody>${rows}</tbody></table></div></section><section class="test-update-map"><div><p class="eyebrow">NO OMISSIONS</p><h2>${c.complete}</h2><p>${c.completeBody}</p></div><ol>${phases}</ol></section><section class="test-code-full"><div><p class="eyebrow">${c.update}</p><code>${esc(lesson.source)}</code></div><pre><code>${esc(update)}</code></pre></section><section class="test-code-compare"><article><p>${c.pure}</p><code>${logicPath}</code><pre><code>${esc(pure)}</code></pre></article><article class="is-after"><p>${c.test}</p><code>${testPath}</code><pre><code>${esc(test)}</code></pre></article></section><section class="test-run"><p>${c.run}</p><code>go test ./internal/lessonlogic</code><strong>✓ PASS / NO WINDOW / NO DRAW</strong></section><section class="test-challenge"><p class="eyebrow">MANUAL CHECK IS STILL REAL</p><h2>${c.manual}</h2><p>${t.manual}</p></section><nav class="test-pager"><a href="${prev}">${c.back}</a><span>${index + 1} / ${lessons.length}</span><a href="${next}">${c.forward}</a></nav>`;
  return shell({ lang, depth: 4, title: { text: t.title, slug: lesson.slug }, desc: t.lead, body });
}

function redirectPage(lang, slug, target) {
  const ja = lang === "ja";
  return shell({ lang, depth: 4, title: { text: ja ? "テスト教材を移動しました" : "Testing lesson moved", slug }, desc: ja ? "LEVEL 01〜12の新しいテスト教材へ移動します。" : "Continue to the new LEVEL 01–12 testing guide.", body: `<section class="test-step-hero"><p class="eyebrow">GUIDE UPDATED</p><h1>${ja ? "LEVEL 01〜12の章へ統合しました。" : "Now part of the LEVEL 01–12 course."}</h1><p><a href="../${target}/">${ja ? "新しい章を開く →" : "OPEN THE NEW LESSON →"}</a></p></section>` });
}

for (const lang of ["ja", "en"]) {
  const dir = path.join(root, "web", lang, "guides", "testing");
  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(path.join(dir, "index.html"), hub(lang));
  lessons.forEach((lesson, i) => {
    const lessonDir = path.join(dir, lesson.slug);
    fs.mkdirSync(lessonDir, { recursive: true });
    fs.writeFileSync(path.join(lessonDir, "index.html"), lessonPage(lang, i));
  });
  for (const [old, target] of Object.entries({ "pure-functions": "tap-target", "table-tests": "timing-meter", "state-transitions": "platformer", "regression-tests": "bullet-hell", "readable-tests": "tap-target" })) {
    const oldDir = path.join(dir, old);
    fs.mkdirSync(oldDir, { recursive: true });
    fs.writeFileSync(path.join(oldDir, "index.html"), redirectPage(lang, old, target));
  }
}

for (const lang of ["ja", "en"]) {
  const file = path.join(root, "web", lang, "index.html");
  let html = fs.readFileSync(file, "utf8");
  const ja = lang === "ja";
  const block = `<!-- testing-guide-home:start -->\n<section class="architecture-promo testing-promo"><div><p class="eyebrow">SPECIAL GUIDE / UNIT TESTING</p><h2>${ja ? "Drawではなく、<br>12本のUpdateをテストする。" : "Test all twelve Updates,<br>not Draw pixels."}</h2><p>${ja ? "LEVEL 01〜12のUpdate全文を読み、Ebitengine非依存の純粋ルールへ分離し、完全なGoテストを書く網羅コースです。" : "Read every LEVEL 01–12 Update, extract Ebitengine-free rules, and write complete runnable Go tests."}</p></div><a href="guides/testing/"><span>READ THE GUIDE</span><strong>${ja ? "ゲームのユニットテスト完全編" : "Complete unit testing for games"}</strong><b>→</b></a></section>\n<!-- testing-guide-home:end -->`;
  const re = /<!-- testing-guide-home:start -->[\s\S]*?<!-- testing-guide-home:end -->/;
  html = re.test(html) ? html.replace(re, block) : html.replace("<!-- visual-effects-home:start -->", `${block}\n<!-- visual-effects-home:start -->`);
  fs.writeFileSync(file, html);
}

console.log("Generated the complete bilingual LEVEL 01–12 unit-testing course.");
