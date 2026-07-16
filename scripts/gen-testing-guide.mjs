/** Generates the bilingual, non-playable game-testing course. */
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");

const lessons = [
  {
    slug: "pure-functions",
    ja: {
      step: "STEP 01", title: "画面からルールを取り出す", lead: "クリック位置と丸の当たり判定を、Ebitengineを知らない小さな関数へ移します。入力と描画はゲーム側、答えを出す計算は純粋ロジック側です。",
      idea: "同じ数字を渡せば、いつでも同じ答えが返る関数なら、画面を開かずに何度でも確かめられます。",
      before: `func (g *game) Update() error {
  px, py := ebiten.CursorPosition()
  dx := float64(px) - g.circleX
  dy := float64(py) - g.circleY
  if math.Hypot(dx, dy) <= g.radius {
    g.score++
  }
  return nil
}`,
      after: `// Ebitengineをimportしないファイル
func PointInCircle(px, py, cx, cy, radius float64) bool {
  dx := px - cx
  dy := py - cy
  return math.Hypot(dx, dy) <= radius
}

// Updateは入力と状態更新だけ
if lessonlogic.PointInCircle(px, py, g.circleX, g.circleY, g.radius) {
  g.score++
}`,
      test: `func TestPointInCircle(t *testing.T) {
  got := PointInCircle(130, 80, 100, 80, 30)
  if !got {
    t.Fatal("円のふちも当たりにしたい")
  }
}`,
      cases: [["中心", "(100, 80)", "true"], ["ふち", "(130, 80)", "true"], ["1px 外", "(131, 80)", "false"]],
      challenge: "半径が0のとき、中心だけは当たりになるでしょうか。予想してからケースを1行足してみよう。",
    },
    en: {
      step: "STEP 01", title: "Lift the rule out of the screen", lead: "Move the circle hit test out of Ebitengine input code and into a tiny function that knows only numbers. The game owns input and drawing; pure logic owns the answer.",
      idea: "A function that always returns the same answer for the same numbers can be checked repeatedly without opening a window.",
      before: `func (g *game) Update() error {
  px, py := ebiten.CursorPosition()
  dx := float64(px) - g.circleX
  dy := float64(py) - g.circleY
  if math.Hypot(dx, dy) <= g.radius {
    g.score++
  }
  return nil
}`,
      after: `// This file does not import Ebitengine.
func PointInCircle(px, py, cx, cy, radius float64) bool {
  dx := px - cx
  dy := py - cy
  return math.Hypot(dx, dy) <= radius
}

// Update only connects input to state.
if lessonlogic.PointInCircle(px, py, g.circleX, g.circleY, g.radius) {
  g.score++
}`,
      test: `func TestPointInCircle(t *testing.T) {
  got := PointInCircle(130, 80, 100, 80, 30)
  if !got {
    t.Fatal("a point on the edge should hit")
  }
}`,
      cases: [["center", "(100, 80)", "true"], ["edge", "(130, 80)", "true"], ["1px outside", "(131, 80)", "false"]],
      challenge: "What should happen when the radius is zero? Predict it, then add one more case.",
    },
  },
  {
    slug: "table-tests",
    ja: {
      step: "STEP 02", title: "境目をテストの表にする", lead: "タイミングメーターの PERFECT / GREAT / GOOD / MISS は、境目でバグが起きやすいルールです。似た確認をコピーせず、入力と期待する答えを表にします。",
      idea: "8はPERFECT、8.1はGREAT。この『ぎりぎり両側』を並べると、< と <= の間違いがすぐ見つかります。",
      before: `func pointsForDistance(distance float64) (int, string) {
  switch {
  case distance <= 8:
    return 100, "PERFECT +100"
  case distance <= 28:
    return 50, "GREAT +50"
  case distance <= 55:
    return 10, "GOOD +10"
  default:
    return 0, "MISS"
  }
}`,
      after: `tests := []struct {
  distance   float64
  wantPoints int
}{
  {distance: 8,   wantPoints: 100},
  {distance: 8.1, wantPoints: 50},
  {distance: 28,  wantPoints: 50},
  {distance: 55.1,wantPoints: 0},
}`,
      test: `for _, tt := range tests {
  points, _ := TimingScore(tt.distance)
  if points != tt.wantPoints {
    t.Errorf("distance %v: got %d, want %d",
      tt.distance, points, tt.wantPoints)
  }
}`,
      cases: [["PERFECTの最後", "8", "100"], ["GREATの最初", "8.1", "50"], ["GOODの最後", "55", "10"], ["MISSの最初", "55.1", "0"]],
      challenge: "マイナスの距離は現実には来ません。関数側で直すか、呼ぶ側の責任にするかを決め、その約束をテスト名に書こう。",
    },
    en: {
      step: "STEP 02", title: "Turn boundaries into a test table", lead: "PERFECT / GREAT / GOOD / MISS has bug-prone edges. Put inputs and expected answers in a table instead of copying the same test many times.",
      idea: "8 is PERFECT while 8.1 is GREAT. Testing both sides of an edge catches a mistaken < or <= immediately.",
      before: `func pointsForDistance(distance float64) (int, string) {
  switch {
  case distance <= 8:
    return 100, "PERFECT +100"
  case distance <= 28:
    return 50, "GREAT +50"
  case distance <= 55:
    return 10, "GOOD +10"
  default:
    return 0, "MISS"
  }
}`,
      after: `tests := []struct {
  distance   float64
  wantPoints int
}{
  {distance: 8,   wantPoints: 100},
  {distance: 8.1, wantPoints: 50},
  {distance: 28,  wantPoints: 50},
  {distance: 55.1,wantPoints: 0},
}`,
      test: `for _, tt := range tests {
  points, _ := TimingScore(tt.distance)
  if points != tt.wantPoints {
    t.Errorf("distance %v: got %d, want %d",
      tt.distance, points, tt.wantPoints)
  }
}`,
      cases: [["last PERFECT", "8", "100"], ["first GREAT", "8.1", "50"], ["last GOOD", "55", "10"], ["first MISS", "55.1", "0"]],
      challenge: "A negative distance should never arrive. Decide which layer owns that promise and make the test name say so.",
    },
  },
  {
    slug: "state-transitions",
    ja: {
      step: "STEP 03", title: "1 tick後の状態を比べる", lead: "アクティブ戦闘RPGでは、tickごとに『ゲージ + 素早さ』を計算します。1000に届いた瞬間だけREADYになる状態遷移を、1歩ずつテストします。",
      idea: "長い戦闘を最後まで再生する代わりに、987→999、988→READYという大事な1コマだけを直接作ります。",
      before: `func (g *game) Update() error {
  for i := range g.runners {
    r := &g.runners[i]
    r.gauge += r.speed
    if r.gauge >= 1000 {
      r.gauge = 0
      r.ready = true
    }
  }
  return nil
}`,
      after: `func AdvanceGauge(gauge, speed, limit int) (int, bool) {
  next := gauge + speed
  if next >= limit {
    return 0, true
  }
  return next, false
}

// Updateは全員に同じルールを適用する
r.gauge, ready = lessonlogic.AdvanceGauge(r.gauge, r.speed, 1000)`,
      test: `gauge, ready := AdvanceGauge(988, 12, 1000)
if gauge != 0 || !ready {
  t.Fatalf("got (%d, %v), want (0, true)",
    gauge, ready)
}`,
      cases: [["まだ途中", "400 + 12", "412 / false"], ["あと1", "987 + 12", "999 / false"], ["ちょうど", "988 + 12", "0 / true"], ["飛び越す", "995 + 12", "0 / true"]],
      challenge: "1000を7だけ飛び越した分を次のゲージへ残すルールなら、戻り値とテストをどう変えるか考えよう。",
    },
    en: {
      step: "STEP 03", title: "Compare the state one tick later", lead: "An active-battle RPG adds speed to a gauge every tick. Test the exact transition that becomes READY at 1000, one step at a time.",
      idea: "Instead of replaying a whole battle, construct only the important frames: 987→999 and 988→READY.",
      before: `func (g *game) Update() error {
  for i := range g.runners {
    r := &g.runners[i]
    r.gauge += r.speed
    if r.gauge >= 1000 {
      r.gauge = 0
      r.ready = true
    }
  }
  return nil
}`,
      after: `func AdvanceGauge(gauge, speed, limit int) (int, bool) {
  next := gauge + speed
  if next >= limit {
    return 0, true
  }
  return next, false
}

// Update applies the rule to each runner.
r.gauge, ready = lessonlogic.AdvanceGauge(r.gauge, r.speed, 1000)`,
      test: `gauge, ready := AdvanceGauge(988, 12, 1000)
if gauge != 0 || !ready {
  t.Fatalf("got (%d, %v), want (0, true)",
    gauge, ready)
}`,
      cases: [["charging", "400 + 12", "412 / false"], ["one short", "987 + 12", "999 / false"], ["exact", "988 + 12", "0 / true"], ["passes limit", "995 + 12", "0 / true"]],
      challenge: "If overflow should carry into the next gauge, how would you change the return value and its test?",
    },
  },
  {
    slug: "regression-tests",
    ja: {
      step: "STEP 04", title: "直したバグを二度と戻さない", lead: "3マッチの特殊ピースは、十字消しや連鎖のように手動確認が難しい仕組みです。バグが起きた盤面を小さなデータで再現し、修正前に失敗するテストとして残します。",
      idea: "再現盤面は『バグを閉じ込めた標本』です。将来コードを整理しても、同じ失敗が戻れば go test が知らせます。",
      before: `// 人の確認だけ
// 1. ゲームを起動
// 2. 何度も交換
// 3. 5個並ぶまで待つ
// 4. 特殊ピースの連鎖を目で見る`,
      after: `board := [8][8]int{
  {1, 2, 3, 4, 5, 1, 2, 3},
  // 必要な並びだけを固定して書く
}

got := specialForMatch(board, matchedCells)
want := specialRainbow
if got != want {
  t.Fatalf("got %v, want %v", got, want)
}`,
      test: `func TestSpecialForMatch(t *testing.T) {
  tests := []struct {
    name  string
    cells []cell
    want  specialKind
  }{
    {"four horizontal", fourInRow(), specialRow},
    {"five", fiveInRow(), specialRainbow},
    {"cross", crossShape(), specialBurst},
  }
  // 同じ入力を、同じ期待値と比べる
}`,
      cases: [["4個・横", "row", "横一列を消す"], ["5個", "rainbow", "同じ色を消す"], ["十字/T字", "burst", "周囲を消す"], ["特殊×特殊", "chain", "連鎖が続く"]],
      challenge: "次に見つけたバグは、直す前に最小の盤面で失敗するテストを書こう。直したあと、そのテストは消さずに残します。",
    },
    en: {
      step: "STEP 04", title: "Never reintroduce a fixed bug", lead: "Match-three special pieces are hard to reproduce by hand. Freeze the failing board as small data and keep a test that fails before the fix.",
      idea: "A reproduction board is a bug specimen. If a later cleanup brings the bug back, go test warns you immediately.",
      before: `// Manual checking only
// 1. Start the game
// 2. Swap many times
// 3. Wait until five pieces align
// 4. Watch the special-piece chain`,
      after: `board := [8][8]int{
  {1, 2, 3, 4, 5, 1, 2, 3},
  // Freeze only the arrangement you need.
}

got := specialForMatch(board, matchedCells)
want := specialRainbow
if got != want {
  t.Fatalf("got %v, want %v", got, want)
}`,
      test: `func TestSpecialForMatch(t *testing.T) {
  tests := []struct {
    name  string
    cells []cell
    want  specialKind
  }{
    {"four horizontal", fourInRow(), specialRow},
    {"five", fiveInRow(), specialRainbow},
    {"cross", crossShape(), specialBurst},
  }
  // Compare the same input with the same expectation.
}`,
      cases: [["four across", "row", "clear one row"], ["five", "rainbow", "clear one color"], ["cross/T", "burst", "clear neighbors"], ["special + special", "chain", "chain continues"]],
      challenge: "For the next bug, write the smallest failing board before the fix. Keep that test after it passes.",
    },
  },
  {
    slug: "readable-tests",
    ja: {
      step: "STEP 05", title: "テストはコードを読みやすくする", lead: "テスト名とAAA（準備・実行・確認）をそろえると、実装を読まなくても関数への期待が分かります。",
      idea: "テストは答え合わせの表です。名前、Arrange、Act、Assertをそろえると、失敗したときにどこを見るかも読めます。",
      before: `func TestRule(t *testing.T) {
  got, _ := pointsForDistance(8)
  if got != 100 { t.Fatal("wrong") }
}`,
      after: `import "github.com/stretchr/testify/assert"

func TestPointsForDistance_PerfectEdge(t *testing.T) {
  // Arrange / 準備
  distance := 8.0

  // Act / 実行
  got, label := pointsForDistance(distance)

  // Assert / 確認
  assert.Equal(t, 100, got)
  assert.Equal(t, "PERFECT +100", label)
}`,
      test: `func TestPointsForDistance_PerfectEdge(t *testing.T) {
  // Arrange
  distance := 8.0
  // Act
  got, _ := pointsForDistance(distance)
  // Assert
  assert.Equal(t, 100, got)
}`,
      cases: [["名前", "TestPointsForDistance_PerfectEdge", "どんな例か分かる"], ["Arrange", "入力を用意", "準備"], ["Act", "関数を1回呼ぶ", "実行"], ["Assert", "期待値と比較", "確認"]],
      challenge: "次のテストを、関数名_条件の名前にして、Arrange・Act・Assertの3つに分けよう。",
    },
    en: {
      step: "STEP 05", title: "Make Tests Read Like Examples", lead: "Named tests and the AAA pattern—Arrange, Act, Assert—let a reader predict the function's promise without reading its implementation.",
      idea: "A test is a small answer sheet. A clear name and three short parts also show where to look when it fails.",
      before: `func TestRule(t *testing.T) {
  got, _ := pointsForDistance(8)
  if got != 100 { t.Fatal("wrong") }
}`,
      after: `import "github.com/stretchr/testify/assert"

func TestPointsForDistance_PerfectEdge(t *testing.T) {
  // Arrange
  distance := 8.0

  // Act
  got, label := pointsForDistance(distance)

  // Assert
  assert.Equal(t, 100, got)
  assert.Equal(t, "PERFECT +100", label)
}`,
      test: `func TestPointsForDistance_PerfectEdge(t *testing.T) {
  // Arrange
  distance := 8.0
  // Act
  got, _ := pointsForDistance(distance)
  // Assert
  assert.Equal(t, 100, got)
}`,
      cases: [["name", "TestPointsForDistance_PerfectEdge", "names the example"], ["Arrange", "prepare input", "setup"], ["Act", "call the rule once", "run"], ["Assert", "compare expected value", "check"]],
      challenge: "Give one older test a name_condition name and split it into Arrange, Act, and Assert.",
    },
  },
];

const esc = (s) => String(s).replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;");

function shell({lang, depth, title, desc, body}) {
  const ja = lang === "ja";
  const prefix = "../".repeat(depth);
  const other = ja ? "en" : "ja";
  const route = depth === 3 ? "guides/testing/" : `guides/testing/${title.slug}/`;
  return `<!doctype html><html lang="${lang}"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover"><meta name="description" content="${esc(desc)}"><title>${esc(title.text)} | Ebi Showcase</title><link rel="stylesheet" href="${prefix}style.css"></head><body class="testing-guide"><header class="nav"><a class="brand" href="${prefix}${lang}/"><span>EBI</span> SHOWCASE</a><nav><a href="${depth === 3 ? "./" : "../"}">TESTING</a><a class="lang" href="${prefix}${other}/${route}" lang="${other}">${ja ? "EN" : "日本語"}</a></nav></header><main>${body}</main><footer><span>EBI SHOWCASE</span><span>GO + EBITENGINE + WASM</span><span>APACHE-2.0</span></footer><script src="${prefix}learn.js"></script></body></html>`;
}

function hub(lang) {
  const ja = lang === "ja";
  const copy = ja ? {
    title: "ゲームをテストできる形にする", desc: "EbitengineのUpdateから純粋なゲームルールを取り出し、Goのユニットテストで画面を開かずに挙動を確かめる5ステップ教材。",
    h1: "遊んで確かめる前に、<br><em>コードで確かめる。</em>", lead: "ゲームは最後に人が遊んで確かめます。でも、当たり判定や得点、ゲージの境目まで毎回目で試す必要はありません。Updateの奥にある『数字を受け取り、次の数字を返すルール』を取り出して、Goに何度でも確認してもらいましょう。",
    eyebrow: "SPECIAL GUIDE / UNIT TESTING", left: "Ebitengineの仕事", mid: "純粋なルール", right: "またゲームへ", input: "入力を読む", rule: "数字から答えを出す", state: "状態を更新する", noWindow: "窓・GPU・タッチ操作なし", command: "$ go test ./...", pass: "ok  internal/lessonlogic", course: "5つの小さな改造", courseLead: "過去の教材を、テストしやすいコードへ本当にリファクタリングします。デモはありません。コードとケース表が、このコースのビジュアルです。", start: "読む →",
  } : {
    title: "Make games testable", desc: "A five-step guide to extracting pure rules from Ebitengine Update and checking behavior with Go unit tests without opening the game.",
    h1: "Check it in code<br><em>before playing it.</em>", lead: "A person should still play the finished game. But you do not need to manually revisit every collision, score edge, and gauge boundary. Lift the number-in/answer-out rules out of Update and let Go check them repeatedly.",
    eyebrow: "SPECIAL GUIDE / UNIT TESTING", left: "Ebitengine layer", mid: "Pure rule", right: "Back to the game", input: "read input", rule: "turn numbers into an answer", state: "update state", noWindow: "no window, GPU, or touch", command: "$ go test ./...", pass: "ok  internal/lessonlogic", course: "Five small refactors", courseLead: "We genuinely refactor earlier lessons into testable code. There is no demo: code and case tables are the visuals.", start: "READ →",
  };
  const cards = lessons.map((lesson, i) => { const t = lesson[lang]; return `<a class="test-course-card" href="${lesson.slug}/"><span>0${i + 1}</span><h3>${t.title}</h3><p>${t.lead}</p><strong>${copy.start}</strong></a>`; }).join("");
  const guideLinks = `<nav class="test-guide-links" aria-label="${ja ? "テスト教材の次のリンク" : "Testing guide links"}"><a href="pure-functions/">${ja ? "最初のステップへ →" : "START STEP 01 →"}</a><a href="../../">${ja ? "← ホームに戻る" : "← BACK TO HOME"}</a></nav>`;
  const body = `<section class="test-hero"><p class="eyebrow">${copy.eyebrow}</p><h1>${copy.h1}</h1><p>${copy.lead}</p></section><section class="test-boundary" aria-label="logic boundary"><div><small>${copy.left}</small><b>${copy.input}</b></div><i>→</i><div class="is-pure"><small>${copy.mid}</small><b>${copy.rule}</b><em>${copy.noWindow}</em></div><i>→</i><div><small>${copy.right}</small><b>${copy.state}</b></div></section><section class="test-terminal"><code>${copy.command}</code><strong>✓ ${copy.pass}</strong></section><section class="test-course"><p class="eyebrow">COURSE MAP</p><h2>${copy.course}</h2><p>${copy.courseLead}</p><div class="test-course-grid">${cards}</div></section>${guideLinks}`;
  return shell({lang, depth: 3, title: {text: copy.title}, desc: copy.desc, body});
}

function lessonPage(lang, index) {
  const lesson = lessons[index];
  const t = lesson[lang];
  const ja = lang === "ja";
  const prev = index === 0 ? "../" : `../${lessons[index - 1].slug}/`;
  const next = index === lessons.length - 1 ? "../" : `../${lessons[index + 1].slug}/`;
  const rows = t.cases.map((r) => `<tr><th>${esc(r[0])}</th><td><code>${esc(r[1])}</code></td><td><code>${esc(r[2])}</code></td></tr>`).join("");
  const copy = ja ? { before: "BEFORE / Updateに全部ある", after: "AFTER / ルールを分ける", test: "TEST / 期待する答えを書く", table: "先に、確かめたい例を並べる", input: "入力・状態", pure: "純粋関数", answer: "次の状態・答え", why: "なぜこれで安心できる？", challenge: "YOUR TURN", run: "このリポジトリで実行", pass: "PASS — ゲーム画面を開かずに確認できました", back: "← 前へ", forward: "次へ →" } : { before: "BEFORE / Everything lives in Update", after: "AFTER / Separate the rule", test: "TEST / Write the expected answer", table: "List the cases first", input: "input + state", pure: "pure function", answer: "next state + answer", why: "Why does this build confidence?", challenge: "YOUR TURN", run: "Run it in this repository", pass: "PASS — checked without opening the game", back: "← PREVIOUS", forward: "NEXT →" };
  const body = `<section class="test-step-hero"><a href="../">TESTING GUIDE</a><p class="eyebrow">${t.step} / ${String(lessons.length).padStart(2, "0")}</p><h1>${t.title}</h1><p>${t.lead}</p></section><section class="test-rule-strip"><span>${copy.input}</span><i>→</i><strong>${copy.pure}</strong><i>→</i><span>${copy.answer}</span></section><section class="test-explain"><div><p class="eyebrow">${copy.why}</p><h2>${t.idea}</h2></div><div class="test-case-table"><p class="eyebrow">${copy.table}</p><table><tbody>${rows}</tbody></table></div></section><section class="test-code-compare"><article><p>${copy.before}</p><pre><code>${esc(t.before)}</code></pre></article><article class="is-after"><p>${copy.after}</p><pre><code>${esc(t.after)}</code></pre></article></section><section class="test-code-focus"><div><p class="eyebrow">${copy.test}</p><h2>_test.go</h2><p>${t.idea}</p></div><pre><code>${esc(t.test)}</code></pre></section><section class="test-run"><p>${copy.run}</p><code>go test ./internal/lessonlogic</code><strong>✓ ${copy.pass}</strong></section><section class="test-challenge"><p class="eyebrow">${copy.challenge}</p><h2>${t.challenge}</h2></section><nav class="test-pager"><a href="${prev}">${copy.back}</a><span>${index + 1} / ${lessons.length}</span><a href="${next}">${copy.forward}</a></nav>`;
  return shell({lang, depth: 4, title: {text: t.title, slug: lesson.slug}, desc: t.lead, body});
}

for (const lang of ["ja", "en"]) {
  const dir = path.join(root, "web", lang, "guides", "testing");
  fs.mkdirSync(dir, {recursive: true});
  fs.writeFileSync(path.join(dir, "index.html"), hub(lang));
  lessons.forEach((lesson, i) => {
    const lessonDir = path.join(dir, lesson.slug);
    fs.mkdirSync(lessonDir, {recursive: true});
    fs.writeFileSync(path.join(lessonDir, "index.html"), lessonPage(lang, i));
  });
}

for (const lang of ["ja", "en"]) {
  const file = path.join(root, "web", lang, "index.html");
  let html = fs.readFileSync(file, "utf8");
  const ja = lang === "ja";
  const block = `<!-- testing-guide-home:start -->\n<section class="architecture-promo testing-promo"><div><p class="eyebrow">SPECIAL GUIDE / UNIT TESTING</p><h2>${ja ? "ゲームを起動せず、<br>ルールを確かめよう。" : "Check game rules<br>without launching the game."}</h2><p>${ja ? "Updateから純粋な計算を取り出し、当たり判定・得点の境目・行動ゲージ・過去のバグをGoのテストで守る5ステップです。" : "A five-step guide to extracting pure calculations from Update and protecting collisions, score edges, action gauges, and past bugs with Go tests."}</p></div><a href="guides/testing/"><span>READ THE GUIDE</span><strong>${ja ? "ゲームのユニットテスト入門" : "Unit testing for games"}</strong><b>→</b></a></section>\n<!-- testing-guide-home:end -->`;
  const re = /<!-- testing-guide-home:start -->[\s\S]*?<!-- testing-guide-home:end -->/;
  html = re.test(html) ? html.replace(re, block) : html.replace("<!-- visual-effects-home:start -->", `${block}\n<!-- visual-effects-home:start -->`);
  fs.writeFileSync(file, html);
}
