#!/usr/bin/env node
/**
 * Polish thin/broken labs on maze-chase + bomb-maze (Codex's final stretch).
 * Replaces non-interactive turn/entities stubs with concept labs.
 */
import { readFileSync, writeFileSync, existsSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;

const specs = {
  "buffered-turn": {
    kind: "input-buffer",
    route: "tracks/maze-chase/buffered-turn",
    ja: {
      eye: "TRY IT / CURRENT vs QUEUED",
      title: "押した向きを予約して、中心で反映",
      body: "方向ボタンは queued だけを変えます。「移動中」にしてから「中心到着」で、予約が current に入ります。壁判定はこの瞬間だけです。",
      hint: "すぐに current を書き換えると、マスの途中で壁にめり込みやすくなります。",
      controls: [
        ["data-lab-dir", "↑", "", "N"],
        ["data-lab-dir", "←", "", "W"],
        ["data-lab-dir", "→", "lab-button-primary", "E"],
        ["data-lab-dir", "↓", "", "S"],
        ["data-lab-move", "マス間を進む"],
        ["data-lab-center", "中心に到着"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-current", "current"],
        ["data-lab-queued", "queued"],
        ["data-lab-note", "結果"],
      ],
      data: 'data-queued="予約した" data-turned="中心で曲がった" data-same="直進" data-moving="移動中"',
      formula: {
        eye: "THE BUFFER RULE",
        lines: ["input → queued", "at center: current = queued if passable"],
        p: "入力の瞬間と向きの確定を分けます。",
      },
    },
    en: {
      eye: "TRY IT / CURRENT vs QUEUED",
      title: "Queue a facing, commit at center",
      body: "Direction buttons only write queued. Move between tiles, then Arrive at center to copy queued into current if the next tile is open.",
      hint: "Writing current immediately makes mid-tile wall clips common.",
      controls: [
        ["data-lab-dir", "↑", "", "N"],
        ["data-lab-dir", "←", "", "W"],
        ["data-lab-dir", "→", "lab-button-primary", "E"],
        ["data-lab-dir", "↓", "", "S"],
        ["data-lab-move", "Move between"],
        ["data-lab-center", "Arrive center"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-current", "current"],
        ["data-lab-queued", "queued"],
        ["data-lab-note", "result"],
      ],
      data: 'data-queued="queued" data-turned="turned at center" data-same="kept going" data-moving="moving"',
      formula: {
        eye: "THE BUFFER RULE",
        lines: ["input → queued", "at center: current = queued if passable"],
        p: "Split the press moment from the facing commit.",
      },
    },
  },
  "patrol-chase": {
    kind: "junction-pick",
    route: "tracks/maze-chase/patrol-chase",
    ja: {
      eye: "TRY IT / AI MODE",
      title: "追跡と散開で交点の選び方を変える",
      body: "交点では「目標」に近い辺を選びます。chase はプレイヤー方向、scatter は隅。モードが経路の意味を切り替えます。",
      hint: "同じ交点ロジックに、目標座標だけ差し込みます。",
      controls: [
        ["data-lab-chase", "追跡モード", "lab-button-primary"],
        ["data-lab-scatter", "散開モード"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-choice", "選択"],
        ["data-lab-note", "結果"],
      ],
      data: 'data-chase="プレイヤー方向へ" data-scatter="隅へ"',
      formula: {
        eye: "MODE → TARGET",
        lines: ["chase: target = player", "scatter: target = corner"],
        p: "交点AIは目標を変えるだけで振る舞いが変わります。",
      },
    },
    en: {
      eye: "TRY IT / AI MODE",
      title: "Chase vs scatter at a junction",
      body: "At a junction, pick the edge closer to a target. Chase aims at the player; scatter aims at a corner. Mode swaps the target.",
      hint: "Same junction picker—only the target changes.",
      controls: [
        ["data-lab-chase", "Chase", "lab-button-primary"],
        ["data-lab-scatter", "Scatter"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-choice", "choice"],
        ["data-lab-note", "result"],
      ],
      data: 'data-chase="toward player" data-scatter="toward corner"',
      formula: {
        eye: "MODE → TARGET",
        lines: ["chase: target = player", "scatter: target = corner"],
        p: "Junction AI changes behavior by swapping targets only.",
      },
    },
  },
  "junction-ai": {
    kind: "junction-pick",
    route: "tracks/maze-chase/junction-ai",
    ja: {
      eye: "TRY IT / JUNCTION",
      title: "交差点で1本だけ選ぶ",
      body: "通路の交点では、入れる方向のうち目標に一番近いものを1つ選びます。Uターン禁止などのルールもここで足せます。",
      hint: "候補を列挙 → スコア付け → 最大を採用、が基本形です。",
      controls: [
        ["data-lab-chase", "目標を右に", "lab-button-primary"],
        ["data-lab-scatter", "目標を上に"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-choice", "採用"],
        ["data-lab-note", "結果"],
      ],
      data: 'data-chase="右を選んだ" data-scatter="上を選んだ"',
      formula: {
        eye: "SCORE EDGES",
        lines: ["candidates = open neighbors", "pick min distance to target"],
        p: "ランダムより、距離スコアの方が読みやすい敵になります。",
      },
    },
    en: {
      eye: "TRY IT / JUNCTION",
      title: "Pick exactly one edge",
      body: "At a junction, list open neighbors and keep the one closest to the target. Anti-U-turn rules plug in here too.",
      hint: "Enumerate → score → take the best.",
      controls: [
        ["data-lab-chase", "Target right", "lab-button-primary"],
        ["data-lab-scatter", "Target up"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-choice", "pick"],
        ["data-lab-note", "result"],
      ],
      data: 'data-chase="picked right" data-scatter="picked up"',
      formula: {
        eye: "SCORE EDGES",
        lines: ["candidates = open neighbors", "pick min distance to target"],
        p: "Distance scoring reads clearer than pure randomness.",
      },
    },
  },
  "timed-bomb": {
    kind: "bomb-timer",
    route: "tracks/bomb-maze/timed-bomb",
    ja: {
      eye: "TRY IT / FUSE",
      title: "設置→カウント→爆発→削除",
      body: "「設置」で爆弾を出し、「1 tick進める」でタイマーを減らします。0で BLAST、もう一度で削除。状態機械の基本です。",
      hint: "本物は約90フレーム点火。ラボは8で体感しやすくしています。",
      controls: [
        ["data-lab-place", "爆弾を置く", "lab-button-primary"],
        ["data-lab-tick", "1ステップ進める"],
        ["data-lab-reset", "消す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-timer", "残り"],
        ["data-lab-state", "状態"],
        ["data-lab-note", "結果"],
      ],
      data: 'data-fuse="8" data-placed="設置" data-tick="カウント" data-boom="爆発！" data-gone="削除" data-none="先に設置"',
      formula: {
        eye: "THE FUSE STATE MACHINE",
        lines: ["timer-- each tick", "timer<=0 → blasting → remove"],
        p: "見た目の丸ではなく、座標・timer・状態を1つの構造体にします。",
      },
    },
    en: {
      eye: "TRY IT / FUSE",
      title: "Place → tick → blast → remove",
      body: "Place arms a bomb; Step drains the fuse. At 0 it BLASTs; another step removes it. Classic state machine.",
      hint: "Real fuse is ~90 frames; this lab uses 8 so you can feel it.",
      controls: [
        ["data-lab-place", "Place bomb", "lab-button-primary"],
        ["data-lab-tick", "Advance step"],
        ["data-lab-reset", "Clear", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-timer", "left"],
        ["data-lab-state", "state"],
        ["data-lab-note", "result"],
      ],
      data: 'data-fuse="8" data-placed="placed" data-tick="tick" data-boom="BOOM" data-gone="removed" data-none="place first"',
      formula: {
        eye: "THE FUSE STATE MACHINE",
        lines: ["timer-- each tick", "timer<=0 → blasting → remove"],
        p: "Not a drawn circle—position, timer, and state in one struct.",
      },
    },
  },
  "cross-blast": {
    kind: "cross-blast",
    route: "tracks/bomb-maze/cross-blast",
    ja: {
      eye: "TRY IT / CROSS ARMS",
      title: "火力ぶんだけ十字に伸ばす",
      body: "中心から上下左右へ、壁に当たるまでマスを進めます。「火力+」で届く距離が伸びます。固い壁で止まるのがポイントです。",
      hint: "各方向は独立した for ループです。",
      controls: [
        ["data-lab-power", "火力 +1", "lab-button-primary"],
        ["data-lab-reset", "火力1へ", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-reach", "火力"],
        ["data-lab-note", "結果"],
      ],
      data: 'data-powered="火力アップ"',
      formula: {
        eye: "FOUR LOOPS",
        lines: ["for each dir in NESW", "walk until wall or power"],
        p: "円形ではなく、グリッドの十字が爆弾の定石です。",
      },
    },
    en: {
      eye: "TRY IT / CROSS ARMS",
      title: "Extend a cross up to power",
      body: "From the center, walk N/E/S/W until a wall or the power runs out. Power+ lengthens each arm. Hard walls stop the flame.",
      hint: "Each direction is its own for-loop.",
      controls: [
        ["data-lab-power", "Power +1", "lab-button-primary"],
        ["data-lab-reset", "Power 1", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-reach", "power"],
        ["data-lab-note", "result"],
      ],
      data: 'data-powered="powered up"',
      formula: {
        eye: "FOUR LOOPS",
        lines: ["for each dir in NESW", "walk until wall or power"],
        p: "Grid crosses beat circles for tile bombs.",
      },
    },
  },
  "chain-explosion": {
    kind: "chain-bomb",
    route: "tracks/bomb-maze/chain-explosion",
    ja: {
      eye: "TRY IT / CHAIN",
      title: "爆風が隣の爆弾に燃え移る",
      body: "「点火」で最初の爆弾が BLAST になり、隣の爆弾へ連鎖します。「消去」で炎を片付けます。キューや再帰で同じことができます。",
      hint: "同じフレームで連鎖を最後まで解決するか、次フレームに延ばすかを選べます。",
      controls: [
        ["data-lab-ignite", "点火", "lab-button-primary"],
        ["data-lab-clear", "炎を消す"],
        ["data-lab-reset", "配置を戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "炎の数"],
        ["data-lab-note", "結果"],
      ],
      data: 'data-chained="連鎖！" data-cleared="消去" data-done="もうない"',
      formula: {
        eye: "BLAST TOUCHES BOMB",
        lines: ["if flame hits bomb", "arm that bomb too"],
        p: "爆発中フラグを見ると、二重起爆を防げます。",
      },
    },
    en: {
      eye: "TRY IT / CHAIN",
      title: "Let blasts ignite neighbor bombs",
      body: "Ignite turns the first bomb into BLAST and spreads to neighbors. Clear removes flames. Queue or recursion both work.",
      hint: "Resolve the whole chain in one tick, or stretch it across ticks.",
      controls: [
        ["data-lab-ignite", "Ignite", "lab-button-primary"],
        ["data-lab-clear", "Clear flames"],
        ["data-lab-reset", "Reset layout", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-count", "flames"],
        ["data-lab-note", "result"],
      ],
      data: 'data-chained="chained!" data-cleared="cleared" data-done="done"',
      formula: {
        eye: "BLAST TOUCHES BOMB",
        lines: ["if flame hits bomb", "arm that bomb too"],
        p: "A blasting flag prevents double-trigger bugs.",
      },
    },
  },
  "escape-ai": {
    kind: "escape-timing",
    route: "tracks/bomb-maze/escape-ai",
    ja: {
      eye: "TRY IT / TIME MAP",
      title: "到着時刻と導火線を比べる",
      body: "AIは「今どのマスが空いているか」だけでなく、「到着までに導火線が残るか」を見ます。ETA+余裕 < 導火線なら通れる経路です。",
      hint: "時間が進むと安全だった道が危険になります。だからtickごとに張り直します。",
      controls: [
        ["data-lab-tick", "時間を進める", "lab-button-primary"],
        ["data-lab-far", "遠回りにする"],
        ["data-lab-reset", "戻す", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-fuse", "導火線"],
        ["data-lab-eta", "到着ETA"],
        ["data-lab-note", "結果"],
      ],
      data: 'data-ok="まだ安全" data-bad="危険→張り直し" data-far="経路が長くなった"',
      formula: {
        eye: "ETA vs FUSE",
        lines: ["eta = steps * framesPerStep", "safe if eta + margin < fuse"],
        p: "空間の最短路に、時間の制約を足したのが回避AIです。",
      },
    },
    en: {
      eye: "TRY IT / TIME MAP",
      title: "Compare arrival time to the fuse",
      body: "AI checks not only open cells, but whether the fuse outlasts the arrival. ETA+margin < fuse means the route is allowed.",
      hint: "Safe paths go deadly as time passes—replan every tick.",
      controls: [
        ["data-lab-tick", "Advance time", "lab-button-primary"],
        ["data-lab-far", "Take a longer path"],
        ["data-lab-reset", "Reset", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-fuse", "fuse"],
        ["data-lab-eta", "ETA"],
        ["data-lab-note", "result"],
      ],
      data: 'data-ok="still safe" data-bad="danger — reroute" data-far="path got longer"',
      formula: {
        eye: "ETA vs FUSE",
        lines: ["eta = steps * framesPerStep", "safe if eta + margin < fuse"],
        p: "Escape AI is shortest path plus a time constraint.",
      },
    },
  },
  "ebi-bomber": {
    kind: "pipeline",
    route: "tracks/bomb-maze/ebi-bomber",
    ja: {
      eye: "TRY IT / FRAME ORDER",
      title: "爆弾ゲームの1 tickを並べる",
      body: "総合はパイプです。入力→移動→爆弾タイマー→爆風→壁/連鎖→敵AI→描画。順番を固定すると「爆風中に歩けるか」が説明できます。",
      hint: "描画は最後。AIは爆風更新のあとに置くことが多いです。",
      controls: [
        ["data-lab-next", "次の工程", "lab-button-primary"],
        ["data-lab-reset", "先頭へ", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-stepn", "工程"],
        ["data-lab-note", "今"],
      ],
      data: 'data-steps="input,move,fuse,blast,walls,ai,draw" data-loop="次フレーム"',
      formula: {
        eye: "FIXED SYSTEMS ORDER",
        lines: ["update bombs before AI", "draw last"],
        p: "順序を変えると、同じ入力でも生死が変わることがあります。",
      },
    },
    en: {
      eye: "TRY IT / FRAME ORDER",
      title: "Order one bomber frame",
      body: "The full game is a pipe: input → move → fuse → blast → walls/chain → AI → draw. A fixed order makes blast-vs-walk rules explainable.",
      hint: "Draw last. AI usually runs after blast updates.",
      controls: [
        ["data-lab-next", "Next job", "lab-button-primary"],
        ["data-lab-reset", "Frame start", "lab-button-quiet"],
      ],
      values: [
        ["data-lab-stepn", "step"],
        ["data-lab-note", "now"],
      ],
      data: 'data-steps="input,move,fuse,blast,walls,ai,draw" data-loop="next tick"',
      formula: {
        eye: "FIXED SYSTEMS ORDER",
        lines: ["update bombs before AI", "draw last"],
        p: "Reordering systems can change who survives the same input.",
      },
    },
  },
};

function btn([attr, label, cls = "", value = ""]) {
  const klass = ` class="lab-button${cls ? ` ${cls}` : ""}"`;
  if (attr === "data-lab-dir" || attr === "data-lab-card") {
    return `<button type="button"${klass} ${attr}="${value}">${label}</button>`;
  }
  return `<button type="button"${klass} ${attr}>${label}</button>`;
}

function labHTML(lang, slug, spec) {
  const c = spec[lang];
  const controls = c.controls.map(btn).join("");
  const values = c.values
    .map(([attr, label]) => {
      let initial = "—";
      if (attr.includes("timer") || attr.includes("reach") || attr.includes("count") || attr.includes("fuse") || attr.includes("eta") || attr.includes("stepn")) initial = "0";
      if (attr.includes("current")) initial = "E";
      if (attr.includes("queued")) initial = "E";
      if (attr.includes("state")) initial = "none";
      if (attr.includes("choice")) initial = "chase:right";
      return `<div><span>${label}</span><strong ${attr}>${initial}</strong></div>`;
    })
    .join("");
  const id = `lab-title-${slug}`;
  const formula = c.formula
    ? `<div class="formula"><p class="eyebrow">${c.formula.eye}</p><div class="formula-lines"><code>${c.formula.lines[0]}</code><span>${lang === "ja" ? "つぎに" : "then"}</span><code>${c.formula.lines[1]}</code></div><p>${c.formula.p}</p></div>`
    : "";
  return `<section class="motion-lab" data-lab="${spec.kind}" ${c.data || ""} aria-labelledby="${id}">
<div class="lab-copy">
<p class="eyebrow">${c.eye}</p>
<h3 id="${id}">${c.title}</h3>
<p>${c.body}</p>
<div class="lab-controls">${controls}</div>
<p class="lab-hint">${c.hint}</p>
</div>
<div class="lab-visual">
<div class="lab-board" data-lab-board role="img" aria-label="lab"></div>
<div class="lab-values" aria-live="polite">${values}</div>
</div>
</section>${formula}`;
}

function patch(html, lang, slug, spec) {
  // Match section or div motion-lab through next code-lesson / why-grid / formula sibling
  const re = /<(section|div) class="motion-lab"[\s\S]*?<\/\1>\s*/;
  if (!re.test(html)) return null;
  let next = html.replace(re, labHTML(lang, slug, spec));
  // If a formula was already present right after, we may have duplicated — remove old adjacent formula only if we injected one
  if (spec[lang].formula) {
    next = next.replace(/(<\/section><div class="formula">[\s\S]*?<\/div>)\s*<div class="formula">[\s\S]*?<\/div>/, "$1");
  }
  return next;
}

let updated = 0;
for (const [slug, spec] of Object.entries(specs)) {
  for (const lang of ["ja", "en"]) {
    const path = join(root, "web", lang, spec.route, "index.html");
    if (!existsSync(path)) {
      console.warn("missing", path);
      continue;
    }
    const html = readFileSync(path, "utf8");
    const next = patch(html, lang, slug, spec);
    if (!next) {
      console.warn("patch failed", lang, slug);
      continue;
    }
    writeFileSync(path, next);
    updated++;
    console.log("polished", lang, slug, "→", spec.kind);
  }
}

const leftovers = [];
for (const lang of ["ja", "en"]) {
  for (const [slug, spec] of Object.entries(specs)) {
    const html = readFileSync(join(root, "web", lang, spec.route, "index.html"), "utf8");
    if (html.includes('data-lab="entities"') || html.includes("lab-action") || html.includes("entity-board")) {
      leftovers.push(`${lang}:${slug}`);
    }
  }
}
console.log(`Updated ${updated}. Stub leftovers: ${leftovers.length ? leftovers.join(", ") : "NONE"}`);
