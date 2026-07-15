# ローカルAIフィードバック巡回

`scripts/ai_feedback_crawler.py` は、公開中のEbi Showcaseページを巡回し、
OpenAI互換のチャットAPIで具体的な改善提案を1件ずつ生成する補助ツールです。

対応バックエンド:

- **LM Studio**（既定）— `http://127.0.0.1:1234/v1`
- **Ollama** — LAN上のホスト（例: `192.168.3.56:11434`）の `/v1`

デフォルトはドライランです。提案を表示するだけで、フォームへは送信しません。
外部フォームへ送るときだけ `--submit` を付けてください。

## 前提

- Python 3.10以降（標準ライブラリだけで動作）
- LM Studio を使う場合: ローカルサーバーを `http://127.0.0.1:1234` で起動し、モデル例は `google/gemma-4-31b-qat`
- Ollama を使う場合: 対象ホストで `ollama serve` が外部から到達でき、モデルが pull 済みであること

## Ollama（LAN）を使う

ホストとポートを指定します（ポート省略時は `11434`）。

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 1 \
  --ollama-host 192.168.3.56 \
  --model qwen3.6:latest
```

`host:port` でも、環境変数でも同じです。

```sh
export OLLAMA_HOST=192.168.3.56:11434
python3 scripts/ai_feedback_crawler.py --once --max-pages 1 --model qwen3.6:latest
```

または OpenAI 互換 URL を直接渡します（ポート `11434` なら provider は自動で ollama）。

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 1 \
  --lm-base-url http://192.168.3.56:11434/v1 \
  --provider ollama \
  --model qwen3.6:latest
```

`--model` を省略し、かつホストにモデルが1つだけのときは、そのモデルを自動選択します。

## LM Studio（localhost）を使う

引数なしが従来どおりです。

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 1
```

## まず1ページ分だけ確認

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 1 \
  --ollama-host 192.168.3.56 --model qwen3.6:latest
```

日本語と英語の入口から内部リンクをたどります。提案は `.cache/ai-feedback-agent/state.sqlite3`
に保存されます。巡回キューも同じDBに保存するため、`--max-pages 24`でも次の起動・次の周期は
続きのページから進みます。同じページが変わらない限り毎回モデルへ投げ直しません。

## 実際にフォームへ送る

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 5 --submit \
  --ollama-host 192.168.3.56 --model qwen3.6:latest
```

`--submit` はGoogle Formsの `formResponse` 以外へは送信しません。送信先と、Google Formsが
必須にしている「対象ページ」フィールドはページ内のフォームから読み取ります。提案は200文字未満に
切り詰め、ページのフォームと同じ隠し項目も送信します。

全ページに共通するレビュー方針は `--instruction` で自由に追加できます。

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 5 \
  --ollama-host 192.168.3.56 --model qwen3.6:latest \
  --instruction="英文のスペルミスに気をつけてレビューしてください"
```

この指示はページ本文とは別の運用者コンテキストとして、各ページリクエストに毎回挿入されます。

## 定期巡回

```sh
python3 scripts/ai_feedback_crawler.py --max-pages 24 --interval-seconds 600 --submit \
  --ollama-host 192.168.3.56 --model qwen3.6:latest
```

停止は `Ctrl-C` です。サイトやフォームに負荷をかけないよう、既定でページ間1秒、巡回間隔10分を
設定しています。まずは `--submit` なしで内容を確認し、重複や不要な提案がないことを見てから有効化してください。

LM Studio の Gemma 4 向けには `chat_template_kwargs.enable_thinking=false` と
`reasoning_effort=none` を付けます。Ollama 接続時はネイティブの `/api/chat` に
`think=false` を付けて最終回答だけを受け取ります（OpenAI 互換 `/v1` だと
thinking モデルで `content` が空になりやすいため）。

ページ本文はモデルへのプロンプトに「信頼できない教材テキスト」として渡し、ページ内に埋め込まれた
指示を実行しないようにしています。APIキーや個人情報はプロンプトにもフォームにも入れません。

## 品質ゲートのレンズ（推奨）

巡回は自由感想より、[`docs/quality-gates/catalog.json`](quality-gates/catalog.json)
の LLM ゲートを拾ってレビューする方が信用できます。

`--lens` を省略した場合は、起動時に LLM 用ゲートからランダムに2件を選び、
その実行中の全ページを同じ2観点でレビューします。選ばれたIDは
`[lenses:random]` として標準エラーへ表示されます。前回と異なる組み合わせ
なら、ページ本文が同じでも新しい観点として再レビューします。

```sh
# 2レンズをランダム選択
python3 scripts/ai_feedback_crawler.py --once --max-pages 3 --force \
  --ollama-host 192.168.3.56 --model qwen3.6:latest

# 観点を固定
python3 scripts/ai_feedback_crawler.py --once --max-pages 3 --force \
  --ollama-host 192.168.3.56 --model qwen3.6:latest \
  --lens loop,authoring
```

`--lens` は `node scripts/check-quality-gates.mjs --lenses …` と同じ catalog を読み、
各ゲートの `prompt_hint` を OPERATOR 指示へ注入します。モデルには
`gate_id` / `verdict` / `evidence` / `fix` の JSON を求めます（執筆ではなく監査向け）。
