# Ebi Showcase ドキュメント全体レビュー – 図・イラスト追加提案

**対象**: `web/ja/` と `web/en/` の HTML ドキュメント全体、および関連する Go ソースコードの解説部分。

---

## 1. 共通の図解提案
- **サイトマップ \& カリキュラムフロー図**
  - **目的**: 学習者が全体の位置付けを一目で把握できるようにする。
  - **配置場所**: 各言語のトップページ (`web/ja/index.html`, `web/en/index.html`) の冒頭に、**SVG** もしくは **PNG** のカリキュラムツリーダイアグラムを挿入。各ステップは **クリック可能リンク** にし、対応ページへジャンプできるようにする。

- **用語集概念マップ**
  - **目的**: `Update`, `Draw`, `GeoM`, `SpriteSheet` などの基礎用語を相互関係で可視化し、どのレッスンで初登場・再利用されるかを示す。
  - **配置場所**: `docs/Glossary.md` のトップに埋め込み、右サイドバーに固定表示（スクロール時に常に見える）。
  - **形式**: **Mermaid** の `graph LR` を用いたインラインコードブロック。例:
    ```mermaid
    graph LR
        Update --> Draw
        Draw --> Layout
        Layout --> GeoM
        GeoM --> SpriteSheet
    ```
    GitHub Pages が Mermaid に対応していない場合は **svg/ mermaid-cli** で画像化し、`web/assets/diagrams/` に保存。

- **デバイス別レイアウト比較図**
  - **目的**: PC, タブレット, スマホの表示サイズ違いで UI がどのように変化するか示す。
  - **配置場所**: 各レッスンの **「デモ」** 前に、**3 カラムの比較画像** (`desktop.png`, `tablet.png`, `mobile.png`) を並べる。
  - **作成方法**: Chrome DevTools のデバイスエミュレーションでスクリーンショットを撮り、`web/assets/layout-comparison/` に保存。自動生成スクリプト `scripts/capture-layout.mjs` を追加すると便利。

---

## 2. 個別レッスン別図解提案（抜粋）
| レッスン | 追加すべき図の種類 | 図の具体的内容・目的 |
|----------|-------------------|----------------------|
| **LEVEL 01** (`web/ja/games/update-draw/index.html`) | **ゲームループフロー図** | `Update → Draw → Layout` の呼び出し順序と、フレームレート制御 (`ebiten.RunGame`) の位置づけを示す。
| **LEVEL 02** (`web/ja/games/velocity/index.html`) | **速度ベクトル図** | 速度ベクトル (`vx, vy`) を矢印で表し、`Update` で位置がどのように変わるかを視覚化。
| **LEVEL 04** (`web/ja/games/collision/index.html`) | **矩形衝突図** | 2 つの矩形が交差している様子を赤線でハイライトし、**左上原点** の座標系を併記。
| **LEVEL 06** (`web/ja/games/animation/index.html`) | **スプライトシート構造図** | スプライトシートの **フレーム格子** と、`Anim("walk-side")` が取得するフレーム順序を示す。
| **LEVEL 08** (`web/ja/games/input/index.html`) | **入力マッピング表** | キーボード・タッチ・マウスの対応を書いた表（例: `Space → Jump`、`Touch → Tap`）。
| **LEVEL 10** (`web/ja/games/score/index.html`) | **スコア保存フロー** | `localStorage` に保存 → 読み込み のデータフロー図。
| **VFX STEP 03** (`web/ja/tracks/visual-effects/geom/index.html`) | **GeoM 変換チェーン図** | `Translate → Rotate → Scale` の適用順序と、**左上原点 → 回転中心** の関係を示す。
| **Platformer STEP 02** (`web/ja/tracks/platformer/jump/index.html`) | **ジャンプ軌道イラスト** | 放物線軌道を描き、**重力定数** と **初速度** の関係式 `y = v0*t - 0.5*g*t²` を併記。
| **Puzzle STEP 01** (`web/ja/tracks/puzzle/rotate/index.html`) | **回転原点図** | `GeoM.Rotate` 前後の座標変化を示す、回転中心を示す赤点付きイラスト。
| **Shooter STEP 03** (`web/ja/tracks/shooter/pool/index.html`) | **オブジェクトプール図** | プールされた弾丸が **再利用** される様子を矢印で表す。
| **Rhythm STEP 02** (`web/ja/tracks/rhythm/bpm/index.html`) | **BPM とフレーム数の関係式** | `interval = 60 / BPM` の計算を **時間軸** にプロットし、ビート同期の視覚化。

*上記は一例です。全レッスンについて同様の視覚化を行うことで、**抽象的な説明が具体的なイメージに変わり、学習ロードがスムーズになります。*

---

## 3. 図の作成・管理指針
1. **形式統一**: 可能な限り **SVG** を使用し、CSS で色やサイズを調整できるようにする。SVG はテキストベースなので、`git diff` が分かりやすい。
2. **ファイル命名規則**: `<lang>_<lesson_slug>_<type>.svg` （例: `ja_level01_flow.svg`）で、言語別に同一画像を共有しやすくする。
3. **埋め込み方法**: HTML では `<img src="../../assets/diagrams/ja_level01_flow.svg" alt="ゲームループの流れ" loading="lazy">` とし、`alt` に簡潔な説明文を入れる。ARIA‑ラベルが必要な場合は `<figure>` と `<figcaption>` を併用。
4. **自動生成**: `scripts/gen-diagrams.mjs` を新規作成し、**Mermaid** コードから SVG を生成する仕組みを導入。CI に `npm run gen-diagrams` を組み込めば、図の変更が自動で反映される。
5. **サイズ最適化**: `svgo` を利用して不要な属性・コメントを除去し、ファイルサイズを数KB に抑える。Web ページの **LCP** 向上に寄与。

---

## 4. フィードバック項目（feedback.md へ追記）
以下は **feedback.md** に直接追記すべき項目のサンプルです。実装時はこのブロック全体を **追記** してください（既存内容を上書きしない）。

```markdown
### 図・イラスト追加提案（全体）
- サイトマップ／カリキュラムフロー図をトップページに配置 → 学習進捗の全体感を提供。
- 用語集概念マップを `docs/Glossary.md` に埋め込み → 用語間の関係が視覚化され、横断的学習が容易になる。
- デバイス別レイアウト比較図を各レッスン冒頭に配置 → レスポンシブ対応の実感が得られる。

### 個別レッスン別図提案（抜粋）
- LEVEL 01: ゲームループフロー図 (`ja_level01_flow.svg`).
- LEVEL 04: 矩形衝突図 (`ja_collision_rect.svg`).
- LEVEL 06: スプライトシート構造図 (`ja_spritesheet_structure.svg`).
- VFX STEP 03: GeoM 変換チェーン図 (`ja_vfx_geom_chain.svg`).
- Platformer STEP 02: ジャンプ軌道イラスト (`ja_platformer_jump_arc.svg`).
- Rhythm STEP 02: BPM とフレーム数の関係式図 (`ja_rhythm_bpm_timing.svg`).

### 図管理ガイドライン
- すべての図は `web/assets/diagrams/` に保存し、命名は `\<lang\>_\<slug\>_\<type\>.svg` とする。
- CI に `npm run gen-diagrams`（Mermaid → SVG）と `svgo` 圧縮を組み込み、プルリクエスト時に自動生成・最適化されるようにする。
- `<img>` タグは `loading="lazy"`、`alt` テキストは必ず記述し、アクセシビリティを担保する。
```

---

## 5. 次のアクション
1. **`scripts/gen-diagrams.mjs`** を作成し、Mermaid → SVG の自動変換パイプラインを実装する（Node.js + `mermaid-cli`）。
2. 既存の **HTML** に上記図を埋め込む PR を作成し、**CI** が通るか確認する。
3. `docs/Glossary.md` に概念マップ用 Mermaid ブロックを追加し、ビルド時に自動で SVG 化されるように設定する。
4. `web/assets/diagrams/` ディレクトリをリポジトリに追加し、`.gitignore` に除外設定を入れずに管理する（サイズは数KB 以内に抑える）。
5. 変更後、`bash scripts/ralph-loop.sh verify` でビルド・テストが成功することを確認。

---

## 6. さらなる概念説明の挿入提案
### 6.1 基礎概念の「なぜ」セクション
- **Update と Draw の違い**: 各 LEVEL の冒頭に、`Update` が「状態を変える」こと、`Draw` が「現在の状態を描画する」ことを、シンプルな擬似コードと共に 2 行程度で説明するブロックを追加。例:
  ```go
  // Update: ゲームロジック、入力、物理計算を実行
  // Draw: 現在フレームの状態を画面に描画
  ```
- **座標系の原点と向き**: `GeoM` 系で座標変換を扱うレッスンでは、左上が (0,0) で右方向が X 正、下方向が Y 正であることを、短い図とともに明記。

### 6.2 アルゴリズムのステップバイステップ解説
- **A* パス探索**（`tactics` トラック）: 探索の「開放リスト」「閉鎖リスト」更新手順を 3 つのサブブロックに分け、各ブロックに疑似コードとフローチャートを挿入。
- **オブジェクトプール**（`shooter`）: プールから取得 → 使用 → 返却 のサイクルを円形矢印で示し、ガーベジコレクション削減の理由を簡潔に説明。

### 6.3 データ駆動型設計の可視化
- **JSON レベルデータ**（`stage-data`）: JSON の構造をツリー図で示し、`internal/lessonlogic` がどのようにデータを取り出すかをコードスニペットで例示。
- **VFX パラメータテーブル**: 各 VFX ステップで使用する `speed`, `size`, `color` などのパラメータを表形式でまとめ、実装時にどこを書き換えるかを明示。

### 6.4 ユーザー体験（UX）向上のためのヒント
- **キーボードショートカット表**: `README.md` の「操作方法」セクションに、`Space`, `Arrow`, `Ctrl+R` などのショートカットを表にまとめ、ページ上部にリンクを貼る。
- **タッチ操作ガイド**: モバイル向けレッスン（例: `tap-target`）では、指でのタップ領域サイズ推奨（最低 48dp）と、タップ時の視覚フィードバック（ハイライト）を説明する小段落を追加。

### 6.5 多言語共通の説明パターン
- **用語の統一**: 日本語では **「更新」**、英語では **“Update”** と表記を統一し、ページ冒頭に「※ここでは `Update` はゲームロジックの更新を指す」等の注釈を入れる。
- **数式のローカライズ**: 「速度 = 距離 ÷ 時間」等の基礎数式は、**日本語** と **英語** 両方の記述を並べ、`<code>` タグでインライン表示する。

### 6.6 学習効果測定のための自己チェックリスト
- 各レッスン末尾に **3〜5問** の小テスト（選択肢または記述）を追加し、正解コードへのリンクを `data-embed-source` で提供。学習者は自己確認でき、理解度向上が期待できる。

---

*本提案は、文章だけでなく適切な「図・表・サンプルコード」を間に挟むことで、概念が直感的に把握できるよう設計しています。ぜひご検討ください。*

**対象**: `web/ja/` と `web/en/` の HTML ドキュメント全体、および関連する Go ソースコードの解説部分。

---

## 1. 共通の図解提案
- **サイトマップ & カリキュラムフロー図**
  - **目的**: 学習者が全体の位置付けを一目で把握できるようにする。
  - **配置場所**: 各言語のトップページ (`web/ja/index.html`, `web/en/index.html`) の冒頭に、**SVG** もしくは **PNG** のカリキュラムツリーダイアグラムを挿入。
  - **内容**: `LEVEL 01‑12` → `Visual Effects Lab` → 各 **Genre Track** (例: Platformer, Shooter, Puzzle …) の流れを矢印で示す。各ステップは **クリック可能リンク** にし、対応ページへジャンプできるようにする。

- **用語集概念マップ**
  - **目的**: `Update`, `Draw`, `GeoM`, `SpriteSheet` などの基礎用語を相互関係で可視化し、どのレッスンで初登場・再利用されるかを示す。
  - **配置場所**: `docs/Glossary.md` のトップに埋め込み、右サイドバーに固定表示（スクロール時に常に見える）。
  - **形式**: **Mermaid** の `graph LR` を用いたインラインコードブロック。例:
    ```mermaid
    graph LR
        Update --> Draw
        Draw --> Layout
        Layout --> GeoM
        GeoM --> SpriteSheet
    ```
    GitHub Pages が Mermaid に対応していない場合は **svg/ mermaid-cli** で画像化し、`web/assets/diagrams/` に保存。

- **デバイス別レイアウト比較図**
  - **目的**: PC, タブレット, スマホの表示サイズ違いで UI がどのように変化するか示す。
  - **配置場所**: 各レッスンの **「デモ」** 前に、**3 カラムの比較画像** (`desktop.png`, `tablet.png`, `mobile.png`) を並べる。
  - **作成方法**: Chrome DevTools のデバイスエミュレーションでスクリーンショットを撮り、`web/assets/layout-comparison/` に保存。自動生成スクリプト `scripts/capture-layout.mjs` を追加すると便利。

---

## 2. 個別レッスン別図解提案（抜粋）
| レッスン | 追加すべき図の種類 | 図の具体的内容・目的 |
|----------|-------------------|----------------------|
| **LEVEL 01** (`web/ja/games/update-draw/index.html`) | **ゲームループフロー図** | `Update → Draw → Layout` の呼び出し順序と、フレームレート制御 (`ebiten.RunGame`) の位置づけを示す。 |
| **LEVEL 02** (`web/ja/games/velocity/index.html`) | **速度ベクトル図** | 速度ベクトル (`vx, vy`) を矢印で表し、`Update` で位置がどのように変わるかを視覚化。 |
| **LEVEL 04** (`web/ja/games/collision/index.html`) | **矩形衝突図** | 2 つの矩形が交差している様子を赤線でハイライトし、**左上原点** の座標系を併記。 |
| **LEVEL 06** (`web/ja/games/animation/index.html`) | **スプライトシート構造図** | スプライトシートの **フレーム格子** と、`Anim("walk-side")` が取得するフレーム順序を示す。 |
| **LEVEL 08** (`web/ja/games/input/index.html`) | **入力マッピング表** | キーボード・タッチ・マウスの対応を書いた表（例: `Space → Jump`、`Touch → Tap`）。 |
| **LEVEL 10** (`web/ja/games/score/index.html`) | **スコア保存フロー** | `localStorage` に保存 → 読み込み のデータフロー図。 |
| **VFX STEP 03** (`web/ja/tracks/visual-effects/geom/index.html`) | **GeoM 変換チェーン図** | `Translate → Rotate → Scale` の適用順序と、**左上原点 -> 回転中心** の関係を示す。 |
| **Platformer STEP 02** (`web/ja/tracks/platformer/jump/index.html`) | **ジャンプ軌道イラスト** | 放物線軌道を描き、**重力定数** と **初速度** の関係式 `y = v0*t - 0.5*g*t²` を併記。 |
| **Puzzle STEP 01** (`web/ja/tracks/puzzle/rotate/index.html`) | **回転原点図** | `GeoM.Rotate` 前後の座標変化を示す、回転中心を示す赤点付きイラスト。 |
| **Shooter STEP 03** (`web/ja/tracks/shooter/pool/index.html`) | **オブジェクトプール図** | プールされた弾丸が **再利用** される様子を矢印で表す。 |
| **Rhythm STEP 02** (`web/ja/tracks/rhythm/bpm/index.html`) | **BPM とフレーム数の関係式** | `interval = 60 / BPM` の計算を **時間軸** にプロットし、ビート同期の視覚化。 |

*上記は一例です。全レッスンについて同様の視覚化を行うことで、**抽象的な説明が具体的なイメージに変わり、学習ロードがスムーズになります。*

---

## 3. 図の作成・管理指針
1. **形式統一**: 可能な限り **SVG** を使用し、CSS で色やサイズを調整できるようにする。SVG はテキストベースなので、`git diff` が分かりやすい。
2. **ファイル命名規則**: `<lang>_<lesson_slug>_<type>.svg` （例: `ja_level01_flow.svg`）で、言語別に同一画像を共有しやすくする。
3. **埋め込み方法**: HTML では `<img src="../../assets/diagrams/ja_level01_flow.svg" alt="ゲームループの流れ" loading="lazy">` とし、`alt` に簡潔な説明文を入れる。ARIA‐ラベルが必要な場合は `<figure>` と `<figcaption>` を併用。
4. **自動生成**: `scripts/gen-diagrams.mjs` を新規作成し、**Mermaid** コードから SVG を生成する仕組みを導入。CI に `npm run gen-diagrams` を組み込めば、図の変更が自動で反映される。
5. **サイズ最適化**: `svgo` を利用して不要な属性・コメントを除去し、ファイルサイズを数KB に抑える。Web ページの **LCP** 向上に寄与。

---

## 4. フィードバック項目（feedback.md へ追記）
以下は **feedback.md** に直接追記すべき項目のサンプルです。実装時はこのブロック全体を **追記** してください（既存内容を上書きしない）。

```markdown
### 図・イラスト追加提案（全体）
- サイトマップ／カリキュラムフロー図をトップページに配置 → 学習進捗の全体感を提供。
- 用語集概念マップを `docs/Glossary.md` に埋め込み → 用語間の関係が視覚化され、横断的学習が容易になる。
- デバイス別レイアウト比較図を各レッスン冒頭に配置 → レスポンシブ対応の実感が得られる。

### 個別レッスン別図提案（抜粋）
- LEVEL 01: ゲームループフロー図 (`ja_level01_flow.svg`).
- LEVEL 04: 矩形衝突図 (`ja_collision_rect.svg`).
- LEVEL 06: スプライトシート構造図 (`ja_spritesheet_structure.svg`).
- VFX STEP 03: GeoM 変換チェーン図 (`ja_vfx_geom_chain.svg`).
- Platformer STEP 02: ジャンプ軌道イラスト (`ja_platformer_jump_arc.svg`).
- Rhythm STEP 02: BPM とフレーム数の関係式図 (`ja_rhythm_bpm_timing.svg`).

### 図管理ガイドライン
- すべての図は `web/assets/diagrams/` に保存し、命名は `<lang>_<slug>_<type>.svg` とする。
- CI に `npm run gen-diagrams`（Mermaid → SVG）と `svgo` 圧縮を組み込み、プルリクエスト時に自動生成・最適化されるようにする。
- `<img>` タグは `loading="lazy"`、`alt` テキストは必ず記述し、アクセシビリティを担保する。
```

---

## 5. 次のアクション
1. **`scripts/gen-diagrams.mjs`** を作成し、Mermaid → SVG の自動変換パイプラインを実装する（Node.js + `mermaid-cli`）。
2. 既存の **HTML** に上記図を埋め込む PR を作成し、**CI** が通るか確認する。
3. `docs/Glossary.md` に概念マップ用 Mermaid ブロックを追加し、ビルド時に自動で SVG 化されるように設定する。
4. `web/assets/diagrams/` ディレクトリをリポジトリに追加し、`.gitignore` に除外設定を入れずに管理する（サイズは数KB 以内に抑える）。
5. 変更後、`bash scripts/ralph-loop.sh verify` でビルド・テストが成功することを確認。

---

*本提案は 2026‑07‑14 時点のコードベースに基づき、**学習者が文章だけでなく視覚情報からも理解できる** ように設計しています。ぜひご検討ください。*
