#!/usr/bin/env node
/**
 * Replace every data-lab="entities" misuse with a concept-specific lab.
 * space-shooter becomes data-lab="bullets" (legitimate list lab, renamed).
 *
 * Usage: node scripts/wipe-entities-labs.mjs
 */
import { readFileSync, writeFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import { gated } from "./curriculum.mjs";

const root = new URL("..", import.meta.url).pathname;

function L(kind, ja, en, formula) {
  return { kind, ja, en, formula };
}

function F(jaEye, jaLines, jaP, enEye, enLines, enP) {
  return {
    ja: { eye: jaEye, lines: jaLines, p: jaP },
    en: { eye: enEye || jaEye, lines: enLines, p: enP },
  };
}

const specs = {
  "space-shooter": L(
    "bullets",
    {
      eye: "TRY IT / BULLET LIST",
      title: "弾を撃って、進めて、消そう",
      body: "「発射」で弾を配列へ追加し、「1フレーム」で全部の y を減らし、画面外を消します。敵も自機も、種類ごとのリストで同じ型を管理します。",
      hint: "本物のゲームも bullets = append / 移動 / フィルタ のくり返しです。",
      controls: [
        ["data-lab-fire", "発射", "lab-button-primary"],
        ["data-lab-step", "1フレーム進める"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-count", "弾の数"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-fire="発射！" data-step="飛翔"',
    },
    {
      eye: "TRY IT / BULLET LIST",
      title: "Spawn, move, cull bullets",
      body: "Fire appends to a slice; Step moves every bullet and drops off-screen ones. Separate lists per kind keep updates simple.",
      hint: "Real games loop append → move → filter the same way.",
      controls: [
        ["data-lab-fire", "Fire", "lab-button-primary"],
        ["data-lab-step", "Advance 1 frame"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [["data-lab-count", "bullets"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-fire="pew" data-step="fly"',
    },
    F("THE LIST RULE", ["append on fire", "move all, drop dead"], "1種類1配列にすると、当たり判定も描画も迷子になりません。",
      "THE LIST RULE", ["append on fire", "move all, drop dead"], "One slice per kind keeps collision and draw obvious."),
  ),
  "ebi-merge": L(
    "preview-next",
    {
      eye: "TRY IT / NEXT QUEUE",
      title: "次に落ちる段をキューで見せる",
      body: "予告は「次の物体」の待ち行列です。先頭を落とし、後ろへ新しい段を足します。得点やゲームオーバー判定とは別の箱です。",
      hint: "次の1個だけ見せるか、2〜3個見せるかも同じキューです。",
      controls: [
        ["data-lab-drop", "先頭を落とす", "lab-button-primary"],
        ["data-lab-enqueue", "予告を足す"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-next", "次"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-enqueued="予告追加" data-dropped="落下 T{n}" data-empty="空"',
    },
    {
      eye: "TRY IT / NEXT QUEUE",
      title: "Queue the next drop tiers",
      body: "The preview is a queue of upcoming pieces. Drop the front, enqueue behind. Scoring and game-over stay in other boxes.",
      hint: "Showing one or three previews is the same queue.",
      controls: [
        ["data-lab-drop", "Drop front", "lab-button-primary"],
        ["data-lab-enqueue", "Enqueue"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [["data-lab-next", "next"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-enqueued="queued" data-dropped="drop T{n}" data-empty="empty"',
    },
    F("PREVIEW ≠ BOARD", ["queue.front → spawn", "queue は盤面と別"], "予告を盤面配列へ混ぜないこと。",
      "PREVIEW ≠ BOARD", ["queue.front → spawn", "keep queue off the board"], "Never mix preview into the board slice."),
  ),
  "intent-status": L(
    "status-ticks",
    {
      eye: "TRY IT / STATUS TURNS",
      title: "効果の残りターンを減らそう",
      body: "毒などの状態は「名前＋残りターン」です。付与で追加し、ターン経過で left--。0で消えます。敵の予告行動も、同じ残り時間の考え方です。",
      hint: "毎ターン開始時にまとめて tick すると忘れません。",
      controls: [
        ["data-lab-add", "毒を付与", "lab-button-primary"],
        ["data-lab-tick", "ターン経過"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-count", "効果数"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-status="POISON" data-added="付与" data-ticked="ターン経過" data-none="状態なし"',
    },
    {
      eye: "TRY IT / STATUS TURNS",
      title: "Tick status durations down",
      body: "A status is a name plus remaining turns. Apply adds one; tick decrements; zero removes it. Enemy intents use the same leftover-time idea.",
      hint: "Tick every status once at turn start.",
      controls: [
        ["data-lab-add", "Apply poison", "lab-button-primary"],
        ["data-lab-tick", "Pass turn"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [["data-lab-count", "effects"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-status="POISON" data-added="applied" data-ticked="turn passed" data-none="no status"',
    },
    F("DURATION FIELD", ["apply: left = N", "tick: left-- ; drop if 0"], "効果ロジックと寿命を分けます。",
      "DURATION FIELD", ["apply: left = N", "tick: left-- ; drop if 0"], "Split effect logic from lifetime."),
  ),
  "card-rewards": L(
    "deck-pick",
    {
      eye: "TRY IT / REWARD PICK",
      title: "3枚から1枚をデッキへ足そう",
      body: "報酬は候補配列です。選んだ1枚だけを自分のデッキへ append。レア度は候補の作り方、編集は配列の操作です。",
      hint: "スキップ（足さない）も同じ画面の選択肢にできます。",
      controls: [["data-lab-reset", "デッキを戻す", "lab-button-quiet"]],
      values: [["data-lab-deck", "デッキ枚数"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-cards="Strike,Guard,Rare Heal" data-picked="{c} を追加"',
    },
    {
      eye: "TRY IT / REWARD PICK",
      title: "Add one of three rewards to the deck",
      body: "Rewards are a candidate list. Only the chosen card is appended to your deck. Rarity is how you build candidates; editing is slice ops.",
      hint: "Skip (add nothing) is just another choice on the same screen.",
      controls: [["data-lab-reset", "Reset deck", "lab-button-quiet"]],
      values: [["data-lab-deck", "deck size"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-cards="Strike,Guard,Rare Heal" data-picked="added {c}"',
    },
    F("PICK ONE", ["show candidates[]", "deck = append(deck, choice)"], "候補と所持デッキは別配列です。",
      "PICK ONE", ["show candidates[]", "deck = append(deck, choice)"], "Candidates and owned deck stay separate slices."),
  ),
  "branching-map": L(
    "map-nodes",
    {
      eye: "TRY IT / GRAPH PATH",
      title: "ノードを選んで道を進めよう",
      body: "マップはグラフです。今いるノードから出る辺だけが次の候補。選ぶたびに path へ履歴を残します。",
      hint: "戻れるマップにするなら、訪れたノード集合も別で持ちます。",
      controls: [["data-lab-reset", "スタートへ", "lab-button-quiet"]],
      values: [["data-lab-path", "経路"], ["data-lab-note", "今地"]],
      board: true,
      data: "",
    },
    {
      eye: "TRY IT / GRAPH PATH",
      title: "Walk a branching node graph",
      body: "The map is a graph. Only edges from the current node are legal next picks. Each choice appends to the path history.",
      hint: "Want revisits? Track a visited set separately.",
      controls: [["data-lab-reset", "Back to start", "lab-button-quiet"]],
      values: [["data-lab-path", "path"], ["data-lab-note", "here"]],
      board: true,
      data: "",
    },
    F("NODES + EDGES", ["choices = edges[current]", "path = append(path, next)"], "全部の部屋を直線配列にしないこと。",
      "NODES + EDGES", ["choices = edges[current]", "path = append(path, next)"], "Don’t flatten every room into one line."),
  ),
  "ebi-ascent": L(
    "pipeline",
    {
      eye: "TRY IT / RUN LOOP",
      title: "マップ→戦闘→報酬の場面をつなぐ",
      body: "大きなゲームは場面のパイプです。「次へ」で今の工程を進め、最後まで行ったら次の周回。状態は場面ごとに分けて持ちます。",
      hint: "セーブは場面の境目で取ると安全です。",
      controls: [
        ["data-lab-next", "次の場面へ", "lab-button-primary"],
        ["data-lab-reset", "最初へ", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "進捗"], ["data-lab-note", "今"]],
      board: true,
      data: 'data-steps="map,battle,reward,rest" data-loop="次のフロアへ"',
    },
    {
      eye: "TRY IT / RUN LOOP",
      title: "Pipe map → battle → reward",
      body: "A run is a pipeline of scenes. Next advances the stage; wrapping starts the next loop. Keep state per scene.",
      hint: "Save on scene boundaries.",
      controls: [
        ["data-lab-next", "Next scene", "lab-button-primary"],
        ["data-lab-reset", "Restart", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "progress"], ["data-lab-note", "now"]],
      board: true,
      data: 'data-steps="map,battle,reward,rest" data-loop="next floor"',
    },
    F("SCENE PIPE", ["currentScene.Update()", "on exit → next scene"], "全部を1つの巨大 Update に書かない。",
      "SCENE PIPE", ["currentScene.Update()", "on exit → next scene"], "Don’t cram every mode into one Update."),
  ),
  "ebi-match": L(
    "stage-goals",
    {
      eye: "TRY IT / STAGE DATA",
      title: "手数と目標スコアを同時に見よう",
      body: "ステージデータは「残り手数」と「目標」のセットです。マッチするたび手数を減らし得点を足す。0手で未達なら失敗です。",
      hint: "盤面レイアウトも同じステージJSONに置けます。",
      controls: [
        ["data-lab-match", "マッチする", "lab-button-primary"],
        ["data-lab-reset", "ステージを戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-moves", "残り手数"], ["data-lab-score", "スコア"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-goal="30" data-moves="10" data-hit="+8点" data-clear="クリア！" data-out="手数切れ"',
    },
    {
      eye: "TRY IT / STAGE DATA",
      title: "Track moves and a score goal",
      body: "Stage data pairs remaining moves with a goal. Each match spends a move and adds score. Zero moves without the goal is failure.",
      hint: "Board layouts can live in the same stage JSON.",
      controls: [
        ["data-lab-match", "Make a match", "lab-button-primary"],
        ["data-lab-reset", "Reset stage", "lab-button-quiet"],
      ],
      values: [["data-lab-moves", "moves"], ["data-lab-score", "score"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-goal="30" data-moves="10" data-hit="+8" data-clear="CLEAR" data-out="no moves"',
    },
    F("MOVES + GOAL", ["moves-- ; score += n", "win if score >= goal"], "ルールとステージ数値を分離します。",
      "MOVES + GOAL", ["moves-- ; score += n", "win if score >= goal"], "Keep rules separate from stage numbers."),
  ),
  "tetromino-shapes": L(
    "shape-cells",
    {
      eye: "TRY IT / SHAPE DATA",
      title: "相対座標の4マスを回転させよう",
      body: "形は「中心からの相対マス」のリストです。形を切り替え、回転で (x,y)→(-y,x) を試します。色や名前はデータ、回転は共通の式です。",
      hint: "O ミノは回転しても見た目が同じになるのが正解です。",
      controls: [
        ["data-lab-next", "次の形", "lab-button-primary"],
        ["data-lab-rot", "90°回転"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-name", "形"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-rotated="回転した"',
    },
    {
      eye: "TRY IT / SHAPE DATA",
      title: "Rotate four relative cells",
      body: "A piece is a list of cells relative to an origin. Swap shapes, then rotate with (x,y)→(-y,x). Names/colors are data; rotation is shared math.",
      hint: "O looking unchanged after rotate is correct.",
      controls: [
        ["data-lab-next", "Next shape", "lab-button-primary"],
        ["data-lab-rot", "Rotate 90°"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [["data-lab-name", "shape"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-rotated="rotated"',
    },
    F("RELATIVE CELLS", ["cells = shape[rot]", "world = origin + cell"], "絶対座標だけで形を持つと回転が地獄です。",
      "RELATIVE CELLS", ["cells = shape[rot]", "world = origin + cell"], "Absolute-only shapes make rotation miserable."),
  ),
  "rotation-kicks": L(
    "kick-try",
    {
      eye: "TRY IT / WALL KICKS",
      title: "ずらし候補を順番に試そう",
      body: "回転が壁に当たったら、(0,0)→(1,0)→(-1,0)…とキック表を順に試します。最初に通った候補で確定。全部だめなら回転キャンセルです。",
      hint: "SRS も同じ「候補を上から試す」構造です。",
      controls: [
        ["data-lab-kick", "次のキックを試す", "lab-button-primary"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-try", "試行"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-kicks="0,0;1,0;-1,0;0,-1" data-ok="kick ({dx},{dy}) OK" data-blocked="kick ({dx},{dy}) 壁" data-fail="全部失敗→キャンセル"',
    },
    {
      eye: "TRY IT / WALL KICKS",
      title: "Try kick offsets in order",
      body: "If a rotate hits a wall, walk a kick table: (0,0), (1,0), (-1,0)… First legal offset wins. None work → cancel the rotate.",
      hint: "SRS is the same “try offsets top-down” shape.",
      controls: [
        ["data-lab-kick", "Try next kick", "lab-button-primary"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [["data-lab-try", "tries"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-kicks="0,0;1,0;-1,0;0,-1" data-ok="kick ({dx},{dy}) OK" data-blocked="kick ({dx},{dy}) blocked" data-fail="all failed → cancel"',
    },
    F("TRY TABLE", ["for offset in kicks", "if fits { apply; break }"], "キック表はデータ、判定は共通関数。",
      "TRY TABLE", ["for offset in kicks", "if fits { apply; break }"], "Kick tables are data; the tester is shared."),
  ),
  "bag-hold-ghost": L(
    "bag-draw",
    {
      eye: "TRY IT / 7-BAG",
      title: "袋から重複なしで取り出そう",
      body: "7種を袋に入れ、空になるまで重複なしで引きます。HOLD は今のピースと入れ替え。公平な乱択の基本形です。",
      hint: "袋が空なら7種を入れ直します。",
      controls: [
        ["data-lab-draw", "袋から引く", "lab-button-primary"],
        ["data-lab-hold", "HOLD と交換"],
        ["data-lab-reset", "袋を詰め直す", "lab-button-quiet"],
      ],
      values: [["data-lab-left", "袋の残り"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-drew="{p} を引いた" data-held="HOLD 交換"',
    },
    {
      eye: "TRY IT / 7-BAG",
      title: "Draw without replacement",
      body: "Fill a bag with seven pieces and draw until empty—no dups mid-bag. HOLD swaps with the current piece. Classic fair randomizer.",
      hint: "Refill all seven when the bag empties.",
      controls: [
        ["data-lab-draw", "Draw", "lab-button-primary"],
        ["data-lab-hold", "Swap HOLD"],
        ["data-lab-reset", "Refill bag", "lab-button-quiet"],
      ],
      values: [["data-lab-left", "left in bag"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-drew="drew {p}" data-held="hold swapped"',
    },
    F("BAG RANDOM", ["if bag empty { refill 7 }", "draw = remove random"], "完全な乱数連打より袋の方が偏りにくいです。",
      "BAG RANDOM", ["if bag empty { refill 7 }", "draw = remove random"], "A bag fights droughts better than raw RNG."),
  ),
  "ebi-blocks": L(
    "pipeline",
    {
      eye: "TRY IT / FRAME JOBS",
      title: "1入力が通る仕事を順番に並べる",
      body: "落下ブロック全体はパイプです。入力→移動→固定→ライン消し→次ピース。同じフレームで順番を守るとバグが減ります。",
      hint: "ゴースト表示は描画専用で、このパイプの外に置けます。",
      controls: [
        ["data-lab-next", "次の工程", "lab-button-primary"],
        ["data-lab-reset", "フレーム先頭へ", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "工程"], ["data-lab-note", "今"]],
      board: true,
      data: 'data-steps="input,fall,lock,clear,spawn" data-loop="次フレーム"',
    },
    {
      eye: "TRY IT / FRAME JOBS",
      title: "Order the jobs one input travels",
      body: "A full stacker is a pipe: input → fall → lock → clear → spawn. Keep that order inside one frame.",
      hint: "Ghost piece can stay draw-only, outside the pipe.",
      controls: [
        ["data-lab-next", "Next job", "lab-button-primary"],
        ["data-lab-reset", "Frame start", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "step"], ["data-lab-note", "now"]],
      board: true,
      data: 'data-steps="input,fall,lock,clear,spawn" data-loop="next frame"',
    },
    F("ORDERED PIPE", ["do jobs in fixed order", "never skip ahead"], "消す前に固定、固定前に移動。",
      "ORDERED PIPE", ["do jobs in fixed order", "never skip ahead"], "Lock before clear; move before lock."),
  ),
  "ally-effects": L(
    "event-queue",
    {
      eye: "TRY IT / CONTACT QUEUE",
      title: "接触イベントを1件ずつ処理",
      body: "接触したらすぐ能力を発動せず、イベントをキューへ入れます。「処理」で先頭から1件解決。同時に複数触れても順番が保てます。",
      hint: "キューが空になるまで Update で消化します。",
      controls: [
        ["data-lab-push", "接触を追加", "lab-button-primary"],
        ["data-lab-pop", "1件処理"],
        ["data-lab-reset", "空にする", "lab-button-quiet"],
      ],
      values: [["data-lab-count", "待ち"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-event="ally" data-pushed="キューへ" data-popped="{e} を解決" data-empty="空"',
    },
    {
      eye: "TRY IT / CONTACT QUEUE",
      title: "Resolve contact events one-by-one",
      body: "On contact, enqueue an event instead of firing instantly. Pop resolves the front. Simultaneous hits keep a stable order.",
      hint: "Drain the queue inside Update until empty.",
      controls: [
        ["data-lab-push", "Queue contact", "lab-button-primary"],
        ["data-lab-pop", "Resolve one"],
        ["data-lab-reset", "Clear", "lab-button-quiet"],
      ],
      values: [["data-lab-count", "waiting"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-event="ally" data-pushed="queued" data-popped="resolve {e}" data-empty="empty"',
    },
    F("QUEUE THEN RESOLVE", ["on contact: push event", "later: pop & apply"], "物理の瞬間と効果解決を分離します。",
      "QUEUE THEN RESOLVE", ["on contact: push event", "later: pop & apply"], "Split the physics moment from effect resolve."),
  ),
  "ebi-strike": L(
    "shot-event-queue",
    {
      eye: "TRY IT / SHOT EVENTS",
      title: "弾の接触とイベントキューを同時に見る",
      body: "「1ショット発射」で弾が軌道を進み、敵A・B・Cへ触れた瞬間にHITイベントが右のキューへ積まれます。「先頭を1件処理」で古いイベントからHPへ反映される様子を比べましょう。",
      hint: "物理は接触した順にpushし、ゲームルールは先頭からpopします。軌道と待ち行列を同時に追ってください。",
      controls: [
        ["data-lab-launch", "1ショット発射", "lab-button-primary"],
        ["data-lab-pop", "先頭を1件処理"],
        ["data-lab-reset", "最初から", "lab-button-quiet"],
      ],
      values: [["data-lab-hits", "接触"], ["data-lab-count", "キュー"], ["data-lab-note", "今の処理"]],
      board: true,
      data: 'data-events="HIT A,HIT B,HIT C,TURN END" data-launched="弾が軌道を進んでいます" data-contact="接触 → {e} をpush" data-finished="停止 → TURN ENDをpush" data-resolved="{e} をpop → 敵{target} HP -1" data-turn="TURN ENDをpop → 次の手番へ" data-empty="キューは空です"',
    },
    {
      eye: "TRY IT / SHOT EVENTS",
      title: "Watch contacts and the event queue together",
      body: "Launch one shot and watch it cross enemies A, B, and C. Each contact pushes a HIT event into the queue on the right. Resolve Front pops the oldest event and applies it to HP.",
      hint: "Physics pushes in contact order; game rules pop from the front. Follow the path and queue at the same time.",
      controls: [
        ["data-lab-launch", "Launch one shot", "lab-button-primary"],
        ["data-lab-pop", "Resolve front"],
        ["data-lab-reset", "Start over", "lab-button-quiet"],
      ],
      values: [["data-lab-hits", "contacts"], ["data-lab-count", "queued"], ["data-lab-note", "now"]],
      board: true,
      data: 'data-events="HIT A,HIT B,HIT C,TURN END" data-launched="shot moving along its path" data-contact="contact → push {e}" data-finished="stopped → push TURN END" data-resolved="pop {e} → enemy {target} HP -1" data-turn="pop TURN END → pass the turn" data-empty="the queue is empty"',
    },
    F("TURN OWNS THE QUEUE", ["physics pushes events", "turn pops events"], "物理と手番ロジックの境界です。",
      "TURN OWNS THE QUEUE", ["physics pushes events", "turn pops events"], "That’s the physics / turn boundary."),
  ),
  "terrain-generation": L(
    "height-layers",
    {
      eye: "TRY IT / HEIGHT FIELD",
      title: "列ごとの高さをノイズで決めよう",
      body: "地形は「列→高さ」です。ノイズでサンプリングし、必要なら削って地層を作ります。seed を固定すれば同じ大地が再現できます。",
      hint: "高さの下を石、上を土、みたいな層分けはあとから載せます。",
      controls: [
        ["data-lab-noise", "ノイズサンプリング", "lab-button-primary"],
        ["data-lab-carve", "1段削る"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-h", "高さ"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-sampled="サンプリング" data-carved="削った"',
    },
    {
      eye: "TRY IT / HEIGHT FIELD",
      title: "Sample a height per column",
      body: "Terrain starts as column→height. Sample with noise, carve if needed. A fixed seed rebuilds the same land.",
      hint: "Layer stone/dirt on top of the height field afterward.",
      controls: [
        ["data-lab-noise", "Sample noise", "lab-button-primary"],
        ["data-lab-carve", "Carve 1"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [["data-lab-h", "heights"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-sampled="sampled" data-carved="carved"',
    },
    F("HEIGHT THEN LAYERS", ["h = noise(x, seed)", "fill cells below h"], "いきなり2Dノイズ全部より、高さからが簡単です。",
      "HEIGHT THEN LAYERS", ["h = noise(x, seed)", "fill cells below h"], "Heights first beat full 2D noise for starters."),
  ),
  "inventory-crafting": L(
    "craft-recipe",
    {
      eye: "TRY IT / RECIPE",
      title: "材料を集めてレシピで作ろう",
      body: "インベントリは個数の辞書、クラフトはレシピ照合です。足りれば材料を減らし成果を足す。スタックも同じ個数管理です。",
      hint: "レシピをデータ表にすると、新しい道具が行追加だけで増えます。",
      controls: [
        ["data-lab-wood", "+木"],
        ["data-lab-string", "+糸"],
        ["data-lab-craft", "クラフト", "lab-button-primary"],
        ["data-lab-reset", "空にする", "lab-button-quiet"],
      ],
      values: [["data-lab-inv", "所持"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-need="材料不足" data-made="弓を作った"',
    },
    {
      eye: "TRY IT / RECIPE",
      title: "Gather materials, craft from a recipe",
      body: "Inventory is counts; crafting is recipe matching. If costs are met, subtract inputs and add the output. Stacks are the same counters.",
      hint: "Recipes as data rows mean new tools are new lines.",
      controls: [
        ["data-lab-wood", "+wood"],
        ["data-lab-string", "+string"],
        ["data-lab-craft", "Craft", "lab-button-primary"],
        ["data-lab-reset", "Clear", "lab-button-quiet"],
      ],
      values: [["data-lab-inv", "inventory"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-need="need materials" data-made="crafted bow"',
    },
    F("PAY THEN GAIN", ["if inv has costs", "inv -= costs; inv += result"], "失敗時に材料だけ減らないよう、先に検査します。",
      "PAY THEN GAIN", ["if inv has costs", "inv -= costs; inv += result"], "Check first so a failed craft never eats items."),
  ),
  "tools-light": L(
    "light-flood",
    {
      eye: "TRY IT / LIGHT BFS",
      title: "光源から明るさを広げよう",
      body: "トーチを置くと中心が明るい値になります。「拡散」で上下左右へ1ずつ減衰してコピー。道具の耐久とは別計算です。",
      hint: "壁で止めたいときは、壁マスへ値を書き込まないだけでできます。",
      controls: [
        ["data-lab-torch", "トーチ設置", "lab-button-primary"],
        ["data-lab-flood", "1回拡散"],
        ["data-lab-reset", "消灯", "lab-button-quiet"],
      ],
      values: [["data-lab-lit", "明るいマス"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-torch="設置" data-flooded="拡散"',
    },
    {
      eye: "TRY IT / LIGHT BFS",
      title: "Flood brightness from a torch",
      body: "Place a torch to set a bright center. Flood copies to neighbors with decay. Tool durability stays a separate system.",
      hint: "Block walls by simply never writing light into them.",
      controls: [
        ["data-lab-torch", "Place torch", "lab-button-primary"],
        ["data-lab-flood", "Flood once"],
        ["data-lab-reset", "Lights out", "lab-button-quiet"],
      ],
      values: [["data-lab-lit", "lit cells"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-torch="placed" data-flooded="flooded"',
    },
    F("DECAY FLOOD", ["center = power", "neighbor = max(old, power-1)"], "光は加算でも上書きでもよいが、ルールを1つに。",
      "DECAY FLOOD", ["center = power", "neighbor = max(old, power-1)"], "Additive or max—pick one rule and stick to it."),
  ),
  "ebi-craft": L(
    "pipeline",
    {
      eye: "TRY IT / WORLD RESTORE",
      title: "保存したい状態を順番に集める",
      body: "ワールド復元は「何をセーブするか」のパイプです。シード→地形→インベントリ→位置。読み込みも同じ順だと壊れません。",
      hint: "画像そのものは保存せず、シードと差分だけにします。",
      controls: [
        ["data-lab-next", "次の項目", "lab-button-primary"],
        ["data-lab-reset", "先頭へ", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "項目"], ["data-lab-note", "今"]],
      board: true,
      data: 'data-steps="seed,chunks,inventory,player" data-loop="セーブ完了→再開"',
    },
    {
      eye: "TRY IT / WORLD RESTORE",
      title: "Collect savable state in order",
      body: "Restoring a world is a pipe of what to save: seed → chunks → inventory → player. Load in the same order.",
      hint: "Save seed + diffs, not raw images.",
      controls: [
        ["data-lab-next", "Next field", "lab-button-primary"],
        ["data-lab-reset", "Start", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "field"], ["data-lab-note", "now"]],
      board: true,
      data: 'data-steps="seed,chunks,inventory,player" data-loop="saved → reload"',
    },
    F("SAVE BOUNDARY", ["write fields in order", "read fields in same order"], "途中のフィールドだけ古い版、を避ける。",
      "SAVE BOUNDARY", ["write fields in order", "read fields in same order"], "Avoid half-old save layouts."),
  ),
  "species-data": L(
    "species-inst",
    {
      eye: "TRY IT / DEF vs INST",
      title: "種族定義から個体を増やそう",
      body: "種族データは変えません（maxHPなど）。個体は定義を参照しつつ、今の HP だけ別持ち。「スポーン」で個体増、「ヒット」で個体だけ減ります。",
      hint: "図鑑は定義の一覧、パーティは個体の一覧です。",
      controls: [
        ["data-lab-spawn", "個体を出す", "lab-button-primary"],
        ["data-lab-hit", "先頭を攻撃"],
        ["data-lab-reset", "個体を消す", "lab-button-quiet"],
      ],
      values: [["data-lab-count", "個体数"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-maxhp="8" data-spawned="スポーン" data-hit="ダメージ" data-none="個体なし"',
    },
    {
      eye: "TRY IT / DEF vs INST",
      title: "Spawn instances from a species def",
      body: "Species data stays constant (maxHP…). Instances reference it but keep their own HP. Spawn adds instances; hit damages only them.",
      hint: "Dex = definitions; party = instances.",
      controls: [
        ["data-lab-spawn", "Spawn", "lab-button-primary"],
        ["data-lab-hit", "Hit first"],
        ["data-lab-reset", "Clear", "lab-button-quiet"],
      ],
      values: [["data-lab-count", "instances"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-maxhp="8" data-spawned="spawned" data-hit="hit" data-none="none"',
    },
    F("DEF ≠ STATE", ["Species.MaxHP", "Instance.HP"], "定義を直接減らすと全個体が壊れます。",
      "DEF ≠ STATE", ["Species.MaxHP", "Instance.HP"], "Mutating the def corrupts every instance."),
  ),
  "party-switch": L(
    "party-swap",
    {
      eye: "TRY IT / ACTIVE INDEX",
      title: "前衛インデックスを切り替えよう",
      body: "パーティは配列、前衛はその添え字です。「交代」で index を回します。強制交代も同じ index 書き換えです。",
      hint: "控えが0なら交代ボタンを押せなくします。",
      controls: [
        ["data-lab-swap", "交代", "lab-button-primary"],
        ["data-lab-reset", "先頭へ", "lab-button-quiet"],
      ],
      values: [["data-lab-active", "前衛"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-party="Ebi,Shell,Coral" data-swapped="前衛 → {p}"',
    },
    {
      eye: "TRY IT / ACTIVE INDEX",
      title: "Rotate the front-line index",
      body: "The party is an array; the front-liner is an index. Swap advances that index. Forced swaps write the same field.",
      hint: "Disable swap when only one member remains.",
      controls: [
        ["data-lab-swap", "Switch", "lab-button-primary"],
        ["data-lab-reset", "First", "lab-button-quiet"],
      ],
      values: [["data-lab-active", "active"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-party="Ebi,Shell,Coral" data-swapped="active → {p}"',
    },
    F("INDEX ONLY", ["active = (active+1) % len", "draw party[active]"], "配列を並べ替えなくても交代できます。",
      "INDEX ONLY", ["active = (active+1) % len", "draw party[active]"], "You can swap without reshuffling the array."),
  ),
  capture: L(
    "capture-roll",
    {
      eye: "TRY IT / CATCH RATE",
      title: "捕獲率を上げてロールしよう",
      body: "捕獲は確率です。エサで率を上げ、「ロール」で 0–99 を振って率未満なら成功。失敗リスクも同じ乱数の反対側です。",
      hint: "所持枠チェックはロールの前に行います。",
      controls: [
        ["data-lab-bait", "エサ +15%", "lab-button-primary"],
        ["data-lab-roll", "ロール"],
        ["data-lab-reset", "率を戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-rate", "捕獲率"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-rate="35" data-baited="エサ" data-caught="捕獲成功 ({r})" data-missed="逃げられた ({r})"',
    },
    {
      eye: "TRY IT / CATCH RATE",
      title: "Boost the rate, then roll",
      body: "Capture is a percentage. Bait raises it; Roll picks 0–99 and succeeds if below the rate. Failure is the other side of the same RNG.",
      hint: "Check party capacity before rolling.",
      controls: [
        ["data-lab-bait", "Bait +15%", "lab-button-primary"],
        ["data-lab-roll", "Roll"],
        ["data-lab-reset", "Reset rate", "lab-button-quiet"],
      ],
      values: [["data-lab-rate", "rate"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-rate="35" data-baited="baited" data-caught="caught! ({r})" data-missed="broke free ({r})"',
    },
    F("RATE THEN RNG", ["rate = base + modifiers", "success if rand < rate"], "表示する率と内部判定を一致させます。",
      "RATE THEN RNG", ["rate = base + modifiers", "success if rand < rate"], "Keep the shown rate identical to the check."),
  ),
  "growth-evolution": L(
    "xp-level",
    {
      eye: "TRY IT / XP CURVE",
      title: "経験値をためてレベルを上げよう",
      body: "必要 XP はレベルで増えます。足りたら level++ して余りを持ち越し。進化は「レベルが閾値以上」の判定を同じ場所に足します。",
      hint: "曲線を変えたいときは need() だけ触ります。",
      controls: [
        ["data-lab-gain", "+4 XP", "lab-button-primary"],
        ["data-lab-reset", "Lv1へ", "lab-button-quiet"],
      ],
      values: [["data-lab-level", "レベル"], ["data-lab-xp", "XP"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-gained="+4 XP" data-evolve="進化できる！"',
    },
    {
      eye: "TRY IT / XP CURVE",
      title: "Bank XP until you level",
      body: "XP needed grows with level. On overflow, level++ and keep the remainder. Evolution is “level >= threshold” beside the same code.",
      hint: "Retune difficulty by editing need() only.",
      controls: [
        ["data-lab-gain", "+4 XP", "lab-button-primary"],
        ["data-lab-reset", "Back to Lv1", "lab-button-quiet"],
      ],
      values: [["data-lab-level", "level"], ["data-lab-xp", "XP"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-gained="+4 XP" data-evolve="ready to evolve"',
    },
    F("CURVE FUNCTION", ["need = level * 10", "while xp >= need { level++ }"], "必要量を表にしても、関数でもよいです。",
      "CURVE FUNCTION", ["need = level * 10", "while xp >= need { level++ }"], "Tables or functions both work for need()."),
  ),
  "ebi-monsters": L(
    "pipeline",
    {
      eye: "TRY IT / ADVENTURE FLOW",
      title: "冒険イベントを順に解決する",
      body: "図鑑・捕獲・育成・バトルをつなぐのもパイプです。今のイベントを終え、次へ進む。セーブ境界はこの境目に置きます。",
      hint: "戦闘中セーブを許すなら、もっと細かい境界が必要です。",
      controls: [
        ["data-lab-next", "次のイベント", "lab-button-primary"],
        ["data-lab-reset", "冒頭へ", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "進捗"], ["data-lab-note", "今"]],
      board: true,
      data: 'data-steps="field,battle,capture,party,save" data-loop="次の冒険へ"',
    },
    {
      eye: "TRY IT / ADVENTURE FLOW",
      title: "Resolve adventure events in order",
      body: "Dex, capture, growth, and battle still form a pipe. Finish the current event, advance, and save on those edges.",
      hint: "Mid-battle saves need finer boundaries.",
      controls: [
        ["data-lab-next", "Next event", "lab-button-primary"],
        ["data-lab-reset", "Intro", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "progress"], ["data-lab-note", "now"]],
      board: true,
      data: 'data-steps="field,battle,capture,party,save" data-loop="next outing"',
    },
    F("EVENT PIPE", ["finish current", "go next / save"], "全部同時に動かさない。",
      "EVENT PIPE", ["finish current", "go next / save"], "Don’t run every system at once."),
  ),
  "pair-rotation": L(
    "kick-try",
    {
      eye: "TRY IT / PAIR KICKS",
      title: "ペア回転のずらしを試そう",
      body: "2個ペアも回転後に壁へ食い込むことがあります。キック表で少しずつずらして、最初に収まる位置を採用します。",
      hint: "4方向それぞれのキック表を持てます。",
      controls: [
        ["data-lab-kick", "補正を試す", "lab-button-primary"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-try", "試行"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-kicks="0,0;1,0;-1,0;0,1" data-ok="補正 ({dx},{dy}) OK" data-blocked="({dx},{dy}) 不可" data-fail="回転キャンセル"',
    },
    {
      eye: "TRY IT / PAIR KICKS",
      title: "Kick a rotating pair",
      body: "A two-cell pair can clip a wall after rotate. Walk kick offsets and keep the first fit.",
      hint: "You can store a kick table per facing.",
      controls: [
        ["data-lab-kick", "Try kick", "lab-button-primary"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [["data-lab-try", "tries"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-kicks="0,0;1,0;-1,0;0,1" data-ok="kick ({dx},{dy}) OK" data-blocked="({dx},{dy}) blocked" data-fail="cancel rotate"',
    },
    F("SAME KICK IDEA", ["rotate pair", "try offsets until fit"], "テトロミノと同じ発想を2マスに縮小しただけです。",
      "SAME KICK IDEA", ["rotate pair", "try offsets until fit"], "Same kick idea, just two cells."),
  ),
  "color-groups": L(
    "bfs-flood",
    {
      eye: "TRY IT / BFS GROUP",
      title: "マスをタップして同色連結を塗ろう",
      body: "4方向に同じ色が続くマスを BFS で集めます。訪問済みを忘れないこと。消す判定は「連結数 ≥ 4」などこの集合サイズで行います。",
      hint: "タップしたマスの色が探索のキーです。",
      controls: [["data-lab-reset", "塗りを消す", "lab-button-quiet"]],
      values: [["data-lab-count", "連結数"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-group="連結 {n} マス"',
    },
    {
      eye: "TRY IT / BFS GROUP",
      title: "Tap a cell to paint its color group",
      body: "BFS gathers 4-way neighbors of the same color. Remember visited. Clear rules use this set’s size (e.g. ≥ 4).",
      hint: "The tapped cell’s color is the search key.",
      controls: [["data-lab-reset", "Clear paint", "lab-button-quiet"]],
      values: [["data-lab-count", "group size"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-group="group size {n}"',
    },
    F("BFS SET", ["queue + visited", "same color → expand"], "再帰DFSでもよいが、深い盤面ではBFSが安全です。",
      "BFS SET", ["queue + visited", "same color → expand"], "DFS works too; BFS is safer on deep boards."),
  ),
  "pellet-map": L(
    "pellet-count",
    {
      eye: "TRY IT / PELLET LAYER",
      title: "残ドットを1個ずつ消そう",
      body: "ドットは迷路と別レイヤーの bool 配列です。通ったら false にし、残り0でクリア。壁データは変えません。",
      hint: "初期配置だけマップデータ、実行中の残りは状態です。",
      controls: [
        ["data-lab-eat", "1個取る", "lab-button-primary"],
        ["data-lab-reset", "配置を戻す", "lab-button-quiet"],
      ],
      values: [["data-lab-left", "残り"], ["data-lab-note", "結果"]],
      board: true,
      data: 'data-ate="取得" data-clear="全消し！"',
    },
    {
      eye: "TRY IT / PELLET LAYER",
      title: "Clear remaining pellets one by one",
      body: "Pellets are a bool layer beside the maze. Stepping on one sets false; zero left clears the stage. Walls stay untouched.",
      hint: "Initial layout is data; remaining counts are runtime state.",
      controls: [
        ["data-lab-eat", "Eat one", "lab-button-primary"],
        ["data-lab-reset", "Restore", "lab-button-quiet"],
      ],
      values: [["data-lab-left", "left"], ["data-lab-note", "result"]],
      board: true,
      data: 'data-ate="ate" data-clear="all clear!"',
    },
    F("LAYER STATE", ["if pellet[y][x] { eat; left-- }", "clear if left == 0"], "迷路文字を直接書き換えて消さない。",
      "LAYER STATE", ["if pellet[y][x] { eat; left-- }", "clear if left == 0"], "Don’t erase pellets by rewriting maze glyphs."),
  ),
  "ebi-maze": L(
    "pipeline",
    {
      eye: "TRY IT / FRAME ORDER",
      title: "1フレームの仕事を順番に処理",
      body: "迷路総合はパイプです。入力→移動→敵AI→ドット→描画準備。順序を守ると、食べた直後の敵接触なども説明しやすくなります。",
      hint: "描画は最後。ここより前で画面を触らない。",
      controls: [
        ["data-lab-next", "次の仕事", "lab-button-primary"],
        ["data-lab-reset", "フレーム先頭", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "工程"], ["data-lab-note", "今"]],
      board: true,
      data: 'data-steps="input,move,enemy,pellets,draw" data-loop="次フレーム"',
    },
    {
      eye: "TRY IT / FRAME ORDER",
      title: "Run one frame’s jobs in order",
      body: "The full maze game is a pipe: input → move → enemy AI → pellets → draw prep. Fixed order makes edge cases explainable.",
      hint: "Draw last—never touch pixels earlier.",
      controls: [
        ["data-lab-next", "Next job", "lab-button-primary"],
        ["data-lab-reset", "Frame start", "lab-button-quiet"],
      ],
      values: [["data-lab-stepn", "step"], ["data-lab-note", "now"]],
      board: true,
      data: 'data-steps="input,move,enemy,pellets,draw" data-loop="next frame"',
    },
    F("FIXED ORDER", ["run systems in list order", "draw is last"], "AIが先かドットが先かで勝敗が変わることがあります。",
      "FIXED ORDER", ["run systems in list order", "draw is last"], "AI-before-pellets can change who wins a tie."),
  ),
};

function btn([attr, label, cls = "", value = ""]) {
  const klass = ` class="lab-button${cls ? ` ${cls}` : ""}"`;
  if (attr === "data-lab-card") {
    return `<button type="button"${klass} data-lab-card="${value}">${label}</button>`;
  }
  return `<button type="button"${klass} ${attr}>${label}</button>`;
}

function labHTML(lang, slug, spec) {
  const c = spec[lang];
  const controls = c.controls.map(btn).join("");
  const values = c.values
    .map(([attr, label]) => {
      let initial = "0";
      if (attr.includes("note") || attr.includes("path") || attr.includes("name") || attr.includes("active") || attr.includes("next") || attr.includes("inv") || attr.includes("h")) initial = "—";
      if (attr.includes("rate")) initial = "35%";
      if (attr.includes("score")) initial = "0/30";
      if (attr.includes("xp")) initial = "0/10";
      if (attr.includes("level")) initial = "1";
      if (attr.includes("moves")) initial = "10";
      if (attr.includes("energy")) initial = "3/3";
      return `<div><span>${label}</span><strong ${attr}>${initial}</strong></div>`;
    })
    .join("");
  const board = c.board ? `<div class="lab-board" data-lab-board role="img" aria-label="lab"></div>` : "";
  const id = `lab-title-${slug}`;
  return `<div class="motion-lab" data-lab="${spec.kind}" ${c.data || ""} aria-labelledby="${id}">
<div class="lab-copy">
<p class="eyebrow">${c.eye}</p>
<h3 id="${id}">${c.title}</h3>
<p>${c.body}</p>
<div class="lab-controls">${controls}</div>
<p class="lab-hint">${c.hint}</p>
</div>
<div class="lab-visual">
${board}
<div class="lab-values" aria-live="polite">${values}</div>
</div>
</div>`;
}

function formulaHTML(lang, formula) {
  if (!formula) return "";
  const f = formula[lang];
  const lines = f.lines
    .map((line, i) =>
      i === 0
        ? `<code>${line}</code>`
        : `<span>${lang === "ja" ? "つぎに" : "then"}</span><code>${line}</code>`,
    )
    .join("");
  return `<div class="formula"><p class="eyebrow">${f.eye}</p><div class="formula-lines">${lines}</div><p>${f.p}</p></div>`;
}

function patch(html, lang, slug, spec) {
  const start = html.search(/<div class="motion-lab"/);
  if (start < 0) return null;
  const after = html.slice(start);
  const endRel = after.search(/<div class="(formula|code-lesson|why-grid)"/);
  if (endRel < 0) return null;
  const end = start + endRel;
  let insert = labHTML(lang, slug, spec);
  if (!/class="formula"/.test(html) && spec.formula) {
    insert += formulaHTML(lang, spec.formula);
  }
  return html.slice(0, start) + insert + html.slice(end);
}

function findRoute(slug) {
  for (const e of gated) {
    if (e.slug === slug) return e.route;
  }
  return null;
}

let updated = 0;
let missing = [];
const requested = new Set(process.argv.slice(2));
for (const [slug, spec] of Object.entries(specs)) {
  if (requested.size && !requested.has(slug)) continue;
  const route = findRoute(slug);
  if (!route) {
    missing.push(slug);
    continue;
  }
  for (const lang of ["ja", "en"]) {
    const path = join(root, "web", lang, route, "index.html");
    if (!existsSync(path)) {
      missing.push(`${lang}:${slug}`);
      continue;
    }
    const html = readFileSync(path, "utf8");
    if (!html.includes('data-lab="entities"') && !requested.has(slug)) {
      console.log("skip (no entities)", lang, route);
      continue;
    }
    const next = patch(html, lang, slug, spec);
    if (!next) {
      console.warn("patch failed", lang, route);
      continue;
    }
    writeFileSync(path, next);
    updated++;
    console.log("wiped", lang, route, "→", spec.kind);
  }
}

// Final sweep: any remaining entities in gated playable?
let leftover = [];
for (const e of gated.filter((x) => x.playable)) {
  for (const lang of ["ja", "en"]) {
    const path = join(root, "web", lang, e.route, "index.html");
    if (!existsSync(path)) continue;
    const html = readFileSync(path, "utf8");
    if (/data-lab="entities"/.test(html)) leftover.push(`${lang}:${e.id}`);
  }
}

console.log(`Updated ${updated}. Missing specs/routes: ${missing.length ? missing.join(", ") : "none"}`);
console.log(`Leftover entities: ${leftover.length ? leftover.join(", ") : "NONE ✅"}`);
