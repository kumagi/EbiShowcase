#!/usr/bin/env node
/**
 * Add the shared Google Form to every Japanese and English content page.
 *
 * The source pages keep the form markup so local previews work too.
 * Running this script again is safe because the marker is idempotent.
 */
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const formURL =
  "https://docs.google.com/forms/d/e/1FAIpQLSdE74SxJYstsQ2pckmG-IIGwgMMlpcp3w7c2bG-RPso-nQLbA/viewform?embedded=true";

function walk(dir, files = []) {
  for (const name of readdirSync(dir)) {
    const full = join(dir, name);
    const stat = statSync(full);
    if (stat.isDirectory()) walk(full, files);
    else if (name === "index.html") files.push(full);
  }
  return files;
}

function block(lang) {
  const japanese = lang === "ja";
  return `
    <section class="feedback-section" aria-labelledby="feedback-title">
      <div class="feedback-copy">
        <p class="eyebrow">${japanese ? "TELL EBI BOY" : "TELL EBI BOY"}</p>
        <h2 id="feedback-title">${japanese ? "ひとことフィードバック" : "A quick note for us"}</h2>
        <p>${japanese ? "わかりにくかったところ、楽しかったところ、追加してほしいことを教えてください。短いひとことでも大歓迎です。" : "Tell us what was fun, confusing, or worth adding. A short note is perfect."}</p>
      </div>
      <div class="feedback-frame">
        <iframe
          title="${japanese ? "Ebi Showcase フィードバックフォーム" : "Ebi Showcase feedback form"}"
          src="${formURL}"
          loading="lazy"
          referrerpolicy="no-referrer-when-downgrade"
        >${japanese ? "フィードバックフォームを読み込めませんでした。フォームを別のタブで開いてください。" : "The feedback form could not be loaded. Please open it in a new tab."}</iframe>
      </div>
    </section>
`;
}

let updated = 0;
const pages = [
  ...walk(join(root, "web", "ja")),
  ...walk(join(root, "web", "en")),
];

for (const file of pages) {
  const lang = file.includes(`${join("web", "ja")}/`) ? "ja" : "en";
  const html = readFileSync(file, "utf8");
  if (html.includes("class=\"feedback-section\"")) continue;
  const mainClose = html.lastIndexOf("</main>");
  if (mainClose < 0) {
    console.error(`No </main> found in ${file}`);
    process.exitCode = 1;
    continue;
  }
  const next = html.slice(0, mainClose) + block(lang) + html.slice(mainClose);
  writeFileSync(file, next);
  updated++;
}

console.log(`Inserted feedback forms into ${updated} page(s).`);
