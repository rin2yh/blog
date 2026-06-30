# OGP (Twitter Card) images

記事ごとの OGP 画像を Hugo のビルド時に自動生成するための素材です。

## 仕組み

1. `base.png` … 固定の背景（テーマ配色のグラデーション＋グラスパネル＋アバター＋
   サイト名）。`cmd/ogp` で生成し、コミットしてあります。
2. `layouts/partials/ogp-image.html` … `base.png` に記事タイトルを
   `images.Text` で重ねて `/post/.../base_hu_*.png` を生成し、その URL を返します。
3. `layouts/partials/head/opengraph.html` … 記事に `image` / `cover.image` が
   指定されていないとき、上記の生成画像を `og:image` / `twitter:image` に使います
   （指定があればそちらを優先、無ければ最後にサイト既定画像へフォールバック）。

CI（Hugo extended）はこの `base.png` とフォントだけで画像を生成するため、
ビルドに Go ツールチェインは不要です。

## 背景を作り直す

配色・レイアウト・サイト名などは `cmd/ogp/main.go` を編集し、リポジトリ
ルートで次を実行すると `base.png` が更新されます。

```sh
go run ./cmd/ogp
```

## フォント

`NotoSansJP-Bold.ttf` は Google Fonts の Noto Sans JP（可変フォント）を
`wght=700` の静的インスタンスに変換したものです。ライセンスは
[SIL Open Font License 1.1](./NotoSansJP-LICENSE.txt)。
