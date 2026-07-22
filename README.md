# Ebi Showcase

Go + [Ebitengine](https://ebitengine.org/) で作られたミニゲーム集です。
リポジトリをcloneしたら、好きなゲームの `main.go` を直接変更し、デスクトップの
ゲームウィンドウですぐに結果を確認できます。

公開版では、同じゲームをWebAssemblyとしてブラウザから遊べます。

- 公開サイト: https://kumagi.github.io/EbiShowcase/
- 使用言語: Go
- ゲームエンジン: Ebitengine
- ライセンス: Apache-2.0

## まず1本動かす

必要なものはGitとGo 1.25以降です。

```sh
git clone https://github.com/kumagi/EbiShowcase.git
cd EbiShowcase
go run ./games/core/tap-target
```

初回だけGoが依存パッケージをダウンロードします。ゲームウィンドウが開いたら、
丸をクリックしてください。終了はターミナルで `Ctrl+C` です。

macOSでコンパイラを求められた場合は `xcode-select --install` を実行します。
Linuxではウィンドウ・OpenGL・音声の開発パッケージが必要です。OS別の正確な手順は
[Ebitengine公式インストールガイド](https://ebitengine.org/ja/documents/install.html)
を参照してください。

## 最初の変更を試す

最初に [`games/core/tap-target/main.go`](games/core/tap-target/main.go) を開きます。
たとえば次のどちらかを変更してみてください。

```go
startSeconds = 30 // 10にすると10秒ゲームになる
```

```go
g.score++ // g.score += 2 にすると1回で2点になる
```

保存したら、もう一度実行します。

```sh
go run ./games/core/tap-target
```

Ebi Showcaseに専用のエディタやプロジェクト生成操作はありません。基本の開発ループは
「`main.go` を保存 → `go run` を再実行」だけです。

## 遊びたいゲームを選ぶ

各ゲームは独立したGoの `main` パッケージです。実行したいディレクトリを
`go run` に渡します。

| ゲーム | 実行コマンド | 向いている変更 |
| --- | --- | --- |
| タップターゲット | `go run ./games/core/tap-target` | 得点、制限時間、円の大きさ |
| 小さなプラットフォーマー | `go run ./games/tracks/platformer/tiny-platformer` | 重力、ジャンプ、足場配置 |
| スプライトアニメーション | `go run ./games/tracks/visual-effects/vfx-walk` | コマ送り、速度、描画位置 |
| ガード・投げ・打撃 | `go run ./games/tracks/fighting/guard-throw` | 当たり判定、技の優先順位 |
| Ebi Quest | `go run ./games/tracks/rpg/ebi-quest` | マップ、クエスト、コマンド戦闘 |
| Ebi Merge | `go run ./games/tracks/merge-physics/ebi-merge` | 円の物理、合成、スコア |

macOSやLinuxで実行可能なゲームを一覧にするには、リポジトリのルートで次を実行します。

```sh
find games -name main.go -print
```

PowerShellでは次のコマンドを使えます。

```powershell
Get-ChildItem games -Recurse -Filter main.go
```

表示された `main.go` の親ディレクトリが `go run` に渡すパスです。

## 自分用のゲームを作る

既存の小さなゲームをコピーすると、空のプロジェクトから始める必要がありません。

```sh
cp -R games/core/tap-target games/my-first-game
go run ./games/my-first-game
```

まずウィンドウタイトル、色、得点ルールを変更し、その後に新しい状態や入力を
足していくのがおすすめです。コピー後もリポジトリ内にあるため、`internal/` の
共有コードや既存アセットをそのまま利用できます。

## コードを読む場所

ほとんどのゲームは、1本の `main.go` に次の5か所があります。

1. `type game struct` — 位置、得点、HPなど、現在のゲーム状態
2. `newGame()` — 最初の状態とステージデータ
3. `Update()` — 入力、移動、衝突、得点などのルール
4. `Draw()` — 現在の状態を画面へ描く処理
5. `main()` — ウィンドウを作ってゲームを起動する入口

ゲームの状態を変える処理は `Update()` またはそこから呼ぶ関数へ書き、`Draw()` は
状態を画面へ描くだけにすると、ルールを追いやすくなります。

画像を使うゲームでは、同じディレクトリの `assets/` や `//go:embed` の指定も確認して
ください。複数ゲームで共有するロジックやアートは [`internal/`](internal/) にあります。

## 変更したゲームを確かめる

まず変更したゲームだけを整形・テストします。

```sh
gofmt -w games/tracks/rpg/ebi-quest/main.go
go test ./games/tracks/rpg/ebi-quest
go run ./games/tracks/rpg/ebi-quest
```

テストファイルがないゲームでも、`go test` はそのパッケージがコンパイルできることを
確認します。リポジトリ全体ではなく、編集中のゲームのパスを指定すると短時間で回せます。

ルールを画面や音声から分けた純粋関数には、同じディレクトリへ `main_test.go` を追加できます。

```go
package main

import "testing"

func TestScoreForHit(t *testing.T) {
    if got := scoreForHit(3); got != 4 {
        t.Fatalf("scoreForHit(3) = %d, want 4", got)
    }
}
```

## ブラウザ版も確認する

普段の編集にはデスクトップ版の `go run` が最速です。WebAssembly、タッチ操作、教材ページ
まで確認したいときだけ、サイト全体をビルドします。追加でNode.jsの現行LTS、Python 3、
Bashが必要です。

```sh
bash scripts/build.sh --fast
python3 -m http.server 8080 --directory dist
```

ブラウザで <http://localhost:8080/> を開きます。WebAssemblyは `file://` から直接開けないため、
必ずローカルHTTPサーバーを使ってください。`build.sh` は全ゲームと全教材を生成するので、
普段の1ゲームの試行より時間がかかります。

## よくある問題

### `go run main.go` では動かない

ファイル1個ではなくパッケージのディレクトリを指定してください。

```sh
go run ./games/tracks/platformer/tiny-platformer
```

同じディレクトリの補助ファイルやOS別ファイルも一緒にコンパイルされます。

### ウィンドウが開かない

まずGoのバージョンを確認します。

```sh
go version
```

次に[Ebitengine公式インストールガイド](https://ebitengine.org/ja/documents/install.html)の
OS別依存パッケージを確認してください。SSHだけのマシンやGUIのないコンテナでは、
デスクトップウィンドウを表示できません。

### 音が鳴らない

ブラウザ版は、ブラウザの自動再生制限により、最初のクリックやキー入力後に音声を開始します。
まずゲーム画面を一度操作してください。

### どこから学べばよいかわからない

公開サイトの[日本語トップページ](https://kumagi.github.io/EbiShowcase/ja/)では、小さなゲームから
ジャンル別の完成ゲームまで順番に遊べます。コード内の用語は
[`docs/Glossary.md`](docs/Glossary.md)でも確認できます。

## リポジトリ運用に参加する場合

このREADMEは、cloneしてゲームコードを動かす人向けです。教材ページの生成、サイト全体の
品質ゲート、公開、フィードバック処理などを変更する場合は、作業前に
[`AGENTS.md`](AGENTS.md)を読んでください。

## License

Ebi Showcaseのコード、文章、オリジナルアセットはApache License 2.0です。
第三者依存・フォント・アセットの表示は
[`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md)に記載しています。
