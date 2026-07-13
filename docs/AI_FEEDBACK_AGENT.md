# ローカルAIフィードバック巡回

`scripts/ai_feedback_crawler.py` は、公開中のEbi Showcaseページを巡回し、
Mac上のLM Studio（OpenAI互換API）で具体的な改善提案を1件ずつ生成する補助ツールです。

デフォルトはドライランです。提案を表示するだけで、フォームへは送信しません。
外部フォームへ送るときだけ `--submit` を付けてください。

## 前提

- LM Studioのローカルサーバーを `http://127.0.0.1:1234` で起動する
- モデルは `google/gemma-4-31b-qat`
- Python 3.10以降（標準ライブラリだけで動作）

## まず1ページ分だけ確認

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 1
```

日本語と英語の入口から内部リンクをたどります。提案は `.cache/ai-feedback-agent/state.sqlite3`
に保存されます。巡回キューも同じDBに保存するため、`--max-pages 24`でも次の起動・次の周期は
続きのページから進みます。同じページが変わらない限り毎回モデルへ投げ直しません。

## 実際にフォームへ送る

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 5 --submit
```

`--submit` はGoogle Formsの `formResponse` 以外へは送信しません。送信先はページ内の
フィードバックフォームから読み取り、提案は200文字未満に切り詰めます。

全ページに共通するレビュー方針は `--instruction` で自由に追加できます。

```sh
python3 scripts/ai_feedback_crawler.py --once --max-pages 5 \
  --instruction="英文のスペルミスに気をつけてレビューしてください"
```

この指示はページ本文とは別の運用者コンテキストとして、Gemmaへの各ページリクエストに毎回挿入されます。

## 定期巡回

```sh
python3 scripts/ai_feedback_crawler.py --max-pages 24 --interval-seconds 600 --submit
```

停止は `Ctrl-C` です。サイトやフォームに負荷をかけないよう、既定でページ間1秒、巡回間隔10分を
設定しています。まずは `--submit` なしで内容を確認し、重複や不要な提案がないことを見てから有効化してください。

Gemma 4の推論トークンで短い回答が途切れないよう、LM Studioへ
`chat_template_kwargs.enable_thinking=false` と `reasoning_effort=none` を渡しています。
モデルやLM Studioのバージョンがこのオプションを知らない場合は、`--model` と
`--lm-base-url` を指定した別のOpenAI互換モデルへ切り替えてください。

ページ本文はモデルへのプロンプトに「信頼できない教材テキスト」として渡し、ページ内に埋め込まれた
指示を実行しないようにしています。APIキーや個人情報はプロンプトにもフォームにも入れません。
