#!/usr/bin/env node
/**
 * Replace mismatched/generic TRY IT labs on thin track articles with
 * concept-specific interactive labs, and add a formula block when missing.
 *
 * Usage: node scripts/enrich-thin-labs.mjs
 */
import { readFileSync, writeFileSync, existsSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;

/** @typedef {{ kind: string, attrs?: string, ja: object, en: object, formula?: { ja: object, en: object } }} Spec */

/** @type {Record<string, Spec>} */
const specs = {
  "merge-rule": {
    kind: "merge-same",
    ja: {
      eye: "TRY IT / SAME TIER",
      title: "同じ段だけ合体させよう",
      body: "左右の段をそろえて「合体」。段が違うと何も起きません。ゲームも「同じ tier のときだけ次の段へ」という1行のルールです。",
      hint: "本物でも接触＋同段の両方がそろった瞬間だけ merge します。",
      controls: [
        ["data-lab-bump-left", "左を上げる"],
        ["data-lab-bump-right", "右を上げる"],
        ["data-lab-merge", "合体！", "lab-button-primary"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-tier", "段"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-mismatch="段が違うので合体しない" data-merge="合体して T{n} に！"',
    },
    en: {
      eye: "TRY IT / SAME TIER",
      title: "Merge only matching tiers",
      body: "Raise left/right until tiers match, then Merge. Different tiers do nothing—same one-line rule as the game.",
      hint: "Real code merges only on contact + equal tier.",
      controls: [
        ["data-lab-bump-left", "Raise left"],
        ["data-lab-bump-right", "Raise right"],
        ["data-lab-merge", "Merge!", "lab-button-primary"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-tier", "tiers"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-mismatch="tiers differ — no merge" data-merge="merged → T{n}"',
    },
    formula: {
      ja: {
        eye: "THE ONE-LINE RULE",
        lines: ["同じ段なら", "次の段の1個にする"],
        p: "接触していても段が違えば無視。条件を1か所にまとめると、バランス調整が楽です。",
      },
      en: {
        eye: "THE ONE-LINE RULE",
        lines: ["if tiers match", "spawn next tier"],
        p: "Contact alone is not enough—equal tiers unlock the merge. One rule keeps balancing simple.",
      },
    },
  },
  "play-a-card": {
    kind: "card-play",
    ja: {
      eye: "TRY IT / RESOLVE",
      title: "カードを使って数字を変えよう",
      body: "攻撃・防御・回復は見た目が違っても、同じ「コストを払って Kind で分岐」です。エナジーが足りないと使えません。ターン終了で敵が殴り、エナジーが戻ります。",
      hint: "ゲーム本体も switch card.Kind で同じ解決をしています。",
      controls: [
        ["data-lab-card", "攻撃 2⚡", "lab-button-primary", "damage"],
        ["data-lab-card", "防御 1⚡", "", "block"],
        ["data-lab-card", "回復 2⚡", "", "heal"],
        ["data-lab-end", "ターン終了"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-energy", "エナジー"],
        ["data-lab-hp", "自分HP"],
        ["data-lab-enemy", "敵HP"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-poor="エナジー不足" data-enemy="敵の攻撃"',
    },
    en: {
      eye: "TRY IT / RESOLVE",
      title: "Play cards to change numbers",
      body: "Attack, block, and heal look different but share one path: pay cost, branch on Kind. Not enough energy? The card won’t play. End turn: foe hits, energy refills.",
      hint: "The WASM game resolves with the same switch on card.Kind.",
      controls: [
        ["data-lab-card", "Strike 2⚡", "lab-button-primary", "damage"],
        ["data-lab-card", "Guard 1⚡", "", "block"],
        ["data-lab-card", "Heal 2⚡", "", "heal"],
        ["data-lab-end", "End turn"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-energy", "energy"],
        ["data-lab-hp", "your HP"],
        ["data-lab-enemy", "foe HP"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-poor="not enough energy" data-enemy="enemy attack"',
    },
    formula: {
      ja: {
        eye: "THE PAY → BRANCH RULE",
        lines: ["energy -= cost", "switch kind { damage / block / heal }"],
        p: "カードを増やすときはデータ1行。解決の仕組みは変えません。",
      },
      en: {
        eye: "THE PAY → BRANCH RULE",
        lines: ["energy -= cost", "switch kind { damage / block / heal }"],
        p: "New cards are new data rows—the resolver stays the same.",
      },
    },
  },
  "turn-energy": {
    kind: "energy-turn",
    ja: {
      eye: "TRY IT / REFILL",
      title: "使って、ターン終了で満タンに戻す",
      body: "1ターンの予算がエナジーです。「使う」で1つ減らし、「ターン終了」で最大まで戻します。カードが強くても、予算を超えて連打はできません。",
      hint: "最大値を変える前に、ターン終了で1だけ回復するルールを足してみましょう。",
      controls: [
        ["data-lab-spend", "1コスト使う", "lab-button-primary"],
        ["data-lab-end", "ターン終了"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-energy", "エナジー"],
        ["data-lab-turn", "ターン"],
        ["data-lab-spent", "使った分"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-max="3" data-empty="もう0" data-spend="1コスト消費" data-refill="満タンに回復"',
    },
    en: {
      eye: "TRY IT / REFILL",
      title: "Spend, then refill on end turn",
      body: "Energy is your turn budget. Spend lowers it; End turn snaps it back to max. Strong cards still can’t outrun the budget.",
      hint: "Before tuning max energy, add a rule that refills 1 energy on end turn.",
      controls: [
        ["data-lab-spend", "Spend 1", "lab-button-primary"],
        ["data-lab-end", "End turn"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-energy", "energy"],
        ["data-lab-turn", "turn"],
        ["data-lab-spent", "spent"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-max="3" data-empty="empty" data-spend="spent 1" data-refill="refilled"',
    },
    formula: {
      ja: {
        eye: "THE TURN BUDGET",
        lines: ["使うたびに energy--", "ターン終了で energy = max"],
        p: "行動の強さはカード、回数の上限はエナジー。役割を分けると調整しやすいです。",
      },
      en: {
        eye: "THE TURN BUDGET",
        lines: ["each play: energy--", "end turn: energy = max"],
        p: "Cards set power; energy caps how often. Split those jobs to tune easily.",
      },
    },
  },
  "find-matches": {
    kind: "match-scan",
    ja: {
      eye: "TRY IT / SCAN",
      title: "3つ以上そろいを光らせよう",
      body: "「スキャン」は縦横に同じ色が3つ以上続く区間を探します。光ったマスが消せる候補。「消して落下」で空いた穴へ上のマスが落ち、新しい色が補充されます。",
      hint: "交換のあとに必ずスキャン→消去→落下を同じ順で回します。",
      controls: [
        ["data-lab-scan", "スキャン", "lab-button-primary"],
        ["data-lab-clear", "消して落下"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "光ったマス"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-found="そろいを発見！" data-none="そろいなし" data-cleared="消去＋重力"',
    },
    en: {
      eye: "TRY IT / SCAN",
      title: "Highlight runs of 3+",
      body: "Scan walks rows and columns for equal runs of length ≥ 3. Lit cells are clear candidates. Clear+fall drops gems into holes and refills.",
      hint: "After every swap, run scan → clear → gravity in that order.",
      controls: [
        ["data-lab-scan", "Scan", "lab-button-primary"],
        ["data-lab-clear", "Clear + fall"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "lit cells"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-found="match found" data-none="no match" data-cleared="cleared + gravity"',
    },
    formula: {
      ja: {
        eye: "THE SCAN RULE",
        lines: ["同じ色が3つ以上続く", "そのマスを消す候補にする"],
        p: "特別な形から書かず、まず「連続カウント」だけにするとバグが減ります。",
      },
      en: {
        eye: "THE SCAN RULE",
        lines: ["run length ≥ 3", "mark those cells"],
        p: "Don’t start with fancy shapes—count runs first and bugs shrink.",
      },
    },
  },
  "clear-and-fall": {
    kind: "match-scan",
    ja: {
      eye: "TRY IT / GRAVITY",
      title: "消したあとに落とそう",
      body: "まずスキャンで候補を光らせ、次に「消して落下」。穴の上にあった色が下へ詰まります。マッチ3は「消す」と「落とす」を別ステップにすると安定します。",
      hint: "落下後にもう一度スキャンすると連鎖になります。",
      controls: [
        ["data-lab-scan", "スキャン", "lab-button-primary"],
        ["data-lab-clear", "消して落下"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "光ったマス"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-found="消せる！" data-none="まだそろっていない" data-cleared="落下した"',
    },
    en: {
      eye: "TRY IT / GRAVITY",
      title: "Clear, then collapse",
      body: "Scan lights candidates; Clear+fall packs colors downward. Keeping clear and gravity as separate steps keeps match-3 stable.",
      hint: "Scan again after a fall to unlock cascades.",
      controls: [
        ["data-lab-scan", "Scan", "lab-button-primary"],
        ["data-lab-clear", "Clear + fall"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "lit cells"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-found="ready to clear" data-none="no match yet" data-cleared="collapsed"',
    },
    formula: {
      ja: {
        eye: "CLEAR THEN FALL",
        lines: ["揃ったマスを空にする", "列ごとに下へ詰める"],
        p: "同時にやろうとすると穴の扱いで混乱します。順番を守ると連鎖も書けます。",
      },
      en: {
        eye: "CLEAR THEN FALL",
        lines: ["empty matched cells", "pack each column down"],
        p: "Doing both at once confuses holes. Keep the order and cascades become easy.",
      },
    },
  },
  "falling-cell": {
    kind: "drop-timer",
    ja: {
      eye: "TRY IT / DROP CLOCK",
      title: "時計が満ちたときだけ1マス落とす",
      body: "「1 tick進める」たびにタイマーが増えます。満タンになると1マス落下。下にブロックがあると LOCK して上から再開。tickごとに落とすと一瞬で床へ着くので、時計で間引きます。",
      hint: "本番は 38 tickで1マス。ラボはわかりやすく 8 tickで再現しています。",
      controls: [
        ["data-lab-step", "1 tick進める", "lab-button-primary"],
        ["data-lab-reset", "時計を0へ", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-timer", "タイマー"],
        ["data-lab-y", "行 y"],
        ["data-lab-note", "出来事"],
      ],
      board: true,
      data: 'data-need="8" data-drop="落下！" data-lock="固定！" data-wait="まだ待つ…" data-top="天井がふさがった"',
    },
    en: {
      eye: "TRY IT / DROP CLOCK",
      title: "Drop one cell only when the clock fills",
      body: "Each Step advances the timer by one tick. When it’s full, the cell moves one row. If blocked below, it LOCKs and respawns up top. Dropping every tick would hit the floor instantly—so we throttle with a clock.",
      hint: "The real game uses 38 ticks; this lab uses 8 so you can feel it sooner.",
      controls: [
        ["data-lab-step", "Advance 1 tick", "lab-button-primary"],
        ["data-lab-reset", "Reset clock", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-timer", "timer"],
        ["data-lab-y", "row y"],
        ["data-lab-note", "event"],
      ],
      board: true,
      data: 'data-need="8" data-drop="DROP" data-lock="LOCK" data-wait="wait…" data-top="TOP BLOCKED"',
    },
    formula: {
      ja: {
        eye: "THE TIMER RULE",
        lines: ["timer++", "timer が満タンなら y++（または固定）"],
        p: "速さを変えたいときは dropFrames を変えるだけ。位置の式は触りません。",
      },
      en: {
        eye: "THE TIMER RULE",
        lines: ["timer++", "when full: y++ (or lock)"],
        p: "Want faster drops? Change dropFrames only—leave the position math alone.",
      },
    },
  },
  "lock-lines": {
    kind: "line-clear",
    ja: {
      eye: "TRY IT / FULL ROW",
      title: "埋まった行だけ消そう",
      body: "黄色い行はマスが全部埋まっています。「行を消す」とその行だけ消えて、上の行が落ちてきます。テトリス型の得点は、この「埋まっているか」判定から始まります。",
      hint: "複数行が同時に埋まっていれば、まとめて消えます。",
      controls: [
        ["data-lab-clear", "行を消す", "lab-button-primary"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-cleared", "消した行"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-cleared="消した！" data-none="埋まっていない"',
    },
    en: {
      eye: "TRY IT / FULL ROW",
      title: "Clear only full rows",
      body: "Gold rows are completely filled. Clear Rows removes them and drops everything above. Scoring in falling-block games starts from this “is the row full?” test.",
      hint: "Several full rows clear together in one pass.",
      controls: [
        ["data-lab-clear", "Clear rows", "lab-button-primary"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-cleared", "cleared"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-cleared="cleared!" data-none="no full rows"',
    },
    formula: {
      ja: {
        eye: "THE FULL-ROW RULE",
        lines: ["行の全マスが true なら消す", "残った行を下へ詰める"],
        p: "消す行の数を得点表に渡すと、シングル〜テトリスの点が作れます。",
      },
      en: {
        eye: "THE FULL-ROW RULE",
        lines: ["if every cell is filled → clear", "pack remaining rows down"],
        p: "Feed the cleared count into a score table for singles through Tetrises.",
      },
    },
  },
  "tile-runner": {
    kind: "tile",
    ja: {
      eye: "TRY IT / LEGAL TILE",
      title: "壁のマスには入れない",
      body: "矢印は行と列を整数で動かします。赤いマスは壁。入れなければ座標は変わりません。本番は通った2つの中心の間だけを滑らかに描きます。",
      hint: "中心に着いてから次の入力を受け付けると、通路からずれません。",
      controls: [
        ["data-lab-up", "↑"],
        ["data-lab-left", "←"],
        ["data-lab-right", "→"],
        ["data-lab-down", "↓"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-pos", "タイル x,y"],
        ["data-lab-face", "向き"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-blocked="その方向は壁"',
    },
    en: {
      eye: "TRY IT / LEGAL TILE",
      title: "Walls reject the step",
      body: "Arrows move row/col by integers. Red cells are walls—illegal moves leave you in place. The real game only lerps between two centers.",
      hint: "Accept a new heading only after arriving at a center to stay in the corridor.",
      controls: [
        ["data-lab-up", "↑"],
        ["data-lab-left", "←"],
        ["data-lab-right", "→"],
        ["data-lab-down", "↓"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-pos", "tile x,y"],
        ["data-lab-face", "facing"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-blocked="wall that way"',
    },
  },
  cascade: {
    kind: "match-scan",
    ja: {
      eye: "TRY IT / CASCADE",
      title: "消す→落とす→もう一度探す",
      body: "スキャンで光らせ、消して落下、もう一度スキャン。落ちたあとに新しいそろいができるのが連鎖です。1回で終わらせないのがポイントです。",
      hint: "本番も while で「そろいがなくなるまで」くり返します。",
      controls: [
        ["data-lab-scan", "スキャン", "lab-button-primary"],
        ["data-lab-clear", "消して落下"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "光ったマス"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-found="連鎖の芽！" data-none="そろいなし（連鎖終了）" data-cleared="落下 — もう一度スキャン"',
    },
    en: {
      eye: "TRY IT / CASCADE",
      title: "Clear → fall → scan again",
      body: "Scan, clear+fall, scan again. New matches after a collapse are the cascade. Don’t stop after one clear.",
      hint: "The real game loops while matches remain.",
      controls: [
        ["data-lab-scan", "Scan", "lab-button-primary"],
        ["data-lab-clear", "Clear + fall"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "lit cells"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-found="cascade seed!" data-none="no match (cascade over)" data-cleared="fell — scan again"',
    },
    formula: {
      ja: {
        eye: "THE CASCADE LOOP",
        lines: ["scan → clear → gravity", "そろいがある間くり返す"],
        p: "得点は連鎖の深さで倍増させると気持ちよくなります。",
      },
      en: {
        eye: "THE CASCADE LOOP",
        lines: ["scan → clear → gravity", "repeat while matches exist"],
        p: "Score multipliers on cascade depth make the loop feel juicy.",
      },
    },
  },
  "special-pieces": {
    kind: "match-scan",
    ja: {
      eye: "TRY IT / MARK SPECIALS",
      title: "消すマスの集合を先に作る",
      body: "特殊ピースも、まず「どのマスを消すか」の集合を作ります。スキャンで集合を見てから消去。爆弾もラインも、最後は同じクリア処理に渡せます。",
      hint: "効果ごとに描画を分けても、消去パイプは1本にまとめると安全です。",
      controls: [
        ["data-lab-scan", "効果範囲を集める", "lab-button-primary"],
        ["data-lab-clear", "まとめて消す"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "対象マス"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-found="効果範囲をマーク" data-none="対象なし" data-cleared="まとめて消去"',
    },
    en: {
      eye: "TRY IT / MARK SPECIALS",
      title: "Build the clear-set first",
      body: "Specials still start as a set of cells to clear. Scan marks the set; then one clear pass handles bombs and lines alike.",
      hint: "Different FX, one clear pipeline keeps edge cases tame.",
      controls: [
        ["data-lab-scan", "Collect affected cells", "lab-button-primary"],
        ["data-lab-clear", "Clear the set"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "targets"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-found="marked effect cells" data-none="none" data-cleared="batch cleared"',
    },
    formula: {
      ja: {
        eye: "SET THEN CLEAR",
        lines: ["消すマスを Set に入れる", "Set をまとめて空にする"],
        p: "特殊効果は「Set をどう埋めるか」だけ変えれば足せます。",
      },
      en: {
        eye: "SET THEN CLEAR",
        lines: ["add cells to a Set", "clear the whole Set"],
        p: "New specials only change how the Set is filled.",
      },
    },
  },
  "deck-cycle": {
    kind: "energy-turn",
    ja: {
      eye: "TRY IT / DRAW CYCLE",
      title: "手札コストを回して山を循環させる",
      body: "カードを使う（コスト消費）→ ターン終了でリソース回復。山札・捨て札の循環も、同じ「ターンの区切り」でまとめて処理すると迷子になりません。",
      hint: "山が空なら捨て札をシャッフルして山に戻す、もターン区切りで行います。",
      controls: [
        ["data-lab-spend", "1枚使う", "lab-button-primary"],
        ["data-lab-end", "ターン終了（山を循環）"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-energy", "手札枠"],
        ["data-lab-turn", "ターン"],
        ["data-lab-spent", "使った枚数"],
        ["data-lab-note", "結果"],
      ],
      board: true,
      data: 'data-max="5" data-empty="手札がない" data-spend="1枚プレイ" data-refill="ドローして手札回復"',
    },
    en: {
      eye: "TRY IT / DRAW CYCLE",
      title: "Spend the hand, refill on the turn edge",
      body: "Play cards (spend) → end turn refills. Deck/discard cycling belongs on that same turn boundary so state stays easy to find.",
      hint: "Empty draw pile? Shuffle discard into draw—also on the turn edge.",
      controls: [
        ["data-lab-spend", "Play 1", "lab-button-primary"],
        ["data-lab-end", "End turn (cycle)"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-energy", "hand slots"],
        ["data-lab-turn", "turn"],
        ["data-lab-spent", "played"],
        ["data-lab-note", "result"],
      ],
      board: true,
      data: 'data-max="5" data-empty="hand empty" data-spend="played 1" data-refill="drew back to full"',
    },
    formula: {
      ja: {
        eye: "THE CYCLE EDGE",
        lines: ["play: hand → discard", "turn end: draw / reshuffle"],
        p: "循環の処理をターンの境目に集めると、戦闘中のバグが減ります。",
      },
      en: {
        eye: "THE CYCLE EDGE",
        lines: ["play: hand → discard", "turn end: draw / reshuffle"],
        p: "Keep cycling on the turn edge and mid-combat bugs drop.",
      },
    },
  },
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
    .map(
      ([attr, label]) =>
        `<div><span>${label}</span><strong ${attr}>${attr.includes("note") || attr.includes("list") ? "—" : attr.includes("tier") ? "1 + 1" : "0"}</strong></div>`,
    )
    .join("");
  const board = c.board
    ? `<div class="lab-board" data-lab-board role="img" aria-label="lab"></div>`
    : "";
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
  const lab = labHTML(lang, slug, spec);
  let next = html.replace(/<div class="motion-lab"[\s\S]*?<\/div>\s*<\/div>\s*<\/div>/, lab);
  // Simpler: replace from motion-lab to the closing before code-lesson / formula
  const start = html.search(/<div class="motion-lab"/);
  if (start < 0) {
    console.warn("no lab", slug, lang);
    return html;
  }
  // Find end: after lab-visual's closing — look for next sibling formula/code-lesson/why-grid
  const after = html.slice(start);
  const endRel = after.search(/<div class="(formula|code-lesson|why-grid)"/);
  if (endRel < 0) {
    console.warn("no end", slug, lang);
    return html;
  }
  const end = start + endRel;
  let insert = lab;
  const hasFormula = /class="formula"/.test(html);
  if (!hasFormula && spec.formula) {
    insert += formulaHTML(lang, spec.formula);
  }
  next = html.slice(0, start) + insert + html.slice(end);
  return next;
}

let updated = 0;
for (const [slug, spec] of Object.entries(specs)) {
  for (const lang of ["ja", "en"]) {
    // Find route under tracks
    const candidates = [
      `tracks/merge-physics/${slug}`,
      `tracks/deckbuilder/${slug}`,
      `tracks/match3/${slug}`,
      `tracks/falling-blocks/${slug}`,
      `tracks/maze-chase/${slug}`,
    ];
    let route = null;
    for (const c of candidates) {
      if (existsSync(join(root, "web", lang, c, "index.html"))) {
        route = c;
        break;
      }
    }
    if (!route) {
      console.warn("missing", lang, slug);
      continue;
    }
    const path = join(root, "web", lang, route, "index.html");
    const html = readFileSync(path, "utf8");
    const next = patch(html, lang, slug, spec);
    if (next !== html) {
      writeFileSync(path, next);
      updated++;
      console.log("updated", lang, route);
    } else {
      console.log("unchanged", lang, route);
    }
  }
}
console.log(`Done. Updated ${updated} page(s).`);
