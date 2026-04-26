---
title: Liquid Glass UI Kit 移行ログ
period: 2026-04-25 〜 2026-04-26
related_commits:
  - af2a1b5 feat: liquid-glass-ui
  - 7cdfb2e fix: tag
  - 0f5e985 fix: theme rename, hover residue, tag casing, search box clarity
  - e62d00c fix: stick right sidebar like the left one on list pages
  - 13f3188 fix: drop article-item hover bg to stop white-card residue while scrolling
---

# Liquid Glass UI Kit 移行ログ

Hugo の旧 Stack テーマから Liquid Glass UI Kit ベースの自作テーマ `stack-liquid-glass` に移行した際の、ハマったポイントとバグ修正履歴。会話トランスクリプト（2026-04-25〜2026-04-26、11 セッション）から事後的に再構成したもの。

## 目的別の入口

| 知りたいこと | 読む順 |
| - | - |
| なぜこの構成になったか / 背景 | この README → [pitfalls.md](./pitfalls.md) |
| 過去にどんな指摘・修正があったか | [timeline.md](./timeline.md) |
| 今後同じ罠を踏まないためのルール | [handoff.md](./handoff.md) |

## ドキュメント一覧

- **[pitfalls.md](./pitfalls.md)** — 構造的に詰まった案件 5 件のディープダイブ（原因と解決の詳細）
- **[timeline.md](./timeline.md)** — Phase 1〜8、計 35 項目の指摘 → 原因 → 対処を時系列で
- **[handoff.md](./handoff.md)** — 引き継ぎチェックリスト（CSS / モバイル UI / Hugo パイプライン / タグ表記 / ビルド確認 / 退避物）

## 出発点

- ユーザが Liquid Glass UI Kit の HTML/CSS（`--font-display: Italiana`、`--font-body: DM Sans`、`.glass` / `.glass-btn` / `.glass-card` 等のトークン体系）を貼り付けて「これを Hugo テーマにしたい」と依頼。
- 元テーマ（Stack）は submodule + `.archive/` に退避、新テーマを `themes/liquid-glass/`（後に `themes/stack-liquid-glass/` へリネーム）として一から実装。
- subagent を活用して partial / CSS / JS を並列生成。
