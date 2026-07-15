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

**はじめに覚えること:** `Update` が入力と状態の更新をすべて行い、`Draw` はいまの `game` を画面へ投影するだけです。画面は常に `game` に追従し、`Draw` の中で `game` を書き換えません。

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
- `docs/AUTHORING_CHECKLIST.md` — 学習者が RULE を書いて確かめられるか
- `docs/quality-gates/` — PLAYABLE / AUTHORING / ADVANCED の機械可読ゲート正本
- `docs/ROADMAP_RALPH_LOOP.md` — Authoring Pass の順序付きチェックリスト
- `scripts/build.sh` — `dist/` を生成するビルド（生成 HTML の更新を含む）
- `scripts/ralph-loop.sh` — カリキュラム進捗と検証
- `scripts/roadmap-ralph-loop.mjs` — 証跡付きロードマップ進捗と完了判定
- `scripts/feedback-sheet.mjs` — フォーム回答のトリアージ
- `scripts/ai_feedback_crawler.py` / `docs/AI_FEEDBACK_AGENT.md` — LM Studio または LAN 上の Ollama で改善提案を巡回生成（送信は明示指定時のみ）
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

## ロードマップ（Authoring Pass）

量のゲート（208 + VFX 29）と、技術ラボ／卒業の *器* は到達済みです。いまの正本は **著者が Go を書けるようにする Authoring Pass** です。

実行順・詳細仕様・完了条件は [`docs/ROADMAP_RALPH_LOOP.md`](docs/ROADMAP_RALPH_LOOP.md) だけを見てください（63項目、上から1件ずつ）。

```sh
node scripts/roadmap-ralph-loop.mjs status       # フェーズ別進捗
node scripts/roadmap-ralph-loop.mjs next         # 次の1件
node scripts/roadmap-ralph-loop.mjs evidence ID  # 証跡ファイルを作成
node scripts/roadmap-ralph-loop.mjs check ID     # 証跡を検査して完了
node scripts/roadmap-ralph-loop.mjs verify       # 順序と証跡を検証
```

要約:

0. **公理** — `Update`＝入力と状態、`Draw`＝投影のみ、画面は `game` に追従
1. **P0** — 「値を変えて」から「1ルール足す」へ契約を切り替え、不一致ページを棚卸しする
2. **P1** — Build Track（空の窓→当たり判定）と Core 前半の YOUR FIRST RULE
3. **P2** — Rhythm ほか U–Y 生成器を二層コードパネル＋RULE 課題へ
4. **P3** — 卒業 starter を穴あき化し、first-30-minutes を書くルートへ
5. **P4** — Authoring メーター、横展開、抜き取り監査、リリース

Playable と Authoring は別メーターです。新ジャンル追加は本パス完遂後の審査まで凍結します。

### 凍結ルール

この Authoring Pass の間は、既存の 208/208 playable ゲートを増やしたり、
新しいジャンルコースを採用したりしません。新規案は候補として記録するだけに
留め、現在のコースが「ブラウザで遊べる」だけでなく「読者が 1 ルールを書いて
確かめられる」状態になることを優先します。凍結解除と新ジャンルの採否は、P4 の
リリース完了後に別途審査します。

この道筋でくり返すゲームループの約束は3つです。

1. `Update` が入力を読み、スコア・位置・衝突などのゲーム状態を更新する
2. `Draw` はその状態を画面へ写すだけで、ルールや入力を処理しない
3. 同じ状態なら同じ画面になる。`Draw` は状態、乱数、スライスを変更しない

したがって YOUR FIRST RULE は `Update` または `Update` が呼ぶ純粋な関数へ足し、`Draw` ではその結果を見える形にします。

## License

Ebi Showcase 自体のコード・文章・オリジナルアセットは Apache License 2.0 です。OFL や BSD など Apache-2.0 と共存できる第三者依存・フォント・アセットは、元のライセンスを維持し、著作権表示と出典を [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md) に記録した場合に限り使用します。
