#!/usr/bin/env bash
#
# 変更された記事ファイル (content/post/**) を Hugo の URL パス /post/<basename>/ に変換する。
# permalink 設定 `post = '/post/:contentbasename/'` (config/_default/permalinks.toml) に依拠する。
#
#   - リーフバンドル content/post/<slug>/index.md -> /post/<slug>/   (親フォルダ名)
#   - 単一ファイル   content/post/2026/foo.md      -> /post/foo/     (拡張子と言語サフィックスを除去)
#   - セクション    _index.md                     -> スキップ
#   - 削除ファイル                                 -> --diff-filter=d で除外
#
# 入力: 環境変数 BASE_SHA, HEAD_SHA
# 出力: 標準出力に /post/<basename>/ を1行ずつ (重複除去済み)
#
set -euo pipefail

: "${BASE_SHA:?BASE_SHA is required}"
: "${HEAD_SHA:?HEAD_SHA is required}"

git diff --name-only --diff-filter=d "${BASE_SHA}...${HEAD_SHA}" -- 'content/post/**' \
| while IFS= read -r f; do
    [[ -z "$f" ]] && continue
    base="$(basename "$f")"
    case "$base" in
      _index.md|_index.*.md)
        continue ;;                                     # セクションページ (記事一覧) はスキップ
      index.md|index.*.md)
        name="$(basename "$(dirname "$f")")" ;;         # リーフバンドル -> 親フォルダ名
      *.md)
        name="${base%.md}" ;;                           # 単一ファイル -> 拡張子除去 (単一言語サイトのため .md のみ)
      *)
        continue ;;                                     # バンドル内の画像等の非 Markdown はスキップ
    esac
    printf '/post/%s/\n' "$name"
  done | sort -u
