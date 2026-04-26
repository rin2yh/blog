---
title: 構造的にハマった案件
parent: ./README.md
---

# 構造的にハマった案件

表面症状と真因が乖離していて時間を溶かしたケース 5 件。指摘 → 原因 → 対処 → 学びの順。

## 1. 「ページネーションがおかしい」の真因はサーバー起動引数だった

- `mise serve` ではなく手動で `hugo server --baseURL http://localhost:1313/` を起動していたため、Hugo が baseURL を `/` と認識し、生成されるリンクは `/blog/page/2/` を指す不整合が発生。結果として 2 ページ目以降が 404。
- レイアウト不具合に見えていたが、`--baseURL` を外して `hugo server` だけにすれば全ページ 200 OK。

**学び:** 開発時は必ず `mise serve`（= 引数なし `hugo server`）。`baseURL` 上書きは prod ビルド時のみ。

## 2. SRI ハッシュ不一致で CSS / JS 全ブロック → ライトモード崩壊に見えた

- `head.html` で `resources.Get | Fingerprint "sha256" | Minify` の順だった。`Fingerprint` 後に Minify するため、**埋め込まれた integrity ハッシュは未圧縮ファイルの sha256**、しかし配信されるのは Minify 後ファイル。ブラウザが SRI 検証で全部弾いていた。
- 表面症状は「ライトモードで文字色が薄い」「全体的に崩れる」だったので長く犯人が見えなかった。
- 修正: `Minify | Fingerprint "sha256"` に順序入れ替え。

**学び:** Fingerprint は必ずパイプラインの最後。

## 3. ホバー後にスクロールすると白い背景が残る

カードに `.glass` のホバースタイル（`background-color: ... 0.2`）を付与していたが、CSS の `:hover` はトラックパッド慣性スクロール中も維持される。スクロール後にカードの上をマウスが通過した状態でフォーカスが残り、結果として「他のカードはガラスなのに 1 枚だけ白くベタ塗り」という残像になっていた。

3 段階で修正:

1. **ScrollHoverGuard (`liquid-glass.js`)** — スクロール中は `<html>` に `is-scrolling` クラスを付け、CSS 側で `.is-scrolling .article-item:hover { ...base値... }` で上書き。MCP（Chrome DevTools）でホバー → スクロールして bg が `0.2 → 0.12` に戻ることを確認。→ コミット `0f5e985`
2. それでも残るケースがあったため、`.article-item:hover` の背景・影変化を完全削除し、ホバーフィードバックは `__title a:hover { color }` のみに縮約。→ コミット `13f3188`
3. `.glass-card`（記事単体ページ用）には `html:not(.is-scrolling)` ガードのみ残置。

**学び:** トラックパッドスクロールでホバー状態がリセットされない問題は CSS だけでは閉じない。`is-scrolling` クラス + JS 監視で上書きするか、ホバーで色以外を変えない設計にする。

## 4. ハンバーガーメニューが「閉じない」

- 当初は `.is-open` トグルだけで動かしていたが、ユーザの実機では「開いたまま」になっているとの報告。MCP では再現せず一度詰まった。
- 「外側タップで閉じない」「body スクロール残り」などの UX 不足が体感的な「効かない」になっていたので、UX 全部入れ替え:
  - document click 委譲で外側タップ → close
  - `body.nav-open::before` で半透明バックドロップ + `backdrop-filter: blur(4px)`
  - `body.nav-open { overflow: hidden }` で背面スクロールロック
  - z-index 整理（メニュー = `var(--z-overlay) + 1`、バックドロップ = `var(--z-overlay)`）
  - トグル自身は `stopPropagation` で外側ハンドラーへの二重反応を抑止

**学び:** モバイルメニューは「開く・閉じる」だけでなく「外側タップ」「背面スクロール」「Escape」「内部リンク遷移時の自動閉」まで含めて初めて「効いている」になる。

## 5. デザイン確認の精度

- 実装後にレイアウト崩れを目視確認していたが、ユーザから複数回「MCP でちゃんと確認すべき」「スクショ撮って確認すべき」と差し戻し。途中から chrome-devtools-mcp で `take_screenshot` → 比較を必ず挟むフローに切替。

**学び:** UI 変更は MCP でビフォーアフターのスクショを撮ってから「完了」と言う。
