/** Shared authoring panels for generated genre-track lessons.
 *
 * A page must never call an invented snippet "REAL GO" when its runnable
 * entry point is a thin Config/Run wrapper. Show that entry first, then the
 * small internal mechanism that explains it.
 */
function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
}

/**
 * @param {{
 *  lang: "ja"|"en", entryPath: string, entryCode: string,
 *  implementationPath: string, implementationCode: string,
 *  rule: { path: string, location: string, action: string, verify: string }
 * }} lesson
 */
export function dualLayerCodeLesson(lesson) {
  const ja = lesson.lang === "ja";
  const ruleLabel = ja ? "YOUR FIRST RULE" : "YOUR FIRST RULE";
  const entryLabel = ja ? "編集する入口" : "EDIT THIS ENTRY";
  const mechanismLabel = ja ? "仕組みの抜粋" : "HOW IT WORKS";
  const edit = ja
    ? `${lesson.rule.path} の ${lesson.rule.location} に ${lesson.rule.action}。${lesson.rule.verify}`
    : `In ${lesson.rule.path}, at ${lesson.rule.location}, ${lesson.rule.action}. ${lesson.rule.verify}`;
  return `<section class="code-lesson">
  <div><p class="eyebrow">${entryLabel}</p><h3>${escapeHTML(lesson.entryPath)}</h3><p>${ja ? "この短い入口が、実際に学習者が開いて変更する Go ファイルです。" : "This short entry is the real Go file the learner opens and changes."}</p></div>
  <pre><code>${escapeHTML(lesson.entryCode)}</code></pre>
</section>
<section class="code-lesson">
  <div><p class="eyebrow">${mechanismLabel}</p><h3>${escapeHTML(lesson.implementationPath)}</h3><p>${ja ? `上の ${escapeHTML(lesson.entryPath)} は開始設定を書く短いファイルです。移動・判定・得点などの処理はこのファイルにあり、次の抜粋で設定がどのルールへ渡るかを追えます。` : `The ${escapeHTML(lesson.entryPath)} file above only supplies the starting configuration. Movement, checks, scoring, and other rules live here; the excerpt below shows where that configuration is used.`}</p></div>
  <pre><code>${escapeHTML(lesson.implementationCode)}</code></pre>
</section>
<section class="why-grid"><article class="challenge"><p class="eyebrow">${ruleLabel}</p><h3>${ja ? "1つルールを書いて確かめよう" : "Write one rule, then verify it"}</h3><p>${escapeHTML(edit)}</p></article></section>`;
}

/** Three different concept cards: data shape, Update order, Draw mapping. */
export function authoringConceptRow(cards) {
  if (!Array.isArray(cards) || cards.length !== 3) throw new Error("authoringConceptRow needs exactly 3 cards");
  return `<div class="concept-row">${cards.map((card, index) => `<article><span class="concept-number">${index + 1}</span><h3>${escapeHTML(card.title)}</h3><p>${escapeHTML(card.body)}</p><code>${escapeHTML(card.code)}</code></article>`).join("")}</div>`;
}
