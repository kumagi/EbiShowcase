#!/usr/bin/env node
/**
 * Add the shared Google Form connection to every Japanese and English content page.
 *
 * The page uses a native, styled form that posts to Google's public formResponse
 * endpoint. The Google Form remains the owner of validation and spreadsheet storage.
 * Running this script again is safe because the marker is idempotent.
 */
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join, relative } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const responseURL =
  "https://docs.google.com/forms/d/e/1FAIpQLSdE74SxJYstsQ2pckmG-IIGwgMMlpcp3w7c2bG-RPso-nQLbA/formResponse";
const pageEntry = "765794446";
const feedbackEntry = "893595607";

function walk(dir, files = []) {
  for (const name of readdirSync(dir)) {
    const full = join(dir, name);
    const stat = statSync(full);
    if (stat.isDirectory()) walk(full, files);
    else if (name === "index.html") files.push(full);
  }
  return files;
}

function escapeHTML(value) {
  return value.replace(/[&<>\"]/g, (character) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[character]));
}

function block(lang, pagePath) {
  const japanese = lang === "ja";
  const pageLabel = escapeHTML(`/${pagePath.replace(/\\/g, "/")}`);
  return `
    <section class="feedback-section" aria-labelledby="feedback-title">
      <div class="feedback-card">
        <div class="feedback-heading">
          <p class="eyebrow">FEEDBACK</p>
          <h2 id="feedback-title">${japanese ? "ひとことフィードバック" : "Quick feedback"}</h2>
        </div>
        <form class="feedback-form" action="${responseURL}" method="POST">
          <input type="hidden" name="entry.${pageEntry}" value="${pageLabel}">
          <label class="feedback-field">
            <span class="sr-only">${japanese ? "フィードバック" : "Feedback"}</span>
            <input class="feedback-message" name="entry.${feedbackEntry}" maxlength="200" required data-sending="${japanese ? "送信中…" : "Sending…"}" data-sent="${japanese ? "送信しました。ありがとうございます！" : "Sent — thank you!"}" data-failed="${japanese ? "送信できませんでした。時間をおいて再試行してください。" : "Could not send. Please try again later."}" placeholder="${japanese ? "ひとこと入力…" : "Write one short note…"}">
          </label>
          <div class="feedback-actions">
            <button type="submit" class="feedback-submit">${japanese ? "送信する" : "Send feedback"}<span>→</span></button>
            <p class="feedback-status" aria-live="polite"></p>
          </div>
          <input type="hidden" name="fvv" value="1">
          <input type="hidden" name="pageHistory" value="0">
        </form>
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
  const pagePath = relative(join(root, "web", lang), file).replace(/\/index\.html$/, "/");
  const html = readFileSync(file, "utf8");
  const marker = '<section class="feedback-section"';
  const existing = html.indexOf(marker);
  if (existing >= 0) {
    const close = html.indexOf("</section>", existing);
    if (close < 0) throw new Error(`Malformed feedback section: ${file}`);
    const prefix = html.slice(0, existing).replace(/[ \t\n]+$/, "\n\n");
    const next = prefix + block(lang, `${lang}/${pagePath}`) + html.slice(close + "</section>".length);
    writeFileSync(file, next);
    updated++;
    continue;
  }
  const mainClose = html.lastIndexOf("</main>");
  if (mainClose < 0) {
    console.error(`No </main> found in ${file}`);
    process.exitCode = 1;
    continue;
  }
  const prefix = html.slice(0, mainClose).replace(/[ \t\n]+$/, "\n\n");
  writeFileSync(file, prefix + block(lang, `${lang}/${pagePath}`) + html.slice(mainClose));
  updated++;
}

console.log(`Updated feedback forms in ${updated} page(s).`);
