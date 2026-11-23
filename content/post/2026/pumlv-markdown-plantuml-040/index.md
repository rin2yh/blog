+++
draft = true
title = 'MarkdownのPlantUMLを選んでプレビューする — pumlv 0.4.0'
categories = ['tech']
tags = ['pumlv', 'PlantUML', 'Markdown', 'ADR']
+++

## はじめに

pumlvの背景や構成、基本的な使い方は、[Zennの記事](https://zenn.dev/rinrin_yuuki/articles/9b69cca81875f6)で紹介しました。この記事では、v0.4.0で追加したMarkdown内のPlantUMLプレビューと、実装時に記録した設計判断を紹介します。

## Markdown内の図をプレビューする

Markdownのコードブロックに `plantuml`、`puml`、`iuml`、`wsd` のいずれかを指定すると、ファイルツリーに図が表示されます。

````markdown
```plantuml
title 認証フロー
利用者 -> Webアプリ: ログイン
Webアプリ -> 認証サービス: 認証要求
```
````

`@startuml` と `@enduml` は省略できます。pumlvが自動で補うためです。図が複数あるMarkdownファイルでは、各図を開始行番号とタイトル付きでファイルツリーに並べます。たとえば、上の図は `L1: 認証フロー` と表示されます。図を選ぶと、そのコードブロックだけをプレビューできます。

コードブロックの追加・編集・削除は、ファイル監視を通じてファイルツリーとプレビューに反映されます。PlantUMLのコードブロックを含まないMarkdownファイルは、ファイルツリーに表示しません。

Markdownの解析にはGoldmarkを使っています。リストや引用の中にあるコードブロックも扱えます。

## 設計判断

### 図の識別に開始行番号を使う

複数のコードブロックを個別に選ぶには、図ごとの選択キーが必要です。永続IDやコード内容のハッシュも候補でしたが、pumlvではMarkdownファイルのパスと開始フェンスの行番号を組み合わせました。永続IDを管理したり、編集前後のブロックを追跡したりせずに済み、識別方法を単純にできます。

行番号を使うため、前にある文章や図を編集して行番号が変わると、選択中の図を同じ図として追えない場合があります。選択した図がなくなったときは、同じファイルに残った先頭の図を選びます。編集をまたいだ同一性を保証せず、実装の単純さを優先しました。この判断の詳細は[ADR 0002](https://github.com/rin2yh/pumlv/blob/v0.4.0/docs/adr/0002-identify-markdown-diagrams-by-line.md)に記録しています。

### 解析結果をRegistryに保持する

ファイルツリーを表示するAPIと、選択した図のソースを返すAPIがあります。リクエストごとにMarkdownを解析すると、図一覧とソースで異なるファイル内容を読む可能性があります。また、ファイル単位で差分更新するRegistryの利点も失われます。

そこで、初回登録時とファイル変更時にMarkdownを解析し、図の行番号・タイトル・正規化したソースをRegistryの同じファイルエントリに保存しました。APIはそのスナップショットを参照し、Watcherは変更されたファイルを更新します。一覧とソースを同じエントリから返すことで、更新途中の情報が混ざらないようにしています。

解析結果をメモリに保持するため、その分のメモリを使います。また、Watcherがファイル変更を反映するまでは直前のスナップショットを返します。リクエストごとに解析するコストを避けるため、変更の反映までに短い遅れが出る点も受け入れました。

この判断の詳細は[ADR 0003](https://github.com/rin2yh/pumlv/blob/v0.4.0/docs/adr/0003-store-markdown-diagrams-in-registry.md)に記録しています。

## おわりに

v0.4.0は[pumlvのリリースページ](https://github.com/rin2yh/pumlv/releases/tag/v0.4.0)から入手できます。インストール方法や起動方法は[pumlvのREADME](https://github.com/rin2yh/pumlv/blob/v0.4.0/README.md)を参照してください。

今回は機能の追加とあわせて、図の識別方法と解析結果の保持方法もADRに残しました。選択状態やキャッシュの扱いを見直すときに、採用した判断とその影響をたどれます。

実装した[Markdownプレビュー対応](https://github.com/rin2yh/pumlv/pull/213)と[複数図への対応](https://github.com/rin2yh/pumlv/pull/231)もご覧ください。最後まで読んでいただき、ありがとうございました！
