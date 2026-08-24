(() => {
  "use strict";

  const STORAGE_KEY = "ebi-showcase-progress-v1";

  const langCutOf = (path) => {
    const ja = path.indexOf("/ja/");
    const en = path.indexOf("/en/");
    return ja >= 0 ? ja : en;
  };

  const lessonIdOfPath = (path) => {
    const cut = langCutOf(path);
    if (cut < 0) return null;
    let id = path.slice(cut);
    if (id.endsWith("index.html")) id = id.slice(0, -"index.html".length);
    return id || null;
  };

  const lessonId = () => lessonIdOfPath(location.pathname);

  const loadData = () => {
    try {
      const raw = JSON.parse(localStorage.getItem(STORAGE_KEY) || "{}");
      return { done: raw.done && typeof raw.done === "object" ? raw.done : {}, last: raw.last || null };
    } catch {
      return { done: {}, last: null };
    }
  };

  const saveData = (data) => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
    } catch {
      /* private mode etc. — progress simply is not persisted */
    }
  };

  const data = loadData();

  const siteBase = () => location.pathname.slice(0, Math.max(langCutOf(location.pathname), 0));

  const hrefToId = (href) => {
    try {
      const url = new URL(href, location.href);
      return lessonIdOfPath(url.pathname);
    } catch {
      return null;
    }
  };

  const isJa = (document.documentElement.lang || "ja").toLowerCase().startsWith("ja");
  const t = {
    resume: isJa ? "つづきから" : "Continue",
    playedCount: (n) => (isJa ? `${n} 本のレッスンをプレイ済み` : `${n} lesson${n === 1 ? "" : "s"} played`),
    doneLabel: isJa ? "プレイ済み" : "Played",
  };

  /* ---------- record this page ---------- */

  const playableFrame =
    document.querySelector("iframe[data-game-src]") ||
    document.querySelector("iframe[src*='/play/']");

  if (playableFrame) {
    const id = lessonId();
    if (id) {
      const title = (document.querySelector("h1") || {}).textContent || document.title;
      data.last = { id, title: title.trim().replace(/\s+/g, " ").slice(0, 60), ts: Date.now() };
      saveData(data);
    }
  }

  const markPlayed = () => {
    const id = lessonId();
    if (!id || data.done[id]) return;
    data.done[id] = Date.now();
    saveData(data);
  };

  /* ---------- render badges + resume card ---------- */

  const CARD_SELECTOR =
    "a.course-card[href], a.vfx-course-card[href], a.test-course-card[href], a.path-step[href], a.track-card[href]";

  const doneBadge = () => {
    const badge = document.createElement("span");
    badge.className = "prog-done";
    badge.setAttribute("role", "img");
    badge.setAttribute("aria-label", t.doneLabel);
    badge.textContent = "✓";
    return badge;
  };

  const renderBadges = () => {
    document.querySelectorAll(CARD_SELECTOR).forEach((card) => {
      const id = hrefToId(card.href);
      if (!id) return;
      if (card.classList.contains("track-card")) {
        const prefix = id.endsWith("/") ? id : `${id}/`;
        let count = 0;
        for (const key of Object.keys(data.done)) {
          if (key.startsWith(prefix)) count += 1;
        }
        if (count > 0) {
          const chip = document.createElement("span");
          chip.className = "prog-count";
          chip.textContent = String(count);
          chip.title = t.playedCount(count);
          const host = card.querySelector("div") || card;
          host.appendChild(chip);
        }
        return;
      }
      if (data.done[id]) {
        card.classList.add("is-played");
        if (card.classList.contains("path-step")) {
          const arrow = card.querySelector("b");
          if (arrow) arrow.textContent = "✓";
        } else {
          card.appendChild(doneBadge());
        }
      }
    });
  };

  const renderResume = () => {
    const slot = document.querySelector("[data-progress-continue]");
    if (!slot) return;
    const played = Object.keys(data.done).length;
    if (!played) return;
    const base = siteBase();
    const resume = document.createElement("a");
    resume.className = "prog-resume";
    resume.href = data.last ? `${base}${data.last.id}` : "#curriculum";
    const label = document.createElement("strong");
    label.textContent = t.resume;
    resume.appendChild(label);
    if (data.last && data.last.title) {
      const where = document.createElement("small");
      where.textContent = data.last.title;
      resume.appendChild(where);
    }
    const count = document.createElement("span");
    count.className = "prog-resume-count";
    count.textContent = t.playedCount(played);
    slot.hidden = false;
    slot.appendChild(resume);
    slot.appendChild(count);
  };

  renderBadges();
  renderResume();

  /* ---------- deferred WASM frame loading (existing behaviour) ---------- */

  const frames = document.querySelectorAll("iframe[data-game-src]");
  if (frames.length === 0) return;

  const load = (frame) => {
    const src = frame.dataset.gameSrc;
    if (!src) return;
    frame.addEventListener(
      "load",
      () => {
        frame.removeAttribute("aria-busy");
        markPlayed();
      },
      { once: true },
    );
    frame.src = src;
    frame.removeAttribute("data-game-src");
  };

  if (!("IntersectionObserver" in window)) {
    frames.forEach(load);
    return;
  }

  const observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue;
        observer.unobserve(entry.target);
        load(entry.target);
      }
    },
    { rootMargin: "200px 0px" },
  );

  frames.forEach((frame) => {
    frame.setAttribute("aria-busy", "true");
    observer.observe(frame);
  });
})();
