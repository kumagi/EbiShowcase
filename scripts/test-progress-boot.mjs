#!/usr/bin/env node
// SPDX-License-Identifier: Apache-2.0

import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import vm from "node:vm";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const bootSource = readFileSync(join(dirname(fileURLToPath(import.meta.url)), "../web/page-boot.js"), "utf8");

const STORAGE_KEY = "ebi-showcase-progress-v1";

function makeStorage({ broken = false } = {}) {
  const map = new Map();
  return {
    map,
    getItem(key) {
      if (broken) throw new Error("denied");
      return map.has(key) ? map.get(key) : null;
    },
    setItem(key, value) {
      if (broken) throw new Error("denied");
      map.set(key, String(value));
    },
  };
}

function makeClassList() {
  const set = new Set();
  return {
    contains: (name) => set.has(name),
    add: (name) => set.add(name),
    remove: (name) => set.delete(name),
  };
}

function makeElement(tag) {
  return {
    tag,
    className: "",
    textContent: "",
    attrs: {},
    children: [],
    hidden: true,
    parent: null,
    classList: makeClassList(),
    dataset: {},
    listeners: {},
    setAttribute(name, value) {
      this.attrs[name] = value;
    },
    getAttribute(name) {
      return this.attrs[name];
    },
    removeAttribute(name) {
      delete this.attrs[name];
    },
    appendChild(child) {
      child.parent = this;
      this.children.push(child);
      return child;
    },
    addEventListener(type, fn) {
      (this.listeners[type] ||= []).push(fn);
    },
    dispatch(type) {
      for (const fn of this.listeners[type] || []) fn();
    },
    querySelector(selector) {
      if (selector === "b" && this.tag === "b") return this;
      return this.children.find((child) =>
        selector === "b" ? child.tag === "b" : selector === "div" ? child.tag === "div" : false,
      );
    },
  };
}

function runPage({ pathname, lang = "ja", storage, dom }) {
  const context = vm.createContext({
    window: {},
    location: { pathname, href: `https://example.test${pathname}` },
    document: {
      documentElement: { lang },
      querySelector: (selector) => dom.select[selector] ?? null,
      querySelectorAll: (selector) => dom.selectAll[selector] ?? [],
      createElement: (tag) => makeElement(tag),
    },
    localStorage: storage,
    URL,
    Date,
  });
  vm.runInContext(bootSource, context);
  return context;
}

const CARD_SELECTOR =
  "a.course-card[href], a.vfx-course-card[href], a.test-course-card[href], a.path-step[href], a.track-card[href]";

/* ---------- 1. playing a lesson records progress ---------- */

const frame = makeElement("iframe");
frame.dataset.gameSrc = "../../../play/tap-target/";
const heading = makeElement("h1");
heading.textContent = "  光る丸をタッチ ";
const lessonStorage = makeStorage();
runPage({
  pathname: "/ja/games/tap-target/",
  storage: lessonStorage,
  dom: {
    select: {
      "iframe[data-game-src]": frame,
      "iframe[src*='/play/']": null,
      h1: heading,
    },
    selectAll: { "iframe[data-game-src]": [frame] },
  },
});
let saved = JSON.parse(lessonStorage.map.get(STORAGE_KEY));
assert.equal(saved.last.id, "/ja/games/tap-target/");
assert.equal(saved.last.title, "光る丸をタッチ");
assert.ok(!saved.done["/ja/games/tap-target/"]);

frame.dispatch("load");
saved = JSON.parse(lessonStorage.map.get(STORAGE_KEY));
assert.equal(saved.done["/ja/games/tap-target/"] > 0, true);

/* ---------- 2. hubs render badges from stored progress ---------- */

const played = JSON.parse(lessonStorage.map.get(STORAGE_KEY));
played.done["/ja/tracks/platformer/tiny-platformer/"] = Date.now();
const hubStorage = makeStorage();
hubStorage.setItem(STORAGE_KEY, JSON.stringify(played));

const courseCard = makeElement("a");
courseCard.href = "../../games/tap-target/";
courseCard.classList.add("course-card");
const playedVfx = makeElement("a");
playedVfx.href = "./vfx-stamp/";
playedVfx.classList.add("vfx-course-card");
const unplayedStep = makeElement("a");
unplayedStep.href = "./moving-platforms/";
unplayedStep.classList.add("path-step");
const arrow = makeElement("b");
arrow.textContent = "→";
unplayedStep.children.push(arrow);
const playedStep = makeElement("a");
playedStep.href = "./tiny-platformer/";
playedStep.classList.add("path-step");
const playedArrow = makeElement("b");
playedArrow.textContent = "→";
playedStep.children.push(playedArrow);
const trackCard = makeElement("a");
trackCard.href = "../platformer/";
trackCard.classList.add("track-card");
const trackBody = makeElement("div");
trackCard.children.push(trackBody);

runPage({
  pathname: "/ja/tracks/platformer/",
  storage: hubStorage,
  dom: {
    select: {},
    selectAll: { [CARD_SELECTOR]: [courseCard, playedVfx, unplayedStep, playedStep, trackCard] },
  },
});

assert.equal(courseCard.classList.contains("is-played"), true);
assert.equal(courseCard.children[0].textContent, "✓");
assert.equal(courseCard.children[0].attrs["aria-label"], "プレイ済み");
assert.equal(unplayedStep.classList.contains("is-played"), false);
assert.equal(playedStep.classList.contains("is-played"), true);
assert.equal(playedArrow.textContent, "✓");
assert.equal(trackCard.classList.contains("is-played"), false);
assert.equal(trackBody.children.length, 1);
assert.equal(trackBody.children[0].className, "prog-count");
assert.equal(trackBody.children[0].textContent, "1");

/* relative hrefs resolve against the current directory */
assert.equal(courseCard.children.length, 1);

/* ---------- 3. home resume card points at the last lesson ---------- */

const resumeSlot = makeElement("div");
resumeSlot.hidden = true;
const homeStorage = makeStorage();
homeStorage.setItem(STORAGE_KEY, JSON.stringify(played));
runPage({
  pathname: "/EbiShowcase/ja/",
  storage: homeStorage,
  dom: {
    select: { "[data-progress-continue]": resumeSlot },
    selectAll: {},
  },
});
assert.equal(resumeSlot.hidden, false);
const resumeLink = resumeSlot.children.find((child) => child.className === "prog-resume");
assert.equal(resumeLink.href, "/EbiShowcase/ja/games/tap-target/");
assert.equal(resumeLink.children[0].textContent, "つづきから");
assert.equal(resumeLink.children[1].textContent, "光る丸をタッチ");
const countLine = resumeSlot.children.find((child) => child.className === "prog-resume-count");
assert.match(countLine.textContent, /^2 本のレッスンをプレイ済み$/);

const enResumeSlot = makeElement("div");
runPage({
  pathname: "/EbiShowcase/en/",
  lang: "en",
  storage: homeStorage,
  dom: {
    select: { "[data-progress-continue]": enResumeSlot },
    selectAll: {},
  },
});
const enResume = enResumeSlot.children.find((child) => child.className === "prog-resume");
assert.equal(enResume.children[0].textContent, "Continue");

/* ---------- 4. denied localStorage never breaks the page ---------- */

runPage({
  pathname: "/ja/games/timing-meter/",
  storage: makeStorage({ broken: true }),
  dom: {
    select: { "iframe[data-game-src]": frame, h1: heading },
    selectAll: { "iframe[data-game-src]": [frame] },
  },
});

console.log("OK — progress recording, badges, resume card, and storage fallback all behave.");
