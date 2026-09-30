## 技術スタック
- [Hugo](https://gohugo.io/)
- [Hugo theme Stack](https://stack.jimmycai.com/)

## 記事の公開日

`draft = false` にして main にマージすると、未指定の `date` が公開時の日本時間で記録され、以後は固定されます。日時を指定したい場合は、手動で `date` を設定してください。

手元のプレビューは `hugo server --environment preview --buildDrafts --buildFuture` を使います。日付のない記事は、一時的にファイル更新日時を表示します。
