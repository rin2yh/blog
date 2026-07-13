+++
date = '2026-03-01T18:28:51+09:00'
draft = true
title = 'Recurrent Effective Go'
categories = []
tags = []
+++

構成(5m)
- 自己紹介 30s
- Effective Goとは 30s
- なぜ、もう一度Effective Goを学ぶのか 30s
- 改めて得た学びの例 1m
- Recurrent Effective Goの用途 2m
- まとめ 30s

## 自己紹介
ニックネーム
趣味：アニメ
Go歴：4年
業務システム、Web API開発でGoを使っている

## Effective Goとは
TODO:書く

## なぜ、もう一度Effective Goを学ぶのか

Go自体を業務で書いてきた年数自体は3年ほどあり、業務Webシステム的な記法ではおおよそ迷いません。
Effective Go自体も一度学んでいます。
しかし、Effective Goの内容でも意外と忘れていることがあります。
理由は至って単純で、業務のプログラムで使われていない構文やベストプラクティスがあるからです。
一人前のGopherになるためにGoの思想をより理解したいと思ったこともあります。

これをもう一度＝リカレント[^1]にEffective Goを見直してみて、学びを得ようという話です。

[^1]リカレントと名付けた理由：https://www.mhlw.go.jp/stf/newpage_18817.html
## 改めて得た学びの例

- package名のベストプラクティス：「short, concise, evocative」

> It's helpful if everyone using the package can use the same name to refer to its contents,
which implies that the package name should be good: short, concise, evocative.

翻訳：パッケージを利用する全員がその中身を同じ名前で参照できると便利です。
そのためには、パッケージ名が短く、簡潔で、内容を連想しやすい、優れたものである必要があります。
By nani\.now

業務で見たパッケージ名の例
- `testUtils`
→internal/pkg/test, pkg/test

https://go.dev/doc/effective_go#package-names

Interface Names
- 1method interface→ method+er：例 Reader
- Read, Write, Close, Flush, Stringは標準的なシグネチャと同じ意味・実装の場合だけ使用する
https://go.dev/doc/effective_go#interface-names

init関数
あまり使えてないので、忘れがち。これまで携わった業務コードではほぼ見なかった。

https://go.dev/doc/effective_go#init

他にあれば書く

## Recurrent Effective Goの用法

1. 部分的に使う（おすすめ）
- Go理解できてないかもなって時
- レビューでEffective Goを参考文献に挙げて理解度が怪しい時、後で見る

2. 全部見る（非推奨）
- タイパが微妙かも
- Effective Goの性質的に、〜〜


## まとめ
改めて、リカレントにEffective Goを学ぶ価値があると思っています。
使っていない文法があればどこで使うか、なぜ存在するかとか考えるきっかけを得ることができます。
流石に全部覚え切ったり、理解し切ったりすることは難しいです。

また、レビュー時など普段使いもできる代物です。
時折、Effective Goいたね、くらいでも良いので思い出してあげてください。

## 宣伝：エンジニアニメというイベントのスタッフしてます。
4/11に劇場版エンジニアニメがあるので、ぜひ来てください。
エンジニアリングとアニメをリンクさせて、よりエンジニアリングの造詣を深めたり
アニメが好きなエンジニアと交流したりできます。

## 参考文献

- Effective Go: https://go.dev/doc/effective_go
