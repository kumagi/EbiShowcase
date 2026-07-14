#!/usr/bin/env node
/**
 * Strict, evidence-backed Ralph loop for docs/ROADMAP_RALPH_LOOP.md.
 * SPDX-License-Identifier: Apache-2.0
 */
import { execFileSync } from "node:child_process";
import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const checklistPath = join(root, "docs", "ROADMAP_RALPH_LOOP.md");
const evidenceRoot = join(root, "docs", "roadmap-evidence");
const command = process.argv[2] || "next";
const argument = process.argv[3];
const full = process.argv.includes("--full");

function load() {
  const markdown = readFileSync(checklistPath, "utf8");
  const tasks = [];
  const pattern = /^- \[([ xX])\] `([^`]+)` — (.+)$/gm;
  let match;
  while ((match = pattern.exec(markdown))) {
    tasks.push({ done: match[1].toLowerCase() === "x", id: match[2], title: match[3], index: match.index });
  }
  if (!tasks.length) throw new Error("ロードマップのタスクを読み取れませんでした。");
  const ids = new Set();
  for (const task of tasks) {
    if (!/^[A-Z0-9-]+$/.test(task.id)) throw new Error(`不正なタスクID: ${task.id}`);
    if (ids.has(task.id)) throw new Error(`重複したタスクID: ${task.id}`);
    ids.add(task.id);
  }
  return { markdown, tasks };
}

function taskByID(tasks, id) {
  const task = tasks.find((item) => item.id === id);
  if (!task) throw new Error(`タスクが見つかりません: ${id}`);
  return task;
}

function phaseOf(id) {
  return id.match(/^P\d+/)?.[0] || "OTHER";
}

function requirements(id) {
  // Authoring Pass (see docs/ROADMAP_RALPH_LOOP.md)
  if (/^P0-0[1-8]$/.test(id) || /^P0-LOOP-0[1-4]$/.test(id)) return ["Implementation", "Automated checks", "Manual review"];
  if (id === "P1-BT-SPEC" || id === "P2-HELP" || id === "P2-HAND-PATTERN") {
    return ["Implementation", "Automated checks", "Manual review"];
  }
  if (/^P1-BT-0[1-4]$/.test(id) || id === "P1-BT-HUB" || id === "P1-BT-VERIFY") {
    return ["Edit target", "Next lines", "RULE challenge", "Desktop", "Phone", "Japanese", "English", "Tests"];
  }
  if (/^P1-CORE-0[1-5]$/.test(id) || id === "P1-CORE-LINK") {
    return ["Edit target", "RULE challenge", "Japanese", "English", "Automated checks"];
  }
  if (/^P2-(RHY|RAY|TD|REV|TOP)-AUDIT$/.test(id)) {
    return ["Mismatch inventory", "Edit targets", "Japanese", "English"];
  }
  if (/^P2-(RHY|RAY|TD|REV|TOP)-GEN$/.test(id)) {
    return ["Dual panel", "Unique concept-row", "RULE challenge", "Japanese", "English", "Automated checks"];
  }
  if (/^P2-(RHY|RAY|TD|REV|TOP)-VERIFY$/.test(id)) {
    return ["Edit target", "RULE challenge", "Desktop", "Phone", "Japanese", "English", "Tests"];
  }
  if (/^P3-(ARC|EXP|PUZ)-STARTER$/.test(id)) {
    return ["Holey starter", "Failing tests", "TODO mapping", "Automated checks"];
  }
  if (/^P3-(ARC|EXP|PUZ)-ARTICLE$/.test(id)) {
    return ["Article", "TODO mapping", "Japanese", "English", "Manual review"];
  }
  if (/^P3-(ARC|EXP|PUZ)-VERIFY$/.test(id) || id === "P3-HUB" || id === "P3-F30" || id === "P3-NAV" || id === "P3-TRACK-CTA") {
    return ["Article", "Starter", "Tests", "Reference game", "Japanese", "English", "Mobile"];
  }
  if (id === "P4-CORE-LATE" || id === "P4-HAND-01" || id === "P4-HAND-02") {
    return ["Edit target", "RULE challenge", "Japanese", "English", "Automated checks"];
  }
  if (id === "P4-METRIC" || id === "P4-CHECKLIST" || id === "P4-COPY-REGRESS") {
    return ["Implementation", "Automated checks", "Manual review"];
  }
  if (id === "P4-SAMPLE") {
    return ["Sample audit", "Edit target", "RULE challenge", "Japanese", "English", "Manual review"];
  }
  if (id === "P4-RELEASE") return ["Full build", "All tests", "All evidence", "Authoring metrics", "Pages artifact"];
  return ["Implementation", "Automated checks", "Manual review"];
}

function evidencePath(id) {
  return join(evidenceRoot, `${id}.md`);
}

function evidenceTemplate(task) {
  const boxes = requirements(task.id).map((item) => `- [ ] ${item}`).join("\n");
  return `# ${task.id} — ${task.title}\n\nStatus: IN PROGRESS\n\n## Required evidence\n\n${boxes}\n\n## Changes\n\n- Files:\n- Behavior:\n\n## Commands and results\n\n\`\`\`text\ncommand\nresult\n\`\`\`\n\n## Manual checks\n\n| Surface | Representative viewport | Input completed | Result / issue |\n| --- | --- | --- | --- |\n| Desktop | 1440 × 900 | Keyboard + pointer | |\n| Tablet | 768 × 1024 | Touch | |\n| Phone | 390 × 844 | Touch | |\n\n- Japanese:\n- English:\n- Readability / accessibility:\n- Screenshots / recordings:\n`;
}

function validateEvidence(task) {
  const path = evidencePath(task.id);
  if (!existsSync(path)) throw new Error(`証跡がありません: docs/roadmap-evidence/${task.id}.md`);
  const text = readFileSync(path, "utf8");
  if (!/^Status:\s*PASS\s*$/mi.test(text)) throw new Error(`${task.id}: StatusをPASSにしてください。`);
  for (const item of requirements(task.id)) {
    const escaped = item.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    if (!new RegExp(`^- \\[x\\] ${escaped}\\s*$`, "mi").test(text)) throw new Error(`${task.id}: 未確認の証跡項目: ${item}`);
  }
}

function firstUnchecked(tasks) {
  return tasks.find((task) => !task.done);
}

function assertSequential(tasks) {
  const first = tasks.findIndex((task) => !task.done);
  if (first >= 0) {
    const outOfOrder = tasks.slice(first + 1).find((task) => task.done);
    if (outOfOrder) throw new Error(`${outOfOrder.id} が先行タスクより先に完了しています。${tasks[first].id}へ戻ってください。`);
  }
}

function replaceCheckbox(markdown, id, checked) {
  const escaped = id.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  const pattern = new RegExp(`^- \\[[ xX]\\] (\\\`${escaped}\\\` — .+)$`, "m");
  if (!pattern.test(markdown)) throw new Error(`チェック欄を更新できません: ${id}`);
  return markdown.replace(pattern, `- [${checked ? "x" : " "}] $1`);
}

function run(name, args) {
  console.log(`$ ${name} ${args.join(" ")}`);
  execFileSync(name, args, { cwd: root, stdio: "inherit" });
}

function runGreenGoTests() {
  run("bash", ["-lc", "go test $(go list ./... | grep -v '/graduation/.*/starter$')"]);
  run("node", ["scripts/check-graduation-starters.mjs"]);
}

function verifyStructure(tasks) {
  assertSequential(tasks);
  for (const task of tasks.filter((item) => item.done)) validateEvidence(task);
  const phases = new Set(tasks.map((task) => phaseOf(task.id)));
  for (const phase of ["P0", "P1", "P2", "P3", "P4"]) if (!phases.has(phase)) throw new Error(`フェーズがありません: ${phase}`);
}

try {
  const state = load();
  const { tasks } = state;
  if (command === "status") {
    const byPhase = {};
    for (const task of tasks) {
      const phase = phaseOf(task.id);
      byPhase[phase] ||= { done: 0, total: 0 };
      byPhase[phase].total++;
      if (task.done) byPhase[phase].done++;
    }
    const done = tasks.filter((task) => task.done).length;
    console.log(JSON.stringify({ done, total: tasks.length, remaining: tasks.length - done, byPhase }, null, 2));
  } else if (command === "next") {
    assertSequential(tasks);
    const task = firstUnchecked(tasks);
    console.log(JSON.stringify(task ? { id: task.id, phase: phaseOf(task.id), title: task.title, evidence: `docs/roadmap-evidence/${task.id}.md` } : { complete: true, total: tasks.length }, null, 2));
  } else if (command === "list") {
    const filter = argument?.toUpperCase();
    for (const task of tasks) if (!filter || phaseOf(task.id) === filter) console.log(`${task.done ? "DONE" : "TODO"}\t${task.id}\t${task.title}`);
  } else if (command === "evidence") {
    const task = taskByID(tasks, argument);
    const next = firstUnchecked(tasks);
    if (!task.done && next?.id !== task.id) throw new Error(`先に ${next.id} を完了してください。`);
    const path = evidencePath(task.id);
    mkdirSync(evidenceRoot, { recursive: true });
    if (!existsSync(path)) writeFileSync(path, evidenceTemplate(task));
    console.log(path);
  } else if (command === "verify-task") {
    validateEvidence(taskByID(tasks, argument));
    console.log(`${argument}: evidence OK`);
  } else if (command === "check") {
    const task = taskByID(tasks, argument);
    const next = firstUnchecked(tasks);
    if (task.done) throw new Error(`${task.id} はすでに完了しています。`);
    if (next?.id !== task.id) throw new Error(`次のタスクは ${next?.id} です。順番に完了してください。`);
    validateEvidence(task);
    run("git", ["diff", "--check"]);
    writeFileSync(checklistPath, replaceCheckbox(state.markdown, task.id, true));
    console.log(`完了: ${task.id}`);
  } else if (command === "uncheck") {
    const task = taskByID(tasks, argument);
    if (!task.done) throw new Error(`${task.id} は未完了です。`);
    writeFileSync(checklistPath, replaceCheckbox(state.markdown, task.id, false));
    console.log(`未完了へ戻しました: ${task.id}`);
  } else if (command === "verify") {
    verifyStructure(tasks);
    run("git", ["diff", "--check"]);
    if (full) {
      runGreenGoTests();
      run("bash", ["scripts/ralph-loop.sh", "verify"]);
    }
    console.log(`Roadmap evidence verified (${tasks.filter((task) => task.done).length}/${tasks.length})${full ? " with full build" : ""}.`);
  } else if (command === "complete") {
    verifyStructure(tasks);
    const next = firstUnchecked(tasks);
    if (next) throw new Error(`ロードマップは未完了です。次: ${next.id}`);
    run("git", ["diff", "--check"]);
    runGreenGoTests();
    run("bash", ["scripts/ralph-loop.sh", "verify"]);
    console.log(`ROADMAP COMPLETE: ${tasks.length}/${tasks.length}`);
  } else {
    throw new Error("usage: roadmap-ralph-loop.mjs {next|status|list [P0-P4]|evidence ID|verify-task ID|check ID|uncheck ID|verify [--full]|complete}  # Authoring Pass");
  }
} catch (error) {
  console.error(error.message);
  process.exit(1);
}
