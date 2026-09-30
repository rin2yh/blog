## 技術スタック
- [Hugo](https://gohugo.io/)
- [Hugo theme Stack](https://stack.jimmycai.com/)

## 記事の公開日

`hugo new`（`mise run post` / `external` / `conf`）で作成した記事には `date` を設定しません。
`draft = false` にして main にマージすると、公開ワークフローが `date` のない記事に日本時間の現在日時を追加します。
ビルドが成功した後、その日付を main にコミットしてからデプロイするため、次回以降のビルドでも公開日は変わりません。
既存の `date` と `draft = true` の記事は変更しません。過去の公開日や予約日時を指定したい場合は、手動で `date` を設定してください。

PRのビルド・スクリーンショットは `preview` 環境を使い、日付がない場合だけファイル更新日時を表示します。
手元でも `hugo server --environment preview --buildDrafts --buildFuture` で確認できます。
通常のビルドは記事に記録された `date` のみを使用します。

公開ワークフローは日付を保存するために main への push を行います。
ブランチ保護を設定する場合は、このワークフローからの push を許可してください。
