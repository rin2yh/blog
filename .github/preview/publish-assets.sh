#!/usr/bin/env bash
#
# スクショ PNG を専用アセットブランチ (既定: pr-preview-assets) に push する。
# ソースの作業ツリーとは競合しないよう $RUNNER_TEMP に隔離クローンして操作する。
# 配置パス: pr-<PR>/<SHA7>/*.png  (SHA をパスに含めることで raw/Camo キャッシュを回避)
#
# 入力 (環境変数):
#   GITHUB_TOKEN        必須  push 用トークン (contents: write が必要)
#   GITHUB_REPOSITORY   必須  owner/repo
#   GITHUB_SERVER_URL   任意  既定 https://github.com
#   ASSET_BRANCH        任意  既定 pr-preview-assets
#   PR                  必須  PR 番号
#   SHA7                必須  短縮 SHA
#   SHOTS_DIR           必須  PNG が入ったディレクトリ
#   RUNNER_TEMP         必須  隔離クローンの作業場所
#
set -euo pipefail

: "${GITHUB_TOKEN:?}" "${GITHUB_REPOSITORY:?}" "${PR:?}" "${SHA7:?}" "${SHOTS_DIR:?}" "${RUNNER_TEMP:?}"
ASSET_BRANCH="${ASSET_BRANCH:-pr-preview-assets}"
SERVER="${GITHUB_SERVER_URL:-https://github.com}"
HOST="${SERVER#https://}"
REMOTE="https://x-access-token:${GITHUB_TOKEN}@${HOST}/${GITHUB_REPOSITORY}.git"

ASSET_DIR="${RUNNER_TEMP}/pr-preview-assets"
rm -rf "${ASSET_DIR}"

if git clone --depth 1 --branch "${ASSET_BRANCH}" "${REMOTE}" "${ASSET_DIR}" 2>/dev/null; then
  echo "Cloned existing branch ${ASSET_BRANCH}"
else
  echo "Branch ${ASSET_BRANCH} not found; creating orphan"
  mkdir -p "${ASSET_DIR}"
  git -C "${ASSET_DIR}" init -q
  git -C "${ASSET_DIR}" checkout -q --orphan "${ASSET_BRANCH}"
fi

cd "${ASSET_DIR}"
git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

# この PR の古い SHA ディレクトリは削除し、最新のみを残す (ブランチ肥大化防止)
if [[ -d "pr-${PR}" ]]; then
  find "pr-${PR}" -mindepth 1 -maxdepth 1 -type d ! -name "${SHA7}" -exec rm -rf {} +
fi

DEST="pr-${PR}/${SHA7}"
mkdir -p "${DEST}"
cp "${SHOTS_DIR}"/*.png "${DEST}/" 2>/dev/null || true

git add -A
if git diff --cached --quiet; then
  echo "No asset changes to push"
  exit 0
fi

git commit -q -m "preview: PR #${PR} @ ${SHA7}"
git push "${REMOTE}" "HEAD:${ASSET_BRANCH}"
echo "Pushed assets to ${ASSET_BRANCH}/${DEST}"
