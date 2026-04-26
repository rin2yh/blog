---
title: 修正履歴の時系列
parent: ./README.md
---

# 修正履歴の時系列

各 Phase は時系列に並んでいる。各行は「指摘 → 原因 / 対処」の 1 行サマリ。深掘りは [pitfalls.md](./pitfalls.md) を参照。

## Phase 1: 初期実装〜公開（feat: liquid-glass-ui, af2a1b5）

| # | 指摘 | 原因 / 対処 |
| - | - | - |
| 1 | ページネーション 404 | `hugo server --baseURL` 引数。引数なしに変更（→ pitfalls #1） |
| 2 | 記事カードに `<a>` ネストエラー | `article-card` 内側のリンクを div に変更 |
| 3 | term ページが空 | taxonomy の kind 分岐を `_default/term.html` で実装 |
| 4 | `.glass-btn` 内 svg がはみ出る | サイズ指定追加 |
| 5 | Archives widget が全件出る | `first 5` で件数制限 |

## Phase 2: ライトモード可読性

| # | 指摘 | 原因 / 対処 |
| - | - | - |
| 6 | 記事詳細の文字が薄い | `_tokens.css` で light の `--color-text-muted` / `--color-text-subtle` を `#1e293b` / `#334155` のソリッドに |
| 7 | TOC・共有ラベルも薄い | `_utils.css` で light モード時 `opacity: 1` を強制 |
| 8 | 全体的に薄字をやめたい | 個別箇所（タイトル / 見出し / メタ / widget）を `#020617` 系に固定 |
| 9 | Fingerprint/Minify 順序 → CSS 全ブロック | パイプラインを `Minify → Fingerprint` に修正（→ pitfalls #2） |

## Phase 3: 一覧 / カード / ホバー挙動

| # | 指摘 | 原因 / 対処 |
| - | - | - |
| 10 | 記事カードのホバー浮き上がり不要 | `translateY` / `scale` の hover を全削除、`GlassParallax` JS も撤去 |
| 11 | RSS 等のリンクが 2 箇所重複 | サイドバーから片方を削除 |
| 12 | メニューバーが記事下に潜る | z-index 修正 |
| 13 | タグチップの高さガタガタ | `tag-cloud` ウィジェットが出現数で `font-size` を 0.81〜1.4rem 動的設定していた → 均一サイズに |
| 14 | ホーム / 一覧 / アーカイブの構造バラバラ | 共通 `article-item.html` に統一、アーカイブは `showReadingTime=false` で再利用 |
| 15 | アーカイブにカテゴリ・タグが出ない | `article-item` パーシャルで両方表示するよう統一 |
| 16 | 外部リンク記事のリンクアイコンが消えた | `helper/external-link.html` を新設し 3 箇所のロジックを集約 |
| 17 | カテゴリ表示でカウントが邪魔 | カウント削除、タグと同じチップスタイルに |
| 18 | ヒーロー不要 | ホームのヒーロー削除 |
| 19 | 長タイトル時のレイアウト崩れ | アイコンとタイトルを同一 `<a>` 配下に置き、折り返し時アイコン取り残しを解消 |

## Phase 4: ナビゲーション / トグル / 検索

| # | 指摘 | 原因 / 対処 |
| - | - | - |
| 20 | ダーク/ライトトグル位置が微妙 | サイドバー底固定 → ブログ名の右隣 → 最終的に QR コード直下に |
| 21 | トグルのラベル変化で状態がわかりづらい | ラベルは「ダークモード」固定、サムの左右位置で ON/OFF を表現 |
| 22 | input 検索アイコンが縦中央にない | flex で center 揃え |
| 23 | 検索窓の値が検索ページに引き継がれない | クエリパラメータ送出を修正 |
| 24 | 検索ボックスが地味で見えづらい | padding / font-size 拡大、border・opacity 強化、重複 `<h3>` 削除 |
| 25 | モバイルでメニューが閉じない | バックドロップ + 外側クリック + scroll lock 追加（→ pitfalls #4） |

## Phase 5: 一覧の右サイドバー

| # | 指摘 | 原因 / 対処 |
| - | - | - |
| 26 | アーカイブ等で右サイドバーが固定されない（左は固定） | `.layout-grid__sidebar` に `position: sticky; top; max-height; overflow-y: auto;`。1024px 以下では `static` に戻す（`e62d00c`） |

## Phase 6: テーマ整理 / 表記ゆれ

| # | 指摘 | 原因 / 対処 |
| - | - | - |
| 27 | テーマ名 `liquid-glass` → `stack-liquid-glass` | ディレクトリリネーム + `hugo.toml` 切替 |
| 28 | タグの `#go` が `#Go` になる | Hugo の `humanize` を打ち消すよう article-item / article-card / article/header / widgets/tag-cloud で表示制御。`#E2E` `#GitHub Actions` 等の混合ケースは保持 |
| 29 | カテゴリラベルが大文字化される | `.article-item__categories .glass-category { text-transform: none }` |
| 30 | ウィジェット見出しが英語 (Categories / Tags) | `i18n/ja.toml` に `[categories]` `[tags]` 追加 |
| 31 | 各ページの英語タイトル | About / Archives / Search / QR Code → 自己紹介 / アーカイブ / 検索 / QR コード |
| 32 | about / QR ページでホバーエフェクトが効いてしまう | `page/single.html` に `article` クラス追加して既存オーバーライドを適用 |
| 33 | アーカイブウィジェット幅が他と揃わない | 外側 `<aside>` に `.glass.widget`、リストを `margin: 0 calc(-1 * var(--space-lg))` でカード端まで広げる |
| 34 | アーカイブウィジェットからアーカイブ該当年に飛ばない | ハッシュリンクを `#year-YYYY` に統一 |

## Phase 7: ホバー残像（最終対応）

| # | 指摘 | 原因 / 対処 |
| - | - | - |
| 35 | ホバー後スクロールで白カード残像 | ScrollHoverGuard 導入 → それでも残ったため `.article-item:hover` の bg/shadow 完全削除（→ pitfalls #3） |

## Phase 8: /simplify レビュー対応

`/simplify` 走行で並列レビュー agent から指摘されたものを反映:

- `partial "article-card.html"` を 4 箇所で `partialCached` 化（initial 0f5e985 で list/taxonomy/index/related、後追いで archives は variant suffix `"archives"` で別キャッシュ）。
- Google Fonts を `media="print" onload="this.media='all'"` で非同期化、レンダリングブロック解消。
- `.article-card` parallax を rAF + rect キャッシュ + 変化検出で最適化。`getBoundingClientRect()` は mouseenter のみに。
- inline `onclick` を `data-copy` 属性 + delegated handler に置換（`donate-card__email`）。GA トラッキング ID を `params.googleAnalytics` から読むよう変更。
- `_default/term.html` 削除（`taxonomy.html` の term ブランチと重複）。7 つの widget 冒頭の `reflect.IsMap` シムを削除し dispatcher 側で常に `dict` 渡しに統一。`assets/css/main.css`（無効 SCSS の dead file）削除。
- 重複していた glass capsule CSS 宣言（site-brand / site-social / pagination / theme-toggle）を削除し HTML 側に `.glass` クラスを付与。
- inline `style=""` 6 箇所を `_utils.css` クラス化（`.icon-frame` 等）。`<img>` に `decoding="async"` 追加。Mermaid CDN を `mermaid@10` に固定。`onsubmit="return false;"` を `search.js` の `preventDefault` に。ripple keyframes を JS 動的注入から `_animations.css` に移動。
- `article-item.html` の `showReadingTime: false` バグ: 当初 `default true` → `hasKey` で判定しようとしたが Hugo に `hasKey` は無く初回ビルド失敗、最終的に `cond (isset . "showReadingTime") .showReadingTime true` で確定。
