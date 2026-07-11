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
- `web/ja/tracks/`・`web/en/tracks/` — 7つの専門コースと43段階の発展教材
- `web/*/guides/game-data/` — アセットデータとゲーム構成を学ぶ横断ガイド
- `examples/data-driven/` — `go:embed`とJSONによるデータ読み込みの実例
- `web/game.html` — 両言語で共有するWASMローダー
- `scripts/build.sh` — `dist/` を生成するビルド
- `.github/workflows/pages.yml` — GitHub Pagesデプロイ

言語ページは検索・共有しやすい独立URLです。ルートへの初回アクセスではブラウザ言語を判定し、画面上で言語を切り替えた後はその選択をブラウザに保存します。

## 操作

クリック、タップ、スペースキー、または上矢印キーで羽ばたきます。ゲームオーバー後に同じ操作でリトライできます。

## License

Apache License 2.0
