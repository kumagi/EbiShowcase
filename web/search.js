// SPDX-License-Identifier: Apache-2.0
/* Client-side search over the static bilingual lesson index. */
(() => {
  "use strict";

  const app = document.querySelector("[data-search-app]");
  if (!app) return;

  const lang = (document.documentElement.lang || "ja").toLowerCase().startsWith("ja") ? "ja" : "en";
  const t =
    lang === "ja"
      ? {
          placeholder: "キーワードで探す（例: 当たり判定、カメラ、ゲージ）",
          empty: "キーワードを入力すると、関連するレッスンが出ます。",
          none: "見つかりませんでした。別の言葉で試してみましょう。",
          all: "すべて",
          hubLabel: "一覧",
          playableLabel: "PLAYABLE",
          results: (n) => `${n} 件`,
        }
      : {
          placeholder: "Search by keyword (e.g. collision, camera, gauge)",
          empty: "Type a keyword to see matching lessons.",
          none: "No matches yet. Try another word.",
          all: "All",
          hubLabel: "HUB",
          playableLabel: "PLAYABLE",
          results: (n) => `${n} result${n === 1 ? "" : "s"}`,
        };

  const input = app.querySelector("[data-search-input]");
  const status = app.querySelector("[data-search-status]");
  const chipsBox = app.querySelector("[data-search-chips]");
  const listBox = app.querySelector("[data-search-results]");

  let entries = [];
  let activeGroup = null;
  const groupOrder = [];
  const groupSeen = new Map();

  const esc = (s) =>
    s.replace(/[&<>"']/g, (ch) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[ch]));

  const renderChips = () => {
    const chip = (label, value) => {
      const b = document.createElement("button");
      b.type = "button";
      b.className = "search-chip" + (activeGroup === value ? " is-active" : "");
      b.textContent = label;
      b.addEventListener("click", () => {
        activeGroup = activeGroup === value ? null : value;
        render();
      });
      return b;
    };
    chipsBox.replaceChildren();
    chipsBox.appendChild(chip(t.all, null));
    for (const label of groupOrder) chipsBox.appendChild(chip(`${label} (${groupSeen.get(label)})`, label));
  };

  const score = (entry, query) => {
    const hay = [entry.t, entry.c || "", entry.d || "", entry.g].join(" ").toLowerCase();
    const idx = hay.indexOf(query);
    if (idx < 0) return -1;
    let rank = 100 - Math.min(idx, 60);
    if (entry.t.toLowerCase().includes(query)) rank += 40;
    if (!entry.hub) rank += 10;
    if (entry.p) rank += 5;
    return rank;
  };

  const render = () => {
    const query = input.value.trim().toLowerCase();
    let rows = entries.filter((e) => !activeGroup || e.g === activeGroup);
    if (query) {
      const words = query.split(/\s+/);
      rows = rows
        .map((e) => ({ e, r: words.reduce((min, w) => Math.min(min, score(e, w)), 1000) }))
        .filter((x) => x.r >= 0)
        .sort((a, b) => b.r - a.r)
        .map((x) => x.e);
    } else {
      rows = [...rows].sort((a, b) => groupOrder.indexOf(a.g) - groupOrder.indexOf(b.g) || a.t.localeCompare(b.t, lang));
    }

    const frag = document.createDocumentFragment();
    const limit = query ? 60 : 200;
    const shown = rows.slice(0, limit);
    for (const e of shown) {
      const li = document.createElement("li");
      const a = document.createElement("a");
      // entry path is /<lang>/<route>/; this page lives at /<lang>/search/
      a.href = `../${e.path.split("/").slice(2).join("/")}`;
      if (e.hub) {
        const badge = document.createElement("small");
        badge.textContent = t.hubLabel;
        badge.className = "search-hub-badge";
        a.appendChild(badge);
      }
      const title = document.createElement("strong");
      title.textContent = e.t;
      a.appendChild(title);
      const meta = document.createElement("span");
      meta.className = "search-meta";
      meta.textContent = [e.g, e.s, e.c].filter(Boolean).join(" · ");
      a.appendChild(meta);
      if (e.d) {
        const d = document.createElement("p");
        d.textContent = e.d;
        a.appendChild(d);
      }
      if (e.p) {
        const p = document.createElement("em");
        p.className = "search-playable";
        p.textContent = t.playableLabel;
        a.appendChild(p);
      }
      li.appendChild(a);
      frag.appendChild(li);
    }
    listBox.replaceChildren(frag);

    if (query || activeGroup) {
      status.textContent = t.results(rows.length);
      if (!rows.length) {
        const li = document.createElement("li");
        li.className = "search-empty";
        li.textContent = t.none;
        listBox.appendChild(li);
      }
    } else {
      status.textContent = t.empty;
    }
  };

  fetch("../../assets/search-index.json")
    .then((res) => res.json())
    .then((data) => {
      entries = data.filter((e) => e.path.startsWith(`/${lang}/`));
      for (const e of entries) {
        if (!groupSeen.has(e.g)) {
          groupSeen.set(e.g, 0);
          groupOrder.push(e.g);
        }
        groupSeen.set(e.g, groupSeen.get(e.g) + 1);
      }
      renderChips();
      render();
      input.addEventListener("input", render);
      const initial = new URLSearchParams(location.search).get("q");
      if (initial) {
        input.value = initial;
        render();
      }
    })
    .catch(() => {
      status.textContent = "index load failed";
    });
})();
