# Ebi Showcase

Ebitengine で作ったミニゲームを、ブラウザの WebAssembly でその場で遊びながら学べる、日英バイリンガルの静的学習サイトです。ゲーム本体はすべて Go + Ebitengine で実装し、HTML / CSS / JavaScript は展示ページと共有 WASM ローダーだけに使います。

公開サイト: https://kumagi.github.io/EbiShowcase/

主人公はオリジナルキャラ **海老・天次郎（えび・てんじろう / Ebi Tenjiroh）** です。有名作品はジャンルの目印として触れるだけで、アート・音楽・名称・シナリオはコピーしません。

## いまの到達点

| 区分 | 状態 |
| --- | --- |
| プレイ可能ゲート（Core + ジャンルトラック） | **208 / 208** |
| Visual Effects Lab（ゲート外） | **29 / 29** |
| Core LEVEL 01–12 | 完成（ゲームループから弾幕まで） |
| ジャンルトラック | **25** コース（各トラックは中間レッスン + 統合ゲーム） |
| 高度品質パス（A–T） | `docs/ADVANCED_QUALITY_CHECKLIST.md` 通過済み |
| 横断ガイド | 環境構築 / ユニットテスト / ゲームデータ構成 |

カリキュラムの集計は次で確認できます。

```sh
bash scripts/ralph-loop.sh status   # playable / total（VFX は別カウント）
bash scripts/ralph-loop.sh next     # 未完成の先頭（完了時は complete）
bash scripts/ralph-loop.sh verify   # 全 WASM 再ビルド + 構造チェック
```

## 学びの道筋

おすすめの順番は次のとおりです。一覧を全部いっぺんに消化する必要はありません。

1. **環境づくり**（任意）— `guides/setup/` で Go を入れ、空の Ebitengine 窓を開く
2. **Core LEVEL 01–12** — `Update` / `Draw` から衝突、タイル、カメラ、敵 AI、弾幕へ
3. **ユニットテスト入門** — `guides/testing/` でルールを純粋関数として確かめる
4. **Visual Effects Lab** — 描画の道具 → 光学トリック → Core への演出載せ直し
5. **ジャンル専門化** — 25 トラックから興味のあるジャンルへ
6. **データと構成** — `guides/game-data/` でルールとデータの分け方を学ぶ

用語の索引は [`docs/Glossary.md`](docs/Glossary.md) です。貢献・エージェント向けの運用契約は [`AGENTS.md`](AGENTS.md) を正とします。

## ローカルで動かす

Go 1.25 以降が必要です。

```sh
go mod download
bash scripts/build.sh
# ローカルで OGP の入力が変わっていなければ再生成を省く
bash scripts/build.sh --fast
python3 -m http.server 8080 --directory dist
```

ブラウザで `http://localhost:8080` を開きます。WASM は `file://` では読み込めないため、ローカル HTTP サーバーが必要です。

SNS 向けの OGP は `node scripts/inject-ogp.mjs` と `go run ./cmd/gen-og-images` で全ページ分を生成します（通常の `build.sh` と CI では常に実行）。ローカルの `build.sh --fast` でもメタデータを必ず再注入し、HTML・OGP生成器・フォント・`SITE_ORIGIN` の入力ハッシュと既存PNGが一致するときだけ重いPNG再生成を省きます。公開 URL のオリジンは環境変数 `SITE_ORIGIN`（省略時 `https://kumagi.github.io/EbiShowcase`）です。

## GitHub Pages へ公開する

1. このディレクトリを GitHub リポジトリへ push します。
2. リポジトリの **Settings → Pages → Build and deployment → Source** を **GitHub Actions** にします。
3. `main` ブランチへ push すると `.github/workflows/pages.yml` がビルド・公開します。

プロジェクトサイト（`https://USER.github.io/REPO/`）でも動くよう、アセットはすべて相対パスです。

## 構成

- `games/core/<slug>/` — Core LEVEL のゲーム実装（Flappy のみ歴史的経緯で `game/main.go`）
- `games/tracks/<track>/<slug>/` — ジャンル／VFX 各ステップのゲーム実装
- `internal/` — 共有アトラス、VFX、レッスン用純ロジック、ジャンル固有ヘルパー
- `web/index.html` — ブラウザ言語を判定する入口
- `web/ja/`・`web/en/` — 日英の独立した教材ページ（共有しやすい URL）
- `web/{ja,en}/games/` — Core LEVEL 01–12
- `web/{ja,en}/tracks/` — Visual Effects Lab + 25 ジャンルトラック
- `web/{ja,en}/guides/` — setup / testing / game-data
- `web/game.html` — 共有 WASM ローダー
- `web/assets/` — サムネ、アトラス、OGP、図解など
- `examples/data-driven/` — `go:embed` と JSON によるデータ読み込みの実例
- `docs/Glossary.md` — フレーム、座標、ラジアン、`iota`、スライスなどの用語集
- `docs/ADVANCED_QUALITY_CHECKLIST.md` — ジャンルトラックの第2品質パス
- `docs/ROADMAP_RALPH_LOOP.md` — Phase 0–4 を完遂する順序付きチェックリスト
- `scripts/build.sh` — `dist/` を生成するビルド（生成 HTML の更新を含む）
- `scripts/ralph-loop.sh` — カリキュラム進捗と検証
- `scripts/roadmap-ralph-loop.mjs` — 証跡付きロードマップ進捗と完了判定
- `scripts/feedback-sheet.mjs` — フォーム回答のトリアージ
- `scripts/ai_feedback_crawler.py` / `docs/AI_FEEDBACK_AGENT.md` — LM Studio で改善提案を巡回生成（送信は明示指定時のみ）
- `.github/workflows/pages.yml` — GitHub Pages デプロイ

`dist/` は生成物で Git 管理外です。手編集しないでください。

言語ページは検索・共有しやすい独立 URL です。ルートへの初回アクセスではブラウザ言語を判定し、画面上で言語を切り替えた後はその選択をブラウザに保存します。

## 操作の目安

教材デモは、入力を同じゲーム内アクションへ変換します。レッスンごとに画面下や記事内の操作説明が正です。

| 入力 | アクション | 主な用途 |
| --- | --- | --- |
| Space / ↑ | ジャンプ・決定・拍のタイミング | フライト、横アクション、会話送り、リズム |
| ← → ↑ ↓ | 移動・選択 | 迷路、倉庫番、戦略、ダンジョン |
| クリック / タップ | 決定・配置・カード選択・射撃 | タッチ教材、TD、デッキ構築、レイキャスト |
| R | リトライ（対応デモ） | 失敗したステージを最初から遊ぶ |

スマホでは、ボタンやタップ判定を最低 48dp（CSS ではおよそ 48px）以上にし、押した瞬間に色や粒で反応を返すのが目安です。キーボード・ポインタ・タッチのすべてで完走できることが、各完成ゲームの品質基準です。

## フィードバックを確認する

各教材ページの末尾には Google Form を埋め込んでいます。回答先のスプレッドシートは、OAuth でログインした Google アカウントから読み書きします。

1. Google Cloud で Sheets API を有効にし、OAuth クライアント（デスクトップアプリ）を作成します。
2. ダウンロードした JSON を `.secrets/` に保存します（このファイルは Git へ入りません）。
3. 次のコマンドを初回だけ実行し、ブラウザで Google ログインとアクセス許可を行います。

```sh
node scripts/feedback-sheet.mjs list
```

4. 次のコマンドで一覧を取得・更新します。

```sh
node scripts/feedback-sheet.mjs list
node scripts/feedback-sheet.mjs pending  # 未対応だけ
node scripts/feedback-sheet.mjs check 12   # 12行目を対応済みにする
node scripts/feedback-sheet.mjs delete 12  # 12行目を削除する
node scripts/feedback-sheet.mjs archive    # 履歴を別タブへ移し、回答シートを見出しだけに戻す
```

`check` は「対応済み」列がなければ自動追加し、✅を記録します。行番号は一覧表示の番号を使ってください。
`archive` はフォームの回答シートを削除せず、履歴を `feedback_archive_YYYYMMDDHHMMSS` タブへコピーしてから旧回答行だけを削除します。フォーム連携と見出しは維持されるため、次の投稿は2行目から始まります。

## ロードマップ

量のゲート（208 + VFX 29）は到達済みです。これからの重心は **深さ・欠けている技術柱・卒業後の道筋** です。新規ジャンルを増やし続けるより、いま公開している教材を信頼できる「自分の1本が書ける」ルートに育てます。

実行順・完了条件・証跡様式の正本は [`docs/ROADMAP_RALPH_LOOP.md`](docs/ROADMAP_RALPH_LOOP.md) です。Phase 0–4 の88項目を必ず上から1件ずつ進めます。

```sh
node scripts/roadmap-ralph-loop.mjs status       # フェーズ別進捗
node scripts/roadmap-ralph-loop.mjs next         # 次の1件
node scripts/roadmap-ralph-loop.mjs evidence ID  # 証跡ファイルを作成
node scripts/roadmap-ralph-loop.mjs check ID     # 証跡を検査して完了
node scripts/roadmap-ralph-loop.mjs verify       # 順序と証跡を検証
```

### 方針

- **ゲート総数を無闇に伸ばさない。** 新しい実験は VFX と同様、サイドカウントやガイドとして切る。
- **公開済みトラックの品質を先に揃える。** とくに拡張枠 U–Y（Reversi / Raycaster / Rhythm / Tower Defense / Top-down Adventure）を A–T と同じ playground quality に載せる。
- **予告済みの技術ギャップを埋める。** Optical Tricks の先にある本物のシェーダ、ほとんどの教材でまだ薄い音、Text/UI、カメラの切り出し。
- **出口を太くする。** 「遊べる → 仕組みが分かる → 自分の短いゲームが書ける」の最後の矢印を卒業プロジェクトと導線改善で支える。

### Phase 1 — 公開済みの信頼を固める（最優先）

- [ ] `docs/ADVANCED_QUALITY_CHECKLIST.md` にトラック U–Y を追加し、最終ゲーム・中間レッスン・日英・操作・リプレイ動機を同水準で通す
- [ ] U–Y の最終ゲームを Desktop + スマホ幅で抜き取り実プレイ監査する
- [ ] README / ホーム表記 / サムネが現状とズレていないかを回帰チェックする

### Phase 2 — 技術ラボ（ジャンル数を増やさず深さを足す）

ゲート外（VFX と同様のサイドカウント）を想定します。

- [ ] **Shader Lab** — `FragmentShader` による blur / 色収差 / 歪みなど。Optical Tricks からの自然な続き
- [ ] **Audio Lab** — SE キュー、ループ BGM、簡易エンベロープ。Rhythm の純 Go 合成以外へ広げる
- [ ] **Text / UI Lab** — CJK フォント埋め込み、ダイアログ、メニューフォーカス（`DebugPrint` 卒業）
- [ ] **Camera Lab** — 追随、デッドゾーン、シェイク、レターボックスを独立して操作できる玩具
- [ ] 4ラボの技術を25個すべてのジャンルトラック完成ゲームへ適用し、Desktop / Tablet / Phone、日英、回帰テストを通す

### Phase 3 — 横断ガイドの第2波と卒業口

- [ ] **セーブとシーン遷移**ガイド（既存の `storage_js.go` パターンを一般化）
- [ ] **プロジェクト分割**ガイド（`internal/lessonlogic` を自分のゲームへ当てはめる）
- [ ] **配布ガイド**（ローカル実行 / GitHub Pages / アセット許諾のチェックリスト）
- [ ] **パフォーマンス入門**（LEVEL 12 弾幕の続き：再利用と計測）
- [ ] **卒業プロジェクト** 3本（解説記事 + コピーして始められる Go スターター + テスト + 完全な参考実装）

### Phase 4 — 導線と継続改善の定着

- [ ] ホームを興味別入口（アクション / パズル / 戦略 / 物語 / 描画）と「初めて向け最短」に分ける
- [ ] フィードバック週次トリアージ（`pending` → 修正 → `check`）を定例化
- [ ] AI 巡回提案（`ai_feedback_crawler.py`）は人がサンプル確認してから `--submit`
- [ ] 最終ゲームのローテーション実プレイ監査を継続
- [ ] 必要なら `ralph-loop.sh status` 横に品質副指標（例: U–Y quality）を足す

### Phase 5 — 新規ジャンルの採用審査（Phase 0–4 の完了後）

Phase 5 はロードマップの完了条件に含めません。Phase 0–4 を完遂した後、既存トラックで教えにくい「計算の型」があるときだけ採用を検討します。候補の例:

- Stealth（視界・気づき）
- Farming / 日次ループ
- Autobattler（編成とシナジー）
- Twin-stick（360° 照準の別解）

有名ジャンルの穴埋めのためだけには増やしません。

### やらないこと（当面）

- プレイ可能ゲートを数百単位で水増しすること
- 有名作品のアート・音・名称・マップの再現
- `dist/` の手編集、生成 HTML の場当たりパッチ
- ライセンス不明な外部アセットの導入

進捗の詳細な作業メモやエージェント向け契約は `AGENTS.md` と `docs/ADVANCED_QUALITY_CHECKLIST.md` を更新して追従します。この README のロードマップは、四半期ごとの方針合わせの入口として使います。

## License

Ebi Showcase 自体のコード・文章・オリジナルアセットは Apache License 2.0 です。OFL や BSD など Apache-2.0 と共存できる第三者依存・フォント・アセットは、元のライセンスを維持し、著作権表示と出典を [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md) に記録した場合に限り使用します。
