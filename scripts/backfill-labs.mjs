#!/usr/bin/env node
/**
 * Inject interactive motion-lab blocks into lesson pages that lack them.
 * Usage: node scripts/backfill-labs.mjs
 */
import { readFileSync, writeFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import { curriculum } from "./curriculum.mjs";

const root = new URL("..", import.meta.url).pathname;

/** @type {Record<string, { kind: string, ja: object, en: object }>} */
const labs = {
  pong: {
    kind: "bounce",
    ja: {
      eye: "TRY IT / BOUNCE",
      title: "壁とパドルではね返してみよう",
      body: "「1フレーム進める」でボールが動きます。左右の壁では vx、パドルでは vy の符号が反転します。",
      step: "1フレーム進める",
      reset: "もとに戻す",
      hint: "速さの大きさはそのまま、向きだけ変わります。",
      values: [
        ["vx", "左右の速さ"],
        ["vy", "上下の速さ"],
        ["出来事", "いま起きた反射"],
      ],
      attrs: 'data-wall="壁" data-paddle="パドル"',
      board: '<div class="lab-ball" data-lab-ball aria-hidden="true"></div>',
      outs: [
        ["data-lab-vx", "6.0"],
        ["data-lab-vy", "4.0"],
        ["data-lab-note", "—"],
      ],
    },
    en: {
      eye: "TRY IT / BOUNCE",
      title: "Bounce off walls and paddles",
      body: "“Next frame” moves the ball. Side walls flip vx; the paddle flips vy.",
      step: "Next frame",
      reset: "Reset",
      hint: "Keep the speed magnitude; only the sign changes.",
      values: [
        ["vx", "horizontal"],
        ["vy", "vertical"],
        ["event", "last bounce"],
      ],
      attrs: 'data-wall="wall" data-paddle="paddle"',
      board: '<div class="lab-ball" data-lab-ball aria-hidden="true"></div>',
      outs: [
        ["data-lab-vx", "6.0"],
        ["data-lab-vy", "4.0"],
        ["data-lab-note", "—"],
      ],
    },
  },
  breakout: {
    kind: "bricks",
    ja: {
      eye: "TRY IT / MANY OBJECTS",
      title: "ブロックをタップして壊そう",
      body: "12個のブロックは同じ配列の要素です。タップで alive を false にし、スコアを足します。",
      reset: "もとに戻す",
      hint: "残りの数と得点が右に出ます。",
      values: [
        ["残ブロック", "alive の数"],
        ["SCORE", "壊すたび +10"],
        ["ルール", "配列をループ"],
      ],
      attrs: "",
      board: '<div class="lab-brick-grid" data-lab-grid></div>',
      outs: [
        ["data-lab-alive", "12"],
        ["data-lab-score", "0"],
        ["", "for i := range"],
      ],
      buttons: [["data-lab-reset", "もとに戻す", "quiet"]],
    },
    en: {
      eye: "TRY IT / MANY OBJECTS",
      title: "Tap bricks to break them",
      body: "Twelve bricks are entries in one array. Tap sets alive to false and adds score.",
      reset: "Reset",
      hint: "Remaining count and score update on the right.",
      values: [
        ["alive", "still standing"],
        ["SCORE", "+10 each"],
        ["rule", "loop the array"],
      ],
      attrs: "",
      board: '<div class="lab-brick-grid" data-lab-grid></div>',
      outs: [
        ["data-lab-alive", "12"],
        ["data-lab-score", "0"],
        ["", "for i := range"],
      ],
      buttons: [["data-lab-reset", "Reset", "quiet"]],
    },
  },
  snake: {
    kind: "snake",
    ja: {
      eye: "TRY IT / GRID HISTORY",
      title: "頭を足して、しっぽを消そう",
      body: "「進む」は新しい頭を先頭に足して末尾を削ります。「食べる」は末尾を残すので体が伸びます。",
      step: "進む",
      eat: "食べる",
      reset: "もとに戻す",
      hint: "体はマス座標のリストです。",
      values: [
        ["長さ", "body の要素数"],
        ["履歴", "頭 → しっぽ"],
        ["ルール", "append + trim"],
      ],
      attrs: 'data-ate="食べた" data-move="移動"',
      board: '<div class="lab-mono-board" data-lab-body>(2,2) (1,2) (0,2)</div>',
      outs: [
        ["data-lab-len", "3"],
        ["data-lab-body", "(2,2) (1,2) (0,2)"],
        ["", "[]point"],
      ],
      buttons: [
        ["data-lab-step", "進む", "primary"],
        ["data-lab-eat", "食べる", ""],
        ["data-lab-reset", "もとに戻す", "quiet"],
      ],
    },
    en: {
      eye: "TRY IT / GRID HISTORY",
      title: "Add a head, drop a tail",
      body: "“Move” prepends a head and trims the tail. “Eat” keeps the tail so the body grows.",
      step: "Move",
      eat: "Eat",
      reset: "Reset",
      hint: "The body is a list of grid points.",
      values: [
        ["length", "body size"],
        ["history", "head → tail"],
        ["rule", "append + trim"],
      ],
      attrs: "",
      board: '<div class="lab-mono-board" data-lab-body>(2,2) (1,2) (0,2)</div>',
      outs: [
        ["data-lab-len", "3"],
        ["data-lab-body", "(2,2) (1,2) (0,2)"],
        ["", "[]point"],
      ],
      buttons: [
        ["data-lab-step", "Move", "primary"],
        ["data-lab-eat", "Eat", ""],
        ["data-lab-reset", "Reset", "quiet"],
      ],
    },
  },
};

// Extended compact definitions for remaining slugs
const compact = {
  "space-shooter": ["entities", "弾を撃って進めて消そう", "Fire, step, and cull bullets", "fire"],
  sokoban: ["push", "箱を1マス押してみよう", "Push the crate one tile", "push"],
  platformer: ["camera", "カメラがプレイヤーを追う", "Camera eases toward the player", "camera"],
  dungeon: ["ai", "距離でパトロールと追跡を切替", "Switch patrol/chase by distance", "ai"],
  "bullet-hell": ["burst", "弾の数で円形弾幕を変える", "Change the circular burst count", "burst"],
  "moving-platforms": ["carry", "足場の差分をプレイヤーへ足す", "Add platform delta to the rider", "carry"],
  "patrol-enemies": ["stomp", "落下中だけ踏みつけ", "Stomp only while falling", "stomp"],
  "scrolling-stage": ["camera", "ワールド座標から画面座標へ", "World to screen via camera", "camera"],
  "powerup-adventure": ["power", "パワー状態で衝突結果が変わる", "Power flag changes collision result", "power"],
  "arena-dodge": ["move8", "斜め移動を正規化しよう", "Normalize diagonal movement", "move8"],
  "auto-turret": ["aim", "一番近い敵へクールダウン発射", "Cooldown nearest foe on cooldown", "aim"],
  swarm: ["pool", "倒した席にすぐ次の敵を上書き", "Overwrite a slot on kill", "pool"],
  "experience-draft": ["draft", "レベルアップで選択モードへ", "Enter draft mode on level-up", "draft"],
  "weapon-evolution": ["evolve", "武器データを差し替えて進化", "Evolve by swapping weapon data", "evolve"],
  "survival-run": ["curve", "秒数が間隔と速度を変える", "Seconds reshape interval and speed", "curve"],
  "tap-counter": ["click", "押すたびに数を足そう", "Add one on each press", "click"],
  "first-shop": ["shop", "足りるときだけ買う", "Buy only when you can afford it", "shop"],
  "idle-factory": ["idle", "時間差分で自動生産", "Produce with delta time", "idle"],
  "offline-bakery": ["save", "不在秒数×生産でオフライン報酬", "Away seconds × rate = offline gain", "save"],
  "village-walk": ["tile", "壁マスには入れない", "Walls block tile steps", "tile"],
  "dialogue-flags": ["flag", "フラグでセリフが変わる", "Flags change dialogue lines", "flag"],
  "command-battle": ["turn", "戦闘状態を1つずつ進める", "Step the battle state machine", "turn"],
  "stats-status": ["damage", "攻撃−防御（強化込み）", "Attack − defense with buffs", "damage"],
  "inventory-shop": ["inv", "お金で装備を買う", "Spend gold to equip an item", "inv"],
  "world-encounters": ["scene", "フィールドと戦闘を切り替える", "Switch field and battle scenes", "scene"],
  "ebi-quest": ["quest", "クエスト番号を保存・読込", "Save and load a quest id", "quest"],
  "frame-attack": ["frames", "発生・持続・硬直をコマ送り", "Step startup / active / recovery", "frames"],
  "hit-reaction": ["react", "ヒットストップとのけぞり", "Hitstop then stun decay", "react"],
  "guard-throw": ["rps", "打撃・ガード・投げの読み合い", "Strike / guard / throw triangle", "rps"],
  "combo-cancel": ["buffer", "早め入力をバッファしてキャンセル", "Buffer early input and cancel", "buffer"],
};

function valueLab(kind, lang, title, attrs = "", extraButtons = [], extraBoard = "", extraValues = null) {
  const isJa = lang === "ja";
  const step = isJa ? "1ステップ" : "Step";
  const reset = isJa ? "もとに戻す" : "Reset";
  const buttons = extraButtons.length
    ? extraButtons
    : [
        ["data-lab-step", isJa ? "進める" : "Step", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ];
  const values = extraValues || [
    [isJa ? "値 A" : "value A", ""],
    [isJa ? "値 B" : "value B", ""],
    [isJa ? "結果" : "result", ""],
  ];
  return { kind, title, attrs, buttons, values, board: extraBoard || '<div class="lab-mono-board">LAB</div>', isJa, step, reset };
}

function renderLab(spec, lang) {
  const isJa = lang === "ja";
  // Full custom specs from `labs`
  if (spec.ja && spec.en) {
    const s = spec[lang];
    const buttons =
      s.buttons ||
      [
        ["data-lab-step", s.step || (isJa ? "進める" : "Step"), "primary"],
        ...(s.eat ? [["data-lab-eat", s.eat, ""]] : []),
        ["data-lab-reset", s.reset, "quiet"],
      ];
    const btnHtml = buttons
      .map(([attr, label, style]) => {
        const cls = style === "primary" ? "lab-button lab-button-primary" : style === "quiet" ? "lab-button lab-button-quiet" : "lab-button";
        return `<button type="button" class="${cls}" ${attr}>${label}</button>`;
      })
      .join("\n            ");
    const valHtml = s.values
      .map((v, i) => {
        const out = s.outs[i];
        const attr = out[0] ? ` ${out[0]}` : "";
        return `<div><span>${v[0]}</span><strong${attr}>${out[1]}</strong><small>${v[1]}</small></div>`;
      })
      .join("\n            ");
    return `
      <div class="motion-lab" data-lab="${spec.kind}" ${s.attrs || ""} aria-labelledby="lab-title-${spec.kind}">
        <div class="lab-copy">
          <p class="eyebrow">${s.eye}</p>
          <h3 id="lab-title-${spec.kind}">${s.title}</h3>
          <p>${s.body}</p>
          <div class="lab-controls">
            ${btnHtml}
          </div>
          <p class="lab-hint">${s.hint}</p>
        </div>
        <div class="lab-visual">
          <div class="lab-board" role="img" aria-label="lab">${s.board}</div>
          <div class="lab-values" aria-live="polite">
            ${valHtml}
          </div>
        </div>
      </div>
`;
  }
  return "";
}

function compactLabHtml(slug, kind, titleJa, titleEn, lang) {
  const isJa = lang === "ja";
  const title = isJa ? titleJa : titleEn;
  const reset = isJa ? "もとに戻す" : "Reset";
  const configs = {
    entities: {
      attrs: "",
      buttons: [
        ["data-lab-fire", isJa ? "発射" : "Fire", "primary"],
        ["data-lab-step", isJa ? "進める" : "Step", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "弾の数" : "shots", "data-lab-count", "0"],
        [isJa ? "リスト" : "list", "data-lab-list", "—"],
        ["y", "", "↑"],
      ],
      board: `<div class="lab-mono-board">${isJa ? "shots []" : "shots []"}</div>`,
      body: isJa
        ? "発射で弾を追加し、進めるで上へ動かします。画面外（y≤0）の弾はリストから消えます。"
        : "Fire appends a shot; Step moves them up. Shots with y≤0 leave the list.",
    },
    push: {
      attrs: `data-blocked="${isJa ? "壁／押せない" : "blocked"}" data-pushed="${isJa ? "箱を押した" : "pushed"}"`,
      buttons: [
        ["data-lab-right", isJa ? "右へ" : "Right", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "マス" : "cells", "data-lab-map", ". P . B ."],
        [isJa ? "結果" : "note", "data-lab-note", "—"],
        ["P/B", "", "player/box"],
      ],
      board: `<div class="lab-mono-board" data-lab-map>. P . B .</div>`,
      body: isJa
        ? "P がプレイヤー、B が箱です。右へ進むとき、先が箱ならその先が空いているときだけ押せます。"
        : "P is the player, B the crate. Moving right pushes only if the far cell is empty.",
    },
    camera: {
      attrs: "",
      buttons: [
        ["data-lab-left", isJa ? "左へ" : "Left", ""],
        ["data-lab-right", isJa ? "右へ" : "Right", "primary"],
        ["data-lab-step", isJa ? "カメラ追従" : "Ease camera", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["camera", "data-lab-cam", "0"],
        ["player", "data-lab-player", "200"],
        ["screen X", "data-lab-screen", "200"],
      ],
      board: `<div class="lab-board" style="position:relative"><div class="lab-actor" data-lab-actor style="top:50%"></div></div>`,
      body: isJa
        ? "プレイヤーを動かし「カメラ追従」で目標へ寄せます。画面座標 = ワールド − カメラです。"
        : "Move the player, then ease the camera. Screen X = world − camera.",
    },
    ai: {
      attrs: `data-patrol="${isJa ? "パトロール" : "patrol"}" data-chase="${isJa ? "追跡" : "chase"}"`,
      buttons: [
        ["data-lab-closer", isJa ? "近づける" : "Closer", "primary"],
        ["data-lab-farther", isJa ? "離す" : "Farther", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "状態" : "state", "data-lab-state", isJa ? "パトロール" : "patrol"],
        [isJa ? "距離" : "distance", "data-lab-dist", "180"],
        ["閾値", "", "120"],
      ],
      board: `<div class="lab-mono-board">${isJa ? "距離 &lt; 120 → 追跡" : "dist &lt; 120 → chase"}</div>`,
      body: isJa
        ? "プレイヤーとの距離が 120 未満で追跡、それ以外はパトロールです。"
        : "Distance under 120 means chase; otherwise patrol.",
    },
    burst: {
      attrs: "",
      buttons: [
        ["data-lab-less", isJa ? "減らす" : "Fewer", ""],
        ["data-lab-more", isJa ? "増やす" : "More", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "弾数" : "count", "data-lab-count", "8"],
        ["式", "", "2π / n"],
        ["cos/sin", "", "velocity"],
      ],
      board: `<div data-lab-board style="width:100%;height:100%;min-height:280px"></div>`,
      body: isJa
        ? "全周を弾の数で割った角度に cos / sin をかけて円形弾幕を作ります。"
        : "Divide the circle by bullet count; cos/sin build the ring.",
    },
    carry: {
      attrs: 'data-dir="1"',
      buttons: [
        ["data-lab-step", isJa ? "足場を動かす" : "Move platform", "primary"],
        ["data-lab-toggle", isJa ? "乗る/降りる" : "Ride toggle", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "足場X" : "plat X", "data-lab-plat", "100"],
        [isJa ? "プレイヤー" : "player", "data-lab-player", "120"],
        ["delta", "data-lab-delta", "0"],
      ],
      board: `<div class="lab-mono-board">${isJa ? "player += plat - prevPlat" : "player += plat - prevPlat"}</div>`,
      body: isJa
        ? "乗っているときだけ、足場の移動差分をプレイヤー座標へ足します。"
        : "While riding, add the platform’s movement delta to the player.",
    },
    stomp: {
      attrs: `data-stomp="${isJa ? "踏みつけ" : "STOMP"}" data-hurt="${isJa ? "ダメージ" : "HURT"}"`,
      buttons: [
        ["data-lab-fall", isJa ? "落下中" : "Falling", "primary"],
        ["data-lab-rise", isJa ? "上昇中" : "Rising", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["vy", "data-lab-vy", "4.0"],
        [isJa ? "判定" : "result", "data-lab-result", isJa ? "踏みつけ" : "STOMP"],
        ["規則", "", "vy > 0"],
      ],
      board: `<div class="lab-mono-board">${isJa ? "上から当たった？" : "Hitting from above?"}</div>`,
      body: isJa
        ? "vy がプラス（落下）のときだけ踏みつけ成功。上昇中はダメージです。"
        : "Only positive vy (falling) counts as a stomp; rising takes damage.",
    },
    power: {
      attrs: `data-yes="${isJa ? "ON" : "ON"}" data-no="OFF" data-stomp="${isJa ? "踏みつけ" : "stomp"}" data-damage="${isJa ? "ダメージ" : "damage"}"`,
      buttons: [
        ["data-lab-toggle", isJa ? "パワー切替" : "Toggle power", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["powered", "data-lab-powered", "OFF"],
        [isJa ? "衝突結果" : "on contact", "data-lab-result", isJa ? "ダメージ" : "damage"],
        ["flag", "", "bool"],
      ],
      board: `<div class="lab-mono-board">${isJa ? "同じ接触 → 状態で分岐" : "same overlap → branch"}</div>`,
      body: isJa
        ? "powered フラグが立つと、敵接触がダメージから踏みつけへ変わります。"
        : "When powered is true, enemy contact becomes a stomp instead of damage.",
    },
    move8: {
      attrs: "",
      buttons: [
        ["data-lab-cardinal", isJa ? "十字方向" : "Cardinal", ""],
        ["data-lab-diag", isJa ? "斜め" : "Diagonal", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "生の速さ" : "raw speed", "data-lab-raw", "5.66"],
        [isJa ? "正規化後" : "normalized", "data-lab-norm", "4.00"],
        ["公式", "", "v / |v|"],
      ],
      board: `<div class="lab-mono-board">dx,dy → hypot → unit</div>`,
      body: isJa
        ? "斜め入力をそのまま足すと √2 倍速になります。長さで割るとどの方向も同じ速さです。"
        : "Raw diagonal input is √2 faster. Divide by length so every direction shares one speed.",
    },
    aim: {
      attrs: "",
      buttons: [
        ["data-lab-step", isJa ? "1フレーム" : "Next frame", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "狙い" : "target", "data-lab-target", "B"],
        ["cooldown", "data-lab-cd", "0"],
        [isJa ? "発射数" : "shots", "data-lab-shots", "0"],
      ],
      board: `<div class="lab-mono-board">${isJa ? "最短距離の敵へ" : "nearest enemy"}</div>`,
      body: isJa
        ? "毎フレームいちばん近い敵を選び、クールダウン間隔でのみ弾を出します。"
        : "Each frame pick the nearest foe; fire only when the cooldown hits zero.",
    },
    pool: {
      attrs: "",
      buttons: [
        ["data-lab-kill", isJa ? "1体倒す" : "Kill one", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "スロット" : "slots", "data-lab-slots", "E0 E1 …"],
        [isJa ? "撃破" : "kills", "data-lab-kills", "0"],
        ["規則", "", "overwrite"],
      ],
      board: `<div class="lab-mono-board">mobs[i] = spawn()</div>`,
      body: isJa
        ? "固定配列の席を消さず、倒したインデックスへ次の敵を上書きします。"
        : "Do not shrink the array — overwrite the defeated index with a new foe.",
    },
    draft: {
      attrs: `data-combat="${isJa ? "戦闘" : "combat"}" data-draft="${isJa ? "選択" : "draft"}"`,
      buttons: [
        ["data-lab-level", isJa ? "レベルアップ" : "Level up", "primary"],
        ['data-lab-card data-card="A"', "A", ""],
        ['data-lab-card data-card="B"', "B", ""],
        ['data-lab-card data-card="C"', "C", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "モード" : "mode", "data-lab-mode", isJa ? "戦闘" : "combat"],
        [isJa ? "選択" : "pick", "data-lab-pick", "—"],
        ["flag", "", "drafting"],
      ],
      board: `<div class="lab-mono-board">${isJa ? "drafting 中は戦闘停止" : "pause while drafting"}</div>`,
      body: isJa
        ? "レベルアップで選択モードへ。カードを選ぶと戦闘に戻ります。"
        : "Level-up enters draft mode. Picking a card returns to combat.",
    },
    evolve: {
      attrs: "",
      buttons: [
        ["data-lab-evolve", isJa ? "進化させる" : "Evolve", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "武器" : "weapon", "data-lab-name", "Ebi Needle"],
        ["count", "data-lab-wcount", "1"],
        ["cooldown", "data-lab-wcd", "32"],
      ],
      board: `<div class="lab-mono-board">weapon = storm</div>`,
      body: isJa
        ? "発射ループはそのまま。weapon 構造体を差し替えるだけで進化します。"
        : "Keep the fire loop; evolution is assigning a new weapon struct.",
    },
    curve: {
      attrs: "",
      buttons: [
        ["data-lab-step", isJa ? "+5秒" : "+5 sec", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["sec", "data-lab-sec", "0"],
        ["interval", "data-lab-interval", "42"],
        ["speed", "data-lab-speed", "0.85"],
      ],
      board: `<div class="lab-mono-board">max(14, 42-sec/2)</div>`,
      body: isJa
        ? "経過秒で出現間隔を短くし、敵速度を上げます。難しさは時間の関数です。"
        : "Elapsed seconds shorten spawns and raise speed. Difficulty is a function of time.",
    },
    click: {
      attrs: "",
      buttons: [
        ["data-lab-tap", isJa ? "タップ" : "Tap", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "回数" : "count", "data-lab-count", "0"],
        ["input", "", "JustPressed"],
        ["+1", "", ""],
      ],
      board: `<div class="lab-mono-board">count++</div>`,
      body: isJa
        ? "押した瞬間だけ数を増やします。押しっぱなしでは増えません。"
        : "Increment only on the press edge — holding does nothing.",
    },
    shop: {
      attrs: `data-bought="${isJa ? "購入した" : "bought"}" data-cant="${isJa ? "足りない" : "not enough"}"`,
      buttons: [
        ["data-lab-earn", isJa ? "稼ぐ +5" : "Earn +5", ""],
        ["data-lab-buy", isJa ? "買う" : "Buy", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "所持" : "gold", "data-lab-gold", "0"],
        [isJa ? "値段" : "cost", "data-lab-cost", "10"],
        [isJa ? "所有" : "owned", "data-lab-owned", "0"],
      ],
      board: `<div class="lab-mono-board">if gold >= cost</div>`,
      body: isJa
        ? "足りるときだけ購入し、値段を上げます。足りなければ何もしません。"
        : "Buy only when affordable, then raise the price. Otherwise do nothing.",
    },
    idle: {
      attrs: "",
      buttons: [
        ["data-lab-step", isJa ? "0.25秒進める" : "Advance 0.25s", "primary"],
        ["data-lab-buy", isJa ? "機械+1" : "Machine +1", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "ポイント" : "points", "data-lab-points", "0.0"],
        [isJa ? "生産" : "rate", "data-lab-rate", "2/s"],
        ["dt", "", "0.25"],
      ],
      board: `<div class="lab-mono-board">points += rate * dt</div>`,
      body: isJa
        ? "経過時間 dt に機械数をかけて自動生産します。"
        : "Offline-style production: points += machines × dt.",
    },
    save: {
      attrs: "",
      buttons: [
        ["data-lab-add", isJa ? "不在 +30秒" : "Away +30s", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "不在" : "away", "data-lab-away", "0s"],
        [isJa ? "報酬" : "gained", "data-lab-gained", "0"],
        ["rate", "", "8/s"],
      ],
      board: `<div class="lab-mono-board">away * rate</div>`,
      body: isJa
        ? "前回終了時刻との差×生産速度が、戻ってきたときの報酬です。"
        : "Reward on return is elapsed away time × production rate.",
    },
    tile: {
      attrs: `data-blocked="${isJa ? "壁だ" : "wall"}"`,
      buttons: [
        ["data-lab-up", "↑", ""],
        ["data-lab-left", "←", ""],
        ["data-lab-right", "→", "primary"],
        ["data-lab-down", "↓", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "位置" : "pos", "data-lab-pos", "1,1"],
        [isJa ? "向き" : "facing", "data-lab-face", "S"],
        [isJa ? "結果" : "note", "data-lab-note", "—"],
      ],
      board: `<div class="lab-mono-board">${isJa ? "壁は (2,1)" : "wall at (2,1)"}</div>`,
      body: isJa
        ? "移動先が壁なら止まり、向きだけ更新します。"
        : "If the next tile is a wall, stay put but still update facing.",
    },
    flag: {
      attrs: `data-yes="true" data-no="false" data-before="${isJa ? "こんにちは。" : "Hello."}" data-after="${isJa ? "見つかったね！" : "You found it!"}"`,
      buttons: [
        ["data-lab-toggle", isJa ? "フラグ切替" : "Toggle flag", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["flag", "data-lab-flag", "false"],
        [isJa ? "セリフ" : "line", "data-lab-text", isJa ? "こんにちは。" : "Hello."],
        ["if", "", "flag"],
      ],
      board: `<div class="lab-mono-board">if flag { … }</div>`,
      body: isJa
        ? "同じNPCでもフラグでセリフが変わります。"
        : "The same NPC can speak different lines based on a flag.",
    },
    turn: {
      attrs: `data-states="${isJa ? "選択,味方,敵,勝敗" : "select,player,enemy,win"}"`,
      buttons: [
        ["data-lab-step", isJa ? "次の状態" : "Next state", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "状態" : "state", "data-lab-state", isJa ? "選択" : "select"],
        ["machine", "", "enum"],
        ["1つずつ", "", ""],
      ],
      board: `<div class="lab-mono-board">select → player → enemy</div>`,
      body: isJa
        ? "戦闘は状態を1つずつ進めます。いまの状態だけが入力を受けます。"
        : "Battles advance one state at a time. Only the current state accepts input.",
    },
    damage: {
      attrs: "",
      buttons: [
        ["data-lab-buff", isJa ? "強化+5" : "Buff +5", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "攻撃" : "atk", "data-lab-atk", "10"],
        [isJa ? "防御" : "def", "data-lab-def", "3"],
        [isJa ? "ダメージ" : "damage", "data-lab-dmg", "7"],
      ],
      board: `<div class="lab-mono-board">max(1, atk+buff-def)</div>`,
      body: isJa
        ? "基本は攻撃−防御。一時強化は攻撃側に足してから計算します。"
        : "Base is attack − defense. Temporary buffs add to attack first.",
    },
    inv: {
      attrs: `data-ok="${isJa ? "装備した" : "equipped"}" data-cant="${isJa ? "お金不足" : "need gold"}"`,
      buttons: [
        ["data-lab-buy", isJa ? "剣を買う 50G" : "Buy sword 50G", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["gold", "data-lab-gold", "100"],
        [isJa ? "装備" : "item", "data-lab-item", "—"],
        [isJa ? "結果" : "note", "data-lab-note", "—"],
      ],
      board: `<div class="lab-mono-board">gold / equip slot</div>`,
      body: isJa
        ? "所持金を減らし、装備スロットへアイテム名を入れます。"
        : "Spend gold and write the item into an equip slot.",
    },
    scene: {
      attrs: `data-field="${isJa ? "フィールド" : "field"}" data-battle="${isJa ? "戦闘" : "battle"}"`,
      buttons: [
        ["data-lab-walk", isJa ? "エンカウント" : "Encounter", "primary"],
        ['data-lab-region data-region="grass"', isJa ? "草原" : "Grass", ""],
        ['data-lab-region data-region="desert"', isJa ? "砂漠" : "Desert", ""],
        ["data-lab-back", isJa ? "フィールドへ" : "Back to field", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "シーン" : "scene", "data-lab-scene", isJa ? "フィールド" : "field"],
        [isJa ? "敵" : "enemy", "data-lab-enemy", "—"],
        ["table", "", "region"],
      ],
      board: `<div class="lab-mono-board">scene + enemyTable</div>`,
      body: isJa
        ? "シーンを戦闘へ切り替え、地域テーブルから敵を抽選します。"
        : "Switch scene to battle and roll an enemy from the region table.",
    },
    quest: {
      attrs: "",
      buttons: [
        ["data-lab-advance", isJa ? "クエスト進める" : "Advance quest", "primary"],
        ["data-lab-save", isJa ? "セーブ" : "Save", ""],
        ["data-lab-load", isJa ? "ロード" : "Load", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["quest", "data-lab-quest", "0"],
        [isJa ? "保存データ" : "saved", "data-lab-saved", "—"],
        ["localStorage", "", ""],
      ],
      board: `<div class="lab-mono-board">save(quest)</div>`,
      body: isJa
        ? "進行は整数のクエスト番号。セーブして読めば続きから再開できます。"
        : "Progress is an integer quest id. Save and load to resume later.",
    },
    frames: {
      attrs: `data-startup="${isJa ? "発生" : "startup"}" data-active="${isJa ? "持続" : "active"}" data-recovery="${isJa ? "硬直" : "recovery"}"`,
      buttons: [
        ["data-lab-step", isJa ? "1F進める" : "Next frame", "primary"],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["frame", "data-lab-frame", "0"],
        [isJa ? "相" : "phase", "data-lab-phase", isJa ? "発生" : "startup"],
        ["1-8 / 9-12 / 13+", "", ""],
      ],
      board: `<div class="lab-mono-board">startup → active → recovery</div>`,
      body: isJa
        ? "技は発生・持続・硬直の3相。持続のあいだだけ攻撃判定が出ます。"
        : "Moves have startup, active, and recovery. Only active frames can hit.",
    },
    react: {
      attrs: "",
      buttons: [
        ["data-lab-hit", isJa ? "ヒット" : "Hit", "primary"],
        ["data-lab-step", isJa ? "1F進める" : "Next frame", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["hitstop", "data-lab-stop", "0"],
        ["stun", "data-lab-stun", "0"],
        ["v2", "data-lab-v", "0.00"],
      ],
      board: `<div class="lab-mono-board">stop → stun → v*=0.86</div>`,
      body: isJa
        ? "ヒットで世界を止め、のけぞり中は速度を減衰させます。"
        : "On hit, freeze the world, then decay knockback during stun.",
    },
    rps: {
      attrs: `data-win="${isJa ? "勝ち" : "WIN"}" data-lose="${isJa ? "負け" : "LOSE"}" data-clash="${isJa ? "相打ち" : "CLASH"}" data-enemy="strike"`,
      buttons: [
        ['data-lab-enemy data-enemy="strike"', isJa ? "敵:打撃" : "Foe:strike", ""],
        ['data-lab-enemy data-enemy="guard"', isJa ? "敵:ガード" : "Foe:guard", ""],
        ['data-lab-enemy data-enemy="throw"', isJa ? "敵:投げ" : "Foe:throw", ""],
        ['data-lab-pick="strike"', isJa ? "打撃" : "Strike", "primary"],
        ['data-lab-pick="guard"', isJa ? "ガード" : "Guard", ""],
        ['data-lab-pick="throw"', isJa ? "投げ" : "Throw", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        [isJa ? "結果" : "result", "data-lab-result", "—"],
        ["RPS", "", "cycle"],
        ["wins()", "", ""],
      ],
      board: `<div class="lab-mono-board">guard &gt; strike &gt; throw</div>`,
      body: isJa
        ? "先に敵の手を選び、自分の手で勝ち負けを見ます。"
        : "Set the foe’s pick, then choose yours to see who wins.",
    },
    buffer: {
      attrs: "",
      buttons: [
        ["data-lab-press", isJa ? "早めに入力" : "Buffer input", "primary"],
        ["data-lab-step", isJa ? "1F進める" : "Next frame", ""],
        ["data-lab-reset", reset, "quiet"],
      ],
      values: [
        ["buffer", "data-lab-buf", "none"],
        ["life", "data-lab-life", "0"],
        [isJa ? "メモ" : "note", "data-lab-note", "—"],
      ],
      board: `<div class="lab-mono-board">bufferLife = 8</div>`,
      body: isJa
        ? "早めの入力を数フレーム覚え、キャンセル窓で次の技へつなぎます。"
        : "Remember an early press for a few frames, then spend it in the cancel window.",
    },
  };

  const c = configs[kind];
  if (!c) return "";
  const btnHtml = c.buttons
    .map(([attr, label, style]) => {
      const cls = style === "primary" ? "lab-button lab-button-primary" : style === "quiet" ? "lab-button lab-button-quiet" : "lab-button";
      // attr may contain multiple attributes
      return `<button type="button" class="${cls}" ${attr}>${label}</button>`;
    })
    .join("\n            ");
  const valHtml = c.values
    .map(([label, attr, initial]) => {
      const a = attr ? ` ${attr}` : "";
      return `<div><span>${label}</span><strong${a}>${initial}</strong><small></small></div>`;
    })
    .join("\n            ");

  return `
      <div class="motion-lab" data-lab="${kind}" ${c.attrs} aria-labelledby="lab-title-${slug}">
        <div class="lab-copy">
          <p class="eyebrow">TRY IT / LAB</p>
          <h3 id="lab-title-${slug}">${title}</h3>
          <p>${c.body}</p>
          <div class="lab-controls">
            ${btnHtml}
          </div>
          <p class="lab-hint">${isJa ? "本物のゲームの公式を、手でゆっくり回す模型です。" : "A slow-motion model of the real game rule."}</p>
        </div>
        <div class="lab-visual">
          <div class="lab-board" role="img" aria-label="lab">${c.board}</div>
          <div class="lab-values" aria-live="polite">
            ${valHtml}
          </div>
        </div>
      </div>
`;
}

function ensureLearnJs(html, depth) {
  const src = `${"../".repeat(depth)}learn.js`;
  if (html.includes("learn.js")) return html;
  if (html.includes("</footer>")) {
    return html.replace("</footer>", `</footer>\n  <script src="${src}"></script>`);
  }
  return html.replace("</body>", `  <script src="${src}"></script>\n</body>`);
}

function inject(html, labHtml) {
  if (html.includes("motion-lab") || html.includes("data-lab=")) return { html, skipped: true };
  if (!html.includes('class="formula"') && !html.includes("class='formula'")) {
    // insert before why-grid or before closing physics section
    if (html.includes("why-grid")) {
      return { html: html.replace(/<div class="why-grid"/, labHtml + '\n      <div class="why-grid"'), skipped: false };
    }
    return { html, skipped: true, reason: "no formula" };
  }
  return {
    html: html.replace(/<div class="formula"/, labHtml + "\n      <div class=\"formula\""),
    skipped: false,
  };
}

let updated = 0;
let skipped = 0;

for (const entry of curriculum) {
  if (!entry.playable || entry.order > 40) continue;
  const slug = entry.slug;
  if (["tap-target", "timing-meter", "catch-stars", "flappy", "tiny-platformer", "growth-curves", "box-viewer"].includes(slug)) {
    // already have labs — still ensure learn.js
    for (const lang of ["ja", "en"]) {
      const file = join(root, "web", lang, entry.route, "index.html");
      if (!existsSync(file)) continue;
      let html = readFileSync(file, "utf8");
      const depth = entry.route.startsWith("games/") ? 3 : 4;
      const next = ensureLearnJs(html, depth);
      if (next !== html) {
        writeFileSync(file, next);
        updated++;
      }
    }
    continue;
  }

  for (const lang of ["ja", "en"]) {
    const file = join(root, "web", lang, entry.route, "index.html");
    if (!existsSync(file)) {
      console.log("missing", file);
      continue;
    }
    let html = readFileSync(file, "utf8");
    let labHtml = "";
    if (labs[slug]) {
      labHtml = renderLab(labs[slug], lang);
    } else if (compact[slug]) {
      const [kind, tJa, tEn] = compact[slug];
      labHtml = compactLabHtml(slug, kind, tJa, tEn, lang);
    } else {
      console.log("no mapping", slug);
      skipped++;
      continue;
    }
    const result = inject(html, labHtml);
    if (result.skipped) {
      skipped++;
      if (result.reason) console.log("skip", slug, lang, result.reason);
      // still ensure learn.js if lab already present
      const depth = entry.route.startsWith("games/") ? 3 : 4;
      const next = ensureLearnJs(result.html, depth);
      if (next !== result.html) writeFileSync(file, next);
      continue;
    }
    const depth = entry.route.startsWith("games/") ? 3 : 4;
    const withJs = ensureLearnJs(result.html, depth);
    writeFileSync(file, withJs);
    updated++;
    console.log("updated", lang, entry.id);
  }
}

console.log(JSON.stringify({ updated, skipped }));
