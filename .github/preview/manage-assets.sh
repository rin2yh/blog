#!/usr/bin/env bash
#
# スクショ PNG を専用アセットブランチ (既定: pr-preview-assets) に置くための git 処理。
# ソースの作業ツリーと競合しないよう $RUNNER_TEMP に隔離クローンして操作する。
# 配置パス: pr-<PR>/<SHA7>/*.png  (SHA をパスに含めることで raw/Camo キャッシュを回避)
#
# 使い方: manage-assets.sh <publish|cleanup>
#   publish  SHOTS_DIR の PNG を pr-<PR>/<SHA7>/ に置いて push (古い SHA は prune)
#   cleanup  pr-<PR>/ を削除して push (PR クローズ時)
#
# 入力 (環境変数):
#   GITHUB_TOKEN       必須  push 用トークン (contents: write が必要)
#   GITHUB_REPOSITORY  必須  owner/repo
#   GITHUB_SERVER_URL  任意  既定 https://github.com
#   ASSET_BRANCH       任意  既定 pr-preview-assets
#   PR                 必須  PR 番号
#   RUNNER_TEMP        必須  隔離クローンの作業場所
#   SHA7               publish時必須  短縮 SHA
#   SHOTS_DIR          publish時必須  PNG が入ったディレクトリ
#
set -euo pipefail

MODE="${1:?usage: manage-assets.sh <publish|cleanup>}"
: "${GITHUB_TOKEN:?}" "${GITHUB_REPOSITORY:?}" "${PR:?}" "${RUNNER_TEMP:?}"

ASSET_BRANCH="${ASSET_BRANCH:-pr-preview-assets}"
HOST="${GITHUB_SERVER_URL:-https://github.com}"
HOST="${HOST#https://}"
# ASSET_REMOTE はローカルテスト用の上書き口 (通常は未設定)
REMOTE="${ASSET_REMOTE:-https://x-access-token:${GITHUB_TOKEN}@${HOST}/${GITHUB_REPOSITORY}.git}"
ASSET_DIR="${RUNNER_TEMP}/pr-preview-assets"

# アセットブランチをクローン (無ければ: publish は orphan 作成 / cleanup は何もせず終了)
clone_asset_branch() {
  rm -rf "${ASSET_DIR}"
  if git clone --depth 1 --branch "${ASSET_BRANCH}" "${REMOTE}" "${ASSET_DIR}" 2>/dev/null; then
    return 0
  fi
  return 1
}

# ステージ済みの変更があるときだけ commit/push する
commit_and_push() {
  git add -A
  if git diff --cached --quiet; then
    echo "No changes to push"
    return 0
  fi
  git config user.name "github-actions[bot]"
  git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
  git commit -q -m "$1"
  git push "${REMOTE}" "HEAD:${ASSET_BRANCH}"
  echo "Pushed to ${ASSET_BRANCH}"
}

case "${MODE}" in
  publish)
    : "${SHA7:?}" "${SHOTS_DIR:?}"
    if ! clone_asset_branch; then
      echo "Branch ${ASSET_BRANCH} not found; creating orphan"
      mkdir -p "${ASSET_DIR}"
      git -C "${ASSET_DIR}" init -q
      git -C "${ASSET_DIR}" checkout -q --orphan "${ASSET_BRANCH}"
    fi
    cd "${ASSET_DIR}"
    # この PR の古い SHA ディレクトリは削除し、最新のみを残す (肥大化防止)
    if [[ -d "pr-${PR}" ]]; then
      find "pr-${PR}" -mindepth 1 -maxdepth 1 -type d ! -name "${SHA7}" -exec rm -rf {} +
    fi
    mkdir -p "pr-${PR}/${SHA7}"
    cp "${SHOTS_DIR}"/*.png "pr-${PR}/${SHA7}/" 2>/dev/null || true
    commit_and_push "preview: PR #${PR} @ ${SHA7}"
    ;;
  cleanup)
    if ! clone_asset_branch; then
      echo "Branch ${ASSET_BRANCH} does not exist; nothing to clean"
      exit 0
    fi
    cd "${ASSET_DIR}"
    if [[ ! -d "pr-${PR}" ]]; then
      echo "No assets for PR #${PR}; nothing to clean"
      exit 0
    fi
    git rm -qr --ignore-unmatch "pr-${PR}"
    commit_and_push "preview: cleanup PR #${PR}"
    ;;
  *)
    echo "Unknown mode: ${MODE} (expected publish|cleanup)" >&2
    exit 1
    ;;
esac
