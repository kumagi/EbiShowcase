/**
 * Generates bilingual Getting Started / environment setup guides.
 * Run: node scripts/gen-setup-guide.mjs
 */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");

const blankMainGo = `package main

import (
\t"image/color"
\t"log"

\t"github.com/hajimehoshi/ebiten/v2"
)

// Game は Ebitengine が求める3つのメソッドを持つ箱です。
type Game struct{}

// Update は数字を進める場所。今は何もしません。
func (g *Game) Update() error {
\treturn nil
}

// Draw は画面を塗る場所。いまは暗い青一色だけ。
func (g *Game) Draw(screen *ebiten.Image) {
\tscreen.Fill(color.RGBA{20, 28, 48, 255})
}

// Layout はゲーム内部の解像度（幅×高さ）を返します。
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
\treturn 640, 480
}

func main() {
\tebiten.SetWindowSize(640, 480)
\tebiten.SetWindowTitle("Empty Window")
\tif err := ebiten.RunGame(&Game{}); err != nil {
\t\tlog.Fatal(err)
\t}
}`;

const blankMainGoEN = `package main

import (
\t"image/color"
\t"log"

\t"github.com/hajimehoshi/ebiten/v2"
)

// Game holds the three methods Ebitengine expects.
type Game struct{}

// Update advances numbers. For now it does nothing.
func (g *Game) Update() error {
\treturn nil
}

// Draw paints the screen. Just a solid dark blue.
func (g *Game) Draw(screen *ebiten.Image) {
\tscreen.Fill(color.RGBA{20, 28, 48, 255})
}

// Layout returns the game's internal resolution.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
\treturn 640, 480
}

func main() {
\tebiten.SetWindowSize(640, 480)
\tebiten.SetWindowTitle("Empty Window")
\tif err := ebiten.RunGame(&Game{}); err != nil {
\t\tlog.Fatal(err)
\t}
}`;

function esc(s) {
  return String(s)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}

function codeBlock(code, opts = {}) {
  const body = `<pre class="setup-code"><code>${esc(code)}</code></pre>`;
  if (!opts.copy) return body;
  return `<div class="full-code setup-full-code">
  <div class="full-code-head">
    <span>${esc(opts.filename || "main.go")}</span>
    <button type="button" class="full-code-copy" data-copy data-copied-label="${esc(opts.copied || "Copied!")}">${esc(opts.copy)}</button>
  </div>
  ${body}
</div>`;
}

function page(lang) {
  const ja = lang === "ja";
  const other = ja ? "en" : "ja";
  const otherLabel = ja ? "EN" : "日本語";
  const home = "../../";
  const css = "../../../style.css";
  const learn = "../../../learn.js";
  const otherHref = `../../../${other}/guides/setup/`;
  const route = `/${lang}/guides/setup/`;

  const t = ja
    ? {
        title: "はじめての環境づくり — 空の窓を開くまで | Ebi Showcase",
        desc: "まっさらな Windows / macOS から Go と Ebitengine を入れ、何もしない空の窓を起動するまでの手順。",
        crumb: "SPECIAL GUIDE",
        eyebrow: "START HERE / ENVIRONMENT",
        h1: "まっさらなパソコンから、<br><em>中身のないゲーム画面ウィンドウ</em>を開くまで。",
        lead: "このサイトのレッスンはブラウザでも遊べます。でも自分の PC でコードを動かしたいときは、まず Go と Ebitengine の準備が必要です。Windows と Mac、どちらの手順も書きます。最後に「何もしない」暗い青い窓が開けば成功です。",
        goalEyebrow: "ゴール",
        goalH: "このページが終わると、できること",
        goals: [
          ["ターミナル（黒い文字の画面）を開ける", "命令を打ち込む場所です"],
          ["go version が数字を返す", "Go が入った証拠です"],
          ["空の Ebitengine 窓が開く", "ゲームの土台が動いた証拠です"],
        ],
        tocEyebrow: "もくじ",
        toc: [
          ["#terminal", "0. ターミナルを開く"],
          ["#install-go", "1. Go を入れる"],
          ["#compiler", "2. Mac だけ：C コンパイラ"],
          ["#check", "3. 動作確認"],
          ["#project", "4. フォルダと go.mod"],
          ["#code", "5. main.go を書く"],
          ["#run", "6. 実行する"],
          ["#troubleshoot", "うまくいかないとき"],
        ],
        step0: {
          id: "terminal",
          n: "00",
          h: "命令を打ち込む場所を開く",
          p: "ブラウザではなく、パソコンに命令を出す画面を使います。名前は OS で違います。",
          winTitle: "Windows",
          winSteps: [
            "キーボードの Windows キーを押す",
            "「powershell」と入力する",
            "「Windows PowerShell」または「ターミナル」を開く",
          ],
          macTitle: "macOS",
          macSteps: [
            "Command + Space キーを同時に押す",
            "「ターミナル」と入力する",
            "「ターミナル」アプリを開く",
          ],
          tip: "文字を入力する画面が開き、カーソルが点滅していれば OK。Windows は黒いことが多く、Mac は白いこともありますが、役割は同じです。これからここにコマンドをコピーして貼り付けます。",
          note: "用語メモ：Windows では「PowerShell」、Mac では「ターミナル」と呼ぶことが多いです。どちらも<strong>同じ役割</strong>——キーボードでパソコンに命令を出す黒い（または白い）窓です。「ターミナル＝Mac専用」ではありません。Windows 11 の「ターミナル」アプリでも PowerShell が使えます。このガイドでは Windows 欄を PowerShell、Mac 欄をターミナルと書きます。",
        },
        step1: {
          id: "install-go",
          n: "01",
          h: "Go をインストールする",
          p: "Ebitengine は Go で書かれたゲームエンジンです。先に Go 本体が必要です。このリポジトリは Go 1.25 以降を使います。",
          winTitle: "Windows",
          winBody: [
            'ブラウザで <a href="https://go.dev/dl/" rel="noreferrer">https://go.dev/dl/</a> を開く',
            "「Microsoft Windows」の <strong>.msi</strong> を選ぶ（だいたい一番上の安定版）",
            "ダウンロードしたファイルを開き、画面の指示どおり「次へ」で進む",
            "終わったら PowerShell（さっき開いた命令の窓）を<strong>いったん閉じて、もう一度開き直す</strong>（パスを読み直すため）",
          ],
          macTitle: "macOS",
          macBody: [
            'ブラウザで <a href="https://go.dev/dl/" rel="noreferrer">https://go.dev/dl/</a> を開く',
            "自分の Mac に合うパッケージを選ぶ。<strong>Apple Silicon（M1/M2/M3…）</strong>なら ARM64、<strong>Intel</strong> なら x86-64 の <strong>.pkg</strong>",
            "ダウンロードした .pkg を開き、画面の指示どおりインストールする",
            "終わったらターミナルを<strong>いったん閉じて、もう一度開き直す</strong>",
          ],
          note: "学校や会社の PC では管理者パスワードが必要なことがあります。入れられないときは保護者や先生に相談してください。",
        },
        step2: {
          id: "compiler",
          n: "02",
          h: "Mac だけ：C コンパイラを入れる",
          p: "Ebitengine は中で C の部品も使います。<strong>Windows の学校 PC はこのステップを飛ばして大丈夫</strong>です。Mac だけ、無料の「コマンドラインツール」——プログラムを作るときに裏側で使う道具一式——を入れます。長い名前を覚える必要はありません。教室の PC が Windows なら、すぐ「3. 動作確認」へ進んでください。",
          winSkip: "Windows の人はこのステップを飛ばして「3. 動作確認」へ進んでください。",
          macCmd: "xcode-select --install",
          macBody: [
            "ターミナルに上のコマンドを貼り付けて Enter",
            "ダイアログが出たら「インストール」を選ぶ（数分〜十数分かかることがあります）",
            "すでに入っている場合は「command line tools are already installed」のようなメッセージで大丈夫です",
          ],
        },
        step3: {
          id: "check",
          n: "03",
          h: "本当に入ったか確かめる",
          p: "ターミナルに次を打ち、Enter を押します。",
          cmd: "go version",
          ok: "例: go version go1.25.0 windows/amd64 のように、go1.25 以上の数字が出れば成功です。",
          fail: "「コマンドが見つかりません」「go は認識されていません」と出たら：インストール後にターミナルを開き直したか確認し、それでもだめなら PC を再起動してからもう一度。",
          bonusTitle: "さらに安心したい人へ（任意）",
          bonusP: "公式の回転サンプルを一度走らせると、画面が出るところまでまとめて確認できます。初回はダウンロードに少し時間がかかります。",
          bonusCmd: "go run github.com/hajimehoshi/ebiten/v2/examples/rotate@latest",
          bonusOk: "ゴファーの絵がくるくる回る窓が開けば、環境はバッチリです。閉じたら次へ進みましょう。",
        },
        step4: {
          id: "project",
          n: "04",
          h: "ゲーム用のフォルダを作る",
          p: "デスクトップでもドキュメントでも構いません。ここでは例として home の直下に ebi-empty を作ります。",
          cmds: [
            ["mkdir ebi-empty", "フォルダを作る"],
            ["cd ebi-empty", "その中へ入る"],
            ["go mod init example.com/ebi-empty", "このフォルダの目次 go.mod を作る"],
          ],
          after: "go.mod は、このゲームの名前（モジュール名）と、外から借りる部品（Ebitengine など）を Go が管理するための『プロジェクトの目次ノート』です。モジュール名は、コードの部品をどの名前で呼ぶかを決める住所のようなもの。これがあることで、Go は必要な外部パッケージを同じ組み合わせで用意できます。小さな go.mod ができていれば OK。名前は後から変えられます。",
        },
        step5: {
          id: "code",
          n: "05",
          h: "main.go を書く（何もしない窓）",
          p: "テキストエディタ（メモ帳、メモ帳++、VS Code、Cursor など）で、いまのフォルダに main.go という名前のファイルを作り、次をそのままコピーします。",
          explain: [
            ["Update", "毎フレームの「数字」担当。今は空。"],
            ["Draw", "毎フレームの「絵」担当。色を一塗り。"],
            ["Layout", "ゲーム内部の幅と高さ。"],
            ["RunGame", "このくり返しを起動するスイッチ。"],
          ],
        },
        step6: {
          id: "run",
          n: "06",
          h: "依存関係を入れて、実行する",
          p: "ターミナルのカレントディレクトリが ebi-empty のまま、次を順番に打ちます。",
          cmds: [
            ["go mod tidy", "Ebitengine などを自動で取ってくる"],
            ["go run .", "プログラムを動かす"],
          ],
          success: "タイトルが Empty Window の、暗い青い長方形の窓が開けば成功です。閉じるボタンで終了できます。中身はまだ何も動きません——それでも「ゲームの土台」はもう動いています。",
          next: "次は、この土台の上に Update / Draw の意味を載せる LEVEL 01 へどうぞ。ブラウザ上のデモでも学べます。",
          nextHref: "../../games/tap-target/",
          nextLabel: "LEVEL 01 を開く →",
        },
        trouble: {
          id: "troubleshoot",
          h: "うまくいかないとき",
          lead: "ここが一番よく止まります。下の表から自分の症状に近いものを選んでください。それでもだめなら、赤いエラー全文をコピーしてページ末尾のフィードバックへ送ってください。",
          items: [
            [
              "PowerShell とターミナル、どっち？",
              "Windows なら「PowerShell」で大丈夫。Mac なら「ターミナル」。どちらも命令を打ち込む窓です。Windows の「ターミナル」アプリを開いても、中身が PowerShell なら問題ありません。",
            ],
            [
              "go: command not found / 認識されない",
              "① Go を入れたあと、命令の窓を<strong>完全に閉じて開き直したか</strong>確認。② Windows ならスタートメニューに「Go」があるか確認。③ それでもだめなら PC を再起動してもう一度 <code>go version</code>。",
            ],
            [
              "xcrun: error: invalid active developer path…（Mac）",
              "ステップ 02 の <code>xcode-select --install</code> をもう一度実行する。途中で止まったら再起動してから再実行。",
            ],
            [
              "go mod tidy や go run でネットエラー",
              "インターネット接続と、会社・学校のプロキシ設定を確認。自宅の回線で試すと通ることがあります。",
            ],
            [
              "窓が一瞬で消える / パニック",
              "命令の窓に赤いエラー全文が出ています。その文を検索するか、このページ末尾のフィードバックへ貼ってください。",
            ],
            [
              "WSL（Windows の Linux）を使っている",
              "公式どおり、実行時に GOOS=windows が必要になることがあります。はじめてなら通常の PowerShell での手順を推奨します。",
            ],
          ],
        },
        footerNote: "手順の詳細は公式ドキュメント（Install / Hello, World!）にもあります。このページは「まっさらな PC 向け」に順序を並べ替えた導入です。つまずいたらまず「うまくいかないとき」へ戻ってください。",
        official: "公式 Install",
        officialHref: "https://ebitengine.org/en/documents/install.html",
      }
    : {
        title: "Getting Started — Open an Empty Window | Ebi Showcase",
        desc: "Step-by-step: install Go and Ebitengine on a fresh Windows or macOS machine, then run a do-nothing empty window.",
        crumb: "SPECIAL GUIDE",
        eyebrow: "START HERE / ENVIRONMENT",
        h1: "From a blank computer<br>to an <em>empty game window</em>.",
        lead: "You can play every lesson in the browser. To run code on your own PC, you need Go and Ebitengine first. This guide covers Windows and Mac. Success looks like a dark-blue window that does nothing—that’s the foundation.",
        goalEyebrow: "GOAL",
        goalH: "When you finish, you can…",
        goals: [
          ["Open a terminal", "Where you type commands"],
          ["See go version print a number", "Proof Go is installed"],
          ["Open an empty Ebitengine window", "Proof the game loop can start"],
        ],
        tocEyebrow: "CONTENTS",
        toc: [
          ["#terminal", "0. Open a terminal"],
          ["#install-go", "1. Install Go"],
          ["#compiler", "2. Mac only: C compiler"],
          ["#check", "3. Verify"],
          ["#project", "4. Folder + go.mod"],
          ["#code", "5. Write main.go"],
          ["#run", "6. Run it"],
          ["#troubleshoot", "Troubleshooting"],
        ],
        step0: {
          id: "terminal",
          n: "00",
          h: "Open a place to type commands",
          p: "You’ll use a terminal—not the browser—to talk to the computer. The name differs by OS.",
          winTitle: "Windows",
          winSteps: [
            "Press the Windows key",
            "Type powershell",
            "Open Windows PowerShell or Terminal",
          ],
          macTitle: "macOS",
          macSteps: [
            "Press Command + Space",
            "Type Terminal",
            "Open the Terminal app",
          ],
          tip: "If you see a blinking cursor, you’re ready. You’ll paste commands there.",
          note: "Name tip: on Windows people often say <strong>PowerShell</strong>; on Mac, <strong>Terminal</strong>. Same job—a window where you type commands. Windows 11’s “Terminal” app can host PowerShell. This guide says PowerShell for Windows and Terminal for Mac.",
        },
        step1: {
          id: "install-go",
          n: "01",
          h: "Install Go",
          p: "Ebitengine is a Go game engine, so install Go first. This repository uses Go 1.25 or later.",
          winTitle: "Windows",
          winBody: [
            'Open <a href="https://go.dev/dl/" rel="noreferrer">https://go.dev/dl/</a>',
            "Download the Microsoft Windows <strong>.msi</strong> (usually the latest stable at the top)",
            "Run the installer and click through Next",
            "When it finishes, <strong>close the PowerShell window and open it again</strong> so PATH refreshes",
          ],
          macTitle: "macOS",
          macBody: [
            'Open <a href="https://go.dev/dl/" rel="noreferrer">https://go.dev/dl/</a>',
            "Pick the right <strong>.pkg</strong>: <strong>ARM64</strong> for Apple Silicon (M1/M2/M3…), <strong>x86-64</strong> for Intel",
            "Run the package and follow the prompts",
            "When it finishes, <strong>quit Terminal and open it again</strong>",
          ],
          note: "School or work PCs may need an admin password. Ask a parent or teacher if you can’t install.",
        },
        step2: {
          id: "compiler",
          n: "02",
          h: "Mac only: install a C compiler",
          p: "Ebitengine also uses a little C under the hood. <strong>Windows school PCs can skip this.</strong> On Mac only, install Apple’s Command Line Tools—the behind-the-scenes tool kit used to build programs. You do not need to memorize the long name. If your classroom PC is Windows, jump to step 3.",
          winSkip: "On Windows, skip this step and go to “3. Verify”.",
          macCmd: "xcode-select --install",
          macBody: [
            "Paste the command above into Terminal and press Enter",
            "In the dialog, choose Install (it can take several minutes)",
            "If tools are already present, a message like “already installed” is fine",
          ],
        },
        step3: {
          id: "check",
          n: "03",
          h: "Check that Go works",
          p: "Type this in the terminal and press Enter:",
          cmd: "go version",
          ok: "Success looks like go version go1.25.0 windows/amd64 — any go1.25+ is fine.",
          fail: "If you see “command not found” or “not recognized”: reopen the terminal after install; if it still fails, reboot once.",
          bonusTitle: "Optional extra check",
          bonusP: "Running the official rotate sample proves a window can open. The first run downloads packages.",
          bonusCmd: "go run github.com/hajimehoshi/ebiten/v2/examples/rotate@latest",
          bonusOk: "A spinning gopher window means your environment is ready. Close it and continue.",
        },
        step4: {
          id: "project",
          n: "04",
          h: "Make a project folder",
          p: "Desktop or Documents is fine. Here we create ebi-empty under your home folder.",
          cmds: [
            ["mkdir ebi-empty", "Create the folder"],
            ["cd ebi-empty", "Enter it"],
            ["go mod init example.com/ebi-empty", "Create the go.mod table of contents"],
          ],
          after: "go.mod is the project’s table-of-contents note. It records the module name—the address used to name your code—and lets Go manage borrowed parts such as Ebitengine. With it, Go can prepare the same external packages again. Seeing the small go.mod file means success; you can rename the module later.",
        },
        step5: {
          id: "code",
          n: "05",
          h: "Write main.go (do-nothing window)",
          p: "In a text editor (Notepad, VS Code, Cursor…), create main.go in that folder and paste this exactly:",
          explain: [
            ["Update", "Per-tick numbers. Empty for now."],
            ["Draw", "Per-frame painting. One solid fill."],
            ["Layout", "Internal width and height."],
            ["RunGame", "Starts the loop."],
          ],
        },
        step6: {
          id: "run",
          n: "06",
          h: "Fetch dependencies and run",
          p: "Still inside ebi-empty, run these in order:",
          cmds: [
            ["go mod tidy", "Download Ebitengine and friends"],
            ["go run .", "Run the program"],
          ],
          success: "A dark-blue window titled Empty Window means success. Close it with the window button. Nothing moves yet—but the game foundation is alive.",
          next: "Next, put meaning into Update / Draw with LEVEL 01. You can also learn from the browser demo.",
          nextHref: "../../games/tap-target/",
          nextLabel: "Open LEVEL 01 →",
        },
        trouble: {
          id: "troubleshoot",
          h: "Troubleshooting",
          lead: "This is where people get stuck most. Match your symptom below. Still stuck? Copy the full red error into the feedback box at the bottom.",
          items: [
            [
              "PowerShell vs Terminal—which one?",
              "On Windows, PowerShell is fine. On Mac, use Terminal. Same job: a window for typing commands. Windows 11’s Terminal app is OK if it opens PowerShell inside.",
            ],
            [
              "go: command not found / not recognized",
              "① After installing Go, fully close and reopen the command window. ② On Windows, check Start for “Go”. ③ Still stuck? Reboot, then try <code>go version</code> again.",
            ],
            [
              "xcrun: error: invalid active developer path… (Mac)",
              "Run step 02’s <code>xcode-select --install</code> again. If it stalled, reboot and retry.",
            ],
            [
              "Network errors on go mod tidy / go run",
              "Check internet and school/work proxies. Home networks often work when campus ones don’t.",
            ],
            [
              "Window flashes and disappears / panic",
              "Read the full red error in the command window. Search it, or paste it into the feedback box below.",
            ],
            [
              "Using WSL (Linux on Windows)",
              "Official docs may require GOOS=windows for go run. Prefer plain PowerShell for this first guide.",
            ],
          ],
        },
        footerNote: "Official Install / Hello, World! docs have more detail. This page reorders them for a blank-machine first run. When stuck, jump back to Troubleshooting first.",
        official: "Official Install",
        officialHref: "https://ebitengine.org/en/documents/install.html",
      };

  const code = ja ? blankMainGo : blankMainGoEN;

  const osCard = (title, bodyHtml) =>
    `<article class="setup-os"><h3>${title}</h3>${bodyHtml}</article>`;

  const ol = (items) =>
    `<ol class="setup-ol">${items.map((x) => `<li>${x}</li>`).join("")}</ol>`;

  const cmdList = (cmds) =>
    `<ol class="setup-cmd-list">${cmds
      .map(
        ([c, note]) =>
          `<li><code>${esc(c)}</code><span>${esc(note)}</span></li>`,
      )
      .join("")}</ol>`;

  return `<!doctype html>
<html lang="${lang}">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
  <meta name="description" content="${esc(t.desc)}">
  <title>${esc(t.title)}</title>
  <link rel="stylesheet" href="${css}">
</head>
<body>
<header class="nav">
  <a class="brand" href="${home}"><span>EBI</span> SHOWCASE</a>
  <nav>
    <a href="#install-go">GO</a>
    <a href="#run">RUN</a>
    <a class="lang" href="${otherHref}" lang="${other}" data-language="${other}">${otherLabel}</a>
  </nav>
</header>
<main>
  <div class="lesson-breadcrumb"><a href="${home}">← CURRICULUM</a><span>${t.crumb}</span></div>

  <section class="data-hero setup-hero">
    <p class="eyebrow">${t.eyebrow}</p>
    <h1>${t.h1}</h1>
    <p>${t.lead}</p>
  </section>

  <section class="setup-goals">
    <div class="guide-heading">
      <p class="eyebrow">${t.goalEyebrow}</p>
      <h2>${t.goalH}</h2>
    </div>
    <div class="setup-goal-grid">
      ${t.goals
        .map(
          ([h, p], i) =>
            `<article><b>0${i + 1}</b><h3>${esc(h)}</h3><p>${esc(p)}</p></article>`,
        )
        .join("")}
    </div>
    <nav class="setup-toc" aria-label="toc">
      <p class="eyebrow">${t.tocEyebrow}</p>
      <ol>${t.toc.map(([href, label]) => `<li><a href="${href}">${esc(label)}</a></li>`).join("")}</ol>
    </nav>
  </section>

  <section class="setup-step" id="${t.step0.id}">
    <p class="setup-n">STEP ${t.step0.n}</p>
    <h2>${t.step0.h}</h2>
    <p class="setup-lead">${t.step0.p}</p>
    <div class="setup-os-grid">
      ${osCard(t.step0.winTitle, ol(t.step0.winSteps))}
      ${osCard(t.step0.macTitle, ol(t.step0.macSteps))}
    </div>
    <p class="setup-tip">${t.step0.tip}</p>
    ${t.step0.note ? `<p class="setup-note">${t.step0.note}</p>` : ""}
  </section>

  <section class="setup-step" id="${t.step1.id}">
    <p class="setup-n">STEP ${t.step1.n}</p>
    <h2>${t.step1.h}</h2>
    <p class="setup-lead">${t.step1.p}</p>
    <div class="setup-os-grid">
      ${osCard(t.step1.winTitle, ol(t.step1.winBody))}
      ${osCard(t.step1.macTitle, ol(t.step1.macBody))}
    </div>
    <p class="setup-note">${t.step1.note}</p>
  </section>

  <section class="setup-step" id="${t.step2.id}">
    <p class="setup-n">STEP ${t.step2.n}</p>
    <h2>${t.step2.h}</h2>
    <p class="setup-lead">${t.step2.p}</p>
    <p class="setup-skip">${t.step2.winSkip}</p>
    ${codeBlock(t.step2.macCmd)}
    ${ol(t.step2.macBody)}
  </section>

  <section class="setup-step" id="${t.step3.id}">
    <p class="setup-n">STEP ${t.step3.n}</p>
    <h2>${t.step3.h}</h2>
    <p class="setup-lead">${t.step3.p}</p>
    ${codeBlock(t.step3.cmd)}
    <p class="setup-ok">${t.step3.ok}</p>
    <p class="setup-fail">${t.step3.fail}</p>
    <div class="setup-bonus">
      <h3>${t.step3.bonusTitle}</h3>
      <p>${t.step3.bonusP}</p>
      ${codeBlock(t.step3.bonusCmd)}
      <p>${t.step3.bonusOk}</p>
    </div>
  </section>

  <section class="setup-step" id="${t.step4.id}">
    <p class="setup-n">STEP ${t.step4.n}</p>
    <h2>${t.step4.h}</h2>
    <p class="setup-lead">${t.step4.p}</p>
    ${cmdList(t.step4.cmds)}
    <p class="setup-ok">${t.step4.after}</p>
  </section>

  <section class="setup-step" id="${t.step5.id}">
    <p class="setup-n">STEP ${t.step5.n}</p>
    <h2>${t.step5.h}</h2>
    <p class="setup-lead">${t.step5.p}</p>
    ${codeBlock(code, {
      copy: ja ? "全文をコピー" : "Copy all",
      copied: ja ? "コピーしました" : "Copied!",
      filename: "main.go",
    })}
    <div class="setup-explain">
      ${t.step5.explain
        .map(([k, v]) => `<div><code>${esc(k)}</code><span>${esc(v)}</span></div>`)
        .join("")}
    </div>
  </section>

  <section class="setup-step" id="${t.step6.id}">
    <p class="setup-n">STEP ${t.step6.n}</p>
    <h2>${t.step6.h}</h2>
    <p class="setup-lead">${t.step6.p}</p>
    ${cmdList(t.step6.cmds)}
    <p class="setup-success">${t.step6.success}</p>
    <p class="setup-lead">${t.step6.next}</p>
    <a class="cta setup-cta" href="${t.step6.nextHref}">${t.step6.nextLabel}</a>
  </section>

  <section class="setup-step setup-trouble" id="${t.trouble.id}">
    <h2>${t.trouble.h}</h2>
    ${t.trouble.lead ? `<p class="setup-lead">${t.trouble.lead}</p>` : ""}
    <div class="setup-trouble-list">
      ${t.trouble.items
        .map(
          ([h, p]) =>
            `<article><h3>${esc(h)}</h3><p>${p}</p></article>`,
        )
        .join("")}
    </div>
    <p class="setup-footer-note">${t.footerNote} <a href="${t.officialHref}" rel="noreferrer">${t.official} ↗</a></p>
  </section>

  <section class="feedback-section" aria-labelledby="feedback-title">
    <div class="feedback-card">
      <div class="feedback-heading">
        <p class="eyebrow">FEEDBACK</p>
        <h2 id="feedback-title">${ja ? "ひとことフィードバック" : "Quick feedback"}</h2>
      </div>
      <form class="feedback-form" action="https://docs.google.com/forms/d/e/1FAIpQLSdE74SxJYstsQ2pckmG-IIGwgMMlpcp3w7c2bG-RPso-nQLbA/formResponse" method="POST">
        <input type="hidden" name="entry.765794446" value="${route}">
        <label class="feedback-field">
          <span class="sr-only">${ja ? "フィードバック" : "Feedback"}</span>
          <input class="feedback-message" name="entry.893595607" maxlength="200" required data-sending="${ja ? "送信中…" : "Sending…"}" data-sent="${ja ? "送信しました。ありがとうございます！" : "Sent — thank you!"}" data-failed="${ja ? "送信できませんでした。時間をおいて再試行してください。" : "Could not send. Please try again later."}" placeholder="${ja ? "ひとこと入力…" : "Write one short note…"}">
        </label>
        <div class="feedback-actions">
          <button type="submit" class="feedback-submit">${ja ? "送信する" : "Send feedback"}<span>→</span></button>
          <p class="feedback-status" aria-live="polite"></p>
        </div>
        <input type="hidden" name="fvv" value="1">
        <input type="hidden" name="pageHistory" value="0">
      </form>
    </div>
  </section>
</main>
<footer><div class="brand"><span>EBI</span> SHOWCASE</div><p>Made with Go + Ebitengine.</p><a href="https://github.com/kumagi/EbiShowcase">VIEW SOURCE ↗</a></footer>
<script src="${learn}"></script>
<script>
  document.querySelectorAll("[data-language]").forEach((a) =>
    a.addEventListener("click", () => localStorage.setItem("ebi-language", a.dataset.language)),
  );
</script>
</body>
</html>
`;
}

for (const lang of ["ja", "en"]) {
  const dir = path.join(root, "web", lang, "guides", "setup");
  fs.mkdirSync(dir, { recursive: true });
  fs.writeFileSync(path.join(dir, "index.html"), page(lang));
  console.log("wrote", path.join("web", lang, "guides", "setup", "index.html"));
}
