// SPDX-License-Identifier: Apache-2.0

const clone = (value) => JSON.parse(JSON.stringify(value));
const keyOf = (x, y) => `${x},${y}`;

async function copyText(value) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(value);
    return;
  }
  const textarea = document.createElement("textarea");
  textarea.value = value;
  textarea.style.position = "fixed";
  textarea.style.opacity = "0";
  document.body.append(textarea);
  textarea.select();
  document.execCommand("copy");
  textarea.remove();
}

const messages = {
  ja: {
    edit: "クリック／タップで直接編集",
    undo: "元に戻す",
    redo: "やり直す",
    reset: "初期例",
    copy: "JSONをコピー",
    download: "JSON保存",
    import: "JSON読込",
    valid: "検査OK：ゲームへ渡せる形です。",
    copied: "JSONをコピーしました。",
    copyFailed: "コピーできませんでした。下のJSON表示から選択してコピーしてください。",
    downloaded: "JSONを保存しました。",
    imported: "JSONを読み込み、同じ編集画面へ復元しました。",
    importFailed: "JSONを読み込めませんでした。形式と値を確認してください。",
    emptyHistory: "これ以上戻せません。",
    jsonTitle: "ゲームへ渡すアセットJSON",
    jsonHint: "この縮小エディタはJSONを作る道具です。本番ではゲーム本体とエディタが同じ型・検査を共有します。",
    palette: "配置するもの",
    inspect: "調整",
    frame: "確認するtick",
    move: "技",
    route: "経路・配置",
    wave: "ウェーブ設定",
    platformerTools: { empty: "消す", ground: "地面", ledge: "足場", coin: "コイン", spawn: "開始", goal: "ゴール", event: "イベント" },
    towerTools: { path: "経路を追加", spawn: "入口", gate: "守る門", build: "建設地点", erase: "消す" },
    phase: { startup: "発生", active: "攻撃判定", recovery: "硬直" },
    fields: { startup: "発生", active: "攻撃判定", recovery: "硬直", damage: "威力", reach: "間合い", enemies: "敵の数", interval: "出現間隔", coins: "初期コイン", lives: "門の耐久", speed: "敵速度倍率" },
    errors: {
      schema: "schemaVersionが対応外です。",
      platformerSize: "横スクロール面のサイズが不正です。",
      platformerSpawn: "開始地点はちょうど1つ必要です。",
      platformerGoal: "ゴールはちょうど1つ必要です。",
      platformerCell: "マスが面の外にあります。",
      fightingMoves: "JABとHEAVYの両方が必要です。",
      fightingFrames: "各技の発生・攻撃判定・硬直は1以上、合計60以下にします。",
      fightingValue: "威力と間合いは正の値にします。",
      towerSpawn: "入口が必要です。",
      towerGate: "守る門が必要です。",
      towerRoute: "経路は入口から門まで、上下左右に隣り合う順で置きます。",
      towerBuild: "建設地点は経路と重ねられません。",
      towerWave: "ウェーブの数・間隔・初期資源は正の値にします。",
    },
  },
  en: {
    edit: "Click or tap to edit directly",
    undo: "Undo",
    redo: "Redo",
    reset: "Reset example",
    copy: "Copy JSON",
    download: "Save JSON",
    import: "Load JSON",
    valid: "Validation passed: this data is ready for a game runtime.",
    copied: "JSON copied.",
    copyFailed: "Could not copy automatically. Select the JSON below and copy it manually.",
    downloaded: "JSON saved.",
    imported: "JSON loaded and restored into the same visual editor.",
    importFailed: "Could not load that JSON. Check its format and values.",
    emptyHistory: "There is nothing else to undo.",
    jsonTitle: "Asset JSON handed to the game",
    jsonHint: "This compact editor creates JSON. In production, the editor and game runtime share the same type and validator.",
    palette: "Place",
    inspect: "Tune",
    frame: "Preview tick",
    move: "Move",
    route: "Route and placement",
    wave: "Wave settings",
    platformerTools: { empty: "Erase", ground: "Ground", ledge: "Ledge", coin: "Coin", spawn: "Spawn", goal: "Goal", event: "Event" },
    towerTools: { path: "Append route", spawn: "Entrance", gate: "Pearl gate", build: "Build slot", erase: "Erase" },
    phase: { startup: "Startup", active: "Active", recovery: "Recovery" },
    fields: { startup: "Startup", active: "Active", recovery: "Recovery", damage: "Damage", reach: "Reach", enemies: "Enemy count", interval: "Spawn interval", coins: "Starting coins", lives: "Gate lives", speed: "Enemy speed scale" },
    errors: {
      schema: "schemaVersion is not supported.",
      platformerSize: "The side-scrolling stage size is invalid.",
      platformerSpawn: "The stage needs exactly one spawn.",
      platformerGoal: "The stage needs exactly one goal.",
      platformerCell: "A cell lies outside the stage.",
      fightingMoves: "Both JAB and HEAVY are required.",
      fightingFrames: "Startup, active, and recovery must be at least 1 and total no more than 60.",
      fightingValue: "Damage and reach must be positive.",
      towerSpawn: "The route needs an entrance.",
      towerGate: "The route needs a pearl gate.",
      towerRoute: "Place the route from entrance to gate in orthogonally adjacent order.",
      towerBuild: "Build slots cannot overlap the route.",
      towerWave: "Wave count, interval, and starting resources must be positive.",
    },
  },
};

export function createPlatformerState() {
  const columns = 18;
  const rows = 8;
  const cells = Array(columns * rows).fill("empty");
  const put = (kind, points) => points.forEach(([x, y]) => { cells[y * columns + x] = kind; });
  put("ground", Array.from({ length: columns }, (_, x) => [x, 7]));
  put("ground", [[5, 6], [6, 6], [11, 6], [12, 6], [13, 6]]);
  put("ledge", [[3, 5], [4, 5], [8, 4], [9, 4], [14, 5], [15, 5]]);
  put("coin", [[4, 4], [8, 3], [12, 5], [15, 4]]);
  put("spawn", [[1, 6]]);
  put("goal", [[17, 6]]);
  put("event", [[10, 6]]);
  return { columns, rows, tileSize: 32, cells };
}

export function exportPlatformerDocument(state) {
  const terrain = [];
  const entities = [];
  state.cells.forEach((kind, index) => {
    if (kind === "empty") return;
    const x = index % state.columns;
    const y = Math.floor(index / state.columns);
    if (kind === "ground" || kind === "ledge") terrain.push({ kind, x, y });
    else if (kind === "event") entities.push({ id: `event-${x}-${y}`, kind, x, y, trigger: "touch", action: "show-tip" });
    else entities.push({ kind, x, y });
  });
  return { schemaVersion: 1, game: "ebi-platformer", id: "coral-run-01", grid: { columns: state.columns, rows: state.rows, tileSize: state.tileSize }, terrain, entities };
}

export function importPlatformerDocument(doc) {
  const state = { columns: doc.grid.columns, rows: doc.grid.rows, tileSize: doc.grid.tileSize, cells: Array(doc.grid.columns * doc.grid.rows).fill("empty") };
  for (const item of [...doc.terrain, ...doc.entities]) state.cells[item.y * state.columns + item.x] = item.kind;
  return state;
}

export function validatePlatformerDocument(doc) {
  const errors = [];
  if (doc.schemaVersion !== 1) errors.push("schema");
  if (!doc.grid || doc.grid.columns < 4 || doc.grid.rows < 4 || doc.grid.tileSize <= 0) return [...errors, "platformerSize"];
  const all = [...(doc.terrain || []), ...(doc.entities || [])];
  if (all.some((item) => item.x < 0 || item.y < 0 || item.x >= doc.grid.columns || item.y >= doc.grid.rows)) errors.push("platformerCell");
  if ((doc.entities || []).filter((item) => item.kind === "spawn").length !== 1) errors.push("platformerSpawn");
  if ((doc.entities || []).filter((item) => item.kind === "goal").length !== 1) errors.push("platformerGoal");
  return errors;
}

export function createFightingState() {
  return {
    selected: "jab",
    frame: 8,
    moves: {
      jab: { id: "jab", startup: 8, active: 5, recovery: 17, damage: 8, reach: 92 },
      heavy: { id: "heavy", startup: 16, active: 7, recovery: 25, damage: 20, reach: 128 },
    },
  };
}

export function exportFightingDocument(state) {
  return { schemaVersion: 1, game: "ebi-fighters", id: "tenjiroh-moves-v1", moves: [clone(state.moves.jab), clone(state.moves.heavy)] };
}

export function importFightingDocument(doc) {
  const byID = Object.fromEntries(doc.moves.map((move) => [move.id, clone(move)]));
  return { selected: "jab", frame: 0, moves: { jab: byID.jab, heavy: byID.heavy } };
}

export function validateFightingDocument(doc) {
  const errors = [];
  if (doc.schemaVersion !== 1) errors.push("schema");
  const moves = doc.moves || [];
  if (!moves.some((move) => move.id === "jab") || !moves.some((move) => move.id === "heavy")) return [...errors, "fightingMoves"];
  if (moves.some((move) => [move.startup, move.active, move.recovery].some((n) => !Number.isFinite(n) || n < 1) || move.startup + move.active + move.recovery > 60)) errors.push("fightingFrames");
  if (moves.some((move) => move.damage <= 0 || move.reach <= 0)) errors.push("fightingValue");
  return errors;
}

export function createTowerDefenseState() {
  return {
    columns: 12,
    rows: 8,
    spawn: { x: 0, y: 2 },
    gate: { x: 11, y: 5 },
    route: [[0, 2], [1, 2], [2, 2], [2, 3], [3, 3], [4, 3], [5, 3], [5, 4], [6, 4], [7, 4], [8, 4], [8, 5], [9, 5], [10, 5], [11, 5]].map(([x, y]) => ({ x, y })),
    buildSlots: [{ x: 1, y: 4 }, { x: 4, y: 1 }, { x: 6, y: 6 }, { x: 9, y: 3 }],
    enemies: 8,
    interval: 55,
    coins: 180,
    lives: 10,
    speed: 1,
  };
}

export function exportTowerDefenseDocument(state) {
  return {
    schemaVersion: 1,
    game: "ebi-defense",
    id: "pearl-gate-editor-01",
    grid: { columns: state.columns, rows: state.rows, tileSize: 48 },
    entrance: clone(state.spawn),
    gate: clone(state.gate),
    route: clone(state.route),
    buildSlots: clone(state.buildSlots),
    waves: [{ id: "tide-runners", enemy: "runner", count: state.enemies, intervalTicks: state.interval, speedScale: state.speed }],
    startingCoins: state.coins,
    lives: state.lives,
  };
}

export function importTowerDefenseDocument(doc) {
  return {
    columns: doc.grid.columns,
    rows: doc.grid.rows,
    spawn: clone(doc.entrance),
    gate: clone(doc.gate),
    route: clone(doc.route),
    buildSlots: clone(doc.buildSlots),
    enemies: doc.waves[0].count,
    interval: doc.waves[0].intervalTicks,
    coins: doc.startingCoins,
    lives: doc.lives,
    speed: doc.waves[0].speedScale,
  };
}

export function validateTowerDefenseDocument(doc) {
  const errors = [];
  if (doc.schemaVersion !== 1) errors.push("schema");
  if (!doc.entrance) errors.push("towerSpawn");
  if (!doc.gate) errors.push("towerGate");
  const route = doc.route || [];
  const validEnds = doc.entrance && doc.gate && route.length >= 2 && keyOf(route[0].x, route[0].y) === keyOf(doc.entrance.x, doc.entrance.y) && keyOf(route.at(-1).x, route.at(-1).y) === keyOf(doc.gate.x, doc.gate.y);
  const adjacent = route.every((point, index) => index === 0 || Math.abs(point.x - route[index - 1].x) + Math.abs(point.y - route[index - 1].y) === 1);
  if (!validEnds || !adjacent) errors.push("towerRoute");
  const routeKeys = new Set(route.map((point) => keyOf(point.x, point.y)));
  if ((doc.buildSlots || []).some((point) => routeKeys.has(keyOf(point.x, point.y)))) errors.push("towerBuild");
  const wave = doc.waves?.[0];
  if (!wave || wave.count <= 0 || wave.intervalTicks <= 0 || wave.speedScale <= 0 || doc.startingCoins < 0 || doc.lives <= 0) errors.push("towerWave");
  return errors;
}

const editorSpecs = {
  platformer: {
    make: createPlatformerState,
    export: exportPlatformerDocument,
    import: importPlatformerDocument,
    validate: validatePlatformerDocument,
    filename: "ebi-platformer-course.json",
  },
  fighting: {
    make: createFightingState,
    export: exportFightingDocument,
    import: importFightingDocument,
    validate: validateFightingDocument,
    filename: "ebi-fighters-moves.json",
  },
  towerDefense: {
    make: createTowerDefenseState,
    export: exportTowerDefenseDocument,
    import: importTowerDefenseDocument,
    validate: validateTowerDefenseDocument,
    filename: "ebi-defense-scenario.json",
  },
};

function actionButton(action, label) {
  return `<button type="button" data-editor-action="${action}">${label}</button>`;
}

function mountEditor(root, kind, lang) {
  const t = messages[lang];
  const spec = editorSpecs[kind];
  let state = spec.make();
  let undoStack = [];
  let redoStack = [];
  let tool = kind === "platformer" ? "ground" : kind === "towerDefense" ? "path" : "jab";
  const storageKey = `ebi-game-data-editor.${kind}`;
  try {
    const saved = localStorage.getItem(storageKey);
    if (saved) state = spec.import(JSON.parse(saved));
  } catch {}

  root.innerHTML = `<div class="asset-editor-head"><div><strong>${t.edit}</strong><span data-editor-status aria-live="polite"></span></div><div class="asset-editor-actions">${actionButton("undo", t.undo)}${actionButton("redo", t.redo)}${actionButton("reset", t.reset)}${actionButton("copy", t.copy)}${actionButton("download", t.download)}<label class="asset-editor-import">${t.import}<input type="file" accept="application/json,.json" data-editor-import></label></div></div><div data-editor-workspace></div><details class="asset-editor-json"><summary>${t.jsonTitle}</summary><p>${t.jsonHint}</p><pre data-editor-json></pre></details>`;
  const workspace = root.querySelector("[data-editor-workspace]");
  const output = root.querySelector("[data-editor-json]");
  const status = root.querySelector("[data-editor-status]");

  const documentValue = () => spec.export(state);
  const persist = () => {
    try { localStorage.setItem(storageKey, JSON.stringify(documentValue())); } catch {}
  };
  const showStatus = (text, bad = false) => {
    status.textContent = text;
    status.classList.toggle("is-error", bad);
  };
  const renderStatus = () => {
    const errors = spec.validate(documentValue());
    showStatus(errors.length ? errors.map((code) => t.errors[code] || code).join(" ") : t.valid, errors.length > 0);
  };
  const commit = (mutate) => {
    undoStack.push(clone(state));
    if (undoStack.length > 80) undoStack.shift();
    redoStack = [];
    mutate(state);
    persist();
    render();
  };

  function renderPlatformer() {
    const tools = Object.entries(t.platformerTools).map(([id, label]) => `<button type="button" class="is-${id}${tool === id ? " is-active" : ""}" data-tool="${id}">${label}</button>`).join("");
    workspace.innerHTML = `<div class="asset-editor-palette"><span>${t.palette}</span>${tools}</div><div class="platformer-editor-board" style="--editor-cols:${state.columns}" role="grid"></div>`;
    const board = workspace.querySelector(".platformer-editor-board");
    state.cells.forEach((kindName, index) => {
      const x = index % state.columns;
      const y = Math.floor(index / state.columns);
      const cell = document.createElement("button");
      cell.type = "button";
      cell.className = `editor-cell is-${kindName}`;
      cell.dataset.index = index;
      cell.setAttribute("role", "gridcell");
      cell.setAttribute("aria-label", `${x}, ${y}: ${t.platformerTools[kindName]}`);
      cell.textContent = ({ coin: "●", spawn: "S", goal: "G", event: "!", ledge: "═" })[kindName] || "";
      cell.addEventListener("click", () => commit((next) => {
        if (tool === "spawn" || tool === "goal") next.cells = next.cells.map((value) => value === tool ? "empty" : value);
        next.cells[index] = tool;
      }));
      board.append(cell);
    });
    workspace.querySelectorAll("[data-tool]").forEach((button) => button.addEventListener("click", () => { tool = button.dataset.tool; render(); }));
  }

  function renderFighting() {
    const move = state.moves[state.selected];
    const total = move.startup + move.active + move.recovery;
    state.frame = Math.min(state.frame, total - 1);
    const timeline = Array.from({ length: total }, (_, frame) => {
      const phase = frame < move.startup ? "startup" : frame < move.startup + move.active ? "active" : "recovery";
      return `<button type="button" class="is-${phase}${frame === state.frame ? " is-current" : ""}" data-frame="${frame}" title="${t.phase[phase]} ${frame + 1}">${frame + 1}</button>`;
    }).join("");
    const field = (name, min, max, step = 1) => `<label><span>${t.fields[name]}</span><input type="range" min="${min}" max="${max}" step="${step}" value="${move[name]}" data-fight-field="${name}"><output>${move[name]}</output></label>`;
    const active = state.frame >= move.startup && state.frame < move.startup + move.active;
    workspace.innerHTML = `<div class="fight-editor-tabs"><span>${t.move}</span><button type="button" data-move="jab" class="${state.selected === "jab" ? "is-active" : ""}">JAB</button><button type="button" data-move="heavy" class="${state.selected === "heavy" ? "is-active" : ""}">HEAVY</button></div><div class="fight-editor-layout"><div><div class="fight-editor-preview"><div class="fighter fighter-player"></div><div class="fight-hitbox${active ? " is-active" : ""}" style="--reach:${Math.min(190, move.reach)}px"></div><div class="fighter fighter-rival"></div><b>${active ? `${t.phase.active} · ${move.damage} DMG` : state.frame < move.startup ? t.phase.startup : t.phase.recovery}</b></div><label class="fight-frame"><span>${t.frame}</span><input type="range" min="0" max="${total - 1}" value="${state.frame}" data-fight-frame><output>${state.frame + 1}/${total}</output></label></div><div class="fight-editor-fields">${field("startup", 1, 30)}${field("active", 1, 16)}${field("recovery", 1, 40)}${field("damage", 1, 40)}${field("reach", 40, 180)}</div></div><div class="fight-timeline" aria-label="${t.frame}">${timeline}</div>`;
    workspace.querySelectorAll("[data-move]").forEach((button) => button.addEventListener("click", () => { state.selected = button.dataset.move; state.frame = 0; render(); }));
    workspace.querySelectorAll("[data-fight-field]").forEach((input) => input.addEventListener("change", () => commit((next) => { next.moves[next.selected][input.dataset.fightField] = Number(input.value); })));
    workspace.querySelector("[data-fight-frame]").addEventListener("input", (event) => { state.frame = Number(event.target.value); render(); });
    workspace.querySelectorAll("[data-frame]").forEach((button) => button.addEventListener("click", () => { state.frame = Number(button.dataset.frame); render(); }));
  }

  function renderTowerDefense() {
    const routeIndex = new Map(state.route.map((point, index) => [keyOf(point.x, point.y), index]));
    const buildKeys = new Set(state.buildSlots.map((point) => keyOf(point.x, point.y)));
    const tools = Object.entries(t.towerTools).map(([id, label]) => `<button type="button" class="is-${id}${tool === id ? " is-active" : ""}" data-tool="${id}">${label}</button>`).join("");
    const field = (name, min, max, step = 1) => `<label><span>${t.fields[name]}</span><input type="range" min="${min}" max="${max}" step="${step}" value="${state[name]}" data-tower-field="${name}"><output>${state[name]}</output></label>`;
    workspace.innerHTML = `<div class="asset-editor-palette"><span>${t.route}</span>${tools}</div><div class="tower-editor-layout"><div class="tower-editor-board" style="--editor-cols:${state.columns}" role="grid"></div><div class="tower-editor-fields"><h4>${t.wave}</h4>${field("enemies", 1, 24)}${field("interval", 20, 120, 5)}${field("coins", 0, 500, 10)}${field("lives", 1, 30)}${field("speed", 0.5, 2, 0.05)}</div></div>`;
    const board = workspace.querySelector(".tower-editor-board");
    for (let y = 0; y < state.rows; y++) {
      for (let x = 0; x < state.columns; x++) {
        const key = keyOf(x, y);
        const index = routeIndex.get(key);
        const isSpawn = state.spawn && key === keyOf(state.spawn.x, state.spawn.y);
        const isGate = state.gate && key === keyOf(state.gate.x, state.gate.y);
        const isBuild = buildKeys.has(key);
        const cell = document.createElement("button");
        cell.type = "button";
        cell.className = `editor-cell${index !== undefined ? " is-path" : ""}${isSpawn ? " is-spawn" : ""}${isGate ? " is-gate" : ""}${isBuild ? " is-build" : ""}`;
        cell.textContent = isSpawn ? "S" : isGate ? "G" : isBuild ? "+" : index !== undefined ? String(index + 1) : "";
        cell.addEventListener("click", () => commit((next) => {
          const same = (point) => point.x === x && point.y === y;
          if (tool === "path" && !next.route.some(same)) next.route.push({ x, y });
          if (tool === "spawn") { next.spawn = { x, y }; next.route = [{ x, y }, ...next.route.filter((point) => !same(point))]; }
          if (tool === "gate") { next.gate = { x, y }; next.route = [...next.route.filter((point) => !same(point)), { x, y }]; }
          if (tool === "build" && !next.buildSlots.some(same)) next.buildSlots.push({ x, y });
          if (tool === "erase") {
            next.route = next.route.filter((point) => !same(point));
            next.buildSlots = next.buildSlots.filter((point) => !same(point));
            if (next.spawn && same(next.spawn)) next.spawn = null;
            if (next.gate && same(next.gate)) next.gate = null;
          }
        }));
        board.append(cell);
      }
    }
    workspace.querySelectorAll("[data-tool]").forEach((button) => button.addEventListener("click", () => { tool = button.dataset.tool; render(); }));
    workspace.querySelectorAll("[data-tower-field]").forEach((input) => input.addEventListener("change", () => commit((next) => { next[input.dataset.towerField] = Number(input.value); })));
  }

  function render() {
    if (kind === "platformer") renderPlatformer();
    if (kind === "fighting") renderFighting();
    if (kind === "towerDefense") renderTowerDefense();
    output.textContent = JSON.stringify(documentValue(), null, 2);
    renderStatus();
  }

  root.querySelectorAll("[data-editor-action]").forEach((button) => button.addEventListener("click", async () => {
    const action = button.dataset.editorAction;
    if (action === "undo") {
      if (!undoStack.length) return showStatus(t.emptyHistory, true);
      redoStack.push(clone(state));
      state = undoStack.pop();
      persist();
      render();
    }
    if (action === "redo") {
      if (!redoStack.length) return showStatus(t.emptyHistory, true);
      undoStack.push(clone(state));
      state = redoStack.pop();
      persist();
      render();
    }
    if (action === "reset") commit(() => { state = spec.make(); });
    if (action === "copy") {
      try {
        await copyText(JSON.stringify(documentValue(), null, 2));
        showStatus(t.copied);
      } catch {
        showStatus(t.copyFailed, true);
      }
    }
    if (action === "download") {
      const url = URL.createObjectURL(new Blob([JSON.stringify(documentValue(), null, 2)], { type: "application/json" }));
      const link = Object.assign(document.createElement("a"), { href: url, download: spec.filename });
      link.click();
      URL.revokeObjectURL(url);
      showStatus(t.downloaded);
    }
  }));

  root.querySelector("[data-editor-import]").addEventListener("change", async (event) => {
    const file = event.target.files?.[0];
    if (!file) return;
    try {
      const doc = JSON.parse(await file.text());
      const errors = spec.validate(doc);
      if (errors.length) throw new Error(errors.join(","));
      undoStack.push(clone(state));
      state = spec.import(doc);
      redoStack = [];
      persist();
      render();
      showStatus(t.imported);
    } catch {
      showStatus(t.importFailed, true);
    }
    event.target.value = "";
  });

  render();
}

if (typeof document !== "undefined") {
  const lang = document.documentElement.lang === "ja" ? "ja" : "en";
  document.querySelectorAll("[data-game-data-editor]").forEach((root) => mountEditor(root, root.dataset.gameDataEditor, lang));
}
