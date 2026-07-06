#!/usr/bin/env bash
#
# スティッキーな PR コメントを投稿/更新する。マーカー <!-- article-preview --> を含む
# 既存コメントがあれば PATCH で更新し、なければ新規 POST する (再実行でコメントが増えない)。
# GitHub REST API を curl で叩くだけなので第三者アクションに依存しない。
#
# 入力 (環境変数):
#   GITHUB_TOKEN      必須  pull-requests: write が必要
#   GITHUB_API_URL    任意  既定 https://api.github.com
#   GITHUB_REPOSITORY 必須  owner/repo
#   PR                必須  PR 番号
#   COMMENT_FILE      必須  コメント本文 Markdown のファイル
#   MARKER            任意  既定 <!-- article-preview -->
#
set -euo pipefail

: "${GITHUB_TOKEN:?}" "${GITHUB_REPOSITORY:?}" "${PR:?}" "${COMMENT_FILE:?}"
API="${GITHUB_API_URL:-https://api.github.com}"
MARKER="${MARKER:-<!-- article-preview -->}"

auth=(-H "Authorization: Bearer ${GITHUB_TOKEN}" -H "Accept: application/vnd.github+json")

# 既存のスティッキーコメント ID を探す
existing_id="$(
  curl -fsSL "${auth[@]}" "${API}/repos/${GITHUB_REPOSITORY}/issues/${PR}/comments?per_page=100" \
    | jq -r --arg m "${MARKER}" 'map(select(.body | contains($m))) | (.[0].id // empty)'
)"

# 本文を JSON にエンコード
payload="$(jq -Rs '{body: .}' < "${COMMENT_FILE}")"

if [[ -n "${existing_id}" ]]; then
  curl -fsSL -X PATCH "${auth[@]}" \
    "${API}/repos/${GITHUB_REPOSITORY}/issues/comments/${existing_id}" \
    -d "${payload}" > /dev/null
  echo "Updated comment ${existing_id}"
else
  curl -fsSL -X POST "${auth[@]}" \
    "${API}/repos/${GITHUB_REPOSITORY}/issues/${PR}/comments" \
    -d "${payload}" > /dev/null
  echo "Created new comment"
fi
