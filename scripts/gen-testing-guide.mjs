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

const aaaCourse = {
  "tap-target": {
    ja: {
      title: "AAA 01：世界を一文で固定する",
      intro: "AAAの最初の仕事は、テストしたい世界を小さくすることです。「中心が(100, 80)、半径30の的に、点(130, 80)を置いた」と一文で言えるArrangeなら、失敗したときに余計な状態を疑わずに済みます。",
      arrange: "表の1行が一つの前提です。点の座標だけを変え、的の中心と半径は固定します。game、入力機器、ウインドウは作りません。",
      act: "PointInCircleを一度だけ呼びます。一回のActに一回の質問だけを置くため、「この点は円の中か」以外の処理は混ざりません。",
      assert: "戻り値一つをwantInsideと比較します。失敗は「円周を内側に含める契約が壊れた」のように、守れなかったルールを直接示します。",
      why: "現行の4サブテストは、どれも小さな座標の前提→PointInCircleを1回→boolを1回比較、というAAAです。最初の章では、この三段を目で追える小ささを基準にします。",
      stories: ["中心は内側になる", "円周上も内側になる", "円周から1px外は外側になる"],
      trap: "newGameを作ってクリックを偽装し、UpdateとDrawまで通して丸の色を見るテストは、何が壊れたのかを曖昧にします。円判定のテストには円判定だけを置きます。",
      practice: "上端の点(100, 50)を加えてください。既存の右端とは向きが違いますが、Arrange・Act・Assertの形は同じなので新しい表の1行にできます。",
    },
    en: {
      title: "AAA 01: describe the whole world in one sentence",
      intro: "AAA first makes the tested world small. An Arrange you can say as “a radius-30 target centered at (100, 80), with the point at (130, 80)” leaves no unrelated state to investigate when the test fails.",
      arrange: "One table row is one premise. Change only the point and keep the target center and radius fixed. Do not create a game, input device, or window.",
      act: "Call PointInCircle exactly once. One Act asks one question, so nothing is mixed with “is this point inside the circle?”",
      assert: "Compare the one result with wantInside. A failure directly identifies a broken contract, such as whether the edge belongs to the circle.",
      why: "All four current subtests follow AAA: small coordinate premise, one PointInCircle call, and one bool comparison. This chapter makes that easy-to-scan shape the baseline.",
      stories: ["the center is inside", "the exact edge is inside", "one pixel past the edge is outside"],
      trap: "A test that creates newGame, fakes a click, runs Update and Draw, and inspects color hides the cause of failure. A circle-rule test contains only the circle rule.",
      practice: "Add the top-edge point (100, 50). It approaches from another direction, but its Arrange, Act, and Assert match the existing cases, so it becomes one more table row.",
    },
  },
  "timing-meter": {
    ja: {
      title: "AAA 02：境界の両側を別の物語にする",
      intro: "境界値テストでは「8ならPERFECT」と「8.1ならGREAT」を一組にします。ただし一つのサブテストで2回Actするのではなく、隣り合う二つのAAAとして書くから、どちら側の契約が壊れたかが名前だけで分かります。",
      arrange: "反射ではpositionとspeed、採点ではdistanceだけを1行に置きます。8、8.1、28、28.1のように境界そのものと直後を隣に並べます。",
      act: "各サブテストはBouncedPositionまたはTimingScoreのどちらか一方を一度だけ呼びます。異なるルールは別のTest関数へ分けます。",
      assert: "位置と速度、得点とラベルは同じActが返す一つの契約です。ただし比較は分け、どちらが違ったかを失敗文で特定します。",
      why: "現行テストは反射4件と採点7件です。通常移動・ちょうど端・端越え、各ランクの境界と直後を独立したAAAにしたため、条件式の&lt;と&lt;=の取り違えを局所化できます。",
      stories: ["右端へちょうど着く間は向きを保つ", "右端を越えたら位置を止めて反転する", "GREATの直後はGOODになる"],
      trap: "一つのt.Runの中で8、8.1、28、28.1を順に呼ぶと、最初の失敗で後続を読めず、サブテスト名も一つしか残りません。境界の片側ごとにAAAを完結させます。",
      practice: "新しく追加した28.1の行を28へ変えてテストを失敗させ、名前・got・wantだけで境界のどちら側を直すべきか説明してから戻してください。",
    },
    en: {
      title: "AAA 02: make each side of a boundary its own story",
      intro: "Boundary tests pair “8 is PERFECT” with “8.1 is GREAT.” They are two neighboring AAA tests—not two Acts in one subtest—so the failed side of the contract is obvious from its name.",
      arrange: "A bounce row holds position and speed; a scoring row holds distance. Put the boundary and the first value beyond it next to each other: 8, 8.1, 28, and 28.1.",
      act: "Each subtest calls either BouncedPosition or TimingScore exactly once. Different rules live in different Test functions.",
      assert: "Position plus speed, or points plus label, form one returned contract. Compare their fields separately so the failure says which part differs.",
      why: "The current suite has four bounce cases and seven scoring cases. Normal motion, exact edge, crossing, and both sides of rank boundaries are independent AAA stories, localizing a mistaken &lt; versus &lt;=.",
      stories: ["landing exactly on the right edge keeps direction", "crossing the edge clamps and reverses", "the first value beyond GREAT becomes GOOD"],
      trap: "Calling 8, 8.1, 28, and 28.1 inside one t.Run loses later evidence after the first failure and leaves only one scenario name. Complete one AAA cycle per side.",
      practice: "Temporarily change the new 28.1 row to 28, observe the failure, and explain the broken side using only the name, got, and want before restoring it.",
    },
  },
  "catch-stars": {
    ja: {
      title: "AAA 03：巨大なループを1個の判断へ縮める",
      intro: "ゲームのUpdateには星の移動、取得、削除、得点、残機が同居します。AAAはその全部を再現せず、「すでに移動した星1個をどう分類するか」という最小のActを見つける設計道具です。",
      arrange: "caughtとyだけで星1個の状況を作ります。特にcaught=trueかつ画面下という行は、二条件が重なったときの優先順位を明示します。",
      act: "ClassifyFallingObjectを一度呼び、スライスの走査や削除はUpdate側の責務として持ち込みません。",
      assert: "boolを二つ調べず、Keep・Caught・Missedの一つの型付き結果を比較します。期待する契約が一語で読めます。",
      why: "5サブテストが通常・下端・落下・取得・優先順位を別々に守っています。AAAがテストしやすい純粋関数の形まで実装を導く例です。",
      stories: ["画面内の星は残す", "画面下の星は落としたと分類する", "取得と落下が同時なら取得を優先する"],
      trap: "星を10個用意してUpdateを回し、残ったスライス長だけをAssertすると、分類・削除・得点のどこが壊れたか分かりません。",
      practice: "y=720の下端ちょうどがKeepになる既存行を説明し、その隣にy=720.1のMissedを追加して境界の両側を完成させてください。",
    },
    en: {
      title: "AAA 03: shrink a giant loop to one decision",
      intro: "The game Update mixes movement, collection, removal, score, and lives. AAA does not recreate all of it; it finds the smallest Act: classify one already-advanced star.",
      arrange: "Use only caught and y to describe one object. The caught-and-below-bottom row explicitly defines precedence when two conditions are true.",
      act: "Call ClassifyFallingObject once. Slice traversal and removal remain Update responsibilities rather than test setup.",
      assert: "Compare one typed Keep, Caught, or Missed result instead of inspecting two booleans. The contract reads as one word.",
      why: "Five subtests independently protect ordinary play, the bottom edge, a miss, a catch, and precedence. AAA is also design pressure toward a testable pure function.",
      stories: ["a visible star stays", "a star below the screen is missed", "a catch wins when catch and miss are both possible"],
      trap: "Creating ten stars, running Update, and asserting only the final slice length cannot identify whether classification, removal, or scoring broke.",
      practice: "Explain why the existing y=720 row is Keep, then add y=720.1 as Missed to complete both sides of the boundary.",
    },
  },
  flappy: {
    ja: {
      title: "AAA 04：1 tickの前後だけを守る",
      intro: "物理は長時間動かすほど原因が見えにくくなります。Arrangeを直前の位置と速度、Actを1 tickの積分、Assertを直後の位置と速度にすると、更新順序そのものが読み取れます。",
      arrange: "position・velocity・gravityを普通の数値で置きます。乱数パイプ、入力、60 tickのループは要りません。",
      act: "IntegrateGravityを一度だけ呼びます。「重力を速度へ足し、その新しい速度で位置を進める」という一回の状態遷移です。",
      assert: "nextPositionとnextVelocityを別々に許容誤差つきで比較します。片方しか見ないテストでは更新順序を守れません。",
      why: "上昇・下降・重力ゼロの3サブテストで、符号が違っても一回の積分契約が保たれることを確認しました。長いシミュレーションより失敗地点が一意です。",
      stories: ["上向き速度は重力で弱まってから位置へ反映される", "下向き速度は増えてから位置へ反映される", "重力ゼロなら速度を保って1 tick進む"],
      trap: "60回Actして最終Yだけを比べると、1 tick目の誤差がどこで生まれたか分からず、途中の間違いが偶然相殺されることもあります。",
      practice: "position=0、velocity=-1、gravity=1の行を加え、次の速度0と次の位置0の両方をAssertしてください。",
    },
    en: {
      title: "AAA 04: protect only one tick before and after",
      intro: "Physics becomes harder to diagnose the longer it runs. Arrange the position and velocity before one tick, Act with one integration, then Assert both values immediately after it.",
      arrange: "Use plain position, velocity, and gravity values. No random pipes, input, or sixty-tick loop is needed.",
      act: "Call IntegrateGravity once: add gravity to velocity, then advance position with that new velocity.",
      assert: "Compare nextPosition and nextVelocity separately with a tolerance. Inspecting only one cannot protect update order.",
      why: "Three subtests cover upward, downward, and zero-gravity motion. A one-tick contract gives each failure one origin instead of hiding it in a long simulation.",
      stories: ["upward speed weakens before position advances", "downward speed grows before position advances", "zero gravity preserves speed for one tick"],
      trap: "Acting sixty times and checking only final Y hides where the first error arose and can even let mistakes cancel each other.",
      practice: "Add position=0, velocity=-1, gravity=1 and Assert both the next velocity of 0 and next position of 0.",
    },
  },
  pong: {
    ja: {
      title: "AAA 05：『何も起きない』も契約にする",
      intro: "得点した場面だけではなく、場内や境界ちょうどで得点しないこともゲームルールです。AAAではno-opをwant=-1として明示し、誤って早く得点する回帰を捕まえます。",
      arrange: "ボールのyを場内、上端ちょうど、上端越え、下端ちょうど、下端越えに置きます。serveの乱数は別の責務なので除きます。",
      act: "ExitScoreを一度呼びます。ボール移動や得点加算まで進めず、「どちらが得点か」だけを質問します。",
      assert: "-1・0・1の戻り値を比較し、no score・player・CPUの三契約を区別します。",
      why: "5サブテストのうち3件は得点しない／境界の契約です。AAAは派手な成功経路だけでなく、変化してはいけない状態を守るためにも使います。",
      stories: ["場内では得点しない", "境界ちょうどではまだ得点しない", "境界を越えた側だけが得点する"],
      trap: "上へ出たらplayer、下へ出たらCPUの2件だけでは、境界ちょうどを&lt;=で誤判定しても見逃します。",
      practice: "min=-20、max=740を別の値へ変えた小さなTest関数を書き、関数が画面定数を内部に隠していないことを確認してください。",
    },
    en: {
      title: "AAA 05: make “nothing happens” a contract",
      intro: "No score while the ball is inside or exactly on an edge is as much a rule as scoring. AAA names that no-op with want=-1 and catches regressions that award a point too early.",
      arrange: "Place y inside, exactly on each boundary, and just beyond each boundary. Serve randomness is another responsibility and stays out.",
      act: "Call ExitScore once. Do not also move the ball or mutate the scoreboard; ask only who scores.",
      assert: "Compare -1, 0, or 1 to distinguish no score, player, and CPU contracts.",
      why: "Three of the five subtests protect no-score or exact-boundary behavior. AAA protects state that must not change, not only exciting success paths.",
      stories: ["inside the field gives no score", "the exact boundary still gives no score", "only crossing the boundary awards a point"],
      trap: "Testing only upper-player and lower-CPU exits misses an accidental &lt;= at the exact edge.",
      practice: "Write a small separate Test function with different min and max values to prove the rule does not hide screen constants.",
    },
  },
  "space-shooter": {
    ja: {
      title: "AAA 06：依存物をArrangeから追い出す",
      intro: "照準テストに敵、プレイヤー、弾生成、乱数まで並べる必要はありません。必要な差分dx・dyとspeedだけをArrangeできる関数境界を作ると、AAAが高速で決定的になります。",
      arrange: "3-4-5の方向、負方向、長さ0の方向を数値で置きます。期待値を人が計算できる値にするのが大切です。",
      act: "AimVelocityを一度だけ呼びます。spawnやslice appendは別のActなので混ぜません。",
      assert: "vxとvyを別々に許容誤差で比較し、長さだけでなく向きの符号まで守ります。",
      why: "正方向・負方向・ゼロ方向の3サブテストへ増やしました。通常例だけでは見えない符号とゼロ除算の契約を、同じAAA形で確認します。",
      stories: ["3-4-5方向を指定speedへ拡大する", "左上方向でも符号を保つ", "同じ位置なら停止ベクトルを返す"],
      trap: "敵をspawnして弾sliceの末尾を探すAssertは、準備も観測も重く、失敗原因を照準以外へ広げます。",
      practice: "dx=0、dy=5、speed=10の垂直方向を加え、vx=0とvy=10をそれぞれAssertしてください。",
    },
    en: {
      title: "AAA 06: drive dependencies out of Arrange",
      intro: "An aiming test needs no enemy, player, bullet spawn, or randomness. A boundary that accepts only dx, dy, and speed makes AAA fast and deterministic.",
      arrange: "Use a 3-4-5 direction, a negative direction, and a zero-length direction. Choose values whose expected result can be calculated by hand.",
      act: "Call AimVelocity once. Spawning and slice append are different Acts and stay out.",
      assert: "Compare vx and vy separately with tolerance, protecting direction signs as well as vector length.",
      why: "The suite now has positive, negative, and zero-direction subtests. The same AAA shape protects sign handling and zero division beyond the happy path.",
      stories: ["scale a 3-4-5 direction to the requested speed", "preserve signs toward the upper left", "return a stopped vector at the same position"],
      trap: "Spawning an enemy and searching the bullet slice makes setup and observation heavy while widening failure causes beyond aiming.",
      practice: "Add dx=0, dy=5, speed=10 and Assert vx=0 and vy=10 separately.",
    },
  },
  breakout: {
    ja: {
      title: "AAA 07：一つのActが返す複数の契約を読む",
      intro: "落球という一つのActionは、次の残機数とゲームオーバー判定を同時に決めます。Actは一度のまま、Assertを二つに分けると、同じ遷移の完全な契約を読みやすく守れます。",
      arrange: "開始残機を3・2・1と並べます。ボール、serve、ブロック配置は残機遷移に不要です。",
      act: "SpendLifeを一度だけ呼びます。戻り値を得るために同じ関数をAssertごとに呼び直しません。",
      assert: "gotLivesとgotGameOverを別々に比較します。一方が正しくても他方が壊れれば、失敗文がそのフィールドを示します。",
      why: "3→2、2→1、1→0の3サブテストで、通常継続から終了境界までを途切れなくしました。複数Assertは複数Actではなく、一つの結果の複数観測です。",
      stories: ["3機から1機失っても続行する", "最後の一つ手前では1機残る", "最後の1機を失うと0機で終了する"],
      trap: "gameOverだけをAssertすると、残機が-1になる実装でも通る可能性があります。逆に残機だけでは終了フラグの回帰を見逃します。",
      practice: "各Assertの片方を一時的に誤ったwantへ変え、失敗メッセージが残機とgameOverを別々に説明できることを確認してください。",
    },
    en: {
      title: "AAA 07: read several contracts from one Act",
      intro: "One life-spend action decides both the next count and game-over state. Keep one Act, then use separate Assertions to protect the complete transition readably.",
      arrange: "Arrange starting lives of 3, 2, and 1. Balls, serves, and bricks are irrelevant to the life transition.",
      act: "Call SpendLife once. Do not call it again for each Assertion.",
      assert: "Compare gotLives and gotGameOver separately. If one field breaks, its own failure identifies it.",
      why: "Three subtests cover 3→2, 2→1, and 1→0 continuously. Multiple Assertions observe one result; they do not create multiple Acts.",
      stories: ["losing one of three lives continues", "the penultimate loss leaves one life", "losing the final life reaches zero and ends"],
      trap: "Asserting only gameOver may allow a count of -1; asserting only lives misses a broken end flag.",
      practice: "Temporarily make each want wrong in turn and verify that the failure messages distinguish lives from gameOver.",
    },
  },
  snake: {
    ja: {
      title: "AAA 08：代表値ではなく変化点を並べる",
      intro: "数式のテストは適当な入力を大量に並べるより、答えが変わる直前と直後をAAAにします。score/3なら2と3、下限なら18と99を選ぶと、仕様の段差が見えます。",
      arrange: "score=0・2・3・18・99を置き、整数除算の最初の段差と速度下限を表します。tickカウンターや蛇の体は要りません。",
      act: "SnakeStepIntervalを一度呼びます。frame%intervalで移動するかどうかはUpdate側の別ルールです。",
      assert: "返ったinterval一つをwantIntervalと比較します。失敗時に曲線のどの段差が変わったかをサブテスト名で追えます。",
      why: "2点では10のまま、3点で9になる行を追加しました。通常値・変化点の両側・下限・下限後を複数AAAで覆うと、式を眺めるだけより意図が強く残ります。",
      stories: ["0点では10 tick", "2点ではまだ10 tick、3点で9 tick", "18点以降は4 tickより短くならない"],
      trap: "0点と99点だけでは中間の整数除算を壊しても通ります。1〜100をループして複雑な式で正解を再実装するAssertも、同じバグを複製します。",
      practice: "score=5と6を隣り合わせで追加し、次の段差も同じAAAの読み方になることを確認してください。",
    },
    en: {
      title: "AAA 08: choose change points, not random examples",
      intro: "For formulas, use AAA immediately before and after answers change. With score/3, choose 2 and 3; around the floor, choose 18 and a high score.",
      arrange: "Use scores 0, 2, 3, 18, and 99 to expose the first integer-division step and the speed floor. No tick counter or snake body is needed.",
      act: "Call SnakeStepInterval once. Whether frame%interval moves is a separate Update rule.",
      assert: "Compare the one interval with wantInterval. The subtest name locates the changed step in the curve.",
      why: "A new pair proves that 2 stays at 10 while 3 becomes 9. Multiple AAA cases preserve normal, both sides of a change, the floor, and values beyond it.",
      stories: ["zero score uses ten ticks", "two stays at ten while three becomes nine", "eighteen and above never go below four"],
      trap: "Only score 0 and 99 miss broken middle steps. Looping 1–100 with a complicated expected-value formula merely reimplements the same bug in Assert.",
      practice: "Add scores 5 and 6 next to each other and confirm that the next step reads with the same AAA shape.",
    },
  },
  sokoban: {
    ja: {
      title: "AAA 09：遷移前・直前・完了を分ける",
      intro: "状態遷移は「途中」と「完了」の2件だけでなく、完了直前も独立したAAAにすると早すぎる遷移を防げます。progress=0.85から0.14進めても0.99で未完了、を明文化します。",
      arrange: "progressを途中0.28、直前0.85、越える0.98に置き、step=0.14は固定します。箱やプレイヤー座標は補間規則に含めません。",
      act: "AdvanceTweenを一度呼びます。Drawによる見た目の補間や次の入力処理まで進めません。",
      assert: "next progressとcompleteを別々に比較し、クランプ値と遷移フラグが同時に正しいことを守ります。",
      why: "3サブテストが途中・1未満ぎりぎり・1越えを作りました。境界直前を足したことで、&gt;=を誤って早める実装を検出できます。",
      stories: ["途中では値だけ進み未完了", "0.99ではまだ未完了", "1を越えると1へクランプして完了"],
      trap: "アニメーションを描いて最終ピクセルだけ比較しても、progressの完了タイミングや入力解禁の契約は分かりません。",
      practice: "progress=0.86を加えてちょうど1へ到達する契約を決めてください。期待を決めてから実行し、現在実装の&gt;=と一致するか確認します。",
    },
    en: {
      title: "AAA 09: separate before, just-before, and complete",
      intro: "A state transition needs more than one in-progress and one complete case. An independent just-before AAA case—0.85 plus 0.14 remains 0.99 and incomplete—prevents early transitions.",
      arrange: "Use progress 0.28, 0.85, and 0.98 with a fixed step of 0.14. Box and player positions do not belong to this interpolation rule.",
      act: "Call AdvanceTween once. Do not also render interpolation or accept the next input.",
      assert: "Compare next progress and complete separately, protecting both the clamped value and transition flag.",
      why: "Three subtests now cover in-progress, just below one, and beyond one. The new near-boundary case catches an implementation that completes too early.",
      stories: ["ordinary progress advances and stays incomplete", "0.99 is still incomplete", "crossing one clamps to one and completes"],
      trap: "Comparing only final animation pixels cannot define when progress completes or input becomes available.",
      practice: "Add progress=0.86 to define exactly reaching one. Choose the expectation before running and compare it with the current >= contract.",
    },
  },
  platformer: {
    ja: {
      title: "AAA 10：異なるActionはTest関数を分ける",
      intro: "横移動と縦移動は同じプレイヤーに起きますが、壊れる理由は別です。AAAのActを一つに保つためTestHorizontalVelocityとTestVerticalVelocityへ分け、各表の中で入力組合せを増やします。",
      arrange: "横はvx・left・right、縦はvy・jump・groundedだけを置きます。地形衝突、コイン、カメラは別の契約です。",
      act: "一つのサブテストはHorizontalVelocityかVerticalVelocityの一方だけを呼びます。横計算→縦計算→衝突までを一つのActと呼びません。",
      assert: "横は速度一つ、縦は速度とleftGroundを比較します。観測対象は公開された結果であり、if文を何行通ったかではありません。",
      why: "横5件・縦4件の9サブテストを、二つの異なるActへ整理しています。複雑なゲームロジックほど、AAAはTest関数を分ける境界線になります。",
      stories: ["左右入力と摩擦は横速度だけを決める", "接地ジャンプは上向き速度と離地を返す", "空中ジャンプは無視しつつ重力は続く"],
      trap: "一つのテストで横速度、縦速度、X衝突、Y衝突、コイン取得まで進めると、失敗したAssertから原因のActへ戻れません。",
      practice: "新しいダッシュ規則を想像し、既存HorizontalVelocityの表へ入る同形ケースか、別のActなのでTestDashVelocityに分けるべきかを理由つきで決めてください。",
    },
    en: {
      title: "AAA 10: split different Actions into different Test functions",
      intro: "Horizontal and vertical motion affect one player but fail for different reasons. Keep one Act by separating TestHorizontalVelocity from TestVerticalVelocity, then vary inputs inside each table.",
      arrange: "Horizontal cases contain vx, left, and right; vertical cases contain vy, jump, and grounded. Terrain, coins, and camera are other contracts.",
      act: "A subtest calls either HorizontalVelocity or VerticalVelocity. Do not rename horizontal calculation, vertical calculation, and collision together as one Act.",
      assert: "Horizontal cases compare one speed; vertical cases compare speed and leftGround. Observe returned behavior, not which internal if statements ran.",
      why: "Nine subtests—five horizontal and four vertical—are organized around two distinct Acts. In larger game logic, AAA becomes a boundary for splitting Test functions.",
      stories: ["left/right input and friction decide only horizontal speed", "a grounded jump returns upward speed and leaving-ground state", "an air jump is ignored while gravity continues"],
      trap: "A test that advances horizontal speed, vertical speed, X collision, Y collision, and a coin cannot trace a failed Assertion back to one Act.",
      practice: "Imagine a dash rule. Decide, with a reason, whether it is another same-shape HorizontalVelocity row or a distinct Act deserving TestDashVelocity.",
    },
  },
  dungeon: {
    ja: {
      title: "AAA 11：現在状態もArrangeして状態機械を守る",
      intro: "AIの答えは距離だけでなく現在のmodeにも依存します。同じ距離200でも徘徊中なら徘徊、追跡中なら追跡という二つのAAAを書くことで、ヒステリシスの『状態を保つ』を守れます。",
      arrange: "current modeとdistanceを一組にします。さらに165・230の境界ちょうどと、その外側164.9・230.1を並べます。",
      act: "EnemyModeを一度呼びます。移動速度を求めるAimVelocityは別のTest関数・別のActとして残します。",
      assert: "next mode一つを比較します。敵の座標やアニメーションではなく、状態機械が約束した遷移だけを観測します。",
      why: "EnemyModeは6サブテストになり、開始・維持・終了と両境界ちょうどを覆います。AimVelocityの3件とは分離され、AI判断と移動計算の失敗が混ざりません。",
      stories: ["165未満なら追跡を始める", "中間距離では現在modeを保つ", "230を越えたら徘徊へ戻る"],
      trap: "近い敵は追跡、遠い敵は徘徊の2件だけでは、中間帯でmodeを毎tick反転させるバグを見逃します。",
      practice: "current=0,distance=230とcurrent=1,distance=165を追加し、境界ちょうどではそれぞれの現在状態を保つことを確かめてください。",
    },
    en: {
      title: "AAA 11: Arrange current state to protect a state machine",
      intro: "AI output depends on current mode as well as distance. At the same distance of 200, separate AAA cases keep wander while wandering and chase while chasing, protecting hysteresis.",
      arrange: "Pair current mode with distance. Put exact boundaries 165 and 230 next to outside values 164.9 and 230.1.",
      act: "Call EnemyMode once. AimVelocity stays in another Test function as another Act.",
      assert: "Compare one next mode. Observe the promised transition rather than enemy coordinates or animation.",
      why: "EnemyMode now has six subtests covering enter, retain, exit, and both exact boundaries. Three AimVelocity cases stay separate, so decision and motion failures do not mix.",
      stories: ["below 165 begins chasing", "the middle band preserves current mode", "above 230 returns to wander"],
      trap: "Only near-chases and far-wanders misses a bug that flips mode every tick inside the middle band.",
      practice: "Add current=0,distance=230 and current=1,distance=165 to prove exact boundaries preserve their respective current modes.",
    },
  },
  "bullet-hell": {
    ja: {
      title: "AAA 12：任意のロジックを契約の網にする",
      intro: "最後は巨大な弾幕Updateを、円同士の接触と画面外判定という二つのActへ分けます。各Actへ通常・境界・境界外・方向違いのAAAを複数置けば、複雑さを弾数ではなく小さな契約の網で制御できます。",
      arrange: "CircleHitには円2個、Outsideには点1個・矩形・marginだけを渡します。弾slice、HP、削除順はそれぞれ別のテスト対象です。",
      act: "一つのサブテストでCircleHitかOutsideのどちらかを一度呼びます。大量の弾をループすることをActにしません。",
      assert: "hit／outsideのbool一つを比較します。Assert内で衝突数を探索したり、期待結果を別の幾何計算で再実装したりしません。",
      why: "接触3件と画面外6件が、重なり・接触ちょうど・分離、内側・margin上・左右上下の外側を守ります。任意のロジックでも、分岐ごとに小さなAAAを増やす習慣が完成します。",
      stories: ["重なる円だけがhitになる", "接するだけならhitにならない", "矩形の左・右・上・下の各margin越えはoutsideになる"],
      trap: "弾を100発作ってUpdate後の残数だけをAssertすると、移動・画面外・衝突・逆順削除・HPのどれが原因か分かりません。",
      practice: "自分のゲームから一つのifを選び、通常・境界ちょうど・境界の直後の3つを先に日本語で書いてください。その後に各文を1行ずつテーブルへ移し、1 Actと直接のAssertで完成させます。",
    },
    en: {
      title: "AAA 12: cover arbitrary logic with a net of contracts",
      intro: "The final chapter splits a giant bullet-hell Update into two Acts: circle contact and outside-screen classification. Multiple normal, boundary, beyond-boundary, and directional AAA cases control complexity through small contracts rather than bullet count.",
      arrange: "CircleHit receives two circles; Outside receives one point, a rectangle, and margin. Bullet slices, HP, and removal order are separate subjects.",
      act: "Each subtest calls either CircleHit or Outside once. Looping over hundreds of bullets is not the Act under test.",
      assert: "Compare one hit or outside bool. Do not search for collision counts or reimplement geometry inside Assert.",
      why: "Three contact cases and six outside cases protect overlap, touching, separation, inside, exact margin, and every outer direction. The course ends with a habit of adding small AAA cases for each branch.",
      stories: ["only overlapping circles hit", "touching edges do not hit", "crossing left, right, top, or bottom margin is outside"],
      trap: "Creating one hundred bullets and asserting only the remaining count mixes movement, bounds, collisions, reverse removal, and HP into one failure.",
      practice: "Choose one if statement from your game. Write normal, exact-boundary, and just-beyond stories in prose first; then turn each sentence into one table row with one Act and a direct Assert.",
    },
  },
};

if (Object.keys(aaaCourse).length !== lessons.length || lessons.some((lesson) => !aaaCourse[lesson.slug])) {
  throw new Error("every testing lesson must include one progressive AAA review");
}
for (const lesson of lessons) {
  for (const lang of ["ja", "en"]) {
    const review = aaaCourse[lesson.slug][lang];
    if (!review || review.stories.length < 3 || !review.arrange || !review.act || !review.assert || !review.why || !review.trap || !review.practice) {
      throw new Error(`${lesson.slug}/${lang}: incomplete AAA review`);
    }
  }
}

const goldenCases = {
  "tap-target": { testName: "InitialTarget", fixture: "g := newGame()", file: "tap-target-initial.png", ja: "固定seedで作った最初の的・HUD・開始表示", en: "the initial seeded target, HUD, and start prompt" },
  "timing-meter": { testName: "StoppedPerfect", fixture: "g := &game{markerX: centerX, speed: 3.2, score: 100, stopped: true}", file: "timing-meter-perfect.png", ja: "中央で止めたPERFECT表示と得点HUD", en: "the PERFECT result and score HUD after stopping at center" },
  "catch-stars": { testName: "InitialBasket", fixture: "g := newGame()", file: "catch-stars-initial.png", ja: "固定seedの最初のカゴ・星・残機HUD", en: "the initial seeded basket, stars, and lives HUD" },
  flappy: { testName: "ReadyState", fixture: "g := newGame()", file: "flappy-ready.png", ja: "開始前の鳥・パイプ・操作案内", en: "the bird, pipes, and input prompt before play starts" },
  pong: { testName: "InitialServe", fixture: "g := newGame()", file: "pong-initial-serve.png", ja: "初期サーブ・パドル・0対0のHUD", en: "the initial serve, paddles, and 0–0 HUD" },
  "space-shooter": { testName: "FirstWave", fixture: "g := newGame()", file: "space-shooter-wave-one.png", ja: "固定seedの第1ウェーブ・自機・HUD", en: "the seeded first wave, player ship, and HUD" },
  breakout: { testName: "InitialBricks", fixture: "g := newGame()", file: "breakout-initial.png", ja: "初期ブロック配置・ボール・残機HUD", en: "the initial brick layout, ball, and lives HUD" },
  snake: { testName: "InitialBoard", fixture: "g := newGame()", file: "snake-initial.png", ja: "固定seedの初期ヘビ・餌・得点HUD", en: "the seeded initial snake, food, and score HUD" },
  sokoban: { testName: "InitialBoard", fixture: "g := newGame()", file: "sokoban-initial.png", ja: "初期の壁・箱・ゴール・プレイヤー", en: "the initial walls, boxes, goals, and player" },
  platformer: { testName: "InitialStage", fixture: "g := newGame()", file: "platformer-initial.png", ja: "開始地点の足場・コイン・プレイヤー", en: "the opening platforms, coins, and player" },
  dungeon: { testName: "InitialRoom", fixture: "g := newGame()", file: "dungeon-initial.png", ja: "最初の部屋・敵・ライフHUD", en: "the initial room, enemies, and life HUD" },
  "bullet-hell": { testName: "InitialBoss", fixture: "g := newGame()", file: "bullet-hell-initial.png", ja: "開始時のボス・自機・HP表示", en: "the opening boss, player, and HP display" },
};

if (Object.keys(goldenCases).length !== lessons.length || lessons.some((lesson) => !goldenCases[lesson.slug])) {
  throw new Error("every testing lesson must include one Draw golden example");
}

const aaaPairExample = `func TestSpendLifeContinuesWhenLivesRemain(t *testing.T) {
	startingLives := 3

	gotLives, gotGameOver := SpendLife(startingLives)

	if gotLives != 2 {
		t.Errorf("lives = %d, want 2", gotLives)
	}
	if gotGameOver {
		t.Error("gameOver = true, want false")
	}
}

func TestSpendLifeEndsAfterLastLife(t *testing.T) {
	startingLives := 1

	gotLives, gotGameOver := SpendLife(startingLives)

	if gotLives != 0 {
		t.Errorf("lives = %d, want 0", gotLives)
	}
	if !gotGameOver {
		t.Error("gameOver = false, want true")
	}
}`;

const readableTestExample = `func TestSpendLife(t *testing.T) {
	testCases := []struct {
		name         string
		lives        int
		wantLives    int
		wantGameOver bool
	}{
		{name: "spending one of three lives continues game", lives: 3, wantLives: 2},
		{name: "spending last life ends game", lives: 1, wantLives: 0, wantGameOver: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotLives, gotGameOver := SpendLife(tc.lives)

			if gotLives != tc.wantLives {
				t.Errorf("lives = %d, want %d", gotLives, tc.wantLives)
			}
			if gotGameOver != tc.wantGameOver {
				t.Errorf("gameOver = %v, want %v", gotGameOver, tc.wantGameOver)
			}
		})
	}
}`;

const goldenHarnessExample = `package main

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type testLoop struct {
	m    *testing.M
	code int
}

func (l *testLoop) Update() error {
	l.code = l.m.Run()
	return ebiten.Termination
}

func (*testLoop) Draw(*ebiten.Image) {}
func (*testLoop) Layout(_, _ int) (int, int) { return 1, 1 }

func TestMain(m *testing.M) {
	loop := &testLoop{m: m}
	if err := ebiten.RunGame(loop); err != nil {
		panic(err)
	}
	os.Exit(loop.code)
}

func assertGolden(t *testing.T, got *ebiten.Image, filename string) {
	t.Helper()

	var actual bytes.Buffer
	if err := png.Encode(&actual, got); err != nil {
		t.Fatal(err)
	}

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filename, actual.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}

	want, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual.Bytes(), want) {
		actualFile := filename + ".actual.png"
		if err := os.WriteFile(actualFile, actual.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Fatalf("Draw output differs from %s; actual image: %s", filename, actualFile)
	}
}`;

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
  const constants = body.includes("floatTolerance") ? "\n\nconst floatTolerance = 1e-9" : "";
  return `package lessonlogic\n\n${imports}${constants}\n\n${body}\n`;
}

const esc = (s) => String(s).replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");

function shell({ lang, depth, title, desc, body }) {
  const ja = lang === "ja";
  const prefix = "../".repeat(depth);
  const other = ja ? "en" : "ja";
  const route = depth === 3 ? "guides/testing/" : `guides/testing/${title.slug}/`;
  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><meta name="description" content="${esc(desc)}"><title>${esc(title.text)} | Ebi Showcase</title><link rel="icon" type="image/png" href="${prefix}assets/favicon.png"><link rel="stylesheet" href="${prefix}style.css"></head><body class="testing-guide"><header class="nav"><a class="brand" href="${prefix}${lang}/"><span>EBI</span> SHOWCASE</a><nav><a href="${depth === 3 ? "./" : "../"}">TESTING</a><a class="lang" href="${prefix}${other}/${route}" lang="${other}">${ja ? "EN" : "日本語"}</a></nav></header><main>${body}</main><footer><span>EBI SHOWCASE</span><span>GO + EBITENGINE</span><span>APACHE-2.0</span></footer><script src="${prefix}learn.js"></script></body></html>`;
}

function hub(lang) {
  const ja = lang === "ja";
  const c = ja ? {
    title: "LEVEL 01〜12のUpdateをユニットテストする", desc: "Drawを原則テストせず、EbitengineのUpdateから純粋なゲームルールを分離してLEVEL 01〜12をGoで網羅的にテストする教材。",
    h1: "Drawを試す前に、<br><em>Updateのルールを守る。</em>",
    lead: "Drawは現在のgameを画面へ投影するだけなので、原則としてユニットテストの対象にしません。入力・時間・乱数・衝突・得点を持つUpdateを、薄い接続役と純粋なルールへ分けます。下の12章では実際のUpdateを読み、AAAを小さな判定から物理、状態遷移、AI、弾幕へ段階的に使い、各ルールを複数の独立したテストで守ります。",
    input: "Updateが入力を読む", pure: "純粋関数へ普通の値を渡す", state: "Updateがgameへ結果を書く", noDraw: "Drawはgameを読むだけ",
    why: "なぜDraw()は原則テストしない？", whyBody: "同じgame・素材・画面サイズなら同じ絵を出す、という参照透過なDrawを守れば、ルールの正しさはgameの状態で確認できます。ピクセル比較はフォント・GPU・OS差で壊れやすく、当たり判定や得点の答えにもなりません。独自シェーダーや画像生成器を守りたい場合は別枠のゴールデン画像テストを使いますが、それはUpdateのユニットテストとは分けます。",
    update: "Update()をどう構造化する？", updateBody: "①Ebitengineから入力を読む、②普通の数値・小さな構造体へ変換、③純粋関数で次の答えを計算、④gameへ結果を保存、の順にします。純粋関数はEbitengineをimportせず、時刻や乱数も引数で受け取ります。自分のGoプロジェクトにinternal/lessonlogicを作ればよく、このサイトのcloneは不要です。",
    foundationTitle: "AAAは、失敗を一つの理由へ絞るための設計。", foundationBody: "Arrangeはテストしたい世界を一文で作り、Actはその世界へ一つの出来事だけを起こし、Assertは外から観測できる契約だけを判定します。// Arrange・// Act・// Assertというラベルコメントは書かず、意味のある変数名と空行で読み順を見せます。ただし空行があるだけではAAAではありません。前提が巨大、Actが複数、Assertが実装の内側を探るなら分解が必要です。下の2本は同じ関数を調べますが、継続と終了を別の物語にしたため、片方が失敗しても原因が混ざりません。",
    arrangeTitle: "01 / ARRANGE — 一つの世界", arrangeBody: "必要な値だけを置きます。入力機器、画面、乱数、巨大なgameが質問に不要なら作りません。前提を一文で説明できなければ、対象がまだ大きすぎます。",
    actTitle: "02 / ACT — 一つの出来事", actBody: "対象の関数や状態遷移を原則一度だけ実行します。二つの異なる操作が必要なら、別のテストか、より小さい関数境界を検討します。",
    assertTitle: "03 / ASSERT — 観測できる契約", assertBody: "gotとwantを直接比べます。複数の戻り値は別々に報告して構いませんが、探索ループや期待値を再計算する複雑な分岐は置きません。",
    foundationSmell: "AAAの匂いチェック：テスト名だけで前提と結果が言えるか。Actは一つか。失敗文だけで壊れた契約が分かるか。通常だけで終わらず、境界・何も起きない経路・優先順位のうち、そのルールに必要な別のAAAを書いたか。",
    qualityTitle: "複数のAAAを書いてから、同じ形だけを表へまとめる。",
    aaaTitle: "AAAレビューを習慣に", aaa: "テストを書き終えたら、前提を一文にできるか、Actが一つか、Assertが結果だけを見ているかを声に出して確認します。さらに通常・境界・直後またはno-opの複数テストがあるかを見ます。",
    dampTitle: "DAMPをDRYより優先", damp: "DAMPはDescriptive And Meaningful Phrasesです。テストでは意味が見える名前や値を少し繰り返して構いません。Assertの中に探索用のfor、switch、複雑な分岐を置かず、一つのgotを一つのwantと比べます。",
    namesTitle: "Goらしい名前", names: "whenFooBar_shouldQuaxはGoの決まりではありません。トップレベルはTestSpendLifeのように対象を名付け、条件と結果はt.Run(\"spending last life ends game\")のようなサブテスト名で表します。これならgo test -runで一件だけ選べます。",
    tableTitle: "同じ形ならテーブル駆動", tables: "境界値のようにArrange・Act・Assertの形が同じで値だけ違うときは、各行が完全な一件になるテーブル駆動テストにします。準備や操作そのものが違うシナリオまで無理に一つの表へ押し込みません。",
    sources: "Go公式の命名説明と、t.Runを使うテーブル駆動テスト",
    goldenTitle: "Drawのゴールデンは、Layoutの表示契約ごとに分ける。",
    goldenBody: "Drawの呼び出し回数、内部座標、個々の描画命令はテストしません。同じgame状態・asset・fontを固定したまま、mobile portrait（390×844）、console 16:9（1920×1080）、PC window（1280×800）の外側サイズをそれぞれLayoutへ渡します。Layoutが返した論理サイズで一枚描き、パターン別の承認済みPNGと比較します。これにより、モバイルの縦積みとタッチ領域、コンソール／TVの安全領域と離れて読めるHUD、PCウインドウの横幅の使い方を、同じゲーム状態で調整できます。入力表示も変えるなら、端末を推測せずcontrol schemeをfixtureへ明示します。同じ論理サイズ・同じ配置を返す固定Layoutなら、同一画像を3枚複製せずdefaultの1枚だけで十分です。ブラウザ側の拡大縮小やCSSの折り返しはDrawではなく実ブラウザのスクリーンショットテストで守ります。差分が出たら.actual.pngを人が見てからUPDATE_GOLDEN=1で更新します。Ebitengineの画素読み出しはメインループ内で行う必要があるため、下のTestMain harness内でテストを実行します。",
    reviewTitle: "1本で満足しない", reviewBody: "テスト件数を水増しするのではなく、通常、境界ちょうど、境界の直後、何も起きない経路、条件が重なる優先順位から必要な物語を選びます。一つの物語につき一つのAAAを完結させます。",
    cards: "12本すべてを、AAAで一つずつ分解する", cardsLead: "各ページでは現行テストを一つずつレビューし、章固有のArrange・Act・Assert、複数シナリオ、避けるべき巨大テスト、次に自分で足す1件まで示します。最後には任意のゲームロジックを小さな契約の網へ分けられます。", read: "全文を読む →",
  } : {
    title: "Unit-test Update across LEVEL 01–12", desc: "A complete guide to leaving Draw untested by default and extracting pure rules from every LEVEL 01–12 Update for Go unit tests.",
    h1: "Protect Update rules<br><em>before testing pictures.</em>",
    lead: "Draw projects the current game and is normally not a unit-test target. Split input, time, randomness, collision, and scoring in Update into a thin adapter plus pure rules. Across twelve chapters, AAA grows from a tiny predicate through physics, transitions, AI, and bullet hell, with every rule protected by several independent tests.",
    input: "Update reads input", pure: "pass plain values to pure rules", state: "Update writes results to game", noDraw: "Draw only reads game",
    why: "Why not unit-test Draw() by default?", whyBody: "When Draw is a referentially transparent projection, the same game, assets, and viewport produce the same picture; rule correctness is visible in game state. Pixel comparisons are brittle across fonts, GPUs, and operating systems and do not prove collision or scoring. Protect a custom shader or image generator with a separate golden-image suite when needed—not with Update unit tests.",
    update: "How should Update() be structured?", updateBody: "① Read Ebitengine input. ② Convert it to plain values or small structs. ③ Calculate the next answer in pure functions. ④ Store results in game. Pure functions do not import Ebitengine and receive time/random values as arguments. Create internal/lessonlogic in your own Go project; cloning this site is not required.",
    foundationTitle: "AAA is design for reducing a failure to one reason.", foundationBody: "Arrange creates the tested world in one sentence, Act causes one event in that world, and Assert judges only an externally observable contract. Do not label the code with // Arrange, // Act, or // Assert comments; meaningful names and blank lines reveal the reading order. Blank lines alone do not make AAA. Huge setup, several Acts, or Assertions that search implementation details still need decomposition. The two tests below call the same rule, but continuation and ending are separate stories, so their failures cannot blur together.",
    arrangeTitle: "01 / ARRANGE — one world", arrangeBody: "Place only required values. Do not create input devices, a screen, randomness, or a giant game when the question does not need them. If the premise does not fit one sentence, the subject is probably too large.",
    actTitle: "02 / ACT — one event", actBody: "Execute the subject function or transition once by default. If the story needs two different operations, consider separate tests or a smaller function boundary.",
    assertTitle: "03 / ASSERT — observable contract", assertBody: "Compare got directly with want. Report several returned fields separately when needed, but do not add search loops or complicated branches that recalculate the answer.",
    foundationSmell: "AAA smell check: does the name state premise and result? Is there one Act? Does the failure identify the broken contract? Beyond the normal path, did you add the boundary, no-op, or precedence AAA stories that this rule actually needs?",
    qualityTitle: "Write several AAA stories, then table only the identical shape.",
    aaaTitle: "Make AAA review a habit", aaa: "After writing a test, say whether its premise fits one sentence, its Act is one event, and its Assert observes only results. Then look for multiple required stories: normal, exact boundary, just beyond, or no-op.",
    dampTitle: "Prefer DAMP over test DRYness", damp: "DAMP means Descriptive And Meaningful Phrases. Repeating a meaningful name or value in test code is fine. Do not put search loops, switches, or complicated branches in Assert; compare one got value with one want value.",
    namesTitle: "Idiomatic Go names", names: "whenFooBar_shouldQuax is not a Go rule. Name the top-level subject, such as TestSpendLife, then put premise and outcome in a subtest such as t.Run(\"spending last life ends game\"). go test -run can select that one case.",
    tableTitle: "Use tables when the shape repeats", tables: "Use a table-driven test when boundary cases share the same Arrange, Act, and Assert and only values change. Each row must be one complete case. Do not force scenarios with different setup or actions into one clever table.",
    sources: "Official Go naming guidance and table-driven subtests with t.Run",
    goldenTitle: "Keep one Draw golden for each Layout contract.",
    goldenBody: "Do not test Draw call counts, internal coordinates, or individual drawing commands. Keep the game state, assets, and fonts fixed, then pass a mobile portrait (390×844), console 16:9 (1920×1080), and PC window (1280×800) outer size through Layout. Render at the logical size returned by Layout and compare each result with its own approved PNG. The same scene can then tune mobile stacking and touch space, console/TV safe areas and distance-readable HUDs, and use of width in a PC window. If input hints differ, put an explicit control scheme in the fixture instead of guessing the device. A fixed Layout that returns the same logical size and arrangement needs one default golden—not three identical copies. Browser scaling and CSS wrapping belong in real-browser screenshot tests rather than Draw tests. Inspect .actual.png before running UPDATE_GOLDEN=1. Ebitengine pixel reads must happen inside its main loop, so the TestMain harness below runs the tests there.",
    reviewTitle: "Never stop at one happy test", reviewBody: "Do not inflate a count. Choose the stories the rule needs from normal, exact boundary, just beyond, no-op, and precedence. Complete one AAA cycle for each story.",
    cards: "Decompose all twelve real games with AAA", cardsLead: "Every page reviews the current tests individually and names its Arrange, Act, Assert, scenario set, giant-test smell, and one next case to add. By the end, arbitrary game logic becomes a net of small contracts.", read: "READ ALL →",
  };
  const cards = lessons.map((lesson) => { const t = lesson[lang]; return `<a class="test-course-card" href="${lesson.slug}/"><span>${lesson.level}</span><h3>${t.title}</h3><p>${t.lead}</p><strong>${c.read}</strong></a>`; }).join("");
  const foundationCards = [["arrange", c.arrangeTitle, c.arrangeBody], ["act", c.actTitle, c.actBody], ["assert", c.assertTitle, c.assertBody]].map(([stage, title, body]) => `<article data-aaa-stage="${stage}"><h3>${title}</h3><p>${body}</p></article>`).join("");
  const qualityCards = [[c.aaaTitle, c.aaa], [c.dampTitle, c.damp], [c.namesTitle, c.names], [c.tableTitle, c.tables], [c.reviewTitle, c.reviewBody]].map(([title, body]) => `<article><h3>${title}</h3><p>${body}</p></article>`).join("");
  const body = `<section class="test-hero"><p class="eyebrow">SPECIAL GUIDE / UNIT TESTING</p><h1>${c.h1}</h1><p>${c.lead}</p></section><section class="test-boundary" aria-label="Update and Draw testing boundary"><div><small>01 / INPUT ADAPTER</small><b>${c.input}</b></div><i>→</i><div class="is-pure"><small>02 / TEST HERE</small><b>${c.pure}</b><em>go test / no window</em></div><i>→</i><div><small>03 / STATE</small><b>${c.state}</b><em>${c.noDraw}</em></div></section><section class="test-explain test-principles"><div><p class="eyebrow">THE DEFAULT</p><h2>${c.why}</h2><p>${c.whyBody}</p></div><div><p class="eyebrow">THE STRUCTURE</p><h2>${c.update}</h2><p>${c.updateBody}</p></div></section><section class="test-quality test-aaa-foundation" data-aaa-foundation><div class="test-quality-copy"><p class="eyebrow">AAA / ONE FAILURE, ONE REASON</p><h2>${c.foundationTitle}</h2><p>${c.foundationBody}</p><div class="test-quality-grid test-aaa-stage-grid">${foundationCards}</div><p class="test-aaa-smell"><strong>AAA REVIEW QUESTIONS</strong>${c.foundationSmell}</p></div><pre><code>${esc(aaaPairExample)}</code></pre></section><section class="test-quality"><div class="test-quality-copy"><p class="eyebrow">MULTIPLE AAA + DAMP + GO STYLE</p><h2>${c.qualityTitle}</h2><div class="test-quality-grid">${qualityCards}</div><p class="test-sources"><a href="https://go.dev/doc/tutorial/add-a-test">${c.sources}</a> · <a href="https://go.dev/blog/subtests">t.Run / subtests</a></p></div><pre><code>${esc(readableTestExample)}</code></pre></section><section class="test-golden"><div><p class="eyebrow">DRAW / VISUAL REGRESSION</p><h2>${c.goldenTitle}</h2><p>${c.goldenBody}</p></div><pre><code>${esc(goldenHarnessExample)}</code></pre></section><section class="test-terminal"><code>$ go test ./internal/lessonlogic</code><strong>✓ rules, not pixels</strong></section><section class="test-course"><p class="eyebrow">LEVEL 01–12 / COMPLETE MAP</p><h2>${c.cards}</h2><p>${c.cardsLead}</p><div class="test-course-grid">${cards}</div></section><nav class="test-guide-links"><a href="tap-target/">${ja ? "LEVEL 01から始める →" : "START LEVEL 01 →"}</a><a href="../../">${ja ? "← ホームへ" : "← HOME"}</a></nav>`;
  return shell({ lang, depth: 3, title: { text: c.title }, desc: c.desc, body });
}

function lessonPage(lang, index) {
  const lesson = lessons[index];
  const t = lesson[lang];
  const ja = lang === "ja";
  const aaa = aaaCourse[lesson.slug][lang];
  const golden = goldenCases[lesson.slug];
  const goldenStem = golden.file.replace(/\.png$/, "");
  const update = updateFunction(lesson.source);
  const pureFunctions = lesson.pure.map((name) => namedFunction(logicFiles, name));
  const tests = lesson.tests.map((name) => namedFunction(testFiles, name));
  const pure = goFile(pureFunctions, supportCode(lesson.support));
  const test = testFile(tests);
  const scenarioCount = (test.match(/\{name:/g) || []).length;
  if (scenarioCount < 3) {
    throw new Error(`${lesson.slug}: every rule needs several independent AAA scenarios`);
  }
  const goldenTest = `func TestDraw${golden.testName}MatchesLayoutGoldens(t *testing.T) {
	layoutCases := []struct {
		name                       string
		outsideWidth, outsideHeight int
		golden                     string
	}{
		{name: "mobile portrait", outsideWidth: 390, outsideHeight: 844, golden: "testdata/${goldenStem}-mobile.png"},
		{name: "console 16:9", outsideWidth: 1920, outsideHeight: 1080, golden: "testdata/${goldenStem}-console.png"},
		{name: "PC window", outsideWidth: 1280, outsideHeight: 800, golden: "testdata/${goldenStem}-pc.png"},
	}

	for _, tc := range layoutCases {
		t.Run(tc.name, func(t *testing.T) {
			${golden.fixture}
			logicalWidth, logicalHeight := g.Layout(tc.outsideWidth, tc.outsideHeight)
			screen := ebiten.NewImage(logicalWidth, logicalHeight)

			g.Draw(screen)

			assertGolden(t, screen, tc.golden)
		})
	}
}`;
  const rows = t.cases.map((r) => `<tr><th>${esc(r[0])}</th><td><code>${esc(r[1])}</code></td><td><code>${esc(r[2])}</code></td></tr>`).join("");
  const phases = lesson.phases[lang].map((phase, i) => `<li><span>${String(i + 1).padStart(2, "0")}</span><p>${phase}</p></li>`).join("");
  const aaaStages = [["arrange", "01 / ARRANGE", aaa.arrange], ["act", "02 / ACT", aaa.act], ["assert", "03 / ASSERT", aaa.assert]].map(([stage, title, body]) => `<article data-aaa-stage="${stage}"><h3>${title}</h3><p>${body}</p></article>`).join("");
  const aaaStories = aaa.stories.map((story, storyIndex) => `<li><span>${String(storyIndex + 1).padStart(2, "0")}</span><p>${story}</p></li>`).join("");
  const prev = index === 0 ? "../" : `../${lessons[index - 1].slug}/`;
  const next = index === lessons.length - 1 ? "../" : `../${lessons[index + 1].slug}/`;
  const c = ja ? {
    complete: "現行Updateの全手順", completeBody: "省略記号はありません。入力から早期return、ループ、勝敗まで、現在ゲームで動く順番です。この中から抽出した純粋関数を直接テストします。",
    update: "REAL GO / Update()全文", pure: "EXTRACTED / 完全な純粋ロジック", test: "TEST / 完全な_test.go", table: "境界を先に決める", manual: "テストのあとも人が確認すること", run: "自分のGoプロジェクトに上のファイルを作って実行", back: "← 前へ", forward: "次へ →",
    aaaEyebrow: `AAA PRACTICE ${lesson.level} / 12`, aaaCount: `${scenarioCount}本のサブテストを一つずつレビュー済みです。各行は「データの数」ではなく、独立した一つのAAAの物語です。`, repetitions: "最低3つの物語で、隣り合う契約を守る", trapTitle: "AAA SMELL / 分解する合図", practiceTitle: "ONE MORE TEST / 次は自分で1本",
    auditTitle: "この章の全テストを、AAAと読みやすさの両方で点検済み。", names: "命名", namesBody: `トップレベルは${lesson.tests.map((name) => `<code>${name}</code>`).join("・")}。各<code>t.Run</code>名が「どういう前提なら何を返すか」を文章で示します。<code>whenFoo_shouldBar</code>形式は使いません。`,
    aaa: "AAAレビュー結果", aaaBody: aaa.why,
    damp: "DAMP", dampBody: "意味のあるケース名とgot / wantを省略しません。Assertには値の探索ループやswitchを置かず、複数戻り値は別々に検証して失敗理由を一つに絞ります。",
    parameterized: "テーブル駆動", parameterizedBody: "同じ関数へ境界値を渡す反復なので、各行が完全な一件になるテーブル駆動テストを使います。外側のforはケース実行であり、Assertの中の探索ではありません。",
    goldenTitle: "同じシーンをLayoutの3契約で守る", goldenBody: `${golden.ja}を一つのfixtureとして固定します。外側サイズだけをmobile portrait・console 16:9・PC windowへ変え、各ケースで<code>Layout</code>が返した論理サイズに<code>Draw</code>して、<code>${goldenStem}-mobile.png</code>・<code>${goldenStem}-console.png</code>・<code>${goldenStem}-pc.png</code>と一枚まるごと比較します。これならゲーム状態の差ではなく画面構成の差だけをレビューできます。座標や描画命令は個別にAssertしません。`, goldenNote: "下の3ケースは、Layoutに縦長・16:9・PCウインドウの実際の分岐がある場合の形です。現在のLayoutが外側サイズを無視して一つの論理画面だけを返すなら、同一PNGを3枚承認せずdefaultの1ケースにします。分岐を追加した時点で3ケースへ広げ、モバイルの縦積み、コンソール／TVの安全領域、PCの横幅を.actual.pngで目視調整します。このコードはゲーム側のmain_test.goへ置き、ガイド冒頭のTestMain / assertGoldenと組み合わせます。",
  } : {
    complete: "Every step in the current Update", completeBody: "There are no ellipses. This is the running order from input and early returns through loops and win/lose. Unit-test the extracted pure rules directly.",
    update: "REAL GO / complete Update()", pure: "EXTRACTED / complete pure logic", test: "TEST / complete _test.go", table: "Choose boundaries first", manual: "What a person still checks", run: "Create these files in your own Go project, then run", back: "← PREVIOUS", forward: "NEXT →",
    aaaEyebrow: `AAA PRACTICE ${lesson.level} / 12`, aaaCount: `${scenarioCount} subtests have been reviewed individually. Each row is an independent AAA story, not merely another data value.`, repetitions: "Protect adjacent contracts with at least three stories", trapTitle: "AAA SMELL / a signal to split", practiceTitle: "ONE MORE TEST / write the next one",
    auditTitle: "Every test in this chapter passes both AAA and readability review.", names: "Naming", namesBody: `Top-level subjects are ${lesson.tests.map((name) => `<code>${name}</code>`).join(" and ")}. Every <code>t.Run</code> name states the premise and expected result in a sentence; there is no whenFoo_shouldBar convention.`,
    aaa: "AAA review result", aaaBody: aaa.why,
    damp: "DAMP", dampBody: "Meaningful case names and got / want values stay visible. Assert contains no search loop or switch, and multiple return values are checked separately so each failure says one thing.",
    parameterized: "Table-driven", parameterizedBody: "These are repeated boundary inputs to the same function, so every row is one complete table-driven case. The outer for runs cases; it is not a search hidden inside Assert.",
    goldenTitle: "Protect the same scene across three Layout contracts", goldenBody: `Freeze ${golden.en} as one fixture. Change only the outer size among mobile portrait, console 16:9, and PC window; in each case, draw at the logical size returned by <code>Layout</code> and compare the whole image with <code>${goldenStem}-mobile.png</code>, <code>${goldenStem}-console.png</code>, or <code>${goldenStem}-pc.png</code>. Reviewers then see layout differences rather than game-state differences. Do not assert individual coordinates or drawing commands.`, goldenNote: "The three cases below are the form to use when Layout has real portrait, 16:9, and PC-window branches. If the current Layout ignores the outer size and exposes one logical screen, keep one default case instead of approving three identical PNGs. Expand the table when the branches are introduced, then inspect .actual.png to tune mobile stacking, console/TV safe areas, and PC width. Put this code in the game's main_test.go and combine it with the TestMain / assertGolden harness from the guide hub.",
  };
  const logicPath = index < 2 ? "internal/lessonlogic/rules.go" : "internal/lessonlogic/core_updates.go";
  const testPath = index < 2 ? "internal/lessonlogic/rules_test.go" : "internal/lessonlogic/core_updates_test.go";
  const auditItems = [[c.names, c.namesBody], [c.aaa, c.aaaBody], [c.damp, c.dampBody], [c.parameterized, c.parameterizedBody]].map(([title, body]) => `<article><h3>${title}</h3><p>${body}</p></article>`).join("");
  const body = `<section class="test-step-hero"><a href="../">TESTING GUIDE</a><p class="eyebrow">LEVEL ${lesson.level} / 12</p><h1>${t.title}</h1><p>${t.lead}</p></section><section class="test-rule-strip"><span>EBITENGINE INPUT</span><i>→</i><strong>PURE RULE + go test</strong><i>→</i><span>GAME STATE</span></section><section class="test-explain"><div><p class="eyebrow">WHY THIS SEAM?</p><h2>${t.idea}</h2><p><strong>Draw:</strong> ${t.manual}</p></div><div class="test-case-table"><p class="eyebrow">${c.table}</p><table><tbody>${rows}</tbody></table></div></section><section class="test-update-map"><div><p class="eyebrow">NO OMISSIONS</p><h2>${c.complete}</h2><p>${c.completeBody}</p></div><ol>${phases}</ol></section><section class="test-code-full"><div><p class="eyebrow">${c.update}</p><code>${esc(lesson.source)}</code></div><pre><code>${esc(update)}</code></pre></section><section class="test-code-compare"><article><p>${c.pure}</p><code>${logicPath}</code><pre><code>${esc(pure)}</code></pre></article><article class="is-after"><p>${c.test}</p><code>${testPath}</code><pre><code>${esc(test)}</code></pre></article></section><section class="test-audit test-aaa-practice" data-aaa-lesson="${lesson.level}"><div><p class="eyebrow">${c.aaaEyebrow}</p><h2>${aaa.title}</h2><p>${aaa.intro}</p><p class="test-aaa-count">${c.aaaCount}</p></div><div class="test-quality-grid test-aaa-stage-grid">${aaaStages}</div><div class="test-aaa-repetitions"><p class="eyebrow">AAA REPETITIONS</p><h3>${c.repetitions}</h3><ol>${aaaStories}</ol></div><aside class="test-aaa-trap"><div><small>${c.trapTitle}</small><p>${aaa.trap}</p></div><div><small>${c.practiceTitle}</small><p>${aaa.practice}</p></div></aside></section><section class="test-audit"><div><p class="eyebrow">TEST REVIEW / ALL CASES</p><h2>${c.auditTitle}</h2></div><div class="test-quality-grid">${auditItems}</div></section><section class="test-golden test-golden-lesson"><div><p class="eyebrow">DRAW / GOLDEN ONLY</p><h2>${c.goldenTitle}</h2><p>${c.goldenBody}</p><p>${c.goldenNote}</p></div><pre><code>${esc(goldenTest)}</code></pre></section><section class="test-run"><p>${c.run}</p><code>go test ./internal/lessonlogic</code><strong>✓ PASS / NO WINDOW / NO DRAW</strong></section><section class="test-challenge"><p class="eyebrow">MANUAL CHECK IS STILL REAL</p><h2>${c.manual}</h2><p>${t.manual}</p></section><nav class="test-pager"><a href="${prev}">${c.back}</a><span>${index + 1} / ${lessons.length}</span><a href="${next}">${c.forward}</a></nav>`;
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
