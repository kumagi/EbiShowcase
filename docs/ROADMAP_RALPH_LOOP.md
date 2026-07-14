# Ebi Showcase roadmap Ralph loop — Authoring Pass

このファイルは **「遊べる説明」から「Go を書いてゲームを作れる」への品質パス** の唯一の正本です。
以前の Phase 0–4（プレイ場品質・技術ラボ・卒業器の整備）は完了済みとし、ここには載せていません。
完遂条件は本ファイルの全チェックが埋まることだけです。

## 問題認識（なぜ今これをやるか）

サイトは 208/208 + VFX 29 で遊べ、ラボ・ガイド・卒業の *器* もある。
しかし典型レッスンは次のパターンに落ちている。

1. 完成デモを遊ぶ
2. 完成系の仕組みを読む（同じ段落の使い回しが多い）
3. REAL GO に短い数式やイディオムが出る
4. CHALLENGE が **定数チューニング**（`YOUR FIRST TUNING` / 「値を変えて」）で終わる

結果、学習者は **観察者・調整者** にはなれるが、**作者** にはなれない。

特に共有エンジン系（Rhythm / Raycaster / Tower Defense / Reversi / Top-down Adventure）では:

- ページの数式 ≠ レッスンの `main.go`（薄い `Run(Config)` / `MakeChart` ラッパー）
- CHALLENGE が「ノートを足せ」と言っても、編集すべき API が書いてない

卒業（`graduation/`）と `first-30-minutes` も、現状はリンク集や薄い `AddStar` 級 starter にとどまり、
「次にキーボードで打つ行」が足りない。

## 完遂時の学習者ゴール

あるレッスンを1本終えた学習者は、少なくとも次ができる。

1. **開くファイル**（`games/.../main.go` または明示された `internal/...`）が分かる
2. そのファイルに **ルールを1つ追加**する（分岐・カウンタ・データ行・状態のどれか）
3. `go test` またはページに書かれた検証手順で、足したルールを確認する

定数のチューニングは補助課題にしてよい。**主課題は RULE である。**

## 固定決定

### Ebitengine の大前提（序盤からくり返す公理）

サイト全体・Build Track・LEVEL 01・以降のすべての教材は、次を **口を酸っぱくして** 共有語彙にする。後続レッスンはこれを「もう一度ゼロから」教え直さないが、忘れる前に短く指差し確認してよい。

1. **`Update` が更新と入力のすべてを担う。**
   スコア・位置・タイマー・AI・衝突・シーン遷移など、`game`（状態）を書き換える処理と、キー／ポインタ／タッチの読取りはここに置く。`Draw` に入力もルールも書かない。
2. **`Draw` は常に、任意の `game` 構造体を画面へ投影するだけに勤しむ。**
   ジャンルごとの見た目ルールに従い、いまのフィールド値をピクセルへ写す。計算の「結果」を決める場所ではない。
3. **画面は常に 1 bit 違わず `game` に追従する。**
   同じ `game` なら同じ絵になる。`Draw` 内で `game` を書き換えてはならない（カウンタ加算、乱数、入力読取り、スライス更新などを Draw に混ぜない）。

パラパラ漫画の比喩はこれに従う: **先にコマの中身（数字）を決め、そのコマを描く。描きながら中身を変えない。**

例外を許すならドキュメントに「なぜ Update 側へ戻せないか」を1文で書く。デフォルトは例外なし。

### その他の固定決定

- プレイ可能ゲート **208/208** と VFX **29** は維持する。Build Track を足す場合はゲート外（または別カウント）とし、無闇に 208 を増やさない。
- Authoring（書けるか）と Playable（遊べるか）は **別メーター**。Authoring 未達でも playable 数を落とさない。
- 生成トラックの HTML は手編集しない。generator を直して再生成する。
- REAL GO は必ず次のいずれか: (A) 学習者が編集する入口コード、(B) 入口＋`internal` 抜粋の二層。ページだけの架空ループは禁止。
- CHALLENGE は既定で RULE。編集先パスと関数名を1文で含める。RULE は原則 **`Update`（または Update が呼ぶ純関数）側** に足す。`Draw` にルールを足させる課題は作らない。
- 卒業プロジェクトは brief 記事 + 穴あき starter + 赤から緑になるテスト + 照合用 reference を必須とする。starter も上記公理に従う。
- 新ジャンル（旧 Phase 5 候補）は、本 Authoring Pass 完遂後の採用審査まで着手しない。
- Apache-2.0 / 互換ライセンス方針、Go 1.25 ベースライン、相対 URL、日英同期は変更しない。

## ループ手順

必ず **未完了の先頭1件だけ** を実装する。

```sh
node scripts/roadmap-ralph-loop.mjs status
node scripts/roadmap-ralph-loop.mjs next
node scripts/roadmap-ralph-loop.mjs evidence P0-01
# 実装と検証
node scripts/roadmap-ralph-loop.mjs check P0-01
node scripts/roadmap-ralph-loop.mjs verify
```

- `check` は証跡ファイルが `Status: PASS` かつ必須項目すべて `[x]` でないと通らない。
- フェーズ境界と完了宣言前は `verify --full`、全完了後は `complete`。
- 後続作業が先行証跡を無効化したら `uncheck` して戻す。

証跡は `docs/roadmap-evidence/<TASK-ID>.md`。旧 Phase の証跡は履歴として残してよいが、本ループの完了判定には使わない。

---

## Phase 0 — 著者契約の固定、Update/Draw 公理、現状棚卸し

サイト全体の約束を「値を変える」から「ルールを足す」へ切り替える。固定決定に記した Update/Draw 公理を、ドキュメントと序盤ページの共有語彙にする作業もここに含む。

### やることの詳細

- **Update/Draw 公理**: 固定決定の3条を AGENTS・Glossary・AUTHORING_CHECKLIST・home/setup/LEVEL01 の beginner-bridge へ刻む。状態の単一の真実は `game` にあり、画面はその投影である。
- **著者 Definition of Done**: 編集ファイル・RULE・検証の三つ。RULE は Update（またはそれが呼ぶ純関数）側。
- **定型コピー**: 「値を変えて」→「1ルール足して確かめる」（実施済みなら回帰）。
- **不一致在庫**: 薄い `main.go` / snippet 不一致に加え、Draw 内での状態書き換えを表にする。
- **ドキュメント接続**: README / AGENTS を本ファイルに合わせる。

### タスク

- [x] `P0-01` — AGENTS.md に Authoring Definition of Done（編集先・RULE・検証）を追記し、トラック教材の契約を TUNING 主から RULE 主へ更新する
- [x] `P0-02` — README のロードマップ節を本 Authoring Pass の要約とコマンド案内だけに書き換え、旧 Phase 1–5 の未チェック一覧を削除する
- [x] `P0-03` — サイト定型文「値を変えて / change values」を「1ルール足して確かめる / add one rule and verify」系へ置換する（生成器・OGP・ガイド・卒業ページのソースを含む）
- [x] `P0-04` — `docs/AUTHORING_CHECKLIST.md` を新設し、レッスン／トラック／卒業の著者品質項目と PLAYABLE との分離を定義する
- [x] `P0-05` — 薄い main.go ラッパーと snippet 不一致ページの棚卸し表を証跡に残す（少なくとも Rhythm / Raycaster / TD / Reversi / Top-down Adventure）
- [x] `P0-LOOP-01` — AGENTS・AUTHORING_CHECKLIST・README 学習道筋に Update/Draw 公理の3条を正式追記し、RULE は Update 側・Draw は投影のみと明記する
- [x] `P0-LOOP-02` — Glossary の概念マップと本文を「画面は game に1bit違わず追従／Draw は game を書かない」に強化する（図解キャプション含む）
- [x] `P0-LOOP-03` — setup・ホーム・LEVEL 01 の beginner-bridge／概念カード／lab 文言を公理の3条で強化する（`inject-beginner-bridges.mjs` と LEVEL 01 本文、日英同期）
- [x] `P0-LOOP-04` — Draw 内で game を書き換えている実装・教材表現を棚卸し、あれば修繕キューを証跡表に残す（典型: Draw での乱数、カウンタ++、入力読取り）
- [x] `P0-06` — `scripts/roadmap-ralph-loop.mjs` の証跡要件を Authoring タスク種別向けに揃え、evidence README を本パス用に更新する
- [x] `P0-07` — ベースライン確認: `ralph-loop.sh status` が 208/208・VFX29、`go test` と構造 verify が緑であることを証跡に記録する
- [x] `P0-08` — 新ジャンル採用審査とゲート水増しを本パス完遂まで凍結する旨を README / AUTHORING_CHECKLIST に明記する

---

## Phase 1 — Build Track（空の窓 → 触るルール）と Core 前半の RULE 化

setup の「空の窓」と LEVEL 01 の「完成ゲーム」の間に、**打鍵しながら積み上げる短いレーン**を置く。
並行して Core 前半の challenge を TUNING から RULE へ差し替える。

### Build Track の設計詳細

推奨はゲート外の 4 ステップ（ブラウザデモでもローカル写経でも可）。各ステップのページは次を必須とする。

1. **いま動いている最小コード**（または埋め込みエントリ）全体が見える
2. **次に足す N 行**が枠で示される（完成品を先に見せすぎない）
3. 検証: 画面で確認できる変化 + 可能なら小さな `go test`
4. 次ステップへの pager
5. **毎STEPで公理を1回指差す**（短文でよい）: 入力と数の変化は Update、絵は Draw、Draw は game を触らない

| STEP | 学習者が書くこと | 終わりの状態 | 公理の焦点 |
| --- | --- | --- | --- |
| 01 | 空の `Update` / `Draw` / `RunGame` | 単色キャンバスが回る | 二つの関数の役割分担だけ |
| 02 | `game` に座標を持ち、`Draw` がその値だけを見て図形を描く | 静止した図形 | Draw は投影。位置の「真実」は game |
| 03 | **Update** で入力を読み `score++` | 数が増える | 入力も加点も Update。Draw はスコアを読むだけ |
| 04 | **Update** で当たり判定し位置を書き換え、Draw は新しい位置を描く | LEVEL 01 の核と同じ型 | ルールは Update。同じ game なら絵は一致する |

LEVEL 01 本体は **完成对照**として残し、「写経の答え合わせ」に使う。Build Track から LEVEL 01 へ、LEVEL 01 から Build Track へ相互リンクする。STEP 02 で「Draw に `x++` を書いてはいけない」反例を1つ示してから正しい書き方へ誘導してよい。

### Core RULE 課題の書き方

`YOUR FIRST TUNING` を残すなら副題へ下げ、主 challenge を例えば次の形にする。

> **YOUR FIRST RULE** — `games/core/tap-target/main.go` の `Update` 内、命中処理の直後に `combo++` を足し、`Draw` はその値を読むだけにして表示せよ。埋め込みマップの ✦ がその行付近。

検証は「ローカルで改変して動かす」か、ルールを `internal/lessonlogic` に切り出してテストするかをレッスンごとに1つ選ぶ。曖昧な「改良してみよう」は禁止。Draw にルールを足す課題は作らない。

### タスク

- [x] `P1-BT-SPEC` — Build Track 4ステップの日英コンテンツ表（slug・次に足す行・検証・公理の焦点・レベル01との関係）を scripts または docs に固定する
- [x] `P1-BT-01` — Build Track STEP 01（空のゲームループ）の Go + 日英ページを実装し、「次に足す行」枠と公理の指差しを付ける
- [x] `P1-BT-02` — Build Track STEP 02（game の値を Draw が投影）を実装し、Draw で状態を進める反例を1つ示してから正しい差分へ導く
- [x] `P1-BT-03` — Build Track STEP 03（Update で入力と得点）を実装し、RULE 課題を Update 内のスコア加算に固定する
- [x] `P1-BT-04` — Build Track STEP 04（Update で当たり判定・再配置）を実装し、Hypot を Update 側で学習者が書く行として明示する
- [x] `P1-BT-HUB` — Build Track ハブとホーム／progress／setup からの導線を日英で追加する（ゲート数は増やさない）。ハブ先頭に公理の3条を置く
- [x] `P1-BT-VERIFY` — Build Track 全ステップを Desktop / Phone で入力確認し、日英・サムネ・構造チェックを通す
- [x] `P1-CORE-01` — LEVEL 01 tap-target の主 challenge を YOUR FIRST RULE（Update 側）にし、編集ファイルと挿入位置を本文に書き、公理の3条を本文で再掲する
- [x] `P1-CORE-02` — LEVEL 02 timing-meter に YOUR FIRST RULE（例: Perfect 連続で bonus）を追加する
- [x] `P1-CORE-03` — LEVEL 03 catch-stars に YOUR FIRST RULE（例: ミス連続で状態変化）を追加する
- [x] `P1-CORE-04` — LEVEL 04 flappy に YOUR FIRST RULE（例: スコア閾で色や速さフラグ）を追加する
- [x] `P1-CORE-05` — LEVEL 05–06（pong / breakout）に各1つの YOUR FIRST RULE を追加する
- [x] `P1-CORE-LINK` — Core 前半ページと Build Track・testing ガイドを相互リンクし、first-30 草案の仮リンクを置く

---

## Phase 2 — 二層コードパネルと生成トラックの RULE 化

共有エンジン系を最優先で直す。ページが示すコードと、学習者が開くファイルを一致させる（または二層で両方見せる）。

### 二層パネルの仕様

各生成レッスンの `code-lesson`（または同等）は次の順で出す。

1. **編集する入口** — その STEP の `games/tracks/.../main.go` 全文、または実際に触る API（例: `MakeChart` / `Taps` / `Config{Step}`）
2. **仕組みの抜粋** — `internal/<pkg>/...` の短い断片 + **リポジトリ相対パス**
3. **CHALLENGE** — 「どのファイルのどの関数に何を足すか」。チューニングだけの課題は置かない（別枠 TUNING は可）

### concept-row の仕様

3枚は互いに複製禁止。必ず次の役割を分ける。

1. **データ形**（何がスライス／構造体に入っているか）
2. **Update 順**（いつ数が変わるか）
3. **Draw 写像**（どの数が画面の何になるか）

lead 文と DEEP DIVE を3回貼るのは不合格。

### トラック別の編集先の例

| トラック | 入口で見せるもの | RULE 課題の例 |
| --- | --- | --- |
| Rhythm | `MakeChart` / `Taps` / `Holds` | ノートを4つ追加、または1つ長押しを足す |
| Raycaster | ミッション設定 or `raycasterui` に渡す地図 | 壁を1つ開ける／敵を1体足す |
| Tower Defense | `Config{Step}` とデータテーブル | ウェーブに敵種を1行足す |
| Reversi | 評価関数やCPU設定の入口 | 角の重みを変えず、別セル評価を1項足す |
| Top-down Adventure | 部屋／鍵フラグの設定口 | 鍵を要するドアを1つ足す |

エンジン本体を毎STEP複製しない。ただし **ページが架空の自前ループだけを「REAL GO」と呼ばない。**

### タスク

- [x] `P2-HELP` — 生成器用の二層 code-lesson ヘルパ（入口・抜粋・パス・RULE challenge フィールド）を scripts に実装する
- [x] `P2-RHY-AUDIT` — Rhythm 全STEPの snippet / main.go / challenge 不一致を証跡表にする
- [x] `P2-RHY-GEN` — `gen-rhythm-track.mjs` を二層パネル・一意 concept-row・RULE challenge 対応に改修し再生成する
- [x] `P2-RHY-VERIFY` — Rhythm 全STEPで編集先が本文から辿れ、Desktop/Phone と日英を確認する
- [x] `P2-RAY-AUDIT` — Raycaster 全STEPの不一致表を作る
- [x] `P2-RAY-GEN` — `gen-raycaster-track.mjs` を同仕様で改修し再生成する
- [x] `P2-RAY-VERIFY` — Raycaster を著者基準で検証する
- [x] `P2-TD-AUDIT` — Tower Defense 全STEPの不一致表を作る
- [x] `P2-TD-GEN` — `gen-tower-defense-track.mjs` を同仕様で改修し再生成する
- [x] `P2-TD-VERIFY` — Tower Defense を著者基準で検証する
- [x] `P2-REV-AUDIT` — Reversi 全STEPの不一致表を作る
- [x] `P2-REV-GEN` — `gen-reversi-track.mjs` を同仕様で改修し再生成する
- [x] `P2-REV-VERIFY` — Reversi を著者基準で検証する
- [x] `P2-TOP-AUDIT` — Top-down Adventure 全STEPの不一致表を作る
- [x] `P2-TOP-GEN` — `gen-topdown-adventure-track.mjs` を同仕様で改修し再生成する
- [x] `P2-TOP-VERIFY` — Top-down Adventure を著者基準で検証する
- [x] `P2-HAND-PATTERN` — 手書き良例（platformer 系）から二層＋RULE の適用パッチ手順を docs に短い手順書として残す

---

## Phase 3 — 卒業制作を教材化し、最初の30分を「書く」ルートにする

卒業を「器がある」から「赤テストを緑にする」体験へ格上げする。
first-30-minutes はリンク集をやめ、時間内に **リポジトリ差分が残る**手順にする。

### arcade-60 の完成要件

starter は `AddStar` だけで終わらせない。最低でも次が穴あきで存在する。

- 入力（キーまたはポインタ）でアクションが起きる場所
- スコア更新
- 制限時間または終了条件
- 終了後の結果表示（DebugPrint でもよいが、TODO で UI Lab 接続を示してよい）
- `go test` が最初は失敗し、指定 TODO を埋めると通る

reference/ は写経禁止・最終照合用と記事に明記する。

exploration-3rooms / puzzle-3stages も同様に、テスト名と TODO が 1:1 になること。

### first-30-minutes の完成要件

| 分 | 成果物 |
| --- | --- |
| 0–5 | setup で空窓がローカル実行できる |
| 5–15 | Build Track の STEP を1つ以上進め、自分の `main.go` 差分がある |
| 15–25 | testing 入門の最小純関数＋テストが通る |
| 25–30 | arcade-60 の最初の赤テストが緑になる |

「読んで終わった」だけでは P3-F30 を完了としない。

### タスク

- [x] `P3-ARC-STARTER` — `graduation/arcade-60/starter` を穴あきの60秒ゲーム骨格＋失敗するテスト群へ拡充する
- [x] `P3-ARC-ARTICLE` — 日英の arcade-60 解説を、テスト名と TODO 位置が対応する手順記事にする
- [x] `P3-ARC-VERIFY` — arcade-60 を starter→テスト緑→reference 照合まで通し、Mobile 幅の記事可読性を確認する
- [x] `P3-EXP-STARTER` — `graduation/exploration-3rooms/starter` を状態・鍵・遷移の穴あき＋テストへ拡充する
- [x] `P3-EXP-ARTICLE` — 日英の exploration brief を TODO/テスト対応の手順記事にする
- [x] `P3-EXP-VERIFY` — exploration-3rooms の著者フローを検証する
- [x] `P3-PUZ-STARTER` — `graduation/puzzle-3stages/starter` をデータ駆動進行の穴あき＋テストへ拡充する
- [x] `P3-PUZ-ARTICLE` — 日英の puzzle brief を TODO/テスト対応の手順記事にする
- [x] `P3-PUZ-VERIFY` — puzzle-3stages の著者フローを検証する
- [x] `P3-HUB` — `web/{ja,en}/graduation/` を前提・クローン手順・4コマンド・3 brief の本文化し、カード置き場で終わらせない
- [x] `P3-F30` — `guides/first-30-minutes` を上表の書くルートへ日英とも書き直し、Build Track と arcade-60 に直結する
- [x] `P3-NAV` — progress / choose-your-path / ホームから「MAKE（書く）」モードが PLAY と並んで見えるようにする
- [x] `P3-TRACK-CTA` — 主要トラックハブ（少なくとも U–Y と platformer）に「この型で自分の1本へ」→ graduation への CTA を付ける

---

## Phase 4 — Authoring ゲートの計測、横展開、リリース

残 Core と代表的な手書きトラックへ横展開し、メーターと監査で戻らないようにする。

### Authoring メーター

`ralph-loop.sh status` または併設スクリプトが、少なくとも次を JSON で出せるとよい。

- playable gated（既存）
- authoring: Build Track 完了、Core RULE 付与数、生成5トラックの二層化、graduation 3本の穴あき化、サンプル監査の合否

### 横展開の下限

全208ページの文章を一夜で書き換えなくてよい。完遂条件は:

1. Build Track + Core 01–12 すべてに YOUR FIRST RULE（または Build との明示的役割分担）
2. 生成5トラック（U–Y）が二層＋RULE 済み
3. 手書き良例トラックを **1コースまるごと**（全STEP）著者仕様に更新し、他トラック用のチェックリスト運用を文書化
4. 無作為5レッスンの抜き取り監査に合格
5. フル verify / Pages 想定ビルドが緑

### タスク

- [x] `P4-CORE-LATE` — LEVEL 07–12 すべてに YOUR FIRST RULE（編集先パス付き）を追加する
- [x] `P4-METRIC` — playable と別に authoring 進捗を表示するスクリプトまたは `ralph-loop.sh` 拡張を入れる
- [x] `P4-HAND-01` — platformer トラック全STEPを二層コードまたは同等の編集先明示＋RULE challenge に更新する
- [x] `P4-HAND-02` — 追加で手書きトラックを1つ（推奨: survivors または match3）同様に更新し、手順書の再利用性を確かめる
- [x] `P4-CHECKLIST` — `AUTHORING_CHECKLIST.md` を最終版にし、ADVANCED_QUALITY との関係（遊べる ≠ 書ける）を明記する
- [x] `P4-SAMPLE` — カリキュラムから無作為5レッスンを抜き、著者基準で監査して結果を証跡に残す（不合格なら当該を直してから PASS）
- [x] `P4-COPY-REGRESS` — 「値を変えて」定型文の回帰検索を零件にし、OGP 再注入後も残っていないことを確認する
- [x] `P4-RELEASE` — `verify --full` と `complete` 前提のフル監査（ビルド、テスト、日英、Build Track、graduation、メーター）を通し、Authoring Pass 完了を宣言する

---

## 完了後（本ファイルの外）

`P4-RELEASE` 完了後のみ検討する。

- 残ジャンルトラックへの Authoring 横展開（継続運用）
- Feedback / AI 巡回を「書けなかったポイント」観点でトリアージ
- 旧 Phase 5 ジャンル候補の採用審査（計算の型が本当に欠けているときだけ）

---

## やらないこと（本パス中）

- プレイ可能ゲートの水増し目的の新規 STEP 大量追加
- 説明文の言い換えだけで RULE 課題を入れないこと
- ページに出ない架空コードを REAL GO と呼ぶこと
- `dist/` の手編集、生成 HTML の場当たりパッチ
- 有名作品のアセット・名称・シナリオの複製
- Authoring 未完了のまま新ジャンル制作に逃げること

## タスクID一覧の数え方

エージェントは `node scripts/roadmap-ralph-loop.mjs status` の `total` を正とする。
上から Sequential に進め、飛ばさない。
