#!/usr/bin/env node
/**
 * Read and triage Ebi Showcase feedback in Google Sheets.
 *
 * Authentication:
 *   Put an installed-app OAuth client JSON under .secrets/.
 *   The first command opens a browser for consent and stores a refresh token.
 *   Service-account JSON is also supported via GOOGLE_APPLICATION_CREDENTIALS.
 *
 * Commands:
 *   node scripts/feedback-sheet.mjs list
 *   node scripts/feedback-sheet.mjs pending
 *   node scripts/feedback-sheet.mjs check 12
 *   node scripts/feedback-sheet.mjs check-range 12 20
 *   node scripts/feedback-sheet.mjs delete 12
 */
import { createSign, randomBytes } from "node:crypto";
import { createServer } from "node:http";
import { execFile } from "node:child_process";
import { existsSync, readdirSync, readFileSync, writeFileSync } from "node:fs";

const spreadsheetId = "1r6jYssPE7AdluEqJ1nqyzzRWrOH8Ncp-zUyK4xnn0lw";
const spreadsheetURL = `https://sheets.googleapis.com/v4/spreadsheets/${spreadsheetId}`;
const credentialsPath = process.env.GOOGLE_APPLICATION_CREDENTIALS || ".secrets/feedback-service-account.json";
const oauthClientPath = process.env.GOOGLE_OAUTH_CLIENT_SECRET;
const oauthTokenPath = ".secrets/feedback-token.json";
const oauthRedirect = "http://localhost:53682/oauth2callback";
const statusHeaders = ["対応済み", "処理済み", "対応", "done", "status"];
const command = process.argv[2] || "list";
const rowArgument = Number(process.argv[3]);
const endRowArgument = Number(process.argv[4]);

const base64url = (value) => Buffer.from(value).toString("base64url");
const json = (value) => JSON.stringify(value);

function makeJWT(credentials) {
  const now = Math.floor(Date.now() / 1000);
  const header = base64url(json({ alg: "RS256", typ: "JWT" }));
  const claims = base64url(json({
    iss: credentials.client_email,
    scope: "https://www.googleapis.com/auth/spreadsheets",
    aud: "https://oauth2.googleapis.com/token",
    iat: now,
    exp: now + 3600,
  }));
  const unsigned = `${header}.${claims}`;
  const signer = createSign("RSA-SHA256");
  signer.update(unsigned);
  return `${unsigned}.${signer.sign(credentials.private_key, "base64url")}`;
}

function findOAuthClient() {
  const candidates = oauthClientPath ? [oauthClientPath] : (existsSync(".secrets") ? readdirSync(".secrets").filter((name) => name.endsWith(".json")).map((name) => `.secrets/${name}`) : []);
  for (const path of candidates) {
    try {
      const data = JSON.parse(readFileSync(path, "utf8"));
      if (data.installed || data.web) return { path, data: data.installed || data.web };
    } catch { /* Ignore non-OAuth JSON files. */ }
  }
  return null;
}

async function exchangeToken(body) {
  const response = await fetch("https://oauth2.googleapis.com/token", {
    method: "POST",
    headers: { "content-type": "application/x-www-form-urlencoded" },
    body: new URLSearchParams(body),
  });
  if (!response.ok) throw new Error(`Google token error ${response.status}: ${await response.text()}`);
  return response.json();
}

async function authorizeWithBrowser(client) {
  const state = randomBytes(20).toString("hex");
  const authURL = new URL(client.auth_uri);
  authURL.search = new URLSearchParams({
    client_id: client.client_id,
    redirect_uri: oauthRedirect,
    response_type: "code",
    access_type: "offline",
    prompt: "consent",
    scope: "https://www.googleapis.com/auth/spreadsheets",
    state,
  });

  const code = await new Promise((resolve, reject) => {
    const server = createServer((request, response) => {
      const url = new URL(request.url, oauthRedirect);
      if (url.pathname !== "/oauth2callback") { response.writeHead(404); response.end(); return; }
      if (url.searchParams.get("state") !== state) { response.writeHead(400); response.end("State mismatch"); server.close(); reject(new Error("OAuth state mismatch")); return; }
      const error = url.searchParams.get("error");
      if (error) { response.writeHead(400); response.end(`OAuth error: ${error}`); server.close(); reject(new Error(`OAuth error: ${error}`)); return; }
      response.writeHead(200, { "content-type": "text/html; charset=utf-8" });
      response.end("<h1>認証できました</h1><p>このタブを閉じてターミナルへ戻ってください。</p>");
      server.close();
      resolve(url.searchParams.get("code"));
    });
    server.on("error", reject);
    server.listen(53682, "127.0.0.1", () => {
      console.log("ブラウザでGoogleにログインし、スプレッドシートへのアクセスを許可してください。");
      console.log(authURL.toString());
      execFile("open", [authURL.toString()], () => {});
    });
  });

  const token = await exchangeToken({
    code,
    client_id: client.client_id,
    client_secret: client.client_secret,
    redirect_uri: oauthRedirect,
    grant_type: "authorization_code",
  });
  writeFileSync(oauthTokenPath, JSON.stringify(token, null, 2));
  return token.access_token;
}

async function getAccessToken() {
  if (process.env.GOOGLE_APPLICATION_CREDENTIALS || existsSync(credentialsPath)) {
    let credentials;
    try { credentials = JSON.parse(readFileSync(credentialsPath, "utf8")); }
    catch { throw new Error(`サービスアカウントJSONが読めません: ${credentialsPath}`); }
    const now = Math.floor(Date.now() / 1000);
    const header = base64url(json({ alg: "RS256", typ: "JWT" }));
    const claims = base64url(json({ iss: credentials.client_email, scope: "https://www.googleapis.com/auth/spreadsheets", aud: "https://oauth2.googleapis.com/token", iat: now, exp: now + 3600 }));
    const unsigned = `${header}.${claims}`;
    const signer = createSign("RSA-SHA256");
    signer.update(unsigned);
    const assertion = `${unsigned}.${signer.sign(credentials.private_key, "base64url")}`;
    const token = await exchangeToken({ grant_type: "urn:ietf:params:oauth:grant-type:jwt-bearer", assertion });
    return token.access_token;
  }

  const found = findOAuthClient();
  if (!found) throw new Error("OAuthクライアントJSONが見つかりません。.secrets に保存してください。");
  if (existsSync(oauthTokenPath)) {
    const token = JSON.parse(readFileSync(oauthTokenPath, "utf8"));
    if (token.refresh_token) {
      const refreshed = await exchangeToken({ client_id: found.data.client_id, client_secret: found.data.client_secret, refresh_token: token.refresh_token, grant_type: "refresh_token" });
      return refreshed.access_token;
    }
  }
  return authorizeWithBrowser(found.data);
}
async function google(path, options = {}) {
  const token = await getAccessToken();
  const response = await fetch(path, {
    ...options,
    headers: { authorization: `Bearer ${token}`, "content-type": "application/json", ...options.headers },
  });
  if (!response.ok) throw new Error(`Google Sheets API ${response.status}: ${await response.text()}`);
  return response.status === 204 ? null : response.json();
}

async function getSheet() {
  const metadata = await google(`${spreadsheetURL}?fields=sheets(properties(sheetId,title,index))`);
  const sheet = metadata.sheets?.find(({ properties }) => properties.sheetId === 1378055629) || metadata.sheets?.[0];
  if (!sheet) throw new Error("スプレッドシートにシートがありません。");
  return sheet.properties;
}

const columnName = (index) => {
  let name = "";
  for (let n = index + 1; n; n = Math.floor((n - 1) / 26)) name = String.fromCharCode(65 + ((n - 1) % 26)) + name;
  return name;
};

async function getRows(sheet) {
  const range = encodeURIComponent(`${sheet.title}!A:Z`);
  const result = await google(`${spreadsheetURL}/values/${range}?majorDimension=ROWS`);
  return result.values || [];
}

async function updateValues(sheet, range, values) {
  const encoded = encodeURIComponent(range);
  await google(`${spreadsheetURL}/values/${encoded}?valueInputOption=USER_ENTERED`, {
    method: "PUT",
    body: JSON.stringify({ range, majorDimension: "ROWS", values }),
  });
}

async function list(sheet, rows) {
  if (!rows.length) {
    console.log("フィードバックはまだありません。");
    return;
  }
  const headers = rows[0];
  console.log(`Sheet: ${sheet.title}`);
  rows.slice(1).forEach((row, index) => {
    const values = headers.map((header, column) => `${header}=${row[column] || ""}`).join(" | ");
    console.log(`${index + 2}: ${values}`);
  });
}

async function listPending(sheet, rows) {
  if (!rows.length) return;
  const headers = rows[0];
  const statusColumn = headers.findIndex((header) => statusHeaders.includes(String(header).trim().toLowerCase()));
  console.log(`Sheet: ${sheet.title}`);
  rows.slice(1).forEach((row, index) => {
    if (statusColumn >= 0 && String(row[statusColumn] || "").trim()) return;
    const values = headers.map((header, column) => `${header}=${String(row[column] || "").replaceAll("\n", "\\n")}`).join(" | ");
    console.log(`${index + 2}: ${values}`);
  });
}

async function findStatusColumn(sheet, rows) {
  const headers = rows[0] || [];
  let index = headers.findIndex((header) => statusHeaders.includes(String(header).trim().toLowerCase()));
  if (index >= 0) return index;
  index = headers.length;
  await updateValues(sheet, `${sheet.title}!${columnName(index)}1`, [["対応済み"]]);
  return index;
}

async function deleteRow(sheet, rowNumber) {
  const token = await getAccessToken();
  const response = await fetch(`${spreadsheetURL}:batchUpdate`, {
    method: "POST",
    headers: { authorization: `Bearer ${token}`, "content-type": "application/json" },
    body: JSON.stringify({ requests: [{ deleteDimension: { range: {
      sheetId: sheet.sheetId,
      dimension: "ROWS",
      startIndex: rowNumber - 1,
      endIndex: rowNumber,
    } } }] }),
  });
  if (!response.ok) throw new Error(`Google Sheets API ${response.status}: ${await response.text()}`);
}

try {
  const sheet = await getSheet();
  const rows = await getRows(sheet);
  if (command === "list") await list(sheet, rows);
  else if (command === "pending") await listPending(sheet, rows);
  else if (command === "check-range") {
    if (!Number.isInteger(rowArgument) || !Number.isInteger(endRowArgument) || rowArgument < 2 || endRowArgument < rowArgument) throw new Error("開始行と終了行を指定してください（例: check-range 12 20）。");
    const column = await findStatusColumn(sheet, rows);
    const values = Array.from({ length: endRowArgument - rowArgument + 1 }, () => ["✅"]);
    await updateValues(sheet, `${sheet.title}!${columnName(column)}${rowArgument}:${columnName(column)}${endRowArgument}`, values);
    console.log(`対応済みにしました: ${sheet.title}!${columnName(column)}${rowArgument}:${columnName(column)}${endRowArgument}`);
  }
  else if (!Number.isInteger(rowArgument) || rowArgument < 2) throw new Error("行番号を指定してください（例: check 12）。1行目は見出しです。");
  else if (command === "check") {
    const column = await findStatusColumn(sheet, rows);
    await updateValues(sheet, `${sheet.title}!${columnName(column)}${rowArgument}`, [["✅"]]);
    console.log(`対応済みにしました: ${sheet.title}!${columnName(column)}${rowArgument}`);
  } else if (command === "delete") {
    await deleteRow(sheet, rowArgument);
    console.log(`削除しました: ${sheet.title} row ${rowArgument}`);
  } else throw new Error(`不明なコマンド: ${command}（list / pending / check / check-range / delete）`);
} catch (error) {
  console.error(error.message);
  process.exit(1);
}
