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
- `web/` — 展示ページとWASMローダー
- `scripts/build.sh` — `dist/` を生成するビルド
- `.github/workflows/pages.yml` — GitHub Pagesデプロイ

## 操作

クリック、タップ、スペースキー、または上矢印キーで羽ばたきます。ゲームオーバー後に同じ操作でリトライできます。

## License

Apache License 2.0
