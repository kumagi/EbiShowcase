#!/usr/bin/env node
// SPDX-License-Identifier: Apache-2.0

import { spawn } from "node:child_process";
import { existsSync, mkdirSync, mkdtempSync, readFileSync, readdirSync, rmSync, statSync, writeFileSync } from "node:fs";
import { createServer } from "node:http";
import { tmpdir } from "node:os";
import { extname, join, relative, resolve, sep } from "node:path";
import { createServer as createNetServer } from "node:net";

const root = new URL("..", import.meta.url).pathname;
const webRoot = join(root, "web");
const args = process.argv.slice(2);
const json = args.includes("--json");
const routeArg = valueAfter("--routes");
const routePattern = routeArg ? new RegExp(routeArg) : null;
const viewportWidth = Number(valueAfter("--width") || 390);
const viewportHeight = Number(valueAfter("--height") || 844);
const screenshotDir = valueAfter("--screenshot-dir");
if (!Number.isInteger(viewportWidth) || viewportWidth < 280 || !Number.isInteger(viewportHeight) || viewportHeight < 480) {
  throw new Error("--width and --height must be integer phone viewport dimensions.");
}

function valueAfter(flag) {
  const index = args.indexOf(flag);
  return index >= 0 ? args[index + 1] : "";
}

function walkHTML(dir, out = []) {
  for (const name of readdirSync(dir)) {
    if (name === "play") continue;
    const path = join(dir, name);
    if (statSync(path).isDirectory()) walkHTML(path, out);
    else if (name.endsWith(".html")) out.push(path);
  }
  return out;
}

function pageURL(file) {
  const local = relative(webRoot, file).split(sep).join("/");
  if (local === "index.html") return "/";
  if (local.endsWith("/index.html")) return `/${local.slice(0, -"index.html".length)}`;
  return `/${local}`;
}

function chromePath() {
  const candidates = [
    process.env.CHROME_BIN,
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
    "/Applications/Chromium.app/Contents/MacOS/Chromium",
    "/usr/bin/google-chrome",
    "/usr/bin/google-chrome-stable",
    "/usr/bin/chromium",
    "/usr/bin/chromium-browser",
  ].filter(Boolean);
  const found = candidates.find(existsSync);
  if (!found) throw new Error("Chrome/Chromium was not found. Set CHROME_BIN to its executable.");
  return found;
}

async function freePort() {
  const server = createNetServer();
  await new Promise((resolveListen, reject) => server.once("error", reject).listen(0, "127.0.0.1", resolveListen));
  const { port } = server.address();
  await new Promise((resolveClose) => server.close(resolveClose));
  return port;
}

function mime(path) {
  return ({
    ".css": "text/css; charset=utf-8",
    ".html": "text/html; charset=utf-8",
    ".js": "text/javascript; charset=utf-8",
    ".mjs": "text/javascript; charset=utf-8",
    ".json": "application/json; charset=utf-8",
    ".png": "image/png",
    ".webp": "image/webp",
    ".svg": "image/svg+xml",
    ".wasm": "application/wasm",
  })[extname(path).toLowerCase()] || "application/octet-stream";
}

async function startWebServer() {
  const server = createServer((request, response) => {
    let pathname;
    try { pathname = decodeURIComponent(new URL(request.url, "http://localhost").pathname); } catch { response.writeHead(400).end(); return; }
    let file = resolve(webRoot, `.${pathname}`);
    if (!file.startsWith(`${webRoot}${sep}`) && file !== webRoot) { response.writeHead(403).end(); return; }
    if (existsSync(file) && statSync(file).isDirectory()) file = join(file, "index.html");
    if (!existsSync(file) || statSync(file).isDirectory()) { response.writeHead(404).end(); return; }
    response.writeHead(200, { "Content-Type": mime(file), "Cache-Control": "no-store" });
    response.end(readFileSync(file));
  });
  await new Promise((resolveListen, reject) => server.once("error", reject).listen(0, "127.0.0.1", resolveListen));
  return server;
}

async function waitForDebugger(port) {
  const endpoint = `http://127.0.0.1:${port}/json/list`;
  for (let attempt = 0; attempt < 100; attempt++) {
    try {
      const pages = await (await fetch(endpoint)).json();
      const page = pages.find((candidate) => candidate.type === "page");
      if (page) return page.webSocketDebuggerUrl;
    } catch {}
    await new Promise((resolveWait) => setTimeout(resolveWait, 50));
  }
  throw new Error("Chrome DevTools did not become ready.");
}

class CDP {
  constructor(url) {
    this.socket = new WebSocket(url);
    this.nextID = 0;
    this.pending = new Map();
    this.waiters = new Map();
    this.socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.id && this.pending.has(message.id)) {
        this.pending.get(message.id)(message);
        this.pending.delete(message.id);
      }
      const queue = this.waiters.get(message.method);
      if (queue?.length) queue.shift()(message.params);
    };
  }

  async open() {
    await new Promise((resolveOpen, reject) => {
      this.socket.onopen = resolveOpen;
      this.socket.onerror = reject;
    });
  }

  send(method, params = {}) {
    return new Promise((resolveMessage) => {
      const id = ++this.nextID;
      this.pending.set(id, resolveMessage);
      this.socket.send(JSON.stringify({ id, method, params }));
    });
  }

  event(method, timeout = 8000) {
    return new Promise((resolveEvent, reject) => {
      const queue = this.waiters.get(method) || [];
      const timer = setTimeout(() => reject(new Error(`Timed out waiting for ${method}`)), timeout);
      queue.push((value) => { clearTimeout(timer); resolveEvent(value); });
      this.waiters.set(method, queue);
    });
  }

  close() { this.socket.close(); }
}

const auditExpression = `new Promise(async resolveAudit => {
  if (document.fonts && document.fonts.ready) await document.fonts.ready;
  await new Promise(done => requestAnimationFrame(() => requestAnimationFrame(done)));
  const viewport = document.documentElement.clientWidth;
  const visible = element => {
    if (element.matches(".sr-only,[aria-hidden='true'],input[type='file']")) return false;
    const style = getComputedStyle(element);
    const rect = element.getBoundingClientRect();
    return style.display !== "none" && style.visibility !== "hidden" && rect.width > 0 && rect.height > 0;
  };
  const label = element => {
    const id = element.id ? "#" + element.id : "";
    const classes = typeof element.className === "string" && element.className.trim()
      ? "." + element.className.trim().split(/\\s+/).slice(0, 3).join(".")
      : "";
    const own = (element.tagName || "node").toLowerCase() + id + classes;
    const section = element.closest && element.closest("section");
    if (!section || section === element) return own;
    const sectionClass = typeof section.className === "string" && section.className.trim()
      ? "." + section.className.trim().split(/\\s+/).slice(0, 2).join(".")
      : "";
    return own + " in section" + sectionClass;
  };
  const outside = [];
  const scrollers = [];
  const innerOverflow = [];
  for (const element of document.querySelectorAll("body *")) {
    if (!visible(element)) continue;
    const rect = element.getBoundingClientRect();
    const style = getComputedStyle(element);
    if ((rect.right > viewport + 1 || rect.left < -1) && style.position !== "fixed") {
      outside.push({ element: label(element), left: Math.round(rect.left), right: Math.round(rect.right), width: Math.round(rect.width) });
    }
    if (element.scrollWidth > element.clientWidth + 1 && /auto|scroll/.test(style.overflowX)) {
      scrollers.push({ element: label(element), client: element.clientWidth, scroll: element.scrollWidth });
    } else if (
      element.scrollWidth > element.clientWidth + 1
      && /hidden|clip/.test(style.overflowX)
      && style.textOverflow !== "ellipsis"
    ) {
      innerOverflow.push({ element: label(element), client: element.clientWidth, scroll: element.scrollWidth, overflow: style.overflowX });
    }
  }
  resolveAudit({
    title: document.title,
    viewport,
    pageWidth: Math.max(document.documentElement.scrollWidth, document.body?.scrollWidth || 0),
    outside: outside.slice(0, 12),
    scrollers: scrollers.slice(0, 12),
    innerOverflow: innerOverflow.sort((a, b) => b.scroll - b.client - (a.scroll - a.client)).slice(0, 12),
  });
})`;

const files = walkHTML(webRoot).filter((file) => !routePattern || routePattern.test(pageURL(file)));
const routes = files.map(pageURL).sort();
const server = await startWebServer();
const webPort = server.address().port;
const debugPort = await freePort();
const profile = mkdtempSync(join(tmpdir(), "ebi-mobile-audit-"));
const chrome = spawn(chromePath(), [
  "--headless=new",
  "--disable-extensions",
  "--disable-background-networking",
  "--no-first-run",
  `--remote-debugging-port=${debugPort}`,
  `--user-data-dir=${profile}`,
  `--window-size=${viewportWidth},${viewportHeight}`,
  "--force-device-scale-factor=1",
  "about:blank",
], { stdio: "ignore" });

const findings = [];
if (screenshotDir) mkdirSync(screenshotDir, { recursive: true });
try {
  const debuggerURL = await waitForDebugger(debugPort);
  const cdp = new CDP(debuggerURL);
  await cdp.open();
  await cdp.send("Page.enable");
  await cdp.send("Runtime.enable");
  await cdp.send("Emulation.setDeviceMetricsOverride", { width: viewportWidth, height: viewportHeight, deviceScaleFactor: 1, mobile: true });
  for (let index = 0; index < routes.length; index++) {
    const route = routes[index];
    const loaded = cdp.event("Page.loadEventFired");
    const navigation = await cdp.send("Page.navigate", { url: `http://127.0.0.1:${webPort}${route}` });
    if (navigation.result?.errorText) throw new Error(`${route}: ${navigation.result.errorText}`);
    await loaded;
    const evaluated = await cdp.send("Runtime.evaluate", { expression: auditExpression, awaitPromise: true, returnByValue: true });
    const result = evaluated.result?.result?.value;
    if (!result) throw new Error(`${route}: layout evaluation failed`);
    if (screenshotDir) {
      const captured = await cdp.send("Page.captureScreenshot", { format: "png", captureBeyondViewport: false });
      const filename = route === "/" ? "index.png" : `${route.replace(/^\/|\/$/g, "").replaceAll("/", "__")}.png`;
      writeFileSync(join(screenshotDir, filename), Buffer.from(captured.result.data, "base64"));
    }
    if (result.pageWidth > result.viewport + 1 || result.outside.length || result.scrollers.length || result.innerOverflow.length) {
      findings.push({ route, ...result });
    }
    if (!json && ((index + 1) % 25 === 0 || index + 1 === routes.length)) console.log(`Audited ${index + 1}/${routes.length} pages; ${findings.length} need fixes.`);
  }
  cdp.close();
} finally {
  chrome.kill("SIGTERM");
  await Promise.race([
    new Promise((resolveExit) => chrome.once("exit", resolveExit)),
    new Promise((resolveWait) => setTimeout(resolveWait, 3000)),
  ]);
  await new Promise((resolveClose) => server.close(resolveClose));
  rmSync(profile, { recursive: true, force: true, maxRetries: 5, retryDelay: 100 });
}

if (json) console.log(JSON.stringify({ pages: routes.length, findings }, null, 2));
else if (findings.length) {
  console.log(`FAIL — ${findings.length}/${routes.length} pages require horizontal movement at ${viewportWidth}px:`);
  for (const finding of findings) {
    console.log(`  ${finding.route} (${finding.pageWidth}px page)`);
    for (const item of finding.outside) console.log(`    outside: ${item.element} [${item.left}, ${item.right}] width=${item.width}`);
    for (const item of finding.scrollers) console.log(`    scroll: ${item.element} ${item.client}→${item.scroll}`);
    for (const item of finding.innerOverflow) console.log(`    inner: ${item.element} ${item.client}→${item.scroll} (${item.overflow})`);
  }
} else console.log(`OK — ${routes.length}/${routes.length} pages fit ${viewportWidth}px without horizontal scrolling.`);

process.exit(findings.length ? 1 : 0);
