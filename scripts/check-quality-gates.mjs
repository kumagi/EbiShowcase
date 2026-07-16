#!/usr/bin/env node
/**
 * Deterministic runners for docs/quality-gates/catalog.json.
 *
 * Usage:
 *   node scripts/check-quality-gates.mjs
 *   node scripts/check-quality-gates.mjs --family structure,site
 *   node scripts/check-quality-gates.mjs --gate authoring.rule-named --sample 12
 *   node scripts/check-quality-gates.mjs --lenses loop,authoring --json
 *   node scripts/check-quality-gates.mjs --list
 *
 * SPDX-License-Identifier: Apache-2.0
 */
import { existsSync, readFileSync, readdirSync, statSync } from "node:fs";
import { join } from "node:path";
import { pathToFileURL } from "node:url";
import { curriculum, gated } from "./curriculum.mjs";

const root = join(new URL("..", import.meta.url).pathname);
const catalogPath = join(root, "docs/quality-gates/catalog.json");
const catalog = JSON.parse(readFileSync(catalogPath, "utf8"));

const args = process.argv.slice(2);
const wantJson = args.includes("--json");
const wantList = args.includes("--list");
const wantStrict = args.includes("--strict");
const sampleIdx = args.indexOf("--sample");
const sampleN = sampleIdx >= 0 ? Number(args[sampleIdx + 1]) : Infinity;
const familyArg = argValue("--family");
const gateArg = argValue("--gate");
const lensesArg = argValue("--lenses");
const appliesArg = argValue("--applies");

function argValue(flag) {
  const i = args.indexOf(flag);
  return i >= 0 ? args[i + 1] : "";
}

function pageKind(entry) {
  if (entry.group === "vfx") return "vfx";
  if (entry.id.startsWith("core/")) return "core";
  if (entry.id.startsWith("tracks/")) return "track";
  return "other";
}

function walkHtml(dir, out = []) {
  if (!existsSync(dir)) return out;
  for (const name of readdirSync(dir)) {
    const path = join(dir, name);
    if (statSync(path).isDirectory()) walkHtml(path, out);
    else if (name.endsWith(".html")) out.push(path);
  }
  return out;
}

function readLesson(lang, route) {
  const file = join(root, "web", lang, route, "index.html");
  if (!existsSync(file)) return { file, html: "", missing: true };
  return { file, html: readFileSync(file, "utf8"), missing: false };
}

function conceptBodies(html) {
  const row = html.match(/class="concept-row"[^>]*>([\s\S]*?)<\/div>/);
  if (!row) return [];
  return [...row[1].matchAll(/<article[\s\S]*?<p>([\s\S]*?)<\/p>/g)].map((m) =>
    m[1].replace(/<[^>]+>/g, "").replace(/\s+/g, " ").trim(),
  );
}

function challengeText(html) {
  const blocks = [...html.matchAll(/class="challenge"[^>]*>([\s\S]*?)<\/article>/g)];
  return blocks.map((m) => m[1].replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim()).join("\n");
}

function hasAxiomWording(html) {
  if (html.includes("update-draw-contract") || html.includes("Update / Draw") || html.includes("Update/Draw")) return true;
  // Authors naturally emphasize Update and Draw with inline markup. Check the
  // words readers see, rather than letting <strong> or <code> split a sentence.
  const text = html.replace(/<[^>]+>/g, " ").replace(/\s+/g, " ");
  const ja =
    /Update.{0,80}(入力|状態|書き換)/.test(text) &&
    /Draw.{0,80}(投影|絵|描)/.test(text) &&
    /(書き換えてはならない|書き換えない|変更しません|書かない|mutat)/i.test(text);
  const en =
    /Update.{0,100}(owns|mutation|mutate|input|changes)/i.test(text) &&
    /Draw.{0,100}(project|pixels|render)/i.test(text) &&
    /(must not write|never write|does not mutate|never mutates|never changes)/i.test(text);
  return ja || en;
}

const findings = [];

function note(gateId, severity, ok, target, detail) {
  findings.push({ gate: gateId, severity, ok, target, detail });
}

function selectedGates() {
  let gates = catalog.gates;
  if (familyArg) {
    const set = new Set(familyArg.split(",").map((s) => s.trim()));
    gates = gates.filter((g) => set.has(g.family));
  }
  if (gateArg) {
    const set = new Set(gateArg.split(",").map((s) => s.trim()));
    gates = gates.filter((g) => set.has(g.id));
  }
  if (appliesArg) {
    const set = new Set(appliesArg.split(",").map((s) => s.trim()));
    gates = gates.filter((g) => g.applies_to.some((a) => set.has(a)));
  }
  if (lensesArg) {
    const set = new Set(lensesArg.split(",").map((s) => s.trim()));
    const all = set.has("all") || set.has("*");
    gates = gates.filter((g) => g.check === "llm" && (all || set.has(g.family)));
  }
  return gates;
}

function lessonTargets() {
  const playable = curriculum.filter((e) => e.playable);
  const buildRoutes = [];
  for (const lang of ["ja", "en"]) {
    const buildRoot = join(root, "web", lang, "build");
    if (!existsSync(buildRoot)) continue;
    for (const name of readdirSync(buildRoot)) {
      if (name === "index.html") continue;
      const route = `build/${name}`;
      if (existsSync(join(buildRoot, name, "index.html"))) {
        buildRoutes.push({ id: `build/${name}`, route, slug: name, group: "build", order: 0, playable: true });
      }
    }
  }
  // Dedupe build by route using ja list only for iteration with langs
  const buildUnique = [];
  const seen = new Set();
  for (const b of buildRoutes) {
    if (seen.has(b.route)) continue;
    seen.add(b.route);
    buildUnique.push(b);
  }
  return [...playable, ...buildUnique].slice(0, Number.isFinite(sampleN) ? sampleN : undefined);
}

function runStructureGates(gates) {
  const ids = new Set(gates.map((g) => g.id));
  const targets = lessonTargets().filter((e) => e.group !== "build");
  for (const entry of targets) {
    const kind = pageKind(entry);
    for (const lang of ["ja", "en"]) {
      const { html, missing, file } = readLesson(lang, entry.route);
      const target = `${lang}/${entry.route}`;
      if (missing) {
        if (ids.has("structure.go-exists") || ids.has("site.bilingual-pair")) {
          note("site.bilingual-pair", "fail", false, target, "missing HTML page");
        }
        continue;
      }
      const checks = [
        ["structure.deep-dive", /DEEP DIVE/],
        ["structure.concept-row", /class="concept-row"/],
        ["structure.motion-lab", /class="motion-lab"|data-lab=/],
        ["structure.code-lesson", /class="code-lesson"/],
        ["structure.why-grid", /class="why-grid"/],
        ["structure.learn-js", /learn\.js/],
      ];
      for (const [id, re] of checks) {
        if (!ids.has(id)) continue;
        if (!gates.find((g) => g.id === id)?.applies_to.includes(kind) && kind !== "core" && kind !== "track" && kind !== "vfx") continue;
        note(id, "fail", re.test(html), target, re.test(html) ? "ok" : "missing marker");
      }
      if (ids.has("structure.play-iframe")) {
        const playOk =
          html.includes(`play/${entry.slug}/`) || html.includes("game.html") || html.includes(`play/${entry.slug}"`);
        note("structure.play-iframe", "fail", playOk, target, playOk ? "ok" : `missing play/${entry.slug}`);
      }
      if (ids.has("structure.core-full-source") && kind === "core" && (entry.order <= 12 || entry.id.startsWith("core/"))) {
        const ok =
          html.includes('class="full-code"') &&
          html.includes("data-embed-source=") &&
          html.includes("data-copy") &&
          /data-embed-slot[\s\S]*package main/.test(html);
        note("structure.core-full-source", "fail", ok, target, ok ? "ok" : "missing full embed/copy");
      }
      if (ids.has("structure.core-update-workshop") && kind === "core" && (entry.order <= 12 || entry.id.startsWith("core/"))) {
        const workshop = html.match(/<!-- core-update-workshop:start -->([\s\S]*?)<!-- core-update-workshop:end -->/);
        const steps = workshop ? (workshop[1].match(/class="update-build-step"/g) || []).length : 0;
        const copies = workshop ? (workshop[1].match(/data-copy/g) || []).length : 0;
        const ok = Boolean(workshop) && steps >= 2 && copies === steps && /Update/.test(workshop[1]);
        note("structure.core-update-workshop", "fail", ok, target, ok ? `${steps} explained blocks` : "missing or incomplete Update workshop");
      }
      if (ids.has("structure.go-exists")) {
        const src = entry.source || (entry.slug === "flappy" ? "game/main.go" : null);
        const path = src ? join(root, src) : entry.group === "track" ? join(root, "games", "tracks", entry.id.split("/")[1], entry.slug, "main.go") : null;
        // curriculum entries already know playable via source
        const ok = entry.playable || (path && existsSync(path));
        note("structure.go-exists", "fail", ok, entry.id, ok ? "ok" : "missing main.go");
      }
      if (ids.has("site.ogp-present")) {
        const ok = /property="og:title"/.test(html) && /name="twitter:card"/.test(html);
        note("site.ogp-present", "warn", ok, target, ok ? "ok" : "missing ogp/twitter");
      }
      if (ids.has("a11y.iframe-title")) {
        const iframes = [...html.matchAll(/<iframe\b[^>]*>/gi)];
        const playFrames = iframes.filter((m) => /lesson-game-frame|play\//.test(m[0]));
        const ok = playFrames.length === 0 || playFrames.every((m) => /\btitle="[^"]+"/.test(m[0]));
        note("a11y.iframe-title", "warn", ok, target, ok ? "ok" : "play iframe missing title");
      }
      if (ids.has("pedagogy.concept-row-ordered")) {
        const bodies = conceptBodies(html);
        let ok = true;
        let detail = "ok";
        if (bodies.length >= 3) {
          if (bodies[0] === bodies[1] || bodies[1] === bodies[2] || bodies[0] === bodies[2]) {
            ok = false;
            detail = "duplicate concept-row paragraphs";
          }
        }
        note("pedagogy.concept-row-ordered", "warn", ok, target, detail);
      }
      if (ids.has("loop.axiom-visible-early") && (kind === "core" || kind === "build") && (entry.order <= 3 || kind === "build" || /tap-target|timing-meter|catch-stars/.test(entry.slug || ""))) {
        const ok = hasAxiomWording(html);
        note("loop.axiom-visible-early", "warn", ok, target, ok ? "ok" : "axiom wording / contract missing");
      }
      if (ids.has("authoring.rule-named") || ids.has("authoring.edit-path") || ids.has("authoring.verify-hint")) {
        const challenge = challengeText(html);
        if (ids.has("authoring.rule-named")) {
          const ok = /YOUR FIRST RULE|1つルール|Write one rule|Write and verify one rule/i.test(html);
          note("authoring.rule-named", "warn", ok, target, ok ? "ok" : "no YOUR FIRST RULE style challenge");
        }
        if (ids.has("authoring.edit-path")) {
          const ok = /(?:games|game)\/[\w./-]+\.go|graduation\/[\w./-]+|build\/[\w./-]+/.test(challenge) || /(?:games|game)\/[\w./-]+\.go/.test(html);
          note("authoring.edit-path", "warn", ok, target, ok ? "ok" : "no games/... path in challenge");
        }
        if (ids.has("authoring.verify-hint")) {
          const ok = /go test|go run|ブラウザ|browser|確かめ|verify/i.test(challenge) || /go test|go run/.test(html);
          note("authoring.verify-hint", "warn", ok, target, ok ? "ok" : "no verify hint");
        }
      }
      if (ids.has("authoring.dual-layer-code") && kind === "track") {
        const ok =
          /編集する入口|Edit entry|entry file|学習者が開いて/.test(html) ||
          (html.match(/class="code-lesson"/g) || []).length >= 2;
        note("authoring.dual-layer-code", "warn", ok, target, ok ? "ok" : "single-layer code only");
      }
      void file;
    }
  }

  // Build track lessons
  if (gates.some((g) => g.applies_to.includes("build"))) {
    for (const entry of lessonTargets().filter((e) => e.group === "build")) {
      for (const lang of ["ja", "en"]) {
        const { html, missing } = readLesson(lang, entry.route);
        const target = `${lang}/${entry.route}`;
        if (missing) continue;
        if (ids.has("loop.axiom-visible-early")) {
          const ok = hasAxiomWording(html);
          note("loop.axiom-visible-early", "warn", ok, target, ok ? "ok" : "axiom wording missing");
        }
        if (ids.has("authoring.rule-named")) {
          const ok = /YOUR FIRST RULE|1つルール|Write one rule|次に足す/i.test(html);
          note("authoring.rule-named", "warn", ok, target, ok ? "ok" : "missing write-next / RULE cue");
        }
      }
    }
  }
}

function runSiteGates(gates) {
  const ids = new Set(gates.map((g) => g.id));
  if (ids.has("site.gated-count")) {
    const vfx = curriculum.filter((e) => e.group === "vfx");
    const ok = gated.length === 208 && vfx.length === 29;
    note("site.gated-count", "fail", ok, "curriculum", `gated=${gated.length} vfx=${vfx.length}`);
  }
  if (ids.has("site.bilingual-pair")) {
    let bad = 0;
    for (const entry of curriculum) {
      for (const lang of ["ja", "en"]) {
        if (!existsSync(join(root, "web", lang, entry.route, "index.html"))) bad++;
      }
    }
    note("site.bilingual-pair", "fail", bad === 0, "curriculum", bad === 0 ? "ok" : `${bad} missing pages`);
  }
  if (ids.has("authoring.copy-no-legacy-tuning-ogp")) {
    const files = walkHtml(join(root, "web"));
    const legacy = /<meta (?:property="og:description"|name="twitter:description")[^>]*値を変えて[^>]*>/;
    const offenders = files.filter((path) => legacy.test(readFileSync(path, "utf8")));
    note(
      "authoring.copy-no-legacy-tuning-ogp",
      "fail",
      offenders.length === 0,
      "web/**/*.html",
      offenders.length ? `${offenders.length} OGP offenders` : "ok",
    );
  }
  if (ids.has("brand.hero-name")) {
    const files = walkHtml(join(root, "web"));
    const offenders = files.filter((path) => /\bEbi Boy\b/.test(readFileSync(path, "utf8")));
    note("brand.hero-name", "fail", offenders.length === 0, "web/**/*.html", offenders.length ? offenders.slice(0, 5).join(", ") : "ok");
  }
  if (ids.has("site.home-thumbnails")) {
    // Delegate detail to check-site-metadata; here only manifest presence.
    const manifest = join(root, "web/assets/home-thumbnails/manifest.json");
    const ok = existsSync(manifest);
    note("site.home-thumbnails", "fail", ok, "home-thumbnails", ok ? "manifest present (run check-site-metadata for deep audit)" : "missing manifest");
  }
}

function lensesPayload(gates) {
  return gates
    .filter((g) => g.check === "llm")
    .map((g) => ({
      id: g.id,
      family: g.family,
      severity: g.severity,
      summary: g.summary,
      prompt_hint: g.prompt_hint || g.summary,
      applies_to: g.applies_to,
    }));
}

function main() {
  const gates = selectedGates();
  if (wantList || lensesArg) {
    const payload = {
      version: catalog.version,
      meters: catalog.meters,
      gates: lensesArg
        ? lensesPayload(gates)
        : gates.map((g) => ({
            id: g.id,
            family: g.family,
            severity: g.severity,
            check: g.check,
            applies_to: g.applies_to,
            summary: g.summary,
          })),
    };
    console.log(JSON.stringify(payload, null, 2));
    return 0;
  }

  const det = gates.filter((g) => g.check === "deterministic");
  const structureLike = det.filter((g) => ["structure", "authoring", "pedagogy", "loop", "a11y", "site"].includes(g.family));
  if (structureLike.length) runStructureGates(structureLike);
  runSiteGates(det);

  const fails = findings.filter((f) => !f.ok && f.severity === "fail");
  const warns = findings.filter((f) => !f.ok && f.severity === "warn");
  const summary = {
    catalog: catalog.version,
    checked_gates: [...new Set(findings.map((f) => f.gate))],
    findings: findings.length,
    fail: fails.length,
    warn: warns.length,
    ok: fails.length === 0 && (!wantStrict || warns.length === 0),
  };

  if (wantJson) {
    console.log(JSON.stringify({ summary, findings, llm_gates: catalog.gates.filter((g) => g.check === "llm").map((g) => g.id) }, null, 2));
  } else {
    console.log(`Quality gates v${catalog.version}: fail=${fails.length} warn=${warns.length} notes=${findings.length}`);
    for (const f of [...fails, ...warns].slice(0, 40)) {
      console.log(`  [${f.severity}] ${f.gate} @ ${f.target}: ${f.detail}`);
    }
    if (fails.length + warns.length > 40) console.log(`  … ${fails.length + warns.length - 40} more`);
    const skipped = catalog.gates.filter((g) => g.check !== "deterministic");
    console.log(`LLM/human gates in catalog (not executed here): ${skipped.length}. Use --lenses <family> for crawler prompts.`);
  }

  if (fails.length) return 1;
  if (wantStrict && warns.length) return 1;
  return 0;
}

if (import.meta.url === pathToFileURL(process.argv[1]).href) {
  process.exit(main());
}

export { catalog, selectedGates, lensesPayload };
