+++
date = '2026-02-01T18:40:22+09:00'
draft = false
title = 'Hugoでdetailsタグを使う方法'
categories = ["tech"]
tags = ["Hugo", "blog"]
+++

## 概要

一言で終わる概要ですが、記載します。

結論、「Hugo公式ドキュメントに記載の方法を使うだけでOK」です。

公式ドキュメント：https://gohugo.io/shortcodes/details/

## 使い方

以下のように記載することで、detailsタグを利用できます。

```md
{{</* details summary="See the details" */>}}
This is a **bold** word.
{{</* /details */>}}
```

例
{{< details summary="See the details" >}}
This is a **bold** word.
{{< /details >}}

他のブログ等でShortcodesを自作して入れるものが多いですが、Hugoのv0.14.0以降は上記の記載だけで導入できます。
カスタマイズはHugoのlayoutsを利用したり、CSSのスタイリングを変更したりといった方法で行うことができます。

他にも、ShortcodesはHugo公式ドキュメントに記載されていて、以前よりも充実していると思います。

https://gohugo.io/shortcodes/

ぜひ参考にしてみてください。

