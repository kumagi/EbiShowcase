# Ebi Showcase

Ebitengineで作ったミニゲームを、WebAssemblyでその場で遊べる静的ショーケースです。ゲーム本体はすべてGo + Ebitengineで実装し、HTML/CSS/JavaScriptは展示ページとWASMローダーだけに使っています。

## ローカルで動かす

Go 1.24以降をインストールして、次を実行します。

```sh
go mod download
bash scripts/build.sh
python3 -m http.server 8080 --directory dist
```

ブラウザで `http://localhost:8080` を開きます。WASMは `file://` では読み込めないため、ローカルHTTPサーバーが必要です。

カリキュラムの整合確認は `bash scripts/ralph-loop.sh verify`（Go 1.24+ が必要）です。

SNS向けの OGP は `node scripts/inject-ogp.mjs` と `go run ./cmd/gen-og-images` で全ページ分を生成します（`build.sh` 内でも実行）。公開URLのオリジンは環境変数 `SITE_ORIGIN`（省略時 `https://kumagi.github.io/EbiShowcase`）です。

## GitHub Pagesへ公開する

1. このディレクトリをGitHubリポジトリへpushします。
2. リポジトリの **Settings → Pages → Build and deployment → Source** を **GitHub Actions** にします。
3. `main` ブランチへpushすると `.github/workflows/pages.yml` がビルド・公開します。

プロジェクトサイト（`https://USER.github.io/REPO/`）でも動くよう、アセットはすべて相対パスです。

## 構成

- `game/main.go` — Ebitengine製ゲーム本体
- `web/index.html` — ブラウザ言語を判定する入口
- `web/ja/`・`web/en/` — 日本語・英語のカリキュラム目次
- `web/ja/games/`・`web/en/games/` — 12段階のゲーム教材ページ
- `web/*/games/flappy/` — 現在プレイできるFlappy Bird教材
- `web/ja/tracks/`・`web/en/tracks/` — 15の専門コースと93段階の発展教材
- `web/*/guides/game-data/` — アセットデータとゲーム構成を学ぶ横断ガイド
- `examples/data-driven/` — `go:embed`とJSONによるデータ読み込みの実例
- `docs/Glossary.md` — フレーム、座標、ラジアン、`iota`、スライスなどの用語集
- `scripts/ai_feedback_crawler.py` / `docs/AI_FEEDBACK_AGENT.md` — LM Studioで改善提案を巡回生成（送信は明示指定時のみ）
- `web/game.html` — 両言語で共有するWASMローダー
- `scripts/build.sh` — `dist/` を生成するビルド
- `.github/workflows/pages.yml` — GitHub Pagesデプロイ

言語ページは検索・共有しやすい独立URLです。ルートへの初回アクセスではブラウザ言語を判定し、画面上で言語を切り替えた後はその選択をブラウザに保存します。

## 操作

クリック、タップ、スペースキー、または上矢印キーで羽ばたきます。ゲームオーバー後に同じ操作でリトライできます。

教材のデモは、入力方法を同じゲーム内アクションへ変換します。

| 入力 | アクション | 主な用途 |
| --- | --- | --- |
| Space / ↑ | ジャンプ・決定 | Flappy、横アクション、会話送り |
| ← → ↑ ↓ | 移動・選択 | 迷路、倉庫番、戦略SLG |
| クリック / タップ | 決定・配置・カード選択 | タッチ、タワーディフェンス、デッキ構築 |
| R | リトライ（対応デモ） | 失敗したステージを最初から遊ぶ |

スマホでは、ボタンやタップ判定を最低48dp（CSSではおよそ48px）以上にし、押した瞬間に色や粒で反応を返すのが目安です。

## フィードバックを確認する

各教材ページの末尾にはGoogle Formを埋め込んでいます。回答先のスプレッドシートは、OAuthでログインしたGoogleアカウントから読み書きします。

1. Google CloudでSheets APIを有効にし、OAuthクライアント（デスクトップアプリ）を作成します。
2. ダウンロードしたJSONを `.secrets/` に保存します（このファイルはGitへ入りません）。
3. 次のコマンドを初回だけ実行し、ブラウザでGoogleログインとアクセス許可を行います。

```sh
node scripts/feedback-sheet.mjs list
```

4. 次のコマンドで一覧を取得・更新します。

```sh
node scripts/feedback-sheet.mjs list
node scripts/feedback-sheet.mjs check 12   # 12行目を対応済みにする
node scripts/feedback-sheet.mjs delete 12   # 12行目を削除する
```

`check` は「対応済み」列がなければ自動追加し、✅を記録します。行番号は一覧表示の番号を使ってください。

## License

Apache License 2.0
