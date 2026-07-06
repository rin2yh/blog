#!/usr/bin/env bash
#
# PR クローズ時にアセットブランチからこの PR のディレクトリ (pr-<PR>/) を削除する。
#
# 入力 (環境変数): GITHUB_TOKEN, GITHUB_REPOSITORY, GITHUB_SERVER_URL(任意),
#                  ASSET_BRANCH(任意), PR, RUNNER_TEMP
#
set -euo pipefail

: "${GITHUB_TOKEN:?}" "${GITHUB_REPOSITORY:?}" "${PR:?}" "${RUNNER_TEMP:?}"
ASSET_BRANCH="${ASSET_BRANCH:-pr-preview-assets}"
SERVER="${GITHUB_SERVER_URL:-https://github.com}"
HOST="${SERVER#https://}"
REMOTE="https://x-access-token:${GITHUB_TOKEN}@${HOST}/${GITHUB_REPOSITORY}.git"

ASSET_DIR="${RUNNER_TEMP}/pr-preview-assets"
rm -rf "${ASSET_DIR}"

if ! git clone --depth 1 --branch "${ASSET_BRANCH}" "${REMOTE}" "${ASSET_DIR}" 2>/dev/null; then
  echo "Branch ${ASSET_BRANCH} does not exist; nothing to clean"
  exit 0
fi

cd "${ASSET_DIR}"
if [[ ! -d "pr-${PR}" ]]; then
  echo "No assets for PR #${PR}; nothing to clean"
  exit 0
fi

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
git rm -qr --ignore-unmatch "pr-${PR}"
if git diff --cached --quiet; then
  echo "No changes to commit"
  exit 0
fi
git commit -q -m "preview: cleanup PR #${PR}"
git push "${REMOTE}" "HEAD:${ASSET_BRANCH}"
echo "Removed pr-${PR} from ${ASSET_BRANCH}"
