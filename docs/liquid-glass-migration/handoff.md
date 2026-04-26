---
title: 引き継ぎチェックリスト
parent: ./README.md
---

# 引き継ぎチェックリスト

過去にハマったポイントをルール化したもの。新規コードや修正時にここを参照する。

## CSS

- ホバーで `background-color` / `box-shadow` を変えるなら `is-scrolling` ガードまでセットで考える。色だけのフィードバックなら不要。背景を変えるなら `<html>.is-scrolling` 中は base に戻すルールを必ず添える（→ [pitfalls #3](./pitfalls.md#3-ホバー後にスクロールすると白い背景が残る)）。
- ライトモードの本文系は `--color-text-muted` / `--color-text-subtle` をアルファ無しの濃いソリッド (`#1e293b` / `#334155`) で。半透明にしない。

## モバイル UI

開閉トグルを実装したら以下 4 点をチェック:

1. 外側タップで閉じる（document 委譲）
2. 背面スクロールロック（`body.nav-open { overflow: hidden }`）
3. Escape で閉じる
4. 内部リンク遷移時の自動 close

これが揃わないと体感的に「効かない」になる（→ [pitfalls #4](./pitfalls.md#4-ハンバーガーメニューが閉じない)）。

## Hugo パイプライン

- **`Fingerprint` は必ずパイプラインの最後**。`Fingerprint → Minify` の順だと SRI ハッシュが未圧縮ファイルのもの、配信は圧縮後となり全ブロックされる（→ [pitfalls #2](./pitfalls.md#2-sri-ハッシュ不一致で-css--js-全ブロック--ライトモード崩壊に見えた)）。
- 開発時は `mise serve`（= 引数なし `hugo server`）。`--baseURL` 上書きは prod ビルド時のみ（→ [pitfalls #1](./pitfalls.md#1-ページネーションがおかしいの真因はサーバー起動引数だった)）。
- Hugo に `hasKey` は無い。キー有無判定は `isset . "key"` を使う。`default` で済むならそれを使う。

## 表記ゆれ

- Hugo の `humanize` は強制大文字化される。タグ・カテゴリの表示は **必ず humanize を打ち消す側で実装**する（既存方針）。`#go` `#slide` は小文字、`#E2E` `#GitHub Actions` は混合ケースで保持。
- ウィジェット見出しは `i18n/ja.toml` の `[categories]` `[tags]` 等で日本語化する。テンプレ内に英語ラベルを直書きしない。

## ビルド・確認フロー

- UI 変更は **chrome-devtools-mcp で `take_screenshot` してビフォーアフター比較してから「完了」と報告**する。MCP 確認なしの自己申告は通らない（過去複数回差し戻しあり）（→ [pitfalls #5](./pitfalls.md#5-デザイン確認の精度)）。
- ビルドは `hugo --gc --minify` で警告ゼロを確認。
- リスト系ページの partial は `partialCached` を使う。variant が異なる場合は suffix で分ける（例: archives は variant `"archives"`）。

## 退避物

- 旧 Stack 関連は `.archive/layouts/`、`.archive/assets-scss/`、`themes/stack/`（submodule）に残置。
- `.gitignore` に `.hugo_cache/` と `.archive/` を追加済み。新たに退避物が増えたら同じく ignore する。
